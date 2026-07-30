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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &faceResource{}
	_ resource.ResourceWithConfigure   = &faceResource{}
	_ resource.ResourceWithImportState = &faceResource{}
)

func NewFaceResource() resource.Resource {
	return &faceResource{}
}

type faceResource struct {
	client *immichclient.Client
}

type faceResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AssetId       types.String `tfsdk:"asset_id"`
	BoundingBoxX1 types.Int64  `tfsdk:"bounding_box_x1"`
	BoundingBoxX2 types.Int64  `tfsdk:"bounding_box_x2"`
	BoundingBoxY1 types.Int64  `tfsdk:"bounding_box_y1"`
	BoundingBoxY2 types.Int64  `tfsdk:"bounding_box_y2"`
	ImageHeight   types.Int64  `tfsdk:"image_height"`
	ImageWidth    types.Int64  `tfsdk:"image_width"`
	PersonId      types.String `tfsdk:"person_id"`
	SourceType    types.String `tfsdk:"source_type"`
	Height        types.Int64  `tfsdk:"height"`
	Width         types.Int64  `tfsdk:"width"`
	X             types.Int64  `tfsdk:"x"`
	Y             types.Int64  `tfsdk:"y"`
}

func (r *faceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_face"
}

func (r *faceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
				Required:    true,
				Description: "Image height in pixels"},
			"image_width": schema.Int64Attribute{
				Required:    true,
				Description: "Image width in pixels"},
			"asset_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Id of Asset"},
			"person_id": schema.StringAttribute{
				Required:    true,
				Description: "Id of Person"},
			"source_type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Source Type"},
			"height": schema.Int64Attribute{
				Required:    true,
				Description: "Face Bounding Box Height"},
			"width": schema.Int64Attribute{
				Required:    true,
				Description: "Face Bounding Box Width"},
			"x": schema.Int64Attribute{
				Required:    true,
				Description: "Face Bounding Box X coordinate"},
			"y": schema.Int64Attribute{
				Required:    true,
				Description: "Face Bounding Box Y coordinate"},
		},
	}
}

func (r *faceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan faceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	face := immichclient.AssetFaceCreateDto{
		AssetId:     plan.AssetId.ValueString(),
		Height:      plan.Height.ValueInt64(),
		Width:       plan.Width.ValueInt64(),
		X:           plan.X.ValueInt64(),
		Y:           plan.Y.ValueInt64(),
		ImageHeight: plan.ImageHeight.ValueInt64(),
		ImageWidth:  plan.ImageWidth.ValueInt64(),
		PersonId:    plan.PersonId.ValueString(),
	}

	newFace, err := r.client.CreateFace(face)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating face",
			"Could not create face, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newFace.Id)
	plan.BoundingBoxX1 = types.Int64Value(newFace.BoundingBoxX1)
	plan.BoundingBoxX2 = types.Int64Value(newFace.BoundingBoxX2)
	plan.BoundingBoxY1 = types.Int64Value(newFace.BoundingBoxY1)
	plan.BoundingBoxY2 = types.Int64Value(newFace.BoundingBoxY2)
	plan.SourceType = types.StringValue(newFace.SourceType)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *faceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state faceResourceModel
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

func (r *faceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan faceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	/* do nothing */

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *faceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state faceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFace(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich Face",
			"Could not delete face, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *faceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *faceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
