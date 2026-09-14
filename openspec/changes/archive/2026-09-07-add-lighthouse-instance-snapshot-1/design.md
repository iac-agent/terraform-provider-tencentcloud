## Context

The existing `tencentcloud_lighthouse_snapshot` resource provides basic CRUD for Lighthouse instance snapshots but only exposes `instance_id` and `snapshot_name` fields. The Lighthouse `CreateInstanceSnapshot` API supports `Tags` on creation, and the `Snapshot` response struct has 12+ fields that are useful to users (disk info, state, progress, timestamps, tags). A new resource is needed to expose these capabilities.

The Lighthouse SDK (`v20200324`) is already vendored at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324/`. The `LightHouseService` in `service_tencentcloud_lighthouse.go` already has `DescribeLighthouseSnapshotById`, `DeleteLighthouseSnapshotById`, and `LighthouseSnapshotStateRefreshFunc` methods that can be reused.

## Goals / Non-Goals

**Goals:**
- Create a new `tencentcloud_lighthouse_instance_snapshot_1` resource with full CRUD
- Support `tags` parameter during snapshot creation via `CreateInstanceSnapshot` API
- Expose all computed fields from the `Snapshot` struct: `snapshot_id`, `disk_usage`, `disk_id`, `disk_size`, `snapshot_state`, `percent`, `latest_operation`, `latest_operation_state`, `latest_operation_request_id`, `created_time`, `tags`
- Support updating `snapshot_name` via `ModifySnapshotAttribute`
- Handle async creation: poll `DescribeSnapshots` until `SnapshotState` is `NORMAL`
- Follow existing code patterns (reference: `resource_tc_lighthouse_snapshot.go`)

**Non-Goals:**
- Do NOT modify the existing `tencentcloud_lighthouse_snapshot` resource
- Do NOT implement `ApplyInstanceSnapshot` (rollback) — this is a separate operation
- Do NOT expose pagination controls (`offset`/`limit`) to users
- Do NOT modify the `LightHouseService` — reuse existing methods

## Decisions

### 1. New resource name: `tencentcloud_lighthouse_instance_snapshot_1`

**Rationale**: The existing resource `tencentcloud_lighthouse_snapshot` already occupies the obvious name. The suffix `_1` follows the convention seen in the requirement to distinguish this enhanced version. The file name will be `resource_tc_lighthouse_instance_snapshot_1.go`.

### 2. Tags schema: TypeList with `key`/`value` sub-fields

The `CreateInstanceSnapshot` API accepts `Tags []*Tag` where `Tag` has `Key` and `Value` fields. The Terraform schema will use:

```hcl
tags {
  key   = "env"
  value = "production"
}
```

Schema definition:
```go
"tags": {
    Optional:    true,
    Type:        schema.TypeList,
    Description: "Tags to associate with the snapshot.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "key": {
                Required:    true,
                Type:        schema.TypeString,
                Description: "Tag key.",
            },
            "value": {
                Required:    true,
                Type:        schema.TypeString,
                Description: "Tag value.",
            },
        },
    },
},
```

Same structure for computed `tags` in the Read function.

### 3. Reuse existing service methods

The following methods from `LightHouseService` are reused:
- `DescribeLighthouseSnapshotById(ctx, snapshotId)` — returns `*lighthouse.Snapshot`
- `DeleteLighthouseSnapshotById(ctx, snapshotId)` — deletes a snapshot
- `LighthouseSnapshotStateRefreshFunc(snapshotId, failStates)` — returns `StateRefreshFunc` for polling

No new service methods are needed.

### 4. Async creation pattern

`CreateInstanceSnapshot` is asynchronous. After calling the API:
1. Check `response.Response.SnapshotId` is not nil/empty
2. Call `d.SetId(snapshotId)` **before** polling — ensures tfstate has the ID even if polling fails
3. Use `BuildStateChangeConf` with `LighthouseSnapshotStateRefreshFunc` to wait for `NORMAL` state
4. Call Read to populate all computed fields

### 5. Computed fields from Snapshot struct

All fields from the `Snapshot` struct are mapped as Computed-only in the schema:

| SDK Field | Schema Field | Type |
|-----------|-------------|------|
| `SnapshotId` | `snapshot_id` | TypeString |
| `DiskUsage` | `disk_usage` | TypeString |
| `DiskId` | `disk_id` | TypeString |
| `DiskSize` | `disk_size` | TypeInt |
| `SnapshotName` | `snapshot_name` | TypeString (also Optional) |
| `SnapshotState` | `snapshot_state` | TypeString |
| `Percent` | `percent` | TypeInt |
| `LatestOperation` | `latest_operation` | TypeString |
| `LatestOperationState` | `latest_operation_state` | TypeString |
| `LatestOperationRequestId` | `latest_operation_request_id` | TypeString |
| `CreatedTime` | `created_time` | TypeString |
| `Tags` | `tags` | TypeList (also Optional) |

### 6. Update: only `snapshot_name` is updatable

`ModifySnapshotAttribute` only supports changing `snapshot_name`. The `instance_id` is ForceNew (snapshot belongs to a specific instance). All other fields are Computed.

### 7. Delete: reuse existing service method

`DeleteSnapshots` accepts `SnapshotIds []*string`. The `DeleteLighthouseSnapshotById` method wraps this with a single ID. No retry/post-delete polling is needed.

## Risks / Trade-offs

- **Risk**: New resource name `_1` suffix may confuse users with the existing `tencentcloud_lighthouse_snapshot` → **Mitigation**: Clear documentation distinguishing the two resources
- **Risk**: Tags cannot be modified after creation (API only supports tags on create) → **Mitigation**: Document this limitation; if user needs to change tags, they must recreate the resource
- **Risk**: Async creation may fail after `d.SetId()` → **Mitigation**: This is the standard pattern; `d.SetId()` before polling ensures `terraform destroy` can clean up even when creation fails mid-way