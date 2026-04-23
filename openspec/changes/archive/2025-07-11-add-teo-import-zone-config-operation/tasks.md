## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_import_zone_config_operation.go` with schema definition (zone_id: Required/ForceNew, content: Required/ForceNew, task_id: Computed) and Timeouts block
- [x] 1.2 Implement Create handler: call `ImportZoneConfig` API, then poll `DescribeZoneConfigImportResult` until status is `success` or `failure`, set ID to `zone_id + FILED_SP + task_id`
- [x] 1.3 Implement Read handler: no-op (return nil)
- [x] 1.4 Implement Delete handler: no-op (return nil)

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_teo_import_zone_config_operation` in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Update `tencentcloud/provider.md` with the new resource entry

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_import_zone_config_operation.md` with description and usage example

## 4. Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_import_zone_config_operation_test.go` with unit tests using gomonkey to mock `ImportZoneConfig` and `DescribeZoneConfigImportResult` API calls (success and failure scenarios)
