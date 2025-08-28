package cloudmanager

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAggregate_basic(t *testing.T) {

	var aggregate aggregateResult
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAggregateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateConfigCreateByWorkingEnvironmentID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAggregateExists("netapp-cloudmanager_aggregate.cl-aggregate1", &aggregate),
				),
			},
			{
				Config: testAccAggregateConfigCreateByWorkingEnvironmentName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAggregateExists("netapp-cloudmanager_aggregate.cl-aggregate2", &aggregate),
				),
			},
			{
				Config: testAccAggregateConfigUpdateAggregate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAggregateExists("netapp-cloudmanager_aggregate.cl-aggregate2", &aggregate),
					resource.TestCheckResourceAttr("netapp-cloudmanager_aggregate.cl-aggregate2", "number_of_disks", "4"),
				),
			},
		},
	})
}

func TestAccAggregate_capacityIncrease(t *testing.T) {

	var aggregate aggregateResult
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAggregateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateConfigCreateForCapacityIncrease(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAggregateExists("netapp-cloudmanager_aggregate.cl-aggregate-capacity", &aggregate),
				),
			},
			{
				Config: testAccAggregateConfigIncreaseCapacity(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAggregateExists("netapp-cloudmanager_aggregate.cl-aggregate-capacity", &aggregate),
					resource.TestCheckResourceAttr("netapp-cloudmanager_aggregate.cl-aggregate-capacity", "increase_capacity_size", "512"),
					resource.TestCheckResourceAttr("netapp-cloudmanager_aggregate.cl-aggregate-capacity", "increase_capacity_unit", "GB"),
				),
			},
		},
	})
}

func testAccCheckAggregateDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "netapp-cloudmanager_aggregate" {
			continue
		}
		clientID := rs.Primary.Attributes["client_id"]
		var aggregate aggregateRequest
		id := rs.Primary.ID
		if aggr, ok := rs.Primary.Attributes["working_environment_id"]; ok {
			aggregate.WorkingEnvironmentID = aggr
		} else if name, ok := rs.Primary.Attributes["working_environment_name"]; ok {
			info, err := client.findWorkingEnvironmentByName(name, clientID, true, "")
			if err != nil {
				aggregate.WorkingEnvironmentID = info.PublicID
			}
		}

		workingEnvDetail, err := client.getWorkingEnvironmentInfo(aggregate.WorkingEnvironmentID, clientID, true, "")
		if err != nil {
			return err
		}
		response, err := client.getAggregate(aggregate, id, workingEnvDetail.WorkingEnvironmentType, clientID, true, "")
		if err == nil {
			if response.Name != "" {
				return fmt.Errorf("aggregate (%s) still exists", id)
			}
		}
	}
	return nil
}

func testAccCheckAggregateExists(name string, aggregate *aggregateResult) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No aggregate ID is set")
		}

		id := rs.Primary.ID
		aggr := aggregateRequest{}

		clientID := rs.Primary.Attributes["client_id"]
		if a, ok := rs.Primary.Attributes["working_environment_id"]; ok {
			aggr.WorkingEnvironmentID = a
		} else if a, ok := rs.Primary.Attributes["working_environment_name"]; ok {
			info, err := client.findWorkingEnvironmentByName(a, clientID, true, "")
			if err != nil {
				return err
			}
			aggr.WorkingEnvironmentID = info.PublicID
		} else {
			return fmt.Errorf("Cannot find working environment")
		}

		workingEnvDetail, err := client.getWorkingEnvironmentInfo(aggr.WorkingEnvironmentID, clientID, true, "")
		if err != nil {
			return err
		}
		response, err := client.getAggregate(aggr, id, workingEnvDetail.WorkingEnvironmentType, clientID, true, "")
		if err != nil {
			return err
		}

		if response.Name != rs.Primary.ID {
			return fmt.Errorf("Resource ID and aggregate ID do not match")
		}

		*aggregate = response

		return nil
	}
}

func testAccAggregateConfigCreateByWorkingEnvironmentID() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate1" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_1"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_id = "vsaworkingenvironment-cfmaavwc"
		number_of_disks = 1
		disk_size_size = 100
		disk_size_unit = "GB"
		capacity_tier = "NONE"
		provider_volume_type = "pd-ssd"
	}
  `
}

func testAccAggregateConfigCreateByWorkingEnvironmentName() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate2" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_2"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "acccvo"
		number_of_disks = 3
		disk_size_size = 100
		disk_size_unit = "GB"
		capacity_tier = "NONE"
		provider_volume_type = "pd-ssd"
	}
  `
}

