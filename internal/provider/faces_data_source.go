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
	_ datasource.DataSource              = &FaceDataSource{}
	_ datasource.DataSourceWithConfigure = &FaceDataSource{}
)

func NewFacesDataSource() datasource.DataSource {
	return &FaceDataSource{}
}

type FaceDataSource struct {
	client *immichclient.Client
}

type FaceDataSourceModel struct {
	AssetId types.String `tfsdk:"asset_id"`
	Face    []FaceModel  `tfsdk:"face"`
}

type FaceModel struct {
	BoundingBoxX1 types.Int64  `tfsdk:"bounding_box_x1"`
	BoundingBoxX2 types.Int64  `tfsdk:"bounding_box_x2"`
	BoundingBoxY1 types.Int64  `tfsdk:"bounding_box_y1"`
	BoundingBoxY2 types.Int64  `tfsdk:"bounding_box_y2"`
	ID            types.String `tfsdk:"id"`
	ImageHeight   types.Int64  `tfsdk:"image_height"`
	ImageWidth    types.Int64  `tfsdk:"image_width"`
	PersonId      types.String `tfsdk:"person_id"`
	PersonName    types.String `tfsdk:"person_name"`
	SourceType    types.String `tfsdk:"source_type"`
}

func (d *FaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_faces"
}

func (d *FaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state FaceDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	faceParams := immichclient.GetFacesParameters{
		Id: state.AssetId.String(),
	}

	Face, err := d.client.GetFaces(faceParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Face",
			err.Error(),
		)
		return
	}

	for _, face := range Face {
		faceState := FaceModel{
			BoundingBoxX1: types.Int64Value(face.BoundingBoxX1),
			BoundingBoxX2: types.Int64Value(face.BoundingBoxX2),
			BoundingBoxY1: types.Int64Value(face.BoundingBoxY1),
			BoundingBoxY2: types.Int64Value(face.BoundingBoxY2),
			ID:            types.StringValue(face.Id),
			ImageHeight:   types.Int64Value(face.ImageHeight),
			ImageWidth:    types.Int64Value(face.ImageWidth),
			PersonId:      types.StringValue(face.Person.Id),
			PersonName:    types.StringValue(face.Person.Name),
			SourceType:    types.StringValue(face.SourceType),
		}

		state.Face = append(state.Face, faceState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *FaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *FaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"asset_id": schema.StringAttribute{
				Required:    true,
				Description: "Face Id"},
			"face": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"bounding_box_x1": schema.Int64Attribute{
							Computed:    true,
							Description: "Bounding box X1 coordinate"},
						"bounding_box_x2": schema.Int64Attribute{
							Computed:    true,
							Description: "Bounding box X2 coordinate"},
						"bounding_box_y1": schema.Int64Attribute{
							Computed:    true,
							Description: "Bounding box Y1 coordinate"},
						"bounding_box_y2": schema.Int64Attribute{
							Computed:    true,
							Description: "Bounding box Y2 coordinate"},
						"image_height": schema.Int64Attribute{
							Computed:    true,
							Description: "Image height in pixels"},
						"image_width": schema.Int64Attribute{
							Computed:    true,
							Description: "Image width in pixels"},
						"person_id": schema.StringAttribute{
							Computed:    true,
							Description: "Id of Person"},
						"person_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of Person"},
						"source_type": schema.StringAttribute{
							Computed:    true,
							Description: "Source Type",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Face ID"},
					},
				},
			},
		},
	}
}
