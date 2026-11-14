package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type Server struct {
	db                   *sql.DB
	sessionSecret        string
	storageDir           string
	mu                   sync.Mutex
	cachedDiscordMembers int
	lastDiscordCheck     time.Time
	lastEventsSync       time.Time
}

type User struct {
	ID            int      `json:"id"`
	DiscordID     string   `json:"discord_id"`
	Username      string   `json:"username"`
	DisplayName   string   `json:"display_name"`
	AvatarURL     string   `json:"avatar_url"`
	IsGuildMember bool     `json:"is_guild_member"`
	IsOGAdmin     bool     `json:"is_og_admin"`
	Roles         []string `json:"roles"`
}

var reUnsafeFilename = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func sanitizeFilename(name string) string {
	base := filepath.Base(name)
	clean := reUnsafeFilename.ReplaceAllString(base, "_")
	if clean == "" || clean == "." || clean == "_" {
		clean = "file"
	}
	return clean
}

func formatByteSize(b int64) string {
	if b >= 1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(b)/(1024.0*1024.0))
	}
	if b >= 1024 {
		return fmt.Sprintf("%d KB", b/1024)
	}
	return fmt.Sprintf("%d B", b)
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "trigamedev-default-secret-key-change-in-prod"
	}
	storageDir := os.Getenv("STORAGE_DIR")
	if storageDir == "" {
		storageDir = "/srv/storage/trigamedev"
	}
	_ = os.MkdirAll(filepath.Join(storageDir, "ebooks"), 0755)
	_ = os.MkdirAll(filepath.Join(storageDir, "images"), 0755)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to configure database: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	s := &Server{
		db:            db,
		sessionSecret: secret,
		storageDir:    storageDir,
	}

	// Ensure schema migrations exist
	s.ensureMigrations()

	mux := http.NewServeMux()

	// Root & Health
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/", s.handleRoot)

	// Static serving for uploaded images
	mux.HandleFunc("/uploads/images/", s.handleServeUploadedImage)

	// Auth & Discord OAuth (supports both /auth/discord/callback and /auth/discort/callback)
	mux.HandleFunc("/auth/me", s.handleAuthMe)
	mux.HandleFunc("/auth/discord/login", s.handleDiscordLogin)
	mux.HandleFunc("/auth/discord/callback", s.handleDiscordCallback)
	mux.HandleFunc("/auth/discort/callback", s.handleDiscordCallback)
	mux.HandleFunc("/auth/logout", s.handleLogout)

	// Public Community Hub Endpoints & Unread Notification Tracking
	mux.HandleFunc("/tgd/overview", s.handleOverview)
	mux.HandleFunc("/tgd/notifications/read", s.handleMarkSectionRead)

	// Image Upload Endpoint (for Key Donations, Ebooks, and CMS)
	mux.HandleFunc("/tgd/upload/image", s.handleUploadImage)

	// Member Key Sharing Portal
	mux.HandleFunc("/tgd/keys/pools", s.handleKeyPools)
	mux.HandleFunc("/tgd/keys/claim", s.handleClaimKey)
	mux.HandleFunc("/tgd/keys/my-keys", s.handleMyKeys)

	// Member Ebook Library & Book Cover/Metadata Scraper
	mux.HandleFunc("/tgd/library/ebooks", s.handleEbooks)
	mux.HandleFunc("/tgd/library/download", s.handleEbookDownload)
	mux.HandleFunc("/tgd/library/book-scrape", s.handleBookScrape)

	// OG Admin CMS, Itch.io One-Shot Scraper & Discord Sync
	mux.HandleFunc("/tgd/cms/item", s.handleCMSItem)
	mux.HandleFunc("/tgd/cms/itch-scrape", s.handleItchScrape)
	mux.HandleFunc("/tgd/cms/discord-sync", s.handleDiscordSync)

	log.Printf("Starting Go API server on :%s", port)
	if err := http.ListenAndServe(":"+port, s.corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) ensureMigrations() {
	stmts := []string{
		`ALTER TABLE trigamedev.users ADD COLUMN IF NOT EXISTS roles_json TEXT DEFAULT '[]'`,
		`ALTER TABLE trigamedev.game_jams ADD COLUMN IF NOT EXISTS jam_type VARCHAR(40) DEFAULT 'monthly'`,
		`ALTER TABLE trigamedev.game_jams ADD COLUMN IF NOT EXISTS winners_json TEXT DEFAULT '[]'`,
		`ALTER TABLE trigamedev.ebooks ADD COLUMN IF NOT EXISTS uploaded_by INT REFERENCES trigamedev.users(id) ON DELETE SET NULL`,
		`CREATE TABLE IF NOT EXISTS trigamedev.user_section_reads (
			user_id INT NOT NULL REFERENCES trigamedev.users(id) ON DELETE CASCADE,
			section VARCHAR(40) NOT NULL,
			last_read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, section)
		)`,
		`UPDATE trigamedev.users SET roles_json = '["OG"]' WHERE is_og_admin = TRUE AND (roles_json IS NULL OR roles_json = '[]' OR roles_json = '')`,
	}
	for _, q := range stmts {
		if _, err := s.db.Exec(q); err != nil {
			log.Printf("Migration note: %v", err)
		}
	}
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Demo-Role")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Sign & verify stateless HMAC token for user ID with configurable TTL
func (s *Server) signToken(userID int, ttl time.Duration) string {
	payload := fmt.Sprintf("%d.%d", userID, time.Now().Add(ttl).Unix())
	mac := hmac.New(sha256.New, []byte(s.sessionSecret))
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return payload + "." + sig
}

func (s *Server) verifyToken(token string) (int, bool) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return 0, false
	}
	payload := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(s.sessionSecret))
	mac.Write([]byte(payload))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return 0, false
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return 0, false
	}
	uid, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}
	return uid, true
}

func (s *Server) currentUser(r *http.Request) *User {
	// 1. Check Authorization: Bearer <token> from localStorage
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		tok := strings.TrimSpace(authHeader[7:])
		if uid, ok := s.verifyToken(tok); ok {
			if u := s.getUserByID(uid); u != nil {
				return u
			}
		}
	}
	// 2. Fallback to HttpOnly session cookie
	if c, err := r.Cookie("tgd_session"); err == nil {
		if uid, ok := s.verifyToken(c.Value); ok {
			return s.getUserByID(uid)
		}
	}
	return nil
}

