---
subcategory: "Application Load Balancer(alb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener_certificates"
sidebar_current: "docs-tencentcloud-data_source-alb_listener_certificates"
description: |-
  Use this data source to query the list of certificates bound to an Application Load Balancer listener.
---

# tencentcloud_alb_listener_certificates

Use this data source to query the list of certificates bound to an Application Load Balancer listener.

## Example Usage

```hcl
data "tencentcloud_alb_listener_certificates" "example" {
  certificate_type = "SVR"
  listener_id      = "lst-d9p3k7wa"
  load_balancer_id = "alb-f8q2xk9m"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_type` - (Required, String) Certificate type. Value: CA or SVR (server certificate).
* `listener_id` - (Required, String) Listener ID, in the format of lst- followed by 8 alphanumeric characters.
* `load_balancer_id` - (Required, String) Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `associated_time` - Certificate binding time.
  * `certificate_id` - Certificate ID.
  * `certificate_type` - Certificate type. Value: CA or SVR (server certificate).
  * `is_default` - Whether it is the default certificate of the listener. Value: true: default certificate; false: extended certificate.
  * `status` - Certificate-listener binding status. Value: Associated (associated), Associating (associating), Disassociating (disassociating), Error (error).
* `certificates` - List of certificates bound to the listener.
* `id` - ID of the data source, a random token.


