---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service_v1"
sidebar_current: "docs-tencentcloud-resource-teo_inference_service_v1"
description: |-
  Provides a resource to create a TEO inference service.
---

# tencentcloud_teo_inference_service_v1

Provides a resource to create a TEO inference service.

## Example Usage

```hcl
resource "tencentcloud_teo_inference_service_v1" "example" {
  zone_id     = "zone-xxxxxx"
  name        = "test-inference-svc"
  listen_port = 8080

  containers {
    image_type = "TCR"
    tcr_config {
      tcr_type    = "Personal"
      image       = "ccr.ccs.tencentyun.com/test/image:v1"
      region_name = "ap-guangzhou"
    }
    environment_variables {
      key   = "MODEL_PATH"
      value = "/models/v1"
    }
  }

  resource_config {
    scaling_mode = "Manual"
    manual_instance_config {
      fixed_instance_count = 1
    }
  }

  affinity_config {
    switch        = "On"
    affinity_mode = "SessionId"
    source        = "Header"
    header_name   = "X-Custom-Session-Id"
  }

  request_paths = ["/predict"]
  description   = "test inference service"
}
```

## Argument Reference

The following arguments are supported:

* `containers` - (Required, List) Container configuration of the inference service. Currently only supports setting 1 container.
* `listen_port` - (Required, Int) Port that the model service needs to listen to. Only integers between 1-65535 are supported.
* `name` - (Required, String, ForceNew) Inference service name. The length is limited to 30 characters. Only lowercase letters, numbers, and hyphens are supported. It must start with a letter and end with a number or letter. Duplicates are not supported.
* `resource_config` - (Required, List) Resource configuration of the inference service.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `affinity_config` - (Optional, List) Inference service affinity configuration.
* `description` - (Optional, String) Description. Length limit is 60 characters.
* `request_paths` - (Optional, Set: [`String`]) Request path list of the inference service. Up to 20 paths.

The `affinity_config` object supports the following:

* `affinity_mode` - (Optional, String) Inference service affinity mode. Valid values: `SessionId`. Default value: `SessionId`.
* `header_name` - (Optional, String) Request header name for passing the session ID. Required when Source is Header. Length limit: 1-64 characters. Only letters, numbers, and hyphens are supported. Default value: `EO-Infer-Session-Id`.
* `source` - (Optional, String) The location where the session ID parameter is passed. Valid values: `Header`. Default value: `Header`.
* `switch` - (Optional, String) Inference service affinity switch. Valid values: `On`, `Off`.

The `auto_scaling_config` object of `resource_config` supports the following:

* `min_instance_count` - (Optional, Int) Minimum number of instances. Will not take effect when scaling policy is configured and within its validity period.
* `scaling_policies` - (Optional, List) Scaling policy list. Up to 5 policies.

The `containers` object supports the following:

* `image_type` - (Required, String) Image type. Valid values: `TCR` (Tencent Cloud Container Registry image).
* `environment_variables` - (Optional, List) Environment variables for the container runtime. Up to 10 variables.
* `startup_command` - (Optional, String) Command executed when the container starts. Defaults to the image's Entrypoint/CMD if not specified. Up to 1024 characters.
* `tcr_config` - (Optional, List) TCR image repository configuration. Required when ImageType is TCR.

The `effective_range` object of `scheduled_scaling_policy` supports the following:

* `effective_type` - (Required, String) Effective type. Valid values: `LongTerm` (long-term), `Custom` (custom start/end dates).
* `end_date` - (Optional, String) Effective end date. Required when EffectiveType is Custom and must not be earlier than StartDate.
* `start_date` - (Optional, String) Effective start date. Required when EffectiveType is Custom.

The `environment_variables` object of `containers` supports the following:

* `key` - (Required, String) Variable name. Only uppercase and lowercase letters, numbers, and underscores are allowed. It must start with a letter or underscore. Length limit is 64 characters.
* `value` - (Required, String) Variable value. Supports any visible characters such as letters, numbers, symbols, etc. Length limit is 2048 characters.

The `manual_instance_config` object of `resource_config` supports the following:

* `fixed_instance_count` - (Required, Int) Fixed instance count.

The `resource_config` object supports the following:

* `scaling_mode` - (Required, String) Scaling mode. Valid values: `Auto` (auto-scale based on request volume), `Manual` (manually set fixed instance count).
* `auto_scaling_config` - (Optional, List) Auto scaling configuration. Required when ScalingMode is Auto.
* `concurrency` - (Optional, Int) Concurrency per instance. Default value is 1.
* `hardware_spec` - (Optional, String) Hardware specification.
* `manual_instance_config` - (Optional, List) Manual instance configuration. Required when ScalingMode is Manual.

The `scaling_policies` object of `auto_scaling_config` supports the following:

* `policy_name` - (Required, String) Policy name. Length limit is 1-30 characters. Must be unique within the same service.
* `policy_type` - (Required, String) Policy type. Cannot be modified after creation. Valid values: `ScheduledScaling`.
* `scheduled_scaling_policy` - (Required, List) Scheduled scaling policy configuration. Required when PolicyType is ScheduledScaling.

The `scheduled_actions` object of `scheduled_scaling_policy` supports the following:

* `cron_expression` - (Required, String) Cron expression describing the trigger time. Uses 5-field standard Cron format: minute hour day month weekday.
* `min_instance_count` - (Required, Int) Minimum number of instances to adjust to after hitting this action.

The `scheduled_scaling_policy` object of `scaling_policies` supports the following:

* `effective_range` - (Required, List) Effective range of the scheduled scaling policy.
* `scheduled_actions` - (Required, List) Scheduled scaling action list. At least 1, up to 10.
* `time_zone` - (Optional, String) Timezone, using IANA timezone identifiers. Defaults to UTC.

The `tcr_config` object of `containers` supports the following:

* `image` - (Required, String) Image address.
* `tcr_type` - (Required, String) TCR service type. Valid values: `Personal`, `Enterprise`.
* `region_name` - (Optional, String) Region name.
* `registry_id` - (Optional, String) Image repository instance ID. Required when TCRType is Enterprise.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time in ISO date format.
* `current_instance_count` - Current number of running instances.
* `inference_url` - Inference access URL, through which the underlying model can be accessed for inference.
* `scaling_status` - Scaling status. Valid values: `Normal`, `ScalingOut`, `ScalingIn`.
* `service_id` - Inference service ID.
* `status` - Inference service status. Valid values: `Deploying`, `Running`, `Stopping`, `Stopped`, `Exception`, `Banned`.
* `update_time` - Last modification time in ISO date format.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

teo inference service can be imported using the id, e.g.

```
terraform import tencentcloud_teo_inference_service_v1.example zone-xxxxxx#svc-xxxxxx
```

