## Context

The `tencentcloud_teo_inference_api_token` resource manages TEO Inference API Tokens via the following cloud APIs:
- `CreateInferenceAPIToken` — creates a new token
- `DeleteInferenceAPIToken` — deletes an existing token
- `DescribeInferenceAPITokens` — queries token list (no `ModifyInferenceAPIToken` exists, so this is a CRD-only resource)

The `DescribeInferenceAPITokens` API supports pagination via `Offset` and `Limit` input parameters and returns `TotalCount` in the response. The resource currently does not expose these fields, limiting the user's ability to control pagination and see the total token count.

The vendor SDK (`tencentcloud-sdk-go/tencentcloud/teo/v20220901`) already contains all required fields in `DescribeInferenceAPITokensRequestParams` and `DescribeInferenceAPITokensResponseParams`.

## Goals / Non-Goals

**Goals:**
- Add `Offset` (Optional, int64, default 0) to the resource schema, passed to `DescribeInferenceAPITokens` request
- Add `Limit` (Optional, int64, default 20, max 100) to the resource schema, passed to `DescribeInferenceAPITokens` request
- Add `TotalCount` (Computed, int64) to the resource schema, populated from `DescribeInferenceAPITokens` response

**Non-Goals:**
- No changes to Create or Delete operations
- No changes to the `InferenceAPIToken` list structure
- No new API endpoints or SDK upgrades required

## Decisions

1. **Offset and Limit as Optional with defaults**: These are set as `Optional: true` with `Default: 0` (Offset) and `Default: 20` (Limit), matching the API defaults. Users can override them in their Terraform configuration.

2. **TotalCount as Computed**: `TotalCount` is read-only from the API response, so it is set as `Computed: true` without `Optional`. Users cannot set it; it is populated by the Read method.

3. **No new pagination loop logic**: Since this is a CRD resource (not a datasource), the Read method reads a single resource by ID. The Offset/Limit parameters are passed through but the resource's Read logic does not implement pagination looping — that pattern is for datasources only.

## Risks / Trade-offs

- **Risk**: If the resource is read by a non-unique filter (e.g., by name), pagination may cause the wrong token to be selected. → **Mitigation**: Ensure the Read method always queries by the resource's unique ID (TokenId), not by paginated listing.
- **Risk**: Adding `TotalCount` as a schema field may cause plan diffs if the value changes between reads. → **Mitigation**: `TotalCount` is purely informational and changes are expected; users should not rely on it for resource identity.