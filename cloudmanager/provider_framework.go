//go:build ignore

package cloudmanager

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
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &NetAppCloudManagerFrameworkProvider{}
)

// NewFrameworkProvider returns a Terraform Plugin Framework compatible provider
func NewFrameworkProvider() provider.Provider {
	return &NetAppCloudManagerFrameworkProvider{}
}

// NetAppCloudManagerFrameworkProvider is the provider implementation.
type NetAppCloudManagerFrameworkProvider struct {
	// version is set during the build, passed as a linker flag
	version string
}

// NetAppCloudManagerFrameworkProviderModel maps provider schema data to a Go type.
type NetAppCloudManagerFrameworkProviderModel struct {
	SaSecretKey types.String `tfsdk:"sa_secret_key"`
	SaClientID  types.String `tfsdk:"sa_client_id"`
	Environment types.String `tfsdk:"environment"`
}

// Metadata returns the provider type name.
func (p *NetAppCloudManagerFrameworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netapp-cloudmanager"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *NetAppCloudManagerFrameworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage NetApp Cloud Manager resources",
		Attributes: map[string]schema.Attribute{
			"sa_secret_key": schema.StringAttribute{
				Description: "Service Account Secret Key for API operations",
				Optional:    true,
				Sensitive:   true,
			},
			"sa_client_id": schema.StringAttribute{
				Description: "Service Account Client ID for API operations",
				Optional:    true,
			},
			"environment": schema.StringAttribute{
				Description: "Environment for API endpoints (e.g., 'prod' or 'staging')",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a CloudManager API client for resource/data source configuration.
func (p *NetAppCloudManagerFrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring NetApp CloudManager client")

	// Retrieve provider data from configuration
	var config NetAppCloudManagerFrameworkProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If specification variables aren't set in the config, check for environment variables
	saSecretKey := os.Getenv("CLOUDMANAGER_SA_SECRET_KEY")
	saClientID := os.Getenv("CLOUDMANAGER_SA_CLIENT_ID")
	environment := os.Getenv("CLOUDMANAGER_ENVIRONMENT")

	if !config.SaSecretKey.IsNull() {
		saSecretKey = config.SaSecretKey.ValueString()
	}

	if !config.SaClientID.IsNull() {
		saClientID = config.SaClientID.ValueString()
	}

	if !config.Environment.IsNull() {
		environment = config.Environment.ValueString()
	}

	// Default values
	if environment == "" {
		environment = "prod"
	}

	// Check required values
	if saSecretKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("sa_secret_key"),
			"Missing Service Account Secret Key",
			"The provider cannot create the CloudManager API client without a secret key.",
		)
	}

	if saClientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("sa_client_id"),
			"Missing Service Account Client ID",
			"The provider cannot create the CloudManager API client without a client ID.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new CloudManager client using the configuration values
	client, err := NewClient(environment, saSecretKey, saClientID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create CloudManager API Client",
			"An unexpected error occurred when creating the CloudManager API client: "+
				err.Error(),
		)
		return
	}

	// Make the CloudManager client available during Resource and DataSource Configure
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured NetApp CloudManager client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *NetAppCloudManagerFrameworkProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCVOAWSResource,
		NewCVOAzureResource,
		NewCVOGCPResource,
		NewCVOOnPremResource,
		NewConnectorAWSResource,
		NewConnectorAzureResource,
		NewConnectorGCPResource,
		NewAggregateResource,
		NewCIFSServerResource,
		NewSnapMirrorResource,
		NewNssAccountResource,
		NewVolumeResource,
		NewAWSFSxResource,
		NewAWSFSxVolumeResource,
		// Add other resources as they are migrated
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *NetAppCloudManagerFrameworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCIFSServerDataSource,
		NewVolumeDataSource,
		NewNssAccountDataSource,
		// Add other data sources as they are migrated
	}
}