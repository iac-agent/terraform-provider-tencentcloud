Provides a resource to create a Config aggregate (account group) config rule.

Example Usage

```hcl
resource "tencentcloud_config_aggregate_config_rule" "example" {
  account_group_id = "ca-Zp9EiU7Pa4Vy0lFFm8XM"
  identifier       = "cam-password-policy-check"
  identifier_type  = "SYSTEM"
  rule_name        = "tf-example-rule"
  resource_type    = ["QCS::CAM::User"]
  risk_level       = 1
  description      = "remark."

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  input_parameter {
    parameter_key = "MinimumLength"
    type          = "Optional"
    value         = "8"
  }

  region_scope = ["ap-guangzhou", "ap-shanghai"]

  tags_scope {
    tag_key   = "env"
    tag_value = "prod"
  }

  exclude_resource_ids_scope = ["ins-xxxxxx"]
}
```

Import

Config aggregate config rule can be imported using the composite id, e.g. `account_group_id#rule_id`.

```
terraform import tencentcloud_config_aggregate_config_rule.example ca-Zp9EiU7Pa4Vy0lFFm8XM#cr-VNsYhja10ewnInqtc0qp
```

-> `account_group_id`, `identifier` and `identifier_type` are immutable: changing any of them forces resource replacement. The Add API does not return the rule id, so the provider correlates the created rule by its (unique) rule name through ListAggregateConfigRules. The `compliance_pack_id`, `compliance_pack_name`, `source_condition`, `manage_input_parameter` and `manage_trigger_type` attributes are read-only details reported by the Describe API for preset rules.
