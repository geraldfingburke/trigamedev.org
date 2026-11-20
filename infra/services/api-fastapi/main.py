import json
import os
import random
from contextlib import asynccontextmanager
from pathlib import Path
from typing import Any

from fastapi import FastAPI, Request, Response
from fastapi.concurrency import run_in_threadpool
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse, PlainTextResponse
from psycopg_pool import ConnectionPool

DATABASE_URL = os.getenv("DATABASE_URL", "")
pool: ConnectionPool | None = None

TAROT_JSON_PATH = Path(__file__).parent / "tarot_card_data.json"
IMAGE_BASE_URL = "https://geraldburke.com/apis/plateautarot/images"

PG_INT32_MIN = -2_147_483_648
PG_INT32_MAX = 2_147_483_647
VARCHAR_MAX_LENGTH = 255


def initialize_database_schemas() -> None:
    if not pool:
        return
    with pool.connection() as conn:
        with conn.cursor() as cur:
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

            cur.execute("CREATE SCHEMA IF NOT EXISTS simple_leaderboard;")
            cur.execute(
                """
                CREATE TABLE IF NOT EXISTS simple_leaderboard.users (
                    game_id SERIAL PRIMARY KEY,
                    email VARCHAR(255) NOT NULL,
                    game_name VARCHAR(255) NOT NULL,
                    created_at TIMESTAMPTZ DEFAULT NOW(),
                    CONSTRAINT uq_leaderboard_users_email_game_name UNIQUE (email, game_name)
                );
                """
            )
            cur.execute(
                "CREATE UNIQUE INDEX IF NOT EXISTS idx_leaderboard_users_email_game_name ON simple_leaderboard.users (email, game_name);"
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
            cur.execute("SELECT COUNT(*) FROM simple_leaderboard.users;")
            user_row = cur.fetchone()
            if (user_row[0] if user_row else 0) == 0:
                cur.execute(
                    """
                    INSERT INTO simple_leaderboard.users (email, game_name)
                    VALUES ('gerald@geraldburke.com', 'Default Leaderboard')
                    ON CONFLICT (email, game_name) DO NOTHING;
                    """
                )
        conn.commit()


@asynccontextmanager
async def lifespan(app: FastAPI):
    global pool
    pool = ConnectionPool(conninfo=DATABASE_URL, min_size=1, max_size=10, open=True)
    try:
        await run_in_threadpool(initialize_database_schemas)
    except Exception as exc:
        print(f"Schema init notice: {exc}")
    yield
    if pool:
        await run_in_threadpool(pool.close)


app = FastAPI(title="Shared FastAPI Backend", lifespan=lifespan)

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


def load_tarot_deck() -> list[dict[str, Any]]:
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
                            "name": row[1],
                            "arcana": row[2],
                            "suit": row[3],
                            "description": row[4],
                            "rank": row[5],
                            "symbols": row[6],
                            "imageURL": row[7],
                            "cardBack": row[8],
                            "id": row[0],
                        }
                        for row in rows
                    ]
    raw = json.loads(TAROT_JSON_PATH.read_text(encoding="utf-8"))
    return [
        {
            "name": card["name"],
            "arcana": card["arcana"],
            "suit": card["suit"],
            "description": card["description"],
            "rank": str(card["rank"]),
            "symbols": card["symbols"],
            "imageURL": f"{IMAGE_BASE_URL}/{card['name']}.jpg",
            "cardBack": f"{IMAGE_BASE_URL}/Card Back.jpg",
            "id": idx,
        }
        for idx, card in enumerate(raw.get("cards", []), start=1)
    ]


def resolve_seed(seed_param: str | None) -> int | str:
    if seed_param is not None and seed_param != "":
        try:
            return int(seed_param)
        except ValueError:
            return seed_param
    return random.randint(1, 2147483647)


def fisher_yates_shuffle(cards: list[dict[str, Any]], seed_val: int | str) -> list[dict[str, Any]]:
    shuffled = list(cards)
    rng = random.Random(seed_val)
    for index in range(len(shuffled)):
        swap_index = rng.randint(0, index)
        shuffled[index], shuffled[swap_index] = shuffled[swap_index], shuffled[index]
    return shuffled


def find_card_by_id(cards: list[dict[str, Any]], id_param: str) -> dict[str, Any] | None:
    try:
        target_id = int(id_param)
    except ValueError:
        return None
    for card in cards:
        if card["id"] == target_id:
            return card
    return None


def filter_cards_by_arcana(cards: list[dict[str, Any]], arcana_param: str) -> list[dict[str, Any]]:
    target = arcana_param.strip().lower()
    return [card for card in cards if card["arcana"].lower() == target]


def filter_cards_by_suit(cards: list[dict[str, Any]], suit_param: str) -> list[dict[str, Any]]:
    target = suit_param.strip().lower()
    return [card for card in cards if card["suit"].lower() == target]


@app.get("/apis/plateautarot")
@app.get("/apis/plateautarot/")
@app.get("/plateautarot/api")
@app.get("/plateautarot/api/")
def get_tarot_cards(
    seed: str | None = None,
    id: str | None = None,
    arcana: str | None = None,
    suit: str | None = None,
    count: str | None = None,
) -> Response:
    seed_val = resolve_seed(seed)
    cards = fisher_yates_shuffle(load_tarot_deck(), seed_val)

    if id is not None:
        matched_card = find_card_by_id(cards, id)
        if matched_card is None:
            return PlainTextResponse("Card not found")
        return JSONResponse(matched_card)

    if arcana is not None:
        cards = filter_cards_by_arcana(cards, arcana)
        if not cards:
            return PlainTextResponse("Invalid Arcana")

    if suit is not None:
        cards = filter_cards_by_suit(cards, suit)
        if not cards:
            return PlainTextResponse("Invalid Suit")

    if count is not None:
        try:
            requested_count = int(count)
        except ValueError:
            return PlainTextResponse("Not a number!")
        if requested_count > 78 or requested_count < 1:
            return PlainTextResponse("Invalid Number")
        cards = cards[: min(requested_count, len(cards))]

    return JSONResponse({"cards": cards, "seed": seed_val})


