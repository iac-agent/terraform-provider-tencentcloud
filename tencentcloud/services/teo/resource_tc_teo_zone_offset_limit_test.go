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

// go test ./tencentcloud/services/teo/ -run "TestTeoZoneOffsetLimit" -v -count=1 -gcflags="all=-l"

// TestTeoZoneOffsetLimit_Schema validates that offset and limit are in the schema
func TestTeoZoneOffsetLimit_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoZone()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "offset")
	assert.Contains(t, res.Schema, "limit")

	offset := res.Schema["offset"]
	assert.Equal(t, schema.TypeInt, offset.Type)
	assert.True(t, offset.Optional)

	limit := res.Schema["limit"]
	assert.Equal(t, schema.TypeInt, limit.Type)
	assert.True(t, limit.Optional)
}

// TestTeoZoneOffsetLimit_Read_DefaultValues tests that Read uses default offset/limit (0, 0) when not set
func TestTeoZoneOffsetLimit_Read_DefaultValues(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	_ = newTeoClientForZoneRead(patches, meta, "zone-12345678", "example.com", "partial", "overseas", "active")

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
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

	// Verify the zone was read correctly
	assert.Equal(t, "zone-12345678", d.Id())
}

// TestTeoZoneOffsetLimit_Read_WithOffsetLimit tests that Read passes offset/limit to DescribeZones
func TestTeoZoneOffsetLimit_Read_WithOffsetLimit(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	var capturedOffset *int64
	var capturedLimit *int64
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		capturedOffset = request.Offset
		capturedLimit = request.Limit
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrString("zone-12345678"),
					ZoneName:      ptrString("example.com"),
					Status:        ptrString("active"),
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

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
		"type":      "partial",
		"area":      "overseas",
		"plan_id":   "edgeone-2kfv1h391n6w",
		"offset":    10,
		"limit":     50,
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify the offset and limit were passed to DescribeZones
	assert.NotNil(t, capturedOffset)
	assert.NotNil(t, capturedLimit)
	assert.Equal(t, int64(10), *capturedOffset)
	assert.Equal(t, int64(50), *capturedLimit)
}

// TestTeoZoneOffsetLimit_Read_OnlyOffset tests that Read passes only offset when limit is not set
func TestTeoZoneOffsetLimit_Read_OnlyOffset(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	var capturedOffset *int64
	var capturedLimit *int64
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		capturedOffset = request.Offset
		capturedLimit = request.Limit
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrString("zone-12345678"),
					ZoneName:      ptrString("example.com"),
					Status:        ptrString("active"),
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

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
		"type":      "partial",
		"area":      "overseas",
		"plan_id":   "edgeone-2kfv1h391n6w",
		"offset":    5,
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify offset was passed but limit uses default (20)
	assert.NotNil(t, capturedOffset)
	assert.NotNil(t, capturedLimit)
	assert.Equal(t, int64(5), *capturedOffset)
	assert.Equal(t, int64(20), *capturedLimit)
}

// TestTeoZoneOffsetLimit_Read_OnlyLimit tests that Read passes only limit when offset is not set
func TestTeoZoneOffsetLimit_Read_OnlyLimit(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	meta.client.Region = "ap-guangzhou"

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	var capturedOffset *int64
	var capturedLimit *int64
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		capturedOffset = request.Offset
		capturedLimit = request.Limit
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:        ptrString("zone-12345678"),
					ZoneName:      ptrString("example.com"),
					Status:        ptrString("active"),
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

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	res := teo.ResourceTencentCloudTeoZone()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
		"type":      "partial",
		"area":      "overseas",
		"plan_id":   "edgeone-2kfv1h391n6w",
		"limit":     100,
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify limit was passed but offset uses default (0)
	assert.NotNil(t, capturedOffset)
	assert.NotNil(t, capturedLimit)
	assert.Equal(t, int64(0), *capturedOffset)
	assert.Equal(t, int64(100), *capturedLimit)
}
