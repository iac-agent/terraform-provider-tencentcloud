package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// dnsRecord49MockMeta implements tccommon.ProviderMeta
type dnsRecord49MockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *dnsRecord49MockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &dnsRecord49MockMeta{}

func newMockMetaDnsRecord49() *dnsRecord49MockMeta {
	return &dnsRecord49MockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDNS49(s string) *string {
	return &s
}

func ptrInt64DNS49(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecord49" -v -count=1 -gcflags="all=-l"

// mockDescribeDnsRecords mocks the service layer query used by Read (via DescribeTeoDnsRecordById).
func mockDescribeDnsRecords(t *testing.T, patches *gomonkey.Patches, teoClient *teov20220901.Client, record *teov20220901.DnsRecord) {
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(1000), *request.Limit)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "id", *request.Filters[0].Name)

		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DNS49(1),
			DnsRecords: []*teov20220901.DnsRecord{record},
			RequestId:  ptrStringDNS49("fake-request-id-read"),
		}
		return resp, nil
	})
}

// TestTeoDnsRecord49_Create_Success tests Create with all optional fields set
func TestTeoDnsRecord49_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Name)
		assert.Equal(t, "www.example.com", *request.Name)
		assert.NotNil(t, request.Type)
		assert.Equal(t, "A", *request.Type)
		assert.NotNil(t, request.Content)
		assert.Equal(t, "1.2.3.4", *request.Content)
		assert.NotNil(t, request.Location)
		assert.Equal(t, "Default", *request.Location)
		assert.NotNil(t, request.TTL)
		assert.Equal(t, int64(300), *request.TTL)
		assert.NotNil(t, request.Weight)
		assert.Equal(t, int64(10), *request.Weight)
		assert.NotNil(t, request.Priority)
		assert.Equal(t, int64(5), *request.Priority)

		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStringDNS49("record-test456"),
			RequestId: ptrStringDNS49("fake-request-id-create"),
		}
		return resp, nil
	})

	mockDescribeDnsRecords(t, patches, teoClient, &teov20220901.DnsRecord{
		ZoneId:     ptrStringDNS49("zone-test123"),
		RecordId:   ptrStringDNS49("record-test456"),
		Name:       ptrStringDNS49("www.example.com"),
		Type:       ptrStringDNS49("A"),
		Location:   ptrStringDNS49("Default"),
		Content:    ptrStringDNS49("1.2.3.4"),
		TTL:        ptrInt64DNS49(300),
		Weight:     ptrInt64DNS49(10),
		Priority:   ptrInt64DNS49(5),
		Status:     ptrStringDNS49("enable"),
		CreatedOn:  ptrStringDNS49("2024-01-01T00:00:00Z"),
		ModifiedOn: ptrStringDNS49("2024-01-01T00:00:00Z"),
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-test123",
		"name":     "www.example.com",
		"type":     "A",
		"content":  "1.2.3.4",
		"location": "Default",
		"ttl":      300,
		"weight":   10,
		"priority": 5,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#record-test456", d.Id())
	assert.Equal(t, "record-test456", d.Get("record_id"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("modified_on"))
}

// TestTeoDnsRecord49_Create_RequiredOnly tests Create with only required fields
func TestTeoDnsRecord49_Create_RequiredOnly(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		assert.Nil(t, request.Location)
		assert.Nil(t, request.TTL)
		assert.Nil(t, request.Weight)
		assert.Nil(t, request.Priority)

		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStringDNS49("record-min789"),
			RequestId: ptrStringDNS49("fake-request-id-create"),
		}
		return resp, nil
	})

	mockDescribeDnsRecords(t, patches, teoClient, &teov20220901.DnsRecord{
		ZoneId:   ptrStringDNS49("zone-test123"),
		RecordId: ptrStringDNS49("record-min789"),
		Name:     ptrStringDNS49("www.example.com"),
		Type:     ptrStringDNS49("A"),
		Content:  ptrStringDNS49("1.2.3.4"),
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#record-min789", d.Id())
}

// TestTeoDnsRecord49_Create_EmptyRecordId tests Create when API returns an empty record id
func TestTeoDnsRecord49_Create_EmptyRecordId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(_ context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrStringDNS49(""),
			RequestId: ptrStringDNS49("fake-request-id-create"),
		}
		return resp, nil
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RecordId is empty")
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord49_Read_Success tests Read back-filling all fields
func TestTeoDnsRecord49_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	mockDescribeDnsRecords(t, patches, teoClient, &teov20220901.DnsRecord{
		ZoneId:     ptrStringDNS49("zone-read123"),
		RecordId:   ptrStringDNS49("record-read456"),
		Name:       ptrStringDNS49("www.example.com"),
		Type:       ptrStringDNS49("CNAME"),
		Location:   ptrStringDNS49("Default"),
		Content:    ptrStringDNS49("target.example.com"),
		TTL:        ptrInt64DNS49(600),
		Weight:     ptrInt64DNS49(-1),
		Priority:   ptrInt64DNS49(0),
		Status:     ptrStringDNS49("enable"),
		CreatedOn:  ptrStringDNS49("2024-01-01T00:00:00Z"),
		ModifiedOn: ptrStringDNS49("2024-06-15T12:00:00Z"),
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "",
		"name":    "",
		"type":    "",
		"content": "",
	})
	d.SetId("zone-read123#record-read456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-read123#record-read456", d.Id())
	assert.Equal(t, "zone-read123", d.Get("zone_id"))
	assert.Equal(t, "www.example.com", d.Get("name"))
	assert.Equal(t, "CNAME", d.Get("type"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, "target.example.com", d.Get("content"))
	assert.Equal(t, 600, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 0, d.Get("priority"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "record-read456", d.Get("record_id"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-06-15T12:00:00Z", d.Get("modified_on"))
}

// TestTeoDnsRecord49_Read_NotFound tests Read when the record does not exist
func TestTeoDnsRecord49_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DNS49(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrStringDNS49("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "",
		"name":    "",
		"type":    "",
		"content": "",
	})
	d.SetId("zone-notfound#record-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord49_Read_BrokenId tests Read with a broken composite id
func TestTeoDnsRecord49_Read_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "",
		"name":    "",
		"type":    "",
		"content": "",
	})
	d.SetId("broken-id-without-separator")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestTeoDnsRecord49_Update_ChangeMutableFields tests Update when mutable fields change
func TestTeoDnsRecord49_Update_ChangeMutableFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		modifyCalled = true

		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-upd123", *request.ZoneId)
		assert.NotNil(t, request.DnsRecords)
		assert.Equal(t, 1, len(request.DnsRecords))

		dnsRecord := request.DnsRecords[0]
		assert.NotNil(t, dnsRecord.RecordId)
		assert.Equal(t, "record-upd456", *dnsRecord.RecordId)
		assert.NotNil(t, dnsRecord.Name)
		assert.Equal(t, "www.example.com", *dnsRecord.Name)
		assert.NotNil(t, dnsRecord.Type)
		assert.Equal(t, "A", *dnsRecord.Type)
		assert.NotNil(t, dnsRecord.Content)
		assert.Equal(t, "5.6.7.8", *dnsRecord.Content)
		assert.NotNil(t, dnsRecord.Location)
		assert.Equal(t, "Default", *dnsRecord.Location)
		assert.NotNil(t, dnsRecord.TTL)
		assert.Equal(t, int64(600), *dnsRecord.TTL)
		// ZoneId/Status/CreatedOn/ModifiedOn must not be passed as input in ModifyDnsRecords
		assert.Nil(t, dnsRecord.ZoneId)
		assert.Nil(t, dnsRecord.Status)
		assert.Nil(t, dnsRecord.CreatedOn)
		assert.Nil(t, dnsRecord.ModifiedOn)

		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringDNS49("fake-request-id-update"),
		}
		return resp, nil
	})

	mockDescribeDnsRecords(t, patches, teoClient, &teov20220901.DnsRecord{
		ZoneId:     ptrStringDNS49("zone-upd123"),
		RecordId:   ptrStringDNS49("record-upd456"),
		Name:       ptrStringDNS49("www.example.com"),
		Type:       ptrStringDNS49("A"),
		Location:   ptrStringDNS49("Default"),
		Content:    ptrStringDNS49("5.6.7.8"),
		TTL:        ptrInt64DNS49(600),
		Status:     ptrStringDNS49("enable"),
		CreatedOn:  ptrStringDNS49("2024-01-01T00:00:00Z"),
		ModifiedOn: ptrStringDNS49("2024-06-15T12:00:00Z"),
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()

	// Build a prior state where content=1.2.3.4, and a new config where content=5.6.7.8,
	// so that d.HasChange("content") returns true inside Update and triggers ModifyDnsRecords.
	state := &terraform.InstanceState{
		ID: "zone-upd123#record-upd456",
		Attributes: map[string]string{
			"id":      "zone-upd123#record-upd456",
			"zone_id": "zone-upd123",
			"name":    "www.example.com",
			"type":    "A",
			"content": "1.2.3.4",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"zone_id":  "zone-upd123",
		"name":     "www.example.com",
		"type":     "A",
		"content":  "5.6.7.8",
		"location": "Default",
		"ttl":      600,
	})

	diff, err := res.Diff(nil, state, rawConfig, meta)
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)
	d.SetId("zone-upd123#record-upd456")

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "5.6.7.8", d.Get("content"))
}