func (s *Server) getUserByID(id int) *User {
	var u User
	var rolesRaw sql.NullString
	err := s.db.QueryRow(`
		SELECT id, discord_id, username, display_name, avatar_url, is_guild_member, is_og_admin, COALESCE(roles_json, '[]')
		FROM trigamedev.users WHERE id = $1
	`, id).Scan(&u.ID, &u.DiscordID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.IsGuildMember, &u.IsOGAdmin, &rolesRaw)
	if err != nil {
		return nil
	}
	roles := []string{}
	if rolesRaw.Valid && rolesRaw.String != "" {
		_ = json.Unmarshal([]byte(rolesRaw.String), &roles)
	}
	if len(roles) == 0 {
		if u.IsOGAdmin {
			roles = []string{"OG"}
		} else if u.IsGuildMember {
			roles = []string{"Member"}
		}
	}
	u.Roles = roles
	return &u
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	var count int
	err := s.db.QueryRow("SELECT NOW(), (SELECT COUNT(*) FROM shared.service_registry)").Scan(&now, &count)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"service":  "api-go",
			"status":   "error",
			"database": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"service":             "api-go",
		"status":              "ok",
		"database":            "connected",
		"db_time":             now.Format(time.RFC3339),
		"registered_services": count,
	})
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": "api-go",
		"message": "TriGameDev & Shared Go Backend API is live",
		"path":    r.URL.Path,
	})
}

func (s *Server) handleServeUploadedImage(w http.ResponseWriter, r *http.Request) {
	fname := filepath.Base(strings.TrimPrefix(r.URL.Path, "/uploads/images/"))
	if fname == "" || fname == "." || fname == "/" {
		http.NotFound(w, r)
		return
	}
	fullPath := filepath.Join(s.storageDir, "images", fname)
	if _, err := os.Stat(fullPath); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeFile(w, r, fullPath)
}

func (s *Server) handleUploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "POST required"})
		return
	}
	u := s.currentUser(r)
	if u == nil || !u.IsGuildMember {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Must be signed in to upload images"})
		return
	}

	if err := r.ParseMultipartForm(15 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Image must be under 15 MB"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing 'file' field"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true, ".svg": true}
	if !allowed[ext] {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Allowed image formats: PNG, JPG, GIF, WEBP, SVG"})
		return
	}

	safeName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), sanitizeFilename(header.Filename))
	dstPath := filepath.Join(s.storageDir, "images", safeName)
	dst, err := os.Create(dstPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save image"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, io.LimitReader(file, 15<<20)); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to write image"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"url": "/uploads/images/" + safeName,
	})
}

func (s *Server) handleAuthMe(w http.ResponseWriter, r *http.Request) {
	u := s.currentUser(r)
	discordConfigured := os.Getenv("DISCORD_CLIENT_ID") != ""
	writeJSON(w, http.StatusOK, map[string]any{
		"authenticated":      u != nil,
		"user":               u,
		"discord_configured": discordConfigured,
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "tgd_session", Value: "", Path: "/", MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "tgd_demo_role", Value: "", Path: "/", MaxAge: -1})
	writeJSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

// Discord OAuth2 Federated Login + Guild & Role Verification
func (s *Server) handleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("DISCORD_CLIENT_ID")
	redirectURI := os.Getenv("DISCORD_REDIRECT_URI")
	if clientID == "" {
		http.Redirect(w, r, "/?error=discord_oauth_not_configured", http.StatusFound)
		return
	}
	remember := r.URL.Query().Get("remember")
	state := "remember"
	if remember == "0" || remember == "false" {
		state = "session"
	}
	u := fmt.Sprintf(
		"https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20guilds.members.read&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)
	http.Redirect(w, r, u, http.StatusFound)
}

func (s *Server) resolveGuildRoleNames(guildID, ogRoleID string, roleIDs []string) []string {
	roleMap := map[string]string{}
	if ogRoleID != "" {
		roleMap[ogRoleID] = "OG"
	}
	// Optional comma-separated role map from env: DISCORD_ROLE_MAP="123:Admin,456:Jam Winner"
	if customMap := os.Getenv("DISCORD_ROLE_MAP"); customMap != "" {
		for _, pair := range strings.Split(customMap, ",") {
			if kv := strings.SplitN(strings.TrimSpace(pair), ":", 2); len(kv) == 2 {
				roleMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
	}
	// If DISCORD_BOT_TOKEN is configured, query live Guild Roles from Discord
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken != "" && guildID != "" {
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/roles", guildID), nil)
		req.Header.Set("Authorization", "Bot "+botToken)
		client := &http.Client{Timeout: 6 * time.Second}
		if resp, err := client.Do(req); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var gRoles []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}
				if json.NewDecoder(resp.Body).Decode(&gRoles) == nil {
					for _, gr := range gRoles {
						if gr.Name != "@everyone" {
							roleMap[gr.ID] = gr.Name
						}
					}
				}
			}
		}
	}

	names := []string{}
	seen := map[string]bool{}
	for _, rID := range roleIDs {
		if name, ok := roleMap[rID]; ok && name != "" && !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		names = append(names, "Member")
	}
	return names
}

