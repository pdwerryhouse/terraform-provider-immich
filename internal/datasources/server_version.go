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
	_ datasource.DataSource              = &serverVersionDataSource{}
	_ datasource.DataSourceWithConfigure = &serverVersionDataSource{}
)

func NewServerVersionDataSource() datasource.DataSource {
	return &serverVersionDataSource{}
}

type serverVersionDataSource struct {
	client *immichclient.Client
}

type serverVersionDataSourceModel struct {
	Major      types.Int64 `tfsdk:"major"`
	Minor      types.Int64 `tfsdk:"minor"`
	Patch      types.Int64 `tfsdk:"patch"`
	Prerelease types.Int64 `tfsdk:"prerelease"`
}

func (d *serverVersionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_version"
}

func (d *serverVersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverVersionDataSourceModel

	serverVersion, err := d.client.GetServerVersion()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Server Version",
			err.Error(),
		)
		return
	}

	state.Major = types.Int64Value(serverVersion.Major)
	state.Minor = types.Int64Value(serverVersion.Minor)
	state.Patch = types.Int64Value(serverVersion.Patch)
	state.Prerelease = types.Int64Value(serverVersion.Prerelease)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serverVersionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Version Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *serverVersionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the Immich server's current version.",
		Attributes: map[string]schema.Attribute{

			"major": schema.Int64Attribute{
				Description: "Major version number.",
				Computed:    true,
			},
			"minor": schema.Int64Attribute{
				Description: "Minor version number.",
				Computed:    true,
			},
			"patch": schema.Int64Attribute{
				Description: "Patch version number.",
				Computed:    true,
			},
			"prerelease": schema.Int64Attribute{
				Description: "Pre-release version number.",
				Computed:    true,
			},
		},
	}
}
