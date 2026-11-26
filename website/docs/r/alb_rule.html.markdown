---
subcategory: "Application Load Balancer(ALB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_alb_rule"
sidebar_current: "docs-tencentcloud-resource-alb_rule"
description: |-
  Provides a resource to create and manage an ALB forwarding rule.
---

# tencentcloud_alb_rule

Provides a resource to create and manage an ALB forwarding rule.

Destroy only deletes the forwarding rule identified by `rule_ids` in `DeleteRules`. The documentation does not state whether it also affects the owning listener, load balancer instance or target groups, so those objects are never deleted by this resource. `dry_run = true` only pre-checks a request and must not be used for a managed rule; a pre-check never creates or changes the cloud object. An Update is rejected when a previously configured optional attribute is removed, because the API does not document whether an omitted value keeps or clears the cloud setting.

## Example Usage

```hcl
resource "tencentcloud_alb_rule" "example" {
  load_balancer_id = "alb-12345678"
  listener_id      = "lst-12345678"
  client_token     = "tf-example-rule-token"

  rules = [
    {
      priority  = 10
      rule_name = "tf-example"
      direction = "Request"

      actions = [
        {
          order = 1
          type  = "TargetGroup"

          target_group_config = {
            target_groups = [
              {
                target_group_id = "lbtg-12345678"
                weight          = 30
              },
            ]

            target_group_sticky_session = {
              enabled = true
              timeout = 1000
            }
          }
        },
      ]

      conditions = [
        {
          type        = "Path"
          path_config = ["/api"]

          host_config         = ["example.com"]
          method_config       = ["GET"]
          source_ip_config    = ["10.0.0.1/32"]
          query_string_config = jsonencode([])
        },
      ]

      tags = [
        {
          tag_key   = "env"
          tag_value = "test"
        },
      ]
    },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, String) Listener ID in the format `lst-` followed by 8 alphanumeric characters. Immutable: changing it forces resource replacement.
* `load_balancer_id` - (Required, String) Load balancer instance ID in the format `alb-` followed by 8 alphanumeric characters. Immutable: changing it forces resource replacement.
* `rules` - (Required, List) Forwarding rule list. This resource manages exactly one rule, so the list must contain exactly one element.
* `client_token` - (Optional, String) Client token used to make the CreateRules request idempotent. ASCII characters only. If omitted, the cloud service uses the request id. Create-only: changing it forces resource replacement.
* `dry_run` - (Optional, Bool) Whether to only pre-check the request. When true the cloud API does not apply the change.

The `rules` object supports the following:

* `actions` - (Required, List) Forwarding action list.
* `conditions` - (Required, List) Forwarding condition list.
* `priority` - (Required, Int) Rule priority. Lower values have higher priority. Valid range: 1 to 10000.
* `direction` - (Optional, String) Rule direction. Valid values: `Request` (client to load balancer), `Response` (backend to load balancer). Defaults to `Request`. Create-only: changing it forces resource replacement.
* `rule_name` - (Optional, String) Forwarding rule name, 1 to 255 characters.
* `tags` - (Optional, List) Tags of the rule. Create-only: the Modify API does not accept tags, so changing them forces resource replacement.

The `actions` object of `rules` supports the following:

* `order` - (Required, Int) Action execution order. Values must not repeat and are executed ascending. Valid range: 1 to 50000.
* `type` - (Required, String) Action type. Valid values: `TargetGroup`, `Redirect`, `FixedResponse`, `Rewrite`, `InsertHeader`, `RemoveHeader`. One of `TargetGroup`, `Redirect` or `FixedResponse` is required and must be placed last.
* `fixed_response_config` - (Optional, Object) Fixed response configuration.
* `insert_header_config` - (Optional, Object) Insert HTTP header configuration.
* `redirect_config` - (Optional, Object) Redirect configuration. Except for `http_code`, not all other settings may use their default values.
* `remove_header_config` - (Optional, Object) Remove HTTP header configuration.
* `rewrite_config` - (Optional, Object) Rewrite configuration.
* `target_group_config` - (Optional, Object) Forwarding target group configuration.

The `fixed_response_config` object of `actions` supports the following:

* `http_code` - (Required, Int) Returned HTTP response code. Supports 2xx, 4xx and 5xx.
* `content_type` - (Optional, String) Format of the returned content. Valid values: `text/plain`, `text/css`, `text/html`, `application/javascript`, `application/json`.
* `content` - (Optional, String) Returned fixed content. ASCII characters only, at most 1KB.

The `insert_header_config` object of `actions` supports the following:

* `key` - (Required, String) HTTP header key, 1 to 40 characters. Supports a-z, A-Z, 0-9, `-` and `_`.
* `value_type` - (Required, String) Type of the HTTP header value. Valid values: `SystemDefined`, `UserDefined`, `ReferenceHeader`.
* `value` - (Required, String) HTTP header value. Its allowed values depend on `value_type`.

The `redirect_config` object of `actions` supports the following:

* `http_code` - (Required, Int) Redirect HTTP code. Valid values: 301, 302, 303, 307, 308.
* `host` - (Optional, String) Redirect host. Defaults to `${host}`.
* `path` - (Optional, String) Redirect path. Defaults to `${path}`.
* `port` - (Optional, String) Redirect port. Defaults to `${port}`. Valid range: 1 to 65535.
* `protocol` - (Optional, String) Redirect protocol. Valid values: `HTTP`, `HTTPS`. Defaults to `${protocol}`.
* `query` - (Optional, String) Redirect query string. Defaults to `${query}`.

The `remove_header_config` object of `actions` supports the following:

* `key` - (Required, String) HTTP header key to remove, 1 to 40 characters.

The `rewrite_config` object of `actions` supports the following:

* `host` - (Optional, String) Rewrite host. Defaults to `${host}`.
* `path` - (Optional, String) Rewrite path. Defaults to `${path}`.
* `query` - (Optional, String) Rewrite query string. Defaults to `${query}`.

The `target_group_config` object of `actions` supports the following:

* `target_groups` - (Required, List) Target group list.
* `target_group_sticky_session` - (Optional, Object) Session affinity across target groups.

The `target_group_sticky_session` object of `target_group_config` supports the following:

* `enabled` - (Required, Bool) Whether session affinity is enabled. Disabled by default.
* `timeout` - (Optional, Int) Session timeout in seconds. Valid range: 1 to 86400. Defaults to 1000.

The `target_groups` object of `target_group_config` supports the following:

* `target_group_id` - (Required, String) Target group ID in the format `lbtg-` followed by 8 alphanumeric characters.
* `weight` - (Optional, Int) Weight. Valid range: 0 to 100. Defaults to 10.

The `conditions` object of `rules` supports the following:

* `type` - (Required, String) Condition type. Valid values: `Host`, `Path`, `Header`, `QueryString`, `Method`, `Cookie`, `SourceIp`.
* `cookie_config` - (Optional, List) Cookie configuration.
* `header_config` - (Optional, Object) HTTP header condition configuration.
* `host_config` - (Optional, Set) Host names. Supports exact, regular expression and wildcard matching.
* `method_config` - (Optional, Set) Request methods. Valid values: `HEAD`, `GET`, `POST`, `OPTIONS`, `PUT`, `PATCH`, `DELETE`.
* `path_config` - (Optional, Set) Forwarding paths. Supports exact, regular expression and wildcard matching.
* `query_string_config` - (Optional, List) Query string configuration.
* `source_ip_config` - (Optional, Set) Source IP CIDRs, for example `x.x.x.x/32` or `x.x.x.x/24`.

The `cookie_config` object of `conditions` supports the following:

* `key` - (Required, String) Cookie key, 1 to 64 characters.
* `value` - (Required, String) Cookie value, 1 to 128 characters.

The `header_config` object of `conditions` supports the following:

* `key` - (Required, String) HTTP header key, 1 to 40 characters.
* `values` - (Required, Set) HTTP header values, 1 to 128 characters each.

The `query_string_config` object of `conditions` supports the following:

* `key` - (Required, String) Query string key, 1 to 16 characters.
* `value` - (Required, String) Query string value, 1 to 128 characters.

The `tags` object of `rules` supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource identifier in the form `<load_balancer_id>#<listener_id>#<rule_id>`.
* `next_token` - Pagination token returned by the last DescribeRules call. Empty on the last page.


## Import

ALB rules can be imported using the composite id `<load_balancer_id>#<listener_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_alb_rule.example alb-12345678#lst-12345678#rule-12345678
```

