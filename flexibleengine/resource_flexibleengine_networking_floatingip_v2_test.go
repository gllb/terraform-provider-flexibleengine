package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/servers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/floatingips"
)

func TestAccNetworkingV2FloatingIP_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIP_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists("flexibleengine_networking_floatingip_v2.fip_1", &fip),
				),
			},
		},
	})
}

func TestAccNetworkingV2FloatingIP_timeout(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIP_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists("flexibleengine_networking_floatingip_v2.fip_1", &fip),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2FloatingIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine floating IP: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_floatingip_v2" {
			continue
		}

		_, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("FloatingIP still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2FloatingIPExists(n string, kp *floatingips.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		found, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("FloatingIP not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckNetworkingV2FloatingIPBoundToCorrectIP(fip *floatingips.FloatingIP, fixed_ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fip.FixedIP != fixed_ip {
			return fmt.Errorf("Floating ip associated with wrong fixed ip")
		}

		return nil
	}
}

func testAccCheckNetworkingV2InstanceFloatingIPAttach(
	instance *servers.Server, fip *floatingips.FloatingIP) resource.TestCheckFunc {

	// When Neutron is used, the Instance sometimes does not know its floating IP until some time
	// after the attachment happened. This can be anywhere from 2-20 seconds. Because of that delay,
	// the test usually completes with failure.
	// However, the Fixed IP is known on both sides immediately, so that can be used as a bridge
	// to ensure the two are now related.
	// I think a better option is to introduce some state changing config in the actual resource.
	return func(s *terraform.State) error {
		for _, networkAddresses := range instance.Addresses {
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "fixed" && address["addr"] == fip.FixedIP {
					return nil
				}
			}
		}
		return fmt.Errorf("Floating IP %+v was not attached to instance %+v", fip, instance)
	}
}

const testAccNetworkingV2FloatingIP_basic = `
resource "flexibleengine_networking_floatingip_v2" "fip_1" {
}
`

var testAccNetworkingV2FloatingIP_fixedip_bind = fmt.Sprintf(`
resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
}

resource "flexibleengine_networking_router_interface_v2" "router_interface_1" {
  router_id = "${flexibleengine_networking_router_v2.router_1.id}"
  subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
}

resource "flexibleengine_networking_router_v2" "router_1" {
  name = "router_1"
  external_gateway = "%s"
}

resource "flexibleengine_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = "${flexibleengine_networking_subnet_v2.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.10"
  }

  fixed_ip {
    subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.20"
  }
}

resource "flexibleengine_networking_floatingip_v2" "fip_1" {
  pool = "%s"
  port_id = "${flexibleengine_networking_port_v2.port_1.id}"
  fixed_ip = "${flexibleengine_networking_port_v2.port_1.fixed_ip.1.ip_address}"
}
`, OS_EXTGW_ID, OS_POOL_NAME)

const testAccNetworkingV2FloatingIP_timeout = `
resource "flexibleengine_networking_floatingip_v2" "fip_1" {
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
