#!/usr/bin/env bash
set -euo pipefail

# This fixture demonstrates intentional suppressions.
# hermesscan:disable-next-line HMS0001 -- fake service has deterministic startup in this demo
sleep 30

# hermesscan:disable-next-line HMS0002 -- fixture intentionally documents PostgreSQL default port
echo "postgres://localhost:5432/app"
