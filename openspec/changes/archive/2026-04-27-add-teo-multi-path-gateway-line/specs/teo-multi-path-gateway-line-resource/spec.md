## ADDED Requirements

### Requirement: Resource schema defines multi-path gateway line fields
The resource `tencentcloud_teo_multi_path_gateway_line` SHALL define the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID
- `gateway_id` (Required, ForceNew, TypeString): 多通道安全网关 ID
- `line_type` (Required, TypeString): 线路类型，取值为 direct/proxy/custom
- `line_address` (Required, TypeString): 线路地址，格式为 ip:port
- `proxy_id` (Optional, TypeString): 四层代理实例 ID，line_type 为 proxy 时必传
- `rule_id` (Optional, TypeString): 转发规则 ID，line_type 为 proxy 时必传
- `line_id` (Computed, TypeString): 线路 ID，由云 API 返回

The resource composite ID SHALL be `{zone_id}#{gateway_id}#{line_id}` separated by `#`.

#### Scenario: Schema fields are correctly defined
- **WHEN** the resource schema is accessed
- **THEN** all fields SHALL have correct types, required/computed settings, and descriptions

#### Scenario: Composite ID format
- **WHEN** a resource is created successfully
- **THEN** the resource ID SHALL be in the format `{zone_id}#{gateway_id}#{line_id}`

### Requirement: Create operation creates a new multi-path gateway line
The Create operation SHALL call `CreateMultiPathGatewayLine` API with zone_id, gateway_id, line_type, line_address, proxy_id (when applicable), and rule_id (when applicable). Upon success, it SHALL set the composite ID from the returned LineId and call Read to sync state.

#### Scenario: Successful creation with proxy line type
- **WHEN** a resource is created with line_type="proxy", gateway_id, zone_id, line_address, proxy_id, and rule_id
- **THEN** the Create function SHALL call CreateMultiPathGatewayLine API with all parameters
- **THEN** the resource ID SHALL be set to `{zone_id}#{gateway_id}#{line_id}`
- **THEN** Read SHALL be called to sync the state

#### Scenario: Successful creation with custom line type
- **WHEN** a resource is created with line_type="custom", gateway_id, zone_id, and line_address
- **THEN** the Create function SHALL call CreateMultiPathGatewayLine API without proxy_id and rule_id
- **THEN** the resource ID SHALL be set to `{zone_id}#{gateway_id}#{line_id}`

#### Scenario: Create API call with retry
- **WHEN** the CreateMultiPathGatewayLine API call fails with a retryable error
- **THEN** the operation SHALL retry using `resource.Retry` with `tccommon.WriteRetryTimeout`

### Requirement: Read operation queries multi-path gateway line details
The Read operation SHALL parse the composite ID to extract zone_id, gateway_id, and line_id, then call `DescribeMultiPathGatewayLine` API. If the resource is not found, it SHALL set the resource ID to empty string.

#### Scenario: Successful read
- **WHEN** Read is called with a valid composite ID
- **THEN** the function SHALL parse zone_id, gateway_id, and line_id from the ID
- **THEN** it SHALL call DescribeMultiPathGatewayLine with ZoneId, GatewayId, and LineId
- **THEN** it SHALL set line_id, line_type, line_address, proxy_id, and rule_id from the response

#### Scenario: Resource not found
- **WHEN** DescribeMultiPathGatewayLine returns no data (resource deleted externally)
- **THEN** the resource ID SHALL be set to empty string
- **THEN** no error SHALL be returned

### Requirement: Update operation modifies multi-path gateway line
The Update operation SHALL call `ModifyMultiPathGatewayLine` API when any mutable field (line_type, line_address, proxy_id, rule_id) has changed. After a successful update, it SHALL call Read to sync state.

#### Scenario: Update line_address
- **WHEN** line_address is changed
- **THEN** the Update function SHALL call ModifyMultiPathGatewayLine with the new line_address

#### Scenario: Update proxy_id and rule_id
- **WHEN** proxy_id or rule_id is changed
- **THEN** the Update function SHALL call ModifyMultiPathGatewayLine with the updated proxy_id and/or rule_id

#### Scenario: No changes detected
- **WHEN** no mutable fields have changed
- **THEN** the Update function SHALL skip the API call and directly call Read

### Requirement: Delete operation removes a multi-path gateway line
The Delete operation SHALL call `DeleteMultiPathGatewayLine` API with zone_id, gateway_id, and line_id parsed from the composite ID.

#### Scenario: Successful deletion
- **WHEN** Delete is called with a valid composite ID
- **THEN** the function SHALL call DeleteMultiPathGatewayLine with ZoneId, GatewayId, and LineId
- **THEN** no error SHALL be returned on success

#### Scenario: Delete API call with retry
- **WHEN** the DeleteMultiPathGatewayLine API call fails with a retryable error
- **THEN** the operation SHALL retry using `resource.Retry` with `tccommon.WriteRetryTimeout`

### Requirement: Resource supports Terraform Import
The resource SHALL support importing existing multi-path gateway lines via `terraform import` using the composite ID format `{zone_id}#{gateway_id}#{line_id}`.

#### Scenario: Import existing resource
- **WHEN** a user runs `terraform import tencentcloud_teo_multi_path_gateway_line.example zone-xxx#gw-xxx#line-2`
- **THEN** the resource SHALL be imported with the correct state by calling Read

### Requirement: Resource registration in provider
The resource SHALL be registered in `provider.go` under `ResourcesMap` with key `tencentcloud_teo_multi_path_gateway_line`, and referenced in `provider.md` under the TEO section.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_multi_path_gateway_line` SHALL be available for use in Terraform configurations

### Requirement: Unit tests using gomonkey mock
The resource SHALL have unit tests in `resource_tc_teo_multi_path_gateway_line_test.go` using gomonkey to mock cloud API calls, covering Create, Read, Update, and Delete operations.

#### Scenario: Unit test coverage
- **WHEN** unit tests are run with `go test`
- **THEN** Create, Read, Update, and Delete operations SHALL be tested with mocked API responses

### Requirement: Markdown example file
The resource SHALL have a `.md` example file at `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line.md` with usage example and import instructions.

#### Scenario: Documentation file exists
- **WHEN** the resource is added
- **THEN** a `.md` example file SHALL exist with a one-line description, Example Usage HCL block, and Import section
