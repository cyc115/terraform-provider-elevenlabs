# Publish Checklist — terraform-provider-elevenlabs

Step-by-step user actions to publish the provider to the official Terraform Registry
for the first time. Complete these **in order** after merging the `cyc_lq-publish-prep` PR.

---

## Step 1 — Generate a dedicated GPG signing key

```bash
cd ~/github/terraform-provider-elevenlabs
bash scripts/release/gen-gpg-key.sh
```

- Enter your release email address when prompted.
- The key is generated **without a passphrase** (`%no-protection`). The private key is
  protected at rest by GitHub Secrets — no passphrase is the standard HashiCorp pattern for CI.
- Note the printed fingerprint.
- This writes two files:
  - `/tmp/elevenlabs-release-public.asc` — public key (upload to registry)
  - `/tmp/elevenlabs-release-private.asc` — private key (store as GitHub secret)

---

## Step 2 — Set the `GPG_PRIVATE_KEY` GitHub secret

```bash
gh secret set GPG_PRIVATE_KEY \
  -R cyc115/terraform-provider-elevenlabs \
  < /tmp/elevenlabs-release-private.asc
```

Verify it was set:

```bash
gh secret list -R cyc115/terraform-provider-elevenlabs
```

---

## Step 3 — Set the `PASSPHRASE` GitHub secret

The key generated in Step 1 has no passphrase. Set the secret to an empty string
(the `crazy-max/ghaction-import-gpg` action still requires the secret to exist):

```bash
printf '' | gh secret set PASSPHRASE -R cyc115/terraform-provider-elevenlabs
```

---

## Step 4 — Upload the public key to the Terraform Registry

1. Open: <https://registry.terraform.io/settings/gpg-keys>
2. Click **Add Key**.
3. Paste the contents of `/tmp/elevenlabs-release-public.asc`.
4. Click **Add GPG Key** to save.

The namespace `cyc115` must already be linked to your GitHub account (it is).

---

## Step 5 — Tag and push to trigger the release workflow

```bash
git tag -a v0.2.0 -m "publish-prep: first Terraform Registry release"
git push origin v0.2.0
```

This triggers `.github/workflows/release.yml`, which:
- Builds binaries for all supported platforms
- Creates a SHA256SUMS file and signs it with your GPG key
- Uploads all artifacts as a GitHub Release

---

## Step 6 — Verify the GitHub release

1. Open: <https://github.com/cyc115/terraform-provider-elevenlabs/actions>
2. Wait for the **Release** workflow to complete (green).
3. Open: <https://github.com/cyc115/terraform-provider-elevenlabs/releases/tag/v0.2.0>
4. Confirm the release contains:
   - `terraform-provider-elevenlabs_0.2.0_SHA256SUMS`
   - `terraform-provider-elevenlabs_0.2.0_SHA256SUMS.sig`
   - Zip archives for all platforms (linux/darwin/freebsd/windows × amd64/arm64/etc.)

---

## Step 7 — Publish to the Terraform Registry

1. Open: <https://registry.terraform.io/publish/provider>
2. Select the `cyc115/terraform-provider-elevenlabs` repository.
3. Verify the GPG key from Step 4 is listed.
4. Click **Publish Provider**.

The registry will index the provider within a few minutes.

---

## Step 8 — Smoke test

In a temporary directory, create a minimal Terraform config:

```hcl
terraform {
  required_providers {
    elevenlabs = {
      source  = "cyc115/elevenlabs"
      version = "0.2.0"
    }
  }
}

provider "elevenlabs" {
  api_key = "sk_test"
}
```

Then run:

```bash
terraform init
```

A successful `terraform init` that downloads the provider from the registry confirms
end-to-end publication.

---

## Post-publish

- Update the version badge in `README.md` to `0.2.0` (or let it float via `latest`).
- Announce in relevant channels if applicable.
- Delete the `/tmp/elevenlabs-release-private.asc` file from your machine:
  ```bash
  rm /tmp/elevenlabs-release-private.asc /tmp/elevenlabs-release-public.asc
  ```
  The private key is now safely stored in GitHub secrets only.