// TestTeoDnsRecord49_Update_NoChange tests Update when no mutable field changed
func TestTeoDnsRecord49_Update_NoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		modifyCalled = true
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrStringDNS49("fake-request-id-update"),
		}
		return resp, nil
	})

	mockDescribeDnsRecords(t, patches, teoClient, &teov20220901.DnsRecord{
		ZoneId:     ptrStringDNS49("zone-same123"),
		RecordId:   ptrStringDNS49("record-same456"),
		Name:       ptrStringDNS49("www.example.com"),
		Type:       ptrStringDNS49("A"),
		Location:   ptrStringDNS49("Default"),
		Content:    ptrStringDNS49("1.2.3.4"),
		TTL:        ptrInt64DNS49(300),
		Status:     ptrStringDNS49("enable"),
		CreatedOn:  ptrStringDNS49("2024-01-01T00:00:00Z"),
		ModifiedOn: ptrStringDNS49("2024-01-01T00:00:00Z"),
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()

	// Build a prior state identical to the new config, so no mutable field changes.
	state := &terraform.InstanceState{
		ID: "zone-same123#record-same456",
		Attributes: map[string]string{
			"id":       "zone-same123#record-same456",
			"zone_id":  "zone-same123",
			"name":     "www.example.com",
			"type":     "A",
			"content":  "1.2.3.4",
			"location": "Default",
			"ttl":      "300",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"zone_id":  "zone-same123",
		"name":     "www.example.com",
		"type":     "A",
		"content":  "1.2.3.4",
		"location": "Default",
		"ttl":      300,
	})

	diff, err := res.Diff(nil, state, rawConfig, meta)
	assert.NoError(t, err)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)
	d.SetId("zone-same123#record-same456")

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled)
	assert.Equal(t, "zone-same123#record-same456", d.Id())
}

