package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-clouding/internal/clouding"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ServerResource{}
var _ resource.ResourceWithImportState = &ServerResource{}

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

// ServerResource defines the resource implementation.
type ServerResource struct {
	client *clouding.API
}

// ServerResourceModel describes the resource data model.
type ServerResourceModel struct {
	Id                            types.String              `tfsdk:"id"`
	Name                          types.String              `tfsdk:"name"`
	Hostname                      types.String              `tfsdk:"hostname"`
	FlavorID                      types.String              `tfsdk:"flavor_id"`
	FirewallID                    types.String              `tfsdk:"firewall_id"`
	AccessConfiguration           *AccessConfigurationModel `tfsdk:"access_configuration"`
	Volume                        *VolumeModel              `tfsdk:"volume"`
	EnablePrivateNetwork          types.Bool                `tfsdk:"enable_private_network"`
	EnableStrictAntiDDoSFiltering types.Bool                `tfsdk:"enable_strict_antiddos_filtering"`
	UserData                      types.String              `tfsdk:"user_data"`
	BackupPreference              *BackupPreferenceModel    `tfsdk:"backup_preference"`
	LastUpdated                   types.String              `tfsdk:"last_updated"`
	Timeouts                      timeouts.Value            `tfsdk:"timeouts"`
}

type AccessConfigurationModel struct {
	SshKeyID     types.String `tfsdk:"ssh_key_id"`
	Password     types.String `tfsdk:"password"`
	SavePassword types.Bool   `tfsdk:"save_password"`
}

type VolumeModel struct {
	Source types.String `tfsdk:"source"`
	Id     types.String `tfsdk:"id"`
	SsdGB  types.Int64  `tfsdk:"ssd_gb"`
	// ShutDownSource types.Bool   `tfsdk:"shutdown_source"`
}

type BackupPreferenceModel struct {
	Slots     types.Int64  `tfsdk:"slots"`
	Frequency types.String `tfsdk:"frequency"`
}

