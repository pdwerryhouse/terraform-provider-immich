package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-immich/internal/client"
)

var (
	_ provider.Provider = &immichProvider{}
)

type immichProviderModel struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"apikey"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &immichProvider{
			version: version,
		}
	}
}

type immichProvider struct {
	version string
}

func (p *immichProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "immich"
	resp.Version = p.version
}

func (p *immichProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"apikey": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *immichProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config immichProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// XXX improve error messages
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Immich host",
			"The Immich host is unknown 1",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown Immich apikey",
			"The Immich apikey is unknown 1",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("IMMICH_HOST")
	apikey := os.Getenv("IMMICH_APIKEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apikey = config.ApiKey.ValueString()
	}

	// XXX improve these
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Immich host",
			"The Immich host is unknown 2",
		)
	}

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown Immich apikey",
			"The Immich apikey is unknown 2",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c, err := client.NewClient(&host, &apikey)

	// XXX improve this
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Immich API Client",
			"An Unexpected error occurred.",
		)
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *immichProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAlbumsDataSource,
	}
}

func (p *immichProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}
