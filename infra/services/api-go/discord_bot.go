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
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

// Discord channel IDs in the TriGameDev Guild (336559457242972160)
const channelSubsID = "1571838890778566716" // #channel-subs

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

// Opt-in notification roles in the TriGameDev Guild
const (
	roleNewGamesID = "1571839186192048269" // Tell Me About New Games (release-feed & key vault)
	roleNewBooksID = "1571839345973788782" // Tell Me About New Books (free-books)
	roleNewJamsID  = "1571839429436379236" // Tell Me About New Jams (monthly-game-jam-announcements)
)

type notifRoleSpec struct {
	Key      string
	RoleID   string
	RoleName string
}

var notificationRoleSpecs = []notifRoleSpec{
	{Key: "games", RoleID: roleNewGamesID, RoleName: "Tell Me About New Games"},
	{Key: "books", RoleID: roleNewBooksID, RoleName: "Tell Me About New Books"},
	{Key: "jams", RoleID: roleNewJamsID, RoleName: "Tell Me About New Jams"},
}

var botSetupReady = false

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

func (s *Server) discordAPIRequest(method, path string, body any) ([]byte, int, error) {
	botToken := strings.TrimSpace(os.Getenv("DISCORD_BOT_TOKEN"))
	if botToken == "" {
		return nil, 0, fmt.Errorf("DISCORD_BOT_TOKEN not configured")
	}
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, "https://discord.com/api/v10"+path, reader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bot "+botToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 65536))
	return respBytes, resp.StatusCode, nil
}

// Resolves the mention prefix and allowed_mentions object for a given section:
// - "news" (and "events" in admin/news) uses @everyone
// - "showcase" and "keys" use <@&1571839186192048269> (Tell Me About New Games)
// - "ebooks" uses <@&1571839345973788782> (Tell Me About New Books)
// - "jams" uses <@&1571839429436379236> (Tell Me About New Jams)
// - "gallery" and "videos" use NO @ directive
func (s *Server) buildMentionForSection(section string) (string, map[string]any) {
	switch section {
	case "news", "events":
		return "@everyone", map[string]any{
			"parse": []string{"everyone"},
		}
	case "showcase", "keys":
		return fmt.Sprintf("<@&%s>", roleNewGamesID), map[string]any{
			"roles": []string{roleNewGamesID},
		}
	case "ebooks":
		return fmt.Sprintf("<@&%s>", roleNewBooksID), map[string]any{
			"roles": []string{roleNewBooksID},
		}
	case "jams":
		return fmt.Sprintf("<@&%s>", roleNewJamsID), map[string]any{
			"roles": []string{roleNewJamsID},
		}
	case "gallery", "videos":
		return "", map[string]any{
			"parse": []string{},
		}
	default:
		return "", map[string]any{
			"parse": []string{},
		}
	}
}

