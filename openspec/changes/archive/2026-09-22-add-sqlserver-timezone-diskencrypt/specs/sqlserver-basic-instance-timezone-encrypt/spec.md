# Spec: SqlServer Basic Instance - TimeZone and DiskEncryptFlag Support

## Overview

This specification defines the requirements for adding `time_zone` and `disk_encrypt_flag` parameters to the `tencentcloud_sqlserver_basic_instance` resource, enabling users to configure timezone settings and disk encryption for SQL Server instances.

## ADDED Requirements

### Requirement: TimeZone Parameter Support

The `tencentcloud_sqlserver_basic_instance` resource MUST support a `time_zone` parameter that allows users to specify the system timezone for SQL Server instances.

#### Scenario: Create Instance with Custom Timezone

- **WHEN** the user specifies `time_zone = "UTC"` in the resource configuration
- **THEN** the instance MUST be created with UTC timezone
- **AND** the timezone MUST be retrievable via `DescribeDBInstances` API
- **AND** the Terraform state MUST show `time_zone = "UTC"`

#### Scenario: Create Instance with Default Timezone

- **WHEN** a user creates a SQL Server instance without specifying `time_zone`
- **THEN** the API default timezone MUST be used ("China Standard Time")
- **AND** the Terraform state MUST be populated with the actual timezone from API
- **AND** no error MUST occur

#### Scenario: Read Existing Instance Timezone

- **WHEN** Terraform reads the instance state
- **THEN** the `time_zone` MUST be retrieved from `DescribeDBInstances` response
- **AND** the value MUST match the instance configuration
- **AND** the state MUST be updated correctly

#### Scenario: Timezone Change Triggers Recreation

- **WHEN** the user changes `time_zone` to "UTC" in configuration
- **THEN** Terraform plan MUST show the instance will be destroyed and recreated
- **AND** the plan MUST indicate `time_zone` change forces replacement
- **AND** applying the plan MUST recreate the instance with new timezone

#### Scenario: Import Populates Timezone

- **WHEN** the user imports the instance with `terraform import`
- **THEN** the Read function MUST be called
- **AND** the `time_zone` MUST be populated from API
- **AND** the state MUST contain the correct timezone value

#### Scenario: Timezone Validation

- **WHEN** a user specifies an invalid timezone value
- **THEN** the API MUST validate the timezone string
- **AND** invalid timezones MUST result in API error
- **AND** error MUST be surfaced to the user

### Requirement: Disk Encryption Flag Support

The `tencentcloud_sqlserver_basic_instance` resource MUST support a `disk_encrypt_flag` parameter that allows users to enable disk encryption for SQL Server instances.

#### Scenario: Create Instance with Encryption Enabled

- **WHEN** the user specifies `disk_encrypt_flag = 1` in the configuration
- **THEN** the instance MUST be created with encryption enabled
- **AND** the encryption status MUST be retrievable via `DescribeDBInstancesAttribute` API
- **AND** the Terraform state MUST show `disk_encrypt_flag = 1`

#### Scenario: Create Instance with Encryption Disabled (Default)

- **WHEN** a user creates a SQL Server instance without specifying `disk_encrypt_flag`
- **THEN** encryption MUST be disabled (default value 0)
- **AND** the Terraform state MUST show `disk_encrypt_flag = 0`
- **AND** no error MUST occur

#### Scenario: Read Existing Instance Encryption Status

- **WHEN** Terraform reads the instance state
- **THEN** the `disk_encrypt_flag` MUST be retrieved via `DescribeDBInstancesAttribute` API
- **AND** the `IsDiskEncryptFlag` field MUST be accessed
- **AND** the value MUST be converted from `*int64` to `int`
- **AND** the state MUST be updated correctly

#### Scenario: Read Operation Handles API Failure Gracefully

- **WHEN** the `DescribeDBInstancesAttribute` API call fails
- **THEN** a warning MUST be logged
- **AND** the Read operation MUST NOT fail entirely
- **AND** other fields MUST still be populated

#### Scenario: Encryption Change Triggers Recreation

- **WHEN** the user changes `disk_encrypt_flag` to 1 in configuration
- **THEN** Terraform plan MUST show the instance will be destroyed and recreated
- **AND** applying the plan MUST recreate the instance with encryption enabled

#### Scenario: Import Populates Encryption Status

- **WHEN** the user imports the instance
- **THEN** the Read function MUST call both APIs
- **AND** the `disk_encrypt_flag` MUST be populated from `DescribeDBInstancesAttribute`
- **AND** the state MUST contain the correct encryption status

#### Scenario: Encryption Flag Validation

- **WHEN** a user specifies an invalid `disk_encrypt_flag` value (e.g., 2)
- **THEN** a validation error MUST occur
- **AND** the error MUST indicate valid values are 0 or 1
- **AND** the resource creation MUST NOT proceed

### Requirement: Combined Parameter Support

The resource MUST support using both `time_zone` and `disk_encrypt_flag` parameters simultaneously without conflicts.

#### Scenario: Create Instance with Both Parameters

- **WHEN** both parameters are provided in configuration
- **THEN** the instance MUST be created with both settings
- **AND** both values MUST be retrievable
- **AND** the state MUST contain both values correctly

#### Scenario: Change Either Parameter Triggers Recreation

- **WHEN** either parameter is changed
- **THEN** the instance MUST be recreated
- **AND** the plan MUST show which parameter changed
- **AND** both parameters MUST be set correctly on new instance

### Requirement: Schema Design

The schema definition MUST follow Terraform and project conventions.

