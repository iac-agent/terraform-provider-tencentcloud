## 1. Schema Definition

- [x] 1.1 Add `cos_private_access` field (TypeString, Optional, Computed) to the `origin` block schema in `resource_tc_teo_zone_setting.go`

## 2. Read Logic

- [x] 2.1 Add nil-check and set logic for `respData.Origin.CosPrivateAccess` in `resourceTencentCloudTeoZoneSettingRead` function, mapping to `originMap["cos_private_access"]`

## 3. Update Logic

- [x] 3.1 Add `cos_private_access` mapping from Terraform schema to `teo.Origin.CosPrivateAccess` in the `resourceTencentCloudTeoZoneSettingUpdate` function, within the `origin` block handling

## 4. Unit Tests

- [x] 4.1 Add unit test cases for `cos_private_access` in `resource_tc_teo_zone_setting_test.go`, covering read with non-nil value, read with nil value, and update scenarios

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_zone_setting.md` to include the `cos_private_access` parameter in the `origin` block example and description
