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
var _ resource.Resource = &ConnectorGCPResource{}
var _ resource.ResourceWithImportState = &ConnectorGCPResource{}

func NewConnectorGCPResource() resource.Resource {
	return &ConnectorGCPResource{}
}

// ConnectorGCPResource defines the resource implementation.
type ConnectorGCPResource struct {
	client *Client
}

// ConnectorGCPResourceModel describes the resource data model.
type ConnectorGCPResourceModel struct {
	Name                      types.String `tfsdk:"name"`
	Region                    types.String `tfsdk:"region"`
	AccountID                 types.String `tfsdk:"account_id"`
	CompanyName               types.String `tfsdk:"company"`
	ProjectID                 types.String `tfsdk:"project_id"`
	NetworkProjectID          types.String `tfsdk:"network_project_id"`
	VPC                       types.String `tfsdk:"vpc_id"`
	SubnetID                  types.String `tfsdk:"subnet_id"`
	ServiceAccountEmail       types.String `tfsdk:"service_account_email"`
	Zone                      types.String `tfsdk:"zone"`
	ClientID                  types.String `tfsdk:"client_id"`
	// Add remaining attributes based on your existing schema
	ID                        types.String `tfsdk:"id"`
}

func (r *ConnectorGCPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector_gcp"
}

func (r *ConnectorGCPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing NetApp CloudManager Connector in GCP",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the connector",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The GCP region where the connector will be created",
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
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The GCP project ID where the connector will be created",
				Required:            true,
			},
			"network_project_id": schema.StringAttribute{
				MarkdownDescription: "The GCP network project ID for the connector",
				Optional:            true,
			},
			"vpc_id": schema.StringAttribute{
				MarkdownDescription: "The GCP VPC where the connector will be created",
				Required:            true,
			},
			"subnet_id": schema.StringAttribute{
				MarkdownDescription: "The GCP subnet where the connector will be created",
				Required:            true,
			},
			"service_account_email": schema.StringAttribute{
				MarkdownDescription: "The GCP service account email for the connector",
				Required:            true,
			},
			"zone": schema.StringAttribute{
				MarkdownDescription: "The GCP zone where the connector will be created",
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

func (r *ConnectorGCPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConnectorGCPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ConnectorGCPResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating Connector GCP: %s", plan.Name.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createConnectorGCP or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorGCPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ConnectorGCPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getConnectorGCP or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorGCPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ConnectorGCPResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state ConnectorGCPResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateConnectorGCP or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorGCPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ConnectorGCPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteConnectorGCP or equivalent method
}

func (r *ConnectorGCPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}