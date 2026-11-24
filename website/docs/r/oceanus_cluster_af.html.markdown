---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_cluster_af"
sidebar_current: "docs-tencentcloud-resource-oceanus_cluster_af"
description: |-
  Provides a resource to manage an Oceanus cluster.
---

# tencentcloud_oceanus_cluster_af

Provides a resource to manage an Oceanus cluster.

## Example Usage

```hcl
resource "tencentcloud_oceanus_cluster_af" "example" {
  cluster_name         = "tf-example-cluster"
  region_id            = 1
  zone_id              = 100002
  login_password       = "Example#2024"
  default_cos_bucket   = "example-1250000000"
  cu                   = 12
  remark               = "created by terraform"
  instance_charge_type = "POSTPAID_BY_SECOND"

  vpc_descriptions = [
    {
      vpc_id    = "vpc-12345678"
      subnet_id = "subnet-12345678"
    }
  ]

  timeouts {
    create = "30m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) Cluster name. Input restrictions: 1-50 English letters, Chinese characters, digits, hyphens (-) or underscores (_) are supported.
* `default_cos_bucket` - (Required, String) Stream compute uses Cloud Object Storage (COS) to store job checkpoints, jar packages or delivered logs. If there is no available bucket in the local region, create one in the COS console. To ensure normal usage, refer to the document for the lifecycle configuration of the corresponding COS bucket.
* `login_password` - (Required, String) Flink UI access password. The username and password are used to log in to the Flink UI of the job; the cluster username is admin.
* `region_id` - (Required, Int) Region ID. It can be queried through DescribeRegions of the region management system with the Product parameter set to oceanus.
* `vpc_descriptions` - (Required, List) Stream compute accesses other cloud product resources in the same region through VPC and elastic network interfaces, and occupies a certain number of subnet IPs. Please ensure that the selected subnet has enough available IPs. If the existing network does not meet your requirements, create a new VPC or subnet in the VPC console.
* `zone_id` - (Required, Int) Availability zone ID. It can be queried through DescribeZones of the region management system with the Product parameter set to oceanus.
* `cluster_type` - (Optional, String) Cluster type. Valid values: MULTI_AZ_CLUSTER (multi availability zone cluster).
* `cu_memory` - (Optional, Int) Core memory ratio. Only [0, 2, 4, 8] are supported.
* `cu` - (Optional, Int) Number of cluster CUs. 12 CU is the minimum computing resource and billing unit of stream compute; 1 CU contains 1 CPU and 4 GB of memory. The management node fee is waived when the number of CUs is greater than or equal to 48.
* `flink_web_ui_network_access_type` - (Optional, String) Flink UI access type. Valid values: NetworkAccess_INTERNAL (intranet access) and NetworkAccess_EXTERNAL (public network access). Default value: NetworkAccess_EXTERNAL.
* `instance_charge_type` - (Optional, String) Cluster billing type. Valid values: PREPAID (monthly subscription) and POSTPAID_BY_SECOND (pay-as-you-go). Default value: POSTPAID_BY_SECOND.
* `period` - (Optional, Int) Purchase duration in months. Value range: [1, 48].
* `remark` - (Optional, String) Cluster description. Input restrictions: 1-50 English letters, Chinese characters, digits, hyphens (-) or underscores (_) are supported.
* `renew_flag` - (Optional, String) Auto-renewal flag. Valid values: NOTIFY_AND_MANUAL_RENEW (notify and renew manually), NOTIFY_AND_AUTO_RENEW (notify and renew automatically) and DISABLE_NOTIFY_AND_MANUAL_RENEW (neither notify nor renew automatically). Default value: NOTIFY_AND_MANUAL_RENEW. The flag is required for monthly subscription clusters (InstanceChargeType set to PREPAID) and is not required for pay-as-you-go clusters.
* `slave_vpc_descriptions` - (Optional, List) Multi availability zone VPC.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the create operation.

The `slave_vpc_descriptions` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.
* `app_id` - (Optional, Int) User AppId. Note: This field may return null, indicating that no valid value can be obtained.
* `owner_uin` - (Optional, String) User UIN. Note: This field may return null, indicating that no valid value can be obtained.

The `vpc_descriptions` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.
* `app_id` - (Optional, Int) User AppId. Note: This field may return null, indicating that no valid value can be obtained.
* `owner_uin` - (Optional, String) User UIN. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the cluster.

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
terraform import tencentcloud_oceanus_cluster_af.example cluster-12345678
```

