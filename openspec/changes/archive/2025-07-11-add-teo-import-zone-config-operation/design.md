## Context

Terraform Provider for TencentCloud currently lacks a resource to import TEO (EdgeOne) zone configurations. The TEO API provides `ImportZoneConfig` which is an asynchronous operation that returns a `TaskId`, and `DescribeZoneConfigImportResult` for polling the task status. This is a one-shot operation (RESOURCE_KIND_OPERATION) where no persistent state needs to be tracked after the import completes.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_teo_import_zone_config_operation` that triggers zone configuration import
- After calling `ImportZoneConfig`, poll `DescribeZoneConfigImportResult` until the async task completes (status is `success` or `failure`)
- Follow the one-shot operation pattern: Create calls the API, Read/Delete are no-ops, no Update
- Register the resource in `provider.go`

**Non-Goals:**
- This resource does NOT manage the zone configuration lifecycle (that's handled by other TEO resources)
- No import support (one-shot operations are not importable)
- No update support (each apply triggers a new import operation)

## Decisions

### 1. Resource Type: One-shot Operation
**Decision**: Use RESOURCE_KIND_OPERATION pattern.
**Rationale**: ImportZoneConfig is a one-time action that imports configuration into a zone. There's no persistent resource to manage after the import completes. The resource ID will be set to `zone_id + FILED_SP + task_id` (using `tccommon.FILED_SP` as separator) to uniquely identify the operation.

### 2. Async Polling
**Decision**: After calling `ImportZoneConfig`, poll `DescribeZoneConfigImportResult` using `helper.Retry()` with `tccommon.ReadRetryTimeout` until the status transitions from `doing` to `success` or `failure`.
**Rationale**: The API documentation states that `ImportZoneConfig` returns a `TaskId` and users must call `DescribeZoneConfigImportResult` to check the result. The task status can be `success`, `failure`, or `doing`. If the import fails, the error message should be surfaced to the user.

### 3. Schema Fields
**Decision**:
- `zone_id` (Required, ForceNew): The zone ID to import configuration into
- `content` (Required, ForceNew): The JSON configuration content to import
- `task_id` (Computed): The task ID returned by the API, used for polling

**Rationale**: These map directly to the `ImportZoneConfig` API input parameters and the returned `TaskId`. Both `zone_id` and `content` are ForceNew because this is a one-shot operation - any change triggers a new import.

### 4. File Naming
**Decision**: `resource_tc_teo_import_zone_config_operation.go`
**Rationale**: Follows the naming convention `resource_tc_<Product>_<Operation>_operation.go` for RESOURCE_KIND_OPERATION resources.

## Risks / Trade-offs

- **[Import failure]** → If the import task fails (status=`failure`), the resource Create will return an error with the failure message from the API, causing Terraform to mark the resource as tainted.
- **[Polling timeout]** → If the import task takes longer than `ReadRetryTimeout`, the Create will fail with a timeout error. Users can adjust the timeout in the Terraform configuration.
- **[Idempotency]** → Each `terraform apply` with this resource will trigger a new import operation, even if the content hasn't changed. This is inherent to the one-shot operation pattern.
