package provider

// PATCH semantics probe result (Task 2.1):
// Sending null for a nested optional field (e.g. conversation.background_music: null) WIPES the
// field in the API response — null is replace-on-write, not preserve. Strategy: always send the
// full ConversationConfig and PlatformSettings objects on every PATCH, never rely on omitempty
// for nested objects. Partial PATCH is NOT safe.

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewConvAIAgentResource() resource.Resource { return &ConvAIAgentResource{} }

type ConvAIAgentResource struct{ client *Client }

type ConvAIAgentModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Tags types.List   `tfsdk:"tags"`

	// ASR
	ASRQuality             types.String `tfsdk:"asr_quality"`
	ASRProvider            types.String `tfsdk:"asr_provider"`
	ASRUserInputAudioFormat types.String `tfsdk:"asr_user_input_audio_format"`
	ASRKeywords            types.List   `tfsdk:"asr_keywords"`

	// Turn
	TurnTimeout              types.Float64 `tfsdk:"turn_timeout"`
	TurnSilenceEndCallTimeout types.Float64 `tfsdk:"turn_silence_end_call_timeout"`
	TurnMode                 types.String  `tfsdk:"turn_mode"`
	TurnEagerness            types.String  `tfsdk:"turn_eagerness"`
	TurnModel                types.String  `tfsdk:"turn_model"`
	TurnSpeculativeTurn      types.Bool    `tfsdk:"turn_speculative_turn"`
	TurnSoftTimeoutSeconds   types.Float64 `tfsdk:"turn_soft_timeout_seconds"`
	TurnSoftTimeoutMessage   types.String  `tfsdk:"turn_soft_timeout_message"`
	TurnSoftTimeoutUseLLM    types.Bool    `tfsdk:"turn_soft_timeout_use_llm"`

	// TTS
	TTSModelID                  types.String  `tfsdk:"tts_model_id"`
	TTSVoiceID                  types.String  `tfsdk:"tts_voice_id"`
	TTSExpressiveMode           types.Bool    `tfsdk:"tts_expressive_mode"`
	TTSOptimizeStreamingLatency types.Int64   `tfsdk:"tts_optimize_streaming_latency"`
	TTSStability                types.Float64 `tfsdk:"tts_stability"`
	TTSSpeed                    types.Float64 `tfsdk:"tts_speed"`
	TTSSimilarityBoost          types.Float64 `tfsdk:"tts_similarity_boost"`
	TTSAgentOutputAudioFormat   types.String  `tfsdk:"tts_agent_output_audio_format"`
	TTSTextNormalisationType    types.String  `tfsdk:"tts_text_normalisation_type"`

	// Conversation
	ConversationMaxDurationSeconds types.Int64  `tfsdk:"conversation_max_duration_seconds"`
	ConversationClientEvents       types.List   `tfsdk:"conversation_client_events"`
	ConversationMonitoringEvents   types.List   `tfsdk:"conversation_monitoring_events"`
	ConversationSourceAttribution  types.Bool   `tfsdk:"conversation_source_attribution"`
	ConversationBGMusicSourceType  types.String `tfsdk:"conversation_bg_music_source_type"`
	ConversationBGMusicSourceID    types.String `tfsdk:"conversation_bg_music_source_id"`
	ConversationBGMusicVolume      types.Float64 `tfsdk:"conversation_bg_music_volume"`
	ConversationBGMusicCrossfade   types.Bool   `tfsdk:"conversation_bg_music_crossfade"`
	ConversationFileInputEnabled   types.Bool   `tfsdk:"conversation_file_input_enabled"`
	ConversationFileInputMaxFiles  types.Int64  `tfsdk:"conversation_file_input_max_files"`

	// VAD
	VADBackgroundVoiceDetection types.Bool `tfsdk:"vad_background_voice_detection"`

	// Agent
	AgentFirstMessage types.String `tfsdk:"agent_first_message"`
	AgentLanguage     types.String `tfsdk:"agent_language"`
	AgentDynamicVars  types.Map    `tfsdk:"agent_dynamic_vars"`

	// Prompt
	PromptText                    types.String  `tfsdk:"prompt_text"`
	PromptLLM                     types.String  `tfsdk:"prompt_llm"`
	PromptReasoningEffort         types.String  `tfsdk:"prompt_reasoning_effort"`
	PromptTemperature             types.Float64 `tfsdk:"prompt_temperature"`
	PromptMaxTokens               types.Int64   `tfsdk:"prompt_max_tokens"`
	PromptTimezone                types.String  `tfsdk:"prompt_timezone"`
	PromptCascadeTimeout          types.Float64 `tfsdk:"prompt_cascade_timeout"`
	PromptEnableParallelToolCalls types.Bool    `tfsdk:"prompt_enable_parallel_tool_calls"`
	PromptIgnoreDefaultPersonality types.Bool   `tfsdk:"prompt_ignore_default_personality"`
	PromptBackupLLMPreference     types.String  `tfsdk:"prompt_backup_llm_preference"`
	PromptRAGEnabled              types.Bool    `tfsdk:"prompt_rag_enabled"`
	PromptRAGEmbeddingModel       types.String  `tfsdk:"prompt_rag_embedding_model"`
	PromptRAGMaxVectorDistance    types.Float64 `tfsdk:"prompt_rag_max_vector_distance"`
	PromptRAGMaxDocumentsLength   types.Int64   `tfsdk:"prompt_rag_max_documents_length"`
	PromptRAGMaxChunks            types.Int64   `tfsdk:"prompt_rag_max_chunks"`

	// Data collection (map of field name → field config)
	DataCollection types.Map `tfsdk:"data_collection"`

	// Webhook override
	WebhookPostCallID      types.String `tfsdk:"webhook_post_call_id"`
	WebhookEvents          types.List   `tfsdk:"webhook_events"`
	WebhookTranscriptFormat types.String `tfsdk:"webhook_transcript_format"`

	// Platform settings — call_limits
	PSAgentConcurrencyLimit types.Int64 `tfsdk:"ps_agent_concurrency_limit"`
	PSDailyLimit            types.Int64 `tfsdk:"ps_daily_limit"`
	PSBurstingEnabled       types.Bool  `tfsdk:"ps_bursting_enabled"`

	// Platform settings — privacy
	PSRecordVoice               types.Bool  `tfsdk:"ps_record_voice"`
	PSRetentionDays             types.Int64 `tfsdk:"ps_retention_days"`
	PSDeleteTranscriptAndPII    types.Bool  `tfsdk:"ps_delete_transcript_and_pii"`
	PSDeleteAudio               types.Bool  `tfsdk:"ps_delete_audio"`
	PSZeroRetentionMode         types.Bool  `tfsdk:"ps_zero_retention_mode"`
	PSHistoryRedactionEnabled   types.Bool  `tfsdk:"ps_history_redaction_enabled"`

	// Platform settings — guardrails
	PSGuardrailsFocusEnabled           types.Bool   `tfsdk:"ps_guardrails_focus_enabled"`
	PSGuardrailsPromptInjectionEnabled types.Bool   `tfsdk:"ps_guardrails_prompt_injection_enabled"`
	PSGuardrailsContentExecutionMode   types.String `tfsdk:"ps_guardrails_content_execution_mode"`
	PSGuardrailsTriggerAction          types.String `tfsdk:"ps_guardrails_trigger_action"`

	// Platform settings — auth
	PSAuthEnableAuth          types.Bool `tfsdk:"ps_auth_enable_auth"`
	PSAuthRequireOriginHeader types.Bool `tfsdk:"ps_auth_require_origin_header"`

	// Platform settings — top-level
	PSAnalysisLLM types.String `tfsdk:"ps_analysis_llm"`
	PSArchived    types.Bool   `tfsdk:"ps_archived"`

	// Coaching
	CoachingType types.String `tfsdk:"coaching_type"`
}

