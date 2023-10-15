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
var _ datasource.DataSource = &ImageDataSource{}

func NewImageDataSource() datasource.DataSource {
	return &ImageDataSource{}
}

// ImageDataSource defines the data source implementation.
type ImageDataSource struct {
	client *clouding.API
}

// ImageDataSourceModel describes the data source data model.
type ImageDataSourceModel struct {
	Id                  types.String             `tfsdk:"id"`
	Name                types.String             `tfsdk:"name"`
	MinimumSizeGb       types.Int64              `tfsdk:"minimum_size_gb"`
	AccessMethods       *ImageAccessMethodsModel `tfsdk:"access_methods"`
	PricePerHour        types.Float64            `tfsdk:"price_per_hour"`
	PricePerMonthApprox types.Float64            `tfsdk:"price_per_month_approx"`
	BillingUnit         types.String             `tfsdk:"billing_unit"`
}

type ImageAccessMethodsModel struct {
	SshKey   types.String `tfsdk:"ssh_key"`
	Password types.String `tfsdk:"password"`
}

func (d *ImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (d *ImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language Image.
		MarkdownDescription: "Image data source allows retrieving a specific image and its access methods by providing the image's unique identifier.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A unique string identifier used to reference a Image.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the image.",
				Computed:            true,
			},
			"minimum_size_gb": schema.Int64Attribute{
				MarkdownDescription: "The minimum size in gigabytes of the image.",
				Computed:            true,
			},
			"access_methods": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The available access methods for the `accessConfiguration` when [creating a new server](https://api.clouding.io/docs#tag/Servers/operation/CreateServer) from this image.",
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
			"price_per_hour": schema.Float64Attribute{
				MarkdownDescription: "The price per hour of the image.",
				Computed:            true,
			},
			"price_per_month_approx": schema.Float64Attribute{
				MarkdownDescription: "The approximate price per month of the image.",
				Computed:            true,
			},
			"billing_unit": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The unit used to bill the image, e.g. "Core" means price per server virtual core.`,
			},
		},
	}
}

func (d *ImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ImageDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	image, err := d.client.GetImageID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving Image",
			err.Error(),
		)
		return
	}

	state.Id = types.StringValue(image.ID)
	state.Name = types.StringValue(image.Name)
	state.MinimumSizeGb = types.Int64Value(image.MinimumSizeGB)
	state.AccessMethods = &ImageAccessMethodsModel{
		SshKey:   types.StringValue(image.AccessMethods.SshKey),
		Password: types.StringValue(image.AccessMethods.Password),
	}
	state.PricePerHour = types.Float64Value(image.PricePerHour)
	state.PricePerMonthApprox = types.Float64Value(image.PricePerMonthApprox)
	state.BillingUnit = types.StringValue(image.BillingUnit)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read image data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
