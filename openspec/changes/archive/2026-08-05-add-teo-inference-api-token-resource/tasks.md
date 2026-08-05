## 1. Service Layer

- [x] 1.1 Add `DescribeTeoInferenceAPITokenById` method in `tencentcloud/services/teo/service_tencentcloud_teo.go` to encapsulate `DescribeInferenceAPITokens` call and match by `TokenId`

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_inference_api_token.go` with Schema definition (`zone_id`, `name`, `token_id`, `content`, `create_time`), Create/Read/Delete functions, and import support
- [x] 2.2 Create `tencentcloud/services/teo/resource_tc_teo_inference_api_token_test.go` with unit tests using gomonkey mock for CRUD scenarios
- [x] 2.3 Create `tencentcloud/services/teo/resource_tc_teo_inference_api_token.md` with example usage documentation

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_inference_api_token` resource in `tencentcloud/provider.go`
- [x] 3.2 Register `tencentcloud_teo_inference_api_token` resource in `tencentcloud/provider.md`

## 4. Code Quality

- [x] 4.1 Run `gofmt` on all modified Go files
- [x] 4.2 Run `make doc` to generate website documentation
- [x] 4.3 Create changelog entry in `.changelog/` directory