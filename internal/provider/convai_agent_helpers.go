package provider

import (
	"context"

	"github.com/cyc115/terraform-provider-elevenlabs/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dataCollectionAttrTypes returns the attribute types for a DataCollectionField map element.
func dataCollectionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":               types.StringType,
		"description":        types.StringType,
		"enum":               types.ListType{ElemType: types.StringType},
		"is_system_provided": types.BoolType,
		"dynamic_variable":   types.StringType,
		"constant_value":     types.StringType,
		"llm":                types.StringType,
	}
}

// stringListToSlice converts a types.List of strings to []string.
func stringListToSlice(ctx context.Context, l types.List, diags *diag.Diagnostics) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	var out []string
	diags.Append(l.ElementsAs(ctx, &out, false)...)
	return out
}

// sliceToStringList converts []string to types.List.
func sliceToStringList(ctx context.Context, s []string, diags *diag.Diagnostics) types.List {
	if s == nil {
		v, d := types.ListValueFrom(ctx, types.StringType, []string{})
		diags.Append(d...)
		return v
	}
	v, d := types.ListValueFrom(ctx, types.StringType, s)
	diags.Append(d...)
	return v
}

// stringMapToTF converts map[string]string to types.Map.
func stringMapToTF(ctx context.Context, m map[string]string, diags *diag.Diagnostics) types.Map {
	if m == nil {
		v, d := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		diags.Append(d...)
		return v
	}
	v, d := types.MapValueFrom(ctx, types.StringType, m)
	diags.Append(d...)
	return v
}

// tfMapToStringMap converts types.Map to map[string]string.
func tfMapToStringMap(ctx context.Context, m types.Map, diags *diag.Diagnostics) map[string]string {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	out := map[string]string{}
	diags.Append(m.ElementsAs(ctx, &out, false)...)
	return out
}

