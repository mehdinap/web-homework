#!/usr/bin/env bash
#
# Prints a JWT for the demo user.
#
#   curl -H "Authorization: Bearer $(./scripts/token.sh)" http://localhost/get
#
# Usage: ./scripts/token.sh [base-url]
set -euo pipefail

BASE_URL="${1:-http://localhost}"
LOGIN_USER="${LOGIN_USER:-student}"
LOGIN_PASSWORD="${LOGIN_PASSWORD:-hunter2}"

curl -fsS -X POST -H 'Content-Type: application/json' \
  -d "{\"username\":\"${LOGIN_USER}\",\"password\":\"${LOGIN_PASSWORD}\"}" \
  "${BASE_URL}/auth/login" \
  | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p'
