## Why

TEO (Tencent EdgeOne) now supports inference API tokens, which allow users to create authentication tokens for EdgeOne's inference API services. The Terraform provider currently lacks support for creating these tokens, preventing users from managing them through Infrastructure as Code. This change adds a new operation resource to enable token creation via Terraform.

## What Changes

- Add a new RESOURCE_KIND_OPERATION resource `tencentcloud_teo_inference_api_token` that calls the `CreateInferenceAPIToken` cloud API
- The resource accepts `zone_id` and `name` as input parameters
- The resource outputs `token_id` and `content` as computed attributes after the API call succeeds
- As an operation resource, it only implements the Create method; Read and Delete are no-ops
- Register the new resource in the TEO provider section

## Capabilities

### New Capabilities
- `teo-inference-api-token`: Provides a Terraform resource to create an inference API token for TEO (EdgeOne) zones

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- New file: `tencentcloud/services/teo/resource_tc_teo_inference_api_token_operation.go`
- New file: `tencentcloud/services/teo/resource_tc_teo_inference_api_token_operation_test.go`
- New file: `tencentcloud/services/teo/resource_tc_teo_inference_api_token_operation.md`
- Modified: `tencentcloud/provider.go` (add resource registration)
- Modified: `tencentcloud/provider.md` (add resource entry)
- Cloud API dependency: `CreateInferenceAPIToken` from `tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already in vendor)