## Why

Terraform Provider for TencentCloud currently has legacy CVM image datasources (`tencentcloud_image` and `tencentcloud_images`) that follow outdated naming conventions and lack several fields from the latest DescribeImages API (e.g., `Tags`, `LicenseType`, `ImageFamily`, `ImageDeprecated`, `CdcCacheStatus`). A new `tencentcloud_cvm_image_4` datasource is needed to follow the current naming convention (`tencentcloud_<product>_<resource>`) and expose the complete Image struct fields from the DescribeImages API.

## What Changes

- Add a new datasource `tencentcloud_cvm_image_4` that calls the CVM DescribeImages API to query image list
- Support input parameters: `image_ids`, `filters`, `instance_type` matching the DescribeImages request
- Expose the full `image_set` output with all fields from the Image struct including `image_id`, `os_name`, `image_type`, `created_time`, `image_name`, `image_description`, `image_size`, `architecture`, `image_state`, `platform`, `image_creator`, `image_source`, `sync_percent`, `is_support_cloudinit`, `snapshot_set`, `tags`, `license_type`, `image_family`, `image_deprecated`, `cdc_cache_status`
- Register the new datasource in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- Add documentation example in `tencentcloud/services/cvm/data_source_tc_cvm_image_4.md`

## Capabilities

### New Capabilities
- `cvm-image-4-datasource`: A new datasource `tencentcloud_cvm_image_4` to query CVM image list via the DescribeImages API, supporting filtering by image IDs, filters, and instance type, returning the complete image set with all available fields.

### Modified Capabilities

## Impact

- New files: `tencentcloud/services/cvm/data_source_tc_cvm_image_4.go`, `tencentcloud/services/cvm/data_source_tc_cvm_image_4_test.go`, `tencentcloud/services/cvm/data_source_tc_cvm_image_4.md`
- Modified files: `tencentcloud/provider.go` (register new datasource), `tencentcloud/provider.md` (document new datasource)
- Dependencies: Uses existing CVM SDK in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312/`
- No breaking changes to existing resources or datasources
