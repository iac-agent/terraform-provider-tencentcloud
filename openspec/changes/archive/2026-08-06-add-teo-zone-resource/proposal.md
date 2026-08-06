## Why

TEO (Tencent EdgeOne) is Tencent Cloud's edge security and acceleration platform. Currently, the Terraform provider lacks a `tencentcloud_teo_zone` resource to manage TEO sites (zones) through Infrastructure as Code. This resource enables users to create, read, update, and delete TEO zones with full lifecycle management, including CNAME/NS/no-domain access types, zone-level configuration, and tag support.

## What Changes

- Add new Terraform resource `tencentcloud_teo_zone` for managing TEO site lifecycle
- Implement CRUD operations using the following TEO cloud APIs:
  - `CreateZone` - create a new TEO site
  - `DescribeZones` - query site details by ID
  - `ModifyZone` - update site configuration (type, area, alias zone name, etc.)
  - `ModifyZoneStatus` - pause/resume the site
  - `ModifyZoneWorkMode` - update configuration group work mode
  - `DeleteZone` - delete the site
- Support zone-level parameters: `zone_name`, `type`, `area`, `plan_id`, `alias_zone_name`, `tags`, `paused`, `work_mode_infos`
- Support computed attributes: `zone_id`, `status`, `ownership_verification`, `name_servers`
- Register the resource in the provider under `tencentcloud_teo_zone`

## Capabilities

### New Capabilities

- `teo-zone-resource`: Full lifecycle management (CRUD) for TEO zones, including site creation with different access types (CNAME/NS/no-domain), zone configuration updates, pause/resume control, work mode management, and tag support.

### Modified Capabilities

<!-- No existing capabilities are modified -->

## Impact

- **Affected code**: `tencentcloud/services/teo/` directory
  - `resource_tc_teo_zone.go` - resource implementation
  - `resource_tc_teo_zone_extension.go` - post-handle and error handling extensions
  - `resource_tc_teo_zone_test.go` - unit tests
  - `resource_tc_teo_zone.md` - documentation
  - `service_tencentcloud_teo.go` - service layer methods (`DescribeTeoZone`, `DescribeTeoZoneById`, `ModifyZoneStatus`)
- **Provider registration**: `tencentcloud/provider.go` and `tencentcloud/provider.md`
- **Dependencies**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already in vendor)