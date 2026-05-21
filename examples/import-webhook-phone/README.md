# Example: import-webhook-phone

Import an existing ElevenLabs workspace webhook and Twilio phone number into Terraform state.

## Quickstart

```bash
# 1. Install the provider
cd ../../ && make install

# 2. Set up variables
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars — fill in all required values

# 3. Initialize and import existing resources
terraform init
terraform import elevenlabs_workspace_webhook.main <webhook_id>
terraform import elevenlabs_convai_phone_number.main <phone_id>

# 4. Verify no drift (plan should show 0 changes after import)
terraform plan
```

## Variables

| Name | Description | Required |
|------|-------------|----------|
| `elevenlabs_api_key` | ElevenLabs API key (`sk_...`). Also accepted via `ELEVENLABS_API_KEY` env var. | Yes |
| `webhook_url` | HTTPS endpoint for ElevenLabs webhook events | Yes |
| `webhook_name` | Display name for the webhook resource | No (default: `my-webhook`) |
| `phone_number` | E.164 phone number to import (e.g. `+15551234567`) | Yes |
| `phone_label` | Human-readable label for the phone number | No (default: `my-phone`) |

## Drift Detection

```bash
terraform plan -refresh-only
```

After import, the plan should be empty. Any non-empty plan means the config drifted from live state.

## Idempotency

```bash
terraform apply   # first run: 0 changes (import only)
terraform apply   # second run: also 0 changes
```

## Destroy

```bash
terraform destroy
```

This deletes the webhook and phone number from ElevenLabs. **Ensure no agents are using them before destroying.**

## Security

> **Warning:** Terraform state may contain sensitive values (e.g. API keys, phone numbers, webhook secrets).
> - `terraform.tfvars` is `.gitignored` — never commit it.
> - HMAC webhook secret is stored in state — use an encrypted remote backend
>   (e.g. S3 + KMS, Terraform Cloud) for team or CI use.
