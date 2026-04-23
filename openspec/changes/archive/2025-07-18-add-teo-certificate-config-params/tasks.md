## 1. Schema Definition

- [x] 1.1 Add `apply_type` parameter (Optional+Computed, TypeString) to the `tencentcloud_teo_certificate_config` resource schema in `resource_tc_teo_certificate_config.go`
- [x] 1.2 Add `client_cert_info` parameter (Optional+Computed, TypeList, MaxItems:1) with `switch` (Required, TypeString) and `cert_infos` (Optional+Computed, TypeList) sub-fields to the resource schema. The `cert_infos` list contains items with `cert_id` (Required, TypeString)

## 2. Update Method (ModifyHostsCertificate API call)

- [x] 2.1 Add `apply_type` parameter to `ModifyHostsCertificate` request construction in `resourceTencentCloudTeoCertificateConfigUpdateOnStart` function in `resource_tc_teo_certificate_config_extension.go`
- [x] 2.2 Add `client_cert_info` parameter to `ModifyHostsCertificate` request construction in `resourceTencentCloudTeoCertificateConfigUpdateOnStart` function, mapping the `switch` and `cert_infos` fields to the `MutualTLS` struct

## 3. Read Method (API response handling)

- [x] 3.1 Add `client_cert_info` read-back logic in `resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0` function to read from `accelerationDomain.Certificate.ClientCertInfo`
- [x] 3.2 Add `apply_type` read-back logic by calling `DescribeHostsSetting` API in the read post-handler and extracting `ApplyType` from `DetailHost.Https.ApplyType`
- [x] 3.3 Add service layer method `DescribeTeoHostsSetting` in `service_tencentcloud_teo.go` to call `DescribeHostsSetting` API

## 4. Unit Tests

- [x] 4.1 Add unit tests for `apply_type` parameter in `resource_tc_teo_certificate_config_test.go` using gomonkey mock
- [x] 4.2 Add unit tests for `client_cert_info` parameter in `resource_tc_teo_certificate_config_test.go` using gomonkey mock
- [x] 4.3 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_certificate_config.md` example file with `apply_type` and `client_cert_info` usage examples
