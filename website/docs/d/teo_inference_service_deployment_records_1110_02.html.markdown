---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service_deployment_records_1110_02"
sidebar_current: "docs-tencentcloud-data_source-teo_inference_service_deployment_records_1110_02"
description: |-
  Use this data source to query the deployment history records of a TEO inference service.
---

# tencentcloud_teo_inference_service_deployment_records_1110_02

Use this data source to query the deployment history records of a TEO inference service.

## Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records_1110_02" "deployment_records" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "svc-2qtuhspyeq8r6"
  sort_by    = "create-time"
  sort_order = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) Inference service ID.
* `zone_id` - (Required, String) Zone ID.
* `sort_by` - (Optional, String) Field used for sorting. Valid values: create-time. Default value: create-time.
* `sort_order` - (Optional, String) Sorting order. Valid values: asc, desc. Default value: desc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

              * `cron_expression` - Cron expression describing when the scheduled scaling action is triggered. 5-field standard Cron format without seconds and years.
              * `effective_type` - Effective range type. Valid values: LongTerm (long-term effective), Custom (custom start and end dates).
              * `end_date` - End date of the effective range. Required when EffectiveType is Custom and must not be earlier than StartDate.
              * `min_instance_count` - Minimum instance count that the inference service is scaled to when this scheduled scaling action is triggered.
              * `start_date` - Start date of the effective range. Required when EffectiveType is Custom.
            * `effective_range` - Effective range of the scheduled scaling policy.
            * `scheduled_actions` - Scheduled scaling action list. At least 1 and at most 10 actions are supported.
            * `time_zone` - Time zone in IANA time zone format used to interpret scheduled actions, for example UTC or Asia/Shanghai. UTC is used when it is not set.
          * `policy_name` - Policy name. 1 to 30 characters. Policy names must be unique within the same service.
          * `policy_type` - Policy type, which cannot be modified after creation. Valid values: ScheduledScaling.
          * `scheduled_scaling_policy` - Scheduled scaling configuration. Required when PolicyType is ScheduledScaling.
        * `cpu_num` - Number of CPU cores allocated to a single instance. Only integer values are supported currently.
        * `disk_size` - Temporary disk size allocated to a single instance, in MB. It must be a multiple of 1024 when specified.
        * `fixed_instance_count` - Fixed instance count.
        * `gpu_num` - Number of GPUs allocated to a single instance. Only integer values are supported currently and the value must be one of AllowedGPUNums of the selected HardwareSpecId.
        * `header_name` - Request header name used to pass the session ID. 1 to 64 characters, only letters, digits and hyphens are allowed. Default value: EO-Infer-Session-Id.
        * `image` - Image address.
        * `key` - Variable name. Only letters, digits and underscores are allowed, and it must start with a letter or an underscore. Maximum 64 characters.
        * `mem_size` - Memory size allocated to a single instance, in MB. It must be a multiple of 1024 when specified.
        * `min_instance_count` - Minimum instance count. It does not take effect when a scaling policy is configured and within its effective period.
        * `region_name` - Region name.
        * `registry_id` - Image repository instance ID. Required when TCRType is Enterprise.
        * `scaling_policies` - Scaling policy list. Up to 5 policies are supported.
        * `source` - Location where the session ID parameter is passed. Valid values: Header. Default value: Header.
        * `tcr_type` - TCR service type. Valid values: Personal, Enterprise.
        * `value` - Variable value. Any visible characters are supported. Maximum 2048 characters.
      * `affinity_mode` - Affinity mode of the inference service. Valid values: SessionId. Default value: SessionId.
      * `auto_scaling_config` - Auto scaling configuration of the inference service. Required when ScalingMode is Auto.
      * `concurrency` - Concurrency of a single instance. Default value: 1.
      * `environment_variables` - Environment variables of the container at runtime. Up to 50 variables are supported.
      * `hardware_config` - Hardware configuration of the inference service. When it is not specified, the default values of the selected HardwareSpecId are used.
      * `hardware_spec_id` - Unique hardware specification ID. It can be obtained from the DescribeInferenceHardwareSpecifications API.
      * `hardware_spec` - Hardware specification identifier. Deprecated, use HardwareSpecId instead.
      * `image_type` - Image type. Valid values: TCR.
      * `manual_instance_config` - Manually configured instance settings of the inference service. Required when ScalingMode is Manual.
      * `scaling_mode` - Scaling mode. Valid values: Auto (adjust instance count automatically based on request volume), Manual (fixed instance count).
      * `session_id_affinity_config` - Affinity configuration based on the session ID. Required when AffinityMode is SessionId.
      * `startup_command` - Command executed when the container starts. The image Entrypoint/CMD is used when it is not set. Maximum 1024 characters.
      * `switch` - Master switch of the inference service affinity. Valid values: On, Off.
      * `tcr_repository_config` - TCR image repository information. Required when ImageType is TCR.
    * `affinity_config` - Affinity configuration of the inference service.
    * `containers` - Container configurations of the inference service.
    * `listen_port` - Port that the model service listens on.
    * `request_paths` - Request path list of the inference service.
    * `resource_config` - Resource configuration of the inference service.
  * `active_status` - Whether this deployment configuration is the currently effective one. Valid values: active (currently effective), inactive (historical or abnormal version).
  * `create_time` - Deployment start time in ISO date format.
  * `duration` - Deployment duration in seconds.
  * `inference_service_config` - Configuration used by this inference service deployment.
  * `operation` - Deployment operation type. Valid values: create, update, resume, stop.
  * `record_id` - Deployment record ID.
  * `status` - Deployment status. Valid values: processing, succeeded, failed.
* `record_set` - Deployment history record list of the inference service.


