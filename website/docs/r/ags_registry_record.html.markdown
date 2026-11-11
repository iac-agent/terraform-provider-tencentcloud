---
subcategory: "AGS(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_registry_record"
sidebar_current: "docs-tencentcloud-resource-ags_registry_record"
description: |-
  Provides a resource to create a record in an AGS (Agent Governance Service) registry.
---

# tencentcloud_ags_registry_record

Provides a resource to create a record in an AGS (Agent Governance Service) registry.

## Example Usage

```hcl
resource "tencentcloud_ags_registry_record" "example" {
  registry_id     = "reg-abcdef01"
  name            = "my-record"
  descriptor_type = "MCP"

  description  = "created by terraform"
  version_name = "v1"

  mcp_source {
    type = "MANUAL"
    descriptors = jsonencode({
      name    = "mcp-server"
      version = "1.0.0"
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `descriptor_type` - (Required, String) Protocol descriptor type. Valid values: MCP, A2A, AGUI, CUSTOM, AGENT_SKILLS. Immutable after record creation.
* `name` - (Required, String) Record name, length 1..255. Unique per tenant and registry by normalized name (case-insensitive); reusable after soft delete. Cannot be changed after creation.
* `registry_id` - (Required, String) ID of the registry the record belongs to.
* `agent_source` - (Optional, Object) Source when descriptor_type is A2A or AGUI; required for those types and forbidden for others.
* `custom_descriptors` - (Optional, String) Required when descriptor_type is CUSTOM and forbidden for other types. Content must be a JSON object string.
* `description` - (Optional, String) Record description, up to 4096 characters. Set an empty string to clear it on update; removing the attribute keeps the previous value.
* `label_mutations` - (Optional, List) Label changes applied in Record update mode, at most 32 entries; label names must be unique within one request.
* `mcp_source` - (Optional, Object) Source when descriptor_type is MCP; required for MCP and forbidden for other types.
* `reason` - (Optional, String) Deletion reason, up to 1024 characters. Only sent by the Delete operation.
* `skill_source` - (Optional, Object) Source when descriptor_type is AGENT_SKILLS; required for AGENT_SKILLS and forbidden for other types.
* `version_name` - (Optional, String) Display name of Revision 1. Only used at creation time.

The `agent_source` object supports the following:

* `type` - (Required, String) Source type. Valid values: MANUAL (submit Agent descriptors JSON text directly), URL_IMPORT (import from a remote Agent Card or AGUI endpoint).
* `descriptors` - (Optional, String) Required when type is MANUAL: generic JSON object text. Validation is performed server side.
* `endpoint_url` - (Optional, String) Agent Card URL (A2A) or Runtime Endpoint URL (AGUI). Required when type is URL_IMPORT, HTTPS.

The `label_mutations` object supports the following:

* `name` - (Required, String) Label name, length 1..63, matching ^[a-z][a-z0-9._-]{0,62}$; stable and latest are reserved.
* `operation` - (Required, String) Operation type. Valid values: SET (create or move a label), DELETE (delete a custom label; reserved labels stable and latest cannot be deleted).
* `reason` - (Optional, String) Change reason, up to 1024 characters.
* `version_id` - (Optional, String) Target version ID. Optional for SET: omitted means unbound (allowed for custom labels, forbidden for stable); forbidden for DELETE.

The `mcp_source` object supports the following:

* `type` - (Required, String) Source type. Valid values: MANUAL (submit MCP descriptors JSON text directly), URL_IMPORT (import from a remote MCP server.json URL).
* `descriptors` - (Optional, String) Required when type is MANUAL: full MCP server.json object as a JSON string. Validation is performed server side.
* `endpoint_url` - (Optional, String) Remote MCP server.json URL (HTTPS). Required when type is URL_IMPORT.

The `skill_source` object supports the following:

* `type` - (Required, String) Source type. Valid values: MANUAL (submit SKILL.md text directly), TAR_PACKAGE (server issues a COS PUT pre-signed URL, the client uploads and the server validates asynchronously).
* `skill_md` - (Optional, String) SKILL.md text; required and non-empty when type is MANUAL, must not be provided when type is TAR_PACKAGE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `action_id` - Action ID.
    * `action_type` - Action type.
    * `actor_sub_account_uin` - Actor sub-account UIN.
    * `actor_type` - Actor type. Valid values: USER, SYSTEM.
    * `actor_uin` - Actor root account UIN.
    * `comment` - Action comment.
    * `create_time` - Occurrence time.
    * `harness_engine_version` - Specific version of the harness engine, for example 1.3.0.
    * `harness_engine` - Harness engine identifier, for example pi.
    * `name` - Metadata name, non-empty lowercase, unique within the array.
    * `request_id` - RequestId of the corresponding cloud API request.
    * `value` - Metadata string value.
  * `app_id` - AppId of the owning tenant.
  * `app_id` - AppId of the owning tenant.
  * `approval_actions` - Historical approval actions of the version.
  * `approval_mode` - Approval mode; locked at creation time.
  * `change_log` - Change reason of the version, up to 4096 characters; immutable.
  * `config_sha256` - SHA-256 of the normalized config content (used for idempotent deduplication).
  * `content_ready_time` - Content ready time, UTC RFC3339. Only returned for TAR_PACKAGE or FILE sources with content_status READY.
  * `content_sha256` - SHA-256 of READY content.
  * `content_size_bytes` - Size of READY content in bytes.
  * `content_status` - Content status.
  * `create_time` - Creation time, ISO 8601 UTC.
  * `create_time` - Creation time.
  * `creator_sub_account_uin` - Creator sub-account UIN.
  * `creator_sub_account_uin` - Creator sub-account UIN; empty when created by the root account.
  * `creator_uin` - Creator root account UIN.
  * `creator_uin` - Creator root account UIN.
  * `description` - Current description of the owning record. Only returned by the version list query.
  * `description` - Record description.
  * `descriptor_type` - Asset type of the owning record. Only returned by the version list query.
  * `descriptor_type` - Protocol descriptor type; immutable after creation.
  * `descriptors` - Protocol descriptor object (JSON string).
  * `label_set` - All label names under the record, including system labels (stable, latest) and custom labels.
  * `label_set` - Label names currently bound to this version.
  * `lifecycle_status` - Lifecycle status. Valid values: ACTIVE, DELETED (soft-delete tombstone).
  * `metadata` - Version metadata Name/Value entries, sorted by name ascending.
  * `name` - Current name of the owning record. Only returned by the version list query.
  * `name` - Record name; may repeat within one registry.
  * `record_id` - ID of the owning record.
  * `record_id` - Record ID, format `rec-` followed by 8 lowercase alphanumeric characters.
  * `registry_id` - ID of the owning registry.
  * `revision` - Incrementing version number, starting from 1.
  * `source_config` - Normalized source config object (JSON string).
  * `source_type` - Content source type.
  * `status` - Version status.
  * `supported_harness_versions` - Harness engines and versions supported by this asset version.
  * `update_time` - Last update time.
  * `update_time` - Last update time.
  * `version_count` - Number of non-deleted versions under the record.
  * `version_id` - Version ID, format `rv-` followed by 8 lowercase alphanumeric characters.
  * `version_name` - Version display name.
* `id` - Resource identifier in the form `<registry_id>#<record_id>`.
* `record_id` - Record ID, format `rec-` followed by 8 lowercase alphanumeric characters.
* `record` - Record metadata and all labels returned by the cloud.
* `resolved_by` - Resolution method. Valid values: DEFAULT_STABLE, LABEL, VERSION_ID.
* `resolved_label` - Label name returned when resolution goes through a label (DEFAULT_STABLE or LABEL).
* `version` - Full version resolved by the Describe operation (defaults to the stable label).

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateRegistryRecord](https://www.tencentcloud.com/search/?wd=CreateRegistryRecord&type=document&lang=en)
  - [DescribeRegistryRecord](https://www.tencentcloud.com/search/?wd=DescribeRegistryRecord&type=document&lang=en)
* Read: [DescribeRegistryRecord](https://www.tencentcloud.com/search/?wd=DescribeRegistryRecord&type=document&lang=en)
* Update:
  - [DescribeRegistryRecord](https://www.tencentcloud.com/search/?wd=DescribeRegistryRecord&type=document&lang=en)
  - [UpdateRegistryRecord](https://www.tencentcloud.com/search/?wd=UpdateRegistryRecord&type=document&lang=en)
* Delete: [DeleteRegistryRecord](https://www.tencentcloud.com/search/?wd=DeleteRegistryRecord&type=document&lang=en)

## Import

ags registry record can be imported using the id, e.g.

```
terraform import tencentcloud_ags_registry_record.example reg-abcdef01#rec-12345678
```

