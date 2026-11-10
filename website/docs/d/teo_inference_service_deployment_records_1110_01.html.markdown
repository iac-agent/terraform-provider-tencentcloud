---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service_deployment_records_1110_01"
sidebar_current: "docs-tencentcloud-data_source-teo_inference_service_deployment_records_1110_01"
description: |-
  Use this data source to query the TEO inference service deployment record list.
---

# tencentcloud_teo_inference_service_deployment_records_1110_01

Use this data source to query the TEO inference service deployment record list.

## Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records_1110_01" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "isvc-2qtuhspy7cr6"
  sort_by    = "create-time"
  sort_order = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) Inference service ID.
* `zone_id` - (Required, String) Zone ID.
* `sort_by` - (Optional, String) Sort field. Valid values: `create-time` (deployment creation time). Default value: `create-time`.
* `sort_order` - (Optional, String) Sort order. Valid values: `asc` (ascending), `desc` (descending). Default value: `desc`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

              * `cron_expression` - Cron expression.
              * `effective_type` - Effective type. Valid values: `LongTerm`, `Custom`.
              * `end_date` - End date.
              * `min_instance_count` - Minimum instance count.
              * `start_date` - Start date.
            * `effective_range` - Effective range.
            * `scheduled_actions` - Scheduled action list.
            * `time_zone` - Time zone.
          * `policy_name` - Policy name.
          * `policy_type` - Policy type. Valid values: `ScheduledScaling`.
          * `scheduled_scaling_policy` - Scheduled scaling configuration.
        * `cpu_num` - Number of CPU cores allocated to a single instance of the inference service.
        * `disk_size` - Temporary disk size allocated to a single instance of the inference service, in MB.
        * `fixed_instance_count` - Fixed instance count.
        * `gpu_num` - Number of GPU cards allocated to a single instance of the inference service.
        * `header_name` - Header name that carries the session ID.
        * `image` - Image address.
        * `key` - Variable name.
        * `mem_size` - Memory size allocated to a single instance of the inference service, in MB.
        * `min_instance_count` - Minimum instance count.
        * `region_name` - Region name.
        * `registry_id` - TCR registry instance ID.
        * `scaling_policies` - Scaling policy list.
        * `source` - Where the session ID is passed. Valid values: `Header`.
        * `tcr_type` - TCR service type. Valid values: `Personal`, `Enterprise`.
        * `value` - Variable value.
      * `affinity_mode` - Affinity mode. Valid values: `SessionId`.
      * `auto_scaling_config` - Auto scaling configuration of the inference service.
      * `concurrency` - Concurrency of a single instance.
      * `environment_variables` - Environment variables of the container.
      * `hardware_config` - Hardware configuration of the inference service.
      * `hardware_spec_id` - Hardware specification unique ID, which can be obtained from DescribeInferenceHardwareSpecifications.
      * `hardware_spec` - Hardware specification identifier. Deprecated, use `hardware_spec_id` instead.
      * `image_type` - Image type. Valid values: `TCR`.
      * `manual_instance_config` - Manual instance configuration of the inference service.
      * `scaling_mode` - Scaling mode. Valid values: `Auto`, `Manual`.
      * `session_id_affinity_config` - Session ID affinity configuration.
      * `startup_command` - Command executed when the container starts. When it is left empty, the image Entrypoint/CMD is used.
      * `switch` - Affinity switch. Valid values: `On`, `Off`.
      * `tcr_repository_config` - TCR repository information.
    * `affinity_config` - Affinity configuration of the inference service.
    * `containers` - Container configuration of the inference service.
    * `listen_port` - Port that the model service needs to listen on.
    * `request_paths` - Request path list of the inference service.
    * `resource_config` - Resource configuration of the inference service.
  * `active_status` - Whether this deployment configuration is the currently active one. Valid values: `active`, `inactive`.
  * `create_time` - Deployment initiation time in ISO date format.
  * `duration` - Deployment duration, in seconds.
  * `inference_service_config` - Configuration of this inference service deployment.
  * `operation` - Deployment operation type. Valid values: `create`, `update`, `resume`, `stop`.
  * `record_id` - Deployment record ID.
  * `status` - Deployment status. Valid values: `processing`, `succeeded`, `failed`.
* `record_set` - List of inference service deployment records.


