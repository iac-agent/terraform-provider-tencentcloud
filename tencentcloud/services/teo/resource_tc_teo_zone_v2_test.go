package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaTeoZoneV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoZoneV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoZoneV2{}

func newMockMetaTeoZoneV2() *mockMetaTeoZoneV2 {
	return &mockMetaTeoZoneV2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringZoneV2(s string) *string {
	return &s
}

func ptrInt64ZoneV2(i int64) *int64 {
	return &i
}

func ptrBoolZoneV2(b bool) *bool {
	return &b
}

func TestTeoZoneV2_Create_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)

	// Mock Create API
	patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext", func(_ context.Context, request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
		assert.NotNil(t, request.Type)
		assert.Equal(t, "full", *request.Type)
		assert.NotNil(t, request.ZoneName)
		assert.Equal(t, "example.com", *request.ZoneName)
		assert.NotNil(t, request.Area)
		assert.Equal(t, "global", *request.Area)
		assert.NotNil(t, request.PlanId)
		assert.Equal(t, "edgeone-2y041pblwaxe", *request.PlanId)
		assert.NotNil(t, request.AliasZoneName)
		assert.Equal(t, "example-alias", *request.AliasZoneName)
		assert.NotNil(t, request.AllowDuplicates)
		assert.True(t, *request.AllowDuplicates)
		assert.NotNil(t, request.JumpStart)
		assert.True(t, *request.JumpStart)
		assert.NotNil(t, request.Tags)
		assert.Equal(t, 1, len(request.Tags))
		assert.Equal(t, "tag-key-1", *request.Tags[0].TagKey)
		assert.Equal(t, "tag-value-1", *request.Tags[0].TagValue)

		resp := teov20220901.NewCreateZoneResponse()
		resp.Response = &teov20220901.CreateZoneResponseParams{
			ZoneId: ptrStringZoneV2("zone-2noz78a8ev6k"),
			OwnershipVerification: &teov20220901.OwnershipVerification{
				DnsVerification: &teov20220901.DnsVerification{
					Subdomain:   ptrStringZoneV2("_teoauthzz"),
					RecordType:  ptrStringZoneV2("TXT"),
					RecordValue: ptrStringZoneV2("teoauth-verify-xxx"),
				},
				FileVerification: &teov20220901.FileVerification{
					Path:    ptrStringZoneV2("/.well-known/teo/verification/z12h416twn.txt"),
					Content: ptrStringZoneV2("teoauth-verify-xxx"),
				},
				NsVerification: &teov20220901.NsVerification{
					NameServers: []*string{ptrStringZoneV2("ns1.teonetdns.com")},
				},
			},
			RequestId: ptrStringZoneV2("fake-request-id-create"),
		}
		return resp, nil
	})

	// Mock Read API (called after Create)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext", func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(100), *request.Limit)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "zone-id", *request.Filters[0].Name)
		assert.Equal(t, "zone-2noz78a8ev6k", *request.Filters[0].Values[0])

		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64ZoneV2(1),
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrStringZoneV2("zone-2noz78a8ev6k"),
					ZoneName:      ptrStringZoneV2("example.com"),
					AliasZoneName: ptrStringZoneV2("example-alias"),
					Area:          ptrStringZoneV2("global"),
					Type:          ptrStringZoneV2("full"),
					Tags: []*teov20220901.Tag{
						{
							TagKey:   ptrStringZoneV2("tag-key-1"),
							TagValue: ptrStringZoneV2("tag-value-1"),
						},
					},
					Resources: []*teov20220901.Resource{
						{
							Id:         ptrStringZoneV2("edgeone-2y041pblwaxe"),
							PayMode:    ptrInt64ZoneV2(0),
							CreateTime: ptrStringZoneV2("2024-01-01T00:00:00Z"),
							Status:     ptrStringZoneV2("normal"),
							Sv: []*teov20220901.Sv{
								{
									Key:   ptrStringZoneV2("zone"),
									Value: ptrStringZoneV2("1"),
								},
							},
							AutoRenewFlag: ptrInt64ZoneV2(1),
							PlanId:        ptrStringZoneV2("edgeone-2y041pblwaxe"),
							Area:          ptrStringZoneV2("global"),
							Group:         ptrStringZoneV2("plan"),
							ZoneNumber:    ptrInt64ZoneV2(1),
							Type:          ptrStringZoneV2("plan"),
						},
					},
					NSDetail: &teov20220901.NSDetail{
						CnameSpeedUp: ptrStringZoneV2("disabled"),
						IsFake:       ptrInt64ZoneV2(0),
						OwnershipVerification: &teov20220901.OwnershipVerification{
							NsVerification: &teov20220901.NsVerification{
								NameServers: []*string{ptrStringZoneV2("ns1.teonetdns.com")},
							},
						},
						OriginalNameServers: []*string{ptrStringZoneV2("ns1.example.com")},
						NameServers:         []*string{ptrStringZoneV2("ns1.teonetdns.com")},
						VanityNameServers: &teov20220901.VanityNameServers{
							Switch:  ptrStringZoneV2("off"),
							Servers: []*string{ptrStringZoneV2("ns1.example.com")},
						},
						VanityNameServersIps: []*teov20220901.VanityNameServersIps{
							{
								Name: ptrStringZoneV2("ns1.example.com"),
								IPv4: ptrStringZoneV2("1.2.3.4"),
							},
						},
					},
					CreatedOn:    ptrStringZoneV2("2024-01-01T00:00:00Z"),
					ModifiedOn:   ptrStringZoneV2("2024-06-15T12:00:00Z"),
					Status:       ptrStringZoneV2("pending"),
					ActiveStatus: ptrStringZoneV2("active"),
					LockStatus:   ptrStringZoneV2("enable"),
					Paused:       ptrBoolZoneV2(false),
					WorkModeInfos: []*teov20220901.ConfigGroupWorkModeInfo{
						{
							ConfigGroupType: ptrStringZoneV2("l7_acceleration"),
							WorkMode:        ptrStringZoneV2("immediate_effect"),
						},
					},
				},
			},
			RequestId: ptrStringZoneV2("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZoneV2()
	res := teo.ResourceTencentCloudTeoZoneV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"type":             "full",
		"zone_name":        "example.com",
		"area":             "global",
		"plan_id":          "edgeone-2y041pblwaxe",
		"alias_zone_name":  "example-alias",
		"allow_duplicates": true,
		"jump_start":       true,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2noz78a8ev6k", d.Id())
	assert.Equal(t, "zone-2noz78a8ev6k", d.Get("zone_id"))
	assert.Equal(t, "full", d.Get("type"))
	assert.Equal(t, "example.com", d.Get("zone_name"))
	assert.Equal(t, "global", d.Get("area"))
	assert.Equal(t, "pending", d.Get("status"))
	assert.Equal(t, false, d.Get("paused"))

	// assert nested structures
	resourcesList := d.Get("resources").([]interface{})
	assert.Equal(t, 1, len(resourcesList))
	resourceMap := resourcesList[0].(map[string]interface{})
	assert.Equal(t, "edgeone-2y041pblwaxe", resourceMap["id"])
	assert.Equal(t, 0, resourceMap["pay_mode"])
	assert.Equal(t, "normal", resourceMap["status"])
	svList := resourceMap["sv"].([]interface{})
	assert.Equal(t, 1, len(svList))
	svMap := svList[0].(map[string]interface{})
	assert.Equal(t, "zone", svMap["key"])

	nsDetailList := d.Get("ns_detail").([]interface{})
	assert.Equal(t, 1, len(nsDetailList))
	nsDetailMap := nsDetailList[0].(map[string]interface{})
	assert.Equal(t, "disabled", nsDetailMap["cname_speed_up"])
	assert.Equal(t, 0, nsDetailMap["is_fake"])
	nsDetailOwnershipVerificationList := nsDetailMap["ownership_verification"].([]interface{})
	assert.Equal(t, 1, len(nsDetailOwnershipVerificationList))
	nsVerificationList := nsDetailOwnershipVerificationList[0].(map[string]interface{})["ns_verification"].([]interface{})
	assert.Equal(t, 1, len(nsVerificationList))
	assert.Equal(t, []interface{}{"ns1.teonetdns.com"}, nsVerificationList[0].(map[string]interface{})["name_servers"])

	vanityNameServersIpsList := nsDetailMap["vanity_name_servers_ips"].([]interface{})
	assert.Equal(t, 1, len(vanityNameServersIpsList))
	assert.Equal(t, "1.2.3.4", vanityNameServersIpsList[0].(map[string]interface{})["i_pv4"])

	tagsList := d.Get("tags").([]interface{})
	assert.Equal(t, 1, len(tagsList))
	assert.Equal(t, "tag-key-1", tagsList[0].(map[string]interface{})["tag_key"])
	assert.Equal(t, "tag-value-1", tagsList[0].(map[string]interface{})["tag_value"])

	workModeInfosList := d.Get("work_mode_infos").([]interface{})
	assert.Equal(t, 1, len(workModeInfosList))
	assert.Equal(t, "l7_acceleration", workModeInfosList[0].(map[string]interface{})["config_group_type"])
}

