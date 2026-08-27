## Why

EdgeOne（TEO）用户需要查询 DDoS 攻击 Top 数据（按协议、攻击类型、攻击源地区等维度统计排行），当前 Terraform Provider 尚未提供对应的数据源，导致用户无法通过声明式配置获取 DDoS 攻击分析数据，只能通过控制台或 SDK 手动查询，影响运维自动化和告警集成场景。

## What Changes

- 新增 `tencentcloud_teo_d_do_s_attack_top_data` 数据源（RESOURCE_KIND_DATASOURCE），调用腾讯云 TEO `DescribeDDoSAttackTopData` API。
- 入参：`start_time`（必填）、`end_time`（必填）、`metric_name`（必填）、`zone_ids`（可选）、`policy_ids`（可选）、`attack_type`（可选）、`protocol_type`（可选）、`port`（可选）、`area`（可选）。
- 出参：`data`（TopEntry 列表，包含 `key` 和 `value` 子字段），其中 `value` 为 `TopEntryValue` 列表（含 `name` 和 `count`）。
- 新增数据源代码文件 `data_source_tc_teo_d_do_s_attack_top_data.go` 及对应的单元测试文件、.md 文档文件。
- 在 `tencentcloud/provider.go` 中注册新数据源。

## Capabilities

### New Capabilities
- `teo-ddos-attack-top-data-datasource`: 提供 TEO DDoS 攻击 Top 数据的查询能力，支持按多种维度（协议、攻击类型、攻击源地区等）查询攻击排行数据。

### Modified Capabilities
<!-- 无，本次为新增数据源，不修改任何现有 capability -->

## Impact

- 代码：
  - `tencentcloud/services/teo/data_source_tc_teo_d_do_s_attack_top_data.go`（新增数据源）
  - `tencentcloud/services/teo/data_source_tc_teo_d_do_s_attack_top_data_test.go`（新增单元测试）
  - `tencentcloud/services/teo/data_source_tc_teo_d_do_s_attack_top_data.md`（新增文档）
  - `tencentcloud/provider.go`（注册数据源）
  - `tencentcloud/provider.md`（注册数据源文档条目）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.DescribeDDoSAttackTopData` 接口，无需变更 vendor。
- 向后兼容：新增数据源，不影响任何现有资源和状态。
- 文档：增加 `website/docs/d/teo_d_do_s_attack_top_data.html.markdown`（由 `make doc` 自动生成）。