package cls_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

// mockMeta implements tccommon.ProviderMeta
type mockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMeta{}

func newMockMeta() *mockMeta {
	return &mockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string {
	return &s
}

func ptrUint64(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/cls/ -run "TestAccTencentCloudClsSearchViewV2" -v -count=1 -gcflags="all=-l"

func TestAccTencentCloudClsSearchViewV2_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateSearchViewWithContext", func(ctx interface{}, request *clsv20201016.CreateSearchViewRequest) (*clsv20201016.CreateSearchViewResponse, error) {
		assert.Equal(t, "logset-123", *request.LogsetId)
		assert.Equal(t, "ap-guangzhou", *request.LogsetRegion)
		assert.Equal(t, "test-view", *request.ViewName)
		assert.Equal(t, "log", *request.ViewType)
		assert.Equal(t, 1, len(request.Topics))
		assert.Equal(t, "ap-guangzhou", *request.Topics[0].Region)
		assert.Equal(t, "logset-123", *request.Topics[0].LogsetId)
		assert.Equal(t, "topic-456", *request.Topics[0].TopicId)
		assert.Equal(t, "test description", *request.Description)
		assert.Equal(t, "my-prefix", *request.ViewIdPrefix)

		resp := clsv20201016.NewCreateSearchViewResponse()
		resp.Response = &clsv20201016.CreateSearchViewResponseParams{
			ViewId:    ptrString("my-prefix-view"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeSearchViews for the Read after Create
	patches.ApplyMethodFunc(clsClient, "DescribeSearchViewsWithContext", func(ctx interface{}, request *clsv20201016.DescribeSearchViewsRequest) (*clsv20201016.DescribeSearchViewsResponse, error) {
		resp := clsv20201016.NewDescribeSearchViewsResponse()
		resp.Response = &clsv20201016.DescribeSearchViewsResponseParams{
			Infos: []*clsv20201016.SearchViewInfo{
				{
					ViewId:       ptrString("my-prefix-view"),
					ViewName:     ptrString("test-view"),
					ViewType:     ptrString("log"),
					LogsetId:     ptrString("logset-123"),
					LogsetRegion: ptrString("ap-guangzhou"),
					Topics: []*clsv20201016.ViewSearchTopic{
						{
							Region:   ptrString("ap-guangzhou"),
							LogsetId: ptrString("logset-123"),
							TopicId:  ptrString("topic-456"),
						},
					},
					Description: ptrString("test description"),
				},
			},
			Total:     ptrUint64(1),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := localcls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "test-view",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-456",
			},
		},
		"description":    "test description",
		"view_id_prefix": "my-prefix",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "my-prefix-view", d.Id())
	assert.Equal(t, "my-prefix-view", d.Get("view_id"))
	assert.Equal(t, "test-view", d.Get("view_name"))
	assert.Equal(t, "log", d.Get("view_type"))
	assert.Equal(t, "logset-123", d.Get("logset_id"))
	assert.Equal(t, "ap-guangzhou", d.Get("logset_region"))
	assert.Equal(t, "test description", d.Get("description"))
}

func TestAccTencentCloudClsSearchViewV2_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateSearchViewWithContext", func(ctx interface{}, request *clsv20201016.CreateSearchViewRequest) (*clsv20201016.CreateSearchViewResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid view name")
	})

	meta := newMockMeta()
	res := localcls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "test-view",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-456",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

func TestAccTencentCloudClsSearchViewV2_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeSearchViewsWithContext", func(ctx interface{}, request *clsv20201016.DescribeSearchViewsRequest) (*clsv20201016.DescribeSearchViewsResponse, error) {
		resp := clsv20201016.NewDescribeSearchViewsResponse()
		resp.Response = &clsv20201016.DescribeSearchViewsResponseParams{
			Infos: []*clsv20201016.SearchViewInfo{
				{
					ViewId:       ptrString("my-prefix-view"),
					ViewName:     ptrString("test-view"),
					ViewType:     ptrString("log"),
					LogsetId:     ptrString("logset-123"),
					LogsetRegion: ptrString("ap-guangzhou"),
					Topics: []*clsv20201016.ViewSearchTopic{
						{
							Region:   ptrString("ap-guangzhou"),
							LogsetId: ptrString("logset-123"),
							TopicId:  ptrString("topic-456"),
						},
					},
					Description: ptrString("test description"),
				},
			},
			Total:     ptrUint64(1),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := localcls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "test-view",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-456",
			},
		},
	})
	d.SetId("my-prefix-view")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "my-prefix-view", d.Id())
	assert.Equal(t, "test-view", d.Get("view_name"))
	assert.Equal(t, "log", d.Get("view_type"))
	assert.Equal(t, "logset-123", d.Get("logset_id"))
	assert.Equal(t, "ap-guangzhou", d.Get("logset_region"))
	assert.Equal(t, "test description", d.Get("description"))
	assert.Equal(t, "my-prefix-view", d.Get("view_id"))
}

