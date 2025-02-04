package cloudmanager

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var clientID = "vAWgtc8ZcLRshb08kybl2Uhh9W0o5ElE"
var AWSVsaName = "acctestawsvsavolume"

func TestAccVolume_basic(t *testing.T) {
	var volume volumeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckGCPVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSVsaVolume(clientID, AWSVsaName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-1", &volume),
					testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-2", &volume),
					testAccCheckVolumeExists("netapp-cloudmanager_volume.cifs-volume-1", &volume),
					testAccCheckVolumeExists("netapp-cloudmanager_volume.iscsi-volume-1", &volume),
				),
			},
			// {
			// 	Config: testAccAzureVsaVolume(clientID, AzureVsaName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-2", &volume),
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-3", &volume),
			// 	),
			// },
			// {
			// 	Config: testAccVolumeConfigCreateWithCapacityTier(clientID, AWSVSAName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-3", &volume)),
			// },
			// {
			// 	Config: testAccHaVolumeConfigCreateWithName(clientID, awsHAWorkingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.nfs-volume-4", &volume)),
			// },
			// {
			// 	Config: testAccVsaCifsVolume(clientID, workingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.cifs-volume-1", &volume)),
			// },
			// {
			// 	Config: testAccVsaIscsiVolume(clientID, workingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.iscsi-volume-1", &volume)),
			// },
			// {
			// 	Config: testAccAzureVolumeConfigCreateNfs(clientID, azureSNWorkingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.azure-nfs-volume-1", &volume)),
			// },
			// {
			// 	Config: testAccAzureVolumeConfigCreateCifs(clientID, azureSNWorkingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.azure-cifs-volume-1", &volume)),
			// },
			// {
			// 	Config: testAccAzureVolumeConfigCreateIscsi(clientID, azureSNWorkingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.azure-iscsi-volume-1", &volume)),
			// },
			// {
			// 	Config: testAccGCPVolumeConfigCreateNfs(clientID, gcpVsaWorkingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.gcp-nfs-volume-1", &volume)),
			// },
			// {
			// 	Config: testAccGCPVolumeConfigCreateCifs(clientID, gcpVsaWorkingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.gcp-cifs-volume-1", &volume)),
			// },
			// {
			// 	Config: testAccGCPVolumeConfigCreateIscsi(clientID, gcpVsaWorkingEnvironmentName),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVolumeExists("netapp-cloudmanager_volume.gcp-iscsi-volume-1", &volume)),
			// },
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
			info, err := client.findWorkingEnvironmentByName(v, rs.Primary.Attributes["client_id"], true, "")
			if err != nil {
				return err
			}
			vol.WorkingEnvironmentID = info.PublicID
		}
		response, err := client.getVolumeByID(vol, rs.Primary.Attributes["client_id"])
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
			info, err := client.getWorkingEnvironmentInfo(v, rs.Primary.Attributes["client_id"], true, "")
			if err != nil {
				return err
			}
			vol.WorkingEnvironmentType = info.WorkingEnvironmentType
		} else if v, ok := rs.Primary.Attributes["working_environment_name"]; ok {
			info, err := client.findWorkingEnvironmentByName(v, rs.Primary.Attributes["client_id"], true, "")
			if err != nil {
				return err
			}
			vol.WorkingEnvironmentID = info.PublicID
			vol.WorkingEnvironmentType = info.WorkingEnvironmentType
		}
		response, err := client.getVolumeByID(vol, rs.Primary.Attributes["client_id"])
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

