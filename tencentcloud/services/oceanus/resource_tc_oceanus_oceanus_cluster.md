Provides a resource to manage an Oceanus cluster.

Example Usage

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name            = "tf-example-oceanus"
  region_id               = 1
  zone_id                 = 100003
  login_password          = "Tf-example-12345"
  default_cos_bucket      = "tf-example-bucket-1250000000"
  cu                      = 12
  remark                  = "created by terraform"
  instance_charge_type    = "POSTPAID_BY_SECOND"
  flink_web_ui_network_access_type = "NetworkAccess_EXTERNAL"

  vpc_descriptions = [
    {
      vpc_id    = "vpc-12345678"
      subnet_id = "subnet-12345678"
    }
  ]
}
```

Import

Oceanus cluster can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_oceanus_cluster.example cluster-12345678
```
