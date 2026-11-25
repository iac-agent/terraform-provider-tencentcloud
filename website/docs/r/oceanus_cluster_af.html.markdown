---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_cluster_af"
sidebar_current: "docs-tencentcloud-resource-oceanus_cluster_af"
description: |-
  Provides a resource to create and manage an Oceanus cluster.
---

# tencentcloud_oceanus_cluster_af

Provides a resource to create and manage an Oceanus cluster.

## Example Usage

```hcl
resource "tencentcloud_oceanus_cluster_af" "example" {
  cluster_name       = "tf-example-cluster"
  region_id          = 1
  zone_id            = 100001
  login_password     = "tf-example-password"
  default_cos_bucket = "tf-test-bucket"
  cu                 = 12
  remark             = "created by terraform"

  vpc_descriptions {
    vpc_id    = "vpc-12345"
    subnet_id = "subnet-12345"
  }

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) Cluster name. It supports 1-50 English letters, Chinese characters, digits, hyphens `-` or underscores `_`.
* `default_cos_bucket` - (Required, String) Stream Compute uses the COS object storage to save job checkpoints, jar packages, delivered logs, and so on. If no bucket is available in the local region, create one in the COS console. To ensure normal use, refer to the documentation for the lifecycle configuration of the corresponding COS bucket.
* `login_password` - (Required, String) Password for accessing the Flink UI. The user name and password are used to log in to the Flink UI for viewing jobs. The cluster user name is admin.
* `region_id` - (Required, Int) Region ID, which can be queried through DescribeRegions of the region management system, with the Product parameter set to oceanus.
* `vpc_descriptions` - (Required, List) Stream Compute accesses other cloud resources in the same region through VPCs and elastic network interfaces, which occupy a certain number of subnet IPs. Make sure the selected subnet has enough available IPs. If the existing network does not meet your requirements, create a new VPC or subnet in the VPC console.
* `zone_id` - (Required, Int) Availability zone ID, which can be queried through DescribeZones of the region management system, with the Product parameter set to oceanus.
* `cluster_type` - (Optional, String) Cluster type. Valid value: `MULTI_AZ_CLUSTER` (multi-availability-zone cluster).
* `cu_memory` - (Optional, Int) CU-to-memory ratio. Only [0, 2, 4, 8] are supported.
* `cu` - (Optional, Int) Number of cluster CUs. 12 CUs are the minimum computing resource and billing unit of Stream Compute; 1 CU contains 1 CPU and 4 GB of memory. Management node fees are waived when the number of CUs is greater than or equal to 48. The target CU must be greater than 12 and satisfy 12 + 7*n (n >= 0).
* `flink_web_ui_network_access_type` - (Optional, String) Flink UI access type. Valid values: `NetworkAccess_INTERNAL` (intranet access), `NetworkAccess_EXTERNAL` (public network access). Default value: `NetworkAccess_EXTERNAL`.
* `instance_charge_type` - (Optional, String) Cluster billing type. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_SECOND` (pay as you go). Default value: `POSTPAID_BY_SECOND`.
* `period` - (Optional, Int) Purchase duration in months. Valid range: [1, 48].
* `remark` - (Optional, String) Cluster description. It supports 1-50 English letters, Chinese characters, digits, hyphens `-` or underscores `_`.
* `renew_flag` - (Optional, String) Auto-renewal flag. Valid values: `NOTIFY_AND_MANUAL_RENEW` (notify and renew manually), `NOTIFY_AND_AUTO_RENEW` (notify and renew automatically), `DISABLE_NOTIFY_AND_MANUAL_RENEW` (no notification and no manual renewal). Default value: `NOTIFY_AND_MANUAL_RENEW`. When InstanceChargeType is PREPAID, corresponding to a monthly-subscription cluster, the auto-renewal flag must be set; pay-as-you-go clusters do not need it.
* `slave_vpc_descriptions` - (Optional, List) VPCs of the multi-availability-zone cluster.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the create/update/delete operations.

The `slave_vpc_descriptions` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.

The `vpc_descriptions` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_read_id` - Cluster ID read back from DescribeClusters; the same value as `id`.
* `create_time` - Cluster creation time.
* `id` - ID of the resource, in the format of cluster-xxx.
* `status_desc` - Cluster status description.
* `status` - Cluster status: 1 not initialized, 3 initializing, 2 running.
* `update_time` - Time of the last operation on the cluster. Default value: `-`.
* `zone` - Availability zone.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateOceanusCluster](https://www.tencentcloud.com/search/?wd=CreateOceanusCluster&type=document&lang=en)
  - [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
* Read: [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
* Update:
  - [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
  - [ScaleOceanusCluster](https://www.tencentcloud.com/search/?wd=ScaleOceanusCluster&type=document&lang=en)
* Delete:
  - [DeleteOceanusCluster](https://www.tencentcloud.com/search/?wd=DeleteOceanusCluster&type=document&lang=en)
  - [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)

## Import

Oceanus cluster can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_cluster_af.example cluster-12345
```

