package cloudmanager

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVolume_basic(t *testing.T) {
	workingEnvironmentName := "awscluster"
	awsHAWorkingEnvironmentName := "awstestclusters"
	azureSNWorkingEnvironmentName := "azuretestcluster"
	clientID := "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
	gcpVsaWorkingEnvironmentName := "gcpcluster"
	var volume volumeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGCPVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeConfigCreateWithName(clientID, workingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-2", &volume),
				),
			},
			{
				Config: testAccVolumeConfigCreateWithCapacityTier(clientID, workingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-3", &volume)),
			},
			{
				Config: testAccHaVolumeConfigCreateWithName(clientID, awsHAWorkingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-4", &volume)),
			},
			{
				Config: testAccVsaCifsVolume(clientID, workingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.cifs-volume-1", &volume)),
			},
			{
				Config: testAccVsaIscsiVolume(clientID, workingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.iscsi-volume-1", &volume)),
			},
			{
				Config: testAccAzureVolumeConfigCreateNfs(clientID, azureSNWorkingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.azure-nfs-volume-1", &volume)),
			},
			{
				Config: testAccAzureVolumeConfigCreateCifs(clientID, azureSNWorkingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.azure-cifs-volume-1", &volume)),
			},
			{
				Config: testAccAzureVolumeConfigCreateIscsi(clientID, azureSNWorkingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.azure-iscsi-volume-1", &volume)),
			},
			{
				Config: testAccGCPVolumeConfigCreateNfs(clientID, gcpVsaWorkingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.gcp-nfs-volume-1", &volume)),
			},
			{
				Config: testAccGCPVolumeConfigCreateCifs(clientID, gcpVsaWorkingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.gcp-cifs-volume-1", &volume)),
			},
			{
				Config: testAccGCPVolumeConfigCreateIscsi(clientID, gcpVsaWorkingEnvironmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.gcp-iscsi-volume-1", &volume)),
			},
		},
	})
}

func testAccCheckGCPVolumeDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "netapp-cloudmanager_volume" {
			continue
		}
		var vol volumeRequest
		vol.ID = rs.Primary.ID
		if v, ok := rs.Primary.Attributes["working_environment_id"]; ok {
			vol.WorkingEnvironmentID = v
		} else if v, ok := rs.Primary.Attributes["working_environment_name"]; ok {
			info, err := client.findWorkingEnvironmentByName(v)
			if err != nil {
				return err
			}
			vol.WorkingEnvironmentID = info.PublicID
		}
		response, err := client.getVolumeByID(vol)
		if err == nil {
			if response.ID != "" {
				return fmt.Errorf("volume (%s) still exists", response.ID)
			}
		}
	}
	return nil
}

func testAccCheckVolumeExists(name string, volume *volumeResponse) resource.TestCheckFunc {
	time.Sleep(20 * time.Second)
	return func(s *terraform.State) error {
		// time.Sleep(20 * time.Second)
		client := testAccProvider.Meta().(*Client)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		var vol volumeRequest
		vol.ID = rs.Primary.ID
		if v, ok := rs.Primary.Attributes["working_environment_id"]; ok {
			vol.WorkingEnvironmentID = v
			info, err := client.getWorkingEnvironmentInfo(v)
			if err != nil {
				return err
			}
			vol.WorkingEnvironmentType = info.WorkingEnvironmentType
		} else if v, ok := rs.Primary.Attributes["working_environment_name"]; ok {
			info, err := client.findWorkingEnvironmentByName(v)
			if err != nil {
				return err
			}
			vol.WorkingEnvironmentID = info.PublicID
			vol.WorkingEnvironmentType = info.WorkingEnvironmentType
		}
		response, err := client.getVolumeByID(vol)
		if err != nil {
			return err
		}

		if response.ID != rs.Primary.ID {
			return fmt.Errorf("Resource ID and volume ID do not match")
		}

		*volume = response
		return nil
	}
}

func testAccVolumeConfigCreateWithName(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "nfs-volume-2" {
		provider = netapp-cloudmanager
		name = "acc_test_vol_2"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		export_policy_name = "export-svm_awscluster-acc_test_vol_2"
		export_policy_type = "custom"
		export_policy_ip = ["10.30.0.0/16"]
		export_policy_nfs_version = ["nfs3", "nfs4"]
		provider_volume_type = "gp2"
		client_id = "%s"
		working_environment_name = "%s"
	}`, clientID, workingEnvironmentName)
}

func testAccVolumeConfigCreateWithCapacityTier(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "nfs-volume-3" {
		provider = netapp-cloudmanager
		name = "acc_test_vol_3"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		export_policy_name = "export-svm_awscluster-acc_test_vol_3"
		export_policy_type = "custom"
		export_policy_ip = ["10.30.0.0/16"]
		export_policy_nfs_version = ["nfs3", "nfs4"]
		capacity_tier = "S3"
		tiering_policy = "auto"
		provider_volume_type = "gp2"
		client_id = "%s"
		working_environment_name = "%s"
	}
	`, clientID, workingEnvironmentName)
}

func testAccHaVolumeConfigCreateWithName(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "nfs-volume-4" {
		provider = netapp-cloudmanager
		name = "acc_ha_test_vol"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		export_policy_type = "custom"
		export_policy_name = "export-svm_awstestclusters-acc_ha_test_vol"
		export_policy_ip = ["10.30.0.1/16"]
		export_policy_nfs_version = ["nfs4"]
		capacity_tier = "S3"
		tiering_policy = "auto"
		provider_volume_type = "gp2"
		client_id = "%s"
		working_environment_name = "%s"
	}
	`, clientID, workingEnvironmentName)
}

func testAccVsaCifsVolume(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "cifs-volume-1" {
		provider = netapp-cloudmanager
		name = "acc_cifs_test_vol_1"
		volume_protocol = "cifs"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		share_name = "share_cifs"
		permission = "full_control"
		users = ["Everyone"]
		capacity_tier= "S3"
		tiering_policy = "auto"
		provider_volume_type = "gp2"
		client_id = "%s"
		working_environment_name = "%s"
	  }
	  `, clientID, workingEnvironmentName)
}

func testAccVsaIscsiVolume(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "iscsi-volume-1" {
		provider = netapp-cloudmanager
		name = "acc_iscsi_test_vol"
		volume_protocol = "iscsi"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = false
		enable_deduplication = false
		snapshot_policy_name = "default"
		igroups = ["test_acc_igroup"]
		os_name = "linux"
		capacity_tier= "S3"
		tiering_policy = "auto"
		provider_volume_type = "gp2"
		client_id = "%s"
		working_environment_name = "%s"
	}
	`, clientID, workingEnvironmentName)
}

// =====================================================================================
// Azure Tests

func testAccAzureVolumeConfigCreateNfs(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "azure-nfs-volume-1" {
		provider = netapp-cloudmanager
		name = "acc_test_azure_vol_1"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		export_policy_type = "custom"
		export_policy_ip = ["10.30.0.0/16"]
		export_policy_nfs_version = ["nfs3", "nfs4"]
		client_id = "%s"
		working_environment_name = "%s"
		provider_volume_type = "Premium_LRS"
	}`, clientID, workingEnvironmentName)
}

func testAccAzureVolumeConfigCreateCifs(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "azure-cifs-volume-1" {
		provider = netapp-cloudmanager
		name = "acc_test_azure_vol_2"
		size = 10
		unit = "GB"
		volume_protocol = "cifs"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		share_name = "acc_test_azure_vol_2_share"
		permission = "full_control"
		users = ["Everyone"]
		client_id = "%s"
		working_environment_name = "%s"
		provider_volume_type = "Premium_LRS"
	}`, clientID, workingEnvironmentName)
}

func testAccAzureVolumeConfigCreateIscsi(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "azure-iscsi-volume-1" {
		provider = netapp-cloudmanager
		name = "acc_test_azure_vol_3"
		volume_protocol = "iscsi"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = false
		enable_deduplication = false
		snapshot_policy_name = "default"
		igroups = ["test_acc_igroup"]
		os_name = "linux"
		provider_volume_type = "Premium_LRS"
		client_id = "%s"
		working_environment_name = "%s"
	}`, clientID, workingEnvironmentName)
}

