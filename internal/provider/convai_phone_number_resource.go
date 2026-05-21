package provider

import (
	"context"

	"github.com/cyc115/terraform-provider-elevenlabs/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewConvAIPhoneNumberResource() resource.Resource { return &ConvAIPhoneNumberResource{} }

type ConvAIPhoneNumberResource struct{ client *Client }

type ConvAIPhoneNumberModel struct {
	ID                types.String `tfsdk:"id"`
	PhoneNumber       types.String `tfsdk:"phone_number"`
	TelephonyProvider types.String `tfsdk:"telephony_provider"`
	Label             types.String `tfsdk:"label"`
	AssignedAgent     types.String `tfsdk:"assigned_agent"`
	SupportsInbound   types.Bool   `tfsdk:"supports_inbound"`
	SupportsOutbound  types.Bool   `tfsdk:"supports_outbound"`
	TwilioAccountSID  types.String `tfsdk:"twilio_account_sid"`
	TwilioAuthToken   types.String `tfsdk:"twilio_auth_token"`
}

func (r *ConvAIPhoneNumberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_phone_number"
}

func (r *ConvAIPhoneNumberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"phone_number": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"telephony_provider": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"label":             schema.StringAttribute{Optional: true, Computed: true},
			"assigned_agent":    schema.StringAttribute{Optional: true, Computed: true},
			"supports_inbound":  schema.BoolAttribute{Computed: true},
			"supports_outbound": schema.BoolAttribute{Computed: true},
			"twilio_account_sid": schema.StringAttribute{Optional: true, Sensitive: true},
			"twilio_auth_token":  schema.StringAttribute{Optional: true, Sensitive: true},
		},
	}
}

func (r *ConvAIPhoneNumberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*Client)
}

func (r *ConvAIPhoneNumberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIPhoneNumberModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := &models.ImportPhoneNumberRequest{
		PhoneNumber: data.PhoneNumber.ValueString(),
		Provider:    data.TelephonyProvider.ValueString(),
		Label:       data.Label.ValueString(),
	}
	if !data.TwilioAccountSID.IsNull() && !data.TwilioAccountSID.IsUnknown() {
		s := data.TwilioAccountSID.ValueString()
		body.TwilioAccountSID = &s
	}
	if !data.TwilioAuthToken.IsNull() && !data.TwilioAuthToken.IsUnknown() {
		s := data.TwilioAuthToken.ValueString()
		body.TwilioAuthToken = &s
	}
	phone, err := r.client.ImportPhone(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("create phone number", err.Error())
		return
	}
	phoneToModel(phone, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIPhoneNumberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIPhoneNumberModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	phone, err := r.client.GetPhone(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read phone number", err.Error())
		return
	}
	phoneToModel(phone, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIPhoneNumberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIPhoneNumberModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	var state ConvAIPhoneNumberModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = state.ID
	label := data.Label.ValueString()
	body := &models.UpdatePhoneNumberRequest{Label: &label}
	if !data.AssignedAgent.IsNull() && !data.AssignedAgent.IsUnknown() {
		s := data.AssignedAgent.ValueString()
		body.AssignedAgent = &s
	}
	phone, err := r.client.UpdatePhone(ctx, data.ID.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("update phone number", err.Error())
		return
	}
	phoneToModel(phone, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIPhoneNumberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIPhoneNumberModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeletePhone(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("delete phone number", err.Error())
	}
}

func (r *ConvAIPhoneNumberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	phone, err := r.client.GetPhone(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("import phone number", err.Error())
		return
	}
	var data ConvAIPhoneNumberModel
	phoneToModel(phone, &data)
	// Twilio credentials not available after import — leave null so plan matches config.
	data.TwilioAccountSID = types.StringNull()
	data.TwilioAuthToken = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func phoneToModel(phone *models.PhoneNumberResponse, data *ConvAIPhoneNumberModel) {
	data.ID = types.StringValue(phone.PhoneNumberID)
	data.PhoneNumber = types.StringValue(phone.PhoneNumber)
	data.TelephonyProvider = types.StringValue(phone.Provider)
	data.Label = types.StringValue(phone.Label)
	data.SupportsInbound = types.BoolValue(phone.SupportsInbound)
	data.SupportsOutbound = types.BoolValue(phone.SupportsOutbound)
	if phone.AssignedAgent != nil {
		data.AssignedAgent = types.StringValue(*phone.AssignedAgent)
	} else {
		data.AssignedAgent = types.StringValue("")
	}
}
