// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

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
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	client *client.Client
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

	users, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
	}

	for _, user := range users {
		userState := usersModel{
			ID:        types.StringValue(user.ID),
			Name:      types.StringValue(user.Name),
			Email:     types.StringValue(user.Email),
			Status:    types.StringValue(user.Status),
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

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"is_admin": schema.BoolAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
						"deleted_at": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
