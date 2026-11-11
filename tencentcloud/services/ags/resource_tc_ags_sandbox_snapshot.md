---
subcategory: ""
layout: "ags"
page_title: "TencentCloud: tencentcloud_ags_sandbox_snapshot"
description: |-
  Provides a resource to create an AGS (TencentCloud Agent Sandbox) sandbox snapshot.
---

# tencentcloud_ags_sandbox_snapshot

Provides a resource to create an AGS (TencentCloud Agent Sandbox) sandbox snapshot.

~> The AGS API provides no modify interface for sandbox snapshots, so every
writable attribute forces replacement: any change to `instance_id`,
`snapshot_name`, `memory`, `tags` or `client_token` destroys the snapshot and
creates a new one.

~> Drift detection is limited: the describe API does not return data that can
be mapped back to `instance_id`, `memory` or `client_token` with confirmed
semantics, so those values are kept in state and never refreshed.

## Example Usage

```hcl
resource "tencentcloud_ags_sandbox_snapshot" "example" {
  instance_id   = "sandbox-xxxxxxxx"
  snapshot_name = "my-snapshot"
  memory        = true

  tags {
    key   = "env"
    value = "test"
  }

  client_token = "unique-idempotency-token"
}
```

## Import

AGS sandbox snapshots can be imported using the snapshot id, e.g.

```shell
terraform import tencentcloud_ags_sandbox_snapshot.example snap-xxxxxxxx
```
