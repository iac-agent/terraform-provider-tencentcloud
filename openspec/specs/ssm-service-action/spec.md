# ssm-service-action Specification

## Purpose
Defines the `tencentcloud_ssm_service` Terraform Plugin Framework action that activates TencentCloud Secrets Manager (SSM) for the invoked account and region via the `CreateService` API, including its schema, invocation semantics, already-enabled idempotency, convergence polling, registration, tests and documentation.

## Requirements

### Requirement: Action Definition and Resource Kind
The system SHALL provide a Terraform Plugin Framework action named `tencentcloud_ssm_service`
(`RESOURCE_KIND_OPERATION`) that activates TencentCloud Secrets Manager (SSM) for the invoked
account and region by calling the `CreateService` API of `ssm/v20190923`.

The action SHALL be implemented with `action.Action` from
`github.com/hashicorp/terraform-plugin-framework/action`, embedding `fw.ActionWithConfigure`, and
SHALL implement only `Metadata`, `Schema` and `Invoke` (plus the promoted `Configure`). It SHALL NOT
implement `Create`, `Read`, `Update`, `Delete`, `ImportState` or `ModifyPlan`, SHALL NOT maintain any
Terraform state, and SHALL NOT expose any output attribute (an action invocation cannot publish values
to state).

#### Scenario: Action type name is the requested Terraform name
- **WHEN** the framework calls `Metadata` on the action
- **THEN** `resp.TypeName` SHALL be `tencentcloud_ssm_service`

#### Scenario: No lifecycle state
- **WHEN** the action is invoked successfully
- **THEN** no resource id SHALL be created, no state entry SHALL be written, and no value SHALL be
  published as an action output

#### Scenario: Destroy does not deactivate SSM
- **WHEN** the configuration containing the action is destroyed
- **THEN** the provider SHALL NOT call any deactivation API, because the SSM API surface used here has
  no deactivate operation and the action holds no state

### Requirement: Action Schema
The system SHALL define the action schema with exactly the following attributes:

- `pay_mode` (String, **Required**): billing mode forwarded to `CreateServiceRequest.PayMode`.
  The only value accepted by the API is `PostPay`.
- `client_token` (String, **Optional + Computed**): idempotency token forwarded to
  `CreateServiceRequest.ClientToken`. When omitted by the user the action SHALL generate one; when
  supplied by the user it SHALL be replayed verbatim on every retry and on later invokes.
- `timeouts` (single nested block, Optional): built with the official
  `terraform-plugin-framework-timeouts` `actiontimeouts.BlockWithOpts`, exposing only the `invoke`
  stage. Its default SHALL be the action default of 10 minutes
  (`defaultSsmServiceInvokeTimeout`).

No response field of `CreateService` (`ResourceId`, `OrderId`) or of `GetServiceStatus`
(`ServiceEnabled`, `InvalidType`, `AccessKeyEscrowEnabled`, `ExpireTime`, `QPSLimit`, `SecretLimit`,
`PayModel`, `RenewFlag`, `ResourceId`, `ResourceRegion`) SHALL be declared as a schema attribute,
because a Framework action cannot publish output values.

#### Scenario: Schema exposes the two API inputs
- **WHEN** the action schema is inspected
- **THEN** it SHALL contain `pay_mode` as a required string attribute, `client_token` as an
  optional-and-computed string attribute, and a `timeouts` block with an `invoke` attribute

#### Scenario: Response fields are not attributes
- **WHEN** the action schema is inspected
- **THEN** it SHALL NOT contain `resource_id`, `order_id`, `service_enabled`, `invalid_type`,
  `access_key_escrow_enabled`, `expire_time`, `qps_limit`, `secret_limit`, `pay_model`,
  `renew_flag` or `resource_region`

#### Scenario: Description states one-time and non-reversible semantics
- **WHEN** the schema `Description` is read
- **THEN** it SHALL state that activation is a one-time operation, that repeated invokes re-trigger the
  remote call, and that destroying the configuration does not deactivate SSM

### Requirement: Invoke Input Validation
The system SHALL validate all inputs at the beginning of `Invoke`, before any cloud API call, using the
Framework tri-state checks (`IsNull()` / `IsUnknown()` / empty string). Every validation failure SHALL
be reported through `resp.Diagnostics.AddError` and SHALL return without calling any API.

- `pay_mode` null, unknown or empty → error summary `Missing pay_mode`.
- `pay_mode` different from `PostPay` → error summary `Invalid pay_mode`, detail explaining that the
  `CreateService` API only supports `PostPay`.
- `client_token` provided but not matching `^[\x21-\x7E]{1,64}$` → error summary `Invalid client_token`,
  detail quoting the API constraint (1–64 visible ASCII characters, no space, no CJK, no control chars).
- Provider client (`a.Client()`) nil → error summary `Provider not configured`.

