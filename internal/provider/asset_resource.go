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
	"encoding/base64"
	"fmt"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &assetResource{}
	_ resource.ResourceWithConfigure   = &assetResource{}
	_ resource.ResourceWithImportState = &assetResource{}
)

func NewAssetResource() resource.Resource {
	return &assetResource{}
}

type assetResource struct {
	client *immichclient.Client
}

type assetResourceModel struct {
	ID             types.String `tfsdk:"id"`
	AssetData      types.String `tfsdk:"asset_data"`
	Filename       types.String `tfsdk:"filename"`
	FileCreatedAt  types.String `tfsdk:"created_at"`
	FileModifiedAt types.String `tfsdk:"modified_at"`
}

func (r *assetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset"
}

func (r *assetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filename": schema.StringAttribute{
				Required: true,
			},
			"asset_data": schema.StringAttribute{
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Required: true,
			},
			"modified_at": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *assetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan assetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	image, err := base64.StdEncoding.DecodeString(plan.AssetData.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading image",
			"Could not create asset, unexpected error: "+err.Error(),
		)
		return
	}

	asset := immichclient.AssetMediaCreateDto{
		AssetData:      image,
		Filename:       plan.Filename.ValueStringPointer(),
		FileCreatedAt:  plan.FileCreatedAt.ValueString(),
		FileModifiedAt: plan.FileModifiedAt.ValueString(),
	}

	newAsset, err := r.client.UploadAsset(asset)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating asset",
			"Could not create asset, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newAsset.Id)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *assetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state assetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	image, err := r.client.DownloadAsset(state.ID.ValueString(), nil, "", "")

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Asset",
			"Could not read Immich asset ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	b64image := base64.StdEncoding.EncodeToString(image)

	state.AssetData = types.StringValue(b64image)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *assetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *assetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state assetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	force_true := true

	assets := immichclient.AssetBulkDeleteDto{
		Force: &force_true,
		Ids:   []string{state.ID.ValueString()},
	}

	err := r.client.DeleteAssets(assets)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich Asset",
			"Could not delete asset, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *assetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *assetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
