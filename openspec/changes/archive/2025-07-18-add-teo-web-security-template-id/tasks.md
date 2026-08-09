## 1. Schema & CRUD 修改

- [x] 1.1 在 `resource_tc_teo_web_security_template.go` 的 Schema 中新增 `template_id` 字段（TypeString, Optional: true, Computed: true），添加描述说明
- [x] 1.2 在 Create 函数中，从 `CreateWebSecurityTemplate` API 响应获取 `TemplateId`，并赋值给 `templateId` 变量和 `d.Set("template_id", templateId)`
- [x] 1.3 在 Read 函数中，从复合 ID 拆分后设置 `d.Set("template_id", templateId)`

## 2. 文档更新

- [x] 2.1 更新 `resource_tc_teo_web_security_template.md` 示例文件，添加 `template_id` 字段说明

## 3. 测试

- [x] 3.1 在 `resource_tc_teo_web_security_template_test.go` 中补充 `template_id` 相关的单元测试用例
