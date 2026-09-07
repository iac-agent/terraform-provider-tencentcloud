## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_lighthouse_instance_snapshot_1.go` with schema definition including all fields: `instance_id` (Required, ForceNew, TypeString), `snapshot_name` (Optional, Computed, TypeString), `tags` (Optional, Computed, TypeList with `key`/`value` sub-fields), `snapshot_id` (Computed, TypeString), `disk_usage` (Computed, TypeString), `disk_id` (Computed, TypeString), `disk_size` (Computed, TypeInt), `snapshot_state` (Computed, TypeString), `percent` (Computed, TypeInt), `latest_operation` (Computed, TypeString), `latest_operation_state` (Computed, TypeString), `latest_operation_request_id` (Computed, TypeString), `created_time` (Computed, TypeString)
- [x] 1.2 Implement `resourceTencentCloudLighthouseInstanceSnapshot1Create` — call `CreateInstanceSnapshot` API with `instance_id`, `snapshot_name`, `tags`; check `SnapshotId` is not nil; call `d.SetId(snapshotId)`; poll via `BuildStateChangeConf` with `LighthouseSnapshotStateRefreshFunc` until `NORMAL`; call Read
- [x] 1.3 Implement `resourceTencentCloudLighthouseInstanceSnapshot1Read` — call `DescribeLighthouseSnapshotById`; if nil, `d.SetId("")` and return; populate all computed fields with nil checks
- [x] 1.4 Implement `resourceTencentCloudLighthouseInstanceSnapshot1Update` — handle `snapshot_name` changes via `ModifySnapshotAttribute`; call Read at end
- [x] 1.5 Implement `resourceTencentCloudLighthouseInstanceSnapshot1Delete` — call `DeleteLighthouseSnapshotById`; no post-delete polling needed

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_lighthouse_instance_snapshot_1` resource in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Add resource entry in `tencentcloud/provider.md` under the Lighthouse section

## 3. Unit Tests

- [x] 3.1 Create `resource_tc_lighthouse_instance_snapshot_1_test.go` with gomonkey mock tests for Create, Read, Update, Delete functions
- [x] 3.2 Test Create: mock `CreateInstanceSnapshot` returning valid `SnapshotId`, mock `DescribeSnapshots` returning `NORMAL` state
- [x] 3.3 Test Read: mock `DescribeSnapshots` returning full snapshot with all fields populated
- [x] 3.4 Test Read with nil response: mock `DescribeSnapshots` returning empty, verify `d.SetId("")` is called
- [x] 3.5 Test Update: mock `ModifySnapshotAttribute` success
- [x] 3.6 Test Delete: mock `DeleteSnapshots` success
- [x] 3.7 Test Create with nil SnapshotId: mock `CreateInstanceSnapshot` returning nil SnapshotId, verify error

## 4. Documentation

- [x] 4.1 Create `resource_tc_lighthouse_instance_snapshot_1.md` with description, example usage, and import section

## 5. Validation

- [x] 5.1 Verify all go files are syntactically correct (no compilation errors)
- [x] 5.2 Verify resource is properly registered in provider.go
- [x] 5.3 Verify all nil checks are present for pointer fields in Read function