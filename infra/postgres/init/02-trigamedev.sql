-- TriGameDev Community Platform Schema (Clean Production Init — No Test Seed Data)
CREATE SCHEMA IF NOT EXISTS trigamedev;

-- 1. Users & Discord Roles
CREATE TABLE IF NOT EXISTS trigamedev.users (
    id SERIAL PRIMARY KEY,
    discord_id VARCHAR(64) UNIQUE NOT NULL,
    username VARCHAR(100) NOT NULL,
    display_name VARCHAR(120) NOT NULL,
    avatar_url TEXT DEFAULT '',
    is_guild_member BOOLEAN DEFAULT FALSE,
    is_og_admin BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Events (Synced from Discord Scheduled Events + OG CMS)
CREATE TABLE IF NOT EXISTS trigamedev.events (
    id SERIAL PRIMARY KEY,
    discord_event_id VARCHAR(64) UNIQUE,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(200) DEFAULT 'TriGameDev Discord Stage',
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    cover_image_url TEXT DEFAULT '',
    event_url TEXT DEFAULT '',
    tag VARCHAR(50) DEFAULT 'Meetup',
    rsvp_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. News & Announcements
CREATE TABLE IF NOT EXISTS trigamedev.news (
    id SERIAL PRIMARY KEY,
    discord_message_id VARCHAR(64) UNIQUE,
    title VARCHAR(200) NOT NULL,
    summary TEXT NOT NULL,
    body_markdown TEXT NOT NULL,
    author_name VARCHAR(100) NOT NULL,
    author_role VARCHAR(60) DEFAULT 'OG Admin',
    tag VARCHAR(50) DEFAULT 'Announcement',
    is_pinned BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Monthly Game Jams
CREATE TABLE IF NOT EXISTS trigamedev.game_jams (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    month_label VARCHAR(60) NOT NULL,
    theme VARCHAR(160) NOT NULL,
    description TEXT NOT NULL,
    itch_url TEXT DEFAULT 'https://itch.io/jams',
    cover_image_url TEXT DEFAULT '',
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(30) DEFAULT 'active',
    participants_count INT DEFAULT 0,
    submissions_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Member Game Showcase
CREATE TABLE IF NOT EXISTS trigamedev.showcase (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    developer_name VARCHAR(120) NOT NULL,
    description TEXT NOT NULL,
    thumbnail_url TEXT DEFAULT '',
    play_url TEXT DEFAULT '',
    source_url TEXT DEFAULT '',
    engine VARCHAR(60) DEFAULT 'Godot',
    jam_id INT REFERENCES trigamedev.game_jams(id) ON DELETE SET NULL,
    views_count INT DEFAULT 0,
    featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 6. Photo Gallery
CREATE TABLE IF NOT EXISTS trigamedev.gallery (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    caption TEXT DEFAULT '',
    image_url TEXT NOT NULL,
    event_label VARCHAR(100) DEFAULT 'Community Meetup',
    uploaded_by VARCHAR(100) DEFAULT 'TriGameDev',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. Video Archive
CREATE TABLE IF NOT EXISTS trigamedev.videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT DEFAULT '',
    video_url TEXT NOT NULL,
    thumbnail_url TEXT DEFAULT '',
    speaker_name VARCHAR(120) DEFAULT 'Community Speaker',
    duration VARCHAR(30) DEFAULT '',
    recorded_at TIMESTAMPTZ DEFAULT NOW()
);

-- 8. Key Sharing Portal: Game Pools
CREATE TABLE IF NOT EXISTS trigamedev.key_pools (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    platform VARCHAR(50) NOT NULL DEFAULT 'Steam',
    store_url TEXT DEFAULT '',
    cover_image_url TEXT DEFAULT '',
    description TEXT DEFAULT '',
    donated_by_name VARCHAR(100) DEFAULT 'Member',
    created_by INT REFERENCES trigamedev.users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 9. Key Sharing Portal: Individual Redeemable Keys
-- Enforces Single Redemption Per User Per Pool via UNIQUE(pool_id, claimed_by)
CREATE TABLE IF NOT EXISTS trigamedev.keys (
    id SERIAL PRIMARY KEY,
    pool_id INT NOT NULL REFERENCES trigamedev.key_pools(id) ON DELETE CASCADE,
    key_code VARCHAR(255) NOT NULL,
    donated_by INT REFERENCES trigamedev.users(id) ON DELETE SET NULL,
    claimed_by INT REFERENCES trigamedev.users(id) ON DELETE SET NULL,
    claimed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_user_claim_per_pool UNIQUE (pool_id, claimed_by)
);

-- 10. Ebook Library (Direct Downloads for Members)
CREATE TABLE IF NOT EXISTS trigamedev.ebooks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    author VARCHAR(150) NOT NULL,
    category VARCHAR(80) NOT NULL DEFAULT 'Game Design',
    description TEXT DEFAULT '',
    cover_image_url TEXT DEFAULT '',
    file_name VARCHAR(255) NOT NULL,
    file_format VARCHAR(20) NOT NULL DEFAULT 'PDF',
    file_size_label VARCHAR(40) DEFAULT '',
    uploaded_by_name VARCHAR(100) DEFAULT 'Member',
    download_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 11. Schema Migrations for Roles, Jam Winners, and Per-User Unread Notifications
ALTER TABLE trigamedev.users ADD COLUMN IF NOT EXISTS roles_json TEXT DEFAULT '[]';
ALTER TABLE trigamedev.game_jams ADD COLUMN IF NOT EXISTS jam_type VARCHAR(40) DEFAULT 'monthly';
ALTER TABLE trigamedev.game_jams ADD COLUMN IF NOT EXISTS winners_json TEXT DEFAULT '[]';

CREATE TABLE IF NOT EXISTS trigamedev.user_section_reads (
    user_id INT NOT NULL REFERENCES trigamedev.users(id) ON DELETE CASCADE,
    section VARCHAR(40) NOT NULL,
    last_read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, section)
);

