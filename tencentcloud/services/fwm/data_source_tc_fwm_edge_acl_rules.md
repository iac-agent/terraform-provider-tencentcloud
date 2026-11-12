Use this data source to query edge ACL rules of Cloud Firewall (CFW).

Example Usage

```hcl
data "tencentcloud_fwm_edge_acl_rules" "example" {
  group_id  = "fwmrg-9z03vakjj8"
  direction = 1
}

output "rule_ids" {
  value = [for rule in data.tencentcloud_fwm_edge_acl_rules.example.data : rule.rule_id]
}
```

Query with filters and pagination:

```hcl
data "tencentcloud_fwm_edge_acl_rules" "example" {
  group_id = "fwmrg-9z03vakjj8"
  filters {
    name   = "RuleAction"
    values = ["accept"]
  }
  limit  = 50
  offset = 0
  order  = "asc"
  by     = "Sequence"
}
```

Argument Reference

The following arguments are supported:

- `group_id` - (Required) Edge ACL rule group ID, the rule group to query rules from.
- `direction` - (Optional) Rule direction: 0-egress, 1-ingress.
- `filters` - (Optional) Filter conditions, support filtering by RuleId, Protocol, RuleAction etc. Each filter contains `name`, `values` and optional `operator_type` (1 means equal).
- `limit` - (Optional) Page size, default 100, max 1000.
- `offset` - (Optional) Pagination offset, default 0.
- `order` - (Optional) Sort order: asc or desc.
- `by` - (Optional) Sort field, support Sequence, RuleId etc.
- `result_output_file` - (Optional) Used to save results.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `total_count` - Total number of rules matched.
- `data` - Edge ACL rule list. Each element contains:
  - `rule_id` - Rule ID.
  - `sequence` - Rule execution order.
  - `direction` - Rule direction: 0-egress, 1-ingress.
  - `source_content` - Source address content.
  - `source_type` - Source address type.
  - `source_name` - Source address name.
  - `target_content` - Target address content.
  - `target_type` - Target address type.
  - `target_name` - Target address name.
  - `port` - Target port.
  - `protocol` - Protocol type.
  - `rule_action` - Rule action: accept, drop, log.
  - `description` - Rule description.
  - `scope` - Rule effective scope.
  - `hit_count` - Rule hit count.
  - `last_hit_time` - Rule last hit time.
  - `create_time` - Rule creation time.
  - `update_time` - Rule latest update time.
