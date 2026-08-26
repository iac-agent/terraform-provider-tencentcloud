## ADDED Requirements

### Requirement: PT-TAGS-001 - Schema Support for Tags
**Priority**: High  
**Type**: Functional

The `tencentcloud_protocol_template` resource MUST support a `tags` parameter in its schema as a list of key-value pairs.

**Acceptance Criteria**:
- `tags` field is of type `TypeList`
- Each element contains `key` (TypeString, Required) and `value` (TypeString, Optional)
- `tags` field is Optional
- `tags` field is ForceNew (immutable on update)
- Description clearly explains the field's purpose

#### Scenario: Define tags schema field
- **WHEN** a user defines the resource in Terraform configuration
- **THEN** they can specify a `tags` parameter as a list of key-value objects
- **AND** the parameter is optional (resource works without tags)

```hcl
resource "tencentcloud_protocol_template" "example" {
  name      = "test-template"
  protocols = ["tcp:80", "udp:443"]

  tags {
    key   = "Environment"
    value = "production"
  }
  tags {
    key   = "Team"
    value = "networking"
  }
}
```

---

### Requirement: PT-TAGS-002 - Create Protocol Template with Tags
**Priority**: High  
**Type**: Functional

When creating a protocol template, any tags specified MUST be passed to the `CreateServiceTemplate` API's `Tags` parameter.

**Acceptance Criteria**:
- Tags are extracted from Terraform schema
- Tags are converted to `[]*vpc.Tag` format
- Tags are included in `CreateServiceTemplateRequest`
- Template is created with all specified tags
- Tags are immediately visible after creation via read

#### Scenario: Create template with tags
- **WHEN** a Terraform configuration with tags specified is applied
- **THEN** the `CreateServiceTemplate` API is called with the `Tags` parameter
- **AND** the template is created successfully
- **AND** all specified tags are applied to the template
- **AND** tags are retrievable via `DescribeServiceTemplates`

#### Scenario: Create template without tags
- **WHEN** a Terraform configuration without tags specified is applied
- **THEN** the template is created successfully
- **AND** no tags are applied
- **AND** the resource operates normally

---

### Requirement: PT-TAGS-003 - Read Protocol Template Tags
**Priority**: High  
**Type**: Functional

When reading a protocol template's state, tags MUST be retrieved from the `DescribeServiceTemplates` API response (`ServiceTemplate.TagSet`) and stored in Terraform state.

**Acceptance Criteria**:
- Tags are read from `ServiceTemplate.TagSet` in the API response
- Each `Tag` object's `Key` and `Value` fields are extracted
- Tags are set in Terraform state as a list of maps
- Nil or empty `TagSet` is handled gracefully

#### Scenario: Read tags from existing template
- **WHEN** Terraform reads the resource state
- **THEN** the `DescribeServiceTemplates` API returns `TagSet` in the response
- **AND** all tags are extracted and stored in Terraform state
- **AND** `terraform show` displays all tags correctly

#### Scenario: Handle template with no tags
- **WHEN** a protocol template exists without any tags
- **THEN** `TagSet` is nil or empty
- **AND** the `tags` field in state is set to empty list
- **AND** no errors are raised

#### Scenario: Handle API response with nil TagSet
- **WHEN** the API returns a template with `TagSet` set to nil
- **THEN** the `tags` field in state is set to empty list
- **AND** no nil pointer dereference errors occur

---

### Requirement: PT-TAGS-004 - Tags Immutability
**Priority**: High  
**Type**: Functional

Tags MUST be immutable on update. Any change to tags MUST trigger resource recreation.

**Acceptance Criteria**:
- `tags` field is marked as `ForceNew: true`
- Changing tags causes `terraform plan` to show resource recreation
- The update function checks for tag changes and returns an error if detected

#### Scenario: Attempt to modify tags
- **WHEN** a user changes tags on an existing protocol template
- **THEN** Terraform plans to recreate the resource
- **AND** a clear message indicates tags are immutable

---

### Requirement: PT-TAGS-005 - Backward Compatibility
**Priority**: High  
**Type**: Non-Functional

Adding tags support MUST NOT break existing protocol template resources or configurations.

**Acceptance Criteria**:
- Existing templates without tags continue to work
- Existing Terraform configurations without tags remain valid
- No forced replacement of existing resources
- Tags field is purely additive

#### Scenario: Existing template without tags
- **WHEN** a protocol template was created before tags support
- **AND** Terraform refreshes state
- **THEN** the template remains functional
- **AND** tags field shows empty list
- **AND** no forced replacement occurs

#### Scenario: Existing configuration without tags
- **WHEN** a Terraform configuration without `tags` field is applied with the new provider version
- **THEN** the configuration remains valid
- **AND** the template operates normally
- **AND** no warnings or errors occur