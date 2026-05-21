terraform {
  required_providers {
    elevenlabs = {
      source  = "registry.terraform.io/cyc115/elevenlabs"
      version = "~> 0.1"
    }
  }
}

provider "elevenlabs" {
  api_key = ""  # falls back to ELEVENLABS_API_KEY env var
}

resource "elevenlabs_workspace_webhook" "lifeos_mike_prod" {
  name        = "lifeos-mike-prod"
  webhook_url = "https://caller.mikechen.info/webhook/elevenlabs"
  auth_type   = "hmac"
  retry_enabled = false
}

resource "elevenlabs_convai_phone_number" "lifeos_mike" {
  phone_number       = "+18166536821"
  telephony_provider = "twilio"
  label              = "lifeos"
}

output "webhook_id" {
  value = elevenlabs_workspace_webhook.lifeos_mike_prod.id
}

output "phone_id" {
  value = elevenlabs_convai_phone_number.lifeos_mike.id
}
