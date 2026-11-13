package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ItchScrapedData struct {
	Kind              string `json:"kind"` // "jam" or "game"
	URL               string `json:"url"`
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`           // Theme (for jam) or Developer Name(s) (for game)
	Description       string `json:"description"`        // Summary / body
	Tag               string `json:"tag"`                // Month label (for jam) or Engine / Made With / Genre (for game)
	ImageURL          string `json:"image_url"`          // Cover image
	SecondaryURL      string `json:"secondary_url"`      // External link (e.g. Steam store link) or Itch URL
	StartDate         string `json:"start_date"`         // RFC3339
	EndDate           string `json:"end_date"`           // RFC3339
	Status            string `json:"status"`             // upcoming, active, voting, completed
	ParticipantsCount int    `json:"participants_count"` // Joined / Ratings
	SubmissionsCount  int    `json:"submissions_count"`  // Entries
}

var (
	reMetaProp  = regexp.MustCompile(`(?i)<meta[^>]+(?:property|name)=["']([^"']+)["'][^>]+content=["']([^"']*)["']`)
	reMetaProp2 = regexp.MustCompile(`(?i)<meta[^>]+content=["']([^"']*)["'][^>]+(?:property|name)=["']([^"']+)["']`)
	reViewJam   = regexp.MustCompile(`ViewJam\([^,]+,\s*(\{[^}]+\})\)`)
	reStatBox   = regexp.MustCompile(`(?s)<div class="stat_box">.*?<div class="stat_value">([0-9,]+)</div><div class="stat_label">([^<]+)</div>`)
	reJamBody   = regexp.MustCompile(`(?s)<div class="jam_content[^"]*">(.*?)</div>`)
	reTableRow  = regexp.MustCompile(`(?s)<tr><td>([^<]+)</td><td>(.*?)</td></tr>`)
	reAnchor    = regexp.MustCompile(`(?s)<a[^>]*>(.*?)</a>`)
	reHref      = regexp.MustCompile(`(?i)href=["']([^"']+)["']`)
	reHTMLTags  = regexp.MustCompile(`<[^>]+>`)
)

func extractMetaMap(rawHTML string) map[string]string {
	m := make(map[string]string)
	for _, match := range reMetaProp.FindAllStringSubmatch(rawHTML, -1) {
		m[strings.ToLower(match[1])] = html.UnescapeString(match[2])
	}
	for _, match := range reMetaProp2.FindAllStringSubmatch(rawHTML, -1) {
		key := strings.ToLower(match[2])
		if _, exists := m[key]; !exists {
			m[key] = html.UnescapeString(match[1])
		}
	}
	return m
}

func cleanHTMLText(s string) string {
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = strings.ReplaceAll(s, "<br/>", "\n")
	s = strings.ReplaceAll(s, "<br />", "\n")
	s = strings.ReplaceAll(s, "</p>", "\n\n")
	s = reHTMLTags.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	lines := strings.Split(s, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return strings.Join(cleaned, "\n")
}

func fetchURLString(targetURL string) (string, error) {
	client := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "TriGameDevBot/1.0 (+https://trigamedev.org)")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("itch.io returned HTTP %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func scrapeItchURL(rawURL string) (*ItchScrapedData, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, fmt.Errorf("invalid URL")
	}
	if !strings.HasSuffix(strings.ToLower(parsed.Host), "itch.io") {
		return nil, fmt.Errorf("URL must be an itch.io jam or game link")
	}

	cleanURL := fmt.Sprintf("https://%s%s", parsed.Host, strings.TrimRight(parsed.Path, "/"))
	isJam := strings.HasPrefix(parsed.Path, "/jam/")

	if isJam {
		return scrapeItchJam(cleanURL)
	}
	return scrapeItchGame(cleanURL)
}

