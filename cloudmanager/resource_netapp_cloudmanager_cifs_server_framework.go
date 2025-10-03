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
var _ resource.Resource = &CIFSServerResource{}
var _ resource.ResourceWithImportState = &CIFSServerResource{}

func NewCIFSServerResource() resource.Resource {
	return &CIFSServerResource{}
}

// CIFSServerResource defines the resource implementation.
type CIFSServerResource struct {
	client *Client
}

// CIFSServerResourceModel describes the resource data model.
type CIFSServerResourceModel struct {
	WorkingEnvironmentID   types.String `tfsdk:"working_environment_id"`
	WorkingEnvironmentName types.String `tfsdk:"working_environment_name"`
	ClientID               types.String `tfsdk:"client_id"`
	DomainName             types.String `tfsdk:"domain"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	DnsIPs                 types.List   `tfsdk:"dns_domain_ips"`
	OrganizationalUnit     types.String `tfsdk:"organizational_unit"`
	NetBIOS                types.String `tfsdk:"netbios"`
	SVMName                types.String `tfsdk:"svm_name"`
	ID                     types.String `tfsdk:"id"`
}

func (r *CIFSServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cifs_server"
}

func (r *CIFSServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing NetApp CloudManager CIFS Server",
		Attributes: map[string]schema.Attribute{
			"working_environment_id": schema.StringAttribute{
				MarkdownDescription: "The working environment ID where the CIFS server will be created",
				Optional:            true,
			},
			"working_environment_name": schema.StringAttribute{
				MarkdownDescription: "The working environment name where the CIFS server will be created",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID for the CIFS server",
				Required:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain name for the CIFS server",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for the CIFS server",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for the CIFS server",
				Required:            true,
				Sensitive:           true,
			},
			"dns_domain_ips": schema.ListAttribute{
				MarkdownDescription: "The DNS domain IPs for the CIFS server",
				ElementType:         types.StringType,
				Required:            true,
			},
			"organizational_unit": schema.StringAttribute{
				MarkdownDescription: "The organizational unit for the CIFS server",
				Optional:            true,
			},
			"netbios": schema.StringAttribute{
				MarkdownDescription: "The NetBIOS name for the CIFS server",
				Optional:            true,
			},
			"svm_name": schema.StringAttribute{
				MarkdownDescription: "The SVM name for the CIFS server",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the CIFS server",
			},
		},
	}
}

func (r *CIFSServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CIFSServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CIFSServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating CIFS Server for environment: %s", plan.WorkingEnvironmentName.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createCIFSServer or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CIFSServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CIFSServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getCIFSServer or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CIFSServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// CIFS server updates may not be supported, so this might need special handling
	// Retrieve values from plan
	var plan CIFSServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state CIFSServerResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateCIFSServer or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CIFSServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state CIFSServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteCIFSServer or equivalent method
}

func (r *CIFSServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}