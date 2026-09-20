## Why

当前 `tencentcloud_teo_l4_proxy` 资源在调用 `DescribeL4Proxy` 接口时未暴露分页参数（`Offset`、`Limit`）和总数统计（`TotalCount`），用户无法通过 Terraform 配置控制分页查询行为，也无法获取四层代理实例的总数信息。该变更补充这些可选和只读参数，提升资源的可用性和完整性。

## What Changes

- 在 `tencentcloud_teo_l4_proxy` 资源中新增 `offset` 可选入参（TypeInt），用于控制 `DescribeL4Proxy` 接口的分页查询偏移量
- 在 `tencentcloud_teo_l4_proxy` 资源中新增 `limit` 可选入参（TypeInt），用于控制 `DescribeL4Proxy` 接口的分页查询限制数目
- 在 `tencentcloud_teo_l4_proxy` 资源中新增 `total_count` Computed 出参（TypeInt），用于暴露 `DescribeL4Proxy` 接口返回的四层代理实例总数

## Capabilities

### New Capabilities
- `teo-l4-proxy-pagination-params`: 为 `tencentcloud_teo_l4_proxy` 资源新增分页查询参数（offset、limit）和总数统计参数（total_count），支持用户在 Read 操作中控制分页行为并获取实例总数

### Modified Capabilities
<!-- None - no existing specs are modified -->

## Impact

- 修改 `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go`：Schema 定义和 Read 方法
- 修改 `tencentcloud/services/teo/service_tencentcloud_teo.go`：`DescribeTeoL4ProxyById` 方法签名和实现
- 修改 `tencentcloud/services/teo/resource_tc_teo_l4_proxy_test.go`：补充单元测试
- 修改 `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md`：文档更新
- 所有新增参数均为 Optional 或 Computed，不影响已有 Terraform 配置和 state，完全向后兼容