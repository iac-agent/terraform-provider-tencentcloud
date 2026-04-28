## Context

The `tencentcloud_cls_cloud_product_log_task_v2` resource manages CLS cloud product log collection tasks. Currently, the resource supports creating, reading, updating, and deleting log collection tasks, but is missing several parameters that the cloud API already supports:

1. **Tags**: The `CreateCloudProductLogCollection` API accepts a `Tags` field (type `[]*Tag` with Key/Value) that allows binding tags to the log topic during creation. This is not exposed in the Terraform resource schema.

2. **is_delete_topic / is_delete_logset**: The `DeleteCloudProductLogCollection` API supports `IsDeleteTopic` (bool) and `IsDeleteLogset` (bool) fields that control cascading deletion of associated topic/logset through the API itself. Currently, the resource uses a terraform-only `force_delete` parameter that triggers separate `DeleteTopic` and `DeleteLogset` API calls after the main delete.

3. **Status**: Both `CreateCloudProductLogCollection` and `DeleteCloudProductLogCollection` API responses include a `Status` field (int64) that indicates the task state. This is not exposed as a computed attribute in the Terraform resource.

Current state:
- Schema has: `instance_id`, `assumer_name`, `log_type`, `cloud_product_region`, `cls_region`, `logset_name`, `topic_name`, `extend`, `logset_id`, `topic_id`, `force_delete`
- Create function passes all existing parameters to the API but not `Tags`
- Delete function uses `force_delete` to make separate DeleteTopic/DeleteLogset calls instead of using API's `IsDeleteTopic`/`IsDeleteLogset`
- Status is not read from any API response

## Goals / Non-Goals

**Goals:**
- Add `tags` parameter to the resource schema (TypeList of Key/Value) and pass it to `CreateCloudProductLogCollectionRequest`
- Add `is_delete_topic` and `is_delete_logset` parameters to the resource schema and pass them to `DeleteCloudProductLogCollectionRequest`
- Add `status` computed parameter to the resource schema
- Maintain full backward compatibility with existing Terraform configurations

**Non-Goals:**
- Removing or deprecating the existing `force_delete` parameter (it still controls the separate DeleteTopic/DeleteLogset behavior)
- Adding tags support for the Modify operation (the `ModifyCloudProductLogCollection` API does not support tags)
- Changing the resource ID format

## Decisions

### 1. Tags Schema Design
**Decision**: Use `TypeList` with `MaxItems: 10` and nested `Key`/`Value` string fields (both required), consistent with how other TencentCloud Terraform resources define tags (e.g., `tencentcloud_cls_alarm_notice`).

**Rationale**: The cloud API defines `Tags` as `[]*Tag` with `Key` and `Value` fields, and limits to 10 tag key-value pairs. Using `TypeList` matches the API structure and is consistent with existing patterns in this provider.

### 2. is_delete_topic / is_delete_logset vs force_delete
**Decision**: Add `is_delete_topic` and `is_delete_logset` as new Optional bool parameters alongside the existing `force_delete` parameter. When `is_delete_topic` is set, pass it to the `DeleteCloudProductLogCollection` API request. The existing `force_delete` behavior (separate DeleteTopic/DeleteLogset calls) remains unchanged.

**Rationale**: The `is_delete_topic`/`is_delete_logset` parameters are passed to the delete API itself, which handles the cascading deletion server-side. The `force_delete` parameter controls additional client-side deletion calls. They serve different purposes and both can coexist. Keeping `force_delete` ensures backward compatibility.

### 3. Status as Computed Attribute
**Decision**: Add `status` as a computed integer attribute (TypeInt, Computed: true) that is read from the create response and the describe response.

**Rationale**: Status is returned by the API but should not be user-configurable. Making it computed allows users to reference the current state of the task in their Terraform configurations.

### 4. Tags in Read Operation
**Decision**: Tags will be read from the describe response if the `CloudProductLogTaskInfo` struct includes tags, or left as-is if the API does not return tags in the describe response.

**Rationale**: Need to verify whether `DescribeCloudProductLogTasks` returns tag information. If not, tags will only be set during creation and not refreshed during read, which is acceptable for create-only tags.

## Risks / Trade-offs

- [Risk] Tags may not be returned by `DescribeCloudProductLogTasks` API → Mitigation: Check the API response struct; if tags aren't returned, the read operation won't update tags (they remain as configured), which is acceptable for create-only tags
- [Risk] Adding `is_delete_topic`/`is_delete_logset` alongside `force_delete` may confuse users about which to use → Mitigation: Clear documentation explaining the difference between the two approaches
- [Risk] `status` value interpretation differs between create and delete responses → Mitigation: Document the status values clearly (-1/0: creating, 1: created, 2: deleting, 3: deleted)
