//go:build ignore

package cloudmanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var _ datasource.DataSource = &CIFSServerDataSource{}

// NewCIFSServerDataSource is a helper function to simplify the provider implementation
func NewCIFSServerDataSource() datasource.DataSource {
	return &CIFSServerDataSource{}
}

// CIFSServerDataSource is the data source implementation
type CIFSServerDataSource struct {
	client *Client
}

// CIFSServerDataSourceModel maps the data source schema data
type CIFSServerDataSourceModel struct {
	WorkingEnvironmentID   types.String `tfsdk:"working_environment_id"`
	WorkingEnvironmentName types.String `tfsdk:"working_environment_name"`
	ClientID               types.String `tfsdk:"client_id"`
	SVMName                types.String `tfsdk:"svm_name"`
	Domain                 types.String `tfsdk:"domain"`
	DnsIPs                 types.List   `tfsdk:"dns_domain_ips"`
	NetBIOS                types.String `tfsdk:"netbios"`
	OrganizationalUnit     types.String `tfsdk:"organizational_unit"`
	ID                     types.String `tfsdk:"id"`
}

// Metadata returns the data source type name
func (d *CIFSServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cifs_server"
}

// Schema defines the schema for the data source
func (d *CIFSServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches CIFS server details",
		Attributes: map[string]schema.Attribute{
			"working_environment_id": schema.StringAttribute{
				Description: "The working environment ID where the CIFS server exists",
				Optional:    true,
			},
			"working_environment_name": schema.StringAttribute{
				Description: "The working environment name where the CIFS server exists",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID for API operations",
				Required:    true,
			},
			"svm_name": schema.StringAttribute{
				Description: "The SVM name for the CIFS server",
				Optional:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The domain of the CIFS server",
				Computed:    true,
			},
			"dns_domain_ips": schema.ListAttribute{
				Description: "The DNS domain IPs of the CIFS server",
				ElementType: types.StringType,
				Computed:    true,
			},
			"netbios": schema.StringAttribute{
				Description: "The NetBIOS name of the CIFS server",
				Computed:    true,
			},
			"organizational_unit": schema.StringAttribute{
				Description: "The organizational unit of the CIFS server",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the CIFS server",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *CIFSServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read fetches the CIFS server from the API
func (d *CIFSServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get the current configuration
	var config CIFSServerDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either working_environment_id or working_environment_name is provided
	if config.WorkingEnvironmentID.IsNull() && config.WorkingEnvironmentName.IsNull() {
		resp.Diagnostics.AddError(
			"Missing required parameter",
			"Either working_environment_id or working_environment_name must be provided",
		)
		return
	}

	// Get CIFS server details from the API using the client
	// This would be implemented using your existing client methods
	
	// Example: client call and populating the model with the response
	// cifsServer, err := d.client.getCIFSServer(...)
	// if err != nil {
	//     resp.Diagnostics.AddError(...)
	//     return
	// }
	
	// Set the computed values in the model
	// config.Domain = types.StringValue(cifsServer.Domain)
	// ...

	// Set the ID
	// config.ID = types.StringValue(...)

	// Set the state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}