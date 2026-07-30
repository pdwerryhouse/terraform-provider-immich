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
	_ datasource.DataSource              = &ActivityDataSource{}
	_ datasource.DataSourceWithConfigure = &ActivityDataSource{}
)

func NewActivitiesDataSource() datasource.DataSource {
	return &ActivityDataSource{}
}

type ActivityDataSource struct {
	client *immichclient.Client
}

type ActivityDataSourceModel struct {
	AlbumId    types.String    `tfsdk:"album_id"`
	Activities []ActivityModel `tfsdk:"activities"`
}

type ActivityModel struct {
	ID        types.String `tfsdk:"id"`
	AssetId   types.String `tfsdk:"asset_id"`
	Comment   types.String `tfsdk:"comment"`
	CreatedAt types.String `tfsdk:"created_at"`
	Type      types.String `tfsdk:"type"`
}

func (d *ActivityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_activities"
}

func (d *ActivityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ActivityDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	Activity, err := d.client.GetActivities(state.AlbumId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Activity",
			err.Error(),
		)
		return
	}

	for _, activity := range Activity {
		activityState := ActivityModel{
			ID:        types.StringValue(activity.Id),
			AssetId:   types.StringValue(activity.AssetId),
			CreatedAt: types.StringValue(activity.CreatedAt),
			Comment:   types.StringValue(activity.Comment),
			Type:      types.StringValue(activity.Type),
		}

		state.Activities = append(state.Activities, activityState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *ActivityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider develops.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ActivityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"album_id": schema.StringAttribute{
				Required:    true,
				Description: "Asset Id"},
			"activities": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"asset_id": schema.StringAttribute{
							Computed:    true,
							Description: "Asset Id"},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Creation date"},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Activity ID"},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "Comment text"},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Reaction Type"},
					},
				},
			},
		},
	}
}
