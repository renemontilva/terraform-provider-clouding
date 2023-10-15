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
var _ datasource.DataSource = &SnapshotDataSource{}

func NewSnapshotDataSource() datasource.DataSource {
	return &SnapshotDataSource{}
}

// SnapshotDataSource defines the data source implementation.
type SnapshotDataSource struct {
	client *clouding.API
}

// SnapshotDataSourceModel describes the data source data model.
type SnapshotDataSourceModel struct {
	Id               types.String        `tfsdk:"id"`
	SizeGb           types.Int64         `tfsdk:"size_gb"`
	Name             types.String        `tfsdk:"name"`
	Description      types.String        `tfsdk:"description"`
	CreatedAt        types.String        `tfsdk:"created_at"`
	SourceServerName types.String        `tfsdk:"source_server_name"`
	Image            *SnapShotImageModel `tfsdk:"image"`
	Cost             *SnapshotCostModel  `tfsdk:"cost"`
}

type SnapShotImageModel struct {
	Id            types.String               `tfsdk:"id"`
	Name          types.String               `tfsdk:"name"`
	AccessMethods SnapshotAccessMethodsModel `tfsdk:"access_methods"`
}

type SnapshotAccessMethodsModel struct {
	SshKey   types.String `tfsdk:"ssh_key"`
	Password types.String `tfsdk:"password"`
}

type SnapshotCostModel struct {
	PricePerHour        types.Float64 `tfsdk:"price_per_hour"`
	PricePerMonthApprox types.Float64 `tfsdk:"price_per_month_approx"`
}

func (d *SnapshotDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

func (d *SnapshotDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Snapshot data source retrieves information about a specific Snapshot based on its unique identifier.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A unique string identifier used to reference a Snapshot.",
				Computed:            false,
				Required:            true,
			},
			"size_gb": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: `The size of the Snapshot in gigabytes. The size of the snapshot restricts the minimum volume size of a server created from the snapshot.`,
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the Snapshot.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the Snapshot.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the Snapshot was created.",
			},
			"source_server_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the server that the Snapshot was created from.",
			},
			"image": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The image that the Snapshot was created from.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "A unique string identifier used to reference an Image.",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the Image.",
					},
					"access_methods": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The access methods of the Image.",
						Attributes: map[string]schema.Attribute{
							"ssh_key": schema.StringAttribute{
								Computed: true,
								MarkdownDescription: `Enum: "not-supported" "optional" "required" "required-with-private-key"` +
									`This is a secure way to access your server over the network. An SSH key pair consists of a public key and a private key.` +
									`When the client attempts to connect to the server, the server checks if the public key matches the private key, and if so, grants access.`,
							},
							"password": schema.StringAttribute{
								Computed: true,
								MarkdownDescription: `Enum: "not-supported" "optional" "required" "required-with-private-key"` +
									`This is a simpler way to access your server by entering a password.`,
							},
						},
					},
				},
			},
			"cost": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The cost of the Snapshot.",
				Attributes: map[string]schema.Attribute{
					"price_per_hour": schema.Float64Attribute{
						Computed:            true,
						MarkdownDescription: "The price per hour of the Snapshot.",
					},
					"price_per_month_approx": schema.Float64Attribute{
						Computed:            true,
						MarkdownDescription: "The approximate price per month of the Snapshot.",
					},
				},
			},
		},
	}
}
func (d *SnapshotDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SnapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SnapshotDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	snapshot, err := d.client.GetSnapshotID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving Snapshot",
			err.Error(),
		)
		return
	}
	// save into the Terraform state.
	state.Id = types.StringValue(snapshot.ID)
	state.Name = types.StringValue(snapshot.Name)
	state.SizeGb = types.Int64Value(snapshot.SizeGb)
	state.Description = types.StringValue(snapshot.Description)
	state.CreatedAt = types.StringValue(snapshot.CreatedAt)
	state.SourceServerName = types.StringValue(snapshot.SourceServeName)
	state.Image = &SnapShotImageModel{
		Id:   types.StringValue(snapshot.Image.ID),
		Name: types.StringValue(snapshot.Image.Name),
		AccessMethods: SnapshotAccessMethodsModel{
			SshKey:   types.StringValue(snapshot.Image.AccessMethods.SshKey),
			Password: types.StringValue(snapshot.Image.AccessMethods.Password),
		},
	}
	state.Cost = &SnapshotCostModel{
		PricePerHour:        types.Float64Value(snapshot.Cost.PricePerHour),
		PricePerMonthApprox: types.Float64Value(snapshot.Cost.PricePerMonthApprox),
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
