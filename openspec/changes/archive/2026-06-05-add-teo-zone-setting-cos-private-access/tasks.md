## 1. Schema and CRUD Code Changes

- [x] 1.1 Add `cos_private_access` field (type: `schema.TypeString`, Optional + Computed, with description) to the `origin` block schema in `tencentcloud/services/teo/resource_tc_teo_zone_setting.go`
- [x] 1.2 Update the Read function to map `respData.Origin.CosPrivateAccess` to `originMap["cos_private_access"]` with nil check, following the existing pattern
- [x] 1.3 Update the Update function to read `cos_private_access` from `originMap` and set `origin.CosPrivateAccess` using `helper.String()`, following the existing pattern for `origin_pull_protocol`

## 2. Unit Tests

- [x] 2.1 Add unit test cases for the `cos_private_access` field in `tencentcloud/services/teo/resource_tc_teo_zone_setting_test.go`, covering read and update scenarios using gomonkey mock

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/resource_tc_teo_zone_setting.md` to include `cos_private_access` in the example usage

## 4. Verification

- [x] 4.1 Run `go test` with `-gcflags=all=-l` on the test file to verify unit tests pass
