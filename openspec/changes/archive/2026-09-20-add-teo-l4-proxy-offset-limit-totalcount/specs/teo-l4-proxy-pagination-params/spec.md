## ADDED Requirements

### Requirement: offset parameter on tencentcloud_teo_l4_proxy resource
The `tencentcloud_teo_l4_proxy` resource SHALL expose an optional `offset` attribute of type `schema.TypeInt` with `Optional: true`. This attribute SHALL be passed to the `DescribeL4Proxy` API request as the `Offset` field. When not specified by the user, the attribute SHALL NOT be set in the API request, allowing the API to use its default value (0).

#### Scenario: offset is passed to DescribeL4Proxy API
- **WHEN** the read function is called on a `tencentcloud_teo_l4_proxy` resource and `offset` is set in the Terraform configuration
- **THEN** the `DescribeL4Proxy` API request SHALL include the `Offset` field with the configured value

#### Scenario: offset is not set when omitted
- **WHEN** the read function is called and `offset` is not set in the Terraform configuration
- **THEN** the `DescribeL4Proxy` API request SHALL NOT include the `Offset` field

### Requirement: limit parameter on tencentcloud_teo_l4_proxy resource
The `tencentcloud_teo_l4_proxy` resource SHALL expose an optional `limit` attribute of type `schema.TypeInt` with `Optional: true`. This attribute SHALL be passed to the `DescribeL4Proxy` API request as the `Limit` field. When not specified by the user, the attribute SHALL NOT be set in the API request, allowing the API to use its default value (20).

#### Scenario: limit is passed to DescribeL4Proxy API
- **WHEN** the read function is called on a `tencentcloud_teo_l4_proxy` resource and `limit` is set in the Terraform configuration
- **THEN** the `DescribeL4Proxy` API request SHALL include the `Limit` field with the configured value

#### Scenario: limit is not set when omitted
- **WHEN** the read function is called and `limit` is not set in the Terraform configuration
- **THEN** the `DescribeL4Proxy` API request SHALL NOT include the `Limit` field

### Requirement: total_count computed attribute on tencentcloud_teo_l4_proxy resource
The `tencentcloud_teo_l4_proxy` resource SHALL expose a computed attribute `total_count` of type `schema.TypeInt` with `Computed: true`. This attribute SHALL reflect the `TotalCount` value from the `DescribeL4Proxy` API response. The attribute SHALL NOT be user-configurable.

#### Scenario: total_count is set after resource read
- **WHEN** the read function is called for an existing `tencentcloud_teo_l4_proxy` resource and the `DescribeL4Proxy` API returns a valid response with `TotalCount`
- **THEN** the `total_count` attribute SHALL be set from `response.Response.TotalCount` in the Terraform state

#### Scenario: total_count is nil in API response
- **WHEN** the `DescribeL4Proxy` API response has a nil `TotalCount` field
- **THEN** the `total_count` attribute SHALL NOT be set in the Terraform state