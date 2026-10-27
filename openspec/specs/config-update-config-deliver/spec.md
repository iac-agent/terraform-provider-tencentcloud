# config-update-config-deliver Specification

## Purpose
Define the Terraform Plugin Framework resource `tencentcloud_config_update_config_deliver` that manages the singleton Config log-delivery setting of an account and region. The resource mirrors the `DescribeConfigDeliver` / `UpdateConfigDeliver` APIs one-to-one, disables delivery on destroy, and follows the repository baseline for client injection without the framework helper package.

## Requirements

### Requirement: Resource registration

The provider SHALL expose a Terraform Plugin Framework resource type named
`tencentcloud_config_update_config_deliver` implemented in
`tencentcloud/services/config/resource_tc_config_update_config_deliver_config.go`.

The resource MUST be registered by appending its factory `NewConfigUpdateConfigDeliverResource` to the
`resourceFactories` slice in `tencentcloud/framework/registry.go`, with the `config` service package
imported there. It MUST NOT be registered in `tencentcloud/provider.go`'s `ResourcesMap`, because
that map is the SDKv2 registry.

`tencentcloud/provider.md` MUST list `tencentcloud_config_update_config_deliver` under the Config product's
Resources section so that `make doc` generates a website page for it.

The pre-existing SDKv2 resource `tencentcloud_config_deliver_config` MUST remain unchanged and registered.

#### Scenario: Resource type resolves

- **WHEN** an operator runs `terraform plan` against a configuration containing a
  `resource "tencentcloud_config_update_config_deliver" "<name>"` block
- **THEN** Terraform resolves the type and produces a plan instead of an "unknown resource type" error

#### Scenario: Registered only on the framework side

- **WHEN** a developer inspects `tencentcloud/provider.go` and `tencentcloud/framework/registry.go`
- **THEN** `tencentcloud_config_update_config_deliver` appears exactly once, in
  `framework/registry.go`'s `resourceFactories`, and does not appear in `provider.go`

#### Scenario: Website documentation generated

- **WHEN** `make doc` runs
- **THEN** a website page for `config_update_config_deliver` is generated from the resource Schema and the
  resource markdown file

### Requirement: Schema mirrors the delivery-setting APIs one-to-one

The resource schema SHALL expose exactly the following attributes, with no additional attributes, no
renamed attributes and no synthetic flags:

| Attribute | Type | Mode | Constraints |
| --- | --- | --- | --- |
| `id` | string | Computed | `UseStateForUnknown` |
| `status` | int64 | Required | delivery switch; `0` = disabled, `1` = enabled |
| `deliver_name` | string | Optional + Computed | delivery service name |
| `target_arn` | string | Optional + Computed | six-segment resource ARN |
| `deliver_prefix` | string | Optional + Computed | log prefix |
| `deliver_type` | string | Optional + Computed | validator `OneOf("COS", "CLS")` |
| `deliver_content_type` | int64 | Optional + Computed | validator `Between(1, 3)` |
| `create_time` | string | Computed | read-only, never sent to the update API |

`status` MUST be `Required`: it decides whether audit logs are delivered, so it cannot be inherited
implicitly from cloud state. It MUST NOT be `Optional` or `Optional + Computed`.

`create_time` exists only in `DescribeConfigDeliver`'s response and MUST be `Computed` only; it MUST NOT
appear in the `UpdateConfigDeliver` request.

No attribute MUST be marked `RequiresReplace`: every field accepted by `UpdateConfigDeliver` is mutable in
place, so a change to any of them produces an in-place update rather than a replacement.

#### Scenario: Required switch is enforced at plan time

- **WHEN** a configuration declares the resource without a `status` argument
- **THEN** `terraform plan` fails with an error naming the missing required `status` argument

#### Scenario: Read-only creation time cannot be set

- **WHEN** a configuration sets `create_time` in a resource block
- **THEN** `terraform plan` fails with a "value for read-only attribute" error

#### Scenario: Enum values are validated before the API call

- **WHEN** a configuration sets `deliver_type = "OSS"` or `deliver_content_type = 5`
- **THEN** `terraform plan` fails with a validation error and no cloud API call is made

#### Scenario: In-place update instead of replacement

- **WHEN** an already-applied resource changes `target_arn`, `deliver_type` or `deliver_prefix`
- **THEN** the plan shows `~ update in-place` and never `-/+ must be replaced`

### Requirement: Singleton state ID and import

Because the delivery setting is a singleton per account and region, the resource state ID SHALL be the fixed
sentinel string `config_update_config_deliver`, not a cloud-returned identifier.

`Delete` MUST NOT clear the ID for the purpose of "no resource found": `DescribeConfigDeliver` has no
not-found state.

`ImportState` SHALL be hand-written (there is no helper package for it in this baseline) and MUST accept
only the sentinel ID: it stores the ID in state and lets the subsequent `Read` hydrate all attributes. An
import with any other ID MUST fail with a diagnostic that names the only valid import ID.

#### Scenario: Create stores the sentinel ID

- **WHEN** the resource is created successfully
- **THEN** the ID recorded in state is exactly `config_update_config_deliver`

#### Scenario: Import hydrates the whole setting

- **GIVEN** the account already has a delivery setting configured out-of-band
- **WHEN** the operator runs `terraform import tencentcloud_config_update_config_deliver.x config_update_config_deliver`
- **THEN** the import succeeds and a following `terraform plan` shows the current `status`, `deliver_name`,
  `target_arn`, `deliver_prefix`, `deliver_type`, `deliver_content_type` and `create_time` populated from the
  cloud response

