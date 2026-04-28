---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_cloud_product_log_task_v2"
sidebar_current: "docs-tencentcloud-resource-cls_cloud_product_log_task_v2"
description: |-
  Provides a resource to create a cls cloud product log task
---

# tencentcloud_cls_cloud_product_log_task_v2

Provides a resource to create a cls cloud product log task

~> **NOTE:** In the destruction of resources, if cascading deletion of logset and topic is required, please set `force_delete` to `true` or use `is_delete_topic` and `is_delete_logset` parameters.

## Example Usage

### Create log delivery using the default newly created logset and topic with tags

```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-0an6hpv3"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"
  force_delete         = true
  tags = {
    "env" = "test"
  }
}
```

### Create log delivery using existing logset and topic with is_delete_topic and is_delete_logset

```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-0an6hpv3"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_id            = "ca5b4f56-1174-4eee-bc4c-69e48e0e8c45"
  topic_id             = "d8177ca9-466b-42f4-a110-5933daf0a83a"
  is_delete_topic      = true
  is_delete_logset     = true
}
```

## Argument Reference

The following arguments are supported:

* `assumer_name` - (Required, String, ForceNew) Cloud product identification, Values: CDS, CWP, CDB, TDSQL-C, MongoDB, TDStore, DCDB, MariaDB, PostgreSQL, BH, APIS.
* `cloud_product_region` - (Required, String, ForceNew) Cloud product region. There are differences in the input format of different log types in different regions. Please refer to the following example:
- CDS(all log type): ap-guangzhou
- CDB-AUDIT: gz
- TDSQL-C-AUDIT: gz
- MongoDB-AUDIT: gz
- MongoDB-SlowLog: ap-guangzhou
- MongoDB-ErrorLog: ap-guangzhou
- TDMYSQL-SLOW: gz
- DCDB(all log type): gz
- MariaDB(all log type): gz
- PostgreSQL(all log type): gz
- BH(all log type): overseas-polaris(Domestic sites overseas)/fsi-polaris(Domestic sites finance)/general-polaris(Domestic sites)/intl-sg-prod(International sites)
- APIS(all log type): gz.
* `cls_region` - (Required, String) CLS target region.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `log_type` - (Required, String, ForceNew) Log type, Values: CDS-AUDIT, CDS-RISK, CDB-AUDIT, TDSQL-C-AUDIT, MongoDB-AUDIT, MongoDB-SlowLog, MongoDB-ErrorLog, TDMYSQL-SLOW, DCDB-AUDIT, DCDB-SLOW, DCDB-ERROR, MariaDB-AUDIT, MariaDB-SLOW, MariaDB-ERROR, PostgreSQL-SLOW, PostgreSQL-ERROR, PostgreSQL-AUDIT, BH-FILELOG, BH-COMMANDLOG, APIS-ACCESS.
* `extend` - (Optional, String) Log configuration extension information, generally used to store additional log delivery configurations.
* `force_delete` - (Optional, Bool) Indicate whether to forcibly delete the corresponding logset and topic. If set to true, it will be forcibly deleted. Default is false.
* `is_delete_logset` - (Optional, Bool) Indicate whether to delete the associated logset when deleting the log collection. If set to true, the associated logset will be deleted. If the logset still has topics, it will not be deleted.
* `is_delete_topic` - (Optional, Bool) Indicate whether to delete the associated topic when deleting the log collection. If set to true, the associated topic will be deleted.
* `logset_id` - (Optional, String, ForceNew) Log set ID.
* `logset_name` - (Optional, String, ForceNew) Log set name, required if `logset_id` is not filled in. If the log set does not exist, it will be automatically created.
* `tags` - (Optional, Map) Tag description list. By specifying this parameter, you can bind tags to the corresponding topic at the same time. Maximum 10 tag key-value pairs are supported, and the same resource can only be bound to the same tag key.
* `topic_id` - (Optional, String, ForceNew) Log theme ID.
* `topic_name` - (Optional, String, ForceNew) The name of the log topic is required when `topic_id` is not filled in. If the log theme does not exist, it will be automatically created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Task status. Valid values: -1 (creating), 0 (creating), 1 (created), 2 (deleting), 3 (deleted).


## Import

cls cloud product log task can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cloud_product_log_task_v2.example postgres-1p7xvpc1#PostgreSQL#PostgreSQL-SLOW#gz
```

