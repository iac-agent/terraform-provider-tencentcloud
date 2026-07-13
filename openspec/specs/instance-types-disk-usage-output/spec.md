## Requirements

### Requirement: DiskUsage top-level computed output parameter
The `tencentcloud_instance_types` data source SHALL expose a top-level computed output parameter `DiskUsage` that indicates the disk usage type (system disk or data disk) from the `DescribeDiskConfigQuota` API response.

#### Scenario: DiskUsage populated from cbs_filter input
- **WHEN** user queries the `tencentcloud_instance_types` data source with `cbs_filter` containing `disk_usage` value
- **THEN** the data source SHALL set the top-level `DiskUsage` computed field to the value provided in `cbs_filter.disk_usage`
- **AND** `DiskUsage` SHALL be `SYSTEM_DISK` when the user specifies `SYSTEM_DISK` as the disk_usage filter
- **AND** `DiskUsage` SHALL be `DATA_DISK` when the user specifies `DATA_DISK` as the disk_usage filter

#### Scenario: DiskUsage not populated without cbs_filter
- **WHEN** user queries the `tencentcloud_instance_types` data source without providing `cbs_filter`
- **THEN** the top-level `DiskUsage` computed field SHALL be empty/null
- **AND** this SHALL not cause any errors in the data source operation

#### Scenario: DiskUsage backward compatibility
- **WHEN** user has an existing Terraform configuration using the `tencentcloud_instance_types` data source that does not reference `DiskUsage`
- **THEN** the configuration SHALL continue to work without modification
- **AND** terraform plan SHALL not show any changes related to the new `DiskUsage` field
- **AND** existing fields including `disk_usage` in `cbs_filter` and `cbs_configs` SHALL remain unchanged
