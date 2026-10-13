---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_list_aggregators"
sidebar_current: "docs-tencentcloud-datasource-config_list_aggregators"
description: |-
  Use this data source to query the list of TencentCloud Config (配置审计) account groups (aggregators).
---

# tencentcloud_config_list_aggregators

Use this data source to query the list of TencentCloud Config (配置审计) account groups (aggregators).

## Example Usage

```hcl
data "tencentcloud_config_list_aggregators" "example" {}

output "total" {
  value = data.tencentcloud_config_list_aggregators.example.total
}

output "items" {
  value = data.tencentcloud_config_list_aggregators.example.items
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Aggregator (account group) list.
  * `account_count` - Count of members in the account group.
  * `account_group_id` - Account group ID.
  * `aggregator_status` - Account group status.
  * `create_time` - Creation time.
  * `description` - Account group description.
  * `member_name` - Account group name.
  * `name` - Account group name.
  * `owner_uin` - Owner user ID.
  * `type` - Account group type, RD: global account group.
* `total` - Total count of aggregators.


