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

// mockMetaDnsRecord17 implements tccommon.ProviderMeta
type mockMetaDnsRecord17 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDnsRecord17) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDnsRecord17{}

func newMockMetaDnsRecord17() *mockMetaDnsRecord17 {
	return &mockMetaDnsRecord17{client: &connectivity.TencentCloudClient{}}
}

func ptrStrDns17(s string) *string {
	return &s
}

func ptrInt64Dns17(i int64) *int64 {
	return &i
}

func newDnsRecord17ForMock() *teov20220901.DnsRecord {
	return &teov20220901.DnsRecord{
		ZoneId:     ptrStrDns17("zone-12345678"),
		RecordId:   ptrStrDns17("record-abcdefgh"),
		Name:       ptrStrDns17("www.example.com"),
		Type:       ptrStrDns17("A"),
		Location:   ptrStrDns17("Default"),
		Content:    ptrStrDns17("1.2.3.4"),
		TTL:        ptrInt64Dns17(300),
		Weight:     ptrInt64Dns17(-1),
		Priority:   ptrInt64Dns17(0),
		Status:     ptrStrDns17("enable"),
		CreatedOn:  ptrStrDns17("2026-01-01T00:00:00+00:00"),
		ModifiedOn: ptrStrDns17("2026-01-02T00:00:00+00:00"),
	}
}

// patchDescribeDnsRecords17 mocks DescribeDnsRecords, records the Limit of each call into gotLimit.
func patchDescribeDnsRecords17(patches *gomonkey.Patches, teoClient *teov20220901.Client, gotLimit *int64, retRecord *teov20220901.DnsRecord, callCount *int) {
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		if gotLimit != nil {
			*gotLimit = *request.Limit
		}
		if callCount != nil {
			*callCount++
		}
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			RequestId: ptrStrDns17("fake-request-id"),
		}
		if retRecord != nil {
			resp.Response.TotalCount = ptrInt64Dns17(1)
			resp.Response.DnsRecords = []*teov20220901.DnsRecord{retRecord}
		}
		return resp, nil
	})
}

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecord17" -v -count=1 -gcflags="all=-l"

// TestTeoDnsRecord17Create_Success tests Create calls CreateDnsRecord and sets composite id
func TestTeoDnsRecord17Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)
	// service layer read path uses UseTeoClient
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "www.example.com", *request.Name)
		assert.Equal(t, "A", *request.Type)
		assert.Equal(t, "1.2.3.4", *request.Content)
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStrDns17("record-abcdefgh"),
			RequestId: ptrStrDns17("fake-request-id"),
		}
		return resp, nil
	})

	var gotLimit int64 = -1
	patchDescribeDnsRecords17(patches, teoClient, &gotLimit, newDnsRecord17ForMock(), nil)

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#record-abcdefgh", d.Id())
	// Limit not configured, defaults to 1000 (cloud api max value)
	assert.Equal(t, int64(1000), gotLimit)
	assert.Equal(t, "record-abcdefgh", d.Get("record_id"))
	assert.Equal(t, "enable", d.Get("status"))
}

// TestTeoDnsRecord17Create_EmptyRecordId tests Create returns error on empty RecordId
func TestTeoDnsRecord17Create_EmptyRecordId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStrDns17(""),
			RequestId: ptrStrDns17("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RecordId is nil or empty")
	assert.Equal(t, "", d.Id(), "id should not be set when RecordId is empty")
}

