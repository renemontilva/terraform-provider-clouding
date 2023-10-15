package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/renemontilva/terraform-provider-clouding/internal/clouding"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SnapShotResource{}
var _ resource.ResourceWithImportState = &SnapShotResource{}

func NewSnapShotResource() resource.Resource {
	return &SnapShotResource{}
}

// SnapShotResource defines the resource implementation.
type SnapShotResource struct {
	client *clouding.API
}

// SnapShotResourceModel describes the resource data model.
type SnapShotResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	ShutDownServer types.Bool   `tfsdk:"shutdown_server"`
	LastUpdated    types.String `tfsdk:"last_updated"`
}

func (r *SnapShotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

func (r *SnapShotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language Firewall.
		MarkdownDescription: `The snapshot server endpoint allows you to create a snapshot of the volume of the server.` +
			`A server snapshot is a point-in-time copy of the current state of a server, including its data, configurations, and settings.` +
			`It is essentially a "picture" of the server at a specific moment in time, which can be used to restore the server to that state or create a new server with that state.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A unique string identifier used to reference a Snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The server id to create the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Snapshot display name. This name is displayed in the UI.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(255),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Snapshot description. The description is displayed in the UI.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(512),
				},
			},
			"shutdown_server": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Shutdown the server before creating the snapshot. This is recommended as it increases stability.",
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *SnapShotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clouding.API)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *clouding.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *SnapShotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SnapShotResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//snapshot := clouding.Snapshot{
	//	Name:           plan.Name.ValueString(),
	//	Description:    plan.Description.ValueString(),
	//	ShutDownServer: true,
	//}
}

func (r *SnapShotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *SnapShotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *SnapShotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *SnapShotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
