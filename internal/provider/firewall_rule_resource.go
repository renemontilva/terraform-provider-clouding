package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-clouding/internal/clouding"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FirewallRuleResource{}
var _ resource.ResourceWithImportState = &FirewallRuleResource{}

func NewFirewallRuleResource() resource.Resource {
	return &FirewallRuleResource{}
}

// FirewallRuleResource defines the resource implementation.
type FirewallRuleResource struct {
	client *clouding.API
}

// FirewallResourceModel describes the resource data model.
type FirewallRuleResourceModel struct {
	Id           types.String `tfsdk:"id"`
	FirewallID   types.String `tfsdk:"firewall_id"`
	SourceIP     types.String `tfsdk:"source_ip"`
	Protocol     types.String `tfsdk:"protocol"`
	Description  types.String `tfsdk:"description"`
	PortRangeMin types.Int64  `tfsdk:"port_range_min"`
	PortRangeMax types.Int64  `tfsdk:"port_range_max"`
}

func (r *FirewallRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (r *FirewallRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language Firewall.
		MarkdownDescription: "The create firewall rule resource allows you to create a new firewall rule to allow or block incoming traffic based on a set of conditions.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A unique string identifier used to reference a Firewall Rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"firewall_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Firewall ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The IP or CIDR that the rule will be applied for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Required: true,
				MarkdownDescription: `A firewall rule protocol is a set of rules and procedures that determine how a firewall handles network traffic.` +
					`Supported protocols are: ah,dccp,egp,esp,gre,hopopt,icmp,igmp,ip,ipip,ospf,pgm,rsvp,sctp,tcp,udp,udplite,vrrp, or any number between 0 and 255 represented as a string.`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A short description of the rule.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(512),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port_range_min": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The minimum port of the port range.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"port_range_max": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The maximum port of the port range.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *FirewallRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FirewallRuleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
	firewallRule := clouding.FirewallRuleID{
		FirewallID: plan.FirewallID.ValueString(),
		FirewallRule: clouding.FirewallRule{
			SourceIP:     plan.SourceIP.ValueString(),
			Protocol:     plan.Protocol.ValueString(),
			Description:  plan.Description.ValueString(),
			PortRangeMin: plan.PortRangeMin.ValueInt64(),
			PortRangeMax: plan.PortRangeMax.ValueInt64(),
		},
	}
	err := r.client.CreateFirewallRule(&firewallRule)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	// Save into the Terraform state.
	plan.Id = types.StringValue(firewallRule.FirewallRule.ID)
	plan.FirewallID = types.StringValue(firewallRule.FirewallID)
	plan.SourceIP = types.StringValue(firewallRule.FirewallRule.SourceIP)
	plan.Protocol = types.StringValue(firewallRule.FirewallRule.Protocol)
	plan.Description = types.StringValue(firewallRule.FirewallRule.Description)
	plan.PortRangeMin = types.Int64Value(firewallRule.FirewallRule.PortRangeMin)
	plan.PortRangeMax = types.Int64Value(firewallRule.FirewallRule.PortRangeMax)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "FirewallRule resource created")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FirewallRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	firewallRule, err := r.client.GetFirewallRule(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Clouding Client Error", fmt.Sprintf("Unable to read firewall id, got error: %s", err))
		return
	}
	// Save updated data into Terraform state
	state.FirewallID = types.StringValue(firewallRule.FirewallID)
	state.Id = types.StringValue(firewallRule.FirewallRule.ID)
	state.SourceIP = types.StringValue(firewallRule.FirewallRule.SourceIP)
	state.Protocol = types.StringValue(firewallRule.FirewallRule.Protocol)
	state.Description = types.StringValue(firewallRule.FirewallRule.Description)
	state.PortRangeMin = types.Int64Value(firewallRule.FirewallRule.PortRangeMin)
	state.PortRangeMax = types.Int64Value(firewallRule.FirewallRule.PortRangeMax)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FirewallResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update FirewallRule on the Clouding API

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FirewallRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFirewallRule(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to delete firewall rule, got error: %s", err))
		return
	}
}

func (r *FirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
