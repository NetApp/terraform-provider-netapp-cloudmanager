//go:build ignore

package cloudmanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var _ datasource.DataSource = &VolumeDataSource{}

// NewVolumeDataSource is a helper function to simplify the provider implementation
func NewVolumeDataSource() datasource.DataSource {
	return &VolumeDataSource{}
}

// VolumeDataSource is the data source implementation
type VolumeDataSource struct {
	client *Client
}

// VolumeDataSourceModel maps the data source schema data
type VolumeDataSourceModel struct {
	WorkingEnvironmentID   types.String `tfsdk:"working_environment_id"`
	WorkingEnvironmentName types.String `tfsdk:"working_environment_name"`
	ClientID               types.String `tfsdk:"client_id"`
	SVMName                types.String `tfsdk:"svm_name"`
	Name                   types.String `tfsdk:"name"`
	Size                   types.Float64 `tfsdk:"size"`
	Unit                   types.String `tfsdk:"unit"`
	AggregateName          types.String `tfsdk:"aggregate_name"`
	SnapshotPolicyName     types.String `tfsdk:"snapshot_policy_name"`
	EnableThinProvisioning types.Bool `tfsdk:"enable_thin_provisioning"`
	EnableCompression      types.Bool `tfsdk:"enable_compression"`
	EnableDeduplication    types.Bool `tfsdk:"enable_deduplication"`
	ExportPolicyName       types.String `tfsdk:"export_policy_name"`
	ExportPolicyType       types.String `tfsdk:"export_policy_type"`
	ExportPolicyIP         types.List `tfsdk:"export_policy_ip"`
	ExportPolicyNfsVersion types.Set `tfsdk:"export_policy_nfs_version"`
	TieringPolicy          types.String `tfsdk:"tiering_policy"`
	CapacityTier           types.String `tfsdk:"capacity_tier"`
	VolumeProtocol         types.String `tfsdk:"volume_protocol"`
	ID                     types.String `tfsdk:"id"`
}

// Metadata returns the data source type name
func (d *VolumeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Schema defines the schema for the data source
func (d *VolumeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches volume details from CloudManager",
		Attributes: map[string]schema.Attribute{
			"working_environment_id": schema.StringAttribute{
				Description: "The working environment ID where the volume exists",
				Optional:    true,
			},
			"working_environment_name": schema.StringAttribute{
				Description: "The working environment name where the volume exists",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID for API operations",
				Required:    true,
			},
			"svm_name": schema.StringAttribute{
				Description: "The SVM name for the volume",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the volume",
				Required:    true,
			},
			"size": schema.Float64Attribute{
				Description: "The size of the volume",
				Computed:    true,
			},
			"unit": schema.StringAttribute{
				Description: "The unit of the volume size (GB, TB, etc.)",
				Computed:    true,
			},
			"aggregate_name": schema.StringAttribute{
				Description: "The aggregate name for the volume",
				Computed:    true,
			},
			"snapshot_policy_name": schema.StringAttribute{
				Description: "The snapshot policy name for the volume",
				Computed:    true,
			},
			"enable_thin_provisioning": schema.BoolAttribute{
				Description: "Whether thin provisioning is enabled",
				Computed:    true,
			},
			"enable_compression": schema.BoolAttribute{
				Description: "Whether compression is enabled",
				Computed:    true,
			},
			"enable_deduplication": schema.BoolAttribute{
				Description: "Whether deduplication is enabled",
				Computed:    true,
			},
			"export_policy_name": schema.StringAttribute{
				Description: "The export policy name for the volume",
				Computed:    true,
			},
			"export_policy_type": schema.StringAttribute{
				Description: "The export policy type for the volume",
				Computed:    true,
			},
			"export_policy_ip": schema.ListAttribute{
				Description: "The export policy IPs for the volume",
				ElementType: types.StringType,
				Computed:    true,
			},
			"export_policy_nfs_version": schema.SetAttribute{
				Description: "The export policy NFS versions for the volume",
				ElementType: types.StringType,
				Computed:    true,
			},
			"tiering_policy": schema.StringAttribute{
				Description: "The tiering policy for the volume",
				Computed:    true,
			},
			"capacity_tier": schema.StringAttribute{
				Description: "The capacity tier for the volume",
				Computed:    true,
			},
			"volume_protocol": schema.StringAttribute{
				Description: "The protocol for the volume",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the volume",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *VolumeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read fetches the volume from the API
func (d *VolumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get the current configuration
	var config VolumeDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either working_environment_id or working_environment_name is provided
	if config.WorkingEnvironmentID.IsNull() && config.WorkingEnvironmentName.IsNull() {
		resp.Diagnostics.AddError(
			"Missing required parameter",
			"Either working_environment_id or working_environment_name must be provided",
		)
		return
	}

	// Get volume details from the API using the client
	// This would be implemented using your existing client methods
	
	// Example: client call and populating the model with the response
	// volume, err := d.client.getVolume(...)
	// if err != nil {
	//     resp.Diagnostics.AddError(...)
	//     return
	// }
	
	// Set the computed values in the model
	// config.Size = types.Float64Value(volume.Size.Size)
	// config.Unit = types.StringValue(volume.Size.Unit)
	// ...

	// Set the ID
	// config.ID = types.StringValue(...)

	// Set the state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}