func TestTeoZoneV2_Create_EmptyZoneId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)

	// Mock Create API with empty ZoneId
	patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext", func(_ context.Context, request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
		resp := teov20220901.NewCreateZoneResponse()
		resp.Response = &teov20220901.CreateZoneResponseParams{
			ZoneId:    ptrStringZoneV2(""),
			RequestId: ptrStringZoneV2("fake-request-id-create"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZoneV2()
	res := teo.ResourceTencentCloudTeoZoneV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"type":      "partial",
		"zone_name": "example-empty.com",
		"area":      "global",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ZoneId is empty")
	assert.Equal(t, "", d.Id())
}

func TestTeoZoneV2_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)

	// Mock Read API
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext", func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(100), *request.Limit)
		assert.Equal(t, int64(0), *request.Offset)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "zone-id", *request.Filters[0].Name)
		assert.Equal(t, "zone-2noz78a8ev6k", *request.Filters[0].Values[0])

		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64ZoneV2(1),
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrStringZoneV2("zone-2noz78a8ev6k"),
					ZoneName:      ptrStringZoneV2("example.com"),
					AliasZoneName: ptrStringZoneV2("example-alias"),
					Area:          ptrStringZoneV2("global"),
					Type:          ptrStringZoneV2("full"),
					Tags: []*teov20220901.Tag{
						{
							TagKey:   ptrStringZoneV2("tag-key-1"),
							TagValue: ptrStringZoneV2("tag-value-1"),
						},
					},
					CreatedOn:    ptrStringZoneV2("2024-01-01T00:00:00Z"),
					ModifiedOn:   ptrStringZoneV2("2024-06-15T12:00:00Z"),
					Status:       ptrStringZoneV2("active"),
					CnameStatus:  ptrStringZoneV2("finished"),
					ActiveStatus: ptrStringZoneV2("active"),
					LockStatus:   ptrStringZoneV2("enable"),
					Paused:       ptrBoolZoneV2(false),
					NameServers:  []*string{ptrStringZoneV2("ns1.teonetdns.com"), ptrStringZoneV2("ns2.teonetdns.com")},
					VanityNameServers: &teov20220901.VanityNameServers{
						Switch:  ptrStringZoneV2("on"),
						Servers: []*string{ptrStringZoneV2("ns1.example.com"), ptrStringZoneV2("ns2.example.com")},
					},
					VanityNameServersIps: []*teov20220901.VanityNameServersIps{
						{
							Name: ptrStringZoneV2("ns1.example.com"),
							IPv4: ptrStringZoneV2("1.2.3.4"),
						},
					},
					WorkModeInfos: []*teov20220901.ConfigGroupWorkModeInfo{
						{
							ConfigGroupType: ptrStringZoneV2("web_security"),
							WorkMode:        ptrStringZoneV2("version_control"),
						},
					},
				},
			},
			RequestId: ptrStringZoneV2("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZoneV2()
	res := teo.ResourceTencentCloudTeoZoneV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"type":      "full",
		"zone_name": "example.com",
		"area":      "global",
	})
	d.SetId("zone-2noz78a8ev6k")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2noz78a8ev6k", d.Id())
	assert.Equal(t, "zone-2noz78a8ev6k", d.Get("zone_id"))
	assert.Equal(t, "example.com", d.Get("zone_name"))
	assert.Equal(t, "example-alias", d.Get("alias_zone_name"))
	assert.Equal(t, "global", d.Get("area"))
	assert.Equal(t, "full", d.Get("type"))
	assert.Equal(t, "active", d.Get("status"))
	assert.Equal(t, "finished", d.Get("cname_status"))
	assert.Equal(t, "active", d.Get("active_status"))
	assert.Equal(t, "enable", d.Get("lock_status"))
	assert.Equal(t, false, d.Get("paused"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-06-15T12:00:00Z", d.Get("modified_on"))
	assert.Equal(t, []interface{}{"ns1.teonetdns.com", "ns2.teonetdns.com"}, d.Get("name_servers"))

	// assert nested structures
	tagsList := d.Get("tags").([]interface{})
	assert.Equal(t, 1, len(tagsList))
	assert.Equal(t, "tag-key-1", tagsList[0].(map[string]interface{})["tag_key"])
	assert.Equal(t, "tag-value-1", tagsList[0].(map[string]interface{})["tag_value"])

	vanityNameServersList := d.Get("vanity_name_servers").([]interface{})
	assert.Equal(t, 1, len(vanityNameServersList))
	vanityNameServersMap := vanityNameServersList[0].(map[string]interface{})
	assert.Equal(t, "on", vanityNameServersMap["switch"])
	assert.Equal(t, []interface{}{"ns1.example.com", "ns2.example.com"}, vanityNameServersMap["servers"])

	vanityNameServersIpsList := d.Get("vanity_name_servers_ips").([]interface{})
	assert.Equal(t, 1, len(vanityNameServersIpsList))
	assert.Equal(t, "ns1.example.com", vanityNameServersIpsList[0].(map[string]interface{})["name"])
	assert.Equal(t, "1.2.3.4", vanityNameServersIpsList[0].(map[string]interface{})["i_pv4"])

	workModeInfosList := d.Get("work_mode_infos").([]interface{})
	assert.Equal(t, 1, len(workModeInfosList))
	assert.Equal(t, "web_security", workModeInfosList[0].(map[string]interface{})["config_group_type"])
	assert.Equal(t, "version_control", workModeInfosList[0].(map[string]interface{})["work_mode"])
}

func TestTeoZoneV2_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)

	// Mock Read API with empty Zones
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext", func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64ZoneV2(0),
			Zones:      []*teov20220901.Zone{},
			RequestId:  ptrStringZoneV2("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZoneV2()
	res := teo.ResourceTencentCloudTeoZoneV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example-notfound.com",
	})
	d.SetId("zone-notfound123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoZoneV2_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)

	// Mock Read API returning nil response
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext", func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		return nil, nil
	})

	meta := newMockMetaTeoZoneV2()
	res := teo.ResourceTencentCloudTeoZoneV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example-nilresp.com",
	})
	d.SetId("zone-nilresp123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoZoneV2_Update_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)

	// Mock Modify API
	patches.ApplyMethodFunc(teoClient, "ModifyZoneWithContext", func(_ context.Context, request *teov20220901.ModifyZoneRequest) (*teov20220901.ModifyZoneResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-2noz78a8ev6k", *request.ZoneId)
		assert.NotNil(t, request.Type)
		assert.Equal(t, "full", *request.Type)
		assert.NotNil(t, request.Area)
		assert.Equal(t, "mainland", *request.Area)
		assert.NotNil(t, request.AliasZoneName)
		assert.Equal(t, "example-alias-updated", *request.AliasZoneName)
		assert.NotNil(t, request.VanityNameServers)
		assert.Equal(t, "on", *request.VanityNameServers.Switch)
		assert.Equal(t, 2, len(request.VanityNameServers.Servers))
		assert.Equal(t, "ns1.example.com", *request.VanityNameServers.Servers[0])
		assert.Equal(t, "ns2.example.com", *request.VanityNameServers.Servers[1])

		resp := teov20220901.NewModifyZoneResponse()
		resp.Response = &teov20220901.ModifyZoneResponseParams{
			RequestId: ptrStringZoneV2("fake-request-id-update"),
		}
		return resp, nil
	})

	// Mock Read API (called after Update)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext", func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64ZoneV2(1),
			Zones: []*teov20220901.Zone{
				{
					ZoneId:       ptrStringZoneV2("zone-2noz78a8ev6k"),
					ZoneName:     ptrStringZoneV2("example.com"),
					Area:         ptrStringZoneV2("mainland"),
					Type:         ptrStringZoneV2("full"),
					Status:       ptrStringZoneV2("active"),
					ActiveStatus: ptrStringZoneV2("active"),
				},
			},
			RequestId: ptrStringZoneV2("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZoneV2()
	res := teo.ResourceTencentCloudTeoZoneV2()

	// construct the ResourceData with the new values
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"type":            "full",
		"zone_name":       "example.com",
		"area":            "mainland",
		"alias_zone_name": "example-alias-updated",
		"vanity_name_servers": []interface{}{
			map[string]interface{}{
				"switch":  "on",
				"servers": []interface{}{"ns1.example.com", "ns2.example.com"},
			},
		},
	})
	d.SetId("zone-2noz78a8ev6k")

	// Force mutable args to be detected as changed (TestResourceDataRaw has no old state)
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		switch key {
		case "type", "zone_name", "area", "alias_zone_name", "vanity_name_servers":
			return true
		}
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2noz78a8ev6k", d.Id())
	assert.Equal(t, "mainland", d.Get("area"))
	assert.Equal(t, "example-alias-updated", d.Get("alias_zone_name"))
}

