package cloudmanager

import (
	"fmt"
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
					resource.TestCheckResourceAttr("netapp-cloudmanager_aggregate.cl-aggregate2", "number_of_disks", "2"),
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
		client.ClientID = rs.Primary.Attributes["client_id"]
		var aggregate aggregateRequest
		id := rs.Primary.ID
		if aggr, ok := rs.Primary.Attributes["working_environment_id"]; ok {
			aggregate.WorkingEnvironmentID = aggr
		} else if name, ok := rs.Primary.Attributes["working_environment_name"]; ok {
			info, err := client.findWorkingEnvironmentByName(name)
			if err != nil {
				aggregate.WorkingEnvironmentID = info.PublicID
			}
		}

		workingEnvDetail, err := client.getWorkingEnvironmentInfo(aggregate.WorkingEnvironmentID)
		if err != nil {
			return err
		}
		response, err := client.getAggregate(aggregate, id, workingEnvDetail.WorkingEnvironmentType)
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

		client.ClientID = rs.Primary.Attributes["client_id"]
		if a, ok := rs.Primary.Attributes["working_environment_id"]; ok {
			aggr.WorkingEnvironmentID = a
		} else if a, ok := rs.Primary.Attributes["working_environment_name"]; ok {
			info, err := client.findWorkingEnvironmentByName(a)
			if err != nil {
				return err
			}
			aggr.WorkingEnvironmentID = info.PublicID
		} else {
			return fmt.Errorf("Cannot find working environment")
		}

		workingEnvDetail, err := client.getWorkingEnvironmentInfo(aggr.WorkingEnvironmentID)
		if err != nil {
			return err
		}
		response, err := client.getAggregate(aggr, id, workingEnvDetail.WorkingEnvironmentType)
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
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_aggregate" "cl-aggregate1" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_1"
		provider_volume_type = "gp2"
		client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		working_environment_id = "VsaWorkingEnvironment-D72MevVC"
		number_of_disks = 1
	  	disk_size_size = 100
	  	disk_size_unit = "GB"
	}
  `)
}

func testAccAggregateConfigCreateByWorkingEnvironmentName() string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_aggregate" "cl-aggregate2" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_2"
		provider_volume_type = "gp2"
		client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		working_environment_name = "testAWS"
		number_of_disks = 1
	  	disk_size_size = 100
	  	disk_size_unit = "GB"
	}
  `)
}

func testAccAggregateConfigUpdateAggregate() string {
	return fmt.Sprintf(`
	resource "netapp-cloudmanager_aggregate" "cl-aggregate2" {
		provider = netapp-cloudmanager
		name = "acc_test_aggr_2"
		provider_volume_type = "gp2"
		client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		working_environment_name = "testAWS"
		number_of_disks = 2
	  	disk_size_size = 100
	  	disk_size_unit = "GB"
	}
  `)
}
