package teo

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoZone1110_02" -v -count=1 -gcflags="all=-l"

func teoZone1110_02Ptr(s string) *string { return &s }

// TestTeoZone1110_02Resource_Metadata validates the resource type name.
func TestTeoZone1110_02Resource_Metadata(t *testing.T) {
	r := &TeoZone1110_02Resource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "tencentcloud_teo_zone_1110_02", resp.TypeName)
}

// TestTeoZone1110_02Resource_Schema validates key attribute flags.
func TestTeoZone1110_02Resource_Schema(t *testing.T) {
	r := &TeoZone1110_02Resource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	type flag struct {
		required, optional, computed bool
	}
	cases := map[string]flag{
		"id":              {computed: true},
		"type":            {optional: true, computed: true},
		"zone_name":       {optional: true, computed: true},
		"area":            {optional: true, computed: true},
		"plan_id":         {optional: true},
		"alias_zone_name": {optional: true, computed: true},
		"tags":            {optional: true},
		"allow_duplicates": {optional: true},
		"jump_start":       {optional: true},
		"vanity_name_servers": {optional: true, computed: true},
		"ownership_verification": {computed: true},
		"zones":                  {computed: true},
	}
	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing in schema", name)
			continue
		}
		assert.Equal(t, want.required, attr.IsRequired(), "attribute %q IsRequired", name)
		assert.Equal(t, want.optional, attr.IsOptional(), "attribute %q IsOptional", name)
		assert.Equal(t, want.computed, attr.IsComputed(), "attribute %q IsComputed", name)
	}

	// nested attributes of vanity_name_servers: switch required, servers optional
	vanityAttr, ok := resp.Schema.Attributes["vanity_name_servers"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("vanity_name_servers is not a SingleNestedAttribute")
	}
	switchAttr, ok := vanityAttr.Attributes["switch"].(schema.StringAttribute)
	if !ok || !switchAttr.Required {
		t.Errorf("vanity_name_servers.switch missing or not required: %+v", vanityAttr.Attributes["switch"])
	}
	serversAttr, ok := vanityAttr.Attributes["servers"].(schema.SetAttribute)
	if !ok || !serversAttr.Optional {
		t.Errorf("vanity_name_servers.servers missing or not optional: %+v", vanityAttr.Attributes["servers"])
	}
}

// TestTeoZone1110_02Resource_ImplementsInterfaces verifies advertised interfaces.
func TestTeoZone1110_02Resource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &TeoZone1110_02Resource{}
	var _ resource.ResourceWithImportState = &TeoZone1110_02Resource{}
}

// TestTeoZone1110_02Resource_Factory validates the registry factory.
func TestTeoZone1110_02Resource_Factory(t *testing.T) {
	got := NewTeoZone1110_02Resource()
	assert.NotNil(t, got)
	_, ok := got.(*TeoZone1110_02Resource)
	assert.True(t, ok)
}

// TestTeoZone1110_02Resource_ConfigureNilProviderData ensures nil ProviderData
// is a silent no-op.
func TestTeoZone1110_02Resource_ConfigureNilProviderData(t *testing.T) {
	r := &TeoZone1110_02Resource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}

// TestBuildTeoZone1110_02CreateRequest_FullFields validates the request
// mapping when every create field is set.
func TestBuildTeoZone1110_02CreateRequest_FullFields(t *testing.T) {
	plan := &TeoZone1110_02ResourceModel{
		Type:            types.StringValue("full"),
		ZoneName:        types.StringValue("example.com"),
		Area:            types.StringValue("overseas"),
		PlanId:          types.StringValue("plan-xxx"),
		AliasZoneName:   types.StringValue("alias-1"),
		Tags:            []teoZoneTagModel{{TagKey: types.StringValue("k1"), TagValue: types.StringValue("v1")}},
		AllowDuplicates: types.BoolValue(true),
		JumpStart:       types.BoolValue(true),
	}

	request := buildTeoZone1110_02CreateRequest(plan)

	assert.NotNil(t, request)
	assert.Equal(t, "full", *request.Type)
	assert.Equal(t, "example.com", *request.ZoneName)
	assert.Equal(t, "overseas", *request.Area)
	assert.Equal(t, "plan-xxx", *request.PlanId)
	assert.Equal(t, "alias-1", *request.AliasZoneName)
	assert.True(t, *request.AllowDuplicates)
	assert.True(t, *request.JumpStart)
	if len(request.Tags) != 1 {
		t.Fatalf("len(request.Tags) = %d, want 1", len(request.Tags))
	}
	assert.Equal(t, "k1", *request.Tags[0].TagKey)
	assert.Equal(t, "v1", *request.Tags[0].TagValue)
}

