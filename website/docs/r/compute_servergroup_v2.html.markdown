---
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: huaweicloudstack_compute_servergroup_v2"
sidebar_current: "docs-huaweicloudstack-resource-compute-servergroup-v2"
description: |-
  Manages a V2 Server Group resource within HuaweiCloudStack.
---

# huaweicloudstack\_compute\_servergroup_v2

Manages a V2 Server Group resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "huaweicloudstack_compute_servergroup_v2" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    If omitted, the `region` argument of the provider is used. Changing
    this creates a new server group.

* `name` - (Required) A unique name for the server group. Changing this creates
    a new server group.

* `policies` - (Required) The set of policies for the server group. Only two
    two policies are available right now, and both are mutually exclusive. See
    the Policies section for more information. Changing this creates a new
    server group.

* `value_specs` - (Optional) Map of additional options.

## Policies

* `affinity` - All instances/servers launched in this group will be hosted on
    the same compute node.

* `anti-affinity` - All instances/servers launched in this group will be
    hosted on different compute nodes.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `policies` - See Argument Reference above.
* `members` - The instances that are part of this server group.

## Import

Server Groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloudstack_compute_servergroup_v2.test-sg 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf
```
