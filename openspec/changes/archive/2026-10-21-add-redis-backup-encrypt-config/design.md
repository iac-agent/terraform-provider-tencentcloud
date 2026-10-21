# Design: Redis Backup Encrypt Config Resource

## Context

TencentDB for Redis (CRS) already has several config-style resources in this provider
(`resource_tc_redis_backup_config.go`, `resource_tc_redis_connection_config.go`), all written in
SDKv2 (`helper/schema`). The new `tencentcloud_redis_backup_encrypt` resource must instead be
written with **Terraform Plugin Framework**, following the `RESOURCE_KIND_CONFIG` template.

Cloud side, only two APIs exist for this feature:

| Purpose | API | Notes |
| --- | --- | --- |
| Read config | `DescribeBackupEncryptConfig` | input: `InstanceId` only |
| Set config | `ModifyBackupEncrypt` | input: `InstanceId`, `BackupEncrypt`, `EncryptKeyType`, `KmsKeyId`, `KmsRegion` |

There is no dedicated Create/Delete API — this is exactly the config pattern:
Create = first `ModifyBackupEncrypt`, Update = `ModifyBackupEncrypt` again, Read =
`DescribeBackupEncryptConfig`, Delete = **disable encryption** (enable-switch semantic).

Constraints:
- The provider muxes SDKv2 and Framework providers, so the Framework resource schema must be
  self-contained (new type name, no SDKv2 counterpart with the same name).
- The existing `crs` package holds SDKv2 code; new Framework resources live in the **same
  package** (`package crs`) per `tencentcloud/framework/README.md`.
- Business code must obtain the cloud client through `*sharedmeta.ProviderMeta` inside `Configure`
  (see `tencentcloud/services/ssm/resource_tc_ssm_secret_version_v2.go`).

## Goals / Non-Goals

**Goals:**
- Declarative management of Redis backup encryption (switch + key type + KMS key ID + KMS region).
- Idempotent Read that reflects cloud state, so drift (manual console change) is detected.
- Safe destroy: never delete the host Redis instance; disable backup encryption instead.
- Import support by host instance ID.

**Non-Goals:**
- Creating/destroying Redis instances (owned by `tencentcloud_redis_instance`).
- Managing KMS keys themselves (owned by the KMS product).
- Any change to the existing SDKv2 `tencentcloud_redis_backup_config` resource.
- Automatic retry/polling of asynchronous task state — `ModifyBackupEncrypt` is **synchronous**
  (response carries only `RequestId`, no `FlowId`/`TaskId`), so no poll loop is required; Read
  immediately after Modify is still performed to verify the cloud state matches the plan.

## Decisions

### D1. Resource kind mapping: `RESOURCE_KIND_CONFIG`

`tencentcloud_redis_backup_encrypt` manages a *setting* of an existing instance, therefore:
- `Metadata.TypeName = "tencentcloud_redis_backup_encrypt"`
- State `id` = host `instance_id`
- `instance_id` is `Required` + `RequiresReplace()`

Alternative considered: modelling it as a `RESOURCE_KIND_GENERAL` resource — rejected, because the
cloud has no create/delete API for the config and the instance lifecycle belongs to another resource.

### D2. Destroy semantic: disable encryption (enable-switch type)

Per the config template's implementation note #3, backup encryption is an **enable switch**:

| Aspect | Decision |
| --- | --- |
| Resource semantic type | Enable-switch (backup encryption on/off) |
| Delete behaviour | Call `ModifyBackupEncrypt` with `BackupEncrypt=false`, then clear state |
| Rationale | Leaving `BackupEncrypt=true` after `terraform destroy` silently keeps the feature and its KMS usage alive — a silent leak. The Redis instance itself is **never** deleted. |
| `EncryptKeyType`/`KmsKeyId`/`KmsRegion` on delete | Omit (cloud keeps the last key configuration; it is inert while encryption is off) |
| Documentation | Schema `Description` and `.md` must state "destroying this resource disables backup encryption on the cloud side; it never deletes the instance" |

Because `backup_encrypt` is `Required`, the plan value is always known; no default is invented.

### D3. Schema design

| Attribute | Type | Framework schema | Notes |
| --- | --- | --- | --- |
| `id` | `types.String` | `Computed` + `UseStateForUnknown` | equals `instance_id` |
| `instance_id` | `types.String` | `Required` + `RequiresReplace` | host instance ID; changing it recreates the resource |
| `backup_encrypt` | `types.Bool` | `Required` | whether backup encryption is enabled |
| `encrypt_key_type` | `types.Int64` | `Optional` + `Computed`, validator `OneOf(0, 1)` | `0` = Tencent Cloud managed key (cloud default), `1` = user-defined KMS key |
| `kms_key_id` | `types.String` | `Optional` + `Computed` | KMS key ID, only meaningful when `encrypt_key_type = 1` |
| `kms_region` | `types.String` | `Optional` + `Computed` | KMS region, only meaningful when `encrypt_key_type = 1` |