// TestBuildTeoZone1110_02CreateRequest_NullFields validates that null fields
// are not set on the request.
func TestBuildTeoZone1110_02CreateRequest_NullFields(t *testing.T) {
	plan := &TeoZone1110_02ResourceModel{
		Type:            types.StringNull(),
		ZoneName:        types.StringNull(),
		Area:            types.StringNull(),
		PlanId:          types.StringNull(),
		AliasZoneName:   types.StringNull(),
		AllowDuplicates: types.BoolNull(),
		JumpStart:       types.BoolNull(),
	}

	request := buildTeoZone1110_02CreateRequest(plan)

	assert.Nil(t, request.Type)
	assert.Nil(t, request.ZoneName)
	assert.Nil(t, request.Area)
	assert.Nil(t, request.PlanId)
	assert.Nil(t, request.AliasZoneName)
	assert.Nil(t, request.AllowDuplicates)
	assert.Nil(t, request.JumpStart)
	assert.Nil(t, request.Tags)
}

// TestTeoZone1110_02Resource_ReadZone_Gone validates the zone-id filter query
// and the "does not exist" signal.
func TestTeoZone1110_02Resource_ReadZone_Gone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(ctx context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			if len(request.Filters) != 1 || *request.Filters[0].Name != "zone-id" {
				return nil, fmt.Errorf("unexpected filters: %+v", request.Filters)
			}
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{TotalCount: teoZone1110_02PtrInt64(0)}
			return resp, nil
		})

	r := &TeoZone1110_02Resource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: mockClient},
	}, &resource.ConfigureResponse{})

	zone, err := r.readZone(context.Background(), "zone-missing")
	assert.NoError(t, err)
	assert.Nil(t, zone)
}

// TestTeoZone1110_02Resource_ReadZone_Success validates the successful query.
func TestTeoZone1110_02Resource_ReadZone_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(ctx context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				TotalCount: teoZone1110_02PtrInt64(1),
				Zones: []*teov20220901.Zone{
					{ZoneId: teoZone1110_02Ptr("zone-abc"), ZoneName: teoZone1110_02Ptr("example.com")},
				},
			}
			return resp, nil
		})

	r := &TeoZone1110_02Resource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: mockClient},
	}, &resource.ConfigureResponse{})

	zone, err := r.readZone(context.Background(), "zone-abc")
	assert.NoError(t, err)
	assert.NotNil(t, zone)
	if zone != nil {
		assert.Equal(t, "zone-abc", *zone.ZoneId)
		assert.Equal(t, "example.com", *zone.ZoneName)
	}
}

// TestTeoZone1110_02Resource_ReadZone_APIError validates API error propagation.
func TestTeoZone1110_02Resource_ReadZone_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(ctx context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=boom")
		})

	r := &TeoZone1110_02Resource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: mockClient},
	}, &resource.ConfigureResponse{})

	zone, err := r.readZone(context.Background(), "zone-abc")
	assert.Error(t, err)
	assert.Nil(t, zone)
	assert.Contains(t, err.Error(), "InternalError")
}

// TestTeoZone1110_02Resource_ReadZone_NilResponse validates that a nil
// response is treated as an API failure instead of "not found".
func TestTeoZone1110_02Resource_ReadZone_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(ctx context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			return nil, nil
		})

	r := &TeoZone1110_02Resource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: mockClient},
	}, &resource.ConfigureResponse{})

	zone, err := r.readZone(context.Background(), "zone-abc")
	assert.Error(t, err)
	assert.Nil(t, zone)
}

