## ADDED Requirements

### Requirement: Resource Schema Definition
The resource `tencentcloud_teo_dns_record_14` SHALL define a schema with the following fields:

**Required fields:**
- `zone_id` (TypeString, ForceNew): TEO 站点 ID
- `name` (TypeString): DNS 记录名称
- `type` (TypeString, ForceNew): DNS 记录类型，取值: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV
- `content` (TypeString): DNS 记录内容

**Optional fields:**
- `location` (TypeString): 解析线路，默认为 Default
- `ttl` (TypeInt): 缓存时间，取值范围 60~86400 秒，默认 300
- `weight` (TypeInt): DNS 记录权重，取值范围 -1~100，默认 -1
- `priority` (TypeInt): MX 记录优先级，取值范围 0~50，默认 0

**Computed fields:**
- `record_id` (TypeString): DNS 记录 ID
- `status` (TypeString): DNS 记录解析状态
- `created_on` (TypeString): 创建时间
- `modified_on` (TypeString): 修改时间

The resource SHALL support import via `schema.ImportStatePassthrough`.

#### Scenario: Schema defines all required and optional fields
- **WHEN** the resource schema is initialized
- **THEN** it SHALL contain zone_id, name, type, content as required fields, location/ttl/weight/priority as optional fields, and record_id/status/created_on/modified_on as computed fields

#### Scenario: zone_id and type are ForceNew
- **WHEN** zone_id or type is changed in the Terraform configuration
- **THEN** Terraform SHALL destroy and recreate the resource

### Requirement: Resource Create
The resource SHALL call `CreateDnsRecord` API to create a DNS record. After successful creation, the resource SHALL set its ID to `zoneId#recordId` format using `tccommon.FILED_SP` as separator.

#### Scenario: Successful DNS record creation
- **WHEN** a valid DNS record configuration is applied
- **THEN** the system SHALL call `CreateDnsRecord` with zone_id, name, type, content, and optional location/ttl/weight/priority
- **AND** the response SHALL contain a non-empty RecordId
- **AND** the resource ID SHALL be set to `zoneId#recordId`
- **AND** a Read operation SHALL be performed to refresh the state

#### Scenario: Create API returns empty RecordId
- **WHEN** the `CreateDnsRecord` API response contains an empty RecordId
- **THEN** the system SHALL return a NonRetryableError

#### Scenario: Create API call fails
- **WHEN** the `CreateDnsRecord` API call fails
- **THEN** the error SHALL be wrapped with `tccommon.RetryError()` for retry handling

### Requirement: Resource Read
The resource SHALL call `DescribeDnsRecords` API to read the current state of a DNS record. It SHALL use the `Filters` parameter with `id` filter to locate the specific record by its record_id.

#### Scenario: Successful DNS record read
- **WHEN** the resource ID is parsed to extract zone_id and record_id
- **THEN** the system SHALL call `DescribeDnsRecords` with ZoneId and Filters containing the record_id
- **AND** it SHALL locate the matching record in the response DnsRecords list
- **AND** it SHALL set all schema fields from the response (only for non-nil response fields)

#### Scenario: DNS record not found
- **WHEN** the `DescribeDnsRecords` response does not contain the target record
- **THEN** the resource SHALL be marked as gone by setting `d.SetId("")`

#### Scenario: Read with pagination
- **WHEN** the DescribeDnsRecords response contains more records than the page limit
- **THEN** the system SHALL use the maximum Limit value (1000) to minimize pagination
- **AND** it SHALL search through the result list for the target record

### Requirement: Resource Update
The resource SHALL call `ModifyDnsRecords` API to update a DNS record. It SHALL construct a DnsRecord object containing the RecordId and all updatable fields (name, content, location, ttl, weight, priority).

#### Scenario: Successful DNS record update
- **WHEN** any updatable field (name, content, location, ttl, weight, priority) is changed
- **THEN** the system SHALL call `ModifyDnsRecords` with ZoneId and a DnsRecord list containing the RecordId and updated field values
- **AND** a Read operation SHALL be performed to refresh the state

#### Scenario: No fields changed
- **WHEN** no updatable fields have changed
- **THEN** the system SHALL NOT call the ModifyDnsRecords API

### Requirement: Resource Delete
The resource SHALL call `DeleteDnsRecords` API to delete a DNS record by passing the record_id in the RecordIds list.

#### Scenario: Successful DNS record deletion
- **WHEN** the resource is destroyed
- **THEN** the system SHALL parse the resource ID to extract zone_id and record_id
- **AND** the system SHALL call `DeleteDnsRecords` with ZoneId and RecordIds containing the record_id
- **AND** the function SHALL return nil on success

#### Scenario: Delete API call fails
- **WHEN** the `DeleteDnsRecords` API call fails
- **THEN** the error SHALL be wrapped with `tccommon.RetryError()` for retry handling

### Requirement: Resource Import
The resource SHALL support import via `schema.ImportStatePassthrough`. When importing, the ID SHALL be in the format `zoneId#recordId`.

#### Scenario: Import existing DNS record
- **WHEN** the user imports a DNS record with ID `zoneId#recordId`
- **THEN** the system SHALL parse the ID to extract zone_id and record_id
- **AND** a Read operation SHALL be performed to populate the state