// dataToCreateRequest maps flat TF model → models.CreateAgentRequest (full object, no omitempty dodge).
func dataToCreateRequest(ctx context.Context, data *ConvAIAgentModel, diags *diag.Diagnostics) *models.CreateAgentRequest {
	req := &models.CreateAgentRequest{
		Name: data.Name.ValueString(),
		Tags: stringListToSlice(ctx, data.Tags, diags),
		ConversationConfig: &models.ConversationConfig{
			ASR: &models.ASRConfig{
				Quality:              data.ASRQuality.ValueString(),
				Provider:             data.ASRProvider.ValueString(),
				UserInputAudioFormat: data.ASRUserInputAudioFormat.ValueString(),
				Keywords:             stringListToSlice(ctx, data.ASRKeywords, diags),
			},
			Turn: &models.TurnConfig{
				TurnTimeout:           data.TurnTimeout.ValueFloat64(),
				SilenceEndCallTimeout: data.TurnSilenceEndCallTimeout.ValueFloat64(),
				Mode:                  data.TurnMode.ValueString(),
				TurnEagerness:         data.TurnEagerness.ValueString(),
				TurnModel:             data.TurnModel.ValueString(),
				SpeculativeTurn:       data.TurnSpeculativeTurn.ValueBool(),
				InterruptionIgnoreTerms: []string{},
				SoftTimeoutConfig: &models.SoftTimeoutConfig{
					TimeoutSeconds:         data.TurnSoftTimeoutSeconds.ValueFloat64(),
					Message:                data.TurnSoftTimeoutMessage.ValueString(),
					UseLLMGeneratedMessage: data.TurnSoftTimeoutUseLLM.ValueBool(),
				},
			},
			TTS: &models.TTSConfig{
				ModelID:                  data.TTSModelID.ValueString(),
				VoiceID:                  data.TTSVoiceID.ValueString(),
				ExpressiveMode:           data.TTSExpressiveMode.ValueBool(),
				OptimizeStreamingLatency: int(data.TTSOptimizeStreamingLatency.ValueInt64()),
				Stability:                data.TTSStability.ValueFloat64(),
				Speed:                    data.TTSSpeed.ValueFloat64(),
				SimilarityBoost:          data.TTSSimilarityBoost.ValueFloat64(),
				AgentOutputAudioFormat:   data.TTSAgentOutputAudioFormat.ValueString(),
				TextNormalisationType:    data.TTSTextNormalisationType.ValueString(),
				SupportedVoices:          []any{},
				SuggestedAudioTags:       []models.AudioTag{},
				PronunciationDictionaryLocators: []any{},
			},
			Conversation: &models.ConversationParams{
				MaxDurationSeconds: int(data.ConversationMaxDurationSeconds.ValueInt64()),
				ClientEvents:       stringListToSlice(ctx, data.ConversationClientEvents, diags),
				MonitoringEvents:   stringListToSlice(ctx, data.ConversationMonitoringEvents, diags),
				SourceAttribution:  data.ConversationSourceAttribution.ValueBool(),
				FileInput: &models.FileInput{
					Enabled:                 data.ConversationFileInputEnabled.ValueBool(),
					MaxFilesPerConversation: int(data.ConversationFileInputMaxFiles.ValueInt64()),
				},
				BackgroundMusic: &models.BackgroundMusic{
					SourceType:    data.ConversationBGMusicSourceType.ValueString(),
					SourceID:      data.ConversationBGMusicSourceID.ValueString(),
					Volume:        data.ConversationBGMusicVolume.ValueFloat64(),
					CrossfadeLoop: data.ConversationBGMusicCrossfade.ValueBool(),
				},
			},
			VAD: &models.VADConfig{
				BackgroundVoiceDetection: data.VADBackgroundVoiceDetection.ValueBool(),
			},
			Agent: &models.AgentConfig{
				FirstMessage: data.AgentFirstMessage.ValueString(),
				Language:     data.AgentLanguage.ValueString(),
				DynamicVariables: &models.DynamicVariables{
					DynamicVariablePlaceholders: tfMapToStringMap(ctx, data.AgentDynamicVars, diags),
				},
				Prompt: &models.PromptConfig{
					Prompt:                   data.PromptText.ValueString(),
					LLM:                      data.PromptLLM.ValueString(),
					ReasoningEffort:          data.PromptReasoningEffort.ValueString(),
					Temperature:              data.PromptTemperature.ValueFloat64(),
					MaxTokens:                int(data.PromptMaxTokens.ValueInt64()),
					Timezone:                 data.PromptTimezone.ValueString(),
					CascadeTimeoutSeconds:    data.PromptCascadeTimeout.ValueFloat64(),
					EnableParallelToolCalls:  data.PromptEnableParallelToolCalls.ValueBool(),
					IgnoreDefaultPersonality: data.PromptIgnoreDefaultPersonality.ValueBool(),
					BackupLLMConfig: func() *models.BackupLLMConfig {
						if p := data.PromptBackupLLMPreference.ValueString(); p != "" {
							return &models.BackupLLMConfig{Preference: p}
						}
						return nil
					}(),
					RAG: &models.RAGConfig{
						Enabled:                    data.PromptRAGEnabled.ValueBool(),
						EmbeddingModel:             data.PromptRAGEmbeddingModel.ValueString(),
						MaxVectorDistance:          data.PromptRAGMaxVectorDistance.ValueFloat64(),
						MaxDocumentsLength:         int(data.PromptRAGMaxDocumentsLength.ValueInt64()),
						MaxRetrievedRAGChunksCount: int(data.PromptRAGMaxChunks.ValueInt64()),
					},
					ToolIDs:            []string{},
					KnowledgeBase:      []any{},
					MCPServerIDs:       []string{},
					NativeMCPServerIDs: []string{},
				},
			},
		},
		PlatformSettings: toPlatformSettings(ctx, data, diags),
	}

	// Data collection
	if !data.DataCollection.IsNull() && !data.DataCollection.IsUnknown() {
		dcMap := map[string]*models.DataCollectionField{}
		dcElements := map[string]DataCollectionFieldModel{}
		diags.Append(data.DataCollection.ElementsAs(ctx, &dcElements, false)...)
		for k, v := range dcElements {
			enumSlice := stringListToSlice(ctx, v.Enum, diags)
			var llmPtr *string
			if !v.LLM.IsNull() && !v.LLM.IsUnknown() && v.LLM.ValueString() != "" {
				s := v.LLM.ValueString()
				llmPtr = &s
			}
			dcMap[k] = &models.DataCollectionField{
				Type:             v.Type.ValueString(),
				Description:      v.Description.ValueString(),
				Enum:             enumSlice,
				IsSystemProvided: v.IsSystemProvided.ValueBool(),
				DynamicVariable:  v.DynamicVariable.ValueString(),
				ConstantValue:    v.ConstantValue.ValueString(),
				LLM:              llmPtr,
			}
		}
		req.PlatformSettings.DataCollection = dcMap
	}

	// Webhook override
	if !data.WebhookPostCallID.IsNull() && !data.WebhookPostCallID.IsUnknown() && data.WebhookPostCallID.ValueString() != "" {
		req.PlatformSettings.WorkspaceOverrides = &models.WorkspaceOverrides{
			Webhooks: &models.WebhookOverride{
				PostCallWebhookID: data.WebhookPostCallID.ValueString(),
				Events:            stringListToSlice(ctx, data.WebhookEvents, diags),
				TranscriptFormat:  data.WebhookTranscriptFormat.ValueString(),
			},
		}
	}

	return req
}