func TestAccTencentCloudClsSearchViewV2_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeSearchViewsWithContext", func(ctx interface{}, request *clsv20201016.DescribeSearchViewsRequest) (*clsv20201016.DescribeSearchViewsResponse, error) {
		resp := clsv20201016.NewDescribeSearchViewsResponse()
		resp.Response = &clsv20201016.DescribeSearchViewsResponseParams{
			Infos:     []*clsv20201016.SearchViewInfo{},
			Total:     ptrUint64(0),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := localcls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "test-view",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-456",
			},
		},
	})
	d.SetId("my-prefix-view")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestAccTencentCloudClsSearchViewV2_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifySearchViewWithContext", func(ctx interface{}, request *clsv20201016.ModifySearchViewRequest) (*clsv20201016.ModifySearchViewResponse, error) {
		assert.Equal(t, "my-prefix-view", *request.ViewId)
		assert.Equal(t, "updated-view", *request.ViewName)
		assert.Equal(t, "metric", *request.ViewType)
		assert.Equal(t, "updated description", *request.Description)

		resp := clsv20201016.NewModifySearchViewResponse()
		resp.Response = &clsv20201016.ModifySearchViewResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeSearchViews for the Read after Update
	patches.ApplyMethodFunc(clsClient, "DescribeSearchViewsWithContext", func(ctx interface{}, request *clsv20201016.DescribeSearchViewsRequest) (*clsv20201016.DescribeSearchViewsResponse, error) {
		resp := clsv20201016.NewDescribeSearchViewsResponse()
		resp.Response = &clsv20201016.DescribeSearchViewsResponseParams{
			Infos: []*clsv20201016.SearchViewInfo{
				{
					ViewId:       ptrString("my-prefix-view"),
					ViewName:     ptrString("updated-view"),
					ViewType:     ptrString("metric"),
					LogsetId:     ptrString("logset-123"),
					LogsetRegion: ptrString("ap-guangzhou"),
					Topics: []*clsv20201016.ViewSearchTopic{
						{
							Region:   ptrString("ap-guangzhou"),
							LogsetId: ptrString("logset-123"),
							TopicId:  ptrString("topic-789"),
						},
					},
					Description: ptrString("updated description"),
				},
			},
			Total:     ptrUint64(1),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := localcls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "updated-view",
		"view_type":     "metric",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-789",
			},
		},
		"description": "updated description",
	})
	d.SetId("my-prefix-view")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "updated-view", d.Get("view_name"))
	assert.Equal(t, "metric", d.Get("view_type"))
	assert.Equal(t, "updated description", d.Get("description"))
}

func TestAccTencentCloudClsSearchViewV2_UpdateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifySearchViewWithContext", func(ctx interface{}, request *clsv20201016.ModifySearchViewRequest) (*clsv20201016.ModifySearchViewResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid view name")
	})

	meta := newMockMeta()
	res := localcls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "updated-view",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-456",
			},
		},
	})
	d.SetId("my-prefix-view")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

func TestAccTencentCloudClsSearchViewV2_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteSearchViewWithContext", func(ctx interface{}, request *clsv20201016.DeleteSearchViewRequest) (*clsv20201016.DeleteSearchViewResponse, error) {
		assert.Equal(t, "my-prefix-view", *request.ViewId)

		resp := clsv20201016.NewDeleteSearchViewResponse()
		resp.Response = &clsv20201016.DeleteSearchViewResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := localcls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "test-view",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-456",
			},
		},
	})
	d.SetId("my-prefix-view")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestAccTencentCloudClsSearchViewV2_DeleteAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteSearchViewWithContext", func(ctx interface{}, request *clsv20201016.DeleteSearchViewRequest) (*clsv20201016.DeleteSearchViewResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=View not found")
	})

	meta := newMockMeta()
	res := localcls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "test-view",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-456",
			},
		},
	})
	d.SetId("my-prefix-view")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

func TestAccTencentCloudClsSearchViewV2_Schema(t *testing.T) {
	res := localcls.ResourceTencentCloudClsSearchViewV2()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "logset_id")
	assert.Contains(t, res.Schema, "logset_region")
	assert.Contains(t, res.Schema, "view_name")
	assert.Contains(t, res.Schema, "view_type")
	assert.Contains(t, res.Schema, "topics")
	assert.Contains(t, res.Schema, "description")
	assert.Contains(t, res.Schema, "view_id_prefix")
	assert.Contains(t, res.Schema, "view_id")

	logsetId := res.Schema["logset_id"]
	assert.Equal(t, schema.TypeString, logsetId.Type)
	assert.True(t, logsetId.Required)
	assert.True(t, logsetId.ForceNew)

	logsetRegion := res.Schema["logset_region"]
	assert.Equal(t, schema.TypeString, logsetRegion.Type)
	assert.True(t, logsetRegion.Required)
	assert.True(t, logsetRegion.ForceNew)

	viewName := res.Schema["view_name"]
	assert.Equal(t, schema.TypeString, viewName.Type)
	assert.True(t, viewName.Required)
	assert.False(t, viewName.ForceNew)

	viewType := res.Schema["view_type"]
	assert.Equal(t, schema.TypeString, viewType.Type)
	assert.True(t, viewType.Required)
	assert.False(t, viewType.ForceNew)

	topics := res.Schema["topics"]
	assert.Equal(t, schema.TypeList, topics.Type)
	assert.True(t, topics.Required)

	description := res.Schema["description"]
	assert.Equal(t, schema.TypeString, description.Type)
	assert.True(t, description.Optional)

	viewIdPrefix := res.Schema["view_id_prefix"]
	assert.Equal(t, schema.TypeString, viewIdPrefix.Type)
	assert.True(t, viewIdPrefix.Optional)
	assert.True(t, viewIdPrefix.ForceNew)

	viewId := res.Schema["view_id"]
	assert.Equal(t, schema.TypeString, viewId.Type)
	assert.True(t, viewId.Computed)
}
