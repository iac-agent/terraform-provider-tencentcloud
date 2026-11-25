---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_oceanus_cluster"
sidebar_current: "docs-tencentcloud-resource-oceanus_oceanus_cluster"
description: |-
  Provides a resource to manage an Oceanus cluster.
---

# tencentcloud_oceanus_oceanus_cluster

Provides a resource to manage an Oceanus cluster.

## Example Usage

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name                     = "tf-example-oceanus"
  region_id                        = 1
  zone_id                          = 100003
  login_password                   = "Tf-example-12345"
  default_cos_bucket               = "tf-example-bucket-1250000000"
  cu                               = 12
  remark                           = "created by terraform"
  instance_charge_type             = "POSTPAID_BY_SECOND"
  flink_web_ui_network_access_type = "NetworkAccess_EXTERNAL"

  vpc_descriptions = [
    {
      vpc_id    = "vpc-12345678"
      subnet_id = "subnet-12345678"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) Cluster name. Limits: 1 to 50 English letters, Chinese characters, digits, hyphens (-) or underscores (_).
* `default_cos_bucket` - (Required, String) Stream compute uses Cloud Object Storage (COS) to store job checkpoints, jar packages, or delivered logs. If there is no available bucket in the local region, create one in the COS console. To ensure normal use, refer to the documentation for the lifecycle configuration of the corresponding COS bucket.
* `login_password` - (Required, String) Flink UI access password. The user name and password are used to log in to the Flink UI of the cluster; the cluster user name is admin.
* `region_id` - (Required, Int) Region ID. It can be queried through DescribeRegions of the region management system, with the Product parameter set to oceanus.
* `vpc_descriptions` - (Required, List) Stream compute accesses other cloud product resources in the same region through VPC and elastic network interfaces, and occupies a certain number of subnet IPs. Make sure the selected subnet has enough available IPs. If the existing network does not meet your requirements, create a new private network or a new subnet in the VPC console.
* `zone_id` - (Required, Int) Availability zone ID. It can be queried through DescribeZones of the region management system, with the Product parameter set to oceanus.
* `cluster_type` - (Optional, String) Cluster type. Valid values: `MULTI_AZ_CLUSTER` (multi-availability-zone cluster).
* `cu_memory` - (Optional, Int) Memory-to-core ratio. Only [0, 2, 4, 8] are supported.
* `cu` - (Optional, Int) Number of CUs of the cluster. 12 CU is the minimum computing resource and billing unit of stream compute; 1 CU contains 1 CPU and 4 GB of memory. When the number of CUs is greater than or equal to 48, the management node fee is waived.
* `flink_web_ui_network_access_type` - (Optional, String) Flink UI access type. Valid values: `NetworkAccess_INTERNAL` (internal network access), `NetworkAccess_EXTERNAL` (public network access). Default value: `NetworkAccess_EXTERNAL`.
* `instance_charge_type` - (Optional, String) Cluster billing type. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_SECOND` (pay-as-you-go). Default value: `POSTPAID_BY_SECOND`.
* `period` - (Optional, Int) Purchase duration in months. Value range: [1, 48].
* `remark` - (Optional, String) Cluster description. Limits: 1 to 50 English letters, Chinese characters, digits, hyphens (-) or underscores (_).
* `renew_flag` - (Optional, String) Auto-renewal flag. Valid values: `NOTIFY_AND_MANUAL_RENEW` (notify and renew manually), `NOTIFY_AND_AUTO_RENEW` (notify and renew automatically), `DISABLE_NOTIFY_AND_MANUAL_RENEW` (do not notify and do not renew automatically). Default value: `NOTIFY_AND_MANUAL_RENEW`. When InstanceChargeType is set to `PREPAID`, the cluster is a monthly-subscription cluster and the auto-renewal flag must be set; pay-as-you-go clusters do not need it.
* `slave_vpc_descriptions` - (Optional, List) VPCs of the multi-availability-zone cluster.

The `slave_vpc_descriptions` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) Private network ID.
* `app_id` - (Optional, Int) User APPID.
* `owner_uin` - (Optional, String) User UIN.

The `vpc_descriptions` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) Private network ID.
* `app_id` - (Optional, Int) User APPID.
* `owner_uin` - (Optional, String) User UIN.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, the Oceanus cluster ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateOceanusCluster](https://www.tencentcloud.com/search/?wd=CreateOceanusCluster&type=document&lang=en)
  - [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
* Read: [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
* Update:
  - [DescribeClusters](https://www.tencentcloud.com/search/?wd=DescribeClusters&type=document&lang=en)
  - [ScaleOceanusCluster](https://www.tencentcloud.com/search/?wd=ScaleOceanusCluster&type=document&lang=en)
* Delete: [DeleteOceanusCluster](https://www.tencentcloud.com/search/?wd=DeleteOceanusCluster&type=document&lang=en)

## Import

Oceanus cluster can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_oceanus_cluster.example cluster-12345678
```