func toPlatformSettings(ctx context.Context, data *ConvAIAgentModel, diags *diag.Diagnostics) *models.PlatformSettings {
	ps := &models.PlatformSettings{
		CallLimits: &models.CallLimits{
			AgentConcurrencyLimit: int(data.PSAgentConcurrencyLimit.ValueInt64()),
			DailyLimit:            int(data.PSDailyLimit.ValueInt64()),
			BurstingEnabled:       data.PSBurstingEnabled.ValueBool(),
		},
		Privacy: &models.Privacy{
			RecordVoice:            data.PSRecordVoice.ValueBool(),
			RetentionDays:          int(data.PSRetentionDays.ValueInt64()),
			DeleteTranscriptAndPII: data.PSDeleteTranscriptAndPII.ValueBool(),
			DeleteAudio:            data.PSDeleteAudio.ValueBool(),
			ZeroRetentionMode:      data.PSZeroRetentionMode.ValueBool(),
		},
		Guardrails: &models.Guardrails{
			Version:         "1",
			Focus:           &models.GuardrailSwitch{IsEnabled: data.PSGuardrailsFocusEnabled.ValueBool()},
			PromptInjection: &models.GuardrailSwitch{IsEnabled: data.PSGuardrailsPromptInjectionEnabled.ValueBool()},
			Content: &models.ContentGuardrail{
				ExecutionMode: data.PSGuardrailsContentExecutionMode.ValueString(),
				TriggerAction: &models.TriggerAction{Type: data.PSGuardrailsTriggerAction.ValueString()},
			},
		},
		Auth: &models.Auth{
			EnableAuth:          data.PSAuthEnableAuth.ValueBool(),
			RequireOriginHeader: data.PSAuthRequireOriginHeader.ValueBool(),
			Allowlist:           []any{},
		},
		AnalysisLLM: data.PSAnalysisLLM.ValueString(),
		Archived:    data.PSArchived.ValueBool(),
	}
	if data.PSHistoryRedactionEnabled.ValueBool() {
		ps.Privacy.ConversationHistoryRedaction = &models.ConversationHistoryRedaction{
			Enabled:  true,
			Entities: []string{},
		}
	}
	return ps
}

