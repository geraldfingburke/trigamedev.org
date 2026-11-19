import json
import os
import random
from contextlib import asynccontextmanager
from pathlib import Path
from typing import Any

from fastapi import FastAPI, Request, Response
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse, PlainTextResponse
from psycopg_pool import ConnectionPool

DATABASE_URL = os.getenv("DATABASE_URL", "")
pool: ConnectionPool | None = None

TAROT_JSON_PATH = Path(__file__).parent / "tarot_card_data.json"
IMAGE_BASE_URL = "https://geraldburke.com/apis/plateautarot/images"


def init_schemas() -> None:
    if not pool:
        return
    with pool.connection() as conn:
        with conn.cursor() as cur:
            # 1. Plateau Tarot PostgreSQL Schema & Table
            cur.execute("CREATE SCHEMA IF NOT EXISTS plateau_tarot;")
            cur.execute(
                """
                CREATE TABLE IF NOT EXISTS plateau_tarot.cards (
                    id INT PRIMARY KEY,
                    name VARCHAR(120) NOT NULL,
                    arcana VARCHAR(40) NOT NULL,
                    suit VARCHAR(40) NOT NULL,
                    description TEXT NOT NULL,
                    rank VARCHAR(40) NOT NULL,
                    symbols TEXT NOT NULL,
                    image_url TEXT NOT NULL,
                    card_back TEXT NOT NULL
                );
                """
            )
            cur.execute("SELECT COUNT(*) FROM plateau_tarot.cards;")
            row = cur.fetchone()
            card_count = row[0] if row else 0
            if card_count == 0 and TAROT_JSON_PATH.exists():
                raw = json.loads(TAROT_JSON_PATH.read_text(encoding="utf-8"))
                for idx, card in enumerate(raw.get("cards", []), start=1):
                    name = card["name"]
                    cur.execute(
                        """
                        INSERT INTO plateau_tarot.cards
                            (id, name, arcana, suit, description, rank, symbols, image_url, card_back)
                        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
                        ON CONFLICT (id) DO NOTHING;
                        """,
                        (
                            idx,
                            name,
                            card["arcana"],
                            card["suit"],
                            card["description"],
                            str(card["rank"]),
                            card["symbols"],
                            f"{IMAGE_BASE_URL}/{name}.jpg",
                            f"{IMAGE_BASE_URL}/Card Back.jpg",
                        ),
                    )

            # 2. Simple Leaderboard PostgreSQL Schema & Tables
            cur.execute("CREATE SCHEMA IF NOT EXISTS simple_leaderboard;")
            cur.execute(
                """
                CREATE TABLE IF NOT EXISTS simple_leaderboard.users (
                    game_id SERIAL PRIMARY KEY,
                    email VARCHAR(255) NOT NULL,
                    game_name VARCHAR(255) NOT NULL,
                    created_at TIMESTAMPTZ DEFAULT NOW()
                );
                """
            )
            cur.execute(
                """
                CREATE TABLE IF NOT EXISTS simple_leaderboard.scores (
                    id SERIAL PRIMARY KEY,
                    game_id INT NOT NULL REFERENCES simple_leaderboard.users(game_id) ON DELETE CASCADE,
                    score INT NOT NULL,
                    user_name VARCHAR(255) NOT NULL,
                    created_at TIMESTAMPTZ DEFAULT NOW()
                );
                """
            )
            cur.execute(
                "CREATE INDEX IF NOT EXISTS idx_leaderboard_scores_game_score ON simple_leaderboard.scores (game_id, score DESC);"
            )
            # Ensure at least game_id = 1 exists for legacy/demo calls
            cur.execute("SELECT COUNT(*) FROM simple_leaderboard.users;")
            u_row = cur.fetchone()
            if (u_row[0] if u_row else 0) == 0:
                cur.execute(
                    """
                    INSERT INTO simple_leaderboard.users (email, game_name)
                    VALUES ('gerald@geraldburke.com', 'Default Leaderboard');
                    """
                )
        conn.commit()


@asynccontextmanager
async def lifespan(app: FastAPI):
    global pool
    pool = ConnectionPool(conninfo=DATABASE_URL, min_size=1, max_size=10, open=True)
    try:
        init_schemas()
    except Exception as exc:
        print(f"Schema init notice: {exc}")
    yield
    if pool:
        pool.close()