func (s *Server) handleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	remember := state != "session"

	clientID := os.Getenv("DISCORD_CLIENT_ID")
	clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	redirectURI := os.Getenv("DISCORD_REDIRECT_URI")
	guildID := os.Getenv("DISCORD_GUILD_ID")
	ogRoleID := os.Getenv("DISCORD_OG_ROLE_ID")

	if code == "" || clientID == "" {
		http.Redirect(w, r, "/?error=missing_code", http.StatusFound)
		return
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	resp, err := http.PostForm("https://discord.com/api/oauth2/token", data)
	if err != nil {
		http.Redirect(w, r, "/?error=token_exchange", http.StatusFound)
		return
	}
	defer resp.Body.Close()

	var tokenRes struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&tokenRes)
	if tokenRes.AccessToken == "" {
		http.Redirect(w, r, "/?error=empty_token", http.StatusFound)
		return
	}

	// Fetch Discord User Profile
	req, _ := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)
	uResp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Redirect(w, r, "/?error=profile_fetch", http.StatusFound)
		return
	}
	defer uResp.Body.Close()

	var dUser struct {
		ID         string `json:"id"`
		Username   string `json:"username"`
		GlobalName string `json:"global_name"`
		Avatar     string `json:"avatar"`
	}
	_ = json.NewDecoder(uResp.Body).Decode(&dUser)

	isMember := true
	isOG := false
	roleNames := []string{"Member"}

	// Verify Guild Membership & Roles in TriGameDev Discord Server
	if guildID != "" {
		mReq, _ := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/users/@me/guilds/%s/member", guildID), nil)
		mReq.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)
		mResp, err := http.DefaultClient.Do(mReq)
		if err == nil {
			defer mResp.Body.Close()
			if mResp.StatusCode == http.StatusOK {
				isMember = true
				var mData struct {
					Roles []string `json:"roles"`
				}
				_ = json.NewDecoder(mResp.Body).Decode(&mData)
				for _, rID := range mData.Roles {
					if ogRoleID != "" && rID == ogRoleID {
						isOG = true
					}
				}
				roleNames = s.resolveGuildRoleNames(guildID, ogRoleID, mData.Roles)
			} else {
				isMember = false
				roleNames = []string{}
			}
		}
	}

	rolesBytes, _ := json.Marshal(roleNames)

	displayName := dUser.GlobalName
	if displayName == "" {
		displayName = dUser.Username
	}
	avatarURL := fmt.Sprintf("https://api.dicebear.com/9.x/bottts-neutral/svg?seed=%s", dUser.Username)
	if dUser.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", dUser.ID, dUser.Avatar)
	}

	var userID int
	err = s.db.QueryRow(`
		INSERT INTO trigamedev.users (discord_id, username, display_name, avatar_url, is_guild_member, is_og_admin, roles_json, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (discord_id) DO UPDATE
		SET username = EXCLUDED.username,
		    display_name = EXCLUDED.display_name,
		    avatar_url = EXCLUDED.avatar_url,
		    is_guild_member = EXCLUDED.is_guild_member,
		    is_og_admin = EXCLUDED.is_og_admin,
		    roles_json = EXCLUDED.roles_json,
		    last_login_at = NOW()
		RETURNING id
	`, dUser.ID, dUser.Username, displayName, avatarURL, isMember, isOG, string(rolesBytes)).Scan(&userID)
	if err != nil {
		http.Redirect(w, r, "/?error=db_upsert", http.StatusFound)
		return
	}

	if remember {
		tok := s.signToken(userID, 90*24*time.Hour)
		http.SetCookie(w, &http.Cookie{
			Name:     "tgd_session",
			Value:    tok,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400 * 90,
		})
		http.Redirect(w, r, "/?auth_token="+url.QueryEscape(tok), http.StatusFound)
		return
	}

	tok := s.signToken(userID, 24*time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "tgd_session",
		Value:    tok,
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/?session_only=1", http.StatusFound)
}

func (s *Server) getDiscordMemberCount() int {
	s.mu.Lock()
	if s.cachedDiscordMembers > 0 && time.Since(s.lastDiscordCheck) < 10*time.Minute {
		val := s.cachedDiscordMembers
		s.mu.Unlock()
		return val
	}
	s.mu.Unlock()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://discord.com/api/v10/invites/6jdhMV6sqt?with_counts=true")
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var inv struct {
				ApproxMembers int `json:"approximate_member_count"`
			}
			if json.NewDecoder(resp.Body).Decode(&inv) == nil && inv.ApproxMembers > 0 {
				s.mu.Lock()
				s.cachedDiscordMembers = inv.ApproxMembers
				s.lastDiscordCheck = time.Now()
				s.mu.Unlock()
				return inv.ApproxMembers
			}
		}
	}

	s.mu.Lock()
	fallback := s.cachedDiscordMembers
	s.mu.Unlock()
	if fallback > 0 {
		return fallback
	}
	return 162
}

func (s *Server) syncDiscordEventsFromGuild(force bool) int {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	guildID := os.Getenv("DISCORD_GUILD_ID")
	if botToken == "" || guildID == "" {
		return 0
	}

	s.mu.Lock()
	if !force && time.Since(s.lastEventsSync) < 5*time.Minute {
		s.mu.Unlock()
		return 0
	}
	s.lastEventsSync = time.Now()
	s.mu.Unlock()

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/scheduled-events?with_user_count=true", guildID), nil)
	req.Header.Set("Authorization", "Bot "+botToken)
	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0
	}
	defer resp.Body.Close()

	var dEvents []struct {
		ID             string    `json:"id"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		ScheduledStart time.Time `json:"scheduled_start_time"`
		UserCount      int       `json:"user_count"`
		Image          string    `json:"image"`
		EntityMetadata struct {
			Location string `json:"location"`
		} `json:"entity_metadata"`
		RecurrenceRule *struct {
			Start time.Time `json:"start"`
		} `json:"recurrence_rule"`
	}
	if json.NewDecoder(resp.Body).Decode(&dEvents) != nil {
		return 0
	}

	synced := 0
	for _, ev := range dEvents {
		loc := ev.EntityMetadata.Location
		if loc == "" {
			loc = "TriGameDev Discord"
		}
		cover := ""
		if ev.Image != "" {
			cover = fmt.Sprintf("https://cdn.discordapp.com/guild-events/%s/%s.png?size=1024", ev.ID, ev.Image)
		}
		_, err := s.db.Exec(`
			INSERT INTO trigamedev.events (discord_event_id, title, description, location, start_time, cover_image_url, rsvp_count, tag)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'Discord Event')
			ON CONFLICT (discord_event_id) DO UPDATE
			SET title = EXCLUDED.title,
			    description = EXCLUDED.description,
			    location = EXCLUDED.location,
			    start_time = EXCLUDED.start_time,
			    cover_image_url = CASE WHEN EXCLUDED.cover_image_url != '' THEN EXCLUDED.cover_image_url ELSE trigamedev.events.cover_image_url END,
			    rsvp_count = EXCLUDED.rsvp_count
		`, ev.ID, ev.Name, ev.Description, loc, ev.ScheduledStart, cover, ev.UserCount)
		if err == nil {
			synced++
		}
		// If this is a recurring Discord event whose initial recurrence start was in the past, preserve that past occurrence too
		if ev.RecurrenceRule != nil && !ev.RecurrenceRule.Start.IsZero() && ev.RecurrenceRule.Start.Before(ev.ScheduledStart) && ev.RecurrenceRule.Start.Before(time.Now()) {
			pastID := ev.ID + "_" + ev.RecurrenceRule.Start.UTC().Format("20060102")
			_, _ = s.db.Exec(`
				INSERT INTO trigamedev.events (discord_event_id, title, description, location, start_time, cover_image_url, rsvp_count, tag)
				VALUES ($1, $2, $3, $4, $5, $6, $7, 'Discord Event')
				ON CONFLICT (discord_event_id) DO NOTHING
			`, pastID, ev.Name, ev.Description, loc, ev.RecurrenceRule.Start, cover, ev.UserCount)
		}
	}
	return synced
}

func (s *Server) handleMarkSectionRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "POST required"})
		return
	}
	u := s.currentUser(r)
	if u == nil {
		writeJSON(w, http.StatusOK, map[string]string{"status": "guest"})
		return
	}
	var req struct {
		Section string `json:"section"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	sec := strings.TrimSpace(req.Section)
	if sec == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "section required"})
		return
	}
	_, _ = s.db.Exec(`
		INSERT INTO trigamedev.user_section_reads (user_id, section, last_read_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, section) DO UPDATE SET last_read_at = NOW()
	`, u.ID, sec)
	writeJSON(w, http.StatusOK, map[string]string{"status": "read", "section": sec})
}

