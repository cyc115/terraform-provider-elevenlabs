#!/usr/bin/env bash
# gen-gpg-key.sh — generate a dedicated GPG key for terraform-provider-elevenlabs releases.
#
# Usage:
#   ./gen-gpg-key.sh              # interactive: prompts for email and passphrase
#   ./gen-gpg-key.sh --dry-run    # prints the key config that would be used; no key generated
#
# Outputs (written to /tmp/):
#   /tmp/elevenlabs-release-public.asc   — public key to upload to registry.terraform.io
#   /tmp/elevenlabs-release-private.asc  — private key for GPG_PRIVATE_KEY GitHub secret
#
# After running, follow the printed instructions to set GitHub secrets and upload the public key.

set -euo pipefail

DRY_RUN=false
for arg in "$@"; do
  [[ "$arg" == "--dry-run" ]] && DRY_RUN=true
done

# --- resolve email ---
if [[ -n "${RELEASE_EMAIL:-}" ]]; then
  EMAIL="$RELEASE_EMAIL"
elif [[ "$DRY_RUN" == "true" ]]; then
  EMAIL="RELEASE_EMAIL@example.com"
else
  read -rp "Release email address for GPG key: " EMAIL
fi

KEY_NAME="Mike Chen / terraform-provider-elevenlabs release"
KEY_COMMENT="elevenlabs-tf-release"

GPG_BATCH_CONFIG="$(cat <<GPGEOF
%echo Generating ElevenLabs provider release key
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: ${KEY_NAME}
Name-Comment: ${KEY_COMMENT}
Name-Email: ${EMAIL}
Expire-Date: 0
%no-protection
%commit
%echo done
GPGEOF
)"

echo ""
echo "=== GPG Key Configuration ==="
echo "  Name:    ${KEY_NAME}"
echo "  Comment: ${KEY_COMMENT}"
echo "  Email:   ${EMAIL}"
echo "  Type:    RSA 4096"
echo "  Expiry:  none"
echo "  Output:  /tmp/elevenlabs-release-{public,private}.asc"
echo ""

if [[ "$DRY_RUN" == "true" ]]; then
  echo "[dry-run] Key generation skipped."
  echo "[dry-run] Batch config that would be used:"
  echo "---"
  echo "$GPG_BATCH_CONFIG"
  echo "---"
  echo "[dry-run] DONE — no key generated, no files written."
  exit 0
fi

# --- generate key ---
echo "Generating key (this may take a moment)..."
echo "$GPG_BATCH_CONFIG" | gpg --batch --gen-key

# --- find fingerprint of newly generated key ---
FINGERPRINT=$(gpg --list-secret-keys --with-colons "${EMAIL}" 2>/dev/null \
  | awk -F: '/^fpr:/ { print $10; exit }')

if [[ -z "$FINGERPRINT" ]]; then
  echo "ERROR: could not determine fingerprint for ${EMAIL}" >&2
  exit 1
fi

echo ""
echo "=== Key Generated ==="
echo "  Fingerprint: ${FINGERPRINT}"
echo ""

# --- export keys ---
gpg --armor --export "${FINGERPRINT}" > /tmp/elevenlabs-release-public.asc
gpg --armor --export-secret-keys "${FINGERPRINT}" > /tmp/elevenlabs-release-private.asc

echo "  Public key:  /tmp/elevenlabs-release-public.asc"
echo "  Private key: /tmp/elevenlabs-release-private.asc"
echo ""
echo "=== Next Steps ==="
echo ""
echo "1. Set GitHub secrets:"
echo "   gh secret set GPG_PRIVATE_KEY -R cyc115/terraform-provider-elevenlabs < /tmp/elevenlabs-release-private.asc"
echo "   gh secret set PASSPHRASE -R cyc115/terraform-provider-elevenlabs"
echo "   (enter an empty passphrase at the prompt, or your chosen passphrase)"
echo ""
echo "2. Upload the PUBLIC key to the Terraform Registry:"
echo "   https://registry.terraform.io/settings/gpg-keys"
echo "   Content: /tmp/elevenlabs-release-public.asc"
echo ""
echo "3. Tag and push to trigger the release workflow:"
echo "   git tag -a v0.2.0 -m 'publish-prep: first registry release'"
echo "   git push origin v0.2.0"
echo ""
echo "Fingerprint for reference: ${FINGERPRINT}"
