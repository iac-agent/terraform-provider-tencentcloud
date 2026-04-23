---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_realtime_log_delivery"
sidebar_current: "docs-tencentcloud-resource-teo_realtime_log_delivery"
description: |-
  Provides a resource to create a TEO realtime log delivery task.
---

# tencentcloud_teo_realtime_log_delivery

Provides a resource to create a TEO realtime log delivery task.

## Example Usage

```hcl
resource "tencentcloud_teo_realtime_log_delivery" "example" {
  zone_id   = "zone-2qtuhspy7cr6"
  task_name = "test-task"
  task_type = "cls"
  entity_list = [
    "domain.example.com",
  ]
  log_type = "domain"
  area     = "mainland"
  fields = [
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

### Query with filters

```hcl
resource "tencentcloud_teo_realtime_log_delivery" "example" {
  zone_id   = "zone-2qtuhspy7cr6"
  task_name = "test-task"
  task_type = "cls"
  entity_list = [
    "domain.example.com",
  ]
  log_type = "domain"
  area     = "mainland"
  fields = [
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

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) Data delivery area, possible values are: `mainland`: within mainland China; `overseas`: worldwide (excluding mainland China).
* `entity_list` - (Required, List: [`String`]) List of entities (seven-layer domain names or four-layer proxy instances) corresponding to real-time log delivery tasks. Example values are as follows: Seven-layer domain name: `domain.example.com`; four-layer proxy instance: sid-2s69eb5wcms7. For values, refer to: `https://cloud.tencent.com/document/api/1552/80690`, `https://cloud.tencent.com/document/api/1552/86336`.
* `fields` - (Required, List: [`String`]) A list of preset fields for delivery.
* `log_type` - (Required, String) Data delivery type, the values are: `domain`: site acceleration log; `application`: four-layer proxy log; `web-rateLiming`: rate limit and CC attack protection log; `web-attack`: managed rule log; `web-rule`: custom rule log; `web-bot`: Bot management log.
* `sample` - (Required, Int) The sampling ratio is in thousandths, with a value range of 1-1000. For example, filling in 605 means the sampling ratio is 60.5%. Leaving it blank means the sampling ratio is 100%.
* `task_name` - (Required, String) The name of the real-time log delivery task. The format is a combination of numbers, English, -, and _. The maximum length is 200 characters.
* `task_type` - (Required, String) The real-time log delivery task type. The possible values are: `cls`: push to Tencent Cloud CLS; `custom_endpoint`: push to a custom HTTP(S) address; `s3`: push to an AWS S3 compatible storage bucket address.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `cls` - (Optional, List) CLS configuration information. This parameter is required when TaskType is cls.
* `custom_endpoint` - (Optional, List) Customize the configuration information of the HTTP service. This parameter is required when TaskType is set to custom_endpoint.
* `custom_fields` - (Optional, List) The list of custom fields delivered supports extracting specified field values from HTTP request headers, response headers, and cookies. Custom field names cannot be repeated and cannot exceed 200 fields.
* `delivery_conditions` - (Optional, List) The filter condition for log delivery. If it is not filled, all logs will be delivered.
* `delivery_status` - (Optional, String) The status of the real-time log delivery task. The values are: `enabled`: enabled; `disabled`: disabled. Leave it blank to keep the original configuration. Not required when creating.
* `filters` - (Optional, List) Filter conditions for querying realtime log delivery tasks. Available filter names: `task-id`, `task-name`, `entity-list`, `task-type`.
* `log_format` - (Optional, List) The output format of log delivery. If it is not filled, it means the default format. The default format logic is as follows: when TaskType is `custom_endpoint`, the default format is an array of multiple JSON objects, each JSON object is a log; when TaskType is `s3`, the default format is JSON Lines; in particular, when TaskType is `cls`, the value of LogFormat.FormatType can only be json, and other parameters in LogFormat will be ignored. It is recommended not to pass LogFormat.
* `s3` - (Optional, List) Configuration information of AWS S3 compatible storage bucket. This parameter is required when TaskType is s3.

The `cls` object supports the following:

* `log_set_id` - (Required, String) Tencent Cloud CLS log set ID.
* `log_set_region` - (Required, String) The region where the Tencent Cloud CLS log set is located.
* `topic_id` - (Required, String) Tencent Cloud CLS log topic ID.

The `conditions` object of `delivery_conditions` supports the following:

* `key` - (Required, String) The key of the filter condition.
* `operator` - (Required, String) Query condition operator, operation types are: `equals`: equal; `notEquals`: not equal; `include`: include; `notInclude`: not include; `startWith`: start with value; `notStartWith`: not start with value; `endWith`: end with value; `notEndWith`: not end with value.
* `value` - (Required, List) The value of the filter condition.

The `custom_endpoint` object supports the following:

* `url` - (Required, String) The custom HTTP interface address for real-time log delivery. Currently, only HTTP/HTTPS protocols are supported.
* `access_id` - (Optional, String) Fill in a custom SecretId to generate an encrypted signature. This parameter is required if the source site requires authentication.
* `access_key` - (Optional, String) Fill in the custom SecretKey to generate the encrypted signature. This parameter is required if the source site requires authentication.
* `compress_type` - (Optional, String) Data compression type, the possible values are: `gzip`: use gzip compression. If it is not filled in, compression is not enabled.
* `headers` - (Optional, List) The custom request header carried when delivering logs. If the header name you fill in is the default header carried by EdgeOne log push, such as Content-Type, then the header value you fill in will overwrite the default value. The header value references a single variable ${batchSize} to obtain the number of logs included in each POST request.
* `protocol` - (Optional, String) When sending logs via POST request, the application layer protocol type used can be: `http`: HTTP protocol; `https`: HTTPS protocol. If not filled in, the protocol type will be parsed according to the filled in URL address.

The `custom_fields` object supports the following:

* `name` - (Required, String) Extract data from the specified location in the HTTP request and response. The values are: `ReqHeader`: extract the specified field value from the HTTP request header; `RspHeader`: extract the specified field value from the HTTP response header; `Cookie`: extract the specified field value from the Cookie.
* `value` - (Required, String) The name of the parameter whose value needs to be extracted, for example: Accept-Language.
* `enabled` - (Optional, Bool) Whether to deliver this field. If left blank, this field will not be delivered.

The `delivery_conditions` object supports the following:

* `conditions` - (Optional, List) Log filtering conditions, the detailed filtering conditions are as follows: - `EdgeResponseStatusCode`: filter according to the status code returned by the EdgeOne node to the client. Supported operators: `equal`, `great`, `less`, `great_equal`, `less_equal`; Value range: any integer greater than or equal to 0; - `OriginResponseStatusCode`: filter according to the origin response status code. Supported operators: `equal`, `great`, `less`, `great_equal`, `less_equal`; Value range: any integer greater than or equal to -1; - `SecurityAction`: filter according to the final disposal action after the request hits the security rule. Supported operators: `equal`; Optional options are as follows: `-`: unknown/miss; `Monitor`: observe; `JSChallenge`: JavaScript challenge; `Deny`: intercept; `Allow`: allow; `BlockIP`: IP ban; `Redirect`: redirect; `ReturnCustomPage`: return to a custom page; `ManagedChallenge`: managed challenge; `Silence`: silent; `LongDelay`: respond after a long wait; `ShortDelay`: respond after a short wait; -`SecurityModule`: filter according to the name of the security module that finally handles the request. Supported operators: `equal`; Optional options: `-`: unknown/missed; `CustomRule`: Web Protection - Custom Rules; `RateLimitingCustomRule`: Web Protection - Rate Limiting Rules; `ManagedRule`: Web Protection - Managed Rules; `L7DDoS`: Web Protection - CC Attack Protection; `BotManagement`: Bot Management - Bot Basic Management; `BotClientReputation`: Bot Management - Client Profile Analysis; `BotBehaviorAnalysis`: Bot Management - Bot Intelligent Analysis; `BotCustomRule`: Bot Management - Custom Bot Rules; `BotActiveDetection`: Bot Management - Active Feature Recognition.

The `filters` object supports the following:

* `name` - (Required, String) The filter field name. Available values: `task-id`, `task-name`, `entity-list`, `task-type`.
* `values` - (Required, List) The filter values. Maximum 20 values per filter.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query.

The `headers` object of `custom_endpoint` supports the following:

* `name` - (Required, String) HTTP header name.
* `value` - (Required, String) HTTP header value.

The `log_format` object supports the following:

* `format_type` - (Required, String) The default output format type for log delivery. The possible values are: `json`: Use the default log output format JSON Lines. The fields in a single log are presented as key-value pairs; `csv`: Use the default log output format csv. Only field values are presented in a single log, without field names.
* `batch_prefix` - (Optional, String) A string to be added before each log delivery batch. Each log delivery batch may contain multiple log records.
* `batch_suffix` - (Optional, String) A string to append after each log delivery batch.
* `field_delimiter` - (Optional, String) In a single log record, a string is inserted between fields as a separator. The possible values are: `	`: tab character; `,`: comma; `;`: semicolon.
* `record_delimiter` - (Optional, String) The string inserted between log records as a separator. The possible values are: `
`: newline character; `	`: tab character; `,`: comma.
* `record_prefix` - (Optional, String) A string to prepend to each log record.
* `record_suffix` - (Optional, String) A string to append to each log record.

The `s3` object supports the following:

* `access_id` - (Required, String) The Access Key ID used to access the bucket.
* `access_key` - (Required, String) The secret key used to access the bucket.
* `bucket` - (Required, String) Bucket name and log storage directory, for example: `your_bucket_name/EO-logs/`. If this directory does not exist in the bucket, it will be created automatically.
* `endpoint` - (Required, String) URLs that do not include bucket names or paths, for example: `https://storage.googleapis.com`, `https://s3.ap-northeast-2.amazonaws.com`, `https://cos.ap-nanjing.myqcloud.com`.
* `region` - (Required, String) The region where the bucket is located, for example: ap-northeast-2.
* `compress_type` - (Optional, String) Data compression type, the values are: gzip: gzip compression. If it is not filled in, compression is not enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `realtime_log_delivery_tasks` - A list of realtime log delivery tasks returned by the API.
  * `area` - Data delivery area.
  * `cls` - CLS configuration information.
    * `log_set_id` - Tencent Cloud CLS log set ID.
    * `log_set_region` - The region where the Tencent Cloud CLS log set is located.
    * `topic_id` - Tencent Cloud CLS log topic ID.
  * `create_time` - Creation time.
  * `custom_endpoint` - Customize the configuration information of the HTTP service.
    * `access_id` - A custom SecretId to generate an encrypted signature.
    * `access_key` - The custom SecretKey to generate the encrypted signature.
    * `compress_type` - Data compression type.
    * `headers` - The custom request header carried when delivering logs.
      * `name` - HTTP header name.
      * `value` - HTTP header value.
    * `protocol` - The application layer protocol type used when sending logs via POST request.
    * `url` - The custom HTTP interface address for real-time log delivery.
  * `custom_fields` - The list of custom fields delivered.
    * `enabled` - Whether to deliver this field.
    * `name` - Extract data from the specified location in the HTTP request and response.
    * `value` - The name of the parameter whose value needs to be extracted.
  * `delivery_conditions` - The filter condition for log delivery.
    * `conditions` - Log filtering conditions.
      * `key` - The key of the filter condition.
      * `operator` - Query condition operator.
      * `value` - The value of the filter condition.
  * `delivery_status` - The status of the real-time log delivery task.
  * `entity_list` - List of entities corresponding to real-time log delivery tasks.
  * `fields` - A list of preset fields for delivery.
  * `log_format` - The output format of log delivery.
    * `batch_prefix` - A string to be added before each log delivery batch.
    * `batch_suffix` - A string to append after each log delivery batch.
    * `field_delimiter` - In a single log record, a string is inserted between fields as a separator.
    * `format_type` - The default output format type for log delivery.
    * `record_delimiter` - The string inserted between log records as a separator.
    * `record_prefix` - A string to prepend to each log record.
    * `record_suffix` - A string to append to each log record.
  * `log_type` - Data delivery type.
  * `s3` - Configuration information of AWS S3 compatible storage bucket.
    * `access_id` - The Access Key ID used to access the bucket.
    * `access_key` - The secret key used to access the bucket.
    * `bucket` - Bucket name and log storage directory.
    * `compress_type` - Data compression type.
    * `endpoint` - URLs that do not include bucket names or paths.
    * `region` - The region where the bucket is located.
  * `sample` - The sampling ratio.
  * `task_id` - Real-time log delivery task ID.
  * `task_name` - The name of the real-time log delivery task.
  * `task_type` - The real-time log delivery task type.
  * `update_time` - Update time.
* `task_id` - Real-time log delivery task ID.


## Import

TEO realtime log delivery can be imported using the id, e.g.

```
terraform import tencentcloud_teo_realtime_log_delivery.example zoneId#taskId
```

