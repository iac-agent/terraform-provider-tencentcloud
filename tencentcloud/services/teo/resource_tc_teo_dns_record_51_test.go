package teo_test

import (
	"context"
	"errors"
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

// mockMetaTeoDnsRecord51 implements tccommon.ProviderMeta
type mockMetaTeoDnsRecord51 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoDnsRecord51) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoDnsRecord51{}

func newMockMetaTeoDnsRecord51() *mockMetaTeoDnsRecord51 {
	return &mockMetaTeoDnsRecord51{client: &connectivity.TencentCloudClient{}}
}

func ptrTeoDnsRecord51String(s string) *string {
	return &s
}

func ptrTeoDnsRecord51Int64(n int64) *int64 {
	return &n
}

// buildTeoDnsRecord51 builds a full DnsRecord element returned by the mocked DescribeDnsRecords.
func buildTeoDnsRecord51(recordId string) *teov20220901.DnsRecord {
	return &teov20220901.DnsRecord{
		ZoneId:     ptrTeoDnsRecord51String("zone-39quuimqg8r6"),
		RecordId:   ptrTeoDnsRecord51String(recordId),
		Name:       ptrTeoDnsRecord51String("a.makn.cn"),
		Type:       ptrTeoDnsRecord51String("A"),
		Location:   ptrTeoDnsRecord51String("Default"),
		Content:    ptrTeoDnsRecord51String("1.2.3.5"),
		TTL:        ptrTeoDnsRecord51Int64(300),
		Weight:     ptrTeoDnsRecord51Int64(-1),
		Priority:   ptrTeoDnsRecord51Int64(0),
		Status:     ptrTeoDnsRecord51String("enable"),
		CreatedOn:  ptrTeoDnsRecord51String("2024-01-01T00:00:00+08:00"),
		ModifiedOn: ptrTeoDnsRecord51String("2024-01-02T00:00:00+08:00"),
	}
}

// newTeoDnsRecord51MockClient patches UseTeoClient and UseTeoV20220901Client to return an empty teo client,
// so that both the service layer (UseTeoClient) and the resource CRUD callbacks (UseTeoV20220901Client)
// hit the same mocked teo client.
func newTeoDnsRecord51MockClient(patches *gomonkey.Patches) *teov20220901.Client {
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecord51().client, "UseTeoClient", teoClient)
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecord51().client, "UseTeoV20220901Client", teoClient)
	return teoClient
}

// mockDescribeDnsRecords51 mocks DescribeDnsRecords to return the given records.
func mockDescribeDnsRecords51(patches *gomonkey.Patches, teoClient *teov20220901.Client, records []*teov20220901.DnsRecord, captured **teov20220901.DescribeDnsRecordsRequest, callCount *int) {
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		if captured != nil {
			*captured = request
		}
		if callCount != nil {
			*callCount++
		}
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrTeoDnsRecord51Int64(int64(len(records))),
			DnsRecords: records,
			RequestId:  ptrTeoDnsRecord51String("fake-request-id"),
		}
		return resp, nil
	})
}

// TestTeoDnsRecord51_Schema validates schema definition.
func TestTeoDnsRecord51_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecord51()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "type")
	assert.Contains(t, res.Schema, "content")
	assert.Contains(t, res.Schema, "location")
	assert.Contains(t, res.Schema, "ttl")
	assert.Contains(t, res.Schema, "weight")
	assert.Contains(t, res.Schema, "priority")
	assert.Contains(t, res.Schema, "record_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "created_on")
	assert.Contains(t, res.Schema, "modified_on")

	assert.True(t, res.Schema["zone_id"].Required)
	assert.True(t, res.Schema["zone_id"].ForceNew)
	assert.True(t, res.Schema["name"].Required)
	assert.True(t, res.Schema["type"].Required)
	assert.True(t, res.Schema["content"].Required)

	assert.True(t, res.Schema["location"].Optional)
	assert.True(t, res.Schema["location"].Computed)
	assert.True(t, res.Schema["ttl"].Optional)
	assert.True(t, res.Schema["ttl"].Computed)
	assert.True(t, res.Schema["weight"].Optional)
	assert.True(t, res.Schema["weight"].Computed)
	assert.True(t, res.Schema["priority"].Optional)
	assert.True(t, res.Schema["priority"].Computed)

	assert.True(t, res.Schema["record_id"].Computed)
	assert.True(t, res.Schema["status"].Computed)
	assert.True(t, res.Schema["created_on"].Computed)
	assert.True(t, res.Schema["modified_on"].Computed)

	assert.NotNil(t, res.Importer)
}

