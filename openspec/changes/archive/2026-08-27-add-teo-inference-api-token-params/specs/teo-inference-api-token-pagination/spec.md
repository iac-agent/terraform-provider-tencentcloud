## ADDED Requirements

### Requirement: Resource supports pagination parameters in DescribeInferenceAPITokens
The `tencentcloud_teo_inference_api_token` resource SHALL expose `Offset` and `Limit` as optional input parameters for the `DescribeInferenceAPITokens` API call in its Read method. The resource SHALL also expose `TotalCount` as a computed output parameter from the API response.

#### Scenario: User sets Offset and Limit
- **WHEN** user configures `offset = 10` and `limit = 50` in the `tencentcloud_teo_inference_api_token` resource
- **THEN** the Read method SHALL pass `Offset=10` and `Limit=50` to the `DescribeInferenceAPITokens` request

#### Scenario: User omits Offset and Limit
- **WHEN** user does not specify `offset` or `limit` in the resource configuration
- **THEN** the Read method SHALL use default values `Offset=0` and `Limit=20` in the `DescribeInferenceAPITokens` request

#### Scenario: TotalCount is populated from API response
- **WHEN** the `DescribeInferenceAPITokens` API returns a response with `TotalCount = 42`
- **THEN** the resource state SHALL set `total_count = 42`

#### Scenario: Backward compatibility
- **WHEN** an existing Terraform configuration that does not specify `offset`, `limit`, or `total_count` is applied
- **THEN** the resource SHALL continue to function without any breaking changes or state migration