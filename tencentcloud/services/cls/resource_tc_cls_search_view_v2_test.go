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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

// mockMeta implements tccommon.ProviderMeta
type searchViewV2MockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *searchViewV2MockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &searchViewV2MockMeta{}

func newSearchViewV2MockMeta() *searchViewV2MockMeta {
	return &searchViewV2MockMeta{client: &connectivity.TencentCloudClient{}}
}

func searchViewV2PtrString(s string) *string {
	return &s
}

func searchViewV2PtrUint64(u uint64) *uint64 {
	return &u
}

// go test ./tencentcloud/services/cls/ -run "TestClsSearchViewV2Create" -v -count=1 -gcflags="all=-l"
func TestClsSearchViewV2Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newSearchViewV2MockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateSearchViewWithContext", func(ctx interface{}, request *clsv20201016.CreateSearchViewRequest) (*clsv20201016.CreateSearchViewResponse, error) {
		assert.Equal(t, "logset-123", *request.LogsetId)
		assert.Equal(t, "ap-guangzhou", *request.LogsetRegion)
		assert.Equal(t, "test-view", *request.ViewName)
		assert.Equal(t, "log", *request.ViewType)
		resp := clsv20201016.NewCreateSearchViewResponse()
		resp.Response = &clsv20201016.CreateSearchViewResponseParams{
			ViewId:    helper.String("view-abc123-view"),
			RequestId: searchViewV2PtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeSearchViewsWithContext", func(ctx interface{}, request *clsv20201016.DescribeSearchViewsRequest) (*clsv20201016.DescribeSearchViewsResponse, error) {
		resp := clsv20201016.NewDescribeSearchViewsResponse()
		resp.Response = &clsv20201016.DescribeSearchViewsResponseParams{
			Infos: []*clsv20201016.SearchViewInfo{
				{
					ViewId:       helper.String("view-abc123-view"),
					ViewName:     helper.String("test-view"),
					ViewType:     helper.String("log"),
					LogsetId:     helper.String("logset-123"),
					LogsetRegion: helper.String("ap-guangzhou"),
					Topics: []*clsv20201016.ViewSearchTopic{
						{
							Region:   helper.String("ap-guangzhou"),
							LogsetId: helper.String("logset-123"),
							TopicId:  helper.String("topic-456"),
						},
					},
					Description: helper.String("test description"),
					CreateTime:  searchViewV2PtrUint64(1700000000),
					UpdateTime:  searchViewV2PtrUint64(1700000001),
				},
			},
			Total:     searchViewV2PtrUint64(1),
			RequestId: searchViewV2PtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newSearchViewV2MockMeta()
	res := cls.ResourceTencentCloudClsSearchViewV2()
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
		"description": "test description",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "view-abc123-view", d.Id())
	assert.Equal(t, "view-abc123-view", d.Get("view_id"))
}

// TestClsSearchViewV2Create_APIError tests Create handles API error
func TestClsSearchViewV2Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newSearchViewV2MockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateSearchViewWithContext", func(ctx interface{}, request *clsv20201016.CreateSearchViewRequest) (*clsv20201016.CreateSearchViewResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid logset_id")
	})

	meta := newSearchViewV2MockMeta()
	res := cls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-invalid",
		"logset_region": "ap-guangzhou",
		"view_name":     "test-view",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-invalid",
				"topic_id":  "topic-456",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestClsSearchViewV2Read_Success tests Read populates state
func TestClsSearchViewV2Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newSearchViewV2MockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeSearchViewsWithContext", func(ctx interface{}, request *clsv20201016.DescribeSearchViewsRequest) (*clsv20201016.DescribeSearchViewsResponse, error) {
		resp := clsv20201016.NewDescribeSearchViewsResponse()
		resp.Response = &clsv20201016.DescribeSearchViewsResponseParams{
			Infos: []*clsv20201016.SearchViewInfo{
				{
					ViewId:       helper.String("view-abc123-view"),
					ViewName:     helper.String("test-view"),
					ViewType:     helper.String("log"),
					LogsetId:     helper.String("logset-123"),
					LogsetRegion: helper.String("ap-guangzhou"),
					Topics: []*clsv20201016.ViewSearchTopic{
						{
							Region:   helper.String("ap-guangzhou"),
							LogsetId: helper.String("logset-123"),
							TopicId:  helper.String("topic-456"),
						},
					},
					Description: helper.String("test description"),
					CreateTime:  searchViewV2PtrUint64(1700000000),
					UpdateTime:  searchViewV2PtrUint64(1700000001),
				},
			},
			Total:     searchViewV2PtrUint64(1),
			RequestId: searchViewV2PtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newSearchViewV2MockMeta()
	res := cls.ResourceTencentCloudClsSearchViewV2()
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
	d.SetId("view-abc123-view")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "view-abc123-view", d.Id())
	assert.Equal(t, "test-view", d.Get("view_name"))
	assert.Equal(t, "log", d.Get("view_type"))
	assert.Equal(t, "logset-123", d.Get("logset_id"))
	assert.Equal(t, "ap-guangzhou", d.Get("logset_region"))
	assert.Equal(t, "test description", d.Get("description"))
}

