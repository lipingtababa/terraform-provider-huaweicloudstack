package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/networks"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2Port_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloudstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloudstack_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_noip(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_noip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloudstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloudstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 1),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_allowedAddressPairs(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrp_port_1, vrrp_port_2, instance_port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_allowedAddressPairs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloudstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("huaweicloudstack_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("huaweicloudstack_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("huaweicloudstack_networking_port_v2.instance_port", &instance_port),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_timeout(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloudstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloudstack_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2PortDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloudstack_networking_port_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Port still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2PortExists(n string, port *ports.Port) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = *found

		return nil
	}
}

func testAccCheckNetworkingV2PortCountFixedIPs(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.FixedIPs) != expected {
			return fmt.Errorf("Expected %d Fixed IPs, got %d", expected, len(port.FixedIPs))
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortCountSecurityGroups(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.SecurityGroups) != expected {
			return fmt.Errorf("Expected %d Security Groups, got %d", expected, len(port.SecurityGroups))
		}

		return nil
	}
}

const testAccNetworkingV2Port_basic = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_noip = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
  }
}
`

const testAccNetworkingV2Port_multipleNoIP = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
  }

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
  }

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
  }
}
`

const testAccNetworkingV2Port_allowedAddressPairs = `
resource "huaweicloudstack_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "huaweicloudstack_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "huaweicloudstack_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${huaweicloudstack_networking_router_v2.vrrp_router.id}"
  subnet_id = "${huaweicloudstack_networking_subnet_v2.vrrp_subnet.id}"
}

resource "huaweicloudstack_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "huaweicloudstack_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "huaweicloudstack_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${huaweicloudstack_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${huaweicloudstack_networking_port_v2.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${huaweicloudstack_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${huaweicloudstack_networking_port_v2.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingV2Port_multipleFixedIPs = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.40"
  }
}
`

const testAccNetworkingV2Port_timeout = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2Port_fixedIPs = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.24"
  }

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_updateSecurityGroups_1 = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_updateSecurityGroups_2 = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
  security_group_ids = ["${huaweicloudstack_networking_secgroup_v2.secgroup_1.id}"]

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_updateSecurityGroups_3 = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "terraform security group acceptance test"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
  security_group_ids = ["${huaweicloudstack_networking_secgroup_v2.secgroup_1.id}"]

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_updateSecurityGroups_4 = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
	security_group_ids = []

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`
