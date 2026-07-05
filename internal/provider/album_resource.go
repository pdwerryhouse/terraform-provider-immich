package provider

import (
	"context"
	"fmt"
	"terraform-provider-immich/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &albumResource{}
	_ resource.ResourceWithConfigure = &albumResource{}
)

func NewAlbumResource() resource.Resource {
	return &albumResource{}
}

type albumResource struct {
	client *client.Client
}

type albumResourceModel struct {
	ID          types.String `tfsdk:"id"`
	AlbumName   types.String `tfsdk:"album_name"`
	Description types.String `tfsdk:"description"`
}

func (r *albumResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_album"
}

func (r *albumResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"album_name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
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

	album, err := r.client.CreateAlbum(plan.AlbumName.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating album",
			"Could not create album, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *albumResource) Read(_ context.Context, _ resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *albumResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *albumResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *albumResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
