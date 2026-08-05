## ADDED Requirements

### Requirement: TEO inference API token creation

The system SHALL support creating a TEO inference API token via the `CreateInferenceAPIToken` cloud API, and return the `token_id` and `content` in the Terraform state.

#### Scenario: Successful token creation

- **WHEN** user provides a valid `zone_id` and `name` (≤30 characters) to create a `tencentcloud_teo_inference_api_token` resource
- **THEN** the system calls `CreateInferenceAPIToken` and stores the returned `token_id` as the Terraform resource ID, and `content` as a sensitive computed attribute

#### Scenario: Duplicate token name

- **WHEN** user attempts to create a token with a name that already exists in the same zone
- **THEN** the system returns an error from the cloud API

### Requirement: TEO inference API token read

The system SHALL support reading a TEO inference API token via the `DescribeInferenceAPITokens` cloud API, matching by `token_id` in the returned list.

#### Scenario: Successful token read

- **WHEN** Terraform reads the state of an existing `tencentcloud_teo_inference_api_token` resource
- **THEN** the system calls `DescribeInferenceAPITokens` with the `zone_id` from state, iterates the `Tokens` list to find the matching `token_id`, and updates all computable attributes in state

#### Scenario: Token not found (deleted externally)

- **WHEN** Terraform reads the state of a `tencentcloud_teo_inference_api_token` resource whose token has been deleted outside of Terraform
- **THEN** the system logs a warning with the resource ID, sets the resource ID to empty string, and returns nil (no error)

### Requirement: TEO inference API token deletion

The system SHALL support deleting a TEO inference API token via the `DeleteInferenceAPIToken` cloud API.

#### Scenario: Successful token deletion

- **WHEN** user runs `terraform destroy` on a `tencentcloud_teo_inference_api_token` resource
- **THEN** the system calls `DeleteInferenceAPIToken` with the `zone_id` and `token_id`, and removes the resource from state

#### Scenario: Token already deleted

- **WHEN** user attempts to delete a token that has already been deleted externally
- **THEN** the system returns an error from the cloud API

### Requirement: Immutable fields

The system SHALL treat all user-provided fields as `ForceNew` since the cloud API does not provide an Update/Modify interface.

#### Scenario: Field change triggers recreation

- **WHEN** user modifies the `name` or `zone_id` of an existing `tencentcloud_teo_inference_api_token` resource
- **THEN** Terraform destroys the existing resource and creates a new one with the updated values