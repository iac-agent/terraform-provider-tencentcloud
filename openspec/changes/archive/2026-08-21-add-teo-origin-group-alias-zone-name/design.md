## Context

`tencentcloud_teo_origin_group` is a `RESOURCE_KIND_GENERAL` Terraform resource that manages the full CRUD lifecycle of a TEO origin group. Its `references` block is a computed list that describes which acceleration domains, rule engines, load balancers, or application proxies reference the origin group.

The vendored TEO SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`) exposes `OriginGroupReference.AliasZoneName` in the `DescribeOriginGroup` response. This field carries the alias zone name of the referenced instance and is useful for cross-zone reference scenarios.

Current state:
- `tencentcloud/services/teo/resource_tc_teo_origin_group.go` defines the `references` nested block with computed `instance_type`, `instance_id`, `instance_name` fields.
- The `Read` method already iterates `respData.References` and populates those computed fields.
- The SDK field `AliasZoneName` is available but not yet exposed in the Terraform schema.

Constraints:
- Backward compatibility: this must be a purely additive computed field; it must not break existing configuration or state.
- The field is read-only (computed); it must not be added to Create/Update/Delete request building.
- `references` is already a `TypeList`/`Computed` block; only the nested field is being added.

## Goals / Non-Goals

**Goals:**
- Expose `alias_zone_name` as a computed string field inside the `references` block of `tencentcloud_teo_origin_group`.
- Populate the field from `OriginGroupReference.AliasZoneName` during the Read operation.
- Cover the new field with gomonkey-based unit tests.

**Non-Goals:**
- Do not change Create/Update/Delete operations.
- Do not change the existing `instance_type`, `instance_id`, `instance_name` behavior.
- Do not add a new resource or data source.
- Do not add `alias_zone_name` to any request (Create/Modify) — the SDK request structs do not accept it.

## Decisions

### Decision 1: Add `alias_zone_name` as a computed field in the `references` nested block

**选择**: 在 `references` 的 `Elem` schema 中新增 `alias_zone_name`，类型为 `schema.TypeString`，属性为 `Computed: true`。

**备选**: 将其作为顶层 computed 字段暴露。

**理由**:
- `AliasZoneName` 在 SDK 中属于 `OriginGroupReference`（`References` 列表的每个元素），语义上必须挂在 `references` 的每个条目下。
- 与其他 computed 字段（`instance_type` / `instance_id` / `instance_name`）保持一致的层级与风格。

### Decision 2: 仅在 Read 中 nil-check 后回填

**选择**: 在 `resourceTencentCloudTeoOriginGroupRead` 的 `references` 循环中，新增：

```go
if references.AliasZoneName != nil {
    referencesMap["alias_zone_name"] = references.AliasZoneName
}
```

**理由**:
- 遵循现有 computed 字段的 nil-check 模式，避免云 API 返回 nil 时写入空值。
- 不修改 `service.DescribeTeoOriginGroupById`，因为 `OriginGroup` 已经包含完整 `References` 数据。

### Decision 3: 使用 gomonkey mock 覆盖单测

**选择**: 在 `resource_tc_teo_origin_group_test.go` 中使用 gomonkey mock `DescribeOriginGroup`（或 service 层读取函数），构造包含 `AliasZoneName` 的响应，断言 state 中 `references.0.alias_zone_name` 正确回填。

**理由**:
- 与资源既有测试方式一致（gomonkey mock 云 API，仅测试业务逻辑）。
- 不依赖真实云资源与 `TF_ACC` 环境。

### Decision 4: 文档仅校验，不强制改示例 HCL

**选择**: 由于 `references` 是 computed 块，示例 HCL 中通常不写 computed 字段；本变更仅需确认 `resource_tc_teo_origin_group.md` 的描述与自动生成的文档一致，不手动添加 computed 字段示例。

**理由**:
- `references` 字段由云端回填，无法在 HCL 中设置；示例中展示该字段反而会造成误导。

## Risks / Trade-offs

- **Risk**: 云 API 返回的 `AliasZoneName` 为 nil（某些引用类型可能不返回） → **Mitigation**: nil-check 跳过，与现有 computed 字段行为一致。
- **Risk**: 该字段与既有 spec `teo-origin-group-reference-zone-fields` 存在重叠 → **Mitigation**: 本变更以独立 capability 记录 `alias_zone_name`，不影响既有 requirement 的其他字段。
- **Trade-off**: 纯 computed 字段无法被用户直接断言 diff → 可接受，computed 字段天然只读，由 state 回填。

## Migration Plan

- 纯增量变更：仅新增 computed 字段，无需 state 迁移。
- 存量资源 refresh 后会自动回填 `alias_zone_name`（若 API 返回）。
- 回滚：移除 schema 字段与 Read 回填逻辑即可，state 中多余值由 Terraform 自动忽略。

## Open Questions

- 无
