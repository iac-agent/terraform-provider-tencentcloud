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

Query with CBS Disk Config Quota

```hcl
data "tencentcloud_instance_types" "with_cbs_filter" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 4
  memory_size       = 8

  cbs_filter {
    disk_types       = ["CLOUD_SSD"]
    disk_charge_type = "PREPAID"
    disk_usage       = "SYSTEM_DISK"
  }
}
```

Query CBS Disk Config Quota with Top-Level Parameters

```hcl
data "tencentcloud_instance_types" "with_cbs_top_level_params" {
  cpu_core_count = 4

  filter {
    name   = "zone"
    values = ["ap-guangzhou-6"]
  }

  cbs_filter {
    disk_charge_type = "PREPAID"
    disk_usage       = "SYSTEM_DISK"
  }

  # Top-level parameters override cbs_filter and derived values for CBS disk config quota query
  disk_types = ["CLOUD_SSD", "CLOUD_HSSD"]
  zones      = ["ap-guangzhou-6", "ap-guangzhou-3"]
  memory     = 8
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

Query CBS Disk Config with Top-Level Parameters

```hcl
data "tencentcloud_instance_types" "cbs_top_level" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 4
  memory_size       = 8

  inquiry_type      = "INQUIRY_CBS_CONFIG"
  disk_charge_type  = "PREPAID"
  instance_families = ["S5"]
}

output "cbs_available" {
  value = data.tencentcloud_instance_types.cbs_top_level.available
}
```

Query CBS Disk Config with Top-Level Parameters and cbs_filter

```hcl
data "tencentcloud_instance_types" "cbs_combined" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 4
  memory_size       = 8

  inquiry_type = "INQUIRY_CVM_CONFIG"
  disk_charge_type = "POSTPAID_BY_HOUR"

  cbs_filter {
    disk_types    = ["CLOUD_SSD"]
    disk_usage    = "DATA_DISK"
  }
}
```