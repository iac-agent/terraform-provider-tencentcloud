## Context

The `tencentcloud_vpc_end_point_service_white_list` resource (in `tencentcloud/services/pls/`) manages a VPC EndPointService white list entry. Its Read path delegates to `VpcService.DescribeVpcEndPointServiceWhiteListById` (in `tencentcloud/services/vpc/service_tencentcloud_vpc.go`), which paginates the `DescribeVpcEndPointServiceWhiteList` API and returns the first matching `*vpc.VpcEndPointServiceUser`.

The `DescribeVpcEndPointServiceWhiteList` API response includes a top-level `TotalCount *uint64` field ("符合条件的白名单个数"), confirmed in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312/models.go` (`DescribeVpcEndPointServiceWhiteListResponseParams`). This field is currently discarded by the service layer, so it is not available to Terraform users.

The same `DescribeVpcEndPointServiceWhiteListById` method is duplicated in `tencentcloud/services/dcg/service_tencentcloud_vpc.go` (a separate `VpcService` type in the dcg package); both copies must be kept consistent.

## Goals / Non-Goals

**Goals:**
- Expose the API `TotalCount` as a computed `total_count` field on the `tencentcloud_vpc_end_point_service_white_list` resource so users can read the total number of matching white list records.
- Keep the change backward compatible: no change to existing schema fields, no state migration, the new field is Computed-only.

**Non-Goals:**
- Do not change the pagination strategy of the service layer query (still returns the first match).
- Do not add `total_count` as an input/filter parameter; it is strictly read-only.
- Do not modify the deprecated `VpcEndpointServiceUserSet` response field usage; keep using `VpcEndPointServiceUserSet` (the non-deprecated one) where the existing code already reads `VpcEndpointServiceUserSet`. (See Risks.)

## Decisions

### Decision 1: Propagate `TotalCount` via an additional return value on `DescribeVpcEndPointServiceWhiteListById`

**Choice:** Change the method signature to return `(endPointServiceWhiteList *vpc.VpcEndPointServiceUser, totalCount *uint64, errRet error)` and update both copies (vpc and dcg packages) plus all callers.

**Alternatives considered:**
- *Alternative A: A separate `DescribeVpcEndPointServiceWhiteListCount` service method.* Rejected because it would issue a second API call for the same data already returned in the first response page, wasting a rate-limited request and risking inconsistency between the record and the count.
- *Alternative B: Store totalCount in a field on a wrapper struct.* Rejected as over-engineering for a single additional scalar; a return value is the established pattern in this codebase.

**Rationale:** `TotalCount` is already present on the first page response. Capturing it during the existing pagination loop (the count from the first page, which represents the total matching count) and returning it avoids extra API calls. Both copies of the method must be updated to keep the duplicated service files in sync.

### Decision 2: `total_count` schema field definition

`total_count` is added as `schema.TypeInt`, `Computed: true`, with a description derived from the API doc ("Total count of matching white list records."). It is read from `*uint64`; conversion uses `helper.UInt64ToInt64` or an equivalent safe cast, guarded by a nil check before `d.Set`.

### Decision 3: Use the first page's `TotalCount`

The pagination loop issues requests with increasing `Offset`. The `TotalCount` field on every page reflects the total matching count (not the page count). To avoid ambiguity, capture `TotalCount` from the first page response only (when `offset == 0`), and use that value for the returned `totalCount`.

## Risks / Trade-offs

- **[Risk] Two copies of `DescribeVpcEndPointServiceWhiteListById` (vpc and dcg) can drift.** → Mitigation: update both copies identically in this change; keep the dcg copy's behavior consistent (returning `totalCount` from its first page as well).
- **[Risk] `TotalCount` semantics across pages.** → Mitigation: capture from the first page only (`offset == 0`), where the count is the total matching count for the query.
- **[Risk] Nil `TotalCount` on some responses.** → Mitigation: guard with a nil check before `d.Set("total_count", ...)` in the resource Read function, consistent with the existing nil-guard pattern for other fields.
- **[Trade-off] Changing a public-ish service method signature touches callers.** → Accepted; only the pls resource (vpc package copy) calls it today, and the dcg copy has no current caller exercising it. Both are updated atomically.
