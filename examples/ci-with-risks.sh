#!/usr/bin/env bash
set -euo pipefail

npm ci &
docker compose up -d
sleep 30
psql -h localhost -p 5432 -c 'select 1'
wait
