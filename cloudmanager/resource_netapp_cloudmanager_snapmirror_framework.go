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
var _ resource.Resource = &SnapMirrorResource{}
var _ resource.ResourceWithImportState = &SnapMirrorResource{}

func NewSnapMirrorResource() resource.Resource {
	return &SnapMirrorResource{}
}

// SnapMirrorResource defines the resource implementation.
type SnapMirrorResource struct {
	client *Client
}

// SnapMirrorResourceModel describes the resource data model.
type SnapMirrorResourceModel struct {
	SourceWorkingEnvironmentID   types.String `tfsdk:"source_working_environment_id"`
	SourceWorkingEnvironmentName types.String `tfsdk:"source_working_environment_name"`
	DestinationWorkingEnvironmentID   types.String `tfsdk:"destination_working_environment_id"`
	DestinationWorkingEnvironmentName types.String `tfsdk:"destination_working_environment_name"`
	SourceVolumeName             types.String `tfsdk:"source_volume_name"`
	DestinationVolumeName        types.String `tfsdk:"destination_volume_name"`
	PolicyName                   types.String `tfsdk:"policy_name"`
	ScheduleName                 types.String `tfsdk:"schedule_name"`
	MaxTransferRate              types.Int64  `tfsdk:"max_transfer_rate"`
	ClientID                     types.String `tfsdk:"client_id"`
	ID                           types.String `tfsdk:"id"`
}

func (r *SnapMirrorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapmirror"
}

func (r *SnapMirrorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing NetApp CloudManager SnapMirror relationships",
		Attributes: map[string]schema.Attribute{
			"source_working_environment_id": schema.StringAttribute{
				MarkdownDescription: "The source working environment ID",
				Optional:            true,
			},
			"source_working_environment_name": schema.StringAttribute{
				MarkdownDescription: "The source working environment name",
				Optional:            true,
			},
			"destination_working_environment_id": schema.StringAttribute{
				MarkdownDescription: "The destination working environment ID",
				Optional:            true,
			},
			"destination_working_environment_name": schema.StringAttribute{
				MarkdownDescription: "The destination working environment name",
				Optional:            true,
			},
			"source_volume_name": schema.StringAttribute{
				MarkdownDescription: "The source volume name",
				Required:            true,
			},
			"destination_volume_name": schema.StringAttribute{
				MarkdownDescription: "The destination volume name",
				Required:            true,
			},
			"policy_name": schema.StringAttribute{
				MarkdownDescription: "The policy name for the SnapMirror relationship",
				Optional:            true,
			},
			"schedule_name": schema.StringAttribute{
				MarkdownDescription: "The schedule name for the SnapMirror relationship",
				Optional:            true,
			},
			"max_transfer_rate": schema.Int64Attribute{
				MarkdownDescription: "The maximum transfer rate for the SnapMirror relationship in KB/s",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the SnapMirror relationship",
			},
		},
	}
}

func (r *SnapMirrorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SnapMirrorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan SnapMirrorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating SnapMirror relationship from %s to %s", plan.SourceVolumeName.ValueString(), plan.DestinationVolumeName.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createSnapMirror or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *SnapMirrorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state SnapMirrorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getSnapMirror or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *SnapMirrorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan SnapMirrorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state SnapMirrorResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateSnapMirror or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *SnapMirrorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state SnapMirrorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteSnapMirror or equivalent method
}

func (r *SnapMirrorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}