## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_teo_l4_proxy.go` with schema definition including all fields (zone_id, proxy_name, area, ipv6, static_ip, accelerate_mainland, d_dos_protection_config, proxy_id, cname, ips, status, l4_proxy_rule_count, update_time)
- [x] 1.2 Implement resource Create method using `CreateL4Proxy` API with retry, composite ID format `zone_id#proxy_id`, and nil check on ProxyId
- [x] 1.3 Implement resource Read method using `DescribeL4Proxy` API with zone_id and proxy-id filter, setting all computed fields from L4Proxy response
- [x] 1.4 Implement resource Update method using `ModifyL4Proxy` API for mutable fields (ipv6, accelerate_mainland), with ForceNew validation for immutable fields
- [x] 1.5 Implement resource Delete method using `DeleteL4Proxy` API with retry
- [x] 1.6 Implement resource Import by parsing composite ID `zone_id#proxy_id`

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_l4_proxy` resource in `tencentcloud/provider.go`
- [x] 2.2 Update `tencentcloud/provider.md` to include the new resource in the documentation index

## 3. Documentation

- [x] 3.1 Create `resource_tc_teo_l4_proxy.md` with example usage and import instructions

## 4. Testing

- [x] 4.1 Create `resource_tc_teo_l4_proxy_test.go` using gomonkey mock for unit testing of CRUD operations

## 5. Validation

- [x] 5.1 Verify code compiles successfully (syntax check only, no `go build`)
- [x] 5.2 Verify all required files are created and provider registration is complete