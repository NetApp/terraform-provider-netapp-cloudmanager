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
var _ resource.Resource = &CVOAWSResource{}
var _ resource.ResourceWithImportState = &CVOAWSResource{}

func NewCVOAWSResource() resource.Resource {
	return &CVOAWSResource{}
}

// CVOAWSResource defines the resource implementation.
type CVOAWSResource struct {
	client *Client
}

// CVOAWSResourceModel describes the resource data model.
type CVOAWSResourceModel struct {
	Name                        types.String `tfsdk:"name"`
	Region                      types.String `tfsdk:"region"`
	WorkspaceID                 types.String `tfsdk:"workspace_id"`
	AccountID                   types.String `tfsdk:"account_id"`
	SubnetID                    types.String `tfsdk:"subnet_id"`
	VPCID                       types.String `tfsdk:"vpc_id"`
	SVM                         types.String `tfsdk:"svm_name"`
	SvmPassword                 types.String `tfsdk:"svm_password"`
	ClientID                    types.String `tfsdk:"client_id"`
	// Add the remaining attributes based on your resource schema
	ID                          types.String `tfsdk:"id"`
}

func (r *CVOAWSResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cvo_aws"
}

func (r *CVOAWSResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing Cloud Volumes ONTAP in AWS",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the working environment",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region where the working environment will be created",
				Required:            true,
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the workspace where the working environment will be created",
				Optional:            true,
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "The NetApp account ID to associate the working environment with",
				Optional:            true,
			},
			"subnet_id": schema.StringAttribute{
				MarkdownDescription: "The AWS subnet ID where the working environment will be created",
				Required:            true,
			},
			"vpc_id": schema.StringAttribute{
				MarkdownDescription: "The AWS VPC ID where the working environment will be created",
				Required:            true,
			},
			"svm_name": schema.StringAttribute{
				MarkdownDescription: "The name of the SVM for the working environment",
				Optional:            true,
			},
			"svm_password": schema.StringAttribute{
				MarkdownDescription: "The password for the SVM",
				Optional:            true,
				Sensitive:           true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID for the working environment",
				Required:            true,
			},
			// Add the remaining schema attributes based on your existing schema

			// ID attribute is automatically added to track the resource
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the working environment",
			},
		},
	}
}

func (r *CVOAWSResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CVOAWSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CVOAWSResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating CVO AWS: %s", plan.Name.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createCVOAWS or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOAWSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CVOAWSResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getCVOAWS or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOAWSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan CVOAWSResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state CVOAWSResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateCVOAWS or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOAWSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state CVOAWSResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteCVOAWS or equivalent method
}

func (r *CVOAWSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}