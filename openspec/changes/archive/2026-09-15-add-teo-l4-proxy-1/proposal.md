## Why

TEO (Tencent EdgeOne) 四层代理实例的 Terraform 资源 `tencentcloud_teo_l4_proxy_1` 需要作为新资源引入，以支持 TEO 四层代理实例的全生命周期管理。当前已有 `tencentcloud_teo_l4_proxy` 资源，本次新增 `tencentcloud_teo_l4_proxy_1` 作为其升级版本，提供更完善的参数映射和接口覆盖。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_l4_proxy_1`，基于 TEO v20220901 SDK 的 `CreateL4Proxy`、`DescribeL4Proxy`、`ModifyL4Proxy`、`DeleteL4Proxy` 四个接口实现完整的 CRUD 操作
- 资源支持 zone_id 和 proxy_id 联合作为资源 ID，支持 Terraform import
- 支持创建时的 DDoS 防护配置（`d_dos_protection_config`，已标记 Deprecated）
- 支持 IPv6、固定 IP、中国大陆网络优化等可选配置
- 不可变参数（创建后不可修改）：`proxy_name`、`area`、`static_ip`、`ddos_protection_config`、`zone_id`
- 可变参数（可通过 ModifyL4Proxy 修改）：`ipv6`、`accelerate_mainland`
- 在 `tencentcloud/provider.go` 中注册新资源

## Capabilities

### New Capabilities
- `teo-l4-proxy-1-resource`: TEO 四层代理实例的 Terraform 资源管理，包括创建、读取、更新、删除和导入功能，覆盖 proxy_id、zone_id、proxy_name、area、ipv6、static_ip、accelerate_mainland、ddos_protection_config 等参数

### Modified Capabilities
<!-- None -->

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_l4_proxy_1.go`（资源 CRUD 实现）
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_l4_proxy_1_test.go`（单元测试）
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_l4_proxy_1.md`（文档模板）
- 修改文件: `tencentcloud/provider.go`（资源注册）
- 新增文件: `.changelog/<PR_NUMBER>.txt`（changelog，由 tfpacer-finalize 阶段生成）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（vendor 中已就绪）