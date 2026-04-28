## 1. Schema and CRUD Implementation

- [x] 1.1 Add `notice_content_id` computed attribute to the resource schema in `tencentcloud/services/cls/resource_tc_cls_notice_content.go`
- [x] 1.2 Update `resourceTencentCloudClsNoticeContentCreate` to set `notice_content_id` from `response.Response.NoticeContentId` after creation
- [x] 1.3 Update `resourceTencentCloudClsNoticeContentRead` to set `notice_content_id` from `respData.NoticeContentId` with nil check
- [x] 1.4 Verify `resourceTencentCloudClsNoticeContentUpdate` properly passes all parameters to ModifyNoticeContent API
- [x] 1.5 Verify `resourceTencentCloudClsNoticeContentDelete` properly passes NoticeContentId to DeleteNoticeContent API

## 2. Service Layer

- [x] 2.1 Verify `DescribeClsNoticeContentById` in `tencentcloud/services/cls/service_tencentcloud_cls.go` returns `NoticeContentId` field in the response

## 3. Tests

- [x] 3.1 Add unit test for `notice_content_id` computed attribute in `tencentcloud/services/cls/resource_tc_cls_notice_content_test.go` using gomonkey mock

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/cls/resource_tc_cls_notice_content.md` to reflect the `notice_content_id` computed attribute
