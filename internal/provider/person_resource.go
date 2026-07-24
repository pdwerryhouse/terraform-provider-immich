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

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &personResource{}
	_ resource.ResourceWithConfigure   = &personResource{}
	_ resource.ResourceWithImportState = &personResource{}
)

func NewPersonResource() resource.Resource {
	return &personResource{}
}

type personResource struct {
	client *client.Client
}

type personResourceModel struct {
	ID         types.String `tfsdk:"id"`
	BirthDate  types.String `tfsdk:"birthdate"`
	Color      types.String `tfsdk:"color"`
	IsFavorite types.Bool   `tfsdk:"is_favorite"`
	IsHidden   types.Bool   `tfsdk:"is_hidden"`
	Name       types.String `tfsdk:"name"`
}

func (r *personResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_person"
}

func (r *personResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create an Immich Person",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"birthdate": schema.StringAttribute{
				Description: "The person's date of birth (YYYY-MM-YY)",
				Optional:    true,
			},
			"color": schema.StringAttribute{
				Description: "The colour displayed for the person, in hex notation (#RRGGBB, default is \"\").",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"is_favorite": schema.BoolAttribute{
				Description: "True if this person is a favourite.",
				Optional:    true,
			},
			"is_hidden": schema.BoolAttribute{
				Description: "True if this person is hidden.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The person's name.",
				Required:    true,
			},
		},
	}
}

func (r *personResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan personResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	person := client.PersonCreateDto{
		BirthDate:  plan.BirthDate.ValueString(),
		Color:      plan.Color.ValueString(),
		IsFavorite: plan.IsFavorite.ValueBool(),
		IsHidden:   plan.IsHidden.ValueBool(),
		Name:       plan.Name.ValueString(),
	}

	newPerson, err := r.client.CreatePerson(person)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating person",
			"Could not create person, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newPerson.Id)
	plan.BirthDate = types.StringValue(newPerson.BirthDate)
	plan.Color = types.StringValue(newPerson.Color)
	plan.IsFavorite = types.BoolValue(newPerson.IsFavorite)
	plan.IsHidden = types.BoolValue(newPerson.IsHidden)
	plan.Name = types.StringValue(newPerson.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *personResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state personResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	person, err := r.client.GetPerson(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Person",
			"Could not read Immich person ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.BirthDate = types.StringValue(person.BirthDate)
	state.Color = types.StringValue(person.Color)
	state.IsFavorite = types.BoolValue(person.IsFavorite)
	state.IsHidden = types.BoolValue(person.IsHidden)
	state.Name = types.StringValue(person.Name)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *personResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan personResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	person := client.PersonUpdateDto{
		BirthDate:  plan.BirthDate.ValueString(),
		Color:      plan.Color.ValueString(),
		IsFavorite: plan.IsFavorite.ValueBool(),
		IsHidden:   plan.IsHidden.ValueBool(),
		Name:       plan.Name.ValueString(),
	}

	updatedPerson, err := r.client.UpdatePerson(plan.ID.ValueString(), person)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich Person",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(updatedPerson.Id)
	plan.BirthDate = types.StringValue(updatedPerson.BirthDate)
	plan.Color = types.StringValue(updatedPerson.Color)
	plan.IsFavorite = types.BoolValue(updatedPerson.IsFavorite)
	plan.IsHidden = types.BoolValue(updatedPerson.IsHidden)
	plan.Name = types.StringValue(updatedPerson.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *personResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state personResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePerson(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Immich Person",
			"Could not delete person, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *personResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider developer.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *personResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
