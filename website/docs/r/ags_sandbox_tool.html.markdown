---
subcategory: "Agent Sandbox(AGS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_tool"
sidebar_current: "docs-tencentcloud-resource-ags_sandbox_tool"
description: |-
  Provides a resource to create and manage an AGS (Agent Sandbox) sandbox tool.
---

# tencentcloud_ags_sandbox_tool

Provides a resource to create and manage an AGS (Agent Sandbox) sandbox tool.

A sandbox tool is the runtime template used to create AGS sandbox instances: it defines the tool type,
the network configuration, the storage mounts and, for custom tools, the container image and its runtime
configuration.

~> **NOTE:** Deleting this resource only removes the sandbox tool definition. The COS buckets, CFS file
systems, agent buckets and volumes referenced by `storage_mounts` are never deleted by this resource;
clean them up separately if they are no longer needed.

~> **NOTE:** The following attributes can only be set on creation. Changing them forces a new resource
because `UpdateSandboxTool` has no corresponding parameter: `tool_name`, `tool_type`, `network_configuration`,
`default_timeout`, `role_arn`, `client_token`, `storage_mounts`, `log_configuration` and `persistent`.

~> **NOTE:** `UpdateSandboxTool` does not document whether an omitted argument keeps or clears the previous
value, so removing `description` from the configuration is refused with an error instead of guessing.

~> **NOTE:** Drift detection is partial, because the describe API does not return everything that was sent:
`client_token` is never returned, and `default_timeout` is only returned as a number of seconds while the
input is a duration string. Values managed by the configuration are therefore kept as configured and are
only replaced by the cloud value when the attribute is not managed (for example right after an import).
Changes made outside of Terraform to those attributes are not detected.

## Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "custom" {
  tool_name = "tf-example-tool"
  tool_type = "custom"

  network_configuration {
    network_mode = "VPC"

    vpc_config {
      subnet_ids         = ["subnet-12345678"]
      security_group_ids = ["sg-12345678"]
    }
  }

  description     = "sandbox tool managed by terraform"
  default_timeout = "10m"
  persistent      = true

  tags {
    key   = "env"
    value = "test"
  }

  storage_mounts {
    name       = "workspace"
    mount_path = "/mnt/workspace"
    read_only  = false

    storage_source {
      cos {
        endpoint    = "cos.ap-guangzhou.myqcloud.com"
        bucket_name = "bucket-1250000000"
        bucket_path = "/data"
      }
    }
  }

  custom_configuration {
    image               = "ccr.ccs.tencentcloudcr.com/tf-example/demo:1.0"
    image_registry_type = "personal"
    command             = ["/bin/sh"]
    args                = ["-c", "python server.py"]

    env {
      name  = "LOG_LEVEL"
      value = "info"
    }

    ports {
      name     = "http"
      port     = 8080
      protocol = "TCP"
    }

    resources {
      cpu     = "1"
      memory  = "2Gi"
      storage = "10Gi"
    }

    probe {
      http_get {
        path   = "/healthz"
        port   = 8080
        scheme = "HTTP"
      }

      ready_timeout_ms  = 30000
      probe_timeout_ms  = 1000
      probe_period_ms   = 5000
      success_threshold = 1
      failure_threshold = 3
    }

    dns_config {
      servers = ["10.0.0.1"]
    }

    security_context {
      privileged = true

      capabilities {
        add  = ["SYS_ADMIN"]
        drop = ["ALL"]
      }
    }
  }

  log_configuration {
    cls_config {
      topic_id = "a1b2c3d4-0000-1111-2222-333344445555"
    }

    log_sources {
      files = ["/logs/app.log"]
    }
  }
}
```

### Example Usage - public network tool with a desktop computer environment

```hcl
resource "tencentcloud_ags_sandbox_tool" "computer" {
  tool_name = "tf-example-computer"
  tool_type = "osworld"

  network_configuration {
    network_mode = "PUBLIC"
  }

  computer_configuration {
    os_world_configuration {
      version = "osworld2"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_configuration` - (Required, List) Network configuration of sandbox instances created from this tool. Changing it forces a new resource.
* `tool_name` - (Required, String) Name of the sandbox tool. Length 1-50 characters, supporting English letters, digits, underscores and hyphens. The name must be unique under the same AppId. Changing it forces a new resource.
* `tool_type` - (Required, String) Type of the sandbox tool. Valid values: `browser`, `code-interpreter`, `mobile`, `osworld`, `custom`, `swebench`, `aio`, `android-world`, `waa`. Changing it forces a new resource.
* `client_token` - (Optional, String) Idempotency token of the create request, at most 64 characters. It is only sent on creation and is never read back from the cloud side. Changing it forces a new resource.
* `computer_configuration` - (Optional, List) Desktop computer environment configuration of the sandbox tool.
* `custom_configuration` - (Optional, List) Custom configuration of the sandbox tool, used when `tool_type` is `custom`.
* `default_timeout` - (Optional, String) Default timeout of a sandbox instance, e.g. `5m`, `300s`, `1h`. The system default (5 minutes) is used when omitted, maximum 24 hours. The cloud side only returns this as a number of seconds, so the refreshed value is normalised to the `Ns` form. Changing it forces a new resource.
* `description` - (Optional, String) Description of the sandbox tool. Maximum length 200 characters.
* `log_configuration` - (Optional, List) Log shipping configuration of the sandbox tool. `UpdateSandboxTool` has no log configuration parameter, so changing it forces a new resource.
* `persistent` - (Optional, Bool) Whether the sandbox is persistent. Changing it forces a new resource.
* `role_arn` - (Optional, String) ARN of the role bound to the sandbox tool. Changing it forces a new resource.
* `storage_mounts` - (Optional, List) Storage mount configuration of sandbox instances created from this tool. `UpdateSandboxTool` has no storage mount parameter, so changing it forces a new resource.
* `tags` - (Optional, List) Tags bound to the sandbox tool.

The `computer_configuration` object supports the following:

* `os_world_configuration` - (Optional, List) Built-in OSWorld configuration.
* `waa_configuration` - (Optional, List) WAA sandbox tool configuration.

The `os_world_configuration` object of `computer_configuration` supports the following:

* `version` - (Optional, String) Built-in OSWorld version. Valid values: `osworld1`, `osworld2`. The cloud side defaults to `osworld1`.

The `waa_configuration` object of `computer_configuration` supports the following:

* `image_id` - (Optional, String) ID of the custom WAA image.

The `custom_configuration` object supports the following:

* `args` - (Optional, Set) Startup arguments.
* `command` - (Optional, Set) Startup command.
* `dns_config` - (Optional, List) DNS configuration of the sandbox.
* `env` - (Optional, List) Environment variables.
* `image_registry_type` - (Optional, String) Image registry type. Valid values: `enterprise`, `personal`, `custom`.
* `image` - (Optional, String) Image address.
* `ports` - (Optional, List) Port configuration.
* `probe` - (Optional, List) Health probe configuration.
* `resources` - (Optional, List) Resource configuration.
* `security_context` - (Optional, List) Security context of the sandbox.

The `dns_config` object of `custom_configuration` supports the following:

* `options` - (Optional, Set) Options, mapping to the `options` directive of resolv.conf.
* `searches` - (Optional, Set) Search domains, mapping to the `search` directive of resolv.conf.
* `servers` - (Optional, Set) DNS server addresses. Valid IP addresses are required. The cloud side defaults to `10.0.0.1`.

The `env` object of `custom_configuration` supports the following:

* `name` - (Optional, String) Name of the environment variable.
* `value` - (Optional, String) Value of the environment variable.

The `ports` object of `custom_configuration` supports the following:

* `name` - (Optional, String) Port name.
* `port` - (Optional, Int) Port number.
* `protocol` - (Optional, String) Protocol.

The `probe` object of `custom_configuration` supports the following:

* `failure_threshold` - (Optional, Int) Number of consecutive failures required to be considered unhealthy.
* `http_get` - (Optional, List) HTTP GET probe configuration.
* `probe_period_ms` - (Optional, Int) Interval between two probes, in milliseconds.
* `probe_timeout_ms` - (Optional, Int) Timeout of a single probe, in milliseconds.
* `ready_timeout_ms` - (Optional, Int) Ready timeout of the health check, in milliseconds.
* `success_threshold` - (Optional, Int) Number of consecutive successes required to be considered healthy.

The `http_get` object of `probe` supports the following:

* `path` - (Optional, String) Request path.
* `port` - (Optional, Int) Request port.
* `scheme` - (Optional, String) Request scheme.

The `resources` object of `custom_configuration` supports the following:

* `cpu` - (Optional, String) CPU resource amount.
* `memory` - (Optional, String) Memory resource amount.
* `storage` - (Optional, String) Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`.

The `security_context` object of `custom_configuration` supports the following:

* `capabilities` - (Optional, List) POSIX capabilities that can be added to or removed from the running container.
* `privileged` - (Optional, Bool) Whether to enable privileged mode. The cloud side defaults to `true`.

The `capabilities` object of `security_context` supports the following:

* `add` - (Optional, Set) Capabilities to add, e.g. `SYS_ADMIN`.
* `drop` - (Optional, Set) Capabilities to remove, e.g. `ALL`.

The `log_configuration` object supports the following:

* `cls_config` - (Optional, List) CLS log shipping configuration.
* `log_sources` - (Optional, List) Log source configuration.

The `cls_config` object of `log_configuration` supports the following:

* `topic_id` - (Optional, String) ID of the CLS log topic used for shipping sandbox tool logs.

The `log_sources` object of `log_configuration` supports the following:

* `files` - (Optional, Set) Log file paths to collect. They must be files under `/logs/`, sub directories are not supported, at most 10 files.

The `network_configuration` object supports the following:

* `network_mode` - (Required, String) Network mode. Valid values: `PUBLIC`, `VPC`, `SANDBOX`.
* `vpc_config` - (Optional, List) VPC network configuration, used when `network_mode` is `VPC`.

The `vpc_config` object of `network_configuration` supports the following:

* `security_group_ids` - (Optional, Set) List of security group IDs.
* `subnet_ids` - (Optional, Set) List of VPC subnet IDs.

The `storage_mounts` object supports the following:

* `mount_path` - (Optional, String) Local mount path inside the sandbox instance.
* `name` - (Optional, String) Name of the storage mount configuration.
* `read_only` - (Optional, Bool) Whether the storage mount is read-only. Defaults to `false` on the cloud side.
* `storage_source` - (Optional, List) Storage source configuration.
* `volume_template` - (Optional, List) Reference of a volume template together with its derivation rules.
* `volume` - (Optional, List) Reference of a volume.

The `storage_source` object of `storage_mounts` supports the following:

* `agent_bucket` - (Optional, List) AgentBucket storage configuration.
* `cfs` - (Optional, List) CFS file system configuration.
* `cos` - (Optional, List) COS bucket configuration.
* `image` - (Optional, List) Image volume configuration.

The `agent_bucket` object of `storage_source` supports the following:

* `access_domain` - (Optional, String) Access domain passed to the AgentBucket.
* `library_id` - (Optional, String) Library ID passed to the AgentBucket.
* `space_id` - (Optional, String) Space ID passed to the AgentBucket.

The `cfs` object of `storage_source` supports the following:

* `file_system_id` - (Optional, String) ID of the CFS file system.
* `path` - (Optional, String) Mount path of the CFS file system.

The `cos` object of `storage_source` supports the following:

* `bucket_name` - (Optional, String) Name of the object storage bucket.
* `bucket_path` - (Optional, String) Path inside the bucket. It must be an absolute path starting with `/`.
* `endpoint` - (Optional, String) Access domain of the object storage.

The `image` object of `storage_source` supports the following:

* `image_registry_type` - (Optional, String) Image registry type. Valid values: `enterprise`, `personal`.
* `reference` - (Optional, String) Image address.
* `sub_path` - (Optional, String) Path inside the image.

The `volume` object of `storage_mounts` supports the following:

* `volume_id` - (Optional, String) Stable volume id, in the form of `vol-xxxxxxxx`.
* `volume_name` - (Optional, String) Volume name, resolved to a stable volume id by the cloud side under the current AppId.

The `volume_template` object of `storage_mounts` supports the following:

* `reuse_key` - (Optional, String) Reuse key. It can be a literal value or `${metadata.xxx}`.
* `volume_template_name` - (Optional, String) Name of the volume template.

The `tags` object supports the following:

* `key` - (Optional, String) Tag key.
* `value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the sandbox tool, in ISO8601 format.
* `id` - ID of the sandbox tool, same as `tool_id`.
* `status_reason` - Reason describing why the sandbox tool is in its current status.
* `status` - Status of the sandbox tool. Valid values: `CREATING`, `ACTIVE`, `DELETING`, `FAILED`.
* `tool_id` - ID of the sandbox tool assigned by the cloud side. It is the resource identity.
* `update_time` - Last update time of the sandbox tool, in ISO8601 format.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateSandboxTool](https://www.tencentcloud.com/search/?wd=CreateSandboxTool&type=document&lang=en)
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
* Read: [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
* Update:
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
  - [UpdateSandboxTool](https://www.tencentcloud.com/search/?wd=UpdateSandboxTool&type=document&lang=en)
* Delete:
  - [DeleteSandboxTool](https://www.tencentcloud.com/search/?wd=DeleteSandboxTool&type=document&lang=en)
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)

## Import

AGS sandbox tool can be imported using the `tool_id`, e.g.
```
$ terraform import tencentcloud_ags_sandbox_tool.custom sdt-xxxxxxxx
```

