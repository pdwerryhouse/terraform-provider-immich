package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-immich/internal/client"
)

var (
	_ datasource.DataSource              = &albumsDataSource{}
	_ datasource.DataSourceWithConfigure = &albumsDataSource{}
)

func NewAlbumsDataSource() datasource.DataSource {
	return &albumsDataSource{}
}

type albumsDataSource struct {
	client *client.Client
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

	albums, err := d.client.GetAlbums()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Album",
			err.Error(),
		)
		return
	}

	for _, album := range albums {
		albumState := albumsModel{
			ID:                    types.StringValue(album.ID),
			Name:                  types.StringValue(album.AlbumName),
			AlbumThumbnailAssetId: types.StringValue(album.AlbumThumbnailAssetId),
			Description:           types.StringValue(album.Description),
			Shared:                types.BoolValue(album.Shared),
			HasSharedLink:         types.BoolValue(album.HasSharedLink),
			Order:                 types.StringValue(album.Order),
			IsActivityEnabled:     types.BoolValue(album.IsActivityEnabled),
			CreatedAt:             types.StringValue(album.CreatedAt),
			UpdatedAt:             types.StringValue(album.UpdatedAt),
			StartDate:             types.StringValue(album.StartDate),
			EndDate:               types.StringValue(album.EndDate),
			OwnerId:               types.StringValue(album.OwnerId),
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

func (d *albumsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"albums": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"album_name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"shared": schema.BoolAttribute{
							Computed: true,
						},
						"order": schema.StringAttribute{
							Computed: true,
						},
						"album_thumbnail_asset_id": schema.StringAttribute{
							Computed: true,
						},
						"has_shared_link": schema.BoolAttribute{
							Computed: true,
						},
						"is_activity_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
						"start_date": schema.StringAttribute{
							Computed: true,
						},
						"end_date": schema.StringAttribute{
							Computed: true,
						},
						"owner_id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
