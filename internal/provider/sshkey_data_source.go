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
var _ datasource.DataSource = &SshkeyDataSource{}

func NewSshkeyDataSource() datasource.DataSource {
	return &SshkeyDataSource{}
}

// SshkeyDataSource defines the data source implementation.
type SshkeyDataSource struct {
	client *clouding.API
}

// SshkeyDataSourceModel describes the data source data model.
type SshkeyDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	FingerPrint   types.String `tfsdk:"fingerprint"`
	PublicKey     types.String `tfsdk:"public_key"`
	PrivateKey    types.String `tfsdk:"private_key"`
	HasPrivateKey types.Bool   `tfsdk:"has_private_key"`
}

func (d *SshkeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sshkey"
}

func (d *SshkeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSH key data source retrieves information about a specific SSH key based on its unique identifier.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A unique string identifier used to reference a SSH key.",
				Computed:            false,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the SSH key.",
				Computed:            true,
			},
			"fingerprint": schema.StringAttribute{
				MarkdownDescription: "The fingerprint of the SSH key.",
				Computed:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "The public key of the SSH key.",
				Computed:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "The private key of the SSH key.",
				Computed:            true,
				Sensitive:           true,
			},
			"has_private_key": schema.BoolAttribute{
				MarkdownDescription: "Whether the SSH key has a private key or not",
				Computed:            true,
			},
		},
	}
}

func (d *SshkeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clouding.API)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *clouding.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SshkeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SshkeyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sshKey, err := d.client.GetSshKeyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get sshkey, got error: %s", err))
		return
	}

	// Set into the Terraform state.
	state.Id = types.StringValue(sshKey.ID)
	state.Name = types.StringValue(sshKey.Name)
	state.FingerPrint = types.StringValue(sshKey.Fingerprint)
	state.PublicKey = types.StringValue(sshKey.PublicKey)
	state.PrivateKey = types.StringValue(sshKey.PrivateKey)
	state.HasPrivateKey = types.BoolValue(sshKey.HasPrivateKey)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read SshKey data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
