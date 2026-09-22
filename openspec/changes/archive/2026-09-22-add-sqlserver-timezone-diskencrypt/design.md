## Context

The `tencentcloud_sqlserver_basic_instance` resource (`RESOURCE_KIND_GENERAL`) manages the full lifecycle of a SQL Server basic instance via the TencentCloud SQL Server SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328`). Currently the resource does not expose two creation-time parameters — `TimeZone` and `DiskEncryptFlag` — even though the `CreateBasicDBInstances` cloud API already accepts them. This change adds those two optional parameters so users can control the instance system timezone and disk encryption at creation time.

The implementation is a straightforward additive change to a single resource file and its service layer. There is no data model migration and no cross-service impact. All new fields are `Optional + Computed + ForceNew`, so the change is fully backward compatible.

The key implementation reference info (SDK struct field definitions, existing schema shape) is captured below so that the implementation phase does not need to re-read the large vendored `models.go`.

## Goals / Non-Goals

**Goals:**
- Add `time_zone` (optional, computed, ForceNew) schema field mapped to `CreateBasicDBInstancesRequest.TimeZone`.
- Add `disk_encrypt_flag` (optional, computed, ForceNew) schema field mapped to `CreateBasicDBInstancesRequest.DiskEncryptFlag`.
- Populate `time_zone` on Read from the existing `DescribeDBInstances` response (`DBInstance.TimeZone`).
- Populate `disk_encrypt_flag` on Read from a separate `DescribeDBInstancesAttribute` response (`IsDiskEncryptFlag`).
- Prevent in-place modification of both parameters via `immutableArgs` (ForceNew enforced).
- Keep the change backward compatible: existing configs/state continue to work unchanged.

**Non-Goals:**
- No update support for timezone or encryption (cloud API does not support modifying these post-creation).
- No changes to other sqlserver resources or data sources.
- No website docs written by hand (generated via `make doc` in the finalize phase only).
- No `go build`/`go vet`/`golint` during code generation (handled by other process steps).

## Decisions

### Decision 1: Both parameters are ForceNew + added to immutableArgs

**Choice:** `time_zone` and `disk_encrypt_flag` are `Optional + Computed + ForceNew`, and both are added to the `immutableArgs` slice in the Update function.

**Rationale:** The `CreateBasicDBInstances` API accepts these only on creation; there is no Modify/Update cloud API that can change timezone or disk encryption for an existing instance. Therefore any change must force recreation. Adding them to `immutableArgs` makes the Update function return an explicit error if a user attempts to change them (consistent with how `collation` is already handled), rather than silently ignoring the change.

**Alternatives considered:**
- Make them required: rejected — breaks backward compatibility for existing users.
- Omit `immutableArgs` entry and rely only on `ForceNew`: rejected — current resource convention uses `immutableArgs` to surface a clear error for ForceNew fields in the Update path, so the new fields must follow suit.

### Decision 2: `disk_encrypt_flag` Read requires a separate API call

**Choice:** Add a new service-layer method `DescribeSqlserverInstanceAttributeById` that calls `DescribeDBInstancesAttribute`, and invoke it in the Read function (wrapped with `tccommon.ReadRetryTimeout` retry). Failure is logged as a warning and does not abort the entire Read.

**Rationale:** The `DescribeDBInstances` response (`DBInstance` struct) does **not** include disk encryption status. Encryption status is only available via `DescribeDBInstancesAttribute` in the `IsDiskEncryptFlag` field. A separate call is therefore mandatory to populate state.

**Alternatives considered:**
- Skip reading `disk_encrypt_flag` (leave state empty): rejected — incomplete resource representation; users cannot verify encryption status.
- Fail the entire Read if the attribute API errors: rejected — degrades resilience; the existing instance fields should still be readable.

### Decision 3: `time_zone` Read reuses the existing DescribeDBInstances call

**Choice:** Read `time_zone` directly from the already-fetched `DBInstance.TimeZone` field — no new API call needed.

**Rationale:** `DescribeDBInstances` already returns `TimeZone` on the `DBInstance` struct, and the Read function already calls `DescribeSqlserverInstanceById`. Reusing the existing response avoids an extra API call.

### Decision 4: Type handling

- `time_zone`: schema `TypeString`. SDK field `TimeZone` is `*string`. On Create, `helper.String(v.(string))`. On Read, `d.Set("time_zone", instance.TimeZone)` with nil check.
- `disk_encrypt_flag`: schema `TypeInt` with `ValidateFunc: ValidateIntegerInRange(0, 1)`. SDK request field `DiskEncryptFlag` is `*int64`; on Create use `helper.IntInt64(v.(int))`. SDK response field `IsDiskEncryptFlag` is `*int64`; on Read convert with `int(*attribute.IsDiskEncryptFlag)` after double nil check.

## Risks / Trade-offs

- **[Extra API call per Read]** `DescribeDBInstancesAttribute` adds one call to each Read. Mitigation: call is necessary for complete state; failure is non-fatal (warning only).
- **[Nil pointer dereference]** `IsDiskEncryptFlag` is `*int64` and may be nil. Mitigation: double nil check (`attribute != nil && attribute.IsDiskEncryptFlag != nil`) before dereference.
- **[BACKWARD COMPATIBILITY]** New optional+computed fields must not break existing config/state. Mitigation: both fields optional; computed values populated on next refresh; no state migration needed.

## Implementation Reference (SDK & Schema)

The following is the reference data the implementation phase needs. It is extracted from the vendored SDK so the large `models.go` does not need to be re-read.

### Create API: `CreateBasicDBInstancesRequest`

Source: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328/models.go` (struct def starts ~line 1019).

