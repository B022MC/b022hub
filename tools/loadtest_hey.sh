#!/usr/bin/env bash

set -euo pipefail

if ! command -v hey >/dev/null 2>&1; then
  echo "error: hey 未安装。macOS 可用 'brew install hey'，Linux 可自行下载二进制。" >&2
  exit 1
fi

BASE_URL="${BASE_URL:-https://your-host.example.com}"
ENDPOINT="${ENDPOINT:-${BASE_URL%/}/v1/chat/completions}"
API_KEY="${API_KEY:-}"
MODEL="${MODEL:-}"
PROMPT="${PROMPT:-Generate exactly 200 short numbered lines without markdown.}"
MAX_TOKENS="${MAX_TOKENS:-256}"
STREAM="${STREAM:-false}"
CONCURRENCY_STEPS="${CONCURRENCY_STEPS:-5 10 20 40 80}"
STAGE_DURATION="${STAGE_DURATION:-60s}"
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-120}"
PAUSE_SECONDS="${PAUSE_SECONDS:-5}"
OUTPUT_DIR="${OUTPUT_DIR:-./tmp/loadtest-hey}"

if [[ -z "${API_KEY}" ]]; then
  echo "error: 请设置 API_KEY 环境变量。" >&2
  exit 1
fi

if [[ -z "${MODEL}" ]]; then
  echo "error: 请设置 MODEL 环境变量。" >&2
  exit 1
fi

json_escape() {
  local value="$1"
  value="${value//\\/\\\\}"
  value="${value//\"/\\\"}"
  value="${value//$'\n'/\\n}"
  value="${value//$'\r'/\\r}"
  value="${value//$'\t'/\\t}"
  printf '%s' "${value}"
}

mkdir -p "${OUTPUT_DIR}"
payload_file="$(mktemp)"
trap 'rm -f "${payload_file}"' EXIT

cat > "${payload_file}" <<EOF
{
  "model": "$(json_escape "${MODEL}")",
  "messages": [
    {
      "role": "user",
      "content": "$(json_escape "${PROMPT}")"
    }
  ],
  "max_tokens": ${MAX_TOKENS},
  "stream": ${STREAM}
}
EOF

echo "endpoint: ${ENDPOINT}"
echo "model: ${MODEL}"
echo "stage_duration: ${STAGE_DURATION}"
echo "concurrency_steps: ${CONCURRENCY_STEPS}"
echo "output_dir: ${OUTPUT_DIR}"
echo

for concurrency in ${CONCURRENCY_STEPS}; do
  report_file="${OUTPUT_DIR}/hey-c${concurrency}.txt"

  echo "==== concurrency=${concurrency} duration=${STAGE_DURATION} ===="
  hey \
    -z "${STAGE_DURATION}" \
    -c "${concurrency}" \
    -m POST \
    -T "application/json" \
    -H "Authorization: Bearer ${API_KEY}" \
    -D "${payload_file}" \
    -t "${TIMEOUT_SECONDS}" \
    "${ENDPOINT}" | tee "${report_file}"

  echo "saved: ${report_file}"
  echo
  sleep "${PAUSE_SECONDS}"
done

echo "done. 建议同时观察："
echo "  1) /api/v1/admin/ops/concurrency"
echo "  2) /api/v1/admin/ops/user-concurrency"
echo "  3) 机器 CPU / 内存 / Redis / PostgreSQL 连接数"
