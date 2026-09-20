## Context

当前 `tencentcloud_teo_l4_proxy` 资源（`resource_tc_teo_l4_proxy.go`）是一个 RESOURCE_KIND_GENERAL 资源，管理 L4 代理实例的完整生命周期（CRUD）。资源的 Read 方法通过 `TeoService.DescribeTeoL4ProxyById()` 调用云 API `DescribeL4Proxy` 获取资源信息。

现有代码在调用 `DescribeL4Proxy` 时仅传递了 `ZoneId` 和 `Filters`（按 proxy-id 过滤），未暴露分页参数 `Offset`/`Limit`，也未将响应的 `TotalCount` 写入 state。

### 涉及的文件

| 文件 | 作用 |
|------|------|
| `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go` | 资源 schema 定义及 CRUD 逻辑 |
| `tencentcloud/services/teo/service_tencentcloud_teo.go` | 服务层，封装 DescribeL4Proxy 调用 |
| `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md` | 资源文档 |
| `tencentcloud/services/teo/resource_tc_teo_l4_proxy_test.go` | 资源单元测试 |

### 现有 Schema 结构（resource_tc_teo_l4_proxy.go 第 27-101 行）

当前 schema 包含以下字段：
- `zone_id` (Required, ForceNew, TypeString) — 站点 ID
- `proxy_id` (Computed, TypeString) — L4 代理实例 ID
- `proxy_name` (Required, TypeString) — 实例名称
- `area` (Optional, TypeString) — 加速区域
- `ipv6` (Optional, TypeString) — 是否开启 IPv6
- `static_ip` (Optional, TypeString) — 是否开启固定 IP
- `accelerate_mainland` (Optional, TypeString) — 是否开启中国大陆网络优化
- `ddos_protection_config` (Optional/Computed, TypeList, Deprecated) — DDoS 保护配置

### 云 API 定义（vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go）

**DescribeL4ProxyRequestParams（行 11481-11496）：**
```go
type DescribeL4ProxyRequestParams struct {
    ZoneId  *string   `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
    Offset  *uint64   `json:"Offset,omitnil,omitempty" name:"Offset"`
    Limit   *uint64   `json:"Limit,omitnil,omitempty" name:"Limit"`
    Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeL4ProxyRequest struct {
    *tchttp.BaseRequest
    ZoneId  *string   `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`
    Offset  *uint64   `json:"Offset,omitnil,omitempty" name:"Offset"`
    Limit   *uint64   `json:"Limit,omitnil,omitempty" name:"Limit"`
    Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}
```

**DescribeL4ProxyResponseParams（行 11540-11549）：**
```go
type DescribeL4ProxyResponseParams struct {
    TotalCount *uint64   `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
    L4Proxies  []*L4Proxy `json:"L4Proxies,omitnil,omitempty" name:"L4Proxies"`
    RequestId  *string   `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeL4ProxyResponse struct {
    *tchttp.BaseResponse
    Response *DescribeL4ProxyResponseParams `json:"Response"`
}
```

**Filter struct（行 16822+）：**
```go
type Filter struct {
    Name   *string   `json:"Name,omitnil,omitempty" name:"Name"`
    Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
}
```

### 现有服务层 DescribeTeoL4ProxyById（service_tencentcloud_teo.go 第 1249-1289 行）

```go
func (me *TeoService) DescribeTeoL4ProxyById(ctx context.Context, zoneId string, proxyId string) (ret *teo.L4Proxy, errRet error) {
    // 构建请求：只设置 ZoneId 和 Filters（proxy-id），未使用 Offset/Limit
    request := teo.NewDescribeL4ProxyRequest()
    request.ZoneId = &zoneId
    filter := &teo.Filter{
        Name:   helper.String("proxy-id"),
        Values: []*string{&proxyId},
    }
    request.Filters = append(request.Filters, filter)
    // ... retry 调用 DescribeL4Proxy ...
    // 返回 L4Proxies[0]
}
```

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_l4_proxy` 资源的 Read 路径中支持传入 `offset`/`limit` 分页参数
- 将 API 响应的 `TotalCount` 通过 `total_count` Computed 属性暴露给用户
- 所有新增参数均为 Optional 或 Computed，保持完全向后兼容

**Non-Goals:**
- 不修改 Create/Update/Delete 路径（这些接口不使用 DescribeL4Proxy）
- 不修改资源现有的所有 schema 字段定义
- 不在 Read 方法之外的函数（如 Create/Update/Delete）中使用这些新参数

## Decisions

### Decision 1: 在 Read 方法中传递 offset/limit，而不是在 Schema 中添加新参数后修改所有 CRUD

**选择**: 只在 Schema 中新增 `offset`、`limit`（Optional）、`total_count`（Computed），仅在 Read 方法中使用。

**理由**: `DescribeL4Proxy` 只是查询接口，offset/limit 参数仅在 Read 路径有意义。Create/Update/Delete 使用的是其他 API（CreateL4Proxy/ModifyL4Proxy/DeleteL4Proxy），不需要这些分页参数。

### Decision 2: 修改 DescribeTeoL4ProxyById 签名增加参数

**选择**: 将 `DescribeTeoL4ProxyById` 的签名从 `(ctx, zoneId, proxyId)` 扩展为 `(ctx, zoneId, proxyId, offset, limit)`。

**理由**: 
- 保持服务层方法职责单一，封装 API 调用细节
- 避免在 resource 文件中直接构造 API 请求，保持分层清晰
- offset/limit 为指针类型，nil 表示不设置（使用 API 默认值）

**备选方案**: 在 resource Read 方法中直接构造 DescribeL4Proxy 请求，绕过 service 层。
- 拒绝理由：破坏现有分层架构，且与现有代码风格不一致。

### Decision 3: offset/limit 使用 TypeInt + uint64 转换

**选择**: Schema 中使用 `schema.TypeInt`，在传递给 API 前通过 `helper.IntUint64()` 转换为 `*uint64`。

**理由**: 与 SDK 中 `Offset *uint64` / `Limit *uint64` 类型匹配，同时 Terraform Plugin SDK v2 对整数的处理以 TypeInt 为主，helper 包提供标准转换。

### Decision 4: total_count 仅在 Read 中设置，不参与资源 ID 计算

**选择**: `total_count` 是 Computed 属性，不在 Create/Import 路径中手动设置。

**理由**: `TotalCount` 是 DescribeL4Proxy 接口根据当前过滤条件返回的统计值，与单个资源实例无直接关系。仅在 Read 方法中通过 API 响应的 `response.Response.TotalCount` 设置即可。

## Risks / Trade-offs

- [Risk] `offset`/`limit` 是查询参数，用户可能在 `terraform plan` 时看到这些参数被标记为 "known after apply"，因为它们对 Read 行为的影响在 plan 阶段不可见。
  → Mitigation: 这些参数标记为 Optional，默认不设置，API 使用默认值（offset=0, limit=20），不影响现有用户。

- [Risk] `total_count` 的值受 offset/limit 之外的过滤条件（如 proxy-id filter）影响，可能与用户直觉不符。
  → Mitigation: 在文档中明确说明 `total_count` 反映的是 `DescribeL4Proxy` 接口在当前过滤条件下的总数，即按 zone_id + proxy-id 过滤后的结果数（通常为 1）。

- [Risk] 修改 `DescribeTeoL4ProxyById` 签名可能导致调用方遗漏更新。
  → Mitigation: 该函数当前仅有 `resourceTencentCloudTeoL4ProxyRead` 一处调用，修改后编译错误即可发现遗漏。