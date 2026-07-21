package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoZoneTagParams" -v -count=1 -gcflags="all=-l"

// TestTeoZoneTagParams_Schema validates that resource_region and service_type are in the schema
func TestTeoZoneTagParams_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoZone()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "resource_region")
	assert.Contains(t, res.Schema, "service_type")

	resourceRegion := res.Schema["resource_region"]
	assert.Equal(t, schema.TypeString, resourceRegion.Type)
	assert.True(t, resourceRegion.Optional)
	assert.True(t, resourceRegion.Computed)

	serviceType := res.Schema["service_type"]
	assert.Equal(t, schema.TypeString, serviceType.Type)
	assert.True(t, serviceType.Optional)
	assert.True(t, serviceType.Computed)
}

// TestTeoZoneTagParams_Read_DefaultValues tests that Read uses default values when resource_region and service_type are not set
func TestTeoZoneTagParams_Read_DefaultValues(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	_ = newTeoClientForZoneRead(patches, meta, "zone-12345678", "example.com", "partial", "overseas", "active")

	var capturedServiceType string
	var capturedRegion string
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		capturedServiceType = serviceType
		capturedRegion = region
		return map[string]string{"key1": "value1"}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
		"type":      "partial",
		"area":      "overseas",
		"plan_id":   "edgeone-2kfv1h391n6w",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify default values are used: serviceType="teo", region=provider region
	assert.Equal(t, "teo", capturedServiceType)
	assert.Equal(t, "ap-guangzhou", capturedRegion)

	// Verify resource_region and service_type are set in state with default values
	assert.Equal(t, "ap-guangzhou", d.Get("resource_region"))
	assert.Equal(t, "teo", d.Get("service_type"))
}

// TestTeoZoneTagParams_Read_CustomValues tests that Read uses custom resource_region and service_type values
func TestTeoZoneTagParams_Read_CustomValues(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	_ = newTeoClientForZoneRead(patches, meta, "zone-12345678", "example.com", "partial", "overseas", "active")

	var capturedServiceType string
	var capturedRegion string
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		capturedServiceType = serviceType
		capturedRegion = region
		return map[string]string{"key1": "value1"}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name":       "example.com",
		"type":            "partial",
		"area":            "overseas",
		"plan_id":         "edgeone-2kfv1h391n6w",
		"resource_region": "ap-beijing",
		"service_type":    "edgeone",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify custom values are used
	assert.Equal(t, "edgeone", capturedServiceType)
	assert.Equal(t, "ap-beijing", capturedRegion)

	// Verify resource_region and service_type are set in state with the custom values
	assert.Equal(t, "ap-beijing", d.Get("resource_region"))
	assert.Equal(t, "edgeone", d.Get("service_type"))
}