// =====================================================================================
// GCP Tests
func testAccGCPVolumeConfigCreateNfs(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "gcp-nfs-volume-1" {
		provider = netapp-cloudmanager
		name = "acc_test_gcp_vol_1"
		volume_protocol = "nfs"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = false
		enable_deduplication = false
		export_policy_type = "custom"
		export_policy_ip = ["10.30.0.0/16"]
		export_policy_nfs_version = ["nfs3", "nfs4"]
		snapshot_policy_name = "default"
		provider_volume_type = "pd-ssd"
		client_id = "%s"
		working_environment_name = "%s"
	}`, clientID, workingEnvironmentName)
}

func testAccGCPVolumeConfigCreateCifs(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "gcp-cifs-volume-1" {
		provider = netapp-cloudmanager
		name = "acc_test_gcp_vol_2"
		size = 10
		unit = "GB"
		volume_protocol = "cifs"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		share_name = "acc_test_gcp_vol_2_share"
		permission = "full_control"
		users = ["Everyone"]
		client_id = "%s"
		working_environment_name = "%s"
		provider_volume_type = "pd-ssd"
	}`, clientID, workingEnvironmentName)
}

func testAccGCPVolumeConfigCreateIscsi(clientID string, workingEnvironmentName string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "gcp-iscsi-volume-1" {
		provider = netapp-cloudmanager
		name = "acc_test_gcp_vol_3"
		volume_protocol = "iscsi"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = false
		enable_deduplication = false
		snapshot_policy_name = "default"
		igroups = ["test_acc_igroup"]
		os_name = "linux"
		provider_volume_type = "pd-ssd"
		client_id = "%s"
		working_environment_name = "%s"
	}`, clientID, workingEnvironmentName)
}
