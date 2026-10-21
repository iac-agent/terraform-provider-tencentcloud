## ADDED Requirements

### Requirement: Resource tencentcloud_teo_zone_v2 manages EdgeOne site lifecycle

The provider SHALL expose a `tencentcloud_teo_zone_v2` resource (RESOURCE_KIND_GENERAL) that manages the full lifecycle of an EdgeOne (TEO) site (Zone) via the teo/v20220901 cloud APIs `CreateZone`, `DescribeZones`, `ModifyZone` and `DeleteZone`. The resource id SHALL be the single cloud-returned `ZoneId`. The resource SHALL support `terraform import` with the ZoneId via `schema.ImportStatePassthrough`.

#### Scenario: Resource is registered in provider
- **WHEN** the provider binary is built and `terraform providers schema` is inspected
- **THEN** resource `tencentcloud_teo_zone_v2` appears in the provider schema registry alongside existing teo resources, and `tencentcloud/provider.md` lists it under the TencentCloud EdgeOne(TEO) Resource section

#### Scenario: Import an existing site
- **WHEN** `terraform import tencentcloud_teo_zone_v2.example zone-2noz78a8ev6k` is executed
- **THEN** the provider reads the site via `DescribeZones` filtered by `zone-id` and populates the resource state with the returned Zone fields

### Requirement: Create creates a site with full CreateZone parameters

The Create operation SHALL call `CreateZone` with optional parameters `type`, `zone_name`, `area`, `plan_id`, `alias_zone_name`, `allow_duplicates`, `jump_start` and the `tags` list (each item with required `tag_key` and `tag_value`). The API call SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with cloud API errors wrapped by `tccommon.RetryError`. After the retry succeeds, the implementation SHALL verify that the response is non-nil and `ZoneId` is non-empty (returning an error otherwise), then set the terraform id to the returned `ZoneId` and invoke the Read operation to populate computed fields including the create-time `ownership_verification` details.

#### Scenario: Create a site with plan binding and tags
- **WHEN** a configuration sets `zone_name`, `type = "full"`, `area = "global"`, `plan_id`, `allow_duplicates = true`, `jump_start = true`, and one `tags { tag_key, tag_value }` block
- **THEN** `CreateZone` is called with all those values mapped to the corresponding request fields (`Type`, `ZoneName`, `Area`, `PlanId`, `AliasZoneName`, `Tags[].TagKey/TagValue`, `AllowDuplicates`, `JumpStart`), the resource id becomes the returned `ZoneId`, and the state contains the returned ownership verification info (dns/file/ns verification details)

#### Scenario: Create returns empty ZoneId
- **WHEN** `CreateZone` succeeds but the response `ZoneId` is nil or empty
- **THEN** the Create operation returns an error and no empty id is written to tfstate

### Requirement: Read fetches the site by zone-id filter with max pagination limit

The Read operation SHALL call `DescribeZones` with `Filters = [{Name: "zone-id", Values: [<id>]}]` and `Limit = 100` (the maximum value documented by the cloud API), wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. If the response is nil, its `Response` is nil, or the `Zones` list is empty, the implementation SHALL first log `log.Printf("[CRUD] teo_zone_v2 id=%s", d.Id())` and then clear the id via `d.SetId("")` without returning an error. Otherwise it SHALL populate the state from the first Zone entry: flat fields (`zone_id`, `zone_name`, `alias_zone_name`, `area`, `type`, `status`, `cname_status`, `active_status`, `lock_status`, `paused`, `created_on`, `modified_on`, `name_servers`, `original_name_servers`), and computed nested lists (`tags` with `tag_key`/`tag_value`, `resources` with billing fields including nested `sv`, `ns_detail`, `cname_detail`, `dns_pod_detail`, `ownership_verification` with `dns_verification`/`file_verification`/`ns_verification` sub-structures, `vanity_name_servers`, `vanity_name_servers_ips`, `work_mode_infos`). Every `d.Set` SHALL be guarded by a nil check on the corresponding response field.

#### Scenario: Read populates full site state
- **WHEN** `DescribeZones` returns one Zone of access type `full` with resources, NSDetail, tags and work mode infos
- **THEN** all mapped fields are set into state, nested lists are flattened into terraform list-of-object structures, and nil response sub-structures are skipped without error

#### Scenario: Read finds the site deleted
- **WHEN** `DescribeZones` returns an empty `Zones` list
- **THEN** the log line `[CRUD] teo_zone_v2 id=<id>` is emitted before the id is cleared, and Read returns success with an empty id so terraform plans re-creation

### Requirement: Update modifies mutable zone attributes via ModifyZone

The Update operation SHALL detect changes on the mutable arguments `type`, `zone_name`, `area`, `alias_zone_name` and `vanity_name_servers` (nested object with required `switch` and optional `servers` list). When any of them changed, it SHALL call `ModifyZone` with `ZoneId` taken from `d.Id()` together with the changed fields, wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and error-wrapped by `tccommon.RetryError`. After a successful update it SHALL invoke the Read operation to refresh state.

#### Scenario: Update area and vanity name servers
- **WHEN** `area` changes to `mainland` and a `vanity_name_servers { switch = "on", servers = [...] }` block is added or changed
- **THEN** `ModifyZone` is called with `ZoneId` from the resource id, the new `Area` and a `VanityNameServers` structure containing `Switch` and `Servers`, followed by a Read to refresh the state

#### Scenario: Update with no mutable change
- **WHEN** no mutable argument has changed
- **THEN** `ModifyZone` is not called and the Read operation still refreshes the state

### Requirement: Delete removes the site via DeleteZone

The Delete operation SHALL call `DeleteZone` with `ZoneId` taken from `d.Id()`, wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with errors wrapped by `tccommon.RetryError` and a nil-response check. Since `DeleteZone` is a synchronous API returning no task id, no polling is required and the operation returns success after the API call succeeds.

#### Scenario: Delete a site
- **WHEN** `terraform destroy` deletes a `tencentcloud_teo_zone_v2` resource
- **THEN** `DeleteZone` is called with the resource id as `ZoneId` and the destroy completes without polling

### Requirement: Unit tests mock cloud APIs with gomonkey

The implementation SHALL provide unit tests in `resource_tc_teo_zone_v2_test.go` that use gomonkey to patch `UseTeoV20220901Client` and the `CreateZoneWithContext`/`DescribeZonesWithContext`/`ModifyZoneWithContext`/`DeleteZoneWithContext` client methods (not the terraform test suite). Tests SHALL cover: create with full parameter passing assertions and ownership verification state assertions; create error on empty ZoneId response; read normal; read not-found clearing the id; read empty response clearing the id; update with parameter assertions; delete with parameter assertions. The generated test code SHALL be compilable in the current environment.

#### Scenario: Mocked create test passes
- **WHEN** the unit test suite for `resource_tc_teo_zone_v2` is compiled and run
- **THEN** the gomonkey-mocked tests verify request parameter mapping and state population without contacting real cloud APIs
