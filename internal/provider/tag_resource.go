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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &tagResource{}
	_ resource.ResourceWithConfigure   = &tagResource{}
	_ resource.ResourceWithImportState = &tagResource{}
)

func NewTagResource() resource.Resource {
	return &tagResource{}
}

type tagResource struct {
	client *immichclient.Client
}

type tagResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Color    types.String `tfsdk:"color"`
	Name     types.String `tfsdk:"name"`
	ParentId types.String `tfsdk:"parent_id"`
}

func (r *tagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *tagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"color": schema.StringAttribute{
				Optional: true,
			},
			"parent_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
		},
	}
}

func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag := immichclient.TagCreateDto{
		Name:     plan.Name.ValueString(),
		Color:    plan.Color.ValueString(),
		ParentId: plan.ParentId.ValueString(),
	}

	newTag, err := r.client.CreateTag(tag)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tag",
			"Could not create tag, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newTag.Id)
	plan.Color = types.StringValue(newTag.Color)
	plan.Name = types.StringValue(newTag.Name)
	plan.ParentId = types.StringValue(newTag.ParentId)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag, err := r.client.GetTagById(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			"Could not read Immich tag ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Color = types.StringValue(tag.Color)
	state.Name = types.StringValue(tag.Name)
	state.ParentId = types.StringValue(tag.ParentId)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag := immichclient.TagUpdateDto{
		Color: plan.Color.ValueString(),
	}

	updatedTag, err := r.client.UpdateTag(plan.ID.ValueString(), tag)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich Tag",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(updatedTag.Id)
	plan.Color = types.StringValue(tag.Color)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTag(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich Tag",
			"Could not delete tag, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *tagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider developer.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
