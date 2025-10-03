//go:build ignore

package cloudmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &VolumeResource{}
var _ resource.ResourceWithImportState = &VolumeResource{}

func NewVolumeResource() resource.Resource {
	return &VolumeResource{}
}

// VolumeResource defines the resource implementation.
type VolumeResource struct {
	client *Client
}

// VolumeResourceModel describes the resource data model.
type VolumeResourceModel struct {
	WorkingEnvironmentID   types.String `tfsdk:"working_environment_id"`
	WorkingEnvironmentName types.String `tfsdk:"working_environment_name"`
	ClientID               types.String `tfsdk:"client_id"`
	SVMName                types.String `tfsdk:"svm_name"`
	Name                   types.String `tfsdk:"name"`
	Size                   types.Float64 `tfsdk:"size"`
	Unit                   types.String `tfsdk:"unit"`
	AggregateName          types.String `tfsdk:"aggregate_name"`
	SnapshotPolicyName     types.String `tfsdk:"snapshot_policy_name"`
	EnableThinProvisioning types.Bool   `tfsdk:"enable_thin_provisioning"`
	EnableCompression      types.Bool   `tfsdk:"enable_compression"`
	EnableDeduplication    types.Bool   `tfsdk:"enable_deduplication"`
	ExportPolicyName       types.String `tfsdk:"export_policy_name"`
	ExportPolicyType       types.String `tfsdk:"export_policy_type"`
	ExportPolicyIP         types.List   `tfsdk:"export_policy_ip"`
	ExportPolicyNfsVersion types.Set    `tfsdk:"export_policy_nfs_version"`
	TieringPolicy          types.String `tfsdk:"tiering_policy"`
	CapacityTier           types.String `tfsdk:"capacity_tier"`
	VolumeProtocol         types.String `tfsdk:"volume_protocol"`
	ShareName              types.String `tfsdk:"share_name"`
	Permission             types.String `tfsdk:"permission"`
	UserName               types.String `tfsdk:"username"`
	UserID                 types.Int64  `tfsdk:"user_id"`
	GroupID                types.Int64  `tfsdk:"group_id"`
	ID                     types.String `tfsdk:"id"`
}

func (r *VolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

func (r *VolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing NetApp CloudManager Volumes",
		Attributes: map[string]schema.Attribute{
			"working_environment_id": schema.StringAttribute{
				MarkdownDescription: "The working environment ID where the volume will be created",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequireOneOf(path.Expressions{
						path.MatchRoot("working_environment_name"),
					}),
				},
			},
			"working_environment_name": schema.StringAttribute{
				MarkdownDescription: "The working environment name where the volume will be created",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequireOneOf(path.Expressions{
						path.MatchRoot("working_environment_id"),
					}),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID to use with the volume",
				Required:            true,
			},
			"svm_name": schema.StringAttribute{
				MarkdownDescription: "The name of the SVM for the volume",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the volume",
				Required:            true,
			},
			"size": schema.Float64Attribute{
				MarkdownDescription: "The size of the volume",
				Required:            true,
			},
			"unit": schema.StringAttribute{
				MarkdownDescription: "The unit of the volume size (GB, TB, etc.)",
				Required:            true,
			},
			"aggregate_name": schema.StringAttribute{
				MarkdownDescription: "The aggregate name for the volume",
				Optional:            true,
			},
			"snapshot_policy_name": schema.StringAttribute{
				MarkdownDescription: "The snapshot policy name for the volume",
				Optional:            true,
			},
			"enable_thin_provisioning": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable thin provisioning",
				Optional:            true,
			},
			"enable_compression": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable compression",
				Optional:            true,
			},
			"enable_deduplication": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable deduplication",
				Optional:            true,
			},
			"export_policy_name": schema.StringAttribute{
				MarkdownDescription: "The export policy name for the volume",
				Optional:            true,
			},
			"export_policy_type": schema.StringAttribute{
				MarkdownDescription: "The export policy type for the volume",
				Optional:            true,
			},
			"export_policy_ip": schema.ListAttribute{
				MarkdownDescription: "The export policy IPs for the volume",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"export_policy_nfs_version": schema.SetAttribute{
				MarkdownDescription: "The export policy NFS versions for the volume",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"tiering_policy": schema.StringAttribute{
				MarkdownDescription: "The tiering policy for the volume",
				Optional:            true,
			},
			"capacity_tier": schema.StringAttribute{
				MarkdownDescription: "The capacity tier for the volume",
				Optional:            true,
			},
			"volume_protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol for the volume (nfs, cifs, iscsi)",
				Optional:            true,
			},
			"share_name": schema.StringAttribute{
				MarkdownDescription: "The share name for CIFS volumes",
				Optional:            true,
			},
			"permission": schema.StringAttribute{
				MarkdownDescription: "The permission for the volume",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for the volume",
				Optional:            true,
			},
			"user_id": schema.Int64Attribute{
				MarkdownDescription: "The user ID for the volume",
				Optional:            true,
			},
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The group ID for the volume",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the volume",
				Computed:            true,
			},
		},
	}
}

func (r *VolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan VolumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creating message
	log.Printf("Creating Volume: %s", plan.Name.ValueString())

	// Implement the creation logic using the existing client methods
	// Convert plan to the format needed by the client
	// Call client.createVolume or equivalent method

	// Set plan ID to the generated ID
	// Update state with the response
	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *VolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state VolumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement read logic
	// Call client.getVolume or equivalent method
	// Update state with the response

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *VolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan VolumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state VolumeResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement update logic
	// Call client.updateVolume or equivalent method

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *VolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state VolumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Implement delete logic
	// Call client.deleteVolume or equivalent method
}

func (r *VolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import implementation
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}