def game_id_exists(game_id: int) -> bool:
    if not pool:
        return False
    with pool.connection() as conn:
        with conn.cursor() as cur:
            cur.execute(
                "SELECT game_id FROM simple_leaderboard.users WHERE game_id = %s;",
                (game_id,),
            )
            return cur.fetchone() is not None


def record_new_user(email: str, game_name: str) -> str:
    assert pool is not None
    with pool.connection() as conn:
        with conn.cursor() as cur:
            cur.execute(
                """
                INSERT INTO simple_leaderboard.users (email, game_name)
                VALUES (%s, %s)
                ON CONFLICT (email, game_name)
                DO UPDATE SET email = EXCLUDED.email
                RETURNING game_id;
                """,
                (email, game_name),
            )
            row = cur.fetchone()
        conn.commit()
    return str(row[0]) if row else "1"


def record_new_score(game_id: int, score: int, user_name: str) -> Response:
    assert pool is not None
    if not game_id_exists(game_id):
        return PlainTextResponse("Invalid GameID")

    with pool.connection() as conn:
        with conn.cursor() as cur:
            cur.execute(
                """
                INSERT INTO simple_leaderboard.scores (game_id, score, user_name)
                VALUES (%s, %s, %s);
                """,
                (game_id, score, user_name),
            )
        conn.commit()
    return PlainTextResponse("Success")


def get_top_scores(game_id: int, count: int = 10) -> list[dict[str, int | str]]:
    assert pool is not None
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
                (game_id, count),
            )
            rows = cur.fetchall()
    return [{"userName": row[0], "score": row[1]} for row in rows]


def get_user_scores(game_id: int, user_name: str, count: int = 10) -> list[dict[str, int]]:
    assert pool is not None
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
                (game_id, user_name, count),
            )
            rows = cur.fetchall()
    return [{row[0]: row[1]} for row in rows]


def parse_int32(raw_value: str) -> int:
    parsed = int(raw_value)
    if parsed < PG_INT32_MIN or parsed > PG_INT32_MAX:
        raise ValueError("Integer out of 32-bit signed range")
    return parsed


def parse_positive_count(count_param: str | None, default: int = 10) -> int:
    if not count_param:
        return default
    parsed = parse_int32(count_param)
    if parsed < 1:
        raise ValueError("Count must be positive")
    return parsed


def is_within_varchar_limit(value: str | None, max_length: int = VARCHAR_MAX_LENGTH) -> bool:
    return bool(value) and len(value) <= max_length


def extract_bounded_string(
    payload: dict[str, Any], key: str, max_length: int = VARCHAR_MAX_LENGTH
) -> str | None:
    value = payload.get(key)
    if not isinstance(value, str):
        return None
    trimmed = value.strip()
    if not trimmed or len(trimmed) > max_length:
        return None
    return trimmed


@app.api_route("/apis/simple-leaderboard-api", methods=["GET", "POST"])
@app.api_route("/apis/simple-leaderboard-api/", methods=["GET", "POST"])
@app.api_route("/apis/simple-leaderboard", methods=["GET", "POST"])
@app.api_route("/apis/simple-leaderboard/", methods=["GET", "POST"])
async def handle_leaderboard_request(request: Request) -> Response:
    if not pool:
        return PlainTextResponse("Database connection unavailable", status_code=500)

    query_params = request.query_params
    action = query_params.get("action")
    game_id_param = query_params.get("gameID")
    score_param = query_params.get("score")
    user_name = query_params.get("userName")
    count_param = query_params.get("count")

    if action == "newScore":
        if not (game_id_param and score_param and is_within_varchar_limit(user_name)):
            return PlainTextResponse("Invalid Params")
        assert user_name is not None
        try:
            game_id = parse_int32(game_id_param)
            score_value = parse_int32(score_param)
        except ValueError:
            return PlainTextResponse("Invalid Params")
        return await run_in_threadpool(record_new_score, game_id, score_value, user_name)

    if action == "topScores":
        if not game_id_param:
            return PlainTextResponse("Invalid params")
        try:
            game_id = parse_int32(game_id_param)
            limit = parse_positive_count(count_param)
        except ValueError:
            return PlainTextResponse("Invalid params")
        scores = await run_in_threadpool(get_top_scores, game_id, limit)
        return JSONResponse(scores)

    if action == "userScores":
        if not (is_within_varchar_limit(user_name) and game_id_param):
            return PlainTextResponse("Invalid params")
        assert user_name is not None
        try:
            game_id = parse_int32(game_id_param)
            limit = parse_positive_count(count_param)
        except ValueError:
            return PlainTextResponse("Invalid params")
        scores = await run_in_threadpool(get_user_scores, game_id, user_name, limit)
        return JSONResponse(scores)

    try:
        payload = await request.json()
    except Exception:
        payload = {}

    if isinstance(payload, dict) and payload.get("action") == "addUser":
        email = extract_bounded_string(payload, "email")
        game_name = extract_bounded_string(payload, "gameName")
        if not (email and game_name):
            return PlainTextResponse("Invalid params")
        registered_game_id = await run_in_threadpool(record_new_user, email, game_name)
        return JSONResponse(registered_game_id)

    return PlainTextResponse("Invalid params")
