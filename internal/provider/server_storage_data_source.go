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
	_ datasource.DataSource              = &serverStorageDataSource{}
	_ datasource.DataSourceWithConfigure = &serverStorageDataSource{}
)

func NewServerStorageDataSource() datasource.DataSource {
	return &serverStorageDataSource{}
}

type serverStorageDataSource struct {
	client *immichclient.Client
}

type serverStorageDataSourceModel struct {
	DiskAvailable       types.String  `tfsdk:"disk_available"`
	DiskAvailableRaw    types.Int64   `tfsdk:"disk_available_raw"`
	DiskSize            types.String  `tfsdk:"disk_size"`
	DiskSizeRaw         types.Int64   `tfsdk:"disk_size_raw"`
	DiskUsagePercentage types.Float64 `tfsdk:"disk_usage_percentage"`
	DiskUse             types.String  `tfsdk:"disk_use"`
	DiskUseRaw          types.Int64   `tfsdk:"disk_use_raw"`
}

func (d *serverStorageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_storage"
}

func (d *serverStorageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverStorageDataSourceModel

	serverStorage, err := d.client.GetStorage()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
	}

	state.DiskAvailable = types.StringValue(serverStorage.DiskAvailable)
	state.DiskAvailableRaw = types.Int64Value(serverStorage.DiskAvailableRaw)
	state.DiskSize = types.StringValue(serverStorage.DiskSize)
	state.DiskSizeRaw = types.Int64Value(serverStorage.DiskSizeRaw)
	state.DiskUsagePercentage = types.Float64Value(serverStorage.DiskUsagePercentage)
	state.DiskUse = types.StringValue(serverStorage.DiskUse)
	state.DiskUseRaw = types.Int64Value(serverStorage.DiskUseRaw)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serverStorageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Storageure Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider develops.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *serverStorageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the current storage utilization information from the Immich server.",
		Attributes: map[string]schema.Attribute{

			"disk_available": schema.StringAttribute{
				Description: "Available disk space (human readable format).",
				Computed:    true,
			},
			"disk_available_raw": schema.Int64Attribute{
				Description: "Available disk space in bytes.",
				Computed:    true,
			},
			"disk_size": schema.StringAttribute{
				Description: "Total disk size (human readable format).",
				Computed:    true,
			},
			"disk_size_raw": schema.Int64Attribute{
				Description: "Total disk size in bytes.",
				Computed:    true,
			},
			"disk_usage_percentage": schema.Float64Attribute{
				Description: "Disk usage percentage (0-100)",
				Computed:    true,
			},
			"disk_use": schema.StringAttribute{
				Description: "Used disk space (human readable format).",
				Computed:    true,
			},
			"disk_use_raw": schema.Int64Attribute{
				Description: "Used disk space in bytes.",
				Computed:    true,
			},
		},
	}
}
