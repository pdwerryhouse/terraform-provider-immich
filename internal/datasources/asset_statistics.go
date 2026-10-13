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
	"regexp"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &assetStatisticsDataSource{}
	_ datasource.DataSourceWithConfigure = &assetStatisticsDataSource{}
)

func NewAssetStatisticsDataSource() datasource.DataSource {
	return &assetStatisticsDataSource{}
}

type assetStatisticsDataSource struct {
	client *immichclient.Client
}

type assetStatisticsDataSourceModel struct {
	IsFavorite types.Bool   `tfsdk:"is_favorite"`
	IsTrashed  types.Bool   `tfsdk:"is_trashed"`
	Visibility types.String `tfsdk:"visibility"`
	Images     types.Int64  `tfsdk:"images"`
	Total      types.Int64  `tfsdk:"total"`
	Videos     types.Int64  `tfsdk:"videos"`
}

func (d *assetStatisticsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_statistics"
}

func (d *assetStatisticsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state assetStatisticsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var isFavorite, isTrashed *bool

	if !state.IsFavorite.IsNull() {
		v := state.IsFavorite.ValueBool()
		isFavorite = &v
	}

	if !state.IsTrashed.IsNull() {
		v := state.IsTrashed.ValueBool()
		isTrashed = &v
	}

	getAssetStatsParams := immichclient.GetAssetStatisticsParameters{
		IsFavorite: isFavorite,
		IsTrashed:  isTrashed,
		Visibility: (*immichclient.AssetVisibility)(state.Visibility.ValueStringPointer()),
	}

	assetStatistics, err := d.client.GetAssetStatistics(&getAssetStatsParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
	}

	state.Images = types.Int64Value(assetStatistics.Images)
	state.Total = types.Int64Value(assetStatistics.Total)
	state.Videos = types.Int64Value(assetStatistics.Videos)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *assetStatisticsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *assetStatisticsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get statistics about Immich assets",
		Attributes: map[string]schema.Attribute{
			"is_favorite": schema.BoolAttribute{
				Description: "If true, show only assets that have been marked as favourites.",
				Optional:    true,
			},
			"is_trashed": schema.BoolAttribute{
				Description: "If true, show only assets that have been put in the bin.",
				Optional:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Filter by visibility state. Options are 'archive', 'timeline', 'hidden' or 'locked'.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(archive|timeline|hidden|locked)?$`),
						"must be 'archive', 'timeline', 'hidden', 'locked' or empty",
					),
				},
			},
			"images": schema.Int64Attribute{
				Description: "Number of image assets.",
				Computed:    true,
			},
			"total": schema.Int64Attribute{
				Description: "Total number of assets.",
				Computed:    true,
			},
			"videos": schema.Int64Attribute{
				Description: "Number of video assets.",
				Computed:    true,
			},
		},
	}
}
