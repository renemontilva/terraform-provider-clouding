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
var _ datasource.DataSource = &BackupDataSource{}

func NewBackupDataSource() datasource.DataSource {
	return &BackupDataSource{}
}

// BackupDataSource defines the data source implementation.
type BackupDataSource struct {
	client *http.Client
}

// BackupDataSourceModel describes the data source data model.
type BackupDataSourceModel struct {
	Id           types.String         `tfsdk:"id"`
	CreatedAt    types.String         `tfsdk:"created_at"`
	ServerId     types.String         `tfsdk:"server_id"`
	ServerName   types.String         `tfsdk:"server_name"`
	VolumeSizeGb types.Number         `tfsdk:"volume_size_gb"`
	Image        ImageDataSourceModel `tfsdk:"image"`
}

func (d *BackupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

func (d *BackupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Backup data source retrieves information about a specific backup based on its unique identifier.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A unique string identifier used to reference a Backup.",
				Computed:            false,
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time when the backup was created.",
				Computed:            true,
			},
			"server_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the server that the backup was created from.",
				Computed:            true,
			},
			"server_name": schema.StringAttribute{
				MarkdownDescription: "The name of the server that the backup was created from.",
				Computed:            true,
			},
			"volume_size_gb": schema.NumberAttribute{
				MarkdownDescription: "The size of the volume in gigabytes.",
				Computed:            true,
			},
			"image": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The image that the backup was created from.",
				Attributes: map[string]schema.Attribute{
					"ssh_key": schema.StringAttribute{
						Computed: true,
						MarkdownDescription: `Enum: ***not-supported*** ***optional*** ***required*** ***required-with-private-key***` +
							`This is a secure way to access your server over the network. An SSH key pair consists of a public key and a private key. When the client attempts to connect to the server, the server checks if the public key matches the private key, and if so, grants access.` +
							`- ***not-supported:*** SSH key is not supported for this image.` +
							`- ***optional:*** Some images may support both SSH key authentication and password authentication. In this case, you can choose to use either method.` +
							`- ***required:*** Some images may require SSH key authentication. This means you'll need to create an SSH key pair and provide its unique identifier when creating the server. You'll also need to have the private key stored on your client machine to access the virtual machine.` +
							`- ***required-with-private-key:*** Some images may require that you use an SSH key with the private key stored in the Clouding servers. In this case, you'll need to either generate an SSH key or provide the private key when creating it.`,
					},
					"password": schema.StringAttribute{
						Computed: true,
						MarkdownDescription: `Enum: ***not-supported*** ***optional*** ***required***` +
							`- ***not-supported:*** Some images may not support password authentication, in which case you'll need to use an SSH key to access the machine.` +
							`- ***optional:*** Some images may allow you to use either password authentication or SSH key authentication.` +
							`- ***required:*** Some images may require a password for authentication. In this case, you'll need to provide a password when creating the server.`,
					},
				},
			},
		},
	}
}

func (d *BackupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BackupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BackupDataSourceModel

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
