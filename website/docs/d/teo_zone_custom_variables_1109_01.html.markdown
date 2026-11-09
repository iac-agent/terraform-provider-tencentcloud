---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_custom_variables_1109_01"
sidebar_current: "docs-tencentcloud-data_source-teo_zone_custom_variables_1109_01"
description: |-
  Use this data source to query zone-level custom variables of a TEO zone, including variable definitions and variable operation rules.
---

# tencentcloud_teo_zone_custom_variables_1109_01

Use this data source to query zone-level custom variables of a TEO zone, including variable definitions and variable operation rules.

## Example Usage

### Query zone-level custom variables

```hcl
data "tencentcloud_teo_zone_custom_variables_1109_01" "example" {
  zone_id = "zone-29371sfdvge2"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

            * `name` - Custom variable name. The custom variable must be defined before it can be used in operations.
            * `value` - Custom variable assignment. Supports constant strings, variables and formulas; Chinese characters are not supported. Up to 1000 characters.
          * `description` - Rule comments.
          * `name` - Action name, which needs to correspond to the parameter struct, e.g. when Name=Set, SetParameters is required. Currently only `Set` is supported.
          * `set_parameters` - Custom variable setting parameters. When there are multiple operations, they are executed in array order. Required when Name is `Set`.
        * `actions` - Actions of the branch. Actions and SubRules cannot both be empty. Valid value of name: `Set`.
        * `condition` - Match condition expression.
        * `name` - Custom variable name. The custom variable must be defined before it can be used in operations.
        * `sub_rules` - Deeper sub-rule list. Rules in this list are executed from top to bottom. Only one layer of SubRules is supported; deeper nesting only exposes rule comments.
        * `value` - Custom variable assignment. Supports constant strings, variables and formulas; Chinese characters are not supported. Up to 1000 characters.
      * `branches` - Sub-rule branches.
      * `description` - Rule comments.
      * `name` - Action name, which needs to correspond to the parameter struct, e.g. when Name=Set, SetParameters is required. Currently only `Set` is supported.
      * `set_parameters` - Custom variable setting parameters. When there are multiple operations, they are executed in array order. Required when Name is `Set`.
    * `actions` - Actions of the branch. Actions and SubRules cannot both be empty. Valid value of name: `Set`.
    * `condition` - Match condition expression.
    * `sub_rules` - Sub-rule list. Rules in this list are executed from top to bottom. SubRules and Actions cannot both be empty, and only one layer of SubRules is supported.
  * `branches` - Sub-rule branches. This list currently only supports one rule entry; multiple entries are invalid.
  * `description` - Rule comments. Multiple comments can be filled in.
  * `description` - Variable description, up to 60 characters.
  * `initial_value` - Initial value of the variable. Supports constant strings, variables and formulas, up to 255 characters.
  * `name` - Variable name. Use the full prefix: user.zone.* for zone-level custom variables, user.rule.* for rule-level custom variables. The custom part after the prefix only supports letters, digits and underscores. Variable names are case-sensitive and up to 50 characters long. The name cannot be modified after creation.
* `custom_variable_operations` - Zone-level custom variable operation rules.
* `custom_variables` - List of zone-level custom variable definitions.