func scrapeItchJam(jamURL string) (*ItchScrapedData, error) {
	rawHTML, err := fetchURLString(jamURL)
	if err != nil {
		return nil, err
	}
	meta := extractMetaMap(rawHTML)

	out := &ItchScrapedData{
		Kind:         "jam",
		URL:          jamURL,
		SecondaryURL: jamURL,
		Title:        meta["og:title"],
		Description:  meta["og:description"],
		ImageURL:     meta["og:image"],
		Status:       "active",
		Tag:          time.Now().Format("January 2006"),
		Subtitle:     "Open Theme",
		StartDate:    time.Now().Format(time.RFC3339),
		EndDate:      time.Now().Add(72 * time.Hour).Format(time.RFC3339),
	}
	if out.ImageURL == "" {
		out.ImageURL = meta["twitter:image"]
	}

	// Parse ViewJam JSON for exact start_date, end_date, voting_end_date
	if m := reViewJam.FindStringSubmatch(rawHTML); len(m) > 1 {
		var vj struct {
			StartDate     string `json:"start_date"`
			EndDate       string `json:"end_date"`
			VotingEndDate string `json:"voting_end_date"`
		}
		if err := json.Unmarshal([]byte(m[1]), &vj); err == nil {
			const layout = "2006-01-02 15:04:05"
			now := time.Now().UTC()
			if st, err := time.Parse(layout, vj.StartDate); err == nil {
				out.StartDate = st.UTC().Format(time.RFC3339)
				out.Tag = st.Format("January 2006")
				if now.Before(st) {
					out.Status = "upcoming"
				}
			}
			if et, err := time.Parse(layout, vj.EndDate); err == nil {
				out.EndDate = et.UTC().Format(time.RFC3339)
				if now.After(et) {
					out.Status = "completed"
					if vj.VotingEndDate != "" {
						if vt, err := time.Parse(layout, vj.VotingEndDate); err == nil && now.Before(vt) {
							out.Status = "voting"
						}
					}
				}
			}
		}
	}

	// Parse stat boxes (Entries, Joined, Ratings)
	for _, sm := range reStatBox.FindAllStringSubmatch(rawHTML, -1) {
		valStr := strings.ReplaceAll(sm[1], ",", "")
		val, _ := strconv.Atoi(valStr)
		label := strings.ToLower(strings.TrimSpace(sm[2]))
		switch label {
		case "entries", "submissions":
			out.SubmissionsCount = val
		case "joined", "participants", "ratings":
			if out.ParticipantsCount == 0 || label == "joined" {
				out.ParticipantsCount = val
			}
		}
	}

	// Parse jam body content & look for "Theme:" if present
	if bm := reJamBody.FindStringSubmatch(rawHTML); len(bm) > 1 {
		bodyText := cleanHTMLText(bm[1])
		if bodyText != "" {
			out.Description = bodyText
			lines := strings.Split(bodyText, "\n")
			for i, line := range lines {
				low := strings.ToLower(strings.TrimSpace(line))
				if strings.HasPrefix(low, "theme:") {
					rest := strings.TrimSpace(line[len("theme:"):])
					if rest != "" {
						out.Subtitle = rest
						break
					} else if i+1 < len(lines) {
						out.Subtitle = strings.TrimSpace(lines[i+1])
						break
					}
				}
			}
		}
	}

	return out, nil
}

func scrapeItchGame(gameURL string) (*ItchScrapedData, error) {
	rawHTML, err := fetchURLString(gameURL)
	if err != nil {
		return nil, err
	}
	meta := extractMetaMap(rawHTML)

	out := &ItchScrapedData{
		Kind:         "game",
		URL:          gameURL,
		SecondaryURL: gameURL,
		Title:        meta["og:title"],
		Description:  meta["og:description"],
		ImageURL:     meta["og:image"],
		Tag:          "Indie",
	}
	if out.Description == "" {
		out.Description = meta["description"]
	}

	// Fetch <gameURL>/data.json for rich structured authors, title, cover, tags
	if dataJSONStr, err := fetchURLString(gameURL + "/data.json"); err == nil {
		var dj struct {
			Title      string `json:"title"`
			CoverImage string `json:"cover_image"`
			Authors    []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"authors"`
			Tags []string `json:"tags"`
		}
		if json.Unmarshal([]byte(dataJSONStr), &dj) == nil {
			if dj.Title != "" {
				out.Title = dj.Title
			}
			if out.ImageURL == "" && dj.CoverImage != "" {
				out.ImageURL = dj.CoverImage
			}
			authors := []string{}
			for _, a := range dj.Authors {
				if a.Name != "" {
					authors = append(authors, a.Name)
				}
			}
			if len(authors) > 0 {
				out.Subtitle = strings.Join(authors, ", ")
			}
			if len(dj.Tags) > 0 && out.Tag == "Indie" {
				out.Tag = dj.Tags[0]
			}
		}
	}

	// Parse info table rows (Made with, Genre, Authors, Links)
	var madeWith, genre string
	for _, rm := range reTableRow.FindAllStringSubmatch(rawHTML, -1) {
		label := strings.ToLower(strings.TrimSpace(rm[1]))
		cellHTML := rm[2]

		anchors := []string{}
		for _, am := range reAnchor.FindAllStringSubmatch(cellHTML, -1) {
			txt := strings.TrimSpace(reHTMLTags.ReplaceAllString(am[1], ""))
			if txt != "" {
				anchors = append(anchors, html.UnescapeString(txt))
			}
		}

		switch label {
		case "made with":
			if len(anchors) > 0 {
				madeWith = strings.Join(anchors, ", ")
			}
		case "genre":
			if len(anchors) > 0 {
				genre = strings.Join(anchors, ", ")
			}
		case "authors":
			if out.Subtitle == "" && len(anchors) > 0 {
				out.Subtitle = strings.Join(anchors, ", ")
			}
		case "links":
			if hm := reHref.FindStringSubmatch(cellHTML); len(hm) > 1 {
				out.SecondaryURL = html.UnescapeString(hm[1])
			}
		}
	}

	if madeWith != "" {
		out.Tag = madeWith
	} else if genre != "" {
		out.Tag = genre
	}

	if out.Title == "" {
		out.Title = strings.TrimSuffix(meta["twitter:title"], " - itch.io")
	}

	return out, nil
}

func (s *Server) handleItchScrape(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" && r.Method == http.MethodPost {
		var body struct {
			URL string `json:"url"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		targetURL = body.URL
	}
	if strings.TrimSpace(targetURL) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Please provide an itch.io URL"})
		return
	}

	data, err := scrapeItchURL(targetURL)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, data)
}
