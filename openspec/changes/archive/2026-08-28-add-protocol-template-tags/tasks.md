## 1. Schema and Resource Definition

- [x] 1.1 Add `tags` parameter (TypeMap, Optional) to `tencentcloud_protocol_template` resource schema in `tencentcloud/services/vpc/resource_tc_protocol_template.go`

## 2. Service Layer - Create

- [x] 2.1 Modify `CreateServiceTemplate` function signature in `tencentcloud/services/vpc/service_tencentcloud_vpc.go` to accept a `tags map[string]interface{}` parameter
- [x] 2.2 In `CreateServiceTemplate`, convert the tags map to `[]*vpc.Tag` and set `request.Tags` before calling the API

## 3. Resource CRUD Functions

- [x] 3.1 In `resourceTencentCloudProtocolTemplateCreate`, extract tags from schema using `helper.GetTags(d, "tags")` and pass them to `vpcProtocol.CreateServiceTemplate`
- [x] 3.2 In `resourceTencentCloudProtocolTemplateRead`, read tags from `template.TagSet` and set them to the state via `d.Set("tags", tagMap)` (only when TagSet is not nil/empty)
- [x] 3.3 In `resourceTencentCloudProtocolTemplateUpdate`, add tag update handling using `svctag.TagService.ModifyTags()` when `d.HasChange("tags")` is true, with proper DiffTags computation

## 4. Unit Tests

- [x] 4.1 Add unit tests in `tencentcloud/services/vpc/resource_tc_protocol_template_test.go` for creating protocol template with tags
- [x] 4.2 Add unit tests for reading protocol template with tags (TagSet conversion to map)
- [x] 4.3 Add unit tests for updating tags on protocol template

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/vpc/resource_tc_protocol_template.md` to include `tags` parameter in example usage
