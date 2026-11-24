Provides a resource to start an AGS sandbox instance.

-> **Note:** Sandbox instances have no update API: every configurable attribute is create-only. Changing any of them destroys the existing instance and starts a new one.

-> **Note:** `timeout`, `source_instance_id` and `snapshot_id` are only accepted by the start API and are not reported back by the describe API, so drift on them cannot be detected.

-> **Note:** Destroying this resource stops the sandbox instance. The instance is also recycled automatically once `timeout` expires.

Example Usage

```hcl

resource "tencentcloud_ags_sandbox_instance" "example" {
  tool_name = "tf-example"
  timeout   = "10m"
  mount_options = [
    {
      name       = "workspace"
      mount_path = "/workspace"
      sub_path   = null
      read_only  = false
    },
  ]
  custom_configuration = {
    args                = ["-port", "8080"]
    command             = ["/bin/go-httpbin"]
    dns_config          = null
    env                 = null
    image               = "ccr.ccs.tencentyun.com/test/go-httpbin:v2.16.1"
    image_registry_type = "personal"
    ports = [
      {
        name     = "http"
        port     = 8080
        protocol = "TCP"
      },
    ]
    probe = {
      failure_threshold = 60
      http_get = {
        path   = "/status/200"
        port   = 8080
        scheme = "HTTP"
      }
      probe_period_ms   = 1000
      probe_timeout_ms  = 1000
      ready_timeout_ms  = 30000
      success_threshold = 1
    }
    resources = {
      cpu     = "1"
      memory  = "1Gi"
      storage = null
    }
    security_context = null
  }
  auth_mode = null
  metadata = [
    {
      name  = "owner"
      value = "tf"
    },
  ]
}
```

Import

AGS sandbox instance can be imported using the sandbox instance id, e.g.

```
$ terraform import tencentcloud_ags_sandbox_instance.example sbi-xxxxxxxx
```
