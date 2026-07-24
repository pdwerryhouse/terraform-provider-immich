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
	_ resource.Resource                = &libraryResource{}
	_ resource.ResourceWithConfigure   = &libraryResource{}
	_ resource.ResourceWithImportState = &libraryResource{}
)

func NewLibraryResource() resource.Resource {
	return &libraryResource{}
}

type libraryResource struct {
	client *client.Client
}

type libraryResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ExclusionPatterns types.List   `tfsdk:"exclusion_patterns"`
	ImportPaths       types.List   `tfsdk:"import_paths"`
	Name              types.String `tfsdk:"name"`
	OwnerId           types.String `tfsdk:"owner_id"`
}

func (r *libraryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_library"
}

func (r *libraryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"exclusion_patterns": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "Exclusion patterns (max 128)",
			},
			"import_paths": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "Import paths (max 128)",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Library name"},
			"owner_id": schema.StringAttribute{
				Required:    true,
				Description: "Owner user ID"},
		},
	}
}

func (r *libraryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan libraryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var exclusionPatterns []string
	diags = plan.ExclusionPatterns.ElementsAs(ctx, &exclusionPatterns, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var importPaths []string
	diags = plan.ImportPaths.ElementsAs(ctx, &importPaths, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	library := client.CreateLibraryDto{
		ExclusionPatterns: exclusionPatterns,
		ImportPaths:       importPaths,
		Name:              plan.Name.ValueString(),
		OwnerId:           plan.OwnerId.ValueString(),
	}

	newLibrary, err := r.client.CreateLibrary(library)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating library",
			"Could not create library, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newLibrary.Id)
	plan.ExclusionPatterns, diags = types.ListValueFrom(ctx, types.StringType, newLibrary.ExclusionPatterns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ImportPaths, _ = types.ListValueFrom(ctx, types.StringType, newLibrary.ImportPaths)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Name = types.StringValue(newLibrary.Name)
	plan.OwnerId = types.StringValue(newLibrary.OwnerId)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *libraryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state libraryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	library, err := r.client.GetLibrary(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Library",
			"Could not read Immich library ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ExclusionPatterns, diags = types.ListValueFrom(ctx, types.StringType, library.ExclusionPatterns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ImportPaths, _ = types.ListValueFrom(ctx, types.StringType, library.ImportPaths)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Name = types.StringValue(library.Name)
	state.OwnerId = types.StringValue(library.OwnerId)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *libraryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan libraryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ExclusionPatterns []string
	diags = plan.ExclusionPatterns.ElementsAs(ctx, &ExclusionPatterns, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ImportPaths []string
	diags = plan.ImportPaths.ElementsAs(ctx, &ImportPaths, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	library := client.UpdateLibraryDto{
		ExclusionPatterns: ExclusionPatterns,
		ImportPaths:       ImportPaths,
		Name:              plan.Name.ValueString(),
	}

	updatedLibrary, err := r.client.UpdateLibrary(plan.ID.ValueString(), library)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich Library",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(updatedLibrary.Id)
	plan.ExclusionPatterns, diags = types.ListValueFrom(ctx, types.StringType, updatedLibrary.ExclusionPatterns)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ImportPaths, _ = types.ListValueFrom(ctx, types.StringType, updatedLibrary.ImportPaths)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Name = types.StringValue(updatedLibrary.Name)
	plan.OwnerId = types.StringValue(updatedLibrary.OwnerId)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *libraryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state libraryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLibrary(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich Library",
			"Could not delete library, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *libraryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *libraryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
