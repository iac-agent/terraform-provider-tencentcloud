## Why

The `area` parameter of the `tencentcloud_teo_l4_proxy` resource is currently marked as `Optional`, but it is a critical configuration that determines the acceleration zone for the L4 proxy instance. In practice, users must always specify the `area` when creating an L4 proxy, and the cloud API documentation also indicates this is a required parameter. Making it optional leads to user confusion and potential misconfigurations where the L4 proxy is created without an explicit acceleration zone, causing unexpected behavior or API errors.

## What Changes

- **BREAKING**: Change the `area` field in `tencentcloud_teo_l4_proxy` resource schema from `Optional: true` to `Required: true`
- Update the Create function to use direct field access (remove `d.GetOk` guard) since the field is now guaranteed to be present
- Update the resource documentation (`.md` file) to reflect that `area` is required in example usage

## Capabilities

### New Capabilities
<!-- No new capabilities introduced -->

### Modified Capabilities
- `teo-l4-proxy`: The `area` attribute changes from optional to required, making it mandatory for all configurations using this resource.

## Impact

- **Code**: `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go` - schema definition and Create function
- **Documentation**: `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md` - update example usage
- **Breaking Change**: Existing `.tf` configurations that omit the `area` field will fail on `terraform plan` / `terraform apply` after this change. Users must add the `area` argument to their configurations.
