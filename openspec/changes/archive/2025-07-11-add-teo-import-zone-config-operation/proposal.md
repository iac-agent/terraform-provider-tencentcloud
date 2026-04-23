## Why

TEO (TencentCloud EdgeOne) users need the ability to import zone configurations via Terraform. The `ImportZoneConfig` API allows bulk importing of site acceleration settings and rule engine configurations, but there is currently no Terraform resource to trigger this operation. Adding a one-shot operation resource enables infrastructure-as-code workflows for zone configuration import.

## What Changes

- Add a new one-shot operation resource `tencentcloud_teo_import_zone_config_operation` that calls the `ImportZoneConfig` API
- After calling `ImportZoneConfig`, poll `DescribeZoneConfigImportResult` until the async import task completes (status becomes `success` or `failure`)
- Register the new resource in `provider.go`

## Capabilities

### New Capabilities
- `teo-import-zone-config-operation`: One-shot operation resource that imports TEO zone configuration via `ImportZoneConfig` API, with async polling via `DescribeZoneConfigImportResult`

### Modified Capabilities

## Impact

- New files: `tencentcloud/services/teo/resource_tc_teo_import_zone_config_operation.go`, test file, doc file
- Modified files: `tencentcloud/provider.go` (resource registration), `tencentcloud/provider.md` (resource documentation)
- API dependency: `teo/v20220901` — `ImportZoneConfig`, `DescribeZoneConfigImportResult`
