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

### Query all sandbox instances:

```hcl
data "tencentcloud_ags_sandbox_instances" "all" {
}
```

### Query sandbox instances by ID:

```hcl
data "tencentcloud_ags_sandbox_instances" "filtered" {
  instance_ids = ["sandbox-8f0b1c2d", "sandbox-9a1b2c3d"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Optional, List) List of sandbox instance IDs to query. If empty, all sandbox instances are queried. Up to 100 IDs are supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

        * `add` - Capabilities to add. Value reference: ["SYS_ADMIN"].
        * `drop` - Capabilities to remove. Value reference: ["ALL"].
        * `path` - Path.
        * `port` - Port.
        * `scheme` - Protocol.
      * `capabilities` - Capabilities refer to POSIX capabilities that can be added to or removed from a running container.
      * `cpu` - CPU resource amount.
      * `failure_threshold` - Health check failure threshold.
      * `http_get` - HTTP GET probe configuration.
      * `image_id` - Custom WAA image ID.
      * `memory` - Memory resource amount.
      * `name` - Environment variable name.
      * `name` - Port name.
      * `options` - Configuration items (corresponding to resolv.conf options).
      * `port` - Port.
      * `privileged` - Whether to enable privilege. Default value: true.
      * `probe_period_ms` - Health check interval.
      * `probe_timeout_ms` - Health check single probe timeout.
      * `protocol` - Protocol.
      * `ready_timeout_ms` - Health check ready timeout.
      * `searches` - Search domains (corresponding to the search directive of resolv.conf).
      * `servers` - DNS server addresses. Format: valid IP addresses. Default value: 10.0.0.1.
      * `storage` - Custom disk size. Enumerated values: 1Gi: 1Gi; 5Gi: 5Gi; 10Gi: 10Gi; 20Gi: 20Gi. Note: this field may return null, indicating that no valid value is obtained.
      * `success_threshold` - Health check success threshold.
      * `value` - Environment variable value.
      * `version` - Built-in OSWorld version. Enumerated values: osworld1: OSWorld v1; osworld2: OSWorld v2. Default value: osworld1.
    * `args` - Start arguments.
    * `command` - Start command.
    * `dns_config` - Sandbox DNS configuration.
    * `env` - Environment variables.
    * `image_digest` - Image digest.
    * `image_registry_type` - Image registry type: `enterprise`, `personal`, `custom`. Enumerated values: enterprise: TCR enterprise container image service; personal: CCR personal container image service.
    * `image` - Image address.
    * `mount_path` - Local mount path of the sandbox instance (optional), inherits the storage configuration of the tool by default.
    * `name` - Metadata name.
    * `name` - Name of the storage configuration in the sandbox tool.
    * `os_world_configuration` - Built-in OSWorld configuration.
    * `ports` - Port configuration.
    * `probe` - Probe configuration.
    * `read_only` - Read/write permission of the sandbox instance storage mount (optional), inherits the storage configuration of the tool by default.
    * `resources` - Resource configuration.
    * `security_context` - Sandbox security context.
    * `sub_path` - Storage mount sub-path of the sandbox instance (optional).
    * `value` - Metadata value.
    * `waa_configuration` - WAA sandbox tool configuration.
  * `auth_mode` - Sandbox access authentication mode. Enumerated values: DEFAULT: default, i.e. TOKEN authentication; TOKEN: Token authentication, i.e. all port access requires a TOKEN; NONE: no authentication, i.e. all port access does not require a TOKEN; PUBLIC: public mode, i.e. the ENVD management port (49983) access requires a TOKEN while other ports do not. Default value: DEFAULT.
  * `computer_configuration` - Configuration of desktop-computer sandbox instances.
  * `create_time` - Creation time (ISO 8601 format).
  * `custom_configuration` - Custom configuration of the sandbox instance.
  * `expires_at` - Expiration time (ISO 8601 format); null means no expiration time.
  * `instance_id` - Unique identifier of the sandbox instance.
  * `metadata` - Metadata of the sandbox instance.
  * `mount_options` - Storage mount options.
  * `network_mode` - Network mode. Enumerated values: PUBLIC: public network access; SANDBOX: no network; INTERNAL_SERVICE: Tencent Cloud internal public services. It can override the network configuration of the tool. But if a tool itself does not support VPC networks, selecting the VPC mode in the instance settings is ineffective.
  * `persistent` - Whether the instance is persistent.
  * `status` - Instance status: STARTING (starting), RUNNING (running), STOPPING (stopping), STOPPED (stopped), STOP_FAILED (stop failed), FAILED (failed), RESIZING (specification updating).
  * `stop_reason` - Stop reason: manual (manual), timeout (timeout), error (error), system (system). Only present when the status is STOPPED, STOP_FAILED or FAILED. When the provider fails to stop the instance, the status is STOP_FAILED and the reason is error.
  * `timeout_seconds` - Timeout in seconds; null means no timeout is set.
  * `tool_id` - ID of the sandbox tool the instance belongs to.
  * `tool_name` - Name of the sandbox tool the instance belongs to.
  * `update_time` - Update time (ISO 8601 format).
* `id` - ID of the data source, a random token.
* `instance_set` - Sandbox instance list.


