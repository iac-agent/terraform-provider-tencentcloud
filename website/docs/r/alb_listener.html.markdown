---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener"
sidebar_current: "docs-tencentcloud-resource-alb_listener"
description: |-
  Provides a resource to manage the forwarding rules of an ALB listener.
---

# tencentcloud_alb_listener

Provides a resource to manage the forwarding rules of an ALB listener.

## Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  listener_id      = "lst-d9p3k7wa"
  load_balancer_id = "alb-f8q2xk9m"

  rules = [{
    priority  = 10
    rule_name = "tf-test-rule"
    direction = "Request"

    actions = [{
      order = 1
      type  = "TargetGroup"
      target_group_config = {
        target_groups = [{
          target_group_id = "lbtg-a1b2c3d4"
          weight          = 10
        }]
      }
    }]

    conditions = [{
      type        = "Host"
      host_config = ["www.example.com"]
    }]

    tags = [{
      tag_key   = "env"
      tag_value = "test"
    }]
  }]
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, String) Listener ID, in the format of lst- followed by 8 alphanumeric characters.
* `load_balancer_id` - (Required, String) Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters.
* `rules` - (Required, List) Forwarding rule list.

The `rules` object supports the following:

* `actions` - (Required, List) Forwarding rule action list.
* `conditions` - (Required, List) Forwarding rule condition list.
* `priority` - (Required, Int) Priority. Smaller values indicate higher priority; values cannot be duplicated. Value range: 1 ~ 10000.
* `direction` - (Optional, String) Forwarding rule direction. Request: client to load balancer; Response: backend server to load balancer. Default Request.
* `rule_name` - (Optional, String) Forwarding rule name, 1 ~ 255 characters, allowing digits, letters, Chinese characters, periods (.), underscores (_) and hyphens (-).
* `tags` - (Optional, List) Tags.

The `actions` object of `rules` supports the following:

* `order` - (Required, Int) Action execution order. Values cannot be duplicated and actions are executed in ascending order. Value range: 1 ~ 50000.
* `type` - (Required, String) Forwarding action type. Valid values: TargetGroup, Redirect, FixedResponse, Rewrite, InsertHeader, RemoveHeader. An action list must contain one of TargetGroup, Redirect or FixedResponse, placed last in execution order.
* `fixed_response_config` - (Optional, Object) Fixed response content configuration.
* `insert_header_config` - (Optional, Object) Insert HTTP header configuration.
* `redirect_config` - (Optional, Object) Redirect configuration. Except for HttpCode, the other settings cannot all use default values.
* `remove_header_config` - (Optional, Object) Remove HTTP header configuration.
* `rewrite_config` - (Optional, Object) Rewrite configuration.
* `target_group_config` - (Optional, Object) Forwarding target group configuration.

The `fixed_response_config` object of `actions` supports the following:

* `http_code` - (Required, Int) Returned HTTP response code. 2xx, 4xx and 5xx are supported.
* `content_type` - (Optional, String) Format of the returned fixed content. Valid values: text/plain, text/css, text/html, application/javascript, application/json.
* `content` - (Optional, String) Returned fixed content. Only ASCII characters are supported, up to 1 KB.

The `insert_header_config` object of `actions` supports the following:

* `key` - (Required, String) Key of the HTTP header to insert, 1 ~ 40 characters, supporting a-z A-Z 0-9 - _. Chinese characters and Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port are not supported.
* `value_type` - (Required, String) HTTP header value type. Valid values: SystemDefined, UserDefined, ReferenceHeader.
* `value` - (Required, String) HTTP header value. When ValueType is SystemDefined, valid values are ClientPort, ClientIp, Protocol and CLBPort. When ValueType is UserDefined, 1 ~ 128 printable characters without double quotes, leading or trailing spaces, or a trailing backslash. When ValueType is ReferenceHeader, references a request header with the same constraints.

The `redirect_config` object of `actions` supports the following:

* `http_code` - (Required, Int) HTTP redirect code. 301, 302, 303, 307 and 308 are supported.
* `host` - (Optional, String) Redirect host address, default ${host}. 3 ~ 128 characters, supporting a-z 0-9 _ . -.
* `path` - (Optional, String) Redirect path, default ${path}. 1 ~ 128 characters, supporting a-z A-Z 0-9 ? = _ . - / : .
* `port` - (Optional, String) Redirect port, default ${port}. Value range: 1 ~ 65535.
* `protocol` - (Optional, String) Redirect protocol. Valid values: HTTP, HTTPS. Default ${protocol}.
* `query` - (Optional, String) Redirect query string, default ${query}. 1 ~ 128 printable characters; #[]{}<> and spaces are not supported.

