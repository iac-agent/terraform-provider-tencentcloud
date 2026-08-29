## ADDED Requirements

### Requirement: Tags parameter passed directly in CreateAlarmNotice request
The system SHALL pass the `tags` parameter directly through `CreateAlarmNoticeRequest.Tags` when creating a CLS alarm notice resource, converting from the Terraform `TypeMap` schema format to `[]*cls.Tag` format.

#### Scenario: Create alarm notice with tags
- **WHEN** user creates a `tencentcloud_cls_alarm_notice` resource with `tags` configured
- **THEN** the system SHALL convert the tags map to `[]*cls.Tag` format and set `CreateAlarmNoticeRequest.Tags` before calling the CLS API

#### Scenario: Create alarm notice without tags
- **WHEN** user creates a `tencentcloud_cls_alarm_notice` resource without `tags`
- **THEN** the system SHALL NOT set `CreateAlarmNoticeRequest.Tags` and proceed with the create request

### Requirement: Tags parameter passed directly in ModifyAlarmNotice request
The system SHALL pass the `tags` parameter directly through `ModifyAlarmNoticeRequest.Tags` when updating a CLS alarm notice resource, converting from the Terraform `TypeMap` schema format to `[]*cls.Tag` format.

#### Scenario: Update alarm notice tags
- **WHEN** user updates the `tags` field of a `tencentcloud_cls_alarm_notice` resource
- **THEN** the system SHALL convert the updated tags map to `[]*cls.Tag` format and set `ModifyAlarmNoticeRequest.Tags` before calling the CLS API

### Requirement: Tags read from AlarmNotice response
The system SHALL read the `tags` parameter from `AlarmNotice.Tags` in the `DescribeAlarmNotices` API response during the Read operation, converting from `[]*cls.Tag` format to Terraform `TypeMap` schema format.

#### Scenario: Read alarm notice with tags
- **WHEN** the system reads a `tencentcloud_cls_alarm_notice` resource that has tags
- **THEN** the system SHALL convert `AlarmNotice.Tags` from `[]*cls.Tag` to `map[string]interface{}` and set it in the Terraform state

#### Scenario: Read alarm notice without tags
- **WHEN** the system reads a `tencentcloud_cls_alarm_notice` resource where `AlarmNotice.Tags` is nil
- **THEN** the system SHALL skip setting tags in the Terraform state and fall back to the tag service for tag reading

### Requirement: Tag format conversion helpers
The system SHALL provide helper functions to convert between `map[string]interface{}` (Terraform TypeMap) and `[]*cls.Tag` (CLS API format).

#### Scenario: Convert TypeMap to cls.Tag slice
- **WHEN** converting a Terraform `tags` map `{"key1": "value1", "key2": "value2"}` to CLS API format
- **THEN** the result SHALL be `[]*cls.Tag{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}}`

#### Scenario: Convert cls.Tag slice to TypeMap
- **WHEN** converting `[]*cls.Tag{{Key: "key1", Value: "value1"}}` to Terraform map format
- **THEN** the result SHALL be `map[string]interface{}{"key1": "value1"}`
