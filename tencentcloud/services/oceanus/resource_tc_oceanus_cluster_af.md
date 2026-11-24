Provides a resource to manage an Oceanus cluster.

Example Usage

```hcl
resource "tencentcloud_oceanus_cluster_af" "example" {
  cluster_name        = "tf-example-cluster"
  region_id           = 1
  zone_id             = 100002
  login_password      = "Example#2024"
  default_cos_bucket  = "example-1250000000"
  cu                  = 12
  remark              = "created by terraform"
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

Import

Oceanus cluster can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_cluster_af.example cluster-12345678
```
