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
	"regexp"
	"terraform-provider-immich/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *client.Client
}

type userResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Email                types.String `tfsdk:"email"`
	Password             types.String `tfsdk:"password"`
	PinCode              types.String `tfsdk:"pincode"`
	IsAdmin              types.Bool   `tfsdk:"is_admin"`
	Notify               types.Bool   `tfsdk:"notify"`
	QuotaSizeInBytes     types.Int64  `tfsdk:"quota_size_in_bytes"`
	ShouldChangePassword types.Bool   `tfsdk:"should_change_password"`
	StorageLabel         types.String `tfsdk:"storage_label"`
	AvatarColor          types.String `tfsdk:"avatar_color"`
	//LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create an Immich user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The user name.",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "The user's email address.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The user's password.",
				Required:    true,
			},
			"is_admin": schema.BoolAttribute{
				Description: "Set to true if this is an user with admin privileges.",
				Optional:    true,
			},
			"notify": schema.BoolAttribute{
				Description: "Set to true if this is an user should be sent a notification email.",
				Optional:    true,
			},
			"quota_size_in_bytes": schema.Int64Attribute{
				Description: "The user's quota.",
				Optional:    true,
			},
			"should_change_password": schema.BoolAttribute{
				Description: "Set to true if the user should be required to change their password upon first login.",
				Optional:    true,
			},
			"storage_label": schema.StringAttribute{
				Description: "The user's storage label.",
				Optional:    true,
			},
			"pincode": schema.StringAttribute{
				Description: "The user's pin code.",
				Optional:    true,
			},
			"avatar_color": schema.StringAttribute{
				Description: "The user's avatar colour. Valid choices: primary, pink, red, yellow, blue, green, purple, orange, gray, amber.",
				Optional:    true,
				Validators: []validator.String{
					// These are example validators from terraform-plugin-framework-validators
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(primary|pink|red|yellow|blue|green|purple|orange|gray|amber)$`),
						"must be 'primary', 'pink', 'red', 'yellow', 'blue', 'green', 'purple, 'orange', 'gray', 'amber'.",
					),
				},
			},
			/*
				"last_updated": schema.StringAttribute{
					Computed: true,
				},
			*/
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := client.UserAdminCreateDto{
		Name:                 plan.Name.ValueString(),
		Email:                plan.Email.ValueString(),
		Password:             plan.Password.ValueString(),
		IsAdmin:              plan.IsAdmin.ValueBool(),
		Notify:               plan.Notify.ValueBool(),
		QuotaSizeInBytes:     plan.QuotaSizeInBytes.ValueInt64(),
		StorageLabel:         plan.StorageLabel.ValueString(),
		ShouldChangePassword: plan.ShouldChangePassword.ValueBool(),
		PinCode:              plan.PinCode.ValueString(),
		AvatarColor:          plan.AvatarColor.ValueString(),
	}

	newUser, err := r.client.CreateUser(user)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newUser.Id)
	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			"Could not read Immich user ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(user.Name)
	state.Email = types.StringValue(user.Email)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := client.UserAdminUpdateDto{
		Name:                 plan.Name.ValueString(),
		Email:                plan.Email.ValueString(),
		Password:             plan.Password.ValueString(),
		IsAdmin:              plan.IsAdmin.ValueBool(),
		QuotaSizeInBytes:     plan.QuotaSizeInBytes.ValueInt64(),
		StorageLabel:         plan.StorageLabel.ValueString(),
		ShouldChangePassword: plan.ShouldChangePassword.ValueBool(),
		PinCode:              plan.PinCode.ValueString(),
		AvatarColor:          plan.AvatarColor.ValueString(),
	}

	updatedUser, err := r.client.UpdateUser(plan.ID.ValueString(), user)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich Album",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(updatedUser.Id)
	//plan.Name = types.StringValue(updatedUser.Name)
	//plan.Email = types.StringValue(updatedUser.Email)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich User",
			"Could not delete album, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider developer.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
