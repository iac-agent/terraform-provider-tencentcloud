## Why

腾讯云 VPC 的协议端口模板（ServiceTemplate）API 支持通过 Tags 参数设置标签（Key/Value），但当前 Terraform 资源 `tencentcloud_protocol_template` 未暴露此功能。需要新增 `Key` 和 `Value` 参数以支持标签管理，提升资源管理灵活性。

## What Changes

- 为 `tencentcloud_protocol_template` 资源新增可选参数 `Key` 和 `Value`（对应云 API 的 Tags 参数）
- 修改资源创建方法以支持传入 Tags 参数到 `CreateServiceTemplate` 接口
- 修改资源读取方法以从 `DescribeServiceTemplates` 接口返回中提取 TagSet 信息
- 修改服务层方法 `CreateServiceTemplate` 和 `DescribeServiceTemplateById` 以支持 Tags 参数传递和返回

## Capabilities

### New Capabilities
- `protocol-template-tags`: 为协议端口模板资源新增标签（Key/Value）参数支持，允许用户在创建和管理协议端口模板时设置标签

### Modified Capabilities

（无现有能力需要修改）

## Impact

- **代码变更**：
  - `tencentcloud/services/vpc/resource_tc_protocol_template.go`：新增参数定义，修改 CRUD 方法
  - `tencentcloud/services/vpc/service_tencentcloud_vpc.go`：修改 `CreateServiceTemplate` 和 `DescribeServiceTemplateById` 方法签名和实现
  
- **API 变更**：使用云 API 已支持的 `Tags` 参数，无新增 API 调用

- **向后兼容性**：新增参数为可选（Optional），不影响现有资源配置和 state，完全向后兼容

- **文档变更**：需要更新资源文档（通过 `make doc` 自动生成）
