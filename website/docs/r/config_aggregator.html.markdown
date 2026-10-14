---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_aggregator"
sidebar_current: "docs-tencentcloud-resource-config_aggregator"
description: |-
  Provides a resource to manage a Tencent Cloud Config (配置审计) aggregator account group.
---

# tencentcloud_config_aggregator

Provides a resource to manage a Tencent Cloud Config (配置审计) aggregator account group.

## Example Usage

```hcl
resource "tencentcloud_config_aggregator" "example" {
  name        = "tf-example-aggregator"
  description = "aggregator created by terraform"
  type        = "CUSTOM"
  owner_uin   = 100000000000

  aggregator_accounts {
    member_uin  = 100000000001
    member_name = "member1"
  }

  aggregator_accounts {
    member_uin  = 100000000002
    member_name = "member2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, String) Account group description.
* `name` - (Required, String) Account group name.
* `owner_uin` - (Required, Int) Account group creator ID.
* `type` - (Required, String) Account group type: `RD` global account group, `CUSTOM` custom account group.
* `aggregator_accounts` - (Optional, List) Account group member information list, up to 100.

The `aggregator_accounts` object supports the following:

* `member_name` - (Required, String) Member name.
* `member_uin` - (Required, Int) Member ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `account_group_id` - Account group ID.
* `aggregator_status` - Creation status.


## Import

Config aggregator can be imported using the `<owner_uin>#<account_group_id>` (the composite id joined by `#`, e.g. `100000000000#aggr-123456`), e.g.

```
terraform import tencentcloud_config_aggregator.example 100000000000#aggr-123456
```

