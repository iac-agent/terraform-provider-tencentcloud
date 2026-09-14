## 1. Schema & Read 逻辑修改

- [x] 1.1 在 `tencentcloud/services/sms/resource_tc_sms_template.go` 的 Schema 中新增 `review_reply` 字段：`TypeString`, `Computed: true`, Description 描述审核回复
- [x] 1.2 在 `resourceTencentCloudSmsTemplateRead` 方法中，在 `template.International` 的 nil check 之后，添加 `template.ReviewReply` 的 nil check 并调用 `d.Set("review_reply", template.ReviewReply)`

## 2. 文档更新

- [x] 2.1 更新 `tencentcloud/services/sms/resource_tc_sms_template.md`，在属性列表中添加 `review_reply` 的说明

## 3. 单元测试补充

- [x] 3.1 在 `tencentcloud/services/sms/resource_tc_sms_template_test.go` 中补充 `review_reply` 字段的单元测试用例，使用 gomonkey mock 云 API 返回值验证 Read 逻辑

## 4. 验证

- [x] 4.1 检查代码正确性：确认 `review_reply` 为 Computed 字段、Read 方法中有 nil check、云 API `DescribeTemplateListStatus` 结构体中确实包含 `ReviewReply` 字段
