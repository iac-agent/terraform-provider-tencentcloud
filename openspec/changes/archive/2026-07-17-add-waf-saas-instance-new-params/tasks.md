## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/waf/resource_tc_waf_saas_instance.go` 的 Schema 中添加 `goods_num` (Optional, TypeInt) 参数
- [x] 1.2 在 `tencentcloud/services/waf/resource_tc_waf_saas_instance.go` 的 Schema 中添加 `pid` (Optional, TypeInt) 参数
- [x] 1.3 在 `tencentcloud/services/waf/resource_tc_waf_saas_instance.go` 的 Schema 中添加 `region_id` (Optional, TypeInt) 参数
- [x] 1.4 在 `tencentcloud/services/waf/resource_tc_waf_saas_instance.go` 的 Schema 中添加 `deal_names` (Computed, TypeList, ElementType: schema.TypeString) 参数

## 2. CRUD 函数修改

- [x] 2.1 修改 Create 方法：当 `goods_num` 被设置时，使用用户指定的值替代硬编码的 1；未设置时保持默认值 1
- [x] 2.2 修改 Create 方法：当 `pid` 被设置时，使用用户指定的值替代 `PID_SAAS[goodsCategory]` 映射值；未设置时保持映射值
- [x] 2.3 修改 Create 方法：当 `region_id` 被设置时，使用用户指定的值替代 `mainlandMode`；未设置时保持从 provider region 推导的逻辑
- [x] 2.4 修改 Create 方法：创建成功后，从 `response.Data.DealNames` 提取值并设置到 `deal_names` 字段
- [x] 2.5 修改 Update 方法：将 `goods_num`, `pid`, `region_id` 加入 `immutableArgs` 数组，使这些参数在更新时不可变更

## 3. 测试

- [x] 3.1 在 `tencentcloud/services/waf/resource_tc_waf_saas_instance_test.go` 中补充 `goods_num` 参数的单元测试用例（使用 gomonkey mock）
- [x] 3.2 在 `tencentcloud/services/waf/resource_tc_waf_saas_instance_test.go` 中补充 `pid` 参数的单元测试用例（使用 gomonkey mock）
- [x] 3.3 在 `tencentcloud/services/waf/resource_tc_waf_saas_instance_test.go` 中补充 `region_id` 参数的单元测试用例（使用 gomonkey mock）
- [x] 3.4 在 `tencentcloud/services/waf/resource_tc_waf_saas_instance_test.go` 中补充 `deal_names` 参数的单元测试用例（使用 gomonkey mock）
- [x] 3.5 运行单元测试（go test -gcflags=all=-l）确保所有新增测试通过

## 4. 文档

- [x] 4.1 更新 `tencentcloud/services/waf/resource_tc_waf_saas_instance.md` 示例文件，添加 `goods_num`, `pid`, `region_id`, `deal_names` 参数的示例