New parameters of interest (both already present in the request struct):

```go
// within type CreateBasicDBInstancesRequest struct { ... }

// <p>系统时区，默认：China Standard Time</p>
TimeZone *string `json:"TimeZone,omitnil,omitempty" name:"TimeZone"`

// <p>磁盘加密标识，0-不加密，1-加密</p>
DiskEncryptFlag *int64 `json:"DiskEncryptFlag,omitnil,omitempty" name:"DiskEncryptFlag"`
```

Notes:
- `TimeZone` is `*string`, optional. API default is `China Standard Time`.
- `DiskEncryptFlag` is `*int64`, optional. `0` = no encryption, `1` = encryption.

### Read API (timezone): `DescribeDBInstances` → `DBInstance`

Source: `vendor/.../sqlserver/v20180328/models.go` (struct `DBInstance` starts ~line 2876).

The `DBInstance` struct contains (among many fields):

```go
// <p>系统字符集排序规则，默认：Chinese_PRC_CI_AS</p>
Collation *string `json:"Collation,omitnil,omitempty" name:"Collation"`

// <p>系统时区，默认：China Standard Time</p>
TimeZone *string `json:"TimeZone,omitnil,omitempty" name:"TimeZone"`
```

The Read function already calls `DescribeSqlserverInstanceById`, which returns an `*sqlserver.DBInstance`. So `time_zone` is read directly from `instance.TimeZone` (with `!= nil` check). No new API call required.

### Read API (disk encryption): `DescribeDBInstancesAttribute`

Source: `vendor/.../sqlserver/v20180328/models.go`:

**Request struct** (~line 5557):
```go
type DescribeDBInstancesAttributeRequest struct {
	*tchttp.BaseRequest

	// 实例ID
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}
```

**Response params struct** (~line 5584):
```go
type DescribeDBInstancesAttributeResponseParams struct {
	// 实例ID
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// ... (regular backup, TDE, SSL, etc. — not relevant here) ...

	// 是否开启磁盘加密，1-开启，0-未开启
	IsDiskEncryptFlag *int64 `json:"IsDiskEncryptFlag,omitnil,omitempty" name:"IsDiskEncryptFlag"`

	// ... other fields ...

	// 唯一请求 ID
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDBInstancesAttributeResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDBInstancesAttributeResponseParams `json:"Response"`
}
```

The response does **not** return a `DiskEncryptFlag` field; it returns `IsDiskEncryptFlag` (`*int64`). Map this to the schema field `disk_encrypt_flag` via `int(*attribute.IsDiskEncryptFlag)`.

### Existing resource schema (relevant portion)

File: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

Existing fields adjacent to where the new fields are added (after `engine_version` and before `period`):

