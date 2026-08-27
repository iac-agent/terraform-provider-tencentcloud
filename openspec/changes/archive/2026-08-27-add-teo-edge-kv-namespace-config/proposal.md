## Why

TEO Edge KV Namespace is a configuration resource that provides KV storage namespace management for EdgeOne edge functions. Users need to be able to query and update KV namespace attributes (such as remarks) through Terraform. The `DescribeEdgeKVNamespaces` and `ModifyEdgeKVNamespace` APIs are already available in the TEO SDK, but there is no corresponding Terraform resource to manage these configurations.

## What Changes

- Add a new Terraform RESOURCE_KIND_CONFIG resource `tencentcloud_teo_edge_kv_namespace` for TEO product
- The resource supports Read (via `DescribeEdgeKVNamespaces`) and Update (via `ModifyEdgeKVNamespace`) operations
- The resource uses `zone_id` + `namespace` as a composite ID for identification
- Read operation returns KV namespace details including capacity, references, and timestamps
- Update operation supports modifying the namespace remark/description

## Capabilities

### New Capabilities
- `teo-edge-kv-namespace-config`: Terraform resource to manage TEO Edge KV namespace configuration, supporting query and update of KV namespace attributes

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- **New files**: `resource_tc_teo_edge_kv_namespace_config.go`, `resource_tc_teo_edge_kv_namespace_config_test.go`, `resource_tc_teo_edge_kv_namespace_config.md`
- **Modified files**: `provider.go` (register new resource), `provider.md` (document new resource)
- **Dependencies**: Uses existing `tencentcloud-sdk-go/tencentcloud/teo/v20220901` SDK (already vendored)
- **No breaking changes** to existing resources or configurations