// agentResponseToModel maps models.AgentResponse → flat TF ConvAIAgentModel.
func agentResponseToModel(ctx context.Context, agent *models.AgentResponse, data *ConvAIAgentModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(agent.AgentID)
	data.Name = types.StringValue(agent.Name)
	data.Tags = sliceToStringList(ctx, agent.Tags, diags)

	cc := agent.ConversationConfig
	if cc == nil {
		return
	}

	if cc.ASR != nil {
		data.ASRQuality = types.StringValue(cc.ASR.Quality)
		data.ASRProvider = types.StringValue(cc.ASR.Provider)
		data.ASRUserInputAudioFormat = types.StringValue(cc.ASR.UserInputAudioFormat)
		data.ASRKeywords = sliceToStringList(ctx, cc.ASR.Keywords, diags)
	}

	if cc.Turn != nil {
		data.TurnTimeout = types.Float64Value(cc.Turn.TurnTimeout)
		data.TurnSilenceEndCallTimeout = types.Float64Value(cc.Turn.SilenceEndCallTimeout)
		data.TurnMode = types.StringValue(cc.Turn.Mode)
		data.TurnEagerness = types.StringValue(cc.Turn.TurnEagerness)
		data.TurnModel = types.StringValue(cc.Turn.TurnModel)
		data.TurnSpeculativeTurn = types.BoolValue(cc.Turn.SpeculativeTurn)
		if cc.Turn.SoftTimeoutConfig != nil {
			data.TurnSoftTimeoutSeconds = types.Float64Value(cc.Turn.SoftTimeoutConfig.TimeoutSeconds)
			data.TurnSoftTimeoutMessage = types.StringValue(cc.Turn.SoftTimeoutConfig.Message)
			data.TurnSoftTimeoutUseLLM = types.BoolValue(cc.Turn.SoftTimeoutConfig.UseLLMGeneratedMessage)
		}
	}

	if cc.TTS != nil {
		data.TTSModelID = types.StringValue(cc.TTS.ModelID)
		data.TTSVoiceID = types.StringValue(cc.TTS.VoiceID)
		data.TTSExpressiveMode = types.BoolValue(cc.TTS.ExpressiveMode)
		data.TTSOptimizeStreamingLatency = types.Int64Value(int64(cc.TTS.OptimizeStreamingLatency))
		data.TTSStability = types.Float64Value(cc.TTS.Stability)
		data.TTSSpeed = types.Float64Value(cc.TTS.Speed)
		data.TTSSimilarityBoost = types.Float64Value(cc.TTS.SimilarityBoost)
		data.TTSAgentOutputAudioFormat = types.StringValue(cc.TTS.AgentOutputAudioFormat)
		data.TTSTextNormalisationType = types.StringValue(cc.TTS.TextNormalisationType)
	}

	if cc.Conversation != nil {
		data.ConversationMaxDurationSeconds = types.Int64Value(int64(cc.Conversation.MaxDurationSeconds))
		data.ConversationClientEvents = sliceToStringList(ctx, cc.Conversation.ClientEvents, diags)
		data.ConversationMonitoringEvents = sliceToStringList(ctx, cc.Conversation.MonitoringEvents, diags)
		data.ConversationSourceAttribution = types.BoolValue(cc.Conversation.SourceAttribution)
		if cc.Conversation.FileInput != nil {
			data.ConversationFileInputEnabled = types.BoolValue(cc.Conversation.FileInput.Enabled)
			data.ConversationFileInputMaxFiles = types.Int64Value(int64(cc.Conversation.FileInput.MaxFilesPerConversation))
		}
		if cc.Conversation.BackgroundMusic != nil {
			data.ConversationBGMusicSourceType = types.StringValue(cc.Conversation.BackgroundMusic.SourceType)
			data.ConversationBGMusicSourceID = types.StringValue(cc.Conversation.BackgroundMusic.SourceID)
			data.ConversationBGMusicVolume = types.Float64Value(cc.Conversation.BackgroundMusic.Volume)
			data.ConversationBGMusicCrossfade = types.BoolValue(cc.Conversation.BackgroundMusic.CrossfadeLoop)
		} else {
			data.ConversationBGMusicSourceType = types.StringValue("")
			data.ConversationBGMusicSourceID = types.StringValue("")
			data.ConversationBGMusicVolume = types.Float64Value(0)
			data.ConversationBGMusicCrossfade = types.BoolValue(false)
		}
	}

	if cc.VAD != nil {
		data.VADBackgroundVoiceDetection = types.BoolValue(cc.VAD.BackgroundVoiceDetection)
	}

	if cc.Agent != nil {
		data.AgentFirstMessage = types.StringValue(cc.Agent.FirstMessage)
		data.AgentLanguage = types.StringValue(cc.Agent.Language)
		if cc.Agent.DynamicVariables != nil {
			data.AgentDynamicVars = stringMapToTF(ctx, cc.Agent.DynamicVariables.DynamicVariablePlaceholders, diags)
		} else {
			v, d := types.MapValueFrom(ctx, types.StringType, map[string]string{})
			diags.Append(d...)
			data.AgentDynamicVars = v
		}

		if cc.Agent.Prompt != nil {
			p := cc.Agent.Prompt
			data.PromptText = types.StringValue(p.Prompt)
			data.PromptLLM = types.StringValue(p.LLM)
			data.PromptReasoningEffort = types.StringValue(p.ReasoningEffort)
			data.PromptTemperature = types.Float64Value(p.Temperature)
			data.PromptMaxTokens = types.Int64Value(int64(p.MaxTokens))
			data.PromptTimezone = types.StringValue(p.Timezone)
			data.PromptCascadeTimeout = types.Float64Value(p.CascadeTimeoutSeconds)
			data.PromptEnableParallelToolCalls = types.BoolValue(p.EnableParallelToolCalls)
			data.PromptIgnoreDefaultPersonality = types.BoolValue(p.IgnoreDefaultPersonality)
			if p.BackupLLMConfig != nil {
				data.PromptBackupLLMPreference = types.StringValue(p.BackupLLMConfig.Preference)
			}
			if p.RAG != nil {
				data.PromptRAGEnabled = types.BoolValue(p.RAG.Enabled)
				data.PromptRAGEmbeddingModel = types.StringValue(p.RAG.EmbeddingModel)
				data.PromptRAGMaxVectorDistance = types.Float64Value(p.RAG.MaxVectorDistance)
				data.PromptRAGMaxDocumentsLength = types.Int64Value(int64(p.RAG.MaxDocumentsLength))
				data.PromptRAGMaxChunks = types.Int64Value(int64(p.RAG.MaxRetrievedRAGChunksCount))
			}
		}
	}

	// Platform settings
	ps := agent.PlatformSettings
	if ps == nil {
		return
	}

	fromPlatformSettings(ctx, ps, data, diags)
}

