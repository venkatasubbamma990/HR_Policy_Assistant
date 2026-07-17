#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source <(tr -d '\r' < .env)
  set +a
fi

resolve_wsl_url() {
  local raw="${OPENAI_BASE_URL:-http://127.0.0.1:1234/v1}"
  if ! grep -qi microsoft /proc/version 2>/dev/null; then
    printf '%s' "$raw"
    return
  fi
  if [[ "${HR_WSL_LLM_DISABLE:-}" == "1" ]]; then
    printf '%s' "$raw"
    return
  fi
  if [[ "$raw" != *127.0.0.1* && "$raw" != *localhost* ]]; then
    printf '%s' "$raw"
    return
  fi

  local gateway port
  gateway="$(awk '$2 == "00000000" { print $3; exit }' /proc/net/route 2>/dev/null || true)"
  if [[ -z "$gateway" ]]; then
    printf '%s' "$raw"
    return
  fi

  gateway="$(printf '%d.%d.%d.%d\n' \
    "0x${gateway:6:2}" "0x${gateway:4:2}" "0x${gateway:2:2}" "0x${gateway:0:2}")"
  port="${HR_WSL_LLM_PORT:-12345}"
  printf 'http://%s:%s/v1' "$gateway" "$port"
}

BASE_URL="$(resolve_wsl_url)"
MODELS_URL="${BASE_URL%/}/models"

echo "Checking LLM endpoint: $MODELS_URL"

if curl -sf --connect-timeout 3 "$MODELS_URL" >/dev/null; then
  echo "OK: LLM server is reachable."
  exit 0
fi

echo "FAIL: cannot reach LLM server at $MODELS_URL"

if grep -qi microsoft /proc/version 2>/dev/null; then
  cat <<'EOF'

WSL cannot reach LM Studio on Windows via localhost.

Fix (one-time, run in PowerShell as Administrator):
  powershell -ExecutionPolicy Bypass -File scripts/setup-wsl-llm.ps1

Then ensure LM Studio is running (Developer -> Start Server) and retry:
  make check-llm

Alternative: enable mirrored networking in %USERPROFILE%\.wslconfig:
  [wsl2]
  networkingMode=mirrored
Then run: wsl --shutdown
EOF
else
  echo "Ensure LM Studio is running (Developer -> Start Server) and OPENAI_BASE_URL is correct."
fi

exit 1