// TestTeoDnsRecord51_CreateSuccess tests Create sets the composite id and refreshes state via Read.
func TestTeoDnsRecord51_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)

	var createRequest *teov20220901.CreateDnsRecordRequest
	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		createRequest = request
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrTeoDnsRecord51String("record-9m3fx1be"),
			RequestId: ptrTeoDnsRecord51String("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeDnsRecords51(patches, teoClient, []*teov20220901.DnsRecord{buildTeoDnsRecord51("record-9m3fx1be")}, nil, nil)

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.makn.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, "zone-39quuimqg8r6#record-9m3fx1be", d.Id())

	assert.NotNil(t, createRequest)
	assert.Equal(t, "zone-39quuimqg8r6", *createRequest.ZoneId)
	assert.Equal(t, "a.makn.cn", *createRequest.Name)
	assert.Equal(t, "A", *createRequest.Type)
	assert.Equal(t, "1.2.3.5", *createRequest.Content)
	assert.Equal(t, "Default", *createRequest.Location)
	assert.Equal(t, int64(300), *createRequest.TTL)
	assert.Equal(t, int64(-1), *createRequest.Weight)
	assert.Equal(t, int64(0), *createRequest.Priority)

	assert.Equal(t, "zone-39quuimqg8r6", d.Get("zone_id"))
	assert.Equal(t, "a.makn.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 0, d.Get("priority"))
	assert.Equal(t, "record-9m3fx1be", d.Get("record_id"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "2024-01-01T00:00:00+08:00", d.Get("created_on"))
	assert.Equal(t, "2024-01-02T00:00:00+08:00", d.Get("modified_on"))
}

// TestTeoDnsRecord51_CreateEmptyRecordId tests Create fails and does not write an id when RecordId is nil.
func TestTeoDnsRecord51_CreateEmptyRecordId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  nil,
			RequestId: ptrTeoDnsRecord51String("fake-request-id"),
		}
		return resp, nil
	})

	describeCalled := false
	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		describeCalled = true
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrTeoDnsRecord51Int64(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrTeoDnsRecord51String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RecordId is nil or empty")
	assert.Equal(t, "", d.Id())
	assert.False(t, describeCalled, "Read should not be invoked when create fails")
}

// TestTeoDnsRecord51_CreateEmptyStringRecordId tests Create fails when RecordId is an empty string.
func TestTeoDnsRecord51_CreateEmptyStringRecordId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrTeoDnsRecord51String(""),
			RequestId: ptrTeoDnsRecord51String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RecordId is nil or empty")
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord51_ReadSuccess tests Read populates all fields when the record is found.
func TestTeoDnsRecord51_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)

	var describeRequest *teov20220901.DescribeDnsRecordsRequest
	mockDescribeDnsRecords51(patches, teoClient, []*teov20220901.DnsRecord{buildTeoDnsRecord51("record-9m3fx1be")}, &describeRequest, nil)

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-39quuimqg8r6#record-9m3fx1be")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, "zone-39quuimqg8r6#record-9m3fx1be", d.Id())
	assert.NotNil(t, describeRequest)
	assert.Equal(t, "zone-39quuimqg8r6", *describeRequest.ZoneId)
	assert.NotNil(t, describeRequest.Limit)
	assert.Equal(t, int64(1000), *describeRequest.Limit)
	assert.Len(t, describeRequest.Filters, 1)
	assert.Equal(t, "id", *describeRequest.Filters[0].Name)
	assert.Equal(t, "record-9m3fx1be", *describeRequest.Filters[0].Values[0])

	assert.Equal(t, "zone-39quuimqg8r6", d.Get("zone_id"))
	assert.Equal(t, "a.makn.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 0, d.Get("priority"))
	assert.Equal(t, "record-9m3fx1be", d.Get("record_id"))
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "2024-01-01T00:00:00+08:00", d.Get("created_on"))
	assert.Equal(t, "2024-01-02T00:00:00+08:00", d.Get("modified_on"))
}

