#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${QA_BASE_URL:-http://127.0.0.1:8080}
echo "Running QA against $BASE_URL"

do_post() {
  local url=$1
  local data=$2
  # Capture body and HTTP status
  resp=$(curl -sS -w "\n%{http_code}" -X POST "$BASE_URL$url" -d "$data") || true
  body=$(echo "$resp" | sed '$d')
  code=$(echo "$resp" | tail -n1)
  echo "$code $body"
  return $code
}

do_get() {
  local url=$1
  resp=$(curl -sS -w "\n%{http_code}" "$BASE_URL$url") || true
  body=$(echo "$resp" | sed '$d')
  code=$(echo "$resp" | tail -n1)
  echo "$code $body"
  return $code
}

echo "1) Create Lesson"
out=$(curl -sS -w "\n%{http_code}" -X POST "$BASE_URL/lessons" -d "parrot is awesome!")
body=$(echo "$out" | sed '$d')
code=$(echo "$out" | tail -n1)
echo "Response: $code $body"
if [[ $code -lt 200 || $code -ge 300 ]]; then
  echo "Create lesson failed"
  exit 1
fi

root_id=$(python - <<PY
import sys, json
print(json.loads(sys.stdin.read())['root_id'])
PY
<<<"$body")
content_id=$(python - <<PY
import sys, json
print(json.loads(sys.stdin.read())['content_id'])
PY
<<<"$body")

echo "root_id=$root_id content_id=$content_id"

echo "2) Create Version"
out=$(curl -sS -w "\n%{http_code}" -X POST "$BASE_URL/lessons/$root_id/versions" -d "parrot v2")
body=$(echo "$out" | sed '$d')
code=$(echo "$out" | tail -n1)
echo "Response: $code $body"
if [[ $code -lt 200 || $code -ge 300 ]]; then
  echo "Create version failed"
  exit 1
fi

version_number=$(python - <<PY
import sys, json
obj = json.loads(sys.stdin.read())
print(obj.get('version_number'))
PY
<<<"$body")
version_content_id=$(python - <<PY
import sys, json
obj = json.loads(sys.stdin.read())
print(obj.get('content_id'))
PY
<<<"$body")
echo "version_number=$version_number version_content_id=$version_content_id"

echo "3) Publish Version"
out=$(curl -sS -w "\n%{http_code}" -X POST "$BASE_URL/lessons/$root_id/versions/$version_number/publish")
body=$(echo "$out" | sed '$d')
code=$(echo "$out" | tail -n1)
echo "Response: $code $body"
if [[ $code -lt 200 || $code -ge 300 ]]; then
  echo "Publish failed"
  exit 1
fi

echo "4) Get Version"
out=$(curl -sS -w "\n%{http_code}" "$BASE_URL/lessons/$root_id/versions/$version_number")
body=$(echo "$out" | sed '$d')
code=$(echo "$out" | tail -n1)
echo "Response: $code $body"
if [[ $code -lt 200 || $code -ge 300 ]]; then
  echo "Get version failed"
  exit 1
fi

echo "5) Diff"
out=$(curl -sS -w "\n%{http_code}" "$BASE_URL/diff?from=$content_id&to=$version_content_id")
body=$(echo "$out" | sed '$d')
code=$(echo "$out" | tail -n1)
echo "Response: $code $body"
if [[ $code -lt 200 || $code -ge 300 ]]; then
  echo "Diff failed"
  exit 1
fi

echo "6) Clone Lesson (from version $version_number)"
out=$(curl -sS -w "\n%{http_code}" -X POST "$BASE_URL/lessons/clone" -H "Content-Type: application/json" -d "{\"from_root_id\":\"$root_id\",\"from_version\":$version_number}")
body=$(echo "$out" | sed '$d')
code=$(echo "$out" | tail -n1)
echo "Response: $code $body"
if [[ $code -lt 200 || $code -ge 300 ]]; then
  echo "Clone failed"
  exit 1
fi

clone_root=$(python - <<PY
import sys, json
print(json.loads(sys.stdin.read())['root_id'])
PY
<<<"$body")
clone_content=$(python - <<PY
import sys, json
print(json.loads(sys.stdin.read())['content_id'])
PY
<<<"$body")
echo "clone_root=$clone_root clone_content=$clone_content"

echo "7) Verify clone independence: create new version on original and diff against clone"
out=$(curl -sS -w "\n%{http_code}" -X POST "$BASE_URL/lessons/$root_id/versions" -d "original changed")
body=$(echo "$out" | sed '$d')
code=$(echo "$out" | tail -n1)
echo "Response: $code $body"
if [[ $code -lt 200 || $code -ge 300 ]]; then
  echo "Create version on original failed"
  exit 1
fi

new_original_content=$(python - <<PY
import sys, json
print(json.loads(sys.stdin.read())['content_id'])
PY
<<<"$body")
echo "new_original_content=$new_original_content"

out=$(curl -sS -w "\n%{http_code}" "$BASE_URL/diff?from=$new_original_content&to=$clone_content")
body=$(echo "$out" | sed '$d')
code=$(echo "$out" | tail -n1)
echo "Response: $code $body"
if [[ $code -lt 200 || $code -ge 300 ]]; then
  echo "Diff after original change failed"
  exit 1
fi

echo "QA flow completed successfully"