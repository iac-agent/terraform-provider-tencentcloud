## Context

TEO (Tencent EdgeOne) 四层代理实例（L4 Proxy）的 Terraform 资源 `tencentcloud_teo_l4_proxy` 已存在，本次新增 `tencentcloud_teo_l4_proxy_1` 作为其独立变体。TEO 是腾讯云边缘安全加速平台，四层代理实例是其核心组件之一，用于管理 TCP/UDP 层级的流量代理。

SDK 中四层代理实例提供了完整的 CRUD 接口链路：
- `CreateL4Proxy`：创建实例，返回 proxyId
- `DescribeL4Proxy`：查询实例列表（支持按 proxy-id 过滤）
- `ModifyL4Proxy`：修改实例（仅允许修改 ipv6、accelerate_mainland）
- `DeleteL4Proxy`：删除实例

此外还有状态管理相关接口：`ModifyL4ProxyStatus`（启停用实例）。

资源 ID 采用 `zone_id` + `proxy_id` 联合格式，以 `tccommon.FILED_SP` (即 `#`) 分隔。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_l4_proxy_1` 资源的完整 CRUD 生命周期管理
- 支持四层代理实例的创建（含 DDoS 防护配置）、查询、修改（ipv6/accelerate_mainland）和删除
- 支持 Terraform import（通过联合 ID：`zone_id#proxy_id`）
- 异步操作通过 StateChangeConf 轮询等待状态稳定
- 遵循与已有 `tencentcloud_teo_l4_proxy` 资源一致的代码模式

**Non-Goals:**
- 不修改已有的 `tencentcloud_teo_l4_proxy` 资源
- 不实现四层代理规则（L4ProxyRule）的管理（那是独立资源）
- 不支持 DDosProtectionConfig 在创建后的修改（云 API 已废弃该字段）

## Decisions

### 1. 资源 ID 设计
- **决策**: 使用 `zone_id#proxy_id` 作为复合 ID
- **理由**: `DescribeL4Proxy` 接口需要 `zone_id` 作为必填入参，仅凭 `proxy_id` 无法查询单个实例。复合 ID 保留了 zone_id 信息，支持 import 和 refresh。

### 2. Schema 设计 - area 字段
- **决策**: `area` 设为 Optional + ForceNew（不可变）
- **理由**: 根据用户需求映射，area 是 `CreateL4Proxy` 的必填参数，但云 API 的 `ModifyL4Proxy` 不支持修改 area。因此设为 ForceNew，修改会触发资源重建。

### 3. proxy_name 不可变性
- **决策**: `proxy_name` 在创建时必填，但创建后不可修改
- **理由**: `ModifyL4Proxy` 不支持修改 ProxyName，且 SDK 注释明确说明"创建完成后不支持修改"。

### 4. static_ip 不可变性
- **决策**: `static_ip` 设为 Optional + ForceNew
- **理由**: `ModifyL4Proxy` 不支持修改 StaticIp，修改需重建资源。

### 5. DDosProtectionConfig 处理
- **决策**: 保留 `d_dos_protection_config` 嵌套结构，但标记为 Deprecated
- **理由**: SDK 中 `CreateL4ProxyRequest.DDosProtectionConfig` 已标记 Deprecated，且 `L4Proxy` 返回体中也标记了 Deprecated。但为了向后兼容，仍需在 schema 中保留并支持读写。

### 6. 异步操作处理
- **决策**: Create 后等待 state 变为 `online`；Delete 前先停用（offline），再删除
- **理由**: 参考已有 `tencentcloud_teo_l4_proxy` 的实现模式。创建后需等待实例就绪；删除前需确保实例已停用。使用 `StateChangeConf` + `teoL4proxyStateRefreshFunc` 实现状态轮询。

### 7. DescribeL4Proxy 接口使用
- **决策**: 通过 `DescribeTeoL4ProxyById` 服务层函数封装查询逻辑，内部使用 Filters（proxy-id）分页查询
- **理由**: 遵循现有代码模式，由 service 层封装查询逻辑，Read 函数只负责从返回值映射到 schema。

### 8. 代码风格
- **决策**: 严格参考 `resource_tc_teo_l4_proxy.go` 的代码风格，包括命名约定、日志格式、错误处理模式
- **理由**: 保持代码库一致性，降低维护成本。

## Risks / Trade-offs

- **风险**: `DDosProtectionConfig` 字段已废弃，未来 SDK 可能移除该字段 → **缓解**: schema 中已标记 Deprecated，用户会被引导不使用该字段
- **风险**: 异步操作等待超时可能导致 terraform apply 失败 → **缓解**: 使用合理的超时时间（10x ReadRetryTimeout），并提供清晰的错误日志
- **风险**: 如果 `proxy_id` 在 response 中为空，会导致 nil pointer → **缓解**: Create 后检查 `response == nil || response.Response == nil || response.Response.ProxyId == nil`