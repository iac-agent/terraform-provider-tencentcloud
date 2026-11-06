Provides a resource to create and manage an AGS (Agent Sandbox) sandbox tool.

A sandbox tool is the runtime template used to create AGS sandbox instances: it defines the tool type,
the network configuration, the storage mounts and, for custom tools, the container image and its runtime
configuration.

~> **NOTE:** Deleting this resource only removes the sandbox tool definition. The COS buckets, CFS file
systems, agent buckets and volumes referenced by `storage_mounts` are never deleted by this resource;
clean them up separately if they are no longer needed.

~> **NOTE:** The following attributes can only be set on creation. Changing them forces a new resource
because `UpdateSandboxTool` has no corresponding parameter: `tool_name`, `tool_type`, `network_configuration`,
`default_timeout`, `role_arn`, `client_token`, `storage_mounts`, `log_configuration` and `persistent`.

~> **NOTE:** `UpdateSandboxTool` does not document whether an omitted argument keeps or clears the previous
value, so removing `description` from the configuration is refused with an error instead of guessing.

~> **NOTE:** Drift detection is partial, because the describe API does not return everything that was sent:
`client_token` is never returned, and `default_timeout` is only returned as a number of seconds while the
input is a duration string. Values managed by the configuration are therefore kept as configured and are
only replaced by the cloud value when the attribute is not managed (for example right after an import).
Changes made outside of Terraform to those attributes are not detected.

Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "custom" {
  tool_name = "tf-example-tool"
  tool_type = "custom"

  network_configuration {
    network_mode = "VPC"

    vpc_config {
      subnet_ids         = ["subnet-12345678"]
      security_group_ids = ["sg-12345678"]
    }
  }

  description     = "sandbox tool managed by terraform"
  default_timeout = "10m"
  persistent      = true

  tags {
    key   = "env"
    value = "test"
  }

  storage_mounts {
    name       = "workspace"
    mount_path = "/mnt/workspace"
    read_only  = false

    storage_source {
      cos {
        endpoint    = "cos.ap-guangzhou.myqcloud.com"
        bucket_name = "bucket-1250000000"
        bucket_path = "/data"
      }
    }
  }

  custom_configuration {
    image               = "ccr.ccs.tencentcloudcr.com/tf-example/demo:1.0"
    image_registry_type = "personal"
    command             = ["/bin/sh"]
    args                = ["-c", "python server.py"]

    env {
      name  = "LOG_LEVEL"
      value = "info"
    }

    ports {
      name     = "http"
      port     = 8080
      protocol = "TCP"
    }

    resources {
      cpu     = "1"
      memory  = "2Gi"
      storage = "10Gi"
    }

    probe {
      http_get {
        path   = "/healthz"
        port   = 8080
        scheme = "HTTP"
      }

      ready_timeout_ms  = 30000
      probe_timeout_ms  = 1000
      probe_period_ms   = 5000
      success_threshold = 1
      failure_threshold = 3
    }

    dns_config {
      servers = ["10.0.0.1"]
    }

    security_context {
      privileged = true

      capabilities {
        add  = ["SYS_ADMIN"]
        drop = ["ALL"]
      }
    }
  }

  log_configuration {
    cls_config {
      topic_id = "a1b2c3d4-0000-1111-2222-333344445555"
    }

    log_sources {
      files = ["/logs/app.log"]
    }
  }
}
```

Example Usage - public network tool with a desktop computer environment

```hcl
resource "tencentcloud_ags_sandbox_tool" "computer" {
  tool_name = "tf-example-computer"
  tool_type = "osworld"

  network_configuration {
    network_mode = "PUBLIC"
  }

  computer_configuration {
    os_world_configuration {
      version = "osworld2"
    }
  }
}
```

Import

AGS sandbox tool can be imported using the `tool_id`, e.g.
```
$ terraform import tencentcloud_ags_sandbox_tool.custom sdt-xxxxxxxx
```