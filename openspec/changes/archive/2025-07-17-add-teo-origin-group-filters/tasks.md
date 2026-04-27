## 1. Schema Definition

- [x] 1.1 Add `filters` parameter to the `tencentcloud_teo_origin_group` resource schema in `tencentcloud/services/teo/resource_tc_teo_origin_group.go` as TypeList with Optional+Computed, containing nested schema with `name` (TypeString, Required), `values` (TypeList of TypeString, Required), and `fuzzy` (TypeBool, Optional) sub-fields

## 2. Service Layer Update

- [x] 2.1 Update `DescribeTeoOriginGroupById` method signature in `tencentcloud/services/teo/service_tencentcloud_teo.go` to accept `zoneId` and `filters` parameters
- [x] 2.2 Implement passing `ZoneId` and `Filters` to the `DescribeOriginGroup` API request in the updated service method
- [x] 2.3 When custom `filters` are provided, use them instead of the default `origin-group-id` filter; when no custom filters are provided, use the default `origin-group-id` filter

## 3. CRUD Function Update

- [x] 3.1 Update `resourceTencentCloudTeoOriginGroupRead` function to read `filters` from the resource schema and convert them to `[]*teo.AdvancedFilter`
- [x] 3.2 Pass `zoneId` and converted `filters` to the updated `DescribeTeoOriginGroupById` service method
- [x] 3.3 Set the `filters` computed value in the read function after the API response

## 4. Unit Tests

- [x] 4.1 Add unit test cases for the `filters` parameter in `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go` using gomonkey mock approach
- [x] 4.2 Test scenario: read with no filters (default behavior)
- [x] 4.3 Test scenario: read with custom filters including name, values, and fuzzy fields
- [x] 4.4 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/teo/resource_tc_teo_origin_group.md` to add `filters` parameter documentation with example usage
