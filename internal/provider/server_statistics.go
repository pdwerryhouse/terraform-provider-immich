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
	_ datasource.DataSource              = &serverStatisticsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverStatisticsDataSource{}
)

func NewServerStatisticsDataSource() datasource.DataSource {
	return &serverStatisticsDataSource{}
}

type serverStatisticsDataSource struct {
	client *immichclient.Client
}

type usageByUserModel struct {
	Photos           types.Int64  `tfsdk:"photos"`
	QuotaSizeInBytes types.Int64  `tfsdk:"quota_size_in_bytes"`
	Usage            types.Int64  `tfsdk:"usage"`
	UsagePhotos      types.Int64  `tfsdk:"usage_photos"`
	UsageVideos      types.Int64  `tfsdk:"usage_videos"`
	UserId           types.String `tfsdk:"user_id"`
	UserName         types.String `tfsdk:"user_name"`
	Videos           types.Int64  `tfsdk:"videos"`
}

type serverStatisticsDataSourceModel struct {
	Photos      types.Int64        `tfsdk:"photos"`
	Usage       types.Int64        `tfsdk:"usage"`
	UsagePhotos types.Int64        `tfsdk:"usage_photos"`
	UsageByUser []usageByUserModel `tfsdk:"usage_by_user"`
	UsageVideos types.Int64        `tfsdk:"usage_videos"`
	Videos      types.Int64        `tfsdk:"videos"`
}

func (d *serverStatisticsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_statistics"
}

func (d *serverStatisticsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverStatisticsDataSourceModel

	serverStatistics, err := d.client.GetServerStatistics()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Server Media Types",
			err.Error(),
		)
		return
	}

	state.Photos = types.Int64Value(serverStatistics.Photos)
	state.Usage = types.Int64Value(serverStatistics.Usage)
	state.UsagePhotos = types.Int64Value(serverStatistics.UsagePhotos)
	state.UsageVideos = types.Int64Value(serverStatistics.UsageVideos)
	state.Videos = types.Int64Value(serverStatistics.Videos)

	var usages []usageByUserModel

	for _, user := range serverStatistics.UsageByUser {
		usage := usageByUserModel{
			Photos:           types.Int64Value(user.Photos),
			QuotaSizeInBytes: types.Int64Value(user.QuotaSizeInBytes),
			Usage:            types.Int64Value(user.Usage),
			UsagePhotos:      types.Int64Value(user.UsagePhotos),
			UsageVideos:      types.Int64Value(user.UsageVideos),
			UserId:           types.StringValue(user.UserId),
			UserName:         types.StringValue(user.UserName),
			Videos:           types.Int64Value(user.Videos),
		}
		usages = append(usages, usage)
	}

	state.UsageByUser = usages

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serverStatisticsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Statistics Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *serverStatisticsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve all media types supported by the Immich server.",
		Attributes: map[string]schema.Attribute{

			"photos": schema.Int64Attribute{
				Description: "Total number of photos.",
				Computed:    true,
			},
			"usage": schema.Int64Attribute{
				Description: "Total storage usage in bytes.",
				Computed:    true,
			},
			"usage_by_user": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"photos": schema.Int64Attribute{
							Description: "Number of photos.",
							Computed:    true,
						},
						"quota_size_in_bytes": schema.Int64Attribute{
							Description: "User quota size in bytes.",
							Computed:    true,
						},
						"usage": schema.Int64Attribute{
							Description: "Total storage usage in bytes.",
							Computed:    true,
						},
						"usage_photos": schema.Int64Attribute{
							Description: "Storage usage for photos in bytes.",
							Computed:    true,
						},
						"usage_videos": schema.Int64Attribute{
							Description: "Storage usage for videos in bytes.",
							Computed:    true,
						},
						"user_id": schema.StringAttribute{
							Description: "User Id.",
							Computed:    true,
						},
						"user_name": schema.StringAttribute{
							Description: "User name.",
							Computed:    true,
						},
						"videos": schema.Int64Attribute{
							Description: "Number of videos.",
							Computed:    true,
						},
					},
				},
			},
			"usage_photos": schema.Int64Attribute{
				Description: "Storage usage for photos in bytes.",
				Computed:    true,
			},
			"usage_videos": schema.Int64Attribute{
				Description: "Storage usage for videos in bytes.",
				Computed:    true,
			},
			"videos": schema.Int64Attribute{
				Description: "Total number of videos.",
				Computed:    true,
			},
		},
	}
}
