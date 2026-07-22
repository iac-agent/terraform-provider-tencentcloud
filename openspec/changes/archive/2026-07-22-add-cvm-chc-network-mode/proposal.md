## Why

CVM产品需要支持通过Terraform切换CHC物理服务器的网络模式（部署网络模式/业务网络模式）。当前缺乏对`ModifyChcNetworkMode`接口的Terraform资源封装，用户无法通过Terraform管理CHC服务器的网络模式切换。

## What Changes

- 新增 `tencentcloud_cvm_chc_network_mode` 资源（RESOURCE_KIND_GENERAL），封装 `ModifyChcNetworkMode` 云API接口
- 支持对CHC服务器网络模式的创建（切换到指定模式）、更新（切换到另一模式）和删除（资源从state中移除）
- 该资源仅有一个Modify接口，属于CRD类型资源，使用`DescribeChcHosts`接口进行Read操作
- 在 `provider.go` 和 `provider.md` 中注册新资源

## Capabilities

### New Capabilities
- `cvm-chc-network-mode-resource`: 封装ModifyChcNetworkMode接口，支持通过Terraform管理CHC服务器网络模式切换，包含chc_ids和network_mode两个参数

### Modified Capabilities

## Impact

- 新增资源文件: `tencentcloud/services/cvm/resource_tc_cvm_chc_network_mode.go`
- 新增测试文件: `tencentcloud/services/cvm/resource_tc_cvm_chc_network_mode_test.go`
- 修改: `tencentcloud/provider.go` - 注册新资源
- 修改: `tencentcloud/provider.md` - 添加新资源文档条目
- 新增文档: `tencentcloud/services/cvm/resource_tc_cvm_chc_network_mode.md`
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312` 中的 `ModifyChcNetworkMode` 和 `DescribeChcHosts` 接口
