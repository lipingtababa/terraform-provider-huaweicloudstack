package huaweicloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkingV2SecGroupRule_importBasic(t *testing.T) {
	resourceName := "huaweicloudstack_networking_secgroup_rule_v2.secgroup_rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRule_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
