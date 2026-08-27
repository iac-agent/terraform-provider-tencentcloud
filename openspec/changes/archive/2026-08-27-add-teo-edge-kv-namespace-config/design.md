## Context

The TEO product already has a general resource `tencentcloud_teo_edge_kv_namespace` (spec: `teo-edge-kv-namespace`) for full CRUD lifecycle management of Edge KV namespaces. This change adds a RESOURCE_KIND_CONFIG variant that manages only the configuration (Read + Update) of an existing namespace, without lifecycle management.

The `DescribeEdgeKVNamespaces` and `ModifyEdgeKVNamespace` APIs are already available in the vendored TEO SDK (`tencentcloud-sdk-go/tencentcloud/teo/v20220901`).

## Goals / Non-Goals

**Goals:**
- Create a RESOURCE_KIND_CONFIG resource `tencentcloud_teo_edge_kv_namespace` that supports Read and Update operations
- Read operation uses `DescribeEdgeKVNamespaces` with `ZoneId` and `Filters` (filtering by namespace name)
- Update operation uses `ModifyEdgeKVNamespace` to update the `remark` field
- Use composite ID `zone_id#namespace` (separated by `tccommon.FILED_SP`)
- Support import using the composite ID format

**Non-Goals:**
- No Create or Delete operations (config resource, not lifecycle resource)
- No pagination parameters exposed to users (handled internally)
- No modification of the existing general `teo-edge-kv-namespace` resource

## Decisions

### Decision 1: Resource Type = RESOURCE_KIND_CONFIG
**Choice**: Implement as RESOURCE_KIND_CONFIG with only Read and Update.
**Rationale**: The requirement specifies managing existing namespace configuration. The existing general resource handles lifecycle. This config resource focuses on reading and updating namespace attributes.
**Alternatives**: Could have been a data source, but the requirement specifies config resource with Update capability.

### Decision 2: Composite ID = zone_id#namespace
**Choice**: Use `zone_id` + `tccommon.FILED_SP` + `namespace` as the resource ID.
**Rationale**: A namespace name is unique within a zone, so the combination of `zone_id` and `namespace` uniquely identifies the resource. This follows the provider's established pattern for composite IDs.

### Decision 3: Read uses DescribeEdgeKVNamespaces with Filters
**Choice**: Use `DescribeEdgeKVNamespaces` with `Filters` containing `namespace` filter to find the target namespace.
**Rationale**: The API does not have a single-namespace lookup endpoint. Using the list endpoint with a filter is the standard approach. Set `Limit` to the maximum value (1000) to avoid pagination issues.

### Decision 4: zone_id and namespace are ForceNew
**Choice**: Mark `zone_id` and `namespace` as `ForceNew` in the schema.
**Rationale**: These fields identify the resource. Changing them means targeting a different namespace, which requires recreation.

### Decision 5: Update is synchronous
**Choice**: Call `ModifyEdgeKVNamespace` and then immediately call `DescribeEdgeKVNamespaces` to read back the updated state.
**Rationale**: The `ModifyEdgeKVNamespace` API is not documented as asynchronous. The response only contains a `RequestId`. Reading back after update ensures state consistency.

## Risks / Trade-offs

- **Risk**: If `DescribeEdgeKVNamespaces` returns multiple namespaces matching the filter, the resource may select the wrong one.
  → **Mitigation**: The `namespace` filter is an exact match, so only one result is expected. The code validates that exactly one result is returned.

- **Risk**: If the `ModifyEdgeKVNamespace` API eventually becomes asynchronous, the current synchronous read-back may read stale data.
  → **Mitigation**: The current API documentation does not indicate async behavior. If this changes, the resource can be updated to add polling logic.