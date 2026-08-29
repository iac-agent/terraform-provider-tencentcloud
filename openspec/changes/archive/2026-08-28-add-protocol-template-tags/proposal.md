## Why

The `tencentcloud_protocol_template` resource currently does not support specifying tags during creation. The underlying cloud API `CreateServiceTemplate` already supports a `Tags` parameter, and the `DescribeServiceTemplates` response returns `TagSet` for each template. Adding tags support enables users to categorize and manage protocol templates with metadata labels, which is consistent with other VPC resources (e.g., `tencentcloud_nat_gateway`, `tencentcloud_eni`, `tencentcloud_route_table`) that already support tags.

## What Changes

- Add a `tags` parameter (optional, TypeList of Key/Value pairs) to the `tencentcloud_protocol_template` resource schema
- Pass tags to the `CreateServiceTemplate` API request during resource creation
- Read tags from the `DescribeServiceTemplates` API response (`TagSet`) during resource read
- Since the `ModifyServiceTemplateAttribute` API does not support tags, tags will be immutable after creation (ForceNew is not needed; tags can be managed via the provider-level default tags mechanism or by recreating the resource)

## Capabilities

### New Capabilities
- `protocol-template-tags`: Add tags (Key/Value) parameter support to the `tencentcloud_protocol_template` resource, enabling tag-based categorization of protocol templates at creation time and reading tags from the Describe API response.

### Modified Capabilities

## Impact

- `tencentcloud/services/vpc/resource_tc_protocol_template.go`: Schema definition, create and read functions
- `tencentcloud/services/vpc/service_tencentcloud_vpc.go`: `CreateServiceTemplate` function signature and implementation to accept tags
- `tencentcloud/services/vpc/resource_tc_protocol_template_test.go`: Unit tests for tags parameter
- `tencentcloud/services/vpc/resource_tc_protocol_template.md`: Documentation example update
- Cloud API dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312` - `Tag` struct already available in vendor