// TestClsSearchViewV2Read_NotFound tests Read sets empty ID when resource not found
func TestClsSearchViewV2Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newSearchViewV2MockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeSearchViewsWithContext", func(ctx interface{}, request *clsv20201016.DescribeSearchViewsRequest) (*clsv20201016.DescribeSearchViewsResponse, error) {
		resp := clsv20201016.NewDescribeSearchViewsResponse()
		resp.Response = &clsv20201016.DescribeSearchViewsResponseParams{
			Infos:     []*clsv20201016.SearchViewInfo{},
			Total:     searchViewV2PtrUint64(0),
			RequestId: searchViewV2PtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newSearchViewV2MockMeta()
	res := cls.ResourceTencentCloudClsSearchViewV2()
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
	d.SetId("view-not-exist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestClsSearchViewV2Update_Success tests Update modifies mutable fields
func TestClsSearchViewV2Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newSearchViewV2MockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifySearchViewWithContext", func(ctx interface{}, request *clsv20201016.ModifySearchViewRequest) (*clsv20201016.ModifySearchViewResponse, error) {
		assert.Equal(t, "view-abc123-view", *request.ViewId)
		assert.Equal(t, "test-view-updated", *request.ViewName)
		resp := clsv20201016.NewModifySearchViewResponse()
		resp.Response = &clsv20201016.ModifySearchViewResponseParams{
			RequestId: searchViewV2PtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeSearchViewsWithContext", func(ctx interface{}, request *clsv20201016.DescribeSearchViewsRequest) (*clsv20201016.DescribeSearchViewsResponse, error) {
		resp := clsv20201016.NewDescribeSearchViewsResponse()
		resp.Response = &clsv20201016.DescribeSearchViewsResponseParams{
			Infos: []*clsv20201016.SearchViewInfo{
				{
					ViewId:       helper.String("view-abc123-view"),
					ViewName:     helper.String("test-view-updated"),
					ViewType:     helper.String("log"),
					LogsetId:     helper.String("logset-123"),
					LogsetRegion: helper.String("ap-guangzhou"),
					Topics: []*clsv20201016.ViewSearchTopic{
						{
							Region:   helper.String("ap-guangzhou"),
							LogsetId: helper.String("logset-123"),
							TopicId:  helper.String("topic-456"),
						},
					},
					Description: helper.String("updated description"),
					CreateTime:  searchViewV2PtrUint64(1700000000),
					UpdateTime:  searchViewV2PtrUint64(1700000002),
				},
			},
			Total:     searchViewV2PtrUint64(1),
			RequestId: searchViewV2PtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newSearchViewV2MockMeta()
	res := cls.ResourceTencentCloudClsSearchViewV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"logset_id":     "logset-123",
		"logset_region": "ap-guangzhou",
		"view_name":     "test-view-updated",
		"view_type":     "log",
		"topics": []interface{}{
			map[string]interface{}{
				"region":    "ap-guangzhou",
				"logset_id": "logset-123",
				"topic_id":  "topic-456",
			},
		},
		"description": "updated description",
	})
	d.SetId("view-abc123-view")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-view-updated", d.Get("view_name"))
}

// TestClsSearchViewV2Delete_Success tests Delete removes the resource
func TestClsSearchViewV2Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newSearchViewV2MockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteSearchViewWithContext", func(ctx interface{}, request *clsv20201016.DeleteSearchViewRequest) (*clsv20201016.DeleteSearchViewResponse, error) {
		assert.Equal(t, "view-abc123-view", *request.ViewId)
		resp := clsv20201016.NewDeleteSearchViewResponse()
		resp.Response = &clsv20201016.DeleteSearchViewResponseParams{
			RequestId: searchViewV2PtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newSearchViewV2MockMeta()
	res := cls.ResourceTencentCloudClsSearchViewV2()
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
	d.SetId("view-abc123-view")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestClsSearchViewV2Delete_APIError tests Delete handles API error
func TestClsSearchViewV2Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newSearchViewV2MockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteSearchViewWithContext", func(ctx interface{}, request *clsv20201016.DeleteSearchViewRequest) (*clsv20201016.DeleteSearchViewResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=View not found")
	})

	meta := newSearchViewV2MockMeta()
	res := cls.ResourceTencentCloudClsSearchViewV2()
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
	d.SetId("view-not-exist")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestClsSearchViewV2_Schema validates schema definition
func TestClsSearchViewV2_Schema(t *testing.T) {
	res := cls.ResourceTencentCloudClsSearchViewV2()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "view_id")
	assert.Contains(t, res.Schema, "logset_id")
	assert.Contains(t, res.Schema, "logset_region")
	assert.Contains(t, res.Schema, "view_name")
	assert.Contains(t, res.Schema, "view_type")
	assert.Contains(t, res.Schema, "topics")
	assert.Contains(t, res.Schema, "view_id_prefix")
	assert.Contains(t, res.Schema, "description")
	assert.Contains(t, res.Schema, "create_time")
	assert.Contains(t, res.Schema, "update_time")

	viewId := res.Schema["view_id"]
	assert.Equal(t, schema.TypeString, viewId.Type)
	assert.True(t, viewId.Computed)

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

	viewType := res.Schema["view_type"]
	assert.Equal(t, schema.TypeString, viewType.Type)
	assert.True(t, viewType.Required)

	topics := res.Schema["topics"]
	assert.Equal(t, schema.TypeList, topics.Type)
	assert.True(t, topics.Required)

	viewIdPrefix := res.Schema["view_id_prefix"]
	assert.Equal(t, schema.TypeString, viewIdPrefix.Type)
	assert.True(t, viewIdPrefix.Optional)
	assert.True(t, viewIdPrefix.ForceNew)

	description := res.Schema["description"]
	assert.Equal(t, schema.TypeString, description.Type)
	assert.True(t, description.Optional)

	createTime := res.Schema["create_time"]
	assert.Equal(t, schema.TypeString, createTime.Type)
	assert.True(t, createTime.Computed)

	updateTime := res.Schema["update_time"]
	assert.Equal(t, schema.TypeString, updateTime.Type)
	assert.True(t, updateTime.Computed)
}
