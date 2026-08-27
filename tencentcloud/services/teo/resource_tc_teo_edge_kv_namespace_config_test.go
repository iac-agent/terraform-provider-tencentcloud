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

type mockMetaEdgeKVNamespaceConfig struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaEdgeKVNamespaceConfig) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaEdgeKVNamespaceConfig{}

func newMockMetaEdgeKVNamespaceConfig() *mockMetaEdgeKVNamespaceConfig {
	return &mockMetaEdgeKVNamespaceConfig{client: &connectivity.TencentCloudClient{}}
}

func ptrStringEKVNC(s string) *string {
	return &s
}

func ptrInt64EKVNC(i int64) *int64 {
	return &i
}

func prepareResourceDataForUpdate(t *testing.T, res *schema.Resource, initialData map[string]interface{}, id string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, res.Schema, initialData)
	d.SetId(id)
	// Set prior state to match initial config, so only explicitly changed fields will show HasChange
	for k, v := range initialData {
		_ = d.Set(k, v)
	}
	return d
}

func TestTeoEdgeKVNamespaceConfig_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespaceConfig().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(1000), *request.Limit)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "namespace", *request.Filters[0].Name)
		assert.Equal(t, "test-ns", *request.Filters[0].Values[0])
		assert.Equal(t, false, *request.Filters[0].Fuzzy)

		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: ptrInt64EKVNC(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace:    ptrStringEKVNC("test-ns"),
					Remark:       ptrStringEKVNC("test remark"),
					Capacity:     ptrInt64EKVNC(1073741824),
					CapacityUsed: ptrInt64EKVNC(512),
					CreatedOn:    ptrStringEKVNC("2024-01-01T00:00:00Z"),
					ModifiedOn:   ptrStringEKVNC("2024-06-15T12:00:00Z"),
				},
			},
			RequestId: ptrStringEKVNC("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespaceConfig()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespaceConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns",
		"remark":    "test remark",
	})
	d.SetId("zone-test123#test-ns")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#test-ns", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "test-ns", d.Get("namespace"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, 1073741824, d.Get("capacity"))
	assert.Equal(t, 512, d.Get("capacity_used"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-06-15T12:00:00Z", d.Get("modified_on"))
}

func TestTeoEdgeKVNamespaceConfig_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespaceConfig().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount:   ptrInt64EKVNC(0),
			KVNamespaces: []*teov20220901.KVNamespace{},
			RequestId:    ptrStringEKVNC("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespaceConfig()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespaceConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns-notexist",
	})
	d.SetId("zone-test123#test-ns-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoEdgeKVNamespaceConfig_Update_Remark(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespaceConfig().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyEdgeKVNamespaceWithContext", func(_ context.Context, request *teov20220901.ModifyEdgeKVNamespaceRequest) (*teov20220901.ModifyEdgeKVNamespaceResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "test-ns", *request.Namespace)
		assert.NotNil(t, request.Remark)
		assert.Equal(t, "updated remark", *request.Remark)

		resp := teov20220901.NewModifyEdgeKVNamespaceResponse()
		resp.Response = &teov20220901.ModifyEdgeKVNamespaceResponseParams{
			RequestId: ptrStringEKVNC("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: ptrInt64EKVNC(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace:    ptrStringEKVNC("test-ns"),
					Remark:       ptrStringEKVNC("updated remark"),
					Capacity:     ptrInt64EKVNC(1073741824),
					CapacityUsed: ptrInt64EKVNC(0),
					CreatedOn:    ptrStringEKVNC("2024-01-01T00:00:00Z"),
					ModifiedOn:   ptrStringEKVNC("2024-06-15T12:00:00Z"),
				},
			},
			RequestId: ptrStringEKVNC("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespaceConfig()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespaceConfig()

	initialData := map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns",
		"remark":    "old remark",
	}
	d := prepareResourceDataForUpdate(t, res, initialData, "zone-test123#test-ns")

	// Change the remark to trigger HasChange
	_ = d.Set("remark", "updated remark")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#test-ns", d.Id())
	assert.Equal(t, "updated remark", d.Get("remark"))
}

func TestTeoEdgeKVNamespaceConfig_Update_ImmutableArgs(t *testing.T) {
	meta := newMockMetaEdgeKVNamespaceConfig()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespaceConfig()

	// Test that changing immutable args returns error
	t.Run("zone_id", func(t *testing.T) {
		initialData := map[string]interface{}{
			"zone_id":   "zone-test123",
			"namespace": "test-ns",
		}
		d := prepareResourceDataForUpdate(t, res, initialData, "zone-test123#test-ns")
		_ = d.Set("zone_id", "zone-new")
		err := res.Update(d, meta)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "zone_id")
		assert.Contains(t, err.Error(), "cannot be changed")
	})

	t.Run("namespace", func(t *testing.T) {
		initialData := map[string]interface{}{
			"zone_id":   "zone-test123",
			"namespace": "test-ns",
		}
		d := prepareResourceDataForUpdate(t, res, initialData, "zone-test123#test-ns")
		_ = d.Set("namespace", "ns-new")
		err := res.Update(d, meta)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "namespace")
		assert.Contains(t, err.Error(), "cannot be changed")
	})
}
