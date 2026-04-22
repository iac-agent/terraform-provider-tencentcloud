## Why

`tencentcloud_teo_l7_acc_rule` 资源中的 `zone_id` 字段当前描述为 "Zone id, required field."，不够明确。用户可能不清楚该字段的具体约束条件，导致传入 null 或空字符串 ""，从而引发 API 调用失败。需要将描述更新为更明确的约束说明，告知用户 zone_id 必须输入有效值，不能为 null 或空字符串 ""。

## What Changes

- 修改 `tencentcloud_teo_l7_acc_rule` 资源中 `zone_id` 字段的 Description
- 将描述从 "Zone id, required field." 改为 "Zone id, which must be a valid value and cannot be null or empty string."
- 同步更新 website/docs/r/teo_l7_acc_rule.html.markdown 中的对应文档

## Capabilities

### New Capabilities

（无）

### Modified Capabilities

- `teo-l7-acc-rule`: 更新 `zone_id` 字段描述，明确 zone_id 必须输入有效值，不能为 null 或空字符串 ""。

## Impact

- 修改文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go` - 更新 zone_id 字段 Description
- 修改文件：`website/docs/r/teo_l7_acc_rule.html.markdown` - 更新 zone_id 字段描述文档
- 无功能逻辑变更，仅文档描述更新