func testAccAggregateConfigUpdateAggregate() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate2" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_2"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "acccvo"
		number_of_disks = 4
		disk_size_size = 100
		disk_size_unit = "GB"
		capacity_tier = "NONE"
		provider_volume_type = "pd-ssd"
	}
  `
}

func testAccAggregateConfigCreateForCapacityIncrease() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate-capacity" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_capacity"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "aws-test-env"
		number_of_disks = 2
		disk_size_size = 100
		disk_size_unit = "GB"
		capacity_tier = "NONE"
		provider_volume_type = "gp3"
	}
  `
}

func testAccAggregateConfigIncreaseCapacity() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate-capacity" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_capacity"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "aws-test-env"
		number_of_disks = 2
		disk_size_size = 100
		disk_size_unit = "GB"
		capacity_tier = "NONE"
		provider_volume_type = "gp3"
		increase_capacity_size = 512
		increase_capacity_unit = "GB"
	}
  `
}

func TestAccAggregate_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccAggregateConfigDiskSizeValidationFail(),
				ExpectError: regexp.MustCompile("disk_size_unit is required when disk_size_size is specified"),
			},
			{
				Config:      testAccAggregateConfigDiskUnitValidationFail(),
				ExpectError: regexp.MustCompile("disk_size_size is required when disk_size_unit is specified"),
			},
			{
				Config:      testAccAggregateConfigInitialEvSizeValidationFail(),
				ExpectError: regexp.MustCompile("initial_ev_aggregate_unit is required when initial_ev_aggregate_size is specified"),
			},
			{
				Config:      testAccAggregateConfigInitialEvUnitValidationFail(),
				ExpectError: regexp.MustCompile("initial_ev_aggregate_size is required when initial_ev_aggregate_unit is specified"),
			},
			{
				Config:      testAccAggregateConfigIncreaseSizeValidationFail(),
				ExpectError: regexp.MustCompile("increase_capacity_unit is required when increase_capacity_size is specified"),
			},
			{
				Config:      testAccAggregateConfigIncreaseUnitValidationFail(),
				ExpectError: regexp.MustCompile("increase_capacity_size is required when increase_capacity_unit is specified"),
			},
		},
	})
}

// Test configurations that should fail validation

func testAccAggregateConfigDiskSizeValidationFail() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate-fail" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_fail"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "acccvo"
		number_of_disks = 1
		disk_size_size = 100
		capacity_tier = "NONE"
		provider_volume_type = "pd-ssd"
	}
  `
}

func testAccAggregateConfigDiskUnitValidationFail() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate-fail" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_fail"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "acccvo"
		number_of_disks = 1
		disk_size_unit = "GB"
		capacity_tier = "NONE"
		provider_volume_type = "pd-ssd"
	}
  `
}

func testAccAggregateConfigInitialEvSizeValidationFail() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate-fail" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_fail"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "aws-test-env"
		number_of_disks = 1
		initial_ev_aggregate_size = 256
		capacity_tier = "NONE"
		provider_volume_type = "gp3"
	}
  `
}

func testAccAggregateConfigInitialEvUnitValidationFail() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate-fail" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_fail"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "aws-test-env"
		number_of_disks = 1
		initial_ev_aggregate_unit = "GB"
		capacity_tier = "NONE"
		provider_volume_type = "gp3"
	}
  `
}

func testAccAggregateConfigIncreaseSizeValidationFail() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate-fail" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_fail"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "aws-test-env"
		number_of_disks = 1
		increase_capacity_size = 512
		capacity_tier = "NONE"
		provider_volume_type = "gp3"
	}
  `
}

func testAccAggregateConfigIncreaseUnitValidationFail() string {
	return `
	resource "netapp-cloudmanager_aggregate" "cl-aggregate-fail" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_fail"
		client_id = "6uOCTkJr78QT51ixCGBTiLMkLglKqoU7"
		working_environment_name = "aws-test-env"
		number_of_disks = 1
		increase_capacity_unit = "GB"
		capacity_tier = "NONE"
		provider_volume_type = "gp3"
	}
  `
}
