Provides a resource to create a TEO realtime log delivery task.

Example Usage

```hcl
resource "tencentcloud_teo_realtime_log_delivery" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  task_name  = "test-task"
  task_type  = "cls"
  entity_list = [
    "domain.example.com",
  ]
  log_type = "domain"
  area     = "mainland"
  fields   = [
    "ServiceID",
    "ConnectTimeStamp",
  ]
  sample = 1000

  cls {
    log_set_id     = "cls-logset-id"
    topic_id       = "cls-topic-id"
    log_set_region = "ap-guangzhou"
  }
}
```

Query with filters

```hcl
resource "tencentcloud_teo_realtime_log_delivery" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  task_name  = "test-task"
  task_type  = "cls"
  entity_list = [
    "domain.example.com",
  ]
  log_type = "domain"
  area     = "mainland"
  fields   = [
    "ServiceID",
    "ConnectTimeStamp",
  ]
  sample = 1000

  cls {
    log_set_id     = "cls-logset-id"
    topic_id       = "cls-topic-id"
    log_set_region = "ap-guangzhou"
  }

  filters {
    name   = "task-type"
    values = ["cls"]
  }
}
```

Import

TEO realtime log delivery can be imported using the id, e.g.

```
terraform import tencentcloud_teo_realtime_log_delivery.example zoneId#taskId
```
