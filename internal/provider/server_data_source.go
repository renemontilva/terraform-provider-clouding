package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ServerDataSource{}

func NewServerDataSource() datasource.DataSource {
	return &ServerDataSource{}
}

// ServerDataSource defines the data source implementation.
type ServerDataSource struct {
	client *http.Client
}

// ServerDataSourceModel describes the data source data model.
type ServerDataSourceModel struct {
	Id                    types.String          `tfsdk:"id"`
	Name                  types.String          `tfsdk:"name"`
	Hostname              types.String          `tfsdk:"hostname"`
	Vcores                types.Number          `tfsdk:"vcores"`
	RamGB                 types.Number          `tfsdk:"ram_gb"`
	Flavor                types.String          `tfsdk:"flavor"`
	VolumeSizeGB          types.Number          `tfsdk:"volume_size_gb"`
	ImageModel            ImageModel            `tfsdk:"image"`
	Status                types.String          `tfsdk:"status"`
	PowerState            types.String          `tfsdk:"power_state"`
	Features              []types.String        `tfsdk:"features"`
	CreatedAt             types.String          `tfsdk:"created_at"`
	DnsAddresses          types.String          `tfsdk:"dns_addresses"`
	PublicIP              types.String          `tfsdk:"public_ip"`
	PrivateIP             types.String          `tfsdk:"private_ip"`
	SshKeyID              types.String          `tfsdk:"ssh_key_id"`
	Firewalls             []FirewallModel       `tfsdk:"firewalls"`
	Snapshots             []SnapshotsModel      `tfsdk:"snapshots"`
	BackupsModel          []BackupsModel        `tfsdk:"backups"`
	BackupPreferenceModel BackupPreferenceModel `tfsdk:"backup_preference"`
	CostModel             CostModel             `tfsdk:"cost"`
}

type ImageModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type FirewallModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type SnapshotsModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
}

type BackupsModel struct {
	ID        types.String `tfsdk:"id"`
	CreatedAt types.String `tfsdk:"created_at"`
	Status    types.String `tfsdk:"status"`
}

type CostModel struct {
	PricePerHour        types.Number `tfsdk:"price_per_hour"`
	PricePerMonthApprox types.Number `tfsdk:"price_per_month_approx"`
}

func (d *ServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *ServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Server data source retrieves specific information about a Server and its associated rules.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A unique string identifier used to reference a Server.",
				Computed:            false,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Server name.",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The Server hostname.",
			},
			"vcores": schema.NumberAttribute{
				MarkdownDescription: "The number of virtual cores allocated for the server.",
				Computed:            true,
			},
			"ram_gb": schema.NumberAttribute{
				MarkdownDescription: "The amount of RAM in GB allocated for the server.",
				Computed:            true,
			},
			"flavor": schema.StringAttribute{
				MarkdownDescription: "The flavor of the server.",
				Computed:            true,
			},
			"volume_size_gb": schema.NumberAttribute{
				MarkdownDescription: "The size of the server's disk in GB.",
				Computed:            true,
			},
			"image": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The image of the server.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "A unique string identifier used to reference a Image.",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Image name.",
					},
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: `Enum: "Creating" "Starting" "Active" "Stopped" "Stopping" "Rebooting" "Resize" "Unarchiving" "Archived" "Archiving" "Pending" "ResettingPassword" "RestoringBackup" "RestoringSnapshot" "Deleted" "Deleting" "Error" "Unknown"` +
					"The status of the server.",
			},
			"power_state": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: `Enum: "NoState" "Running" "Paused" "Shutdown" "Crashed" "Suspended"` +
					"The power state of the server.",
			},
			"features": schema.ListAttribute{
				Computed: true,
				MarkdownDescription: "The features that are applied to the server. The possible features are:" +
					"- **AllowSmtpOut:** The Allow SMTP Out feature allows the server to send emails. This feature is disabled by default. To enable it use the Allow server SMTP out endpoint." +
					"- **AntiDDoSNetworkFilter:** The [strict Anti-DDoS filtering](https://help.clouding.io/hc/en-us/articles/6310749915036) can protect the server under constant DDoS attacks by filtering out incoming malicious traffic. This feature can only be enabled during server creation by setting the value of enableStrictAntiDDoSFiltering to true. This feature cannot be disabled." +
					"- **Backups:** Periodic backups are enabled for this server. To configure the backup strategy of the server use the [Configure server backups](https://api.clouding.io/docs#tag/Servers/operation/ConfigureServerBackups) endpoint." +
					"- **PrivateNetwork:** The server has a second network interface connected to the private network of the user that is isolated from the public internet. To enable this feature use the [Enable server private network](https://api.clouding.io/docs#tag/Servers/operation/EnableServerPrivateNetwork) endpoint.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the server was created.",
			},
			"dns_addresses": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The DNS addresses of the server. The DNS address points to the public IP of the server.",
			},
			"public_ip": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The public IP of the server.",
			},
			"private_ip": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The private IP of the server. The private IP is only available if the server has the PrivateNetwork feature enabled.",
			},
			"ssh_key_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The SSH key ID of the server.",
				Sensitive:           true,
			},
			"firewalls": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of all firewall profiles attached to this server.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A unique string identifier used to reference a Firewall.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Firewall name.",
						},
					},
				},
			},
			"snapshots": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of all snapshots generated from this server.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A unique string identifier used to reference a Snapshot.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Snapshot name.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the snapshot was created.",
						},
					},
				},
			},
			"backups": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of all backups generated from this server.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A unique string identifier used to reference a Backup.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the backup was created.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the backup.",
						},
					},
				},
			},
			"backup_preference": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The backup strategy of the server.",
				Attributes: map[string]schema.Attribute{
					"slots": schema.NumberAttribute{
						Computed:            true,
						MarkdownDescription: "The number of backups maintained.",
					},
					"frequency": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: `The frequency of the backups. The possible values are: oneDay, twoDays, threeDays, fourDays, fiveDays, sixDays, oneWeek.`,
					},
				},
			},
			"cost": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The cost of the server.",
				Attributes: map[string]schema.Attribute{
					"price_per_hour": schema.NumberAttribute{
						Computed:            true,
						MarkdownDescription: "The hourly price of the server.",
					},
					"price_per_month_approx": schema.NumberAttribute{
						Computed:            true,
						MarkdownDescription: "The approximate monthly price of the server.",
					},
				},
			},
		},
	}
}

func (d *ServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ServerDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
