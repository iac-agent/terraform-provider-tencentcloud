## ADDED Requirements

### Requirement: Resource schema supports filters input parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL include a `filters` schema field of type `TypeList` with `Optional: true`. Each filter item SHALL contain `name` (TypeString, Required) and `values` (TypeList of TypeString, Required), mapping to the `Filter` struct (`Name`, `Values`) in the `DescribeL7AccRules` API request.

#### Scenario: User specifies filters to query L7 acceleration rules
- **WHEN** a user configures `filters` with `name = "rule-id"` and `values = ["rule-abc123"]` in the `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the Read function SHALL pass these filters to the `DescribeL7AccRules` API request's `Filters` field

#### Scenario: User does not specify filters
- **WHEN** a user does not configure `filters` in the `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the Read function SHALL use the existing behavior (filtering by rule-id from the resource ID)

### Requirement: Resource schema supports rules computed output parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL include a `rules` schema field of type `TypeList` with `Computed: true`. Each rule item SHALL contain `status` (TypeString, Computed), `rule_id` (TypeString, Computed), `rule_name` (TypeString, Computed), `description` (TypeList of TypeString, Computed), `branches` (TypeList, Computed with same nested schema as the top-level `branches`), and `rule_priority` (TypeInt, Computed), mapping to the `RuleEngineItem` struct in the `DescribeL7AccRules` API response.

#### Scenario: Read function populates rules from API response
- **WHEN** the Read function calls `DescribeL7AccRules` and receives a response with `Rules` data
- **THEN** the `rules` computed field SHALL be populated with the flattened list of `RuleEngineItem` entries from the response

#### Scenario: API response contains no rules
- **WHEN** the Read function calls `DescribeL7AccRules` and the response contains an empty `Rules` list
- **THEN** the `rules` computed field SHALL be set to an empty list

### Requirement: Service layer method for DescribeL7AccRules with filters
A new service method `DescribeTeoL7AccRuleByFilters` SHALL be added to `TeoService` in `service_tencentcloud_teo.go`. This method SHALL accept `zoneId` and `filters` (of type `[]*teov20220901.Filter`) as parameters, construct a `DescribeL7AccRules` request with the provided filters and `Limit` set to 1000 (the API maximum), and return the response parameters.

#### Scenario: Calling DescribeL7AccRules with custom filters
- **WHEN** `DescribeTeoL7AccRuleByFilters` is called with a zoneId and filters
- **THEN** the method SHALL construct a `DescribeL7AccRules` request with the given `ZoneId` and `Filters`, set `Limit` to 1000, and return the `DescribeL7AccRulesResponseParams`

#### Scenario: Calling DescribeL7AccRules with empty filters
- **WHEN** `DescribeTeoL7AccRuleByFilters` is called with a zoneId and nil/empty filters
- **THEN** the method SHALL construct a `DescribeL7AccRules` request with the given `ZoneId` only, set `Limit` to 1000, and return the response

### Requirement: Backward compatibility with existing resource behavior
The addition of `filters` and `rules` parameters SHALL NOT break existing Terraform configurations. The `filters` parameter is optional and defaults to the existing behavior. The `rules` parameter is computed-only and does not require user configuration.

#### Scenario: Existing Terraform configuration without new parameters
- **WHEN** an existing `tencentcloud_teo_l7_acc_rule_v2` resource configuration does not include `filters` or `rules`
- **THEN** the resource SHALL behave identically to before, with the existing single-rule read logic using `zone_id` and `rule_id`