func (s *Server) computeUnreadCounts(userID int) map[string]int {
	reads := map[string]time.Time{}
	rows, err := s.db.Query(`SELECT section, last_read_at FROM trigamedev.user_section_reads WHERE user_id = $1`, userID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var sec string
			var t time.Time
			if rows.Scan(&sec, &t) == nil {
				reads[sec] = t
			}
		}
	}

	epoch := time.Unix(0, 0).UTC()
	getSince := func(sec string) time.Time {
		if t, ok := reads[sec]; ok {
			return t
		}
		return epoch
	}

	unread := map[string]int{}
	var cJams, cEvents, cNews, cShowcase, cGallery, cVideos, cKeys, cLibrary int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.game_jams WHERE created_at > $1`, getSince("jams")).Scan(&cJams)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.events WHERE created_at > $1`, getSince("events")).Scan(&cEvents)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.news WHERE published_at > $1`, getSince("news")).Scan(&cNews)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.showcase WHERE created_at > $1`, getSince("showcase")).Scan(&cShowcase)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.gallery WHERE created_at > $1`, getSince("media")).Scan(&cGallery)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.videos WHERE recorded_at > $1`, getSince("media")).Scan(&cVideos)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.key_pools WHERE created_at > $1`, getSince("keys")).Scan(&cKeys)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.ebooks WHERE created_at > $1`, getSince("library")).Scan(&cLibrary)

	unread["jams"] = cJams
	unread["events"] = cEvents
	unread["news"] = cNews
	unread["showcase"] = cShowcase
	unread["media"] = cGallery + cVideos
	unread["keys"] = cKeys
	unread["library"] = cLibrary
	return unread
}

