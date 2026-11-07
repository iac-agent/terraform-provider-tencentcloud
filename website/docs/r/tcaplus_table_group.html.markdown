---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_table_group"
sidebar_current: "docs-tencentcloud-resource-tcaplus_table_group"
description: |-
  Provides a resource to create a TcaplusDB table group
---

# tencentcloud_tcaplus_table_group

Provides a resource to create a TcaplusDB table group

## Example Usage

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

  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

### Create a TcaplusDB table group with a user specified table group id

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

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of the TcaplusDB cluster to which the table group belongs.
* `table_group_name` - (Required, String) Table group name, which can consist of Chinese characters, English letters or digits, with a maximum length of 32 characters.
* `resource_tags` - (Optional, Set) List of tags of the table group.
* `table_group_id` - (Optional, String, ForceNew) Table group ID. It can be specified by the user but must be unique within the cluster. If not specified, the API assigns an auto-incremented value. Changing this forces a new resource to be created.

The `resource_tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Creation time of the table group.
* `enable_mysql` - Whether MySQL load balancing is enabled. 0: not enabled, 1: enabling, 2: enabled.
* `mysql_conn_ip` - VIP of the MySQL load balancing.
* `mysql_conn_port` - VPort of the MySQL load balancing.
* `table_count` - Number of tables contained in the table group.
* `total_size` - Total storage size of the tables contained in the table group, in MB.
* `txh_backup_expire_day` - Number of days after which the Txh backup file of the tables expires and is deleted.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateTableGroup](https://www.tencentcloud.com/document/api/1016/35027)
  - [DescribeTableGroupTags](https://www.tencentcloud.com/document/api/1016/36900)
  - [DescribeTableGroups](https://www.tencentcloud.com/document/api/1016/35025)
* Read:
  - [DescribeTableGroupTags](https://www.tencentcloud.com/document/api/1016/36900)
  - [DescribeTableGroups](https://www.tencentcloud.com/document/api/1016/35025)
* Update:
  - [DescribeTableGroupTags](https://www.tencentcloud.com/document/api/1016/36900)
  - [DescribeTableGroups](https://www.tencentcloud.com/document/api/1016/35025)
  - [DescribeTasks](https://www.tencentcloud.com/document/api/1016/35045)
  - [ModifyTableGroupName](https://www.tencentcloud.com/document/api/1016/35024)
  - [ModifyTableGroupTags](https://www.tencentcloud.com/document/api/1016/36899)
* Delete: [DeleteTableGroup](https://www.tencentcloud.com/document/api/1016/35026)

## Import

TcaplusDB table group can be imported using the `cluster_id#table_group_id`, e.g.

```
terraform import tencentcloud_tcaplus_table_group.example 5516511420#109
```

