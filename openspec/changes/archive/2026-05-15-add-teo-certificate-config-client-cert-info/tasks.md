## 1. Schema Definition

- [x] 1.1 Add `client_cert_info` parameter to `ResourceTencentCloudTeoCertificateConfig()` schema in `resource_tc_teo_certificate_config.go`, with TypeList, MaxItems: 1, Optional + Computed, containing `switch` (TypeString, Required) and `cert_infos` (TypeList, Optional + Computed) sub-fields. The `cert_infos` sub-field shall contain: `cert_id` (TypeString, Required), `alias` (TypeString, Computed), `type` (TypeString, Computed), `expire_time` (TypeString, Computed), `deploy_time` (TypeString, Computed), `sign_algo` (TypeString, Computed), `status` (TypeString, Computed).

## 2. Write Path (Update/Create)

- [x] 2.1 Update `resourceTencentCloudTeoCertificateConfigUpdateOnStart()` in `resource_tc_teo_certificate_config_extension.go` to read `client_cert_info` from ResourceData and construct `ClientCertInfo` (MutualTLS) in the `ModifyHostsCertificate` request. When `client_cert_info` is set, populate `ClientCertInfo.Switch` and `ClientCertInfo.CertInfos` (only `CertId` needed for input per API docs).

## 3. Read Path

- [x] 3.1 Update `resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0()` in `resource_tc_teo_certificate_config_extension.go` to read `ClientCertInfo` from `AccelerationDomainCertificate` and populate `client_cert_info` in terraform state. Handle nil checks for `ClientCertInfo`, `Switch`, `CertInfos`, and each `CertificateInfo` field before setting.

## 4. Unit Tests

- [x] 4.1 Add unit test cases for `client_cert_info` in `resource_tc_teo_certificate_config_test.go` using gomonkey mock approach: test that the schema correctly accepts `client_cert_info` with `switch` and `cert_infos`, test the read path populates state correctly when `ClientCertInfo` is present, and test the read path when `ClientCertInfo` is nil.

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_certificate_config.md` to add `client_cert_info` usage example in the Example Usage section.