#### Scenario: Schema Properties for time_zone

- **WHEN** the schema definition for `time_zone` is created
- **THEN** it MUST have TypeString, Optional, Computed, ForceNew properties
- **AND** it MUST have a clear description mentioning common values

#### Scenario: Schema Properties for disk_encrypt_flag

- **WHEN** the schema definition for `disk_encrypt_flag` is created
- **THEN** it MUST have TypeInt, Optional, Computed, ForceNew, Default: 0 properties
- **AND** it MUST have ValidateFunc enforcing range (0, 1)
- **AND** it MUST have a clear description explaining values

#### Scenario: Immutable Args Configuration

- **WHEN** both parameters are ForceNew and a user attempts to modify them
- **THEN** the immutableArgs list MUST include both parameters
- **AND** Terraform MUST prevent in-place updates
- **AND** recreation MUST be required

### Requirement: API Integration

The implementation MUST correctly integrate with TencentCloud SQL Server APIs.

#### Scenario: Create API Parameter Mapping

- **WHEN** creating an instance with `CreateBasicDBInstances` API
- **THEN** `time_zone` (schema) MUST map to `TimeZone` (API)
- **AND** `disk_encrypt_flag` (schema) MUST map to `DiskEncryptFlag` (API)
- **AND** correct type conversions MUST be applied

#### Scenario: Read API - TimeZone from DescribeDBInstances

- **WHEN** reading instance state with `DescribeDBInstances` API
- **THEN** the `TimeZone` field MUST be accessed from `DBInstance` struct
- **AND** nil pointer check MUST be performed
- **AND** value MUST be set in state

#### Scenario: Read API - DiskEncryptFlag from DescribeDBInstancesAttribute

- **WHEN** reading instance state
- **THEN** a separate `DescribeDBInstancesAttribute` API call MUST be made
- **AND** the `IsDiskEncryptFlag` field MUST be accessed
- **AND** type conversion from `*int64` to `int` MUST be performed
- **AND** nil checks MUST be performed
- **AND** value MUST be set in state

### Requirement: Error Handling

The implementation MUST handle errors gracefully and provide useful error messages.

#### Scenario: DescribeDBInstancesAttribute API Failure

- **WHEN** the `DescribeDBInstancesAttribute` API call fails during read
- **THEN** a warning MUST be logged
- **AND** the error MUST NOT fail the entire Read operation
- **AND** other fields MUST still be populated

#### Scenario: Nil Pointer Safety

- **WHEN** accessing pointer fields in API responses
- **THEN** nil checks MUST be performed
- **AND** no nil pointer dereferences MUST occur
- **AND** missing fields MUST be handled gracefully

### Requirement: Testing

Comprehensive tests MUST be provided to verify all functionality.

#### Scenario: Acceptance Test - Timezone

- **WHEN** the timezone acceptance test runs
- **THEN** it MUST create instance with custom timezone
- **AND** verify timezone in state
- **AND** test import functionality
- **AND** verify ForceNew behavior

#### Scenario: Acceptance Test - Disk Encryption

- **WHEN** the disk encryption acceptance test runs
- **THEN** it MUST create instance with encryption enabled
- **AND** verify `disk_encrypt_flag=1` in state
- **AND** test import functionality
- **AND** verify ForceNew behavior

#### Scenario: Acceptance Test - Both Parameters

- **WHEN** the combined parameters acceptance test runs
- **THEN** it MUST create instance with both timezone and encryption
- **AND** verify both values in state
- **AND** test import
- **AND** verify no conflicts

### Requirement: Documentation

Complete and accurate documentation MUST be provided.

#### Scenario: Resource Documentation File

- **WHEN** the resource documentation file is created
- **THEN** it MUST include usage examples with timezone
- **AND** usage examples with encryption
- **AND** common timezone values
- **AND** encryption implications
- **AND** ForceNew behavior notes
- **AND** import examples

#### Scenario: Website Documentation Generation

- **WHEN** documentation is generated with `make doc`
- **THEN** the website docs MUST include both parameters in Argument Reference
- **AND** correct descriptions
- **AND** proper formatting

### Requirement: Backward Compatibility

The changes MUST be fully backward compatible with existing configurations.

#### Scenario: Existing Configurations Work Unchanged

- **WHEN** an existing Terraform configuration without new parameters is applied with the updated provider
- **THEN** no changes MUST be required
- **AND** the instance MUST continue to work
- **AND** the state MUST be populated with computed values

#### Scenario: State Migration Not Required

- **WHEN** existing Terraform state files are used with the updated provider
- **THEN** no state migration MUST be required
- **AND** the state format MUST remain compatible
- **AND** Terraform refresh MUST populate new fields

## API Mappings

### Create Operation: CreateBasicDBInstances

| Schema Field | API Field | Type Conversion | Required | Default |
|--------------|-----------|-----------------|----------|---------|
| `time_zone` | `TimeZone` | string -> *string | No | "China Standard Time" (API) |
| `disk_encrypt_flag` | `DiskEncryptFlag` | int -> *int64 | No | 0 |

### Read Operation: DescribeDBInstances

| API Field | Schema Field | Type Conversion | Always Present |
|-----------|--------------|-----------------|----------------|
| `DBInstance.TimeZone` | `time_zone` | *string -> string | Yes |

### Read Operation: DescribeDBInstancesAttribute

| API Field | Schema Field | Type Conversion | Always Present |
|-----------|--------------|-----------------|----------------|
| `IsDiskEncryptFlag` | `disk_encrypt_flag` | *int64 -> int | Yes |