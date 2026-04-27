## ADDED Requirements

### Requirement: Resource schema definition
The system SHALL define a Terraform resource `tencentcloud_teo_multi_path_gateway_line` with the following schema fields:
- `zone_id`: TypeString, Required, ForceNew — 站点 ID
- `gateway_id`: TypeString, Required, ForceNew — 多通道安全加速网关 ID
- `line_type`: TypeString, Required — 线路类型，取值有 direct/proxy/custom
- `line_address`: TypeString, Required — 线路地址，格式为 ip:port
- `proxy_id`: TypeString, Optional — 四层代理实例 ID，当 LineType 为 proxy 时必传
- `rule_id`: TypeString, Optional — 转发规则 ID，当 LineType 为 proxy 时必传
- `line_id`: TypeString, Computed — 线路 ID，由云 API 返回

The resource SHALL support Import.

#### Scenario: Schema fields defined correctly
- **WHEN** the resource is registered
- **THEN** the schema contains zone_id (Required, ForceNew), gateway_id (Required, ForceNew), line_type (Required), line_address (Required), proxy_id (Optional), rule_id (Optional), line_id (Computed)

#### Scenario: Resource supports import
- **WHEN** a user imports an existing multi-path gateway line resource
- **THEN** the resource state SHALL be populated by reading the current configuration from the cloud API

### Requirement: Create multi-path gateway line
The system SHALL create a multi-path gateway line by calling `CreateMultiPathGatewayLine` API with zone_id, gateway_id, line_type, line_address, proxy_id, and rule_id. The returned line_id SHALL be stored in the Terraform state.

The resource ID SHALL be composed as `zone_id + FIELD_SP + gateway_id + FIELD_SP + line_id`.

After creation, the system SHALL call Read to refresh the resource state.

#### Scenario: Successful creation of custom line
- **WHEN** a user creates a `tencentcloud_teo_multi_path_gateway_line` resource with line_type="custom", zone_id, gateway_id, and line_address
- **THEN** the system calls CreateMultiPathGatewayLine with the provided parameters
- **AND** the resource ID is set to zone_id + FIELD_SP + gateway_id + FIELD_SP + line_id
- **AND** the state is refreshed by calling Read

#### Scenario: Successful creation of proxy line
- **WHEN** a user creates a `tencentcloud_teo_multi_path_gateway_line` resource with line_type="proxy", zone_id, gateway_id, line_address, proxy_id, and rule_id
- **THEN** the system calls CreateMultiPathGatewayLine with all provided parameters including proxy_id and rule_id

### Requirement: Read multi-path gateway line
The system SHALL read a multi-path gateway line by calling `DescribeMultiPathGatewayLine` API. The request SHALL use zone_id, gateway_id, and line_id obtained from `d.Get()` (not from `d.Id()` parsing).

The response contains a `Line` object of type `MultiPathGatewayLine` with fields: LineId, LineType, LineAddress, ProxyId, RuleId. These SHALL be mapped to the corresponding Terraform schema fields.

#### Scenario: Successful read
- **WHEN** the system reads a `tencentcloud_teo_multi_path_gateway_line` resource
- **THEN** it calls DescribeMultiPathGatewayLine with zone_id, gateway_id, and line_id from d.Get()
- **AND** maps Line.LineType → line_type, Line.LineAddress → line_address, Line.ProxyId → proxy_id, Line.RuleId → rule_id, Line.LineId → line_id

#### Scenario: Resource not found
- **WHEN** DescribeMultiPathGatewayLine returns a resource-not-found error
- **THEN** the resource SHALL be removed from the Terraform state

### Requirement: Update multi-path gateway line
The system SHALL update a multi-path gateway line by calling `ModifyMultiPathGatewayLine` API with zone_id, gateway_id, line_id (from d.Get()), and the modified fields: line_type, line_address, proxy_id, rule_id.

After update, the system SHALL call Read to refresh the resource state.

#### Scenario: Successful update of line_type
- **WHEN** a user updates the line_type field of a `tencentcloud_teo_multi_path_gateway_line` resource
- **THEN** the system calls ModifyMultiPathGatewayLine with zone_id, gateway_id, line_id, and the new line_type value
- **AND** refreshes the state by calling Read

#### Scenario: Successful update of proxy_id and rule_id
- **WHEN** a user updates the proxy_id and rule_id fields
- **THEN** the system calls ModifyMultiPathGatewayLine with all fields including the new proxy_id and rule_id

### Requirement: Delete multi-path gateway line
The system SHALL delete a multi-path gateway line by calling `DeleteMultiPathGatewayLine` API with zone_id, gateway_id, and line_id obtained from `d.Get()`.

#### Scenario: Successful deletion
- **WHEN** a user deletes a `tencentcloud_teo_multi_path_gateway_line` resource
- **THEN** the system calls DeleteMultiPathGatewayLine with zone_id, gateway_id, and line_id from d.Get()
- **AND** the resource is removed from the Terraform state

### Requirement: Resource registration in provider
The system SHALL register the `tencentcloud_teo_multi_path_gateway_line` resource in `tencentcloud/provider.go` with the key `"tencentcloud_teo_multi_path_gateway_line"`, and add corresponding documentation entry in `tencentcloud/provider.md`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_multi_path_gateway_line` is available for use in Terraform configurations

### Requirement: Error handling with retry
The system SHALL use `helper.Retry()` with `tccommon.ReadRetryTimeout` for Read and Delete operations. If the API call fails, the error SHALL be wrapped with `tccommon.RetryError()` and returned.

#### Scenario: Read operation retries on transient failure
- **WHEN** DescribeMultiPathGatewayLine fails with a transient error
- **THEN** the system retries using helper.Retry() with tccommon.ReadRetryTimeout
- **AND** wraps the error with tccommon.RetryError() if all retries are exhausted

#### Scenario: Delete operation retries on transient failure
- **WHEN** DeleteMultiPathGatewayLine fails with a transient error
- **THEN** the system retries using helper.Retry() with tccommon.ReadRetryTimeout
- **AND** wraps the error with tccommon.RetryError() if all retries are exhausted

### Requirement: Unit test coverage
The system SHALL provide unit tests in `resource_tc_teo_multi_path_gateway_line_test.go` using gomonkey to mock cloud API calls, covering Create, Read, Update, and Delete operations.

#### Scenario: Unit tests exist for CRUD operations
- **WHEN** running `go test` with `-gcflags=all=-l` on the test file
- **THEN** all unit tests for Create, Read, Update, and Delete pass successfully
