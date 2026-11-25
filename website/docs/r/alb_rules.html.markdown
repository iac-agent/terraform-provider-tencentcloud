---
subcategory: "Application Load Balancer(alb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_rules"
sidebar_current: "docs-tencentcloud-resource-alb_rules"
description: |-
  Provides a resource to manage a set of Application Load Balancer forwarding rules.
---

# tencentcloud_alb_rules

Provides a resource to manage a set of Application Load Balancer forwarding rules.

## Example Usage

```hcl
resource "tencentcloud_alb_rules" "example" {
  listener_id      = "lst-12345678"
  load_balancer_id = "alb-12345678"

  rules = [
    {
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
              enabled = false
              timeout = 1000
            }
          }
        }
      ]
      conditions = [
        {
          type        = "Host"
          host_config = ["example.com"]
        }
      ]
      priority  = 1
      rule_name = "tf-test-rule"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, String) The listener ID, in the format lst- followed by 8 alphanumeric characters.
* `load_balancer_id` - (Required, String) The load balancer instance ID, in the format alb- followed by 8 alphanumeric characters.
* `rules` - (Required, List) Forwarding rule list.
* `client_token` - (Optional, String) The client token used to guarantee request idempotency. Generate a unique value on your client so that it is unique across requests. ClientToken supports only ASCII characters. If it is not specified, the API request RequestId is used as the ClientToken. Each API request has a different RequestId.
* `dry_run` - (Optional, Bool) Whether to only perform a pre-check of this request.

The `rules` object supports the following:

* `actions` - (Required, List) Forwarding rule action list.
* `conditions` - (Required, List) Forwarding rule condition list.
* `priority` - (Required, Int) The priority. A smaller value means higher priority. Priorities cannot be duplicated. Valid values: 1 ~ 10000.
* `direction` - (Optional, String) The direction of the forwarding rule. Request: requests from the client to the load balancer. Response: responses from the backend server to the load balancer. Default value: Request.
* `rule_name` - (Optional, String) The name of the forwarding rule, 1 ~ 255 characters long. It can contain digits, uppercase and lowercase letters, Chinese characters, periods (.), underscores (_) and hyphens (-).
* `tags` - (Optional, List) Tags.

The `actions` object of `rules` supports the following:

* `order` - (Required, Int) The order in which forwarding actions are executed. It cannot be repeated and actions are executed in ascending order. Valid values: 1 ~ 50000.
* `type` - (Required, String) The forwarding action type. Valid values: TargetGroup: forward to a target group. Redirect: redirect. FixedResponse: return a fixed response. Rewrite: rewrite. InsertHeader: write an HTTP header. RemoveHeader: delete an HTTP header. A forwarding action must include one of TargetGroup, Redirect and FixedResponse, and it must be executed last.
* `fixed_response_config` - (Optional, Object) Fixed response content configuration.
* `insert_header_config` - (Optional, Object) Insert HTTP header configuration.
* `redirect_config` - (Optional, Object) Redirection configuration. Except for HttpCode, the other configurations cannot all use default values.
* `remove_header_config` - (Optional, Object) Remove HTTP header configuration.
* `rewrite_config` - (Optional, Object) Rewrite configuration.
* `target_group_config` - (Optional, Object) Forwarding target group configuration.

The `fixed_response_config` object of `actions` supports the following:

* `http_code` - (Required, Int) The HTTP response code to return. 2xx, 4xx and 5xx are supported.
* `content_type` - (Optional, String) The format of the fixed content to return. Valid values: text/plain, text/css, text/html, application/javascript or application/json.
* `content` - (Optional, String) The fixed content to return. Only ASCII characters are supported, up to 1 KB.

The `insert_header_config` object of `actions` supports the following:

* `key` - (Required, String) The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.
* `value_type` - (Required, String) The value type of the HTTP header. Valid values: SystemDefined: a system-defined header. UserDefined: a user-defined header. ReferenceHeader: a header referenced from the request.
* `value` - (Required, String) The value of the HTTP header. When ValueType is SystemDefined, valid values are ClientPort: client port, ClientIp: client IP address, Protocol: protocol of the client request, CLBPort: load balancer listener port. When ValueType is UserDefined, it is 1 ~ 128 printable characters, does not support ", cannot start or end with a space, and cannot end with \. When ValueType is ReferenceHeader, it references a header in the request, 1 ~ 128 printable characters, does not support ", cannot start or end with a space, and cannot end with \.

The `redirect_config` object of `actions` supports the following:

* `http_code` - (Required, Int) The HTTP response code to return. 2xx, 4xx and 5xx are supported.
* `host` - (Optional, String) The host address for redirect/rewrite, with a default value of $${host}. It is 3 ~ 128 characters long. Supported characters: a-z 0-9 _ . -.
* `path` - (Optional, String) The path for redirect/rewrite, with a default value of $${path}. It is 1 ~ 128 characters long. Supported characters: a-z A-Z 0-9 ? = _ . - / :.
* `port` - (Optional, String) The port for redirection, with a default value of $${port}. Valid values: 1 ~ 65535.
* `protocol` - (Optional, String) The protocol for redirection. Valid values: HTTP, HTTPS. Default value: $${protocol}.
* `query` - (Optional, String) The query string for redirect/rewrite, with a default value of $${query}. It is 1 ~ 128 characters long and supports printable characters, but not #[]{}|<>& or spaces.

The `remove_header_config` object of `actions` supports the following:

* `key` - (Required, String) The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.

The `rewrite_config` object of `actions` supports the following:

* `host` - (Optional, String) The host address for redirect/rewrite, with a default value of $${host}. It is 3 ~ 128 characters long. Supported characters: a-z 0-9 _ . -.
* `path` - (Optional, String) The path for redirect/rewrite, with a default value of $${path}. It is 1 ~ 128 characters long. Supported characters: a-z A-Z 0-9 ? = _ . - / :.
* `query` - (Optional, String) The query string for redirect/rewrite, with a default value of $${query}. It is 1 ~ 128 characters long and supports printable characters, but not #[]{}|<>& or spaces.

The `target_group_config` object of `actions` supports the following:

* `target_groups` - (Required, List) Target group list.
* `target_group_sticky_session` - (Optional, Object) Sticky session between target groups.

The `target_group_sticky_session` object of `target_group_config` supports the following:

* `enabled` - (Required, Bool) Whether to enable sticky sessions. Disabled by default.
* `timeout` - (Optional, Int) The timeout period in seconds. Valid values: 1 ~ 86400. Default value: 1000.

The `target_groups` object of `target_group_config` supports the following:

* `target_group_id` - (Required, String) The target group ID, in the format lbtg- followed by 8 alphanumeric characters.
* `weight` - (Optional, Int) The weight. Valid values: [0, 100]. Default value: 10.

The `conditions` object of `rules` supports the following:

* `type` - (Required, String) The forwarding action type. Valid values: TargetGroup: forward to a target group. Redirect: redirect. FixedResponse: return a fixed response. Rewrite: rewrite. InsertHeader: write an HTTP header. RemoveHeader: delete an HTTP header. A forwarding action must include one of TargetGroup, Redirect and FixedResponse, and it must be executed last.
* `cookie_config` - (Optional, List) Cookie configuration.
* `header_config` - (Optional, Object) HTTP header configuration.
* `host_config` - (Optional, List) Host name. A host configuration can appear only once in a rule and is 3 ~ 128 characters long, supporting exact, regular expression and wildcard matching. It cannot start or end with a period (.) or underscore (_). Exact match supported character set: a-z 0-9 . - _. Regular expression match: starts with a tilde (~); supported characters: a-z 0-9 . - ? = ~ _ - + \ ^ * ! $ & | ( ) [ ]. Wildcard match: * matches multiple characters and ? matches one character; supported characters: a-z 0-9 . - _ * ?.
* `method_config` - (Optional, List) The request method. Valid values: HEAD, GET, POST, OPTIONS, PUT, PATCH, DELETE.
* `path_config` - (Optional, List) The forwarding path, 1 ~ 128 characters long, supporting exact, regular expression and wildcard matching. Exact match supported character set: a-z A-Z 0-9 . - _ / = :. Regular expression match: starts with ~ (case-sensitive) or ~* (case-insensitive); supported characters: a-z A-Z 0-9 . - _ / = ? ~ ^ * $ : ( ) [ ] + |. Wildcard match: * matches multiple characters and ? matches one character; supported characters: a-z A-Z 0-9 . - _ / = :.
* `query_string_config` - (Optional, List) Query string configuration.
* `source_ip_config` - (Optional, List) Source IP matching configuration, in CIDR format: IP address x.x.x.x/32 or IP range x.x.x.x/24.

The `cookie_config` object of `conditions` supports the following:

* `key` - (Required, String) The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.
* `value` - (Required, String) The value of the HTTP header. When ValueType is SystemDefined, valid values are ClientPort: client port, ClientIp: client IP address, Protocol: protocol of the client request, CLBPort: load balancer listener port. When ValueType is UserDefined, it is 1 ~ 128 printable characters, does not support ", cannot start or end with a space, and cannot end with \. When ValueType is ReferenceHeader, it references a header in the request, 1 ~ 128 printable characters, does not support ", cannot start or end with a space, and cannot end with \.

The `header_config` object of `conditions` supports the following:

* `key` - (Required, String) The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.
* `values` - (Required, List) The value of the HTTP header, 1 ~ 128 characters long, supporting printable characters. " is not supported; it cannot start or end with a space and cannot end with \.

The `query_string_config` object of `conditions` supports the following:

* `key` - (Required, String) The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.
* `value` - (Required, String) The value of the HTTP header. When ValueType is SystemDefined, valid values are ClientPort: client port, ClientIp: client IP address, Protocol: protocol of the client request, CLBPort: load balancer listener port. When ValueType is UserDefined, it is 1 ~ 128 printable characters, does not support ", cannot start or end with a space, and cannot end with \. When ValueType is ReferenceHeader, it references a header in the request, 1 ~ 128 printable characters, does not support ", cannot start or end with a space, and cannot end with \.

The `tags` object of `rules` supports the following:

* `tag_key` - (Required, String) The tag key.
* `tag_value` - (Required, String) The tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.

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

Alb rules can be imported using the id, e.g.

```
terraform import tencentcloud_alb_rules.example rule-12345678
```