func testAccAWSVsaVolume(clientID string, weName string) string {
	cvoCreation := fmt.Sprintf(`
	resource "netapp-cloudmanager_cvo_aws" "cvo-aws" {
		provider = netapp-cloudmanager
		name = "%s"
		region = "us-east-1"
		subnet_id = "subnet-1"
		vpc_id = "vpc-1"
		svm_password = "netapp1!"
		client_id = "%s"
		writing_speed_state = "NORMAL"
	  }`, weName, clientID)

	vol1 := fmt.Sprintf(`
	  resource "netapp-cloudmanager_volume" "nfs-volume-1" {
		depends_on = [netapp-cloudmanager_cifs_server.cl-cifs]
		provider = netapp-cloudmanager
		name = "acc_test_vol_1"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		export_policy_name = "export-svm_acctestawsvsavolume-acc_test_vol_1"
		export_policy_type = "custom"
		export_policy_ip = ["10.30.0.0/16"]
		export_policy_nfs_version = ["nfs3", "nfs4"]
		provider_volume_type = "gp2"
		client_id = "%s"
		working_environment_name = "%s"
	}`, clientID, weName)

	vol2 := fmt.Sprintf(`
	resource "netapp-cloudmanager_volume" "nfs-volume-2" {
		depends_on = [netapp-cloudmanager_cifs_server.cl-cifs]
		provider = netapp-cloudmanager
		name = "acc_test_vol_2"
		size = 10
		unit = "GB"
		enable_thin_provisioning = true
		enable_compression = true
		enable_deduplication = true
		snapshot_policy_name = "default"
		export_policy_name = "export-svm_acctestawsvsavolume-acc_test_vol_2"
		export_policy_type = "custom"
		export_policy_ip = ["10.30.0.0/16"]
		export_policy_nfs_version = ["nfs3", "nfs4"]
		capacity_tier = "S3"
		tiering_policy = "auto"
		provider_volume_type = "gp2"
		client_id = "%s"
		working_environment_name = "%s"
	}
	`, clientID, weName)

	cifsSetUp := fmt.Sprintf(`
		resource "netapp-cloudmanager_cifs_server" "cl-cifs" {
			depends_on = [netapp-cloudmanager_cvo_aws.cvo-aws]
			provider = netapp-cloudmanager
			domain = "test.com"
			username = "admin"
			password = "abcde"
			dns_domain = "test.com"
			ip_addresses = ["1.0.0.2"]
			organizational_unit = "CN=Computers"
			client_id = "%s"
			working_environment_name = "%s"
			netbios = "%s"
			is_workgroup = false
		}
	`, clientID, weName, weName)

	vol3 := fmt.Sprintf(`
		resource "netapp-cloudmanager_volume" "cifs-volume-1" {
			depends_on = [netapp-cloudmanager_cifs_server.cl-cifs]
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
	  `, clientID, weName)

	vol4 := fmt.Sprintf(`
		resource "netapp-cloudmanager_volume" "iscsi-volume-1" {
			depends_on = [netapp-cloudmanager_cifs_server.cl-cifs]
			provider = netapp-cloudmanager
			name = "acc_iscsi_test_vol"
			volume_protocol = "iscsi"
			size = 10
			unit = "GB"
			enable_thin_provisioning = true
			enable_compression = true
			enable_deduplication = true
			snapshot_policy_name = "default"
			igroups = ["test_igroup"]
			initiator {
			alias = "test_alias"
			iqn = "iqn.1995-08.com.example:string"
			}
			os_name = "linux"
			capacity_tier= "S3"
			tiering_policy = "auto"
			provider_volume_type = "gp2"
			client_id = "%s"
			working_environment_name = "%s"
		}
	`, clientID, weName)

	return (cvoCreation + cifsSetUp + vol1 + vol2 + vol3 + vol4)
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

// =====================================================================================
// Azure Tests

func testAccAzureHaVolume(clientID string, weName string) string {

	cvoAzureHa := fmt.Sprintf(`
		resource "netapp-cloudmanager_cvo_azure" "cvo-azure-ha" {
			provider = netapp-cloudmanager
			name = "accTestCVOAzureHa"
			location = "westus"
			subscription_id = "1"
			subnet_id = "subnet1"
			vnet_id = "Vnetid1"
			vnet_resource_group = "rg"
			cidr = "10.0.0.0/24"
			data_encryption_type = "AZURE"
			svm_password = "Netapp1!"
			ontap_version = "latest"
			use_latest_version = true
			license_type = "ha-capacity-paygo"
			client_id = "vAWgtc8ZcLRshb08kybl2Uhh9W0o5ElE"
			workspace_id = "workspace-4kAJAvNk"
			nss_account = "c9ce2a6a-400d-4012-83b3-1e4dae16e3e4"
			writing_speed_state = "NORMAL"
			cloud_provider_account =  "InstanceProfile"
			capacity_package_name =  "Professional"
			is_ha = true
		}`)

	return cvoAzureHa
}

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
