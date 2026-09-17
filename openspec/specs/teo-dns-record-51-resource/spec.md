## ADDED Requirements

### Requirement: Resource schema defines the teo dns record fields
The `tencentcloud_teo_dns_record_51` resource SHALL expose the following top-level schema fields (flattened, no extra wrapping list level):
- `zone_id` (Required, ForceNew, TypeString): 站点 ID。
- `name` (Required, TypeString): DNS 记录名，中文/韩文/日文域名需转换为 punycode 后输入。
- `type` (Required, TypeString): DNS 记录类型，取值 A / AAAA / MX / CNAME / TXT / NS / CAA / SRV。
- `content` (Required, TypeString): DNS 记录内容，根据 Type 填入对应内容。
- `location` (Optional, Computed, TypeString): 解析线路，不指定默认为 Default。
- `ttl` (Optional, Computed, TypeInt): 缓存时间，范围 60~86400 秒，默认 300。
- `weight` (Optional, Computed, TypeInt): 记录权重，范围 -1~100，默认 -1 表示不设置权重。
- `priority` (Optional, Computed, TypeInt): MX 记录优先级，范围 0~50，默认 0。
- `record_id` (Computed, TypeString): DNS 记录 ID。
- `status` (Computed, TypeString): 解析状态，enable（已生效）/ disable（已停用）。
- `created_on` (Computed, TypeString): 创建时间。
- `modified_on` (Computed, TypeString): 修改时间。

#### Scenario: Schema includes all CRUD fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include `zone_id`, `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority`, `record_id`, `status`, `created_on`, and `modified_on` with the types and optionality listed above

#### Scenario: Zone id change forces recreation
- **WHEN** a user changes `zone_id` in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated (ForceNew)

#### Scenario: Optional fields with cloud defaults
- **WHEN** a user omits `location`, `ttl`, `weight`, or `priority` in the configuration
- **THEN** the schema SHALL accept it and the Read operation SHALL populate the values from the cloud API response (Computed behavior) without a plan diff

### Requirement: Resource create operation
The Create operation SHALL build a `CreateDnsRecord` request from the schema (ZoneId, Name, Type, Content, Location, TTL, Weight, Priority), call `CreateDnsRecordWithContext` wrapped in `resource.Retry(tccommon.WriteRetryTimeout)` with `tccommon.RetryError` error wrapping, verify the response is non-empty (response, Response, RecordId non-nil and non-empty string, otherwise return `NonRetryableError` without writing an ID), then set the composite ID `zoneId#recordId` (separator `tccommon.FILED_SP`) outside the retry block, and finally invoke the Read operation.

#### Scenario: Successful create
- **WHEN** a user applies a `tencentcloud_teo_dns_record_51` resource with required fields
- **THEN** `CreateDnsRecord` is called with ZoneId/Name/Type/Content (and Location/TTL/Weight/Priority when configured), the resource ID is set to `zoneId#recordId`, and state is refreshed via Read

#### Scenario: Create returns empty record id
- **WHEN** the `CreateDnsRecord` API returns a nil Response or an empty RecordId
- **THEN** the Create operation SHALL return a `NonRetryableError` and SHALL NOT write an ID into state

#### Scenario: Create API failure retries
- **WHEN** the `CreateDnsRecord` API returns a retryable error
- **THEN** the operation SHALL retry within `tccommon.WriteRetryTimeout` and wrap errors with `tccommon.RetryError`

### Requirement: Resource read operation
The Read operation SHALL parse the composite ID `zoneId#recordId` from `d.Id()` (error when the split does not yield exactly 2 parts), then query the record through the teo service layer using `DescribeDnsRecords` with `Filters = [{Name: "id", Values: [recordId]}]` and `Limit` set to the cloud-API-documented maximum (1000), wrapped in `resource.Retry(tccommon.ReadRetryTimeout)`. When the record is found, each schema field SHALL be set only if the corresponding response field is non-nil. When the record is not found (empty result), the Read operation SHALL first log `log.Printf("[CRUD] teo_dns_record_51 id=%s", d.Id())` to preserve context and then call `d.SetId("")` and return nil.

#### Scenario: Read existing record
- **WHEN** Terraform refreshes state for an existing `tencentcloud_teo_dns_record_51` resource
- **THEN** `DescribeDnsRecords` is called with the `id` filter for the parsed recordId and Limit=1000, and all fields including `record_id`, `status`, `created_on`, `modified_on` are populated from the matched `DnsRecord` element

#### Scenario: Read record deleted out of band
- **WHEN** the DNS record no longer exists on the cloud side
- **THEN** the Read operation SHALL log `[CRUD] teo_dns_record_51 id=<id>` before clearing, then set `d.SetId("")` and return nil so Terraform plans a recreate

#### Scenario: Broken composite id
- **WHEN** `d.Id()` does not contain exactly two parts separated by `tccommon.FILED_SP`
- **THEN** the Read operation SHALL return an `id is broken` error

