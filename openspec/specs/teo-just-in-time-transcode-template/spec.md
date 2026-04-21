## ADDED Requirements

### Requirement: Resource schema definition
The resource `tencentcloud_teo_just_in_time_transcode_template` SHALL define the following schema fields:

- `zone_id` (Required, ForceNew, string): 站点ID
- `template_name` (Required, ForceNew, string): 即时转码模板名称
- `comment` (Optional, ForceNew, string): 模板描述信息
- `video_stream_switch` (Optional, ForceNew, string): 启用视频流开关，取值 on/off
- `audio_stream_switch` (Optional, ForceNew, string): 启用音频流开关，取值 on/off
- `video_template` (Optional, ForceNew, TypeList, MaxItems: 1): 视频流配置参数
- `audio_template` (Optional, ForceNew, TypeList, MaxItems: 1): 音频流配置参数
- `template_id` (Computed, string): 即时转码模板唯一标识
- `create_time` (Computed, string): 模板创建时间
- `update_time` (Computed, string): 模板最后修改时间
- `type` (Computed, string): 模板类型

The `video_template` nested block SHALL contain:
- `codec` (Optional, ForceNew, string): 视频流编码格式，可选值 H.264/H.265
- `fps` (Optional, ForceNew, float): 视频帧率，取值范围 [0, 30]
- `bitrate` (Optional, ForceNew, int): 视频流码率，取值 0 和 [128, 10000]
- `resolution_adaptive` (Optional, ForceNew, string): 分辨率自适应，可选值 open/close
- `width` (Optional, ForceNew, int): 视频流宽度最大值，取值 0 和 [128, 1920]
- `height` (Optional, ForceNew, int): 视频流高度最大值，取值 0 和 [128, 1080]
- `fill_type` (Optional, ForceNew, string): 填充方式，可选值 stretch/black/white/gauss

The `audio_template` nested block SHALL contain:
- `codec` (Optional, ForceNew, string): 音频流编码格式，可选值 libfdk_aac
- `audio_channel` (Optional, ForceNew, int): 音频通道数，可选值 2

#### Scenario: Schema fields match cloud API parameters
- **WHEN** the resource schema is defined
- **THEN** all Create API input parameters have corresponding schema fields with ForceNew set
- **AND** the Create API output parameter TemplateId is mapped to a Computed field template_id
- **AND** the Describe API output fields (CreateTime, UpdateTime, Type) are mapped to Computed fields

#### Scenario: No Update function is defined
- **WHEN** the resource is defined
- **THEN** the schema has no Update function (only Create, Read, Delete)
- **AND** all non-Computed fields have ForceNew set to true

### Requirement: Resource Create operation
The resource SHALL implement the Create function that calls `CreateJustInTimeTranscodeTemplate` cloud API.

- The composite resource ID SHALL be `zoneId#templateId` (joined with `tccommon.FILED_SP`)
- After successful creation, the Read function SHALL be called to refresh the state
- The function SHALL use `tccommon.ReadRetryTimeout` for retry logic
- If the API call fails, the error SHALL be wrapped with `tccommon.RetryError()`

#### Scenario: Successful creation of transcode template
- **WHEN** `tencentcloud_teo_just_in_time_transcode_template` resource is created with valid parameters
- **THEN** the Create function calls `CreateJustInTimeTranscodeTemplate` with the correct parameters
- **AND** the resource ID is set to `zoneId#templateId`
- **AND** the Read function is called to refresh the state

#### Scenario: Creation with video and audio templates
- **WHEN** the resource is created with video_template and audio_template blocks
- **THEN** VideoTemplateInfo and AudioTemplateInfo are properly constructed from the nested schema
- **AND** the API request includes both VideoTemplate and AudioTemplate parameters

### Requirement: Resource Read operation
The resource SHALL implement the Read function that calls `DescribeJustInTimeTranscodeTemplates` cloud API.

- The composite ID SHALL be split by `tccommon.FILED_SP` to get zoneId and templateId
- The Describe request SHALL use Filter with Name="template-id" and Values=[templateId] to query the specific template
- The Limit SHALL be set to 1000 (max value per API spec) and Offset to 0
- If the template is not found in the response, `d.SetId("")` SHALL be called to mark the resource as gone
- The function SHALL use `tccommon.ReadRetryTimeout` for retry logic

#### Scenario: Successful read of existing template
- **WHEN** the Read function is called for an existing resource
- **THEN** the Describe API is called with ZoneId and Filter template-id
- **AND** the template data is found in TemplateSet
- **AND** the state is updated with template_id, template_name, comment, type, video_stream_switch, audio_stream_switch, video_template, audio_template, create_time, update_time

#### Scenario: Resource not found during read
- **WHEN** the Read function is called but the template no longer exists
- **THEN** the Describe API returns an empty TemplateSet
- **AND** `d.SetId("")` is called to mark the resource as removed

### Requirement: Resource Delete operation
The resource SHALL implement the Delete function that calls `DeleteJustInTimeTranscodeTemplates` cloud API.

- The composite ID SHALL be split by `tccommon.FILED_SP` to get zoneId and templateId
- The Delete request SHALL include ZoneId and TemplateIds=[templateId]
- The function SHALL use `tccommon.ReadRetryTimeout` for retry logic

#### Scenario: Successful deletion of transcode template
- **WHEN** the resource is destroyed
- **THEN** the Delete function calls `DeleteJustInTimeTranscodeTemplates` with ZoneId and TemplateIds
- **AND** the resource is removed from the Terraform state

### Requirement: Resource Import support
The resource SHALL support import via `schema.ImportStatePassthrough`.

- The import ID format SHALL be `zoneId#templateId`
- After import, the Read function SHALL be called to populate the state

#### Scenario: Import existing template
- **WHEN** `terraform import tencentcloud_teo_just_in_time_transcode_template.example zoneId#templateId` is executed
- **THEN** the resource ID is set to `zoneId#templateId`
- **AND** the Read function populates all state attributes

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` and listed in `tencentcloud/provider.md`.

- The resource key SHALL be `tencentcloud_teo_just_in_time_transcode_template`
- The resource function SHALL be `teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()`

#### Scenario: Resource is available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_just_in_time_transcode_template` is available as a resource type
- **AND** it appears in the provider documentation

### Requirement: Unit tests with gomonkey mock
The resource SHALL have unit tests using gomonkey to mock cloud API calls.

- Tests SHALL cover Create, Read, and Delete operations
- Tests SHALL NOT use Terraform test suite (no TestAcc prefix)
- Tests SHALL use `go test` to run

#### Scenario: Unit test for Create operation
- **WHEN** unit tests are executed
- **THEN** Create function is tested with mocked API response returning TemplateId
- **AND** the resource ID is correctly set to `zoneId#templateId`

#### Scenario: Unit test for Read operation
- **WHEN** unit tests are executed
- **THEN** Read function is tested with mocked Describe API response
- **AND** state fields are correctly populated from the response

#### Scenario: Unit test for Delete operation
- **WHEN** unit tests are executed
- **THEN** Delete function is tested with mocked API call
- **AND** the Delete request contains correct ZoneId and TemplateIds

### Requirement: Documentation file
The resource SHALL have a .md documentation file following the project's documentation format.

- The file SHALL be located at `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.md`
- The first line SHALL be a one-sentence description mentioning TEO
- The file SHALL include an Example Usage section with HCL configuration
- The file SHALL include an Import section with the import format

#### Scenario: Documentation file is created
- **WHEN** the resource is implemented
- **THEN** a .md file exists with description, example usage, and import sections
- **AND** the example includes video_template and audio_template configuration
