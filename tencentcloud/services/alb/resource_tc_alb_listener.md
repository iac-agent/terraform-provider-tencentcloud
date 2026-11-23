---
subcategory: ""
---

# tencentcloud_alb_listener

Provide a resource to create an ALB listener.

## Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-xxxxxxxx"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "example-listener"

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "lbtg-xxxxxxxx"
        weight          = 10
      }
      target_group_sticky_session {
        enabled = false
        timeout = 1000
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) ID of the ALB instance, in the format of `alb-xxxxxxxx`.
* `listener_port` - (Required, ForceNew) Front-end port of the ALB instance. Value range: 1-65535.
* `listener_protocol` - (Required, ForceNew) Listener protocol. Valid values: `HTTP`, `HTTPS` and `QUIC`.
* `default_actions` - (Required) List of default forwarding rule actions. Currently, only 1 default forwarding rule action is supported per listener.
  * `type` - (Required) Type of the forwarding action. When creating a listener, only forwarding to target group is supported.
  * `target_group_config` - (Required) Forwarding target group configuration. Only a single target group is supported.
    * `target_groups` - (Required) List of target groups.
      * `target_group_id` - (Required) Target group ID, in the format of `lbtg-xxxxxxxx`.
      * `weight` - (Optional) Weight. Value range: [0, 100]. Default is 10.
    * `target_group_sticky_session` - (Optional) Sticky session among target groups.
      * `enabled` - (Required) Whether to enable sticky session. Default is disabled.
      * `timeout` - (Optional) Timeout in seconds. Value range: 1-86400. Default is 1000.
* `ca_certificate_ids` - (Optional) List of CA certificate IDs of the listener. Currently, only 1 CA certificate is supported. Required when `ca_enabled` is `true`.
* `ca_enabled` - (Optional) Whether to enable mutual authentication. Default is `false`.
* `certificate_ids` - (Optional) List of server certificate IDs.
* `gzip_enabled` - (Optional) Whether to enable Gzip compression. Default is `true`.
* `http2enabled` - (Optional) Whether to enable HTTP/2. HTTP protocol defaults to `false`, HTTPS protocol defaults to `true`. Only HTTPS protocol supports this parameter.
* `idle_timeout` - (Optional) Idle connection timeout in seconds. Value range: 1-600. Default is 15.
* `listener_name` - (Optional) Custom listener name. Length is 1-255 characters, which can contain Chinese characters, letters, digits, hyphens (-), slashes (/), periods (.) and underscores (_).
* `request_timeout` - (Optional) Request timeout in seconds. Value range: 1-600. Default is 60. If the backend server does not respond within the timeout, the ALB aborts waiting and returns HTTP 504 to the client.
* `security_policy_id` - (Optional) Security policy ID, in the format of `tls-xxxxxxxx`.
* `tags` - (Optional) List of tags. Up to 20 tags are supported. Each tag contains `tag_key` (required) and `tag_value` (required).
* `x_forwarded_for_config` - (Optional) X-Forwarded-For configuration.
  * `x_forwarded_for_alb_id_enabled` - (Optional) Whether to get the ALB instance ID via the ALB-ID header.
  * `x_forwarded_for_client_src_port_enabled` - (Optional) Whether to get the client port via the X-Forwarded-Client-srcport header.
  * `x_forwarded_for_host_enabled` - (Optional) Whether to get the client domain via the X-Forwarded-Host header.
  * `x_forwarded_for_mode` - (Optional) How to handle the X-Forwarded-For (XFF) HTTP header. Valid values: `append` (default), `remove` and `passthrough`.
  * `x_forwarded_for_port_enabled` - (Optional) Whether to get the listener port via the X-Forwarded-Port header.
  * `x_forwarded_for_proto_enabled` - (Optional) Whether to get the listener protocol via the X-Forwarded-Proto header.
  * `x_tencent_client_idn_enabled` - (Optional) Whether to access the issuer of the client certificate $ssl_client_i_dn via the X-Tencent-Client-IDN header.
  * `x_tencent_client_sdn_enabled` - (Optional) Whether to access the subject of the client certificate $ssl_client_s_dn via the X-Tencent-Client-SDN header.
  * `x_tencent_client_serial_enabled` - (Optional) Whether to access the serial number of the client certificate $ssl_client_serial via the X-Tencent-Client-Serial header.
  * `x_tencent_client_verify_enabled` - (Optional) Whether to access the verification result of the client certificate $ssl_client_verify via the X-Tencent-Client-Verify header.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, in the format of `loadBalancerId#listenerId`.

## Import

ALB listener can be imported using the loadBalancerId#listenerId, e.g.

```
terraform import tencentcloud_alb_listener.example alb-xxxxxxxx#lst-xxxxxxxx
```