### Requirement: Resource update operation
The Update operation SHALL parse the composite ID, and when any of `name`, `type`, `content`, `location`, `ttl`, `weight`, `priority` has a change, build a single-element `DnsRecords` array for the `ModifyDnsRecords` request containing RecordId plus only the mutable fields (ZoneId, Status, CreatedOn, ModifiedOn SHALL NOT be set on the DnsRecord element because the cloud API treats them as output-only and ignores them as inputs). The request's top-level `ZoneId` SHALL be set from the parsed ID. The `ModifyDnsRecordsWithContext` call SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout)` with `tccommon.RetryError` wrapping, and the Update operation SHALL end by invoking Read.

#### Scenario: Update mutable fields
- **WHEN** a user changes `content` (or any of name/type/location/ttl/weight/priority) in the configuration
- **THEN** `ModifyDnsRecords` is called once with a single `DnsRecords` element containing RecordId and all mutable fields, and the top-level ZoneId

#### Scenario: No mutable field changed
- **WHEN** none of the mutable fields change and no ForceNew field changes
- **THEN** the Update operation SHALL NOT call `ModifyDnsRecords` and SHALL only refresh state via Read

#### Scenario: Update API failure
- **WHEN** `ModifyDnsRecords` returns a retryable error
- **THEN** the operation SHALL retry within `tccommon.WriteRetryTimeout` and wrap errors with `tccommon.RetryError`

### Requirement: Resource delete operation
The Delete operation SHALL parse the composite ID and call `DeleteDnsRecords` with `ZoneId` and `RecordIds = [recordId]`, wrapped in `resource.Retry(tccommon.WriteRetryTimeout)` with `tccommon.RetryError` wrapping, and SHALL verify the response is non-nil (returning `NonRetryableError` on a nil Response).

#### Scenario: Successful delete
- **WHEN** a user destroys the `tencentcloud_teo_dns_record_51` resource
- **THEN** `DeleteDnsRecords` is called with the parsed ZoneId and a single-element RecordIds array

#### Scenario: Delete API failure
- **WHEN** `DeleteDnsRecords` returns a retryable error
- **THEN** the operation SHALL retry within `tccommon.WriteRetryTimeout` and wrap errors with `tccommon.RetryError`

### Requirement: Resource import
The `tencentcloud_teo_dns_record_51` resource SHALL support `terraform import` via `schema.ImportStatePassthrough` using the composite ID `zoneId#recordId`.

#### Scenario: Import existing record
- **WHEN** a user runs `terraform import tencentcloud_teo_dns_record_51.example zone-2noz78a8ev6k#record-xxxx`
- **THEN** the resource state SHALL be populated by the Read operation using the parsed zoneId and recordId

### Requirement: Provider registration and documentation
The provider SHALL register `tencentcloud_teo_dns_record_51` in the `ResourcesMap` of `tencentcloud/provider.go` (bound to `teo.ResourceTencentCloudTeoDnsRecord51()`), list the resource name in `tencentcloud/provider.md` under the TEO category, and ship a resource example document `tencentcloud/services/teo/resource_tc_teo_dns_record_51.md` containing a one-sentence description (mentioning TEO, "Provides a resource to ..."), an Example Usage HCL block, and an Import section explaining the composite `zoneId#recordId` ID. The document SHALL NOT contain `Argument Reference` or `Attribute Reference` sections (auto-generated).

#### Scenario: Resource is usable after registration
- **WHEN** a user writes `resource "tencentcloud_teo_dns_record_51" "example" { ... }`
- **THEN** the provider SHALL resolve the resource from ResourcesMap and execute its CRUD lifecycle

#### Scenario: Documentation example covers import
- **WHEN** the resource markdown is rendered into website docs
- **THEN** the Example Usage shows all required and optional fields, and the Import section documents the `terraform import ... <zoneId>#<recordId>` syntax

### Requirement: Unit tests with gomonkey mocks
The change SHALL add `tencentcloud/services/teo/resource_tc_teo_dns_record_51_test.go` using gomonkey mocks (not the Terraform acceptance test suite) to unit-test the business logic: Create success (composite ID set), Create empty RecordId error, Read hit and miss paths, Update triggering ModifyDnsRecords with correct request fields, Update without changes skipping the API call, Delete success, and composite ID parse errors. The generated test code SHALL compile under the current environment.

#### Scenario: Mocked create sets composite id
- **WHEN** the unit test mocks `CreateDnsRecordWithContext` returning a valid RecordId and runs the Create function
- **THEN** the resource ID equals `zoneId#recordId` and state fields are populated via the mocked Read

#### Scenario: Mocked update calls modify once
- **WHEN** the unit test changes `content` and runs the Update function with mocked `ModifyDnsRecordsWithContext`
- **THEN** the modify request contains one DnsRecords element with RecordId and the new content, and its ZoneId-only output fields are not set
