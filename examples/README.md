# Examples

Runnable Terraform examples for the `cyc115/elevenlabs` provider.

| Example | Description |
|---------|-------------|
| [clone-restaurant-agent](./clone-restaurant-agent/) | Create a new ConvAI agent cloned from an existing configuration |
| [import-webhook-phone](./import-webhook-phone/) | Import an existing workspace webhook and Twilio phone number into Terraform state |
| [replicate-agent](./replicate-agent/) | Replicate an agent configuration as a new Terraform-managed resource |

## Prerequisites

All examples require the provider binary installed locally:

```bash
# From repo root
make install
```

This installs `terraform-provider-elevenlabs` into `~/.terraform.d/plugins/`.

## Common Workflow

```bash
cd <example-dir>
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars — add your ElevenLabs API key and any required values
terraform init
terraform plan
terraform apply
```

## Security Notes

- **Never commit `terraform.tfvars`** — it contains secrets. Each example's `.gitignore` blocks it.
- **Terraform state is sensitive** — local state files (`terraform.tfstate`) are plaintext and may contain API keys or PII. Use an encrypted remote backend for team or CI use.
- **Always destroy unused resources** — ElevenLabs bills per agent; run `terraform destroy` when done.
