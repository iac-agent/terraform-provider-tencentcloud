## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.go` with schema definition, CRUD functions (Create/Read/Update/Delete), and Import support
- [x] 1.2 Register `tencentcloud_teo_l7_acc_rules` resource in `tencentcloud/provider.go`

## 2. Testing

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules_test.go` with unit tests using gomonkey mock for API calls

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_l7_acc_rules.md` with example usage and import instructions

## 4. Verification

- [ ] 4.1 Run `gofmt` on all modified Go files
- [ ] 4.2 Run `make doc` to generate website documentation
- [ ] 4.3 Create changelog file in `.changelog/` directory