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

package datasources

import (
	"context"
	"fmt"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &serverConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &serverConfigDataSource{}
)

func NewServerConfigDataSource() datasource.DataSource {
	return &serverConfigDataSource{}
}

type serverConfigDataSource struct {
	client *immichclient.Client
}

type serverConfigDataSourceModel struct {
	ExternalDomain   types.String `tfsdk:"external_domain"`
	IsInitialized    types.Bool   `tfsdk:"is_initialized"`
	IsOnboarded      types.Bool   `tfsdk:"is_onboarded"`
	LoginPageMessage types.String `tfsdk:"login_page_message"`
	MaintenanceMode  types.Bool   `tfsdk:"maintenance_mode"`
	MapDarkStyleUrl  types.String `tfsdk:"map_dark_style_url"`
	MapLightStyleUrl types.String `tfsdk:"map_light_style_url"`
	MinFaces         types.Int64  `tfsdk:"min_faces"`
	OauthButtonText  types.String `tfsdk:"oauth_button_text"`
	PublicUsers      types.Bool   `tfsdk:"public_users"`
	TrashDays        types.Int64  `tfsdk:"trash_days"`
	UserDeleteDelay  types.Int64  `tfsdk:"user_delete_delay"`
}

func (d *serverConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_config"
}

func (d *serverConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverConfigDataSourceModel

	serverConfig, err := d.client.GetServerConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
	}

	state.ExternalDomain = types.StringValue(serverConfig.ExternalDomain)
	state.IsInitialized = types.BoolValue(serverConfig.IsInitialized)
	state.IsOnboarded = types.BoolValue(serverConfig.IsOnboarded)
	state.LoginPageMessage = types.StringValue(serverConfig.LoginPageMessage)
	state.MaintenanceMode = types.BoolValue(serverConfig.MaintenanceMode)
	state.MapDarkStyleUrl = types.StringValue(serverConfig.MapDarkStyleUrl)
	state.MapLightStyleUrl = types.StringValue(serverConfig.MapLightStyleUrl)
	state.MinFaces = types.Int64Value(serverConfig.MinFaces)
	state.OauthButtonText = types.StringValue(serverConfig.OauthButtonText)
	state.PublicUsers = types.BoolValue(serverConfig.PublicUsers)
	state.TrashDays = types.Int64Value(serverConfig.TrashDays)
	state.UserDeleteDelay = types.Int64Value(serverConfig.UserDeleteDelay)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serverConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *serverConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"external_domain": schema.StringAttribute{
				Computed: true,
			},
			"is_initialized": schema.BoolAttribute{
				Computed: true,
			},
			"is_onboarded": schema.BoolAttribute{
				Computed: true,
			},
			"login_page_message": schema.StringAttribute{
				Computed: true,
			},
			"maintenance_mode": schema.BoolAttribute{
				Computed: true,
			},
			"map_dark_style_url": schema.StringAttribute{
				Computed: true,
			},
			"map_light_style_url": schema.StringAttribute{
				Computed: true,
			},
			"min_faces": schema.Int64Attribute{
				Computed: true,
			},
			"oauth_button_text": schema.StringAttribute{
				Computed: true,
			},
			"public_users": schema.BoolAttribute{
				Computed: true,
			},
			"trash_days": schema.Int64Attribute{
				Computed: true,
			},
			"user_delete_delay": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}
