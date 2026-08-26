## Why

The `tencentcloud_protocol_template` resource currently does not support tags, preventing users from managing tag metadata on protocol templates. The `CreateServiceTemplate` API already supports passing `Tags` (a list of `Key`/`Value` pairs), and the `DescribeServiceTemplates` API returns `TagSet` in the response. Adding tags support enables users to organize and categorize protocol templates using Tencent Cloud's tag management system.

## What Changes

- Add `tags` parameter (TypeList, Optional) to the `tencentcloud_protocol_template` resource schema, where each tag element contains `key` (TypeString) and `value` (TypeString) sub-fields
- Pass tags from Terraform configuration to `CreateServiceTemplate` API during resource creation
- Read tags from `DescribeServiceTemplates` API response (`ServiceTemplate.TagSet`) and set them in Terraform state during resource read
- Mark tags as immutable on update (the `ModifyServiceTemplateAttribute` API does not support tags)

## Capabilities

### New Capabilities
- `protocol-template-tags`: Support for tags on the `tencentcloud_protocol_template` resource, including creating protocol templates with tags and reading tags back into state

### Modified Capabilities
<!-- None - this is a new capability being added to an existing resource -->

## Impact

- **Code**: `tencentcloud/services/vpc/resource_tc_protocol_template.go` (schema + CRUD logic), `tencentcloud/services/vpc/service_tencentcloud_vpc.go` (service layer: `CreateServiceTemplate`, `DescribeServiceTemplateById`), `tencentcloud/services/vpc/resource_tc_protocol_template_test.go` (unit tests), `tencentcloud/services/vpc/resource_tc_protocol_template.md` (documentation)
- **SDK**: Uses existing `Tag` struct and `Tags`/`TagSet` fields already present in `vpc/v20170312` SDK
- **Backward Compatibility**: Fully backward compatible — `tags` is an optional field, existing configurations without tags continue to work