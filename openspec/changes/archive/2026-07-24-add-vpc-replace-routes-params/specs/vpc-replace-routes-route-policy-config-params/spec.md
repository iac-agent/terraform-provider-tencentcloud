# Spec Delta: VPC Replace Routes With Route Policy Config Parameters

**Capability**: `vpc-replace-routes-route-policy-config-params`
**Change ID**: `add-vpc-replace-routes-params`
**Type**: ADDED

---

## ADDED Requirements

### Requirement: Write-Path Parameter Wiring (ReplaceRoutesWithRoutePolicy)

The `tencentcloud_vpc_replace_routes_with_route_policy_config` resource SHALL correctly map its schema fields to the `ReplaceRoutesWithRoutePolicy` cloud API request:

**API Mapping (write path):**

| Terraform Attribute | API Field (ReplaceRoutesWithRoutePolicyRequest) | Type |
|---------------------|--------------------------------------------------|------|
| `route_table_id` | `RouteTableId` | *string |
| `routes` | `Routes` | []*ReplaceRoutesWithRoutePolicyRoute |

Each `routes` block SHALL map to a `ReplaceRoutesWithRoutePolicyRoute` with:
- `route_item_id` → `RouteItemId` (*string)
- `force_match_policy` → `ForceMatchPolicy` (*bool)

**Rationale**: Ensures the existing write-path parameters are correctly wired to the create/update request, as confirmed against the vendored SDK.

#### Scenario: Create resource with route_table_id and routes

**Given** a user defines a `tencentcloud_vpc_replace_routes_with_route_policy_config` resource with `route_table_id` and at least one `routes` block
**When** the resource is created
**Then** `ReplaceRoutesWithRoutePolicy` SHALL be called with `RouteTableId` and `Routes` populated from the schema
**AND** the resource id SHALL be set to the `route_table_id` value

**Example Configuration**:
```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }
}
```

**Acceptance Criteria**:
- `route_table_id` is `Required` and `ForceNew`
- `routes` is `Required`
- The API call includes both fields

#### Scenario: Update routes in-place

**Given** an existing resource with `force_match_policy = true`
**When** the user changes `force_match_policy` to `false`
**Then** `ReplaceRoutesWithRoutePolicy` SHALL be called again with the updated `Routes`
**AND** the resource SHALL NOT be recreated (only `route_table_id` is `ForceNew`)

---

### Requirement: Read-Path Query Parameters (DescribeRouteTables)

The `tencentcloud_vpc_replace_routes_with_route_policy_config` resource SHALL expose optional query-input parameters that influence the `DescribeRouteTables` read call used to refresh state. These parameters describe how to query, NOT what the resource is, and SHALL NOT be set back into state from the response.

**API Mapping (read path):**

| Terraform Attribute | API Field (DescribeRouteTablesRequest) | Type | Notes |
|---------------------|----------------------------------------|------|-------|
| `name` | `Filters[].Name` | string | Populates the `Name` of a `Filter` entry |
| `values` | `Filters[].Values` | []*string | Populates the `Values` of the same `Filter` entry |
| `need_router_info` | `NeedRouterInfo` | *bool | Toggles router info in response |
| `route_table_ids` | `RouteTableIds` | []*string | Mutually exclusive with `Filters` |
| `limit` | `Limit` | *string | Max 100 per cloud API |

**Rationale**: The cloud API `DescribeRouteTables` supports richer query capabilities than the current read path exposes. Surfacing these parameters gives users finer control over route table queries during state refresh.

#### Scenario: Read with custom filter name and values

**Given** a user sets `name = "route-table-name"` and `values = ["my-rtb"]`
**When** the resource `Read` method runs
**Then** the `DescribeRouteTables` request SHALL include a `Filter{Name: "route-table-name", Values: ["my-rtb"]}`
**AND** the `route-table-id` filter (derived from `d.Id()`) SHALL also be present when `route_table_ids` is not set

**Example Configuration**:
```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  name   = "route-table-name"
  values = ["my-rtb"]
}
```

**Acceptance Criteria**:
- `name` and `values` are `Optional`
- When both are set, a single `Filter` is appended to `request.Filters`
- The fields are NOT written back to state from the response

#### Scenario: Read with explicit route_table_ids (Filters mutually exclusive)

**Given** a user sets `route_table_ids = ["rtb-aaa", "rtb-bbb"]`
**When** the resource `Read` method runs
**Then** the `DescribeRouteTables` request SHALL set `RouteTableIds = ["rtb-aaa", "rtb-bbb"]`
**AND** NO `Filters` SHALL be added (cloud API forbids specifying both `RouteTableIds` and `Filters`, including the `route-table-id` filter from `d.Id()`)

**Example Configuration**:
```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  route_table_ids = ["rtb-olsbhnyc"]
}
```

**Acceptance Criteria**:
- `route_table_ids` is `Optional`
- When set, the request uses `RouteTableIds` and omits all `Filters`
- The field is NOT written back to state from the response

#### Scenario: Read with need_router_info disabled

**Given** a user sets `need_router_info = false`
**When** the resource `Read` method runs
**Then** the `DescribeRouteTables` request SHALL set `NeedRouterInfo = false`

**Example Configuration**:
```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  need_router_info = false
}
```

