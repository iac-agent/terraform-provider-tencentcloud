## Why

The `tencentcloud_waf_saas_instance` resource currently hardcodes `GoodsNum` to 1, derives `Pid` from the `goods_category` parameter via a static mapping (`PID_SAAS`), and derives `RegionId` from the provider region. This limits flexibility for users who need to specify custom Pid values, purchase multiple goods instances, or control the RegionId explicitly. Additionally, the `DealNames` (order numbers) from the creation response are not exposed, preventing users from tracking the billing orders associated with their WAF SaaS instance purchase.

## What Changes

- Add `goods_num` (TypeInt, Optional) parameter to `tencentcloud_waf_saas_instance` resource schema, mapping to `request.Goods.GoodsNum` in the `GenerateDealsAndPayNew` API
- Add `pid` (TypeInt, Optional) parameter to `tencentcloud_waf_saas_instance` resource schema, mapping to `request.Goods.GoodsDetail.Pid` in the `GenerateDealsAndPayNew` API
- Add `region_id` (TypeInt, Optional) parameter to `tencentcloud_waf_saas_instance` resource schema, mapping to `request.Goods.RegionId` in the `GenerateDealsAndPayNew` API
- Add `deal_names` (TypeList, Computed) parameter to `tencentcloud_waf_saas_instance` resource schema, mapping to `response.Data.DealNames` in the `GenerateDealsAndPayNew` API response

## Capabilities

### New Capabilities
- `waf-saas-instance-new-params`: Add new parameters (goods_num, pid, region_id, deal_names) to the tencentcloud_waf_saas_instance resource

### Modified Capabilities

## Impact

- `tencentcloud/services/waf/resource_tc_waf_saas_instance.go`: Schema and CRUD logic changes
- `tencentcloud/services/waf/resource_tc_waf_saas_instance_test.go`: Unit test updates
- `tencentcloud/services/waf/resource_tc_waf_saas_instance.md`: Documentation updates
