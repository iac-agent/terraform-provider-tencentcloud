---
subcategory: "Application Load Balancer(ALB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener"
sidebar_current: "docs-tencentcloud-resource-alb_listener"
description: |-
  Provides a resource to manage an Application Load Balancer (ALB) listener.
---

# tencentcloud_alb_listener

Provides a resource to manage an Application Load Balancer (ALB) listener.

## Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "alb-12345678"
  listener_name     = "tf-test-listener"
  listener_port     = 80
  listener_protocol = "HTTP"

  default_actions {
    type = "TargetGroup"

    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }
    }
  }

  gzip_enabled    = true
  idle_timeout    = 15
  request_timeout = 60
}
```

## Argument Reference

The following arguments are supported:

* `default_actions` - (Required, List) List of default forwarding actions. Currently a listener supports exactly one default forwarding action, and the only supported action type is forwarding to a target group.
* `listener_port` - (Required, Int) Front-end port used by the listener. Valid values: 1 to 65535. It cannot be changed after creation.
* `listener_protocol` - (Required, String) Protocol of the listener. Valid values: `HTTP`, `HTTPS` and `QUIC`. It cannot be changed after creation.
* `load_balancer_id` - (Required, String) ID of the ALB instance this listener belongs to, in the form `alb-xxxxxxxx`. It cannot be changed after creation.
* `ca_certificate_ids` - (Optional, Set) List of CA certificate IDs configured on the listener. Currently only one CA certificate is supported, and it is required when `ca_enabled` is `true`.
* `ca_enabled` - (Optional, Bool) Whether mutual authentication (two-way authentication) is enabled. Default value: `false`.
* `certificate_ids` - (Optional, Set) List of server certificate IDs of the listener.
* `client_token` - (Optional, String) Client token used to guarantee the idempotency of the create request. It must be unique among different requests and can only contain ASCII characters. It is only used when the listener is created.
* `gzip_enabled` - (Optional, Bool) Whether Gzip compression is enabled. Default value: `true`.
* `http2enabled` - (Optional, Bool) Whether HTTP/2 is enabled. It is only supported by the `HTTPS` protocol. Default value: `false` for HTTP and `true` for HTTPS.
* `idle_timeout` - (Optional, Int) Connection idle timeout in seconds. Valid values: 1 to 600. Default value: 15.
* `listener_name` - (Optional, String) Name of the listener. Length is 1 to 255 characters; it can contain Chinese characters, letters, digits, hyphens (-), slashes (/), periods (.) and underscores (_).
* `request_timeout` - (Optional, Int) Connection request timeout in seconds. Valid values: 1 to 600. Default value: 60.
* `security_policy_id` - (Optional, String) ID of the TLS security policy, in the form `tls-xxxxxxxx`.
* `tags` - (Optional, List) Tags of the listener. Up to 20 tags are supported. The listener modify API does not accept tags, so changing them forces the listener to be recreated.
* `x_forwarded_for_config` - (Optional, List) X-Forwarded-For configuration of the listener.

The `default_actions` object supports the following:

* `target_group_config` - (Required, List) Target group configuration of the forwarding action. Only a single target group is supported when the listener is created.
* `type` - (Required, String) Type of the forwarding action. Only `TargetGroup` (forward to a target group) is supported for the default action of a listener.

The `target_group_config` object of `default_actions` supports the following:

* `target_groups` - (Required, List) List of target groups.
* `target_group_sticky_session` - (Optional, List) Session persistence between the client and the target group.

The `target_group_sticky_session` object of `target_group_config` supports the following:

* `enabled` - (Required, Bool) Whether session persistence is enabled. It is disabled by default.
* `timeout` - (Optional, Int) Session persistence timeout in seconds. Valid values: 1 to 86400. Default value: 1000.

The `target_groups` object of `target_group_config` supports the following:

* `target_group_id` - (Required, String) ID of the target group, in the form `lbtg-xxxxxxxx`.
* `weight` - (Optional, Int) Weight of the target group. Valid values: 0 to 100. Default value: 10.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

The `x_forwarded_for_config` object supports the following:

* `x_forwarded_for_alb_id_enabled` - (Optional, Bool) Whether the ALB instance ID is obtained through the ALB-ID header.
* `x_forwarded_for_client_src_port_enabled` - (Optional, Bool) Whether the client port accessing the ALB instance is obtained through the X-Forwarded-Client-srcport header.
* `x_forwarded_for_host_enabled` - (Optional, Bool) Whether the domain name of the client accessing the ALB instance is obtained through the X-Forwarded-Host header.
* `x_forwarded_for_mode` - (Optional, String) How the X-Forwarded-For (XFF) HTTP header is handled. Valid values: `append` (default, append the client IP to the end of the header), `remove` (drop the header) and `passthrough` (forward the header unchanged).
* `x_forwarded_for_port_enabled` - (Optional, Bool) Whether the listener port of the ALB instance is obtained through the X-Forwarded-Port header.
* `x_forwarded_for_proto_enabled` - (Optional, Bool) Whether the listener protocol of the ALB instance is obtained through the X-Forwarded-Proto header.
* `x_tencent_client_idn_enabled` - (Optional, Bool) Whether the issuer `$ssl_client_i_dn` of the client certificate is obtained through the X-Tencent-Client-IDN header.
* `x_tencent_client_sdn_enabled` - (Optional, Bool) Whether the subject `$ssl_client_s_dn` of the client certificate is obtained through the X-Tencent-Client-SDN header.
* `x_tencent_client_serial_enabled` - (Optional, Bool) Whether the serial number `$ssl_client_serial` of the client certificate is obtained through the X-Tencent-Client-Serial header.
* `x_tencent_client_verify_enabled` - (Optional, Bool) Whether the verification result `$ssl_client_verify` of the client certificate is obtained through the X-Tencent-Client-Verify header.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the listener in ISO 8601 format.
* `id` - Resource identifier in the form `<load_balancer_id>#<listener_id>`.
* `listener_id` - ID of the listener, in the form `lst-xxxxxxxx`.
* `listener_status` - Status of the listener. Valid values: `Active` (running), `Provisioning` (being created), `Configuring` (being modified) and `ProvisionFailed` (failed to be created).
* `modify_time` - Last modification time of the listener in ISO 8601 format.

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

