## Context

Tencent Cloud Config (`config`, API version `v20220802`) provides account groups (账号组) that
aggregate compliance data from multiple member accounts. Each account group has a delivery setting
(投递设置) that controls how Config audit logs are shipped to COS or CLS.

This change adds a Terraform **RESOURCE_KIND_CONFIG** resource
`tencentcloud_config_update_aggregate_config_deliver`. RESOURCE_KIND_CONFIG semantics: the setting
always exists once configured; Terraform only **reads (R)** and **updates (U)** it. There is no
create or delete cloud API — Create performs the first Update, and Delete is intentionally omitted
from the cloud lifecycle (no-op).

The existing `tencentcloud_config_deliver_config` resource (the account-level singleton, files
`resource_tc_config_deliver_config.go`) is the closest analog. The account-group variant differs by
requiring `account_group_id` and by exposing the extra `deliver_uin` member-account field. This
design reuses that resource's structure verbatim and adapts the API calls.

## Goals / Non-Goals

**Goals:**
- Provide `tencentcloud_config_update_aggregate_config_deliver` with Read + Update (Create = first
  Update, Delete = no-op).
- Register the resource in the provider and generate its docs.
- Keep the implementation style aligned with `resource_tc_config_deliver_config.go` and
  `resource_tc_igtm_strategy.go`.

**Non-Goals:**
- No create/delete cloud API support (none exists).
- No changes to the existing `tencentcloud_config_deliver_config` resource.
- No data source for this setting.

## Cloud API Reference (from vendor)

Package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`

### DescribeAggregateConfigDeliver (Read)

`vendor/.../config/v20220802/models.go` lines 1873–1943:

```go
// Request (line 1878)
type DescribeAggregateConfigDeliverRequest struct {
	*tchttp.BaseRequest

	// 账号组ID
	AccountGroupId *string `json:"AccountGroupId,omitnil,omitempty" name:"AccountGroupId"`
}