**Acceptance Criteria**:
- `need_router_info` is `Optional` of type `bool`
- When set, `request.NeedRouterInfo` is populated
- The field is NOT written back to state from the response

#### Scenario: Read with custom limit

**Given** a user sets `limit = 50`
**When** the resource `Read` method runs
**Then** the `DescribeRouteTables` request SHALL set `Limit = "50"`
**AND** values greater than 100 SHALL be clamped to 100 (cloud API maximum)

**Example Configuration**:
```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  limit = 50
}
```

**Acceptance Criteria**:
- `limit` is `Optional` of type `int`
- When set, `request.Limit` is populated with the string form of the value
- Values > 100 are clamped to 100
- The field is NOT written back to state from the response

#### Scenario: Omit all new query parameters (backward compatible)

**Given** a user does NOT set any of `name`, `values`, `need_router_info`, `route_table_ids`, `limit`
**When** the resource `Read` method runs
**Then** the read path SHALL behave exactly as before (filter by `route-table-id` from `d.Id()`, internal limit of 100)
**AND** no errors SHALL be raised

**Acceptance Criteria**:
- All new parameters are `Optional`
- Existing configurations continue to work without modification
- No state migration required

---

### Requirement: Dedicated Read Helper (Non-Breaking)

The resource `Read` method SHALL use a dedicated read helper (or inline request construction) for `DescribeRouteTables` rather than mutating the widely-used shared `VpcService.DescribeRouteTables(ctx, routeTableId, routeTableName, vpcId, tags, associationMain, tagKey)` signature.

**Rationale**: The shared helper is called from at least five other resources/data sources. Changing its signature would cascade changes and risk regressions. A dedicated read path keeps this change isolated.

#### Scenario: Shared helper signature unchanged

**Given** the change is applied
**When** other callers (`data_source_tc_route_table`, `data_source_tc_vpc_route_tables`, `resource_tc_route_table_association`, `resource_tc_vpc`) invoke `VpcService.DescribeRouteTables`
**Then** the helper signature and behavior SHALL remain identical (no compile errors, no behavior change)

**Acceptance Criteria**:
- `VpcService.DescribeRouteTables` signature is unchanged
- A new method (e.g. `DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig`) or inline construction handles the new parameters
- The new helper wraps the API call in `tccommon.ReadRetryTimeout` retry and returns `tccommon.RetryError` on failure

---

### Requirement: Empty-Response Handling on Read

When the `DescribeRouteTables` read call returns an empty result, the resource `Read` method SHALL preserve the id in logs before clearing state.

#### Scenario: Resource not found during read

**Given** the `DescribeRouteTables` response is empty (route table deleted out-of-band)
**When** the resource `Read` method runs
**Then** the method SHALL first log `log.Printf("[CRUD] tencentcloud_vpc_replace_routes_with_route_policy_config id=%s", d.Id())`
**AND** then call `d.SetId("")` and return `nil`

**Acceptance Criteria**:
- The id is logged BEFORE `d.SetId("")` is called
- The method returns `nil` (not an error) for a missing resource

---

### Requirement: Mock-Based Unit Tests

The change SHALL include mock-based unit tests (using `gomonkey`) for the new parameters. The tests SHALL NOT rely on the Terraform acceptance test suite.

#### Scenario: Unit tests cover new read-path parameters

**Given** the test file `resource_tc_vpc_replace_routes_with_route_policy_config_test.go`
**When** `go test -gcflags=all=-l` is run against the new test cases
**Then** the following scenarios SHALL be covered:
- Read with `name` + `values` filters — verify the request carries a `Filter{Name, Values}`
- Read with `route_table_ids` set — verify the request uses `RouteTableIds` and NO `Filters`
- Read with `need_router_info = false` — verify `request.NeedRouterInfo == false`
- Read with `limit = 50` — verify `request.Limit == "50"`

**Acceptance Criteria**:
- Tests use `gomonkey` to mock the VPC client's `DescribeRouteTablesWithContext`
- Tests run with `go test -gcflags=all=-l ./tencentcloud/services/vpc/ -run TestVpcReplaceRoutesWithRoutePolicyConfig`
- All new tests pass

---

## MODIFIED Requirements

None. This change adds new parameters without modifying existing spec-level requirements.

---

## REMOVED Requirements

None. This is a pure addition with no deprecations or removals.

---

## Cross-References

### API References
- **TencentCloud API**: VPC Service v20170312
  - `ReplaceRoutesWithRoutePolicy`: maps `route_table_id` and `routes`
  - `DescribeRouteTables`: maps `Filters.Name`, `Filters.Values`, `NeedRouterInfo`, `RouteTableIds`, `Limit`

### Implementation Files
- `tencentcloud/services/vpc/resource_tc_vpc_replace_routes_with_route_policy_config.go` — schema + Read/Create/Update
- `tencentcloud/services/vpc/service_tencentcloud_vpc.go` — dedicated read helper
- `tencentcloud/services/vpc/resource_tc_vpc_replace_routes_with_route_policy_config_test.go` — mock-based unit tests
- `tencentcloud/services/vpc/resource_tc_vpc_replace_routes_with_route_policy_config.md` — documentation

### Backward Compatibility
This change is fully backward compatible:
- All new parameters are `Optional`
- Existing configurations and state remain valid
- No state migration required
- The shared `DescribeRouteTables` helper signature is unchanged
