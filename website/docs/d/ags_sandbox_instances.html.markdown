---
subcategory: "Agent Sandbox(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_instances"
sidebar_current: "docs-tencentcloud-data_source-ags_sandbox_instances"
description: |-
  Use this data source to query the list of AGS sandbox instances.
---

# tencentcloud_ags_sandbox_instances

Use this data source to query the list of AGS sandbox instances.

## Example Usage

```hcl
data "tencentcloud_ags_sandbox_instances" "example" {
  instance_ids = ["ags-fake-instance-id"]
}
```

### Query all sandbox instances

```hcl
data "tencentcloud_ags_sandbox_instances" "all" {
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Optional, List) List of sandbox instance ids to query. If it is empty, all instances are queried. At most 100 ids are supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

        * `add` - Capabilities to add. Reference value: `["SYS_ADMIN"]`.
        * `drop` - Capabilities to remove. Reference value: `["ALL"]`.
        * `path` - Path.
        * `port` - Port.
        * `scheme` - Protocol.
      * `capabilities` - POSIX capabilities that can be added to or removed from the running container.
      * `cpu` - Amount of cpu resources.
      * `failure_threshold` - Failure threshold of the health check.
      * `http_get` - HTTP GET probe configuration.
      * `image_id` - Id of the custom WAA image.
      * `memory` - Amount of memory resources.
      * `name` - Name of the environment variable.
      * `name` - Name of the port.
      * `options` - Configuration items, corresponding to the options of resolv.conf.
      * `port` - Port number.
      * `privileged` - Whether privileged mode is enabled. Default value: `true`.
      * `probe_period_ms` - Interval in milliseconds between two health check probes.
      * `probe_timeout_ms` - Timeout in milliseconds of a single health check probe.
      * `protocol` - Protocol.
      * `ready_timeout_ms` - Timeout in milliseconds for the health check to become ready.
      * `searches` - Search domains, corresponding to the `search` directive of resolv.conf.
      * `servers` - DNS server addresses. Valid IP addresses are required. Default value: `10.0.0.1`.
      * `storage` - Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`. Note: this field may return `null`, which means no valid value is obtained.
      * `success_threshold` - Success threshold of the health check.
      * `value` - Value of the environment variable.
      * `version` - Version of the built-in OSWorld. Valid values: `osworld1` (OSWorld v1), `osworld2` (OSWorld v2). Default value: `osworld1`.
    * `args` - Startup arguments.
    * `command` - Startup command.
    * `dns_config` - Sandbox DNS configuration.
    * `env` - Environment variables.
    * `image_digest` - Image digest.
    * `image_registry_type` - Image registry type. Valid values: `enterprise`, `personal`, `custom`. Enumerated values: `enterprise`: TCR enterprise container image service; `personal`: CCR personal container image service.
    * `image` - Image address.
    * `mount_path` - Local mount path inside the sandbox instance (optional). By default it inherits the storage configuration of the tool.
    * `name` - Name of the metadata.
    * `name` - Name of the storage configuration in the sandbox tool.
    * `os_world_configuration` - Configuration of the built-in OSWorld.
    * `ports` - Port configuration.
    * `probe` - Probe configuration.
    * `read_only` - Read and write permission of the sandbox instance storage mount (optional). By default it inherits the storage configuration of the tool.
    * `resources` - Resource configuration.
    * `security_context` - Sandbox security context.
    * `sub_path` - Sub path of the sandbox instance storage mount (optional).
    * `value` - Value of the metadata.
    * `waa_configuration` - Configuration of the WAA sandbox tool.
  * `auth_mode` - Authentication mode of sandbox access. Valid values: `DEFAULT` (default, i.e. TOKEN authentication), `TOKEN` (all port accesses must carry the TOKEN), `NONE` (all port accesses do not need the TOKEN), `PUBLIC` (public mode, i.e. the ENVD management port (49983) access must carry the TOKEN while other ports do not need it). Default value: `DEFAULT`.
  * `computer_configuration` - Configuration of desktop-like sandbox instances.
  * `create_time` - Creation time of the instance in ISO 8601 format.
  * `custom_configuration` - Custom configuration of the sandbox instance.
  * `expires_at` - Expiration time in ISO 8601 format. `null` means no expiration time.
  * `instance_id` - Unique identifier of the sandbox instance.
  * `metadata` - Metadata of the sandbox instance.
  * `mount_options` - Storage mount options.
  * `network_mode` - Network mode. Valid values: `PUBLIC` (public network access), `SANDBOX` (no network), `INTERNAL_SERVICE` (Tencent Cloud internal public services). It can override the network configuration of the tool level. However, if a tool itself does not support VPC networks, selecting the VPC mode in the instance settings has no effect.
  * `persistent` - Whether the instance is persistent.
  * `status` - Status of the instance. Valid values: `STARTING`, `RUNNING`, `STOPPING`, `STOPPED`, `STOP_FAILED`, `FAILED`, `RESIZING`.
  * `stop_reason` - Reason why the instance was stopped. Valid values: `manual`, `timeout`, `error`, `system`. It is only set when the status is `STOPPED`, `STOP_FAILED` or `FAILED`. When the provider stops the instance and the stop fails, the status is `STOP_FAILED` and the reason is `error`.
  * `timeout_seconds` - Timeout in seconds. `null` means no timeout is set.
  * `tool_id` - Id of the sandbox tool the instance belongs to.
  * `tool_name` - Name of the sandbox tool the instance belongs to.
  * `update_time` - Last update time of the instance in ISO 8601 format.
* `id` - Id of the data source, a random token.
* `instance_set` - List of sandbox instances.


