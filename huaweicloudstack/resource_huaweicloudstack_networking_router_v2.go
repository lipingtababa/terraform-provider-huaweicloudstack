package huaweicloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/routers"
)

func resourceNetworkingRouterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingRouterV2Create,
		Read:   resourceNetworkingRouterV2Read,
		Update: resourceNetworkingRouterV2Update,
		Delete: resourceNetworkingRouterV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"distributed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"external_gateway": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         false,
				DiffSuppressFunc: suppressExternalGateway,
			},
			"enable_snat": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceNetworkingRouterV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	createOpts := RouterCreateOpts{
		routers.CreateOpts{
			Name:     d.Get("name").(string),
			TenantID: d.Get("tenant_id").(string),
		},
		MapValueSpecs(d),
	}

	if asuRaw, ok := d.GetOk("admin_state_up"); ok {
		asu := asuRaw.(bool)
		createOpts.AdminStateUp = &asu
	}

	if dRaw, ok := d.GetOk("distributed"); ok {
		d := dRaw.(bool)
		createOpts.Distributed = &d
	}

	externalGateway := d.Get("external_gateway").(string)
	if externalGateway != "" {
		gatewayInfo := routers.GatewayInfo{
			NetworkID: externalGateway,
		}
		createOpts.GatewayInfo = &gatewayInfo
	}

	if esRaw, ok := d.GetOk("enable_snat"); ok {
		if externalGateway == "" {
			return fmt.Errorf("setting enable_snat requires external_gateway to be set")
		}
		es := esRaw.(bool)
		createOpts.GatewayInfo.EnableSNAT = &es
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := routers.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack Neutron router: %s", err)
	}
	log.Printf("[INFO] Router ID: %s", n.ID)

	log.Printf("[DEBUG] Waiting for HuaweiCloudStack Neutron Router (%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForRouterActive(networkingClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(n.ID)

	return resourceNetworkingRouterV2Read(d, meta)
}

func resourceNetworkingRouterV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	n, err := routers.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HuaweiCloudStack Neutron Router: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Router %s: %+v", d.Id(), n)

	d.Set("name", n.Name)
	d.Set("admin_state_up", n.AdminStateUp)
	d.Set("distributed", n.Distributed)
	d.Set("tenant_id", n.TenantID)
	d.Set("external_gateway", n.GatewayInfo.NetworkID)
	d.Set("enable_snat", n.GatewayInfo.EnableSNAT)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingRouterV2Update(d *schema.ResourceData, meta interface{}) error {
	routerId := d.Id()
	osMutexKV.Lock(routerId)
	defer osMutexKV.Unlock(routerId)

	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	var updateOpts routers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Gateway settings
	var updateGatewaySettings bool
	var externalGateway string
	gatewayInfo := routers.GatewayInfo{}

	if v := d.Get("external_gateway").(string); v != "" {
		externalGateway = v
	}

	if externalGateway != "" {
		gatewayInfo.NetworkID = externalGateway
	}
	if d.HasChange("external_gateway") {
		updateGatewaySettings = true
	}

	if d.HasChange("enable_snat") {
		updateGatewaySettings = true
		if externalGateway == "" {
			return fmt.Errorf("setting enable_snat requires external_gateway to be set")
		}

		enableSNAT := d.Get("enable_snat").(bool)
		gatewayInfo.EnableSNAT = &enableSNAT
	}

	if updateGatewaySettings {
		updateOpts.GatewayInfo = &gatewayInfo
	}

	log.Printf("[DEBUG] Updating Router %s with options: %+v", d.Id(), updateOpts)

	_, err = routers.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloudStack Neutron Router: %s", err)
	}

	return resourceNetworkingRouterV2Read(d, meta)
}

func resourceNetworkingRouterV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForRouterDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloudStack Neutron Router: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForRouterActive(networkingClient *golangsdk.ServiceClient, routerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := routers.Get(networkingClient, routerId).Extract()
		if err != nil {
			return nil, r.Status, err
		}

		log.Printf("[DEBUG] HuaweiCloudStack Neutron Router: %+v", r)
		return r, r.Status, nil
	}
}

func waitForRouterDelete(networkingClient *golangsdk.ServiceClient, routerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloudStack Router %s.\n", routerId)

		r, err := routers.Get(networkingClient, routerId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloudStack Router %s", routerId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = routers.Delete(networkingClient, routerId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloudStack Router %s", routerId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] HuaweiCloudStack Router %s still active.\n", routerId)
		return r, "ACTIVE", nil
	}
}
