package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecord22" -v -count=1 -gcflags="all=-l"

// TestTeoDnsRecord22_Create_Success tests Create calls API and sets composite ID
func TestTeoDnsRecord22_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrString("record-abc123"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString("zone-1234567890"),
					RecordId:   ptrString("record-abc123"),
					Name:       ptrString("a.example.com"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.5"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64DNS22(300),
					Weight:     ptrInt64DNS22(-1),
					Priority:   ptrInt64DNS22(5),
					Status:     ptrString("enable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord22()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "a.example.com",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 5,
		"status":   "enable",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#record-abc123", d.Id())
}

// TestTeoDnsRecord22_Create_APIError tests Create handles API error
func TestTeoDnsRecord22_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord22()
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

// TestTeoDnsRecord22_Read_Success tests Read retrieves DNS record data
func TestTeoDnsRecord22_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString("zone-1234567890"),
					RecordId:   ptrString("record-abc123"),
					Name:       ptrString("a.example.com"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.5"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64DNS22(300),
					Weight:     ptrInt64DNS22(-1),
					Priority:   ptrInt64DNS22(5),
					Status:     ptrString("enable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord22()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890#record-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "a.example.com", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 5, d.Get("priority"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "record-abc123", d.Get("record_id"))
}

// TestTeoDnsRecord22_Read_NotFound tests Read handles record not found
func TestTeoDnsRecord22_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord22()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890#record-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord22_Update_ModifyDnsRecords tests Update with DNS record modification
func TestTeoDnsRecord22_Update_ModifyDnsRecords(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString("zone-1234567890"),
					RecordId:   ptrString("record-abc123"),
					Name:       ptrString("a.example.com"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.6"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64DNS22(300),
					Weight:     ptrInt64DNS22(-1),
					Priority:   ptrInt64DNS22(5),
					Status:     ptrString("enable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-01T01:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord22()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "a.example.com",
		"type":     "A",
		"content":  "1.2.3.6",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 5,
		"status":   "enable",
	})
	d.SetId("zone-1234567890#record-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord22_Update_ModifyDnsRecordsStatus tests Update with status modification
func TestTeoDnsRecord22_Update_ModifyDnsRecordsStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString("zone-1234567890"),
					RecordId:   ptrString("record-abc123"),
					Name:       ptrString("a.example.com"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.5"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64DNS22(300),
					Weight:     ptrInt64DNS22(-1),
					Priority:   ptrInt64DNS22(5),
					Status:     ptrString("disable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-01T01:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord22()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "a.example.com",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 5,
		"status":   "disable",
	})
	d.SetId("zone-1234567890#record-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord22_Delete_Success tests Delete removes DNS record
func TestTeoDnsRecord22_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord22()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890#record-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord22_Delete_APIError tests Delete handles API error
func TestTeoDnsRecord22_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord22()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890#record-abc123")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecord22_Schema validates schema definition
func TestTeoDnsRecord22_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecord22()

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

	assert.Contains(t, res.Schema, "type")
	typeField := res.Schema["type"]
	assert.Equal(t, schema.TypeString, typeField.Type)
	assert.True(t, typeField.Required)

	assert.Contains(t, res.Schema, "content")
	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Required)

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

	assert.Contains(t, res.Schema, "record_id")
	recordId := res.Schema["record_id"]
	assert.Equal(t, schema.TypeString, recordId.Type)
	assert.True(t, recordId.Computed)
}

func ptrInt64DNS22(i int64) *int64 {
	return &i
}
