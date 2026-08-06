## 1. Resource Schema Definition

- [x] 1.1 Define resource schema in `resource_tc_teo_zone.go` with all parameters: `zone_name`, `type`, `area`, `plan_id`, `alias_zone_name`, `tags`, `paused`, `status`, `ownership_verification`, `name_servers`, `work_mode_infos`, `resource_region`, `service_type`
- [x] 1.2 Register `ResourceTencentCloudTeoZone()` in `tencentcloud/provider.go` under `ResourcesMap`
- [x] 1.3 Add resource entry in `tencentcloud/provider.md`

## 2. Service Layer Implementation

- [x] 2.1 Implement `DescribeTeoZone` in `service_tencentcloud_teo.go` for querying zone by ID using `DescribeZones` API with `zone-id` filter and pagination
- [x] 2.2 Implement `DescribeTeoZoneById` in `service_tencentcloud_teo.go` for the resource read handler
- [x] 2.3 Implement `ModifyZoneStatus` in `service_tencentcloud_teo.go` for pause/resume operations

## 3. CRUD Implementation

- [x] 3.1 Implement `resourceTencentCloudTeoZoneCreate` using `CreateZone` API with retry, post-create polling for zone readiness, and tag attachment
- [x] 3.2 Implement `resourceTencentCloudTeoZoneRead` using `DescribeTeoZoneById` to populate all schema attributes including `ownership_verification`, `plan_id` (from Resources), `tags`, and `work_mode_infos`
- [x] 3.3 Implement `resourceTencentCloudTeoZoneUpdate` with three separate update paths: `ModifyZone` for config changes, `ModifyZoneStatus` for pause/resume, and `ModifyZoneWorkMode` for work mode changes, plus tag updates
- [x] 3.4 Implement `resourceTencentCloudTeoZoneDelete` using `DeleteZone` API with pre-delete pause check

## 4. Extension Handlers

- [x] 4.1 Implement `resourceTencentCloudTeoZoneCreateRequestOnError0` for error handling during create (non-retryable on `ResourceInUse` errors)
- [x] 4.2 Implement `resourceTencentCloudTeoZoneCreatePostHandleResponse0` for post-create polling until zone exits `pending` status
- [x] 4.3 Implement `resourceTencentCloudTeoZoneReadPostHandleResponse0` for extracting `plan_id` from `Resources` field
- [x] 4.4 Implement `resourceTencentCloudTeoZoneUpdatePostRequest1` for post-update status polling
- [x] 4.5 Implement `resourceTencentCloudTeoZoneDeletePostFillRequest0` for pre-delete zone pause check

## 5. Documentation

- [x] 5.1 Create `resource_tc_teo_zone.md` with example usage (basic, custom tag params, version control mode) and import instructions
- [x] 5.2 Run `make doc` to generate `website/docs/` documentation

## 6. Testing

- [x] 6.1 Implement `TestAccTencentCloudTeoZone_basic` acceptance test covering create, import, update (alias_zone_name, paused), and work mode changes
- [x] 6.2 Implement sweep function `testSweepZone` for test cleanup

## 7. Verification

- [x] 7.1 Verify `go vet` passes for all changed files
- [x] 7.2 Verify `gofmt` formatting is applied to all changed Go files
- [x] 7.3 Verify provider compiles successfully with `go build`