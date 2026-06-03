### Requirement: Datasource schema for tencentcloud_cvm_image_4
The datasource SHALL define a schema with the following input parameters:
- `image_ids`: Optional, TypeList of TypeString. List of image IDs to query. Mutually exclusive with `filters`.
- `filters`: Optional, TypeList of nested Resource with `name` (Required, TypeString) and `values` (Required, TypeSet of TypeString). Filter conditions for the query. Mutually exclusive with `image_ids`.
- `instance_type`: Optional, TypeString. Instance type for compatibility check, e.g., `SA5.MEDIUM2`.
- `result_output_file`: Optional, TypeString. Used to save results.

The datasource SHALL define the following computed output parameter:
- `image_set`: Computed, TypeList of nested Resource containing all Image struct fields.

#### Scenario: Query images by image IDs
- **WHEN** user provides `image_ids` parameter with one or more image IDs
- **THEN** the datasource SHALL call DescribeImages API with ImageIds set to the provided values and return matching images in `image_set`

#### Scenario: Query images by filters
- **WHEN** user provides `filters` parameter with name and values pairs
- **THEN** the datasource SHALL call DescribeImages API with Filters set to the provided filter conditions and return matching images in `image_set`

#### Scenario: Query images with instance type
- **WHEN** user provides `instance_type` parameter along with either `image_ids` or `filters`
- **THEN** the datasource SHALL call DescribeImages API with InstanceType set to the provided value for compatibility check

#### Scenario: Query images without any filter
- **WHEN** user does not provide `image_ids`, `filters`, or `instance_type`
- **THEN** the datasource SHALL call DescribeImages API without any filter and return all available images in `image_set`

### Requirement: Image set output fields
The `image_set` output SHALL contain the following computed fields for each image:
- `image_id`: TypeString. Image ID.
- `os_name`: TypeString. OS name of the image.
- `image_type`: TypeString. Image type (PUBLIC_IMAGE, PRIVATE_IMAGE, SHARED_IMAGE, etc.).
- `created_time`: TypeString. Image creation time.
- `image_name`: TypeString. Image name.
- `image_description`: TypeString. Image description.
- `image_size`: TypeInt. Image size in GiB.
- `architecture`: TypeString. Architecture (x86_64, arm, i386).
- `image_state`: TypeString. Image state (CREATING, NORMAL, CREATEFAILED, SYNCING, IMPORTING, IMPORTFAILED).
- `platform`: TypeString. Source platform.
- `image_creator`: TypeString. Image creator.
- `image_source`: TypeString. Image source (OFFICIAL, CREATE_IMAGE, EXTERNAL_IMPORT).
- `sync_percent`: TypeInt. Sync percentage.
- `is_support_cloudinit`: TypeBool. Whether cloud-init is supported.
- `snapshot_set`: TypeList of nested Resource with `snapshot_id` (TypeString), `disk_usage` (TypeString), `disk_size` (TypeInt).
- `tags`: TypeList of nested Resource with `key` (TypeString) and `value` (TypeString).
- `license_type`: TypeString. License type (TencentCloud, BYOL).
- `image_family`: TypeString. Image family.
- `image_deprecated`: TypeBool. Whether the image is deprecated.
- `cdc_cache_status`: TypeString. CDC image cache status.

#### Scenario: All image fields are populated from API response
- **WHEN** DescribeImages API returns an Image with all fields populated
- **THEN** each field in `image_set` SHALL be set from the corresponding API response field, with nil checks to avoid setting empty values

#### Scenario: Image fields with nil values
- **WHEN** DescribeImages API returns an Image with some nil fields
- **THEN** the datasource SHALL skip setting those fields (not call Set for nil values) to avoid overwriting with zero values

### Requirement: Read function with retry and pagination
The datasource read function SHALL:
1. Use `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()`
2. Create context using `tccommon.NewResourceLifeCycleHandleFuncContext()`
3. Build a `paramMap` from schema inputs
4. Call the service method inside `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `tccommon.RetryError()` for error wrapping
5. Handle pagination with Limit set to maximum value (100)
6. Set the datasource ID using `helper.BuildToken()`
7. Handle `result_output_file` output with `tccommon.WriteToFile()`

#### Scenario: Successful read with pagination
- **WHEN** DescribeImages API returns more than 100 images
- **THEN** the datasource SHALL paginate through all results using Offset and Limit=100 until all images are collected

#### Scenario: API call failure with retry
- **WHEN** DescribeImages API call fails with a retryable error
- **THEN** the datasource SHALL retry the call up to ReadRetryTimeout

### Requirement: Service layer method
A new service method `DescribeCvmImage4ByFilter` SHALL be added to the CVM service that:
1. Accepts a `paramMap map[string]interface{}`
2. Creates a DescribeImagesRequest from the paramMap
3. Handles ImageIds, Filters, and InstanceType from paramMap
4. Handles pagination with Limit=100
5. Returns the complete list of Image objects

#### Scenario: Service method with ImageIds
- **WHEN** paramMap contains "ImageIds" key
- **THEN** the request SHALL set ImageIds to the provided values

#### Scenario: Service method with Filters
- **WHEN** paramMap contains "Filters" key
- **THEN** the request SHALL set Filters to the provided SDK Filter structs

#### Scenario: Service method with InstanceType
- **WHEN** paramMap contains "InstanceType" key
- **THEN** the request SHALL set InstanceType to the provided value

### Requirement: Provider registration
The datasource SHALL be registered in `tencentcloud/provider.go` with the name `"tencentcloud_cvm_image_4"` and documented in `tencentcloud/provider.md`.

#### Scenario: Datasource is available in Terraform
- **WHEN** the provider is loaded
- **THEN** the datasource `tencentcloud_cvm_image_4` SHALL be available for use in Terraform configurations

### Requirement: Unit tests with gomonkey
Unit tests SHALL be created using gomonkey mock approach (not Terraform test suite) with:
- Mock for the CVM DescribeImages API call
- Test cases for querying by image_ids, filters, and instance_type
- Tests run with `go test -gcflags=all=-l`

#### Scenario: Unit test for read by image_ids
- **WHEN** the read function is called with image_ids set
- **THEN** the test SHALL verify that DescribeImages is called with correct ImageIds and the results are properly mapped

#### Scenario: Unit test for read by filters
- **WHEN** the read function is called with filters set
- **THEN** the test SHALL verify that DescribeImages is called with correct Filters and the results are properly mapped
