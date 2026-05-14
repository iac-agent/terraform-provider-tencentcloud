## Requirements

### Requirement: Resource Schema Definition
The resource `tencentcloud_teo_dns_record_18` SHALL define the following schema fields:

- `zone_id` (TypeString, Required, ForceNew): 站点 ID，创建后不可变更
- `name` (TypeString, Required): DNS 记录名称
- `type` (TypeString, Required): DNS 记录类型（A, AAAA, MX, CNAME, TXT, NS, CAA, SRV）
- `content` (TypeString, Required): DNS 记录内容
- `location` (TypeString, Optional, Computed): 解析路线，默认为 DEFAULT
- `ttl` (TypeInt, Optional, Computed): 缓存时间，范围 60-86400，默认 300 秒
- `weight` (TypeInt, Optional, Computed): DNS 记录权重，范围 -1 到 100，默认 -1
- `priority` (TypeInt, Optional, Computed): MX 优先级，范围 0-50，默认 0
- `record_id` (TypeString, Computed): DNS 记录 ID，由创建接口返回
- `status` (TypeString, Computed): DNS 记录状态（enable/disable）
- `created_on` (TypeString, Computed): 创建时间
- `modified_on` (TypeString, Computed): 修改时间

The resource SHALL support import via `schema.ImportStatePassthrough`.

#### Scenario: Schema fields are correctly defined
- **WHEN** the resource schema is registered
- **THEN** all required fields (zone_id, name, type, content) are present with Required=true
- **AND** zone_id has ForceNew=true
- **AND** optional fields (location, ttl, weight, priority) have Computed=true
- **AND** computed fields (record_id, status, created_on, modified_on) have Computed=true and Optional=false

### Requirement: Create DNS Record
The resource SHALL create a DNS record by calling the `CreateDnsRecord` API with zone_id, name, type, content, location, ttl, weight, and priority parameters. Upon success, the resource SHALL set its ID to `zone_id + FILED_SP + record_id` using the returned RecordId.

#### Scenario: Successful creation with all fields
- **WHEN** a DNS record is created with zone_id, name="www", type="A", content="1.2.3.4", location="DEFAULT", ttl=300, weight=-1, priority=0
- **THEN** the CreateDnsRecord API is called with all specified parameters
- **AND** the resource ID is set to "zone_id#record_id"
- **AND** the Read function is called to refresh the state

#### Scenario: Successful creation with required fields only
- **WHEN** a DNS record is created with only zone_id, name, type, and content
- **THEN** the CreateDnsRecord API is called without optional parameters
- **AND** the resource ID is set to "zone_id#record_id"

#### Scenario: Create API returns empty response
- **WHEN** the CreateDnsRecord API returns a nil response or nil RecordId
- **THEN** a NonRetryableError SHALL be returned

### Requirement: Read DNS Record
The resource SHALL read a DNS record by calling the `DescribeTeoDnsRecordById` service method with zone_id and record_id parsed from the composite ID. If the record is not found, the resource SHALL set its ID to empty string to signal deletion.

#### Scenario: Successful read
- **WHEN** the resource reads its state with a valid composite ID "zone_id#record_id"
- **THEN** the DescribeTeoDnsRecordById method is called with zone_id and record_id
- **AND** all non-nil response fields are set on the resource state

#### Scenario: Record not found
- **WHEN** the DescribeTeoDnsRecordById returns nil
- **THEN** the resource ID is set to empty string
- **AND** no error is returned

#### Scenario: Broken composite ID
- **WHEN** the resource ID does not contain exactly 2 parts separated by FILED_SP
- **THEN** an error SHALL be returned

### Requirement: Update DNS Record
The resource SHALL update a DNS record by calling the `ModifyDnsRecords` API when any mutable field (name, type, content, location, ttl, weight, priority) has changed. The request SHALL include a DnsRecord struct with the RecordId and all mutable fields.

#### Scenario: Successful update of mutable fields
- **WHEN** the content field is changed from "1.2.3.4" to "5.6.7.8"
- **THEN** the ModifyDnsRecords API is called with ZoneId and a DnsRecord containing RecordId and all mutable fields
- **AND** the Read function is called to refresh the state

#### Scenario: No mutable fields changed
- **WHEN** no mutable fields have changed
- **THEN** the ModifyDnsRecords API is NOT called
- **AND** the Read function is still called to refresh the state

### Requirement: Delete DNS Record
The resource SHALL delete a DNS record by calling the `DeleteDnsRecords` API with ZoneId and RecordIds (containing the single record_id).

#### Scenario: Successful deletion
- **WHEN** a DNS record is deleted with composite ID "zone_id#record_id"
- **THEN** the DeleteDnsRecords API is called with ZoneId=zone_id and RecordIds=["record_id"]

### Requirement: Resource Registration
The resource SHALL be registered in `tencentcloud/provider.go` with the name `tencentcloud_teo_dns_record_18` and the factory function `ResourceTencentCloudTeoDnsRecord18`. The resource SHALL also be listed in `tencentcloud/provider.md`.

#### Scenario: Resource is registered in provider
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_dns_record_18` is available for use in Terraform configurations

### Requirement: Unit Tests
The resource SHALL have unit tests using gomonkey mock approach. Tests SHALL cover Create, Read, Update, and Delete operations with mocked cloud API responses.

#### Scenario: Unit tests pass with mocked API
- **WHEN** unit tests are run with `go test -gcflags=all=-l`
- **THEN** all tests for the resource pass successfully

### Requirement: Resource Documentation
The resource SHALL have a markdown documentation file at `gendoc/resource/tencentcloud_teo_dns_record_18.md` following the standard format: one-line description, example usage with jsonencode() for JSON fields, and import section (since this is a RESOURCE_KIND_GENERAL resource).

#### Scenario: Documentation file exists and is valid
- **WHEN** the documentation is generated
- **THEN** the file contains a description mentioning TEO, example usage, and import instructions specifying the composite ID format
