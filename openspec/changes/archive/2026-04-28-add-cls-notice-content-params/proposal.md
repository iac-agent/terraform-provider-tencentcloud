## Why

The `tencentcloud_cls_notice_content` resource currently lacks support for the full set of parameters available in the CLS NoticeContent CRUD APIs. Specifically, the resource needs to expose the `notice_content_id` as a computed attribute so users can reference it in their Terraform configurations, and ensure all parameters from the Create/Read/Update/Delete APIs are properly mapped to the Terraform schema.

## What Changes

- Add `notice_content_id` as a computed attribute to the `tencentcloud_cls_notice_content` resource schema, sourced from `CreateNoticeContent` response and `DescribeNoticeContents` response
- Ensure the `name` (Required), `type` (Optional), and `notice_contents` (Optional) parameters are properly supported across all CRUD operations (Create, Read, Update, Delete) as defined by the CLS cloud APIs
- The `notice_contents` parameter supports nested configuration with `type`, `trigger_content`, and `recovery_content` sub-blocks
- Each `trigger_content`/`recovery_content` sub-block supports `title`, `content`, and `headers` fields

## Capabilities

### New Capabilities
- `cls-notice-content-resource`: Full CRUD support for CLS Notice Content resource with all API parameters including name, type, notice_contents, and computed notice_content_id

### Modified Capabilities

## Impact

- `tencentcloud/services/cls/resource_tc_cls_notice_content.go` - Add `notice_content_id` computed attribute, ensure all CRUD parameter mappings
- `tencentcloud/services/cls/resource_tc_cls_notice_content_test.go` - Add unit tests for new parameters
- `tencentcloud/services/cls/resource_tc_cls_notice_content.md` - Update documentation
- `tencentcloud/provider.go` - Resource already registered, verify registration
