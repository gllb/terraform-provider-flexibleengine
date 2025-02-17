package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/dcs/v1/availablezones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDcsAZV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsAZV1Read,
		DeprecationMessage: "this data source is deprecated. " +
			"You can use the availability zone code directly, e.g. eu-west-0b.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsAZV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV1Client, err := config.dcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating dcs key client: %s", err)
	}

	v, err := availablezones.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Dcs az : %+v", v)
	var filteredAZs []availablezones.AvailableZone
	if v.RegionID == GetRegion(d, config) {
		AZs := v.AvailableZones
		for _, newAZ := range AZs {
			if newAZ.ResourceAvailability != "true" {
				continue
			}

			name := d.Get("name").(string)
			if name != "" && newAZ.Name != name {
				continue
			}

			port := d.Get("port").(string)
			if port != "" && newAZ.Port != port {
				continue
			}
			filteredAZs = append(filteredAZs, newAZ)
		}
	}

	if len(filteredAZs) < 1 {
		return fmt.Errorf("Not found any available zones")
	}

	az := filteredAZs[0]
	log.Printf("[DEBUG] Dcs az : %+v", az)

	d.SetId(az.ID)
	d.Set("code", az.Code)
	d.Set("name", az.Name)
	d.Set("port", az.Port)

	return nil
}
