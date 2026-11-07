Provides a resource to create a TcaplusDB table group

Example Usage

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = "vpc-i5yyodl9"
  subnet_id                = "subnet-hhi88a58"
  password                 = "Password@2026"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_table_group" "example" {
  cluster_id      = tencentcloud_tcaplus_cluster.example.id
  table_group_name = "tf_example_table_group"

  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

Create a TcaplusDB table group with a user specified table group id

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = "vpc-i5yyodl9"
  subnet_id                = "subnet-hhi88a58"
  password                 = "Password@2026"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_table_group" "example" {
  cluster_id       = tencentcloud_tcaplus_cluster.example.id
  table_group_name = "tf_example_table_group"
  table_group_id   = "109"

  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

Import

TcaplusDB table group can be imported using the `cluster_id#table_group_id`, e.g.

```
terraform import tencentcloud_tcaplus_table_group.example 5516511420#109
```
