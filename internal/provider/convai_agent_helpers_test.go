package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringListToSlice_NilList(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	result := stringListToSlice(ctx, types.ListNull(types.StringType), &diags)
	if result != nil {
		t.Fatalf("expected nil for null list, got %v", result)
	}
}

func TestStringListToSlice_EmptyList(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	l, d := types.ListValueFrom(ctx, types.StringType, []string{})
	diags.Append(d...)
	result := stringListToSlice(ctx, l, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}
}

func TestStringListToSlice_NonEmpty(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	l, d := types.ListValueFrom(ctx, types.StringType, []string{"a", "b", "c"})
	diags.Append(d...)
	result := stringListToSlice(ctx, l, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(result) != 3 || result[0] != "a" || result[1] != "b" || result[2] != "c" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestSliceToStringList_Nil(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	l := sliceToStringList(ctx, nil, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	// nil input → empty list
	if l.IsNull() || l.IsUnknown() {
		t.Fatalf("expected non-null non-unknown list for nil input")
	}
	elems := l.Elements()
	if len(elems) != 0 {
		t.Fatalf("expected empty list elements, got %d", len(elems))
	}
}

func TestSliceToStringList_RoundTrip(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	input := []string{"x", "y"}
	l := sliceToStringList(ctx, input, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	result := stringListToSlice(ctx, l, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics on round-trip: %v", diags)
	}
	if len(result) != 2 || result[0] != "x" || result[1] != "y" {
		t.Fatalf("round-trip failed: %v", result)
	}
}

func TestStringMapToTF_NilMap(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	m := stringMapToTF(ctx, nil, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if m.IsNull() || m.IsUnknown() {
		t.Fatalf("expected empty map, not null/unknown")
	}
	if len(m.Elements()) != 0 {
		t.Fatalf("expected empty map elements")
	}
}

func TestStringMapToTF_RoundTrip(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	input := map[string]string{"k1": "v1", "k2": "v2"}
	tfm := stringMapToTF(ctx, input, &diags)
	if diags.HasError() {
		t.Fatalf("stringMapToTF diagnostics: %v", diags)
	}
	result := tfMapToStringMap(ctx, tfm, &diags)
	if diags.HasError() {
		t.Fatalf("tfMapToStringMap diagnostics: %v", diags)
	}
	if len(result) != 2 || result["k1"] != "v1" || result["k2"] != "v2" {
		t.Fatalf("round-trip failed: %v", result)
	}
}

func TestTFMapToStringMap_NullMap(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	result := tfMapToStringMap(ctx, types.MapNull(types.StringType), &diags)
	if result != nil {
		t.Fatalf("expected nil for null map, got %v", result)
	}
}

// TestBackupLLMConfigNilGuard verifies that an empty PromptBackupLLMPreference
// produces nil BackupLLMConfig in the request, not {"preference":""}.
func TestBackupLLMConfigNilGuard(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	data := minimalAgentModel(ctx, &diags)
	if diags.HasError() {
		t.Fatalf("setup diagnostics: %v", diags)
	}
	data.PromptBackupLLMPreference = types.StringValue("")

	req := dataToCreateRequest(ctx, &data, &diags)
	if diags.HasError() {
		t.Fatalf("dataToCreateRequest diagnostics: %v", diags)
	}
	if req.ConversationConfig.Agent.Prompt.BackupLLMConfig != nil {
		t.Fatalf("expected nil BackupLLMConfig for empty preference, got %+v",
			req.ConversationConfig.Agent.Prompt.BackupLLMConfig)
	}
}

func TestBackupLLMConfigSetWhenNonEmpty(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics

	data := minimalAgentModel(ctx, &diags)
	if diags.HasError() {
		t.Fatalf("setup diagnostics: %v", diags)
	}
	data.PromptBackupLLMPreference = types.StringValue("auto")

	req := dataToCreateRequest(ctx, &data, &diags)
	if diags.HasError() {
		t.Fatalf("dataToCreateRequest diagnostics: %v", diags)
	}
	if req.ConversationConfig.Agent.Prompt.BackupLLMConfig == nil {
		t.Fatal("expected non-nil BackupLLMConfig for non-empty preference")
	}
	if req.ConversationConfig.Agent.Prompt.BackupLLMConfig.Preference != "auto" {
		t.Fatalf("expected preference 'auto', got '%s'",
			req.ConversationConfig.Agent.Prompt.BackupLLMConfig.Preference)
	}
}

// minimalAgentModel returns a ConvAIAgentModel with required list/map fields
// initialized so dataToCreateRequest does not panic.
func minimalAgentModel(ctx context.Context, diags *diag.Diagnostics) ConvAIAgentModel {
	emptyStrList, d := types.ListValueFrom(ctx, types.StringType, []string{})
	diags.Append(d...)
	emptyStrMap, d2 := types.MapValueFrom(ctx, types.StringType, map[string]string{})
	diags.Append(d2...)
	emptyDC, d3 := types.MapValue(
		types.ObjectType{AttrTypes: dataCollectionAttrTypes()},
		map[string]attr.Value{},
	)
	diags.Append(d3...)

	return ConvAIAgentModel{
		Name: types.StringValue("test"),
		Tags: emptyStrList,

		ASRKeywords:                 emptyStrList,
		ConversationClientEvents:    emptyStrList,
		ConversationMonitoringEvents: emptyStrList,
		AgentDynamicVars:            emptyStrMap,
		WebhookEvents:               emptyStrList,
		DataCollection:              emptyDC,
	}
}
