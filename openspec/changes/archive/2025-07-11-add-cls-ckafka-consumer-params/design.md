## Context

The `tencentcloud_cls_ckafka_consumer` Terraform resource manages CLS (Cloud Log Service) CKafka consumer delivery tasks. The current resource schema supports `topic_id`, `need_content`, `content`, `ckafka`, and `compression` parameters, but the CLS cloud API (cls/v20201016) now supports additional parameters that are not yet exposed in the Terraform provider.

The cloud API supports:
- **CreateConsumer**: TopicId, NeedContent, Content, Ckafka, Compression, RoleArn, ExternalId, AdvancedConfig
- **ModifyConsumer**: TopicId, Effective, NeedContent, Content, Ckafka, Compression, RoleArn, ExternalId, AdvancedConfig
- **DescribeConsumer** (response): Effective, NeedContent, Content, Ckafka, Compression

The missing parameters are: `effective` (only in ModifyConsumer input and DescribeConsumer output), `role_arn`, `external_id`, and `advanced_config` (with nested fields `partition_hash_status` and `partition_fields`).

Reference resource for code style: `tencentcloud_igtm_strategy`

## Goals / Non-Goals

**Goals:**
- Add `effective`, `role_arn`, `external_id`, and `advanced_config` parameters to the `tencentcloud_cls_ckafka_consumer` resource schema
- Update CRUD functions to correctly pass new parameters to cloud API calls
- Read `effective` from DescribeConsumer response and set it in state
- Add unit tests for the new parameter handling logic using gomonkey mock approach
- Update resource documentation (.md file)

**Non-Goals:**
- Changing existing parameter behavior or schema (backward compatibility required)
- Adding new cloud API calls or modifying service layer functions beyond what's needed for new parameters
- Modifying the `tencentcloud_cls_ckafka_consumer` resource ID logic (still uses `topic_id`)

## Decisions

1. **`effective` parameter is Optional + Computed**: Since `effective` is only settable via ModifyConsumer (not CreateConsumer), it should be Optional but also Computed so that the Read operation can populate it from the DescribeConsumer response. It will be set in the update flow after creation if specified.

2. **`advanced_config` as TypeList with MaxItems 1**: Following the existing pattern for `content` and `ckafka` in the same resource, `advanced_config` uses TypeList with MaxItems 1 and nested schema for its sub-fields.

3. **`role_arn` and `external_id` as TypeString Optional**: These are simple string parameters that map directly to the cloud API request fields. They are Optional since the cloud API marks them as omitempty.

4. **`effective` added to mutableArgs in update**: The `effective` field is mutable and should be included in the mutableArgs list in the update function, triggering a ModifyConsumer call when changed.

5. **`advanced_config` sub-fields handling**: `partition_hash_status` (TypeBool, Optional) and `partition_fields` (TypeSet of strings, Optional) map directly to `AdvancedConsumerConfiguration.PartitionHashStatus` and `AdvancedConsumerConfiguration.PartitionFields` in the SDK.

## Risks / Trade-offs

- [Risk] `effective` is not in CreateConsumer API → Mitigation: Mark as Optional + Computed. After create, if user specifies `effective = true`, the update flow will be triggered to call ModifyConsumer with Effective=true.
- [Risk] DescribeConsumer response does not include `role_arn`, `external_id`, `advanced_config` → Mitigation: These fields will only be set from Terraform state (not refreshed from API read). If the API doesn't return them, they won't be overwritten on read.
- [Risk] Backward compatibility → Mitigation: All new fields are Optional with no changes to existing fields, so existing Terraform configurations will continue to work without modification.
