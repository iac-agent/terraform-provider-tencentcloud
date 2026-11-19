Provides a resource to create an aggregate (account group) config rule of Config.

Example Usage

System preset aggregate rule with scheduled trigger

```hcl
resource "tencentcloud_config_aggregate_rule" "example" {
  identifier       = "cam-user-group-bound"
  identifier_type  = "SYSTEM"
  rule_name        = "tf-example-aggregate-rule"
  resource_type    = ["QCS::CAM::User"]
  account_group_id = "agr-xxxxxx"
  risk_level       = 3
  description      = "tf example aggregate config rule"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }
}
```

Aggregate rule with input parameters

```hcl
resource "tencentcloud_config_aggregate_rule" "example_with_params" {
  identifier       = "cam-user-mfa-check"
  identifier_type  = "SYSTEM"
  rule_name        = "tf-example-aggregate-rule-params"
  resource_type    = ["QCS::CAM::User"]
  account_group_id = "agr-xxxxxx"
  risk_level       = 2

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  input_parameter {
    parameter_key = "maxMemorySize"
    type          = "Require"
    value         = "512"
  }
}
```

Import

Aggregate config rule can be imported using the accountGroupId#ruleId, e.g.

```
terraform import tencentcloud_config_aggregate_rule.example agr-xxxxxx#cr-3xhsd76j603v0a8ma0i73
```
