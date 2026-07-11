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
	_ datasource.DataSource              = &serverAboutDataSource{}
	_ datasource.DataSourceWithConfigure = &serverAboutDataSource{}
)

func NewServerAboutDataSource() datasource.DataSource {
	return &serverAboutDataSource{}
}

type serverAboutDataSource struct {
	client *client.Client
}

type serverAboutDataSourceModel struct {
	Build                      types.String `tfsdk:"build"`
	BuildImage                 types.String `tfsdk:"build_image"`
	BuildImageUrl              types.String `tfsdk:"build_image_url"`
	BuildUrl                   types.String `tfsdk:"build_url"`
	Exiftool                   types.String `tfsdk:"exiftool"`
	Ffmpeg                     types.String `tfsdk:"ffmpeg"`
	Imagemagick                types.String `tfsdk:"imagemagick"`
	Libvips                    types.String `tfsdk:"libvips"`
	Licensed                   types.Bool   `tfsdk:"licensed"`
	Nodejs                     types.String `tfsdk:"nodejs"`
	Repository                 types.String `tfsdk:"repository"`
	RepositoryUrl              types.String `tfsdk:"repository_url"`
	SourceCommit               types.String `tfsdk:"source_commit"`
	SourceRef                  types.String `tfsdk:"source_ref"`
	SourceUrl                  types.String `tfsdk:"source_url"`
	ThirdPartyBugFeatureUrl    types.String `tfsdk:"third_party_bug_feature_url"`
	ThirdPartyDocumentationUrl types.String `tfsdk:"third_party_documentation_url"`
	ThirdPartySourceUrl        types.String `tfsdk:"third_party_source_url"`
	ThirdPartySupportUrl       types.String `tfsdk:"third_party_support_url"`
	Version                    types.String `tfsdk:"version"`
	VersionUrl                 types.String `tfsdk:"version_url"`
}

func (d *serverAboutDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_about"
}

func (d *serverAboutDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverAboutDataSourceModel

	serverAbout, err := d.client.GetServerAbout()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
	}

	state.Build = types.StringValue(serverAbout.Build)
	state.BuildImage = types.StringValue(serverAbout.BuildImage)
	state.BuildImageUrl = types.StringValue(serverAbout.BuildImageUrl)
	state.BuildUrl = types.StringValue(serverAbout.BuildUrl)
	state.Exiftool = types.StringValue(serverAbout.Exiftool)
	state.Ffmpeg = types.StringValue(serverAbout.Ffmpeg)
	state.Imagemagick = types.StringValue(serverAbout.Imagemagick)
	state.Libvips = types.StringValue(serverAbout.Libvips)
	state.Licensed = types.BoolValue(serverAbout.Licensed)
	state.Nodejs = types.StringValue(serverAbout.Nodejs)
	state.Repository = types.StringValue(serverAbout.Repository)
	state.RepositoryUrl = types.StringValue(serverAbout.RepositoryUrl)
	state.SourceCommit = types.StringValue(serverAbout.SourceCommit)
	state.SourceRef = types.StringValue(serverAbout.SourceRef)
	state.SourceUrl = types.StringValue(serverAbout.SourceUrl)
	state.ThirdPartyBugFeatureUrl = types.StringValue(serverAbout.ThirdPartyBugFeatureUrl)
	state.ThirdPartyDocumentationUrl = types.StringValue(serverAbout.ThirdPartyDocumentationUrl)
	state.ThirdPartySourceUrl = types.StringValue(serverAbout.ThirdPartySourceUrl)
	state.ThirdPartySupportUrl = types.StringValue(serverAbout.ThirdPartySupportUrl)
	state.Version = types.StringValue(serverAbout.Version)
	state.VersionUrl = types.StringValue(serverAbout.VersionUrl)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serverAboutDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *serverAboutDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"build": schema.StringAttribute{
				Computed: true,
			},
			"build_image": schema.StringAttribute{
				Computed: true,
			},
			"build_image_url": schema.StringAttribute{
				Computed: true,
			},
			"build_url": schema.StringAttribute{
				Computed: true,
			},
			"exiftool": schema.StringAttribute{
				Computed: true,
			},
			"ffmpeg": schema.StringAttribute{
				Computed: true,
			},
			"imagemagick": schema.StringAttribute{
				Computed: true,
			},
			"libvips": schema.StringAttribute{
				Computed: true,
			},
			"licensed": schema.BoolAttribute{
				Computed: true,
			},
			"nodejs": schema.StringAttribute{
				Computed: true,
			},
			"repository": schema.StringAttribute{
				Computed: true,
			},
			"repository_url": schema.StringAttribute{
				Computed: true,
			},
			"source_commit": schema.StringAttribute{
				Computed: true,
			},
			"source_ref": schema.StringAttribute{
				Computed: true,
			},
			"source_url": schema.StringAttribute{
				Computed: true,
			},
			"third_party_bug_feature_url": schema.StringAttribute{
				Computed: true,
			},
			"third_party_documentation_url": schema.StringAttribute{
				Computed: true,
			},
			"third_party_source_url": schema.StringAttribute{
				Computed: true,
			},
			"third_party_support_url": schema.StringAttribute{
				Computed: true,
			},
			"version": schema.StringAttribute{
				Computed: true,
			},
			"version_url": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
