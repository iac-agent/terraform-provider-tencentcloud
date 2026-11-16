Use this data source to query edge ACL rules of Cloud Firewall (CFW).

Example Usage

```hcl
data "tencentcloud_fwm_edge_acl_rules" "example" {
  group_id  = "fwmrg-9z03vakjj8"
  direction = 1
}

output "rule_ids" {
  value = [for r in data.tencentcloud_fwm_edge_acl_rules.example.rules : r.rule_id]
}
```

Query with filters

```hcl
data "tencentcloud_fwm_edge_acl_rules" "filtered" {
  group_id  = "fwmrg-9z03vakjj8"
  direction = 1
  filters {
    name   = "RuleId"
    values = ["2fe0b492-5623-4646-b40e-7b0908307ca8"]
  }
}
```

Argument Reference

The following arguments are supported:

- `group_id` - (Optional) Edge ACL rule group ID.
- `direction` - (Optional) Rule direction: 0-egress, 1-ingress.
- `filters` - (Optional) Filter list. Each filter contains `name`, `values` and `operator_type` (1-equal, 7-in, 9-fuzzy match). Supported filter names: RuleId, Direction, Protocol, RuleAction, etc.
- `limit` - (Optional) Page size, default 100, max 1000.
- `offset` - (Optional) Page offset, default 0.
- `order` - (Optional) Sort order: asc-ascending, desc-descending.
- `by` - (Optional) Sort field, supports Sequence, RuleId, etc.
- `result_output_file` - (Optional) Used to save results.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `total_count` - Total number of rules matching the filter.
- `rules` - A list of edge ACL rules. Each element contains the following attributes:
  - `rule_id` - Rule ID.
  - `sequence` - Rule execution sequence.
  - `direction` - Rule direction: 0-egress, 1-ingress.
  - `source_content` - Source address content.
  - `source_type` - Source address type.
  - `source_name` - Source address name (returned when the type is template/instance/tag).
  - `target_content` - Target address content.
  - `target_type` - Target address type.
  - `target_name` - Target address name (returned when the type is template/instance/tag).
  - `port` - Port.
  - `protocol` - Protocol type.
  - `rule_action` - Rule action: accept-pass, drop-block, log-observe.
  - `description` - Rule description.
  - `scope` - Effective scope of the rule.
  - `param_template_id` - Parameter template ID.
  - `invalid` - Whether the rule is invalid: 0-valid, 1-invalid.
  - `create_time` - Rule creation time.
  - `update_time` - Rule latest update time.
