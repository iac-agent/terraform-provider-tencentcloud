## 1. 资源Schema与CRUD实现

- [x] 1.1 创建 `tencentcloud/services/cvm/resource_tc_cvm_chc_network_mode.go` 文件，定义资源schema（chc_ids: TypeList/ForceNew, network_mode: TypeString）和ResourceTencentCloudCvmChcNetworkMode函数
- [x] 1.2 实现resourceTencentCloudCvmChcNetworkModeCreate函数，调用ModifyChcNetworkMode API，使用chc_ids拼接作为资源ID
- [x] 1.3 实现resourceTencentCloudCvmChcNetworkModeRead函数，调用DescribeChcHosts API验证CHC主机存在，network_mode从state保留
- [x] 1.4 实现resourceTencentCloudCvmChcNetworkModeUpdate函数，当network_mode变化时调用ModifyChcNetworkMode API
- [x] 1.5 实现resourceTencentCloudCvmChcNetworkModeDelete函数，仅从state移除资源，不调用云API

## 2. Provider注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 tencentcloud_cvm_chc_network_mode 资源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 tencentcloud_cvm_chc_network_mode 资源条目

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/cvm/resource_tc_cvm_chc_network_mode.md` 文档文件，包含一句话描述（提及CVM产品）、Example Usage（HCL配置示例）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/cvm/resource_tc_cvm_chc_network_mode_test.go` 测试文件，使用mock(gomonkey)方式编写CRUD操作的单元测试
- [x] 4.2 使用 `go test -gcflags=all=-l` 运行单元测试并确保通过

## 5. 代码验证

- [x] 5.1 检查所有函数返回的error，确保error被正确处理
- [x] 5.2 检查Create函数中对API返回值空值的判断（Response为nil、ChcIds为空等）
- [x] 5.3 检查Read函数中对空响应的处理，先打印log再SetId("")
- [x] 5.4 检查CRUD代码中retry块的使用符合规范（retry块内仅调用API，设置id等操作在retry外）
