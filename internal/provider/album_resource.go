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
	"time"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &albumResource{}
	_ resource.ResourceWithConfigure   = &albumResource{}
	_ resource.ResourceWithImportState = &albumResource{}
)

func NewAlbumResource() resource.Resource {
	return &albumResource{}
}

type albumResource struct {
	client *immichclient.Client
}

type albumResourceModel struct {
	ID          types.String `tfsdk:"id"`
	AlbumName   types.String `tfsdk:"album_name"`
	Description types.String `tfsdk:"description"`
	Order       types.String `tfsdk:"order"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *albumResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_album"
}

func (r *albumResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"album_name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"order": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *albumResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan albumResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	album := immichclient.CreateAlbumDto{
		AlbumName:   plan.AlbumName.ValueString(),
		Description: plan.Description.ValueString(),
	}

	newAlbum, err := r.client.CreateAlbum(album)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating album",
			"Could not create album, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newAlbum.Id)
	plan.AlbumName = types.StringValue(newAlbum.AlbumName)
	plan.Description = types.StringValue(newAlbum.Description)
	plan.Order = types.StringValue(newAlbum.Order)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *albumResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state albumResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	album, err := r.client.GetAlbum(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Album",
			"Could not read Immich album ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.AlbumName = types.StringValue(album.AlbumName)
	state.Description = types.StringValue(album.Description)
	state.Order = types.StringValue(album.Order)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *albumResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan albumResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	album := immichclient.UpdateAlbumDto{
		AlbumName:   plan.AlbumName.ValueString(),
		Description: plan.Description.ValueString(),
		Order:       plan.Order.ValueString(),
	}

	updatedAlbum, err := r.client.UpdateAlbum(plan.ID.ValueString(), album)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich Album",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(updatedAlbum.Id)
	plan.AlbumName = types.StringValue(updatedAlbum.AlbumName)
	plan.Description = types.StringValue(updatedAlbum.Description)
	plan.Order = types.StringValue(updatedAlbum.Order)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *albumResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state albumResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAlbum(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich Album",
			"Could not delete album, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *albumResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *albumResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
