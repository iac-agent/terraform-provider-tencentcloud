---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener"
sidebar_current: "docs-tencentcloud-resource-alb_listener"
description: |-
  Provides a resource to manage an ALB listener.
---

# tencentcloud_alb_listener

Provides a resource to manage an ALB listener.

## Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-12345678"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-test-listener"

  default_actions {
    type = "ForwardGroup"

    target_group_config {
      target_groups {
        target_group_id = "lbtg-abcdefgh"
        weight          = 10
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `default_actions` - (Required, List) List of default forwarding rule actions. Currently a listener supports adding only 1 default forwarding rule action.
* `listener_port` - (Required, Int) Port used by the frontend of the load balancer instance. Value range: 1~65535.
* `listener_protocol` - (Required, String) Listener protocol. Valid values: `HTTP`, `HTTPS`.
* `load_balancer_id` - (Required, String) Load balancer instance ID, in the form `alb-` followed by 8 letters or digits.
* `ca_certificate_ids` - (Optional, List) List of CA certificate IDs configured for the listener. Currently a listener supports adding only 1 CA certificate. Required when `ca_enabled` is `true`.
* `ca_enabled` - (Optional, Bool) Whether to enable mutual authentication. Valid values: `true` (enabled), `false` (not enabled, default).
* `certificate_ids` - (Optional, List) List of server certificate IDs.
* `client_token` - (Optional, String) Client token used to ensure idempotency of the request. Generate a value on your client side and make sure it is unique among different requests. Only ASCII characters are supported.
* `gzip_enabled` - (Optional, Bool) Whether to enable Gzip compression. Valid values: `true` (enabled, default), `false` (not enabled).
* `http2_enabled` - (Optional, Bool) Whether to enable the HTTP/2 feature. Default value: `false` for the HTTP protocol and `true` for the HTTPS protocol. Only the HTTPS protocol supports this parameter.
* `idle_timeout` - (Optional, Int) Connection idle timeout in seconds. Value range: 1~600. Default value: 15. If no request is received within the timeout, the load balancer disconnects the current connection and creates a new one when the next request arrives.
* `listener_name` - (Optional, String) Custom listener name. Length: 1~255 characters, must be a Chinese or harmless string, may contain Chinese characters, letters, digits, hyphens (-), forward slashes (/), periods (.) and underscores (_).
* `request_timeout` - (Optional, Int) Connection request timeout in seconds. Value range: 1~600. Default value: 60. If the backend server does not return a response within the timeout, the load balancer gives up waiting and returns HTTP 504 to the client.
* `security_policy_id` - (Optional, String) Security policy ID, in the form `tls-` followed by 8 letters or digits.
* `tags` - (Optional, List) Tag list. Maximum: 20 tags.
* `x_forwarded_for_config` - (Optional, Object) X-Forwarded-For configuration.

The `default_actions` object supports the following:

* `target_group_config` - (Required, Object) Forwarding target group configuration. When creating a listener, the target group configuration of a forwarding action supports only a single target group.
* `type` - (Required, String) Forwarding action type. When creating a listener, the default forwarding action type supports only forwarding to a target group.

The `target_group_config` object of `default_actions` supports the following:

* `target_groups` - (Required, List) Target group list.
* `target_group_sticky_session` - (Optional, Object) Session persistence between target groups.

The `target_group_sticky_session` object of `target_group_config` supports the following:

* `enabled` - (Required, Bool) Whether to enable session persistence. Disabled by default.
* `timeout` - (Optional, Int) Timeout in seconds. Value range: 1~86400. Default value: 1000.

The `target_groups` object of `target_group_config` supports the following:

* `target_group_id` - (Required, String) Target group ID, in the form `lbtg-` followed by 8 letters or digits.
* `weight` - (Optional, Int) Weight. Value range: [0, 100]. Default value: 10.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

The `x_forwarded_for_config` object supports the following:

* `x_forwarded_for_alb_id_enabled` - (Optional, Bool) Whether to obtain the load balancer instance ID through the ALB-ID header. `true`: yes. `false`: no.
* `x_forwarded_for_client_src_port_enabled` - (Optional, Bool) Whether to obtain the port of the client that accesses the load balancer instance through the X-Forwarded-Client-srcport header. `true`: yes. `false`: no.
* `x_forwarded_for_host_enabled` - (Optional, Bool) Whether to enable obtaining the domain name of the client that accesses the load balancer instance through the X-Forwarded-Host header. `true`: yes. `false`: no.
* `x_forwarded_for_mode` - (Optional, String) How to handle the X-Forwarded-For (XFF) HTTP header. `append`: append mode (default), appends the real IP address of the client to the end of the X-Forwarded-For header and keeps the original XFF chain. `remove`: remove mode, removes the X-Forwarded-For header and does not pass it to the backend server. `passthrough`: passthrough mode, keeps the X-Forwarded-For header unchanged and passes it to the backend server without modification.
* `x_forwarded_for_port_enabled` - (Optional, Bool) Whether to obtain the listening port of the load balancer instance through the X-Forwarded-Port header. `true`: yes. `false`: no.
* `x_forwarded_for_proto_enabled` - (Optional, Bool) Whether to obtain the listening protocol of the load balancer instance through the X-Forwarded-Proto header. `true`: yes. `false`: no.
* `x_tencent_client_idn_enabled` - (Optional, Bool) Whether to access the issuer of the client certificate `$ssl_client_i_dn` through the X-Tencent-Client-IDN header. `true`: yes. `false`: no.
* `x_tencent_client_sdn_enabled` - (Optional, Bool) Whether to access the subject `$ssl_client_s_dn` of the client certificate through the X-Tencent-Client-SDN header. `true`: yes. `false`: no.
* `x_tencent_client_serial_enabled` - (Optional, Bool) Whether to access the serial number `$ssl_client_serial` of the client certificate through the X-Tencent-Client-Serial header. `true`: yes. `false`: no.
* `x_tencent_client_verify_enabled` - (Optional, Bool) Whether to access the verification result `$ssl_client_verify` of the client certificate through the X-Tencent-Client-Verify header. `true`: yes. `false`: no.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, in the form `<load_balancer_id>#<listener_id>`.
* `listener_id` - Listener ID, in the form `lst-` followed by 8 letters or digits.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateListener](https://www.tencentcloud.com/document/api/1311/84262)
  - [DescribeListeners](https://www.tencentcloud.com/document/api/1311/84257)
* Read: [DescribeListeners](https://www.tencentcloud.com/document/api/1311/84257)
* Update:
  - [DescribeListeners](https://www.tencentcloud.com/document/api/1311/84257)
  - [ModifyListenerAttributes](https://www.tencentcloud.com/document/api/1311/84255)
* Delete: [DeleteListener](https://www.tencentcloud.com/document/api/1311/84261)

## Import

ALB listener can be imported using the composite id `load_balancer_id#listener_id`, e.g.

```
terraform import tencentcloud_alb_listener.example alb-12345678#lst-12345678
```

