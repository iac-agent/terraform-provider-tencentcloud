---
subcategory: "Application Load Balancer(ALB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_listener"
sidebar_current: "docs-tencentcloud-resource-alb_listener"
description: |-
  Provides a resource to manage a listener forwarding rule.
---

# tencentcloud_alb_listener

Provides a resource to manage a listener forwarding rule.

## Example Usage

```hcl
resource "tencentcloud_alb_listener" "example" {
  listener_id      = "lst-12345678"
  load_balancer_id = "alb-12345678"
  dry_run          = false

  rules = [
    {
      priority  = 10
      rule_name = "tf-test-rule"

      conditions = [
        {
          type        = "Host"
          host_config = ["example.com"]
        }
      ]

      actions = [
        {
          order = 1
          type  = "TargetGroup"
          target_group_config = {
            target_groups = [
              {
                target_group_id = "lbtg-12345678"
                weight          = 10
              }
            ]
            target_group_sticky_session = {
              enabled = true
              timeout = 1000
            }
          }
        }
      ]
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, String) The listener ID, in the format lst- followed by 8 alphanumeric characters.
* `load_balancer_id` - (Required, String) The load balancer instance ID, in the format alb- followed by 8 alphanumeric characters.
* `rules` - (Required, List) The forwarding rule list. The resource manages a single rule; a list with exactly one rule is expected.
* `client_token` - (Optional, String) The client token used to ensure request idempotency. Generate a unique value on the client so different requests carry different tokens. Only ASCII characters are supported. If unspecified, the API uses the request ID as the client token.
* `dry_run` - (Optional, Bool) Whether to only perform a pre-check for this request.

The `rules` object supports the following:

* `actions` - (Required, List) The forwarding rule action list.
* `conditions` - (Optional, List) The forwarding rule condition list.
* `direction` - (Optional, String) The direction of the forwarding rule. Request: from the client to the load balancer; Response: from the backend server to the load balancer. Default Request.
* `priority` - (Optional, Int) The priority. A smaller value means a higher priority; values cannot repeat. Valid values: 1 to 10000.
* `rule_name` - (Optional, String) The forwarding rule name, 1 to 255 characters long, supporting digits, upper and lower case letters, Chinese characters, periods (.), underscores (_) and hyphens (-).
* `tags` - (Optional, List) The tags.

The `actions` object of `rules` supports the following:

* `order` - (Required, Int) The execution order of the forwarding action. Orders cannot repeat; actions are executed in ascending order. Valid values: 1 to 50000.
* `type` - (Required, String) The forwarding action type. Valid values: TargetGroup (forward to a target group), Redirect (redirect), FixedResponse (return a fixed response), Rewrite (rewrite), InsertHeader (write an HTTP header), RemoveHeader (remove an HTTP header). A forwarding action must include TargetGroup, Redirect, FixedResponse or Rewrite.
* `fixed_response_config` - (Optional, Object) The fixed response content configuration.
* `insert_header_config` - (Optional, Object) The inserted HTTP Header configuration.
* `redirect_config` - (Optional, Object) The redirection configuration. Except for HttpCode, the other configurations cannot all use default values.
* `remove_header_config` - (Optional, Object) The removed HTTP Header configuration.
* `rewrite_config` - (Optional, Object) The rewrite configuration.
* `target_group_config` - (Optional, Object) The forwarding target group configuration.

The `fixed_response_config` object of `actions` supports the following:

* `http_code` - (Required, Int) The HTTP response code returned, supporting 2xx, 4xx, 5xx.
* `content_type` - (Optional, String) The format of the fixed content. Valid values: text/plain, text/css, text/html, application/javascript, or application/json.
* `content` - (Optional, String) The fixed content returned. Only ASCII characters are supported, up to 1 KB.

The `insert_header_config` object of `actions` supports the following:

* `key` - (Required, String) The key of the HTTP Header to insert, 1 to 40 characters long, supporting a-z A-Z 0-9 - _ . Chinese characters and Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-* are not supported.
* `value_type` - (Required, String) The type of the HTTP Header value. Valid values: SystemDefined (system-defined header), UserDefined (user-defined header), ReferenceHeader (reference a header from the request).
* `value` - (Required, String) The value of the HTTP Header. When ValueType is SystemDefined, valid values are ClientPort (client port), ClientIp (client IP address), Protocol (client request protocol), CLBPort (load balancer listening port). When ValueType is UserDefined, 1 to 128 printable characters.

The `redirect_config` object of `actions` supports the following:

* `http_code` - (Required, Int) The HTTP code for redirection. Supported values: 301, 302, 303, 307, 308.
* `host` - (Optional, String) The redirection host, default ${host}. 3 to 128 characters long, supporting a-z 0-9 _ . -.
* `path` - (Optional, String) The redirection path, default ${path}. 1 to 128 characters long, supporting a-z A-Z 0-9 ? = _ . - / : .
* `port` - (Optional, String) The redirection port, default ${port}. Valid values: 1 to 65535.
* `protocol` - (Optional, String) The redirection protocol. Valid values: HTTP, HTTPS, default ${protocol}.
* `query` - (Optional, String) The redirection query string, default ${query}. 1 to 128 printable characters, excluding #[ ]{}|<> & and spaces.

The `remove_header_config` object of `actions` supports the following:

* `key` - (Required, String) The key of the HTTP Header to remove, 1 to 40 characters long, supporting a-z A-Z 0-9 - _ . Chinese characters and Host, Cookie are not supported.

The `rewrite_config` object of `actions` supports the following:

* `host` - (Optional, String) The rewrite host, default ${host}. 3 to 128 characters long, supporting a-z 0-9 _ . -.
* `path` - (Optional, String) The rewrite path, default ${path}. 1 to 128 characters long, supporting a-z A-Z 0-9 ? = _ . - / : .
* `query` - (Optional, String) The rewrite query string, default ${query}. 1 to 128 printable characters, excluding #[ ]{}|<> & and spaces.

The `target_group_config` object of `actions` supports the following:

* `target_groups` - (Required, List) The target group list.
* `target_group_sticky_session` - (Optional, Object) Sticky sessions between target groups.

The `target_group_sticky_session` object of `target_group_config` supports the following:

* `enabled` - (Required, Bool) Whether to enable sticky sessions between target groups. Disabled by default.
* `timeout` - (Optional, Int) The timeout in seconds, ranging from 1 to 86400, default 1000.

The `target_groups` object of `target_group_config` supports the following:

* `target_group_id` - (Required, String) The target group ID, in the format lbtg- followed by 8 alphanumeric characters.
* `weight` - (Optional, Int) The weight, ranging from 0 to 100, default 10.

The `conditions` object of `rules` supports the following:

* `type` - (Required, String) The forwarding condition type. Valid values: Host, Path, Header, QueryString, Method, Cookie, SourceIp.
* `cookie_config` - (Optional, List) The Cookie configuration.
* `header_config` - (Optional, Object) The HTTP Header configuration.
* `host_config` - (Optional, List) The host name. A host configuration can appear only once in a rule, 3 to 128 characters long, supporting exact, regex and wildcard matching. It cannot start or end with a period (.) or underscore (_).
* `method_config` - (Optional, List) The request method. Valid values: HEAD, GET, POST, OPTIONS, PUT, PATCH, DELETE.
* `path_config` - (Optional, List) The forwarding path, 1 to 128 characters long, supporting exact, regex and wildcard matching.
* `query_string_config` - (Optional, List) The query string configuration.
* `source_ip_config` - (Optional, List) The source IP matching configuration in CIDR format, for example x.x.x.x/32 or x.x.x.x/24.

The `cookie_config` object of `conditions` supports the following:

* `key` - (Required, String) The key of the Cookie, 1 to 64 characters long, supporting letters, digits and underscores.
* `value` - (Required, String) The value of the Cookie, 1 to 128 printable characters.

The `header_config` object of `conditions` supports the following:

* `key` - (Required, String) The key of the HTTP Header, 1 to 40 characters long, supporting a-z A-Z 0-9 - _ . Chinese characters, Host and Cookie are not supported.
* `values` - (Required, List) The values of the HTTP Header, 1 to 128 printable characters. Double quotes are not supported; leading and trailing spaces are not allowed, and a trailing backslash is not allowed.

The `query_string_config` object of `conditions` supports the following:

* `key` - (Required, String) The key of the query string, 1 to 16 characters long. Printable characters are supported, excluding spaces and #[]{}\|<>&. * matches multiple characters and ? matches a single character.
* `value` - (Required, String) The value of the query string, 1 to 128 characters long. Printable characters are supported, excluding spaces and #[]{}\|<>&. * matches multiple characters and ? matches a single character.

The `tags` object of `rules` supports the following:

* `tag_key` - (Required, String) The key of the tag.
* `tag_value` - (Required, String) The value of the tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, the ID of the first forwarding rule created.

The `rules` object exports the following:

* `create_time` - The creation time.
* `modify_time` - The last modification time.
* `rule_id` - The forwarding rule ID, in the format rule- followed by 8 alphanumeric characters.
* `status` - The forwarding rule status. Provisioning: creating; Active: running; Configuring: configuring.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateRules](https://www.tencentcloud.com/document/api/1311/84277)
  - [DescribeRules](https://www.tencentcloud.com/document/api/1311/84275)
* Read: [DescribeRules](https://www.tencentcloud.com/document/api/1311/84275)
* Update:
  - [DescribeRules](https://www.tencentcloud.com/document/api/1311/84275)
  - [ModifyRulesAttributes](https://www.tencentcloud.com/document/api/1311/84274)
* Delete: [DeleteRules](https://www.tencentcloud.com/document/api/1311/84276)

## Import

ALB listener can be imported using the id, e.g.

```
terraform import tencentcloud_alb_listener.example rule-12345678
```

