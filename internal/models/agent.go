package models

type CreateAgentRequest struct {
	Name               string              `json:"name"`
	ConversationConfig *ConversationConfig `json:"conversation_config,omitempty"`
	PlatformSettings   *PlatformSettings   `json:"platform_settings,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
}

type UpdateAgentRequest = CreateAgentRequest

type AgentResponse struct {
	AgentID            string              `json:"agent_id"`
	Name               string              `json:"name"`
	ConversationConfig *ConversationConfig `json:"conversation_config"`
	PlatformSettings   *PlatformSettings   `json:"platform_settings"`
	Tags               []string            `json:"tags"`
}

type ConversationConfig struct {
	ASR             *ASRConfig          `json:"asr,omitempty"`
	Turn            *TurnConfig         `json:"turn,omitempty"`
	TTS             *TTSConfig          `json:"tts,omitempty"`
	Conversation    *ConversationParams `json:"conversation,omitempty"`
	VAD             *VADConfig          `json:"vad,omitempty"`
	Agent           *AgentConfig        `json:"agent,omitempty"`
	LanguagePresets map[string]any      `json:"language_presets,omitempty"`
}

type ASRConfig struct {
	Quality              string   `json:"quality,omitempty"`
	Provider             string   `json:"provider,omitempty"`
	UserInputAudioFormat string   `json:"user_input_audio_format,omitempty"`
	Keywords             []string `json:"keywords"`
}

type TurnConfig struct {
	TurnTimeout               float64            `json:"turn_timeout,omitempty"`
	InitialWaitTime           *float64           `json:"initial_wait_time"`
	SilenceEndCallTimeout     float64            `json:"silence_end_call_timeout,omitempty"`
	Mode                      string             `json:"mode,omitempty"`
	TurnEagerness             string             `json:"turn_eagerness,omitempty"`
	SpellingPatience          string             `json:"spelling_patience,omitempty"`
	SpeculativeTurn           bool               `json:"speculative_turn,omitempty"`
	RetranscribeOnTurnTimeout bool               `json:"retranscribe_on_turn_timeout,omitempty"`
	TurnModel                 string             `json:"turn_model,omitempty"`
	InterruptionIgnoreTerms   []string           `json:"interruption_ignore_terms"`
	SoftTimeoutConfig         *SoftTimeoutConfig `json:"soft_timeout_config,omitempty"`
}

type SoftTimeoutConfig struct {
	TimeoutSeconds         float64 `json:"timeout_seconds"`
	Message                string  `json:"message"`
	UseLLMGeneratedMessage bool    `json:"use_llm_generated_message"`
}

type TTSConfig struct {
	ModelID                        string     `json:"model_id,omitempty"`
	VoiceID                        string     `json:"voice_id,omitempty"`
	SupportedVoices                []any      `json:"supported_voices"`
	ExpressiveMode                 bool       `json:"expressive_mode,omitempty"`
	SuggestedAudioTags             []AudioTag `json:"suggested_audio_tags"`
	AgentOutputAudioFormat         string     `json:"agent_output_audio_format,omitempty"`
	OptimizeStreamingLatency       int        `json:"optimize_streaming_latency,omitempty"`
	Stability                      float64    `json:"stability,omitempty"`
	Speed                          float64    `json:"speed,omitempty"`
	SimilarityBoost                float64    `json:"similarity_boost,omitempty"`
	TextNormalisationType          string     `json:"text_normalisation_type,omitempty"`
	PronunciationDictionaryLocators []any     `json:"pronunciation_dictionary_locators"`
}

type AudioTag struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

type ConversationParams struct {
	TextOnly            bool             `json:"text_only"`
	MaxDurationSeconds  int              `json:"max_duration_seconds,omitempty"`
	ClientEvents        []string         `json:"client_events"`
	FileInput           *FileInput       `json:"file_input,omitempty"`
	MonitoringEnabled   bool             `json:"monitoring_enabled"`
	MonitoringEvents    []string         `json:"monitoring_events"`
	BackgroundMusic     *BackgroundMusic `json:"background_music,omitempty"`
	SourceAttribution   bool             `json:"source_attribution"`
}

type FileInput struct {
	Enabled                 bool `json:"enabled"`
	MaxFilesPerConversation int  `json:"max_files_per_conversation,omitempty"`
}

type BackgroundMusic struct {
	SourceType    string  `json:"source_type,omitempty"`
	SourceID      string  `json:"source_id,omitempty"`
	Volume        float64 `json:"volume,omitempty"`
	CrossfadeLoop bool    `json:"crossfade_loop"`
}

type VADConfig struct {
	BackgroundVoiceDetection bool `json:"background_voice_detection"`
}

type AgentConfig struct {
	FirstMessage                     string            `json:"first_message,omitempty"`
	Language                         string            `json:"language,omitempty"`
	HinglishMode                     bool              `json:"hinglish_mode"`
	DynamicVariables                 *DynamicVariables `json:"dynamic_variables,omitempty"`
	DisableFirstMessageInterruptions bool              `json:"disable_first_message_interruptions"`
	MaxConversationDurationMessage   string            `json:"max_conversation_duration_message"`
	Prompt                           *PromptConfig     `json:"prompt,omitempty"`
}

type DynamicVariables struct {
	DynamicVariablePlaceholders map[string]string `json:"dynamic_variable_placeholders"`
}

type PromptConfig struct {
	Prompt                   string           `json:"prompt,omitempty"`
	LLM                      string           `json:"llm,omitempty"`
	ReasoningEffort          string           `json:"reasoning_effort,omitempty"`
	Temperature              float64          `json:"temperature,omitempty"`
	MaxTokens                int              `json:"max_tokens,omitempty"`
	ToolIDs                  []string         `json:"tool_ids"`
	BuiltInTools             map[string]any   `json:"built_in_tools,omitempty"`
	EnableParallelToolCalls  bool             `json:"enable_parallel_tool_calls"`
	KnowledgeBase            []any            `json:"knowledge_base"`
	MCPServerIDs             []string         `json:"mcp_server_ids"`
	NativeMCPServerIDs       []string         `json:"native_mcp_server_ids"`
	IgnoreDefaultPersonality bool             `json:"ignore_default_personality"`
	RAG                      *RAGConfig       `json:"rag,omitempty"`
	Timezone                 string           `json:"timezone,omitempty"`
	BackupLLMConfig          *BackupLLMConfig `json:"backup_llm_config,omitempty"`
	CascadeTimeoutSeconds    float64          `json:"cascade_timeout_seconds,omitempty"`
	Tools                    []map[string]any `json:"tools,omitempty"`
}

type RAGConfig struct {
	Enabled                   bool    `json:"enabled"`
	EmbeddingModel            string  `json:"embedding_model,omitempty"`
	OptionalRAGEnabled        bool    `json:"optional_rag_enabled"`
	MaxVectorDistance         float64 `json:"max_vector_distance,omitempty"`
	MaxDocumentsLength        int     `json:"max_documents_length,omitempty"`
	MaxRetrievedRAGChunksCount int    `json:"max_retrieved_rag_chunks_count,omitempty"`
}

type BackupLLMConfig struct {
	Preference string `json:"preference,omitempty"`
}

// --- PlatformSettings ---

type PlatformSettings struct {
	DataCollection     map[string]*DataCollectionField `json:"data_collection,omitempty"`
	WorkspaceOverrides *WorkspaceOverrides             `json:"workspace_overrides,omitempty"`
	CallLimits         *CallLimits                     `json:"call_limits,omitempty"`
	Privacy            *Privacy                        `json:"privacy,omitempty"`
	Guardrails         *Guardrails                     `json:"guardrails,omitempty"`
	Auth               *Auth                           `json:"auth,omitempty"`
	AnalysisLLM        string                          `json:"analysis_llm,omitempty"`
	SummaryLanguage    *string                         `json:"summary_language"`
	Archived           bool                            `json:"archived"`
	Overrides          *ClientOverrides                `json:"overrides,omitempty"`
	Evaluation         *Evaluation                     `json:"evaluation,omitempty"`
	Safety             *Safety                         `json:"safety,omitempty"`
	CoachingSettings   *CoachingSettings               `json:"coaching_settings,omitempty"`
}

type DataCollectionField struct {
	Type             string   `json:"type,omitempty"`
	Description      string   `json:"description,omitempty"`
	Enum             []string `json:"enum"`
	IsSystemProvided bool     `json:"is_system_provided"`
	DynamicVariable  string   `json:"dynamic_variable"`
	ConstantValue    string   `json:"constant_value"`
	LLM              *string  `json:"llm"`
}

type WorkspaceOverrides struct {
	Webhooks *WebhookOverride `json:"webhooks,omitempty"`
}

type WebhookOverride struct {
	PostCallWebhookID string   `json:"post_call_webhook_id,omitempty"`
	Events            []string `json:"events"`
	TranscriptFormat  string   `json:"transcript_format,omitempty"`
	SendAudio         *bool    `json:"send_audio"`
}

type CallLimits struct {
	AgentConcurrencyLimit int  `json:"agent_concurrency_limit"`
	DailyLimit            int  `json:"daily_limit,omitempty"`
	BurstingEnabled       bool `json:"bursting_enabled"`
}

type Privacy struct {
	RecordVoice                  bool                          `json:"record_voice"`
	RetentionDays                int                           `json:"retention_days"`
	DeleteTranscriptAndPII       bool                          `json:"delete_transcript_and_pii"`
	DeleteAudio                  bool                          `json:"delete_audio"`
	ApplyToExistingConversations bool                          `json:"apply_to_existing_conversations"`
	ZeroRetentionMode            bool                          `json:"zero_retention_mode"`
	ConversationHistoryRedaction *ConversationHistoryRedaction `json:"conversation_history_redaction,omitempty"`
}

type ConversationHistoryRedaction struct {
	Enabled  bool     `json:"enabled"`
	Entities []string `json:"entities"`
}

type Guardrails struct {
	Version         string           `json:"version,omitempty"`
	Focus           *GuardrailSwitch `json:"focus,omitempty"`
	PromptInjection *GuardrailSwitch `json:"prompt_injection,omitempty"`
	Content         *ContentGuardrail `json:"content,omitempty"`
}

type GuardrailSwitch struct {
	IsEnabled bool `json:"is_enabled"`
}

type ContentGuardrail struct {
	ExecutionMode string         `json:"execution_mode,omitempty"`
	Config        *ContentConfig `json:"config,omitempty"`
	TriggerAction *TriggerAction `json:"trigger_action,omitempty"`
}

type ContentConfig struct {
	Sexual              *ContentCategory `json:"sexual,omitempty"`
	Violence            *ContentCategory `json:"violence,omitempty"`
	Harassment          *ContentCategory `json:"harassment,omitempty"`
	SelfHarm            *ContentCategory `json:"self_harm,omitempty"`
	Profanity           *ContentCategory `json:"profanity,omitempty"`
	ReligionOrPolitics  *ContentCategory `json:"religion_or_politics,omitempty"`
	MedicalAndLegalInfo *ContentCategory `json:"medical_and_legal_information,omitempty"`
}

type ContentCategory struct {
	IsEnabled bool   `json:"is_enabled"`
	Threshold string `json:"threshold,omitempty"`
}

type TriggerAction struct {
	Type string `json:"type,omitempty"`
}

type Auth struct {
	EnableAuth          bool  `json:"enable_auth"`
	Allowlist           []any `json:"allowlist"`
	RequireOriginHeader bool  `json:"require_origin_header"`
}

type ClientOverrides struct {
	ConversationConfigOverride map[string]any `json:"conversation_config_override,omitempty"`
}

type Evaluation struct {
	Criteria []any `json:"criteria"`
}

type Safety struct {
	IsBlockedIVC          bool `json:"is_blocked_ivc"`
	IsBlockedNonIVC       bool `json:"is_blocked_non_ivc"`
	IgnoreSafetyEvaluation bool `json:"ignore_safety_evaluation"`
}

type CoachingSettings struct {
	Type         string `json:"type,omitempty"`
	MemoryBaseID string `json:"memory_base_id,omitempty"`
}
