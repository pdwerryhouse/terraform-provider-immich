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

package resources

//import (
//	"context"
//	"fmt"
//	"regexp"
//	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"
//
//	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
//	"github.com/hashicorp/terraform-plugin-framework/path"
//	"github.com/hashicorp/terraform-plugin-framework/resource"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//)
//
//var (
//	_ resource.Resource                = &albumOrderResource{}
//	_ resource.ResourceWithConfigure   = &albumOrderResource{}
//	_ resource.ResourceWithImportState = &albumOrderResource{}
//)
//
//func NewAlbumOrderResource() resource.Resource {
//	return &albumOrderResource{}
//}
//
//type albumOrderResource struct {
//	client *immichclient.Client
//}
//
//type albumOrderResourceModel struct {
//	AlbumId types.String `tfsdk:"album_id"`
//	Order   types.String `tfsdk:"order"`
//}
//
//func (r *albumOrderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
//	resp.TypeName = req.ProviderTypeName + "_album_order"
//}
//
//func (r *albumOrderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
//	resp.Schema = schema.Schema{
//		Attributes: map[string]schema.Attribute{
//			"album_id": schema.StringAttribute{
//				Required: true,
//			},
//			"order": schema.StringAttribute{
//				Required: true,
//				Validators: []validator.String{
//					// These are example validators from terraform-plugin-framework-validators
//					stringvalidator.RegexMatches(
//						regexp.MustCompile(`^(asc|desc)$`),
//						"must be 'asc' or 'desc'.",
//					),
//				},
//			},
//		},
//	}
//}
//
//func (r *albumOrderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
//	var plan albumOrderResourceModel
//	diags := req.Plan.Get(ctx, &plan)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	albumOrder := client.AlbumOrderUpdate{
//		Order: plan.Order.ValueString(),
//	}
//
//	_, err := r.client.UpdateAlbumOrder(plan.AlbumId.ValueString(), albumOrder)
//	if err != nil {
//		resp.Diagnostics.AddError(
//			"Error creating albumOrder",
//			"Could not create albumOrder, unexpected error: "+err.Error(),
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
//func (r *albumOrderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
//	var state albumOrderResourceModel
//	diags := req.State.Get(ctx, &state)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	album, err := r.client.GetAlbum(state.AlbumId.ValueString())
//	if err != nil {
//		resp.Diagnostics.AddError(
//			"Error Reading AlbumOrder",
//			"Could not read Immich albumOrder ID "+state.AlbumId.ValueString()+": "+err.Error(),
//		)
//		return
//	}
//
//	state.Order = types.StringValue(album.Order)
//
//	diags = resp.State.Set(ctx, &state)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//}
//
//func (r *albumOrderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
//	var plan albumOrderResourceModel
//	diags := req.Plan.Get(ctx, &plan)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	albumOrder := client.AlbumOrderUpdate{
//		Order: plan.Order.ValueString(),
//	}
//
//	updatedAlbumOrder, err := r.client.UpdateAlbumOrder(plan.AlbumId.ValueString(), albumOrder)
//
//	if err != nil {
//		resp.Diagnostics.AddError(
//			"Error Updating Immich AlbumOrder",
//			"Could not update order, unexpected error: "+err.Error(),
//		)
//		return
//	}
//
//	plan.Order = types.StringValue(updatedAlbumOrder.Order)
//
//	diags = resp.State.Set(ctx, plan)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//}
//
//func (r *albumOrderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
//	var state albumOrderResourceModel
//	diags := req.State.Get(ctx, &state)
//	resp.Diagnostics.Append(diags...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//}
//
//func (r *albumOrderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
//func (r *albumOrderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
//	// Retrieve import ID and save to id attribute
//	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
//}
//
