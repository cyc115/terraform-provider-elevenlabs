provider "elevenlabs" {
  api_key = var.elevenlabs_api_key
}

resource "elevenlabs_workspace_webhook" "main" {
  name          = var.webhook_name
  webhook_url   = var.webhook_url
  auth_type     = "hmac"
  retry_enabled = true
}

resource "elevenlabs_convai_phone_number" "main" {
  phone_number       = var.phone_number
  telephony_provider = "twilio"
  label              = var.phone_label
}

output "webhook_id" {
  value = elevenlabs_workspace_webhook.main.id
}

output "phone_id" {
  value = elevenlabs_convai_phone_number.main.id
}
