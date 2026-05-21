package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New() provider.Provider { return &ElevenLabsProvider{} }

type ElevenLabsProvider struct{}

type ElevenLabsProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

func (p *ElevenLabsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "elevenlabs"
}

func (p *ElevenLabsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provider for the [ElevenLabs](https://elevenlabs.io) API. Manages ConvAI agents, workspace webhooks, and phone numbers.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "ElevenLabs API key. Can also be set via the `ELEVENLABS_API_KEY` environment variable.",
			},
		},
	}
}

func (p *ElevenLabsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ElevenLabsProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiKey := config.APIKey.ValueString()
	if apiKey == "" {
		apiKey = os.Getenv("ELEVENLABS_API_KEY")
	}
	c := NewClient(apiKey)
	resp.ResourceData = c
	resp.DataSourceData = c
}

func (p *ElevenLabsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewConvAIAgentResource,
		NewWorkspaceWebhookResource,
		NewConvAIPhoneNumberResource,
	}
}

func (p *ElevenLabsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
