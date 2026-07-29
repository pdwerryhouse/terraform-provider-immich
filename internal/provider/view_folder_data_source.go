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
	"terraform-provider-immich/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &viewFolderPathDataSource{}
	_ datasource.DataSourceWithConfigure = &viewFolderPathDataSource{}
)

func NewViewFolderDataSource() datasource.DataSource {
	return &viewFolderPathDataSource{}
}

type viewFolderPathDataSource struct {
	client *client.Client
}

type viewFolderPathDataSourceModel struct {
	Path   types.String `tfsdk:"path"`
	Assets []assetModel `tfsdk:"assets"`
}

// XXX add all these and more

type assetModel struct {
	Checksum         types.String `tfsdk:"checksum"`
	CreatedAt        types.String `tfsdk:"created_at"`
	ID               types.String `tfsdk:"id"`
	DuplicateId      types.String `tfsdk:"duplicate_id"`
	Duration         types.Int64  `tfsdk:"duration"`
	FileCreatedAt    types.String `tfsdk:"file_created_at"`
	FileModifiedAt   types.String `tfsdk:"file_modified_at"`
	HasMetadata      types.Bool   `tfsdk:"has_metadata"`
	Height           types.Int64  `tfsdk:"height"`
	IsArchived       types.Bool   `tfsdk:"is_archived"`
	IsEdited         types.Bool   `tfsdk:"is_edited"`
	IsFavorite       types.Bool   `tfsdk:"is_favorite"`
	IsOffline        types.Bool   `tfsdk:"is_offline"`
	IsTrashed        types.Bool   `tfsdk:"is_trashed"`
	LibraryId        types.String `tfsdk:"library_id"`
	LivePhotoVideoId types.String `tfsdk:"live_photo_video_id"`
	LocalDateTime    types.String `tfsdk:"local_date_time"`
	OriginalFileName types.String `tfsdk:"original_file_name"`
	OriginalMimeType types.String `tfsdk:"original_mime_type"`
	OriginalPath     types.String `tfsdk:"original_path"`
	OwnerId          types.String `tfsdk:"owner_id"`
	Resized          types.Bool   `tfsdk:"resized"`
	Thumbhash        types.String `tfsdk:"thumbhash"`
	Type             types.String `tfsdk:"type"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	Visibility       types.String `tfsdk:"visibility"`
	Width            types.Int64  `tfsdk:"width"`
}

func (d *viewFolderPathDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view_folders"
}

func (d *viewFolderPathDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state viewFolderPathDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	assets, err := d.client.GetAssetsByOriginalPath(state.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Assets",
			err.Error(),
		)
		return
	}

	for _, asset := range assets {
		assetState := assetModel{
			Checksum:         types.StringValue(asset.Checksum),
			CreatedAt:        types.StringValue(asset.CreatedAt),
			DuplicateId:      types.StringValue(asset.DuplicateId),
			Duration:         types.Int64Value(asset.Duration),
			FileCreatedAt:    types.StringValue(asset.FileCreatedAt),
			FileModifiedAt:   types.StringValue(asset.FileModifiedAt),
			HasMetadata:      types.BoolValue(asset.HasMetadata),
			ID:               types.StringValue(asset.Id),
			Height:           types.Int64Value(asset.Height),
			IsArchived:       types.BoolValue(asset.IsArchived),
			IsEdited:         types.BoolValue(asset.IsEdited),
			IsFavorite:       types.BoolValue(asset.IsFavorite),
			IsOffline:        types.BoolValue(asset.IsOffline),
			IsTrashed:        types.BoolValue(asset.IsTrashed),
			LibraryId:        types.StringValue(asset.LibraryId),
			LivePhotoVideoId: types.StringValue(asset.LivePhotoVideoId),
			LocalDateTime:    types.StringValue(asset.LocalDateTime),
			OriginalFileName: types.StringValue(asset.OriginalFileName),
			OriginalMimeType: types.StringValue(asset.OriginalMimeType),
			OriginalPath:     types.StringValue(asset.OriginalPath),
			OwnerId:          types.StringValue(asset.OwnerId),
			Resized:          types.BoolValue(asset.Resized),
			Thumbhash:        types.StringValue(asset.Thumbhash),
			Type:             types.StringValue(asset.Type),
			UpdatedAt:        types.StringValue(asset.UpdatedAt),
			Visibility:       types.StringValue(asset.Visibility),
			Width:            types.Int64Value(asset.Width),
		}

		state.Assets = append(state.Assets, assetState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *viewFolderPathDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider develops.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *viewFolderPathDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required:    true,
				Description: "Path"},
			"assets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"checksum": schema.StringAttribute{
							Computed:    true,
							Description: "Checksum"},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Time Created"},
						"duplicate_id": schema.StringAttribute{
							Computed:    true,
							Description: "Duplicate ID"},
						"duration": schema.Int64Attribute{
							Computed:    true,
							Description: "Duration"},
						"file_created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Time File Created"},
						"file_modified_at": schema.StringAttribute{
							Computed:    true,
							Description: "Time File Modified"},
						"has_metadata": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether has Metadata"},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID"},
						"height": schema.Int64Attribute{
							Computed:    true,
							Description: "Height"},
						"is_archived": schema.BoolAttribute{
							Computed:    true,
							Description: "Is archived"},
						"is_edited": schema.BoolAttribute{
							Computed:    true,
							Description: "Is edited"},
						"is_favorite": schema.BoolAttribute{
							Computed:    true,
							Description: "Is favourite"},
						"is_offline": schema.BoolAttribute{
							Computed:    true,
							Description: "Is offline"},
						"is_trashed": schema.BoolAttribute{
							Computed:    true,
							Description: "Is trashed"},
						"library_id": schema.StringAttribute{
							Computed:    true,
							Description: "Library Id"},
						"live_photo_video_id": schema.StringAttribute{
							Computed:    true,
							Description: "Live Photo Video Id"},
						"local_date_time": schema.StringAttribute{
							Computed:    true,
							Description: "Local Date and Time"},
						"original_file_name": schema.StringAttribute{
							Computed:    true,
							Description: "Original File Name"},
						"original_mime_type": schema.StringAttribute{
							Computed:    true,
							Description: "Original Mime Type"},
						"original_path": schema.StringAttribute{
							Computed:    true,
							Description: "Original Path"},
						"owner_id": schema.StringAttribute{
							Computed:    true,
							Description: "Owner Id"},
						"resized": schema.BoolAttribute{
							Computed:    true,
							Description: "Resized"},
						"thumbhash": schema.StringAttribute{
							Computed:    true,
							Description: "Thumbnail hash"},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type"},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "Updated at"},
						"visibility": schema.StringAttribute{
							Computed:    true,
							Description: "Visibility"},
						"width": schema.Int64Attribute{
							Computed:    true,
							Description: "Width"},
					},
				},
			},
		},
	}
}
