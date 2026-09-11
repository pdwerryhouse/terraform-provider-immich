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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &partnerResource{}
	_ resource.ResourceWithConfigure   = &partnerResource{}
	_ resource.ResourceWithImportState = &partnerResource{}
)

func NewPartnerResource() resource.Resource {
	return &partnerResource{}
}

type partnerResource struct {
	client *immichclient.Client
}

type partnerResourceModel struct {
	ID         types.String `tfsdk:"id"`
	InTimeline types.Bool   `tfsdk:"in_timeline"`
}

func (r *partnerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner"
}

func (r *partnerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"in_timeline": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *partnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan partnerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	partner := immichclient.PartnerCreateDto{
		SharedWithId: plan.ID.ValueString(),
	}

	newPartner, err := r.client.CreatePartner(partner)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating partner",
			"Could not create partner, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newPartner.Id)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *partnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state partnerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	Partner, err := r.client.GetPartners(immichclient.PartnerDirectionsharedby)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Partner",
			err.Error(),
		)
		return
	}

	var partnerState partnerResourceModel

	for _, partner := range Partner {
		if partner.Id == state.ID.ValueString() {
			partnerState = partnerResourceModel{
				ID: types.StringValue(partner.Id),
			}
			break
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Partner",
			"Could not read Immich partner ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = partnerState.ID

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *partnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan partnerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	partner := immichclient.PartnerUpdateDto{
		InTimeline: plan.InTimeline.ValueBool(),
	}

	updatedPartner, err := r.client.UpdatePartner(plan.ID.ValueString(), partner)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich Partner",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(updatedPartner.Id)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *partnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state partnerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePartner(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich Partner",
			"Could not delete partner, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *partnerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *partnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
