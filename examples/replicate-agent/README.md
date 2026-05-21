# Example: replicate-agent

Replicate an ElevenLabs ConvAI agent configuration as a new Terraform-managed resource. Useful for blue/green agent deploys or environment promotion.

## Quickstart

```bash
# 1. Install the provider
cd ../../ && make install

# 2. Set up variables
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars — set elevenlabs_api_key to your real sk_... key

# 3. Preview and apply
terraform init
terraform plan
terraform apply
```

## Variables

| Name | Description | Required |
|------|-------------|----------|
| `elevenlabs_api_key` | ElevenLabs API key (`sk_...`). Also accepted via `ELEVENLABS_API_KEY` env var. | Yes |

## Drift Detection

```bash
terraform plan -refresh-only
```

## Idempotency

A second `terraform apply` immediately after the first should report **0 changes**.

> **Note on prompt whitespace:** this example preserves a trailing newline in `prompt_text`
> (no `trimspace`) to match the ElevenLabs API's stored value exactly. If you see a perpetual
> diff on `prompt_text`, check whether the live agent's prompt ends with a newline.

## Destroy

```bash
terraform destroy
```

This removes the replica agent from ElevenLabs. The source agent is not affected.

## Security

> **Warning:** Terraform state may contain sensitive values (e.g. API keys, dynamic vars).
> - `terraform.tfvars` is `.gitignored` — never commit it.
> - For team or CI use, store state in an encrypted remote backend
>   (e.g. S3 + KMS, Terraform Cloud). Local state is plaintext.
