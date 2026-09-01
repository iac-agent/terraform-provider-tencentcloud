# teo-l4-proxy-rule-params Specification

## Purpose
TBD - created by archiving change add-teo-l4-proxy-rule-params. Update Purpose after archive.
## Requirements
### Requirement: L4 proxy rule supports BuId field
The `tencentcloud_teo_l4_proxy_rule` resource's `l4_proxy_rules` block SHALL support an optional `bu_id` field that accepts a string value representing the business unit ID. The field SHALL be writable in Create and Update operations, and readable from the Describe response.

#### Scenario: Create L4 proxy rule with BuId
- **WHEN** a user creates a `tencentcloud_teo_l4_proxy_rule` resource with `bu_id` set in the `l4_proxy_rules` block
- **THEN** the Create API call SHALL include `BuId` in the `L4ProxyRules` request parameter
- **AND** the Read response SHALL populate `bu_id` in the Terraform state

#### Scenario: Create L4 proxy rule without BuId
- **WHEN** a user creates a `tencentcloud_teo_l4_proxy_rule` resource without specifying `bu_id`
- **THEN** the Create API call SHALL NOT include `BuId` in the `L4ProxyRules` request parameter
- **AND** the Read response SHALL populate `bu_id` from the API response if present

#### Scenario: Update L4 proxy rule BuId
- **WHEN** a user modifies the `bu_id` field in an existing `tencentcloud_teo_l4_proxy_rule` resource
- **THEN** the Update API call SHALL include the new `BuId` value in the `L4ProxyRules` request parameter

#### Scenario: Read L4 proxy rule with BuId from API
- **WHEN** the Describe API returns a `BuId` value in the `L4ProxyRules` response
- **THEN** the Terraform state SHALL reflect the `bu_id` value

### Requirement: L4 proxy rule supports RemoteAuth read-only fields
The `tencentcloud_teo_l4_proxy_rule` resource's `l4_proxy_rules` block SHALL support a Computed-only `remote_auth` nested block containing `switch`, `address`, and `server_faulty_behavior` fields. These fields SHALL be populated from the Describe response only and SHALL NOT be writable in Create or Update operations.

#### Scenario: Read L4 proxy rule with RemoteAuth enabled
- **WHEN** the Describe API returns a `RemoteAuth` object with `Switch` set to "on", `Address` set to "example.auth.com:8888", and `ServerFaultyBehavior` set to "reject"
- **THEN** the Terraform state SHALL reflect `switch` = "on", `address` = "example.auth.com:8888", and `server_faulty_behavior` = "reject" in the `remote_auth` block

#### Scenario: Read L4 proxy rule without RemoteAuth
- **WHEN** the Describe API returns a `RemoteAuth` value of nil
- **THEN** the `remote_auth` block SHALL NOT be set in the Terraform state

#### Scenario: RemoteAuth fields are not writable
- **WHEN** a user specifies `remote_auth` in the `l4_proxy_rules` block of a Terraform configuration
- **THEN** the fields SHALL be ignored in Create and Update API calls
- **AND** on the next Read, the values SHALL be overwritten by the API response

