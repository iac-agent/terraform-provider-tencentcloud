## Context

TEO (Tencent EdgeOne) provides the `CreateInferenceAPIToken` API (package `teo/v20220901`) to create inference API tokens for EdgeOne zones. This is a synchronous API that accepts `ZoneId` and `Name` and returns `TokenId` and `Content`.

The Terraform provider currently has no resource to manage inference API tokens. This change adds a new `RESOURCE_KIND_OPERATION` resource that wraps the `CreateInferenceAPIToken` API.

Operation resources in this provider follow a specific pattern: they only implement the Create method (calling the cloud API), while Read and Delete are no-ops. All parameters use `ForceNew: true`. The same pattern is used by existing TEO operation resources like `tencentcloud_teo_create_cls_index_operation`, `tencentcloud_teo_check_cname_status_operation`, `tencentcloud_teo_identify_zone_operation`, and `tencentcloud_teo_l7_acc_rule_priority_operation`.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_teo_inference_api_token` that creates an inference API token via the `CreateInferenceAPIToken` API
- Input parameters: `zone_id` (Required, ForceNew), `name` (Required, ForceNew)
- Output (Computed) attributes: `token_id`, `content`
- Follow the established operation resource pattern (Create only, Read/Delete no-ops)
- Register the resource in the TEO provider section of `provider.go` and `provider.md`
- Provide unit tests with mocked API calls using gomonkey
- Provide documentation following the `.md` file format

**Non-Goals:**
- No Read/Update/Delete operations (operation resources are create-only)
- No import support (operation resources do not support import)
- No state management beyond the Terraform lifecycle

## Decisions

### Decision 1: Follow existing operation resource pattern
**Choice**: Create the resource following the exact pattern of `resource_tc_teo_create_cls_index_operation.go` and `resource_tc_teo_check_cname_status_operation.go`.

**Rationale**: All TEO operation resources follow the same structure. Consistency reduces maintenance burden and makes the code predictable for other developers.

**Alternatives considered**: None - this is the established pattern.

### Decision 2: Use `zone_id` as the resource ID (SetId)
**Choice**: Set `d.SetId(zoneId)` after successful API call.

**Rationale**: The operation resource pattern uses the primary identifier (zone_id) as the resource ID. This is consistent with other TEO operation resources like `create_cls_index_operation` and `check_cname_status_operation`.

### Decision 3: Set output attributes from API response
**Choice**: After the API call succeeds, set `token_id` and `content` as Computed attributes from the response.

**Rationale**: The API returns `TokenId` and `Content` in the response, and users need to access these values. Setting them as Computed attributes allows referencing them via `resource.tencentcloud_teo_inference_api_token.example.token_id` and `resource.tencentcloud_teo_inference_api_token.example.content`.

### Decision 4: Use gomonkey for unit tests
**Choice**: Mock the API call using gomonkey instead of Terraform's acceptance test framework.

**Rationale**: For new operation resources, the project guidelines specify using gomonkey to mock cloud API calls, enabling unit tests that don't require actual cloud credentials.

### Decision 5: No async polling needed
**Choice**: This is a synchronous API, no polling is required after calling `CreateInferenceAPIToken`.

**Rationale**: The API documentation does not indicate this is an asynchronous operation, and the response immediately returns the created token information.

## Risks / Trade-offs

- **[Risk] Token content is sensitive**: The `content` field contains the actual API token, which will be stored in Terraform state in plaintext. → **Mitigation**: Document this in the resource description and mark `content` as `Sensitive: true` in the schema. This is consistent with how other token/secret resources handle sensitive data.
- **[Risk] No way to delete tokens via Terraform**: The cloud API does not provide a `DeleteInferenceAPIToken` or `DescribeInferenceAPIToken` API, so tokens can only be created, not managed after creation. → **Mitigation**: This is the nature of an operation resource - clearly document this limitation.