---
subcategory: "Application Load Balancer(alb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener_certificates"
sidebar_current: "docs-tencentcloud-data_source-alb_listener_certificates"
description: |-
  ALB listener certificates data source provides the list of certificates bound to an ALB listener.
---

# tencentcloud_alb_listener_certificates

ALB listener certificates data source provides the list of certificates bound to an ALB listener.

## Example Usage

```hcl
data "tencentcloud_alb_listener_certificates" "example" {
  load_balancer_id = "alb-12345678"
  listener_id      = "lst-12345678"
  certificate_type = "SVR"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_type` - (Required, String) Certificate type. Values: `CA` or `SVR` (server certificate).
* `listener_id` - (Required, String) ID of the listener, in the format of lst- followed by 8 alphanumeric characters.
* `load_balancer_id` - (Required, String) ID of the ALB instance, in the format of alb- followed by 8 alphanumeric characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `associated_time` - Time when the certificate was bound.
  * `certificate_id` - ID of the certificate.
  * `certificate_type` - Certificate type. Values: `CA` or `SVR` (server certificate).
  * `is_default` - Whether the certificate is the default certificate of the listener. Values: `true` (default certificate), `false` (additional certificate).
  * `status` - Binding status of the certificate and the listener. Values: `Associated` (associated), `Associating` (associating), `Disassociating` (disassociating), `Error` (error).
* `certificates` - List of certificates bound to the listener.
* `id` - ID of the data source, a random token.


