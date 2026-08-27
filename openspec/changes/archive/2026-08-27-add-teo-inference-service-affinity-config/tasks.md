## 1. Resource Schema & CRUD Implementation

- [x] 1.1 Create `resource_tc_teo_inference_service_v1.go` with schema definition including `affinity_config` nested block with `switch`, `affinity_mode`, `source`, `header_name` sub-fields, and implement Create/Read/Update/Delete CRUD functions
- [x] 1.2 Register `tencentcloud_teo_inference_service_v1` resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`

## 2. Unit Tests

- [x] 2.1 Create `resource_tc_teo_inference_service_v1_test.go` with gomonkey-based unit tests for all CRUD functions

## 3. Documentation

- [x] 3.1 Create `resource_tc_teo_inference_service_v1.md` example documentation file

## 4. Code Quality

- [ ] 4.1 Run gofmt on all modified Go files
- [ ] 4.2 Generate website documentation via `make doc`
- [ ] 4.3 Create changelog entry in `.changelog/` directory