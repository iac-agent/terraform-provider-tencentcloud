# teo-edge-kv-namespace-config Specification

## Purpose
TBD - created by archiving change add-teo-edge-kv-namespace-config. Update Purpose after archive.
## Requirements
### Requirement: Read Edge KV Namespace Config
The system SHALL read the Edge KV namespace configuration by calling `DescribeEdgeKVNamespaces` API with `ZoneId` and a `Filters` parameter filtering by `namespace` name. The system SHALL set `Limit` to 1000 (the API maximum). The system SHALL set computed fields from the API response only when the corresponding response fields are not nil.

#### Scenario: Successful read
- **WHEN** the resource exists and `DescribeEdgeKVNamespaces` returns a matching namespace
- **THEN** the system sets `zone_id`, `namespace`, `remark`, `capacity`, `capacity_used`, `references` (with `reference_type`, `reference_id`, `reference_name`, `zone_id`, `zone_name`, `alias_zone_name`), `created_on`, and `modified_on` from the response, checking each field for nil before setting

#### Scenario: Resource not found
- **WHEN** `DescribeEdgeKVNamespaces` returns an empty `KVNamespaces` list or no matching namespace
- **THEN** the system logs `[CRUD] teo_edge_kv_namespace id=<current_id>` before calling `d.SetId("")` to remove the resource from state

#### Scenario: API returns nil response
- **WHEN** `DescribeEdgeKVNamespaces` returns a nil response or nil `Response`
- **THEN** the system returns a `NonRetryableError` to trigger retry, preserving the current resource ID

### Requirement: Update Edge KV Namespace Config
The system SHALL allow updating the `remark` field of an existing Edge KV namespace by calling the `ModifyEdgeKVNamespace` API. The `zone_id` and `namespace` fields SHALL be ForceNew, meaning changes to these fields trigger resource recreation. After a successful update, the system SHALL read back the updated state.

#### Scenario: Update remark
- **WHEN** user changes the `remark` field in the Terraform configuration
- **THEN** the system calls `ModifyEdgeKVNamespace` API with `zone_id`, `namespace`, and the new `remark` value, then reads back the updated state

#### Scenario: Change zone_id or namespace triggers recreation
- **WHEN** user changes `zone_id` or `namespace` in the Terraform configuration
- **THEN** Terraform destroys the existing resource and creates a new one (ForceNew behavior)

#### Scenario: Update API returns nil response
- **WHEN** the `ModifyEdgeKVNamespace` API returns a nil or empty response
- **THEN** the system returns a `NonRetryableError`

### Requirement: Import Edge KV Namespace Config
The system SHALL support importing an existing Edge KV namespace using the composite ID format `zone_id#namespace`.

#### Scenario: Successful import
- **WHEN** user runs `terraform import tencentcloud_teo_edge_kv_namespace.example zone_id#namespace`
- **THEN** the system parses the composite ID, calls `DescribeEdgeKVNamespaces` to read the resource state, and populates all fields

### Requirement: Provider Registration
The system SHALL register `tencentcloud_teo_edge_kv_namespace` in `tencentcloud/provider.go` resource map and add the corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Resource available after registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_edge_kv_namespace` is available as a valid resource type

### Requirement: Unit Tests with Gomonkey Mock
The system SHALL provide unit tests in `resource_tc_teo_edge_kv_namespace_config_test.go` using gomonkey to mock cloud API calls. Tests SHALL cover Read and Update operations and SHALL pass with `go test -gcflags=all=-l`.

#### Scenario: Test Read and Update operations
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** all tests pass, verifying the resource's Read and Update logic through mocked API responses