The `remove_header_config` object of `actions` supports the following:

* `key` - (Required, String) Key of the HTTP header to remove, 1 ~ 40 characters, supporting a-z A-Z 0-9 - _.

The `rewrite_config` object of `actions` supports the following:

* `host` - (Optional, String) Rewrite host address, default ${host}. 3 ~ 128 characters, supporting a-z 0-9 _ . -.
* `path` - (Optional, String) Rewrite path, default ${path}. 1 ~ 128 characters, supporting a-z A-Z 0-9 ? = _ . - / : .
* `query` - (Optional, String) Rewrite query string, default ${query}. 1 ~ 128 printable characters; #[]{}<> and spaces are not supported.

The `target_group_config` object of `actions` supports the following:

* `target_groups` - (Required, List) Target group list.
* `target_group_sticky_session` - (Optional, Object) Sticky session between target groups.

The `target_group_sticky_session` object of `target_group_config` supports the following:

* `enabled` - (Required, Bool) Whether to enable sticky sessions. Disabled by default.
* `timeout` - (Optional, Int) Timeout in seconds, value range 1 ~ 86400, default 1000.

The `target_groups` object of `target_group_config` supports the following:

* `target_group_id` - (Required, String) Target group ID, in the format of lbtg- followed by 8 alphanumeric characters.
* `weight` - (Optional, Int) Weight, value range [0, 100], default 10.

The `conditions` object of `rules` supports the following:

* `type` - (Required, String) Forwarding condition type. Valid values: Host, Path, Header, QueryString, Method, Cookie, SourceIp.
* `cookie_config` - (Optional, List) Cookie configuration.
* `header_config` - (Optional, Object) HTTP header configuration.
* `host_config` - (Optional, List) Host names. A host can be configured only once per rule; 3 ~ 128 characters, supporting exact, regular-expression and wildcard matching.
* `method_config` - (Optional, List) Request methods. Valid values: HEAD, GET, POST, OPTIONS, PUT, PATCH, DELETE.
* `path_config` - (Optional, List) Forwarding paths, 1 ~ 128 characters, supporting exact, regular-expression and wildcard matching.
* `query_string_config` - (Optional, List) Query string configuration.
* `source_ip_config` - (Optional, List) Source IP match configuration in CIDR format, for example x.x.x.x/32 or x.x.x.x/24.

The `cookie_config` object of `conditions` supports the following:

* `key` - (Required, String) Cookie key, 1 ~ 64 characters, supporting letters, digits and underscores.
* `value` - (Required, String) Cookie value, 1 ~ 128 printable characters.

The `header_config` object of `conditions` supports the following:

* `key` - (Required, String) HTTP header key, 1 ~ 40 characters, supporting a-z A-Z 0-9 - _.
* `values` - (Required, List) HTTP header values, 1 ~ 128 printable characters each.

The `query_string_config` object of `conditions` supports the following:

* `key` - (Required, String) Query string key, 1 ~ 16 characters; printable characters except spaces and #[]{}|\<>&. * and ? are supported as wildcards.
* `value` - (Required, String) Query string value, 1 ~ 128 characters; printable characters except spaces and #[]{}|\<>&. * and ? are supported as wildcards.

The `tags` object of `rules` supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, taken from the first rule ID returned by CreateRules.

The `rules` object exports the following:

* `create_time` - Creation time.
* `modify_time` - Last modification time.
* `rule_id` - Forwarding rule ID, in the format of rule- followed by 8 alphanumeric characters.
* `status` - Forwarding rule status. Provisioning: creating; Active: running; Configuring: configuring.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateRules](https://www.tencentcloud.com/document/api/1311/84277)
* Update: [ModifyRulesAttributes](https://www.tencentcloud.com/document/api/1311/84274)
* Delete: [DeleteRules](https://www.tencentcloud.com/document/api/1311/84276)

## Import

ALB listener can be imported using the id, e.g.

```
terraform import tencentcloud_alb_listener.example rule-a1b2c3d4
```

