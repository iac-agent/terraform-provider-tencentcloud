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
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaTeoDnsRecord48 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoDnsRecord48) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoDnsRecord48{}

func newMockMetaTeoDnsRecord48() *mockMetaTeoDnsRecord48 {
	return &mockMetaTeoDnsRecord48{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringTeo(s string) *string { return &s }
func ptrInt64Teo(i int64) *int64    { return &i }

// go test ./tencentcloud/services/teo/ -run "TestUnitTeoDnsRecord48" -v -count=1 -gcflags="all=-l"

func TestUnitTeoDnsRecord48Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecord48()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	createCalled := false
	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		createCalled = true
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Equal(t, "a.example.cn", *request.Name)
		assert.Equal(t, "A", *request.Type)
		assert.Equal(t, "1.2.3.5", *request.Content)
		assert.Equal(t, "Default", *request.Location)
		assert.Equal(t, int64(300), *request.TTL)
		assert.Equal(t, int64(-1), *request.Weight)
		assert.Equal(t, int64(0), *request.Priority)
		resp := &teov20220901.CreateDnsRecordResponse{}
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStringTeo("record-12345"),
			RequestId: ptrStringTeo("fake-request-id"),
		}
		return resp, nil
	})

	// Mock the service DescribeTeoDnsRecordById used by Read after Create
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)
	describeCalled := false
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		describeCalled = true
		resp := &teov20220901.DescribeDnsRecordsResponse{}
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64Teo(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringTeo("zone-39quuimqg8r6"),
					RecordId:   ptrStringTeo("record-12345"),
					Name:       ptrStringTeo("a.example.cn"),
					Type:       ptrStringTeo("A"),
					Content:    ptrStringTeo("1.2.3.5"),
					Location:   ptrStringTeo("Default"),
					TTL:        ptrInt64Teo(300),
					Weight:     ptrInt64Teo(-1),
					Priority:   ptrInt64Teo(0),
					Status:     ptrStringTeo("enable"),
					CreatedOn:  ptrStringTeo("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringTeo("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringTeo("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoDnsRecord48()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.example.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.True(t, createCalled)
	assert.True(t, describeCalled)
	assert.Equal(t, "zone-39quuimqg8r6#record-12345", d.Id())
	assert.Equal(t, "record-12345", d.Get("record_id"))
}

func TestUnitTeoDnsRecord48CreateEmptyRecordId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecord48()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := &teov20220901.CreateDnsRecordResponse{}
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStringTeo(""),
			RequestId: ptrStringTeo("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoDnsRecord48()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.example.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

func TestUnitTeoDnsRecord48Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecord48()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := &teov20220901.DescribeDnsRecordsResponse{}
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64Teo(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringTeo("zone-39quuimqg8r6"),
					RecordId:   ptrStringTeo("record-12345"),
					Name:       ptrStringTeo("a.example.cn"),
					Type:       ptrStringTeo("A"),
					Content:    ptrStringTeo("1.2.3.5"),
					Location:   ptrStringTeo("Default"),
					TTL:        ptrInt64Teo(300),
					Weight:     ptrInt64Teo(-1),
					Priority:   ptrInt64Teo(0),
					Status:     ptrStringTeo("enable"),
					CreatedOn:  ptrStringTeo("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringTeo("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringTeo("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoDnsRecord48()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.example.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})
	d.SetId("zone-39quuimqg8r6#record-12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-39quuimqg8r6#record-12345", d.Id())
	assert.Equal(t, "zone-39quuimqg8r6", d.Get("zone_id"))
	assert.Equal(t, "a.example.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "record-12345", d.Get("record_id"))
	assert.Equal(t, "enable", d.Get("status"))
}

func TestUnitTeoDnsRecord48ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecord48()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := &teov20220901.DescribeDnsRecordsResponse{}
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64Teo(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrStringTeo("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoDnsRecord48()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.example.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})
	d.SetId("zone-39quuimqg8r6#record-12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestUnitTeoDnsRecord48Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecord48()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Len(t, request.DnsRecords, 1)
		assert.Equal(t, "record-12345", *request.DnsRecords[0].RecordId)
		assert.Equal(t, "1.2.3.6", *request.DnsRecords[0].Content)
		resp := &teov20220901.ModifyDnsRecordsResponse{}
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringTeo("fake-request-id"),
		}
		return resp, nil
	})

	// Mock Read after Update
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := &teov20220901.DescribeDnsRecordsResponse{}
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64Teo(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringTeo("zone-39quuimqg8r6"),
					RecordId:   ptrStringTeo("record-12345"),
					Name:       ptrStringTeo("a.example.cn"),
					Type:       ptrStringTeo("A"),
					Content:    ptrStringTeo("1.2.3.6"),
					Location:   ptrStringTeo("Default"),
					TTL:        ptrInt64Teo(300),
					Weight:     ptrInt64Teo(-1),
					Priority:   ptrInt64Teo(0),
					Status:     ptrStringTeo("enable"),
					CreatedOn:  ptrStringTeo("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringTeo("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringTeo("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoDnsRecord48()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.example.cn",
		"type":     "A",
		"content":  "1.2.3.6",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})
	d.SetId("zone-39quuimqg8r6#record-12345")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "1.2.3.6", d.Get("content"))
}

func TestUnitTeoDnsRecord48Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecord48()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		deleteCalled = true
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Len(t, request.RecordIds, 1)
		assert.Equal(t, "record-12345", *request.RecordIds[0])
		resp := &teov20220901.DeleteDnsRecordsResponse{}
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrStringTeo("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoDnsRecord48()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.example.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})
	d.SetId("zone-39quuimqg8r6#record-12345")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}
