---
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: huaweicloudstack_networking_router_v2"
sidebar_current: "docs-huaweicloudstack-resource-networking-router-v2"
description: |-
  Manages a V2 router resource within HuaweiCloudStack.
---

# huaweicloudstack\_networking\_router_v2

Manages a V2 router resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "huaweicloudstack_networking_router_v2" "router_1" {
  name                = "my_router"
  admin_state_up      = true
  external_network_id = "f67f0d72-0ddf-11e4-9d95-e1f29f417e2f"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
    A networking client is needed to create a router. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    router.

* `name` - (Optional) A unique name for the router. Changing this
    updates the `name` of an existing router.

* `admin_state_up` - (Optional) Administrative up/down status for the router
    (must be "true" or "false" if provided). Changing this updates the
    `admin_state_up` of an existing router.

* `distributed` - (Optional) Indicates whether or not to create a
    distributed router. The default policy setting in Neutron restricts
    usage of this property to administrative users only.

* `external_network_id` - (Optional) The network UUID of an external gateway
    for the router. A router with an external gateway is required if any
    compute instances or load balancers will be using floating IPs. Changing
    this updates the external gateway of the router.

* `enable_snat` - (Optional) Enable Source NAT for the router. Valid values are
    "true" or "false". An `external_network_id` has to be set in order to
    set this property. Changing this updates the `enable_snat` of the router.

* `external_fixed_ip` - (Optional) An external fixed IP for the router. This
    can be repeated. The structure is described below. An `external_network_id`
    has to be set in order to set this property. Changing this updates the
    external fixed IPs of the router.

* `tenant_id` - (Optional) The owner of the floating IP. Required if admin wants
    to create a router for another tenant. Changing this creates a new router.

* `value_specs` - (Optional) Map of additional driver-specific options.

* `availability_zone_hints` -  (Optional) An availability zone is used to make 
    network resources highly available. Used for resources with high availability so that they are scheduled on different availability zones. Changing
    this creates a new router.

The `external_fixed_ip` block supports:

* `subnet_id` - (Optional) Subnet in which the fixed IP belongs to.

* `ip_address` - (Optional) The IP address to set on the router.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the router.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `external_network_id` - See Argument Reference above.
* `enable_snat` - See Argument Reference above.
* `external_fixed_ip` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `value_specs` - See Argument Reference above.
* `availability_zone_hints` - See Argument Reference above.

## Import

Routers can be imported using the `id`, e.g.

```
$ terraform import huaweicloudstack_networking_router_v2.router_1 014395cd-89fc-4c9b-96b7-13d1ee79dad2
```
