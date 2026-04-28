## Context

The `tencentcloud_cls_notice_content` resource already exists in the Terraform provider with basic CRUD support for `name`, `type`, and `notice_contents` parameters. The resource uses the CLS SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`) for API calls. The current implementation lacks the `notice_content_id` computed attribute, which is needed for cross-resource references in Terraform configurations.

The CLS NoticeContent APIs provide:
- **CreateNoticeContent**: Creates a notice content template, returns `NoticeContentId`
- **DescribeNoticeContents**: Queries notice content templates by filters, returns `NoticeContentTemplate` list
- **ModifyNoticeContent**: Updates name, type, and notice_contents fields
- **DeleteNoticeContent**: Deletes by NoticeContentId

The `NoticeContentTemplate` response struct includes additional computed fields: `Flag`, `Uin`, `SubUin`, `CreateTime`, `UpdateTime`.

## Goals / Non-Goals

**Goals:**
- Add `notice_content_id` as a computed attribute so users can reference it in their Terraform configurations
- Ensure all CRUD operations properly map parameters between Terraform schema and cloud API
- Maintain backward compatibility with existing configurations

**Non-Goals:**
- Adding computed fields from the response (`Flag`, `Uin`, `SubUin`, `CreateTime`, `UpdateTime`) is out of scope for this change
- Modifying the data source `tencentcloud_cls_notice_contents` is out of scope
- Changing the resource ID format is out of scope

## Decisions

1. **Add `notice_content_id` as Computed attribute**: The `notice_content_id` field will be added as a `Computed: true, Type: schema.TypeString` attribute. This allows users to reference the ID in other resources without needing to use the resource ID directly. The value is sourced from `response.Response.NoticeContentId` on create and `NoticeContentTemplate.NoticeContentId` on read.

2. **Keep existing schema fields unchanged**: The existing `name` (Required), `type` (Optional), and `notice_contents` (Optional) fields remain as-is to maintain backward compatibility.

3. **Reuse existing service layer**: The `DescribeClsNoticeContentById` method already exists in the service layer and returns `*cls.NoticeContentTemplate`, which contains all needed fields including `NoticeContentId`.

4. **Read operation populates notice_content_id**: In the Read function, after calling `DescribeClsNoticeContentById`, the `notice_content_id` field is set from `respData.NoticeContentId` if it is not nil, following the nil-check pattern used for other fields.

## Risks / Trade-offs

- [Risk] Adding `notice_content_id` as a computed field duplicates the resource ID → Mitigation: This is a common pattern in Terraform providers where the API-returned ID differs from the resource ID format, and it enables cleaner cross-resource references
- [Risk] Existing state files don't have `notice_content_id` populated → Mitigation: As a computed field, Terraform will automatically populate it on the next read/refresh
