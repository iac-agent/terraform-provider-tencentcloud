package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaTeoDnsRecordV3 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoDnsRecordV3) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoDnsRecordV3{}

func newMockMetaTeoDnsRecordV3() *mockMetaTeoDnsRecordV3 {
	return &mockMetaTeoDnsRecordV3{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDRV3(s string) *string {
	return &s
}

func ptrInt64DRV3(i int64) *int64 {
	return &i
}

func mockTeoDnsRecordV3() *teov20220901.DnsRecord {
	return &teov20220901.DnsRecord{
		ZoneId:     ptrStringDRV3("zone-test123"),
		RecordId:   ptrStringDRV3("record-123"),
		Name:       ptrStringDRV3("a.makn.cn"),
		Type:       ptrStringDRV3("A"),
		Location:   ptrStringDRV3("Default"),
		Content:    ptrStringDRV3("1.2.3.5"),
		TTL:        ptrInt64DRV3(300),
		Weight:     ptrInt64DRV3(-1),
		Priority:   ptrInt64DRV3(5),
		Status:     ptrStringDRV3("enable"),
		CreatedOn:  ptrStringDRV3("2024-01-01T00:00:00Z"),
		ModifiedOn: ptrStringDRV3("2024-01-02T00:00:00Z"),
	}
}

func mockDescribeDnsRecordsV3(t *testing.T, teoClient *teov20220901.Client, patches *gomonkey.Patches, record *teov20220901.DnsRecord) {
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)

		resp := teov20220901.NewDescribeDnsRecordsResponse()
		if record == nil {
			resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
				TotalCount: ptrInt64DRV3(0),
				DnsRecords: []*teov20220901.DnsRecord{},
				RequestId:  ptrStringDRV3("fake-request-id-read"),
			}
			return resp, nil
		}

		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DRV3(1),
			DnsRecords: []*teov20220901.DnsRecord{record},
			RequestId:  ptrStringDRV3("fake-request-id-read"),
		}
		return resp, nil
	})
}

func TestTeoDnsRecordV3_Create_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Name)
		assert.Equal(t, "a.makn.cn", *request.Name)
		assert.NotNil(t, request.Type)
		assert.Equal(t, "A", *request.Type)
		assert.NotNil(t, request.Content)
		assert.Equal(t, "1.2.3.5", *request.Content)

		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStringDRV3("record-123"),
			RequestId: ptrStringDRV3("fake-request-id-create"),
		}
		return resp, nil
	})

	mockDescribeDnsRecordsV3(t, teoClient, patches, mockTeoDnsRecordV3())

	meta := newMockMetaTeoDnsRecordV3()
	res := teo.ResourceTencentCloudTeoDnsRecordV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-test123",
		"name":     "a.makn.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 5,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#record-123", d.Id())
	assert.Equal(t, "record-123", d.Get("record_id"))
	assert.Equal(t, "enable", d.Get("status"))
}

func TestTeoDnsRecordV3_Create_RecordIdNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  nil,
			RequestId: ptrStringDRV3("fake-request-id-create"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecordV3()
	res := teo.ResourceTencentCloudTeoDnsRecordV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

func TestTeoDnsRecordV3_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoClient", teoClient)
	mockDescribeDnsRecordsV3(t, teoClient, patches, mockTeoDnsRecordV3())

	meta := newMockMetaTeoDnsRecordV3()
	res := teo.ResourceTencentCloudTeoDnsRecordV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123#record-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#record-123", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "record-123", d.Get("record_id"))
	assert.Equal(t, "a.makn.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 5, d.Get("priority"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get("modified_on"))
}

func TestTeoDnsRecordV3_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoClient", teoClient)
	mockDescribeDnsRecordsV3(t, teoClient, patches, nil)

	meta := newMockMetaTeoDnsRecordV3()
	res := teo.ResourceTencentCloudTeoDnsRecordV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123#record-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoDnsRecordV3_Update_Change(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoClient", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		modifyCalled = true
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, 1, len(request.DnsRecords))
		assert.NotNil(t, request.DnsRecords[0].RecordId)
		assert.Equal(t, "record-123", *request.DnsRecords[0].RecordId)
		assert.NotNil(t, request.DnsRecords[0].TTL)
		assert.Equal(t, int64(600), *request.DnsRecords[0].TTL)

		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringDRV3("fake-request-id-update"),
		}
		return resp, nil
	})

	record := mockTeoDnsRecordV3()
	record.TTL = ptrInt64DRV3(600)
	mockDescribeDnsRecordsV3(t, teoClient, patches, record)

	meta := newMockMetaTeoDnsRecordV3()
	res := teo.ResourceTencentCloudTeoDnsRecordV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-test123",
		"name":     "a.makn.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      600,
		"weight":   -1,
		"priority": 5,
	})
	d.SetId("zone-test123#record-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, 600, d.Get("ttl"))
}

func TestTeoDnsRecordV3_Update_NoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoClient", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		modifyCalled = true
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringDRV3("fake-request-id-update"),
		}
		return resp, nil
	})

	mockDescribeDnsRecordsV3(t, teoClient, patches, mockTeoDnsRecordV3())

	meta := newMockMetaTeoDnsRecordV3()
	res := teo.ResourceTencentCloudTeoDnsRecordV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-test123#record-123")

	patches.ApplyMethodFunc(d, "HasChange", func(_ string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled)
}

func TestTeoDnsRecordV3_Delete_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordV3().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, 1, len(request.RecordIds))
		assert.Equal(t, "record-123", *request.RecordIds[0])

		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrStringDRV3("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecordV3()
	res := teo.ResourceTencentCloudTeoDnsRecordV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123#record-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
