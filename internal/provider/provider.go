// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-clouding/internal/clouding"
)

// Ensure CloudingProvider satisfies various provider interfaces.
var _ provider.Provider = &CloudingProvider{}

// CloudingProvider defines the provider implementation.
type CloudingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CloudingProvider{
			version: version,
		}
	}
}

// CloudingProviderModel describes the provider data model.
type CloudingProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *CloudingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "clouding"
	resp.Version = p.version
}

// Schema returns the schema for the provider, which includes the available resource types and their attributes.
func (p *CloudingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `**Disclaimer: Unofficial Terraform Provider for Clouding.io.**
			The Clouding provider is used to interact with the resources supported by Clouding API.
			The provider needs to be configured with the proper credentials before it can be used.`,
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "The token of the clouding API.",
				Optional:            true,
			},
		},
	}
}

// Configure configures the CloudingProvider with the given configuration in provider.tf.
// It validates the required fields and sets up the client configuration for data sources and resources.
func (p *CloudingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring the Clouding provider")

	var config CloudingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Clouding API token",
			"The provider requires a valid Clouding API token, there is an unknown value for Clouding API."+
				"Either in the provider configuration set the value statically or use CLOUDING_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Set the default value for the Clouding API with the environment variables
	token := os.Getenv("CLOUDING_TOKEN")

	// Set the value for the Clouding API with the provider configuration and override the environment variables
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the values are empty, return an error
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Empty Clouding API token",
			"The provider requires a valid Clouding API token, there is an empty value for Clouding API."+
				"Either in the provider configuration set the value statically or use CLOUDING_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "clouding_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "clouding_token")

	tflog.Debug(ctx, "Creating the Clouding provider")

	// Example client configuration for data sources and resources
	client, err := clouding.NewAPI(token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating the Clouding API client",
			"An unexpected error occurred while creating the Clouding API client."+
				"Please verify the Clouding API token is valid and try again."+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error:"+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CloudingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFirewallResource,
		NewFirewallRuleResource,
		NewServerResource,
		NewSshKeyResource,
	}
}

func (p *CloudingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFirewallDataSource,
		NewImageDataSource,
		NewSnapshotDataSource,
		NewSshkeyDataSource,
	}
}