#### Scenario: Unsupported import ID rejected

- **WHEN** the operator runs `terraform import ... some-other-id`
- **THEN** the import fails with a diagnostic stating that `config_update_config_deliver` is the only
  supported import ID

### Requirement: Create and Update share one modify path

`Create` and `Update` SHALL both call the same private helper which invokes `UpdateConfigDeliver`; there MUST
be no separate creation API call.

The modify path MUST set `Status` unconditionally (it is `Required`) and MUST set each other request field
only when its planned value is neither null nor unknown.

`Create` MUST adopt pre-existing cloud state rather than attempting to detect and reject it: applying to an
account whose delivery setting is already configured SHALL overwrite that setting.

Both operations MUST run within a retry loop bounded by a project retry timeout
(`tccommon.WriteRetryTimeout`), and MUST surface the API error verbatim when retries are exhausted, including
the resource ID in the message.

#### Scenario: First apply performs an update against an empty setting

- **GIVEN** the account has never configured delivery
- **WHEN** the resource is applied with `status = 1` and a `target_arn`
- **THEN** `UpdateConfigDeliver` is called once and the resource is created in state with the sentinel ID

#### Scenario: Omitted optional fields are not sent

- **WHEN** the resource is applied with only `status = 0` and no other arguments
- **THEN** the `UpdateConfigDeliver` request carries only `Status`

#### Scenario: Subsequent change updates in place

- **WHEN** an applied resource changes `deliver_prefix`
- **THEN** exactly one `UpdateConfigDeliver` call is made carrying `Status` and the new `DeliverPrefix`, and
  no replacement occurs

### Requirement: Read reflects the cloud response nil-safely

`Read` SHALL call `DescribeConfigDeliver` with no request parameters and write every attribute from the
response.

A response field that is a nil pointer MUST be written as a null Terraform value, not as the Go zero value
(`0` or `""`), so that `Optional + Computed` attributes do not produce spurious diffs.

If `response` or `response.Response` is nil, `Read` MUST fail with a diagnostic and MUST NOT remove the
resource from state, because an empty body is an API anomaly and not evidence that the setting was deleted.

`Read` MUST log the resource ID when the describe call fails, referring to the resource as
`config_update_config_deliver`.

#### Scenario: Nil response fields become null

- **GIVEN** `DescribeConfigDeliver` returns `DeliverPrefix = nil`
- **WHEN** `Read` runs
- **THEN** `deliver_prefix` is null in state rather than an empty string

#### Scenario: Empty response body is an error, not a deletion

- **WHEN** `DescribeConfigDeliver` returns a nil `Response` payload
- **THEN** `Read` reports an error diagnostic and the resource stays in state with its ID intact

#### Scenario: Drift is detected

- **GIVEN** an applied resource with `status = 1`
- **WHEN** the setting is disabled out-of-band and `terraform plan` runs
- **THEN** the plan shows `status` changing from `1` back to the configured value

### Requirement: Destroy turns delivery off

This resource is an enable switch. `Delete` SHALL call `UpdateConfigDeliver` with `Status = 0` so that
destroying the resource disables Config log delivery on the cloud side.

`Delete` MUST NOT delete, reset or otherwise mutate the Config service itself, and MUST NOT attempt to delete
any COS bucket or CLS topic referenced by `target_arn`.

`Delete` MUST send only `Status = 0`; it MUST NOT clear the remaining fields in the request.

The turn-off behaviour MUST be described in both the `Schema.Description` and the resource markdown file.

#### Scenario: Destroy disables delivery

- **GIVEN** an applied resource with `status = 1`
- **WHEN** `terraform destroy` runs
- **THEN** `UpdateConfigDeliver` is called with `Status = 0`, the resource is removed from state, and the COS
  bucket / CLS topic referenced by `target_arn` still exists

#### Scenario: Destroy is safe when already disabled

- **GIVEN** an applied resource whose cloud-side `status` is already `0`
- **WHEN** `terraform destroy` runs
- **THEN** the destroy succeeds (the disable call is idempotent) and does not fail the run

### Requirement: Client injection follows the baseline without the framework helper package

`tencentcloud/framework/fw` does not exist in this repository, so the resource MUST NOT import it and MUST NOT
embed `fw.ResourceWithConfigure`, `fw.WithImportByID` or `fw.WithNoUpdate`.

The resource SHALL hold a `ConfigService` client field and implement `Configure` by hand: return silently when
`req.ProviderData` is nil, type-assert it to `*sharedmeta.ProviderMeta`, report a diagnostic naming the
unexpected type when the assertion fails, and otherwise construct the service with
`NewConfigService(meta.Client)`.

The resource MUST keep compile-time assertions for `resource.Resource`, `resource.ResourceWithConfigure` and
`resource.ResourceWithImportState`.

#### Scenario: Compiled without the helper package

- **WHEN** the provider package is compiled
- **THEN** the resource file contains no import of `tencentcloud/framework/fw` and the build succeeds

#### Scenario: Misconfigured provider data is reported

- **WHEN** `req.ProviderData` has a type other than `*sharedmeta.ProviderMeta`
- **THEN** `Configure` adds an "Unexpected ProviderData type" diagnostic instead of panicking
