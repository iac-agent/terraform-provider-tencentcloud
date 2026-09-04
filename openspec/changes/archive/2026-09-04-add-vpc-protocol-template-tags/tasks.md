## 1. Service Layer Changes

- [x] 1.1 Update `CreateServiceTemplate` function signature in `tencentcloud/services/vpc/service_tencentcloud_vpc.go` to accept a `tags` parameter (e.g. `tags []*vpc.Tag` or `tags map[string]interface{}`)
- [x] 1.2 In `CreateServiceTemplate`, build `request.Tags` as a list of `*vpc.Tag` (each with `Key` and `Value`) from the `tags` parameter when it is non-empty, and leave `request.Tags` unset when empty
- [x] 1.3 Verify the `CreateServiceTemplateRequest` SDK struct exposes the `Tags []*Tag` field (confirmed in `vendor/.../vpc/v20170312/models.go`)

## 2. Resource Schema Changes

- [x] 2.1 Add a `tags` schema field (`TypeMap`, `Optional`, `Description: "Tags of the protocol template."`) to `ResourceTencentCloudProtocolTemplate()` in `tencentcloud/services/vpc/resource_tc_protocol_template.go`

## 3. Create Function Changes

- [x] 3.1 Read `tags` from the schema data (`d.GetOk("tags")`) in `resourceTencentCloudProtocolTemplateCreate`
- [x] 3.2 Convert the tags map to `[]*vpc.Tag` and pass it to the updated `CreateServiceTemplate` service function
- [x] 3.3 Ensure existing create flow (id extraction, `d.SetId`) remains unchanged

## 4. Read Function Changes

- [x] 4.1 In `resourceTencentCloudProtocolTemplateRead`, after obtaining the template from `DescribeServiceTemplateById`, read `template.TagSet`
- [x] 4.2 Build a `map[string]string` from `TagSet` entries, guarding each `tag.Key` and `tag.Value` against nil before populating the map, and set it via `d.Set("tags", tags)` when non-empty

## 5. Update Function Changes

- [x] 5.1 Add imports for `svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"` and the helper/common packages as needed in `resource_tc_protocol_template.go`
- [x] 5.2 When `d.HasChange("tags")`, instantiate `svctag.NewTagService(client)`, compute `replaceTags, deleteTags` via `svctag.DiffTags(oldValue, newValue)`, build the resource name with `tccommon.BuildTagResourceName("vpc", "service", client.Region, d.Id())`, and call `tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)`
- [x] 5.3 Keep the existing `name`/`protocols` change handling that calls `ModifyServiceTemplateAttribute` unchanged

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/vpc/resource_tc_protocol_template.md` with a `tags` usage example in the Example Usage HCL block

## 7. Tests

- [x] 7.1 Add test cases in `tencentcloud/services/vpc/resource_tc_protocol_template_test.go` (using the existing terraform acceptance test suite) covering tags in the create flow and tags update flow

## 8. Validation

- [x] 8.1 Verify the code compiles successfully
- [x] 8.2 Verify no lint errors
