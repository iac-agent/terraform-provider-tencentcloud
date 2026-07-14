## Why

The `tencentcloud_instance_types` data source currently hardcodes the `InquiryType` parameter as `"INQUIRY_CVM_CONFIG"` when calling the `DescribeDiskConfigQuota` API via the CBS service, and the `DiskChargeType` and `InstanceFamilies` parameters are only available as nested fields within the `cbs_filter` block with limited flexibility. Users need the ability to explicitly specify these parameters to control the query behavior of the CBS disk configuration API, enabling more flexible CBS configuration queries such as querying CBS-only configurations (`INQUIRY_CBS_CONFIG`) or specifying specific instance families and disk charge types independently.

## What Changes

- Add `InquiryType` parameter to the `cbs_filter` nested schema block, allowing users to specify the inquiry type (`INQUIRY_CBS_CONFIG` or `INQUIRY_CVM_CONFIG`) for the `DescribeDiskConfigQuota` API call. Currently this is hardcoded to `"INQUIRY_CVM_CONFIG"` and cannot be changed by users.
- Add `InstanceFamilies` parameter to the `cbs_filter` nested schema block, allowing users to specify instance families for filtering. Currently `InstanceFamilies` is derived from the instance type results' `family` field, which limits flexibility when users want to query CBS configurations for specific instance families independently.
- Modify the `DescribeDiskConfigQuota` call in the CBS service layer to pass user-provided `InquiryType` and `InstanceFamilies` values instead of hardcoded or derived values.

## Capabilities

### New Capabilities
- `instance-types-datasource-cbs-params`: Adds three new input parameters (`InquiryType`, `DiskChargeType` as a top-level cbs_filter parameter instead of nested, and `InstanceFamilies`) to the `cbs_filter` block of the `tencentcloud_instance_types` data source, enabling flexible CBS disk configuration queries.

### Modified Capabilities
<!-- No existing spec-level behavior changes -->

## Impact

### Affected Code
- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add `InquiryType` and `InstanceFamilies` fields to the `cbs_filter` schema, update data source read logic to pass these parameters
- `tencentcloud/services/cbs/service_tencentcloud_cbs.go` - Modify `DescribeDiskConfigQuota` method signature and implementation to accept `InquiryType` and `InstanceFamilies` as explicit parameters
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new parameters and examples

### API Changes
- `DescribeDiskConfigQuota` (CBS API) - The method call will now pass user-provided `InquiryType` and `InstanceFamilies` instead of hardcoded values

### Dependencies
None - uses existing CBS SDK structures that already support these parameters

### Breaking Changes
None - all changes are additive. The default behavior remains unchanged (InquiryType defaults to `INQUIRY_CVM_CONFIG` when not specified)
