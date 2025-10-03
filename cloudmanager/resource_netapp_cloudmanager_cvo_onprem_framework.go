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
var _ resource.Resource = &CVOOnPremResource{}
var _ resource.ResourceWithImportState = &CVOOnPremResource{}

func NewCVOOnPremResource() resource.Resource {
	return &CVOOnPremResource{}
}

// CVOOnPremResource defines the resource implementation.
type CVOOnPremResource struct {
	client *Client
}

// CVOOnPremResourceModel describes the resource data model.
type CVOOnPremResourceModel struct {
	Name                      types.String `tfsdk:"name"`
	WorkspaceID               types.String `tfsdk:"workspace_id"`
	ClusterAddress            types.String `tfsdk:"cluster_address"`
	ClusterUserName           types.String `tfsdk:"cluster_user_name"`
	ClusterPassword           types.String `tfsdk:"cluster_password"`
	ClientID                  types.String `tfsdk:"client_id"`
	// Add remaining attributes based on your existing schema
	ID                        types.String `tfsdk:"id"`
}

func (r *CVOOnPremResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cvo_onprem"
}

func (r *CVOOnPremResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing Cloud Volumes ONTAP On-Premises",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the working environment",
				Required:            true,
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the workspace where the working environment will be created",
				Optional:            true,
			},
			"cluster_address": schema.StringAttribute{
				MarkdownDescription: "The cluster management LIF address",
				Required:            true,
			},
			"cluster_user_name": schema.StringAttribute{
				MarkdownDescription: "The username for the cluster",
				Required:            true,
			},
			"cluster_password": schema.StringAttribute{
				MarkdownDescription: "The password for the cluster",
				Required:            true,
				Sensitive:           true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID for the working environment",
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

func (r *CVOOnPremResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CVOOnPremResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CVOOnPremResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating CVO OnPrem: %s", plan.Name.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createCVOOnPrem or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOOnPremResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CVOOnPremResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getCVOOnPrem or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOOnPremResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan CVOOnPremResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state CVOOnPremResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateCVOOnPrem or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CVOOnPremResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state CVOOnPremResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteCVOOnPrem or equivalent method
}

func (r *CVOOnPremResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}