## Context

The `tencentcloud_ses_template` resource is a RESOURCE_KIND_GENERAL resource with full CRUD operations:
- `CreateEmailTemplate` → creates a template (returns `TemplateID`)
- `UpdateEmailTemplate` → updates a template
- `GetEmailTemplate` → reads template details (returns `TemplateContent`, `TemplateName`, `TemplateStatus`)
- `DeleteEmailTemplate` → deletes a template

Currently, the resource reads `TemplateContent` and `TemplateName` from the `GetEmailTemplate` response but does not expose `TemplateStatus` (the template review status: 0=approved, 1=pending, 2=rejected). This field is already available in the vendor SDK (`GetEmailTemplateResponseParams.TemplateStatus *uint64`) and is returned by the API, but is not used by the Terraform resource.

## Goals / Non-Goals

**Goals:**
- Add a new computed output attribute `template_status` (TypeInt) to the resource schema
- Populate `template_status` in the Read function from `templateResponse.TemplateStatus`
- Update documentation to describe the new attribute

**Non-Goals:**
- No changes to Create/Update/Delete functions — `TemplateStatus` is read-only
- No changes to the service layer (`service_tencentcloud_ses.go`) — the `DescribeSesTemplate` function already returns the full `GetEmailTemplateResponseParams` struct which includes `TemplateStatus`
- No changes to the vendor SDK
- No breaking changes or schema modifications to existing fields

## Decisions

1. **Schema field type**: `schema.TypeInt` — matching the `*uint64` type in the SDK. Use `Computed: true` since this is a read-only output field.

2. **No ForceNew**: Since this is computed, it does not trigger re-creation.

3. **No `d.SetId("")` override**: The existing Read function already handles nil response correctly by checking `templateResponse == nil`. The new field will only be set if `templateResponse` is non-nil.

4. **Backward compatibility**: This is a purely additive change — existing configurations and state files remain valid. The new field will appear in state after the next Read operation.

## Risks / Trade-offs

- **[Low Risk] Type mismatch**: The SDK uses `*uint64` which maps safely to `schema.TypeInt`. No truncation or sign issues.
- **[Low Risk] API returns nil**: If `TemplateStatus` is nil in the API response, the `d.Set` call must guard against nil. We will add a nil check before calling `d.Set("template_status", ...)`.