// TestTeoZone1110_02Resource_DeleteZone_Success validates delete invocation.
func TestTeoZone1110_02Resource_DeleteZone_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DeleteZoneWithContext",
		func(ctx context.Context, request *teov20220901.DeleteZoneRequest) (*teov20220901.DeleteZoneResponse, error) {
			assert.Equal(t, "zone-abc", *request.ZoneId)
			resp := teov20220901.NewDeleteZoneResponse()
			resp.Response = &teov20220901.DeleteZoneResponseParams{RequestId: teoZone1110_02Ptr("req-1")}
			return resp, nil
		})

	r := &TeoZone1110_02Resource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: mockClient},
	}, &resource.ConfigureResponse{})

	assert.NoError(t, r.deleteZone(context.Background(), "zone-abc"))
}

// TestTeoZone1110_02Resource_DeleteZone_APIError validates delete error propagation.
func TestTeoZone1110_02Resource_DeleteZone_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DeleteZoneWithContext",
		func(ctx context.Context, request *teov20220901.DeleteZoneRequest) (*teov20220901.DeleteZoneResponse, error) {
			return nil, fmt.Errorf("[TencentCloudSDKError] Code=OperationDenied, Message=zone is locked")
		})

	r := &TeoZone1110_02Resource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: mockClient},
	}, &resource.ConfigureResponse{})

	err := r.deleteZone(context.Background(), "zone-abc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OperationDenied")
}