#### Scenario: Missing pay_mode is rejected
- **WHEN** `pay_mode` is null, unknown or an empty string
- **THEN** `Invoke` SHALL add the `Missing pay_mode` error diagnostic and SHALL NOT call `CreateService`

#### Scenario: Unsupported pay_mode is rejected
- **WHEN** `pay_mode` is set to any value other than `PostPay`
- **THEN** `Invoke` SHALL add the `Invalid pay_mode` error diagnostic and SHALL NOT call `CreateService`

#### Scenario: Malformed client_token is rejected
- **WHEN** `client_token` is set but is longer than 64 characters, empty, or contains a space,
  a non-ASCII character or a control character
- **THEN** `Invoke` SHALL add the `Invalid client_token` error diagnostic and SHALL NOT call `CreateService`

#### Scenario: Provider client unavailable
- **WHEN** the injected provider client is nil
- **THEN** `Invoke` SHALL add the `Provider not configured` error diagnostic and SHALL NOT call any API

### Requirement: CreateService Invocation
The system SHALL implement a service-layer method
`SsmService.CreateSsmService(ctx context.Context, payMode, clientToken *string)` in
`tencentcloud/services/ssm/service_tencentcloud_ssm.go` that calls
`ssm.CreateServiceRequest` / `UseSsmClient().CreateService`, following the SDKv2 conventions of that file
(`ratelimit.Check(request.GetAction())`, the `[CRITAL]` failure `defer` log, and the `[DEBUG]` success log).
The `Invoke` handler SHALL call it inside `resource.Retry(tccommon.WriteRetryTimeout, …)` and wrap failures
with `tccommon.RetryError`, without nested retries.

The `client_token` value SHALL be resolved exactly once, before entering the retry block, so every retry of
the same invoke reuses the identical token. When the user omits `client_token`, the action SHALL generate
`ssm-<helper.BuildToken()>`.

#### Scenario: Successful activation request
- **WHEN** `Invoke` runs with `pay_mode = "PostPay"` and no `client_token`
- **THEN** the action SHALL generate one token of the form `ssm-<random>`, call `CreateService` once with
  `PayMode = "PostPay"` and that `ClientToken`, and SHALL NOT produce any error diagnostic

#### Scenario: Returned identifiers are logged, not stored
- **WHEN** `CreateService` returns `Response.ResourceId` and/or `Response.OrderId`
- **THEN** the action SHALL log them (nil-safe) and SHALL NOT attempt to persist them in state

#### Scenario: Transient create failure is retried
- **WHEN** `CreateService` returns a retryable error
- **THEN** the service/action SHALL retry within `tccommon.WriteRetryTimeout` using the same `ClientToken`
  for every attempt

#### Scenario: Non-retryable create failure surfaces as a diagnostic
- **WHEN** `CreateService` fails after retries are exhausted with an error other than
  `ResourceInUse.ServiceAlreadyEnabled`
- **THEN** `Invoke` SHALL add an error diagnostic whose summary names the failed activation and whose
  detail contains the API error

#### Scenario: Empty response is handled safely
- **WHEN** `CreateService` returns a nil response or a nil `Response` field without an error
- **THEN** the action SHALL treat it as a failure diagnostic rather than dereferencing nil

### Requirement: Already-Enabled Idempotency
The system SHALL treat the error code `ResourceInUse.ServiceAlreadyEnabled`
(`ssm.RESOURCEINUSE_SERVICEALREADYENABLED`) as "already at the target state", not as a failure. The check
SHALL be implemented as a package-level helper so that it is unit-testable, and SHALL be detected both from
a `*tcerr.TencentCloudSDKError` code and from the error string as a fallback.

#### Scenario: Already-enabled error continues to convergence
- **WHEN** `CreateService` returns `ResourceInUse.ServiceAlreadyEnabled`
- **THEN** the action SHALL log the condition, report it through `resp.SendProgress`, and continue to the
  `GetServiceStatus` convergence check without producing an error diagnostic

#### Scenario: Already-enabled error is recognised by a helper
- **WHEN** the helper is given an error carrying code `ResourceInUse.ServiceAlreadyEnabled`
- **THEN** it SHALL return true, and it SHALL return false for any other error

### Requirement: Convergence Wait
The system SHALL poll `GetServiceStatus` (through the existing
`SsmService.DescribeSsmServiceStatusByFilter`) after a successful or already-enabled `CreateService`
call, and SHALL only complete successfully once `ServiceEnabled == true`.

The polling budget SHALL come from the `timeouts.invoke` value, defaulting to
`defaultSsmServiceInvokeTimeout` (10 minutes). The poll interval SHALL be a package-level `var`
(`ssmServicePollInterval`, default 5 seconds) so that unit tests can shorten it. Progress SHALL be
reported through `resp.SendProgress` guarded by a nil check. No additional retry wrapper SHALL be nested
around the polling call.

