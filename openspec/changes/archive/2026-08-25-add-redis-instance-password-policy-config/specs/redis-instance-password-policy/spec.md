## ADDED Requirements

### Requirement: Resource schema for redis instance password policy

The system SHALL provide a Terraform resource named `tencentcloud_redis_instance_password_policy` of kind RESOURCE_KIND_CONFIG that manages the password complexity policy of an existing TencentCloud Redis instance. The resource schema SHALL expose the following top-level fields (flattened, no nested `password_policy` container):

- `instance_id` (string, Required): the Redis instance ID whose password policy is managed.
- `enabled` (bool, Required): whether the instance-level password complexity policy is enabled (`true` enables enforcement on all password changes; `false` disables it).
- `min_letter_count` (int, Optional): minimum number of letter (upper/lower case) characters. Range [1, 16], default 1.
- `min_digit_count` (int, Optional): minimum number of digit characters. Range [1, 16], default 1.
- `min_special_count` (int, Optional): minimum number of special characters. Range [1, 16], default 1.
- `min_length` (int, Optional): minimum total password length. Range [8, 64], default 8; MUST be greater than or equal to the sum of `min_letter_count`, `min_digit_count` and `min_special_count`.

The resource ID SHALL be the `instance_id`, because the password policy is bound one-to-one to a Redis instance.

#### Scenario: Required fields enforced in schema
- **WHEN** a user defines a `tencentcloud_redis_instance_password_policy` resource without `instance_id` or without `enabled`
- **THEN** Terraform SHALL report a validation error indicating `instance_id` and `enabled` are required.

#### Scenario: Resource ID equals instance id
- **WHEN** the resource is created with `instance_id = "crs-xxxx"`
- **THEN** the resource SHALL set its ID to `crs-xxxx`.

### Requirement: Read redis instance password policy

The system SHALL implement a Read operation that calls the cloud API `DescribeInstancePasswordPolicy` with the `InstanceId` obtained from the resource ID, reads the returned `PasswordPolicy` object, and populates the Terraform state with `enabled`, `min_letter_count`, `min_digit_count`, `min_special_count` and `min_length`.

The Read operation SHALL use `tccommon.ReadRetryTimeout` for retry and wrap cloud API errors with `tccommon.RetryError`. Before clearing the resource ID due to an empty result, the system SHALL log the current resource ID using the resource name `redis_instance_password_policy` to preserve forensic context. Each field SHALL be set only when the corresponding response field is not nil.

#### Scenario: Successful read populates state
- **WHEN** `DescribeInstancePasswordPolicy` returns a non-empty `PasswordPolicy` with `Enabled=true`, `MinLetterCount=2`, `MinDigitCount=1`, `MinSpecialCount=1`, `MinLength=8`
- **THEN** the system SHALL set `enabled=true`, `min_letter_count=2`, `min_digit_count=1`, `min_special_count=1`, `min_length=8` into state.

#### Scenario: Empty read preserves log context
- **WHEN** `DescribeInstancePasswordPolicy` returns an empty/nil `PasswordPolicy`
- **THEN** the system SHALL log `[CRUD]` with the resource name and current ID, then clear the resource ID from state.

### Requirement: Create redis instance password policy

The system SHALL implement a Create operation that sets the resource ID to the provided `instance_id` and then writes the password policy by delegating to the update logic (calling cloud API `ModifyInstancePasswordPolicy`). Create SHALL verify the cloud API response is non-empty before completing.

#### Scenario: Successful create writes policy
- **WHEN** a user creates the resource with `instance_id="crs-xxxx"`, `enabled=true` and optional min_* values
- **THEN** the system SHALL set the ID to `crs-xxxx`, call `ModifyInstancePasswordPolicy` with `InstanceId` and a `PasswordPolicy` containing the provided values, and then perform a Read to refresh state.

### Requirement: Update redis instance password policy

The system SHALL implement an Update operation that, when any schema field changes, calls cloud API `ModifyInstancePasswordPolicy` with the `InstanceId` from the resource ID and a `PasswordPolicy` object populated from the schema:
- `Enabled` SHALL always be sent (it is required by the cloud API).
- `MinLetterCount`, `MinDigitCount`, `MinSpecialCount`, `MinLength` SHALL be sent only when the corresponding schema field is set.
The update SHALL use `tccommon.WriteRetryTimeout` for retry, wrap cloud API errors with `tccommon.RetryError`, and perform a Read afterward to reconcile state.

#### Scenario: Successful update of enabled flag
- **WHEN** the user changes `enabled` from `false` to `true` for an existing resource
- **THEN** the system SHALL call `ModifyInstancePasswordPolicy` with `Enabled=true` plus any set min_* values, then refresh state from `DescribeInstancePasswordPolicy`.

### Requirement: Delete is state-only removal

The system SHALL implement a Delete operation that removes the resource from Terraform state without destroying or resetting the password policy on the Redis instance. Delete SHALL NOT call any cloud API, because the password policy is bound to the instance and the cloud API does not provide an independent reset semantic.

#### Scenario: terraform destroy does not alter instance policy
- **WHEN** the user runs `terraform destroy` on a `tencentcloud_redis_instance_password_policy` resource
- **THEN** the system SHALL remove the resource from state and return success without changing the instance's actual password policy.

### Requirement: Import existing password policy

The system SHALL support importing the resource via `terraform import tencentcloud_redis_instance_password_policy.<name> <instance_id>`, using `schema.ImportStatePassthrough`. After import, a Read SHALL populate the remaining fields from `DescribeInstancePasswordPolicy`.

#### Scenario: Import by instance id
- **WHEN** the user runs `terraform import tencentcloud_redis_instance_password_policy.example crs-xxxx`
- **THEN** the resource ID SHALL become `crs-xxxx` and a subsequent Read SHALL fill in `enabled` and min_* values from the cloud.
