## ADDED Requirements

### Requirement: Resource Schema Definition
The `tencentcloud_teo_dns_record_22` resource SHALL define the following schema fields:
- `zone_id` (TypeString, Required, ForceNew): 站点 ID
- `name` (TypeString, Required): DNS 记录名称
- `type` (TypeString, Required): DNS 记录类型
- `content` (TypeString, Required): DNS 记录内容
- `location` (TypeString, Optional, Computed): DNS 记录解析线路
- `ttl` (TypeInt, Optional, Computed): 缓存时间
- `weight` (TypeInt, Optional, Computed): DNS 记录权重
- `priority` (TypeInt, Optional, Computed): MX 记录优先级
- `status` (TypeString, Optional, Computed): DNS 记录解析状态
- `created_on` (TypeString, Computed): 创建时间
- `modified_on` (TypeString, Computed): 修改时间
- `record_id` (TypeString, Computed): DNS 记录 ID

#### Scenario: Schema fields match cloud API parameters
- **WHEN** the resource schema is defined
- **THEN** all CreateDnsRecord input parameters (ZoneId, Name, Type, Content, Location, TTL, Weight, Priority) SHALL have corresponding schema fields
- **AND** all DnsRecord output-only fields (Status, CreatedOn, ModifiedOn, RecordId) SHALL be Computed fields

#### Scenario: zone_id is ForceNew
- **WHEN** zone_id is changed after resource creation
- **THEN** the resource SHALL be recreated (old resource deleted, new resource created)

### Requirement: Resource Create Operation
The resource SHALL call `CreateDnsRecord` API to create a DNS record, mapping schema fields to request parameters.

#### Scenario: Successful DNS record creation
- **WHEN** creating a teo_dns_record_22 resource with zone_id, name, type, and content
- **THEN** the system SHALL call CreateDnsRecord with the provided parameters
- **AND** set the resource ID to `zone_id#record_id` format using tccommon.FILED_SP separator
- **AND** call Read to refresh the resource state after creation

#### Scenario: Create with optional parameters
- **WHEN** creating a resource with optional parameters (location, ttl, weight, priority)
- **THEN** the system SHALL include these parameters in the CreateDnsRecord request

#### Scenario: Create API call failure
- **WHEN** the CreateDnsRecord API call fails
- **THEN** the system SHALL retry with tccommon.WriteRetryTimeout and wrap the error using tccommon.RetryError

#### Scenario: Create returns empty RecordId
- **WHEN** the CreateDnsRecord response has an empty RecordId
- **THEN** the system SHALL return a NonRetryableError

### Requirement: Resource Read Operation
The resource SHALL call `DescribeDnsRecords` API to read the DNS record state, using RecordId filter to locate the specific record.

#### Scenario: Successful DNS record read
- **WHEN** reading a teo_dns_record_22 resource with valid composite ID (zone_id#record_id)
- **THEN** the system SHALL parse the ID to extract zone_id and record_id
- **AND** call DescribeDnsRecords with ZoneId and filter by RecordId
- **AND** set all schema fields from the response if the corresponding response field is not nil

#### Scenario: Resource not found
- **WHEN** the DescribeDnsRecords response does not contain the target record
- **THEN** the system SHALL set the resource ID to empty string to mark it as deleted

#### Scenario: Read with nil response fields
- **WHEN** a DnsRecord response field is nil
- **THEN** the system SHALL NOT call d.Set() for that field

### Requirement: Resource Update Operation
The resource SHALL support updating mutable fields by calling `ModifyDnsRecords` for content changes and `ModifyDnsRecordsStatus` for status changes.

#### Scenario: Update mutable fields
- **WHEN** any of name, type, content, location, ttl, weight, or priority is changed
- **THEN** the system SHALL call ModifyDnsRecords with ZoneId and a DnsRecord containing the RecordId and updated fields

#### Scenario: Update status field
- **WHEN** the status field is changed
- **THEN** the system SHALL call ModifyDnsRecordsStatus with ZoneId and either RecordsToEnable or RecordsToDisable based on the new status value

#### Scenario: Update API call failure
- **WHEN** the ModifyDnsRecords or ModifyDnsRecordsStatus API call fails
- **THEN** the system SHALL retry with tccommon.WriteRetryTimeout and wrap the error using tccommon.RetryError

### Requirement: Resource Delete Operation
The resource SHALL call `DeleteDnsRecords` API to delete the DNS record.

#### Scenario: Successful DNS record deletion
- **WHEN** deleting a teo_dns_record_22 resource
- **THEN** the system SHALL parse the composite ID to extract zone_id and record_id
- **AND** call DeleteDnsRecords with ZoneId and RecordIds containing the record_id

#### Scenario: Delete API call failure
- **WHEN** the DeleteDnsRecords API call fails
- **THEN** the system SHALL retry with tccommon.WriteRetryTimeout and wrap the error using tccommon.RetryError

### Requirement: Resource Import Support
The resource SHALL support import with composite ID format `zone_id#record_id`.

#### Scenario: Import with valid composite ID
- **WHEN** importing a teo_dns_record_22 resource with ID `zone_id#record_id`
- **THEN** the system SHALL parse the ID and call Read to populate the resource state

### Requirement: Resource Registration
The `tencentcloud_teo_dns_record_22` resource SHALL be registered in `provider.go` and `provider.md`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_dns_record_22` SHALL be available and mapped to `ResourceTencentCloudTeoDnsRecord22` function

### Requirement: Unit Tests
The resource SHALL have unit tests using gomonkey mock approach, covering Create, Read, Update, and Delete operations.

#### Scenario: Unit tests cover CRUD operations
- **WHEN** running `go test -gcflags=all=-l` on the test file
- **THEN** all Create, Read, Update, and Delete mock tests SHALL pass