func (s *Server) handleOverview(w http.ResponseWriter, r *http.Request) {
	// Opportunistically sync Discord events in background if bot token is configured
	go s.syncDiscordEventsFromGuild(false)

	// 1. Events (both upcoming and past, ordered newest first)
	events := []map[string]any{}
	eRows, err := s.db.Query(`
		SELECT id, title, description, location, start_time, cover_image_url, tag, rsvp_count
		FROM trigamedev.events ORDER BY start_time DESC LIMIT 40
	`)
	if err == nil {
		defer eRows.Close()
		for eRows.Next() {
			var id, rsvp int
			var title, desc, loc, img, tag string
			var st time.Time
			if err := eRows.Scan(&id, &title, &desc, &loc, &st, &img, &tag, &rsvp); err == nil {
				events = append(events, map[string]any{
					"id": id, "title": title, "description": desc, "location": loc,
					"start_time": st.Format(time.RFC3339), "cover_image_url": img, "tag": tag, "rsvp_count": rsvp,
				})
			}
		}
	}

	// 2. News
	news := []map[string]any{}
	nRows, err := s.db.Query(`
		SELECT id, title, summary, body_markdown, author_name, author_role, tag, is_pinned, published_at
		FROM trigamedev.news ORDER BY is_pinned DESC, published_at DESC LIMIT 20
	`)
	if err == nil {
		defer nRows.Close()
		for nRows.Next() {
			var id int
			var title, summary, body, author, role, tag string
			var pinned bool
			var pub time.Time
			if err := nRows.Scan(&id, &title, &summary, &body, &author, &role, &tag, &pinned, &pub); err == nil {
				news = append(news, map[string]any{
					"id": id, "title": title, "summary": summary, "body_markdown": body,
					"author_name": author, "author_role": role, "tag": tag, "is_pinned": pinned,
					"published_at": pub.Format("Jan 02, 2006"),
				})
			}
		}
	}

	// 3. Game Jams (including winners_json and jam_type)
	jams := []map[string]any{}
	jRows, err := s.db.Query(`
		SELECT id, title, month_label, theme, description, itch_url, cover_image_url, start_date, end_date, status, participants_count, submissions_count, COALESCE(jam_type, 'monthly'), COALESCE(winners_json, '[]')
		FROM trigamedev.game_jams ORDER BY start_date DESC LIMIT 30
	`)
	if err == nil {
		defer jRows.Close()
		for jRows.Next() {
			var id, parts, subs int
			var title, month, theme, desc, itch, img, status, jamType, winnersRaw string
			var sd, ed time.Time
			if err := jRows.Scan(&id, &title, &month, &theme, &desc, &itch, &img, &sd, &ed, &status, &parts, &subs, &jamType, &winnersRaw); err == nil {
				var winners []ItchJamWinner
				_ = json.Unmarshal([]byte(winnersRaw), &winners)
				if winners == nil {
					winners = []ItchJamWinner{}
				}
				jams = append(jams, map[string]any{
					"id": id, "title": title, "month_label": month, "theme": theme,
					"description": desc, "itch_url": itch, "cover_image_url": img,
					"start_date": sd.Format(time.RFC3339), "end_date": ed.Format(time.RFC3339),
					"status": status, "participants_count": parts, "submissions_count": subs,
					"jam_type": jamType, "winners": winners,
				})
			}
		}
	}

	// 4. Showcase ("Games by Our People") — randomized order
	showcase := []map[string]any{}
	sRows, err := s.db.Query(`
		SELECT id, title, developer_name, description, thumbnail_url, play_url, source_url, engine, views_count, featured
		FROM trigamedev.showcase ORDER BY RANDOM() LIMIT 24
	`)
	if err == nil {
		defer sRows.Close()
		for sRows.Next() {
			var id, views int
			var title, dev, desc, thumb, play, src, engine string
			var feat bool
			if err := sRows.Scan(&id, &title, &dev, &desc, &thumb, &play, &src, &engine, &views, &feat); err == nil {
				showcase = append(showcase, map[string]any{
					"id": id, "title": title, "developer_name": dev, "description": desc,
					"thumbnail_url": thumb, "play_url": play, "source_url": src,
					"engine": engine, "views_count": views, "featured": feat,
				})
			}
		}
	}

	// 5. Gallery — randomized order
	gallery := []map[string]any{}
	gRows, err := s.db.Query(`
		SELECT id, title, caption, image_url, event_label, uploaded_by
		FROM trigamedev.gallery ORDER BY RANDOM() LIMIT 24
	`)
	if err == nil {
		defer gRows.Close()
		for gRows.Next() {
			var id int
			var title, cap, img, ev, up string
			if err := gRows.Scan(&id, &title, &cap, &img, &ev, &up); err == nil {
				gallery = append(gallery, map[string]any{
					"id": id, "title": title, "caption": cap, "image_url": img,
					"event_label": ev, "uploaded_by": up,
				})
			}
		}
	}

	// 6. Videos — randomized order
	videos := []map[string]any{}
	vRows, err := s.db.Query(`
		SELECT id, title, description, video_url, thumbnail_url, speaker_name, duration
		FROM trigamedev.videos ORDER BY RANDOM() LIMIT 24
	`)
	if err == nil {
		defer vRows.Close()
		for vRows.Next() {
			var id int
			var title, desc, vurl, thumb, spk, dur string
			if err := vRows.Scan(&id, &title, &desc, &vurl, &thumb, &spk, &dur); err == nil {
				videos = append(videos, map[string]any{
					"id": id, "title": title, "description": desc, "video_url": vurl,
					"thumbnail_url": thumb, "speaker_name": spk, "duration": dur,
				})
			}
		}
	}

	// Stats summary
	var totalKeysAvailable, totalProjects, totalEbooks, totalShowcase, totalJams int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.keys WHERE claimed_by IS NULL`).Scan(&totalKeysAvailable)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.key_pools`).Scan(&totalProjects)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.ebooks`).Scan(&totalEbooks)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.showcase`).Scan(&totalShowcase)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM trigamedev.game_jams`).Scan(&totalJams)

	unread := map[string]int{}
	if u := s.currentUser(r); u != nil {
		unread = s.computeUnreadCounts(u.ID)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"events":   events,
		"news":     news,
		"jams":     jams,
		"showcase": showcase,
		"gallery":  gallery,
		"videos":   videos,
		"unread":   unread,
		"stats": map[string]any{
			"total_members":  s.getDiscordMemberCount(),
			"total_jams":     totalJams,
			"released_games": totalShowcase,
			"available_keys": totalKeysAvailable,
			"total_projects": totalProjects,
			"total_ebooks":   totalEbooks,
			"total_showcase": totalShowcase,
		},
	})
}

// Member Key Sharing Portal: List Projects, Donate Keys, or Remove Own Project
func (s *Server) handleKeyPools(w http.ResponseWriter, r *http.Request) {
	u := s.currentUser(r)

	if r.Method == http.MethodGet {
		uid := 0
		if u != nil {
			uid = u.ID
		}
		rows, err := s.db.Query(`
			SELECT
				p.id, p.title, p.platform, p.store_url, p.cover_image_url, p.description, p.donated_by_name, COALESCE(p.created_by, 0),
				COUNT(k.id) FILTER (WHERE k.claimed_by IS NULL) AS available_count,
				COUNT(k.id) AS total_count,
				BOOL_OR(k.claimed_by = $1 AND $1 > 0) AS claimed_by_me
			FROM trigamedev.key_pools p
			LEFT JOIN trigamedev.keys k ON k.pool_id = p.id
			GROUP BY p.id
			ORDER BY available_count DESC, p.id DESC
		`, uid)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		defer rows.Close()

		pools := []map[string]any{}
		for rows.Next() {
			var id, createdBy, avail, total int
			var title, platform, store, cover, desc, donor string
			var claimedByMe sql.NullBool
			if err := rows.Scan(&id, &title, &platform, &store, &cover, &desc, &donor, &createdBy, &avail, &total, &claimedByMe); err == nil {
				canDelete := u != nil && (u.IsOGAdmin || (createdBy > 0 && createdBy == u.ID) || strings.EqualFold(donor, u.DisplayName))
				pools = append(pools, map[string]any{
					"id":              id,
					"title":           title,
					"platform":        platform,
					"store_url":       store,
					"cover_image_url": cover,
					"description":     desc,
					"donated_by_name": donor,
					"created_by":      createdBy,
					"can_delete":      canDelete,
					"available_count": avail,
					"total_count":     total,
					"claimed_by_user": claimedByMe.Bool,
				})
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"pools": pools})
		return
	}

	if r.Method == http.MethodDelete {
		if u == nil {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Sign in required"})
			return
		}
		poolID, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if poolID <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Valid project id required"})
			return
		}
		var createdBy int
		var donor string
		err := s.db.QueryRow(`SELECT COALESCE(created_by, 0), donated_by_name FROM trigamedev.key_pools WHERE id = $1`, poolID).Scan(&createdBy, &donor)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
			return
		}
		if !u.IsOGAdmin && createdBy != u.ID && !strings.EqualFold(donor, u.DisplayName) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Only the member who shared these keys can remove this project"})
			return
		}
		_, _ = s.db.Exec(`DELETE FROM trigamedev.key_pools WHERE id = $1`, poolID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
		return
	}

	if r.Method == http.MethodPost {
		if u == nil || !u.IsGuildMember {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Must be a logged-in TriGameDev Discord member to donate keys"})
			return
		}
		var req struct {
			PoolID        int      `json:"pool_id"`
			Title         string   `json:"title"`
			Platform      string   `json:"platform"`
			StoreURL      string   `json:"store_url"`
			CoverImageURL string   `json:"cover_image_url"`
			Description   string   `json:"description"`
			Keys          []string `json:"keys"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
			return
		}

		cleanKeys := []string{}
		for _, k := range req.Keys {
			trimmed := strings.TrimSpace(k)
			if trimmed != "" {
				cleanKeys = append(cleanKeys, trimmed)
			}
		}
		if len(cleanKeys) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Please provide at least one key code"})
			return
		}

		tx, err := s.db.Begin()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		defer tx.Rollback()

		poolID := req.PoolID
		if poolID == 0 {
			if req.Platform == "" {
				req.Platform = "Steam"
			}
			if req.CoverImageURL == "" {
				req.CoverImageURL = "/brand/filler.png"
			}
			err = tx.QueryRow(`
				INSERT INTO trigamedev.key_pools (title, platform, store_url, cover_image_url, description, donated_by_name, created_by)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id
			`, req.Title, req.Platform, req.StoreURL, req.CoverImageURL, req.Description, u.DisplayName, u.ID).Scan(&poolID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
		}

		for _, code := range cleanKeys {
			_, err := tx.Exec(`
				INSERT INTO trigamedev.keys (pool_id, key_code, donated_by)
				VALUES ($1, $2, $3)
			`, poolID, code, u.ID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		s.notifyDiscordWebhook(fmt.Sprintf("[KEY VAULT] **%s** donated **%d key(s)** for **%s** (%s).", u.DisplayName, len(cleanKeys), req.Title, req.Platform))

		writeJSON(w, http.StatusCreated, map[string]any{
			"status":     "donated",
			"pool_id":    poolID,
			"keys_added": len(cleanKeys),
		})
	}
}

// Single-Redemption Atomic Key Claim: ties key permanently to requester and removes from project
func (s *Server) handleClaimKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "POST required"})
		return
	}
	u := s.currentUser(r)
	if u == nil || !u.IsGuildMember {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Please log in with your TriGameDev Discord account to claim keys."})
		return
	}

	var req struct {
		PoolID int `json:"pool_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PoolID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Valid pool_id required"})
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	// Check if user already claimed a key from this project
	var existingCode string
	err = tx.QueryRow(`
		SELECT key_code FROM trigamedev.keys WHERE pool_id = $1 AND claimed_by = $2
	`, req.PoolID, u.ID).Scan(&existingCode)
	if err == nil {
		writeJSON(w, http.StatusConflict, map[string]any{
			"error":    "You have already redeemed a key from this project (limit 1 per member).",
			"key_code": existingCode,
		})
		return
	}

	// Atomically lock and claim 1 unclaimed key
	var keyID int
	var keyCode string
	err = tx.QueryRow(`
		SELECT id, key_code
		FROM trigamedev.keys
		WHERE pool_id = $1 AND claimed_by IS NULL
		ORDER BY id ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`, req.PoolID).Scan(&keyID, &keyCode)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusGone, map[string]string{"error": "All keys for this project have already been claimed!"})
		return
	} else if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_, err = tx.Exec(`
		UPDATE trigamedev.keys
		SET claimed_by = $1, claimed_at = NOW()
		WHERE id = $2
	`, u.ID, keyID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "claimed",
		"pool_id":  req.PoolID,
		"key_id":   keyID,
		"key_code": keyCode,
	})
}

// My Keys Section: shows all keys permanently claimed by the current user
func (s *Server) handleMyKeys(w http.ResponseWriter, r *http.Request) {
	u := s.currentUser(r)
	if u == nil {
		writeJSON(w, http.StatusOK, map[string]any{"claimed_keys": []any{}})
		return
	}

	rows, err := s.db.Query(`
		SELECT k.id, k.pool_id, k.key_code, k.claimed_at, p.title, p.platform, p.store_url, p.cover_image_url, p.donated_by_name
		FROM trigamedev.keys k
		JOIN trigamedev.key_pools p ON p.id = k.pool_id
		WHERE k.claimed_by = $1
		ORDER BY k.claimed_at DESC
	`, u.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	claimed := []map[string]any{}
	for rows.Next() {
		var id, poolID int
		var code, title, platform, store, cover, donor string
		var claimedAt time.Time
		if err := rows.Scan(&id, &poolID, &code, &claimedAt, &title, &platform, &store, &cover, &donor); err == nil {
			claimed = append(claimed, map[string]any{
				"id":              id,
				"pool_id":         poolID,
				"key_code":        code,
				"claimed_at":      claimedAt.Format("Jan 02, 2006 15:04"),
				"title":           title,
				"platform":        platform,
				"store_url":       store,
				"cover_image_url": cover,
				"donated_by_name": donor,
			})
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"claimed_keys": claimed})
}

// Ebook Library: Browse, Upload, or Remove Own Ebooks
func (s *Server) handleEbooks(w http.ResponseWriter, r *http.Request) {
	u := s.currentUser(r)

	if r.Method == http.MethodGet {
		rows, err := s.db.Query(`
			SELECT id, title, author, category, description, cover_image_url, file_name, file_format, file_size_label, uploaded_by_name, COALESCE(uploaded_by, 0), download_count
			FROM trigamedev.ebooks ORDER BY id DESC
		`)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		defer rows.Close()

		ebooks := []map[string]any{}
		for rows.Next() {
			var id, uploadedBy, dlCount int
			var title, author, cat, desc, cover, fname, ffmt, fsize, upBy string
			if err := rows.Scan(&id, &title, &author, &cat, &desc, &cover, &fname, &ffmt, &fsize, &upBy, &uploadedBy, &dlCount); err == nil {
				canDelete := u != nil && (u.IsOGAdmin || (uploadedBy > 0 && uploadedBy == u.ID) || strings.EqualFold(upBy, u.DisplayName))
				ebooks = append(ebooks, map[string]any{
					"id":               id,
					"title":            title,
					"author":           author,
					"category":         cat,
					"description":      desc,
					"cover_image_url":  cover,
					"file_name":        fname,
					"file_format":      ffmt,
					"file_size_label":  fsize,
					"uploaded_by_name": upBy,
					"uploaded_by":      uploadedBy,
					"can_delete":       canDelete,
					"download_count":   dlCount,
				})
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"ebooks": ebooks})
		return
	}

	if r.Method == http.MethodDelete {
		if u == nil {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Sign in required"})
			return
		}
		bookID, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if bookID <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Valid ebook id required"})
			return
		}
		var uploadedBy int
		var upBy, fileName string
		err := s.db.QueryRow(`SELECT COALESCE(uploaded_by, 0), uploaded_by_name, file_name FROM trigamedev.ebooks WHERE id = $1`, bookID).Scan(&uploadedBy, &upBy, &fileName)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Ebook not found"})
			return
		}
		if !u.IsOGAdmin && uploadedBy != u.ID && !strings.EqualFold(upBy, u.DisplayName) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Only the member who shared this ebook can remove it"})
			return
		}
		_, _ = s.db.Exec(`DELETE FROM trigamedev.ebooks WHERE id = $1`, bookID)
		if fileName != "" {
			_ = os.Remove(filepath.Join(s.storageDir, "ebooks", filepath.Base(fileName)))
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
		return
	}

	if r.Method == http.MethodPost {
		if u == nil || !u.IsGuildMember {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "Member login required"})
			return
		}

		var title, author, category, description, coverImageURL, fileName, fileFormat, fileSizeLabel string

		ct := r.Header.Get("Content-Type")
		if strings.HasPrefix(ct, "multipart/form-data") {
			if err := r.ParseMultipartForm(100 << 20); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File too large (max 100 MB)"})
				return
			}
			title = strings.TrimSpace(r.FormValue("title"))
			author = strings.TrimSpace(r.FormValue("author"))
			category = strings.TrimSpace(r.FormValue("category"))
			description = strings.TrimSpace(r.FormValue("description"))
			coverImageURL = strings.TrimSpace(r.FormValue("cover_image_url"))
			fileFormat = strings.TrimSpace(r.FormValue("file_format"))

			if file, header, err := r.FormFile("file"); err == nil {
				defer file.Close()
				safeFile := fmt.Sprintf("%d_%s", time.Now().Unix(), sanitizeFilename(header.Filename))
				dstPath := filepath.Join(s.storageDir, "ebooks", safeFile)
				dst, err := os.Create(dstPath)
				if err != nil {
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to store ebook file"})
					return
				}
				written, copyErr := io.Copy(dst, io.LimitReader(file, 100<<20))
				dst.Close()
				if copyErr != nil {
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to write ebook file"})
					return
				}
				fileName = safeFile
				fileSizeLabel = formatByteSize(written)
				ext := strings.TrimPrefix(strings.ToUpper(filepath.Ext(header.Filename)), ".")
				if ext != "" && fileFormat == "" {
					fileFormat = ext
				}
				if title == "" {
					t, a := cleanBookQuery(header.Filename)
					title = t
					if author == "" {
						author = a
					}
				}
			}
		} else {
			var req struct {
				Title         string `json:"title"`
				Author        string `json:"author"`
				Category      string `json:"category"`
				Description   string `json:"description"`
				CoverImageURL string `json:"cover_image_url"`
				FileName      string `json:"file_name"`
				FileFormat    string `json:"file_format"`
				FileSizeLabel string `json:"file_size_label"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
				return
			}
			title = req.Title
			author = req.Author
			category = req.Category
			description = req.Description
			coverImageURL = req.CoverImageURL
			fileName = req.FileName
			fileFormat = req.FileFormat
			fileSizeLabel = req.FileSizeLabel
		}

		if title == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title required"})
			return
		}
		if author == "" {
			author = "Unknown Author"
		}
		if category == "" {
			category = "Game Design"
		}
		if coverImageURL == "" {
			coverImageURL = "/brand/filler.png"
		}
		if fileFormat == "" {
			fileFormat = "PDF"
		}
		if fileName == "" {
			fileName = sanitizeFilename(strings.ToLower(title)) + "." + strings.ToLower(fileFormat)
		}
		if fileSizeLabel == "" {
			fileSizeLabel = "5.0 MB"
		}

		var id int
		err := s.db.QueryRow(`
			INSERT INTO trigamedev.ebooks (title, author, category, description, cover_image_url, file_name, file_format, file_size_label, uploaded_by_name, uploaded_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`, title, author, category, description, coverImageURL, fileName, fileFormat, fileSizeLabel, u.DisplayName, u.ID).Scan(&id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		s.notifyDiscordWebhook(fmt.Sprintf("[EBOOK LIBRARY] **%s** uploaded a new ebook: **%s** by *%s* (%s).", u.DisplayName, title, author, fileFormat))
		writeJSON(w, http.StatusCreated, map[string]any{"status": "created", "id": id})
	}
}

