package config_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

type mockMetaConfigListAggregators struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaConfigListAggregators) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaConfigListAggregators{}

func newMockMetaConfigListAggregators() *mockMetaConfigListAggregators {
	return &mockMetaConfigListAggregators{client: &connectivity.TencentCloudClient{}}
}

func ptrStringConfigAggregator(s string) *string {
	return &s
}

func ptrUint64ConfigAggregator(i uint64) *uint64 {
	return &i
}

// TestConfigListAggregatorsDS_ReadSuccess tests data source Read with successful response
func TestConfigListAggregatorsDS_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigListAggregators().client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregatorsWithContext", func(_ context.Context, request *configv20220802.ListAggregatorsRequest) (*configv20220802.ListAggregatorsResponse, error) {
		resp := configv20220802.NewListAggregatorsResponse()
		resp.Response = &configv20220802.ListAggregatorsResponseParams{
			Total: ptrUint64ConfigAggregator(1),
			Items: []*configv20220802.Aggregator{
				{
					Name:             ptrStringConfigAggregator("group-a"),
					Description:      ptrStringConfigAggregator("test group"),
					OwnerUin:         ptrUint64ConfigAggregator(100000000001),
					CreateTime:       ptrStringConfigAggregator("2024-01-01 00:00:00"),
					AccountCount:     ptrUint64ConfigAggregator(3),
					Type:             ptrStringConfigAggregator("RD"),
					AccountGroupId:   ptrStringConfigAggregator("aggr-xxxxxxxx"),
					AggregatorStatus: ptrUint64ConfigAggregator(1),
					MemberName:       ptrStringConfigAggregator("group-a"),
				},
			},
			RequestId: ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigListAggregators()
	res := config.DataSourceTencentCloudConfigListAggregators()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	diags := res.ReadContext(context.TODO(), d, meta)
	assert.False(t, diags.HasError())
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, 1, d.Get("total"))

	items := d.Get("items").([]interface{})
	assert.Equal(t, 1, len(items))

	item0 := items[0].(map[string]interface{})
	assert.Equal(t, "group-a", item0["name"])
	assert.Equal(t, "test group", item0["description"])
	assert.Equal(t, 100000000001, item0["owner_uin"])
	assert.Equal(t, "2024-01-01 00:00:00", item0["create_time"])
	assert.Equal(t, 3, item0["account_count"])
	assert.Equal(t, "RD", item0["type"])
	assert.Equal(t, "aggr-xxxxxxxx", item0["account_group_id"])
	assert.Equal(t, 1, item0["aggregator_status"])
	assert.Equal(t, "group-a", item0["member_name"])
}

// TestConfigListAggregatorsDS_ReadNilResponse tests data source Read with nil response
func TestConfigListAggregatorsDS_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigListAggregators().client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregatorsWithContext", func(_ context.Context, request *configv20220802.ListAggregatorsRequest) (*configv20220802.ListAggregatorsResponse, error) {
		resp := configv20220802.NewListAggregatorsResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaConfigListAggregators()
	res := config.DataSourceTencentCloudConfigListAggregators()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	diags := res.ReadContext(context.TODO(), d, meta)
	assert.True(t, diags.HasError())
}