func fromPlatformSettings(ctx context.Context, ps *models.PlatformSettings, data *ConvAIAgentModel, diags *diag.Diagnostics) {
	data.PSAnalysisLLM = types.StringValue(ps.AnalysisLLM)
	data.PSArchived = types.BoolValue(ps.Archived)

	if ps.CallLimits != nil {
		data.PSAgentConcurrencyLimit = types.Int64Value(int64(ps.CallLimits.AgentConcurrencyLimit))
		data.PSDailyLimit = types.Int64Value(int64(ps.CallLimits.DailyLimit))
		data.PSBurstingEnabled = types.BoolValue(ps.CallLimits.BurstingEnabled)
	}

	if ps.Privacy != nil {
		data.PSRecordVoice = types.BoolValue(ps.Privacy.RecordVoice)
		data.PSRetentionDays = types.Int64Value(int64(ps.Privacy.RetentionDays))
		data.PSDeleteTranscriptAndPII = types.BoolValue(ps.Privacy.DeleteTranscriptAndPII)
		data.PSDeleteAudio = types.BoolValue(ps.Privacy.DeleteAudio)
		data.PSZeroRetentionMode = types.BoolValue(ps.Privacy.ZeroRetentionMode)
		if ps.Privacy.ConversationHistoryRedaction != nil {
			data.PSHistoryRedactionEnabled = types.BoolValue(ps.Privacy.ConversationHistoryRedaction.Enabled)
		} else {
			data.PSHistoryRedactionEnabled = types.BoolValue(false)
		}
	}

	if ps.Guardrails != nil {
		if ps.Guardrails.Focus != nil {
			data.PSGuardrailsFocusEnabled = types.BoolValue(ps.Guardrails.Focus.IsEnabled)
		}
		if ps.Guardrails.PromptInjection != nil {
			data.PSGuardrailsPromptInjectionEnabled = types.BoolValue(ps.Guardrails.PromptInjection.IsEnabled)
		}
		if ps.Guardrails.Content != nil {
			data.PSGuardrailsContentExecutionMode = types.StringValue(ps.Guardrails.Content.ExecutionMode)
			if ps.Guardrails.Content.TriggerAction != nil {
				data.PSGuardrailsTriggerAction = types.StringValue(ps.Guardrails.Content.TriggerAction.Type)
			}
		}
	}

	if ps.Auth != nil {
		data.PSAuthEnableAuth = types.BoolValue(ps.Auth.EnableAuth)
		data.PSAuthRequireOriginHeader = types.BoolValue(ps.Auth.RequireOriginHeader)
	}

	if ps.CoachingSettings != nil {
		data.CoachingType = types.StringValue(ps.CoachingSettings.Type)
	} else {
		data.CoachingType = types.StringValue("")
	}

	// Webhook override
	if ps.WorkspaceOverrides != nil && ps.WorkspaceOverrides.Webhooks != nil {
		wh := ps.WorkspaceOverrides.Webhooks
		data.WebhookPostCallID = types.StringValue(wh.PostCallWebhookID)
		data.WebhookEvents = sliceToStringList(ctx, wh.Events, diags)
		data.WebhookTranscriptFormat = types.StringValue(wh.TranscriptFormat)
	} else {
		data.WebhookPostCallID = types.StringValue("")
		data.WebhookEvents = sliceToStringList(ctx, []string{}, diags)
		data.WebhookTranscriptFormat = types.StringValue("")
	}

	// Data collection
	if len(ps.DataCollection) > 0 {
		attrTypes := dataCollectionAttrTypes()
		dcObjMap := map[string]attr.Value{}
		for k, f := range ps.DataCollection {
			if f == nil {
				continue
			}
			enumVal := sliceToStringList(ctx, f.Enum, diags)
			llmStr := ""
			if f.LLM != nil {
				llmStr = *f.LLM
			}
			objVal, d := types.ObjectValue(attrTypes, map[string]attr.Value{
				"type":               types.StringValue(f.Type),
				"description":        types.StringValue(f.Description),
				"enum":               enumVal,
				"is_system_provided": types.BoolValue(f.IsSystemProvided),
				"dynamic_variable":   types.StringValue(f.DynamicVariable),
				"constant_value":     types.StringValue(f.ConstantValue),
				"llm":                types.StringValue(llmStr),
			})
			diags.Append(d...)
			dcObjMap[k] = objVal
		}
		mapVal, d := types.MapValue(types.ObjectType{AttrTypes: attrTypes}, dcObjMap)
		diags.Append(d...)
		data.DataCollection = mapVal
	} else {
		mapVal, d := types.MapValue(types.ObjectType{AttrTypes: dataCollectionAttrTypes()}, map[string]attr.Value{})
		diags.Append(d...)
		data.DataCollection = mapVal
	}
}
