# resource_tc_oceanus_oceanus_cluster

 Provides a resource to create an Oceanus cluster.

## Example Usage

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name        = "tf-example-cluster"
  region_id           = 1
  zone_id             = 100003
  login_password      = "YourPassword123"
  default_cos_bucket  = "example-bucket-1250000000"
  cu                  = 12
  cu_memory           = 4
  instance_charge_type = "POSTPAID_BY_SECOND"

  vpc_descriptions {
    vpc_id    = "vpc-xxx"
    subnet_id = "subnet-xxx"
  }
}
```

Multi-AZ cluster example:

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name        = "tf-example-multiaz-cluster"
  region_id           = 1
  zone_id             = 100003
  login_password      = "YourPassword123"
  default_cos_bucket  = "example-bucket-1250000000"
  cluster_type        = "MULTI_AZ_CLUSTER"
  instance_charge_type = "POSTPAID_BY_SECOND"

  vpc_descriptions {
    vpc_id    = "vpc-xxx"
    subnet_id = "subnet-xxx"
  }

  slave_vpc_descriptions {
    vpc_id    = "vpc-yyy"
    subnet_id = "subnet-yyy"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) Cluster name. Supports 1-50 english letters, chinese characters, numbers, hyphens or underscores.
* `region_id` - (Required, ForceNew) Region ID, which can be queried via DescribeRegions with product oceanus.
* `zone_id` - (Required, ForceNew) Zone ID, which can be queried via DescribeZones with product oceanus.
* `login_password` - (Required) Flink UI access password. The cluster user name is admin.
* `vpc_descriptions` - (Required, ForceNew) VPC descriptions. Stream compute accesses other cloud resources in the same region via VPC and ENI and occupies subnet IPs.
  * `vpc_id` - (Required) VPC ID.
  * `subnet_id` - (Required) Subnet ID.
  * `app_id` - (Optional) User AppId.
  * `owner_uin` - (Optional) User UIN.
* `default_cos_bucket` - (Required, ForceNew) Default COS bucket used by stream compute to save checkpoints, jars or logs.
* `cu` - (Optional) Cluster CU number. 12 CU is the minimum compute resource and billing unit, 1 CU contains 1 CPU and 4GB memory. Management node fee is waived when CU >= 48.
* `remark` - (Optional) Cluster description. Supports 1-50 english letters, chinese characters, numbers, hyphens or underscores.
* `period` - (Optional) Purchase period in months, range [1, 48]. Only valid for PREPAID clusters.
* `instance_charge_type` - (Optional, ForceNew) Cluster charge type. Valid values: PREPAID, POSTPAID_BY_SECOND. Default: POSTPAID_BY_SECOND.
* `cluster_type` - (Optional, ForceNew) Cluster type. Valid value: MULTI_AZ_CLUSTER.
* `renew_flag` - (Optional) Auto renewal flag. Valid values: NOTIFY_AND_MANUAL_RENEW, NOTIFY_AND_AUTO_RENEW, DISABLE_NOTIFY_AND_MANUAL_RENEW. Default: NOTIFY_AND_MANUAL_RENEW. Only valid for PREPAID clusters.
* `flink_web_ui_network_access_type` - (Optional, ForceNew) Flink UI network access type. Valid values: NetworkAccess_INTERNAL, NetworkAccess_EXTERNAL. Default: NetworkAccess_EXTERNAL.
* `slave_vpc_descriptions` - (Optional, ForceNew) Multi-AZ VPC descriptions. Only valid for MULTI_AZ_CLUSTER.
  * `vpc_id` - (Required) VPC ID.
  * `subnet_id` - (Required) Subnet ID.
  * `app_id` - (Optional) User AppId.
  * `owner_uin` - (Optional) User UIN.
* `cu_memory` - (Optional, ForceNew) CU memory ratio, only supports [0, 2, 4, 8].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the cluster.
* `status` - Cluster status. 1: uninitialized, 2: running, 3: initializing.
* `create_time` - Cluster create time.
* `update_time` - Cluster last operation time.

## Import

Oceanus cluster can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_oceanus_cluster.example cluster-xxx
```
