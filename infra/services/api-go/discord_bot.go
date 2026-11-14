package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var reNonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\u2019", "")
	s = reNonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "post"
	}
	return s
}

// Default Discord channel IDs in the TriGameDev Guild (336559457242972160)
var defaultSectionChannels = map[string]string{
	"jams":     "1480332771166654485", // jams / monthly-game-jam-announcements
	"events":   "1571142618363076619", // admin / news
	"news":     "1571142618363076619", // admin / news
	"showcase": "1571143126964379698", // games / release-feed
	"gallery":  "1571143218597335083", // random / gallery
	"videos":   "1571143218597335083", // random / gallery
	"keys":     "1560744377586819162", // games / free-games
	"ebooks":   "1571143288935944315", // random / free-books
}

type discordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	URL         string              `json:"url,omitempty"`
	Color       int                 `json:"color,omitempty"`
	Fields      []discordEmbedField `json:"fields,omitempty"`
	Image       *discordEmbedMedia  `json:"image,omitempty"`
	Footer      *discordEmbedFooter `json:"footer,omitempty"`
}

type discordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type discordEmbedMedia struct {
	URL string `json:"url"`
}

type discordEmbedFooter struct {
	Text string `json:"text"`
}

func resolvePublicURL(raw string) string {
	u := strings.TrimSpace(raw)
	if u == "" || u == "/brand/filler.png" {
		return ""
	}
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	if strings.HasPrefix(u, "/") {
		return "https://trigamedev.org" + u
	}
	return "https://trigamedev.org/" + u
}

