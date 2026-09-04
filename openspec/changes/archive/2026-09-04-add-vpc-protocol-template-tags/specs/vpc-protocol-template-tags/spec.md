## ADDED Requirements

### Requirement: Tags on protocol template creation
The `tencentcloud_protocol_template` resource SHALL support an optional `tags` parameter (TypeMap, where map keys are tag keys and map values are tag values) that is passed to the `CreateServiceTemplate` API as `Tags` (a list of `*vpc.Tag`, each with `Key` and `Value`). When tags are not specified, `Tags` SHALL NOT be set in the API request.

#### Scenario: Create protocol template with tags
- **WHEN** a user specifies `tags = { "env" = "prod", "team" = "infra" }` in the `tencentcloud_protocol_template` resource configuration
- **THEN** the provider SHALL pass `Tags` containing both tag key-value pairs in the `CreateServiceTemplate` API request

#### Scenario: Create protocol template without tags
- **WHEN** a user does NOT specify `tags` in the `tencentcloud_protocol_template` resource configuration
- **THEN** the provider SHALL NOT set `Tags` in the `CreateServiceTemplate` API request

### Requirement: Read protocol template tags
The provider SHALL refresh the `tags` parameter of `tencentcloud_protocol_template` from the `DescribeServiceTemplates` API response (`ServiceTemplateSet[].TagSet`). Each tag's `Key` and `Value` SHALL be read into the tags map. The provider SHALL guard against nil `Key`/`Value` before populating the map.

#### Scenario: Read existing protocol template with tags
- **WHEN** the provider reads an existing `tencentcloud_protocol_template` resource that has tags bound
- **THEN** `tags` SHALL be refreshed from the `DescribeServiceTemplates` API response (`TagSet[].Key` and `TagSet[].Value`)

#### Scenario: Read existing protocol template without tags
- **WHEN** the provider reads an existing `tencentcloud_protocol_template` resource that has no tags bound
- **THEN** `tags` SHALL be empty and the provider SHALL NOT error

### Requirement: Update protocol template tags via shared tag service
Because the `ModifyServiceTemplateAttribute` API does not accept a `Tags` field, the provider SHALL reconcile tag changes on Update using the shared tag service. When `tags` changes, the provider SHALL compute the diff with `svctag.DiffTags` and apply it with `tagService.ModifyTags` using a resource name built from `tccommon.BuildTagResourceName("vpc", "service", region, d.Id())`.

#### Scenario: Add a new tag on update
- **WHEN** a user adds a new tag key-value pair to the `tags` map on an existing `tencentcloud_protocol_template`
- **THEN** the provider SHALL reconcile the tags via the shared tag service so the new tag is bound to the protocol template

#### Scenario: Remove a tag on update
- **WHEN** a user removes a tag key from the `tags` map on an existing `tencentcloud_protocol_template`
- **THEN** the provider SHALL reconcile the tags via the shared tag service so the removed tag is unbound from the protocol template

#### Scenario: Tags unchanged on update
- **WHEN** a user updates `name` or `protocols` but does NOT change `tags`
- **THEN** the provider SHALL NOT call the shared tag service for tags and SHALL only call `ModifyServiceTemplateAttribute` for the changed fields
