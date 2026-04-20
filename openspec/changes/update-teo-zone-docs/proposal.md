## Why

The `type` parameter in `tencentcloud_teo_zone` resource is missing documentation for two new access types: `dnsPodAccess` (DNSPod access) and `ai` (AI access). Users cannot discover these valid options through the current documentation, which may lead to confusion or underutilization of available features.

## What Changes

- Update the `Description` field of the `type` parameter in `tencentcloud_teo_zone` resource schema
- Add documentation for two new access types: `dnsPodAccess` and `ai`
- Maintain backward compatibility - no breaking changes

## Capabilities

### New Capabilities
None - this is a documentation update for an existing parameter.

### Modified Capabilities
- `teo-zone-type-parameter`: The `type` parameter in `tencentcloud_teo_zone` resource now supports two additional access types (`dnsPodAccess` and `ai`), requiring documentation update to reflect the complete list of valid values.

## Impact

- **Affected Code**: `tencentcloud/services/teo/resource_tc_teo_zone.go` - schema description field
- **Documentation**: Website documentation will be auto-generated via `make doc` command
- **Backward Compatibility**: Fully backward compatible - no breaking changes
- **Dependencies**: No new dependencies required
