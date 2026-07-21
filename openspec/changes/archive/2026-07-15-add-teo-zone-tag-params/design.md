## Context

The `tencentcloud_teo_zone` resource manages TEO (EdgeOne) zones in Terraform. Currently, when performing tag operations (reading tags via `DescribeResourceTagsByResourceIds` and modifying tags via `ModifyResourceTags`), the resource hardcodes:
- `serviceType = "teo"` (used in `DescribeResourceTags` call and QCS resource name)
- `resourceRegion = tcClient.Region` (used in `DescribeResourceTags` call and QCS resource name)

TEO is a global service, and the tag API's `DescribeResourceTagsByResourceIds` request supports `ResourceRegion` and `ServiceType` as input parameters. Some scenarios require users to specify custom values for these parameters (e.g., when the resource is in a specific region different from the provider's default region, or when the service type identifier differs).

The SDK already supports these parameters:
- `DescribeResourceTagsByResourceIdsRequest` has `ResourceRegion *string` and `ServiceType *string` fields
- `ModifyResourceTagsRequest` uses a `Resource` string in QCS six-segment format which encodes service type and region

## Goals / Non-Goals

**Goals:**
- Add `ResourceRegion` and `ServiceType` as optional parameters to the `tencentcloud_teo_zone` resource schema
- Use these parameters when calling tag-related APIs (DescribeResourceTagsByResourceIds, ModifyResourceTags) in the resource's CRUD methods
- Maintain backward compatibility: if not specified, fall back to current hardcoded values

**Non-Goals:**
- Changing the tag service's `DescribeResourceTags` or `ModifyTags` method signatures (we will handle the parameters at the resource level)
- Adding `ResourceRegion`/`ServiceType` to other TEO resources (only `tencentcloud_teo_zone` is in scope)
- Modifying any TEO-specific API calls (CreateZone, DescribeZones, ModifyZone, etc.) — these parameters only affect tag operations

## Decisions

### Decision 1: Schema design for new parameters
Both `ResourceRegion` and `ServiceType` will be added as `Optional` + `Computed` string fields in the resource schema. When not explicitly set by the user, the Read operation will populate them with the default values currently used in the code (provider region for `ResourceRegion`, `"teo"` for `ServiceType`).

**Rationale**: Using `Optional` + `Computed` ensures backward compatibility — existing configurations that don't specify these parameters will continue to work with the same defaults, while the Read operation will reflect the effective values in the state.

**Alternative considered**: Making them `Optional` only (without `Computed`). This was rejected because the Read operation needs to reflect what values were actually used, and `Computed` allows Terraform to show the effective defaults in the state.

### Decision 2: How to pass ResourceRegion and ServiceType to tag operations
For `DescribeResourceTags`: Currently the `DescribeResourceTags` method in the tag service takes `serviceType`, `resourceType`, `region`, `resourceId` as parameters. We will read `ResourceRegion` and `ServiceType` from the schema and pass them to the existing `DescribeResourceTags` method.

For `ModifyResourceTags` (used in Create and Update): Currently the resource constructs a QCS resource name string. We will read `ResourceRegion` and `ServiceType` from the schema and use them in the `BuildTagResourceName` call instead of hardcoding the values.

**Rationale**: This approach minimizes changes to the tag service layer and keeps the logic contained within the teo zone resource. The tag service's method signatures remain unchanged.

### Decision 3: Default values
- `ResourceRegion` defaults to `tcClient.Region` (same as current behavior)
- `ServiceType` defaults to `"teo"` (same as current behavior)

**Rationale**: These are the values currently hardcoded in the resource, so using them as defaults ensures zero breaking changes.

## Risks / Trade-offs

- [Risk] Users who already have `tencentcloud_teo_zone` resources in their state will see a plan diff if they run `terraform plan` after upgrading, because the new `Computed` fields will be populated in the state on the next Read. → Mitigation: Using `Computed: true` means Terraform will auto-populate these without requiring user input, and no plan diff should occur since the Read will set the same default values that were previously hardcoded.

- [Risk] If a user specifies an incorrect `ResourceRegion` or `ServiceType`, tag operations may fail or return incorrect results. → Mitigation: Document the expected values clearly in the schema descriptions and the .md documentation file.
