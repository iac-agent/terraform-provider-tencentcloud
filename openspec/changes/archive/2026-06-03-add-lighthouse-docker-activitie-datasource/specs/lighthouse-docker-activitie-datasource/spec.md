## ADDED Requirements

### Requirement: Query Docker Activities by Instance ID
The system SHALL allow users to query Docker activities for a specific Lighthouse instance using the `instance_id` parameter.

#### Scenario: Query Docker activities by instance ID
- **WHEN** the user defines a data source with `instance_id` set
- **THEN** the system SHALL call the `DescribeDockerActivities` API with the specified InstanceId
- **AND** return matching Docker activities in `docker_activity_set`

#### Scenario: Query without instance_id
- **WHEN** the user defines a data source without `instance_id`
- **THEN** the system SHALL call the API without InstanceId filter
- **AND** return all accessible Docker activities

### Requirement: Filter Docker Activities by Activity IDs
The system SHALL allow users to filter Docker activities by a list of activity IDs using the `activity_ids` parameter.

#### Scenario: Filter by specific activity IDs
- **WHEN** the user provides `activity_ids` as a list of strings
- **THEN** the system SHALL pass ActivityIds to the API
- **AND** return only the Docker activities matching those IDs

### Requirement: Filter Docker Activities by Time Range
The system SHALL allow users to filter Docker activities by creation time range using `created_time_begin` and `created_time_end` parameters.

#### Scenario: Filter by time range
- **WHEN** the user provides both `created_time_begin` and `created_time_end`
- **THEN** the system SHALL pass CreatedTimeBegin and CreatedTimeEnd to the API
- **AND** return only Docker activities created within the specified time range

#### Scenario: Filter with only begin time
- **WHEN** the user provides only `created_time_begin`
- **THEN** the system SHALL pass CreatedTimeBegin to the API
- **AND** return Docker activities created after the specified time

### Requirement: Return Complete Docker Activity Information
The system SHALL return complete Docker activity information in the `docker_activity_set` computed attribute, mapping all fields from the `DockerActivity` SDK struct.

#### Scenario: All fields populated for a Docker activity
- **WHEN** the API returns a DockerActivity with all fields present
- **THEN** each item in `docker_activity_set` SHALL include:
  - `activity_id` (string): Activity ID
  - `activity_name` (string): Activity name
  - `activity_state` (string): Activity state (INIT/OPERATING/SUCCESS/FAILED)
  - `activity_command_output` (string): Activity command output (base64 encoded)
  - `container_ids` (list of string): Container ID list
  - `created_time` (string): Creation time (ISO 8601)
  - `end_time` (string): End time (ISO 8601)

#### Scenario: Handle nullable fields gracefully
- **WHEN** the API returns a DockerActivity with some nil fields
- **THEN** the system SHALL skip setting those nil fields instead of causing panics
- **AND** the data source SHALL NOT crash on null fields

### Requirement: Automatic Pagination
The system SHALL automatically handle pagination when querying Docker activities, hiding Offset and Limit from users.

#### Scenario: Many Docker activities returned
- **WHEN** the API returns more results than one page can hold (Limit=100)
- **THEN** the service layer SHALL automatically paginate through all results
- **AND** return the complete list to the user

#### Scenario: Few Docker activities returned
- **WHEN** the API returns results that fit within one page
- **THEN** the system SHALL return results without additional pagination

### Requirement: Result Output File
The system SHALL support saving query results to a JSON file via the `result_output_file` parameter.

#### Scenario: Save results to file
- **WHEN** the user specifies `result_output_file`
- **THEN** the system SHALL write the `docker_activity_set` data to the specified file

### Requirement: Provider Registration
The system SHALL register the new data source in the Terraform provider.

#### Scenario: Data source is available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_lighthouse_docker_activitie` SHALL be available as a data source
- **AND** it SHALL be listed in `provider.md`

### Requirement: Unit Tests with Mock
The system SHALL include unit tests using gomonkey mock for the data source read function.

#### Scenario: Mock-based unit tests
- **WHEN** the test suite is run with `go test -gcflags=all=-l`
- **THEN** tests SHALL mock the cloud API calls
- **AND** verify the data source correctly maps request parameters and response fields
