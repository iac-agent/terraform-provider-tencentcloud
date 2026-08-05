## Requirements

### Requirement: Resource schema definition
The `tencentcloud_teo_inference_api_token_v7` resource SHALL have the following schema fields:

- `zone_id` (TypeString, Required, ForceNew): The zone ID where the inference API token belongs to.
- `name` (TypeString, Required, ForceNew): The name of the inference API token, with a maximum length of 30 characters.
- `token_id` (TypeString, Computed): The unique ID of the inference API token, returned by the Create API.
- `content` (TypeString, Computed, Sensitive): The content of the inference API token, returned by the Create API. This is sensitive information.

#### Scenario: Schema validation on create
- **WHEN** a user creates a `tencentcloud_teo_inference_api_token_v7` resource with `zone_id` and `name`
- **THEN** the resource SHALL be created successfully and `token_id` and `content` SHALL be computed from the API response

#### Scenario: Schema validation with missing required fields
- **WHEN** a user creates a `tencentcloud_teo_inference_api_token_v7` resource without `zone_id` or `name`
- **THEN** Terraform SHALL report a validation error before making any API call

### Requirement: Create inference API token
The resource SHALL call `CreateInferenceAPIToken` API to create a new inference API token.

#### Scenario: Successful token creation
- **WHEN** the Create method is called with valid `zone_id` and `name`
- **THEN** the system SHALL call `CreateInferenceAPIToken` API with the provided parameters
- **AND** the system SHALL set the resource ID to the returned `TokenId`
- **AND** the system SHALL set `content` from the API response

#### Scenario: Create API returns empty response
- **WHEN** the Create method calls `CreateInferenceAPIToken` API and receives a response with nil or empty `TokenId`
- **THEN** the system SHALL return a `NonRetryableError` to prevent writing an empty resource ID

#### Scenario: Create API call fails
- **WHEN** the Create method calls `CreateInferenceAPIToken` API and the API returns an error
- **THEN** the system SHALL retry the API call using `resource.Retry` with `tccommon.ReadRetryTimeout`
- **AND** the system SHALL wrap the error with `tccommon.RetryError` for retryable errors

### Requirement: Read inference API token
The resource SHALL call `DescribeInferenceAPITokens` API to query the token's current state.

#### Scenario: Token exists
- **WHEN** the Read method is called with a valid `TokenId` in the resource ID
- **THEN** the system SHALL call `DescribeInferenceAPITokens` API with the `zone_id` from state
- **AND** the system SHALL traverse the returned `Tokens` list to find the matching `TokenId`
- **AND** the system SHALL set `zone_id`, `name`, `token_id`, and `content` (if non-nil) in the state

#### Scenario: Token does not exist
- **WHEN** the Read method calls `DescribeInferenceAPITokens` and the target `TokenId` is not found in the response
- **THEN** the system SHALL log a warning message with the resource ID
- **AND** the system SHALL call `d.SetId("")` to remove the resource from state

#### Scenario: Read API returns empty response
- **WHEN** the Read method calls `DescribeInferenceAPITokens` and receives a nil or empty response
- **THEN** the system SHALL log a warning message with the resource ID
- **AND** the system SHALL call `d.SetId("")` to remove the resource from state

### Requirement: Delete inference API token
The resource SHALL call `DeleteInferenceAPIToken` API to delete the token.

#### Scenario: Successful token deletion
- **WHEN** the Delete method is called with a valid `TokenId`
- **THEN** the system SHALL call `DeleteInferenceAPIToken` API with `zone_id` and `token_id`
- **AND** the system SHALL succeed even if the API returns an error (resource may already be deleted externally)

#### Scenario: Delete API call fails
- **WHEN** the Delete method calls `DeleteInferenceAPIToken` API and the API returns an error
- **THEN** the system SHALL retry the API call using `resource.Retry` with `tccommon.ReadRetryTimeout`
- **AND** the system SHALL wrap the error with `tccommon.RetryError` for retryable errors

### Requirement: Resource import support
The resource SHALL support Terraform import via `terraform import tencentcloud_teo_inference_api_token_v7.foo <token_id>`.

#### Scenario: Import existing token
- **WHEN** a user runs `terraform import tencentcloud_teo_inference_api_token_v7.foo <token_id>`
- **THEN** the system SHALL use `schema.ImportStatePassthrough` to set the resource ID to the provided token_id
- **AND** the subsequent Read operation SHALL populate the remaining state fields from the API

### Requirement: No update operation
The resource SHALL NOT support update operations since the API only provides Create, Read, and Delete.

#### Scenario: Attempt to modify resource
- **WHEN** a user modifies `zone_id` or `name` in an existing resource configuration
- **THEN** Terraform SHALL detect that the resource must be destroyed and recreated (ForceNew behavior)
- **AND** the system SHALL NOT call any update API