Provides a resource to create a Oceanus cluster

Example Usage

```hcl
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name       = "tf_example_cluster"
  region_id          = 1
  zone_id            = 100003
  login_password     = "Password@123"
  default_cos_bucket = "tf-example-bucket-1250000000"
  vpc_descriptions {
    vpc_id    = "vpc-xxxxxx"
    subnet_id = "subnet-xxxxxx"
  }
  instance_charge_type           = "POSTPAID_BY_SECOND"
  flink_web_ui_network_access_type = "NetworkAccess_EXTERNAL"
}
```

Import

Oceanus cluster can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_oceanus_cluster.example cluster-xxxxxx
```
