## Context

Terraform Provider for TencentCloud 目前未支持 TEO（EdgeOne）多通道安全加速网关线路资源的管理。TEO 多通道安全加速网关线路用于配置网关的不同线路类型（直连线路、EdgeOne 四层代理线路、自定义线路），支持线路的创建、查询、修改和删除操作。

当前 TEO 服务下已有多种 Terraform 资源（如 `tencentcloud_teo_multi_path_gateway` 等），但缺少线路级别的管理资源。本次新增 `tencentcloud_teo_multi_path_gateway_line` 资源，遵循现有 TEO 资源的代码组织模式和架构风格。

云 API 接口信息（SDK 包: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`）：
- CreateMultiPathGatewayLine: 创建线路，入参 ZoneId/GatewayId/LineType/LineAddress/ProxyId/RuleId，出参 LineId
- DescribeMultiPathGatewayLine: 查询线路详情，入参 ZoneId/GatewayId/LineId，出参 Line（MultiPathGatewayLine 类型，含 LineId/LineType/LineAddress/ProxyId/RuleId）
- ModifyMultiPathGatewayLine: 修改线路，入参 ZoneId/GatewayId/LineId/LineType/LineAddress/ProxyId/RuleId
- DeleteMultiPathGatewayLine: 删除线路，入参 ZoneId/GatewayId/LineId

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_multi_path_gateway_line` 资源，支持多通道安全加速网关线路的完整 CRUD 生命周期管理
- 使用复合 ID（zone_id + gateway_id + line_id）标识资源实例
- 遵循现有 TEO 资源的代码风格和架构模式（参考 `tencentcloud_igtm_strategy`）
- 在 Read/Update/Delete 方法中从 `d.Get()` 获取 ID 各字段，而非从 `d.Id()` 解析
- 支持资源导入（Import）
- 提供单元测试覆盖核心逻辑

**Non-Goals:**
- 不新增数据源（Datasource）
- 不修改现有 TEO 资源的 schema 或行为
- 不处理异步接口轮询（本次涉及的四个接口均为同步接口）

## Decisions

### 1. 复合 ID 格式
**决策**: 使用 `zone_id + tccommon.FIELD_SP + gateway_id + tccommon.FIELD_SP + line_id` 作为资源 ID

**理由**: 线路资源需要 ZoneId、GatewayId 和 LineId 三个字段唯一确定，这与 TEO 其他资源（如 L4Proxy 规则等）的复合 ID 模式一致。在 Read/Update/Delete 方法中，通过 `d.Get()` 分别获取各字段值，而非解析 `d.Id()` 字符串。

### 2. Schema 字段设计
**决策**:
- `zone_id`: Required, ForceNew — 站点 ID，创建后不可变
- `gateway_id`: Required, ForceNew — 网关 ID，创建后不可变
- `line_type`: Required — 线路类型（direct/proxy/custom）
- `line_address`: Required — 线路地址，格式为 ip:port
- `proxy_id`: Optional — 四层代理实例 ID，当 LineType 为 proxy 时必传
- `rule_id`: Optional — 转发规则 ID，当 LineType 为 proxy 时必传
- `line_id`: Computed — 线路 ID，创建后由云 API 返回

**理由**: zone_id 和 gateway_id 在 Create/Delete/Describe/Modify 接口中均作为路径参数，且不支持修改，因此设为 ForceNew。line_id 由云 API 创建时返回，设为 Computed。proxy_id 和 rule_id 仅在 LineType 为 proxy 时需要，设为 Optional。

### 3. Update 方法实现
**决策**: 在 Update 方法中调用 ModifyMultiPathGatewayLine 接口，传入所有可修改字段（line_type, line_address, proxy_id, rule_id）

**理由**: ModifyMultiPathGatewayLine 接口支持修改 LineType、LineAddress、ProxyId、RuleId 四个字段，采用全量更新的方式，确保 Terraform state 与云端一致。

### 4. 错误处理与重试
**决策**: 在 Read/Delete 操作中使用 `helper.Retry()` + `tccommon.ReadRetryTimeout` 进行最终一致性重试

**理由**: 遵循现有 TEO 资源的重试模式，Create 操作后可能需要短暂等待资源生效，Read 操作需要处理资源尚未就绪的情况。

### 5. 单元测试策略
**决策**: 使用 gomonkey mock 云 API 调用，编写单元测试覆盖 CRUD 逻辑

**理由**: 根据 Go 代码生成要求，新增资源使用 mock（gomonkey）方式进行业务逻辑的单元测试，不使用 Terraform 测试套件。

## Risks / Trade-offs

- **[线路类型约束]** → 直连线路（direct）不支持创建、编辑和删除；proxy 类型线路不支持删除。这些约束由云 API 侧控制，Terraform 侧不做额外校验，依赖云 API 返回错误信息。
- **[复合 ID 解析]** → 使用 `d.Get()` 获取 ID 各字段而非解析 `d.Id()`，避免了分隔符冲突的风险，但需要确保 Import 时正确设置 state 中的各字段。
- **[ProxyId/RuleId 条件必传]** → 当 LineType 为 proxy 时，ProxyId 和 RuleId 为必传参数。Terraform schema 中将它们设为 Optional，通过 cloud API 侧校验来确保参数完整性。
