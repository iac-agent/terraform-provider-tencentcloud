## ADDED Requirements

### Requirement: goods_num parameter in waf_saas_instance resource
The `tencentcloud_waf_saas_instance` resource SHALL include a `goods_num` parameter of type `TypeInt`, marked as `Optional`. When set, this value SHALL be used as `request.Goods.GoodsNum` in the `GenerateDealsAndPayNew` API call during resource creation. When not set, the value SHALL default to `1` (matching the existing hardcoded behavior). This parameter SHALL be immutable after creation (included in `immutableArgs` in the Update method).

#### Scenario: goods_num is specified during creation
- **WHEN** a user sets `goods_num = 2` in the terraform configuration
- **THEN** the `GenerateDealsAndPayNew` API request SHALL set `Goods.GoodsNum` to 2

#### Scenario: goods_num is not specified during creation
- **WHEN** a user does not set `goods_num` in the terraform configuration
- **THEN** the `GenerateDealsAndPayNew` API request SHALL set `Goods.GoodsNum` to 1 (preserving existing behavior)

#### Scenario: goods_num is changed after creation
- **WHEN** a user attempts to change `goods_num` after resource creation
- **THEN** the Update method SHALL return an error stating that `goods_num` cannot be changed

### Requirement: pid parameter in waf_saas_instance resource
The `tencentcloud_waf_saas_instance` resource SHALL include a `pid` parameter of type `TypeInt`, marked as `Optional`. When set, this value SHALL override the Pid derived from the `goods_category` mapping (`PID_SAAS`) and be used as `request.Goods.GoodsDetail.Pid` in the `GenerateDealsAndPayNew` API call during resource creation. When not set, the Pid SHALL continue to be derived from `goods_category` via the `PID_SAAS` mapping (preserving existing behavior). This parameter SHALL be immutable after creation (included in `immutableArgs` in the Update method).

#### Scenario: pid is specified during creation
- **WHEN** a user sets `pid = 1000827` in the terraform configuration
- **THEN** the `GenerateDealsAndPayNew` API request SHALL set `Goods.GoodsDetail.Pid` to 1000827, overriding the value derived from `goods_category`

#### Scenario: pid is not specified during creation
- **WHEN** a user does not set `pid` in the terraform configuration and sets `goods_category = "premium_saas"`
- **THEN** the `GenerateDealsAndPayNew` API request SHALL set `Goods.GoodsDetail.Pid` to the value from `PID_SAAS["premium_saas"]` (preserving existing behavior)

#### Scenario: pid is changed after creation
- **WHEN** a user attempts to change `pid` after resource creation
- **THEN** the Update method SHALL return an error stating that `pid` cannot be changed

### Requirement: region_id parameter in waf_saas_instance resource
The `tencentcloud_waf_saas_instance` resource SHALL include a `region_id` parameter of type `TypeInt`, marked as `Optional`. When set, this value SHALL override the RegionId derived from the provider region and be used as `request.Goods.RegionId` in the `GenerateDealsAndPayNew` API call during resource creation. When not set, the RegionId SHALL continue to be derived from the provider region (`REGION_ID_MAINLAND` for ap-guangzhou, `REGION_ID_NON_MAINLAND` for ap-seoul) (preserving existing behavior). This parameter SHALL be immutable after creation (included in `immutableArgs` in the Update method).

#### Scenario: region_id is specified during creation
- **WHEN** a user sets `region_id = 1` in the terraform configuration
- **THEN** the `GenerateDealsAndPayNew` API request SHALL set `Goods.RegionId` to 1, overriding the value derived from the provider region

#### Scenario: region_id is not specified during creation
- **WHEN** a user does not set `region_id` and the provider region is `ap-guangzhou`
- **THEN** the `GenerateDealsAndPayNew` API request SHALL set `Goods.RegionId` to `REGION_ID_MAINLAND` (preserving existing behavior)

#### Scenario: region_id is changed after creation
- **WHEN** a user attempts to change `region_id` after resource creation
- **THEN** the Update method SHALL return an error stating that `region_id` cannot be changed

### Requirement: deal_names parameter in waf_saas_instance resource
The `tencentcloud_waf_saas_instance` resource SHALL include a `deal_names` parameter of type `TypeList` with `ElementType: schema.TypeString`, marked as `Computed`. After successful resource creation via the `GenerateDealsAndPayNew` API, the `response.Data.DealNames` values SHALL be set to the `deal_names` field in the terraform state. Since the `DescribeInstances` API does not return deal names, this field SHALL NOT be populated during the Read operation.

#### Scenario: deal_names populated after creation
- **WHEN** the `GenerateDealsAndPayNew` API returns `response.Data.DealNames = ["20250101-order-001", "20250101-order-002"]`
- **THEN** the `deal_names` field in terraform state SHALL contain ["20250101-order-001", "20250101-order-002"]

#### Scenario: deal_names after refresh
- **WHEN** a terraform refresh is performed on the resource
- **THEN** the `deal_names` field SHALL remain empty (not populated by the DescribeInstances Read API), and terraform SHALL detect a drift for this computed field