// DataCollectionFieldModel mirrors models.DataCollectionField for TF schema.
type DataCollectionFieldModel struct {
	Type             types.String `tfsdk:"type"`
	Description      types.String `tfsdk:"description"`
	Enum             types.List   `tfsdk:"enum"`
	IsSystemProvided types.Bool   `tfsdk:"is_system_provided"`
	DynamicVariable  types.String `tfsdk:"dynamic_variable"`
	ConstantValue    types.String `tfsdk:"constant_value"`
	LLM              types.String `tfsdk:"llm"`
}

func (r *ConvAIAgentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_agent"
}

func (r *ConvAIAgentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Required: true},
			"tags": schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},

			// ASR
			"asr_quality":               schema.StringAttribute{Optional: true, Computed: true},
			"asr_provider":              schema.StringAttribute{Optional: true, Computed: true},
			"asr_user_input_audio_format": schema.StringAttribute{Optional: true, Computed: true},
			"asr_keywords":              schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},

			// Turn
			"turn_timeout":                   schema.Float64Attribute{Optional: true, Computed: true},
			"turn_silence_end_call_timeout":   schema.Float64Attribute{Optional: true, Computed: true},
			"turn_mode":                       schema.StringAttribute{Optional: true, Computed: true},
			"turn_eagerness":                  schema.StringAttribute{Optional: true, Computed: true},
			"turn_model":                      schema.StringAttribute{Optional: true, Computed: true},
			"turn_speculative_turn":           schema.BoolAttribute{Optional: true, Computed: true},
			"turn_soft_timeout_seconds":       schema.Float64Attribute{Optional: true, Computed: true},
			"turn_soft_timeout_message":       schema.StringAttribute{Optional: true, Computed: true},
			"turn_soft_timeout_use_llm":       schema.BoolAttribute{Optional: true, Computed: true},

			// TTS
			"tts_model_id":                   schema.StringAttribute{Optional: true, Computed: true},
			"tts_voice_id":                   schema.StringAttribute{Optional: true, Computed: true},
			"tts_expressive_mode":            schema.BoolAttribute{Optional: true, Computed: true},
			"tts_optimize_streaming_latency": schema.Int64Attribute{Optional: true, Computed: true},
			"tts_stability":                  schema.Float64Attribute{Optional: true, Computed: true},
			"tts_speed":                      schema.Float64Attribute{Optional: true, Computed: true},
			"tts_similarity_boost":           schema.Float64Attribute{Optional: true, Computed: true},
			"tts_agent_output_audio_format":  schema.StringAttribute{Optional: true, Computed: true},
			"tts_text_normalisation_type":    schema.StringAttribute{Optional: true, Computed: true},

			// Conversation
			"conversation_max_duration_seconds": schema.Int64Attribute{Optional: true, Computed: true},
			"conversation_client_events":        schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"conversation_monitoring_events":    schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"conversation_source_attribution":   schema.BoolAttribute{Optional: true, Computed: true},
			"conversation_bg_music_source_type": schema.StringAttribute{Optional: true, Computed: true},
			"conversation_bg_music_source_id":   schema.StringAttribute{Optional: true, Computed: true},
			"conversation_bg_music_volume":      schema.Float64Attribute{Optional: true, Computed: true},
			"conversation_bg_music_crossfade":   schema.BoolAttribute{Optional: true, Computed: true},
			"conversation_file_input_enabled":   schema.BoolAttribute{Optional: true, Computed: true},
			"conversation_file_input_max_files": schema.Int64Attribute{Optional: true, Computed: true},

			// VAD
			"vad_background_voice_detection": schema.BoolAttribute{Optional: true, Computed: true},

			// Agent
			"agent_first_message": schema.StringAttribute{Optional: true, Computed: true},
			"agent_language":      schema.StringAttribute{Optional: true, Computed: true},
			"agent_dynamic_vars":  schema.MapAttribute{ElementType: types.StringType, Optional: true, Computed: true},

			// Prompt
			"prompt_text":                      schema.StringAttribute{Optional: true, Computed: true},
			"prompt_llm":                       schema.StringAttribute{Optional: true, Computed: true},
			"prompt_reasoning_effort":          schema.StringAttribute{Optional: true, Computed: true},
			"prompt_temperature":               schema.Float64Attribute{Optional: true, Computed: true},
			"prompt_max_tokens":                schema.Int64Attribute{Optional: true, Computed: true},
			"prompt_timezone":                  schema.StringAttribute{Optional: true, Computed: true},
			"prompt_cascade_timeout":           schema.Float64Attribute{Optional: true, Computed: true},
			"prompt_enable_parallel_tool_calls": schema.BoolAttribute{Optional: true, Computed: true},
			"prompt_ignore_default_personality": schema.BoolAttribute{Optional: true, Computed: true},
			"prompt_backup_llm_preference":     schema.StringAttribute{Optional: true, Computed: true},
			"prompt_rag_enabled":               schema.BoolAttribute{Optional: true, Computed: true},
			"prompt_rag_embedding_model":       schema.StringAttribute{Optional: true, Computed: true},
			"prompt_rag_max_vector_distance":   schema.Float64Attribute{Optional: true, Computed: true},
			"prompt_rag_max_documents_length":  schema.Int64Attribute{Optional: true, Computed: true},
			"prompt_rag_max_chunks":            schema.Int64Attribute{Optional: true, Computed: true},

			// Data collection
			"data_collection": schema.MapNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type":               schema.StringAttribute{Optional: true, Computed: true},
						"description":        schema.StringAttribute{Optional: true, Computed: true},
						"enum":               schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
						"is_system_provided": schema.BoolAttribute{Optional: true, Computed: true},
						"dynamic_variable":   schema.StringAttribute{Optional: true, Computed: true},
						"constant_value":     schema.StringAttribute{Optional: true, Computed: true},
						"llm":                schema.StringAttribute{Optional: true, Computed: true},
					},
				},
			},

			// Webhook override
			"webhook_post_call_id":      schema.StringAttribute{Optional: true, Computed: true},
			"webhook_events":            schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"webhook_transcript_format": schema.StringAttribute{Optional: true, Computed: true},

			// Platform settings — call_limits
			"ps_agent_concurrency_limit": schema.Int64Attribute{Optional: true, Computed: true},
			"ps_daily_limit":             schema.Int64Attribute{Optional: true, Computed: true},
			"ps_bursting_enabled":        schema.BoolAttribute{Optional: true, Computed: true},

			// Platform settings — privacy
			"ps_record_voice":                schema.BoolAttribute{Optional: true, Computed: true},
			"ps_retention_days":              schema.Int64Attribute{Optional: true, Computed: true},
			"ps_delete_transcript_and_pii":   schema.BoolAttribute{Optional: true, Computed: true},
			"ps_delete_audio":                schema.BoolAttribute{Optional: true, Computed: true},
			"ps_zero_retention_mode":         schema.BoolAttribute{Optional: true, Computed: true},
			"ps_history_redaction_enabled":   schema.BoolAttribute{Optional: true, Computed: true},

			// Platform settings — guardrails
			"ps_guardrails_focus_enabled":            schema.BoolAttribute{Optional: true, Computed: true},
			"ps_guardrails_prompt_injection_enabled": schema.BoolAttribute{Optional: true, Computed: true},
			"ps_guardrails_content_execution_mode":   schema.StringAttribute{Optional: true, Computed: true},
			"ps_guardrails_trigger_action":           schema.StringAttribute{Optional: true, Computed: true},

			// Platform settings — auth
			"ps_auth_enable_auth":           schema.BoolAttribute{Optional: true, Computed: true},
			"ps_auth_require_origin_header": schema.BoolAttribute{Optional: true, Computed: true},

			// Platform settings — top-level
			"ps_analysis_llm": schema.StringAttribute{Optional: true, Computed: true},
			"ps_archived":     schema.BoolAttribute{Optional: true, Computed: true},

			// Coaching (read-only: API sets this; not configurable on create/update)
			"coaching_type": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *ConvAIAgentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *ConvAIAgentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIAgentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := dataToCreateRequest(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := r.client.CreateAgent(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("create agent", err.Error())
		return
	}
	// Create response only contains agent_id; GET for the full state.
	agent, err := r.client.GetAgent(ctx, created.AgentID)
	if err != nil {
		resp.Diagnostics.AddError("read agent after create", err.Error())
		return
	}
	agentResponseToModel(ctx, agent, &data, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIAgentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIAgentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	agent, err := r.client.GetAgent(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read agent", err.Error())
		return
	}
	agentResponseToModel(ctx, agent, &data, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIAgentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIAgentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	var state ConvAIAgentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = state.ID
	body := dataToCreateRequest(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	agent, err := r.client.UpdateAgent(ctx, data.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("update agent", err.Error())
		return
	}
	agentResponseToModel(ctx, agent, &data, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIAgentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIAgentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteAgent(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete agent", err.Error())
	}
}

func (r *ConvAIAgentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	agent, err := r.client.GetAgent(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("import agent", err.Error())
		return
	}
	var data ConvAIAgentModel
	agentResponseToModel(ctx, agent, &data, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
