## ADDED Requirements

### Requirement: Create TEO Zone

The system SHALL support creating a new TEO zone via the `CreateZone` API. The user MUST provide at minimum `zone_name`, `type`, `area`, and `plan_id`. The system SHALL return the `zone_id` and `ownership_verification` information after creation. The system SHALL poll until the zone status transitions from `pending` to a stable state.

#### Scenario: Create a CNAME access zone
- **WHEN** user creates a `tencentcloud_teo_zone` with `type = "partial"`, `zone_name = "example.com"`, `area = "overseas"`, and a valid `plan_id`
- **THEN** the system calls `CreateZone` API and returns a `zone_id` in the format `zone-xxxxxxxx`
- **AND** the system polls until zone status is no longer `pending`
- **AND** the zone is recorded in Terraform state with the returned `zone_id`

#### Scenario: Create a NS access zone
- **WHEN** user creates a `tencentcloud_teo_zone` with `type = "full"`, `zone_name = "example.com"`, `area = "global"`, and a valid `plan_id`
- **THEN** the system calls `CreateZone` API and returns a `zone_id` and `ownership_verification` with NS verification details
- **AND** the zone is recorded in Terraform state

#### Scenario: Create a no-domain access zone
- **WHEN** user creates a `tencentcloud_teo_zone` with `type = "noDomainAccess"`, `area` left empty, and a valid `plan_id`
- **THEN** the system calls `CreateZone` API successfully
- **AND** `ownership_verification` is empty since no-domain access doesn't require verification

### Requirement: Read TEO Zone

The system SHALL read the current state of a TEO zone using the `DescribeZones` API filtered by `zone-id`. The system SHALL populate all computed and configured attributes: `zone_name`, `type`, `area`, `alias_zone_name`, `plan_id`, `status`, `paused`, `ownership_verification`, `name_servers`, `tags`, and `work_mode_infos`. If the zone is not found, the system SHALL clear the resource ID from state.

#### Scenario: Read an existing zone
- **WHEN** Terraform reads a `tencentcloud_teo_zone` with a valid `zone_id`
- **THEN** the system calls `DescribeZones` with `zone-id` filter
- **AND** populates all schema attributes from the API response
- **AND** reports the current state to Terraform

#### Scenario: Read a deleted zone
- **WHEN** Terraform reads a `tencentcloud_teo_zone` with a `zone_id` that no longer exists
- **THEN** the system calls `DescribeZones` and receives an empty result
- **AND** clears the resource ID via `d.SetId("")`
- **AND** logs a warning that the zone was not found

### Requirement: Update TEO Zone Configuration

The system SHALL support updating zone configuration attributes (`type`, `alias_zone_name`, `area`) via the `ModifyZone` API. The system SHALL only call the API when at least one of these mutable fields has changed.

#### Scenario: Update zone type
- **WHEN** user changes `type` from `"partial"` to `"full"`
- **THEN** the system calls `ModifyZone` API with the new type
- **AND** the update succeeds

#### Scenario: Update alias zone name
- **WHEN** user changes `alias_zone_name` from `"alias1"` to `"alias2"`
- **THEN** the system calls `ModifyZone` API with the new alias zone name
- **AND** the update succeeds

#### Scenario: No changes to mutable fields
- **WHEN** none of `type`, `alias_zone_name`, or `area` have changed
- **THEN** the system SHALL NOT call `ModifyZone` API

### Requirement: Update TEO Zone Status

The system SHALL support pausing and resuming a TEO zone via the `ModifyZoneStatus` API. The `paused` attribute SHALL be a writable boolean: `true` to pause, `false` to resume.

#### Scenario: Pause a zone
- **WHEN** user sets `paused = true` on a running zone
- **THEN** the system calls `ModifyZoneStatus` API with `Paused = true`
- **AND** polls until the zone's `ActiveStatus` is `"paused"` or `"inactive"`

#### Scenario: Resume a zone
- **WHEN** user sets `paused = false` on a paused zone
- **THEN** the system calls `ModifyZoneStatus` API with `Paused = false`
- **AND** polls until the zone's `ActiveStatus` is not `"paused"` or `"inactive"`

### Requirement: Update TEO Zone Work Mode

The system SHALL support updating configuration group work modes via the `ModifyZoneWorkMode` API. The `work_mode_infos` attribute SHALL accept a list of configuration group and work mode pairs.

#### Scenario: Set work mode for a configuration group
- **WHEN** user sets `work_mode_infos` with `config_group_type = "l7_acceleration"` and `work_mode = "version_control"`
- **THEN** the system calls `ModifyZoneWorkMode` API with the specified work mode configuration
- **AND** the update succeeds

### Requirement: Delete TEO Zone

The system SHALL support deleting a TEO zone via the `DeleteZone` API. Before deletion, the system SHALL automatically pause the zone if it is not already paused. The zone SHALL be removed from Terraform state after successful deletion.

#### Scenario: Delete a paused zone
- **WHEN** user destroys a `tencentcloud_teo_zone` that is already paused
- **THEN** the system calls `DeleteZone` API directly
- **AND** the zone is removed from Terraform state

#### Scenario: Delete a running zone
- **WHEN** user destroys a `tencentcloud_teo_zone` that is not paused
- **THEN** the system first calls `ModifyZoneStatus` to pause the zone
- **AND** then calls `DeleteZone` API
- **AND** the zone is removed from Terraform state

### Requirement: Import TEO Zone

The system SHALL support importing an existing TEO zone into Terraform state using the `zone_id` as the import identifier.

#### Scenario: Import by zone ID
- **WHEN** user runs `terraform import tencentcloud_teo_zone.example zone-xxxxxxxx`
- **THEN** the system reads the zone details via `DescribeZones`
- **AND** populates the Terraform state with all attributes

### Requirement: TEO Zone Tag Management

The system SHALL support tag management on TEO zones via the standard TencentCloud tag service. Tags SHALL be managed as a `TypeMap` attribute. The system SHALL support creating, updating, and deleting tags.

#### Scenario: Create zone with tags
- **WHEN** user creates a `tencentcloud_teo_zone` with `tags = { env = "production" }`
- **THEN** the system creates the zone and applies the tags via the tag service
- **AND** the tags are visible in the read response

#### Scenario: Update zone tags
- **WHEN** user changes tags from `{ env = "production" }` to `{ env = "staging" }`
- **THEN** the system calls the tag service to modify the tags
- **AND** the new tags are reflected in state