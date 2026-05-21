variable "elevenlabs_api_key" {
  description = "ElevenLabs API key (env: ELEVENLABS_API_KEY). Must start with 'sk_'."
  type        = string
  sensitive   = true
  validation {
    condition     = length(var.elevenlabs_api_key) > 0 && can(regex("^sk_", var.elevenlabs_api_key))
    error_message = "elevenlabs_api_key must not be empty and must start with 'sk_'."
  }
}

variable "agent_name" {
  description = "Name for the cloned ConvAI agent (must be unique within your workspace)"
  type        = string
  default     = "my-restaurant-agent"
}

variable "webhook_post_call_id" {
  description = "Workspace webhook ID to fire after each call (leave empty to disable)"
  type        = string
  default     = ""
}
