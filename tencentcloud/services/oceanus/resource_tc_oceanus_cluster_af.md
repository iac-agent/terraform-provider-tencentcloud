Provides a resource to create and manage an Oceanus cluster.

Example Usage

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

Import

Oceanus cluster can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_cluster_af.example cluster-12345
```
