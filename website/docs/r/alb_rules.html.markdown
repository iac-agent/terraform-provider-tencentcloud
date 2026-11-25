---
subcategory: "Application Load Balancer(ALB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_rules"
sidebar_current: "docs-tencentcloud-resource-alb_rules"
description: |-
  Provides a resource to manage a forwarding rule of an ALB listener.
---

# tencentcloud_alb_rules

Provides a resource to manage a forwarding rule of an ALB listener.

## Example Usage

```hcl
resource "tencentcloud_alb_rules" "example" {
  load_balancer_id = "alb-12345678"
  listener_id      = "lst-12345678"
  rule_name        = "tf-test-rule"
  priority         = 10
  direction        = "Request"

  actions {
    order = 1

    type = "TargetGroup"

    target_group_config {
      target_groups {
        target_group_id = "lbtg-12345678"
        weight          = 10
      }

      target_group_sticky_session {
        enabled = true
        timeout = 1000
      }
    }
  }

  conditions {
    type        = "Host"
    host_config = ["www.example.com"]
  }

  tags {
    tag_key   = "env"
    tag_value = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `actions` - (Required, List) Forwarding action list.
* `conditions` - (Required, List) Forwarding condition list.
* `listener_id` - (Required, String) Listener ID, in the format of lst- followed by 8 alphanumeric characters.
* `load_balancer_id` - (Required, String) Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters.
* `priority` - (Required, Int) Priority. A smaller value means a higher priority. Value range: 1 ~ 10000.
* `direction` - (Optional, String) Direction of the forwarding rule. Request: request direction from the client to the load balancer, Response: response direction from the backend server to the load balancer. Default value: Request.
* `rule_name` - (Optional, String) Name of the forwarding rule, 1 ~ 255 characters. Digits, upper and lower case letters, Chinese characters, periods (.), underscores (_) and hyphens (-) are supported.
* `tags` - (Optional, List) Tag list.

The `actions` object supports the following:

* `order` - (Required, Int) Execution order of the forwarding action. It must be unique and actions are executed in ascending order. Value range: 1 ~ 50000.
* `type` - (Required, String) Type of the forwarding action. Valid values: TargetGroup: forward to the target group; Redirect: redirect; FixedResponse: return fixed content; Rewrite: rewrite; InsertHeader: insert an HTTP header; RemoveHeader: remove an HTTP header. A forwarding action must contain one of TargetGroup, Redirect and FixedResponse, and it must be executed last.
* `fixed_response_config` - (Optional, Object) Configuration of the fixed response content.
* `insert_header_config` - (Optional, Object) Configuration of the inserted HTTP header.
* `redirect_config` - (Optional, Object) Redirect configuration. Except for HttpCode, the other configurations cannot all use the default values.
* `remove_header_config` - (Optional, Object) Configuration of the removed HTTP header.
* `rewrite_config` - (Optional, Object) Rewrite configuration.
* `target_group_config` - (Optional, Object) Configuration of the forwarding target group.

The `fixed_response_config` object of `actions` supports the following:

* `http_code` - (Required, Int) HTTP response code returned. 2xx, 4xx and 5xx are supported.
* `content_type` - (Optional, String) Format of the fixed content returned. Valid values: text/plain, text/css, text/html, application/javascript, application/json.
* `content` - (Optional, String) Fixed content returned. Only ASCII characters are supported, with a maximum of 1KB.

The `insert_header_config` object of `actions` supports the following:

* `key` - (Required, String) Key of the HTTP header to insert, 1 ~ 40 characters. Supported characters: a-z A-Z 0-9 - _ . Chinese characters are not supported, and Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host, x-forwarded-port are not supported.
* `value_type` - (Required, String) Type of the HTTP header value. Valid values: SystemDefined: system defined header; UserDefined: user defined header; ReferenceHeader: references a header of the request.
* `value` - (Required, String) Type of the HTTP header value. When ValueType is SystemDefined, valid values are ClientPort: client port, ClientIp: client IP address, Protocol: protocol of the client request, CLBPort: listening port of the load balancer instance. When ValueType is UserDefined, 1 ~ 128 printable characters, double quotes are not supported, the value cannot start or end with a space, and cannot end with \\. When ValueType is ReferenceHeader, references a header of the request, 1 ~ 128 printable characters, double quotes are not supported, the value cannot start or end with a space, and cannot end with \\.

The `redirect_config` object of `actions` supports the following:

* `http_code` - (Required, Int) HTTP code of the redirect. 301, 302, 303, 307 and 308 are supported.
* `host` - (Optional, String) Host address of the redirect, default value ${host}. 3 ~ 128 characters, supported characters: a-z 0-9 _ . -.
* `path` - (Optional, String) Path of the redirect, default value ${path}. 1 ~ 128 characters, supported characters: a-z A-Z 0-9 ? = _ . - / : .
* `port` - (Optional, String) Port of the redirect, default value ${port}. Value range: 1 ~ 65535.
* `protocol` - (Optional, String) Protocol of the redirect. Valid values: HTTP, HTTPS. Default value: ${protocol}.
* `query` - (Optional, String) Query string of the redirect, default value ${query}. 1 ~ 128 characters, printable characters are supported, #[]{}|<> and spaces are not supported.

The `remove_header_config` object of `actions` supports the following:

* `key` - (Required, String) Key of the HTTP header to remove, 1 ~ 40 characters. Supported characters: a-z A-Z 0-9 - _ . Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host, x-forwarded-port and server are not supported.

The `rewrite_config` object of `actions` supports the following:

* `host` - (Optional, String) Host address of the rewrite, default value ${host}. 3 ~ 128 characters, supported characters: a-z 0-9 _ . -.
* `path` - (Optional, String) Path of the rewrite, default value ${path}. 1 ~ 128 characters, supported characters: a-z A-Z 0-9 ? = _ . - / : .
* `query` - (Optional, String) Query string of the rewrite, default value ${query}. 1 ~ 128 characters, printable characters are supported, #[]{}|<> and spaces are not supported.

The `target_group_config` object of `actions` supports the following:

* `target_groups` - (Required, List) Target group list.
* `target_group_sticky_session` - (Optional, Object) Session persistence between target groups.

The `target_group_sticky_session` object of `target_group_config` supports the following:

* `enabled` - (Required, Bool) Whether to enable session persistence. Disabled by default.
* `timeout` - (Optional, Int) Timeout period in seconds. Value range: 1 ~ 86400. Default value: 1000.

The `target_groups` object of `target_group_config` supports the following:

* `target_group_id` - (Required, String) Target group ID, in the format of lbtg- followed by 8 alphanumeric characters.
* `weight` - (Optional, Int) Weight. Value range: [0, 100]. Default value: 10.

The `conditions` object supports the following:

* `type` - (Required, String) Type of the forwarding condition. Valid values: Host: host; Path: path; Header: HTTP header field; QueryString: HTTP query string; Method: request method; Cookie: cookie; SourceIp: source IP.
* `cookie_config` - (Optional, List) Cookie configuration.
* `header_config` - (Optional, Object) HTTP header configuration.
* `host_config` - (Optional, List) Host name. Host configuration can appear only once in a rule, 3 ~ 128 characters, exact match, regular expression match and wildcard match are supported. It cannot start or end with a period (.) or an underscore (_). Exact match supported characters: a-z 0-9 . - _ . Regular expression match starts with a tilde (~), supported characters: a-z 0-9 . - ? = ~ _ - + \ ^ * ! $ & | ( ) [ ] . Wildcard match uses an asterisk (*) for multiple characters and a question mark (?) for a single character, supported characters: a-z 0-9 . - _ * ?.
* `method_config` - (Optional, List) Request method. Valid values: HEAD, GET, POST, OPTIONS, PUT, PATCH, DELETE.
* `path_config` - (Optional, List) Forwarding path, 1 ~ 128 characters, exact match, regular expression match and wildcard match are supported. Exact match supported characters: a-z A-Z 0-9 . - _ / = :. Regular expression match must start with ~, where ~ means case-sensitive and ~* means case-insensitive, supported characters: a-z A-Z 0-9 . - _ / = ? ~ ^ * $ : ( ) [ ] + |. Wildcard match uses * for multiple characters and ? for a single character, supported characters: a-z A-Z 0-9 . - _ / = :.
* `query_string_config` - (Optional, List) Query string configuration.
* `source_ip_config` - (Optional, List) Source IP matching configuration. CIDR format, IP address x.x.x.x/32, IP range x.x.x.x/24.

The `cookie_config` object of `conditions` supports the following:

* `key` - (Required, String) Cookie key, 1 ~ 64 characters. Letters, digits and underscores are supported.
* `value` - (Required, String) Cookie value, 1 ~ 128 characters. Printable characters are supported.

The `header_config` object of `conditions` supports the following:

* `key` - (Required, String) Key of the HTTP header, 1 ~ 40 characters. Supported characters: a-z A-Z 0-9 - _ . Chinese characters, Host and Cookie are not supported.
* `values` - (Required, List) Values of the HTTP header, 1 ~ 128 characters. Printable characters are supported. Double quotes are not supported, the value cannot start or end with a space, and cannot end with \\.

The `query_string_config` object of `conditions` supports the following:

* `key` - (Required, String) Key of the query string, 1 ~ 16 characters. Printable characters are supported, spaces and #[]{}\|<> are not supported. * matches multiple characters and ? matches a single character.
* `value` - (Required, String) Value of the query string, 1 ~ 128 characters. Printable characters are supported, spaces and #[]{}\|<> are not supported. * matches multiple characters and ? matches a single character.

The `tags` object supports the following:

* `tag_key` - (Required, String) Key of the tag.
* `tag_value` - (Required, String) Value of the tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, in the form `<load_balancer_id>#<listener_id>#<rule_id>`.

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

ALB rules can be imported using the composite id load_balancer_id#listener_id#rule_id, e.g.

```
terraform import tencentcloud_alb_rules.example alb-12345678#lst-12345678#rule-12345678
```

