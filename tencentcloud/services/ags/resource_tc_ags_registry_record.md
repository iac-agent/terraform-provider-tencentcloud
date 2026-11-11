---
subcategory: ""
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_registry_record"
description: |-
  Provides a resource to create a record in an AGS (Agent Governance Service) registry.
---

# tencentcloud_ags_registry_record

Provides a resource to create a record in an AGS (Agent Governance Service) registry.

## Example Usage

```hcl
resource "tencentcloud_ags_registry_record" "example" {
  registry_id     = "reg-abcdef01"
  name            = "my-record"
  descriptor_type = "MCP"

  description  = "created by terraform"
  version_name = "v1"

  mcp_source {
    type        = "MANUAL"
    descriptors = jsonencode({
      name    = "mcp-server"
      version = "1.0.0"
    })
  }
}
```

## Import

ags registry record can be imported using the id, e.g.

```shell
terraform import tencentcloud_ags_registry_record.example reg-abcdef01#rec-12345678
```
