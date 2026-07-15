Use this data source to query instances type.

Example Usage

```hcl
data "tencentcloud_instance_types" "example" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 4
  memory_size       = 8
}
```

Complete Example

```hcl
data "tencentcloud_instance_types" "example" {
  cpu_core_count   = 4
  memory_size      = 8
  exclude_sold_out = true

  filter {
    name   = "instance-family"
    values = ["SA2"]
  }

  filter {
    name   = "zone"
    values = ["ap-guangzhou-6"]
  }
}
```

Query with Network and Performance Requirements

```hcl
data "tencentcloud_instance_types" "high_network" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 8
  memory_size       = 16
}

output "instance_details" {
  value = [for instance in data.tencentcloud_instance_types.high_network.instance_types : {
    type             = instance.instance_type
    type_name        = instance.type_name
    network_card     = instance.network_card
    bandwidth        = instance.instance_bandwidth
    pps              = instance.instance_pps
    cpu_type         = instance.cpu_type
    frequency        = instance.frequency
    status_category  = instance.status_category
  }]
}
```

Query GPU Instances

```hcl
data "tencentcloud_instance_types" "gpu_instances" {
  gpu_core_count = 1
  
  filter {
    name   = "zone"
    values = ["ap-guangzhou-6"]
  }
}

output "gpu_details" {
  value = [for instance in data.tencentcloud_instance_types.gpu_instances.instance_types : {
    type       = instance.instance_type
    gpu_count  = instance.gpu_count
    fpga       = instance.fpga
  }]
}
```

Query with Local Disk Support

```hcl
data "tencentcloud_instance_types" "local_disk" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 4
}

output "local_disk_types" {
  value = [for instance in data.tencentcloud_instance_types.local_disk.instance_types : 
    instance.local_disk_type_list if length(instance.local_disk_type_list) > 0
  ]
}
```

Query Price Information

```hcl
data "tencentcloud_instance_types" "with_pricing" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 2
  memory_size       = 4
}

output "pricing_info" {
  value = [for instance in data.tencentcloud_instance_types.with_pricing.instance_types : {
    type  = instance.instance_type
    price = length(instance.price) > 0 ? instance.price[0] : null
  }]
}
```

Query with CBS Pricing Information

```hcl
data "tencentcloud_instance_types" "cbs_pricing" {
  cpu_core_count   = 2
  memory_size      = 2
  exclude_sold_out = true

  filter {
    name   = "instance-family"
    values = ["S6"]
  }

  filter {
    name   = "zone"
    values = ["ap-guangzhou-6"]
  }

  cbs_filter {
    disk_types       = ["CLOUD_SSD"]
    disk_charge_type = "PREPAID"
    disk_usage       = "SYSTEM_DISK"
  }
}

output "cbs_pricing_info" {
  value = [for instance in data.tencentcloud_instance_types.cbs_pricing.instance_types : {
    type        = instance.instance_type
    cbs_configs = [for c in instance.cbs_configs : {
      disk_type              = c.disk_type
      charge_unit            = c.charge_unit
      discount_price         = c.discount_price
      discount_price_high    = c.discount_price_high
      original_price         = c.original_price
      original_price_high    = c.original_price_high
      unit_price             = c.unit_price
      unit_price_discount    = c.unit_price_discount
      unit_price_discount_high = c.unit_price_discount_high
      unit_price_high        = c.unit_price_high
    }]
  }]
}
```