// Direct Ebook Download for Logged-In Members
func (s *Server) handleEbookDownload(w http.ResponseWriter, r *http.Request) {
	u := s.currentUser(r)
	if u == nil || !u.IsGuildMember {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Must be logged in with Discord as a TriGameDev member to download ebooks."})
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Valid ebook id required"})
		return
	}

	var title, author, desc, fileName, fileFormat string
	err = s.db.QueryRow(`
		UPDATE trigamedev.ebooks
		SET download_count = download_count + 1
		WHERE id = $1
		RETURNING title, author, description, file_name, file_format
	`, id).Scan(&title, &author, &desc, &fileName, &fileFormat)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Ebook not found"})
		return
	}

	diskPath := filepath.Join(s.storageDir, "ebooks", filepath.Base(fileName))
	if f, err := os.Open(diskPath); err == nil {
		defer f.Close()
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = io.Copy(w, f)
		return
	}

	content := fmt.Sprintf("====================================================================\n"+
		"TRIGAMEDEV COMMUNITY EBOOK LIBRARY - DIRECT DOWNLOAD EDITION\n"+
		"====================================================================\n\n"+
		"Title:       %s\n"+
		"Author:      %s\n"+
		"Format:      %s\n"+
		"Downloaded:  %s UTC by %s\n\n"+
		"Synopsis:\n%s\n\n"+
		"--------------------------------------------------------------------\n"+
		"TriGameDev.org Member Archive (129.121.165.252)\n",
		title, author, fileFormat, time.Now().UTC().Format(time.RFC822), u.DisplayName, desc)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName+".txt"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(content))
}

