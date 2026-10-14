# 技术设计: 修改 tencentcloud_dts_sync_job 资源的 tags 字段为 Set 类型

## Context

`tencentcloud_dts_sync_job` 资源（位于 `tencentcloud/services/dts/resource_tc_dts_sync_job.go`）的 `tags` 字段当前定义为：

```go
"tags": {
    Type:        schema.TypeList,
    Optional:    true,
    Computed:    true,
    ForceNew:    true,
    Description: "Tag information.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "tag_key":   {Type: schema.TypeString, Optional: true, Description: "Tag key."},
            "tag_value": {Type: schema.TypeString, Optional: true, Description: "Tag value."},
        },
    },
},
```

云 API 侧现状（vendor SDK `dts/v20211206/models.go`）：
- `CreateSyncJobRequest.Tags []*TagItem`：创建时按用户传入顺序透传给云 API
- `SyncJobInfo.Tags []*TagItem`（DescribeSyncJobs 响应）：查询返回的 tags 顺序由云端决定，**不保证与创建时传入顺序一致**

当前 Read 方法将云 API 返回的 tags 顺序原样写入 state。当返回顺序与用户配置顺序不同时，List 类型的顺序敏感比较产生永久性 diff。

约束：
- Terraform Plugin SDK v2
- vendor 模式管理依赖，不引入新依赖
- tags 为 `ForceNew: true`，不涉及 Update 路径

## Goals / Non-Goals

**Goals:**

- 将 `tags` 字段从 `schema.TypeList` 改为 `schema.TypeSet`，消除顺序敏感性，修复 plan 永久 diff 问题
- Create 路径正确读取 Set 类型参数并组装 `[]*dts.TagItem` 请求
- Read 路径将云 API 返回的 tags 正确写回 state（Set 语义）
- 保持 HCL 配置语法完全兼容（`tags { tag_key = ... tag_value = ... }` 写法不变）
- 更新单元测试与文档

**Non-Goals:**

- 不修改数据源 `tencentcloud_dts_sync_jobs` 的 tags 输出结构（独立 schema，本变更不涉及）
- 不修改服务层 `service_tencentcloud_dts.go`（DescribeDtsSyncJob 仅透传 SDK 结构）
- 不修改 vendor 下的云 SDK
- 不支持 tags 的 Update（保持 ForceNew 语义，与现状一致）
- 不处理 tencentcloud_dts_sync_job 以外的其他 DTS 资源

## Decisions

### 决策 1: 使用 `schema.TypeSet` 而非保持 List + 排序

**选择**: `Type: schema.TypeSet`，`Elem` 保持 `*schema.Resource{...}`（含 `tag_key`/`tag_value` 子字段）不变。

**理由**:
- Set 是 Terraform 表达无序集合的规范方式，语义上与"标签集合"天然匹配
- 读取端排序方案不可行：无法定义稳定的排序键使云 API 返回顺序与用户配置顺序一致（用户配置顺序本身任意）
- Set 中 `tag_key + tag_value` 组合作为哈希键天然去重，与腾讯云标签键唯一性约束一致
- 先例：项目内已有同类变更（`modify-mongodb-availability-zone-list-to-set`），社区惯例也是将无序返回的字段改为 Set

**备选方案（否决）**:
- 保持 List 并在 Read 后按 tag_key 排序：仅解决单向不一致（用户配置若按 value 排序仍会 diff），且引入人为排序规则
- 在 Create 前重排用户配置：不可行，用户配置顺序无法预知云 API 返回顺序

### 决策 2: Create 方法中的类型读取方式

`v.([]interface{})` 改为 `v.(*schema.Set).List()`，其余组装 `dts.TagItem` 的逻辑不变：

```go
if v, ok := d.GetOk("tags"); ok {
    for _, item := range v.(*schema.Set).List() {
        dMap := item.(map[string]interface{})
        tagItem := dts.TagItem{}
        if v, ok := dMap["tag_key"]; ok {
            tagItem.TagKey = helper.String(v.(string))
        }
        if v, ok := dMap["tag_value"]; ok {
            tagItem.TagValue = helper.String(v.(string))
        }
        request.Tags = append(request.Tags, &tagItem)
    }
}
```

**理由**: `d.GetOk` 对 TypeSet 返回 `*schema.Set`，`.List()` 转为 `[]interface{}` 后逐元素按 map 处理，循环体零改动，风险最小。

### 决策 3: Read 方法写回逻辑保持结构不变

现有逻辑组装 `[]interface{}{map[string]interface{}{"tag_key": ..., "tag_value": ...}}` 后 `d.Set("tags", tagsList)`。`d.Set` 对 Set 类型字段接受切片输入并由 SDK 按 Set 语义存储，无需显式构造 `*schema.Set`。

**理由**: 复用现有代码路径，改动面最小；SDK 的 `d.Set` 会基于 schema 定义自动处理类型转换。

### 决策 4: 测试与文档的更新范围

- `resource_tc_dts_sync_job_test.go`：`resource.TestCheckResourceAttrSet("...", "tags.#")` 在 Set 类型下依然有效（`tags.#` 是 Set 的计数索引），现有断言无需修改；验收测试配置 HCL 写法不变。本资源属于存量资源修改，遵循项目规则使用 terraform 测试套件补充测试用例（不引入 gomonkey）
- `resource_tc_dts_sync_job.md`：示例中 `tags` 块的 HCL 写法在 Set 下语义不变，文件保持现状；Argument Reference 部分由 `make doc` 工具自动生成（收尾阶段统一执行）

### 决策 5: 兼容性与状态迁移

- HCL 兼容：Set 与 List 的块语法完全相同，用户配置零改动
- State 迁移：SDK 在 Read 时按新 schema 写回，自动完成 List→Set 的状态迁移；首次 plan 可能显示一次性 tags 格式转换 diff，apply（或纯 refresh）后收敛
- 无需编写自定义 StateUpgraders：字段结构（tag_key/tag_value）未变，仅容器类型变化，SDK 原生处理

## Risks / Trade-offs

- **[风险] schema 类型变更为破坏性变更，存量用户升级后首次 plan 出现一次性 diff** → 在 changelog 中明确说明；diff 仅涉及 tags 的表示形式，apply 后收敛，无实际云端资源变更
- **[风险] Set 哈希去重导致用户配置重复 tag_key 时静默收敛** → 腾讯云标签键本身唯一，重复配置属于用户错误，Set 去重行为与云端约束一致
- **[风险] 遗漏对 Set 类型的读取断言导致运行时 panic** → Create/Read 路径均使用 `v.(*schema.Set).List()` 标准写法，且与项目内既有 Set 字段实现模式一致；单元测试覆盖 tags 用例
- **[代价] tags 内元素顺序在 plan/state 展示中由哈希决定，可读性略降** → 换取顺序无关的正确 diff 行为，符合 Terraform 惯例

## Migration Plan

1. 合并代码变更（schema + Create 读取逻辑）
2. 用户升级 Provider 后首次 `terraform plan`：可能显示 tags 的一次性格式转换 diff
3. `terraform apply`（或 `terraform plan -refresh-only`）后 state 收敛，后续 plan 无 diff
4. 回滚策略：回退 Provider 版本即可，state 中 Set 表示可被旧版本 List schema 重新读取（元素结构未变）
