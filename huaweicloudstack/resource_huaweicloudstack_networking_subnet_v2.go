package huaweicloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/subnets"
)

func resourceNetworkingSubnetV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingSubnetV2Create,
		Read:   resourceNetworkingSubnetV2Read,
		Update: resourceNetworkingSubnetV2Update,
		Delete: resourceNetworkingSubnetV2Delete,
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
			"network_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressAsuDiff,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"allocation_pools": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"no_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
				ForceNew: true,
			},
			"enable_dhcp": {
				Type:         schema.TypeBool,
				Optional:     true,
				ForceNew:     false,
				Default:      true,
				ValidateFunc: validateTrueOnly,
			},
			"dns_nameservers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"host_routes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr": {
							Type:     schema.TypeString,
							Required: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingSubnetV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	log.Printf("[DEBUG] Checking Subnet 0 %s: network_id: %s", d.Id(), d.Get("network_id").(string))
	_, nid := ExtractValSFromNid(d.Get("network_id").(string))
	createOpts := SubnetCreateOpts{
		subnets.CreateOpts{
			NetworkID:       nid,
			CIDR:            d.Get("cidr").(string),
			Name:            d.Get("name").(string),
			TenantID:        d.Get("tenant_id").(string),
			AllocationPools: resourceSubnetAllocationPoolsV2(d),
			DNSNameservers:  resourceSubnetDNSNameserversV2(d),
			HostRoutes:      resourceSubnetHostRoutesV2(d),
			EnableDHCP:      nil,
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Checking Subnet 1 %s: network_id: %s", d.Id(), d.Get("network_id").(string))
	noGateway := d.Get("no_gateway").(bool)
	gatewayIP := d.Get("gateway_ip").(string)

	if gatewayIP != "" && noGateway {
		return fmt.Errorf("Both gateway_ip and no_gateway cannot be set")
	}

	if gatewayIP != "" {
		createOpts.GatewayIP = &gatewayIP
	}

	if noGateway {
		disableGateway := ""
		createOpts.GatewayIP = &disableGateway
	}

	enableDHCP := d.Get("enable_dhcp").(bool)
	createOpts.EnableDHCP = &enableDHCP

	log.Printf("[DEBUG] Checking Subnet 2 %s: network_id: %s", d.Id(), d.Get("network_id").(string))
	if v, ok := d.GetOk("ip_version"); ok {
		ipVersion := resourceNetworkingSubnetV2DetermineIPVersion(v.(int))
		createOpts.IPVersion = ipVersion
	}

	s, err := subnets.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack Neutron subnet: %s", err)
	}

	log.Printf("[DEBUG] Waiting for Subnet (%s) to become available", s.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForSubnetActive(networkingClient, s.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(s.ID)

	log.Printf("[DEBUG] Checking Subnet 3 %s: network_id: %s", d.Id(), d.Get("network_id").(string))
	log.Printf("[DEBUG] Created Subnet %s: %#v", s.ID, s)
	return resourceNetworkingSubnetV2Read(d, meta)
}

func resourceNetworkingSubnetV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Checking Subnet 4 %s: network_id: %s", d.Id(), d.Get("network_id").(string))
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	s, err := subnets.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "subnet")
	}

	log.Printf("[DEBUG] Retrieved Subnet %s: %#v", d.Id(), s)

	asu, _ := ExtractValSFromNid(d.Get("network_id").(string))
	nid := FormatNidFromValS(asu, s.NetworkID)
	d.Set("network_id", nid)
	d.Set("cidr", s.CIDR)
	d.Set("ip_version", s.IPVersion)
	d.Set("name", s.Name)
	d.Set("tenant_id", s.TenantID)
	d.Set("gateway_ip", s.GatewayIP)
	d.Set("dns_nameservers", s.DNSNameservers)
	d.Set("host_routes", s.HostRoutes)
	d.Set("enable_dhcp", s.EnableDHCP)
	d.Set("network_id", s.NetworkID)

	// Set the allocation_pools
	var allocationPools []map[string]interface{}
	for _, v := range s.AllocationPools {
		pool := make(map[string]interface{})
		pool["start"] = v.Start
		pool["end"] = v.End

		allocationPools = append(allocationPools, pool)
	}
	d.Set("allocation_pools", allocationPools)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingSubnetV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	// Check if both gateway_ip and no_gateway are set
	if _, ok := d.GetOk("gateway_ip"); ok {
		noGateway := d.Get("no_gateway").(bool)
		if noGateway {
			return fmt.Errorf("Both gateway_ip and no_gateway cannot be set.")
		}
	}

	var updateOpts subnets.UpdateOpts

	noGateway := d.Get("no_gateway").(bool)
	gatewayIP := d.Get("gateway_ip").(string)

	if gatewayIP != "" && noGateway {
		return fmt.Errorf("Both gateway_ip and no_gateway cannot be set")
	}

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("gateway_ip") {
		updateOpts.GatewayIP = nil
		if v, ok := d.GetOk("gateway_ip"); ok {
			gatewayIP := v.(string)
			updateOpts.GatewayIP = &gatewayIP
		}
	}

	if d.HasChange("no_gateway") {
		if d.Get("no_gateway").(bool) {
			gatewayIP := ""
			updateOpts.GatewayIP = &gatewayIP
		}
	}

	if d.HasChange("dns_nameservers") {
		updateOpts.DNSNameservers = resourceSubnetDNSNameserversV2(d)
	}

	if d.HasChange("host_routes") {
		updateOpts.HostRoutes = resourceSubnetHostRoutesV2(d)
	}

	if d.HasChange("enable_dhcp") {
		v := d.Get("enable_dhcp").(bool)
		updateOpts.EnableDHCP = &v
	}

	if d.HasChange("allocation_pools") {
		updateOpts.AllocationPools = resourceSubnetAllocationPoolsV2(d)
	}

	log.Printf("[DEBUG] Updating Subnet %s with options: %+v", d.Id(), updateOpts)

	_, err = subnets.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloudStack Neutron Subnet: %s", err)
	}

	return resourceNetworkingSubnetV2Read(d, meta)
}