Notes:
- `encrypt_key_type` / `kms_key_id` / `kms_region` use `Optional + Computed` because the cloud has
  defaults (`EncryptKeyType` default `0`) and back-fills the actual values; do **not** hardcode
  provider-side `Default`.
- Request-side reads must guard with `!IsNull() && !IsUnknown()`; response-side writes use
  `types.XxxPointerValue(...)` so a nil SDK pointer becomes `Null` naturally.
- No `List`/`Set`/`Nested` blocks are needed — all fields are scalars.

### D4. Lifecycle callbacks

| Callback | Behaviour |
| --- | --- |
| `Create` | `modifyBackupEncrypt(plan)` → on success set `plan.Id = plan.InstanceId` → `resp.State.Set` → then `Read`-equivalent refresh so state matches cloud |
| `Read` | If `instance_id` empty (post-import) backfill from `id`; call `DescribeBackupEncryptConfig` with `tccommon.ReadRetryTimeout` retry; on `ResourceNotFound.InstanceNotExists` → `log.Printf` + `resp.State.RemoveResource(ctx)`; if `response.Response == nil` → `log.Printf("[CRUD] redis_backup_encrypt id=%s")` + `AddError`; otherwise write back the four fields |
| `Update` | `modifyBackupEncrypt(plan)` with the same helper used by Create → `resp.State.Set(&plan)` → refresh |
| `Delete` | Call `ModifyBackupEncrypt(InstanceId=id, BackupEncrypt=false)`; do **not** touch the instance |
| `ImportState` | `resp.State.SetAttribute(ctx, path.Root("id"), req.ID)`; `instance_id` is backfilled in `Read` — import ID is the plain instance ID (no `FILED_SP` needed because there is exactly one config per instance) |

Both Create and Update share one private `modifyBackupEncrypt` helper — the signature trait of the
config pattern.

### D5. Error handling and retry

- All cloud calls are wrapped in `resource.Retry(tccommon.WriteRetryTimeout | tccommon.ReadRetryTimeout, func() *resource.RetryError { ... })`
  with `tccommon.RetryError(err)` inside — consistent with the rest of the `crs` package.
- Only the API call lives inside the retry block; `plan.Id` assignment and `State.Set` happen
  **after** the retry block, as required by the provider rules.
- All errors are reported through `resp.Diagnostics.AddError(...)`; never `panic`.
- `ModifyBackupEncrypt` is synchronous, so no async polling loop is needed. (If the cloud later
  returns a task identifier, the flow must be: set `Id` first, then poll `Describe...` — noted here
  for future changes only.)

### D6. Registration and docs

1. `tencentcloud/framework/registry.go`
   - add import `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"`
   - add `crs.NewCrsRedisBackupEncryptConfigResource,` to `resourceFactories`
   (`tencentcloud/provider.go` is **not** touched — Framework resources are registered only here.)
2. `tencentcloud/provider.md` → under `TencentDB for Redis(crs)` → `Resource`, append
   `tencentcloud_redis_backup_encrypt`.
3. `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_config.md` — one-line description
   (starts with "Provides a resource to ..."), `Example Usage` (with `jsonencode()` if any JSON
   string value appears — none here), `Import` section (config resources do have Import).

## Reference: Cloud API definitions (copied from vendor)

