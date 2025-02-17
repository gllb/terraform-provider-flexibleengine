package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/volumeattach"
)

func TestAccComputeV2VolumeAttach_basic(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttach_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists("flexibleengine_compute_volume_attach_v2.va_1", &va),
				),
			},
		},
	})
}

func TestAccComputeV2VolumeAttach_device(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttach_device,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists("flexibleengine_compute_volume_attach_v2.va_1", &va),
				),
			},
		},
	})
}

func TestAccComputeV2VolumeAttach_timeout(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttach_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists("flexibleengine_compute_volume_attach_v2.va_1", &va),
				),
			},
		},
	})
}

func testAccCheckComputeV2VolumeAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_compute_volume_attach_v2" {
			continue
		}

		instanceId, volumeId, err := parseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = volumeattach.Get(computeClient, instanceId, volumeId).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2VolumeAttachExists(n string, va *volumeattach.VolumeAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine compute client: %s", err)
		}

		instanceId, volumeId, err := parseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := volumeattach.Get(computeClient, instanceId, volumeId).Extract()
		if err != nil {
			return err
		}

		if found.ServerID != instanceId || found.VolumeID != volumeId {
			return fmt.Errorf("VolumeAttach not found")
		}

		*va = *found

		return nil
	}
}

func testAccCheckComputeV2VolumeAttachDevice(
	va *volumeattach.VolumeAttachment, device string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if va.Device != device {
			return fmt.Errorf("Requested device of volume attachment (%s) does not match: %s",
				device, va.Device)
		}

		return nil
	}
}

var testAccComputeV2VolumeAttach_basic = fmt.Sprintf(`
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_compute_volume_attach_v2" "va_1" {
  instance_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
  volume_id = "${flexibleengine_blockstorage_volume_v2.volume_1.id}"
}
`, OS_NETWORK_ID)

var testAccComputeV2VolumeAttach_device = fmt.Sprintf(`
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_compute_volume_attach_v2" "va_1" {
  instance_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
  volume_id = "${flexibleengine_blockstorage_volume_v2.volume_1.id}"
  device = "/dev/vdc"
}
`, OS_NETWORK_ID)

var testAccComputeV2VolumeAttach_timeout = fmt.Sprintf(`
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "flexibleengine_compute_volume_attach_v2" "va_1" {
  instance_id = "${flexibleengine_compute_instance_v2.instance_1.id}"
  volume_id = "${flexibleengine_blockstorage_volume_v2.volume_1.id}"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_NETWORK_ID)
