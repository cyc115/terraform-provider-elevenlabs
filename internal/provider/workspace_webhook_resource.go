package provider

import (
	"context"

	"github.com/cyc115/terraform-provider-elevenlabs/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewWorkspaceWebhookResource() resource.Resource { return &WorkspaceWebhookResource{} }

type WorkspaceWebhookResource struct{ client *Client }

type WorkspaceWebhookModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	WebhookURL                 types.String `tfsdk:"webhook_url"`
	AuthType                   types.String `tfsdk:"auth_type"`
	RetryEnabled               types.Bool   `tfsdk:"retry_enabled"`
	IsDisabled                 types.Bool   `tfsdk:"is_disabled"`
	IsAutoDisabled             types.Bool   `tfsdk:"is_auto_disabled"`
	CreatedAtUnix              types.Int64  `tfsdk:"created_at_unix"`
	Secret                     types.String `tfsdk:"secret"`
}

func (r *WorkspaceWebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_webhook"
}

func (r *WorkspaceWebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an ElevenLabs workspace webhook. Webhooks receive event callbacks (e.g. post-call transcripts) from ElevenLabs.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, MarkdownDescription: "ElevenLabs webhook ID (assigned on create)."},
			"name":            schema.StringAttribute{Required: true, MarkdownDescription: "Display name for the webhook."},
			"webhook_url":     schema.StringAttribute{Required: true},
			"auth_type":       schema.StringAttribute{Optional: true, Computed: true},
			"retry_enabled":   schema.BoolAttribute{Optional: true, Computed: true},
			"is_disabled":     schema.BoolAttribute{Optional: true, Computed: true},
			"is_auto_disabled": schema.BoolAttribute{Computed: true},
			"created_at_unix": schema.Int64Attribute{Computed: true},
			// Secret is write-once: set on Create response, never overwritten on Read/Import.
			"secret": schema.StringAttribute{Computed: true, Sensitive: true},
		},
	}
}

func (r *WorkspaceWebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *WorkspaceWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkspaceWebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := &models.CreateWebhookRequest{
		Name:         data.Name.ValueString(),
		WebhookURL:   data.WebhookURL.ValueString(),
		AuthType:     data.AuthType.ValueString(),
		RetryEnabled: data.RetryEnabled.ValueBool(),
	}
	wh, err := r.client.CreateWebhook(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("create webhook", err.Error())
		return
	}
	webhookToModel(wh, &data, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkspaceWebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	wh, err := r.client.GetWebhook(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read webhook", err.Error())
		return
	}
	// Preserve secret — Read never returns it; keep what's in state.
	webhookToModel(wh, &data, false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkspaceWebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	var state WorkspaceWebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = state.ID
	name := data.Name.ValueString()
	url := data.WebhookURL.ValueString()
	authType := data.AuthType.ValueString()
	retryEnabled := data.RetryEnabled.ValueBool()
	body := &models.UpdateWebhookRequest{
		Name:         &name,
		WebhookURL:   &url,
		AuthType:     &authType,
		RetryEnabled: &retryEnabled,
	}
	wh, err := r.client.UpdateWebhook(ctx, data.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("update webhook", err.Error())
		return
	}
	data.Secret = state.Secret // preserve secret
	webhookToModel(wh, &data, false)
	data.Secret = state.Secret
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkspaceWebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteWebhook(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete webhook", err.Error())
	}
}

func (r *WorkspaceWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	wh, err := r.client.GetWebhook(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("import webhook", err.Error())
		return
	}
	var data WorkspaceWebhookModel
	webhookToModel(wh, &data, false)
	// Secret is unknown after import — acceptable per spec AC-4.
	data.Secret = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func webhookToModel(wh *models.WebhookResponse, data *WorkspaceWebhookModel, setSecret bool) {
	data.ID = types.StringValue(wh.WebhookID)
	data.Name = types.StringValue(wh.Name)
	data.WebhookURL = types.StringValue(wh.WebhookURL)
	data.AuthType = types.StringValue(wh.AuthType)
	data.RetryEnabled = types.BoolValue(wh.RetryEnabled)
	data.IsDisabled = types.BoolValue(wh.IsDisabled)
	data.IsAutoDisabled = types.BoolValue(wh.IsAutoDisabled)
	data.CreatedAtUnix = types.Int64Value(wh.CreatedAtUnix)
	if setSecret {
		data.Secret = types.StringValue(wh.Secret)
	}
}
