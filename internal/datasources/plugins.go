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

package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"
)

var (
	_ datasource.DataSource              = &pluginsDataSource{}
	_ datasource.DataSourceWithConfigure = &pluginsDataSource{}
)

func NewPluginsDataSource() datasource.DataSource {
	return &pluginsDataSource{}
}

type pluginsDataSource struct {
	client *immichclient.Client
}

type pluginsDataSourceModel struct {
	Plugin []pluginsModel `tfsdk:"plugins"`
}

type pluginMethodModel struct {
	Description   types.String   `tfsdk:"description"`
	HostFunctions types.Bool     `tfsdk:"host_functions"`
	Key           types.String   `tfsdk:"key"`
	Name          types.String   `tfsdk:"name"`
	Title         types.String   `tfsdk:"title"`
	Types         []types.String `tfsdk:"types"`
	UiHints       []types.String `tfsdk:"ui_hints"`
	// XXX omit Schema for the moment
}

type pluginsModel struct {
	ID          types.String        `tfsdk:"id"`
	CreatedAt   types.String        `tfsdk:"created_at"`
	Author      types.String        `tfsdk:"author"`
	Description types.String        `tfsdk:"description"`
	Methods     []pluginMethodModel `tfsdk:"methods"`
	Name        types.String        `tfsdk:"name"`
	Title       types.String        `tfsdk:"title"`
	Version     types.String        `tfsdk:"version"`
	UpdatedAt   types.String        `tfsdk:"updated_at"`
}

func (d *pluginsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugins"
}

func (d *pluginsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state pluginsDataSourceModel

	plugins, err := d.client.SearchPlugins(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich Plugin",
			err.Error(),
		)
		return
	}

	for _, plugin := range plugins {
		pluginMethods := []pluginMethodModel{}

		for _, method := range plugin.Methods {
			var pluginMethod pluginMethodModel

			pluginMethod.Description = types.StringValue(method.Description)
			pluginMethod.HostFunctions = types.BoolValue(method.HostFunctions)
			pluginMethod.Key = types.StringValue(method.Key)
			pluginMethod.Name = types.StringValue(method.Name)
			pluginMethod.Title = types.StringValue(method.Title)

			pluginMethodTypes := []types.String{}

			for _, methodType := range method.Types {
				pluginMethodType := types.StringValue(string(methodType))
				pluginMethodTypes = append(pluginMethodTypes, pluginMethodType)
			}

			pluginMethod.Types = pluginMethodTypes

			pluginMethod.Types = pluginMethodTypes

			pluginMethodUiHints := []types.String{}

			for _, methodUiHint := range method.UiHints {
				pluginMethodUiHint := types.StringValue(string(methodUiHint))
				pluginMethodUiHints = append(pluginMethodUiHints, pluginMethodUiHint)
			}

			pluginMethod.UiHints = pluginMethodUiHints

			pluginMethods = append(pluginMethods, pluginMethod)
		}

		pluginstate := pluginsModel{
			ID:          types.StringValue(plugin.Id),
			Author:      types.StringValue(plugin.Author),
			CreatedAt:   types.StringValue(plugin.CreatedAt),
			Methods:     pluginMethods,
			Name:        types.StringValue(plugin.Name),
			Description: types.StringValue(plugin.Description),
			Title:       types.StringValue(plugin.Title),
			Version:     types.StringValue(plugin.Version),
			UpdatedAt:   types.StringValue(plugin.UpdatedAt),
		}

		state.Plugin = append(state.Plugin, pluginstate)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *pluginsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *pluginsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of plugins.",
		Attributes: map[string]schema.Attribute{
			"plugins": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Plugin Id.",
							Computed:    true,
						},
						"author": schema.StringAttribute{
							Description: "The plugin's author.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The time when the plugin was created.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The plugin's description.",
							Computed:    true,
						},
						"methods": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"description": schema.StringAttribute{
										Description: "Method description.",
										Computed:    true,
									},
									"host_functions": schema.BoolAttribute{
										Description: "Has host functions.",
										Computed:    true,
									},
									"key": schema.StringAttribute{
										Description: "Method key.",
										Computed:    true,
									},
									"name": schema.StringAttribute{
										Description: "Method name.",
										Computed:    true,
									},
									"types": schema.ListAttribute{
										Description: "The plugin's types.",
										ElementType: types.StringType,
										Computed:    true,
									},
									"title": schema.StringAttribute{
										Description: "Method title.",
										Computed:    true,
									},
									"ui_hints": schema.ListAttribute{
										Description: "The plugin's user interface hints.",
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
						},
						"name": schema.StringAttribute{
							Description: "The name of the plugin.",
							Computed:    true,
						},
						"title": schema.StringAttribute{
							Description: "The title of the plugin.",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "The time when the plugin was last updated.",
							Computed:    true,
						},
						"version": schema.StringAttribute{
							Description: "The plugin's version.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