// TestTeoZoneTagParams_Create_UsesTagParams tests that Create uses resource_region and service_type in tag ModifyTags
func TestTeoZoneTagParams_Create_UsesTagParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	// Mock UseTeoClient for CreateZoneWithContext and DescribeZones
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	// Mock CreateZoneWithContext
	patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext", func(_ context.Context, request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
		resp := teov20220901.NewCreateZoneResponse()
		resp.Response = &teov20220901.CreateZoneResponseParams{
			ZoneId:    ptrString("zone-12345678"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeZones to return "pending" for the post-handle retry check,
	// and a complete zone for the read-after-create
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrString("zone-12345678"),
					ZoneName:      ptrString("example.com"),
					Status:        ptrString("pending"),
					Type:          ptrString("partial"),
					Paused:        ptrBool(false),
					Area:          ptrString("overseas"),
					ActiveStatus:  ptrString("active"),
					NameServers:   []*string{ptrString("ns1.example.com"), ptrString("ns2.example.com")},
					AliasZoneName: ptrString("alias-test"),
					Resources: []*teov20220901.Resource{
						{
							PlanId: ptrString("edgeone-2kfv1h391n6w"),
						},
					},
					OwnershipVerification: &teov20220901.OwnershipVerification{
						DnsVerification: &teov20220901.DnsVerification{
							Subdomain:   ptrString("_verify1"),
							RecordType:  ptrString("TXT"),
							RecordValue: ptrString("verify-value"),
						},
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock ModifyTags to capture resourceName (which contains serviceType and region)
	var capturedResourceName string
	patches.ApplyMethodFunc(&svctag.TagService{}, "ModifyTags", func(ctx context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		capturedResourceName = resourceName
		return nil
	})

	// Mock DescribeResourceTags for the read-after-create
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{"DoNotMove": "TF-Test"}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name":       "example.com",
		"type":            "partial",
		"area":            "overseas",
		"plan_id":         "edgeone-2kfv1h391n6w",
		"tags":            map[string]interface{}{"DoNotMove": "TF-Test"},
		"resource_region": "ap-beijing",
		"service_type":    "edgeone",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify ModifyTags was called with correct QCS resource name that includes custom service_type and region
	// QCS format: qcs::service_type:region:uin/:resourceType/id
	assert.Contains(t, capturedResourceName, "edgeone")
	assert.Contains(t, capturedResourceName, "ap-beijing")
	assert.Contains(t, capturedResourceName, "zone")
	assert.Contains(t, capturedResourceName, d.Id())
}

// TestTeoZoneTagParams_Create_DefaultTagParams tests that Create uses defaults when resource_region and service_type are not set
func TestTeoZoneTagParams_Create_DefaultTagParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext", func(_ context.Context, request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
		resp := teov20220901.NewCreateZoneResponse()
		resp.Response = &teov20220901.CreateZoneResponseParams{
			ZoneId:    ptrString("zone-12345678"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrString("zone-12345678"),
					ZoneName:      ptrString("example.com"),
					Status:        ptrString("pending"),
					Type:          ptrString("partial"),
					Paused:        ptrBool(false),
					Area:          ptrString("overseas"),
					ActiveStatus:  ptrString("active"),
					NameServers:   []*string{ptrString("ns1.example.com"), ptrString("ns2.example.com")},
					AliasZoneName: ptrString("alias-test"),
					Resources: []*teov20220901.Resource{
						{
							PlanId: ptrString("edgeone-2kfv1h391n6w"),
						},
					},
					OwnershipVerification: &teov20220901.OwnershipVerification{
						DnsVerification: &teov20220901.DnsVerification{
							Subdomain:   ptrString("_verify1"),
							RecordType:  ptrString("TXT"),
							RecordValue: ptrString("verify-value"),
						},
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	var capturedResourceName string
	patches.ApplyMethodFunc(&svctag.TagService{}, "ModifyTags", func(ctx context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		capturedResourceName = resourceName
		return nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{"DoNotMove": "TF-Test"}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
		"type":      "partial",
		"area":      "overseas",
		"plan_id":   "edgeone-2kfv1h391n6w",
		"tags":      map[string]interface{}{"DoNotMove": "TF-Test"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify ModifyTags used default values: serviceType="teo", region=provider region
	assert.Contains(t, capturedResourceName, "teo")
	assert.Contains(t, capturedResourceName, "ap-guangzhou")
}

// TestTeoZoneTagParams_Update_UsesTagParams tests that Update uses resource_region and service_type in ModifyTags
func TestTeoZoneTagParams_Update_UsesTagParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	_ = newTeoClientForZoneRead(patches, meta, "zone-12345678", "example.com", "partial", "overseas", "active")

	var capturedResourceName string
	patches.ApplyMethodFunc(&svctag.TagService{}, "ModifyTags", func(ctx context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		capturedResourceName = resourceName
		return nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{"key1": "value1"}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name":       "example.com",
		"type":            "partial",
		"area":            "overseas",
		"plan_id":         "edgeone-2kfv1h391n6w",
		"resource_region": "ap-beijing",
		"service_type":    "edgeone",
		"tags":            map[string]interface{}{"key1": "value1"},
	})
	d.SetId("zone-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify ModifyTags was called with QCS resource name containing custom values
	assert.Contains(t, capturedResourceName, "edgeone")
	assert.Contains(t, capturedResourceName, "ap-beijing")
}

// TestTeoZoneTagParams_Update_DefaultTagParams tests that Update uses default values when resource_region and service_type are not set
func TestTeoZoneTagParams_Update_DefaultTagParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	_ = newTeoClientForZoneRead(patches, meta, "zone-12345678", "example.com", "partial", "overseas", "active")

	var capturedResourceName string
	patches.ApplyMethodFunc(&svctag.TagService{}, "ModifyTags", func(ctx context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		capturedResourceName = resourceName
		return nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{"key1": "value1"}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
		"type":      "partial",
		"area":      "overseas",
		"plan_id":   "edgeone-2kfv1h391n6w",
		"tags":      map[string]interface{}{"key1": "value1"},
	})
	d.SetId("zone-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify ModifyTags used default values: serviceType="teo", region=provider region
	assert.Contains(t, capturedResourceName, "teo")
	assert.Contains(t, capturedResourceName, "ap-guangzhou")
}

// newTeoClientForZoneRead creates a mocked TEO client that returns a valid zone for DescribeZones
// and mocks other zone API methods needed for update flows
func newTeoClientForZoneRead(patches *gomonkey.Patches, meta *mockMeta, zoneId, zoneName, zoneType, zoneArea, zoneStatus string) *teov20220901.Client {
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrString(zoneId),
					ZoneName:      ptrString(zoneName),
					Status:        ptrString(zoneStatus),
					Type:          ptrString(zoneType),
					Paused:        ptrBool(false),
					Area:          ptrString(zoneArea),
					ActiveStatus:  ptrString("active"),
					NameServers:   []*string{ptrString("ns1.example.com"), ptrString("ns2.example.com")},
					AliasZoneName: ptrString("alias-test"),
					Resources: []*teov20220901.Resource{
						{
							PlanId: ptrString("edgeone-2kfv1h391n6w"),
						},
					},
					OwnershipVerification: &teov20220901.OwnershipVerification{
						DnsVerification: &teov20220901.DnsVerification{
							Subdomain:   ptrString("_verify1"),
							RecordType:  ptrString("TXT"),
							RecordValue: ptrString("verify-value"),
						},
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock ModifyZoneWithContext for update flows
	patches.ApplyMethodFunc(teoClient, "ModifyZoneWithContext", func(_ context.Context, request *teov20220901.ModifyZoneRequest) (*teov20220901.ModifyZoneResponse, error) {
		resp := teov20220901.NewModifyZoneResponse()
		resp.Response = &teov20220901.ModifyZoneResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock ModifyZoneStatusWithContext for update flows
	patches.ApplyMethodFunc(teoClient, "ModifyZoneStatusWithContext", func(_ context.Context, request *teov20220901.ModifyZoneStatusRequest) (*teov20220901.ModifyZoneStatusResponse, error) {
		resp := teov20220901.NewModifyZoneStatusResponse()
		resp.Response = &teov20220901.ModifyZoneStatusResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	return teoClient
}

func ptrBool(b bool) *bool {
	return &b
}
