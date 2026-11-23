---
subcategory: "Application Load Balancer(ALB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener"
sidebar_current: "docs-tencentcloud-resource-alb_listener"
description: |-
  Provides a resource to create an alb listener
---

# tencentcloud_alb_listener

Provides a resource to create an alb listener

## Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-xxxxxxxx"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-example"

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

  tags {
    tag_key   = "test_key"
    tag_value = "test_value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `default_actions` - (Required, List) List of default forwarding rule actions. Currently a listener supports only 1 default forwarding rule action.
* `listener_port` - (Required, Int) Port used by the frontend of the load balancer instance. Valid value range: 1~65535.
* `listener_protocol` - (Required, String) Listener protocol. Valid values: HTTP, HTTPS or QUIC.
* `load_balancer_id` - (Required, String) ID of the load balancer instance, format is alb- followed by 8 alphanumeric characters.
* `ca_certificate_ids` - (Optional, List) List of CA certificate IDs configured for the listener. Currently only 1 CA certificate is supported. Required when `ca_enabled` is true.
* `ca_enabled` - (Optional, Bool) Whether to enable mutual authentication. true: enable. false (default): disable.
* `certificate_ids` - (Optional, List) List of server certificate IDs.
* `client_token` - (Optional, String) Client token used to ensure request idempotency. Only ASCII characters are supported.
* `gzip_enabled` - (Optional, Bool) Whether to enable Gzip compression. true (default): yes. false: no.
* `http2enabled` - (Optional, Bool) Whether to enable HTTP/2. Default false for HTTP protocol, default true for HTTPS protocol. Only HTTPS protocol supports this parameter.
* `idle_timeout` - (Optional, Int) Connection idle timeout in seconds. Valid value range: 1~600. Default is 15.
* `listener_name` - (Optional, String) Custom listener name. Length is 1~255 characters, can contain Chinese characters, letters, digits, hyphens (-), slashes (/), periods (.) and underscores (_).
* `request_timeout` - (Optional, Int) Connection request timeout in seconds. Valid value range: 1~600. Default is 60.
* `security_policy_id` - (Optional, String) Security policy ID, format is tls- followed by 8 alphanumeric characters.
* `tags` - (Optional, List) List of tags. Up to 20 supported.
* `x_forwarded_for_config` - (Optional, List) X-Forwarded-For configuration.

The `default_actions` object supports the following:

* `type` - (Required, String) Forwarding action type. Only forwarding to target group is supported when creating a listener.
* `target_group_config` - (Required, List) Forwarding target group configuration. Only a single target group is supported.

The `target_group_config` object supports the following:

* `target_groups` - (Required, List) List of target groups.
* `target_group_sticky_session` - (Optional, List) Sticky session between target groups.

The `target_groups` object supports the following:

* `target_group_id` - (Required, String) Target group ID, format is lbtg- followed by 8 alphanumeric characters.
* `weight` - (Optional, Int) Weight, valid value range: [0, 100], default is 10.

The `target_group_sticky_session` object supports the following:

* `enabled` - (Required, Bool) Whether to enable sticky session, default is disabled.
* `timeout` - (Optional, Int) Timeout in seconds, valid value range: 1~86400, default is 1000.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

The `x_forwarded_for_config` object supports the following:

* `x_forwarded_for_alb_id_enabled` - (Optional, Bool) Whether to obtain the load balancer instance ID via the ALB-ID header field.
* `x_forwarded_for_client_src_port_enabled` - (Optional, Bool) Whether to obtain the client port via the X-Forwarded-Client-srcport header field.
* `x_forwarded_for_host_enabled` - (Optional, Bool) Whether to obtain the client domain via the X-Forwarded-Host header field.
* `x_forwarded_for_mode` - (Optional, String) Specifies how to handle the X-Forwarded-For (XFF) HTTP header field. Valid values: append, remove, passthrough.
* `x_forwarded_for_port_enabled` - (Optional, Bool) Whether to obtain the listener port via the X-Forwarded-Port header field.
* `x_forwarded_for_proto_enabled` - (Optional, Bool) Whether to obtain the listener protocol via the X-Forwarded-Proto header field.
* `x_tencent_client_idn_enabled` - (Optional, Bool) Whether to access the issuer of the client certificate $ssl_client_i_dn via the X-Tencent-Client-IDN header.
* `x_tencent_client_sdn_enabled` - (Optional, Bool) Whether to access the subject of the client certificate $ssl_client_s_dn via the X-Tencent-Client-SDN header.
* `x_tencent_client_serial_enabled` - (Optional, Bool) Whether to access the serial number of the client certificate $ssl_client_serial via the X-Tencent-Client-Serial header.
* `x_tencent_client_verify_enabled` - (Optional, Bool) Whether to access the verification result of the client certificate $ssl_client_verify via the X-Tencent-Client-Verify header.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource (listener ID).
* `listener_id` - Listener ID, format is lst- followed by 8 alphanumeric characters.
* `listener_status` - Listener status. Valid values: Active, Provisioning, Configuring, ProvisionFailed.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateListener](https://cloud.tencent.com/document/api/1822/133705)
* Read: [DescribeListeners](https://cloud.tencent.com/document/api/1822/133700)
* Update: [ModifyListenerAttributes](https://cloud.tencent.com/document/api/1822/133698)
* Delete: [DeleteListener](https://cloud.tencent.com/document/api/1822/133704)

## Import

alb listener can be imported using the id, e.g.

```
terraform import tencentcloud_alb_listener.example alb-xxxxxxxx:lst-xxxxxxxx
```
