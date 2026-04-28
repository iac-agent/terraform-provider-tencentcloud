## ADDED Requirements

### Requirement: Cls notice content resource CRUD operations
The system SHALL provide a Terraform resource `tencentcloud_cls_notice_content` that manages the full lifecycle of a CLS Notice Content template through Create, Read, Update, and Delete operations using the CLS cloud APIs.

#### Scenario: Create notice content with all parameters
- **WHEN** a user creates a `tencentcloud_cls_notice_content` resource with `name` (required string), `type` (optional int, 0=Chinese, 1=English), and `notice_contents` (optional list)
- **THEN** the system SHALL call `CreateNoticeContent` API with the provided parameters and set the resource ID from `response.Response.NoticeContentId`

#### Scenario: Read notice content by ID
- **WHEN** the system reads a `tencentcloud_cls_notice_content` resource
- **THEN** the system SHALL call `DescribeNoticeContents` API with a filter on `noticeContentId` matching the resource ID, and populate all schema fields from the response `NoticeContentTemplate`

#### Scenario: Update notice content mutable fields
- **WHEN** a user updates `name`, `type`, or `notice_contents` fields
- **THEN** the system SHALL call `ModifyNoticeContent` API with `NoticeContentId`, `Name`, `Type`, and `NoticeContents` parameters

#### Scenario: Delete notice content
- **WHEN** a user deletes a `tencentcloud_cls_notice_content` resource
- **THEN** the system SHALL call `DeleteNoticeContent` API with `NoticeContentId` parameter

### Requirement: Notice content ID computed attribute
The system SHALL expose `notice_content_id` as a computed string attribute on the `tencentcloud_cls_notice_content` resource, sourced from `CreateNoticeContent` response and `DescribeNoticeContents` response.

#### Scenario: Notice content ID is set after creation
- **WHEN** a `tencentcloud_cls_notice_content` resource is created
- **THEN** the `notice_content_id` attribute SHALL be populated from the `response.Response.NoticeContentId` field

#### Scenario: Notice content ID is read from API
- **WHEN** a `tencentcloud_cls_notice_content` resource is read
- **THEN** the `notice_content_id` attribute SHALL be populated from the `NoticeContentTemplate.NoticeContentId` field in the DescribeNoticeContents response

### Requirement: Notice contents nested block schema
The `notice_contents` attribute SHALL be a TypeList with nested schema containing `type` (required string), `trigger_content` (optional block), and `recovery_content` (optional block).

#### Scenario: Notice contents with trigger and recovery content
- **WHEN** a user specifies `notice_contents` with `type`, `trigger_content`, and `recovery_content` sub-blocks
- **THEN** the system SHALL properly map all nested fields to the `NoticeContent` SDK struct for Create and Update operations

#### Scenario: Trigger content and recovery content sub-blocks
- **WHEN** a user specifies `trigger_content` or `recovery_content` sub-blocks
- **THEN** each sub-block SHALL support `title` (optional string), `content` (optional string), and `headers` (optional set of strings) fields, mapped to the `NoticeContentInfo` SDK struct

### Requirement: Resource import support
The `tencentcloud_cls_notice_content` resource SHALL support import using the notice content ID.

#### Scenario: Import existing notice content
- **WHEN** a user imports a `tencentcloud_cls_notice_content` resource using its ID
- **THEN** the system SHALL read the resource state from the DescribeNoticeContents API and populate all schema fields
