## 1. Schema Modification

- [x] 1.1 Add `offset` Optional TypeInt field to `tencentcloud_teo_l4_proxy` resource schema in `resource_tc_teo_l4_proxy.go`
- [x] 1.2 Add `limit` Optional TypeInt field to `tencentcloud_teo_l4_proxy` resource schema in `resource_tc_teo_l4_proxy.go`
- [x] 1.3 Add `total_count` Computed TypeInt field to `tencentcloud_teo_l4_proxy` resource schema in `resource_tc_teo_l4_proxy.go`

## 2. Service Layer Update

- [x] 2.1 Modify `DescribeTeoL4ProxyById` signature to accept `offset` and `limit` parameters (`*uint64`) in `service_tencentcloud_teo.go`
- [x] 2.2 Set `request.Offset` and `request.Limit` from new parameters when non-nil in `DescribeTeoL4ProxyById`
- [x] 2.3 Return `TotalCount` from API response alongside `L4Proxy` in `DescribeTeoL4ProxyById`

## 3. Read Function Update

- [x] 3.1 Update `resourceTencentCloudTeoL4ProxyRead` to read `offset`/`limit` from resource data and pass to `DescribeTeoL4ProxyById`
- [x] 3.2 Set `total_count` in Terraform state from the `TotalCount` returned by `DescribeTeoL4ProxyById`

## 4. Unit Tests

- [x] 4.1 Add test cases for `offset`/`limit` parameter passing in `resource_tc_teo_l4_proxy_test.go`
- [x] 4.2 Add test cases for `total_count` being set in `resource_tc_teo_l4_proxy_test.go`

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_l4_proxy.md` with examples showing `offset`, `limit`, and `total_count` usage