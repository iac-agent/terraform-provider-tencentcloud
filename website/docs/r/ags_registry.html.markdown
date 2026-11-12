---
subcategory: "Agent Gateway Service(AGS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_registry"
sidebar_current: "docs-tencentcloud-resource-ags_registry"
description: |-
  Provides a resource to create an AGS (Agent Gateway Service) registry.
---

# tencentcloud_ags_registry

Provides a resource to create an AGS (Agent Gateway Service) registry.

-> **Note:** All configurable attributes are immutable on the cloud side (there is no Modify API). Changing `name`, `description`, `approval_mode` or `tags` forces the registry to be recreated.

## Example Usage

### Create a registry with tags

```hcl
resource "tencentcloud_ags_registry" "example" {
  name          = "tf-example-registry"
  description   = "created by terraform"
  approval_mode = "AUTO"

  tags {
    key   = "env"
    value = "prod"
  }

  tags {
    key   = "team"
    value = "devops"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the registry, unique within the same AppId and region, 1-255 characters long. Immutable after creation; changing it forces replacement.
* `approval_mode` - (Optional, String) Approval mode, decided at creation and cannot be changed afterwards; omitted defaults to AUTO. Valid values: `AUTO` (auto approve), `MANUAL` (manual approval). Case sensitive. Immutable after creation; changing it forces replacement.
* `description` - (Optional, String) Description text, up to 4096 characters. Immutable after creation; changing it forces replacement.
* `tags` - (Optional, List) Tencent Cloud custom tags bound at creation; keys must be unique, up to 10 tags. Immutable after creation; changing it forces replacement.

The `tags` object supports the following:

* `key` - (Required, String) Tag key. Reserved prefixes such as `qcs:`, `project` and project reserved prefixes are not allowed, and leading/trailing spaces are not allowed.
* `value` - (Required, String) Tag value. Leading/trailing spaces are not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `key` - Tag key.
    * `value` - Tag value.
  * `app_id` - AppId of the owning account.
  * `approval_mode` - Approval mode. Valid values: `AUTO`, `MANUAL`. Decided at creation and cannot be changed.
  * `create_time` - Creation time in ISO 8601 UTC, e.g. `2026-08-11T10:00:00Z`.
  * `creator_sub_account_uin` - UIN of the creator sub account; empty string when created directly by the main account.
  * `creator_uin` - UIN of the creator main account.
  * `description` - Description of the registry.
  * `name` - Registry name, unique within the same AppId and region.
  * `published_record_count` - Number of records bound to stable labels.
  * `record_count` - Total number of records under the registry.
  * `region` - Tencent Cloud region where the registry resides, e.g. `ap-guangzhou`.
  * `registry_id` - Registry ID, in the format `reg-` plus 8 lowercase letters/digits.
  * `status` - Status of the registry. Valid values: `ACTIVE`, `ARCHIVED`.
  * `tags` - Authoritative custom tags, stably sorted by key and value; an empty array is returned when there are no tags.
  * `update_time` - Time of the most recent update, ISO 8601 UTC.
* `id` - ID of the registry, same as registry_id.
* `registry_id` - ID of the created registry, in the format `reg-` plus 8 lowercase letters/digits.
* `registry` - Detailed information of the registry.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateRegistry](https://www.tencentcloud.com/search/?wd=CreateRegistry&type=document&lang=en)
  - [DescribeRegistry](https://www.tencentcloud.com/search/?wd=DescribeRegistry&type=document&lang=en)
* Read: [DescribeRegistry](https://www.tencentcloud.com/search/?wd=DescribeRegistry&type=document&lang=en)
* Update:
  - [GetPasswordRules](https://www.tencentcloud.com/document/api/598/81237)
  - [UpdatePasswordRules](https://www.tencentcloud.com/document/api/598/81238)
* Delete: [DeleteRegistry](https://www.tencentcloud.com/search/?wd=DeleteRegistry&type=document&lang=en)

## Import

AGS registry can be imported using the id, e.g.

```sh
terraform import tencentcloud_ags_registry.example reg-abcd1234
```

