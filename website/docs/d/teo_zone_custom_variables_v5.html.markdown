---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_custom_variables_v5"
sidebar_current: "docs-tencentcloud-data_source-teo_zone_custom_variables_v5"
description: |-
  Use this data source to query the zone-level custom variables of TencentCloud EdgeOne (TEO), including variable definitions and variable operation rules. When no custom variables are configured, both the variable list and the operation rule list are empty.
---

# tencentcloud_teo_zone_custom_variables_v5

Use this data source to query the zone-level custom variables of TencentCloud EdgeOne (TEO), including variable definitions and variable operation rules. When no custom variables are configured, both the variable list and the operation rule list are empty.

## Example Usage

### Query zone-level custom variables by zone ID

```hcl
data "tencentcloud_teo_zone_custom_variables_v5" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

### Query zone-level custom variables referenced by a zone resource

```hcl
data "tencentcloud_teo_zone_custom_variables_v5" "example" {
  zone_id = tencentcloud_teo_zone.example.id
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

                * `name` - Custom variable name. The custom variable must be defined before it can be used in operations.
                * `value` - Custom variable value. Constant strings, variables and formulas are supported, while Chinese is not supported. The length cannot exceed 1000 characters.
              * `name` - Operation name. The name needs to correspond to the parameter struct, for example, when Name=Set, SetParameters is required. Currently only `Set` is supported.
              * `set_parameters` - Custom variable setting parameters. When multiple operations exist in this parameter, they are executed in the order of the array. Required when Name is `Set`. Note: this field may return null, which indicates a failure to obtain a valid value.
            * `actions` - Operations. Note: Actions and SubRules cannot both be empty. Note: this field may return null, which indicates a failure to obtain a valid value.
            * `condition` - Match condition.
            * `name` - Custom variable name. The custom variable must be defined before it can be used in operations.
            * `value` - Custom variable value. Constant strings, variables and formulas are supported, while Chinese is not supported. The length cannot exceed 1000 characters.
          * `branches` - Sub-rule branches. The recursive branch structure of the cloud API is truncated at this level because only one level of SubRules is currently supported. Note: this field may return null, which indicates a failure to obtain a valid value.
          * `description` - Rule remarks. Note: this field may return null, which indicates a failure to obtain a valid value.
          * `name` - Operation name. The name needs to correspond to the parameter struct, for example, when Name=Set, SetParameters is required. Currently only `Set` is supported.
          * `set_parameters` - Custom variable setting parameters. When multiple operations exist in this parameter, they are executed in the order of the array. Required when Name is `Set`. Note: this field may return null, which indicates a failure to obtain a valid value.
        * `actions` - Operations. Note: Actions and SubRules cannot both be empty. Note: this field may return null, which indicates a failure to obtain a valid value.
        * `condition` - Match condition.
        * `name` - Custom variable name. The custom variable must be defined before it can be used in operations.
        * `sub_rules` - Sub-rule list. When multiple rules exist in this list, they are executed from top to bottom. Note: SubRules and Actions cannot both be empty, and only one level of SubRules is currently supported. Note: this field may return null, which indicates a failure to obtain a valid value.
        * `value` - Custom variable value. Constant strings, variables and formulas are supported, while Chinese is not supported. The length cannot exceed 1000 characters.
      * `branches` - Sub-rule branches. Note: this field may return null, which indicates a failure to obtain a valid value.
      * `description` - Rule remarks. Note: this field may return null, which indicates a failure to obtain a valid value.
      * `name` - Operation name. The name needs to correspond to the parameter struct, for example, when Name=Set, SetParameters is required. Currently only `Set` is supported.
      * `set_parameters` - Custom variable setting parameters. When multiple operations exist in this parameter, they are executed in the order of the array. Required when Name is `Set`. Note: this field may return null, which indicates a failure to obtain a valid value.
    * `actions` - Operations. Note: Actions and SubRules cannot both be empty. Note: this field may return null, which indicates a failure to obtain a valid value.
    * `condition` - Match condition.
    * `sub_rules` - Sub-rule list. When multiple rules exist in this list, they are executed from top to bottom. Note: SubRules and Actions cannot both be empty, and only one level of SubRules is currently supported. Note: this field may return null, which indicates a failure to obtain a valid value.
  * `branches` - Sub-rule branches. Currently only one rule is supported in this list; additional rules are invalid. Note: this field may return null, which indicates a failure to obtain a valid value.
  * `description` - Rule remarks. Multiple remarks can be filled in. Note: this field may return null, which indicates a failure to obtain a valid value.
  * `description` - Variable description. The length cannot exceed 60 characters.
  * `initial_value` - Initial value of the variable. Constant strings, variables and formulas are supported. The length cannot exceed 255 characters.
  * `name` - Variable name. The full prefix is required: user.zone.* indicates a zone-level custom variable, and user.rule.* indicates a rule-level custom variable. The custom part after the prefix only supports upper and lower case letters, numbers and underscores. Variable names are case-sensitive and can be up to 50 characters long. The name cannot be modified after creation.
* `custom_variable_operations` - Zone-level custom variable operation rules. Empty when no operation rules are configured.
* `custom_variables` - List of zone-level custom variables. Empty when no custom variables are configured.


