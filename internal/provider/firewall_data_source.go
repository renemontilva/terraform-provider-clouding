package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-clouding/internal/clouding"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FirewallDataSource{}

func NewFirewallDataSource() datasource.DataSource {
	return &FirewallDataSource{}
}

// FirewallDataSource defines the data source implementation.
type FirewallDataSource struct {
	client *clouding.API
}

// FirewallDataSourceModel describes the data source data model.
type FirewallDataSourceModel struct {
	Id          types.String              `tfsdk:"id"`
	Name        types.String              `tfsdk:"name"`
	Description types.String              `tfsdk:"description"`
	Rules       []FirewallRuleModel       `tfsdk:"rules"`
	Attachments []FirewallAttachmentModel `tfsdk:"attachments"`
}

type FirewallRuleModel struct {
	Id           types.String `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
	Protocol     types.String `tfsdk:"protocol"`
	PortRangeMin types.Int64  `tfsdk:"port_range_min"`
	PortRangeMax types.Int64  `tfsdk:"port_range_max"`
	SourceIP     types.String `tfsdk:"source_ip"`
	Enable       types.Bool   `tfsdk:"enable"`
}

type FirewallAttachmentModel struct {
	FirewallID   types.String `tfsdk:"firewall_id"`
	FirewallName types.String `tfsdk:"firewall_name"`
}

func (d *FirewallDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall"
}

func (d *FirewallDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language Firewall.
		MarkdownDescription: "Firewall data source retrieves specific information about a Firewall and its associated rules.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A unique string identifier used to reference a Firewall.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Firewall display name. This name is displayed in the UI.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The Firewall description. The description is displayed in the UI.",
				Computed:            true,
			},
			"rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A unique string identifier used to reference a Firewall rule.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Firewall rule description. The description is displayed in the UI.",
						},
						"protocol": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Firewall rule protocol. The protocol is displayed in the UI.",
						},
						"port_range_min": schema.NumberAttribute{
							Computed:            true,
							MarkdownDescription: "The Firewall rule port range min. The port range min is displayed in the UI.",
						},
						"port_range_max": schema.NumberAttribute{
							Computed:            true,
							MarkdownDescription: "The Firewall rule port range max. The port range max is displayed in the UI.",
						},
						"source_ip": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Firewall rule source IP. The source IP is displayed in the UI.",
						},
						"enable": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "The Firewall rule enable. The enable is displayed in the UI.",
						},
					},
				},
			},
			"attachments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"firewall_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A unique string identifier used to reference a Firewall.",
						},
						"firewall_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Firewall display name. This name is displayed in the UI.",
						},
					},
				},
			},
		},
	}
}

func (d *FirewallDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clouding.API)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *FirewallDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state FirewallDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.GetFirewallID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving Firewall",
			err.Error(),
		)
		return
	}

	// Set the values from the API response into the model
	var firewallRules []FirewallRuleModel
	var firewallAttachments []FirewallAttachmentModel
	for _, rule := range response.Rules {
		firewallRule := FirewallRuleModel{
			Id:           types.StringValue(rule.ID),
			Description:  types.StringValue(rule.Description),
			Protocol:     types.StringValue(rule.Protocol),
			PortRangeMin: types.Int64Value(rule.PortRangeMin),
			PortRangeMax: types.Int64Value(rule.PortRangeMax),
			SourceIP:     types.StringValue(rule.SourceIP),
			Enable:       types.BoolValue(rule.Enabled),
		}
		firewallRules = append(firewallRules, firewallRule)
	}

	for _, attachment := range response.Attachments {
		firewallAttachment := FirewallAttachmentModel{
			FirewallID:   types.StringValue(attachment.ServerID),
			FirewallName: types.StringValue(attachment.ServerName),
		}
		firewallAttachments = append(firewallAttachments, firewallAttachment)
	}

	state.Name = types.StringValue(response.Name)
	state.Description = types.StringValue(response.Description)
	state.Rules = firewallRules
	state.Attachments = firewallAttachments
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
