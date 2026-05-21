variable "elevenlabs_api_key" {
  description = "ElevenLabs API key (env: ELEVENLABS_API_KEY). Must start with 'sk_'."
  type        = string
  sensitive   = true
  validation {
    condition     = length(var.elevenlabs_api_key) > 0 && can(regex("^sk_", var.elevenlabs_api_key))
    error_message = "elevenlabs_api_key must not be empty and must start with 'sk_'."
  }
}

variable "webhook_name" {
  description = "Display name for the workspace webhook"
  type        = string
  default     = "my-webhook"
}

variable "webhook_url" {
  description = "HTTPS endpoint that receives ElevenLabs webhook events"
  type        = string
  validation {
    condition     = can(regex("^https://", var.webhook_url))
    error_message = "webhook_url must be an https:// URL."
  }
}

variable "phone_number" {
  description = "E.164 phone number to import (e.g. +15551234567)"
  type        = string
  validation {
    condition     = can(regex("^\\+[0-9]{7,15}$", var.phone_number))
    error_message = "phone_number must be in E.164 format (e.g. +15551234567)."
  }
}

variable "phone_label" {
  description = "Human-readable label for the phone number resource"
  type        = string
  default     = "my-phone"
}
