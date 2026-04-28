## Context

The `tencentcloud_cls_ckafka_consumer` Terraform resource manages CLS (Cloud Log Service) CKafka consumer delivery tasks. Currently, the resource supports basic parameters (topic_id, need_content, content, ckafka, compression) but is missing several parameters that the cloud API already supports:

- **effective**: Controls whether the delivery task is active (available in ModifyConsumer and DescribeConsumer)
- **role_arn**: Role ARN for cross-account CKafka access (available in CreateConsumer and ModifyConsumer)
- **external_id**: External ID for role assumption (available in CreateConsumer and ModifyConsumer)
- **advanced_config**: Advanced configuration including partition hashing (available in CreateConsumer and ModifyConsumer)

The cloud API structs confirm these fields exist in the SDK (`cls/v20201016/models.go`):
- `CreateConsumerRequest`: has `RoleArn *string`, `ExternalId *string`, `AdvancedConfig *AdvancedConsumerConfiguration`
- `ModifyConsumerRequest`: has `Effective *bool`, `RoleArn *string`, `ExternalId *string`, `AdvancedConfig *AdvancedConsumerConfiguration`
- `DescribeConsumerResponseParams`: has `Effective *bool` (does NOT return RoleArn, ExternalId, AdvancedConfig)
- `AdvancedConsumerConfiguration`: has `PartitionHashStatus *bool`, `PartitionFields []*string`

Current resource file: `tencentcloud/services/cls/resource_tc_cls_ckafka_consumer.go`

## Goals / Non-Goals

**Goals:**
- Add `effective`, `role_arn`, `external_id`, and `advanced_config` parameters to the resource schema
- Update create, read, and update functions to handle new parameters
- Maintain backward compatibility - all new parameters are Optional
- Update unit tests for new parameter logic
- Update documentation (.md file)

**Non-Goals:**
- Do not modify existing parameter behavior or schema definitions
- Do not add parameters not supported by the cloud API
- Do not change the resource ID format (still uses topic_id)

## Decisions

### Decision 1: `effective` parameter handling
- `effective` is set via ModifyConsumer (not CreateConsumer) and returned by DescribeConsumer
- In the create function: do not set `effective` since CreateConsumer does not support it; after creation, the task may default to a specific effective state
- In the read function: read `effective` from DescribeConsumer response and set it in state
- In the update function: include `effective` in the ModifyConsumer request when changed
- Add `effective` to the `mutableArgs` list in the update function

**Rationale**: The cloud API only allows setting `effective` through ModifyConsumer, not CreateConsumer. The resource should still expose this parameter so users can control task effectiveness after creation.

### Decision 2: Write-only parameters (role_arn, external_id)
- `role_arn` and `external_id` are input-only parameters (not returned by DescribeConsumer)
- These parameters MUST be marked as Optional in the schema
- In the read function: these cannot be populated from the API response, so they rely on Terraform's state management
- In the update function: these should be included in the ModifyConsumer request when changed

**Rationale**: Since DescribeConsumer does not return these fields, the Terraform state will preserve the last configured values. This is a common pattern for write-only parameters in Terraform providers.

### Decision 3: `advanced_config` parameter structure
- `advanced_config` is a TypeList with MaxItems:1, containing sub-fields:
  - `partition_hash_status` (bool, Optional): Ckafka partition hash status, default false
  - `partition_fields` (TypeSet of string, Optional): List of fields for hash calculation, max 5 fields
- In the create and update functions: map to `AdvancedConsumerConfiguration` struct
- In the read function: since DescribeConsumer does not return `AdvancedConfig`, this is a write-only parameter like role_arn/external_id

**Rationale**: Follows the same pattern as the existing `content` and `ckafka` TypeList parameters in this resource.

### Decision 4: Update mutableArgs list
- Add `effective`, `role_arn`, `external_id`, `advanced_config` to the `mutableArgs` list in the update function
- This ensures the ModifyConsumer API is called when any of these parameters change

## Risks / Trade-offs

- **[Write-only parameters]** → `role_arn`, `external_id`, and `advanced_config` are not returned by DescribeConsumer. If the API silently changes these values, Terraform will not detect the drift. Mitigation: This is an inherent limitation of the API design; document this behavior clearly.
- **[effective not in CreateConsumer]** → The `effective` parameter cannot be set during resource creation. Users must update the resource after creation to set effectiveness. Mitigation: This follows the API's design where effective defaults to a value and can be modified later.
