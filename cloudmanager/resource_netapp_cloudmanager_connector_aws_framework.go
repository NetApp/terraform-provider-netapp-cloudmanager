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
var _ resource.Resource = &ConnectorAWSResource{}
var _ resource.ResourceWithImportState = &ConnectorAWSResource{}

func NewConnectorAWSResource() resource.Resource {
	return &ConnectorAWSResource{}
}

// ConnectorAWSResource defines the resource implementation.
type ConnectorAWSResource struct {
	client *Client
}

// ConnectorAWSResourceModel describes the resource data model.
type ConnectorAWSResourceModel struct {
	Name                      types.String `tfsdk:"name"`
	Region                    types.String `tfsdk:"region"`
	AccountID                 types.String `tfsdk:"account_id"`
	CompanyName               types.String `tfsdk:"company"`
	KeyName                   types.String `tfsdk:"key_name"`
	InstanceType              types.String `tfsdk:"instance_type"`
	SubnetID                  types.String `tfsdk:"subnet_id"`
	SecurityGroupID           types.String `tfsdk:"security_group_id"`
	IAMRole                   types.String `tfsdk:"iam_role"`
	ClientID                  types.String `tfsdk:"client_id"`
	// Add remaining attributes based on your existing schema
	ID                        types.String `tfsdk:"id"`
}

func (r *ConnectorAWSResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector_aws"
}

func (r *ConnectorAWSResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing NetApp CloudManager Connector in AWS",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the connector",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The AWS region where the connector will be created",
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
			"key_name": schema.StringAttribute{
				MarkdownDescription: "The AWS key pair name to use for the connector instance",
				Required:            true,
			},
			"instance_type": schema.StringAttribute{
				MarkdownDescription: "The AWS instance type for the connector",
				Required:            true,
			},
			"subnet_id": schema.StringAttribute{
				MarkdownDescription: "The AWS subnet ID for the connector",
				Required:            true,
			},
			"security_group_id": schema.StringAttribute{
				MarkdownDescription: "The AWS security group ID for the connector",
				Required:            true,
			},
			"iam_role": schema.StringAttribute{
				MarkdownDescription: "The IAM role for the connector",
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

func (r *ConnectorAWSResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConnectorAWSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ConnectorAWSResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating Connector AWS: %s", plan.Name.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createConnectorAWS or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorAWSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ConnectorAWSResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getConnectorAWS or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorAWSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ConnectorAWSResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state ConnectorAWSResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateConnectorAWS or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ConnectorAWSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ConnectorAWSResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteConnectorAWS or equivalent method
}

func (r *ConnectorAWSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}