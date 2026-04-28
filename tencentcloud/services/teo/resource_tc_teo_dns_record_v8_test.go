package teo_test

import (
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

// mockMetaForV8 implements tccommon.ProviderMeta
type mockMetaForV8 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForV8) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForV8{}

func newMockMetaForV8() *mockMetaForV8 {
	return &mockMetaForV8{client: &connectivity.TencentCloudClient{}}
}

func ptrInt64(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestDnsRecordV8" -v -count=1 -gcflags="all=-l"

// TestDnsRecordV8_Create_Success tests Create calls API and sets ID
func TestDnsRecordV8_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForV8().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx interface{}, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrString("record-abcdefghij"),
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
					ZoneId:     ptrString("zone-1234567890"),
					RecordId:   ptrString("record-abcdefghij"),
					Name:       ptrString("a.makn.cn"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.5"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(5),
					Status:     ptrString("enable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV8()
	res := teo.ResourceTencentCloudTeoDnsRecordV8()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-1234567890",
		"name":     "a.makn.cn",
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
	assert.Equal(t, "zone-1234567890#record-abcdefghij", d.Id())
}

// TestDnsRecordV8_Create_APIError tests Create handles API error
func TestDnsRecordV8_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForV8().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx interface{}, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForV8()
	res := teo.ResourceTencentCloudTeoDnsRecordV8()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestDnsRecordV8_Read_Success tests Read retrieves DNS record data
func TestDnsRecordV8_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForV8().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrString("zone-1234567890"),
					RecordId:   ptrString("record-abcdefghij"),
					Name:       ptrString("a.makn.cn"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.5"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(5),
					Status:     ptrString("enable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV8()
	res := teo.ResourceTencentCloudTeoDnsRecordV8()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "a.makn.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, "enable", d.Get("status"))
}

// TestDnsRecordV8_Read_NotFound tests Read handles DNS record not found
func TestDnsRecordV8_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForV8().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV8()
	res := teo.ResourceTencentCloudTeoDnsRecordV8()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestDnsRecordV8_Update_ContentFields tests Update modifies content fields
func TestDnsRecordV8_Update_ContentFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForV8().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx interface{}, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
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
					ZoneId:     ptrString("zone-1234567890"),
					RecordId:   ptrString("record-abcdefghij"),
					Name:       ptrString("a.makn.cn"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.6"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(5),
					Status:     ptrString("enable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV8()
	res := teo.ResourceTencentCloudTeoDnsRecordV8()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.6",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestDnsRecordV8_Update_Status tests Update modifies status
func TestDnsRecordV8_Update_Status(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForV8().client, "UseTeoV20220901Client", teoClient)

	// Mock ModifyDnsRecordsWithContext since TestResourceDataRaw has no old state,
	// d.HasChange() returns true for all fields with values set
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx interface{}, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(ctx interface{}, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
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
					ZoneId:     ptrString("zone-1234567890"),
					RecordId:   ptrString("record-abcdefghij"),
					Name:       ptrString("a.makn.cn"),
					Type:       ptrString("A"),
					Content:    ptrString("1.2.3.5"),
					Location:   ptrString("Default"),
					TTL:        ptrInt64(300),
					Weight:     ptrInt64(-1),
					Priority:   ptrInt64(5),
					Status:     ptrString("disable"),
					CreatedOn:  ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrString("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV8()
	res := teo.ResourceTencentCloudTeoDnsRecordV8()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
		"status":  "disable",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestDnsRecordV8_Delete_Success tests Delete removes DNS record
func TestDnsRecordV8_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForV8().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx interface{}, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForV8()
	res := teo.ResourceTencentCloudTeoDnsRecordV8()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestDnsRecordV8_Delete_APIError tests Delete handles API error
func TestDnsRecordV8_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForV8().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx interface{}, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMetaForV8()
	res := teo.ResourceTencentCloudTeoDnsRecordV8()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-1234567890#record-abcdefghij")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestDnsRecordV8_Schema validates schema definition
func TestDnsRecordV8_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecordV8()

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

	assert.Contains(t, res.Schema, "status")
	statusField := res.Schema["status"]
	assert.Equal(t, schema.TypeString, statusField.Type)
	assert.True(t, statusField.Optional)
	assert.True(t, statusField.Computed)

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
