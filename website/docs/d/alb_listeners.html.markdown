---
subcategory: "Application Load Balancer(alb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listeners"
sidebar_current: "docs-tencentcloud-data_source-alb_listeners"
description: |-
  Use this data source to query the list of ALB listeners.
---

# tencentcloud_alb_listeners

Use this data source to query the list of ALB listeners.

## Example Usage

```hcl
data "tencentcloud_alb_listeners" "example" {
  load_balancer_id = "alb-f8q2xk9m"
  listener_ids     = ["lst-a1b2c3d4"]

  filters = [
    {
      name   = "Protocol"
      values = ["HTTP"]
    },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, String) Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters.
* `filters` - (Optional, List) Filter conditions, with a maximum of 20. The following fields are supported: Protocol (protocol type), Tags (tags).
* `listener_ids` - (Optional, List) List of listener IDs. An ID is in the format of lst- followed by 8 alphanumeric characters.

The `filters` object supports the following:

* `name` - (Required, String) Name of the filter.
* `values` - (Required, List) Values of the filter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
    * `x_forwarded_for_alb_id_enabled` - Whether to obtain the load balancer instance ID through the ALB-ID header field. true: yes; false: no.
    * `x_forwarded_for_client_src_port_enabled` - Whether to obtain the client port accessing the load balancer through the X-Forwarded-Client-srcport header field. true: yes; false: no.
    * `x_forwarded_for_host_enabled` - Whether to enable obtaining the domain name accessing the load balancer through the X-Forwarded-Host header field. true: yes; false: no.
    * `x_forwarded_for_mode` - How to handle the X-Forwarded-For (XFF) HTTP header field. append: append mode (default), appending the client real IP address to the end of the X-Forwarded-For header while preserving the original XFF chain information; remove: removal mode, removing the X-Forwarded-For header field and not passing it to the backend server; passthrough: passthrough mode, keeping X-Forwarded-For unchanged and passing it to the backend server without any modification.
    * `x_forwarded_for_port_enabled` - Whether to obtain the listener port of the load balancer through the X-Forwarded-Port header field. true: yes; false: no.
    * `x_forwarded_for_proto_enabled` - Whether to obtain the listener protocol of the load balancer through the X-Forwarded-Proto header field. true: yes; false: no.
    * `x_tencent_client_idn_enabled` - Whether to access the issuer of the client certificate ($ssl_client_i_dn) through the X-Tencent-Client-IDN header. true: yes; false: no.
    * `x_tencent_client_sdn_enabled` - Whether to access the subject of the client certificate ($ssl_client_s_dn) through the X-Tencent-Client-SDN header. true: yes; false: no.
    * `x_tencent_client_serial_enabled` - Whether to access the serial number of the client certificate ($ssl_client_serial) through the X-Tencent-Client-Serial header. true: yes; false: no.
    * `x_tencent_client_verify_enabled` - Whether to access the verification result of the client certificate ($ssl_client_verify) through the X-Tencent-Client-Verify header. true: yes; false: no.
  * `ca_enable` - Whether mutual authentication is enabled.
  * `create_time` - Creation time of the listener instance, in ISO 8601 format (for example, 2025-01-01T08:30:00+08:00).
  * `gzip_enabled` - Whether Gzip compression is enabled.
  * `http2_enable` - Whether HTTP/2 is enabled.
  * `idle_timeout` - Idle timeout period.
  * `listener_id` - Listener ID, in the format of lst- followed by 8 alphanumeric characters.
  * `listener_name` - Listener name.
  * `listener_port` - Listener port.
  * `listener_protocol` - Listener protocol.
  * `listener_status` - Listener status. Valid values: Active (running), Provisioning (creating), Configuring (modifying configuration), ProvisionFailed (creation failed).
  * `modify_time` - Last modification time of the listener instance, in ISO 8601 format (for example, 2025-01-01T08:30:00+08:00).
  * `request_timeout` - Connection request timeout period.
  * `tags` - Tags.
  * `tls_security_policy_id` - Security policy ID.
  * `x_forwarded_for_config` - XForwardedFor configuration.
* `id` - ID of the data source, a random token.
* `listeners` - Listener information.


