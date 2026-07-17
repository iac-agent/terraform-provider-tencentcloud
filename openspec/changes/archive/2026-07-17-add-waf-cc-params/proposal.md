## Why

tencentcloud_waf_cc 资源当前缺少 CC 规则的定时执行相关参数（JobType、StartDateTime、TimeTZone）以及 DeleteCCRule 接口的 Data 出参，导致用户无法通过 Terraform 配置 CC 规则的定时/周期执行策略，也无法获取删除操作的返回数据。云 API 已支持这些参数，需要在 Terraform 侧补齐。

## What Changes

- 为 tencentcloud_waf_cc 资源新增 Optional 参数 `job_type`（TypeString），对应 UpsertCCRule 接口的 request.JobType 和 DescribeCCRuleList 返回的 CCRuleItems.JobType，用于配置规则执行方式（TimedJob 定时执行 / CronJob 周期执行）
- 为 tencentcloud_waf_cc 资源新增 Optional 参数 `start_date_time`（TypeInt），对应 UpsertCCRule 接口的 request.JobDateTime.Timed[0].StartDateTime 和 DescribeCCRuleList 返回的 CCRuleItems.JobDateTime.Timed[0].StartDateTime，用于配置定时执行的开始时间戳（秒）
- 为 tencentcloud_waf_cc 资源新增 Optional 参数 `time_t_zone`（TypeString），对应 UpsertCCRule 接口的 request.JobDateTime.TimeTZone 和 DescribeCCRuleList 返回的 CCRuleItems.JobDateTime.TimeTZone，用于配置时区
- 为 tencentcloud_waf_cc 资源新增 Computed 参数 `data`（TypeString），对应 DeleteCCRule 接口出参 response.Data

## Capabilities

### New Capabilities
- `waf-cc-scheduled-execution`: 为 tencentcloud_waf_cc 资源新增定时执行相关参数（job_type、start_date_time、time_t_zone）及 DeleteCCRule 出参 data

### Modified Capabilities

## Impact

- 受影响文件：tencentcloud/services/waf/resource_tc_waf_cc.go（schema 定义、Create/Read/Update 逻辑）
- 受影响文件：tencentcloud/services/waf/resource_tc_waf_cc_test.go（补充单元测试）
- 受影响文件：tencentcloud/services/waf/resource_tc_waf_cc.md（文档更新）
- 依赖：云 API 已在 vendor 中支持 JobType、JobDateTime（含 Timed/TimeTZone）等字段
- 向后兼容：所有新增参数均为 Optional/Computed，不影响现有配置和 state
