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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &apiKeyResource{}
	_ resource.ResourceWithConfigure   = &apiKeyResource{}
	_ resource.ResourceWithImportState = &apiKeyResource{}
)

func NewApiKeyResource() resource.Resource {
	return &apiKeyResource{}
}

type apiKeyResource struct {
	client *client.Client
}

type apiKeyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Permissions types.List   `tfsdk:"permissions"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Secret      types.String `tfsdk:"secret"`
}

func (r *apiKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *apiKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "API key name",
				Required:    true,
			},
			"permissions": schema.ListAttribute{
				Description: "Permissions as a list of strings (see https://api.immich.app/models/Permission) ",
				ElementType: types.StringType,
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation date.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update date.",
				Computed:    true,
			},
			"secret": schema.StringAttribute{
				Description: "API key secret (warning: stored in cleartext in state)",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	diags = plan.Permissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := client.ApiKeyCreateDto{
		Name:        plan.Name.ValueString(),
		Permissions: permissions,
	}

	newApiKey, err := r.client.CreateApiKey(apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating apiKey",
			"Could not create apiKey, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newApiKey.ApiKey.Id)
	plan.Name = types.StringValue(newApiKey.ApiKey.Name)
	plan.CreatedAt = types.StringValue(newApiKey.ApiKey.CreatedAt)
	plan.UpdatedAt = types.StringValue(newApiKey.ApiKey.UpdatedAt)
	plan.Secret = types.StringValue(newApiKey.Secret)

	plan.Permissions, diags = types.ListValueFrom(ctx, types.StringType, newApiKey.ApiKey.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *apiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.client.GetApiKey(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ApiKey",
			"Could not read Immich apiKey ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(apiKey.Name)
	state.UpdatedAt = types.StringValue(apiKey.UpdatedAt)
	state.CreatedAt = types.StringValue(apiKey.CreatedAt)

	state.Permissions, diags = types.ListValueFrom(ctx, types.StringType, apiKey.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *apiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permissions []string
	diags = plan.Permissions.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := client.ApiKeyUpdateDto{
		Name:        plan.Name.ValueString(),
		Permissions: permissions,
	}

	updatedApiKey, err := r.client.UpdateApiKey(plan.ID.ValueString(), apiKey)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich ApiKey",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(updatedApiKey.Id)
	plan.Name = types.StringValue(updatedApiKey.Name)
	plan.CreatedAt = types.StringValue(updatedApiKey.CreatedAt)
	plan.UpdatedAt = types.StringValue(updatedApiKey.UpdatedAt)

	plan.Permissions, diags = types.ListValueFrom(ctx, types.StringType, updatedApiKey.Permissions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *apiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApiKey(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich ApiKey",
			"Could not delete apiKey, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *apiKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
