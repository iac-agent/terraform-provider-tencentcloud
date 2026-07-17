## ADDED Requirements

### Requirement: WAF CC rule supports job_type parameter
tencentcloud_waf_cc 资源 SHALL 支持 `job_type` 参数（Optional, TypeString），用于配置规则执行方式。取值为 TimedJob（定时执行）或 CronJob（周期执行）。该参数 SHALL 在 Create 和 Update 操作中传入 UpsertCCRule 接口的 request.JobType，并在 Read 操作中从 DescribeCCRuleList 返回的 CCRuleItems.JobType 读取。

#### Scenario: Create CC rule with job_type specified
- **WHEN** user creates a tencentcloud_waf_cc resource with job_type = "TimedJob"
- **THEN** the resource SHALL pass job_type to UpsertCCRule request.JobType
- **THEN** after read, the resource SHALL have job_type set to "TimedJob"

#### Scenario: Create CC rule without job_type
- **WHEN** user creates a tencentcloud_waf_cc resource without specifying job_type
- **THEN** the resource SHALL NOT pass JobType to UpsertCCRule request
- **THEN** the cloud API default behavior SHALL apply

#### Scenario: Update CC rule with job_type change
- **WHEN** user updates a tencentcloud_waf_cc resource changing job_type from "TimedJob" to "CronJob"
- **THEN** the resource SHALL pass the new job_type to UpsertCCRule request.JobType

### Requirement: WAF CC rule supports start_date_time parameter
tencentcloud_waf_cc 资源 SHALL 支持 `start_date_time` 参数（Optional, TypeInt），用于配置定时执行的开始时间戳（秒）。该参数对应云 API 的 request.JobDateTime.Timed[0].StartDateTime（uint64）。在 Create 和 Update 操作中，当 start_date_time 有值时，SHALL 构造 JobDateTime 结构体并设置 Timed[0].StartDateTime。在 Read 操作中，SHALL 从 CCRuleItems.JobDateTime.Timed[0].StartDateTime 读取值。

#### Scenario: Create CC rule with start_date_time specified
- **WHEN** user creates a tencentcloud_waf_cc resource with start_date_time = 1700000000
- **THEN** the resource SHALL construct JobDateTime with Timed[0].StartDateTime = 1700000000 and pass to UpsertCCRule
- **THEN** after read, the resource SHALL have start_date_time set to 1700000000

#### Scenario: Create CC rule without start_date_time
- **WHEN** user creates a tencentcloud_waf_cc resource without specifying start_date_time
- **THEN** the resource SHALL NOT construct JobDateTime.Timed items for StartDateTime

#### Scenario: Read CC rule with start_date_time from API
- **WHEN** the DescribeCCRuleList API returns CCRuleItems with JobDateTime.Timed[0].StartDateTime = 1700000000
- **THEN** the resource SHALL set start_date_time = 1700000000 in terraform state

#### Scenario: Read CC rule without JobDateTime.Timed from API
- **WHEN** the DescribeCCRuleList API returns CCRuleItems with JobDateTime nil or Timed empty
- **THEN** the resource SHALL NOT set start_date_time

### Requirement: WAF CC rule supports time_t_zone parameter
tencentcloud_waf_cc 资源 SHALL 支持 `time_t_zone` 参数（Optional, TypeString），用于配置时区。该参数对应云 API 的 request.JobDateTime.TimeTZone。在 Create 和 Update 操作中，当 time_t_zone 有值时，SHALL 构造或更新 JobDateTime 结构体并设置 TimeTZone。在 Read 操作中，SHALL 从 CCRuleItems.JobDateTime.TimeTZone 读取值。

#### Scenario: Create CC rule with time_t_zone specified
- **WHEN** user creates a tencentcloud_waf_cc resource with time_t_zone = "UTC+8"
- **THEN** the resource SHALL construct JobDateTime with TimeTZone = "UTC+8" and pass to UpsertCCRule
- **THEN** after read, the resource SHALL have time_t_zone set to "UTC+8"

#### Scenario: Create CC rule without time_t_zone
- **WHEN** user creates a tencentcloud_waf_cc resource without specifying time_t_zone
- **THEN** the resource SHALL NOT set JobDateTime.TimeTZone

#### Scenario: Read CC rule with time_t_zone from API
- **WHEN** the DescribeCCRuleList API returns CCRuleItems with JobDateTime.TimeTZone = "UTC+8"
- **THEN** the resource SHALL set time_t_zone = "UTC+8" in terraform state

### Requirement: WAF CC rule supports data computed parameter
tencentcloud_waf_cc 资源 SHALL 支持 `data` 参数（Computed, TypeString），对应 DeleteCCRule 接口出参 response.Data。

#### Scenario: Delete CC rule returns data
- **WHEN** a tencentcloud_waf_cc resource is deleted and DeleteCCRule returns Data in response
- **THEN** the Data value SHALL be available as a computed attribute on the resource
