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
var _ resource.Resource = &AWSFSxResource{}
var _ resource.ResourceWithImportState = &AWSFSxResource{}

func NewAWSFSxResource() resource.Resource {
	return &AWSFSxResource{}
}

// AWSFSxResource defines the resource implementation.
type AWSFSxResource struct {
	client *Client
}

// AWSFSxResourceModel describes the resource data model.
type AWSFSxResourceModel struct {
	Name                types.String `tfsdk:"name"`
	ClientID            types.String `tfsdk:"client_id"`
	Region              types.String `tfsdk:"region"`
	WorkspaceID         types.String `tfsdk:"workspace_id"`
	PrimarySubnetID     types.String `tfsdk:"primary_subnet_id"`
	SecondarySubnetID   types.String `tfsdk:"secondary_subnet_id"`
	FQDN                types.String `tfsdk:"fqdn"`
	FilesystemID        types.String `tfsdk:"filesystem_id"`
	StorageCapacityGB   types.Int64  `tfsdk:"storage_capacity_gb"`
	ThroughputCapacity  types.Int64  `tfsdk:"throughput_capacity"`
	ThroughputCapacity2 types.Int64  `tfsdk:"throughput_capacity_2"`
	AWSCredentialsName  types.String `tfsdk:"aws_credentials_name"`
	KMSKeyID            types.String `tfsdk:"kms_key_id"`
	ID                  types.String `tfsdk:"id"`
}

func (r *AWSFSxResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_fsx"
}

func (r *AWSFSxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing AWS FSx for ONTAP in NetApp Cloud Manager",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the AWS FSx for ONTAP file system",
				Required:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID for API operations",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The AWS region where the FSx file system will be created",
				Required:            true,
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The workspace ID for the FSx file system",
				Required:            true,
			},
			"primary_subnet_id": schema.StringAttribute{
				MarkdownDescription: "The primary subnet ID for the FSx file system",
				Required:            true,
			},
			"secondary_subnet_id": schema.StringAttribute{
				MarkdownDescription: "The secondary subnet ID for the FSx file system",
				Required:            true,
			},
			"fqdn": schema.StringAttribute{
				MarkdownDescription: "The fully qualified domain name of the FSx file system",
				Computed:            true,
			},
			"filesystem_id": schema.StringAttribute{
				MarkdownDescription: "The FSx file system ID",
				Computed:            true,
			},
			"storage_capacity_gb": schema.Int64Attribute{
				MarkdownDescription: "The storage capacity of the FSx file system in GB",
				Required:            true,
			},
			"throughput_capacity": schema.Int64Attribute{
				MarkdownDescription: "The throughput capacity of the FSx file system",
				Required:            true,
			},
			"throughput_capacity_2": schema.Int64Attribute{
				MarkdownDescription: "The secondary throughput capacity of the FSx file system",
				Optional:            true,
			},
			"aws_credentials_name": schema.StringAttribute{
				MarkdownDescription: "The name of the AWS credentials to use",
				Required:            true,
			},
			"kms_key_id": schema.StringAttribute{
				MarkdownDescription: "The KMS key ID for encryption",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the FSx file system",
				Computed:            true,
			},
		},
	}
}

func (r *AWSFSxResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AWSFSxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan AWSFSxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the creation
	log.Printf("Creating AWS FSx for ONTAP file system: %s", plan.Name.ValueString())

	// Create the FSx file system request payload
	awsFSxRequest := &AWSFSxRequest{
		Name:               plan.Name.ValueString(),
		Region:             plan.Region.ValueString(),
		WorkspaceID:        plan.WorkspaceID.ValueString(),
		PrimarySubnetID:    plan.PrimarySubnetID.ValueString(),
		SecondarySubnetID:  plan.SecondarySubnetID.ValueString(),
		StorageCapacityGB:  int(plan.StorageCapacityGB.ValueInt64()),
		ThroughputCapacity: int(plan.ThroughputCapacity.ValueInt64()),
		AWSCredentialsName: plan.AWSCredentialsName.ValueString(),
	}

	// Optional fields
	if !plan.ThroughputCapacity2.IsNull() {
		awsFSxRequest.ThroughputCapacity2 = int(plan.ThroughputCapacity2.ValueInt64())
	}
	
	if !plan.KMSKeyID.IsNull() {
		awsFSxRequest.KMSKeyID = plan.KMSKeyID.ValueString()
	}

	// Call the API
	fsxResponse, err := r.client.createAWSFSx(plan.ClientID.ValueString(), awsFSxRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AWS FSx for ONTAP file system",
			"Could not create the file system: "+err.Error(),
		)
		return
	}

	// Update the state with the response
	plan.ID = types.StringValue(fsxResponse.ID)
	plan.FQDN = types.StringValue(fsxResponse.FQDN)
	plan.FilesystemID = types.StringValue(fsxResponse.FilesystemID)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *AWSFSxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state AWSFSxResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the read operation
	log.Printf("Reading AWS FSx for ONTAP file system: %s", state.ID.ValueString())

	// Call the API to get the FSx file system details
	fsxResponse, err := r.client.getAWSFSx(state.ClientID.ValueString(), state.ID.ValueString())
	if err != nil {
		// Check if the resource no longer exists
		if isResourceNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading AWS FSx for ONTAP file system",
			"Could not read the file system details: "+err.Error(),
		)
		return
	}

	// Update the state
	state.Name = types.StringValue(fsxResponse.Name)
	state.Region = types.StringValue(fsxResponse.Region)
	state.FQDN = types.StringValue(fsxResponse.FQDN)
	state.FilesystemID = types.StringValue(fsxResponse.FilesystemID)
	state.StorageCapacityGB = types.Int64Value(int64(fsxResponse.StorageCapacityGB))
	state.ThroughputCapacity = types.Int64Value(int64(fsxResponse.ThroughputCapacity))
	
	if fsxResponse.ThroughputCapacity2 > 0 {
		state.ThroughputCapacity2 = types.Int64Value(int64(fsxResponse.ThroughputCapacity2))
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *AWSFSxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan AWSFSxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state AWSFSxResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the update
	log.Printf("Updating AWS FSx for ONTAP file system: %s", state.ID.ValueString())

	// Create the update request
	updateRequest := &AWSFSxUpdateRequest{
		Name:               plan.Name.ValueString(),
		StorageCapacityGB:  int(plan.StorageCapacityGB.ValueInt64()),
		ThroughputCapacity: int(plan.ThroughputCapacity.ValueInt64()),
	}

	// Optional fields
	if !plan.ThroughputCapacity2.IsNull() {
		updateRequest.ThroughputCapacity2 = int(plan.ThroughputCapacity2.ValueInt64())
	}

	// Call the API
	_, err := r.client.updateAWSFSx(plan.ClientID.ValueString(), state.ID.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating AWS FSx for ONTAP file system",
			"Could not update the file system: "+err.Error(),
		)
		return
	}

	// Update the state with the plan
	state.Name = plan.Name
	state.StorageCapacityGB = plan.StorageCapacityGB
	state.ThroughputCapacity = plan.ThroughputCapacity
	state.ThroughputCapacity2 = plan.ThroughputCapacity2

	// Set updated state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *AWSFSxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state AWSFSxResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log the delete operation
	log.Printf("Deleting AWS FSx for ONTAP file system: %s", state.ID.ValueString())

	// Call the API to delete the FSx file system
	err := r.client.deleteAWSFSx(state.ClientID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting AWS FSx for ONTAP file system",
			"Could not delete the file system: "+err.Error(),
		)
		return
	}
}

func (r *AWSFSxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import using the resource ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}