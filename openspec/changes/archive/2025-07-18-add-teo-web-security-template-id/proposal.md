## Why

The `tencentcloud_teo_web_security_template` resource currently lacks the `template_id` parameter in its Terraform schema, even though `template_id` is already referenced in the update function's `immutableArgs` list. The `template_id` is a core identifier returned by the `CreateWebSecurityTemplate` API and is required by the `DescribeWebSecurityTemplate`, `ModifyWebSecurityTemplate`, and `DeleteWebSecurityTemplate` APIs. Exposing it as a schema field allows users to reference the template ID in other Terraform resources and data sources.

## What Changes

- Add `template_id` as a computed optional field in the `tencentcloud_teo_web_security_template` resource schema
- The `template_id` is generated server-side by `CreateWebSecurityTemplate` API and returned in the response
- In the read function, set `template_id` from the composite ID (zoneId#templateId)
- The `template_id` should be immutable (already enforced in update function's `immutableArgs`)

## Capabilities

### New Capabilities
- `teo-web-security-template-id`: Add template_id parameter to the teo_web_security_template resource, exposing the server-generated template identifier as a Terraform schema field

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_web_security_template.go` - Schema definition, read function
- `tencentcloud/services/teo/resource_tc_teo_web_security_template.md` - Documentation
- `tencentcloud/services/teo/resource_tc_teo_web_security_template_test.go` - Unit tests
