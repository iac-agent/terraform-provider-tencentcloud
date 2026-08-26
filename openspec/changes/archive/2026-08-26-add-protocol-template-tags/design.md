## Context

The `tencentcloud_protocol_template` resource (`tencentcloud/services/vpc/resource_tc_protocol_template.go`) manages VPC protocol templates (service templates). It currently supports `name` and `protocols` fields. The `CreateServiceTemplate` API already accepts `Tags []*Tag` (a list of key-value pairs), and the `DescribeServiceTemplates` API returns `TagSet []*Tag` in the `ServiceTemplate` response. Adding tags support enables users to organize protocol templates using Tencent Cloud's tag system.

The `ModifyServiceTemplateAttribute` API does **not** support tags, so tags can only be set during creation and read back — they cannot be updated in-place. This is acceptable for Terraform as tags are typically part of the resource identity.

## Goals / Non-Goals

**Goals:**
- Add an optional `tags` parameter (TypeList) to the resource schema, with each element containing `key` (TypeString) and `value` (TypeString) sub-fields
- Pass tags from Terraform configuration to `CreateServiceTemplate` API on resource creation
- Read tags from `DescribeServiceTemplates` API response (`ServiceTemplate.TagSet`) and set them in Terraform state
- Mark tags as immutable on update (tags are ForceNew)

**Non-Goals:**
- In-place tag updates via `ModifyServiceTemplateAttribute` (not supported by the API)
- Import-time tag reconciliation (tags are read-only after creation)
- Tag support on the `protocol_template_group` resource (separate change)

## Decisions

1. **Tags as TypeList with key/value sub-fields**: The tags parameter is a `TypeList` where each element is a map with `key` and `value` sub-fields. This mirrors the `[]*Tag` structure in the SDK and provides explicit schema validation.

2. **Tags are ForceNew (immutable on update)**: Since `ModifyServiceTemplateAttribute` does not support tags, tags are marked as ForceNew. Any change to tags will trigger resource recreation. The `"tags"` field is added to `immutableArgs` in the update function.

3. **Service layer changes**: The `CreateServiceTemplate` function signature is extended to accept `tags map[string]string`, and the `DescribeServiceTemplateById` function returns tags alongside the template. This keeps the service layer as the single point of SDK interaction.

4. **Read path**: Tags are read from `response.ServiceTemplateSet[].TagSet` in the `DescribeServiceTemplates` response. The `DescribeServiceTemplateById` function already returns `*vpc.ServiceTemplate` which contains `TagSet`.

## Risks / Trade-offs

- **Risk**: Tags are immutable — changing tags forces recreation of the protocol template. This could disrupt existing security group rules that reference the template. → **Mitigation**: Document this behavior clearly. Users should plan tag names carefully before creating the resource.
- **Risk**: The `DescribeServiceTemplates` API might not return tags for templates created before tag support was added. → **Mitigation**: Handle nil/empty `TagSet` gracefully — set tags to empty list in state.