## Context

The `tencentcloud_waf_saas_instance` resource manages WAF SaaS edition instances via the `GenerateDealsAndPayNew` API. Currently:

- `GoodsNum` is hardcoded to `1` in the Create method (line 168 of resource_tc_waf_saas_instance.go)
- `Pid` is automatically derived from `goods_category` via the `PID_SAAS` mapping table
- `RegionId` is automatically derived from the provider region (`REGION_ID_MAINLAND` or `REGION_ID_NON_MAINLAND`)
- `DealNames` (order numbers) from the `GenerateDealsAndPayNew` response are not exposed to users

All four new parameters map to the `GenerateDealsAndPayNew` API. The other CRUD APIs (DescribeInstances, ModifyInstanceName, ModifyInstanceRenewFlag, SwitchElasticMode, ModifyInstanceQpsLimit) are not affected by these parameters.

## Goals / Non-Goals

**Goals:**
- Add `goods_num` (Optional, TypeInt) to allow users to specify the number of goods instances
- Add `pid` (Optional, TypeInt) to allow users to specify a custom product ID
- Add `region_id` (Optional, TypeInt) to allow users to specify the region ID explicitly
- Add `deal_names` (Computed, TypeList of TypeString) to expose order numbers from the creation response
- Maintain full backward compatibility: when new optional parameters are not set, behavior must be identical to the current implementation

**Non-Goals:**
- Changing the existing `goods_category` parameter or its mapping logic
- Adding update support for `goods_num`, `pid`, `region_id` (these are create-only parameters)
- Modifying any other WAF resources or data sources

## Decisions

### Decision 1: goods_num as Optional with default behavior
- **Choice**: Add `goods_num` as Optional (no Default). When not set, use the existing hardcoded value of 1.
- **Rationale**: Backward compatibility is preserved. Users who don't set this parameter get the same behavior as before.

### Decision 2: pid as Optional overriding PID_SAAS mapping
- **Choice**: Add `pid` as Optional. When set, it overrides the Pid value derived from `goods_category`. When not set, the existing `PID_SAAS[goodsCategory]` mapping is used.
- **Rationale**: Allows users to specify custom Pid values while maintaining backward compatibility. The existing mapping covers standard use cases, but some users may need non-standard Pid values.

### Decision 3: region_id as Optional overriding provider region derivation
- **Choice**: Add `region_id` as Optional. When set, it overrides the RegionId value derived from the provider region. When not set, the existing logic (REGION_ID_MAINLAND or REGION_ID_NON_MAINLAND based on provider region) is used.
- **Rationale**: Allows users to specify RegionId explicitly while maintaining backward compatibility.

### Decision 4: deal_names as Computed TypeList
- **Choice**: Add `deal_names` as Computed TypeList of TypeString. Populate it from `response.Data.DealNames` after successful creation.
- **Rationale**: DealNames is `[]*string` in the API response. Since it's only available during creation and not from DescribeInstances, it should be Computed. It cannot be read back, so it will be empty on refresh.

### Decision 5: Immutable args for new create-only parameters
- **Choice**: Add `goods_num`, `pid`, `region_id` to the `immutableArgs` list in the Update method.
- **Rationale**: These parameters are only used in the `GenerateDealsAndPayNew` API (creation). There are no update APIs for them, so they must be immutable after creation.

## Risks / Trade-offs

- **[Risk] deal_names not refreshable** → The `DescribeInstances` API does not return deal names. On refresh, the `deal_names` field will become empty. This is acceptable since deal names are only relevant during the creation process and can be obtained from the initial creation response in the state.
- **[Risk] pid override may break goods_category consistency** → If a user sets a custom `pid` that doesn't match their `goods_category`, the billing order may not match expectations. Mitigation: Document clearly that `pid` overrides the default mapping from `goods_category`.
