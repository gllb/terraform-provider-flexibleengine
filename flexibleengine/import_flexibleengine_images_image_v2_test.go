package flexibleengine

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// PASS, but Taking 15m+ in time...
func TestAccImagesImageV2_importBasic(t *testing.T) {
	resourceName := "flexibleengine_images_image_v2.image_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region",
					"local_file_path",
					"image_cache_path",
					"image_source_url",
				},
			},
		},
	})
}
