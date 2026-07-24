## Context

The `tencentcloud_vpc_replace_routes_with_route_policy_config` resource (RESOURCE_KIND_GENERAL) currently manages two parameters — `route_table_id` and `routes` — both of which are wired into the `ReplaceRoutesWithRoutePolicy` cloud API. The resource's `Read` method refreshes state by calling the shared `VpcService.DescribeRouteTables` helper, which today only supports filtering by `route-table-id`, `route-table-name`, `vpc-id`, `association.main`, `tag-key`, and tags. It does not expose `NeedRouterInfo`, explicit `RouteTableIds`, or a configurable `Limit`.

The cloud API `DescribeRouteTables` (vpc/v20170312) exposes additional request fields:
- `Filters []*Filter` (each `Filter` has `Name *string` and `Values []*string`)
- `RouteTableIds []*string`
- `Limit *string` (max 100)
- `NeedRouterInfo *bool`

The `ReplaceRoutesWithRoutePolicy` request already exposes `RouteTableId *string` and `Routes []*ReplaceRoutesWithRoutePolicyRoute`, both of which are already in the resource schema. This change confirms the wiring of those two write-path fields and introduces the missing read-path query parameters (`Name`, `Values`, `NeedRouterInfo`, `RouteTableIds`, `Limit`).

## Goals / Non-Goals

**Goals:**
- Ensure `route_table_id` and `routes` are correctly mapped from the Terraform schema to the `ReplaceRoutesWithRoutePolicy` request (create/update path).
- Expose `Name`, `Values`, `NeedRouterInfo`, `RouteTableIds`, and `Limit` as optional schema parameters so the `DescribeRouteTables` read path can use them.
- Keep the change fully backward compatible: all new fields are `Optional` and the existing `route_table_id`/`routes` behavior is unchanged.
- Provide mock-based unit tests (gomonkey) for the new parameters without relying on the Terraform acceptance test suite.

**Non-Goals:**
- Refactoring the shared `VpcService.DescribeRouteTables` helper signature in a way that breaks other callers (`data_source_tc_route_table`, `data_source_tc_vpc_route_tables`, `resource_tc_route_table_association`, `resource_tc_vpc`). The existing helper signature MUST be preserved; new read-path parameters will be delivered through a dedicated read helper or through additional optional fields handled inside the resource's `Read` method.
- Changing the resource kind or its import behavior.
- Exposing `Offset` as a schema parameter (pagination is handled internally by the helper).

## Decisions

### Decision 1: Add a dedicated read helper instead of mutating the shared `DescribeRouteTables` signature
**Choice:** Add a new method (e.g. `VpcService.DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig`) that builds its own `DescribeRouteTablesRequest`, or construct the request inline inside the resource's `Read` method, rather than extending the existing `DescribeRouteTables(ctx, routeTableId, routeTableName, vpcId, tags, associationMain, tagKey)` signature.

**Rationale:** The existing helper is called from at least five other resources/data sources. Changing its signature would cascade changes across the codebase and risk regressions. A dedicated read path keeps this change isolated and low-risk.

**Alternatives considered:**
- Add variadic options to the existing helper — rejected because it complicates the widely-used signature and makes call sites harder to read.
- Add new positional parameters — rejected because the helper already has six positional parameters; adding more harms readability.

### Decision 2: Map `Name`/`Values` to a `Filters`-shaped block rather than two flat top-level scalars
**Choice:** Because the requirement specifies `request.Filters.Name` → `Name` and `request.Filters.Values` → `Values`, and a single `Filter` carries both a `Name` and a list of `Values`, represent these as schema fields that populate a `Filter` entry in the read request. When both are provided, a `Filter{Name, Values}` is appended to `request.Filters`.

**Rationale:** Mirrors the cloud API's `Filter` shape and avoids ambiguity. `route_table_id` already drives the `route-table-id` filter for state refresh, so the user-supplied `Name`/`Values` provide an additional filter dimension.

**Alternatives considered:**
- Expose a generic `filters` list block — rejected because the requirement explicitly names `Name` and `Values` as the target schema fields.

### Decision 3: `RouteTableIds` and `Limit` are read-only query inputs, not part of the resource state
**Choice:** `RouteTableIds`, `Limit`, `Name`, `Values`, and `NeedRouterInfo` are `Optional` input parameters used only to influence the `DescribeRouteTables` read call. They are not set back into state from the response (they are query inputs, not resource attributes).

**Rationale:** These fields describe how to query, not what the resource is. Setting them from state would be semantically incorrect and could cause diff churn.

### Decision 4: `Limit` capped at 100
**Choice:** When the user provides `Limit`, clamp/validate to the cloud API maximum of 100 (the API annotation states "最大值为100"). If unset, the read helper uses its existing internal limit of 100.

**Rationale:** Matches the cloud API constraint documented in `DescribeRouteTablesRequest.Limit`.

## Risks / Trade-offs

- **Risk:** New read-path parameters are ignored if the user does not set them, and the existing `route-table-id` filter (derived from `d.Id()`) still drives state refresh. → **Mitigation:** Preserve the existing `route-table-id` filter behavior unconditionally; user-supplied `Name`/`Values` only add to the filter list.
- **Risk:** `Filters` and `RouteTableIds` are mutually exclusive in the cloud API ("参数不支持同时指定RouteTableIds和Filters"). → **Mitigation:** Document this in schema descriptions; if the user provides `RouteTableIds`, skip adding the `route-table-id`/`Name`/`Values` filters and query by `RouteTableIds` directly.
- **Risk:** Adding a dedicated read helper duplicates some pagination logic. → **Mitigation:** Keep the helper minimal and reuse `fillFilter` where possible; duplication is confined to the new resource's read path.
