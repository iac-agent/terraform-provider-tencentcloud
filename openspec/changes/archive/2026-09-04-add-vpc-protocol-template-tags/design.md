## Context

The `tencentcloud_protocol_template` resource currently supports creating protocol port templates with only `name` and `protocols` parameters. The VPC `CreateServiceTemplate` API also accepts a `Tags []*Tag` parameter to bind tags at creation time, and the `DescribeServiceTemplates` API returns the bound tags via `ServiceTemplateSet[].TagSet`. The Terraform resource does not expose tags, so users cannot manage them through Terraform.

**Current state:**
- Resource file: `tencentcloud/services/vpc/resource_tc_protocol_template.go`
- Service layer: `tencentcloud/services/vpc/service_tencentcloud_vpc.go` (`CreateServiceTemplate` function)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312`
- Sibling reference: `tencentcloud/services/vpc/resource_tc_address_extra_template.go` (already implements the same tags pattern for address templates)

**API behavior analysis:**

| API | Tags in Request | Tags in Response |
|-----|-----------------|------------------|
| `CreateServiceTemplate` | Yes (`Tags []*Tag`, each with `Key` and `Value`) | N/A |
| `DescribeServiceTemplates` | N/A | Yes (`ServiceTemplateSet[].TagSet []*Tag`) |
| `ModifyServiceTemplateAttribute` | No (`Tags` field does NOT exist) | N/A |
| `DeleteServiceTemplate` | No | N/A |

**Key constraint:** Tags can be set at creation via `CreateServiceTemplate`, but the `ModifyServiceTemplateAttribute` API does not accept a `Tags` field. Therefore, tag updates after creation must be handled through the shared tag service (`tag.ModifyResourcesTag` API), consistent with how `tencentcloud_address_extra_template` handles tag updates.

## Goals / Non-Goals

**Goals:**
- Add a `tags` (TypeMap, Optional) parameter to `tencentcloud_protocol_template`, where map keys are tag keys and map values are tag values
- Pass tags to the `CreateServiceTemplate` API request as `request.Tags` (list of `*vpc.Tag` with `Key`/`Value`) when specified by the user
- Read tags from `DescribeServiceTemplates` API response (`TagSet`) to support state refresh and import
- Reconcile tag changes on Update using the shared tag service (`svctag.NewTagService`, `svctag.DiffTags`, `tagService.ModifyTags`) with `tccommon.BuildTagResourceName("vpc", "service", region, d.Id())`, because `ModifyServiceTemplateAttribute` does not support tags
- Update the `CreateServiceTemplate` service function signature to accept and pass a tags parameter
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Adding tags to the `tencentcloud_protocol_templates` datasource or `tencentcloud_protocol_template_group` resource (out of scope)
- Changing any existing schema field behavior

## Decisions

### Decision 1: Use `tags` as TypeMap (consistent with sibling resources)

**Rationale:** The sibling `tencentcloud_address_extra_template` resource models tags as `TypeMap` where keys are tag keys and values are tag values. Using the same model keeps the provider consistent and matches the established pattern for VPC template resources. The cloud API `Tag` struct has `Key` and `Value` string fields, which maps naturally to a map.

### Decision 2: Tag updates use the shared tag service, not `ModifyServiceTemplateAttribute`

**Rationale:** The `ModifyServiceTemplateAttribute` API does NOT accept a `Tags` field (verified in the SDK `models.go`). Therefore, when tags change after creation, the provider reconciles them through the shared tag service: `svctag.NewTagService`, `svctag.DiffTags` to compute the diff, and `tagService.ModifyTags` to apply. The resource name is built with `tccommon.BuildTagResourceName("vpc", "service", region, d.Id())`. This is the exact same approach used by `tencentcloud_address_extra_template` (which uses resource type `"address"`).

### Decision 3: Resource type for tags is `"service"`

**Rationale:** The protocol template maps to the VPC "service template" concept (API uses `ServiceTemplate`). Following the convention where the address template uses resource type `"address"` in `BuildTagResourceName`, the service/protocol template uses resource type `"service"`. The service type is `"vpc"`, consistent with all other VPC resources.

### Decision 4: Update `CreateServiceTemplate` service function signature

**Rationale:** The existing `CreateServiceTemplate(ctx, name, services)` service function does not accept tags. To wire tags through Create, the signature is extended to accept a `tags` parameter (e.g., `tags map[string]interface{}` or `tags []*vpc.Tag`). Passing the raw schema map and converting inside the service function keeps the resource Create handler simple. This mirrors how the resource Create builds the request directly in `tencentcloud_address_extra_template`, but here we extend the existing service-layer wrapper to minimize churn.

## Risks / Trade-offs

- **[Risk] Incorrect resource type string in `BuildTagResourceName`**: If the resource type string (`"service"`) does not match what the tag API expects, tag reconciliation will fail at runtime.
  - **Mitigation:** The string follows the established VPC convention (address template uses `"address"`). The tag service returns an error if the resource cannot be found, so a mismatch surfaces as a clear API error rather than silent data loss.

- **[Risk] Imported resources do not refresh tags correctly**: The Read function must parse `TagSet` from the Describe response.
  - **Mitigation:** Read guards each `Tag.Key`/`Tag.Value` against nil before populating the map, consistent with `tencentcloud_address_extra_template`.

- **[Trade-off] Tag updates require an extra API call**: Because `ModifyServiceTemplateAttribute` does not support tags, a separate `tag.ModifyResourcesTag` call is needed on Update when tags change. This is unavoidable given the API constraints and is consistent with sibling resources.
