---
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: huaweicloudstack_networking_network_v2"
sidebar_current: "docs-huaweicloudstack-datasource-networking-network-v2"
description: |-
  Get information on an HuaweiCloudStack Network.
---

# huaweicloudstack\_networking\_network\_v2

Use this data source to get the ID of an available HuaweiCloudStack network.

## Example Usage

```hcl
data "huaweicloudstack_networking_network_v2" "network" {
  name = "tf_test_network"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve networks ids. If omitted, the
  `region` argument of the provider is used.

* `network_id` - (Optional) The ID of the network.

* `name` - (Optional) The name of the network.

* `status` - (Optional) The status of the network.

* `matching_subnet_cidr` - (Optional) The CIDR of a subnet within the network.

* `tenant_id` - (Optional) The owner of the network.

* `availability_zone_hints` - (Optional) The availability zone candidates for the network.


## Attributes Reference

`id` is set to the ID of the found network. In addition, the following attributes
are exported:

* `admin_state_up` - (Optional) The administrative state of the network.
* `name` - See Argument Reference above.
* `region` - See Argument Reference above.
* `shared` - (Optional)  Specifies whether the network resource can be accessed
    by any tenant or not.
* `availability_zone_hints` - (Optional) The availability zone candidates for the network.
