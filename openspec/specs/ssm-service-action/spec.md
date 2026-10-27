# ssm-service-action Specification

## Purpose
TBD - created by archiving change add-ssm-service-action. Update Purpose after archive.
## Requirements
### Requirement: Subscribe the current account to SSM

The action `tencentcloud_ssm_service` SHALL subscribe the current account to Secrets Manager by calling
`CreateService` of `ssm/v20190923`.
As a Terraform action it MUST keep no state, MUST implement only `Invoke` (plus `Metadata` / `Schema` /
`Configure`), and MUST NOT expose any `Read` / `Update` / `Delete` lifecycle.

#### Scenario: Invoke with valid inputs

- **WHEN** the user runs `terraform apply -invoke=action.tencentcloud_ssm_service.<name>` with a valid
  `pay_mode` and `client_token`
- **THEN** the action calls `CreateService` with `PayMode` and `ClientToken` taken from the configuration
  wrapped in the project retry mechanism (`tccommon.WriteRetryTimeout` + `tccommon.RetryError`)
- **AND** the same `client_token` value is reused across every retry of that single invocation
- **AND** the returned `resource_id` and `order_id` are reported to the operator

#### Scenario: Required inputs are missing or invalid

- **WHEN** `pay_mode` or `client_token` is null, unknown or empty, or `pay_mode` is not `PostPay`, or
  `client_token` is longer than 64 characters
- **THEN** the action returns an error diagnostic before calling any cloud API

#### Scenario: Subscription already exists

- **WHEN** `CreateService` fails with `ResourceInUse.ServiceAlreadyEnabled`
- **THEN** the action treats the desired end state as already reached, records it in the log and continues
  to the status-verification step instead of failing

#### Scenario: Cloud API failure

- **WHEN** `CreateService` fails with any other error
- **THEN** the action returns an error diagnostic containing the API name and the failure reason

### Requirement: Verify subscription effectiveness before reporting success

The action MUST poll `GetServiceStatus` until `ServiceEnabled` is `true` before reporting success, because
`CreateService` is asynchronous and its response does not imply the service is usable.

#### Scenario: Subscription converges

- **WHEN** `GetServiceStatus` returns `ServiceEnabled = true` during polling
- **THEN** the action reports the observed subscription snapshot and returns without diagnostics

#### Scenario: Subscription does not converge

- **WHEN** `ServiceEnabled` never becomes `true` within the polling budget derived from
  `tccommon.WriteRetryTimeout`
- **THEN** the action returns an error diagnostic that includes the last observed `ServiceEnabled`,
  `invalid_type` and `resource_id` values

### Requirement: Expose the subscription snapshot as invoke-time observations

The action SHALL report the read-only fields of `GetServiceStatus` — `service_enabled`, `invalid_type`,
`access_key_escrow_enabled`, `expire_time`, `qps_limit`, `secret_limit`, `pay_model`, `renew_flag`,
`resource_id` and `resource_region` — as observations of the invocation. Because an action has no state
and its only output channel is `InvokeResponse.Diagnostics` / `SendProgress`, these fields SHALL NOT be
declared as Terraform schema attributes.

#### Scenario: Snapshot reporting

- **WHEN** the action has verified that the subscription is effective
- **THEN** it reports the ten observation values through progress messages
- **AND** its schema contains exactly the two input attributes `pay_mode` and `client_token`

#### Scenario: Nil-safe conversion

- **WHEN** any `GetServiceStatus` response field is nil
- **THEN** the action renders it as an empty/unknown marker in the reported snapshot instead of panicking

### Requirement: Unsubscribing is out of scope and documented

The action SHALL document that no API can revert the subscription, that `terraform destroy` therefore does
not unsubscribe, and that every invoke performs a real API call.

#### Scenario: Documentation contract

- **WHEN** the action's schema `Description` and its example `.md` are reviewed
- **THEN** both state that repeated invokes result in repeated operations against the cloud API
- **AND** both state that destroying the action configuration does not roll back the subscription

### Requirement: The action is registered as the first framework action of this provider

The action factory `NewSsmService()` SHALL be registered in `actionFactories` of
`tencentcloud/framework/registry.go`, and the type name SHALL be listed under the SSM node of
`tencentcloud/provider.md`. `tencentcloud/provider.go` SHALL NOT be modified.

#### Scenario: Registration

- **WHEN** the provider starts
- **THEN** `framework.Actions()` returns exactly one factory, `ssm.NewSsmService`
- **AND** the resulting `Metadata.TypeName` is `tencentcloud_ssm_service`

#### Scenario: Client injection without the fw helper package

- **WHEN** the framework runtime calls `Configure` with `req.ProviderData == nil`
- **THEN** the action returns early without diagnostics
- **WHEN** it is called with a value that is not `*sharedmeta.ProviderMeta`
- **THEN** the action emits an `Unexpected ProviderData type` error diagnostic and does not panic

