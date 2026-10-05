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
	"fmt"
	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &viewFolderUniquePathsDataSource{}
	_ datasource.DataSourceWithConfigure = &viewFolderUniquePathsDataSource{}
)

func NewViewFolderUniquePathsDataSource() datasource.DataSource {
	return &viewFolderUniquePathsDataSource{}
}

type viewFolderUniquePathsDataSource struct {
	client *immichclient.Client
}

type viewFolderUniquePathsDataSourceModel struct {
	UniquePaths types.List `tfsdk:"paths"`
}

func (d *viewFolderUniquePathsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view_folders_unique_paths"
}

func (d *viewFolderUniquePathsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state viewFolderUniquePathsDataSourceModel

	viewFolderUniquePaths, err := d.client.GetUniqueOriginalPaths()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
	}

	paths, diags := types.ListValueFrom(ctx, types.StringType, viewFolderUniquePaths)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.UniquePaths = paths

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *viewFolderUniquePathsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *viewFolderUniquePathsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"paths": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}
