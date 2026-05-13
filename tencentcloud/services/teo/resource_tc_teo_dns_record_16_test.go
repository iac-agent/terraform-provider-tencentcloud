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

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecord16" -v -count=1 -gcflags="all=-l"

// mockMeta16 implements tccommon.ProviderMeta
type mockMeta16 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMeta16) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMeta16{}

func newMockMeta16() *mockMeta16 {
	return &mockMeta16{client: &connectivity.TencentCloudClient{}}
}

func ptrString16(s string) *string {
	return &s
}

func ptrInt64(i int64) *int64 {
	return &i
}

// TestTeoDnsRecord16_Create_Success tests Create calls API and sets composite ID
func TestTeoDnsRecord16_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta16().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrString16("record-abcdefghij"),
			RequestId: ptrString16("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString16("zone-1234567890"),
					RecordId:   ptrString16("record-abcdefghij"),
					Name:       ptrString16("a.example.com"),
					Type:       ptrString16("A"),
					Content:    ptrString16("1.2.3.5"),
					Location:   ptrString16("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(0),
					Status:     ptrString16("enable"),
					CreatedOn:  ptrString16("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString16("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString16("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta16()
	res := teo.ResourceTencentCloudTeoDnsRecord16()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "a.example.com",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#record-abcdefghij", d.Id())
	assert.Equal(t, "a.example.com", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
}

// TestTeoDnsRecord16_Create_APIError tests Create handles API error
func TestTeoDnsRecord16_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta16().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta16()
	res := teo.ResourceTencentCloudTeoDnsRecord16()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoDnsRecord16_Read_Success tests Read retrieves DNS record data
func TestTeoDnsRecord16_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta16().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString16("zone-1234567890"),
					RecordId:   ptrString16("record-abcdefghij"),
					Name:       ptrString16("a.example.com"),
					Type:       ptrString16("A"),
					Content:    ptrString16("1.2.3.5"),
					Location:   ptrString16("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(0),
					Status:     ptrString16("enable"),
					CreatedOn:  ptrString16("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString16("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString16("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta16()
	res := teo.ResourceTencentCloudTeoDnsRecord16()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "a.example.com", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("modified_on"))
}

// TestTeoDnsRecord16_Read_NotFound tests Read handles record not found
func TestTeoDnsRecord16_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta16().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrString16("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta16()
	res := teo.ResourceTencentCloudTeoDnsRecord16()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord16_Update_Success tests Update calls ModifyDnsRecords
func TestTeoDnsRecord16_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta16().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrString16("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString16("zone-1234567890"),
					RecordId:   ptrString16("record-abcdefghij"),
					Name:       ptrString16("a.example.com"),
					Type:       ptrString16("A"),
					Content:    ptrString16("1.2.3.6"),
					Location:   ptrString16("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(0),
					Status:     ptrString16("enable"),
					CreatedOn:  ptrString16("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString16("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString16("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta16()
	res := teo.ResourceTencentCloudTeoDnsRecord16()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "a.example.com",
		"type":     "A",
		"content":  "1.2.3.6",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord16_Update_APIError tests Update handles API error
func TestTeoDnsRecord16_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta16().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid record")
	})

	meta := newMockMeta16()
	res := teo.ResourceTencentCloudTeoDnsRecord16()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "a.example.com",
		"type":     "A",
		"content":  "1.2.3.6",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoDnsRecord16_Delete_Success tests Delete removes DNS record
func TestTeoDnsRecord16_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta16().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrString16("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta16()
	res := teo.ResourceTencentCloudTeoDnsRecord16()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord16_Delete_APIError tests Delete handles API error
func TestTeoDnsRecord16_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta16().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMeta16()
	res := teo.ResourceTencentCloudTeoDnsRecord16()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecord16_Schema validates schema definition
func TestTeoDnsRecord16_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecord16()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

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
	assert.False(t, location.ForceNew)

	assert.Contains(t, res.Schema, "ttl")
	ttl := res.Schema["ttl"]
	assert.Equal(t, schema.TypeInt, ttl.Type)
	assert.True(t, ttl.Optional)
	assert.True(t, ttl.Computed)
	assert.False(t, ttl.ForceNew)

	assert.Contains(t, res.Schema, "weight")
	weight := res.Schema["weight"]
	assert.Equal(t, schema.TypeInt, weight.Type)
	assert.True(t, weight.Optional)
	assert.True(t, weight.Computed)
	assert.False(t, weight.ForceNew)

	assert.Contains(t, res.Schema, "priority")
	priority := res.Schema["priority"]
	assert.Equal(t, schema.TypeInt, priority.Type)
	assert.True(t, priority.Optional)
	assert.True(t, priority.Computed)
	assert.False(t, priority.ForceNew)

	// Check computed-only fields - status should NOT be Optional
	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Computed)
	assert.False(t, status.Optional)

	assert.Contains(t, res.Schema, "created_on")
	createdOn := res.Schema["created_on"]
	assert.Equal(t, schema.TypeString, createdOn.Type)
	assert.True(t, createdOn.Computed)

	assert.Contains(t, res.Schema, "modified_on")
	modifiedOn := res.Schema["modified_on"]
	assert.Equal(t, schema.TypeString, modifiedOn.Type)
	assert.True(t, modifiedOn.Computed)
}