func (r *ServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides a range of operations for managing Clouding virtual machines.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A unique string identifier used to reference a Server.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the server.",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname of the server. It should be a valid hostname according to the [domain names RFC](https://www.rfc-editor.org/rfc/rfc1035). This value cannot be changed.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9]+$`),
						"Hostname must be lowercase alphanumeric characters only",
					),
				},
			},
			"flavor_id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the desired flavor size. Flavors are pre-defined configurations of CPU and RAM. The list of available flavors can be retrieved from the [flavor sizes](https://api.clouding.io/docs#tag/Sizes/operation/ListAllFlavors) endpoint.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"firewall_id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the initial firewall that will be attached to the server. Firewalls can be attached or detached after server creation.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"access_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "When creating a server, you need to choose a method to access it. The two options are SSH key authentication and password authentication. The availability and requirements of these methods depend on the accessMethods of the volume's source.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"ssh_key_id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier of the SSH key. The availability of this method depends on the accessMethods of the volume's source.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Default: null" +
							"The password that will be used by the new server." +
							"The availability of this method depends on the accessMethods of the volume's source.",
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"save_password": schema.BoolAttribute{
						MarkdownDescription: "Default: false" +
							"If true, the password will be stored in our database." +
							"This will enable password retrieval. If the password is not saved, you will not be able to retrieve your password. You will still be able to change the password.",
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			"volume": schema.SingleNestedAttribute{
				MarkdownDescription: "The volume configuration and origin.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"source": schema.StringAttribute{
						MarkdownDescription: "Enum: ```image``` ```backup``` ```snapshot``` ```server``` " +
							"This property is used to specify the source of the volume of the new server.",
						Required: true,
						Validators: []validator.String{
							stringvalidator.Any(
								stringvalidator.OneOf("image", "backup", "snapshot", "server"),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier of the volume's source. This property is used in conjunction with the sourceand it can be from an [image](https://api.clouding.io/docs#tag/Images/operation/ListAllImages), [backup](https://api.clouding.io/docs#tag/Backups/operation/ListAllBackups), [snapshot](https://api.clouding.io/docs#tag/Snapshots/operation/ListAllSnapshots) or [server](https://api.clouding.io/docs#tag/Servers/operation/ListAllServers).",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"ssd_gb": schema.Int64Attribute{
						MarkdownDescription: "Minimum: >=5" +
							"The size of the volume in gigabytes. The minimum size depends on the source. For example if the source is snapshot and the snapshot is 20 gigabytes, this property should be set to minimum 20 gigabytes. The list of available volume sizes can be retrieved from the volume sizes endpoint.",
						Required: true,
						Validators: []validator.Int64{
							int64validator.AtLeast(5),
						},
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
				},
			},
			"enable_private_network": schema.BoolAttribute{
				MarkdownDescription: "Default: false" +
					"If true, the server will have second network interface connected to the private network of the user that is isolated from the public internet.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"enable_strict_antiddos_filtering": schema.BoolAttribute{
				MarkdownDescription: "Default: false" +
					"If true, [strict Anti-DDoS filtering](https://help.clouding.io/hc/en-us/articles/6310749915036) will be enabled, which may impact some network protocols. It is only recommended for server under constant DDoS attacks. If your server is not under constant attacks, we recommend leaving this option disabled and rely on our standard Anti-DDoS filtering which is always enabled. This feature cannot be disabled after the server is created.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"user_data": schema.StringAttribute{
				MarkdownDescription: "Default: null" +
					"Can be used to specify scripts/commands that the server will execute during the first startup. [More information](https://help.clouding.io/hc/en-us/articles/4801240126620) ",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_preference": schema.SingleNestedAttribute{
				MarkdownDescription: "The backup strategy of the server.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"slots": schema.Int64Attribute{
						MarkdownDescription: "[2..30]" +
							"The number of backups that will be kept.",
						Optional: true,
						Validators: []validator.Int64{
							int64validator.AtLeast(2),
							int64validator.AtMost(30),
						},
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"frequency": schema.StringAttribute{
						MarkdownDescription: "Enum: ```OneDay``` ```TwoDays``` ```ThreeDays``` ```FourDays``` ```FiveDays``` ```SixDays``` ```OneWeek``` " +
							"How often backups will be created.",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.Any(
								stringvalidator.OneOf("OneDay", "TwoDays", "ThreeDays", "FourDays", "FiveDays", "SixDays", "OneWeek"),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			"last_updated": schema.StringAttribute{
				MarkdownDescription: "The datetime of the last update.",
				Computed:            true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

func (r *ServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// timeout
	createTimeout, diags := plan.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	// provider client data and make a call using it.
	var backupPreference *clouding.BackupPreference
	var volume *clouding.Volume
	var accessConfiguration *clouding.AccessConfiguration
	if plan.BackupPreference != nil {
		backupPreference = &clouding.BackupPreference{
			Slots:     plan.BackupPreference.Slots.ValueInt64(),
			Frequency: plan.BackupPreference.Frequency.ValueString(),
		}
	}
	if plan.Volume != nil {
		volume = &clouding.Volume{
			Source: plan.Volume.Source.ValueString(),
			ID:     plan.Volume.Id.ValueString(),
			SsdGb:  plan.Volume.SsdGB.ValueInt64(),
			//ShutDownSource: plan.Volume.ShutDownSource.ValueBool(),
		}
	}
	if plan.AccessConfiguration != nil {
		accessConfiguration = &clouding.AccessConfiguration{
			SshKey:       plan.AccessConfiguration.SshKeyID.ValueString(),
			Password:     plan.AccessConfiguration.Password.ValueString(),
			SavePassword: plan.AccessConfiguration.SavePassword.ValueBool(),
		}
	}

	server := clouding.Server{
		Name:                          plan.Name.ValueString(),
		Hostname:                      plan.Hostname.ValueString(),
		FlavorID:                      plan.FlavorID.ValueString(),
		FirewallID:                    plan.FirewallID.ValueString(),
		AccessConfiguration:           accessConfiguration,
		Volume:                        volume,
		EnablePrivateNetwork:          plan.EnablePrivateNetwork.ValueBool(),
		EnableStrictAntiDDoSFiltering: plan.EnableStrictAntiDDoSFiltering.ValueBool(),
		UserData:                      plan.UserData.ValueString(),
		BackupPreference:              backupPreference,
	}
	err := r.client.CreateServer(&server)
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to create server, got error: %s", err))
		return
	}

	if server.Action.ID == "" {
		resp.Diagnostics.AddError("Clouding API Error", "Server Action ID response is empty")
		return
	}

	// Wait for server action to complete
	// Check the status of the action every 5 seconds, timeout after either 20 minutes or the value defined in the timeouts attribute
	// There should be a better way to do this
	err = r.client.WaitForAction(ctx, &server.Action, 5*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to wait for server action, got error: %s", err))
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Server resource action completed at: %s", server.Action.CompletedAt))

	// Save into the Terraform state.
	plan.Id = types.StringValue(server.ID)
	plan.Name = types.StringValue(server.Name)
	plan.Hostname = types.StringValue(server.Hostname)
	plan.FlavorID = types.StringValue(server.FlavorID)
	plan.FirewallID = types.StringValue(server.FirewallID)
	if server.AccessConfiguration != nil {
		plan.AccessConfiguration = &AccessConfigurationModel{
			SshKeyID:     types.StringValue(server.AccessConfiguration.SshKeyID),
			Password:     types.StringValue(server.AccessConfiguration.Password),
			SavePassword: types.BoolValue(server.AccessConfiguration.SavePassword),
		}
	}
	if server.Volume != nil {
		plan.Volume = &VolumeModel{
			Source: types.StringValue(server.Volume.Source),
			Id:     types.StringValue(server.Volume.ID),
			SsdGB:  types.Int64Value(server.Volume.SsdGb),
		}
	}
	//plan.Volume.ShutDownSource = types.BoolValue(server.Volume.ShutDownSource)
	plan.EnablePrivateNetwork = types.BoolValue(server.EnablePrivateNetwork)
	plan.EnableStrictAntiDDoSFiltering = types.BoolValue(server.EnableStrictAntiDDoSFiltering)
	plan.UserData = types.StringValue(server.UserData)
	if server.BackupPreference != nil {
		plan.BackupPreference.Slots = types.Int64Value(server.BackupPreference.Slots)
		plan.BackupPreference.Frequency = types.StringValue(server.BackupPreference.Frequency)
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "Server resource created")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var volume *clouding.Volume
	var accessConfiguration *clouding.AccessConfiguration
	if state.Volume != nil {
		volume = &clouding.Volume{
			Source: state.Volume.Source.ValueString(),
			ID:     state.Volume.Id.ValueString(),
			SsdGb:  state.Volume.SsdGB.ValueInt64(),
		}
	}
	if state.AccessConfiguration != nil {
		accessConfiguration = &clouding.AccessConfiguration{
			Password:     state.AccessConfiguration.Password.ValueString(),
			SshKeyID:     state.AccessConfiguration.SshKeyID.ValueString(),
			SavePassword: state.AccessConfiguration.SavePassword.ValueBool(),
		}
	}

	server := clouding.Server{
		ID:                  state.Id.ValueString(),
		Volume:              volume,
		AccessConfiguration: accessConfiguration,
	}
	err := r.client.GetServerID(&server)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server, got error: %s", err))
		return
	}

	// Overwrite server into Terraform state
	state.Id = types.StringValue(server.ID)
	state.Name = types.StringValue(server.Name)
	state.Hostname = types.StringValue(server.Hostname)
	state.FlavorID = types.StringValue(server.FlavorID)
	state.FirewallID = types.StringValue(server.FirewallID)
	if server.AccessConfiguration != nil {
		state.AccessConfiguration = &AccessConfigurationModel{
			SshKeyID:     types.StringValue(server.AccessConfiguration.SshKeyID),
			Password:     types.StringValue(server.AccessConfiguration.Password),
			SavePassword: types.BoolValue(server.AccessConfiguration.SavePassword),
		}
	}
	if server.Volume != nil {
		state.Volume = &VolumeModel{
			Source: types.StringValue(server.Volume.Source),
			Id:     types.StringValue(server.Volume.ID),
			SsdGB:  types.Int64Value(server.Volume.SsdGb),
		}
	}
	state.EnablePrivateNetwork = types.BoolValue(server.EnablePrivateNetwork)
	state.EnableStrictAntiDDoSFiltering = types.BoolValue(server.EnableStrictAntiDDoSFiltering)
	state.UserData = types.StringValue(server.UserData)
	if server.BackupPreference != nil {
		state.BackupPreference = &BackupPreferenceModel{
			Slots:     types.Int64Value(server.BackupPreference.Slots),
			Frequency: types.StringValue(server.BackupPreference.Frequency),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Server on the Clouding API
	err := r.client.UpdateServerName(plan.Id.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to update server, got error: %s", err))
		return
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete Server on the Clouding API
	action, err := r.client.DeleteServer(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to delete server, got error: %s", err))
		return
	}
	err = r.client.WaitForAction(ctx, &action, 5*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Clouding API Error", fmt.Sprintf("Unable to wait for server action, got error: %s", err))
		return
	}
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
