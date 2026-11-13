#!/usr/bin/env bash
set -euo pipefail

# Encrypted PostgreSQL -> GitHub Automated Backup Script
# Dumps all databases & roles from infra-postgres, compresses with gzip,
# encrypts with age, and pushes to private GitHub repo geraldfingburke/vps-db-backups.

BACKUP_ROOT="/opt/infra/backups"
REPO_DIR="${BACKUP_ROOT}/repo"
LOCAL_DIR="${BACKUP_ROOT}/local"
AGE_KEY_FILE="${BACKUP_ROOT}/keys/age.key"
SSH_KEY="/root/.ssh/id_ed25519_github_backup"
GITHUB_REPO="git@github.com:geraldfingburke/vps-db-backups.git"
TIMESTAMP=$(date -u +"%Y-%m-%d_%H-%M-%S_UTC")
DATE_DAY=$(date -u +"%Y-%m-%d")

mkdir -p "${LOCAL_DIR}"

# Extract public key from age key file
AGE_PUB_KEY=$(grep -o 'age1[a-z0-9]*' "${AGE_KEY_FILE}")

export GIT_SSH_COMMAND="ssh -i ${SSH_KEY} -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new"

# Initialize or pull backup repo
if [ ! -d "${REPO_DIR}/.git" ]; then
    rm -rf "${REPO_DIR}"
    git clone "${GITHUB_REPO}" "${REPO_DIR}"
    cd "${REPO_DIR}"
    git config user.name "VPS Backup Bot"
    git config user.email "backup@trigamedev.org"
else
    cd "${REPO_DIR}"
    git pull --rebase origin main || true
fi

mkdir -p "${REPO_DIR}/daily" "${REPO_DIR}/latest"

# 1. Dump entire Postgres cluster (roles + schemas + data)
RAW_DUMP="${LOCAL_DIR}/pg_cluster_${TIMESTAMP}.sql.gz"
ENC_DUMP="${LOCAL_DIR}/pg_cluster_${TIMESTAMP}.sql.gz.age"

docker exec infra-postgres pg_dumpall -U postgres | gzip -9 > "${RAW_DUMP}"

# 2. Encrypt with age
age -r "${AGE_PUB_KEY}" -o "${ENC_DUMP}" "${RAW_DUMP}"

# Remove unencrypted dump immediately so only encrypted backups rest on disk
rm -f "${RAW_DUMP}"

# 3. Copy into Git repo (keep daily snapshot + latest pointer)
cp "${ENC_DUMP}" "${REPO_DIR}/daily/pg_cluster_${DATE_DAY}.sql.gz.age"
cp "${ENC_DUMP}" "${REPO_DIR}/latest/pg_cluster_latest.sql.gz.age"

# 4. Prune local backups older than 14 days and repo daily backups older than 30 days
find "${LOCAL_DIR}" -name "*.age" -type f -mtime +14 -delete
find "${REPO_DIR}/daily" -name "*.age" -type f -mtime +30 -delete

# Write a helpful README in the backup repo if missing
cat > "${REPO_DIR}/README.md" << 'EOF'
# VPS Encrypted PostgreSQL Backups

Automated nightly backups from VPS (`129.121.165.252`).
All dumps are compressed with `gzip` and encrypted with [`age`](https://github.com/FiloSottile/age).

## How to Restore / Decrypt

```bash
age -d -i age.key latest/pg_cluster_latest.sql.gz.age | gunzip > restore.sql
psql -U postgres -f restore.sql
```
EOF

# 5. Commit and push to GitHub
cd "${REPO_DIR}"
git add -A
if ! git diff --cached --quiet; then
    git commit -m "Automated encrypted backup: ${TIMESTAMP}"
    git branch -M main
    git push -u origin main
    echo "[${TIMESTAMP}] Backup encrypted and pushed to GitHub successfully."
else
    echo "[${TIMESTAMP}] No changes to push."
fi
