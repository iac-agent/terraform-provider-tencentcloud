Use this data source to query Docker activities for a Lighthouse instance.

Example Usage

Query Docker activities by instance ID

```hcl
data "tencentcloud_lighthouse_docker_activitie" "example" {
  instance_id = "lhins-12345678"
}
```

Query Docker activities by instance ID and activity IDs

```hcl
data "tencentcloud_lighthouse_docker_activitie" "example" {
  instance_id  = "lhins-12345678"
  activity_ids = ["lhda-12345678", "lhda-87654321"]
}
```

Query Docker activities by time range

```hcl
data "tencentcloud_lighthouse_docker_activitie" "example" {
  instance_id        = "lhins-12345678"
  created_time_begin = 1717200000
  created_time_end   = 1719800000
}
```
