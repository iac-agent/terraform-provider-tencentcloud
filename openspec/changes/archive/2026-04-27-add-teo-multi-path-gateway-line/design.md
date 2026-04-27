## Context

Terraform Provider for TencentCloud 当前支持 TEO (EdgeOne) 产品线的多种资源管理，但缺少对多通道安全加速网关线路（MultiPathGatewayLine）的管理支持。该线路资源属于 TEO 多通道安全加速网关的子资源，用于配置网关下的加速线路。

云 API 已在 vendor 中提供完整的 CRUD 接口：
- `CreateMultiPathGatewayLine` - 创建线路
- `DescribeMultiPathGatewayLine` - 查询线路详情
- `ModifyMultiPathGatewayLine` - 修改线路
- `DeleteMultiPathGatewayLine` - 删除线路

所有接口均为同步接口，不需要异步轮询。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_multi_path_gateway_line` 资源，支持完整的 CRUD 生命周期管理
- 使用复合 ID 格式 `{zone_id}#{gateway_id}#{line_id}` 唯一标识资源实例
- 支持 Terraform Import 导入已有线路资源
- 遵循项目现有的代码风格和模式（参考 tencentcloud_igtm_strategy）
- 提供单元测试（使用 gomonkey mock 方式）
- 生成对应的 .md 示例文件

**Non-Goals:**
- 不实现多通道安全加速网关（MultiPathGateway）本身的管理（属于其他资源）
- 不实现密钥管理（SecretKey）相关功能
- 不实现异步轮询逻辑（所有接口均为同步）

## Decisions

### 1. 复合 ID 格式：`{zone_id}#{gateway_id}#{line_id}`
**选择**: 使用三段式复合 ID，以 `#` 分隔。
**理由**: TEO 多通道安全加速网关线路的 API 操作需要三个标识参数（ZoneId、GatewayId、LineId），且 Read/Update/Delete 接口均需要这三个参数。采用复合 ID 可避免额外存储，Read 时直接从 ID 解析即可。
**替代方案**: 将 gateway_id 和 line_id 存储为 Computed 字段 — 增加复杂度，不如直接从 ID 解析简洁。

### 2. Schema 字段设计
**选择**:
- `zone_id` - Required, ForceNew（站点 ID，创建后不可变）
- `gateway_id` - Required, ForceNew（网关 ID，创建后不可变）
- `line_type` - Required（线路类型：direct/proxy/custom）
- `line_address` - Required（线路地址，格式 ip:port）
- `proxy_id` - Optional（四层代理实例 ID，proxy 类型时必传）
- `rule_id` - Optional（转发规则 ID，proxy 类型时必传）
- `line_id` - Computed（线路 ID，由创建接口返回）

**理由**: 根据 CreateMultiPathGatewayLine 和 ModifyMultiPathGatewayLine 的入参设计。ZoneId 和 GatewayId 在所有接口中一致，创建后不可修改，设为 ForceNew。LineType 和 LineAddress 在 Modify 接口中支持修改，设为可更新字段。ProxyId 和 RuleId 仅在 LineType 为 proxy 时需要，设为 Optional。LineId 由创建接口返回，设为 Computed。

### 3. 代码文件组织
**选择**: 资源代码放在 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line.go`。
**理由**: 遵循项目现有的按服务拆分目录的组织模式，teo 服务已有独立目录。

### 4. 单元测试策略
**选择**: 使用 gomonkey mock 云 API 进行单元测试。
**理由**: 按照项目要求，新增资源使用 mock 方式而非 Terraform 测试套件。

## Risks / Trade-offs

- **[API 参数校验]** proxy_id 和 rule_id 仅在 line_type 为 proxy 时必传，但 Terraform schema 级别无法动态校验必填 → 在 Create/Update 函数中添加逻辑校验，当 line_type 为 proxy 时检查 proxy_id 和 rule_id 是否已设置
- **[Read 接口返回结构]** DescribeMultiPathGatewayLine 返回的 Line 字段是 MultiPathGatewayLine 对象而非列表 → Read 时直接从 Line 对象中提取字段，无需遍历
- **[Update 接口]** ModifyMultiPathGatewayLine 不支持修改 line_type 为 direct 类型的 line_address → 遵循云 API 行为，由 API 侧返回错误
