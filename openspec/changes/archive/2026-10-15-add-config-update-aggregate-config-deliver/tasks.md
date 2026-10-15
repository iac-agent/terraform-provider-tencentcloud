## 1. Service Layer

- [x] 1.1 Append `DescribeAggregateConfigDeliver(ctx context.Context, accountGroupId string) (*configv20220802.DescribeAggregateConfigDeliverResponseParams, error)` to `tencentcloud/services/config/service_tencentcloud_config.go`
- [x] 1.2 Build `NewDescribeAggregateConfigDeliverRequest()`, set `AccountGroupId`; wrap the call in `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `ratelimit.Check`; return `NonRetryableError` when `result == nil || result.Response == nil`

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver.go` with `ResourceTencentCloudConfigUpdateAggregateConfigDeliver()` schema
- [x] 2.2 Schema: `account_group_id` (String, Required, ForceNew), `status` (Int, Required), `deliver_name`/`target_arn`/`deliver_prefix`/`deliver_type` (String, Optional), `deliver_uin`/`deliver_content_type` (Int, Optional), `create_time` (String, Computed); add `Importer` with `schema.ImportStatePassthrough`
- [x] 2.3 Implement Create handler: `d.SetId(account_group_id)` then delegate to the Update handler
- [x] 2.4 Implement Read handler: call `service.DescribeAggregateConfigDeliver(ctx, d.Id())`; if response nil, log id and `d.SetId("")`; otherwise set each non-nil field
- [x] 2.5 Implement Update handler: build `NewUpdateAggregateConfigDeliverRequest()`, set `AccountGroupId` from `d.Id()` and all schema fields; call `UpdateAggregateConfigDeliverWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError`; then call Read
- [x] 2.6 Implement Delete handler as a no-op returning nil

## 3. Provider Registration

- [x] 3.1 Register `"tencentcloud_config_update_aggregate_config_deliver"` in `tencentcloud/provider.go` ResourcesMap (adjacent to other config resources)
- [x] 3.2 Add `tencentcloud_config_update_aggregate_config_deliver` entry to `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver.md` with one-line description, example usage (including Import section) per gendoc conventions; do not add Argument/Attribute Reference sections

## 5. Tests

- [x] 5.1 Create `tencentcloud/services/config/resource_tc_config_update_aggregate_config_deliver_test.go` using the Terraform acceptance test suite (matching `resource_tc_config_deliver_config_test.go` style) covering basic apply and update of `status`

## 6. Verification

- [x] 6.1 Confirm all new schema fields map to existing cloud API params (`DescribeAggregateConfigDeliver` read, `UpdateAggregateConfigDeliver` update) and that no non-existent params are referenced
- [x] 6.2 Confirm `snake_case` naming, `tccommon` retry patterns and error handling follow provider conventions
