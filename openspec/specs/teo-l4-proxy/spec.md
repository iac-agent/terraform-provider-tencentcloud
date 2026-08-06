## Requirements

### Requirement: Resource Create
The system SHALL support creating a TEO L4 proxy instance via the `tencentcloud_teo_l4_proxy` resource. The resource MUST call the `CreateL4Proxy` API with required fields (zone_id, proxy_name, area) and optional fields (ipv6, static_ip, accelerate_mainland, d_dos_protection_config). The resource ID SHALL be set to `<zone_id>#<proxy_id>` after successful creation.

#### Scenario: Create L4 proxy with required fields only
- **WHEN** user applies a `tencentcloud_teo_l4_proxy` resource with `zone_id`, `proxy_name`, and `area`
- **THEN** the provider calls `CreateL4Proxy` API and sets the resource ID to `zone_id#proxy_id`

#### Scenario: Create L4 proxy with all optional fields
- **WHEN** user applies a `tencentcloud_teo_l4_proxy` resource with all optional fields (ipv6, static_ip, accelerate_mainland, d_dos_protection_config)
- **THEN** the provider calls `CreateL4Proxy` API with all specified parameters and sets the resource ID

#### Scenario: Create returns empty proxy_id
- **WHEN** `CreateL4Proxy` API returns a response with nil or empty ProxyId
- **THEN** the provider SHALL return a NonRetryableError

### Requirement: Resource Read
The system SHALL support reading a TEO L4 proxy instance via the `DescribeL4Proxy` API. The Read method MUST parse the composite ID into zone_id and proxy_id, then call DescribeL4Proxy with zone_id and a proxy-id filter. All computed fields (proxy_name, area, cname, ips, status, ipv6, static_ip, accelerate_mainland, d_dos_protection_config, l4_proxy_rule_count, update_time) SHALL be set from the API response.

#### Scenario: Read existing L4 proxy
- **WHEN** the provider reads a `tencentcloud_teo_l4_proxy` resource with a valid composite ID
- **THEN** the provider calls `DescribeL4Proxy` with zone_id and proxy-id filter, and sets all computed fields

#### Scenario: Read non-existent L4 proxy
- **WHEN** `DescribeL4Proxy` returns an empty L4Proxies list
- **THEN** the provider SHALL log the id and call `d.SetId("")` to remove the resource from state

#### Scenario: Read with nil response
- **WHEN** `DescribeL4Proxy` returns nil response or nil L4Proxies
- **THEN** the provider SHALL log the id and call `d.SetId("")` to remove the resource from state

### Requirement: Resource Update
The system SHALL support updating mutable fields (ipv6, accelerate_mainland) via the `ModifyL4Proxy` API. Update MUST check for changes in mutable fields and only call the API when at least one mutable field has changed.

#### Scenario: Update ipv6 and accelerate_mainland
- **WHEN** user modifies `ipv6` and `accelerate_mainland` in an existing resource
- **THEN** the provider calls `ModifyL4Proxy` with the updated values

#### Scenario: No changes to mutable fields
- **WHEN** user only modifies ForceNew fields (proxy_name, area, static_ip)
- **THEN** the provider SHALL return an error indicating these fields are immutable

### Requirement: Resource Delete
The system SHALL support deleting a TEO L4 proxy instance via the `DeleteL4Proxy` API. After deletion, the provider SHALL read the resource to confirm deletion.

#### Scenario: Delete existing L4 proxy
- **WHEN** user runs `terraform destroy` on a `tencentcloud_teo_l4_proxy` resource
- **THEN** the provider calls `DeleteL4Proxy` with zone_id and proxy_id, then confirms deletion via Read

### Requirement: Resource Import
The system SHALL support importing existing TEO L4 proxy instances using the composite ID format `zone_id#proxy_id`.

#### Scenario: Import existing L4 proxy
- **WHEN** user runs `terraform import tencentcloud_teo_l4_proxy.example zone-xxx#proxy-xxx`
- **THEN** the provider parses the composite ID, calls `DescribeL4Proxy` to populate state, and the resource is managed by Terraform

### Requirement: Retry on API failure
The system SHALL use retry logic with `tccommon.ReadRetryTimeout` for all cloud API calls. API errors SHALL be wrapped with `tccommon.RetryError()` to enable retry behavior.

#### Scenario: Retry on transient API error
- **WHEN** a cloud API call fails with a transient error
- **THEN** the provider retries the call up to the timeout duration

### Requirement: Schema field validation
The system SHALL define schema fields with appropriate types, constraints, and ForceNew flags matching the API capabilities.

#### Scenario: Required fields are enforced
- **WHEN** user omits `zone_id`, `proxy_name`, or `area`
- **THEN** Terraform SHALL report a validation error during plan

#### Scenario: Immutable fields trigger recreation
- **WHEN** user modifies `proxy_name`, `area`, or `static_ip`
- **THEN** Terraform SHALL plan to destroy and recreate the resource