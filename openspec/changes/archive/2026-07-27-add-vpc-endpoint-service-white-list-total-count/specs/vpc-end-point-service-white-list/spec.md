## ADDED Requirements

### Requirement: Total Count Field
The `tencentcloud_vpc_end_point_service_white_list` resource SHALL expose a computed `total_count` field representing the total number of matching white list records returned by the `DescribeVpcEndPointServiceWhiteList` API.

#### Scenario: Schema definition of total_count
- **WHEN** the resource schema is defined
- **THEN** `total_count` is declared as `schema.TypeInt`, `Computed: true`, with a description indicating it is the total count of matching white list records

#### Scenario: Read populates total_count from API response
- **WHEN** the resource Read function queries the `DescribeVpcEndPointServiceWhiteList` API and the response `TotalCount` is non-nil
- **THEN** the `total_count` field is set in the Terraform state from the response `TotalCount` value

#### Scenario: Read handles nil TotalCount safely
- **WHEN** the response `TotalCount` is nil
- **THEN** the Read function does not panic and leaves `total_count` unset in state

### Requirement: Service Layer Returns Total Count
The `VpcService.DescribeVpcEndPointServiceWhiteListById` method SHALL return the `TotalCount` from the `DescribeVpcEndPointServiceWhiteList` first-page response alongside the existing `*vpc.VpcEndPointServiceUser` record, so the resource Read function can populate `total_count` without an extra API call.

#### Scenario: Service layer captures TotalCount from first page
- **WHEN** the service layer paginates the `DescribeVpcEndPointServiceWhiteList` API
- **THEN** the `TotalCount` value from the first page (`Offset == 0`) is returned to the caller

#### Scenario: Backward compatible return signature
- **WHEN** the method signature is updated to include `totalCount *uint64`
- **THEN** all callers (including the pls resource Read function and the duplicated dcg copy) are updated to accept and use the new return value

### Requirement: Unit Test Coverage
The new `total_count` field SHALL be covered by unit tests using gomonkey mocks, without relying on the Terraform acceptance test suite.

#### Scenario: Unit test verifies total_count is set
- **WHEN** the `DescribeVpcEndPointServiceWhiteList` API is mocked to return a response with `TotalCount` set
- **THEN** the resource Read function sets `total_count` in state to the mocked value

#### Scenario: Unit test verifies nil TotalCount safety
- **WHEN** the mocked API response has `TotalCount` nil
- **THEN** the Read function completes without error and does not set `total_count`
