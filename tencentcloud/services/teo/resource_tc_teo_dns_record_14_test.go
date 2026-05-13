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

// mockMeta implements tccommon.ProviderMeta
type mockMeta14 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMeta14) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMeta14{}

func newMockMeta14() *mockMeta14 {
	return &mockMeta14{client: &connectivity.TencentCloudClient{}}
}

func ptrString14(s string) *string { return &s }
func ptrInt6414(i int64) *int64    { return &i }

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecord14" -v -count=1 -gcflags="all=-l"

// TestTeoDnsRecord14_Create_Success tests Create calls API and sets ID
func TestTeoDnsRecord14_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrString14("record-abcdefghij"),
			RequestId: ptrString14("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt6414(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString14("zone-1234567890"),
					RecordId:   ptrString14("record-abcdefghij"),
					Name:       ptrString14("a.example.com"),
					Type:       ptrString14("A"),
					Content:    ptrString14("1.2.3.5"),
					Location:   ptrString14("Default"),
					TTL:        ptrInt6414(300),
					Weight:     ptrInt6414(-1),
					Priority:   ptrInt6414(0),
					Status:     ptrString14("enable"),
					CreatedOn:  ptrString14("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString14("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString14("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta14()
	res := teo.ResourceTencentCloudTeoDnsRecord14()
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
}

// TestTeoDnsRecord14_Create_APIError tests Create handles API error
func TestTeoDnsRecord14_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta14()
	res := teo.ResourceTencentCloudTeoDnsRecord14()
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

// TestTeoDnsRecord14_Create_EmptyRecordId tests Create handles empty RecordId
func TestTeoDnsRecord14_Create_EmptyRecordId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrString14(""),
			RequestId: ptrString14("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta14()
	res := teo.ResourceTencentCloudTeoDnsRecord14()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestTeoDnsRecord14_Read_Success tests Read retrieves DNS record data
func TestTeoDnsRecord14_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt6414(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString14("zone-1234567890"),
					RecordId:   ptrString14("record-abcdefghij"),
					Name:       ptrString14("a.example.com"),
					Type:       ptrString14("A"),
					Content:    ptrString14("1.2.3.5"),
					Location:   ptrString14("Default"),
					TTL:        ptrInt6414(300),
					Weight:     ptrInt6414(-1),
					Priority:   ptrInt6414(0),
					Status:     ptrString14("enable"),
					CreatedOn:  ptrString14("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString14("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString14("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta14()
	res := teo.ResourceTencentCloudTeoDnsRecord14()
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
	assert.Equal(t, "record-abcdefghij", d.Get("record_id"))
}

// TestTeoDnsRecord14_Read_NotFound tests Read handles DNS record not found
func TestTeoDnsRecord14_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt6414(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrString14("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta14()
	res := teo.ResourceTencentCloudTeoDnsRecord14()
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

// TestTeoDnsRecord14_Update_Success tests Update calls ModifyDnsRecords
func TestTeoDnsRecord14_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrString14("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt6414(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString14("zone-1234567890"),
					RecordId:   ptrString14("record-abcdefghij"),
					Name:       ptrString14("a.example.com"),
					Type:       ptrString14("A"),
					Content:    ptrString14("2.3.4.5"),
					Location:   ptrString14("Default"),
					TTL:        ptrInt6414(300),
					Weight:     ptrInt6414(-1),
					Priority:   ptrInt6414(0),
					Status:     ptrString14("enable"),
					CreatedOn:  ptrString14("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString14("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString14("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta14()
	res := teo.ResourceTencentCloudTeoDnsRecord14()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "a.example.com",
		"type":     "A",
		"content":  "2.3.4.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	// Simulate a change in content
	_ = d.Set("content", "2.3.4.5")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord14_Delete_Success tests Delete removes DNS record
func TestTeoDnsRecord14_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrString14("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta14()
	res := teo.ResourceTencentCloudTeoDnsRecord14()
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

// TestTeoDnsRecord14_Delete_APIError tests Delete handles API error
func TestTeoDnsRecord14_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta14().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMeta14()
	res := teo.ResourceTencentCloudTeoDnsRecord14()
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

// TestTeoDnsRecord14_Schema validates schema definition
func TestTeoDnsRecord14_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecord14()

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

	assert.Contains(t, res.Schema, "name")
	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.False(t, name.ForceNew)

	assert.Contains(t, res.Schema, "type")
	typeField := res.Schema["type"]
	assert.Equal(t, schema.TypeString, typeField.Type)
	assert.True(t, typeField.Required)
	assert.True(t, typeField.ForceNew)

	assert.Contains(t, res.Schema, "content")
	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Required)
	assert.False(t, content.ForceNew)

	// Check optional fields WITHOUT ForceNew
	assert.Contains(t, res.Schema, "location")
	location := res.Schema["location"]
	assert.Equal(t, schema.TypeString, location.Type)
	assert.True(t, location.Optional)
	assert.False(t, location.ForceNew)

	assert.Contains(t, res.Schema, "ttl")
	ttl := res.Schema["ttl"]
	assert.Equal(t, schema.TypeInt, ttl.Type)
	assert.True(t, ttl.Optional)
	assert.False(t, ttl.ForceNew)

	assert.Contains(t, res.Schema, "weight")
	weight := res.Schema["weight"]
	assert.Equal(t, schema.TypeInt, weight.Type)
	assert.True(t, weight.Optional)
	assert.False(t, weight.ForceNew)

	assert.Contains(t, res.Schema, "priority")
	priority := res.Schema["priority"]
	assert.Equal(t, schema.TypeInt, priority.Type)
	assert.True(t, priority.Optional)
	assert.False(t, priority.ForceNew)

	// Check computed fields
	assert.Contains(t, res.Schema, "record_id")
	recordId := res.Schema["record_id"]
	assert.Equal(t, schema.TypeString, recordId.Type)
	assert.True(t, recordId.Computed)

	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Computed)

	assert.Contains(t, res.Schema, "created_on")
	assert.Contains(t, res.Schema, "modified_on")
}