// TestTeoDnsRecord49_Delete_Success tests Delete with a single record id
func TestTeoDnsRecord49_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-del123", *request.ZoneId)
		assert.NotNil(t, request.RecordIds)
		assert.Equal(t, 1, len(request.RecordIds))
		assert.Equal(t, "record-del456", *request.RecordIds[0])

		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrStringDNS49("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-del123",
		"name":    "www.example.com",
		"type":    "A",
		"content": "1.2.3.4",
	})
	d.SetId("zone-del123#record-del456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecord49_Import tests import by reading a resource set via composite ID
func TestTeoDnsRecord49_Import(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDnsRecord49().client, "UseTeoV20220901Client", teoClient)

	mockDescribeDnsRecords(t, patches, teoClient, &teov20220901.DnsRecord{
		ZoneId:     ptrStringDNS49("zone-imp123"),
		RecordId:   ptrStringDNS49("record-imp456"),
		Name:       ptrStringDNS49("import.example.com"),
		Type:       ptrStringDNS49("A"),
		Location:   ptrStringDNS49("Default"),
		Content:    ptrStringDNS49("9.9.9.9"),
		TTL:        ptrInt64DNS49(300),
		Status:     ptrStringDNS49("enable"),
		CreatedOn:  ptrStringDNS49("2024-01-01T00:00:00Z"),
		ModifiedOn: ptrStringDNS49("2024-01-01T00:00:00Z"),
	})

	meta := newMockMetaDnsRecord49()
	res := teo.ResourceTencentCloudTeoDnsRecord49()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "",
		"name":    "",
		"type":    "",
		"content": "",
	})
	d.SetId("zone-imp123#record-imp456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-imp123#record-imp456", d.Id())
	assert.Equal(t, "zone-imp123", d.Get("zone_id"))
	assert.Equal(t, "record-imp456", d.Get("record_id"))
}
