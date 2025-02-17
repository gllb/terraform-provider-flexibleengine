package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2Network_importBasic(t *testing.T) {
	resourceName := "flexibleengine_networking_network_v2.network_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Network_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
