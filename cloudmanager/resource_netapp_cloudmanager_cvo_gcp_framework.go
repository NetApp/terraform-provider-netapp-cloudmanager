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
var _ resource.Resource = &CVOGCPResource{}
var _ resource.ResourceWithImportState = &CVOGCPResource{}

func NewCVOGCPResource() resource.Resource {
	return &CVOGCPResource{}
}

// CVOGCPResource defines the resource implementation.
type CVOGCPResource struct {
	client *Client
}

// CVOGCPResourceModel describes the resource data model.
type CVOGCPResourceModel struct {
	Name                      types.String `tfsdk:"name"`
	Region                    types.String `tfsdk:"region"`
	WorkspaceID               types.String `tfsdk:"workspace_id"`
	AccountID                 types.String `tfsdk:"account_id"`
	SubnetID                  types.String `tfsdk:"subnet_id"`
	NetworkProjectID          types.String `tfsdk:"network_project_id"`
	VPCID                     types.String `tfsdk:"vpc_id"`
	SVM                       types.String `tfsdk:"svm_name"`
	SvmPassword               types.String `tfsdk:"svm_password"`
	ClientID                  types.String `tfsdk:"client_id"`
	ProjectID                 types.String `tfsdk:"project_id"`
	// Add remaining attributes based on your existing schema
	ID                        types.String `tfsdk:"id"`
}

func (r *CVOGCPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cvo_gcp"
}

func (r *CVOGCPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing Cloud Volumes ONTAP in GCP",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the working environment",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The GCP region where the working environment will be created",
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
				MarkdownDescription: "The GCP subnet ID where the working environment will be created",
				Required:            true,
			},
			"network_project_id": schema.StringAttribute{
				MarkdownDescription: "The GCP network project ID",
				Optional:            true,
			},
			"vpc_id": schema.StringAttribute{
				MarkdownDescription: "The GCP VPC ID where the working environment will be created",
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
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The GCP project ID where the working environment will be created",
				Required:            true,
			},
			// Add remaining schema attributes based on your existing schema

			// ID attribute is automatically added to track the resource
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the working environment",
			},
		},
	}
}

func (r *CVOGCPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CVOGCPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CVOGCPResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating CVO GCP: %s", plan.Name.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createCVOGCP or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOGCPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CVOGCPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getCVOGCP or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOGCPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan CVOGCPResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state CVOGCPResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateCVOGCP or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOGCPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state CVOGCPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteCVOGCP or equivalent method
}

func (r *CVOGCPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}