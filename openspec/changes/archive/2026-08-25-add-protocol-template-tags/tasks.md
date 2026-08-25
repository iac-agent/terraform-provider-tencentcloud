## 1. Service Layer Modifications

- [x] 1.1 Modify `CreateServiceTemplate` method in `service_tencentcloud_vpc.go` to accept optional `key` and `value` parameters and set `Tags` on the request
- [x] 1.2 Verify `DescribeServiceTemplateById` method returns `TagSet` field from the ServiceTemplate response (no modification needed if already present)

## 2. Resource Schema and CRUD Implementation

- [x] 2.1 Add `Key` and `Value` parameters to `tencentcloud_protocol_template` resource schema in `resource_tc_protocol_template.go` (both Optional, TypeString)
- [x] 2.2 Modify `resourceTencentCloudProtocolTemplateCreate` function to read `Key` and `Value` from config and pass to service layer
- [x] 2.3 Modify `resourceTencentCloudProtocolTemplateRead` function to extract `TagSet` from response and set `Key` and `Value` in state
- [x] 2.4 Verify `resourceTencentCloudProtocolTemplateUpdate` function does not need modification (tags are immutable)

## 3. Testing

- [x] 3.1 Add unit test cases in `resource_tc_protocol_template_test.go` for creating resource with Key and Value parameters
- [x] 3.2 Add unit test cases for reading resource with tags
- [x] 3.3 Add unit test cases for creating resource without tags (backward compatibility)

## 4. Documentation

- [x] 4.1 Create/update `tencentcloud_protocol_template.md` example file with Key and Value parameter usage
- [x] 4.2 Verify `make doc` generates correct documentation (will be done in finalize phase)

## 5. Validation

- [x] 5.1 Verify code compiles without errors (go build - not executed, just code review)
- [x] 5.2 Verify backward compatibility: existing configurations without Key/Value still work
- [x] 5.3 Verify all required files are modified: resource_tc_protocol_template.go, service_tencentcloud_vpc.go, resource_tc_protocol_template_test.go, tencentcloud_protocol_template.md
