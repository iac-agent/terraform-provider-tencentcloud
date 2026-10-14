package config_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	configsvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

type mockMetaConfigAggregator struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaConfigAggregator) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaConfigAggregator{}

func newMockMetaConfigAggregator() *mockMetaConfigAggregator {
	return &mockMetaConfigAggregator{client: &connectivity.TencentCloudClient{}}
}

func ptrStringConfigAggregator(s string) *string {
	return &s
}

func ptrUint64ConfigAggregator(v uint64) *uint64 {
	return &v
}

// TestConfigAggregator_Create verifies Create passes fields and sets composite id.
func TestConfigAggregator_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigAggregator().client, "UseConfigV20220802Client", client)

	var capturedCreate *configv20220802.CreateAggregatorRequest
	patches.ApplyMethodFunc(client, "CreateAggregatorWithContext", func(_ interface{}, request *configv20220802.CreateAggregatorRequest) (*configv20220802.CreateAggregatorResponse, error) {
		capturedCreate = request
		resp := configv20220802.NewCreateAggregatorResponse()
		resp.Response = &configv20220802.CreateAggregatorResponseParams{
			AccountGroupId: ptrStringConfigAggregator("aggr-123456"),
			RequestId:      ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(client, "DescribeAggregatorWithContext", func(_ interface{}, request *configv20220802.DescribeAggregatorRequest) (*configv20220802.DescribeAggregatorResponse, error) {
		resp := configv20220802.NewDescribeAggregatorResponse()
		resp.Response = &configv20220802.DescribeAggregatorResponseParams{
			Name:        ptrStringConfigAggregator("tf_aggregator"),
			Description: ptrStringConfigAggregator("created by tf"),
			Type:        ptrStringConfigAggregator("CUSTOM"),
			AggregatorAccounts: []*configv20220802.AggregatorAccount{
				{
					MemberUin:  ptrUint64ConfigAggregator(10001),
					MemberName: ptrStringConfigAggregator("member1"),
				},
			},
			AggregatorStatus: ptrUint64ConfigAggregator(1),
			RequestId:        ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigAggregator()
	res := configsvc.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "tf_aggregator",
		"description": "created by tf",
		"type":        "CUSTOM",
		"owner_uin":   10000,
		"aggregator_accounts": []interface{}{
			map[string]interface{}{
				"member_uin":  10001,
				"member_name": "member1",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "10000"+tccommon.FILED_SP+"aggr-123456", d.Id())

	assert.NotNil(t, capturedCreate)
	assert.Equal(t, "tf_aggregator", *capturedCreate.Name)
	assert.Equal(t, "created by tf", *capturedCreate.Description)
	assert.Equal(t, "CUSTOM", *capturedCreate.Type)
	assert.Len(t, capturedCreate.AggregatorAccounts, 1)
	assert.Equal(t, uint64(10001), *capturedCreate.AggregatorAccounts[0].MemberUin)

	// Read back into state.
	assert.Equal(t, "tf_aggregator", d.Get("name").(string))
	assert.Equal(t, 1, d.Get("aggregator_status").(int))
	assert.Equal(t, "aggr-123456", d.Get("account_group_id").(string))
}

// TestConfigAggregator_CreateEmptyId verifies Create fails when AccountGroupId is empty.
func TestConfigAggregator_CreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigAggregator().client, "UseConfigV20220802Client", client)
	patches.ApplyMethodFunc(client, "CreateAggregatorWithContext", func(_ interface{}, request *configv20220802.CreateAggregatorRequest) (*configv20220802.CreateAggregatorResponse, error) {
		resp := configv20220802.NewCreateAggregatorResponse()
		resp.Response = &configv20220802.CreateAggregatorResponseParams{
			RequestId: ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigAggregator()
	res := configsvc.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "tf_aggregator",
		"description": "created by tf",
		"type":        "CUSTOM",
		"owner_uin":   10000,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestConfigAggregator_Read verifies Read fills state from composite id.
func TestConfigAggregator_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigAggregator().client, "UseConfigV20220802Client", client)

	var capturedRequest *configv20220802.DescribeAggregatorRequest
	patches.ApplyMethodFunc(client, "DescribeAggregatorWithContext", func(_ interface{}, request *configv20220802.DescribeAggregatorRequest) (*configv20220802.DescribeAggregatorResponse, error) {
		capturedRequest = request
		resp := configv20220802.NewDescribeAggregatorResponse()
		resp.Response = &configv20220802.DescribeAggregatorResponseParams{
			Name:        ptrStringConfigAggregator("tf_aggregator"),
			Description: ptrStringConfigAggregator("created by tf"),
			Type:        ptrStringConfigAggregator("CUSTOM"),
			AggregatorAccounts: []*configv20220802.AggregatorAccount{
				{
					MemberUin:  ptrUint64ConfigAggregator(10001),
					MemberName: ptrStringConfigAggregator("member1"),
				},
			},
			AggregatorStatus: ptrUint64ConfigAggregator(1),
			RequestId:        ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigAggregator()
	res := configsvc.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "old",
		"description": "old",
		"type":        "CUSTOM",
		"owner_uin":   10000,
	})
	d.SetId("10000" + tccommon.FILED_SP + "aggr-123456")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "aggr-123456", *capturedRequest.AccountGroupId)
	assert.Equal(t, uint64(10000), *capturedRequest.OwnerUin)

	assert.Equal(t, "tf_aggregator", d.Get("name").(string))
	assert.Equal(t, "created by tf", d.Get("description").(string))
	assert.Equal(t, 1, d.Get("aggregator_status").(int))
	accounts := d.Get("aggregator_accounts").([]interface{})
	assert.Len(t, accounts, 1)
	assert.Equal(t, 10001, accounts[0].(map[string]interface{})["member_uin"].(int))
}

// TestConfigAggregator_ReadNotFound verifies Read clears id when resource is gone.
func TestConfigAggregator_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigAggregator().client, "UseConfigV20220802Client", client)
	patches.ApplyMethodFunc(client, "DescribeAggregatorWithContext", func(_ interface{}, request *configv20220802.DescribeAggregatorRequest) (*configv20220802.DescribeAggregatorResponse, error) {
		resp := configv20220802.NewDescribeAggregatorResponse()
		resp.Response = &configv20220802.DescribeAggregatorResponseParams{
			RequestId: ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigAggregator()
	res := configsvc.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "tf_aggregator",
		"description": "created by tf",
		"type":        "CUSTOM",
		"owner_uin":   10000,
	})
	d.SetId("10000" + tccommon.FILED_SP + "aggr-123456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestConfigAggregator_Update verifies Update sends fields from id split.
func TestConfigAggregator_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigAggregator().client, "UseConfigV20220802Client", client)

	var capturedUpdate *configv20220802.UpdateAggregatorRequest
	patches.ApplyMethodFunc(client, "UpdateAggregatorWithContext", func(_ interface{}, request *configv20220802.UpdateAggregatorRequest) (*configv20220802.UpdateAggregatorResponse, error) {
		capturedUpdate = request
		resp := configv20220802.NewUpdateAggregatorResponse()
		resp.Response = &configv20220802.UpdateAggregatorResponseParams{
			RequestId: ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(client, "DescribeAggregatorWithContext", func(_ interface{}, request *configv20220802.DescribeAggregatorRequest) (*configv20220802.DescribeAggregatorResponse, error) {
		resp := configv20220802.NewDescribeAggregatorResponse()
		resp.Response = &configv20220802.DescribeAggregatorResponseParams{
			Name:             ptrStringConfigAggregator("tf_aggregator_new"),
			Description:      ptrStringConfigAggregator("updated"),
			Type:             ptrStringConfigAggregator("CUSTOM"),
			AggregatorStatus: ptrUint64ConfigAggregator(1),
			RequestId:        ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigAggregator()
	res := configsvc.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "tf_aggregator_new",
		"description": "updated",
		"type":        "CUSTOM",
		"owner_uin":   10000,
	})
	d.SetId("10000" + tccommon.FILED_SP + "aggr-123456")

	// Simulate change detection (TestResourceDataRaw treats configured keys as changed).
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "name" || key == "description"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedUpdate)
	assert.Equal(t, "tf_aggregator_new", *capturedUpdate.Name)
	assert.Equal(t, "updated", *capturedUpdate.Description)
	assert.Equal(t, "aggr-123456", *capturedUpdate.AccountGroupId)
	assert.Equal(t, uint64(10000), *capturedUpdate.OwnerUin)

	assert.Equal(t, "tf_aggregator_new", d.Get("name").(string))
}

// TestConfigAggregator_Delete verifies Delete sends id fields.
func TestConfigAggregator_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockMetaConfigAggregator().client, "UseConfigV20220802Client", client)

	var capturedDelete *configv20220802.DeleteAggregatorsRequest
	patches.ApplyMethodFunc(client, "DeleteAggregatorsWithContext", func(_ interface{}, request *configv20220802.DeleteAggregatorsRequest) (*configv20220802.DeleteAggregatorsResponse, error) {
		capturedDelete = request
		resp := configv20220802.NewDeleteAggregatorsResponse()
		resp.Response = &configv20220802.DeleteAggregatorsResponseParams{
			RequestId: ptrStringConfigAggregator("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfigAggregator()
	res := configsvc.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "tf_aggregator",
		"description": "created by tf",
		"type":        "CUSTOM",
		"owner_uin":   10000,
	})
	d.SetId("10000" + tccommon.FILED_SP + "aggr-123456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedDelete)
	assert.Equal(t, "aggr-123456", *capturedDelete.AccountGroupId)
	assert.Equal(t, uint64(10000), *capturedDelete.OwnerUin)
}

// TestConfigAggregator_Schema validates the schema definition.
func TestConfigAggregator_Schema(t *testing.T) {
	res := configsvc.ResourceTencentCloudConfigAggregator()

	for _, key := range []string{"name", "description", "type", "owner_uin"} {
		field, ok := res.Schema[key]
		assert.True(t, ok, "%s should exist in schema", key)
		assert.True(t, field.Required, "%s should be required", key)
	}

	accounts, ok := res.Schema["aggregator_accounts"]
	assert.True(t, ok)
	assert.True(t, accounts.Optional)
	assert.Equal(t, schema.TypeList, accounts.Type)

	for _, key := range []string{"account_group_id", "aggregator_status"} {
		field, ok := res.Schema[key]
		assert.True(t, ok, "%s should exist in schema", key)
		assert.True(t, field.Computed, "%s should be computed", key)
	}
}
