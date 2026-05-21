#!/usr/bin/env bash
# Tests for gen-gpg-key.sh --dry-run behaviour.
# RED before the --dry-run flag exists; GREEN after.
#
# Run:  bash scripts/release/gen-gpg-key_test.sh

set -euo pipefail

SCRIPT="$(dirname "$0")/gen-gpg-key.sh"
PASS=0
FAIL=0

assert_contains() {
  local label="$1" expected="$2" actual="$3"
  if echo "$actual" | grep -qF "$expected"; then
    echo "  PASS: $label"
    PASS=$((PASS+1))
  else
    echo "  FAIL: $label"
    echo "        expected substring: $expected"
    echo "        actual output:      $(echo "$actual" | head -3)"
    FAIL=$((FAIL+1))
  fi
}

assert_not_contains() {
  local label="$1" unexpected="$2" actual="$3"
  if ! echo "$actual" | grep -qF "$unexpected"; then
    echo "  PASS: $label"
    PASS=$((PASS+1))
  else
    echo "  FAIL: $label — output unexpectedly contained: $unexpected"
    FAIL=$((FAIL+1))
  fi
}

echo "=== gen-gpg-key.sh --dry-run tests ==="
echo ""

OUTPUT=$(RELEASE_EMAIL="test@example.com" bash "$SCRIPT" --dry-run 2>&1)

assert_contains "prints dry-run marker"      "[dry-run] Key generation skipped." "$OUTPUT"
assert_contains "prints key config header"   "=== GPG Key Configuration ==="     "$OUTPUT"
assert_contains "shows RSA 4096"             "RSA 4096"                           "$OUTPUT"
assert_contains "shows email from env"       "test@example.com"                   "$OUTPUT"
assert_contains "shows done marker"          "[dry-run] DONE"                     "$OUTPUT"
assert_contains "shows batch config"         "Key-Type: RSA"                      "$OUTPUT"
assert_contains "shows key length"           "Key-Length: 4096"                   "$OUTPUT"
assert_contains "shows name-real"            "Mike Chen"                          "$OUTPUT"
assert_not_contains "no files written"       "/tmp/elevenlabs-release-public.asc: created" "$OUTPUT"
assert_not_contains "gpg not called"         "Generating key"                     "$OUTPUT"

echo ""
echo "Results: ${PASS} passed, ${FAIL} failed"
[[ $FAIL -eq 0 ]] && exit 0 || exit 1
