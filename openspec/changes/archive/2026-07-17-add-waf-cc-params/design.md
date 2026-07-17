## Context

tencentcloud_waf_cc 是 WAF 产品的 CC 防护规则资源，使用 UpsertCCRule 接口进行创建和更新，DescribeCCRuleList 接口进行读取，DeleteCCRule 接口进行删除。当前资源缺少定时执行相关参数（JobType、StartDateTime、TimeTZone）及 DeleteCCRule 出参 Data。

云 API 结构：
- UpsertCCRule 接口支持 `JobType`（*string）和 `JobDateTime`（*JobDateTime）入参
- JobDateTime 是嵌套结构体，包含 `Timed`（[]*TimedJob）和 `TimeTZone`（*string）
- TimedJob 包含 `StartDateTime`（*uint64）和 `EndDateTime`（*uint64）
- DescribeCCRuleList 返回的 CCRuleItems 已包含 JobType 和 JobDateTime 字段
- DeleteCCRule 出参包含 Data（*string）

当前资源 schema 未包含这些参数，需要平铺嵌套结构到 terraform schema 顶层。

## Goals / Non-Goals

**Goals:**
- 为 tencentcloud_waf_cc 资源新增 job_type、start_date_time、time_t_zone 三个 Optional 参数
- 为 tencentcloud_waf_cc 资源新增 data Computed 出参（DeleteCCRule 返回值）
- 在 Create/Update 中正确构造 JobDateTime 嵌套结构并传入 UpsertCCRule
- 在 Read 中从 CCRuleItems.JobDateTime 中读取并平铺到 terraform state
- 补充单元测试验证新增参数的正确性
- 更新资源文档

**Non-Goals:**
- 不支持 CronJob 周期执行的详细参数（Days/WDays/StartTime/EndTime），需求仅涉及 TimedJob 的 StartDateTime
- 不支持 EndDateTime 参数（需求未涉及）
- 不修改资源的 ID 格式或生命周期逻辑

## Decisions

1. **嵌套结构平铺方案**：JobDateTime 是嵌套结构体（含 Timed 数组和 TimeTZone 字段），需求仅涉及 Timed[0].StartDateTime 和 TimeTZone。采用平铺方式将 start_date_time 和 time_t_zone 作为顶层 schema 字段，避免引入嵌套 block 增加用户配置复杂度。

2. **start_date_time 类型选择**：云 API 中 StartDateTime 为 *uint64（秒级时间戳），terraform 中使用 TypeInt 表示。当用户配置了 start_date_time 时，在 Create/Update 中构造 JobDateTime.Timed 数组，将 StartDateTime 设为用户指定值，EndDateTime 不设置（需求未涉及）。

3. **data 参数设计**：DeleteCCRule 出参 Data 为 *string，在 terraform 中作为 Computed TypeString 参数。由于 Delete 操作后资源已被销毁，data 的实际用途有限，但按需求需要暴露该字段。

4. **JobDateTime 构造时机**：仅在 job_type 或 start_date_time 或 time_t_zone 有值时才构造 JobDateTime 结构体传入请求，避免发送空结构体。

## Risks / Trade-offs

- [Risk] StartDateTime 仅取 Timed 数组第一个元素 → 若云 API 返回多个 Timed 项则丢失信息 → 当前云 API 实际只返回一个 Timed 项，风险可控
- [Risk] 平铺 JobDateTime 嵌套结构限制了未来扩展性 → 若后续需支持 CronJob 参数可能需要调整 schema → 当前按需求最小化实现，后续可按需添加
