## ADDED Requirements

### Requirement: Resource schema definition
The `tencentcloud_teo_plan_v2` resource SHALL define the following schema fields:
- `plan_type` (TypeString, Required, ForceNew): Plan type with values "personal", "basic", "standard", "enterprise"
- `auto_use_voucher` (TypeString, Optional, ForceNew): Whether to auto-use vouchers, values "true"/"false", only effective for prepaid plans
- `prepaid_plan_param` (TypeList, Optional, ForceNew, MaxItems: 1): Nested block containing:
  - `period` (TypeInt, Optional): Subscription period in months (1-12, 24, 36), default 1
  - `renew_flag` (TypeString, Optional): Auto-renewal flag, values "on"/"off", default "off"
- `renew_flag` (TypeString, Optional): Top-level auto-renewal flag for ModifyPlan, values "on"/"off"
- `plan_id` (TypeString, Computed): Plan ID from CreatePlan response
- `deal_name` (TypeString, Computed): Order number from CreatePlan response
- `area` (TypeString, Computed): Service area
- `status` (TypeString, Computed): Plan status
- `pay_mode` (TypeString, Computed): Payment mode
- `enabled_time` (TypeString, Computed): Plan effective time
- `expired_time` (TypeString, Computed): Plan expiry time

#### Scenario: Schema fields match API parameters
- **WHEN** the resource schema is defined
- **THEN** all CreatePlan input parameters (plan_type, auto_use_voucher, prepaid_plan_param) SHALL be present as schema fields
- **AND** CreatePlan output parameters (plan_id, deal_name) SHALL be present as computed fields
- **AND** ModifyPlan input parameters (plan_id, renew_flag) SHALL be supported in the update path

### Requirement: Create operation
The resource SHALL create a TE0 plan using the CreatePlan API.

#### Scenario: Successful plan creation
- **WHEN** `tencentcloud_teo_plan_v2` resource is created with plan_type set
- **THEN** the CreatePlan API SHALL be called with plan_type, auto_use_voucher (if provided), and prepaid_plan_param (if provided)
- **AND** plan_id from the response SHALL be set as the resource ID
- **AND** deal_name from the response SHALL be stored in state

#### Scenario: CreatePlan returns nil response
- **WHEN** CreatePlan API returns nil response
- **THEN** a NonRetryableError SHALL be returned

#### Scenario: CreatePlan returns empty plan_id
- **WHEN** CreatePlan API returns response with nil PlanId
- **THEN** a NonRetryableError SHALL be returned

### Requirement: Read operation
The resource SHALL read plan information using the DescribePlans API via the existing `TeoService.DescribeTeoPlansById` service method.

#### Scenario: Plan exists
- **WHEN** the resource is read and the plan exists
- **THEN** all computed fields (plan_id, area, status, pay_mode, enabled_time, expired_time) SHALL be populated from the DescribePlans response
- **AND** plan_type SHALL be set from the response
- **AND** nil response fields SHALL be skipped (not set)

#### Scenario: Plan not found
- **WHEN** the resource is read and the plan does not exist
- **THEN** the resource ID SHALL be cleared (d.SetId(""))
- **AND** no error SHALL be returned

### Requirement: Update operation
The resource SHALL support updating the `renew_flag` field via the ModifyPlan API.

#### Scenario: Update renew_flag
- **WHEN** renew_flag field changes
- **THEN** the ModifyPlan API SHALL be called with plan_id and RenewFlag struct containing Switch field set to the new value

#### Scenario: No update needed
- **WHEN** no mutable fields have changed
- **THEN** no API call SHALL be made
- **AND** the Read operation SHALL still be called to sync state

#### Scenario: ForceNew fields changed
- **WHEN** plan_type, auto_use_voucher, or prepaid_plan_param changes
- **THEN** Terraform SHALL force resource recreation (these fields are ForceNew)

### Requirement: Delete operation
The resource SHALL delete a plan using the DestroyPlan API.

#### Scenario: Successful plan deletion
- **WHEN** `tencentcloud_teo_plan_v2` resource is destroyed
- **THEN** the DestroyPlan API SHALL be called with plan_id

### Requirement: Resource import
The resource SHALL support import via plan_id.

#### Scenario: Import existing plan
- **WHEN** `terraform import tencentcloud_teo_plan_v2.xxx <plan_id>` is executed
- **THEN** the resource state SHALL be populated by reading the plan via DescribePlans API

### Requirement: Provider registration
The resource SHALL be registered in provider.go and provider.md.

#### Scenario: Resource registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_plan_v2` SHALL be available as a resource
- **AND** it SHALL map to `teo.ResourceTencentCloudTeoPlanV2()`

### Requirement: Unit tests
The resource SHALL have unit tests using gomonkey mock approach.

#### Scenario: Unit test coverage
- **WHEN** unit tests are run
- **THEN** Create, Read, Update, Delete operations SHALL be tested with mocked cloud API calls
- **AND** tests SHALL use `go test -gcflags=all=-l` to run

### Requirement: Documentation
The resource SHALL have a .md documentation file following the project format.

#### Scenario: Documentation exists
- **WHEN** the resource is implemented
- **THEN** a `resource_tc_teo_plan_v2.md` file SHALL exist in the teo service directory
- **AND** it SHALL include a one-line description with "TEO" product name
- **AND** it SHALL include Example Usage section
- **AND** it SHALL include Import section
- **AND** it SHALL NOT include Argument Reference or Attribute Reference sections
