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

// mockMetaForDnsRecord50 implements tccommon.ProviderMeta
type mockMetaForDnsRecord50 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForDnsRecord50) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForDnsRecord50{}

func newMockMetaForDnsRecord50() *mockMetaForDnsRecord50 {
	return &mockMetaForDnsRecord50{client: &connectivity.TencentCloudClient{}}
}

func ptrStrDR50(s string) *string {
	return &s
}

func ptrInt64DR50(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecord50" -v -count=1 -gcflags="all=-l"

// TestTeoDnsRecord50_Create_Success tests Create calls API, sets composite ID and backfills attributes
func TestTeoDnsRecord50_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "www.example.com", *request.Name)
		assert.Equal(t, "A", *request.Type)
		assert.Equal(t, "1.2.3.4", *request.Content)
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStrDR50("record-abcdefgh"),
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, int64(1000), *request.Limit)
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "id", *request.Filters[0].Name)
		assert.Len(t, request.Filters[0].Values, 1)
		assert.Equal(t, "record-abcdefgh", *request.Filters[0].Values[0])
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DR50(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStrDR50("zone-12345678"),
					RecordId:   ptrStrDR50("record-abcdefgh"),
					Name:       ptrStrDR50("www.example.com"),
					Type:       ptrStrDR50("A"),
					Location:   ptrStrDR50("Default"),
					Content:    ptrStrDR50("1.2.3.4"),
					TTL:        ptrInt64DR50(300),
					Weight:     ptrInt64DR50(-1),
					Priority:   ptrInt64DR50(0),
					Status:     ptrStrDR50("enable"),
					CreatedOn:  ptrStrDR50("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStrDR50("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#record-abcdefgh", d.Id())
	assert.Equal(t, "record-abcdefgh", d.Get("record_id"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "www.example.com", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.4", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 0, d.Get("priority"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("modified_on"))
}

// TestTeoDnsRecord50_Create_WithOptionalParams tests Create passes optional params to API
func TestTeoDnsRecord50_Create_WithOptionalParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		assert.Equal(t, "CMCC", *request.Location)
		assert.Equal(t, int64(600), *request.TTL)
		assert.Equal(t, int64(30), *request.Weight)
		assert.Equal(t, int64(10), *request.Priority)
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStrDR50("record-with-optional"),
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, _ *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DR50(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStrDR50("zone-12345678"),
					RecordId:   ptrStrDR50("record-with-optional"),
					Name:       ptrStrDR50("www.example.com"),
					Type:       ptrStrDR50("A"),
					Location:   ptrStrDR50("CMCC"),
					Content:    ptrStrDR50("1.2.3.4"),
					TTL:        ptrInt64DR50(600),
					Weight:     ptrInt64DR50(30),
					Priority:   ptrInt64DR50(10),
					Status:     ptrStrDR50("enable"),
					CreatedOn:  ptrStrDR50("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStrDR50("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-12345678",
		"name":     "www.example.com",
		"type":     "A",
		"content":  "1.2.3.4",
		"location": "CMCC",
		"ttl":      600,
		"weight":   30,
		"priority": 10,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#record-with-optional", d.Id())
}

// TestTeoDnsRecord50_Create_EmptyRecordId tests Create fails when API returns empty record id
func TestTeoDnsRecord50_Create_EmptyRecordId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, _ *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStrDR50(""),
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record id is empty")
	assert.Equal(t, "", d.Id(), "no empty id should be written into state")
}

// TestTeoDnsRecord50_Create_APIError tests Create handles API error
func TestTeoDnsRecord50_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, _ *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoDnsRecord50_Read_Success tests Read populates state from API
func TestTeoDnsRecord50_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, int64(1000), *request.Limit)
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "id", *request.Filters[0].Name)
		assert.Equal(t, "record-abcdefgh", *request.Filters[0].Values[0])
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DR50(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStrDR50("zone-12345678"),
					RecordId:   ptrStrDR50("record-abcdefgh"),
					Name:       ptrStrDR50("www.example.com"),
					Type:       ptrStrDR50("A"),
					Location:   ptrStrDR50("Default"),
					Content:    ptrStrDR50("1.2.3.4"),
					TTL:        ptrInt64DR50(300),
					Weight:     ptrInt64DR50(-1),
					Priority:   ptrInt64DR50(0),
					Status:     ptrStrDR50("enable"),
					CreatedOn:  ptrStrDR50("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStrDR50("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#record-abcdefgh", d.Id())
	assert.Equal(t, "record-abcdefgh", d.Get("record_id"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 0, d.Get("priority"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("modified_on"))
}

// TestTeoDnsRecord50_Read_WithQueryParams tests Read sends user-configured filters/sort/match alongside the id filter
func TestTeoDnsRecord50_Read_WithQueryParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		assert.Len(t, request.Filters, 2)
		assert.Equal(t, "id", *request.Filters[0].Name)
		assert.Equal(t, "record-abcdefgh", *request.Filters[0].Values[0])
		assert.Equal(t, "type", *request.Filters[1].Name)
		assert.Equal(t, "A", *request.Filters[1].Values[0])
		if request.Filters[1].Fuzzy != nil {
			assert.False(t, *request.Filters[1].Fuzzy)
		}
		assert.Equal(t, "created-on", *request.SortBy)
		assert.Equal(t, "desc", *request.SortOrder)
		assert.Equal(t, "all", *request.Match)
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DR50(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStrDR50("zone-12345678"),
					RecordId:   ptrStrDR50("record-abcdefgh"),
					Name:       ptrStrDR50("www.example.com"),
					Type:       ptrStrDR50("A"),
					Location:   ptrStrDR50("Default"),
					Content:    ptrStrDR50("1.2.3.4"),
					TTL:        ptrInt64DR50(300),
					Weight:     ptrInt64DR50(-1),
					Priority:   ptrInt64DR50(0),
					Status:     ptrStrDR50("enable"),
					CreatedOn:  ptrStrDR50("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStrDR50("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"name":       "www.example.com",
		"type":       "A",
		"content":    "1.2.3.4",
		"filters":    []interface{}{map[string]interface{}{"name": "type", "values": []interface{}{"A"}}},
		"sort_by":    "created-on",
		"sort_order": "desc",
		"match":      "all",
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#record-abcdefgh", d.Id())
}

// TestTeoDnsRecord50_Read_NotFound tests Read handles record not found
func TestTeoDnsRecord50_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, _ *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DR50(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-12345678#record-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord50_Read_BrokenId tests Read returns error on broken composite id
func TestTeoDnsRecord50_Read_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	// No API call should be made for a broken id.
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, _ *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		t.Fatalf("DescribeDnsRecordsWithContext should not be called when id is broken")
		return nil, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("broken-id-without-separator")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestTeoDnsRecord50_Read_APIError tests Read handles API error
func TestTeoDnsRecord50_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, _ *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=server error")
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}

// TestTeoDnsRecord50_Update_Success tests Update calls ModifyDnsRecords with single-element list
func TestTeoDnsRecord50_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Len(t, request.DnsRecords, 1)
		assert.Equal(t, "record-abcdefgh", *request.DnsRecords[0].RecordId)
		assert.Equal(t, "updated.example.com", *request.DnsRecords[0].Name)
		assert.Equal(t, "A", *request.DnsRecords[0].Type)
		assert.Equal(t, "5.6.7.8", *request.DnsRecords[0].Content)
		assert.Equal(t, "CMCC", *request.DnsRecords[0].Location)
		assert.Equal(t, int64(600), *request.DnsRecords[0].TTL)
		assert.Equal(t, int64(30), *request.DnsRecords[0].Weight)
		assert.Equal(t, int64(10), *request.DnsRecords[0].Priority)
		// Ignored fields must not be set inside the DnsRecord element.
		assert.Nil(t, request.DnsRecords[0].ZoneId)
		assert.Nil(t, request.DnsRecords[0].Status)
		assert.Nil(t, request.DnsRecords[0].CreatedOn)
		assert.Nil(t, request.DnsRecords[0].ModifiedOn)
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DR50(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStrDR50("zone-12345678"),
					RecordId:   ptrStrDR50("record-abcdefgh"),
					Name:       ptrStrDR50("updated.example.com"),
					Type:       ptrStrDR50("A"),
					Location:   ptrStrDR50("CMCC"),
					Content:    ptrStrDR50("5.6.7.8"),
					TTL:        ptrInt64DR50(600),
					Weight:     ptrInt64DR50(30),
					Priority:   ptrInt64DR50(10),
					Status:     ptrStrDR50("enable"),
					CreatedOn:  ptrStrDR50("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStrDR50("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-12345678",
		"name":     "updated.example.com",
		"type":     "A",
		"content":  "5.6.7.8",
		"location": "CMCC",
		"ttl":      600,
		"weight":   30,
		"priority": 10,
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "updated.example.com", d.Get("name"))
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get("modified_on"))
}

// TestTeoDnsRecord50_Update_APIError tests Update handles API error
func TestTeoDnsRecord50_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, _ *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecord50_Delete_Success tests Delete calls DeleteDnsRecords with single-element RecordIds
func TestTeoDnsRecord50_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Len(t, request.RecordIds, 1)
		assert.Equal(t, "record-abcdefgh", *request.RecordIds[0])
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrStrDR50("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord50_Delete_APIError tests Delete handles API error
func TestTeoDnsRecord50_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForDnsRecord50().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, _ *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Record not found")
	})

	meta := newMockMetaForDnsRecord50()
	res := teo.ResourceTencentCloudTeoDnsRecord50()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecord50_Schema validates schema definition
func TestTeoDnsRecord50_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecord50()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "type")
	assert.Contains(t, res.Schema, "content")
	assert.Contains(t, res.Schema, "location")
	assert.Contains(t, res.Schema, "ttl")
	assert.Contains(t, res.Schema, "weight")
	assert.Contains(t, res.Schema, "priority")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "sort_by")
	assert.Contains(t, res.Schema, "sort_order")
	assert.Contains(t, res.Schema, "match")
	assert.Contains(t, res.Schema, "record_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "created_on")
	assert.Contains(t, res.Schema, "modified_on")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	for _, field := range []string{"name", "type", "content"} {
		s := res.Schema[field]
		assert.Equal(t, schema.TypeString, s.Type)
		assert.True(t, s.Required)
		assert.False(t, s.ForceNew)
	}

	for _, field := range []string{"location", "sort_by", "sort_order", "match"} {
		s := res.Schema[field]
		assert.Equal(t, schema.TypeString, s.Type)
		assert.True(t, s.Optional)
	}
	assert.True(t, res.Schema["location"].Computed)

	for _, field := range []string{"ttl", "weight", "priority"} {
		s := res.Schema[field]
		assert.Equal(t, schema.TypeInt, s.Type)
		assert.True(t, s.Optional)
		assert.True(t, s.Computed)
	}

	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)
	assert.False(t, filters.Computed)
	filterElem, ok := filters.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.Contains(t, filterElem.Schema, "name")
	assert.Contains(t, filterElem.Schema, "values")
	assert.Contains(t, filterElem.Schema, "fuzzy")
	assert.True(t, filterElem.Schema["name"].Required)
	assert.True(t, filterElem.Schema["values"].Required)
	assert.True(t, filterElem.Schema["fuzzy"].Optional)

	for _, field := range []string{"record_id", "status", "created_on", "modified_on"} {
		s := res.Schema[field]
		assert.Equal(t, schema.TypeString, s.Type)
		assert.True(t, s.Computed)
		assert.False(t, s.Required)
		assert.False(t, s.Optional)
	}
}
