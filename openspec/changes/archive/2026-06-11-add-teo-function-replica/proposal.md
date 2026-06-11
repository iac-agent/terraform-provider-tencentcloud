## Why

Terraform Provider for TencentCloud currently lacks support for managing TEO (TencentCloud EdgeOne) edge function replicas. Users need to create, read, update, and delete function replicas through Terraform to enable versioned edge function deployments with different content per replica. This resource allows Infrastructure-as-Code management of TEO edge function replicas.

## What Changes

- Add new Terraform resource `tencentcloud_teo_function_replica` (RESOURCE_KIND_GENERAL) with full CRUD support:
  - **Create**: Call `CreateFunctionReplica` API to create an edge function replica with specified zone, function, name, content, and remark.
  - **Read**: Call `DescribeFunctionReplicas` API to query the replica and sync state.
  - **Update**: Call `ModifyFunctionReplica` API to update replica content and remark.
  - **Delete**: Call `DeleteFunctionReplica` API to delete the replica.
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add resource documentation markdown for `make doc` generation.
- Add unit tests with gomonkey mocks.

## Capabilities

### New Capabilities
- `teo-function-replica`: Manage TEO edge function replicas lifecycle (create, read, update, delete) via Terraform resource `tencentcloud_teo_function_replica`

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- New files: `tencentcloud/services/teo/resource_tc_teo_function_replica.go`, corresponding test and doc files
- Modified files: `tencentcloud/provider.go` (resource registration), `tencentcloud/provider.md` (resource listing)
- Cloud API dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already vendored)
- Resource uses composite ID: `zone_id#function_id#replica_name` with `tccommon.FILED_SP` separator
