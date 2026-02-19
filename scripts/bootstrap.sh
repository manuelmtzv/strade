#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

copy_if_missing() {
  local src="$1"
  local dest="$2"

  if [[ -f "$dest" ]]; then
    echo "[skip] $dest already exists"
    return
  fi

  cp "$src" "$dest"
  echo "[ok] created $dest from $src"
}

generate_password() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -base64 24 | tr -d '=+/' | cut -c1-24
  else
    tr -dc 'A-Za-z0-9' </dev/urandom | head -c 24
  fi
}

ensure_redis_password() {
  local env_file="$1"

  if [[ ! -f "$env_file" ]]; then
    return
  fi

  if grep -q '^REDIS_PASSWORD=' "$env_file"; then
    local current
    current="$(grep '^REDIS_PASSWORD=' "$env_file" | head -n1 | cut -d'=' -f2-)"
    if [[ -n "$current" ]]; then
      echo "[skip] REDIS_PASSWORD already set in $env_file"
      return
    fi

    local password
    password="$(generate_password)"
    sed -i "s/^REDIS_PASSWORD=$/REDIS_PASSWORD=${password}/" "$env_file"
    echo "[ok] REDIS_PASSWORD generated in $env_file"
    return
  fi

  echo "REDIS_PASSWORD=$(generate_password)" >> "$env_file"
  echo "[ok] REDIS_PASSWORD appended to $env_file"
}

copy_if_missing .env.example .env.dev
copy_if_missing .env.example .env

ensure_redis_password .env

echo "Done. Review .env and .env.dev before starting services."
