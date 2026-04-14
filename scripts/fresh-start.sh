#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

PRUNE=false
if [[ "${1:-}" == "--prune" ]]; then
  PRUNE=true
fi

echo "[1/4] Stopping and removing containers, networks, and volumes..."
docker compose down -v --remove-orphans

echo "[2/4] Rebuilding images from scratch (no cache)..."
docker compose build --no-cache

echo "[3/4] Starting fresh stack..."
docker compose up -d --force-recreate

if [[ "$PRUNE" == "true" ]]; then
  echo "[4/4] Pruning dangling Docker resources..."
  docker image prune -f
  docker container prune -f
  docker network prune -f
else
  echo "[4/4] Skipping prune (pass --prune to enable)."
fi

echo "Done. Fresh stack is up."
