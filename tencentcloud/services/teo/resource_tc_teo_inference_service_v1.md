Provides a resource to create a TEO inference service.

Example Usage

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

Import

teo inference service can be imported using the id, e.g.

```
terraform import tencentcloud_teo_inference_service_v1.example zone-xxxxxx#svc-xxxxxx
```