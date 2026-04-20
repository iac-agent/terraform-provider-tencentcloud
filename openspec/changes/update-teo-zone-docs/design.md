## Context

This is a documentation-only update for the `tencentcloud_teo_zone` Terraform resource. The `type` parameter currently documents four access types (`partial`, `full`, `noDomainAccess`), but two additional access types (`dnsPodAccess` and `ai`) are valid but undocumented. The implementation already supports these values in the underlying API, but users cannot discover them through the Terraform schema documentation.

## Goals / Non-Goals

**Goals:**
- Update the schema description field to document all six valid access types
- Maintain backward compatibility with existing configurations
- Auto-generate updated website documentation

**Non-Goals:**
- Modify the parameter validation logic (already supports all types)
- Change API calls or business logic
- Add new parameters or change schema structure

## Decisions

**Decision 1: Direct schema description update**
- **Rationale**: The simplest approach is to update the `Description` field in the schema. This is a non-breaking change that improves documentation without affecting functionality.
- **Alternatives considered**:
  - Adding validation: Not needed since API already accepts all types
  - Creating a new parameter: Unnecessary, the current parameter already supports all types

**Decision 2: Documentation generation via `make doc`**
- **Rationale**: The project uses `make doc` to auto-generate website documentation from schema definitions. Running this command will update the documentation files automatically.
- **Alternatives considered**:
  - Manually editing website docs: Against project convention, would be overwritten

## Risks / Trade-offs

**Risk: Incomplete API documentation**
- **Mitigation**: Before finalizing, verify the complete list of valid types from Tencent Cloud API documentation

**Risk: User confusion about new types**
- **Mitigation**: Clear description format with all types listed and their purposes
