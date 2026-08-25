## ADDED Requirements

### Requirement: Resource schema for tencentcloud_teo_inference_api_token
The `tencentcloud_teo_inference_api_token` resource SHALL expose the following schema fields:
- `zone_id` (string, Required, ForceNew): 站点 ID
- `name` (string, Required, ForceNew): 推理 API Token 名称
- `token_id` (string, Computed): 推理 API Token ID
- `content` (string, Computed): 推理 API Token 内容
- `create_time` (string, Computed): 创建时间（ISO 日期格式）

The resource id SHALL be a composite of `zone_id` and `token_id` joined by `tccommon.FILED_SP`.

#### Scenario: Schema declares required input fields
- **WHEN** the resource schema is registered
- **THEN** `zone_id` and `name` SHALL be `Required` and ForceNew

#### Scenario: Schema declares computed output fields
- **WHEN** the resource schema is registered
- **THEN** `token_id`, `content`, and `create_time` SHALL be `Computed`

### Requirement: Create a TEO inference API token
The resource Create operation SHALL call `CreateInferenceAPIToken` with `ZoneId` and `Name`, validate the returned `TokenId`/`Content` are non-empty, then set the composite id and invoke Read to populate computed fields.

#### Scenario: Successful creation
- **WHEN** `terraform apply` creates the resource with valid `zone_id` and `name`
- **THEN** the provider SHALL call `CreateInferenceAPIToken` and set `id = "<zone_id>#<token_id>"`

#### Scenario: Empty TokenId in response
- **WHEN** `CreateInferenceAPIToken` returns an empty `TokenId`
- **THEN** the provider SHALL return a `NonRetryableError` (after logging the logId and relevant fields) and SHALL NOT write an empty id

### Requirement: Read a TEO inference API token
The resource Read operation SHALL call `DescribeInferenceAPITokens` (via a service-layer helper that filters by `TokenId`), and on match populate `token_id`, `name`, `content`, `create_time`. On empty result, it SHALL first log `[CRUD] teo_inference_api_token id=<id>` then `d.SetId("")`.

#### Scenario: Token exists
- **WHEN** Read is called for an existing token
- **THEN** the provider SHALL populate `token_id`, `name`, `content`, and `create_time` from the matched list item, skipping `set` calls for any nil response fields

#### Scenario: Token not found
- **WHEN** `DescribeInferenceAPITokens` returns an empty token list (token already deleted)
- **THEN** the provider SHALL log `[CRUD] teo_inference_api_token id=<id>` and then call `d.SetId("")`

### Requirement: Delete a TEO inference API token
The resource Delete operation SHALL parse `zone_id` and `token_id` from the composite id, call `DeleteInferenceAPIToken` with `ZoneId` and `TokenId`, and return `tccommon.RetryError` on API failure.

#### Scenario: Successful deletion
- **WHEN** `terraform destroy` deletes the resource
- **THEN** the provider SHALL call `DeleteInferenceAPIToken` with the parsed `ZoneId` and `TokenId`

### Requirement: No in-place update (CRD-only)
Because the cloud API exposes no `ModifyInferenceAPIToken`, the resource SHALL NOT perform in-place updates. The Update operation SHALL place top-level mutable-intent fields (other than the id) into the `immutableArgs` array and SHALL return an error if any of them changed.

#### Scenario: Changing name triggers ForceNew
- **WHEN** the user changes `name` in the configuration
- **THEN** Terraform SHALL treat it as ForceNew (destroy + create), not an in-place update

#### Scenario: immutableArgs detection
- **WHEN** the Update operation detects a change in an `immutableArgs` field
- **THEN** the provider SHALL return an error and SHALL NOT call any non-existent modify API
