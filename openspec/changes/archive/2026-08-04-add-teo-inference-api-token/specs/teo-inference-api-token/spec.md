# teo-inference-api-token Specification

## Purpose

定义 TEO 服务创建推理 API Token operation 资源的行为规范。

## ADDED Requirements

### Requirement: Resource Schema Definition

Resource 必须支持以下参数：

**Input Parameters:**
- `zone_id` (String, Required): 站点 ID
- `name` (String, Required): 推理 API Token 的名称，长度限制不超过 30 个字符

**Output Attributes:**
- `token_id` (String, Computed): 推理 API Token ID
- `content` (String, Computed, Sensitive): 推理 API Token 内容

所有输入参数必须设置为 ForceNew，即参数变更时重新创建资源。`content` 属性必须标记为 Sensitive。

#### Scenario: Successful resource creation with required parameters

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-12345678"
  name    = "my-token"
}
```

- **WHEN** user provides valid `zone_id` and `name`
- **THEN** the resource creates an inference API token for the specified zone
- **AND** sets the resource ID to the `zone_id`
- **AND** outputs `token_id` and `content` as computed attributes

#### Scenario: Handle missing required parameters

- **WHEN** user omits `zone_id` parameter
- **THEN** the resource returns validation error indicating `zone_id` is required

- **WHEN** user omits `name` parameter
- **THEN** the resource returns validation error indicating `name` is required

#### Scenario: Handle invalid parameter values

- **WHEN** user provides empty `zone_id`
- **THEN** the resource returns validation error indicating `zone_id` cannot be empty

- **WHEN** user provides empty `name`
- **THEN** the resource returns validation error indicating `name` cannot be empty

### Requirement: Create Operation

Resource 必须调用 `CreateInferenceAPIToken` API 接口创建推理 API Token。

#### Scenario: Successful API call

- **WHEN** the resource is created with valid parameters
- **THEN** it calls `CreateInferenceAPITokenWithContext` with the provided `zone_id` and `name`
- **AND** uses `tccommon.WriteRetryTimeout` for retry logic on transient failures
- **AND** sets the resource ID to `zone_id`
- **AND** sets `token_id` and `content` from the API response
- **AND** calls Read operation to update state
- **AND** logs API calls using standard logging patterns with `tccommon.LogElapsed`
- **AND** performs consistency check with `tccommon.InconsistentCheck`

#### Scenario: API returns empty response

- **WHEN** the API returns nil response or nil Response field
- **THEN** the resource returns a `NonRetryableError` indicating the API returned an empty response
- **AND** does NOT set the resource ID to avoid writing empty state

#### Scenario: API rate limiting handling

- **WHEN** API returns rate limiting error (429)
- **THEN** the resource retries the request according to `tccommon.WriteRetryTimeout`
- **AND** logs retry attempts for debugging

#### Scenario: API authentication failure

- **WHEN** API returns authentication error due to invalid credentials
- **THEN** the resource returns the original API error to help with troubleshooting
- **AND** does not retry (authentication errors are not transient)

#### Scenario: Invalid zone_id or name

- **WHEN** API returns error indicating invalid `zone_id` or `name`
- **THEN** the resource returns the original API error message
- **AND** helps user identify which parameter is invalid

#### Scenario: Network timeout

- **WHEN** API call times out due to network issues
- **THEN** the resource retries according to retry policy
- **AND** returns timeout error after exhausting retries

### Requirement: Read Operation

Resource 必须提供空实现的 Read 操作（operation 类型资源特性）。

#### Scenario: Read operation returns immediately

- **WHEN** the resource is read
- **THEN** the Read function returns nil without any API calls
- **AND** does not update any state attributes

### Requirement: Delete Operation

Resource 必须提供空实现的 Delete 操作（operation 类型资源特性）。

#### Scenario: Delete operation returns immediately

- **WHEN** the resource is deleted
- **THEN** the Delete function returns nil without any API calls
- **AND** removes the resource from Terraform state

### Requirement: No Update Operation

Resource 必须不支持 Update 操作，所有参数变更触发重新创建。

#### Scenario: Parameter change triggers recreation

- **WHEN** user changes `zone_id` or `name` parameter
- **THEN** Terraform destroys the existing resource instance
- **AND** creates a new resource instance with the updated parameters
- **AND** calls CreateInferenceAPIToken API again with new parameters

#### Scenario: No direct update method

- **WHEN** user attempts to use Terraform's `terraform apply` to update parameters
- **THEN** the resource performs destroy+create cycle instead of update

### Requirement: Error Handling and Logging

Resource 必须正确处理错误并提供充分的日志信息。

#### Scenario: API error logging

- **WHEN** an API call fails
- **THEN** the resource logs the error with request and response details
- **AND** includes the API action name for easier debugging
- **AND** returns the error to user

#### Scenario: Operation duration logging

- **WHEN** any CRUD operation is executed
- **THEN** the resource logs the operation duration using `tccommon.LogElapsed`
- **AND** includes the resource name and operation type in log message

#### Scenario: Consistency check

- **WHEN** any CRUD operation is executed
- **THEN** the resource performs consistency check using `tccommon.InconsistentCheck`
- **AND** detects and reports state inconsistencies

### Requirement: Unit Testing

Resource 必须提供完整的单元测试覆盖，使用 gomonkey 进行 mock。

#### Scenario: Test successful creation

- **WHEN** unit test runs with valid mock API response
- **THEN** the test verifies CreateInferenceAPIToken is called with correct parameters
- **AND** the test verifies resource ID is set correctly
- **AND** the test verifies `token_id` and `content` are set from response
- **AND** the test verifies retry logic on transient errors

#### Scenario: Test API error handling

- **WHEN** unit test simulates API error
- **THEN** the test verifies error is returned correctly
- **AND** the test verifies error message is preserved

#### Scenario: Test missing required parameters

- **WHEN** unit test runs with missing required parameters
- **THEN** the test verifies validation error is returned
- **AND** the test verifies error message indicates which parameter is missing

#### Scenario: Test empty API response

- **WHEN** unit test simulates nil API response
- **THEN** the test verifies NonRetryableError is returned
- **AND** the test verifies resource ID is not set

### Requirement: Documentation

Resource 必须提供清晰的使用文档。

#### Scenario: Documentation includes parameters

- **WHEN** users read the resource documentation
- **THEN** they see clear descriptions of `zone_id` and `name` parameters
- **AND** they see descriptions of `token_id` and `content` computed attributes
- **AND** they see usage examples showing how to create the resource
- **AND** they understand this is an operation resource (create-only)

#### Scenario: Documentation includes notes

- **WHEN** users read the resource documentation
- **THEN** they see notes about the operation resource characteristics
- **AND** they understand that Read and Delete operations are empty implementations
- **AND** they understand that parameter changes trigger recreation
- **AND** they are warned that `content` (the token) is sensitive and stored in state