## Why

The existing `tencentcloud_lighthouse_snapshot` resource is minimal — it only exposes `instance_id` and `snapshot_name` fields, lacking support for tags and rich computed attributes (snapshot state, disk info, creation time, etc.). Users need a more complete snapshot resource that covers all fields available in the Lighthouse `Snapshot` API response, including tag management during creation and comprehensive read-only status fields.

## What Changes

- Add a new Terraform resource `tencentcloud_lighthouse_instance_snapshot_1` for Lighthouse instance snapshots
- Support the same CRUD APIs as the existing snapshot resource (`CreateInstanceSnapshot`, `DescribeSnapshots`, `ModifySnapshotAttribute`, `DeleteSnapshots`)
- **Additional capabilities over the existing resource**:
  - Support `tags` parameter during snapshot creation
  - Expose all computed output fields from the Snapshot struct: `snapshot_id`, `disk_usage`, `disk_id`, `disk_size`, `snapshot_state`, `percent`, `latest_operation`, `latest_operation_state`, `latest_operation_request_id`, `created_time`, `tags`
  - The `instance_id` is ForceNew (snapshot cannot be moved between instances)
  - The `snapshot_name` is updatable via `ModifySnapshotAttribute`

## Capabilities

### New Capabilities
- `lighthouse-instance-snapshot-1`: Full CRUD resource for Lighthouse instance snapshots with tags support and comprehensive computed snapshot attributes

### Modified Capabilities
<!-- No existing capabilities are modified. This is a new resource. -->

## Impact

- **New file**: `tencentcloud/services/lighthouse/resource_tc_lighthouse_instance_snapshot_1.go` — resource CRUD implementation
- **New file**: `tencentcloud/services/lighthouse/resource_tc_lighthouse_instance_snapshot_1_test.go` — unit tests with gomonkey mock
- **New file**: `tencentcloud/services/lighthouse/resource_tc_lighthouse_instance_snapshot_1.md` — documentation
- **Modified file**: `tencentcloud/provider.go` — register new resource in ResourcesMap
- **Modified file**: `tencentcloud/provider.md` — add resource to documentation index
- **Reuses existing service methods**: `DescribeLighthouseSnapshotById`, `DeleteLighthouseSnapshotById`, `LighthouseSnapshotStateRefreshFunc` from `service_tencentcloud_lighthouse.go`
- **No breaking changes**: Existing `tencentcloud_lighthouse_snapshot` resource is untouched
- **Dependencies**: Uses existing SDK `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324` (already vendored)