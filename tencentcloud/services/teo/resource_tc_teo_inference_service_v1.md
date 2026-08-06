Provides a resource to create a TEO inference service

Example Usage

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
    scaling_mode = "Auto"
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

Import

TEO inference service can be imported using the composite id (zone_id#service_id), e.g.

```
terraform import tencentcloud_teo_inference_service_v1.example zone-xxxxxxxx#service-xxxxxxxx
```