//go:build ignore

package cloudmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ConnectorAzureResource{}
var _ resource.ResourceWithImportState = &ConnectorAzureResource{}

func NewConnectorAzureResource() resource.Resource {
	return &ConnectorAzureResource{}
}

// ConnectorAzureResource defines the resource implementation.
type ConnectorAzureResource struct {
	client *Client
}

// ConnectorAzureResourceModel describes the resource data model.
type ConnectorAzureResourceModel struct {
	Name                      types.String `tfsdk:"name"`
	Location                  types.String `tfsdk:"location"`
	AccountID                 types.String `tfsdk:"account_id"`
	CompanyName               types.String `tfsdk:"company"`
	ResourceGroup             types.String `tfsdk:"resource_group"`
	VnetID                    types.String `tfsdk:"vnet_id"`
	SubnetID                  types.String `tfsdk:"subnet_id"`
	NetworkSecurityGroup      types.String `tfsdk:"network_security_group_name"`
	ClientID                  types.String `tfsdk:"client_id"`
	// Add remaining attributes based on your existing schema
	ID                        types.String `tfsdk:"id"`
}

func (r *ConnectorAzureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector_azure"
}

func (r *ConnectorAzureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing NetApp CloudManager Connector in Azure",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the connector",
				Required:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "The Azure location where the connector will be created",
				Required:            true,
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "The NetApp account ID to associate with the connector",
				Required:            true,
			},
			"company": schema.StringAttribute{
				MarkdownDescription: "The company name for the connector",
				Required:            true,
			},
			"resource_group": schema.StringAttribute{
				MarkdownDescription: "The Azure resource group where the connector will be created",
				Required:            true,
			},
			"vnet_id": schema.StringAttribute{
				MarkdownDescription: "The Azure VNet ID where the connector will be created",
				Required:            true,
			},
			"subnet_id": schema.StringAttribute{
				MarkdownDescription: "The Azure subnet ID where the connector will be created",
				Required:            true,
			},
			"network_security_group_name": schema.StringAttribute{
				MarkdownDescription: "The Azure network security group name for the connector",
				Required:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID for the connector",
				Required:            true,
			},
			// Add remaining schema attributes based on your existing schema

			// ID attribute is automatically added to track the resource
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the connector",
			},
		},
	}
}

func (r *ConnectorAzureResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ConnectorAzureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ConnectorAzureResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating Connector Azure: %s", plan.Name.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createConnectorAzure or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorAzureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ConnectorAzureResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getConnectorAzure or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorAzureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ConnectorAzureResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state ConnectorAzureResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateConnectorAzure or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorAzureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ConnectorAzureResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteConnectorAzure or equivalent method
}

func (r *ConnectorAzureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}