app = FastAPI(title="Shared FastAPI Backend (Plateau Tarot & Simple Leaderboard)", lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
def health():
    db_status = "connected"
    db_time = None
    services = []
    try:
        if pool:
            with pool.connection() as conn:
                with conn.cursor() as cur:
                    cur.execute("SELECT NOW()::text")
                    row = cur.fetchone()
                    if row:
                        db_time = row[0]
                    cur.execute("SELECT service_name, runtime FROM shared.service_registry ORDER BY id")
                    services = [{"name": r[0], "runtime": r[1]} for r in cur.fetchall()]
    except Exception as e:
        db_status = f"error: {e}"

    return {
        "service": "api-fastapi",
        "status": "ok",
        "database": db_status,
        "db_time": db_time,
        "registered_services": services,
    }


@app.get("/")
def root():
    return {
        "service": "api-fastapi",
        "message": "Shared FastAPI Backend is live",
        "apis": ["/apis/plateautarot/", "/apis/simple-leaderboard-api/"],
    }


# ==============================================================================
# 1. PLATEAU TAROT API (FastAPI + PostgreSQL)
# Preserves 100% compatibility with https://github.com/geraldfingburke/plateau-tarot-api
# ==============================================================================
def _load_tarot_deck() -> list[dict[str, Any]]:
    if pool:
        with pool.connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT id, name, arcana, suit, description, rank, symbols, image_url, card_back
                    FROM plateau_tarot.cards
                    ORDER BY id ASC;
                    """
                )
                rows = cur.fetchall()
                if rows:
                    return [
                        {
                            "name": r[1],
                            "arcana": r[2],
                            "suit": r[3],
                            "description": r[4],
                            "rank": r[5],
                            "symbols": r[6],
                            "imageURL": r[7],
                            "cardBack": r[8],
                            "id": r[0],
                        }
                        for r in rows
                    ]
    # Fallback to JSON file if DB is unavailable
    raw = json.loads(TAROT_JSON_PATH.read_text(encoding="utf-8"))
    deck = []
    for idx, c in enumerate(raw.get("cards", []), start=1):
        deck.append(
            {
                "name": c["name"],
                "arcana": c["arcana"],
                "suit": c["suit"],
                "description": c["description"],
                "rank": str(c["rank"]),
                "symbols": c["symbols"],
                "imageURL": f"{IMAGE_BASE_URL}/{c['name']}.jpg",
                "cardBack": f"{IMAGE_BASE_URL}/Card Back.jpg",
                "id": idx,
            }
        )
    return deck


@app.get("/apis/plateautarot")
@app.get("/apis/plateautarot/")
@app.get("/plateautarot/api")
@app.get("/plateautarot/api/")
def plateau_tarot_handler(
    seed: str | None = None,
    id: str | None = None,
    arcana: str | None = None,
    suit: str | None = None,
    count: str | None = None,
) -> Response:
    cards = _load_tarot_deck()

    if seed is not None and seed != "":
        try:
            seed_val: int | str = int(seed)
        except ValueError:
            seed_val = seed
    else:
        seed_val = random.randint(1, 2147483647)

    # Apply Fisher-Yates shuffle to deck using the seed
    rng = random.Random(seed_val)
    for i in range(len(cards)):
        r = rng.randint(0, i)
        cards[i], cards[r] = cards[r], cards[i]

    # 1. Filter by id
    if id is not None:
        try:
            target_id = int(id)
        except ValueError:
            return PlainTextResponse("Card not found")
        for card in cards:
            if card["id"] == target_id:
                return JSONResponse(card)
        return PlainTextResponse("Card not found")

    # 2. Filter by arcana
    if arcana is not None:
        target_arcana = arcana.strip().lower()
        selection = [c for c in cards if c["arcana"].lower() == target_arcana]
        if len(selection) < 1:
            return PlainTextResponse("Invalid Arcana")
        cards = selection

    # 3. Filter by suit
    if suit is not None:
        target_suit = suit.strip().lower()
        selection = [c for c in cards if c["suit"].lower() == target_suit]
        if len(selection) < 1:
            return PlainTextResponse("Invalid Suit")
        cards = selection

    # 4. Limit by count
    if count is not None:
        try:
            n = int(count)
        except ValueError:
            return PlainTextResponse("Not a number!")
        if n > 78 or n < 1 or n > len(cards):
            if n > 78 or n < 1:
                return PlainTextResponse("Invalid Number")
            n = len(cards)
        cards = cards[:n]

    return JSONResponse({"cards": cards, "seed": seed_val})


# ==============================================================================
# 2. SIMPLE LEADERBOARD API (FastAPI + PostgreSQL)
# Preserves 100% compatibility with https://github.com/geraldfingburke/simple-leaderboard-api
# ==============================================================================
@app.api_route("/apis/simple-leaderboard-api", methods=["GET", "POST"])
@app.api_route("/apis/simple-leaderboard-api/", methods=["GET", "POST"])
@app.api_route("/apis/simple-leaderboard", methods=["GET", "POST"])
@app.api_route("/apis/simple-leaderboard/", methods=["GET", "POST"])
async def simple_leaderboard_handler(request: Request) -> Response:
    if not pool:
        return PlainTextResponse("Database connection unavailable", status_code=500)

    qp = request.query_params
    action = qp.get("action")
    game_id_raw = qp.get("gameID")
    score_raw = qp.get("score")
    user_name = qp.get("userName")
    count_raw = qp.get("count")

    # 1. action == "newScore"
    if action == "newScore":
        if game_id_raw and score_raw and user_name:
            try:
                game_id = int(game_id_raw)
                score_val = int(score_raw)
            except ValueError:
                return PlainTextResponse("Invalid Params")

            with pool.connection() as conn:
                with conn.cursor() as cur:
                    cur.execute(
                        "SELECT game_id FROM simple_leaderboard.users WHERE game_id = %s;",
                        (game_id,),
                    )
                    if not cur.fetchone():
                        return PlainTextResponse("Invalid GameID")

                    cur.execute(
                        """
                        INSERT INTO simple_leaderboard.scores (game_id, score, user_name)
                        VALUES (%s, %s, %s);
                        """,
                        (game_id, score_val, user_name),
                    )
                conn.commit()
            return PlainTextResponse("Success")
        return PlainTextResponse("Invalid Params")

    # 2. action == "topScores"
    if action == "topScores":
        if not game_id_raw:
            return PlainTextResponse("Invalid params")
        try:
            game_id = int(game_id_raw)
            limit = int(count_raw) if count_raw else 10
        except ValueError:
            return PlainTextResponse("Invalid params")

        with pool.connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT user_name, score
                    FROM simple_leaderboard.scores
                    WHERE game_id = %s
                    ORDER BY score DESC
                    LIMIT %s;
                    """,
                    (game_id, limit),
                )
                rows = cur.fetchall()
        return JSONResponse([{"userName": r[0], "score": r[1]} for r in rows])

    # 3. action == "userScores"
    if action == "userScores":
        if not user_name or not game_id_raw:
            return PlainTextResponse("Invalid params")
        try:
            game_id = int(game_id_raw)
            limit = int(count_raw) if count_raw else 10
        except ValueError:
            return PlainTextResponse("Invalid params")

        with pool.connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT user_name, score
                    FROM simple_leaderboard.scores
                    WHERE game_id = %s AND user_name = %s
                    ORDER BY score DESC
                    LIMIT %s;
                    """,
                    (game_id, user_name, limit),
                )
                rows = cur.fetchall()
        return JSONResponse([{r[0]: r[1]} for r in rows])

    # 4. POST JSON body {"action": "addUser", "email": "...", "gameName": "..."}
    try:
        data = await request.json()
    except Exception:
        data = {}

    if isinstance(data, dict) and data.get("action") == "addUser":
        email = (data.get("email") or "").strip()
        game_name = (data.get("gameName") or "").strip()
        if not email or not game_name:
            return PlainTextResponse("Invalid params")

        with pool.connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT game_id
                    FROM simple_leaderboard.users
                    WHERE email = %s AND game_name = %s;
                    """,
                    (email, game_name),
                )
                existing = cur.fetchone()
                if existing:
                    return JSONResponse(existing[0])

                cur.execute(
                    """
                    INSERT INTO simple_leaderboard.users (email, game_name)
                    VALUES (%s, %s)
                    RETURNING game_id;
                    """,
                    (email, game_name),
                )
                new_row = cur.fetchone()
            conn.commit()
        return JSONResponse(new_row[0] if new_row else 1)

    return PlainTextResponse("Invalid params")
