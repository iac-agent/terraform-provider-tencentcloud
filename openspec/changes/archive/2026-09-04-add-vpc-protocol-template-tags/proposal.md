## Why

The VPC `CreateServiceTemplate` API accepts a `Tags` parameter to bind tags to protocol port templates at creation time, and the `DescribeServiceTemplates` API returns the bound tags via `TagSet`. However, the Terraform resource `tencentcloud_protocol_template` does not expose tags, so users cannot manage tags on protocol port templates through Terraform — they must use the console or API directly.

## What Changes

- Add a `tags` (TypeMap, Optional) parameter to the `tencentcloud_protocol_template` resource so users can specify tags when creating a protocol port template. The tags map keys are tag keys and map values are tag values, consistent with the sibling `tencentcloud_address_extra_template` resource.
- In the Create flow, pass the user-provided tags to the `CreateServiceTemplate` API request as `request.Tags` (a list of `*vpc.Tag` with `Key` and `Value`).
- In the Read flow, refresh `tags` from the `DescribeServiceTemplates` API response (`ServiceTemplateSet[].TagSet[].Key` and `TagSet[].Value`).
- In the Update flow, when `tags` changes, use the shared tag service (`svctag.NewTagService`, `svctag.DiffTags`, `tagService.ModifyTags`) with `tccommon.BuildTagResourceName("vpc", "service", region, d.Id())` to reconcile tags, because the `ModifyServiceTemplateAttribute` API does not accept a `Tags` field.

## Capabilities

### New Capabilities
- `vpc-protocol-template-tags`: Enable the optional `tags` parameter on the `tencentcloud_protocol_template` resource, allowing users to bind tags to a protocol port template at creation and reconcile tag changes on update via the shared tag service.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/vpc/resource_tc_protocol_template.go` — add `tags` schema field, wire tags through Create (pass to `CreateServiceTemplate` request), add Read support (set tags from `TagSet`), add Update support (reconcile tags via shared tag service)
  - `tencentcloud/services/vpc/service_tencentcloud_vpc.go` — update `CreateServiceTemplate` service function to accept and pass a `tags` parameter to the `CreateServiceTemplate` API request
  - `tencentcloud/services/vpc/resource_tc_protocol_template_test.go` — add test cases covering tags in create and update flows
  - `tencentcloud/services/vpc/resource_tc_protocol_template.md` — update documentation example with `tags` usage
- **API behavior:** The `CreateServiceTemplate` API supports a `Tags []*Tag` input parameter; the `DescribeServiceTemplates` API returns `ServiceTemplateSet[].TagSet []*Tag`. The `ModifyServiceTemplateAttribute` API does NOT accept tags, so tag updates after creation are handled through the shared tag service (`tag.ModifyResourcesTag` API).
- **Backward compatibility:** fully backward compatible — the new `tags` parameter is Optional and defaults to not being set, so existing configurations continue to work unchanged.
