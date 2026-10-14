# 提案: 修改 tencentcloud_dts_sync_job 资源的 tags 字段为 Set 类型

## Why

`tencentcloud_dts_sync_job` 资源中的 `tags` 字段当前定义为 `schema.TypeList`。用户在 Terraform 配置中以固定顺序声明 tags，但 DTS 云 API 的查询接口（DescribeSyncJobs）返回的 `Tags` 数组顺序是无序的（与创建时传入顺序不一致）。由于 List 类型对顺序敏感，Read 后写回 state 的 tags 顺序与配置文件顺序不一致，导致 `terraform plan` 永远显示 tags 存在 diff（无止境的非空 plan），用户无法收敛状态。将 `tags` 改为 `schema.TypeSet` 可以消除顺序敏感性，从根本上修复该问题。

## What Changes

- **BREAKING**：`tencentcloud_dts_sync_job` 资源 schema 中的 `tags` 字段类型从 `schema.TypeList` 修改为 `schema.TypeSet`（子字段 `tag_key`、`tag_value` 保持不变，仍为 `schema.TypeString`）
- 修改 Create 方法（`resourceTencentCloudDtsSyncJobCreate`）中读取 `tags` 参数的逻辑：`v.([]interface{})` 改为 `v.(*schema.Set).List()`
- 修改 Read 方法（`resourceTencentCloudDtsSyncJobRead`）中写回 `tags` 的逻辑（组装 `tagsList` 后 `d.Set` 的方式保持不变，Set 类型自动忽略顺序）
- 更新单元测试文件 `resource_tc_dts_sync_job_test.go` 中与 tags 相关的断言（`tags.#` 仍然有效，保持测试兼容）
- 更新资源文档 `resource_tc_dts_sync_job.md` 中的示例（HCL 语法层面 `tags` 块写法不变，文档由 `make doc` 重新生成）

说明：
- HCL 配置语法无需变更——用户配置文件中 `tags { tag_key = ... tag_value = ... }` 的写法在 List 与 Set 下完全相同
- 升级 Provider 后首次 `terraform plan` 会看到一次性的状态格式转换（List→Set），apply 后收敛，不会产生实际资源变更
- 该字段为 `ForceNew: true`，不涉及 Update 逻辑改动；`tag_key` + `tag_value` 组合在 Set 中天然去重，与腾讯云标签键唯一的约束一致

## Capabilities

### New Capabilities

- `dts-sync-job-tags-set`: tencentcloud_dts_sync_job 资源 tags 字段使用 Set 类型的行为规范——tags 作为无序集合处理，Read 写回 state 时不因云 API 返回顺序与配置顺序不一致而产生 diff

### Modified Capabilities

<!-- 本变更仅修改 tencentcloud_dts_sync_job 资源 tags 字段的类型实现，openspec/specs/ 下无该资源已有 spec，无需修改现有 spec -->

## Impact

- **代码文件**：
  - `tencentcloud/services/dts/resource_tc_dts_sync_job.go`（schema 定义 + Create 中 tags 读取逻辑）
  - `tencentcloud/services/dts/resource_tc_dts_sync_job_test.go`（测试断言核对）
- **文档文件**：
  - `tencentcloud/services/dts/resource_tc_dts_sync_job.md`（示例核对，Argument Reference 由工具自动生成）
- **不受影响**：
  - `tencentcloud/services/dts/data_source_tc_dts_sync_jobs.go`（数据源有独立的 tags 输出结构，不在本变更范围）
  - `tencentcloud/services/dts/service_tencentcloud_dts.go`（服务层 DescribeDtsSyncJob 仅透传 SDK 结构，无需修改）
  - `tencentcloud/provider.go` / `tencentcloud/provider.md`（资源已注册，无新增资源）
- **云 API 依赖**：无新增/变更 API 调用。仅涉及现有 `CreateSyncJob`（入参 `Tags []*TagItem`）与 `DescribeSyncJobs`（出参 `SyncJobInfo.Tags []*TagItem`）的参数组装方式，SDK 结构不变
- **兼容性影响**：属于 schema 类型变更（List→Set），state 中已存储的 tags 会发生一次性格式转换；HCL 配置语法完全兼容
