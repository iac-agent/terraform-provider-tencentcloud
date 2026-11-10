---
subcategory: "AGS(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_tool"
sidebar_current: "docs-tencentcloud-resource-ags_sandbox_tool"
description: |-
  Provides a resource to create a sandbox tool of AGS (Agent Development Platform)
---

# tencentcloud_ags_sandbox_tool

Provides a resource to create a sandbox tool of AGS (Agent Development Platform)

~> **Note:** The `tool_name`, `tool_type`, `default_timeout`, `client_token`, `role_arn`, `storage_mounts`, `log_configuration` and `persistent` attributes are immutable: the `UpdateSandboxTool` API does not accept them, so changing any of them forces resource replacement.

## Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "my-tool"
  tool_type = "browser"

  network_configuration {
    network_mode = "PUBLIC"
  }

  description = "sandbox tool example"
  persistent  = false

  tags {
    key   = "env"
    value = "prod"
  }
}

resource "tencentcloud_ags_sandbox_tool" "vpc_custom" {
  tool_name = "my-custom-tool"
  tool_type = "custom"

  network_configuration {
    network_mode = "VPC"

    vpc_config {
      subnet_ids         = ["subnet-xxxx"]
      security_group_ids = ["sg-xxxx"]
    }
  }

  custom_configuration {
    image               = "registry.example.com/my-image:1.0"
    image_registry_type = "personal"
    command             = ["/bin/sh"]
    args                = ["-c", "run"]

    env {
      name  = "FOO"
      value = "bar"
    }

    ports {
      name     = "http"
      port     = 8080
      protocol = "TCP"
    }

    resources {
      cpu    = "2"
      memory = "4Gi"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_configuration` - (Required, List) Network configuration.
* `tool_name` - (Required, String) Sandbox tool name, 1-50 characters, supports English letters, digits, underscores and hyphens. The name must be unique under the same AppId. Immutable: the Update API does not accept it, so changing it forces replacement.
* `tool_type` - (Required, String) Sandbox tool type. Valid values: `browser`, `code-interpreter`, `mobile`, `osworld`, `custom`, `swebench`, `aio`, `android-world`, `waa`. Immutable: the Update API does not accept it, so changing it forces replacement.
* `client_token` - (Optional, String) Idempotency token for the create call, up to 64 characters. Immutable after creation.
* `computer_configuration` - (Optional, List) Desktop computer environment sandbox configuration.
* `custom_configuration` - (Optional, List) Custom sandbox tool configuration.
* `default_timeout` - (Optional, String) Default timeout, supports formats such as `5m`, `300s`, `1h`. Defaults to the system value (5 minutes) when unset, maximum 24 hours. Immutable: the Update API does not accept it, so changing it forces replacement. Note: the Describe API returns the value in seconds (default_timeout_seconds), so drift detection on this field is limited.
* `description` - (Optional, String) Sandbox tool description, up to 200 characters.
* `log_configuration` - (Optional, List) Sandbox tool log pushing configuration. Immutable: the Update API does not accept it, so changing it forces replacement.
* `persistent` - (Optional, Bool) Whether the sandbox is resident (persistent). Immutable: the Update API does not accept it, so changing it forces replacement.
* `role_arn` - (Optional, String) Role ARN. Immutable: the Update API does not accept it, so changing it forces replacement.
* `storage_mounts` - (Optional, List) Storage mount configurations of the sandbox tool. Immutable: the Update API does not accept it, so changing it forces replacement.
* `tags` - (Optional, List) Tag specification, binds tags to the sandbox tool.

The `computer_configuration` object supports the following:

* `os_world_configuration` - (Optional, List) Built-in OSWorld configuration.
* `waa_configuration` - (Optional, List) WAA sandbox tool configuration.

The `os_world_configuration` object of `computer_configuration` supports the following:

* `version` - (Optional, String) Built-in OSWorld version. Valid values: `osworld1`, `osworld2`. Defaults to `osworld1`.

The `waa_configuration` object of `computer_configuration` supports the following:

* `image_id` - (Optional, String) Custom WAA image ID.

The `custom_configuration` object supports the following:

* `args` - (Optional, Set) Start arguments.
* `command` - (Optional, Set) Start command.
* `dns_config` - (Optional, List) Sandbox DNS configuration.
* `env` - (Optional, List) Environment variables.
* `image_registry_type` - (Optional, String) Image registry type. Valid values: `enterprise`, `personal`, `custom`.
* `image` - (Optional, String) Image address.
* `ports` - (Optional, List) Port configurations.
* `probe` - (Optional, List) Probe configuration.
* `resources` - (Optional, List) Resource configuration.
* `security_context` - (Optional, List) Sandbox security context.

The `dns_config` object of `custom_configuration` supports the following:

* `servers` - (Required, Set) DNS server addresses, valid IP addresses. Defaults to `10.0.0.1`.
* `options` - (Optional, Set) Options (the `options` directive of resolv.conf).
* `searches` - (Optional, Set) Search domains (the `search` directive of resolv.conf).

The `env` object of `custom_configuration` supports the following:

* `name` - (Optional, String) Environment variable name.
* `value` - (Optional, String) Environment variable value.

The `ports` object of `custom_configuration` supports the following:

* `name` - (Optional, String) Port name.
* `port` - (Optional, Int) Port number.
* `protocol` - (Optional, String) Protocol.

The `probe` object of `custom_configuration` supports the following:

* `failure_threshold` - (Optional, Int) Health check failure threshold.
* `http_get` - (Optional, List) HTTP GET probe configuration.
* `probe_period_ms` - (Optional, Int) Health check interval in milliseconds.
* `probe_timeout_ms` - (Optional, Int) Health check single probe timeout in milliseconds.
* `ready_timeout_ms` - (Optional, Int) Health check readiness timeout in milliseconds.
* `success_threshold` - (Optional, Int) Health check success threshold.

The `http_get` object of `probe` supports the following:

* `path` - (Optional, String) Path.
* `port` - (Optional, Int) Port number.
* `scheme` - (Optional, String) Scheme.

The `resources` object of `custom_configuration` supports the following:

* `cpu` - (Optional, String) CPU resource amount.
* `gpu_resource_group_id` - (Optional, String) Associated GPU resource group ID.
* `memory` - (Optional, String) Memory resource amount.
* `storage` - (Optional, String) Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`.

The `security_context` object of `custom_configuration` supports the following:

* `capabilities` - (Optional, List) POSIX capabilities added to or removed from the running container.
* `privileged` - (Optional, Bool) Whether privileged mode is enabled. Defaults to true.

The `capabilities` object of `security_context` supports the following:

* `add` - (Optional, Set) Capabilities to add, e.g. `["SYS_ADMIN"]`.
* `drop` - (Optional, Set) Capabilities to drop, e.g. `["ALL"]`.

The `log_configuration` object supports the following:

* `cls_config` - (Optional, List) CLS log pushing configuration.
* `log_sources` - (Optional, List) Log source configuration.

The `cls_config` object of `log_configuration` supports the following:

* `topic_id` - (Optional, String) CLS log topic ID used by the sandbox tool log pushing.

The `log_sources` object of `log_configuration` supports the following:

* `files` - (Optional, Set) Log file paths to collect, must be files under the `/logs/` directory, subdirectories are not supported, up to 10 files.

The `network_configuration` object supports the following:

* `network_mode` - (Required, String) Network mode. Valid values: `PUBLIC`, `VPC`, `SANDBOX`.
* `vpc_config` - (Optional, List) VPC network configuration.

The `vpc_config` object of `network_configuration` supports the following:

* `security_group_ids` - (Optional, Set) Security group ID list.
* `subnet_ids` - (Optional, Set) VPC subnet ID list.

The `storage_mounts` object supports the following:

* `mount_path` - (Optional, String) Local mount path inside the sandbox instance.
* `name` - (Optional, String) Storage mount configuration name.
* `read_only` - (Optional, Bool) Whether the mount is read-only, defaults to false.
* `storage_source` - (Optional, List) Storage configuration. This field may return null.
* `volume_template` - (Optional, List) VolumeTemplate reference and derivation rules. This field may return null.
* `volume` - (Optional, List) Volume reference. This field may return null.

The `storage_source` object of `storage_mounts` supports the following:

* `agent_bucket` - (Optional, List) AgentBucket storage configuration.
* `cfs` - (Optional, List) File storage (CFS) configuration.
* `cos` - (Optional, List) Object storage (COS) bucket configuration.
* `image` - (Optional, List) Image volume configuration.

The `agent_bucket` object of `storage_source` supports the following:

* `access_domain` - (Optional, String) AgentBucket AccessDomain.
* `library_id` - (Optional, String) AgentBucket LibraryID.
* `space_id` - (Optional, String) AgentBucket spaceId.

The `cfs` object of `storage_source` supports the following:

* `file_system_id` - (Optional, String) CFS resource ID.
* `path` - (Optional, String) CFS mount path.

The `cos` object of `storage_source` supports the following:

* `bucket_name` - (Optional, String) Object storage bucket name.
* `bucket_path` - (Optional, String) Object storage bucket path, must be an absolute path starting with `/`.
* `endpoint` - (Optional, String) Object storage access domain.

The `image` object of `storage_source` supports the following:

* `digest` - (Optional, String) Image digest. Not required on requests.
* `image_registry_type` - (Optional, String) Image registry type. Valid values: `enterprise`, `personal`.
* `reference` - (Optional, String) Image address.
* `sub_path` - (Optional, String) Path inside the image.

The `volume` object of `storage_mounts` supports the following:

* `volume_id` - (Optional, String) Stable Volume ID, format `vol-xxxxxxxx`.
* `volume_name` - (Optional, String) Volume name; resolved server-side to a stable VolumeId under the current AppId.

The `volume_template` object of `storage_mounts` supports the following:

* `volume_template_name` - (Required, String) VolumeTemplate name.
* `reuse_key` - (Optional, String) Reuse key, a literal or `${metadata.xxx}`.

The `tags` object supports the following:

* `key` - (Optional, String) Tag key.
* `value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the sandbox tool, ISO8601 format.
* `id` - Resource identifier, the sandbox tool ID.
* `status_reason` - Reason why the sandbox tool is in its current status.
* `status` - Sandbox tool status. Valid values: `CREATING`, `ACTIVE`, `DELETING`, `FAILED`.
* `tool_id` - The created sandbox tool ID.
* `update_time` - Last update time of the sandbox tool, ISO8601 format.


## Import

AGS sandbox tools can be imported using the tool id, e.g.

```
terraform import tencentcloud_ags_sandbox_tool.example tool-1234567890abcdef
```

