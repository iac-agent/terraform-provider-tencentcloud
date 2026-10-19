## Requirements

### Requirement: Create Lighthouse instance snapshot
The system SHALL support creating a Lighthouse instance snapshot via the `CreateInstanceSnapshot` API. The resource MUST accept `instance_id` (Required, ForceNew) and optionally `snapshot_name` and `tags`. After the API returns a `SnapshotId`, the system SHALL set the Terraform resource ID and then poll the `DescribeSnapshots` API until `SnapshotState` reaches `NORMAL`.

#### Scenario: Create snapshot with required fields only
- **WHEN** user provides only `instance_id`
- **THEN** the system creates the snapshot, sets the resource ID, and waits for `NORMAL` state

#### Scenario: Create snapshot with tags
- **WHEN** user provides `instance_id` and `tags` with key-value pairs
- **THEN** the system creates the snapshot with the specified tags bound to it

#### Scenario: Create snapshot with custom name
- **WHEN** user provides `instance_id` and `snapshot_name`
- **THEN** the system creates the snapshot with the specified name

#### Scenario: Create fails with empty SnapshotId
- **WHEN** the `CreateInstanceSnapshot` API returns a response with nil or empty `SnapshotId`
- **THEN** the system returns a non-retryable error

### Requirement: Read Lighthouse instance snapshot
The system SHALL read the snapshot's current state via the `DescribeSnapshots` API. All computed fields from the `Snapshot` struct MUST be populated: `snapshot_id`, `disk_usage`, `disk_id`, `disk_size`, `snapshot_state`, `percent`, `latest_operation`, `latest_operation_state`, `latest_operation_request_id`, `created_time`, `tags`. The `snapshot_name` and `tags` fields SHALL also be populated from the API response.

#### Scenario: Read existing snapshot
- **WHEN** the resource ID is set to a valid snapshot ID
- **THEN** the system populates all computed fields from the API response

#### Scenario: Read deleted snapshot
- **WHEN** the snapshot no longer exists (API returns empty result)
- **THEN** the system logs the situation and calls `d.SetId("")` to remove the resource from state

### Requirement: Update Lighthouse instance snapshot name
The system SHALL support updating the `snapshot_name` via the `ModifySnapshotAttribute` API. No other fields are updatable. The `instance_id` is ForceNew and cannot be changed.

#### Scenario: Update snapshot name
- **WHEN** user changes `snapshot_name` in the Terraform configuration
- **THEN** the system calls `ModifySnapshotAttribute` with the new name and reads back the updated state

#### Scenario: Attempt to change instance_id
- **WHEN** user changes `instance_id` in the Terraform configuration
- **THEN** Terraform forces resource recreation (ForceNew)

### Requirement: Delete Lighthouse instance snapshot
The system SHALL delete the snapshot via the `DeleteSnapshots` API.

#### Scenario: Delete snapshot
- **WHEN** user removes the resource from configuration or runs `terraform destroy`
- **THEN** the system calls `DeleteSnapshots` with the snapshot ID and removes the resource from state

### Requirement: Async creation preserves resource ID
The system SHALL call `d.SetId()` with the snapshot ID immediately after a successful `CreateInstanceSnapshot` API call, before starting the async polling loop. This ensures that even if the polling fails, the Terraform state contains the resource ID, allowing `terraform destroy` to clean up the partially-created resource.

#### Scenario: Creation succeeds but polling times out
- **WHEN** the `CreateInstanceSnapshot` API returns a valid `SnapshotId` but the polling loop times out before reaching `NORMAL`
- **THEN** the resource ID is preserved in state, enabling `terraform destroy` to delete the snapshot