## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_teo_l4_proxy_1.go` with schema definition and CRUD functions (ResourceTencentCloudTeoL4Proxy1), following the existing `resource_tc_teo_l4_proxy.go` pattern
- [x] 1.2 Create `resource_tc_teo_l4_proxy_1_extension.go` with post-handle (state waiting after create) and pre-delete (stop instance before delete) logic, reusing `teoL4proxyStateRefreshFunc` from the existing extension file
- [x] 1.3 Register `tencentcloud_teo_l4_proxy_1` in `tencentcloud/provider.go` alongside existing TEO resources

## 2. Unit Tests

- [x] 2.1 Create `resource_tc_teo_l4_proxy_1_test.go` with unit tests using gomonkey to mock TEO API calls, covering create, read, update, delete, and import scenarios

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_l4_proxy_1.md` with example usage and import instructions

## 4. Provider Registration

- [x] 4.1 Add `tencentcloud_teo_l4_proxy_1` entry to `tencentcloud/provider.md`

## 5. Finalization (via tfpacer-finalize skill)

- [ ] 5.1 Run `gofmt` on all modified Go files
- [ ] 5.2 Run `make doc` to generate website documentation
- [ ] 5.3 Create `.changelog/<PR_NUMBER>.txt` with changelog entry