// TestRefreshTeoZone1110_02Model validates the model refresh against a fully
// populated cloud zone object, including every leaf field and the zones list.
func TestRefreshTeoZone1110_02Model(t *testing.T) {
	ctx := context.Background()
	zone := &teov20220901.Zone{
		ZoneId:        teoZone1110_02Ptr("zone-abc"),
		ZoneName:      teoZone1110_02Ptr("example.com"),
		AliasZoneName: teoZone1110_02Ptr("alias-1"),
		Area:          teoZone1110_02Ptr("overseas"),
		Type:          teoZone1110_02Ptr("full"),
		Tags: []*teov20220901.Tag{
			{TagKey: teoZone1110_02Ptr("k1"), TagValue: teoZone1110_02Ptr("v1")},
		},
		Resources: []*teov20220901.Resource{
			{
				Id:            teoZone1110_02Ptr("res-1"),
				PayMode:       teoZone1110_02PtrInt64(0),
				CreateTime:    teoZone1110_02Ptr("2024-01-01T00:00:00Z"),
				EnableTime:    teoZone1110_02Ptr("2024-01-02T00:00:00Z"),
				ExpireTime:    teoZone1110_02Ptr("2025-01-01T00:00:00Z"),
				Status:        teoZone1110_02Ptr("normal"),
				Sv: []*teov20220901.Sv{
					{
						Key:       teoZone1110_02Ptr("zone"),
						Value:     teoZone1110_02Ptr("1"),
						Pack:      teoZone1110_02Ptr("zone"),
						ProtectionSpecs: teoZone1110_02Ptr("cm_30G"),
					},
				},
				AutoRenewFlag: teoZone1110_02PtrInt64(1),
				PlanId:        teoZone1110_02Ptr("plan-1"),
				Area:          teoZone1110_02Ptr("global"),
				Group:         teoZone1110_02Ptr("plan"),
				ZoneNumber:    teoZone1110_02PtrInt64(2),
				Type:          teoZone1110_02Ptr("vodeo"),
			},
		},
		NSDetail: &teov20220901.NSDetail{
			CnameSpeedUp:         teoZone1110_02Ptr("enabled"),
			IsFake:               teoZone1110_02PtrInt64(0),
			OriginalNameServers:  []*string{teoZone1110_02Ptr("ns1.old.com")},
			NameServers:          []*string{teoZone1110_02Ptr("ns1.eo.com")},
			VanityNameServers: &teov20220901.VanityNameServers{
				Switch:  teoZone1110_02Ptr("on"),
				Servers: []*string{teoZone1110_02Ptr("ns1.custom.com")},
			},
			VanityNameServersIps: []*teov20220901.VanityNameServersIps{
				{Name: teoZone1110_02Ptr("ns1.custom.com"), IPv4: teoZone1110_02Ptr("1.2.3.4")},
			},
			OwnershipVerification: &teov20220901.OwnershipVerification{
				NsVerification: &teov20220901.NsVerification{
					NameServers: []*string{teoZone1110_02Ptr("ns1.eo.com")},
				},
			},
		},
		CreatedOn:    teoZone1110_02Ptr("2024-01-01T00:00:00Z"),
		ModifiedOn:   teoZone1110_02Ptr("2024-01-03T00:00:00Z"),
		Status:       teoZone1110_02Ptr("pending"),
		CnameStatus:  teoZone1110_02Ptr("pending"),
		ActiveStatus: teoZone1110_02Ptr("inactive"),
		LockStatus:   teoZone1110_02Ptr("enable"),
		Paused:       teoZone1110_02PtrBool(false),
		IsFake:       teoZone1110_02PtrInt64(0),
		CnameSpeedUp: teoZone1110_02Ptr("enabled"),
		OwnershipVerification: &teov20220901.OwnershipVerification{
			DnsVerification: &teov20220901.DnsVerification{
				Subdomain:   teoZone1110_02Ptr("_teo"),
				RecordType:  teoZone1110_02Ptr("TXT"),
				RecordValue: teoZone1110_02Ptr("verify-me"),
			},
			FileVerification: &teov20220901.FileVerification{
				Path:    teoZone1110_02Ptr("/.well-known/teo-verification/z.txt"),
				Content: teoZone1110_02Ptr("content"),
			},
		},
		VanityNameServers: &teov20220901.VanityNameServers{
			Switch:  teoZone1110_02Ptr("off"),
			Servers: []*string{teoZone1110_02Ptr("ns1.custom.com")},
		},
		VanityNameServersIps: []*teov20220901.VanityNameServersIps{
			{Name: teoZone1110_02Ptr("ns1.custom.com"), IPv4: teoZone1110_02Ptr("5.6.7.8")},
		},
		WorkModeInfos: []*teov20220901.ConfigGroupWorkModeInfo{
			{ConfigGroupType: teoZone1110_02Ptr("l7_acceleration"), WorkMode: teoZone1110_02Ptr("version_control")},
		},
	}

	model := &TeoZone1110_02ResourceModel{
		Id:              types.StringNull(),
		Type:            types.StringValue("full"),
		ZoneName:        types.StringValue("example.com"),
		Area:            types.StringNull(),
		PlanId:          types.StringValue("plan-keep"),
		AllowDuplicates: types.BoolValue(true),
		JumpStart:       types.BoolNull(),
	}

	err := refreshTeoZone1110_02Model(ctx, model, zone)
	assert.NoError(t, err)

	assert.Equal(t, "full", model.Type.ValueString())
	assert.Equal(t, "example.com", model.ZoneName.ValueString())
	assert.Equal(t, "overseas", model.Area.ValueString())
	assert.Equal(t, "alias-1", model.AliasZoneName.ValueString())

	// create-only fields must be preserved (no read-back from DescribeZones)
	assert.Equal(t, "plan-keep", model.PlanId.ValueString())
	assert.True(t, model.AllowDuplicates.ValueBool())

	// ownership_verification object
	assert.False(t, model.OwnershipVerification.IsNull())
	ovAttrs, diag := model.OwnershipVerification.ToObjectValue(ctx)
	assert.False(t, diag.HasError())
	dnsObj := ovAttrs.AttributeValues()["dns_verification"].(types.Object)
	assert.False(t, dnsObj.IsNull())

	// zones list: one element, validate leaf fields
	assert.False(t, model.Zones.IsNull())
	assert.Equal(t, 1, len(model.Zones.Elements()))
	var details []teoZoneDetailModel
	d := model.Zones.ElementsAs(ctx, &details, false)
	assert.False(t, d.HasError())
	detail := details[0]
	assert.Equal(t, "zone-abc", detail.ZoneId.ValueString())
	assert.Equal(t, "example.com", detail.ZoneName.ValueString())
	assert.Equal(t, "alias-1", detail.AliasZoneName.ValueString())
	assert.Equal(t, "overseas", detail.Area.ValueString())
	assert.Equal(t, "full", detail.Type.ValueString())
	assert.Equal(t, "pending", detail.Status.ValueString())
	assert.Equal(t, "pending", detail.CnameStatus.ValueString())
	assert.Equal(t, "inactive", detail.ActiveStatus.ValueString())
	assert.Equal(t, "enable", detail.LockStatus.ValueString())
	assert.False(t, detail.Paused.ValueBool())
	assert.Equal(t, int64(0), detail.IsFake.ValueInt64())
	assert.Equal(t, "enabled", detail.CnameSpeedUp.ValueString())

	// tags
	if len(detail.Tags) != 1 {
		t.Fatalf("len(detail.Tags) = %d, want 1", len(detail.Tags))
	}
	assert.Equal(t, "k1", detail.Tags[0].TagKey.ValueString())
	assert.Equal(t, "v1", detail.Tags[0].TagValue.ValueString())

	// resources
	if len(detail.Resources) != 1 {
		t.Fatalf("len(detail.Resources) = %d, want 1", len(detail.Resources))
	}
	res := detail.Resources[0]
	assert.Equal(t, "res-1", res.Id.ValueString())
	assert.Equal(t, int64(0), res.PayMode.ValueInt64())
	assert.Equal(t, int64(1), res.AutoRenewFlag.ValueInt64())
	assert.Equal(t, "plan-1", res.PlanId.ValueString())
	assert.Equal(t, "global", res.Area.ValueString())
	assert.Equal(t, "plan", res.Group.ValueString())
	assert.Equal(t, int64(2), res.ZoneNumber.ValueInt64())
	assert.Equal(t, "vodeo", res.Type.ValueString())
	if len(res.Sv) != 1 {
		t.Fatalf("len(res.Sv) = %d, want 1", len(res.Sv))
	}
	assert.Equal(t, "zone", res.Sv[0].Key.ValueString())
	assert.Equal(t, "1", res.Sv[0].Value.ValueString())
	assert.Equal(t, "zone", res.Sv[0].Pack.ValueString())
	assert.True(t, res.Sv[0].InstanceId.IsNull())
	assert.Equal(t, "cm_30G", res.Sv[0].ProtectionSpecs.ValueString())

	// NS detail
	assert.NotNil(t, detail.NSDetail)
	assert.Equal(t, "enabled", detail.NSDetail.CnameSpeedUp.ValueString())
	assert.Equal(t, int64(0), detail.NSDetail.IsFake.ValueInt64())
	var originalServers []string
	d = detail.NSDetail.OriginalNameServers.ElementsAs(ctx, &originalServers, false)
	assert.False(t, d.HasError())
	assert.Equal(t, []string{"ns1.old.com"}, originalServers)
	var nameServers []string
	d = detail.NSDetail.NameServers.ElementsAs(ctx, &nameServers, false)
	assert.False(t, d.HasError())
	assert.Equal(t, []string{"ns1.eo.com"}, nameServers)

	// vanity detail object inside ns_detail
	vanityObj := detail.NSDetail.VanityNameServers
	assert.False(t, vanityObj.IsNull())
	vanityVals, vdiag := vanityObj.ToObjectValue(ctx)
	assert.False(t, vdiag.HasError())
	assert.Equal(t, "on", vanityVals.AttributeValues()["switch"].(types.String).ValueString())
	vanityServers := vanityVals.AttributeValues()["servers"].(types.Set)
	var servers []string
	d = vanityServers.ElementsAs(ctx, &servers, false)
	assert.False(t, d.HasError())
	assert.Equal(t, []string{"ns1.custom.com"}, servers)

	// ns_detail vanity ips
	if len(detail.NSDetail.VanityNameServersIps) != 1 {
		t.Fatalf("len(detail.NSDetail.VanityNameServersIps) = %d, want 1", len(detail.NSDetail.VanityNameServersIps))
	}
	assert.Equal(t, "ns1.custom.com", detail.NSDetail.VanityNameServersIps[0].Name.ValueString())
	assert.Equal(t, "1.2.3.4", detail.NSDetail.VanityNameServersIps[0].IPv4.ValueString())

	// cname / dnspod detail absent for full type
	assert.Nil(t, detail.CNAMEDetail)
	assert.Nil(t, detail.DNSPodDetail)

	// legacy vanity / ips / work modes
	assert.False(t, detail.VanityNameServers.IsNull())
	if len(detail.VanityNameServersIps) != 1 {
		t.Fatalf("len(detail.VanityNameServersIps) = %d, want 1", len(detail.VanityNameServersIps))
	}
	assert.Equal(t, "5.6.7.8", detail.VanityNameServersIps[0].IPv4.ValueString())
	if len(detail.WorkModeInfos) != 1 {
		t.Fatalf("len(detail.WorkModeInfos) = %d, want 1", len(detail.WorkModeInfos))
	}
	assert.Equal(t, "l7_acceleration", detail.WorkModeInfos[0].ConfigGroupType.ValueString())
	assert.Equal(t, "version_control", detail.WorkModeInfos[0].WorkMode.ValueString())

	// legacy sets
	var legacyOriginal []string
	d = detail.OriginalNameServers.ElementsAs(ctx, &legacyOriginal, false)
	assert.False(t, d.HasError())
	assert.Equal(t, []string{"ns1.old.com"}, legacyOriginal)
	var legacyNameServers []string
	d = detail.NameServers.ElementsAs(ctx, &legacyNameServers, false)
	assert.False(t, d.HasError())
	assert.Equal(t, []string{"ns1.eo.com"}, legacyNameServers)

	// legacy ownership verification with dns + file verification
	assert.False(t, detail.OwnershipVerification.IsNull())
}

