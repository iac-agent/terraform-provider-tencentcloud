## Context

The `tencentcloud_protocol_template` resource manages VPC protocol port templates via the TencentCloud VPC API. Currently it supports `name` and `protocols` parameters, but does not support tags. The cloud API `CreateServiceTemplate` already supports a `Tags` field (type `[]*Tag`, where `Tag` has `Key` and `Value` string fields), and the `DescribeServiceTemplates` response includes `TagSet` (type `[]*Tag`) in each `ServiceTemplate` object.

Other VPC resources (e.g., `tencentcloud_nat_gateway`, `tencentcloud_eni`, `tencentcloud_route_table`) already support tags using `schema.TypeMap` with `helper.GetTags()` / tag service for CRUD. However, for `protocol_template`, the `CreateServiceTemplate` API natively accepts `Tags` in the request, and the `DescribeServiceTemplates` response returns `TagSet` directly on the resource object. The `ModifyServiceTemplateAttribute` API does NOT support updating tags.

## Goals / Non-Goals

**Goals:**
- Add `tags` parameter (optional, `schema.TypeMap`) to `tencentcloud_protocol_template` resource schema
- Pass tags to `CreateServiceTemplate` API during resource creation
- Read tags from `TagSet` in `DescribeServiceTemplates` response during resource read
- Support tag updates via the tag service (`svctag.TagService`) since `ModifyServiceTemplateAttribute` does not support tags

**Non-Goals:**
- Changing existing `name` or `protocols` parameter behavior
- Adding any other new parameters beyond `tags`
- Modifying the `protocol_template_group` resource

## Decisions

### 1. Tags schema type: `schema.TypeMap`

**Decision**: Use `schema.TypeMap` for the `tags` parameter, consistent with other VPC resources (`tencentcloud_nat_gateway`, `tencentcloud_eni`, `tencentcloud_route_table`).

**Rationale**: All VPC resources in this provider use `schema.TypeMap` for tags. This provides a simple key-value interface. Although the cloud API uses `[]*Tag` (Key/Value pairs), `schema.TypeMap` is the established convention.

**Alternative**: `schema.TypeList` with nested `Key`/`Value` fields — rejected because it's inconsistent with the existing VPC resource patterns.

### 2. Tag creation: Native API `Tags` field

**Decision**: Pass tags directly via `request.Tags` in `CreateServiceTemplate`.

**Rationale**: The `CreateServiceTemplate` API natively supports `Tags` in its request body. Using the native API field is simpler and more reliable than using the tag service after creation.

### 3. Tag reading: From `TagSet` in DescribeServiceTemplates response

**Decision**: Read tags from `template.TagSet` returned by `DescribeServiceTemplates`, converting `[]*Tag` to `map[string]interface{}` for `d.Set("tags", ...)`.

**Rationale**: The `TagSet` field is directly available on the `ServiceTemplate` response object, so we can read tags without an extra API call to the tag service.

### 4. Tag update: Via tag service

**Decision**: Use `svctag.TagService.ModifyTags()` for tag updates, since `ModifyServiceTemplateAttribute` does not support tags.

**Rationale**: The tag service provides a standard mechanism for tag CRUD on TencentCloud resources. The resource name format follows `tccommon.BuildTagResourceName("vpc", "service-template", region, templateId)`. This is consistent with how other VPC resources handle tag updates.

**Alternative**: Making tags `ForceNew` — rejected because it would require resource recreation for tag changes, which is unnecessarily destructive.

## Risks / Trade-offs

- [Tag service resource name format uncertain] → Need to verify the correct resource prefix and type for the tag service. Will use `tccommon.BuildTagResourceName("vpc", "service-template", region, templateId)` based on the VPC product convention, but this may need adjustment during implementation.
- [Tags not visible in ModifyServiceTemplateAttribute] → If users expect tags to be part of the modify flow, they need to understand that tag updates go through a separate API. This is consistent with other VPC resources.