func TestTeoZoneV2_Update_NoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)

	// Mock Modify API, should NOT be called
	patches.ApplyMethodFunc(teoClient, "ModifyZoneWithContext", func(_ context.Context, request *teov20220901.ModifyZoneRequest) (*teov20220901.ModifyZoneResponse, error) {
		t.Errorf("ModifyZoneWithContext should not be called when no mutable arg changed")
		resp := teov20220901.NewModifyZoneResponse()
		return resp, nil
	})

	// Mock Read API (called after Update)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext", func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64ZoneV2(1),
			Zones: []*teov20220901.Zone{
				{
					ZoneId:   ptrStringZoneV2("zone-2noz78a8ev6k"),
					ZoneName: ptrStringZoneV2("example.com"),
					Area:     ptrStringZoneV2("global"),
					Type:     ptrStringZoneV2("full"),
					Status:   ptrStringZoneV2("active"),
				},
			},
			RequestId: ptrStringZoneV2("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZoneV2()
	res := teo.ResourceTencentCloudTeoZoneV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"type":      "full",
		"zone_name": "example.com",
		"area":      "global",
	})
	d.SetId("zone-2noz78a8ev6k")

	// Force all args to be detected as unchanged (TestResourceDataRaw has no old state)
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2noz78a8ev6k", d.Id())
	assert.Equal(t, "global", d.Get("area"))
}

func TestTeoZoneV2_Delete_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoZoneV2().client, "UseTeoV20220901Client", teoClient)

	// Mock Delete API
	patches.ApplyMethodFunc(teoClient, "DeleteZoneWithContext", func(_ context.Context, request *teov20220901.DeleteZoneRequest) (*teov20220901.DeleteZoneResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-2noz78a8ev6k", *request.ZoneId)

		resp := teov20220901.NewDeleteZoneResponse()
		resp.Response = &teov20220901.DeleteZoneResponseParams{
			RequestId: ptrStringZoneV2("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoZoneV2()
	res := teo.ResourceTencentCloudTeoZoneV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
	})
	d.SetId("zone-2noz78a8ev6k")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
