package cloudmanager

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFSXVolume_basic(t *testing.T) {
	fileSystemID := "fs-wji22bngfx__3_1_183_4"
	clientID := "vAWgtc8ZcLRshb08kybl2Uhh9W0o5ElE"
	var volume volumeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFSXVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFSXVolumeConfigCreate(clientID, fileSystemID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFSXVolumeExists("netapp-cloudmanager_aws_fsx_volume.nfs-volume-2", &volume),
				),
			},
		},
	})
}

func testAccCheckFSXVolumeDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range state.RootModule().Resources {
		fmt.Println(rs.Type)
		if rs.Type != "netapp-cloudmanager_aws_fsx_volume" {
			continue
		}
		var vol volumeRequest
		vol.ID = rs.Primary.ID
		vol.FileSystemID = rs.Primary.Attributes["file_system_id"]
		vol.Name = rs.Primary.Attributes["name"]
		var svm string
		if v, ok := rs.Primary.Attributes["svm_name"]; ok {
			svm = v
		} else {
			weInfo, err := client.getFSXWorkingEnvironmentInfo(rs.Primary.Attributes["tenant_id"], vol.FileSystemID, rs.Primary.Attributes["client_id"], true, "")
			if err != nil {
				return fmt.Errorf("Cannot find working environment")
			}
			svm = weInfo.SvmName

		}
		vol.SvmName = svm
		response, err := client.getVolumeByID(vol, rs.Primary.Attributes["client_id"], true, "")
		if err == nil {
			if response.ID != "" {
				return fmt.Errorf("volume (%s) still exists", response.ID)
			}
		}
	}
	return nil
}

func testAccCheckFSXVolumeExists(name string, volume *volumeResponse) resource.TestCheckFunc {
	//time.Sleep(20 * time.Second)
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
		if v, ok := rs.Primary.Attributes["file_system_id"]; ok {
			vol.FileSystemID = v
			var svm string
			if v2, ok := rs.Primary.Attributes["svm_name"]; ok {
				svm = v2
			} else {
				weInfo, err := client.getFSXWorkingEnvironmentInfo(rs.Primary.Attributes["tenant_id"], v, rs.Primary.Attributes["client_id"], true, "")
				if err != nil {
					return fmt.Errorf("Cannot find working environment")
				}
				svm = weInfo.SvmName

			}

			volume.SvmName = svm
		}
		response, err := client.getVolumeByID(vol, rs.Primary.Attributes["client_id"], true, "")
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

func testAccFSXVolumeConfigCreate(clientID string, fileSystemID string) string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_aws_fsx_volume" "nfs-volume-2" {
		provider = netapp-cloudmanager
		name = "acc_test_vol_2"
		size = 10
		unit = "GB"
		svm_name = "svm_default"
		enable_storage_efficiency = true
		snapshot_policy_name = "default"
		export_policy_type = "custom"
		export_policy_ip = ["10.30.0.0/16"]
		export_policy_nfs_version = ["nfs3", "nfs4"]
		volume_protocol = "nfs"
		client_id = "%s"
		file_system_id = "%s"
		tenant_id = "account-j3aZttuL"
	}`, clientID, fileSystemID)
}
