package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecordV7" -v -count=1 -gcflags="all=-l"

// TestTeoDnsRecordV7_Create_Success tests Create calls API and sets composite ID
func TestTeoDnsRecordV7_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStringForV7("record-abcdefghij"),
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64ForV7(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringForV7("zone-1234567890"),
					RecordId:   ptrStringForV7("record-abcdefghij"),
					Name:       ptrStringForV7("www"),
					Type:       ptrStringForV7("A"),
					Content:    ptrStringForV7("1.2.3.4"),
					Location:   ptrStringForV7("Default"),
					TTL:        ptrInt64ForV7(300),
					Weight:     ptrInt64ForV7(-1),
					Priority:   ptrInt64ForV7(0),
					Status:     ptrStringForV7("enable"),
					CreatedOn:  ptrStringForV7("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringForV7("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV7()
	res := teo.ResourceTencentCloudTeoDnsRecordV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "www",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#record-abcdefghij", d.Id())
}

// TestTeoDnsRecordV7_Create_APIError tests Create handles API error
func TestTeoDnsRecordV7_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForV7()
	res := teo.ResourceTencentCloudTeoDnsRecordV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "www",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoDnsRecordV7_Read_Success tests Read retrieves DNS record data
func TestTeoDnsRecordV7_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64ForV7(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringForV7("zone-1234567890"),
					RecordId:   ptrStringForV7("record-abcdefghij"),
					Name:       ptrStringForV7("www"),
					Type:       ptrStringForV7("A"),
					Content:    ptrStringForV7("1.2.3.4"),
					Location:   ptrStringForV7("Default"),
					TTL:        ptrInt64ForV7(300),
					Weight:     ptrInt64ForV7(-1),
					Priority:   ptrInt64ForV7(5),
					Status:     ptrStringForV7("enable"),
					CreatedOn:  ptrStringForV7("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringForV7("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV7()
	res := teo.ResourceTencentCloudTeoDnsRecordV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "www",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "www", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.4", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, "enable", d.Get("status"))
}

// TestTeoDnsRecordV7_Read_NotFound tests Read handles record not found
func TestTeoDnsRecordV7_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64ForV7(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV7()
	res := teo.ResourceTencentCloudTeoDnsRecordV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "www",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecordV7_Update_Content tests Update with content changes
func TestTeoDnsRecordV7_Update_Content(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64ForV7(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringForV7("zone-1234567890"),
					RecordId:   ptrStringForV7("record-abcdefghij"),
					Name:       ptrStringForV7("www"),
					Type:       ptrStringForV7("A"),
					Content:    ptrStringForV7("5.6.7.8"),
					Location:   ptrStringForV7("Default"),
					TTL:        ptrInt64ForV7(300),
					Weight:     ptrInt64ForV7(-1),
					Priority:   ptrInt64ForV7(0),
					Status:     ptrStringForV7("enable"),
					CreatedOn:  ptrStringForV7("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringForV7("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV7()
	res := teo.ResourceTencentCloudTeoDnsRecordV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "www",
		"type":    "A",
		"content": "5.6.7.8",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecordV7_Update_Status tests Update with status change
func TestTeoDnsRecordV7_Update_Status(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64ForV7(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringForV7("zone-1234567890"),
					RecordId:   ptrStringForV7("record-abcdefghij"),
					Name:       ptrStringForV7("www"),
					Type:       ptrStringForV7("A"),
					Content:    ptrStringForV7("1.2.3.4"),
					Location:   ptrStringForV7("Default"),
					TTL:        ptrInt64ForV7(300),
					Weight:     ptrInt64ForV7(-1),
					Priority:   ptrInt64ForV7(0),
					Status:     ptrStringForV7("disable"),
					CreatedOn:  ptrStringForV7("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringForV7("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV7()
	res := teo.ResourceTencentCloudTeoDnsRecordV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "www",
		"type":    "A",
		"content": "1.2.3.4",
		"status":  "disable",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecordV7_Delete_Success tests Delete removes DNS record
func TestTeoDnsRecordV7_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrStringForV7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV7()
	res := teo.ResourceTencentCloudTeoDnsRecordV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "www",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecordV7_Delete_APIError tests Delete handles API error
func TestTeoDnsRecordV7_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMetaForV7()
	res := teo.ResourceTencentCloudTeoDnsRecordV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "www",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecordV7_Schema validates schema definition
func TestTeoDnsRecordV7_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecordV7()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields with ForceNew
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check required fields without ForceNew
	assert.Contains(t, res.Schema, "name")
	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.False(t, name.ForceNew)

	assert.Contains(t, res.Schema, "type")
	typeField := res.Schema["type"]
	assert.Equal(t, schema.TypeString, typeField.Type)
	assert.True(t, typeField.Required)
	assert.False(t, typeField.ForceNew)

	assert.Contains(t, res.Schema, "content")
	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Required)
	assert.False(t, content.ForceNew)

	// Check optional fields
	assert.Contains(t, res.Schema, "location")
	location := res.Schema["location"]
	assert.Equal(t, schema.TypeString, location.Type)
	assert.True(t, location.Optional)
	assert.True(t, location.Computed)

	assert.Contains(t, res.Schema, "ttl")
	ttl := res.Schema["ttl"]
	assert.Equal(t, schema.TypeInt, ttl.Type)
	assert.True(t, ttl.Optional)
	assert.True(t, ttl.Computed)

	assert.Contains(t, res.Schema, "weight")
	weight := res.Schema["weight"]
	assert.Equal(t, schema.TypeInt, weight.Type)
	assert.True(t, weight.Optional)
	assert.True(t, weight.Computed)

	assert.Contains(t, res.Schema, "priority")
	priority := res.Schema["priority"]
	assert.Equal(t, schema.TypeInt, priority.Type)
	assert.True(t, priority.Optional)
	assert.True(t, priority.Computed)

	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Optional)
	assert.True(t, status.Computed)

	// Check computed fields
	assert.Contains(t, res.Schema, "created_on")
	createdOn := res.Schema["created_on"]
	assert.Equal(t, schema.TypeString, createdOn.Type)
	assert.True(t, createdOn.Computed)

	assert.Contains(t, res.Schema, "modified_on")
	modifiedOn := res.Schema["modified_on"]
	assert.Equal(t, schema.TypeString, modifiedOn.Type)
	assert.True(t, modifiedOn.Computed)
}

// mockMetaForV7 implements tccommon.ProviderMeta
type mockMetaForV7 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForV7) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForV7{}

func newMockMetaForV7() *mockMetaForV7 {
	return &mockMetaForV7{client: &connectivity.TencentCloudClient{}}
}

func ptrStringForV7(s string) *string {
	return &s
}

func ptrInt64ForV7(i int64) *int64 {
	return &i
}
