## 1. Service Layer Changes

- [x] 1.1 Modify `CreateServiceTemplate` in `service_tencentcloud_vpc.go` to accept a `tags map[string]string` parameter and populate `request.Tags` with `[]*vpc.Tag{Key, Value}` from the map
- [x] 1.2 Modify `DescribeServiceTemplateById` in `service_tencentcloud_vpc.go` to return a `tags map[string]string` from `template.TagSet`, handling nil/empty `TagSet` gracefully

## 2. Resource Schema and CRUD

- [x] 2.1 Add `tags` parameter (TypeList, Optional, ForceNew) with `key` (TypeString, Required) and `value` (TypeString, Optional) sub-fields to the resource schema in `resource_tc_protocol_template.go`
- [x] 2.2 Update `resourceTencentCloudProtocolTemplateCreate` to extract tags from schema, convert to `map[string]string`, and pass to `CreateServiceTemplate`
- [x] 2.3 Update `resourceTencentCloudProtocolTemplateRead` to set tags from `DescribeServiceTemplateById` return value into Terraform state, handling nil/empty tags
- [x] 2.4 Update `resourceTencentCloudProtocolTemplateUpdate` to add `"tags"` to `immutableArgs` array, returning error if tags change is detected

## 3. Unit Tests

- [x] 3.1 Add unit test cases in `resource_tc_protocol_template_test.go` for tags: create with tags, create without tags, read with tags, read without tags, update tags (expect error)

## 4. Documentation

- [x] 4.1 Update `resource_tc_protocol_template.md` to include tags in the example usage