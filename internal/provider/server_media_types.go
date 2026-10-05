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
	_ datasource.DataSource              = &serverMediaTypesDataSource{}
	_ datasource.DataSourceWithConfigure = &serverMediaTypesDataSource{}
)

func NewServerMediaTypesDataSource() datasource.DataSource {
	return &serverMediaTypesDataSource{}
}

type serverMediaTypesDataSource struct {
	client *immichclient.Client
}

type serverMediaTypesDataSourceModel struct {
	Image   []types.String `tfsdk:"image"`
	Sidecar []types.String `tfsdk:"sidecar"`
	Video   []types.String `tfsdk:"video"`
}

func (d *serverMediaTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_media_types"
}

func (d *serverMediaTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverMediaTypesDataSourceModel

	serverMediaTypes, err := d.client.GetSupportedMediaTypes()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Server Media Types",
			err.Error(),
		)
		return
	}

	var images []types.String

	for _, image := range serverMediaTypes.Image {
		images = append(images, types.StringValue(image))
	}
	state.Image = images

	var sidecars []types.String

	for _, sidecar := range serverMediaTypes.Sidecar {
		sidecars = append(images, types.StringValue(sidecar))
	}
	state.Sidecar = sidecars

	var videos []types.String

	for _, video := range serverMediaTypes.Video {
		videos = append(images, types.StringValue(video))
	}
	state.Video = videos

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serverMediaTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source MediaTypesure Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider develops.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *serverMediaTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve all media types supported by the Immich server.",
		Attributes: map[string]schema.Attribute{

			"image": schema.ListAttribute{
				Description: "Supported image MIME types",
				ElementType: types.StringType,
				Computed:    true,
			},
			"sidecar": schema.ListAttribute{
				Description: "Supported sidecar MIME types",
				ElementType: types.StringType,
				Computed:    true,
			},
			"video": schema.ListAttribute{
				Description: "Supported video MIME types",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}
