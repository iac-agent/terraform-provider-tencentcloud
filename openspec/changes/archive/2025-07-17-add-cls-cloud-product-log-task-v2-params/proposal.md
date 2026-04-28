## Why

The `tencentcloud_cls_cloud_product_log_task_v2` resource currently lacks support for tags, `is_delete_topic`, `is_delete_logset`, and `status` parameters. The cloud API `CreateCloudProductLogCollection` supports a `Tags` field that allows binding tags to the log topic during creation, but this is not exposed in the Terraform resource. Similarly, the `DeleteCloudProductLogCollection` API supports `IsDeleteTopic` and `IsDeleteLogset` fields for more granular control over cascading deletion, replacing the current terraform-only `force_delete` approach. The `Status` field from the create and delete responses is also not exposed as a computed attribute. These additions will bring the Terraform resource to parity with the cloud API capabilities.

## What Changes

- Add `tags` parameter (TypeList of Key/Value objects) to the resource schema, mapped to `CreateCloudProductLogCollectionRequest.Tags`
- Add `is_delete_topic` parameter (TypeBool, Optional) to the resource schema, mapped to `DeleteCloudProductLogCollectionRequest.IsDeleteTopic`
- Add `is_delete_logset` parameter (TypeBool, Optional) to the resource schema, mapped to `DeleteCloudProductLogCollectionRequest.IsDeleteLogset`
- Add `status` computed parameter (TypeInt) to the resource schema, mapped to `CreateCloudProductLogCollectionResponse.Response.Status` and `DeleteCloudProductLogCollectionResponse.Response.Status`
- Update Create function to pass `tags` to the API request and read `status` from the response
- Update Delete function to pass `is_delete_topic` and `is_delete_logset` to the API request
- Update Read function to read `status` from the describe response
- Update the resource .md documentation file with new parameter examples

## Capabilities

### New Capabilities
- `cls-cloud-product-log-task-v2-params`: Adds tags, is_delete_topic, is_delete_logset, and status parameters to the tencentcloud_cls_cloud_product_log_task_v2 resource

### Modified Capabilities

## Impact

- Affected files:
  - `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go` - Schema and CRUD logic changes
  - `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2_test.go` - Unit test updates
  - `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.md` - Documentation updates
  - `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_extension.go` - May need extension updates
- APIs: `CreateCloudProductLogCollection`, `DeleteCloudProductLogCollection`, `DescribeCloudProductLogTasks`
- Backward compatible: All new fields are Optional or Computed, existing configurations remain valid
