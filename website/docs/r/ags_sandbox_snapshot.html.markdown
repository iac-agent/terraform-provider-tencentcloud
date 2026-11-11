---
subcategory: "AGS(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_snapshot"
sidebar_current: "docs-tencentcloud-resource-ags_sandbox_snapshot"
description: |-
  ---
subcategory: ""
layout: "ags"
page_title: "TencentCloud: tencentcloud_ags_sandbox_snapshot"
description: |-
  Provides a resource to create an AGS (TencentCloud Agent Sandbox) sandbox snapshot.
---
---

# tencentcloud_ags_sandbox_snapshot

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

### ## Import

AGS sandbox snapshots can be imported using the snapshot id, e.g.

```hcl
shell
terraform import tencentcloud_ags_sandbox_snapshot.example snap-xxxxxxxx
```

## Argument Reference

The following arguments are supported:

* `client_token` - (Optional, String) Idempotency token used when creating the snapshot.
* `instance_id` - (Optional, String) ID of the sandbox instance from which the snapshot is created.
* `memory` - (Optional, Bool) Whether the snapshot includes the memory state of the sandbox instance.
* `snapshot_name` - (Optional, String) Name of the sandbox snapshot.
* `tags` - (Optional, List) Tags of the sandbox snapshot.

The `tags` object supports the following:

* `key` - (Optional, String) Tag key.
* `value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the sandbox snapshot in ISO8601 format.
* `id` - ID of the sandbox snapshot, the same value as snapshot_id returned by the Create API.
* `status_reason` - Reason why the sandbox snapshot is in its current status.
* `status` - Status of the sandbox snapshot. Valid values: ACTIVE, DELETING, DELETED.
* `update_time` - Last update time of the sandbox snapshot in ISO8601 format.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateSandboxSnapshot](https://www.tencentcloud.com/search/?wd=CreateSandboxSnapshot&type=document&lang=en)
  - [DescribeSandboxSnapshotList](https://www.tencentcloud.com/search/?wd=DescribeSandboxSnapshotList&type=document&lang=en)
* Read: [DescribeSandboxSnapshotList](https://www.tencentcloud.com/search/?wd=DescribeSandboxSnapshotList&type=document&lang=en)
* Update:
  - [GetPasswordRules](https://www.tencentcloud.com/document/api/598/81237)
  - [UpdatePasswordRules](https://www.tencentcloud.com/document/api/598/81238)
* Delete: [DeleteSandboxSnapshot](https://www.tencentcloud.com/search/?wd=DeleteSandboxSnapshot&type=document&lang=en)

