Provides a resource to manage a Tencent Cloud Config (配置审计) aggregator account group.

Example Usage

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

Import

Config aggregator can be imported using the `<owner_uin>#<account_group_id>` (the composite id joined by `#`, e.g. `100000000000#aggr-123456`), e.g.

```
terraform import tencentcloud_config_aggregator.example 100000000000#aggr-123456
```
