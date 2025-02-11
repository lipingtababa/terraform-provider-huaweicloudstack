---
layout: "huaweicloudstack"
page_title: "Provider: HuaweiCloudStack"
sidebar_current: "docs-huaweicloudstack-index"
description: |-
  The HuaweiCloudStack provider is used to interact with the many resources supported by HuaweiCloudStack. The provider needs to be configured with the proper credentials before it can be used.
---

# HuaweiCloudStack Provider

The HuaweiCloudStack provider is used to interact with the
many resources supported by HuaweiCloudStack. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the HuaweiCloudStack Provider
provider "huaweicloudstack" {
  user_name   = "${var.user_name}"
  password    = "${var.password}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  region      = "${var.region}"
  auth_url    = "https://iam.myhwclouds.com:443/v3"
}

# Create a web server
resource "huaweicloudstack_compute_instance_v2" "test-server" {
  # ...
}
```

## Authentication

This provider offers 2 means for authentication.

- Username + Password
- Token

### Username + Password

```hcl
provider "huaweicloudstack" {
  user_name   = "${var.user_name}"
  password    = "${var.password}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  region      = "${var.region}"
  auth_url    = "https://iam.myhwclouds.com:443/v3"
}
```

### Token

```hcl
provider "huaweicloudstack" {
  token       = "${var.token}"
  domain_name = "${var.domain_name}"
  tenant_name = "${var.tenant_name}"
  region      = "${var.region}"
  auth_url    = "https://iam.myhwclouds.com:443/v3"
}
```


## Configuration Reference

The following arguments are supported:

* `access_key` - (Optional) The access key of the HuaweiCloudStack to use.
  If omitted, the `OS_ACCESS_KEY` environment variable is used.

* `secret_key` - (Optional) The secret key of the HuaweiCloudStack to use.
  If omitted, the `OS_SECRET_KEY` environment variable is used.

* `auth_url` - (Required) The Identity authentication URL. If omitted, the
  `OS_AUTH_URL` environment variable is used.

* `region` - (Optional) The region of the HuaweiCloudStack to use. If omitted,
  the `OS_REGION_NAME` environment variable is used. If `OS_REGION_NAME` is
  not set, then no region will be used. It should be possible to omit the
  region in single-region HuaweiCloudStack environments, but this behavior may vary
  depending on the HuaweiCloudStack environment being used.

* `user_name` - (Optional) The Username to login with. If omitted, the
  `OS_USERNAME` environment variable is used.

* `user_id` - (Optional) The User ID to login with. If omitted, the
  `OS_USER_ID` environment variable is used.

* `tenant_id` - (Optional) The ID of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_ID` or
  `OS_PROJECT_ID` environment variables are used.

* `tenant_name` - (Optional) The Name of the Tenant (Identity v2) or Project
  (Identity v3) to login with. If omitted, the `OS_TENANT_NAME` or
  `OS_PROJECT_NAME` environment variable are used.

* `password` - (Optional) The Password to login with. If omitted, the
  `OS_PASSWORD` environment variable is used.

* `token` - (Optional; Required if not using `user_name` and `password`)
  A token is an expiring, temporary means of access issued via the Keystone
  service. By specifying a token, you do not have to specify a username/password
  combination, since the token was already created by a username/password out of
  band of Terraform. If omitted, the `OS_AUTH_TOKEN` environment variable is used.

* `domain_id` - (Optional) The ID of the Domain to scope to (Identity v3). If
  If omitted, the following environment variables are checked (in this order):
  `OS_USER_DOMAIN_ID`, `OS_PROJECT_DOMAIN_ID`, `OS_DOMAIN_ID`.

* `domain_name` - (Optional) The Name of the Domain to scope to (Identity v3).
  If omitted, the following environment variables are checked (in this order):
  `OS_USER_DOMAIN_NAME`, `OS_PROJECT_DOMAIN_NAME`, `OS_DOMAIN_NAME`,
  `DEFAULT_DOMAIN`.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `OS_INSECURE` environment variable is used.

* `cacert_file` - (Optional) Specify a custom CA certificate when communicating
  over SSL. You can specify either a path to the file or the contents of the
  certificate. If omitted, the `OS_CACERT` environment variable is used.

* `cert` - (Optional) Specify client certificate file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the certificate. If omitted the `OS_CERT` environment variable is used.

* `key` - (Optional) Specify client private key file for SSL client
  authentication. You can specify either a path to the file or the contents of
  the key. If omitted the `OS_KEY` environment variable is used.

* `endpoint_type` - (Optional) Specify which type of endpoint to use from the
  service catalog. It can be set using the OS_ENDPOINT_TYPE environment
  variable. If not set, public endpoints is used.

* `delegated_project` - (Optional) The name of delegated project (Identity v3).

## Additional Logging

This provider has the ability to log all HTTP requests and responses between
Terraform and the HuaweiCloudStack which is useful for troubleshooting and
debugging.

To enable these logs, set the `OS_DEBUG` environment variable to `1` along
with the usual `TF_LOG=DEBUG` environment variable:

```shell
$ OS_DEBUG=1 TF_LOG=DEBUG terraform apply
```

If you submit these logs with a bug report, please ensure any sensitive
information has been scrubbed first!

## Testing and Development

In order to run the Acceptance Tests for development, the following environment
variables must also be set:

* `OS_REGION_NAME` - The region in which to create the server instance.

* `OS_IMAGE_ID` or `OS_IMAGE_NAME` - a UUID or name of an existing image in
    Glance.

* `OS_FLAVOR_ID` or `OS_FLAVOR_NAME` - an ID or name of an existing flavor.

* `OS_POOL_NAME` - The name of a Floating IP pool.

* `OS_NETWORK_ID` - The UUID of a network in your test environment.

* `OS_EXTGW_ID` - The UUID of the external gateway.

You should be able to use any HuaweiCloudStack environment to develop on as long as the
above environment variables are set.
