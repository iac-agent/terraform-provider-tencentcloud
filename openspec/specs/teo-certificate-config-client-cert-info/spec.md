## Requirements

### Requirement: client_cert_info parameter in teo_certificate_config resource
The `tencentcloud_teo_certificate_config` resource SHALL support a `client_cert_info` parameter that configures the client-side CA certificate for edge mutual TLS authentication in EdgeOne (TEO). The parameter SHALL be of type `TypeList` with `MaxItems: 1`, Optional and Computed, containing a nested resource with:
- `switch` (Required, TypeString): Mutual TLS configuration switch, values: `on` or `off`.
- `cert_infos` (Optional, TypeList): Mutual TLS certificate list, where each element is a nested resource with:
  - `cert_id` (Required, TypeString): Certificate ID from the SSL certificate list.
  - `alias` (Computed, TypeString): Alias of the certificate.
  - `type` (Computed, TypeString): Type of the certificate.
  - `expire_time` (Computed, TypeString): Certificate expiration time.
  - `deploy_time` (Computed, TypeString): Certificate deployment time.
  - `sign_algo` (Computed, TypeString): Signature algorithm.

#### Scenario: Create resource with client_cert_info enabled
- **WHEN** a user creates a `tencentcloud_teo_certificate_config` resource with `client_cert_info` configured (switch = "on", cert_infos with cert_id)
- **THEN** the resource SHALL call `ModifyHostsCertificate` with the `ClientCertInfo` field populated, and the resource state SHALL reflect the configured client certificate info

#### Scenario: Create resource without client_cert_info
- **WHEN** a user creates a `tencentcloud_teo_certificate_config` resource without specifying `client_cert_info`
- **THEN** the resource SHALL NOT include `ClientCertInfo` in the `ModifyHostsCertificate` request, and the API SHALL keep the original configuration

#### Scenario: Update client_cert_info from disabled to enabled
- **WHEN** a user updates `client_cert_info` from switch "off" to switch "on" with cert_infos
- **THEN** the resource SHALL call `ModifyHostsCertificate` with the updated `ClientCertInfo`, and the state SHALL reflect the new configuration

#### Scenario: Read client_cert_info from API response
- **WHEN** the resource read flow is executed
- **THEN** the resource SHALL read `ClientCertInfo` from the API response and populate the `client_cert_info` parameter in the Terraform state, including all computed sub-fields (alias, type, expire_time, deploy_time, sign_algo)

#### Scenario: Delete resource with client_cert_info
- **WHEN** a user deletes a `tencentcloud_teo_certificate_config` resource that has `client_cert_info` configured
- **THEN** the resource SHALL be deleted without error (the delete operation is a no-op for this resource type)
