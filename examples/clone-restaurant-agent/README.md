# Example: clone-restaurant-agent

Clone an existing ElevenLabs ConvAI agent into a new Terraform-managed resource.

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

## Destroy

```bash
terraform destroy
```

This removes the cloned agent from ElevenLabs. The source agent is not affected.

## Security

> **Warning:** Terraform state may contain sensitive values (e.g. API keys, dynamic vars).
> - `terraform.tfvars` is `.gitignored` — never commit it.
> - For team or CI use, store state in an encrypted remote backend
>   (e.g. S3 + KMS, Terraform Cloud). Local state is plaintext.