// TestTeoDnsRecord17Read_DefaultLimit tests Read calls DescribeDnsRecords with default Limit 1000
func TestTeoDnsRecord17Read_DefaultLimit(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	var (
		gotLimit   int64 = -1
		gotZoneId  string
		gotFilterN string
		gotFilterV string
		gotOffset  int64 = -1
	)
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		gotLimit = *request.Limit
		gotOffset = *request.Offset
		gotZoneId = *request.ZoneId
		if len(request.Filters) > 0 {
			gotFilterN = *request.Filters[0].Name
			gotFilterV = *request.Filters[0].Values[0]
		}
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64Dns17(1),
			DnsRecords: []*teov20220901.DnsRecord{newDnsRecord17ForMock()},
			RequestId:  ptrStrDns17("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Limit not set in config, default 1000
	assert.Equal(t, int64(1000), gotLimit)
	assert.Equal(t, int64(0), gotOffset)
	assert.Equal(t, "zone-12345678", gotZoneId)
	assert.Equal(t, "id", gotFilterN)
	assert.Equal(t, "record-abcdefgh", gotFilterV)
	assert.Equal(t, "www.example.com", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.4", d.Get("content"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, "record-abcdefgh", d.Get("record_id"))
}

// TestTeoDnsRecord17Read_UserLimit tests Read calls DescribeDnsRecords with user-configured Limit
func TestTeoDnsRecord17Read_UserLimit(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	var gotLimit int64 = -1
	patchDescribeDnsRecords17(patches, teoClient, &gotLimit, newDnsRecord17ForMock(), nil)

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"Limit": 500,
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Limit configured by user, passed through to the api request
	assert.Equal(t, int64(500), gotLimit)
}

// TestTeoDnsRecord17Read_NotFound tests Read clears id when DescribeDnsRecords returns empty list
func TestTeoDnsRecord17Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patchDescribeDnsRecords17(patches, teoClient, nil, nil, nil)

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-12345678#record-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord17Read_BrokenId tests Read returns error on broken composite id
func TestTeoDnsRecord17Read_BrokenId(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("broken-id-without-separator")

	err := res.Read(d, newMockMetaDnsRecord17())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestTeoDnsRecord17Update_ModifyDnsRecords tests Update calls ModifyDnsRecords for updatable fields
func TestTeoDnsRecord17Update_ModifyDnsRecords(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	modifyCalled := 0
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		modifyCalled++
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Len(t, request.DnsRecords, 1)
		record := request.DnsRecords[0]
		assert.Equal(t, "record-abcdefgh", *record.RecordId)
		assert.Equal(t, "www.example.com", *record.Name)
		assert.Equal(t, "A", *record.Type)
		assert.Equal(t, "1.2.3.5", *record.Content)
		// output-only fields must not be sent
		assert.Nil(t, record.ZoneId)
		assert.Nil(t, record.Status)
		assert.Nil(t, record.CreatedOn)
		assert.Nil(t, record.ModifiedOn)
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStrDns17("fake-request-id"),
		}
		return resp, nil
	})

	var gotLimit int64 = -1
	patchDescribeDnsRecords17(patches, teoClient, &gotLimit, newDnsRecord17ForMock(), nil)

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, modifyCalled, "ModifyDnsRecordsWithContext should be called exactly once")
}

// TestTeoDnsRecord17Update_StatusToDisable tests Update calls ModifyDnsRecordsStatus when status changes
func TestTeoDnsRecord17Update_StatusToDisable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	modifyStatusCalled := 0
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		modifyStatusCalled++
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "record-abcdefgh", *request.RecordsToDisable[0])
		assert.Len(t, request.RecordsToEnable, 0)
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStrDns17("fake-request-id"),
		}
		return resp, nil
	})

	var gotLimit int64 = -1
	patchDescribeDnsRecords17(patches, teoClient, &gotLimit, newDnsRecord17ForMock(), nil)

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"status": "disable",
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, modifyStatusCalled, "ModifyDnsRecordsStatusWithContext should be called exactly once")
}

// TestTeoDnsRecord17Update_LimitOnlyNoWriteApi tests Update with only Limit change triggers no write api call
func TestTeoDnsRecord17Update_LimitOnlyNoWriteApi(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	// any write api call would panic the test with a clear failure
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		return nil, fmt.Errorf("ModifyDnsRecordsWithContext should not be called when only Limit changes")
	})
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		return nil, fmt.Errorf("ModifyDnsRecordsStatusWithContext should not be called when only Limit changes")
	})

	var gotLimit int64 = -1
	patchDescribeDnsRecords17(patches, teoClient, &gotLimit, newDnsRecord17ForMock(), nil)

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	// only `Limit` is configured, all updatable fields are absent from the diff,
	// Update should skip both write branches and only refresh via Read.
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"Limit": 500,
	})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, int64(500), gotLimit, "Read after Update should use the user-configured Limit")
}

// TestTeoDnsRecord17Delete_Success tests Delete calls DeleteDnsRecords with ZoneId and RecordIds
func TestTeoDnsRecord17Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDnsRecord17()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	deleteCalled := 0
	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		deleteCalled++
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "record-abcdefgh", *request.RecordIds[0])
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrStrDns17("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecord17()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-12345678#record-abcdefgh")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, deleteCalled, "DeleteDnsRecordsWithContext should be called exactly once")
}
