---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service_v1"
sidebar_current: "docs-tencentcloud-resource-teo_inference_service_v1"
description: |-
  Provides a resource to create a TEO inference service
---

# tencentcloud_teo_inference_service_v1

Provides a resource to create a TEO inference service

## Example Usage

```hcl
resource "tencentcloud_teo_inference_service_v1" "example" {
  zone_id     = "zone-xxxxxxxx"
  name        = "my-inference-service"
  listen_port = 8080
  description = "example inference service"

  containers {
    image_type = "TCR"

    tcr_repository_config {
      tcr_type    = "Personal"
      image       = "ccr.ccs.tencentyun.com/namespace/image:v1.0.0"
      region_name = "ap-guangzhou"
    }

    startup_command = "python app.py"

    environment_variables {
      key   = "MODEL_PATH"
      value = "/models/v1"
    }
  }

  resource_config {
    scaling_mode  = "Auto"
    hardware_spec = "GPU.A10.1"

    auto_scaling_config {
      min_instance_count = 1

      scaling_policies {
        policy_name = "peak-hours"
        policy_type = "ScheduledScaling"

        scheduled_scaling_policy {
          scheduled_actions {
            cron_expression    = "0 9 * * 1-5"
            min_instance_count = 3
          }

          effective_range {
            effective_type = "LongTerm"
          }

          time_zone = "Asia/Shanghai"
        }
      }
    }

    concurrency = 10
  }

  request_paths = ["/predict", "/health"]
}
```

## Argument Reference

The following arguments are supported:

* `containers` - (Required, List) Container configuration for the inference service. Currently, only 1 container is supported.
* `listen_port` - (Required, Int) The port that the inference service listens on. Only supports integers between 1-65535.
* `name` - (Required, String, ForceNew) Inference service name. The name is up to 30 characters, only supports lowercase letters, numbers, hyphens, starts with a letter, ends with a number or letter, and does not support duplicates.
* `resource_config` - (Required, List) Resource configuration for the inference service.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `description` - (Optional, String) Description. Max length: 60 characters.
* `operation` - (Optional, String) Operation to perform on the inference service. Valid values: `Stop`, `Resume`. This field is not persisted to state.
* `request_paths` - (Optional, Set: [`String`]) Request path list for the inference service. Up to 20 paths.

The `auto_scaling_config` object of `resource_config` supports the following:

* `min_instance_count` - (Required, Int) Minimum number of instances.
* `scaling_policies` - (Optional, List) Scaling policy list. Up to 5 policies.

The `containers` object supports the following:

* `image_type` - (Required, String) Image type. Valid values: `TCR` (Tencent Cloud Container Registry).
* `environment_variables` - (Optional, List) Environment variables for the container runtime. Up to 10 variables.
* `startup_command` - (Optional, String) Command executed when the container starts. Defaults to the image's Entrypoint/CMD if not specified. Max length: 1024 characters.
* `tcr_repository_config` - (Optional, List) TCR repository configuration. Required when ImageType is TCR.

The `effective_range` object of `scheduled_scaling_policy` supports the following:

* `effective_type` - (Required, String) Effective type. Valid values: `LongTerm`, `Custom`.
* `end_date` - (Optional, String) End date for the effective range. Required when EffectiveType is Custom and must not be earlier than StartDate.
* `start_date` - (Optional, String) Start date for the effective range. Required when EffectiveType is Custom.

The `environment_variables` object of `containers` supports the following:

* `key` - (Required, String) Variable name. Only uppercase and lowercase letters, numbers, and underscores, must start with a letter or underscore. Max length: 64 characters.
* `value` - (Required, String) Variable value. Any visible characters. Max length: 2048 characters.

The `manual_instance_config` object of `resource_config` supports the following:

* `fixed_instance_count` - (Required, Int) Fixed instance count.

The `resource_config` object supports the following:

* `hardware_spec` - (Required, String) Hardware specification. Note: This field can only be set during creation and cannot be modified afterwards.
* `scaling_mode` - (Required, String) Scaling mode. Valid values: `Auto` (auto scaling based on request volume), `Manual` (fixed instance count).
* `auto_scaling_config` - (Optional, List) Auto scaling configuration. Required when ScalingMode is Auto.
* `concurrency` - (Optional, Int) Concurrency per instance. Default: 1.
* `manual_instance_config` - (Optional, List) Manual instance configuration. Required when ScalingMode is Manual.

The `scaling_policies` object of `auto_scaling_config` supports the following:

* `policy_name` - (Required, String) Policy name. Length: 1-30 characters. Must be unique within the same service.
* `policy_type` - (Required, String) Policy type. Valid values: `ScheduledScaling` (scheduled scaling).
* `scheduled_scaling_policy` - (Optional, List) Scheduled scaling policy configuration. Required when PolicyType is ScheduledScaling.

The `scheduled_actions` object of `scheduled_scaling_policy` supports the following:

* `cron_expression` - (Required, String) Cron expression for triggering the scheduled scaling action. Uses 5-field standard Cron format: minute hour day month weekday.
* `min_instance_count` - (Required, Int) Minimum number of instances to adjust to when this scheduled scaling action is triggered.

The `scheduled_scaling_policy` object of `scaling_policies` supports the following:

* `effective_range` - (Required, List) Effective range for the scheduled scaling policy.
* `scheduled_actions` - (Required, List) Scheduled scaling action list. At least 1, up to 10.
* `time_zone` - (Optional, String) Timezone for the scheduled actions, e.g., UTC, Asia/Shanghai. Defaults to UTC.

The `tcr_repository_config` object of `containers` supports the following:

* `image` - (Required, String) Image address.
* `region_name` - (Required, String) Region name.
* `tcr_type` - (Required, String) TCR service type. Valid values: `Personal`, `Enterprise`.
* `registry_id` - (Optional, String) Registry instance ID. Required when TCRType is Enterprise.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time in ISO date format.
* `inference_url` - Inference access URL.
* `service_id` - Inference service ID.
* `status` - Inference service status. Valid values: `Deploying`, `Running`, `Stopping`, `Stopped`, `Exception`, `Banned`.
* `update_time` - Last modification time in ISO date format.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `3m`) Used when creating the resource.
* `update` - (Defaults to `3m`) Used when updating the resource.
* `delete` - (Defaults to `3m`) Used when deleting the resource.

## Import

TEO inference service can be imported using the composite id (zone_id#service_id), e.g.

```
terraform import tencentcloud_teo_inference_service_v1.example zone-xxxxxxxx#service-xxxxxxxx
```

