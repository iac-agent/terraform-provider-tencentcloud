## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_teo_dns_record_v2` with the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID
- `name` (Required, TypeString): DNS 记录名
- `type` (Required, TypeString): DNS 记录类型
- `content` (Required, TypeString): DNS 记录内容
- `location` (Optional, Computed, TypeString): DNS 记录解析线路
- `ttl` (Optional, Computed, TypeInt): 缓存时间（秒）
- `weight` (Optional, Computed, TypeInt): DNS 记录权重
- `priority` (Optional, Computed, TypeInt): MX 记录优先级
- `record_id` (Computed, TypeString): DNS 记录 ID

#### Scenario: Schema defines all fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include zone_id, name, type, content, location, ttl, weight, priority, and record_id with correct types and constraints

#### Scenario: zone_id is ForceNew
- **WHEN** zone_id is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

#### Scenario: Immutable fields are not ForceNew
- **WHEN** name, type, content, location, ttl, weight, or priority is changed in the Terraform configuration
- **THEN** the system SHALL invoke the Update function (which SHALL reject the change) rather than silently recreating

### Requirement: Resource Create Operation
The system SHALL implement Create by calling `CreateDnsRecord` with ZoneId, Name, Type, Content, Location, TTL, Weight, and Priority, then set the resource ID to `{zoneId}#{recordId}` from the returned RecordId.

#### Scenario: Successful create
- **WHEN** valid required fields are provided
- **THEN** the system SHALL call CreateDnsRecord
- **AND** set the resource ID to `{zoneId}#{recordId}`

#### Scenario: Create returns empty response
- **WHEN** CreateDnsRecord returns a nil response or nil RecordId
- **THEN** the system SHALL return a NonRetryableError instead of writing an empty ID

### Requirement: Resource Read Operation
The system SHALL implement Read by parsing `{zoneId}#{recordId}`, querying `DescribeDnsRecords` with a filter on `id`, and populating schema fields from the returned record.

#### Scenario: Successful read
- **WHEN** the record exists
- **THEN** the system SHALL populate zone_id, name, type, content, location, ttl, weight, priority, and record_id

#### Scenario: Record not found
- **WHEN** DescribeDnsRecords returns no matching record
- **THEN** the system SHALL log the id and clear the resource ID (`d.SetId("")`)

#### Scenario: Nil fields skipped
- **WHEN** a returned field is nil
- **THEN** the system SHALL NOT call d.Set for that field

### Requirement: Resource Update Operation (immutable)
The system SHALL implement Update as immutable: it SHALL reject changes to name, type, content, location, ttl, weight, and priority by returning an error.

#### Scenario: Immutable field changed
- **WHEN** any of name, type, content, location, ttl, weight, or priority is changed
- **THEN** the system SHALL return an error indicating the argument cannot be changed

#### Scenario: No immutable field changed
- **WHEN** no immutable field changed
- **THEN** the system SHALL call Read to refresh state

### Requirement: Resource Delete Operation
The system SHALL implement Delete by calling `DeleteDnsRecords` with ZoneId and RecordIds=[recordId].

#### Scenario: Successful delete
- **WHEN** the resource is destroyed
- **THEN** the system SHALL call DeleteDnsRecords with the parsed zoneId and recordId

#### Scenario: Delete API retry
- **WHEN** DeleteDnsRecords returns a retryable error
- **THEN** the system SHALL retry using tccommon.WriteRetryTimeout and return tccommon.RetryError on failure

### Requirement: Unit Tests
The system SHALL provide unit tests in `resource_tc_teo_dns_record_v2_test.go` using gomonkey to mock cloud API calls, covering Create, Read, Delete, and Update immutability.

#### Scenario: Create/Read/Delete tests pass
- **WHEN** go test is run with -gcflags=all=-l on the test file
- **THEN** all test cases for Create, Read, and Delete SHALL pass

#### Scenario: Update immutability test
- **WHEN** a test changes an immutable field
- **THEN** the Update SHALL return an error without calling any modify API

### Requirement: Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_dns_record_v2.md` with description, example usage, and import section.

#### Scenario: Documentation exists
- **WHEN** the resource is created
- **THEN** a .md file SHALL exist with a one-line description mentioning TEO, example usage, and import section showing the composite ID `{zoneId}#{recordId}`
