package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = (*mantaProvider)(nil)
var _ provider.ProviderWithFunctions = (*mantaProvider)(nil)

type mantaProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

type mantaProvider struct {
	endpoint string
}

func New() provider.Provider {
	return &mantaProvider{}
}

func (p *mantaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "manta"
}

func (p *mantaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *mantaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config mantaProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Endpoint.IsNull() {
		p.endpoint = config.Endpoint.ValueString()
	}
}

func (p *mantaProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *mantaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *mantaProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewDeepMergeFunction,
		NewIsPalindromeFunction,
		NewMaskFunction,
		NewSemverCompareFunction,
		NewTruncateFunction,
	}
}
