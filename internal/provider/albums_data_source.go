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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"
)

var (
	_ datasource.DataSource              = &albumsDataSource{}
	_ datasource.DataSourceWithConfigure = &albumsDataSource{}
)

func NewAlbumsDataSource() datasource.DataSource {
	return &albumsDataSource{}
}

type albumsDataSource struct {
	client *immichclient.Client
}

type albumsDataSourceModel struct {
	Album []albumsModel `tfsdk:"albums"`
}

type albumsModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"album_name"`
	AlbumThumbnailAssetId types.String `tfsdk:"album_thumbnail_asset_id"`
	Description           types.String `tfsdk:"description"`
	Shared                types.Bool   `tfsdk:"shared"`
	HasSharedLink         types.Bool   `tfsdk:"has_shared_link"`
	Order                 types.String `tfsdk:"order"`
	IsActivityEnabled     types.Bool   `tfsdk:"is_activity_enabled"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
	StartDate             types.String `tfsdk:"start_date"`
	EndDate               types.String `tfsdk:"end_date"`
	OwnerId               types.String `tfsdk:"owner_id"`
}

func (d *albumsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_albums"
}

func (d *albumsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state albumsDataSourceModel

	albums, err := d.client.GetAllAlbums(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Album",
			err.Error(),
		)
		return
	}

	for _, album := range albums {
		albumState := albumsModel{
			ID:                    types.StringValue(album.Id),
			Name:                  types.StringValue(album.AlbumName),
			AlbumThumbnailAssetId: types.StringValue(album.AlbumThumbnailAssetId),
			Description:           types.StringValue(album.Description),
			Shared:                types.BoolValue(album.Shared),
			HasSharedLink:         types.BoolValue(album.HasSharedLink),
			IsActivityEnabled:     types.BoolValue(album.IsActivityEnabled),
			CreatedAt:             types.StringValue(album.CreatedAt),
			UpdatedAt:             types.StringValue(album.UpdatedAt),
			StartDate:             types.StringValue(album.StartDate),
			EndDate:               types.StringValue(album.EndDate),
		}

		state.Album = append(state.Album, albumState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *albumsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *albumsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of albums.",
		Attributes: map[string]schema.Attribute{
			"albums": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Album Id.",
							Computed:    true,
						},
						"album_name": schema.StringAttribute{
							Description: "The name of the album.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The album's description.",
							Computed:    true,
						},
						"shared": schema.BoolAttribute{
							Description: "If true, this is a shared album.",
							Computed:    true,
						},
						"order": schema.StringAttribute{
							Description: "The album order: asc or desc.",
							Computed:    true,
						},
						"album_thumbnail_asset_id": schema.StringAttribute{
							Description: "The id of the album's thumbnail asset.",
							Computed:    true,
						},
						"has_shared_link": schema.BoolAttribute{
							Description: "If true, this album has a shared link.",
							Computed:    true,
						},
						"is_activity_enabled": schema.BoolAttribute{
							Description: "If true, comments and likes can be left on the album.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Time at which the album was created.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Time at which the album was last updated.",
							Computed:    true,
						},
						"start_date": schema.StringAttribute{
							Description: "Date and time of earliest photo in the album.",
							Computed:    true,
						},
						"end_date": schema.StringAttribute{
							Description: "Date and time of latest photo in the album.",
							Computed:    true,
						},
						"owner_id": schema.StringAttribute{
							Description: "Id of the album's owner.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
