# Design: reply-exactly-ok

## Context

The incoming request is:

> reply exactly: OK

It carries two constraints that fully determine the design:

1. **Output constraint** — the assistant must reply with exactly `OK`.
2. **Process constraint** — the task must still go through the openspec propose workflow, producing a complete change under `openspec/changes/` that passes `openspec validate --strict`.

There is no third constraint asking for any Tencent Cloud Terraform resource, data source or field. Therefore the design problem is not "how to model a cloud object in Terraform", it is "how to represent a no-op request in openspec without smuggling in unintended provider changes".

## Design Goals

| Goal | Rationale |
|---|---|
| Zero provider change | The request asks for a reply, not a cloud resource. Any schema/CRUD edit would be unrequested scope creep and could break existing TF state. |
| Zero cloud API dependency | Nothing in the request maps to a Tencent Cloud API, so there is nothing to cross-check against `vendor/`. |
| Fully validatable change | The change must still satisfy `openspec validate --strict` so the workflow does not special-case it. |
| Self-describing tasks | A reviewer must be able to see, from `tasks.md` alone, that "no code change" is the intended outcome and not an oversight. |

## Non-Goals

- Adding a `tencentcloud_*` resource or data source.
- Adding a service-layer wrapper in `tencentcloud/services/**`.
- Registering anything in `tencentcloud/provider.go` / `tencentcloud/provider.md`.
- Producing `website/docs/**` content.
- Producing `.changelog/**` entries.

## Decision: model the request as a workflow/traceability capability

Since there is no cloud object to model, the capability introduced by this change is named after what the change actually is: traceability of an openspec workflow request. Naming it after a service (e.g. `add-xxx-resource`) would be misleading, and would invite a future reader to search for a nonexistent resource.

**Alternatives considered and rejected:**

1. *Skip the openspec workflow and only reply* — rejected: the task explicitly requires a complete change under `openspec/changes/`.
2. *Invent a trivial resource so the change "looks normal"* — rejected: it would add unrequested public API surface to the provider, and cannot be justified by any vendored cloud API, violating the constraint that only cloud-API-supported behaviour may be added to the provider.
3. *Create an empty change with only `.openspec.yaml`* — rejected: `openspec validate --strict` requires `proposal.md`, `tasks.md` and at least one spec delta with `### Requirement:` and `#### Scenario:` blocks.

## File Layout

| File | Action | Purpose |
|---|---|---|
| `openspec/changes/reply-exactly-ok/proposal.md` | New | Why / What / API mapping / Capabilities / Impact, stating the change is non-functional |
| `openspec/changes/reply-exactly-ok/design.md` | New (this file) | Design reasoning, rejected alternatives, file layout |
| `openspec/changes/reply-exactly-ok/tasks.md` | New | Task checklist proving no code, vendor or docs change is intended |
| `openspec/changes/reply-exactly-ok/specs/openspec-workflow-request-traceability/spec.md` | New | `## ADDED Requirements` with `### Requirement:` + `#### Scenario:` blocks |
| `tencentcloud/provider.go` | Unchanged | — |
| `tencentcloud/provider.md` | Unchanged | — |
| `tencentcloud/services/**` | Unchanged | — |
| `vendor/**` | Unchanged | — |
| `website/**` | Unchanged | — |
| `.changelog/**` | Unchanged | — |

## Verification Strategy

Because there is no code, the standard `go build` / `go vet` / `make doc` verification chain is not applicable and must not be introduced as new work. Verification is instead:

1. `pnpm dlx @fission-ai/openspec@1.2.0 validate reply-exactly-ok --strict` exits 0.
2. `git status --porcelain` shows changes confined to `openspec/changes/reply-exactly-ok/**`.
3. The spec delta contains at least one `### Requirement:` and every requirement contains at least one `#### Scenario:`.
4. The provided reply body is the literal string `OK` with no surrounding prose.

## Backward Compatibility

Terraform configurations and state files are untouched: no schema field is added, removed or re-typed, and no provider registration changes, so no configuration can start failing or produce a new diff because of this change.

## Risks

| Risk | Mitigation |
|---|---|
| A future reader mistakes the change for a missing resource | `proposal.md` "Explicitly out of scope" section and `tasks.md` both state that no code change is intended |
| Reviewer expects `make doc` output | `design.md` and `tasks.md` explain that `make doc` is only required when a resource or data source changes |
| Accidental scope creep while "completing" the change | Tasks are limited to documentation-of-intent and validation steps |
