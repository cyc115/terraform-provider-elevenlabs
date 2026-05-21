# terraform-provider-elevenlabs

Terraform provider for [ElevenLabs Conversational AI](https://elevenlabs.io/conversational-ai) — agents, webhooks, phone numbers.

## Resources

- `elevenlabs_convai_agent` — full CRUD for ConvAI agents including `conversation_config` and `platform_settings`
- `elevenlabs_workspace_webhook` — workspace-level post-call webhooks
- `elevenlabs_convai_phone_number` — Twilio phone numbers imported into ElevenLabs

## Dev setup

```hcl
# ~/.terraformrc
provider_installation {
  dev_overrides {
    "registry.terraform.io/cyc115/elevenlabs" = "/Users/<you>/.terraform.d/plugins/registry.terraform.io/cyc115/elevenlabs/0.0.1-dev/darwin_arm64"
  }
  direct {}
}
```

```bash
make install
```

## Provider config

```hcl
provider "elevenlabs" {
  api_key = var.elevenlabs_api_key  # or ELEVENLABS_API_KEY env var
}
```

## License

MIT
