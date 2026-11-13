-- TriGameDev Community Platform Schema & Seed Data
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
    rsvp_count INT DEFAULT 12,
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
    status VARCHAR(30) DEFAULT 'active', -- upcoming, active, voting, completed
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
    views_count INT DEFAULT 100,
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
    duration VARCHAR(30) DEFAULT '42:15',
    recorded_at TIMESTAMPTZ DEFAULT NOW()
);

-- 8. Key Sharing Portal: Game Pools
CREATE TABLE IF NOT EXISTS trigamedev.key_pools (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    platform VARCHAR(50) NOT NULL DEFAULT 'Steam', -- Steam, Itch.io, GOG, Epic
    store_url TEXT DEFAULT '',
    cover_image_url TEXT DEFAULT '',
    description TEXT DEFAULT '',
    donated_by_name VARCHAR(100) DEFAULT 'Anonymous Member',
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
    file_size_label VARCHAR(40) DEFAULT '4.2 MB',
    uploaded_by_name VARCHAR(100) DEFAULT 'OG Archive',
    download_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================================
-- INITIAL SEED DATA (Only inserted if tables are empty)
-- ============================================================================

INSERT INTO trigamedev.users (id, discord_id, username, display_name, avatar_url, is_guild_member, is_og_admin)
VALUES
    (1, '100000000000000001', 'geraldfingburke', 'Gerald Burke (OG)', 'https://api.dicebear.com/9.x/bottts-neutral/svg?seed=Gerald', TRUE, TRUE),
    (2, '100000000000000002', 'synthwave_dev', 'Mira Vance (Member)', 'https://api.dicebear.com/9.x/avataaars/svg?seed=Mira', TRUE, FALSE)
ON CONFLICT (discord_id) DO NOTHING;

SELECT setval('trigamedev.users_id_seq', (SELECT COALESCE(MAX(id), 1) FROM trigamedev.users));

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM trigamedev.game_jams LIMIT 1) THEN
        INSERT INTO trigamedev.game_jams (title, month_label, theme, description, itch_url, cover_image_url, start_date, end_date, status, participants_count, submissions_count)
        VALUES
        ('TriGameDev Monthly Jam #42', 'November 2026', 'SIGNAL LOST // LOFI HORIZONS', 'Build a game in 72 hours exploring atmospheric frequencies, coastal synth setups, or broken transmissions. Solo or teams up to 4!', 'https://itch.io/jams', '/inspo/sunset-studio.jpg', NOW() - INTERVAL '2 days', NOW() + INTERVAL '5 days', 'active', 48, 14),
        ('TriGameDev Monthly Jam #41', 'October 2026', 'CLOCKWORK PHANTOMS', 'Our spooky October jam focused on mechanical puzzles and haunted automata.', 'https://itch.io/jams', '/inspo/rooftop-twilight.jpg', NOW() - INTERVAL '35 days', NOW() - INTERVAL '28 days', 'completed', 62, 23);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM trigamedev.showcase LIMIT 1) THEN
        INSERT INTO trigamedev.showcase (title, developer_name, description, thumbnail_url, play_url, source_url, engine, views_count, featured)
        VALUES
        ('Frequency 88.4: Coastal Drift', 'gerald_b & synth_mira', 'Tune your analog receiver across a dying archipelago while composing chill synthwave beats to guide lost ships home.', '/inspo/sunset-studio.jpg', 'https://itch.io', 'https://github.com/geraldfingburke', 'Godot 4.3', 1420, TRUE),
        ('Rooftop Telemetry', 'pixel_crow', 'A narrative parkour adventure across twilight antennas and suburban powerlines. Built for the TriGameDev Monthly Jam.', '/inspo/rooftop-twilight.jpg', 'https://itch.io', 'https://github.com/geraldfingburke', 'Godot 4.3', 980, TRUE),
        ('Afterhours Cyber-Lounge', 'jinsang_fan', 'Isometric point-and-click detective story set in a cluttered neon apartment full of strange artifacts and cassette tapes.', '/inspo/lofi-room.jpg', 'https://itch.io', 'https://github.com/geraldfingburke', 'Custom / Go + WebGL', 2150, TRUE),
        ('Chromancer Micro', 'gerald_b', 'Fast-paced color-alchemical arena survival crafted during MicroJam 040.', '/inspo/sunset-studio.jpg', 'https://itch.io/jam/micro-jam-040/rate/3609527', 'https://github.com/geraldfingburke/cromancer', 'Godot', 640, FALSE);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM trigamedev.events LIMIT 1) THEN
        INSERT INTO trigamedev.events (discord_event_id, title, description, location, start_time, end_time, cover_image_url, tag, rsvp_count)
        VALUES
        ('evt_101', 'Playtest Night & Pizza: Monthly Jam #42 Builds', 'Bring your WIP or finished jam builds! We screen-share in Discord Stage & meet in-person at the downtown co-op lab.', 'Discord Stage + Tri-Cities Co-Op', NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days 3 hours', '/inspo/sunset-studio.jpg', 'Playtest', 34),
        ('evt_102', 'Shader & Lofi Audio Workshop in Godot 4', 'Live deep-dive into CRT post-processing shaders, reactive audio buses, and procedural skyboxes.', 'Discord Voice #workshop-stage', NOW() + INTERVAL '6 days', NOW() + INTERVAL '6 days 2 hours', '/inspo/rooftop-twilight.jpg', 'Workshop', 29),
        ('evt_103', 'Coffee & Code Saturday Morning Hangout', 'Casual co-working session. Bring whatever bug or sprite sheet you are wrestling with this week.', 'Discord Lounge + Local Roasters', NOW() + INTERVAL '9 days', NOW() + INTERVAL '9 days 3 hours', '/inspo/lofi-room.jpg', 'Hangout', 19);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM trigamedev.news LIMIT 1) THEN
        INSERT INTO trigamedev.news (discord_message_id, title, summary, body_markdown, author_name, author_role, tag, is_pinned)
        VALUES
        ('msg_201', 'Welcome to the New TriGameDev Community Hub & VPS!', 'We have officially migrated to our dedicated VPS with Discord OAuth, a Member Key Vault, and Direct-Download Ebook Library.', 'Welcome home, devs! Our new platform is tightly linked with the **TriGameDev Discord**. Log in with your Discord account to claim donated Steam/Itch keys, download game dev ebooks, or—if you carry the **OG Role**—publish directly from the built-in CMS.', 'Gerald Burke', 'OG Admin', 'Infrastructure', TRUE),
        ('msg_202', 'Monthly Game Jam #42 Theme Announced: SIGNAL LOST', 'You have 72 hours to build something unforgettable around analog signals, twilight horizons, or lost transmissions.', 'Submissions are open until Sunday at 11:59 PM EST. Remember to post your progress gifs in `#jam-wip` on Discord—our webhook bot automatically highlights new submissions right here!', 'Mira Vance', 'Jam Host', 'Game Jam', FALSE),
        ('msg_203', '15 New Steam Keys & 4 Shader Books Added to the Vault', 'Community members dropped a fresh batch of indie Steam keys and Godot 4 reference guides.', 'Head over to the **Key Portal** and **Ebook Library** tabs while logged in. Remember: each key pool allows **one redemption per member** so everyone gets a fair shot!', 'TriGameDev Bot', 'Discord Sync', 'Community', FALSE);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM trigamedev.gallery LIMIT 1) THEN
        INSERT INTO trigamedev.gallery (title, caption, image_url, event_label, uploaded_by)
        VALUES
        ('Sunset Synth Audio Rig', 'Testing custom MIDI controllers for the November Audio Jam.', '/inspo/sunset-studio.jpg', 'Audio Jam Setup', 'gerald_b'),
        ('Twilight Rooftop Concept Art', 'Environment study painted during the Tuesday Art Sprint.', '/inspo/rooftop-twilight.jpg', 'Art Sprint', 'synth_mira'),
        ('Late Night Debug Lounge', 'When the scope creep hits at 2 AM on Sunday before submission.', '/inspo/lofi-room.jpg', 'Monthly Jam #41', 'pixel_crow'),
        ('Shader Graph Experiments', 'Volumetric clouds and retro dithering pass in Godot 4.', '/inspo/sunset-studio.jpg', 'Tech Demo', 'gerald_b'),
        ('Controller Playtest Booth', 'Four-player couch co-op stress test at the monthly meetup.', '/inspo/rooftop-twilight.jpg', 'October Meetup', 'og_alex'),
        ('Cassette UI Mockups', 'Diegetic inventory system for Afterhours Cyber-Lounge.', '/inspo/lofi-room.jpg', 'UI Showcase', 'jinsang_fan'),
        ('Retro CRT Filter Pass', 'Scanline and bloom tuning on ultrawide monitors.', '/inspo/sunset-studio.jpg', 'Shader Workshop', 'synth_mira');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM trigamedev.videos LIMIT 1) THEN
        INSERT INTO trigamedev.videos (title, description, video_url, thumbnail_url, speaker_name, duration)
        VALUES
        ('Building Custom GDScript Tools & MCP Workflows', 'How to supercharge your solo dev pipeline in Godot 4 with custom editor plugins.', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', '/inspo/sunset-studio.jpg', 'Gerald Burke', '38:40'),
        ('Composing Lofi & Synthwave Loops for Indie Games', 'From 4-bar chord progressions in Ableton to seamless adaptive loops in FMOD.', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', '/inspo/rooftop-twilight.jpg', 'Mira Vance', '51:10'),
        ('Juice & Game Feel in 30 Minutes', 'Screen shake, hit-stop, tweening curves, and particle bursts that make combat pop.', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', '/inspo/lofi-room.jpg', 'Pixel Crow', '29:55');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM trigamedev.key_pools LIMIT 1) THEN
        INSERT INTO trigamedev.key_pools (id, title, platform, store_url, cover_image_url, description, donated_by_name, created_by)
        VALUES
        (1, 'Aseprite (Pixel Art & Animated Sprites)', 'Steam', 'https://store.steampowered.com', '/inspo/sunset-studio.jpg', 'The gold-standard animated sprite editor & pixel art tool for indie game devs.', 'Gerald Burke', 1),
        (2, 'Hyper Light Drifter: Collector Edition', 'Steam', 'https://store.steampowered.com', '/inspo/rooftop-twilight.jpg', 'Extra bundle keys donated for TriGameDev members studying top-down dash combat and pixel atmosphere.', 'Mira Vance', 2),
        (3, 'Retro Synth & SFX Complete Asset Pack', 'Itch.io', 'https://itch.io', '/inspo/lofi-room.jpg', 'Over 400 royalty-free analog synth loops, UI blips, and ambient rain textures.', 'Gerald Burke', 1),
        (4, 'Celeste (Indie Platformer Study)', 'Steam', 'https://store.steampowered.com', '/inspo/rooftop-twilight.jpg', 'Donated spare keys from the charity bundle—essential study material for platformer feel!', 'Pixel Crow', 1);

        PERFORM setval('trigamedev.key_pools_id_seq', (SELECT COALESCE(MAX(id), 1) FROM trigamedev.key_pools));

        INSERT INTO trigamedev.keys (pool_id, key_code, donated_by, claimed_by, claimed_at)
        VALUES
        (1, 'ASEP-9942-LOFI-7731-TGD1', 1, NULL, NULL),
        (1, 'ASEP-8821-SYNTH-4410-TGD2', 1, NULL, NULL),
        (1, 'ASEP-3319-PIXEL-9012-TGD3', 1, 2, NOW() - INTERVAL '1 day'),
        (2, 'HLDX-7741-NEON-2209-DRFT', 2, NULL, NULL),
        (2, 'HLDX-5510-CYBR-8832-SWD2', 2, NULL, NULL),
        (3, 'ITCH-SYNTH-PACK-2026-ALPHA', 1, NULL, NULL),
        (3, 'ITCH-SYNTH-PACK-2026-BETA', 1, NULL, NULL),
        (3, 'ITCH-SYNTH-PACK-2026-GAMMA', 1, NULL, NULL),
        (4, 'CLST-4490-MTN1-8823-DASH', 1, 1, NOW() - INTERVAL '3 days'),
        (4, 'CLST-1102-STRW-9941-JUMP', 1, NULL, NULL);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM trigamedev.ebooks LIMIT 1) THEN
        INSERT INTO trigamedev.ebooks (title, author, category, description, cover_image_url, file_name, file_format, file_size_label, uploaded_by_name, download_count)
        VALUES
        ('Making Games as a Hobby: Complete Companion Edition', 'Gerald Burke', 'Game Design', 'Practical, no-nonsense guide to scoping, finishing, and shipping indie games without burning out.', '/inspo/sunset-studio.jpg', 'making-games-as-a-hobby.pdf', 'PDF', '6.4 MB', 'Gerald Burke (OG)', 84),
        ('The Book of Shaders & Lofi Post-Processing', 'TriGameDev Community Notes', 'Shaders & VFX', 'Annotated GLSL & Godot 4 spatial/canvas_item shader recipes: CRT curvature, dithering, bloom, and water.', '/inspo/rooftop-twilight.jpg', 'trigamedev-shader-grimoire.pdf', 'PDF', '11.8 MB', 'Mira Vance', 67),
        ('Game Programming Patterns (Indie Field Reference)', 'Community Study Edition', 'Architecture', 'State machines, component patterns, command queues, and spatial partitions tailored for Godot & Go.', '/inspo/lofi-room.jpg', 'game-programming-patterns-notes.epub', 'EPUB', '3.9 MB', 'Gerald Burke (OG)', 52),
        ('Procedural Audio & Synth Design for Devs', 'TriGameDev Audio Guild', 'Audio & Music', 'How to synthesize footsteps, lasers, rain, and lofi chords procedurally at runtime.', '/inspo/sunset-studio.jpg', 'procedural-audio-handbook.pdf', 'PDF', '8.1 MB', 'Mira Vance', 41);
    END IF;
END $$;
