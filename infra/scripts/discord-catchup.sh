#!/usr/bin/env bash
# ==============================================================================
# TriGameDev Discord Bot Catch-Up Drip Script
# Posts 1 unannounced existing entry across all categories every 6 hours
# until the database is 100% caught up.
# Usage:
#   ./discord-catchup.sh          # Respects 6-hour cooldown
#   ./discord-catchup.sh --force  # Immediately posts 1 unannounced entry
#   ./discord-catchup.sh --status # Prints remaining unannounced counts
# ==============================================================================
set -euo pipefail

MODE="${1:-normal}"

if [ "$MODE" = "--status" ]; then
  docker exec -i infra-postgres psql -U postgres -d main -c "
    SELECT section, COUNT(*) AS remaining_unannounced FROM (
      SELECT 'jams' AS section FROM trigamedev.game_jams WHERE COALESCE(discord_announced, FALSE) = FALSE
      UNION ALL
      SELECT 'events' FROM trigamedev.events WHERE COALESCE(discord_announced, FALSE) = FALSE AND start_time >= NOW()
      UNION ALL
      SELECT 'news' FROM trigamedev.news WHERE COALESCE(discord_announced, FALSE) = FALSE
      UNION ALL
      SELECT 'showcase' FROM trigamedev.showcase WHERE COALESCE(discord_announced, FALSE) = FALSE
      UNION ALL
      SELECT 'gallery' FROM trigamedev.gallery WHERE COALESCE(discord_announced, FALSE) = FALSE
      UNION ALL
      SELECT 'videos' FROM trigamedev.videos WHERE COALESCE(discord_announced, FALSE) = FALSE
      UNION ALL
      SELECT 'keys' FROM trigamedev.key_pools WHERE COALESCE(discord_announced, FALSE) = FALSE
      UNION ALL
      SELECT 'ebooks' FROM trigamedev.ebooks WHERE COALESCE(discord_announced, FALSE) = FALSE
    ) q GROUP BY section ORDER BY section;
  "
  exit 0
fi

FORCE_FLAG="0"
if [ "$MODE" = "--force" ]; then
  FORCE_FLAG="1"
fi

curl -s -X POST "http://127.0.0.1:8080/tgd/internal/discord-catchup?force=${FORCE_FLAG}"
echo ""
