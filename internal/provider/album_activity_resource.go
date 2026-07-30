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

//import (
//	"context"
//	"fmt"
//	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"
//
//	"github.com/hashicorp/terraform-plugin-framework/path"
//	"github.com/hashicorp/terraform-plugin-framework/resource"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//)
//
//var (
//	_ resource.Resource                = &albumActivityResource{}
//	_ resource.ResourceWithConfigure   = &albumActivityResource{}
//	_ resource.ResourceWithImportState = &albumActivityResource{}
//)
//
//func NewAlbumActivityResource() resource.Resource {
//	return &albumActivityResource{}
//}
//
//type albumActivityResource struct {
//	client *immichclient.Client
//}
//
//type albumActivityResourceModel struct {
//	AlbumId types.String `tfsdk:"album_id"`
//	Enabled types.Bool   `tfsdk:"enabled"`
//}
//
//func (r *albumActivityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
//	resp.TypeName = req.ProviderTypeName + "_album_activity"
//}
//
//func (r *albumActivityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
//	resp.Schema = schema.Schema{
//		Attributes: map[string]schema.Attribute{
//			"album_id": schema.StringAttribute{
//				Required: true,
//			},
//			"enabled": schema.BoolAttribute{
//				Required: true,
//			},
//		},
//	}
//}
//
//func (r *albumActivityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
//	var plan albumActivityResourceModel
//	diags := req.Plan.Get(ctx, &plan)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	albumActivity := client.AlbumActivityUpdate{
//		IsActivityEnabled: plan.Enabled.ValueBool(),
//	}
//
//	_, err := r.client.UpdateAlbumActivity(plan.AlbumId.ValueString(), albumActivity)
//	if err != nil {
//		resp.Diagnostics.AddError(
//			"Error creating albumActivity",
//			"Could not create albumActivity, unexpected error: "+err.Error(),
//		)
//		return
//	}
//
//	diags = resp.State.Set(ctx, plan)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//}
//
//func (r *albumActivityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
//	var state albumActivityResourceModel
//	diags := req.State.Get(ctx, &state)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	album, err := r.client.GetAlbum(state.AlbumId.ValueString())
//	if err != nil {
//		resp.Diagnostics.AddError(
//			"Error Reading AlbumActivity",
//			"Could not read Immich albumActivity ID "+state.AlbumId.ValueString()+": "+err.Error(),
//		)
//		return
//	}
//
//	state.Enabled = types.BoolValue(album.IsActivityEnabled)
//
//	diags = resp.State.Set(ctx, &state)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//}
//
//func (r *albumActivityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
//	var plan albumActivityResourceModel
//	diags := req.Plan.Get(ctx, &plan)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	albumActivity := client.AlbumActivityUpdate{
//		IsActivityEnabled: plan.Enabled.ValueBool(),
//	}
//
//	updatedAlbumActivity, err := r.client.UpdateAlbumActivity(plan.AlbumId.ValueString(), albumActivity)
//
//	if err != nil {
//		resp.Diagnostics.AddError(
//			"Error Updating Immich AlbumActivity",
//			"Could not update order, unexpected error: "+err.Error(),
//		)
//		return
//	}
//
//	plan.Enabled = types.BoolValue(updatedAlbumActivity.IsActivityEnabled)
//
//	diags = resp.State.Set(ctx, plan)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//}
//
//func (r *albumActivityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
//	var state albumActivityResourceModel
//	diags := req.State.Get(ctx, &state)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//}
//
//func (r *albumActivityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
//	if req.ProviderData == nil {
//		return
//	}
//
//	client, ok := req.ProviderData.(*immichclient.Client)
//
//	if !ok {
//		resp.Diagnostics.AddError(
//			"Unexpected Resource Configure Type",
//			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider developer.", req.ProviderData),
//		)
//
//		return
//	}
//
//	r.client = client
//}
//
//func (r *albumActivityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
//	// Retrieve import ID and save to id attribute
//	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
//}
//
