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

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecord15" -v -count=1 -gcflags="all=-l"

// TestTeoDnsRecord15_Create_Success tests Create calls API and sets composite ID
func TestTeoDnsRecord15_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "a.example.com", *request.Name)
		assert.Equal(t, "A", *request.Type)
		assert.Equal(t, "1.2.3.5", *request.Content)
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrString("rec-abcdefghij"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString("zone-12345678"),
					RecordId:   ptrString("rec-abcdefghij"),
					Name:       ptrString("a.example.com"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.5"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(0),
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
	res := teo.ResourceTencentCloudTeoDnsRecord15()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-12345678",
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
	assert.Equal(t, "zone-12345678#rec-abcdefghij", d.Id())
}

// TestTeoDnsRecord15_Create_APIError tests Create handles API error
func TestTeoDnsRecord15_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord15()
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

// TestTeoDnsRecord15_Read_Success tests Read retrieves DNS record data
func TestTeoDnsRecord15_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString("zone-12345678"),
					RecordId:   ptrString("rec-abcdefghij"),
					Name:       ptrString("a.example.com"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.5"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(0),
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
	res := teo.ResourceTencentCloudTeoDnsRecord15()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-12345678#rec-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "a.example.com", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "rec-abcdefghij", d.Get("record_id"))
}

// TestTeoDnsRecord15_Read_NotFound tests Read handles DNS record not found
func TestTeoDnsRecord15_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord15()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-12345678#rec-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord15_Update_Success tests Update calls ModifyDnsRecords API
func TestTeoDnsRecord15_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, 1, len(request.DnsRecords))
		assert.Equal(t, "rec-abcdefghij", *request.DnsRecords[0].RecordId)
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString("zone-12345678"),
					RecordId:   ptrString("rec-abcdefghij"),
					Name:       ptrString("a.example.com"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.6"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(0),
					Status:     ptrString("enable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord15()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-12345678",
		"name":     "a.example.com",
		"type":     "A",
		"content":  "1.2.3.6",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})
	d.SetId("zone-12345678#rec-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord15_Delete_Success tests Delete calls DeleteDnsRecords API
func TestTeoDnsRecord15_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, 1, len(request.RecordIds))
		assert.Equal(t, "rec-abcdefghij", *request.RecordIds[0])
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord15()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-12345678#rec-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord15_Delete_APIError tests Delete handles API error
func TestTeoDnsRecord15_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoDnsRecord15()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "a.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-12345678#rec-abcdefghij")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecord15_Schema validates schema definition
func TestTeoDnsRecord15_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecord15()

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
	nameField := res.Schema["name"]
	assert.Equal(t, schema.TypeString, nameField.Type)
	assert.True(t, nameField.Required)
	assert.False(t, nameField.ForceNew)

	assert.Contains(t, res.Schema, "type")
	typeField := res.Schema["type"]
	assert.Equal(t, schema.TypeString, typeField.Type)
	assert.True(t, typeField.Required)
	assert.False(t, typeField.ForceNew)

	assert.Contains(t, res.Schema, "content")
	contentField := res.Schema["content"]
	assert.Equal(t, schema.TypeString, contentField.Type)
	assert.True(t, contentField.Required)
	assert.False(t, contentField.ForceNew)

	// Check optional fields
	assert.Contains(t, res.Schema, "location")
	locationField := res.Schema["location"]
	assert.Equal(t, schema.TypeString, locationField.Type)
	assert.True(t, locationField.Optional)
	assert.True(t, locationField.Computed)

	assert.Contains(t, res.Schema, "ttl")
	ttlField := res.Schema["ttl"]
	assert.Equal(t, schema.TypeInt, ttlField.Type)
	assert.True(t, ttlField.Optional)
	assert.True(t, ttlField.Computed)

	assert.Contains(t, res.Schema, "weight")
	weightField := res.Schema["weight"]
	assert.Equal(t, schema.TypeInt, weightField.Type)
	assert.True(t, weightField.Optional)
	assert.True(t, weightField.Computed)

	assert.Contains(t, res.Schema, "priority")
	priorityField := res.Schema["priority"]
	assert.Equal(t, schema.TypeInt, priorityField.Type)
	assert.True(t, priorityField.Optional)
	assert.True(t, priorityField.Computed)

	// Check computed fields
	assert.Contains(t, res.Schema, "record_id")
	recordIdField := res.Schema["record_id"]
	assert.Equal(t, schema.TypeString, recordIdField.Type)
	assert.True(t, recordIdField.Computed)

	assert.Contains(t, res.Schema, "status")
	statusField := res.Schema["status"]
	assert.Equal(t, schema.TypeString, statusField.Type)
	assert.True(t, statusField.Computed)

	assert.Contains(t, res.Schema, "created_on")
	assert.Contains(t, res.Schema, "modified_on")
}

func ptrInt64(i int64) *int64 {
	return &i
}
