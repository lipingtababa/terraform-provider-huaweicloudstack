---
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: huaweicloudstack_networking_floatingip_v2"
sidebar_current: "docs-huaweicloudstack-resource-networking-floatingip-v2"
description: |-
  Manages a V2 floating IP resource within HuaweiCloudStack.
---

# huaweicloudstack\_networking\_floatingip_v2

Manages a V2 floating IP resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "huaweicloudstack_networking_floatingip_v2" "floatip_1" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a floating IP that can be used with
    another networking resource, such as a load balancer. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    floating IP (which may or may not have a different address).

* `pool` - (Optional) The name of the pool from which to obtain the floating
    IP. Only admin_external_net is valid. Changing this creates a new floating IP.

* `port_id` - (Optional) ID of an existing port with at least one IP address to
    associate with this floating IP.

* `tenant_id` - (Optional) The target tenant ID in which to allocate the floating
    IP, if you specify this together with a port_id, make sure the target port
    belongs to the same tenant. Changing this creates a new floating IP (which
    may or may not have a different address)

* `fixed_ip` - Fixed IP of the port to associate with this floating IP. Required if
the port has multiple fixed IPs.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `pool` - See Argument Reference above.
* `address` - The actual floating IP address itself.
* `port_id` - ID of associated port.
* `tenant_id` - the ID of the tenant in which to create the floating IP.
* `fixed_ip` - The fixed IP which the floating IP maps to.

## Import

Floating IPs can be imported using the `id`, e.g.

```
$ terraform import huaweicloudstack_networking_floatingip_v2.floatip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