// OG Role CMS Endpoint for creating/deleting public content across all Post Types
func (s *Server) handleCMSItem(w http.ResponseWriter, r *http.Request) {
	u := s.currentUser(r)
	if u == nil || !u.IsOGAdmin {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "Requires OG Role in the TriGameDev Discord Server"})
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Section           string          `json:"section"` // jams, showcase, news, events, gallery, videos
			JamType           string          `json:"jam_type"`
			Title             string          `json:"title"`
			Subtitle          string          `json:"subtitle"`
			Description       string          `json:"description"`
			Tag               string          `json:"tag"`
			URL               string          `json:"url"`
			SecondaryURL      string          `json:"secondary_url"`
			ImageURL          string          `json:"image_url"`
			StartDate         string          `json:"start_date"`
			EndDate           string          `json:"end_date"`
			Status            string          `json:"status"`
			ParticipantsCount int             `json:"participants_count"`
			SubmissionsCount  int             `json:"submissions_count"`
			Winners           []ItchJamWinner `json:"winners"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
			return
		}

		// Auto-scrape if a Jam or Showcase item is submitted with a Steam or Itch.io URL and missing fields
		lowURL := strings.ToLower(req.URL)
		if (req.Section == "jams" || req.Section == "showcase") && (strings.Contains(lowURL, "itch.io") || strings.Contains(lowURL, "steampowered.com") || strings.Contains(lowURL, "steamcommunity.com")) && (req.Title == "" || (req.Section == "jams" && len(req.Winners) == 0)) {
			if scraped, err := scrapeItchURL(req.URL); err == nil {
				if req.Title == "" {
					req.Title = scraped.Title
				}
				if req.Subtitle == "" {
					req.Subtitle = scraped.Subtitle
				}
				if req.Description == "" {
					req.Description = scraped.Description
				}
				if req.Tag == "" {
					req.Tag = scraped.Tag
				}
				if req.ImageURL == "" {
					req.ImageURL = scraped.ImageURL
				}
				if req.SecondaryURL == "" {
					req.SecondaryURL = scraped.SecondaryURL
				}
				if req.StartDate == "" {
					req.StartDate = scraped.StartDate
				}
				if req.EndDate == "" {
					req.EndDate = scraped.EndDate
				}
				if req.Status == "" {
					req.Status = scraped.Status
				}
				if req.ParticipantsCount == 0 {
					req.ParticipantsCount = scraped.ParticipantsCount
				}
				if req.SubmissionsCount == 0 {
					req.SubmissionsCount = scraped.SubmissionsCount
				}
				if len(req.Winners) == 0 && len(scraped.Winners) > 0 {
					req.Winners = scraped.Winners
				}
			}
		}

		if req.Title == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
			return
		}
		if req.Section != "news" && req.ImageURL == "" {
			req.ImageURL = "/brand/filler.png"
		}

		var err error
		switch req.Section {
		case "news":
			if req.Tag == "" {
				req.Tag = "Announcement"
			}
			_, err = s.db.Exec(`
				INSERT INTO trigamedev.news (title, summary, body_markdown, author_name, author_role, tag)
				VALUES ($1, $2, $3, $4, 'OG', $5)
			`, req.Title, req.Subtitle, req.Description, u.DisplayName, req.Tag)
		case "events":
			st := time.Now().Add(72 * time.Hour)
			if req.StartDate != "" {
				if parsed, pErr := time.Parse(time.RFC3339, req.StartDate); pErr == nil {
					st = parsed
				} else if parsed, pErr := time.Parse("2006-01-02T15:04", req.StartDate); pErr == nil {
					st = parsed
				}
			}
			loc := req.Subtitle
			if loc == "" {
				loc = "TriGameDev Discord"
			}
			tag := req.Tag
			if tag == "" {
				tag = "Meetup"
			}
			_, err = s.db.Exec(`
				INSERT INTO trigamedev.events (title, description, location, start_time, cover_image_url, tag)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, req.Title, req.Description, loc, st, req.ImageURL, tag)
		case "jams":
			st := time.Now()
			et := time.Now().Add(48 * time.Hour)
			if parsedSt, pErr := time.Parse(time.RFC3339, req.StartDate); pErr == nil {
				st = parsedSt
			}
			if parsedEt, pErr := time.Parse(time.RFC3339, req.EndDate); pErr == nil {
				et = parsedEt
			}
			status := req.Status
			if status == "" {
				status = "active"
			}
			jamType := req.JamType
			if jamType == "" {
				jamType = "monthly"
			}
			winnersBytes, _ := json.Marshal(req.Winners)
			_, err = s.db.Exec(`
				INSERT INTO trigamedev.game_jams (title, month_label, theme, description, itch_url, cover_image_url, start_date, end_date, status, participants_count, submissions_count, jam_type, winners_json)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			`, req.Title, req.Tag, req.Subtitle, req.Description, req.URL, req.ImageURL, st, et, status, req.ParticipantsCount, req.SubmissionsCount, jamType, string(winnersBytes))
		case "showcase":
			srcURL := req.SecondaryURL
			if srcURL == "" {
				srcURL = req.URL
			}
			_, err = s.db.Exec(`
				INSERT INTO trigamedev.showcase (title, developer_name, description, thumbnail_url, play_url, source_url, engine, featured)
				VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE)
			`, req.Title, req.Subtitle, req.Description, req.ImageURL, req.URL, srcURL, req.Tag)
		case "gallery":
			evLabel := req.Tag
			if evLabel == "" {
				evLabel = "Community Photo"
			}
			_, err = s.db.Exec(`
				INSERT INTO trigamedev.gallery (title, caption, image_url, event_label, uploaded_by)
				VALUES ($1, $2, $3, $4, $5)
			`, req.Title, req.Description, req.ImageURL, evLabel, u.DisplayName)
		case "videos":
			speaker := req.Subtitle
			if speaker == "" {
				speaker = u.DisplayName
			}
			_, err = s.db.Exec(`
				INSERT INTO trigamedev.videos (title, description, video_url, thumbnail_url, speaker_name, duration)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, req.Title, req.Description, req.URL, req.ImageURL, speaker, req.Tag)
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unknown post type"})
			return
		}

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		s.notifyDiscordWebhook(fmt.Sprintf("[CMS] **%s** published a new **%s** entry on TriGameDev.org: **%s**", u.DisplayName, req.Section, req.Title))
		writeJSON(w, http.StatusCreated, map[string]string{"status": "published"})
		return
	}

	if r.Method == http.MethodDelete {
		section := r.URL.Query().Get("section")
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		tableMap := map[string]string{
			"news":     "trigamedev.news",
			"events":   "trigamedev.events",
			"jams":     "trigamedev.game_jams",
			"showcase": "trigamedev.showcase",
			"gallery":  "trigamedev.gallery",
			"videos":   "trigamedev.videos",
		}
		tbl, ok := tableMap[section]
		if !ok || id <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid section or id"})
			return
		}
		_, _ = s.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = $1", tbl), id)
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// Discord Scheduled Events Sync
func (s *Server) handleDiscordSync(w http.ResponseWriter, r *http.Request) {
	u := s.currentUser(r)
	if u == nil || !u.IsOGAdmin {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "OG Role required"})
		return
	}
	syncedEvents := s.syncDiscordEventsFromGuild(true)

	writeJSON(w, http.StatusOK, map[string]any{
		"status":        "synced",
		"synced_events": syncedEvents,
		"timestamp":     time.Now().Format(time.RFC3339),
	})
}

func (s *Server) notifyDiscordWebhook(content string) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		return
	}
	go func() {
		body, _ := json.Marshal(map[string]string{"content": content})
		_, _ = http.Post(webhookURL, "application/json", bytes.NewReader(body))
	}()
}
