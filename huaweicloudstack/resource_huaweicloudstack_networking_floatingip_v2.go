package huaweicloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/networks"
	"github.com/huaweicloud/golangsdk/pagination"
)

func resourceNetworkingFloatingIPV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkFloatingIPV2Create,
		Read:   resourceNetworkFloatingIPV2Read,
		Update: resourceNetworkFloatingIPV2Update,
		Delete: resourceNetworkFloatingIPV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_POOL_NAME", nil),
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkFloatingIPV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack network client: %s", err)
	}

	poolID, err := getNetworkID(d, meta, d.Get("pool").(string))
	if err != nil {
		return fmt.Errorf("Error retrieving floating IP pool name: %s", err)
	}
	if len(poolID) == 0 {
		return fmt.Errorf("No network found with name: %s", d.Get("pool").(string))
	}
	createOpts := FloatingIPCreateOpts{
		floatingips.CreateOpts{
			FloatingNetworkID: poolID,
			PortID:            d.Get("port_id").(string),
			TenantID:          d.Get("tenant_id").(string),
			FixedIP:           d.Get("fixed_ip").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	floatingIP, err := floatingips.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error allocating floating IP: %s", err)
	}

	log.Printf("[DEBUG] Waiting for HuaweiCloudStack Neutron Floating IP (%s) to become available.", floatingIP.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForFloatingIPActive(networkingClient, floatingIP.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(floatingIP.ID)

	return resourceNetworkFloatingIPV2Read(d, meta)
}

func resourceNetworkFloatingIPV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack network client: %s", err)
	}

	floatingIP, err := floatingips.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "floating IP")
	}

	d.Set("address", floatingIP.FloatingIP)
	d.Set("port_id", floatingIP.PortID)
	d.Set("fixed_ip", floatingIP.FixedIP)
	poolName, err := getNetworkName(d, meta, floatingIP.FloatingNetworkID)
	if err != nil {
		return fmt.Errorf("Error retrieving floating IP pool name: %s", err)
	}
	d.Set("pool", poolName)
	d.Set("tenant_id", floatingIP.TenantID)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkFloatingIPV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack network client: %s", err)
	}

	var updateOpts floatingips.UpdateOpts

	if d.HasChange("port_id") {
		portID := d.Get("port_id").(string)
		updateOpts.PortID = &portID
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = floatingips.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating floating IP: %s", err)
	}

	return resourceNetworkFloatingIPV2Read(d, meta)
}

func resourceNetworkFloatingIPV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack network client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForFloatingIPDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloudStack Neutron Floating IP: %s", err)
	}

	d.SetId("")
	return nil
}

func getNetworkID(d *schema.ResourceData, meta interface{}, networkName string) (string, error) {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return "", fmt.Errorf("Error creating HuaweiCloudStack network client: %s", err)
	}

	opts := networks.ListOpts{Name: networkName}
	pager := networks.List(networkingClient, opts)
	networkID := ""

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}

		for _, n := range networkList {
			if n.Name == networkName {
				networkID = n.ID
				return false, nil
			}
		}

		return true, nil
	})

	return networkID, err
}

func getNetworkName(d *schema.ResourceData, meta interface{}, networkID string) (string, error) {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return "", fmt.Errorf("Error creating HuaweiCloudStack network client: %s", err)
	}

	opts := networks.ListOpts{ID: networkID}
	pager := networks.List(networkingClient, opts)
	networkName := ""

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}

		for _, n := range networkList {
			if n.ID == networkID {
				networkName = n.Name
				return false, nil
			}
		}

		return true, nil
	})

	return networkName, err
}

func waitForFloatingIPActive(networkingClient *golangsdk.ServiceClient, fId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		f, err := floatingips.Get(networkingClient, fId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloudStack Neutron Floating IP: %+v", f)
		if f.Status == "DOWN" || f.Status == "ACTIVE" {
			return f, "ACTIVE", nil
		}

		return f, "", nil
	}
}

func waitForFloatingIPDelete(networkingClient *golangsdk.ServiceClient, fId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloudStack Floating IP %s.\n", fId)

		f, err := floatingips.Get(networkingClient, fId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloudStack Floating IP %s", fId)
				return f, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault500); ok {
				log.Printf("[DEBUG] Got 500 error when delting HuaweiCloudStack Floating IP %s, it should be stream control on API server, try again later", fId)
				return f, "ACTIVE", nil
			}
			return f, "ACTIVE", err
		}

		err = floatingips.Delete(networkingClient, fId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloudStack Floating IP %s", fId)
				return f, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault500); ok {
				log.Printf("[DEBUG] Got 500 error when delting HuaweiCloudStack Floating IP %s, it should be stream control on API server, try again later", fId)
				return f, "ACTIVE", nil
			}
			return f, "ACTIVE", err
		}

		log.Printf("[DEBUG] HuaweiCloudStack Floating IP %s still active.\n", fId)
		return f, "ACTIVE", nil
	}
}
