## 1. Schema 定义

- [x] 1.1 在 tencentcloud/services/waf/resource_tc_waf_cc.go 的 Schema 中新增 `job_type`（Optional, TypeString）参数，描述规则执行方式
- [x] 1.2 在 tencentcloud/services/waf/resource_tc_waf_cc.go 的 Schema 中新增 `start_date_time`（Optional, TypeInt）参数，描述定时执行的开始时间戳（秒）
- [x] 1.3 在 tencentcloud/services/waf/resource_tc_waf_cc.go 的 Schema 中新增 `time_t_zone`（Optional, TypeString）参数，描述时区
- [x] 1.4 在 tencentcloud/services/waf/resource_tc_waf_cc.go 的 Schema 中新增 `data`（Computed, TypeString）参数，对应 DeleteCCRule 出参

## 2. Create 函数修改

- [x] 2.1 在 resourceTencentCloudWafCcCreate 中读取 job_type 参数并设置到 request.JobType
- [x] 2.2 在 resourceTencentCloudWafCcCreate 中读取 start_date_time 和 time_t_zone 参数，构造 JobDateTime 结构体（含 Timed[0].StartDateTime 和 TimeTZone）并设置到 request.JobDateTime

## 3. Read 函数修改

- [x] 3.1 在 resourceTencentCloudWafCcRead 中从 CCRuleItems 读取 JobType 并设置到 d.Set("job_type", ...)
- [x] 3.2 在 resourceTencentCloudWafCcRead 中从 CCRuleItems.JobDateTime.Timed[0].StartDateTime 读取并设置到 d.Set("start_date_time", ...)
- [x] 3.3 在 resourceTencentCloudWafCcRead 中从 CCRuleItems.JobDateTime.TimeTZone 读取并设置到 d.Set("time_t_zone", ...)

## 4. Update 函数修改

- [x] 4.1 在 resourceTencentCloudWafCcUpdate 中读取 job_type 参数并设置到 request.JobType
- [x] 4.2 在 resourceTencentCloudWafCcUpdate 中读取 start_date_time 和 time_t_zone 参数，构造 JobDateTime 结构体并设置到 request.JobDateTime

## 5. Delete 函数修改

- [x] 5.1 在 resourceTencentCloudWafCcDelete 中处理 DeleteCCRule 返回的 Data 出参并设置到 d.Set("data", ...)

## 6. 单元测试

- [x] 6.1 在 tencentcloud/services/waf/resource_tc_waf_cc_test.go 中补充 job_type、start_date_time、time_t_zone、data 参数的单元测试用例，使用 gomonkey mock 云 API

## 7. 文档更新

- [x] 7.1 更新 tencentcloud/services/waf/resource_tc_waf_cc.md 文档，补充新增参数的 Example Usage