func (s *Server) postToDiscordChannel(section string, content string, embed *discordEmbed) error {
	botToken := strings.TrimSpace(os.Getenv("DISCORD_BOT_TOKEN"))
	if botToken == "" {
		return fmt.Errorf("DISCORD_BOT_TOKEN not configured")
	}
	channelID := defaultSectionChannels[section]
	if channelID == "" {
		return fmt.Errorf("unknown Discord channel for section %s", section)
	}

	if !strings.Contains(content, "@here") {
		content = "@here " + content
	}

	payload := map[string]any{
		"content": content,
		"allowed_mentions": map[string]any{
			"parse": []string{"everyone"}, // Discord API uses "everyone" parse group to enable @here and @everyone
		},
	}
	if embed != nil {
		if embed.Color == 0 {
			embed.Color = 0xEF4444 // TriGameDev Red
		}
		if embed.Footer == nil {
			embed.Footer = &discordEmbedFooter{Text: "TriGameDev.org"}
		}
		payload["embeds"] = []any{embed}
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID), bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("discord API status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// Announce a specific entry by section and ID, and mark discord_announced = TRUE
func (s *Server) announceEntryByID(section string, id int) error {
	switch section {
	case "jams":
		var title, monthLabel, theme, itchURL, coverURL string
		var parts, subs int
		err := s.db.QueryRow(`
			SELECT title, month_label, theme, itch_url, cover_image_url, participants_count, submissions_count
			FROM trigamedev.game_jams WHERE id = $1
		`, id).Scan(&title, &monthLabel, &theme, &itchURL, &coverURL, &parts, &subs)
		if err != nil {
			return err
		}
		siteURL := fmt.Sprintf("https://trigamedev.org/jams/%s", slugify(title))
		fields := []discordEmbedField{
			{Name: "Jammers", Value: fmt.Sprintf("%d", parts), Inline: true},
			{Name: "Entries", Value: fmt.Sprintf("%d", subs), Inline: true},
			{Name: "View on TriGameDev", Value: siteURL, Inline: false},
		}
		if itchURL != "" {
			fields = append(fields, discordEmbedField{Name: "Itch.io Jam Page", Value: itchURL, Inline: false})
		}
		emb := &discordEmbed{
			Title:       title,
			URL:         siteURL,
			Description: fmt.Sprintf("**Theme:** %s\n**Month:** %s", theme, monthLabel),
			Fields:      fields,
		}
		if img := resolvePublicURL(coverURL); img != "" {
			emb.Image = &discordEmbedMedia{URL: img}
		}
		if err := s.postToDiscordChannel("jams", fmt.Sprintf("@here **%s** is on TriGameDev.org!\n%s", title, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.game_jams SET discord_announced = TRUE WHERE id = $1`, id)
		return nil

	case "events":
		var title, desc, loc, coverURL, tag string
		var st time.Time
		err := s.db.QueryRow(`
			SELECT title, description, location, start_time, cover_image_url, tag
			FROM trigamedev.events WHERE id = $1
		`, id).Scan(&title, &desc, &loc, &st, &coverURL, &tag)
		if err != nil {
			return err
		}
		siteURL := fmt.Sprintf("https://trigamedev.org/events/%s", slugify(title))
		emb := &discordEmbed{
			Title:       title,
			URL:         siteURL,
			Description: desc,
			Fields: []discordEmbedField{
				{Name: "When", Value: st.Local().Format("Mon, Jan 02, 2006 3:04 PM"), Inline: true},
				{Name: "Location", Value: loc, Inline: true},
				{Name: "Event Link", Value: siteURL, Inline: false},
			},
		}
		if img := resolvePublicURL(coverURL); img != "" {
			emb.Image = &discordEmbedMedia{URL: img}
		}
		if err := s.postToDiscordChannel("events", fmt.Sprintf("@here Upcoming Event: **%s**\n%s", title, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.events SET discord_announced = TRUE WHERE id = $1`, id)
		return nil

	case "news":
		var title, summary, body, author, tag string
		err := s.db.QueryRow(`
			SELECT title, summary, body_markdown, author_name, tag
			FROM trigamedev.news WHERE id = $1
		`, id).Scan(&title, &summary, &body, &author, &tag)
		if err != nil {
			return err
		}
		siteURL := fmt.Sprintf("https://trigamedev.org/news/%s", slugify(title))
		desc := body
		if strings.TrimSpace(summary) != "" && summary != body {
			desc = fmt.Sprintf("**%s**\n\n%s", summary, body)
		}
		emb := &discordEmbed{
			Title:       title,
			URL:         siteURL,
			Description: desc,
			Fields: []discordEmbedField{
				{Name: "Posted By", Value: author, Inline: true},
				{Name: "Tag", Value: tag, Inline: true},
				{Name: "Read on TriGameDev", Value: siteURL, Inline: false},
			},
		}
		if err := s.postToDiscordChannel("news", fmt.Sprintf("@here New Announcement: **%s**\n%s", title, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.news SET discord_announced = TRUE WHERE id = $1`, id)
		return nil

	case "showcase":
		var title, dev, desc, thumb, playURL, engine, actionLabel string
		err := s.db.QueryRow(`
			SELECT title, developer_name, description, thumbnail_url, play_url, engine, COALESCE(action_label, '')
			FROM trigamedev.showcase WHERE id = $1
		`, id).Scan(&title, &dev, &desc, &thumb, &playURL, &engine, &actionLabel)
		if err != nil {
			return err
		}
		siteURL := fmt.Sprintf("https://trigamedev.org/games/%s", slugify(title))
		if actionLabel == "" {
			if strings.Contains(strings.ToLower(playURL), "steampowered.com") {
				actionLabel = "Buy on Steam"
			} else {
				actionLabel = "Play Now on Itch.io"
			}
		}
		emb := &discordEmbed{
			Title:       title,
			URL:         siteURL,
			Description: desc,
			Fields: []discordEmbedField{
				{Name: "Developer", Value: dev, Inline: true},
				{Name: "Tag / Engine", Value: engine, Inline: true},
				{Name: "View on TriGameDev", Value: siteURL, Inline: false},
				{Name: actionLabel, Value: playURL, Inline: false},
			},
		}
		if img := resolvePublicURL(thumb); img != "" {
			emb.Image = &discordEmbedMedia{URL: img}
		}
		if err := s.postToDiscordChannel("showcase", fmt.Sprintf("@here Check out **%s** by **%s** on Games by Our People!\n%s", title, dev, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.showcase SET discord_announced = TRUE WHERE id = $1`, id)
		return nil

	case "gallery":
		var title, caption, imgURL, evLabel, uploadedBy string
		err := s.db.QueryRow(`
			SELECT title, caption, image_url, event_label, uploaded_by
			FROM trigamedev.gallery WHERE id = $1
		`, id).Scan(&title, &caption, &imgURL, &evLabel, &uploadedBy)
		if err != nil {
			return err
		}
		siteURL := fmt.Sprintf("https://trigamedev.org/media/%s", slugify(title))
		emb := &discordEmbed{
			Title:       title,
			URL:         siteURL,
			Description: caption,
			Fields: []discordEmbedField{
				{Name: "Album", Value: evLabel, Inline: true},
				{Name: "Shared By", Value: uploadedBy, Inline: true},
				{Name: "View on TriGameDev", Value: siteURL, Inline: false},
			},
		}
		if img := resolvePublicURL(imgURL); img != "" {
			emb.Image = &discordEmbedMedia{URL: img}
		}
		if err := s.postToDiscordChannel("gallery", fmt.Sprintf("@here New photo added to the TriGameDev Gallery: **%s**\n%s", title, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.gallery SET discord_announced = TRUE WHERE id = $1`, id)
		return nil

	case "videos":
		var title, desc, videoURL, thumbURL, speaker, duration string
		err := s.db.QueryRow(`
			SELECT title, description, video_url, thumbnail_url, speaker_name, duration
			FROM trigamedev.videos WHERE id = $1
		`, id).Scan(&title, &desc, &videoURL, &thumbURL, &speaker, &duration)
		if err != nil {
			return err
		}
		siteURL := fmt.Sprintf("https://trigamedev.org/media/%s", slugify(title))
		emb := &discordEmbed{
			Title:       title,
			URL:         siteURL,
			Description: desc,
			Fields: []discordEmbedField{
				{Name: "Creator", Value: speaker, Inline: true},
				{Name: "View on TriGameDev", Value: siteURL, Inline: false},
				{Name: "Direct Video Link", Value: videoURL, Inline: false},
			},
		}
		if img := resolvePublicURL(thumbURL); img != "" {
			emb.Image = &discordEmbedMedia{URL: img}
		}
		if err := s.postToDiscordChannel("videos", fmt.Sprintf("@here New video posted to TriGameDev: **%s**\n%s", title, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.videos SET discord_announced = TRUE WHERE id = $1`, id)
		return nil

	case "keys":
		var title, platform, desc, coverURL, donor string
		var totalKeys int
		err := s.db.QueryRow(`
			SELECT p.title, p.platform, p.description, p.cover_image_url, p.donated_by_name,
			       (SELECT COUNT(*) FROM trigamedev.keys k WHERE k.pool_id = p.id AND k.claimed_by IS NULL)
			FROM trigamedev.key_pools p WHERE p.id = $1
		`, id).Scan(&title, &platform, &desc, &coverURL, &donor, &totalKeys)
		if err != nil {
			return err
		}
		siteURL := fmt.Sprintf("https://trigamedev.org/keys/%s", slugify(title))
		emb := &discordEmbed{
			Title:       title,
			URL:         siteURL,
			Description: desc,
			Fields: []discordEmbedField{
				{Name: "Platform", Value: platform, Inline: true},
				{Name: "Donated By", Value: donor, Inline: true},
				{Name: "Available Keys", Value: fmt.Sprintf("%d", totalKeys), Inline: true},
				{Name: "Claim Key on TriGameDev", Value: siteURL, Inline: false},
			},
		}
		if img := resolvePublicURL(coverURL); img != "" {
			emb.Image = &discordEmbedMedia{URL: img}
		}
		if err := s.postToDiscordChannel("keys", fmt.Sprintf("@here **%s** shared keys for **%s** (%s) in the Key Vault!\n%s", donor, title, platform, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.key_pools SET discord_announced = TRUE WHERE id = $1`, id)
		return nil

	case "ebooks":
		var title, author, category, desc, coverURL, fmtLabel, sizeLabel, uploadedBy string
		err := s.db.QueryRow(`
			SELECT title, author, category, description, cover_image_url, file_format, file_size_label, uploaded_by_name
			FROM trigamedev.ebooks WHERE id = $1
		`, id).Scan(&title, &author, &category, &desc, &coverURL, &fmtLabel, &sizeLabel, &uploadedBy)
		if err != nil {
			return err
		}
		siteURL := fmt.Sprintf("https://trigamedev.org/library/%s", slugify(title))
		emb := &discordEmbed{
			Title:       title,
			URL:         siteURL,
			Description: desc,
			Fields: []discordEmbedField{
				{Name: "Author", Value: author, Inline: true},
				{Name: "Category", Value: category, Inline: true},
				{Name: "Format", Value: fmt.Sprintf("%s (%s)", fmtLabel, sizeLabel), Inline: true},
				{Name: "Shared By", Value: uploadedBy, Inline: true},
				{Name: "Download on TriGameDev", Value: siteURL, Inline: false},
			},
		}
		if img := resolvePublicURL(coverURL); img != "" {
			emb.Image = &discordEmbedMedia{URL: img}
		}
		if err := s.postToDiscordChannel("ebooks", fmt.Sprintf("@here New book in the Ebook Library: **%s** by *%s*!\n%s", title, author, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.ebooks SET discord_announced = TRUE WHERE id = $1`, id)
		return nil
	}

	return fmt.Errorf("unknown section: %s", section)
}

// Returns true if there are still unannounced entries in the 6-hour catch-up queue
func (s *Server) hasUnannouncedEntries() bool {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT id FROM trigamedev.game_jams WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT id FROM trigamedev.events WHERE COALESCE(discord_announced, FALSE) = FALSE AND start_time >= NOW()
			UNION ALL
			SELECT id FROM trigamedev.news WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT id FROM trigamedev.showcase WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT id FROM trigamedev.gallery WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT id FROM trigamedev.videos WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT id FROM trigamedev.key_pools WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT id FROM trigamedev.ebooks WHERE COALESCE(discord_announced, FALSE) = FALSE
		) q
	`).Scan(&count)
	return err == nil && count > 0
}

// Pick 1 unannounced entry across the entire database and post it to Discord
func (s *Server) processSingleCatchupPost(force bool) (string, int, error) {
	if !force {
		var lastRun time.Time
		err := s.db.QueryRow(`SELECT updated_at FROM trigamedev.bot_state WHERE key = 'last_catchup_post_at'`).Scan(&lastRun)
		if err == nil && time.Since(lastRun) < 6*time.Hour {
			return "", 0, nil
		}
	}

	var section string
	var id int
	err := s.db.QueryRow(`
		SELECT section, id FROM (
			SELECT 'jams' AS section, id FROM trigamedev.game_jams WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT 'events' AS section, id FROM trigamedev.events WHERE COALESCE(discord_announced, FALSE) = FALSE AND start_time >= NOW()
			UNION ALL
			SELECT 'news' AS section, id FROM trigamedev.news WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT 'showcase' AS section, id FROM trigamedev.showcase WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT 'gallery' AS section, id FROM trigamedev.gallery WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT 'videos' AS section, id FROM trigamedev.videos WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT 'keys' AS section, id FROM trigamedev.key_pools WHERE COALESCE(discord_announced, FALSE) = FALSE
			UNION ALL
			SELECT 'ebooks' AS section, id FROM trigamedev.ebooks WHERE COALESCE(discord_announced, FALSE) = FALSE
		) unannounced
		ORDER BY RANDOM()
		LIMIT 1
	`).Scan(&section, &id)
	if err != nil {
		// Nothing left to catch up!
		return "", 0, nil
	}

	if err := s.announceEntryByID(section, id); err != nil {
		return section, id, err
	}

	_, _ = s.db.Exec(`
		INSERT INTO trigamedev.bot_state (key, updated_at)
		VALUES ('last_catchup_post_at', NOW())
		ON CONFLICT (key) DO UPDATE SET updated_at = NOW()
	`)
	log.Printf("[Discord Catchup] Announced %s #%d to Discord", section, id)
	return section, id, nil
}

// Background loop that checks every 5 minutes if 6 hours have elapsed since the last catch-up post
func (s *Server) startCatchupWorker() {
	go func() {
		time.Sleep(15 * time.Second)
		for {
			if sec, id, err := s.processSingleCatchupPost(false); err != nil {
				log.Printf("[Discord Catchup Error] %s #%d: %v", sec, id, err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}
