## ADDED Requirements

### Requirement: Literal reply compliance

When a request's entire scope is an instruction to reply with a specific literal string, the response SHALL contain exactly that string, with no prefix, suffix, explanation, markdown decoration, or trailing commentary.

#### Scenario: Requested literal is emitted verbatim
- **WHEN** a request states `reply exactly: OK`
- **THEN** the response body is the two-character string `OK`
- **AND** the response contains no additional sentences, heading markers, code fences, or list bullets

#### Scenario: Literal is not paraphrased
- **WHEN** the requested literal is `OK`
- **THEN** the response is not replaced by synonyms or expansions such as `Okay`, `Got it`, or `OK.`
- **AND** the response does not re-state or describe the instruction it is obeying

### Requirement: No provider mutation

A request whose scope is a literal reply instruction SHALL NOT cause any change to the Terraform provider's public surface. The set of registered resource types, the set of registered data source types, and every resource/data source schema MUST remain identical before and after the change.

#### Scenario: No resource or data source is registered
- **WHEN** the change is applied
- **THEN** `tencentcloud/provider.go` contains the same `ResourcesMap` and `DataSourcesMap` entries as before
- **AND** `tencentcloud/provider.md` is unmodified
- **AND** the number of `tencentcloud_*` types exposed by the provider is unchanged

#### Scenario: No schema field is added
- **WHEN** a developer diffs the change
- **THEN** no file under `tencentcloud/services/**` is added or modified
- **AND** no `schema.Schema` definition is created, edited, or re-typed
- **AND** no resource `Timeouts` block or `Importer` is introduced or altered

#### Scenario: No cloud API dependency is introduced
- **WHEN** a developer inspects the change
- **THEN** no Tencent Cloud API is invoked and no `service_tencentcloud_*.go` wrapper is added
- **AND** no package under `vendor/` is added, removed, or upgraded
- **AND** `go.mod` and `go.sum` are unmodified

#### Scenario: Existing configurations and state keep working
- **WHEN** an operator runs `terraform plan` on an existing configuration after this change
- **THEN** the plan shows no new diff caused by this change
- **AND** no previously valid configuration becomes invalid

### Requirement: Change traceability and validity

The request SHALL be captured as a complete openspec change so that the intent ("literal reply only, no provider change") is reviewable without reading chat history. The change directory MUST contain `proposal.md`, `design.md`, `tasks.md`, and at least one spec delta under `specs/<capability>/spec.md`.

#### Scenario: Change documents are present
- **WHEN** a reviewer lists `openspec/changes/reply-exactly-ok/`
- **THEN** the directory contains `proposal.md`, `design.md`, `tasks.md`, and `specs/openspec-workflow-request-traceability/spec.md`

#### Scenario: Strict validation passes
- **WHEN** `pnpm dlx @fission-ai/openspec@1.2.0 validate reply-exactly-ok --strict` is executed from the repository root
- **THEN** the command exits with status 0
- **AND** every `### Requirement:` heading is accompanied by at least one `#### Scenario:` block

#### Scenario: Scope of the change is self-evident
- **WHEN** a reviewer reads `proposal.md`
- **THEN** an "Explicitly out of scope" section states that `tencentcloud/provider.go`, `tencentcloud/services/**`, `vendor/**`, `website/**`, and `.changelog/**` are intentionally untouched

#### Scenario: Documentation pipeline is not triggered
- **WHEN** the change is applied without any resource or data source modification
- **THEN** no file under `website/**` is added or edited by this change
- **AND** `make doc` is not required to make the change consistent