// TestTeoDnsRecord51_ReadNotFound tests Read clears the id when the record is not found.
func TestTeoDnsRecord51_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)
	mockDescribeDnsRecords51(patches, teoClient, []*teov20220901.DnsRecord{}, nil, nil)

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-39quuimqg8r6#record-9m3fx1be")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecord51_UpdateModify tests Update calls ModifyDnsRecords with a single element containing only mutable fields.
func TestTeoDnsRecord51_UpdateModify(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)

	var modifyRequest *teov20220901.ModifyDnsRecordsRequest
	modifyCallCount := 0
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		modifyRequest = request
		modifyCallCount++
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrTeoDnsRecord51String("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeDnsRecords51(patches, teoClient, []*teov20220901.DnsRecord{buildTeoDnsRecord51("record-9m3fx1be")}, nil, nil)

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()

	// Build a prior state where content is the old value, and a new config where
	// only content changes, so that d.HasChange("content") returns true inside Update.
	state := &terraform.InstanceState{
		ID: "zone-39quuimqg8r6#record-9m3fx1be",
		Attributes: map[string]string{
			"id":       "zone-39quuimqg8r6#record-9m3fx1be",
			"zone_id":  "zone-39quuimqg8r6",
			"name":     "a.makn.cn",
			"type":     "A",
			"content":  "1.2.3.5",
			"location": "Default",
			"ttl":      "300",
			"weight":   "-1",
			"priority": "0",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.makn.cn",
		"type":     "A",
		"content":  "1.2.3.4",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})

	diff, err := res.Diff(nil, state, rawConfig, meta)
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, 1, modifyCallCount)
	assert.NotNil(t, modifyRequest)
	assert.Equal(t, "zone-39quuimqg8r6", *modifyRequest.ZoneId)
	assert.Len(t, modifyRequest.DnsRecords, 1)

	elem := modifyRequest.DnsRecords[0]
	assert.Equal(t, "record-9m3fx1be", *elem.RecordId)
	assert.Equal(t, "a.makn.cn", *elem.Name)
	assert.Equal(t, "A", *elem.Type)
	assert.Equal(t, "1.2.3.4", *elem.Content)
	assert.Nil(t, elem.ZoneId)
	assert.Nil(t, elem.Status)
	assert.Nil(t, elem.CreatedOn)
	assert.Nil(t, elem.ModifiedOn)
}

// TestTeoDnsRecord51_UpdateNoChange tests Update skips ModifyDnsRecords when no mutable field changes.
func TestTeoDnsRecord51_UpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)

	modifyCallCount := 0
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.ModifyDnsRecordsRequest) (*teov20220901.ModifyDnsRecordsResponse, error) {
		modifyCallCount++
		resp := teov20220901.NewModifyDnsRecordsResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsResponseParams{
			RequestId: ptrTeoDnsRecord51String("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeDnsRecords51(patches, teoClient, []*teov20220901.DnsRecord{buildTeoDnsRecord51("record-9m3fx1be")}, nil, nil)

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()

	// Build a prior state identical to the new config so that no mutable field changes.
	state := &terraform.InstanceState{
		ID: "zone-39quuimqg8r6#record-9m3fx1be",
		Attributes: map[string]string{
			"id":       "zone-39quuimqg8r6#record-9m3fx1be",
			"zone_id":  "zone-39quuimqg8r6",
			"name":     "a.makn.cn",
			"type":     "A",
			"content":  "1.2.3.5",
			"location": "Default",
			"ttl":      "300",
			"weight":   "-1",
			"priority": "0",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"zone_id":  "zone-39quuimqg8r6",
		"name":     "a.makn.cn",
		"type":     "A",
		"content":  "1.2.3.5",
		"location": "Default",
		"ttl":      300,
		"weight":   -1,
		"priority": 0,
	})

	diff, err := res.Diff(nil, state, rawConfig, meta)
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, modifyCallCount, "ModifyDnsRecords should not be called when no mutable field changed")
}

// TestTeoDnsRecord51_DeleteSuccess tests Delete calls DeleteDnsRecords with the parsed ZoneId and RecordIds.
func TestTeoDnsRecord51_DeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)

	var deleteRequest *teov20220901.DeleteDnsRecordsRequest
	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		deleteRequest = request
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrTeoDnsRecord51String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-39quuimqg8r6#record-9m3fx1be")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, deleteRequest)
	assert.Equal(t, "zone-39quuimqg8r6", *deleteRequest.ZoneId)
	assert.Len(t, deleteRequest.RecordIds, 1)
	assert.Equal(t, "record-9m3fx1be", *deleteRequest.RecordIds[0])
}

// TestTeoDnsRecord51_DeleteApiError tests Delete returns the wrapped error when the API fails.
func TestTeoDnsRecord51_DeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := newTeoDnsRecord51MockClient(patches)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		return nil, errors.New("internal error")
	})

	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
	})
	d.SetId("zone-39quuimqg8r6#record-9m3fx1be")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}

// TestTeoDnsRecord51_BrokenId tests Read/Update/Delete return an "id is broken" error when the composite id is invalid.
func TestTeoDnsRecord51_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// No API call should happen for a broken id, so no API mocks are applied.
	meta := newMockMetaTeoDnsRecord51()
	res := teo.ResourceTencentCloudTeoDnsRecord51()

	// Read with broken id.
	dRead := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
	})
	dRead.SetId("broken-id-without-separator")
	err := res.Read(dRead, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")

	// Update with broken id.
	dUpdate := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
	})
	dUpdate.SetId("broken#id#with#extra#separators")
	err = res.Update(dUpdate, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")

	// Delete with broken id.
	dDelete := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-39quuimqg8r6",
	})
	dDelete.SetId("broken-id-without-separator")
	err = res.Delete(dDelete, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}
