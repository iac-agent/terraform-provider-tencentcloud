## 1. Schema Definition

- [x] 1.1 Add `tags` parameter to `ResourceTencentCloudClsCloudProductLogTaskV2` schema as TypeMap, consistent with other CLS resources (cls_topic, cls_dashboard, etc.)
- [x] 1.2 Add `is_delete_topic` parameter to schema as TypeBool, Optional, default false
- [x] 1.3 Add `is_delete_logset` parameter to schema as TypeBool, Optional, default false
- [x] 1.4 Add `status` parameter to schema as TypeInt, Computed

## 2. Create Function

- [x] 2.1 Add tags handling in `resourceTencentCloudClsCloudProductLogTaskV2Create`: convert schema tags map to `[]*clsv20201016.Tag` and set `request.Tags`
- [x] 2.2 Read `status` from `CreateCloudProductLogCollection` response and set it in the resource state after successful creation

## 3. Read Function

- [x] 3.1 Read `status` from `DescribeCloudProductLogTasks` response (`respData.Tasks[0].Status`) and set it in the resource state (with nil check)
- [x] 3.2 Read `tags` from the describe response if available - CloudProductLogTaskInfo does NOT have Tags field, so tags cannot be refreshed from describe response (create-only tags, preserved in state)

## 4. Delete Function

- [x] 4.1 Add `is_delete_topic` handling in `resourceTencentCloudClsCloudProductLogTaskV2Delete`: set `request.IsDeleteTopic` when the parameter is provided
- [x] 4.2 Add `is_delete_logset` handling in `resourceTencentCloudClsCloudProductLogTaskV2Delete`: set `request.IsDeleteLogset` when the parameter is provided

## 5. Update Function

- [x] 5.1 Add `tags` to the `immutableArgs` list in `resourceTencentCloudClsCloudProductLogTaskV2Update` since `ModifyCloudProductLogCollection` does not support tags
- [x] 5.2 Add `is_delete_topic` and `is_delete_logset` to the `immutableArgs` list in the update function

## 6. Unit Tests

- [x] 6.1 Add unit test cases for the new `tags` parameter in `resource_tc_cls_cloud_product_log_task_v2_test.go` using gomonkey mock approach
- [x] 6.2 Add unit test cases for the new `is_delete_topic` and `is_delete_logset` parameters
- [x] 6.3 Add unit test cases for the `status` computed attribute
- [x] 6.4 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## 7. Documentation

- [x] 7.1 Update `resource_tc_cls_cloud_product_log_task_v2.md` with new parameter examples (tags, is_delete_topic, is_delete_logset)
