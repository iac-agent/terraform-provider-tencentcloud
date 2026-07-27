## 1. Service Layer Changes

- [x] 1.1 Update `DescribeVpcEndPointServiceWhiteListById` in `tencentcloud/services/vpc/service_tencentcloud_vpc.go` to return an additional `totalCount *uint64` value captured from the first-page (`Offset == 0`) `DescribeVpcEndPointServiceWhiteList` response `TotalCount`.
- [x] 1.2 Update the duplicated `DescribeVpcEndPointServiceWhiteListById` in `tencentcloud/services/dcg/service_tencentcloud_vpc.go` with the identical signature/behavior change to keep both copies consistent.

## 2. Resource Schema & CRUD Changes

- [x] 2.1 Add `total_count` field to the `ResourceTencentCloudVpcEndPointServiceWhiteList()` schema in `tencentcloud/services/pls/resource_tc_vpc_end_point_service_white_list.go`: `schema.TypeInt`, `Computed: true`, description "Total count of matching white list records.".
- [x] 2.2 Update `resourceTencentCloudVpcEndPointServiceWhiteListRead` to receive the new `totalCount` return value from the service layer and set `total_count` in state with a nil guard (e.g. `if totalCount != nil { _ = d.Set("total_count", int(*totalCount)) }`).

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/pls/resource_tc_vpc_end_point_service_white_list.md` to document the new `total_count` computed field (run `make doc` during finalize to regenerate `website/docs/` markdown; do not edit `website/` files manually).

## 4. Unit Tests

- [x] 4.1 Add gomonkey-based unit tests in `tencentcloud/services/pls/resource_tc_vpc_end_point_service_white_list_test.go` covering: (a) `total_count` is set when the mocked API response has a non-nil `TotalCount`, and (b) Read does not panic and leaves `total_count` unset when `TotalCount` is nil.
- [x] 4.2 Run `go test -gcflags=all=-l` on the affected test files to confirm the new unit tests pass.

## 5. Verification

- [x] 5.1 Verify the change compiles and the service layer signature updates do not break other callers (grep for `DescribeVpcEndPointServiceWhiteListById` usages).
- [x] 5.2 Verify backward compatibility: existing schema fields are unchanged; `total_count` is Computed-only.
