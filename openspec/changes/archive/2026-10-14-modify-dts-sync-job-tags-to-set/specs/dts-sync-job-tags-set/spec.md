# Spec Delta: dts-sync-job-tags-set

## ADDED Requirements

### Requirement: tencentcloud_dts_sync_job 的 tags 字段 SHALL 为无序集合（Set）类型

`tencentcloud_dts_sync_job` 资源 schema 中的 `tags` 字段 SHALL 定义为 `schema.TypeSet`，元素为包含 `tag_key`（string）与 `tag_value`（string）子字段的资源块。tags 字段 SHALL 保持 `Optional + Computed + ForceNew` 属性不变。

#### Scenario: schema 类型定义为 Set

- **WHEN** 检查 `resourceTencentCloudDtsSyncJob` 的 schema 定义
- **THEN** `tags` 字段的 Type 为 `schema.TypeSet`，Elem 为包含 `tag_key`、`tag_value` 两个 string 子字段的 `*schema.Resource`，且 Optional/Computed/ForceNew 属性与变更前一致

#### Scenario: HCL 配置语法保持兼容

- **WHEN** 用户使用既有配置写法 `tags { tag_key = "createBy" tag_value = "Terraform" }` 创建 tencentcloud_dts_sync_job 资源
- **THEN** 配置可以被正常解析，无需修改任何 HCL 语法

### Requirement: 创建同步任务时 SHALL 以 Set 语义读取 tags 参数

`tencentcloud_dts_sync_job` 资源的 Create 方法 SHALL 通过 `v.(*schema.Set).List()` 读取 tags 参数，遍历每个元素组装 `dts.TagItem`（TagKey/TagValue）并传入 `CreateSyncJob` 请求的 `Tags` 字段。每个 tag_key + tag_value 组合 SHALL 在请求中仅出现一次（Set 去重语义）。

#### Scenario: 创建时正确组装 tags 请求

- **WHEN** 用户配置了多个 tags（如 `{tag_key: "aaa", tag_value: "bbb"}` 与 `{tag_key: "ccc", tag_value: "ddd"}`）并创建 tencentcloud_dts_sync_job
- **THEN** CreateSyncJob 请求的 Tags 数组包含全部配置的 TagItem 对象，每个对象的 TagKey/TagValue 与配置一致

#### Scenario: 重复 tag 配置被去重

- **WHEN** 用户在配置中声明了两个完全相同的 tag_key + tag_value 组合
- **THEN** CreateSyncJob 请求的 Tags 数组中该组合仅出现一次

### Requirement: 读取同步任务时 tags 写回 state SHALL 不受云 API 返回顺序影响

`tencentcloud_dts_sync_job` 资源的 Read 方法 SHALL 将 DescribeSyncJobs 返回的 `SyncJobInfo.Tags`（`[]*TagItem`）逐项转换为 `{tag_key, tag_value}` map 列表并通过 `d.Set("tags", ...)` 写回 state。由于 tags 为 Set 类型，terraform 比较 state 与配置时 SHALL 忽略元素顺序，从而当云 API 返回的 tags 顺序与用户配置顺序不一致时 SHALL NOT 产生 diff。

#### Scenario: 云 API 返回顺序与配置顺序不一致时不产生 diff

- **WHEN** tencentcloud_dts_sync_job 创建成功后执行 Read，且 DescribeSyncJobs 返回的 Tags 顺序与用户配置中声明的 tags 顺序不同
- **THEN** tags 以 Set 语义写回 state，`terraform plan` 不显示 tags 字段的任何变更

#### Scenario: 云端返回空 tags

- **WHEN** DescribeSyncJobs 返回的 SyncJobInfo.Tags 为 nil 或空数组
- **THEN** Read 方法不向 state 写入 tags 数据（沿用现有 nil 判断逻辑），不产生错误

### Requirement: tags 字段 SHALL 保持 ForceNew 语义

`tencentcloud_dts_sync_job` 资源的 `tags` 字段 SHALL 保持 `ForceNew: true`。变更 tags 配置时 SHALL 触发资源重建（先隔离销毁旧任务再创建新任务），而非调用更新接口。

#### Scenario: 修改 tags 触发重建

- **WHEN** 用户修改现有 tencentcloud_dts_sync_job 资源的 tags 配置并执行 apply
- **THEN** terraform 计划销毁并重建该资源，不调用任何针对 tags 的更新 API
