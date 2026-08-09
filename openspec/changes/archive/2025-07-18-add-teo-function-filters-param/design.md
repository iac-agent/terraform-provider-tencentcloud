## Context

The `tencentcloud_teo_function` resource manages EdgeOne (TEO) functions in the TencentCloud Terraform Provider. It currently supports the full CRUD lifecycle using the `CreateFunction`, `DescribeFunctions`, `ModifyFunction`, and `DeleteFunction` APIs.

The READ operation currently uses `DescribeTeoFunctionById` which calls the `DescribeFunctions` API with only `ZoneId` and `FunctionIds` parameters. The API also supports a `Filters` parameter (type: `[]*Filter`) that enables filtering by function name (fuzzy match) and remark (fuzzy match), but this capability is not exposed in the Terraform resource.

### Current Schema Fields

| Field | Type | Mode | Description |
|-------|------|------|-------------|
| `zone_id` | string | Required, ForceNew | Zone ID |
| `function_id` | string | Computed | Function ID |
| `name` | string | Required | Function name |
| `remark` | string | Optional | Function description |
| `content` | string | Required | Function JavaScript code |
| `domain` | string | Computed | Default domain |
| `create_time` | string | Computed | Creation time |
| `update_time` | string | Computed | Modification time |

### API Details

The `DescribeFunctionsRequest` struct supports:
- `ZoneId` (*string) - Zone ID (required)
- `FunctionIds` ([]*string) - Filter by function ID list
- `Filters` ([]*Filter) - Filter conditions with `Name` and `Values` fields
  - Supported filter names: `name` (fuzzy match by function name), `remark` (fuzzy match by function description)
- `Offset` (*int64) - Pagination offset
- `Limit` (*int64) - Pagination limit (max: 200)

The `Filter` struct has:
- `Name` (*string) - Field name to filter on
- `Values` ([]*string) - Filter values for the field

## Goals / Non-Goals

**Goals:**
- Add the `filters` parameter to the `tencentcloud_teo_function` resource schema
- Pass the `filters` parameter to the `DescribeFunctions` API during the READ operation
- Maintain full backward compatibility with existing configurations

**Non-Goals:**
- Changing the existing `zone_id`, `function_id`, or other existing schema fields
- Adding pagination parameters (`offset`, `limit`) to the resource schema
- Modifying the CREATE, UPDATE, or DELETE operations
- Creating a new data source for TEO functions

## Decisions

### Decision 1: `filters` as Optional schema field

**Choice**: Add `filters` as an Optional TypeList of nested blocks in the resource schema.

**Rationale**: The `filters` parameter is optional in the `DescribeFunctions` API. When not specified, the API returns all functions filtered by `ZoneId` and `FunctionIds` as it does today. Making it Optional in the Terraform schema ensures backward compatibility.

**Alternative considered**: Making `filters` a Computed-only field — rejected because filters are input parameters for the query, not output from the API.

### Decision 2: Nested block structure for filters

**Choice**: Define `filters` as a TypeList with nested `name` (TypeString, Required) and `values` (TypeList of TypeString, Required) blocks.

**Rationale**: This directly maps to the SDK's `Filter` struct which has `Name` (*string) and `Values` ([]*string). Using TypeList rather than TypeSet preserves order, and nested blocks allow natural HCL configuration.

### Decision 3: Pass filters through service layer

**Choice**: Update the `DescribeTeoFunctionById` service method to accept an optional `filters` parameter, or create a separate service method for filtered queries.

**Rationale**: The existing `DescribeTeoFunctionById` method is used in the READ operation. Adding `filters` support to this method keeps the change minimal. However, since the current method is specifically designed to look up by ID, it may be cleaner to pass filters directly in the READ function when calling the API.

**Selected approach**: Modify the READ function to pass `filters` to the `DescribeFunctions` API request alongside `ZoneId` and `FunctionIds`. The service layer method `DescribeTeoFunctionById` already constructs and sends the `DescribeFunctionsRequest`; we update it to also accept and set `Filters`.

### Decision 4: filters in READ operation only

**Choice**: The `filters` parameter is only used in the READ operation, not in CREATE, UPDATE, or DELETE.

**Rationale**: `Filters` is a query parameter for the `DescribeFunctions` API only. It has no meaning in the `CreateFunction`, `ModifyFunction`, or `DeleteFunction` APIs. The filters are used to narrow down the query results when reading the function state.

### Decision 5: filters should be treated as immutable

**Choice**: Add `filters` to the `immutableArgs` list in the Update function so that changing filters requires resource recreation.

**Rationale**: Since `filters` is a query parameter (not a resource property that the API stores), changing filters after creation could potentially return different results. However, for a GENERAL resource, `filters` is used to help locate the resource during READ. If filters change, the resource identity doesn't change — only how we query it. Given the resource already has `FunctionIds` for precise identification, filters serve as a supplementary query aid. We should add it to `immutableArgs` to avoid confusion, since changing filters mid-lifecycle could cause the READ to find different functions.

## Risks / Trade-offs

- **[Risk] Filters may return multiple results** → The READ operation already uses `FunctionIds` for precise identification, so filters are supplementary. The current code already handles the case where `DescribeFunctions` returns a list by taking the first result. Adding filters does not change this behavior.

- **[Risk] Backward compatibility** → Since `filters` is Optional and not used in existing configurations, there is no backward compatibility risk. The existing READ logic continues to work as before when `filters` is not specified.

- **[Trade-off] filters as immutable** → Marking `filters` as immutable means users need to recreate the resource to change filters. This is acceptable because filters are primarily query aids, not resource properties.
