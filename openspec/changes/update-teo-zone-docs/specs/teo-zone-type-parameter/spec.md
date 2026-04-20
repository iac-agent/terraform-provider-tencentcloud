## MODIFIED Requirements

### Requirement: Type parameter documentation
The `type` parameter in `tencentcloud_teo_zone` resource SHALL document all valid access types in its schema description, including `partial` (CNAME access), `full` (NS access), `noDomainAccess` (No domain access), `dnsPodAccess` (DNSPod access), and `ai` (AI access).

#### Scenario: User reads schema documentation
- **WHEN** user views the `type` parameter description in schema or generated documentation
- **THEN** all six access types are clearly listed with their descriptions
- **AND** the default value (`partial`) is mentioned

#### Scenario: User specifies dnsPodAccess type
- **WHEN** user sets `type = "dnsPodAccess"` in their Terraform configuration
- **THEN** the configuration is accepted
- **AND** the resource is created with DNSPod access type

#### Scenario: User specifies ai type
- **WHEN** user sets `type = "ai"` in their Terraform configuration
- **THEN** the configuration is accepted
- **AND** the resource is created with AI access type