Package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412`

### `DescribeBackupEncryptConfig`

Client methods (`vendor/.../redis/v20180412/client.go` L1775-L1815):

```go
func NewDescribeBackupEncryptConfigRequest() *DescribeBackupEncryptConfigRequest
func (c *Client) DescribeBackupEncryptConfig(request *DescribeBackupEncryptConfigRequest) (*DescribeBackupEncryptConfigResponse, error)
func (c *Client) DescribeBackupEncryptConfigWithContext(ctx context.Context, request *DescribeBackupEncryptConfigRequest) (*DescribeBackupEncryptConfigResponse, error)
```

Request (`models.go` L2340-L2351):

```go
type DescribeBackupEncryptConfigRequestParams struct {
	// Instance ID.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}

type DescribeBackupEncryptConfigRequest struct {
	*tchttp.BaseRequest
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}
```

Response (`models.go` L2372-L2393):

```go
type DescribeBackupEncryptConfigResponseParams struct {
	// Whether backups are encrypted. true: encrypted, false: not encrypted.
	BackupEncrypt *bool `json:"BackupEncrypt,omitnil,omitempty" name:"BackupEncrypt"`
	// Encryption key type. 0: Tencent Cloud managed (default); 1: user-defined KMS key.
	EncryptKeyType *int64 `json:"EncryptKeyType,omitnil,omitempty" name:"EncryptKeyType"`
	// KMS key ID, format xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx.
	KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
	// KMS region, format ap-xxxx.
	KmsRegion *string `json:"KmsRegion,omitnil,omitempty" name:"KmsRegion"`
	// Unique request ID.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeBackupEncryptConfigResponse struct {
	*tchttp.BaseResponse
	Response *DescribeBackupEncryptConfigResponseParams `json:"Response"`
}
```

### `ModifyBackupEncrypt`

Client methods (`client.go` L5587-L5625):

```go
func NewModifyBackupEncryptRequest() *ModifyBackupEncryptRequest
func (c *Client) ModifyBackupEncrypt(request *ModifyBackupEncryptRequest) (*ModifyBackupEncryptResponse, error)
func (c *Client) ModifyBackupEncryptWithContext(ctx context.Context, request *ModifyBackupEncryptRequest) (*ModifyBackupEncryptResponse, error)
```

Request (`models.go` L8697-L8732):

```go
type ModifyBackupEncryptRequestParams struct {
	// Instance ID.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
	// Whether to enable backup encryption. Default: false.
	BackupEncrypt *bool `json:"BackupEncrypt,omitnil,omitempty" name:"BackupEncrypt"`
	// Encryption key type. 0: Tencent Cloud managed (default); 1: user-defined KMS key. Default: 0.
	EncryptKeyType *int64 `json:"EncryptKeyType,omitnil,omitempty" name:"EncryptKeyType"`
	// KMS key ID, format xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx.
	KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
	// KMS region, format ap-xxxx. See the KMS region list.
	KmsRegion *string `json:"KmsRegion,omitnil,omitempty" name:"KmsRegion"`
}

type ModifyBackupEncryptRequest struct {
	*tchttp.BaseRequest
	// same five fields as above
}
```

Response (`models.go` L8757-L8766) — **no task/flow identifier**, i.e. synchronous:

```go
type ModifyBackupEncryptResponseParams struct {
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyBackupEncryptResponse struct {
	*tchttp.BaseResponse
	Response *ModifyBackupEncryptResponseParams `json:"Response"`
}
```

## Reference: Framework template structure (RESOURCE_KIND_CONFIG)

Source: `~/.codebuddy/skills/tfpacer-tf-templates/references/3-resource-config.md` and
`references/6-field-types-and-computed.md`.

### Skeleton (English identifiers as shipped)

```go
package crs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var (
	_ resource.Resource                = &CrsRedisBackupEncryptConfigResource{}
	_ resource.ResourceWithConfigure   = &CrsRedisBackupEncryptConfigResource{}
	_ resource.ResourceWithImportState = &CrsRedisBackupEncryptConfigResource{}
)

func NewCrsRedisBackupEncryptConfigResource() resource.Resource {
	return &CrsRedisBackupEncryptConfigResource{}
}

type CrsRedisBackupEncryptConfigResource struct {
	client *connectivity.TencentCloudClient // or RedisService, mirroring ssm reference
}

type CrsRedisBackupEncryptConfigResourceModel struct {
	Id             types.String `tfsdk:"id"` // = host instance ID
	InstanceId     types.String `tfsdk:"instance_id"`
	BackupEncrypt  types.Bool   `tfsdk:"backup_encrypt"`
	EncryptKeyType types.Int64  `tfsdk:"encrypt_key_type"`
	KmsKeyId       types.String `tfsdk:"kms_key_id"`
	KmsRegion      types.String `tfsdk:"kms_region"`
}

func (r *CrsRedisBackupEncryptConfigResource) Metadata(...) // TypeName = "tencentcloud_redis_backup_encrypt"
func (r *CrsRedisBackupEncryptConfigResource) Schema(...)   // see D3
func (r *CrsRedisBackupEncryptConfigResource) Configure(...) // *sharedmeta.ProviderMeta -> client
func (r *CrsRedisBackupEncryptConfigResource) Create(...)    // modifyBackupEncrypt -> plan.Id = plan.InstanceId -> State.Set -> refresh
func (r *CrsRedisBackupEncryptConfigResource) Read(...)      // DescribeBackupEncryptConfig (+ retry, not-found handling)
func (r *CrsRedisBackupEncryptConfigResource) Update(...)    // modifyBackupEncrypt -> State.Set -> refresh
func (r *CrsRedisBackupEncryptConfigResource) Delete(...)    // ModifyBackupEncrypt(BackupEncrypt=false)
func (r *CrsRedisBackupEncryptConfigResource) ImportState(...) // SetAttribute(id, req.ID)
func (r *CrsRedisBackupEncryptConfigResource) modifyBackupEncrypt(ctx, m *Model) error // shared by Create/Update
```

### Template key points checklist (must be satisfied)

1. Create and Update **share** one private `modifyBackupEncrypt` helper — both call `ModifyBackupEncrypt`.
2. State ID = host resource ID; one config instance per host.
3. Delete semantic decided and documented: **enable-switch → disable** (see D2). `README`/`.md`
   and Schema `Description` must say "destroy disables backup encryption; the instance is never deleted".
4. Host ID uses `RequiresReplace()`; when `Describe` returns `ResourceNotFound.InstanceNotExists`,
   `log.Printf` the id and call `resp.State.RemoveResource(ctx)`.
5. `Optional + Computed` fields judged with `IsNull()`/`IsUnknown()` on the request side; response
   side uses `types.XxxPointerValue(...)`.
6. Naming: file `resource_tc_redis_backup_encrypt_config.go`, TF name
   `tencentcloud_redis_backup_encrypt`, factory `NewCrsRedisBackupEncryptConfigResource()`,
   registered in `resourceFactories`.
7. From the provider rules (in addition to the template): only the API call lives inside
   `resource.Retry`; `d.SetId`/`State.Set` after the retry block; every returned error checked;
   errors surfaced through `resp.Diagnostics.AddError`.

## Implementation plan sketch

```
tencentcloud/services/crs/
├── resource_tc_redis_backup_encrypt_config.go        # new Framework resource (package crs)
├── resource_tc_redis_backup_encrypt_config_test.go   # gomonkey unit tests, no TF acceptance suite
└── resource_tc_redis_backup_encrypt_config.md        # description + Example Usage + Import

tencentcloud/framework/registry.go                    # +1 import (crs), +1 factory entry
tencentcloud/provider.md                              # +1 line under TencentDB for Redis(crs) / Resource
```

### Service layer

Either call the SDK client directly inside the resource (as the Framework ssm reference does) or
add two thin helpers to `RedisService` in `service_tencentcloud_redis.go`
(`DescribeBackupEncryptConfig(ctx, instanceId)`, `ModifyBackupEncrypt(ctx, params)`). Preferred:
add the helpers so the resource only orchestrates state and retry, matching existing `crs` style.

## Risks / Trade-offs

| Risk | Mitigation |
| --- | --- |
| Disabling encryption on destroy may surprise users who adopted the resource only to read/keep the setting | Documented explicitly in Schema `Description`, `.md` and `provider.md` entry; the behaviour follows the provider-wide enable-switch rule |
| `EncryptKeyType` semantics when `BackupEncrypt=false` are not meaningful; cloud may return stale values | Fields are `Optional + Computed`, so Read always mirrors the cloud and no spurious diff is produced |
| Disabling encryption with in-flight encrypted backups | Out of provider scope; rely on cloud-side validation errors, surfaced verbatim via `AddError` |
| `kms_key_id` / `kms_region` only valid with `encrypt_key_type = 1` | Documented; cloud API rejects invalid combinations and the error is returned to the user |
| Framework + SDKv2 mux requires no name collision | `tencentcloud_redis_backup_encrypt` is a brand-new type name; `framework/provider_test.go` mux sanity test covers collisions |
| Unit tests cannot hit real cloud APIs | gomonkey mocks `UseRedisClient()` / `DescribeBackupEncryptConfig` / `ModifyBackupEncrypt`; tests assert request field mapping and state transitions only |

## Migration Plan

Purely additive; no state migration. Existing users are unaffected until they add
`tencentcloud_redis_backup_encrypt` to their configuration. Rollback = delete the resource block
(note: destroy disables encryption, so re-enable it in the console or via `backup_encrypt = true`
before removing the block if encryption must be preserved).

## Open Questions

None blocking. Two items were confirmed from the vendor SDK:
1. `ModifyBackupEncrypt` returns no task identifier → synchronous, no polling needed.
2. `DescribeBackupEncryptConfig` returns exactly the four mapped fields → 1:1 mapping with schema.
