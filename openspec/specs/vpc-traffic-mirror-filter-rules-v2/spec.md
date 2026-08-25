# vpc-traffic-mirror-filter-rules-v2 Specification

## Purpose
TBD - created by archiving change add-vpc-traffic-mirror-filter-rules-v2. Update Purpose after archive.
## Requirements
### Requirement: Resource schema defines traffic mirror filter rules fields
The `tencentcloud_vpc_traffic_mirror_filter_rules_v2` resource SHALL expose the following schema fields:
- `traffic_mirror_id` (Required, TypeString, ForceNew): The traffic mirror instance ID
- `ingress_filter_rules` (Optional, TypeList): List of ingress filter rules
- `egress_filter_rules` (Optional, TypeList): List of egress filter rules
- Each filter rule item SHALL contain: `src_net`, `dst_net`, `protocol`, `src_port`, `dst_port`, `traffic_mirror_filter_rule_id` (Computed), `priority` (TypeInt), `action`, `description`, `created_time` (Computed)

#### Scenario: Resource schema is correctly defined
- **WHEN** the resource is registered in the provider
- **THEN** all fields are present with correct types, Required/Optional/Computed/ForceNew attributes

### Requirement: Create filter rules via CreateTrafficMirrorFilterRules API
The resource SHALL call `CreateTrafficMirrorFilterRules` API on resource creation. The request SHALL include `TrafficMirrorId`, and optionally `IngressFilterRules` and/or `EgressFilterRules`. On success, the resource ID SHALL be set to the `traffic_mirror_id`.

#### Scenario: Create with both ingress and egress rules
- **WHEN** user defines both `ingress_filter_rules` and `egress_filter_rules`
- **THEN** the API is called with both rule lists and the resource ID is set to the traffic_mirror_id

#### Scenario: Create with empty rules returns error
- **WHEN** user defines neither `ingress_filter_rules` nor `egress_filter_rules`
- **THEN** the provider SHALL return an error indicating at least one rule is required

### Requirement: Read filter rules via DescribeTrafficMirrorFilterRules API
The resource SHALL call `DescribeTrafficMirrorFilterRules` API on resource read. The request SHALL include `TrafficMirrorId` from the resource ID. On success, the response rules SHALL be mapped back to the resource schema fields. If the response is empty (no rules found), the resource SHALL be removed from state by calling `d.SetId("")`.

#### Scenario: Read existing rules
- **WHEN** the resource exists and has rules
- **THEN** `ingress_filter_rules` and `egress_filter_rules` are populated from the API response

#### Scenario: Read when traffic mirror no longer exists
- **WHEN** the API returns empty response or the traffic mirror is not found
- **THEN** the resource is removed from Terraform state

### Requirement: Update filter rules via ModifyTrafficMirrorFilterRules API
The resource SHALL call `ModifyTrafficMirrorFilterRules` API on resource update. The request SHALL include `TrafficMirrorId` and the full set of `IngressFilterRules` and/or `EgressFilterRules` from the new configuration.

#### Scenario: Update rule set
- **WHEN** user modifies filter rules in the configuration
- **THEN** the API is called with the complete new rule set, replacing all existing rules

### Requirement: Delete filter rules via DeleteTrafficMirrorFilterRules API
The resource SHALL call `DeleteTrafficMirrorFilterRules` API on resource deletion. The request SHALL include `TrafficMirrorId` and the current rule IDs (`IngressFilterRuleIds` and `EgressFilterRuleIds`) read from state.

#### Scenario: Delete all rules
- **WHEN** the resource is destroyed
- **THEN** the API is called with all rule IDs from current state, and the resource is removed from Terraform state

### Requirement: API calls use retry with ReadRetryTimeout
All API calls (Create, Read, Update, Delete) SHALL use `tccommon.ReadRetryTimeout` with `helper.Retry()` for eventual consistency. On API errors, the error SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: API call fails with transient error
- **WHEN** an API call returns a retryable error
- **THEN** the operation is retried up to the timeout duration

### Requirement: Resource is registered in provider
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md` with the resource name `tencentcloud_vpc_traffic_mirror_filter_rules_v2`.

#### Scenario: Resource is available in the provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_vpc_traffic_mirror_filter_rules_v2` is available as a managed resource

