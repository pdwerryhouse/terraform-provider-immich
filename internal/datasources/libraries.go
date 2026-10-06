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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"
)

var (
	_ datasource.DataSource              = &LibraryDataSource{}
	_ datasource.DataSourceWithConfigure = &LibraryDataSource{}
)

func NewLibrariesDataSource() datasource.DataSource {
	return &LibraryDataSource{}
}

type LibraryDataSource struct {
	client *immichclient.Client
}

type LibraryDataSourceModel struct {
	Library []LibraryModel `tfsdk:"library"`
}

type LibraryModel struct {
	AssetCount        types.Int64  `tfsdk:"asset_count"`
	CreatedAt         types.String `tfsdk:"created_at"`
	ExclusionPatterns types.List   `tfsdk:"exclusion_patterns"`
	Id                types.String `tfsdk:"id"`
	ImportPaths       types.List   `tfsdk:"import_paths"`
	Name              types.String `tfsdk:"name"`
	OwnerId           types.String `tfsdk:"owner_id"`
	RefreshedAt       types.String `tfsdk:"refreshed_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

func (d *LibraryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_libraries"
}

func (d *LibraryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LibraryDataSourceModel

	Library, err := d.client.GetAllLibraries()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Library",
			err.Error(),
		)
		return
	}

	for _, library := range Library {
		ep, diags := types.ListValueFrom(ctx, types.StringType, library.ExclusionPatterns)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		ip, diags := types.ListValueFrom(ctx, types.StringType, library.ImportPaths)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		libraryState := LibraryModel{
			AssetCount:        types.Int64Value(library.AssetCount),
			CreatedAt:         types.StringValue(library.CreatedAt),
			ExclusionPatterns: ep,
			Id:                types.StringValue(library.Id),
			ImportPaths:       ip,
			Name:              types.StringValue(library.Name),
			OwnerId:           types.StringValue(library.OwnerId),
			RefreshedAt:       types.StringValue(library.RefreshedAt),
			UpdatedAt:         types.StringValue(library.UpdatedAt),
		}

		state.Library = append(state.Library, libraryState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *LibraryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *LibraryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"library": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"asset_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of assets"},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Creation date"},
						"exclusion_patterns": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Exclusion patterns",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Library ID"},
						"import_paths": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Import paths",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Library name"},
						"owner_id": schema.StringAttribute{
							Computed:    true,
							Description: "Owner user ID"},
						"refreshed_at": schema.StringAttribute{
							Computed:    true,
							Description: "Last refresh date"},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "Last update date"},
					},
				},
			},
		},
	}
}
