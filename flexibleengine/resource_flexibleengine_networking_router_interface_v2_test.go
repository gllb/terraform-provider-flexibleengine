package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/networks"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/chnsz/golangsdk/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2RouterInterface_basic_subnet(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterface_basic_subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("flexibleengine_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("flexibleengine_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("flexibleengine_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_basic_port(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterface_basic_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("flexibleengine_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("flexibleengine_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2PortExists("flexibleengine_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2RouterInterfaceExists("flexibleengine_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_timeout(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterface_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("flexibleengine_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("flexibleengine_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("flexibleengine_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("flexibleengine_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterInterfaceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "flexibleengine_networking_router_interface_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router interface still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RouterInterfaceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating FlexibleEngine networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router interface not found")
		}

		return nil
	}
}

const testAccNetworkingV2RouterInterface_basic_subnet = `
resource "flexibleengine_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_router_interface_v2" "int_1" {
  subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
  router_id = "${flexibleengine_networking_router_v2.router_1.id}"
}

resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2RouterInterface_basic_port = `
resource "flexibleengine_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_router_interface_v2" "int_1" {
  router_id = "${flexibleengine_networking_router_v2.router_1.id}"
  port_id = "${flexibleengine_networking_port_v2.port_1.id}"
}

resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
}

resource "flexibleengine_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.1"
  }
}
`

const testAccNetworkingV2RouterInterface_timeout = `
resource "flexibleengine_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_router_interface_v2" "int_1" {
  subnet_id = "${flexibleengine_networking_subnet_v2.subnet_1.id}"
  router_id = "${flexibleengine_networking_router_v2.router_1.id}"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}

resource "flexibleengine_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "flexibleengine_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${flexibleengine_networking_network_v2.network_1.id}"
}
`
