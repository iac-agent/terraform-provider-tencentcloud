# Tasks: add-config-aggregator-resource

## 1. 资源实现

- [x] 1.1 新建 `tencentcloud/services/config/resource_tc_config_aggregator.go`：
  - `ResourceTencentCloudConfigAggregator()`，Schema 按 design.md（name/description/type/owner_uin/aggregator_accounts 必填可选如 design；account_group_id/aggregator_status 为 computed）。
  - Create：构造 `CreateAggregatorRequest`（Name/Description/Type/AggregatorAccounts），retry(WriteRetryTimeout) 调 `UseConfigV20220802Client().CreateAggregatorWithContext`；校验 Response/AccountGroupId 非空否则 NonRetryableError；`d.SetId(ownerUin + FILED_SP + accountGroupId)` 后调用 Read。
  - Read：拆分联合 id，retry(ReadRetryTimeout) 调 `DescribeAggregatorWithContext`，判空回填各字段，资源不存在按规范打日志后 SetId("")。
  - Update：HasChange 检查后构造 `UpdateAggregatorRequest`（Name/Description/AggregatorAccounts + id 中的 AccountGroupId/OwnerUin），retry 调 `UpdateAggregatorWithContext`。
  - Delete：构造 `DeleteAggregatorsRequest`（id 中 AccountGroupId/OwnerUin），retry 调 `DeleteAggregatorsWithContext`。
  - Importer: `schema.ImportStatePassthrough`。

## 2. 单元测试

- [x] 2.1 新建 `tencentcloud/services/config/resource_tc_config_aggregator_test.go`：使用 gomonkey mock 云 API（Create/Describe/Update/Delete），对 CRUD 业务逻辑做单元测试，不使用 terraform 测试套件，不实际执行。

## 3. 注册与文档

- [x] 3.1 `tencentcloud/provider.go` 中注册 `"tencentcloud_config_aggregator": config.ResourceTencentCloudConfigAggregator()`（参照 igtm_strategy 注册位置），同步更新 `tencentcloud/provider.md`。
- [x] 3.2 新建 `tencentcloud/services/config/resource_tc_config_aggregator.md`：一句话描述（"Provides a resource to ..."，含 config 产品名）、Example Usage（jsonencode 如适用）、Import 部分说明联合 id `<owner_uin><分隔符><account_group_id>`；不写 Argument/Attribute Reference。

## 4. 正确性检查

- [x] 4.1 核对所有入参存在于对应云 API request struct（CreateAggregator 无 OwnerUin，Describe/Update/Delete 必传 OwnerUin），出参字段与 response struct 一致；确保代码可编译，不执行 go build/vet/test。
