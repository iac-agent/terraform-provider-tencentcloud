---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_chc_network_mode"
sidebar_current: "docs-tencentcloud-resource-cvm_chc_network_mode"
description: |-
  Provides a resource to switch CVM CHC physical server network mode.
---

# tencentcloud_cvm_chc_network_mode

Provides a resource to switch CVM CHC physical server network mode.

## Example Usage

```hcl
resource "tencentcloud_cvm_chc_network_mode" "example" {
  chc_ids      = ["chc-1a2b3c4d"]
  network_mode = "DEPLOY"
}
```

## Argument Reference

The following arguments are supported:

* `chc_ids` - (Required, List: [`String`], ForceNew) CHC physical server ID list.
* `network_mode` - (Required, String) Network mode to switch to. Valid values: DEPLOY (deploy network mode), BUSINESS (business network mode).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



