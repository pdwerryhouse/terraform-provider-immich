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
	_ resource.Resource                = &activityResource{}
	_ resource.ResourceWithConfigure   = &activityResource{}
	_ resource.ResourceWithImportState = &activityResource{}
)

func NewActivityResource() resource.Resource {
	return &activityResource{}
}

type activityResource struct {
	client *immichclient.Client
}

type activityResourceModel struct {
	ID        types.String `tfsdk:"id"`
	AlbumId   types.String `tfsdk:"album_id"`
	AssetId   types.String `tfsdk:"asset_id"`
	Comment   types.String `tfsdk:"comment"`
	CreatedAt types.String `tfsdk:"created_at"`
	Type      types.String `tfsdk:"type"`
}

func (r *activityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_activity"
}

func (r *activityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"album_id": schema.StringAttribute{
				Required:    true,
				Description: "Album Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"asset_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Asset Id",
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Comment Text",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Type of activity (like or comment)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Creation date"},
		},
	}
}

func (r *activityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan activityResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	activity := immichclient.ActivityCreateDto{
		AlbumId: plan.AlbumId.ValueString(),
		AssetId: plan.AssetId.ValueStringPointer(),
		Comment: plan.Comment.ValueStringPointer(),
		Type:    immichclient.ReactionType(plan.Type.ValueString()),
	}

	newActivity, err := r.client.CreateActivity(activity)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating activity",
			"Could not create activity, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newActivity.Id)
	plan.AssetId = types.StringValue(newActivity.AssetId)
	plan.Comment = types.StringValue(newActivity.Comment)
	plan.CreatedAt = types.StringValue(newActivity.CreatedAt)
	plan.Type = types.StringValue(string(newActivity.Type))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *activityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state activityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	/* do nothing */

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *activityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	/* do nothing */
}

func (r *activityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state activityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteActivity(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich Activity",
			"Could not delete activity, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *activityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *activityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
