---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service_deployment_records_1110_03"
sidebar_current: "docs-tencentcloud-data_source-teo_inference_service_deployment_records_1110_03"
description: |-
  Use this data source to query the deployment records of a TEO inference service.
---

# tencentcloud_teo_inference_service_deployment_records_1110_03

Use this data source to query the deployment records of a TEO inference service.

## Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records_1110_03" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "svc-1a2b3c4d"
  sort_by    = "create-time"
  sort_order = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) Inference service ID.
* `zone_id` - (Required, String) Zone ID.
* `sort_by` - (Optional, String) Field used to sort the deployment records. Valid value: `create-time` (deployment creation time). The cloud default is `create-time`.
* `sort_order` - (Optional, String) Sorting order. Valid values: `asc` (ascending), `desc` (descending). The cloud default is `desc`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

              * `cron_expression` - Cron expression describing when the scheduled action is triggered. Standard five-field cron format, seconds and years are not supported.
              * `effective_type` - Validity type. Valid values: `LongTerm`, `Custom`.
              * `end_date` - End date of the validity range. Set when EffectiveType is `Custom`.
              * `min_instance_count` - Minimum number of instances the inference service is scaled to when this action is triggered.
              * `start_date` - Start date of the validity range. Set when EffectiveType is `Custom`.
            * `effective_range` - Validity range of the scheduled scaling policy.
            * `scheduled_actions` - Scheduled scaling actions.
            * `time_zone` - Time zone (IANA time zone name) used to evaluate the scheduled actions. The cloud default is `UTC`.
          * `policy_name` - Policy name.
          * `policy_type` - Policy type. Valid value: `ScheduledScaling`.
          * `scheduled_scaling_policy` - Scheduled scaling configuration. Required when PolicyType is `ScheduledScaling`.
        * `cpu_num` - Number of CPU cores allocated to a single instance.
        * `disk_size` - Temporary disk size allocated to a single instance, in MB.
        * `fixed_instance_count` - Fixed number of instances.
        * `gpu_num` - Number of GPU cards allocated to a single instance.
        * `header_name` - Name of the request header carrying the session ID.
        * `image` - Image address.
        * `key` - Variable name.
        * `mem_size` - Memory size allocated to a single instance, in MB.
        * `min_instance_count` - Minimum number of instances.
        * `region_name` - Region name.
        * `registry_id` - Image registry instance ID. Required when TCRType is `Enterprise`.
        * `scaling_policies` - Scaling policy list.
        * `source` - Where the session ID parameter is passed. Valid value: `Header`.
        * `tcr_type` - TCR service type. Valid values: `Personal`, `Enterprise`.
        * `value` - Variable value.
      * `affinity_mode` - Affinity mode. Valid value: `SessionId`.
      * `auto_scaling_config` - Auto scaling configuration. Required when ScalingMode is `Auto`.
      * `concurrency` - Concurrency of a single instance.
      * `environment_variables` - Environment variables of the container runtime.
      * `hardware_config` - Hardware configuration of the inference service.
      * `hardware_spec_id` - Unique identifier of the hardware specification.
      * `hardware_spec` - Hardware specification identifier. Deprecated, use `hardware_spec_id` instead.
      * `image_type` - Image type. Valid value: `TCR`.
      * `manual_instance_config` - Manual instance configuration. Required when ScalingMode is `Manual`.
      * `scaling_mode` - Scaling mode. Valid values: `Auto`, `Manual`.
      * `session_id_affinity_config` - Affinity configuration used when AffinityMode is `SessionId`.
      * `startup_command` - Command executed when the container starts. When it is not set, the Entrypoint/CMD of the image is used.
      * `switch` - Master switch of the inference service affinity. Valid values: `On`, `Off`.
      * `tcr_repository_config` - TCR image repository information.
    * `affinity_config` - Affinity configuration of the inference service.
    * `containers` - Container configuration of the inference service.
    * `listen_port` - Port the model service listens on.
    * `request_paths` - Request paths of the inference service.
    * `resource_config` - Resource configuration of the inference service.
  * `active_status` - Whether this deployment configuration is the currently effective one. Valid values: `active`, `inactive`.
  * `create_time` - Deployment start time in ISO date format.
  * `duration` - Deployment duration, in seconds.
  * `inference_service_config` - Configuration of this inference service deployment.
  * `operation` - Deployment operation type. Valid values: `create`, `update`, `resume`, `stop`.
  * `record_id` - Deployment record ID.
  * `status` - Deployment status. Valid values: `processing`, `succeeded`, `failed`.
* `record_set` - Deployment records of the inference service.
* `total_count` - Total number of deployment records.


