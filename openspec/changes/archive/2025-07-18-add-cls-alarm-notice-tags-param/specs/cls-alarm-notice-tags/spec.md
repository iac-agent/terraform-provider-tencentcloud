## ADDED Requirements

### Requirement: Tags passed through CLS API on Create
The system SHALL pass the `tags` parameter directly in the `CreateAlarmNotice` API request, converting the Terraform `TypeMap` tags to `[]*cls.Tag` format (each element containing `Key` and `Value` fields). The system SHALL NOT use a separate tag service call after resource creation.

#### Scenario: Create alarm notice with tags
- **WHEN** user creates a `tencentcloud_cls_alarm_notice` resource with tags configured
- **THEN** the system SHALL include the tags in the `CreateAlarmNotice` request's `Tags` field
- **AND** the system SHALL NOT call `tagService.ModifyTags` after creation

#### Scenario: Create alarm notice without tags
- **WHEN** user creates a `tencentcloud_cls_alarm_notice` resource without tags configured
- **THEN** the system SHALL create the alarm notice without tags
- **AND** the `Tags` field in the request SHALL be empty or nil

### Requirement: Tags read from CLS API response on Read
The system SHALL read tags from the `AlarmNotice.Tags` field in the CLS API response during the Read operation, converting `[]*cls.Tag` to `map[string]string` format. The system SHALL NOT use the separate tag service's `DescribeResourceTags` for reading tags.

#### Scenario: Read alarm notice with tags
- **WHEN** the system reads a `tencentcloud_cls_alarm_notice` resource that has tags
- **THEN** the system SHALL extract tags from `AlarmNotice.Tags` in the CLS API response
- **AND** the system SHALL convert each `Tag` element's `Key` and `Value` to a `map[string]string` entry
- **AND** the system SHALL set the result via `d.Set("tags", tagsMap)`

#### Scenario: Read alarm notice without tags
- **WHEN** the system reads a `tencentcloud_cls_alarm_notice` resource that has no tags
- **THEN** the system SHALL skip setting the tags field (AlarmNotice.Tags is nil)

### Requirement: Tags passed through CLS API on Update
The system SHALL pass the `tags` parameter directly in the `ModifyAlarmNotice` API request when tags have changed, converting the Terraform `TypeMap` tags to `[]*cls.Tag` format. The system SHALL NOT use the separate tag service's `ModifyTags` for updating tags.

#### Scenario: Update alarm notice tags
- **WHEN** user updates the `tags` field of a `tencentcloud_cls_alarm_notice` resource
- **THEN** the system SHALL include the updated tags in the `ModifyAlarmNotice` request's `Tags` field
- **AND** the system SHALL NOT call `tagService.ModifyTags`

#### Scenario: Update alarm notice without tag changes
- **WHEN** user updates other fields of a `tencentcloud_cls_alarm_notice` resource without changing tags
- **THEN** the system SHALL NOT include tags in the `ModifyAlarmNotice` request (unless other fields also changed)
- **AND** the system SHALL detect `d.HasChange("tags")` as false

### Requirement: Tags schema unchanged
The `tags` schema field SHALL remain as `TypeMap`, `Optional`, with `TypeString` elements. No schema migration SHALL be required.

#### Scenario: Existing Terraform configurations
- **WHEN** user has existing Terraform configurations using the `tags` field
- **THEN** the configurations SHALL continue to work without modification
- **AND** the state SHALL be automatically updated on the next refresh/plan
