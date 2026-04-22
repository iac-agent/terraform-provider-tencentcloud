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

// mockMetaDnsRecordV5 implements tccommon.ProviderMeta
type mockMetaDnsRecordV5 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDnsRecordV5) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDnsRecordV5{}

func newMockMetaDnsRecordV5() *mockMetaDnsRecordV5 {
	return &mockMetaDnsRecordV5{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDnsV5(s string) *string {
	return &s
}

func ptrInt64DnsV5(i int64) *int64 {
	return &i
}

// mockDescribeDnsRecords creates a standard DescribeDnsRecords mock response
func mockDescribeDnsRecords(patches *gomonkey.Patches, teoClient *teov20220901.Client, name, status string) {
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DnsV5(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringDnsV5("zone-39quuimqg8r6"),
					RecordId:   ptrStringDnsV5("record-abc123"),
					Name:       ptrStringDnsV5(name),
					Type:       ptrStringDnsV5("A"),
					Content:    ptrStringDnsV5("1.2.3.5"),
					Location:   ptrStringDnsV5("Default"),
					TTL:        ptrInt64DnsV5(300),
					Weight:     ptrInt64DnsV5(-1),
					Priority:   ptrInt64DnsV5(5),
					Status:     ptrStringDnsV5(status),
					CreatedOn:  ptrStringDnsV5("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringDnsV5("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringDnsV5("fake-request-id"),
		}
		return resp, nil
	})
}

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecordV5" -v -count=1 -gcflags="all=-l"

// TestTeoDnsRecordV5_CreateSuccess tests Create with all required and optional fields
func TestTeoDnsRecordV5_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	// Patch UseTeoV20220901Client (used by Create directly)
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoV20220901Client", teoClient)
	// Patch UseTeoClient (used by Read via TeoService)
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoClient", teoClient)

	// Patch CreateDnsRecordWithContext
	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Equal(t, "a.makn.cn", *request.Name)
		assert.Equal(t, "A", *request.Type)
		assert.Equal(t, "1.2.3.5", *request.Content)

		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStringDnsV5("record-abc123"),
			RequestId: ptrStringDnsV5("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeDnsRecords (called by Read via TeoService)
	mockDescribeDnsRecords(patches, teoClient, "a.makn.cn", "enable")

	meta := newMockMetaDnsRecordV5()
	res := teo.ResourceTencentCloudTeoDnsRecordV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
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
	assert.Equal(t, "zone-39quuimqg8r6#record-abc123", d.Id())
	assert.Equal(t, "record-abc123", d.Get("record_id"))
	assert.Equal(t, "a.makn.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 5, d.Get("priority"))
	assert.Equal(t, "enable", d.Get("status"))
}

// TestTeoDnsRecordV5_CreateAPIError tests Create handles API error
func TestTeoDnsRecordV5_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaDnsRecordV5()
	res := teo.ResourceTencentCloudTeoDnsRecordV5()
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

// TestTeoDnsRecordV5_ReadSuccess tests Read with existing record
func TestTeoDnsRecordV5_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoClient", teoClient)

	mockDescribeDnsRecords(patches, teoClient, "a.makn.cn", "enable")

	meta := newMockMetaDnsRecordV5()
	res := teo.ResourceTencentCloudTeoDnsRecordV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-39quuimqg8r6#record-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-39quuimqg8r6#record-abc123", d.Id())
	assert.Equal(t, "record-abc123", d.Get("record_id"))
	assert.Equal(t, "a.makn.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 5, d.Get("priority"))
	assert.Equal(t, "enable", d.Get("status"))
}

// TestTeoDnsRecordV5_ReadDeleted tests Read with non-existing record (deleted)
func TestTeoDnsRecordV5_ReadDeleted(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DnsV5(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrStringDnsV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDnsRecordV5()
	res := teo.ResourceTencentCloudTeoDnsRecordV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-39quuimqg8r6#record-deleted")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecordV5_UpdateMutableFields tests Update for mutable fields via ModifyDnsRecords
func TestTeoDnsRecordV5_UpdateMutableFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoClient", teoClient)

	// Patch ModifyDnsRecordsWithContext
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Equal(t, 1, len(request.DnsRecords))
		assert.Equal(t, "record-abc123", *request.DnsRecords[0].RecordId)
		assert.Equal(t, "b.makn.cn", *request.DnsRecords[0].Name)

		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringDnsV5("fake-request-id"),
		}
		return resp, nil
	})

	// Patch ModifyDnsRecordsStatusWithContext (may also be called if HasChange("status") is true)
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringDnsV5("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeDnsRecords (called by Read via TeoService after update)
	mockDescribeDnsRecords(patches, teoClient, "b.makn.cn", "enable")

	meta := newMockMetaDnsRecordV5()
	res := teo.ResourceTencentCloudTeoDnsRecordV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "b.makn.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 5,
		"status":   "enable",
	})
	d.SetId("zone-39quuimqg8r6#record-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecordV5_UpdateStatus tests Update for status field via ModifyDnsRecordsStatus
func TestTeoDnsRecordV5_UpdateStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoClient", teoClient)

	// Patch ModifyDnsRecordsWithContext (may be called if HasChange returns true for mutable args)
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringDnsV5("fake-request-id"),
		}
		return resp, nil
	})

	// Patch ModifyDnsRecordsStatusWithContext
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)

		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringDnsV5("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeDnsRecords (called by Read after update)
	mockDescribeDnsRecords(patches, teoClient, "a.makn.cn", "disable")

	meta := newMockMetaDnsRecordV5()
	res := teo.ResourceTencentCloudTeoDnsRecordV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.makn.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 5,
		"status":   "disable",
	})
	d.SetId("zone-39quuimqg8r6#record-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecordV5_DeleteSuccess tests Delete with valid record
func TestTeoDnsRecordV5_DeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Equal(t, 1, len(request.RecordIds))
		assert.Equal(t, "record-abc123", *request.RecordIds[0])

		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrStringDnsV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDnsRecordV5()
	res := teo.ResourceTencentCloudTeoDnsRecordV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-39quuimqg8r6#record-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecordV5_DeleteAPIError tests Delete handles API error
func TestTeoDnsRecordV5_DeleteAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecordV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMetaDnsRecordV5()
	res := teo.ResourceTencentCloudTeoDnsRecordV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-39quuimqg8r6#record-abc123")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecordV5_Schema validates schema definition
func TestTeoDnsRecordV5_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecordV5()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Required fields
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)

	typeField := res.Schema["type"]
	assert.Equal(t, schema.TypeString, typeField.Type)
	assert.True(t, typeField.Required)

	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Required)

	// Optional+Computed fields
	location := res.Schema["location"]
	assert.Equal(t, schema.TypeString, location.Type)
	assert.True(t, location.Optional)
	assert.True(t, location.Computed)

	ttl := res.Schema["ttl"]
	assert.Equal(t, schema.TypeInt, ttl.Type)
	assert.True(t, ttl.Optional)
	assert.True(t, ttl.Computed)

	weight := res.Schema["weight"]
	assert.Equal(t, schema.TypeInt, weight.Type)
	assert.True(t, weight.Optional)
	assert.True(t, weight.Computed)

	priority := res.Schema["priority"]
	assert.Equal(t, schema.TypeInt, priority.Type)
	assert.True(t, priority.Optional)
	assert.True(t, priority.Computed)

	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Optional)
	assert.True(t, status.Computed)

	// Computed fields
	recordId := res.Schema["record_id"]
	assert.Equal(t, schema.TypeString, recordId.Type)
	assert.True(t, recordId.Computed)
	assert.False(t, recordId.Optional)

	createdOn := res.Schema["created_on"]
	assert.Equal(t, schema.TypeString, createdOn.Type)
	assert.True(t, createdOn.Computed)
	assert.False(t, createdOn.Optional)

	modifiedOn := res.Schema["modified_on"]
	assert.Equal(t, schema.TypeString, modifiedOn.Type)
	assert.True(t, modifiedOn.Computed)
	assert.False(t, modifiedOn.Optional)
}
