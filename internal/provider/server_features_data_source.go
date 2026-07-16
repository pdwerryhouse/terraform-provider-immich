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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &serverFeaturesDataSource{}
	_ datasource.DataSourceWithConfigure = &serverFeaturesDataSource{}
)

func NewServerFeaturesDataSource() datasource.DataSource {
	return &serverFeaturesDataSource{}
}

type serverFeaturesDataSource struct {
	client *client.Client
}

type serverFeaturesDataSourceModel struct {
	ConfigFile          types.Bool `tfsdk:"config_file"`
	DuplicateDetection  types.Bool `tfsdk:"duplicate_detection"`
	Email               types.Bool `tfsdk:"email"`
	FacialRecognition   types.Bool `tfsdk:"facial_recognition"`
	ImportFaces         types.Bool `tfsdk:"import_faces"`
	Map                 types.Bool `tfsdk:"map"`
	Oauth               types.Bool `tfsdk:"oauth"`
	OauthAutoLaunch     types.Bool `tfsdk:"oauth_auto_launch"`
	Ocr                 types.Bool `tfsdk:"ocr"`
	PasswordLogin       types.Bool `tfsdk:"password_login"`
	RealtimeTranscoding types.Bool `tfsdk:"realtime_transcoding"`
	ReverseGeocoding    types.Bool `tfsdk:"reverse_geocoding"`
	Search              types.Bool `tfsdk:"search"`
	Sidecar             types.Bool `tfsdk:"sidecar"`
	SmartSearch         types.Bool `tfsdk:"smart_search"`
	Trash               types.Bool `tfsdk:"trash"`
}

func (d *serverFeaturesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_features"
}

func (d *serverFeaturesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverFeaturesDataSourceModel

	serverFeatures, err := d.client.GetServerFeatures()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
	}

	state.ConfigFile = types.BoolValue(serverFeatures.ConfigFile)
	state.DuplicateDetection = types.BoolValue(serverFeatures.DuplicateDetection)
	state.Email = types.BoolValue(serverFeatures.Email)
	state.FacialRecognition = types.BoolValue(serverFeatures.FacialRecognition)
	state.ImportFaces = types.BoolValue(serverFeatures.ImportFaces)
	state.Map = types.BoolValue(serverFeatures.Map)
	state.Oauth = types.BoolValue(serverFeatures.Oauth)
	state.OauthAutoLaunch = types.BoolValue(serverFeatures.OauthAutoLaunch)
	state.Ocr = types.BoolValue(serverFeatures.Ocr)
	state.PasswordLogin = types.BoolValue(serverFeatures.PasswordLogin)
	state.RealtimeTranscoding = types.BoolValue(serverFeatures.RealtimeTranscoding)
	state.ReverseGeocoding = types.BoolValue(serverFeatures.ReverseGeocoding)
	state.Search = types.BoolValue(serverFeatures.Search)
	state.Sidecar = types.BoolValue(serverFeatures.Sidecar)
	state.SmartSearch = types.BoolValue(serverFeatures.SmartSearch)
	state.Trash = types.BoolValue(serverFeatures.Trash)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serverFeaturesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider develops.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *serverFeaturesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config_file": schema.BoolAttribute{
				Computed: true,
			},
			"duplicate_detection": schema.BoolAttribute{
				Computed: true,
			},
			"email": schema.BoolAttribute{
				Computed: true,
			},
			"facial_recognition": schema.BoolAttribute{
				Computed: true,
			},
			"import_faces": schema.BoolAttribute{
				Computed: true,
			},
			"map": schema.BoolAttribute{
				Computed: true,
			},
			"oauth": schema.BoolAttribute{
				Computed: true,
			},
			"oauth_auto_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ocr": schema.BoolAttribute{
				Computed: true,
			},
			"password_login": schema.BoolAttribute{
				Computed: true,
			},
			"realtime_transcoding": schema.BoolAttribute{
				Computed: true,
			},
			"reverse_geocoding": schema.BoolAttribute{
				Computed: true,
			},
			"search": schema.BoolAttribute{
				Computed: true,
			},
			"sidecar": schema.BoolAttribute{
				Computed: true,
			},
			"smart_search": schema.BoolAttribute{
				Computed: true,
			},
			"trash": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}