#### Scenario: Service enabled is reported as success
- **WHEN** `GetServiceStatus` returns `ServiceEnabled == true` within the budget
- **THEN** `Invoke` SHALL return without error diagnostics

#### Scenario: Terminal invalid state fails fast
- **WHEN** `GetServiceStatus` returns `ServiceEnabled == false` together with `InvalidType == 2`
  (suspended due to arrears) or `InvalidType == 3` (resource released)
- **THEN** `Invoke` SHALL stop polling immediately and add an error diagnostic describing the terminal state

#### Scenario: Timeout reports the last observed snapshot
- **WHEN** the polling budget is exhausted while `ServiceEnabled` is still false
- **THEN** `Invoke` SHALL add an error diagnostic whose summary states that SSM did not converge and whose
  detail contains the last observed `service_enabled`, `invalid_type` and `resource_id` values

#### Scenario: Nil response during polling is not fatal
- **WHEN** a polling attempt returns a nil response or nil `Response` without an error
- **THEN** the action SHALL continue polling until the budget is exhausted instead of dereferencing nil

### Requirement: Action Registration
The system SHALL register the action factory `ssm.NewSsmServiceAction()` in
`tencentcloud/framework/registry.go`'s `actionFactories` slice, and SHALL NOT register the action in the
SDKv2 `tencentcloud/provider.go` `ResourcesMap`.

The factory SHALL be named `NewSsmServiceAction` because the `ssm` package already declares
`NewSsmService(client *connectivity.TencentCloudClient) SsmService` and Go does not allow overloads.

#### Scenario: Action factory is collected
- **WHEN** the framework provider collects action factories
- **THEN** `frameworkActions()` SHALL include `ssm.NewSsmServiceAction`

#### Scenario: Action is not an SDKv2 resource
- **WHEN** the SDKv2 provider resources map is built
- **THEN** it SHALL NOT contain `tencentcloud_ssm_service`

### Requirement: Unit Tests
The system SHALL provide `tencentcloud/services/ssm/action_tc_ssm_service_test.go` using gomonkey to mock
the SSM cloud API (`UseSsmClient`, `CreateService`, `GetServiceStatus`) instead of the Terraform acceptance
test suite, covering only business logic.

#### Scenario: Successful invoke
- **WHEN** a test invokes the action with `pay_mode = "PostPay"` and mocks `CreateService` to succeed and
  `GetServiceStatus` to return `ServiceEnabled = true`
- **THEN** the test SHALL assert the request carried `PayMode = "PostPay"` with a non-empty `ClientToken`,
  and that no error diagnostic was produced

#### Scenario: Already-enabled invoke succeeds
- **WHEN** a test mocks `CreateService` to return the `ResourceInUse.ServiceAlreadyEnabled` error and
  `GetServiceStatus` to return `ServiceEnabled = true`
- **THEN** the test SHALL assert no error diagnostic was produced

#### Scenario: Create API error is reported
- **WHEN** a test mocks `CreateService` to return an arbitrary error
- **THEN** the test SHALL assert an error diagnostic containing the API error is produced

#### Scenario: Validation failures skip the API
- **WHEN** a test invokes with a missing or non-`PostPay` `pay_mode`, or a malformed `client_token`
- **THEN** the test SHALL assert an error diagnostic is produced and that `CreateService` was never called

#### Scenario: Convergence timeout is reported
- **WHEN** a test mocks `GetServiceStatus` to always return `ServiceEnabled = false` with a shortened poll
  interval and a short `invoke` timeout
- **THEN** the test SHALL assert an error diagnostic is produced and that it contains the last observed
  snapshot

### Requirement: Documentation
The system SHALL provide `tencentcloud/services/ssm/action_tc_ssm_service.md` with a one-line description
mentioning SSM (Secrets Manager), an `Example Usage` section using the `action` block with a `config` block
and the `terraform apply -invoke=...` command, and a note that actions require Terraform 1.14 or later.
The file SHALL NOT contain an `Import` section and SHALL NOT contain hand-written
`Argument Reference` / `Attribute Reference` sections, which are generated by the documentation tooling.
The website document `website/docs/a/ssm_service.html.markdown` SHALL be produced by `make doc`, never
hand-written.

#### Scenario: Example file has the required sections
- **WHEN** the action is delivered
- **THEN** `action_tc_ssm_service.md` SHALL exist with a one-line description mentioning SSM, an
  `Example Usage` block containing `action "tencentcloud_ssm_service" "example" { config { ... } }`,
  a Terraform 1.14+ note, and no `Import` section

#### Scenario: provider.md lists the action
- **WHEN** `tencentcloud/provider.md` is read
- **THEN** under the `Secrets Manager(SSM)` node there SHALL be an `Action` section listing
  `tencentcloud_ssm_service`
