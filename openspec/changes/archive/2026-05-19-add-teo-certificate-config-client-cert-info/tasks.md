## 1. Schema Definition

- [x] 1.1 Add `client_cert_info` parameter (TypeList, MaxItems: 1, Optional, Computed) to the resource schema in `resource_tc_teo_certificate_config.go`, with sub-fields `switch` (Required, TypeString) and `cert_infos` (Optional, Computed, TypeList)
- [x] 1.2 Add `cert_infos` sub-schema with fields: `cert_id` (Required, TypeString), `alias` (Computed, TypeString), `type` (Computed, TypeString), `expire_time` (Computed, TypeString), `deploy_time` (Computed, TypeString), `sign_algo` (Computed, TypeString), `status` (Computed, TypeString)

## 2. Read Logic

- [x] 2.1 Update `resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0` in `resource_tc_teo_certificate_config_extension.go` to parse `ClientCertInfo` from `AccelerationDomainCertificate` response, handling nil checks
- [x] 2.2 Parse `ClientCertInfo.Switch` and set `client_cert_info.switch` in state
- [x] 2.3 Parse `ClientCertInfo.CertInfos` list and set `client_cert_info.cert_infos` with all sub-fields (`cert_id`, `alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status`), handling nil checks for each field

## 3. Update Logic

- [x] 3.1 Update `resourceTencentCloudTeoCertificateConfigUpdateOnStart` in `resource_tc_teo_certificate_config_extension.go` to construct and set `request.ClientCertInfo` when `client_cert_info` is specified
- [x] 3.2 Build `MutualTLS` struct with `Switch` and `CertInfos` (only `CertId` needed for write) from the terraform schema data

## 4. Documentation

- [x] 4.1 Update `resource_tc_teo_certificate_config.md` to add `client_cert_info` example usage

## 5. Unit Tests

- [x] 5.1 Add unit tests for the new `client_cert_info` parameter in `resource_tc_teo_certificate_config_test.go` using gomonkey mock approach, covering create with client_cert_info, read with ClientCertInfo in response, and update with client_cert_info
