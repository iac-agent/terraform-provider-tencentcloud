## 1. 需求分析和设计

- [ ] 1.1 详细分析云 API 接口，确定需要新增的具体参数
- [ ] 1.2 评估新增参数对现有资源和数据源的影响
- [ ] 1.3 设计参数在 Terraform schema 中的定义方式

## 2. 资源代码修改

- [ ] 2.1 修改 `tencentcloud/services/teo/resource_tc_teo_origin_acl.go` 文件，在 schema 中新增参数定义
- [ ] 2.2 更新资源的 Create 函数，支持新增参数的设置
- [ ] 2.3 更新资源的 Read 函数，支持新增参数的读取
- [ ] 2.4 更新资源的 Update 函数，支持新增参数的更新
- [ ] 2.5 更新资源的 Delete 函数（如需要）

## 3. 数据源代码修改

- [ ] 3.1 修改 `tencentcloud/services/teo/data_source_tc_teo_origin_acl.go` 文件，在 schema 中新增参数定义
- [ ] 3.2 更新数据源的 Read 函数，支持新增参数的读取和展示

## 4. 服务层修改

- [ ] 4.1 检查 `tencentcloud/services/teo/service_tencentcloud_teo.go` 是否需要修改以支持新增参数
- [ ] 4.2 更新相关的服务层函数（如需要）

## 5. 测试代码更新

- [ ] 5.1 更新 `tencentcloud/services/teo/resource_tc_teo_origin_acl_test.go` 文件，新增测试用例
- [ ] 5.2 更新 `tencentcloud/services/teo/data_source_tc_teo_origin_acl_test.go` 文件，新增测试用例
- [ ] 5.3 使用 mock 方法对云 API 进行 mock 处理，只进行业务代码逻辑的单元测试

## 6. 文档更新

- [ ] 6.1 更新 `tencentcloud/services/teo/resource_tc_teo_origin_acl.md` 文件，新增参数说明和示例
- [ ] 6.2 更新 `tencentcloud/services/teo/data_source_tc_teo_origin_acl.md` 文件，新增参数说明和示例
- [ ] 6.3 确保文档格式符合要求：一句话描述、Example Usage、Import（如适用）

## 7. 提供者注册检查

- [ ] 7.1 检查 `tencentcloud/provider.go` 中是否已正确注册资源和数据源
- [ ] 7.2 检查 `tencentcloud/provider.md` 中是否已正确注册资源和数据源

## 8. 代码验证

- [ ] 8.1 检查所有函数返回的 error 是否正确处理
- [ ] 8.2 检查代码是否符合 Terraform Provider 的开发规范
- [ ] 8.3 检查新增参数在云 API 接口中的存在性和一致性

## 9. 最终检查

- [ ] 9.1 确保所有修改的文件都已保存
- [ ] 9.2 检查代码的正确性和可编译性
- [ ] 9.3 准备提交代码变更
