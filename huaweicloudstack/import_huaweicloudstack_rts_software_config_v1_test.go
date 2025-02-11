package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRtsSoftwareConfigV1_importBasic(t *testing.T) {
	resourceName := "huaweicloudstack_rts_software_config_v1.config_1"
	var stackName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRtsSoftwareConfigV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsSoftwareConfigV1_basic(stackName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
