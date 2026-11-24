Provides a resource to start an AGS sandbox instance.

-> **Note:** All configurable attributes are create-only: `StartSandboxInstance` has no matching update API, so changing any of them destroys the existing sandbox instance and starts a new one.

-> **Note:** `timeout` and `client_token` are only accepted by the create API. The describe API reports the timeout in seconds (`timeout_seconds`) instead of the original duration string and never echoes back a client token, so those values are kept as they are and drift on them cannot be detected.

-> **Note:** Destroying this resource stops the sandbox instance. The cloud side reclaims the instance once it is stopped, and instances with a `timeout` are reclaimed automatically when the timeout expires.

Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "tf-example"
  tool_type = "custom"
  network_configuration = {
    network_mode = "SANDBOX"
    vpc_config   = null
  }
}

resource "tencentcloud_ags_sandbox_instance" "example" {
  tool_id = tencentcloud_ags_sandbox_tool.example.id
  timeout = "10m"
  auth_mode = "TOKEN"

  mount_options = [
    {
      name       = "workspace"
      mount_path = "/workspace"
      sub_path   = "/sub"
      read_only  = true
    },
  ]

  custom_configuration = {
    image               = "ccr.ccs.tencentyun.com/test/go-httpbin:v2.16.1"
    image_registry_type = "personal"
    command             = ["/bin/go-httpbin"]
    args                = ["-port", "8080"]
    env                 = null
    ports = [
      {
        name     = "http"
        port     = 8080
        protocol = "TCP"
      },
    ]
    resources = {
      cpu     = "1"
      memory  = "2Gi"
      storage = "5Gi"
    }
    probe             = null
    dns_config        = null
    security_context  = null
  }

  metadata = [
    {
      name  = "k1"
      value = "v1"
    },
  ]
}
```

Import

Sandbox instances can be imported using the `instance_id`, e.g.

```
$ terraform import tencentcloud_ags_sandbox_instance.example inst-12345678
```
