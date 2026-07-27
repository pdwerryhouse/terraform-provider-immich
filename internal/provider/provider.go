// Copyright (C) 2026 Paul Dwerryhouse <paul@dwerryhouse.com.au>
//
// This file is part of terraform-provider-immich.
//
// terraform-provider-immich is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// terraform-provider-immich is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with terraform-provider-immich.  If not, see <https://www.gnu.org/licenses/>.

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
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"apikey"`
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
		Description: "Terraform provider for Immich",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The Immich API Endpoint. May also be provided via IMMICH_ENDPOINT.",
				Optional:    true,
			},
			"apikey": schema.StringAttribute{
				Description: "The Immich API Key. May also be provided via IMMICH_API_KEY.",
				Optional:    true,
				Sensitive:   true,
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
	if config.Endpoint.IsUnknown() {
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

	endpoint := os.Getenv("IMMICH_ENDPOINT")
	apikey := os.Getenv("IMMICH_API_KEY")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apikey = config.ApiKey.ValueString()
	}

	// XXX improve these
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Immich endpoint",
			"The Immich endpoint is unknown 2",
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

	c, err := client.NewClient(&endpoint, &apikey)

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
		NewActivitiesDataSource,
		NewAlbumsDataSource,
		NewApiKeysDataSource,
		NewFacesDataSource,
		NewLibrariesDataSource,
		NewUsersDataSource,
		NewPartnerDataSource,
		NewPeopleDataSource,
		NewServerAboutDataSource,
		NewServerConfigDataSource,
		NewServerFeaturesDataSource,
		NewServerStorageDataSource,
		NewTagsDataSource,
	}
}

func (p *immichProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewActivityResource,
		NewAlbumResource,
		//NewAlbumActivityResource,
		//NewAlbumOrderResource,
		NewApiKeyResource,
		NewFaceResource,
		NewLibraryResource,
		NewPartnerResource,
		NewPersonResource,
		NewUserResource,
		NewTagResource,
	}
}