// TestRefreshTeoZone1110_02Model_MinimalZone validates null-safe refresh for a
// sparse zone object (no details, no tags, no servers).
func TestRefreshTeoZone1110_02Model_MinimalZone(t *testing.T) {
	ctx := context.Background()
	zone := &teov20220901.Zone{
		ZoneId: teoZone1110_02Ptr("zone-min"),
		Type:   teoZone1110_02Ptr("partial"),
	}

	model := &TeoZone1110_02ResourceModel{}
	err := refreshTeoZone1110_02Model(ctx, model, zone)
	assert.NoError(t, err)

	assert.Equal(t, "partial", model.Type.ValueString())
	assert.True(t, model.OwnershipVerification.IsNull())
	assert.False(t, model.Zones.IsNull())

	var details []teoZoneDetailModel
	d := model.Zones.ElementsAs(ctx, &details, false)
	assert.False(t, d.HasError())
	assert.Nil(t, details[0].NSDetail)
	assert.Nil(t, details[0].CNAMEDetail)
	assert.Nil(t, details[0].DNSPodDetail)
	assert.True(t, details[0].OriginalNameServers.IsNull())
	assert.True(t, details[0].NameServers.IsNull())
}

// TestTeoZoneStringSetFromSdk validates empty/nil maps to null set and values
// round-trip.
func TestTeoZoneStringSetFromSdk(t *testing.T) {
	ctx := context.Background()

	nullSet, err := teoZoneStringSetFromSdk(ctx, nil)
	assert.NoError(t, err)
	assert.True(t, nullSet.IsNull())

	emptySet, err := teoZoneStringSetFromSdk(ctx, []*string{})
	assert.NoError(t, err)
	assert.True(t, emptySet.IsNull())

	set, err := teoZoneStringSetFromSdk(ctx, []*string{teoZone1110_02Ptr("a"), teoZone1110_02Ptr("b")})
	assert.NoError(t, err)
	assert.False(t, set.IsNull())
	var values []string
	d := set.ElementsAs(ctx, &values, false)
	assert.False(t, d.HasError())
	assert.Len(t, values, 2)
}

// TestTeoZoneOwnershipObjectFromSdk_Nil validates nil maps to null object.
func TestTeoZoneOwnershipObjectFromSdk_Nil(t *testing.T) {
	obj, err := teoZoneOwnershipObjectFromSdk(context.Background(), nil)
	assert.NoError(t, err)
	assert.True(t, obj.IsNull())
}

func teoZone1110_02PtrInt64(i int64) *int64   { return &i }
func teoZone1110_02PtrBool(b bool) *bool      { return &b }
