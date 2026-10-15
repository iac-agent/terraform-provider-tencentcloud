## Why

Tencent Cloud Config supports account groups (账号组) that aggregate resources across multiple
member accounts. Administrators need to configure a delivery (投递设置) for each account group so
that Config audit logs are shipped to COS or CLS for centralized compliance auditing. Today this
account-group level delivery setting has no Terraform support, forcing users to configure it
manually in the console and making it impossible to manage as code.

## What Changes

- Add a new Terraform resource `tencentcloud_config_update_aggregate_config_deliver`
  (RESOURCE_KIND_CONFIG, config-only lifecycle) for the account-group delivery setting.
- Resource Read uses `DescribeAggregateConfigDeliver` (input `account_group_id`).
- Resource Create/Update use `UpdateAggregateConfigDeliver` — there is no dedicated Create API,
  so Create behaves as the first Update (idempotent set).
- Delete is a no-op: the account-group delivery setting always exists once it is configured;
  there is no delete API, so the resource only manages reading and updating.
- Add a new service-layer method `DescribeAggregateConfigDeliver(ctx, accountGroupId)` to
  `service_tencentcloud_config.go`.
- Register the resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add example documentation `resource_tc_config_update_aggregate_config_deliver.md`.

## Capabilities

### New Capabilities
- `config-update-aggregate-config-deliver`: Manage the Tencent Cloud Config account-group delivery
  setting (投递设置) as a Terraform RESOURCE_KIND_CONFIG resource, supporting read and update of the
  delivery name, target ARN, delivery status, prefix, type, cross-account uin and content type.

### Modified Capabilities
<!-- No existing spec requirements change. -->

## Impact

- **New files**
  - `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver.go`
  - `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver.md`
  - `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver_test.go`
- **Modified files**
  - `tencentcloud/services/config/service_tencentcloud_config.go` (append service method)
  - `tencentcloud/provider.go` (register resource)
  - `tencentcloud/provider.md` (register resource doc entry)
- **Cloud APIs**
  - `config` v20220802 `DescribeAggregateConfigDeliver`
  - `config` v20220802 `UpdateAggregateConfigDeliver`
- No breaking changes; purely additive.
