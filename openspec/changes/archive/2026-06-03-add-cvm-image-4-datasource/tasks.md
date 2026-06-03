## 1. Service Layer

- [x] 1.1 Add `DescribeCvmImage4ByFilter` method to `tencentcloud/services/cvm/service_tencentcloud_cvm.go` that accepts `paramMap map[string]interface{}`, creates DescribeImagesRequest with ImageIds/Filters/InstanceType from paramMap, handles pagination with Limit=100, and returns `[]*cvm.Image`

## 2. Datasource Schema and Read Function

- [x] 2.1 Create `tencentcloud/services/cvm/data_source_tc_cvm_image_4.go` with schema definition including input parameters: `image_ids` (TypeList of TypeString, Optional), `filters` (TypeList of nested Resource with `name` and `values`, Optional), `instance_type` (TypeString, Optional), `result_output_file` (TypeString, Optional), and computed output `image_set` (TypeList of nested Resource with all Image struct fields)
- [x] 2.2 Implement `dataSourceTencentCloudCvmImage4Read` function following the igtm_instance_list pattern: use defer LogElapsed/InconsistentCheck, build paramMap from schema, call service method inside resource.Retry with RetryError, map response fields with nil checks, set ID with helper.BuildToken(), handle result_output_file

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_cvm_image_4` datasource in `tencentcloud/provider.go`
- [x] 3.2 Add `tencentcloud_cvm_image_4` entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/cvm/data_source_tc_cvm_image_4.md` with description, example usage (query by image_ids and query by filters), following the documentation format from gendoc/README.md and other datasource .md files

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/cvm/data_source_tc_cvm_image_4_test.go` with gomonkey mock tests for: read by image_ids, read by filters, read by instance_type, handling nil fields in response
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` and ensure all tests pass
