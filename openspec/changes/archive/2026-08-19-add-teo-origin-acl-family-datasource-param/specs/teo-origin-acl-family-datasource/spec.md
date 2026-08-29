## ADDED Requirements

### Requirement: origin_acl_family field in tencentcloud_teo_origin_acl data source

The `tencentcloud_teo_origin_acl` data source SHALL expose `origin_acl_family` as a Computed string field inside the `origin_acl_info` block.

#### Scenario: Query origin ACL with OriginACLFamily returned
- **WHEN** user queries the `tencentcloud_teo_origin_acl` data source with a valid `zone_id`
- **AND** the `DescribeOriginACL` API returns `OriginACLInfo.OriginACLFamily` with a non-nil value (e.g., "mlc")
- **THEN** the data source SHALL set `origin_acl_info.0.origin_acl_family` to the returned value (e.g., "mlc")

#### Scenario: Query origin ACL with OriginACLFamily as nil
- **WHEN** user queries the `tencentcloud_teo_origin_acl` data source with a valid `zone_id`
- **AND** the `DescribeOriginACL` API returns `OriginACLInfo.OriginACLFamily` as nil
- **THEN** the data source SHALL NOT set `origin_acl_info.0.origin_acl_family` (field remains empty string)

#### Scenario: Schema validation
- **WHEN** user examines the data source schema
- **THEN** the `origin_acl_info` block SHALL contain a field named `origin_acl_family` of type `TypeString` with `Computed: true`

### Requirement: Data source Read method handles OriginACLFamily

The `dataSourceTencentCloudTeoOriginAclRead` function SHALL correctly read and set the `OriginACLFamily` field from the API response.

#### Scenario: Successful read with OriginACLFamily
- **WHEN** the `DescribeTeoOriginAclByFilter` service method returns a response with `OriginACLInfo.OriginACLFamily` set
- **THEN** the Read method SHALL set `origin_acl_family` in the `origin_acl_info` map to the response value
- **AND** the method SHALL call `d.Set("origin_acl_info", ...)` with the updated map

#### Scenario: Read with nil OriginACLFamily
- **WHEN** the `DescribeTeoOriginAclByFilter` service method returns a response with `OriginACLInfo.OriginACLFamily` as nil
- **THEN** the Read method SHALL NOT include `origin_acl_family` key in the `origin_acl_info` map
- **AND** the method SHALL still successfully set `origin_acl_info` with other fields

## Implementation Verification

### Requirement: Code implementation exists and is correct

The implementation SHALL exist in the data source file and follow Terraform Provider patterns.

#### Scenario: Schema definition exists
- **WHEN** examining `tencentcloud/services/teo/data_source_tc_teo_origin_acl.go`
- **THEN** the schema SHALL define `origin_acl_family` as `TypeString` with `Computed: true` inside the `origin_acl_info` block (around line 309-313)

#### Scenario: Read method implementation exists
- **WHEN** examining the `dataSourceTencentCloudTeoOriginAclRead` function
- **THEN** the function SHALL check `if respData.OriginACLInfo.OriginACLFamily != nil`
- **AND** SHALL set `originACLInfoMap["origin_acl_family"]` to `respData.OriginACLInfo.OriginACLFamily`

#### Scenario: Cloud API support exists
- **WHEN** examining `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`
- **THEN** the `OriginACLInfo` struct SHALL contain field `OriginACLFamily *string` with appropriate JSON tag
