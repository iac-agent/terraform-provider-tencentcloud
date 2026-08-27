## Context

当前 TEO 产品线在 Terraform Provider 中已有若干数据源（如 `tencentcloud_teo_environments`、`tencentcloud_teo_zone_available_plans` 等），但缺少 DDoS 攻击 Top 数据查询能力。腾讯云 TEO SDK（`teov20220901`）已提供 `DescribeDDoSAttackTopData` 接口，可直接调用。

当前状态：
- SDK 已 vendored，`DescribeDDoSAttackTopData` 接口可用
- 请求结构：`DescribeDDoSAttackTopDataRequest` 包含 `StartTime`、`EndTime`、`MetricName`、`ZoneIds`、`PolicyIds`、`AttackType`、`ProtocolType`、`Port`、`Limit`、`Area`
- 响应结构：`DescribeDDoSAttackTopDataResponse` 包含 `Data []*TopEntry`、`TotalCount *uint64`
- `TopEntry` 包含 `Key *string` 和 `Value []*TopEntryValue`
- `TopEntryValue` 包含 `Name *string` 和 `Count *int64`

约束：
- 数据源为 RESOURCE_KIND_DATASOURCE，仅需实现 Read 方法
- 必须遵循 Provider 现有模式：`defer tccommon.LogElapsed()`、`defer tccommon.InconsistentCheck()`、`resource.Retry` 重试
- 数据源不应暴露分页参数（Limit）给用户，内部设置合理默认值
- 向后兼容：新增数据源，不影响现有资源

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_d_do_s_attack_top_data` 数据源，支持通过 Terraform 查询 DDoS 攻击 Top 数据
- 支持所有 API 入参（必填的 `start_time`、`end_time`、`metric_name` 及可选的 `zone_ids`、`policy_ids`、`attack_type`、`protocol_type`、`port`、`area`）
- 返回完整的 TopEntry 数据结构（含 `key` 和 `value` 嵌套列表）
- 通过 `result_output_file` 支持结果导出

**Non-Goals:**
- 不暴露 `Limit` 分页参数给用户（内部硬编码为合理值）
- 不暴露 `TotalCount` 为独立字段（可通过 `data` 列表长度推导）
- 不修改任何现有 TEO 资源或数据源

## Decisions

### Decision 1: 数据源文件命名为 `data_source_tc_teo_d_do_s_attack_top_data.go`

**选择**：文件名与 Terraform 资源名保持一致：`data_source_tc_teo_d_do_s_attack_top_data.go`。

**理由**：遵循 Provider 现有命名约定，数据源文件名格式为 `data_source_tc_<product>_<name>.go`。

### Decision 2: Schema 设计中 `data` 使用 TypeList 嵌套 TypeList

**选择**：
```go
"data": {
    Type:     schema.TypeList,
    Computed: true,
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "key": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "value": {
                Type:     schema.TypeList,
                Computed: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "count": {
                            Type:     schema.TypeInt,
                            Computed: true,
                        },
                    },
                },
            },
        },
    },
},
```

**备选**：将 `value` 扁平化为 TypeSet 或使用 JSON 字符串。

**理由**：TypeList 嵌套 TypeList 与云 API 响应结构（`[]*TopEntry` → `[]*TopEntryValue`）直接对应，且 Terraform 能够正确 set 和读取嵌套列表。参考现有 `data_source_tc_teo_environments.go` 中 `current_config_group_version_infos` 的嵌套 List 模式。

### Decision 3: `Limit` 参数内部硬编码为 100

**选择**：Read 函数中调用 API 时，将 `Limit` 设为 `helper.IntInt64(100)`（或更高值），不暴露给用户。

**理由**：
- 遵循 Provider 规范："数据源分页:不暴露 limit/offset 参数给用户,内部实现自动分页获取所有数据"
- 该 API 为 Top N 查询（无 Offset 参数，无法分页），设置较大的 Limit 以获取尽可能多的数据
- 100 是合理的最大值，API 默认值仅为 10

### Decision 4: ID 生成策略

**选择**：使用 `start_time`、`end_time`、`metric_name` 拼接生成 ID，格式为 `startTime#endTime#metricName`。

**备选**：使用固定 `"top_data"` 或时间戳。

**理由**：数据源 ID 需要唯一标识本次查询，组合三个必填参数可保证同一次查询幂等。同时这些参数在 API 层面足够区分不同查询。

### Decision 5: `policy_ids` 使用 TypeSet 而非 TypeList

**选择**：`zone_ids` 和 `policy_ids` 使用 `schema.TypeSet`。

**理由**：集合类型语义上更合适（ID 集合无顺序要求），且与 Provider 中其他类似参数（如 `zone_ids`）保持一致。

### Decision 6: Read 函数中调用 API 后处理空响应

**选择**：
```go
if response == nil || response.Response == nil || len(response.Response.Data) == 0 {
    log.Printf("[DATASOURCE] tencentcloud_teo_d_do_s_attack_top_data read empty, skip SetId")
    return fmt.Errorf("...")
}
```

**备选**：return nil（允许空结果）。

**理由**：遵循 RESOURCE_KIND_DATASOURCE 规范第 14 条：若 API 返回空，应返回 `NonRetryableError` 而非 `d.SetId("")`，避免因 API 短暂波动导致本地 state 中 id 被清空。

## Risks / Trade-offs

- **Risk**：API 的 `Limit` 参数未来可能增加默认值限制，导致硬编码 100 被拒绝 → **Mitigation**：若 API 返回错误，retry 机制会捕获并通过 `RetryError` 包装返回，用户可感知；后续可调整 Limit 值
- **Risk**：`metric_name` 参数取值较多且未来可能扩展，若用户传入无效值 → **Mitigation**：API 层面会返回参数错误，不需要 Provider 层面做 ValidateFunc 校验
- **Trade-off**：`policy_ids` 在 API 中类型为 `[]*int64`，Terraform 中 TypeSet 的元素类型为 TypeInt，需在调用时做类型转换 → 可接受，转换为 `[]*int64` 是标准模式

## Migration Plan

- 无需迁移：本次为纯新增数据源，不影响任何现有资源和 state
- 部署：新增 Go 文件、测试文件、.md 文档，在 `provider.go` 中注册即可
- 回滚：删除新增文件、移除注册代码即可

## Open Questions

- 无