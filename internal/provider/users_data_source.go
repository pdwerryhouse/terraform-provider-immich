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
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	client *immichclient.Client
}

type usersDataSourceModel struct {
	User []usersModel `tfsdk:"users"`
}

type usersModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Email     types.String `tfsdk:"email"`
	IsAdmin   types.Bool   `tfsdk:"is_admin"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	DeletedAt types.String `tfsdk:"deleted_at"`
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel

	users, err := d.client.SearchUsersAdmin(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
	}

	for _, user := range users {
		userState := usersModel{
			ID:        types.StringValue(user.Id),
			Name:      types.StringValue(user.Name),
			Email:     types.StringValue(user.Email),
			Status:    types.StringValue(string(user.Status)),
			IsAdmin:   types.BoolValue(user.IsAdmin),
			CreatedAt: types.StringValue(user.CreatedAt),
			UpdatedAt: types.StringValue(user.UpdatedAt),
			DeletedAt: types.StringValue(user.DeletedAt),
		}

		state.User = append(state.User, userState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of users.",
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "User Id.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the user.",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "The user's email address.",
							Computed:    true,
						},
						"is_admin": schema.BoolAttribute{
							Description: "If true, the user is an administrator.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The user's status.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The date when the user was created.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "The date when the user was last modified.",
							Computed:    true,
						},
						"deleted_at": schema.StringAttribute{
							Description: "The date when the user was deleted.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
