package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCSBSBackupPolicyV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1DataSourceID("data.flexibleengine_csbs_backup_policy_v1.csbs_policy"),
					resource.TestCheckResourceAttr("data.flexibleengine_csbs_backup_policy_v1.csbs_policy", "name", "backup-policy"),
					resource.TestCheckResourceAttr("data.flexibleengine_csbs_backup_policy_v1.csbs_policy", "status", "suspended"),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find backup data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("backup data source ID not set ")
		}

		return nil
	}
}

var testAccCSBSBackupPolicyV1DataSource_basic = fmt.Sprintf(`
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "flexibleengine_csbs_backup_policy_v1" "backup_policy_v1" {
	name  = "backup-policy"
  	resource {
      id = "${flexibleengine_compute_instance_v2.instance_1.id}"
      type = "OS::Nova::Server"
      name = "resource4"
  	}
  	scheduled_operation {
      name ="mybackup"
      enabled = true
      operation_type ="backup"
      max_backups = "2"
      trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  	}
}

data "flexibleengine_csbs_backup_policy_v1" "csbs_policy" {  
  id = "${flexibleengine_csbs_backup_policy_v1.backup_policy_v1.id}"
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID)
