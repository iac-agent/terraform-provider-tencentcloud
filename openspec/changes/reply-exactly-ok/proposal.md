## Why

The requested scope for this change is a single, self-contained instruction: `reply exactly: OK`. This change records that instruction as an openspec proposal so it is tracked, validated and archived like any other change, instead of living only in chat history.

Because the instruction contains no requirement to add, modify or remove any Tencent Cloud Terraform resource or data source, this proposal is intentionally **non-functional**: it must not introduce, change or deprecate any resource schema, any CRUD implementation, any service-layer wrapper, any provider registration, or any documentation page. The provider's behaviour must remain byte-for-byte identical before and after this change.

The proposal only exists to satisfy the openspec workflow (proposal / design / tasks / spec delta) for a request that has no code impact.

## What Changes

- Nothing in the Terraform provider's runtime behaviour changes.
- No new resource or data source is registered.
- No existing resource or data source schema is modified, and no `Optional` field is added.
- No new Tencent Cloud API is called; no vendored SDK package under `vendor/` is added or upgraded.
- The only artifacts produced are the openspec planning documents of this change (`proposal.md`, `design.md`, `tasks.md`, `specs/openspec-workflow-request-traceability/spec.md`).
- The requested reply text `OK` is required to be emitted verbatim in the assistant's response for this request.

### Explicitly out of scope

- Any change to `tencentcloud/provider.go` or `tencentcloud/provider.md` (the expected result of this step is the openspec change only, not provider registration).
- Any change under `tencentcloud/services/**`.
- Any change under `website/` (generated only via `make doc` in the finalize step).
- Any change under `.changelog/` (created only in the finalize step).

## API Mapping

No API change is involved.

| Previous API | New API | Notes |
|---|---|---|
| None | None | No Tencent Cloud API is added, replaced or deprecated by this change. |

Because there is no API surface change, no `vendor/` SDK capability check is required. The change is verified as feasible by confirming it requires zero cloud API support at all.

## Capabilities

### New Capabilities

- `openspec-workflow-request-traceability`: Introduces requirement-level, verifiable statements that a request whose entire scope is a literal reply instruction is captured by an openspec change that performs no provider mutation, while still producing the requested reply verbatim.

### Modified Capabilities

<!-- None: no existing capability's requirement-level behaviour is altered. -->

## Impact

- **New code**: none.
- **Modified code**: none.
- **Modified vendor**: none.
- **Modified website/docs**: none (must not be hand-edited; `make doc` is not required because no resource or data source changes).
- **New openspec artifacts**:
  - `openspec/changes/reply-exactly-ok/proposal.md`
  - `openspec/changes/reply-exactly-ok/design.md`
  - `openspec/changes/reply-exactly-ok/tasks.md`
  - `openspec/changes/reply-exactly-ok/specs/openspec-workflow-request-traceability/spec.md`
- **Risk**: none for provider users. Terraform state and existing HCL configurations are untouched, so backward compatibility is trivially preserved.
