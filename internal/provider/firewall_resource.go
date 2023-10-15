package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-clouding/internal/clouding"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FirewallResource{}
var _ resource.ResourceWithImportState = &FirewallResource{}

func NewFirewallResource() resource.Resource {
	return &FirewallResource{}
}

// FirewallResource defines the resource implementation.
type FirewallResource struct {
	client *clouding.API
}

// FirewallResourceModel describes the resource data model.
type FirewallResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *FirewallResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall"
}

func (r *FirewallResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language Firewall.
		MarkdownDescription: "Create a Firewall with empty rules and attachments.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A unique string identifier used to reference a Firewall.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Firewall display name. This name is displayed in the UI.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Firewall description. The description is displayed in the UI.",
			},
			"last_updated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Firewall datetime update",
			},
		},
	}
}

func (r *FirewallResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
	firewall := clouding.Firewall{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}
	err := r.client.CreateFirewall(&firewall)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create firewall resource, got error: %s", err))
		return
	}

	// save into the Terraform state.
	plan.Id = types.StringValue(firewall.ID)
	plan.Name = types.StringValue(firewall.Name)
	plan.Description = types.StringValue(firewall.Description)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Firewall resource created")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	firewall, err := r.client.GetFirewallID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Clouding Client Error", fmt.Sprintf("Unable to read firewall id, got error: %s", err))
		return
	}
	state.Name = types.StringValue(firewall.Name)
	state.Description = types.StringValue(firewall.Description)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Firewall on the Clouding API
	firewall := clouding.Firewall{
		NewName:        plan.Name.ValueString(),
		NewDescription: plan.Description.ValueString(),
	}
	err := r.client.UpdateFirewall(plan.Id.ValueString(), firewall)
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to update firewall, got error: %s", err))
		return
	}

	// Fetch the updated Firewall from the Clouding API
	firewall, err = r.client.GetFirewallID(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to read firewall id, got error: %s", err))
		return
	}

	// Update the Terraform state with the updated values
	plan.Name = types.StringValue(firewall.Name)
	plan.Description = types.StringValue(firewall.Description)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFirewall(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to delete firewall, got error: %s", err))
		return
	}
}

func (r *FirewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