// ResponseParams (line 1905)
type DescribeAggregateConfigDeliverResponseParams struct {
	// 投递名称
	DeliverName *string `json:"DeliverName,omitnil,omitempty" name:"DeliverName"`
	// 资源六段式
	TargetArn *string `json:"TargetArn,omitnil,omitempty" name:"TargetArn"`
	// 投递状态 DeliverStatus：0 关闭  1 开启
	Status *uint64 `json:"Status,omitnil,omitempty" name:"Status"`
	// 创建时间
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`
	// 日志前缀
	DeliverPrefix *string `json:"DeliverPrefix,omitnil,omitempty" name:"DeliverPrefix"`
	// 投递类型
	DeliverType *string `json:"DeliverType,omitnil,omitempty" name:"DeliverType"`
	// 支持跨账号投递的成员账号uin，只能是委派管理员。默认为0，即投递到管理员账号下
	DeliverUin *int64 `json:"DeliverUin,omitnil,omitempty" name:"DeliverUin"`
	// 1：配置变更 2： 资源列表 3：全部
	DeliverContentType *uint64 `json:"DeliverContentType,omitnil,omitempty" name:"DeliverContentType"`
	// 唯一请求 ID
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

// Response (line 1940)
type DescribeAggregateConfigDeliverResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAggregateConfigDeliverResponseParams `json:"Response"`
}
```

Client method: `DescribeAggregateConfigDeliverWithContext(ctx, request)`
(`client.go` line 1089). It requires credential and returns
`*DescribeAggregateConfigDeliverResponse`.

### UpdateAggregateConfigDeliver (Update / Create)

`vendor/.../config/v20220802/models.go` lines 5200–5306:

```go
// Request (line 5229)
type UpdateAggregateConfigDeliverRequest struct {
	*tchttp.BaseRequest

	// 0 关闭  1 开启
	Status *uint64 `json:"Status,omitnil,omitempty" name:"Status"`
	// 账号组ID
	AccountGroupId *string `json:"AccountGroupId,omitnil,omitempty" name:"AccountGroupId"`
	// 投递服务名称
	DeliverName *string `json:"DeliverName,omitnil,omitempty" name:"DeliverName"`
	// 资源六段式
	// COS：qcs::cos:$region:$account:prefix/$appid/$BucketName
	// CLS: qcs::cls:$region:$account:cls/topicId
	TargetArn *string `json:"TargetArn,omitnil,omitempty" name:"TargetArn"`
	// 资源前缀
	DeliverPrefix *string `json:"DeliverPrefix,omitnil,omitempty" name:"DeliverPrefix"`
	// 投递类型  COS CLS
	DeliverType *string `json:"DeliverType,omitnil,omitempty" name:"DeliverType"`
	// 支持跨账号投递的成员账号uin，只能是委派管理员。默认为0，即投递到管理员账号下
	DeliverUin *int64 `json:"DeliverUin,omitnil,omitempty" name:"DeliverUin"`
	// 1：配置变更 2： 资源列表 3：全选
	DeliverContentType *uint64 `json:"DeliverContentType,omitnil,omitempty" name:"DeliverContentType"`
}

// ResponseParams (line 5287) — only RequestId
type UpdateAggregateConfigDeliverResponseParams struct {
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

// Response (line 5292)
type UpdateAggregateConfigDeliverResponse struct {
	*tchttp.BaseResponse
	Response *UpdateAggregateConfigDeliverResponseParams `json:"Response"`
}
```

Client method: `UpdateAggregateConfigDeliverWithContext(ctx, request)` (`client.go` line 3163).

Notes:
- `Status` in the Update request is documented as required by the requirement mapping, but the Go
  struct marks it `omitempty`; in the provider it is `Required: true` in schema.
- Both APIs are **synchronous** (no FlowId/TaskId is returned), so no async polling is required.
  After Update we call Read to refresh state — the single retry on Update is sufficient.

## Decisions

- **Resource file name** → `resource_tc_config_update_aggregate_config_deliver.go`
  (RESOURCE_KIND_CONFIG naming: `resource_tc_<Product>_<资源名>_config.go`).
- **Resource ID = `account_group_id`.** The account-group delivery setting is uniquely identified by
  its account group. Using `account_group_id` as the id makes the Read/Update call natural
  (pass `d.Id()` as `AccountGroupId`) and keeps id stable. `account_group_id` is `ForceNew`.
  (Alternative: `helper.BuildToken()` singleton id as in `config_deliver_config`; rejected because
  the account group id is a natural, stable, importable key.)
- **Create = first Update.** `resourceTencentCloudConfigUpdateAggregateConfigDeliverCreate` sets
  `d.SetId(account_group_id)` then delegates to Update. Since the deliver setting is created
  implicitly on the first Update call, we set the id before calling Update.
- **Delete = no-op.** There is no delete API; the setting persists. Return nil.
- **Read nil handling.** If `response.Response == nil`, log and `d.SetId("")` (setting no longer
  exists) per provider rules; if only individual fields are nil, skip setting those fields.
- **Schema types** map the SDK pointer types:
  - `Status` (`*uint64`) → `schema.TypeInt` (Required)
  - `DeliverContentType` (`*uint64`) → `schema.TypeInt` (Optional)
  - `DeliverUin` (`*int64`) → `schema.TypeInt` (Optional)
  - `AccountGroupId`, `DeliverName`, `TargetArn`, `DeliverPrefix`, `DeliverType`, `CreateTime` → String

### Schema table

| Field | Type | Mode | API field |
|---|---|---|---|
| `account_group_id` | String | Required, ForceNew | Describe: req.AccountGroupId / Update: req.AccountGroupId |
| `status` | Int | Required | Describe: resp.Status / Update: req.Status |
| `deliver_name` | String | Optional | Describe: resp.DeliverName / Update: req.DeliverName |
| `target_arn` | String | Optional | Describe: resp.TargetArn / Update: req.TargetArn |
| `deliver_prefix` | String | Optional | Describe: resp.DeliverPrefix / Update: req.DeliverPrefix |
| `deliver_type` | String | Optional | Describe: resp.DeliverType / Update: req.DeliverType |
| `deliver_uin` | Int | Optional | Describe: resp.DeliverUin / Update: req.DeliverUin |
| `deliver_content_type` | Int | Optional | Describe: resp.DeliverContentType / Update: req.DeliverContentType |
| `create_time` | String | Computed | Describe: resp.CreateTime |

### Service layer (append to `service_tencentcloud_config.go`)

```go
func (me *ConfigService) DescribeAggregateConfigDeliver(ctx context.Context, accountGroupId string) (ret *configv20220802.DescribeAggregateConfigDeliverResponseParams, errRet error)
```
- Builds `NewDescribeAggregateConfigDeliverRequest()`, sets `AccountGroupId`.
- Wraps the call in `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `ratelimit.Check`.
- Returns `NonRetryableError` when `result == nil || result.Response == nil`.
- Model after existing `DescribeConfigDeliver` (service file lines 579–615).

### Code-style reference: `resource_tc_config_deliver_config.go`

The new resource is a near-copy with these deltas:
- Add `account_group_id` (Required, ForceNew) and `deliver_uin` (Optional Int) to schema.
- CRUD funcs read `account_group_id` from schema / `d.Id()`.
- Create sets `d.SetId(account_group_id value)` then calls Update.
- Update builds `NewUpdateAggregateConfigDeliverRequest()`, sets all fields including `AccountGroupId`.
- Read calls `service.DescribeAggregateConfigDeliver(ctx, d.Id())`.
- Importer: `schema.ImportStatePassthrough` (id = account_group_id).

## Risks / Trade-offs

- [Account-group delivery setting may not exist yet] → Create calls Update first, which creates it;
  Read is only called after Update.
- [Setting removed out-of-band] → Read logs and clears id (`d.SetId("")`) so Terraform plans a
  recreate on next apply.
- [`account_group_id` change] → marked ForceNew; changing it destroys+recreates the tf resource
  (cloud setting for the old group is untouched, matching no-op delete semantics).
- [Delete cannot remove the cloud setting] → documented as RESOURCE_KIND_CONFIG; Delete is no-op.
  Users can set `status = 0` before removing the resource if they want to disable delivery.

## Migration Plan

Purely additive. No migration needed. New resource becomes available after provider build.

## Open Questions

None — API surface and schema fully determined from vendor definitions and the requirement mapping.
