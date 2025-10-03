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
var _ datasource.DataSource = &NssAccountDataSource{}

// NewNssAccountDataSource is a helper function to simplify the provider implementation
func NewNssAccountDataSource() datasource.DataSource {
	return &NssAccountDataSource{}
}

// NssAccountDataSource is the data source implementation
type NssAccountDataSource struct {
	client *Client
}

// NssAccountDataSourceModel maps the data source schema data
type NssAccountDataSourceModel struct {
	ClientID    types.String `tfsdk:"client_id"`
	Username    types.String `tfsdk:"username"`
	ID          types.String `tfsdk:"id"`
}

// Metadata returns the data source type name
func (d *NssAccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nss_account"
}

// Schema defines the schema for the data source
func (d *NssAccountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches NetApp Support Site (NSS) account details",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The client ID for API operations",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username for the NSS account",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the NSS account",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *NssAccountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read fetches the NSS account from the API
func (d *NssAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get the current configuration
	var config NssAccountDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get NSS account details from the API using the client
	// This would be implemented using your existing client methods
	
	// Example: client call and populating the model with the response
	// nssAccount, err := d.client.getNssAccount(config.ClientID.ValueString(), config.Username.ValueString())
	// if err != nil {
	//     resp.Diagnostics.AddError(...)
	//     return
	// }
	
	// Set the ID
	// config.ID = types.StringValue(nssAccount.ID)

	// Set the state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}