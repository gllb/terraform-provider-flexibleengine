package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sdrs/v1/protectiongroups"
)

func TestAccSdrsProtectiongroupV1_basic(t *testing.T) {
	var group protectiongroups.Group

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSdrs(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSdrsProtectiongroupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSdrsProtectiongroupV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectiongroupV1Exists("flexibleengine_sdrs_protectiongroup_v1.group_1", &group),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_protectiongroup_v1.group_1", "name", "group_1"),
				),
			},
			{
				Config: testAccSdrsProtectiongroupV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectiongroupV1Exists("flexibleengine_sdrs_protectiongroup_v1.group_1", &group),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_protectiongroup_v1.group_1", "name", "group_updated"),
				),
			},
			{
				Config: testAccSdrsProtectiongroupV1_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSdrsProtectiongroupV1Exists("flexibleengine_sdrs_protectiongroup_v1.group_1", &group),
					resource.TestCheckResourceAttr(
						"flexibleengine_sdrs_protectiongroup_v1.group_1", "enable", "false"),
				),
			},
		},
	})
}

func testAccCheckSdrsProtectiongroupV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sdrsClient, err := config.sdrsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_sdrs_protectiongroup_v1" {
			continue
		}

		_, err := protectiongroups.Get(sdrsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SDRS protectiongroup still exists")
		}
	}

	return nil
}

func testAccCheckSdrsProtectiongroupV1Exists(n string, group *protectiongroups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sdrsClient, err := config.sdrsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine SDRS client: %s", err)
		}

		found, err := protectiongroups.Get(sdrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("SDRS protectiongroup not found")
		}

		*group = *found

		return nil
	}
}

var testAccSdrsProtectiongroupV1_basic = fmt.Sprintf(`
data "flexibleengine_sdrs_domain_v1" "domain_1" {
	name = "SDRS_HypeDomain01"
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
	name = "group_1"
	description = "test description"
	source_availability_zone = "eu-west-0a"
	target_availability_zone = "eu-west-0b"
	domain_id = data.flexibleengine_sdrs_domain_v1.domain_1.id
	source_vpc_id = "%s"
	dr_type = "migration"
}

resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "volume for replication pair"
  availability_zone = "eu-west-0a"
  size = 1
}

resource "flexibleengine_sdrs_replication_pair_v1" "replication_1" {
  name = "replication_1"
  description = "test replication pair"
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  volume_id = flexibleengine_blockstorage_volume_v2.volume_1.id
  delete_target_volume = true
}`, OS_VPC_ID)

var testAccSdrsProtectiongroupV1_update = fmt.Sprintf(`
data "flexibleengine_sdrs_domain_v1" "domain_1" {
	name = "SDRS_HypeDomain01"
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
	name = "group_updated"
	description = "test description"
	source_availability_zone = "eu-west-0a"
	target_availability_zone = "eu-west-0b"
	domain_id = data.flexibleengine_sdrs_domain_v1.domain_1.id
	source_vpc_id = "%s"
	dr_type = "migration"
	enable = true
}

resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "volume for replication pair"
  availability_zone = "eu-west-0a"
  size = 1
}

resource "flexibleengine_sdrs_replication_pair_v1" "replication_1" {
  name = "replication_1"
  description = "test replication pair"
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  volume_id = flexibleengine_blockstorage_volume_v2.volume_1.id
  delete_target_volume = true
}`, OS_VPC_ID)

var testAccSdrsProtectiongroupV1_disable = fmt.Sprintf(`
data "flexibleengine_sdrs_domain_v1" "domain_1" {
	name = "SDRS_HypeDomain01"
}

resource "flexibleengine_sdrs_protectiongroup_v1" "group_1" {
	name = "group_updated"
	description = "test description"
	source_availability_zone = "eu-west-0a"
	target_availability_zone = "eu-west-0b"
	domain_id = data.flexibleengine_sdrs_domain_v1.domain_1.id
	source_vpc_id = "%s"
	dr_type = "migration"
	enable = false
}

resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  description = "volume for replication pair"
  availability_zone = "eu-west-0a"
  size = 1
}

resource "flexibleengine_sdrs_replication_pair_v1" "replication_1" {
  name = "replication_1"
  description = "test replication pair"
  group_id = flexibleengine_sdrs_protectiongroup_v1.group_1.id
  volume_id = flexibleengine_blockstorage_volume_v2.volume_1.id
  delete_target_volume = true
}`, OS_VPC_ID)
