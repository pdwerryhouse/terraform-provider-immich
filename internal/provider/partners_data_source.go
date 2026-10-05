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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"
)

var (
	_ datasource.DataSource              = &PartnerDataSource{}
	_ datasource.DataSourceWithConfigure = &PartnerDataSource{}
)

func NewPartnerDataSource() datasource.DataSource {
	return &PartnerDataSource{}
}

type PartnerDataSource struct {
	client *immichclient.Client
}

type PartnerDataSourceModel struct {
	Partner []PartnerModel `tfsdk:"partner"`
}

type PartnerModel struct {
	AvatarColor      types.String `tfsdk:"avatar_color"`
	Email            types.String `tfsdk:"email"`
	Id               types.String `tfsdk:"id"`
	InTimeline       types.Bool   `tfsdk:"in_timeline"`
	Name             types.String `tfsdk:"name"`
	ProfileChangedAt types.String `tfsdk:"profile_changed_at"`
	ProfileImagePath types.String `tfsdk:"profile_image_path"`
}

func (d *PartnerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partners"
}

func (d *PartnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state PartnerDataSourceModel

	Partner, err := d.client.GetPartners(immichclient.GetPartnersParameters{Direction: immichclient.PartnerDirectionSharedBy})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Partner",
			err.Error(),
		)
		return
	}

	for _, partner := range Partner {
		partnerState := PartnerModel{
			AvatarColor:      types.StringValue(string(partner.AvatarColor)),
			Email:            types.StringValue(partner.Email),
			Id:               types.StringValue(partner.Id),
			InTimeline:       types.BoolValue(partner.InTimeline),
			Name:             types.StringValue(partner.Name),
			ProfileChangedAt: types.StringValue(partner.ProfileChangedAt),
			ProfileImagePath: types.StringValue(partner.ProfileImagePath),
		}

		state.Partner = append(state.Partner, partnerState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *PartnerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *PartnerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"partner": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"avatar_color": schema.StringAttribute{
							Computed:    true,
							Description: "User's Avatar Colour"},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "User email"},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "User ID"},
						"in_timeline": schema.BoolAttribute{
							Computed:    true,
							Description: "Show in timeline"},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "User name"},
						"profile_changed_at": schema.StringAttribute{
							Computed:    true,
							Description: "Profile change date"},
						"profile_image_path": schema.StringAttribute{
							Computed:    true,
							Description: "Profile image path"},
					},
				},
			},
		},
	}
}
