## Context

TEO（Tencent Edge One）提供四层代理服务，目前 Terraform Provider 中已存在 TEO 相关资源（如 `tencentcloud_teo_zone`、`tencentcloud_teo_dns_record` 等），但缺少 `tencentcloud_teo_l4_proxy` 资源来管理四层代理实例。云 API 已提供完整的 CreateL4Proxy、DescribeL4Proxy、ModifyL4Proxy、DeleteL4Proxy 接口，vendor 中已包含这些接口的定义。

参考资源：`tencentcloud/services/teo/` 目录下的现有 TEO 资源（如 `resource_tc_teo_dns_record.go`），遵循项目统一的 CRUD 模式、retry 机制、错误处理、复合 ID 格式等约定。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_l4_proxy` 资源的完整 CRUD 操作
- 使用 `zone_id#proxy_id` 复合 ID 格式，支持 terraform import
- 遵循项目现有的代码规范和错误处理模式
- 为可修改字段（ipv6、accelerate_mainland）提供 Update 支持
- 为不可修改字段（proxy_name、area、static_ip）标记 ForceNew

**Non-Goals:**
- 不支持通过 ModifyL4Proxy 修改 DDosProtectionConfig（API CreateL4Proxy 中已标记此字段为 deprecated）
- 不修改已有 TEO 资源的 schema 或行为
- 不创建 datasource 资源（本变更仅包含 RESOURCE_KIND_GENERAL 资源）

## Decisions

### 1. 复合 ID 格式：`zone_id` + `FILED_SP` + `proxy_id`
- **理由**: DescribeL4Proxy 接口需要同时传入 zone_id 和 filters（含 proxy-id），两者都是定位四层代理实例的必要条件。使用复合 ID 确保 import 后能正确解析。
- **替代方案**: 单独使用 proxy_id 作为 ID → 不可行，因为 Read 时缺少 zone_id 无法调用 DescribeL4Proxy。

### 2. 参数 readonly 策略
- **proxy_name**: 仅 Create 时传入，ForceNew。API 注释明确"创建完成后不支持修改"。
- **area**: 仅 Create 时传入，ForceNew。ModifyL4Proxy 接口不支持修改加速区域。
- **static_ip**: 仅 Create 时传入，ForceNew。ModifyL4Proxy 接口不支持修改固定 IP 配置。
- **d_dos_protection_config**: 仅 Create 时传入，Optional。ModifyL4Proxy 不支持修改此配置，且 CreateL4Proxy 中已标记此字段为 deprecated。
- **ipv6**: Create + Update 支持，ModifyL4Proxy 接口支持修改。
- **accelerate_mainland**: Create + Update 支持，ModifyL4Proxy 接口支持修改。

### 3. Read 实现策略
- 使用 DescribeL4Proxy 接口，通过 zone_id + Filters（proxy-id）精确查询单个实例
- 设置 Limit 为最大值（1000），避免分页问题
- 检查 Response 是否为 nil、L4Proxies 是否为空，若为空则 SetId("") 表示资源已被删除
- 所有 computed 字段（proxy_name, area, cname, ips, status, static_ip, l4_proxy_rule_count, update_time）在 Read 时设置

### 4. Retry 策略
- 所有云 API 调用使用 `tccommon.ReadRetryTimeout` 作为超时
- 调用失败时使用 `tccommon.RetryError()` 包装错误
- Retry 块内仅执行 API 调用，设置 ID 等成功操作放到 retry 块外

### 5. 错误处理
- 使用 `defer tccommon.LogElapsed()` 记录耗时
- 使用 `defer tccommon.InconsistentCheck()` 检查一致性
- Create 成功后检查 ProxyId 是否为空，若为空返回 NonRetryableError
- Read 中若 API 返回空响应或空列表，SetId("") 前打印 log.Printf 保留现场

### 6. Schema 结构
- `zone_id`: Required, ForceNew, TypeString
- `proxy_id`: Computed, TypeString（由 API 返回）
- `proxy_name`: Required, ForceNew, TypeString
- `area`: Required, ForceNew, TypeString
- `ipv6`: Optional, TypeString（取值 on/off）
- `static_ip`: Optional, ForceNew, TypeString（取值 on/off）
- `accelerate_mainland`: Optional, TypeString（取值 on/off）
- `d_dos_protection_config`: Optional, ForceNew, TypeList（嵌套 DDosProtectionConfig 结构）
- `cname`: Computed, TypeString
- `ips`: Computed, TypeList
- `status`: Computed, TypeString
- `l4_proxy_rule_count`: Computed, TypeInt
- `update_time`: Computed, TypeString

## Risks / Trade-offs

- [ModifyL4Proxy 仅返回 RequestId] → Update 后通过 Read 方法验证修改结果，确保状态一致性
- [DeleteL4Proxy 仅返回 RequestId] → Delete 后通过 Read 方法轮询确认资源已删除
- [DescribeL4Proxy 通过 Filters 匹配] → 使用 proxy-id filter 精确匹配，设置 Limit 为最大值避免遗漏
- [DDosProtectionConfig 已废弃] → 仅支持在 Create 时传入，Read 时读取但此字段来自 API 结果（标记为 deprecated），schema 中加入 Deprecated 说明
