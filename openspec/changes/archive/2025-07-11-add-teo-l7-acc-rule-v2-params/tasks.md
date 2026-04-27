## 1. Schema 定义

- [x] 1.1 在 `resource_tc_teo_l7_acc_rule_v2.go` 的 Schema 中添加 `filters` 可选参数（TypeList，包含 name 和 values 子字段），参考 SDK `Filter` 结构

## 2. Service 层更新

- [x] 2.1 更新 `service_tencentcloud_teo.go` 中的 `DescribeTeoL7AccRuleById` 方法，增加可选 filters 参数支持：当传入 filters 时使用自定义过滤条件，当未传入时保持现有 rule-id 过滤逻辑

## 3. Read 逻辑更新

- [x] 3.1 在 `ResourceTencentCloudTeoL7AccRuleV2Read` 中读取 `filters` 参数并传递给 service 层方法
- [x] 3.2 确保 filters 参数未指定时 Read 操作保持原有行为（向后兼容）

## 4. 单元测试

- [x] 4.1 在 `resource_tc_teo_l7_acc_rule_v2_test.go` 中添加 `filters` 参数的单元测试用例，使用 gomonkey mock 云 API

## 5. 文档更新

- [x] 5.1 更新 `resource_tc_teo_l7_acc_rule_v2.md` 示例文件，添加 filters 参数的说明和示例用法