func resourceNetworkingSubnetV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSubnetDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloudStack Neutron Subnet: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceSubnetAllocationPoolsV2(d *schema.ResourceData) []subnets.AllocationPool {
	rawAPs := d.Get("allocation_pools").([]interface{})
	aps := make([]subnets.AllocationPool, len(rawAPs))
	for i, raw := range rawAPs {
		rawMap := raw.(map[string]interface{})
		aps[i] = subnets.AllocationPool{
			Start: rawMap["start"].(string),
			End:   rawMap["end"].(string),
		}
	}
	return aps
}

func resourceSubnetDNSNameserversV2(d *schema.ResourceData) []string {
	rawDNSN := d.Get("dns_nameservers").(*schema.Set)
	dnsn := make([]string, rawDNSN.Len())
	for i, raw := range rawDNSN.List() {
		dnsn[i] = raw.(string)
	}
	return dnsn
}

func resourceSubnetHostRoutesV2(d *schema.ResourceData) []subnets.HostRoute {
	rawHR := d.Get("host_routes").([]interface{})
	hr := make([]subnets.HostRoute, len(rawHR))
	for i, raw := range rawHR {
		rawMap := raw.(map[string]interface{})
		hr[i] = subnets.HostRoute{
			DestinationCIDR: rawMap["destination_cidr"].(string),
			NextHop:         rawMap["next_hop"].(string),
		}
	}
	return hr
}

func resourceNetworkingSubnetV2DetermineIPVersion(v int) golangsdk.IPVersion {
	var ipVersion golangsdk.IPVersion
	switch v {
	case 4:
		ipVersion = golangsdk.IPv4
	case 6:
		ipVersion = golangsdk.IPv6
	}

	return ipVersion
}

func waitForSubnetActive(networkingClient *golangsdk.ServiceClient, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := subnets.Get(networkingClient, subnetId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloudStack Neutron Subnet: %+v", s)
		return s, "ACTIVE", nil
	}
}

func waitForSubnetDelete(networkingClient *golangsdk.ServiceClient, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloudStack Subnet %s.\n", subnetId)

		s, err := subnets.Get(networkingClient, subnetId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloudStack Subnet %s", subnetId)
				return s, "DELETED", nil
			}
			return s, "ACTIVE", err
		}

		err = subnets.Delete(networkingClient, subnetId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloudStack Subnet %s", subnetId)
				return s, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return s, "ACTIVE", nil
				}
			}
			return s, "ACTIVE", err
		}

		log.Printf("[DEBUG] HuaweiCloudStack Subnet %s still active.\n", subnetId)
		return s, "ACTIVE", nil
	}
}
