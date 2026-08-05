## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.go` with schema, CRUD functions, and import support
- [x] 1.2 Register resource `tencentcloud_teo_l7_acc_rules` in `tencentcloud/provider.go`
- [x] 1.3 Register resource `tencentcloud_teo_l7_acc_rules` in `tencentcloud/provider.md`

## 2. Unit Tests

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules_test.go` with gomonkey-based unit tests for CRUD operations

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.md` with example usage and import instructions

## 4. Finalization

- [x] 4.1 Run `gofmt` on all modified Go files (delegated to tfpacer-finalize skill)
- [x] 4.2 Run `make doc` to generate website documentation (delegated to tfpacer-finalize skill)
- [x] 4.3 Create `.changelog/` entry for the change (delegated to tfpacer-finalize skill)