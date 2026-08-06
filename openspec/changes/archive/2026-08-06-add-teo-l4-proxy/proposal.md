## Why

TEO（Tencent Edge One）四层代理实例是边缘加速的关键基础设施，目前 Terraform Provider 缺少对四层代理实例生命周期的管理能力。用户无法通过 IaC 方式创建、修改、查询和删除四层代理实例，需要手动在控制台操作，不利于自动化运维和配置管理。

## What Changes

- 新增 `tencentcloud_teo_l4_proxy` 资源，支持四层代理实例的完整生命周期管理（创建、查询、修改、删除）
- 支持配置四层代理实例的加速区域、IPv6 访问、固定 IP、中国大陆网络优化等参数
- 支持 DDoS 防护配置（创建时指定）
- 支持通过 `zone_id` 和 `proxy_id` 联合 ID 导入已有资源

## Capabilities

### New Capabilities
- `teo-l4-proxy`: 提供 TEO 四层代理实例的 Terraform 资源管理能力，包括创建、读取、更新和删除四层代理实例

### Modified Capabilities
<!-- None -->

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_l4_proxy_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（更新文档索引）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已存在于 vendor）