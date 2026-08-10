## Requirements

### Requirement: Resource schema definition
The `tencentcloud_teo_inference_api_token_v9` resource MUST define the following schema fields:
- `zone_id` (TypeString, Required, ForceNew): The site ID for the inference API token.
- `name` (TypeString, Required, ForceNew): The name of the inference API token, limited to 30 characters.
- `token_id` (TypeString, Computed): The inference API token ID, returned by the cloud API after creation.
- `content` (TypeString, Computed, Sensitive): The inference API token content, returned only once during creation.

#### Scenario: Schema validation
- **WHEN** a user defines a `tencentcloud_teo_inference_api_token_v9` resource with `zone_id` and `name`
- **THEN** the Terraform schema MUST accept the configuration and validate the required fields

### Requirement: Create inference API token
The resource MUST call `CreateInferenceAPIToken` API to create a new inference API token. The Create function MUST:
- Accept `zone_id` and `name` as input parameters
- Call `CreateInferenceAPIToken` with retry logic using `tccommon.ReadRetryTimeout`
- On success, set `d.SetId(tokenId)` with the returned `TokenId`
- Set `token_id` and `content` from the API response
- Check that the API response is not nil and `TokenId` is not empty before setting the resource ID

#### Scenario: Successful token creation
- **WHEN** a user creates a `tencentcloud_teo_inference_api_token_v9` resource with valid `zone_id` and `name`
- **THEN** the system SHALL call `CreateInferenceAPIToken` and set the resource ID to the returned `TokenId`

#### Scenario: Create API returns empty response
- **WHEN** `CreateInferenceAPIToken` returns a nil response or empty `TokenId`
- **THEN** the system SHALL return a `NonRetryableError` and NOT set an empty resource ID

### Requirement: Read inference API token
The resource MUST call `DescribeInferenceAPITokens` API to query the token details. The Read function MUST:
- Use `zone_id` and `token_id` (from resource ID) to query
- Call `DescribeInferenceAPITokens` with `Limit=100` (maximum) and iterate through `Tokens` to find the matching token
- If the token is found, set `zone_id`, `name`, `token_id` in state
- If the token is not found (response is nil or Tokens list is empty), log the id and call `d.SetId("")`
- NOT overwrite `content` in state (Content is only returned once during creation)

#### Scenario: Token exists
- **WHEN** the Read function is called for an existing token
- **THEN** the system SHALL find the token in the `DescribeInferenceAPITokens` response and update the state with `zone_id`, `name`, and `token_id`

#### Scenario: Token not found
- **WHEN** the Read function is called and `DescribeInferenceAPITokens` returns no matching token
- **THEN** the system SHALL log the id and call `d.SetId("")` to remove the resource from state

### Requirement: Delete inference API token
The resource MUST call `DeleteInferenceAPIToken` API to delete a token. The Delete function MUST:
- Accept `zone_id` from state and `token_id` from resource ID
- Call `DeleteInferenceAPIToken` with retry logic using `tccommon.ReadRetryTimeout`
- Return nil on success (standard Terraform delete behavior)

#### Scenario: Successful token deletion
- **WHEN** a user deletes a `tencentcloud_teo_inference_api_token_v9` resource
- **THEN** the system SHALL call `DeleteInferenceAPIToken` with the correct `zone_id` and `token_id`

#### Scenario: Delete API returns error
- **WHEN** `DeleteInferenceAPIToken` returns an error
- **THEN** the system SHALL wrap the error with `tccommon.RetryError` and return it for retry

### Requirement: Provider registration
The resource MUST be registered in the TencentCloud provider so users can reference it as `tencentcloud_teo_inference_api_token_v9` in their Terraform configurations.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** the `tencentcloud_teo_inference_api_token_v9` resource SHALL be registered and available for use