## ADDED Requirements

### Requirement: Import zone configuration operation resource
The system SHALL provide a Terraform resource `tencentcloud_teo_import_zone_config_operation` that calls the `ImportZoneConfig` API to import TEO zone configuration and polls `DescribeZoneConfigImportResult` until the import task completes.

#### Scenario: Successful zone configuration import
- **WHEN** a user creates a `tencentcloud_teo_import_zone_config_operation` resource with valid `zone_id` and `content`
- **THEN** the system SHALL call `ImportZoneConfig` with the provided parameters
- **AND** the system SHALL poll `DescribeZoneConfigImportResult` using the returned `TaskId` and `ZoneId` until the status is `success`
- **AND** the system SHALL set `task_id` to the returned TaskId value
- **AND** the system SHALL set the resource ID to `zone_id + FILED_SP + task_id`

#### Scenario: Import task failure
- **WHEN** the `ImportZoneConfig` API is called successfully but the async import task fails (status=`failure`)
- **THEN** the system SHALL return an error containing the failure message from `DescribeZoneConfigImportResult`

#### Scenario: Import task polling timeout
- **WHEN** the import task does not complete within `ReadRetryTimeout`
- **THEN** the system SHALL return a retry timeout error

### Requirement: Schema definition
The resource SHALL define the following schema fields:
- `zone_id` (Required, String, ForceNew): The zone ID
- `content` (Required, String, ForceNew): The JSON configuration content to import
- `task_id` (Computed, String): The task ID returned by the API

#### Scenario: Zone ID is required
- **WHEN** a user creates the resource without specifying `zone_id`
- **THEN** Terraform SHALL return a required field error

#### Scenario: Content is required
- **WHEN** a user creates the resource without specifying `content`
- **THEN** Terraform SHALL return a required field error

### Requirement: Read handler is no-op
The Read handler SHALL be a no-op that returns nil, as this is a one-shot operation with no persistent state to read.

#### Scenario: Read after create
- **WHEN** Terraform calls Read after Create
- **THEN** the system SHALL return nil without making any API calls

### Requirement: Delete handler is no-op
The Delete handler SHALL be a no-op that returns nil, as there is no cloud resource to tear down.

#### Scenario: Resource destruction
- **WHEN** a user destroys the resource
- **THEN** the system SHALL return nil without making any API calls

### Requirement: No update support
The resource SHALL NOT support update operations. All schema fields that affect the API call SHALL be marked as ForceNew.

#### Scenario: Field modification triggers recreation
- **WHEN** a user modifies `zone_id` or `content`
- **THEN** Terraform SHALL force resource recreation

### Requirement: Provider registration
The resource SHALL be registered in `provider.go` with the name `tencentcloud_teo_import_zone_config_operation`.

#### Scenario: Resource is available in provider
- **WHEN** a user references `tencentcloud_teo_import_zone_config_operation` in a Terraform configuration
- **THEN** the provider SHALL recognize and handle the resource

### Requirement: Unit tests with mock
The resource SHALL have unit tests using gomonkey to mock the cloud API calls, following the project's testing conventions for new resources.

#### Scenario: Unit test covers create success
- **WHEN** the unit test runs with mocked `ImportZoneConfig` returning success and mocked `DescribeZoneConfigImportResult` returning status `success`
- **THEN** the create function SHALL complete without error

#### Scenario: Unit test covers create failure
- **WHEN** the unit test runs with mocked `ImportZoneConfig` returning success but mocked `DescribeZoneConfigImportResult` returning status `failure`
- **THEN** the create function SHALL return an error

### Requirement: Documentation
The resource SHALL have a `.md` documentation file with usage example, following the project's documentation conventions.

#### Scenario: Documentation exists
- **WHEN** the resource is implemented
- **THEN** a `resource_tc_teo_import_zone_config_operation.md` file SHALL exist in the `tencentcloud/services/teo/` directory with a description, example usage, and import section (if applicable)
