package huaweicloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkingV2SecGroup_importBasic(t *testing.T) {
	resourceName := "huaweicloudstack_networking_secgroup_v2.secgroup_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroup_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
