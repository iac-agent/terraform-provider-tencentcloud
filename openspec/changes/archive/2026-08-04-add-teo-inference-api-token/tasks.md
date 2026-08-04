## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_inference_api_token_operation.go` with schema definition (zone_id, name as required ForceNew inputs; token_id, content as computed outputs with content marked Sensitive)
- [x] 1.2 Implement Create function that calls `CreateInferenceAPIToken` API with retry, validates response is not empty, sets resource ID and computed attributes
- [x] 1.3 Implement Read function as no-op (returns nil)
- [x] 1.4 Implement Delete function as no-op (returns nil)

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_inference_api_token` in `tencentcloud/provider.go` in the TEO resource section
- [x] 2.2 Add resource entry in `tencentcloud/provider.md` in the TEO section

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_inference_api_token_operation.md` with usage example and parameter descriptions

## 4. Unit Testing

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_inference_api_token_operation_test.go` with gomonkey-based unit tests covering successful creation, API errors, empty response, and missing parameters

## 5. Verification

- [x] 5.1 Verify code compiles successfully (Go syntax check)
- [x] 5.2 Verify all files are correctly placed and named