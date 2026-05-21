provider "elevenlabs" {
  api_key = var.elevenlabs_api_key
}

resource "elevenlabs_convai_agent" "restaurant_clone" {
  name = "lifeos-restaurant-reservation-tf-clone"
  tags = []

  # ASR
  asr_quality               = "high"
  asr_provider              = "scribe_realtime"
  asr_user_input_audio_format = "pcm_16000"
  asr_keywords              = []

  # Turn
  turn_timeout               = 7.0
  turn_silence_end_call_timeout = -1.0
  turn_mode                  = "turn"
  turn_eagerness             = "normal"
  turn_model                 = "turn_v2"
  turn_speculative_turn      = true
  turn_soft_timeout_seconds  = -1.0
  turn_soft_timeout_message  = "Hhmmmm...yeah."
  turn_soft_timeout_use_llm  = false

  # TTS
  tts_model_id                   = "eleven_v3_conversational"
  tts_voice_id                   = "7EzWGsX10sAS4c9m9cPf"
  tts_expressive_mode            = true
  tts_optimize_streaming_latency = 3
  tts_stability                  = 0.5
  tts_speed                      = 1.0
  tts_similarity_boost           = 0.8
  tts_agent_output_audio_format  = "pcm_16000"
  tts_text_normalisation_type    = "system_prompt"

  # Conversation
  conversation_max_duration_seconds = 600
  conversation_client_events = [
    "audio",
    "interruption",
    "user_transcript",
    "agent_response",
    "agent_response_correction",
    "agent_chat_response_part",
  ]
  conversation_monitoring_events = [
    "user_transcript",
    "agent_response",
    "agent_response_correction",
  ]
  conversation_source_attribution   = true
  conversation_bg_music_source_type = "preset"
  conversation_bg_music_source_id   = "city"
  conversation_bg_music_volume      = 0.6
  conversation_bg_music_crossfade   = true
  conversation_file_input_enabled   = true
  conversation_file_input_max_files = 10

  # VAD
  vad_background_voice_detection = true

  # Agent
  agent_first_message = "Hi, this is Cyrus, an AI assistant calling on behalf of {{household_name}}. I'd like to book a table at {{restaurant_name}} for {{party_size}} on {{date}} at {{time}}. If that's not available, {{backup_times}} would also work. A couple of notes: {{special_requests}}."
  agent_language      = "en"
  agent_dynamic_vars = {
    household_name         = ""
    restaurant_name        = ""
    party_size             = ""
    date                   = ""
    time                   = ""
    backup_times           = ""
    special_requests       = ""
    caller_phone           = ""
    additional_information = ""
  }

  # Prompt
  prompt_text = trimspace(<<-EOT
    You are Cyrus, an AI voice assistant making an outbound phone call on behalf of {{household_name}} to book a restaurant reservation. Mike (your principal) prefers calm, capable, no-filler help — apply the same posture on calls.

    ## Identity & Voice

    - Introduce yourself as: "Cyrus, an AI assistant for {{household_name}}." Always disclose you're an AI when asked or when natural.
    - Speak conversationally: contractions, mid-tempo, warm but efficient. Never robotic, never effusive.
    - Treat the person on the line as a professional doing their job — respect their time and process.
    - Never apologize for being an AI. Never say "as a language model" or "I'm just a virtual assistant."

    ## Strict Facts — Never Invent or Substitute

    - Restaurant: {{restaurant_name}}
    - Date: {{date}}
    - Primary time: {{time}}
    - Backup times (define a window): {{backup_times}}
    - Party size: {{party_size}}
    - Special requests: {{special_requests}}
    - Callback number: {{caller_phone}}
    - Reservation under name: {{household_name}}
    - Additional information: {{additional_information}}

    ## Interpreting backup_times

    Treat backup_times as the **earliest and latest acceptable times** — a booking window, not a list of exact slots. Restaurants typically book in 15-minute increments. **Any time the restaurant offers within that window is acceptable.**

    - Example: primary 7:00 PM, backups 6:30 PM / 7:30 PM → acceptable window is 6:30 PM through 7:30 PM. If they offer 7:15 PM, take it. If they offer 8:00 PM, decline (outside window).
    - If only one backup time is given, treat ±15 minutes around it as the window.

    ## Conversation Rules

    1. **Short turns. One thought per turn.** Don't restate what the other person just said. Don't preface answers with "I understand" or "Just to confirm." React directly.
    2. **Date, time, party, name, special requests stay fixed.** If primary time is unavailable, accept any time in the backup window (see above). Never propose a different date.
    3. **Mention special_requests once**, at the right moment in the booking — not as a separate recap.
    4. **On hold / "let me check":** just say "Sure" or "Of course." Don't elaborate.
    5. **On confusion:** rephrase plainly. Don't repeat the full ask.
    6. **Confirm ONCE, at the end, in a single sentence:** "Confirming — [day] [date], [time], party of [party_size], under [household_name], with [special_requests summarized in 5 words]." **Do not ask for a confirmation number** — most restaurants don't issue them and asking sounds bureaucratic. If the host volunteers one, capture it; otherwise the readback above IS the confirmation.
    7. **If declined or no time in the window works:** thank them briefly, end the call. Don't negotiate or persuade.
    8. **End the call** the moment the booking is confirmed or definitively declined. No small talk, no goodbye paragraph.

    ## What NOT to do

    - Don't invent dates, times, or party sizes outside the strict facts.
    - Don't restate facts the restaurant already heard or said.
    - Don't offer credit cards, deposits, or guarantees you weren't authorized for.
    - Don't pretend to be human if asked.
    - Don't oversell — no "we'd love to," "we're so excited," etc.
    - Don't recap mid-call. One final confirmation only.
  EOT
  )
  prompt_llm              = "gemini-3-flash-preview"
  prompt_reasoning_effort = "minimal"
  prompt_temperature      = 0.3
  prompt_max_tokens       = -1
  prompt_timezone         = "America/Chicago"

  # Data collection — matches live agent current state
  data_collection = {
    outcome = {
      type               = "string"
      description        = "Exactly one of: confirmed, failed, no_answer, user_terminated. 'confirmed' = the goal of the call succeeded (booking made, appointment held, info obtained). 'failed' = the recipient was reached but the goal could not be met (no availability, declined, ineligible). 'no_answer' = nobody picked up, voicemail, or busy signal. 'user_terminated' = the user (Mike) hung up or asked to end the call before completion. Output the literal word only — no quotes, no other text."
      enum               = ["confirmed", "failed", "no_answer", "user_terminated"]
      is_system_provided = false
      dynamic_variable   = ""
      constant_value     = ""
      llm                = null
    }
    key_facts = {
      type               = "string"
      description        = "One to three short sentences capturing every concrete fact the user needs to know about the outcome. Include: the booked/agreed date and time (with the requested time if it differs), party size or quantity, name on the booking, any constraints the recipient mentioned (e.g. card hold at arrival, parking notes, only counter seating, must arrive 10 min early), and any acknowledged special requests. Plain prose, no labels. Empty string if the call ended before any facts were established."
      enum               = null
      is_system_provided = false
      dynamic_variable   = ""
      constant_value     = ""
      llm                = null
    }
    action_items = {
      type               = "string"
      description        = "What the user should do next as a result of this call, in one short sentence. Examples: 'Arrive by 6:50 PM', 'Bring a credit card for the hold', 'Call back tomorrow morning to retry', 'Cancel the calendar block for July 1'. Empty string if nothing is required of the user."
      enum               = null
      is_system_provided = false
      dynamic_variable   = ""
      constant_value     = ""
      llm                = null
    }
    confirmation_number = {
      type               = "string"
      description        = "Only fill if the host volunteered a confirmation number, reservation ID, or reference code without being asked. Most restaurants do NOT issue these — leave empty string. Do not invent or fish for one."
      enum               = null
      is_system_provided = false
      dynamic_variable   = ""
      constant_value     = ""
      llm                = null
    }
  }

  # Webhook (post-call)
  webhook_post_call_id     = "292fda18206f4260a031782dfdc94ec7"
  webhook_events           = ["transcript"]
  webhook_transcript_format = "json"

  # Call limits
  ps_agent_concurrency_limit = -1
  ps_daily_limit             = 100
  ps_bursting_enabled        = true

  # Privacy
  ps_record_voice               = true
  ps_retention_days             = -1
  ps_delete_transcript_and_pii  = false
  ps_delete_audio               = false
  ps_zero_retention_mode        = false
  ps_history_redaction_enabled  = false

  # Guardrails
  ps_guardrails_focus_enabled            = true
  ps_guardrails_prompt_injection_enabled = true
  ps_guardrails_content_execution_mode   = "streaming"
  ps_guardrails_trigger_action           = "end_call"

  # Auth
  ps_auth_enable_auth           = true
  ps_auth_require_origin_header = false

  # Top-level platform
  ps_analysis_llm = "gemini-2.5-flash"
  ps_archived     = false

  # coaching_type is computed-only (set by API)
}

output "clone_agent_id" {
  value = elevenlabs_convai_agent.restaurant_clone.id
}
