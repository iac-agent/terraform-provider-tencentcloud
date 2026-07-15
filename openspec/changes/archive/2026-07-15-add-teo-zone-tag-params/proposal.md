## Why

The `tencentcloud_teo_zone` resource currently hardcodes `serviceType="teo"`, `resourceType="zone"`, and `resourceRegion=tcClient.Region` when calling tag-related APIs (DescribeResourceTagsByResourceIds for reading tags and ModifyResourceTags for creating/updating tags). For TEO (EdgeOne), which is a global service, the tag API requires specific `ResourceRegion` and `ServiceType` parameters that may differ from the provider's default region. Users need the ability to configure these parameters explicitly to ensure correct tag operations for their TEO zone resources.

## What Changes

- Add `ResourceRegion` (string, optional) parameter to the `tencentcloud_teo_zone` resource schema, which will be passed to the `DescribeResourceTagsByResourceIds` request as `ResourceRegion` when reading resource tags
- Add `ServiceType` (string, optional) parameter to the `tencentcloud_teo_zone` resource schema, which will be passed to the `DescribeResourceTagsByResourceIds` request as `ServiceType` when reading resource tags
- Update the `resourceTencentCloudTeoZoneRead` method to use `ResourceRegion` and `ServiceType` from schema when calling `tagService.DescribeResourceTags`
- Update the `resourceTencentCloudTeoZoneCreate` and `resourceTencentCloudTeoZoneUpdate` methods to use `ResourceRegion` and `ServiceType` from schema when constructing the QCS resource name for tag operations via `ModifyResourceTags`

## Capabilities

### New Capabilities
- `teo-zone-tag-params`: Adds `ResourceRegion` and `ServiceType` configurable parameters to the `tencentcloud_teo_zone` resource, enabling users to specify custom region and service type values for tag API operations instead of using hardcoded defaults

### Modified Capabilities
<!-- No existing capabilities are being modified at the spec level -->

## Impact

- Affected code: `tencentcloud/services/teo/resource_tc_teo_zone.go` (schema definition, CRUD methods for tag operations)
- Affected code: `tencentcloud/services/teo/resource_tc_teo_zone.md` (documentation update)
- Affected APIs: `DescribeResourceTagsByResourceIds` (tag/v20180813), `ModifyResourceTags` (tag/v20180813) — these APIs already support `ResourceRegion` and `ServiceType` parameters
- Backward compatibility: Both new parameters are optional with computed defaults (falling back to existing hardcoded values), so existing Terraform configurations remain unaffected