```go
"engine_version": {
    Type:        schema.TypeString,
    ForceNew:    true,
    Optional:    true,
    Default:     "2008R2",
    Description: "Version of the SQL Server basic database engine. ...",
},
"time_zone": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    ForceNew:    true,
    Description: "System timezone for the SQL Server instance. Default is `China Standard Time`. This setting cannot be changed after creation.",
},
"disk_encrypt_flag": {
    Type:         schema.TypeInt,
    Optional:     true,
    Computed:     true,
    ForceNew:     true,
    ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
    Description:  "Disk encryption flag. `0` - Disabled (default), `1` - Enabled. Disk encryption cannot be changed after instance creation.",
},
"period": {
    Type:         schema.TypeInt,
    Optional:     true,
    Default:      1,
    ValidateFunc: tccommon.ValidateIntegerInRange(1, 48),
    Description:  "Purchase instance period, the default value is 1, which means one month. The value does not exceed 48.",
},
```

New schema fields to add (confirming design):

```go
"time_zone": {
    Type:        schema.TypeString,
    Optional:    true,
    Computed:    true,
    ForceNew:    true,
    Description: "System timezone for the SQL Server instance. Default is `China Standard Time`. This setting cannot be changed after creation.",
},
"disk_encrypt_flag": {
    Type:         schema.TypeInt,
    Optional:     true,
    Computed:     true,
    ForceNew:     true,
    ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
    Description:  "Disk encryption flag. `0` - Disabled (default), `1` - Enabled. Disk encryption cannot be changed after instance creation.",
},
```

### Existing Update immutableArgs

File: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`, function `resourceTencentCloudSqlserverBasicInstanceUpdate`:

```go
immutableArgs := []string{"collation", "time_zone", "disk_encrypt_flag"}
```

Both new fields must appear in this list (alongside the existing `collation`).

### Service-layer method reference

File: `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go`

**New method to add** — `DescribeSqlserverInstanceAttributeById`:

```go
func (me *SqlserverService) DescribeSqlserverInstanceAttributeById(ctx context.Context, instanceId string) (
	attribute *sqlserver.DescribeDBInstancesAttributeResponseParams,
	errRet error,
) {
	logId := tccommon.GetLogId(ctx)
	request := sqlserver.NewDescribeDBInstancesAttributeRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeDBInstancesAttribute(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		attribute = response.Response
	}
	return
}
```

**Modify** — `CreateSqlserverBasicInstance` (add paramMap handling for the two new parameters after the existing `collation` assignment):

```go
// time_zone
if v, ok := paramMap["time_zone"]; ok {
	request.TimeZone = helper.String(v.(string))
}

// disk_encrypt_flag
if v, ok := paramMap["disk_encrypt_flag"]; ok {
	request.DiskEncryptFlag = helper.IntInt64(v.(int))
}
```

### Resource-layer Create/Read reference

File: `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go`

**Create** — add to paramMap (after `collation`):

```go
// time_zone
if v, ok := d.GetOk("time_zone"); ok {
	paramMap["time_zone"] = v.(string)
}

// disk_encrypt_flag
if v, ok := d.GetOkExists("disk_encrypt_flag"); ok {
	paramMap["disk_encrypt_flag"] = v.(int)
}
```

**Read** — add after existing field sets:

```go
// time_zone
if instance.TimeZone != nil {
	_ = d.Set("time_zone", instance.TimeZone)
}

// Get disk encryption flag from attributes API
var attribute *sqlserver.DescribeDBInstancesAttributeResponseParams
outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
	attribute, inErr = sqlserverService.DescribeSqlserverInstanceAttributeById(ctx, instanceId)
	if inErr != nil {
		return tccommon.RetryError(inErr)
	}
	return nil
})
if outErr != nil {
	log.Printf("[WARN]%s describe sqlserver instance attribute failed, reason: %v", logId, outErr)
}

// disk_encrypt_flag
if attribute != nil && attribute.IsDiskEncryptFlag != nil {
	_ = d.Set("disk_encrypt_flag", int(*attribute.IsDiskEncryptFlag))
}
```

### Cloud API field summary

| Operation | Cloud API | SDK Request Field | SDK Response Field | Type |
|-----------|-----------|------------------|-------------------|------|
| Create | CreateBasicDBInstances | `TimeZone` | — | `*string` |
| Create | CreateBasicDBInstances | `DiskEncryptFlag` | — | `*int64` |
| Read (timezone) | DescribeDBInstances | — | `DBInstance.TimeZone` | `*string` |
| Read (encrypt) | DescribeDBInstancesAttribute | `InstanceId` | `IsDiskEncryptFlag` | `*int64` |