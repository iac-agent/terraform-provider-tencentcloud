## ADDED Requirements

### Requirement: Tags parameter in protocol template schema
The `tencentcloud_protocol_template` resource SHALL include a `tags` parameter of type `schema.TypeMap` with `Optional: true`. Each key-value pair represents a tag where the key and value are both strings.

#### Scenario: Creating a protocol template with tags
- **WHEN** a user creates a `tencentcloud_protocol_template` resource with the `tags` parameter specified (e.g., `tags = { "env" = "prod" }`)
- **THEN** the resource SHALL pass the tags to the `CreateServiceTemplate` API request's `Tags` field as a slice of `vpc.Tag` objects, and the tags SHALL be persisted on the created protocol template

#### Scenario: Creating a protocol template without tags
- **WHEN** a user creates a `tencentcloud_protocol_template` resource without specifying the `tags` parameter
- **THEN** the resource SHALL create the protocol template without tags, and the `Tags` field SHALL NOT be set on the API request

### Requirement: Reading tags from DescribeServiceTemplates
The `tencentcloud_protocol_template` resource's read function SHALL read tags from the `TagSet` field of the `ServiceTemplate` object returned by `DescribeServiceTemplates`, and set them to the `tags` attribute in the Terraform state.

#### Scenario: Reading a protocol template with tags
- **WHEN** the read function is called for a protocol template that has tags
- **THEN** the function SHALL convert the `TagSet` (slice of `vpc.Tag` with `Key` and `Value` fields) to a `map[string]interface{}` and call `d.Set("tags", tagMap)`

#### Scenario: Reading a protocol template without tags
- **WHEN** the read function is called for a protocol template that has no tags (TagSet is nil or empty)
- **THEN** the function SHALL NOT call `d.Set("tags", ...)` for the tags field, leaving it as the default empty value

### Requirement: Updating tags via tag service
The `tencentcloud_protocol_template` resource's update function SHALL support tag changes using the tag service (`svctag.TagService`), since the `ModifyServiceTemplateAttribute` API does not support tags.

#### Scenario: Updating tags on a protocol template
- **WHEN** a user updates the `tags` parameter of a `tencentcloud_protocol_template` resource
- **THEN** the update function SHALL detect the change via `d.HasChange("tags")`, compute the tags to add and remove using `svctag.DiffTags()`, and call `tagService.ModifyTags()` with the appropriate resource name

#### Scenario: No tag change during update
- **WHEN** a user updates other parameters of a `tencentcloud_protocol_template` resource but not `tags`
- **THEN** the update function SHALL NOT call the tag service for tag modifications
