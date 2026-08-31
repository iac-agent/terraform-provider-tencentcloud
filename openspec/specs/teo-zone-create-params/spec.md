## Requirements

### Requirement: allow_duplicates schema parameter
The `tencentcloud_teo_zone` resource SHALL include an `allow_duplicates` schema field of type `schema.TypeBool` with `Optional: true` and `ForceNew: true`. This field controls whether duplicate zone access is allowed during zone creation. The default value SHALL be `false` when not specified.

#### Scenario: allow_duplicates is passed to CreateZone API
- **WHEN** a `tencentcloud_teo_zone` resource is created with `allow_duplicates` set to `true`
- **THEN** the `CreateZone` API request's `AllowDuplicates` field SHALL be set to `true`

#### Scenario: allow_duplicates is not specified
- **WHEN** a `tencentcloud_teo_zone` resource is created without specifying `allow_duplicates`
- **THEN** the `CreateZone` API request's `AllowDuplicates` field SHALL NOT be set (nil)

#### Scenario: allow_duplicates change forces resource recreation
- **WHEN** the `allow_duplicates` field is changed on an existing `tencentcloud_teo_zone` resource
- **THEN** Terraform SHALL destroy and recreate the resource

### Requirement: jump_start schema parameter
The `tencentcloud_teo_zone` resource SHALL include a `jump_start` schema field of type `schema.TypeBool` with `Optional: true` and `ForceNew: true`. This field controls whether to skip existing DNS record scanning during zone creation. The default value SHALL be `false` when not specified.

#### Scenario: jump_start is passed to CreateZone API
- **WHEN** a `tencentcloud_teo_zone` resource is created with `jump_start` set to `true`
- **THEN** the `CreateZone` API request's `JumpStart` field SHALL be set to `true`

#### Scenario: jump_start is not specified
- **WHEN** a `tencentcloud_teo_zone` resource is created without specifying `jump_start`
- **THEN** the `CreateZone` API request's `JumpStart` field SHALL NOT be set (nil)

#### Scenario: jump_start change forces resource recreation
- **WHEN** the `jump_start` field is changed on an existing `tencentcloud_teo_zone` resource
- **THEN** Terraform SHALL destroy and recreate the resource

### Requirement: file_verification sub-block in ownership_verification
The `ownership_verification` schema block of `tencentcloud_teo_zone` SHALL include a `file_verification` sub-block of type `schema.TypeList` with `Computed: true` and `MaxItems: 1`. The sub-block SHALL contain:
- `path` (TypeString, Computed): The URL path for file verification
- `content` (TypeString, Computed): The content to write to the verification file

#### Scenario: file_verification is populated from CreateZone response
- **WHEN** a `tencentcloud_teo_zone` resource is created and the `CreateZone` API response includes `OwnershipVerification.FileVerification`
- **THEN** the `file_verification` sub-block SHALL be populated with `path` and `content` from the response

#### Scenario: file_verification is populated from DescribeZones response
- **WHEN** a `tencentcloud_teo_zone` resource is read and the `DescribeZones` API response includes `Zone.OwnershipVerification.FileVerification`
- **THEN** the `file_verification` sub-block SHALL be populated with `path` and `content` from the response

#### Scenario: file_verification is nil in API response
- **WHEN** the API response's `OwnershipVerification.FileVerification` is nil
- **THEN** the `file_verification` sub-block SHALL be set to an empty list `[]interface{}{}`

### Requirement: ns_verification sub-block in ownership_verification
The `ownership_verification` schema block of `tencentcloud_teo_zone` SHALL include a `ns_verification` sub-block of type `schema.TypeList` with `Computed: true` and `MaxItems: 1`. The sub-block SHALL contain:
- `name_servers` (TypeList of TypeString, Computed): The list of DNS server addresses for NS verification

#### Scenario: ns_verification is populated from CreateZone response
- **WHEN** a `tencentcloud_teo_zone` resource is created and the `CreateZone` API response includes `OwnershipVerification.NsVerification`
- **THEN** the `ns_verification` sub-block SHALL be populated with `name_servers` from the response

#### Scenario: ns_verification is populated from DescribeZones response
- **WHEN** a `tencentcloud_teo_zone` resource is read and the `DescribeZones` API response includes `Zone.OwnershipVerification.NsVerification`
- **THEN** the `ns_verification` sub-block SHALL be populated with `name_servers` from the response

#### Scenario: ns_verification is nil in API response
- **WHEN** the API response's `OwnershipVerification.NsVerification` is nil
- **THEN** the `ns_verification` sub-block SHALL be set to an empty list `[]interface{}{}`