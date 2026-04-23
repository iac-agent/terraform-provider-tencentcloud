## 1. Schema Definition

- [x] 1.1 Add `client_cert_info` parameter to the schema in `tencentcloud/services/teo/resource_tc_teo_certificate_config.go`, as a `TypeList` with `MaxItems: 1`, Optional and Computed, containing nested `switch` (Required, TypeString) and `cert_infos` (Optional, TypeList of nested resources with cert_id Required + alias/type/expire_time/deploy_time/sign_algo Computed)

## 2. Update Logic (Create/Update)

- [x] 2.1 Update `resourceTencentCloudTeoCertificateConfigUpdateOnStart` in `tencentcloud/services/teo/resource_tc_teo_certificate_config_extension.go` to populate `request.ClientCertInfo` (type `*MutualTLS`) when `client_cert_info` is set in the Terraform configuration, mapping `switch` to `MutualTLS.Switch` and `cert_infos` to `MutualTLS.CertInfos`

## 3. Read Logic

- [x] 3.1 Update `resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0` in `tencentcloud/services/teo/resource_tc_teo_certificate_config_extension.go` to read `ClientCertInfo` from the API response (`AccelerationDomain.Certificate.ClientCertInfo`) and populate the `client_cert_info` field in the Terraform state

## 4. Documentation

- [x] 4.1 Update the `.md` example file `tencentcloud/services/teo/resource_tc_teo_certificate_config.md` to include the new `client_cert_info` parameter in the example usage

## 5. Unit Tests

- [x] 5.1 Add unit tests in `tencentcloud/services/teo/resource_tc_teo_certificate_config_test.go` using gomonkey mock to verify the `client_cert_info` parameter is correctly handled in create, read, and update flows
