---
subcategory: "Application Load Balancer(ALB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener"
sidebar_current: "docs-tencentcloud-resource-alb_listener"
description: |-
  Provides a resource to create an ALB listener.
---

# tencentcloud_alb_listener

Provides a resource to create an ALB listener.

## Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-12345678"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-example"

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }
      target_group_sticky_session {
        enabled = false
        timeout = 1000
      }
    }
  }

  gzip_enabled = true
  idle_timeout = 15
}
```

HTTPS listener with certificates:

```hcl
resource "tencentcloud_alb_listener" "https_example" {
  load_balancer_id   = "alb-12345678"
  listener_port      = 443
  listener_protocol  = "HTTPS"
  listener_name      = "tf-example-https"
  certificate_ids    = ["cert-12345678"]
  security_policy_id = "tls-12345678"
  http2enabled       = true

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }
    }
  }

  x_forwarded_for_config {
    x_forwarded_for_mode = "append"
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, String, ForceNew) ID of the ALB instance, in the format `alb-` followed by 8 alphanumeric characters.
* `listener_port` - (Required, Int, ForceNew) Port used by the frontend of the load balancer instance. Valid values: 1~65535.
* `listener_protocol` - (Required, String, ForceNew) Listener protocol. Valid values: `HTTP`, `HTTPS` or `QUIC`.
* `default_actions` - (Required, List) List of default forwarding rule actions. Currently a listener only supports adding 1 default forwarding rule action.
  * `type` - (Required, String) Forwarding action type. When creating a listener, the default forwarding action type only supports forwarding to target group.
  * `target_group_config` - (Required, List) Forwarding target group configuration. When creating a listener, the target group configuration in the forwarding action only supports a single target group.
    * `target_groups` - (Required, List) Target group list.
      * `target_group_id` - (Required, String) Target group ID, in the format `lbtg-` followed by 8 alphanumeric characters.
      * `weight` - (Optional, Int) Weight, valid range: [0, 100], default is 10.
    * `target_group_sticky_session` - (Optional, List) Sticky session between target groups.
      * `enabled` - (Required, Bool) Whether to enable sticky session, disabled by default.
      * `timeout` - (Optional, Int) Timeout in seconds, valid range: 1~86400, default value: 1000.
* `ca_certificate_ids` - (Optional, List) List of CA certificate IDs configured for the listener. Currently a listener only supports adding 1 CA certificate. Required when `ca_enabled` is true.
* `ca_enabled` - (Optional, Bool) Whether to enable mutual authentication. Valid values: `true` (enabled), `false` (disabled, default).
* `certificate_ids` - (Optional, List) List of server certificate IDs.
* `client_token` - (Optional, String) Client token used to ensure request idempotency. Only ASCII characters are supported.
* `gzip_enabled` - (Optional, Bool) Whether to enable Gzip compression. Valid values: `true` (default), `false`.
* `http2enabled` - (Optional, Bool) Whether to enable HTTP/2. HTTP protocol defaults to `false`, HTTPS protocol defaults to `true`. Only HTTPS protocol supports this parameter.
* `idle_timeout` - (Optional, Int) Connection idle timeout in seconds. Valid range: 1~600. Default value: 15.
* `listener_name` - (Optional, String) Custom listener name. Length is 1~255 characters; it can contain Chinese characters, letters, digits, hyphens (-), slashes (/), periods (.) and underscores (_).
* `request_timeout` - (Optional, Int) Connection request timeout in seconds. Valid values: 1~600. Default value: 60.
* `security_policy_id` - (Optional, String) Security policy ID, in the format `tls-` followed by 8 alphanumeric characters.
* `tags` - (Optional, Map) Tag list. Up to 20 tags are supported.
* `x_forwarded_for_config` - (Optional, List) X-Forwarded-For configuration.
  * `x_forwarded_for_alb_id_enabled` - (Optional, Bool) Whether to obtain the load balancer instance ID through the ALB-ID header field.
  * `x_forwarded_for_client_src_port_enabled` - (Optional, Bool) Whether to obtain the client port accessing the load balancer instance through the X-Forwarded-Client-srcport header field.
  * `x_forwarded_for_host_enabled` - (Optional, Bool) Whether to enable obtaining the client domain accessing the load balancer instance through the X-Forwarded-Host header field.
  * `x_forwarded_for_mode` - (Optional, String) Specifies how to handle the X-Forwarded-For (XFF) HTTP header field. Valid values: `append` (default), `remove`, `passthrough`.
  * `x_forwarded_for_port_enabled` - (Optional, Bool) Whether to obtain the listener port of the load balancer instance through the X-Forwarded-Port header field.
  * `x_forwarded_for_proto_enabled` - (Optional, Bool) Whether to obtain the listener protocol of the load balancer instance through the X-Forwarded-Proto header field.
  * `x_tencent_client_idn_enabled` - (Optional, Bool) Whether to access the issuer of the client certificate $ssl_client_i_dn through the X-Tencent-Client-IDN header.
  * `x_tencent_client_sdn_enabled` - (Optional, Bool) Whether to access the subject of the client certificate $ssl_client_s_dn through the X-Tencent-Client-SDN header.
  * `x_tencent_client_serial_enabled` - (Optional, Bool) Whether to access the serial number of the client certificate $ssl_client_serial through the X-Tencent-Client-Serial header.
  * `x_tencent_client_verify_enabled` - (Optional, Bool) Whether to access the verification result of the client certificate $ssl_client_verify through the X-Tencent-Client-Verify header.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource in the format `loadBalancerId#listenerId`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateListener](https://cloud.tencent.com/document/api/1822/133705)
* Read: [DescribeListeners](https://cloud.tencent.com/document/api/1822/133700)
* Update: [ModifyListenerAttributes](https://cloud.tencent.com/document/api/1822/133698)
* Delete: [DeleteListener](https://cloud.tencent.com/document/api/1822/133704)

## Import

ALB listener can be imported using the loadBalancerId#listenerId, e.g.

```
terraform import tencentcloud_alb_listener.example alb-12345678#lst-12345678
```
