## Context

The `tencentcloud_teo_web_security_template` resource manages TEO (TencentCloud EdgeOne) web security policy templates. Currently, the resource schema includes `zone_id`, `template_name`, and `security_policy`, but the `template_id` — the server-generated unique identifier for the template — is not exposed as a schema field. The `template_id` is only available internally as part of the composite resource ID (format: `zoneId#templateId`).

The `template_id` is already referenced in the update function's `immutableArgs` list, confirming it is expected to be an immutable field. The `CreateWebSecurityTemplate` API returns `TemplateId` in its response, and `DescribeWebSecurityTemplate`, `ModifyWebSecurityTemplate`, and `DeleteWebSecurityTemplate` all accept `TemplateId` as a required input parameter.

## Goals / Non-Goals

**Goals:**
- Add `template_id` as an optional computed field in the resource schema
- Ensure `template_id` is populated from the `CreateWebSecurityTemplate` API response after creation
- Ensure `template_id` is set in the read function from the composite ID
- Maintain backward compatibility with existing Terraform configurations

**Non-Goals:**
- Changing the composite ID format or structure
- Modifying any existing schema fields
- Adding any other new parameters beyond `template_id`

## Decisions

1. **Schema type: Optional + Computed**
   - `template_id` will be defined as `Optional: true, Computed: true` because it is generated server-side by the `CreateWebSecurityTemplate` API. Users do not need to provide it during creation, but it will be computed and stored in state after creation.

2. **No ForceNew on template_id**
   - Since `template_id` is server-generated, it does not need `ForceNew: true`. The value is set after creation from the API response.

3. **Immutable enforcement**
   - The existing `immutableArgs` list in the update function already includes `template_id`, so immutability is already enforced at the application level.

4. **Read function: set template_id from composite ID**
   - In the read function, after splitting the composite ID, `template_id` will be set via `d.Set("template_id", templateId)`.

5. **Create function: set template_id from API response**
   - After the `CreateWebSecurityTemplate` API call succeeds, the response contains `TemplateId`. This value will be stored and set in the state.

## Risks / Trade-offs

- [Risk] Adding a computed field could cause state drift if users manually set `template_id` → Mitigation: The field is Optional+Computed, so Terraform handles this gracefully
- [Risk] Backward compatibility must be maintained → Mitigation: Adding an Optional+Computed field is a non-breaking change; existing configurations will continue to work without specifying `template_id`
