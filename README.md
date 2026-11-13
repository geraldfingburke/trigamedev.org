# TriGameDev.org — Community Platform & Shared VPS Infrastructure

This repository contains the **TriGameDev.org** community web application alongside the shared VPS Docker infrastructure (`129.121.165.252`).

## Stack Overview
* **Frontend (`web/`):** Vue 3 + Vite + Tailwind CSS v4 (Lofi Sunset & Twilight Synth dashboard aesthetic)
* **Backend (`infra/services/api-go/`):** Go 1.23 REST API with Discord OAuth2, Guild/OG Role verification, Single-Redemption Key Sharing Pool (`SELECT ... FOR UPDATE SKIP LOCKED`), Direct-Download Ebook Library, and OG Admin CMS.
* **Auxiliary Python API (`infra/services/api-fastapi/`):** Python 3.12 FastAPI with `psycopg3` pool.
* **Database (`infra/postgres/`):** Shared PostgreSQL 17 (`trigamedev` schema) with automated 6-hour `age`-encrypted backups pushed to a private GitHub backup repository.
* **Reverse Proxy (`infra/caddy/`):** Caddy 2 Alpine serving static sites and automatic HTTPS.

## Local Development

```bash
cd web
npm install
npm run dev
```

## Production Deployment

```bash
cd web && npm run build
scp -i ~/.ssh/vps -r dist/* root@129.121.165.252:/opt/infra/sites/trigamedev.org/
scp -i ~/.ssh/vps -r dist/* root@129.121.165.252:/opt/infra/sites/default/
```
