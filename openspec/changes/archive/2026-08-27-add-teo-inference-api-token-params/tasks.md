## 1. Schema Changes

- [x] 1.1 Add `offset` field (Optional, TypeInt, Default: 0) to the resource schema in `resource_tc_teo_inference_api_token.go`
- [x] 1.2 Add `limit` field (Optional, TypeInt, Default: 20) to the resource schema in `resource_tc_teo_inference_api_token.go`
- [x] 1.3 Add `total_count` field (Computed, TypeInt) to the resource schema in `resource_tc_teo_inference_api_token.go`

## 2. Read Method Changes

- [x] 2.1 Update the Read method to pass `Offset` and `Limit` from schema to the `DescribeInferenceAPITokens` request
- [x] 2.2 Update the Read method to set `TotalCount` from the `DescribeInferenceAPITokens` response into the resource state

## 3. Documentation and Tests

- [x] 3.1 Update the resource .md documentation file to include examples with `offset`, `limit`, and `total_count`
- [x] 3.2 Add unit test cases for the new parameters in `resource_tc_teo_inference_api_token_test.go`

## 4. Verification

- [x] 4.1 Verify the code compiles correctly with `go build`
- [x] 4.2 Run `make doc` to regenerate website documentation