## MODIFIED Requirements

### Requirement: zone_id 字段描述约束
当用户使用 `tencentcloud_teo_l7_acc_rule` 资源时，`zone_id` 字段的描述 SHALL 明确告知用户该字段必须输入有效值，不能为 null 或空字符串 ""。

#### Scenario: 查看字段描述
- **WHEN** 用户查看 `zone_id` 字段的描述信息
- **THEN** 描述应说明 "Zone id, which must be a valid value and cannot be null or empty string."

#### Scenario: 传入 null 值
- **WHEN** 用户传入 null 作为 zone_id 值
- **THEN** Terraform 在 plan 阶段提示该字段为 Required，不允许为空

#### Scenario: 传入空字符串
- **WHEN** 用户传入空字符串 "" 作为 zone_id 值
- **THEN** 云 API 调用失败，返回参数错误
