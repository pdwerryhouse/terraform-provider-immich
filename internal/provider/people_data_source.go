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
	_ datasource.DataSource              = &peopleDataSource{}
	_ datasource.DataSourceWithConfigure = &peopleDataSource{}
)

func NewPeopleDataSource() datasource.DataSource {
	return &peopleDataSource{}
}

type peopleDataSource struct {
	client *immichclient.Client
}

type peopleDataSourceModel struct {
	User []peopleModel `tfsdk:"people"`
}

type peopleModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	BirthDate     types.String `tfsdk:"birthdate"`
	ThumbnailPath types.String `tfsdk:"thumbnail_path"`
	IsHidden      types.Bool   `tfsdk:"is_hidden"`
	IsFavorite    types.Bool   `tfsdk:"is_favorite"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func (d *peopleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_people"
}

func (d *peopleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state peopleDataSourceModel

	people, err := d.client.GetAllPeople(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Person",
			err.Error(),
		)
		return
	}

	for _, person := range people.People {
		peopletate := peopleModel{
			ID:            types.StringValue(person.Id),
			Name:          types.StringValue(person.Name),
			BirthDate:     types.StringValue(person.BirthDate),
			ThumbnailPath: types.StringValue(person.ThumbnailPath),
			IsHidden:      types.BoolValue(person.IsHidden),
			IsFavorite:    types.BoolValue(person.IsFavorite),
			UpdatedAt:     types.StringValue(person.UpdatedAt),
		}

		state.User = append(state.User, peopletate)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *peopleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *peopleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"people": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"birthdate": schema.StringAttribute{
							Computed: true,
						},
						"is_hidden": schema.BoolAttribute{
							Computed: true,
						},
						"is_favorite": schema.BoolAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
						"thumbnail_path": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
