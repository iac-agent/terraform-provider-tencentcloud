## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_teo_edge_kv_namespace_config.go` with RESOURCE_KIND_CONFIG schema (Read + Update only)
- [x] 1.2 Implement Read function using `DescribeEdgeKVNamespaces` API with namespace filter, Limit=1000, pagination support
- [x] 1.3 Implement Update function using `ModifyEdgeKVNamespace` API with read-back after update
- [x] 1.4 Implement Import support using composite ID `zone_id#namespace`

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_edge_kv_namespace` in `tencentcloud/provider.go` resource map
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `resource_tc_teo_edge_kv_namespace_config.md` with example usage and import section

## 4. Unit Tests

- [x] 4.1 Create `resource_tc_teo_edge_kv_namespace_config_test.go` with gomonkey mocks for Read and Update operations