func (s *Server) postToDiscordChannel(section string, headline string, embed *discordEmbed) error {
	channelID := defaultSectionChannels[section]
	if channelID == "" {
		return fmt.Errorf("unknown Discord channel for section %s", section)
	}

	mentionPrefix, allowedMentions := s.buildMentionForSection(section)
	content := strings.TrimSpace(mentionPrefix + " " + headline)

	payload := map[string]any{
		"content":          content,
		"allowed_mentions": allowedMentions,
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

	respBytes, status, err := s.discordAPIRequest("POST", fmt.Sprintf("/channels/%s/messages", channelID), payload)
	if err != nil {
		return err
	}
	if status >= 300 {
		return fmt.Errorf("discord API status %d: %s", status, string(respBytes))
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
		if err := s.postToDiscordChannel("jams", fmt.Sprintf("**%s** is on TriGameDev.org!\n%s", title, siteURL), emb); err != nil {
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
		if err := s.postToDiscordChannel("events", fmt.Sprintf("Upcoming Event: **%s**\n%s", title, siteURL), emb); err != nil {
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
		if err := s.postToDiscordChannel("news", fmt.Sprintf("New Announcement: **%s**\n%s", title, siteURL), emb); err != nil {
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
		if err := s.postToDiscordChannel("showcase", fmt.Sprintf("Check out **%s** by **%s** on Games by Our People!\n%s", title, dev, siteURL), emb); err != nil {
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
		fullImg := resolvePublicURL(imgURL)
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
		if fullImg != "" {
			emb.Image = &discordEmbedMedia{URL: fullImg}
		}
		msgText := fmt.Sprintf("**%s** (shared by %s)\n%s", title, uploadedBy, siteURL)
		if fullImg != "" {
			msgText = fmt.Sprintf("**%s** (shared by %s)\n%s\n%s", title, uploadedBy, siteURL, fullImg)
		}
		if err := s.postToDiscordChannel("gallery", msgText, emb); err != nil {
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
				{Name: "Video Link", Value: videoURL, Inline: false},
			},
		}
		if img := resolvePublicURL(thumbURL); img != "" {
			emb.Image = &discordEmbedMedia{URL: img}
		}
		if err := s.postToDiscordChannel("videos", fmt.Sprintf("**%s** by **%s**\n%s\n%s", title, speaker, siteURL, videoURL), emb); err != nil {
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
		if err := s.postToDiscordChannel("keys", fmt.Sprintf("**%s** shared keys for **%s** (%s) in the Key Vault!\n%s", donor, title, platform, siteURL), emb); err != nil {
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
		if err := s.postToDiscordChannel("ebooks", fmt.Sprintf("New book in the Ebook Library: **%s** by *%s*!\n%s", title, author, siteURL), emb); err != nil {
			return err
		}
		_, _ = s.db.Exec(`UPDATE trigamedev.ebooks SET discord_announced = TRUE WHERE id = $1`, id)
		return nil
	}

	return fmt.Errorf("unknown section: %s", section)
}

// Pick 1 unannounced entry across the database (excluding news, which only triggers on new posts) and post it to Discord
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

// Ensures the interactive subscription panel post exists in #channel-subs (1571838890778566716)
func (s *Server) ensureDiscordNotificationsSetup() (map[string]any, error) {
	// Get bot's own user ID
	var botUserID = "1482097872521400582"
	if meBytes, meStatus, meErr := s.discordAPIRequest("GET", "/users/@me", nil); meErr == nil && meStatus == 200 {
		var me struct {
			ID string `json:"id"`
		}
		if json.Unmarshal(meBytes, &me) == nil && me.ID != "" {
			botUserID = me.ID
		}
	}

	// Check if our interactive button post already exists in #channel-subs
	msgsBytes, mStatus, mErr := s.discordAPIRequest("GET", fmt.Sprintf("/channels/%s/messages?limit=20", channelSubsID), nil)
	existingMsgID := ""
	if mErr == nil && mStatus == 200 {
		var msgs []struct {
			ID     string `json:"id"`
			Author struct {
				ID string `json:"id"`
			} `json:"author"`
			Components []any `json:"components"`
		}
		if json.Unmarshal(msgsBytes, &msgs) == nil {
			for _, m := range msgs {
				if m.Author.ID == botUserID && len(m.Components) > 0 {
					existingMsgID = m.ID
					break
				}
			}
		}
	}

	panelPayload := map[string]any{
		"embeds": []any{
			map[string]any{
				"title":       "CHANNEL SUBSCRIPTIONS",
				"description": "Click the buttons below to subscribe or unsubscribe from notification pings when new content is posted on **TriGameDev.org**.",
				"color":       0xEF4444,
				"fields": []map[string]any{
					{"name": "Tell Me About New Games", "value": "Notifications for `#release-feed` and `#free-games` (Key Vault)", "inline": false},
					{"name": "Tell Me About New Books", "value": "Notifications for `#free-books` (Ebook Library)", "inline": false},
					{"name": "Tell Me About New Jams", "value": "Notifications for `#monthly-game-jam-announcements`", "inline": false},
				},
				"footer": map[string]any{
					"text": "TriGameDev.org",
				},
			},
		},
		"components": []any{
			map[string]any{
				"type": 1, // ActionRow
				"components": []any{
					map[string]any{"type": 2, "style": 2, "label": "Tell Me About New Games", "custom_id": "tgd_notif:games"},
					map[string]any{"type": 2, "style": 2, "label": "Tell Me About New Books", "custom_id": "tgd_notif:books"},
					map[string]any{"type": 2, "style": 2, "label": "Tell Me About New Jams", "custom_id": "tgd_notif:jams"},
				},
			},
		},
	}

	if existingMsgID == "" {
		postBytes, pStatus, pErr := s.discordAPIRequest("POST", fmt.Sprintf("/channels/%s/messages", channelSubsID), panelPayload)
		if pErr != nil || (pStatus != 200 && pStatus != 201) {
			return nil, fmt.Errorf("failed to post subscription panel in #channel-subs (HTTP %d: %s)", pStatus, string(postBytes))
		}
		var posted struct {
			ID string `json:"id"`
		}
		_ = json.Unmarshal(postBytes, &posted)
		existingMsgID = posted.ID
		log.Printf("[Discord Setup] Posted interactive subscription panel (%s) in #channel-subs", existingMsgID)
	} else {
		_, _, _ = s.discordAPIRequest("PATCH", fmt.Sprintf("/channels/%s/messages/%s", channelSubsID, existingMsgID), panelPayload)
	}

	botSetupReady = true
	return map[string]any{
		"channel_subs_id":  channelSubsID,
		"panel_message_id": existingMsgID,
		"roles": map[string]string{
			"games": roleNewGamesID,
			"books": roleNewBooksID,
			"jams":  roleNewJamsID,
		},
	}, nil
}

// Handles Discord Button Interaction (toggles the corresponding notification role on the clicking member)
func (s *Server) handleDiscordComponentInteraction(rawEvent json.RawMessage) {
	var inter struct {
		ID      string `json:"id"`
		Token   string `json:"token"`
		Type    int    `json:"type"`
		GuildID string `json:"guild_id"`
		Data    struct {
			CustomID string `json:"custom_id"`
		} `json:"data"`
		Member struct {
			Roles []string `json:"roles"`
			User  struct {
				ID       string `json:"id"`
				Username string `json:"username"`
			} `json:"user"`
		} `json:"member"`
	}
	if err := json.Unmarshal(rawEvent, &inter); err != nil {
		return
	}
	if inter.Type != 3 || !strings.HasPrefix(inter.Data.CustomID, "tgd_notif:") {
		return
	}

	roleKey := strings.TrimPrefix(inter.Data.CustomID, "tgd_notif:")
	roleID := ""
	roleName := ""
	for _, spec := range notificationRoleSpecs {
		if spec.Key == roleKey {
			roleID = spec.RoleID
			roleName = spec.RoleName
			break
		}
	}

	respondEphemeral := func(msg string) {
		_, _, _ = s.discordAPIRequest("POST", fmt.Sprintf("/interactions/%s/%s/callback", inter.ID, inter.Token), map[string]any{
			"type": 4, // CHANNEL_MESSAGE_WITH_SOURCE
			"data": map[string]any{
				"content": msg,
				"flags":   64, // EPHEMERAL
			},
		})
	}

	if roleID == "" || inter.GuildID == "" || inter.Member.User.ID == "" {
		respondEphemeral("Unknown notification role.")
		return
	}

	hasRole := false
	for _, r := range inter.Member.Roles {
		if r == roleID {
			hasRole = true
			break
		}
	}

	memberRolePath := fmt.Sprintf("/guilds/%s/members/%s/roles/%s", inter.GuildID, inter.Member.User.ID, roleID)
	if hasRole {
		respBytes, status, err := s.discordAPIRequest("DELETE", memberRolePath, nil)
		if err != nil || status >= 300 {
			log.Printf("[Discord Role Toggle] Failed to remove role %s from %s (HTTP %d: %s)", roleID, inter.Member.User.ID, status, string(respBytes))
			respondEphemeral("Could not update your role — please ensure the `TGDAuth` role has the **Manage Roles** permission enabled in Discord Server Settings.")
			return
		}
		respondEphemeral(fmt.Sprintf("Removed **%s** — you are unsubscribed.", roleName))
	} else {
		respBytes, status, err := s.discordAPIRequest("PUT", memberRolePath, nil)
		if err != nil || status >= 300 {
			log.Printf("[Discord Role Toggle] Failed to add role %s to %s (HTTP %d: %s)", roleID, inter.Member.User.ID, status, string(respBytes))
			respondEphemeral("Could not update your role — please ensure the `TGDAuth` role has the **Manage Roles** permission enabled in Discord Server Settings.")
			return
		}
		respondEphemeral(fmt.Sprintf("Added **%s** — you are now subscribed!", roleName))
	}
}

// Persistent Discord Gateway v10 WebSocket loop to receive button clicks (INTERACTION_CREATE)
func (s *Server) startDiscordGateway() {
	go func() {
		for {
			botToken := strings.TrimSpace(os.Getenv("DISCORD_BOT_TOKEN"))
			if botToken == "" {
				time.Sleep(30 * time.Second)
				continue
			}

			s.runGatewaySession(botToken)
			time.Sleep(5 * time.Second)
		}
	}()
}

func (s *Server) runGatewaySession(botToken string) {
	conn, _, err := websocket.DefaultDialer.Dial("wss://gateway.discord.gg/?v=10&encoding=json", nil)
	if err != nil {
		log.Printf("[Discord Gateway] Dial error: %v", err)
		return
	}
	defer conn.Close()

	var writeMu sync.Mutex
	sendJSON := func(v any) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		return conn.WriteJSON(v)
	}

	type gatewayPayload struct {
		Op int             `json:"op"`
		D  json.RawMessage `json:"d"`
		S  *int64          `json:"s"`
		T  string          `json:"t"`
	}

	// Read HELLO (op 10)
	var hello gatewayPayload
	if err := conn.ReadJSON(&hello); err != nil || hello.Op != 10 {
		return
	}
	var helloData struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	}
	_ = json.Unmarshal(hello.D, &helloData)
	if helloData.HeartbeatInterval <= 0 {
		helloData.HeartbeatInterval = 41250
	}

	// Send IDENTIFY (op 2)
	_ = sendJSON(map[string]any{
		"op": 2,
		"d": map[string]any{
			"token":   botToken,
			"intents": 1, // GUILDS
			"properties": map[string]string{
				"os":      "linux",
				"browser": "trigamedev-bot",
				"device":  "trigamedev-bot",
			},
		},
	})

	stopHB := make(chan struct{})
	defer close(stopHB)

	var lastSeq *int64
	go func() {
		ticker := time.NewTicker(time.Duration(helloData.HeartbeatInterval) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopHB:
				return
			case <-ticker.C:
				if err := sendJSON(map[string]any{"op": 1, "d": lastSeq}); err != nil {
					return
				}
			}
		}
	}()

	for {
		var msg gatewayPayload
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}
		if msg.S != nil {
			seq := *msg.S
			lastSeq = &seq
		}
		switch msg.Op {
		case 0: // Dispatch
			if msg.T == "INTERACTION_CREATE" {
				go s.handleDiscordComponentInteraction(msg.D)
			}
		case 1: // Heartbeat request
			_ = sendJSON(map[string]any{"op": 1, "d": lastSeq})
		case 7, 9: // Reconnect / Invalid Session
			return
		}
	}
}

// Background loop that posts 1 unannounced entry every 6 hours and ensures the #channel-subs panel is posted
func (s *Server) startCatchupWorker() {
	s.startDiscordGateway()

	go func() {
		time.Sleep(3 * time.Second)
		if _, err := s.ensureDiscordNotificationsSetup(); err != nil {
			log.Printf("[Discord Setup Notice] %v", err)
		}
		for {
			if !botSetupReady {
				if _, err := s.ensureDiscordNotificationsSetup(); err == nil {
					log.Printf("[Discord Setup] #channel-subs panel is now active!")
				}
			}
			if sec, id, err := s.processSingleCatchupPost(false); err != nil {
				log.Printf("[Discord Catchup Error] %s #%d: %v", sec, id, err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}
