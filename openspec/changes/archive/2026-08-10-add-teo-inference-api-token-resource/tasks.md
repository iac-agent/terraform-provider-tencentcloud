## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_teo_inference_api_token_v9.go` with resource schema definition containing `zone_id`, `name`, `token_id`, and `content` fields
- [x] 1.2 Implement Create function: call `CreateInferenceAPIToken` with retry, validate response, set resource ID to `TokenId`
- [x] 1.3 Implement Read function: call `DescribeInferenceAPITokens` with `ZoneId` and `Limit=100`, iterate `Tokens` to find matching token by `TokenId`, handle not-found case
- [x] 1.4 Implement Delete function: call `DeleteInferenceAPIToken` with `ZoneId` and `TokenId` with retry

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_inference_api_token_v9` resource in `tencentcloud/provider.go`

## 3. Documentation

- [x] 3.1 Create `resource_tc_teo_inference_api_token_v9.md` with Example Usage section following gendoc format

## 4. Testing

- [x] 4.1 Create `resource_tc_teo_inference_api_token_v9_test.go` with unit tests using gomonkey mocks for CRUD operations
