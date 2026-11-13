import os
from contextlib import asynccontextmanager
from fastapi import FastAPI
from psycopg_pool import ConnectionPool

DATABASE_URL = os.getenv("DATABASE_URL", "")
pool: ConnectionPool | None = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    global pool
    pool = ConnectionPool(conninfo=DATABASE_URL, min_size=1, max_size=10, open=True)
    yield
    if pool:
        pool.close()

app = FastAPI(title="Shared FastAPI Backend", lifespan=lifespan)

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
    }
