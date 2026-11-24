Provides a resource to manage an Oceanus cluster.

Example Usage

```hcl
resource "tencentcloud_oceanus_cluster" "example" {
  cluster_name         = "tf-example"
  region_id            = 1
  zone_id              = 100003
  login_password       = "Password@123"
  default_cos_bucket   = "tf-example-1258469122"
  cu                   = 12
  cu_memory            = 4
  remark               = "created by terraform"
  instance_charge_type = "POSTPAID_BY_SECOND"

  vpc_descriptions {
    vpc_id    = "vpc-4uwe5pmj"
    subnet_id = "subnet-1h1h1h1h"
  }
}
```

Import

Oceanus cluster can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_cluster.example cluster-1a2b3c4d
```
