# terraform-provider-elevenlabs

![Version](https://img.shields.io/badge/version-0.2.0-blue)

Terraform provider for [ElevenLabs Conversational AI](https://elevenlabs.io/conversational-ai) — agents, webhooks, phone numbers.

## Installation

```hcl
terraform {
  required_providers {
    elevenlabs = {
      source  = "cyc115/elevenlabs"
      version = "~> 0.2"
    }
  }
}
```

```bash
terraform init
```

The provider is published at [registry.terraform.io/providers/cyc115/elevenlabs](https://registry.terraform.io/providers/cyc115/elevenlabs).

## Provider configuration

```hcl
provider "elevenlabs" {
  api_key = var.elevenlabs_api_key   # or set ELEVENLABS_API_KEY env var
}
```

## Resources

- [`elevenlabs_convai_agent`](docs/resources/convai_agent.md) — full CRUD for ConvAI agents
- [`elevenlabs_workspace_webhook`](docs/resources/workspace_webhook.md) — workspace-level post-call webhooks
- [`elevenlabs_convai_phone_number`](docs/resources/convai_phone_number.md) — Twilio phone numbers imported into ElevenLabs

See the [docs/](docs/) directory or the [Terraform Registry page](https://registry.terraform.io/providers/cyc115/elevenlabs) for full attribute documentation.

## Examples

See [examples/](examples/) for ready-to-run configurations:

| Example | Description |
|---------|-------------|
| [clone-restaurant-agent](examples/clone-restaurant-agent/) | Create a new ConvAI agent |
| [import-webhook-phone](examples/import-webhook-phone/) | Import webhook + phone into state |
| [replicate-agent](examples/replicate-agent/) | Replicate an agent configuration |

## Resource: `elevenlabs_convai_agent`

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string (computed) | Agent ID |
| `name` | string (required) | Agent name |
| `tags` | list(string) | Tags |
| `asr_quality` | string | ASR quality (`high`, `low`) |
| `asr_provider` | string | ASR provider |
| `asr_user_input_audio_format` | string | User audio format |
| `asr_keywords` | list(string) | ASR hint keywords |
| `turn_timeout` | number | Turn timeout seconds |
| `turn_silence_end_call_timeout` | number | Silence-to-hangup seconds (-1 = disabled) |
| `turn_mode` | string | `turn` or `interruptible` |
| `turn_eagerness` | string | `low`, `normal`, `high` |
| `turn_model` | string | Turn detection model |
| `turn_speculative_turn` | bool | Enable speculative turns |
| `turn_soft_timeout_seconds` | number | Soft timeout seconds |
| `turn_soft_timeout_message` | string | Soft timeout filler message |
| `turn_soft_timeout_use_llm` | bool | Use LLM for soft timeout message |
| `tts_model_id` | string | TTS model |
| `tts_voice_id` | string | Voice ID |
| `tts_expressive_mode` | bool | Enable expressive mode |
| `tts_optimize_streaming_latency` | number | Latency optimization level (0-4) |
| `tts_stability` | number | Voice stability (0-1) |
| `tts_speed` | number | Speech speed |
| `tts_similarity_boost` | number | Similarity boost (0-1) |
| `tts_agent_output_audio_format` | string | Agent audio format |
| `tts_text_normalisation_type` | string | Text normalization mode |
| `conversation_max_duration_seconds` | number | Max call duration |
| `conversation_client_events` | list(string) | Client event types to stream |
| `conversation_monitoring_events` | list(string) | Monitoring event types |
| `conversation_source_attribution` | bool | Enable source attribution |
| `conversation_bg_music_source_type` | string | Background music type (`preset`, `url`) |
| `conversation_bg_music_source_id` | string | Background music ID or URL |
| `conversation_bg_music_volume` | number | Background music volume (0-1) |
| `conversation_bg_music_crossfade` | bool | Crossfade loop |
| `conversation_file_input_enabled` | bool | Allow file uploads |
| `conversation_file_input_max_files` | number | Max files per conversation |
| `vad_background_voice_detection` | bool | Background voice detection |
| `agent_first_message` | string | Opening message from agent |
| `agent_language` | string | Primary language code |
| `agent_dynamic_vars` | map(string) | Dynamic variable placeholders |
| `prompt_text` | string | System prompt |
| `prompt_llm` | string | LLM model ID |
| `prompt_reasoning_effort` | string | `minimal`, `low`, `medium`, `high` |
| `prompt_temperature` | number | LLM temperature |
| `prompt_max_tokens` | number | Max tokens (-1 = unlimited) |
| `prompt_timezone` | string | Agent timezone (IANA) |
| `prompt_cascade_timeout` | number | LLM cascade timeout seconds |
| `prompt_enable_parallel_tool_calls` | bool | Parallel tool calls |
| `prompt_ignore_default_personality` | bool | Suppress default personality |
| `prompt_backup_llm_preference` | string | Backup LLM strategy |
| `prompt_rag_enabled` | bool | Enable RAG |
| `prompt_rag_embedding_model` | string | RAG embedding model |
| `prompt_rag_max_vector_distance` | number | RAG max vector distance |
| `prompt_rag_max_documents_length` | number | RAG max document length |
| `prompt_rag_max_chunks` | number | RAG max retrieved chunks |
| `data_collection` | map(object) | Post-call data collection fields |
| `webhook_post_call_id` | string | Post-call webhook ID |
| `webhook_events` | list(string) | Webhook event types |
| `webhook_transcript_format` | string | Transcript format (`json`, `text`) |
| `ps_agent_concurrency_limit` | number | Concurrent call limit (-1 = unlimited) |
| `ps_daily_limit` | number | Daily call limit |
| `ps_bursting_enabled` | bool | Allow bursting above daily limit |
| `ps_record_voice` | bool | Record voice |
| `ps_retention_days` | number | Conversation retention days (-1 = forever) |
| `ps_delete_transcript_and_pii` | bool | Delete transcript and PII |
| `ps_delete_audio` | bool | Delete audio recordings |
| `ps_zero_retention_mode` | bool | Zero retention mode |
| `ps_history_redaction_enabled` | bool | Conversation history redaction |
| `ps_guardrails_focus_enabled` | bool | Focus guardrail |
| `ps_guardrails_prompt_injection_enabled` | bool | Prompt injection guardrail |
| `ps_guardrails_content_execution_mode` | string | Content guardrail mode (`streaming`) |
| `ps_guardrails_trigger_action` | string | Guardrail trigger action (`end_call`) |
| `ps_auth_enable_auth` | bool | Require auth for widget calls |
| `ps_auth_require_origin_header` | bool | Require Origin header |
| `ps_analysis_llm` | string | Post-call analysis LLM |
| `ps_archived` | bool | Archive agent |
| `coaching_type` | string | Coaching type (`coached`) |

### `data_collection` object schema

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Field type (`string`) |
| `description` | string | LLM extraction instruction |
| `enum` | list(string) | Allowed values (empty = free text) |
| `is_system_provided` | bool | System-provided field |
| `dynamic_variable` | string | Dynamic variable binding |
| `constant_value` | string | Constant value override |
| `llm` | string | LLM override for extraction |

## Resource: `elevenlabs_workspace_webhook`

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string (computed) | Webhook ID |
| `name` | string (required) | Webhook name |
| `webhook_url` | string (required) | HTTPS endpoint URL |
| `auth_type` | string (computed) | Auth type (`hmac`) |
| `retry_enabled` | bool (computed) | Retry on failure |
| `is_disabled` | bool (computed) | Disabled state |
| `is_auto_disabled` | bool (computed) | Auto-disabled after repeated failures |
| `created_at_unix` | number (computed) | Creation timestamp |
| `secret` | string (computed, sensitive) | HMAC signing secret (only available on create) |

> **Note:** `secret` is written to state on create only and is not available after import.

## Resource: `elevenlabs_convai_phone_number`

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string (computed) | Phone number ID |
| `phone_number` | string (required, forces replace) | E.164 phone number |
| `telephony_provider` | string (required, forces replace) | Provider (`twilio`) |
| `label` | string | Human-readable label |
| `assigned_agent` | string | Agent ID assigned to this number |
| `supports_inbound` | bool (computed) | Inbound call support |
| `supports_outbound` | bool (computed) | Outbound call support |
| `twilio_account_sid` | string (sensitive) | Twilio account SID |
| `twilio_auth_token` | string (sensitive) | Twilio auth token |

---

## Advanced — local development

To test against a locally built binary before the provider is published (or to use
a development build after publication):

```hcl
# ~/.terraformrc
provider_installation {
  dev_overrides {
    "registry.terraform.io/cyc115/elevenlabs" = "/Users/<you>/.terraform.d/plugins/registry.terraform.io/cyc115/elevenlabs/0.1.1/darwin_arm64"
  }
  direct {}
}
```

```bash
make install   # builds and installs to the dev_overrides path
```

> **Note:** When using `dev_overrides`, `terraform init` will skip version locking
> for this provider. Remove the `dev_overrides` block to use the registry version.

## License

MIT
