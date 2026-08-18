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

type mockMetaTeoDnsRecordV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoDnsRecordV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoDnsRecordV2{}

func newMockMetaTeoDnsRecordV2() *mockMetaTeoDnsRecordV2 {
	return &mockMetaTeoDnsRecordV2{client: &connectivity.TencentCloudClient{}}
}

func ptrTeoDnsRecordV2String(s string) *string { return &s }

func ptrTeoDnsRecordV2Int64(n int64) *int64 { return &n }

func TestTeoDnsRecordV2_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecordV2()

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
}

func TestTeoDnsRecordV2_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecordV2()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrTeoDnsRecordV2String("record-abc123"),
			RequestId: ptrTeoDnsRecordV2String("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrTeoDnsRecordV2Int64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:   ptrTeoDnsRecordV2String("zone-abc123"),
					RecordId: ptrTeoDnsRecordV2String("record-abc123"),
					Name:     ptrTeoDnsRecordV2String("a.makn.cn"),
					Type:     ptrTeoDnsRecordV2String("A"),
					Content:  ptrTeoDnsRecordV2String("1.2.3.5"),
					Location: ptrTeoDnsRecordV2String("Default"),
					TTL:      ptrTeoDnsRecordV2Int64(300),
					Weight:   ptrTeoDnsRecordV2Int64(-1),
					Priority: ptrTeoDnsRecordV2Int64(5),
				},
			},
			RequestId: ptrTeoDnsRecordV2String("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecordV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-abc123",
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
	assert.Equal(t, "zone-abc123#record-abc123", d.Id())
	assert.Equal(t, "zone-abc123", d.Get("zone_id"))
	assert.Equal(t, "a.makn.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 5, d.Get("priority"))
	assert.Equal(t, "record-abc123", d.Get("record_id"))
}

func TestTeoDnsRecordV2_CreateEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecordV2()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		return nil, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecordV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-abc123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoDnsRecordV2_CreateEmptyRecordId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecordV2()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateDnsRecordWithContext", func(ctx context.Context, request *teov20220901.CreateDnsRecordRequest) (*teov20220901.CreateDnsRecordResponse, error) {
		resp := teov20220901.NewCreateDnsRecordResponse()
		resp.Response = &teov20220901.CreateDnsRecordResponseParams{
			RecordId:  ptrTeoDnsRecordV2String(""),
			RequestId: ptrTeoDnsRecordV2String("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecordV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-abc123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoDnsRecordV2_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecordV2()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrTeoDnsRecordV2Int64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:   ptrTeoDnsRecordV2String("zone-abc123"),
					RecordId: ptrTeoDnsRecordV2String("record-abc123"),
					Name:     ptrTeoDnsRecordV2String("a.makn.cn"),
					Type:     ptrTeoDnsRecordV2String("A"),
					Content:  ptrTeoDnsRecordV2String("1.2.3.5"),
					Location: ptrTeoDnsRecordV2String("Default"),
					TTL:      ptrTeoDnsRecordV2Int64(300),
					Weight:   ptrTeoDnsRecordV2Int64(-1),
					Priority: ptrTeoDnsRecordV2Int64(5),
				},
			},
			RequestId: ptrTeoDnsRecordV2String("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecordV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-abc123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-abc123#record-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-abc123#record-abc123", d.Id())
	assert.Equal(t, "zone-abc123", d.Get("zone_id"))
	assert.Equal(t, "a.makn.cn", d.Get("name"))
	assert.Equal(t, "A", d.Get("type"))
	assert.Equal(t, "1.2.3.5", d.Get("content"))
	assert.Equal(t, "Default", d.Get("location"))
	assert.Equal(t, 300, d.Get("ttl"))
	assert.Equal(t, -1, d.Get("weight"))
	assert.Equal(t, 5, d.Get("priority"))
	assert.Equal(t, "record-abc123", d.Get("record_id"))
}

func TestTeoDnsRecordV2_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecordV2()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrTeoDnsRecordV2Int64(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrTeoDnsRecordV2String("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecordV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-abc123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-abc123#record-not-found")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoDnsRecordV2_DeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecordV2()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteDnsRecordsWithContext", func(ctx context.Context, request *teov20220901.DeleteDnsRecordsRequest) (*teov20220901.DeleteDnsRecordsResponse, error) {
		resp := teov20220901.NewDeleteDnsRecordsResponse()
		resp.Response = &teov20220901.DeleteDnsRecordsResponseParams{
			RequestId: ptrTeoDnsRecordV2String("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecordV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-abc123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-abc123#record-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestTeoDnsRecordV2_UpdateImmutableChanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teo.ResourceTencentCloudTeoDnsRecordV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-abc123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-abc123#record-abc123")

	// Simulate HasChange for an immutable field.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "name"
	})

	err := res.Update(d, newMockMetaTeoDnsRecordV2())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "argument `name` cannot be changed")
}

func TestTeoDnsRecordV2_UpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoDnsRecordV2()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecords", func(request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrTeoDnsRecordV2Int64(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:   ptrTeoDnsRecordV2String("zone-abc123"),
					RecordId: ptrTeoDnsRecordV2String("record-abc123"),
					Name:     ptrTeoDnsRecordV2String("a.makn.cn"),
					Type:     ptrTeoDnsRecordV2String("A"),
					Content:  ptrTeoDnsRecordV2String("1.2.3.5"),
				},
			},
			RequestId: ptrTeoDnsRecordV2String("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoDnsRecordV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-abc123",
		"name":    "a.makn.cn",
		"type":    "A",
		"content": "1.2.3.5",
	})
	d.SetId("zone-abc123#record-abc123")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-abc123#record-abc123", d.Id())
	assert.Equal(t, "a.makn.cn", d.Get("name"))
}
