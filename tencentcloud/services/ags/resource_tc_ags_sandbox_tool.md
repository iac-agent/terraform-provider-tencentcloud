Provides a resource to create a sandbox tool of AGS (Agent Development Platform)

~> **Note:** The `tool_name`, `tool_type`, `default_timeout`, `client_token`, `role_arn`, `storage_mounts`, `log_configuration` and `persistent` attributes are immutable: the `UpdateSandboxTool` API does not accept them, so changing any of them forces resource replacement.

Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "my-tool"
  tool_type = "browser"

  network_configuration {
    network_mode = "PUBLIC"
  }

  description = "sandbox tool example"
  persistent  = false

  tags {
    key   = "env"
    value = "prod"
  }
}

resource "tencentcloud_ags_sandbox_tool" "vpc_custom" {
  tool_name = "my-custom-tool"
  tool_type = "custom"

  network_configuration {
    network_mode = "VPC"

    vpc_config {
      subnet_ids         = ["subnet-xxxx"]
      security_group_ids = ["sg-xxxx"]
    }
  }

  custom_configuration {
    image               = "registry.example.com/my-image:1.0"
    image_registry_type = "personal"
    command             = ["/bin/sh"]
    args                = ["-c", "run"]

    env {
      name  = "FOO"
      value = "bar"
    }

    ports {
      name     = "http"
      port     = 8080
      protocol = "TCP"
    }

    resources {
      cpu    = "2"
      memory = "4Gi"
    }
  }
}
```

Import

AGS sandbox tools can be imported using the tool id, e.g.

```
terraform import tencentcloud_ags_sandbox_tool.example tool-1234567890abcdef
```
