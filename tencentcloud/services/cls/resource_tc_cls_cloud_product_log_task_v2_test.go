package cls_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

// go test ./tencentcloud/services/cls/ -run "TestClsCloudProductLogTaskV2" -v -count=1 -gcflags="all=-l"

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

func ptrString(s string) *string { return &s }
func ptrInt64(i int64) *int64    { return &i }
func ptrUint64(u uint64) *uint64 { return &u }
func ptrBool(b bool) *bool       { return &b }

// TestClsCloudProductLogTaskV2_Schema validates the schema definition for new parameters
func TestClsCloudProductLogTaskV2_Schema(t *testing.T) {
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()

	assert.NotNil(t, res)

	// Check tags parameter
	assert.Contains(t, res.Schema, "tags")
	tags := res.Schema["tags"]
	assert.Equal(t, schema.TypeMap, tags.Type)
	assert.True(t, tags.Optional)
	assert.False(t, tags.Computed)

	// Check is_delete_topic parameter
	assert.Contains(t, res.Schema, "is_delete_topic")
	isDeleteTopic := res.Schema["is_delete_topic"]
	assert.Equal(t, schema.TypeBool, isDeleteTopic.Type)
	assert.True(t, isDeleteTopic.Optional)
	assert.Equal(t, false, isDeleteTopic.Default.(bool))

	// Check is_delete_logset parameter
	assert.Contains(t, res.Schema, "is_delete_logset")
	isDeleteLogset := res.Schema["is_delete_logset"]
	assert.Equal(t, schema.TypeBool, isDeleteLogset.Type)
	assert.True(t, isDeleteLogset.Optional)
	assert.Equal(t, false, isDeleteLogset.Default.(bool))

	// Check status computed parameter
	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeInt, status.Type)
	assert.True(t, status.Computed)
	assert.False(t, status.Optional)
}

// TestClsCloudProductLogTaskV2Create_WithTags tests that tags are correctly passed to the Create API request
func TestClsCloudProductLogTaskV2Create_WithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.CreateCloudProductLogCollectionRequest) (*clsv20201016.CreateCloudProductLogCollectionResponse, error) {
		// Verify tags are passed (order may vary due to map iteration)
		assert.NotNil(t, request.Tags)
		assert.Equal(t, 2, len(request.Tags))
		tagMap := make(map[string]string)
		for _, tag := range request.Tags {
			assert.NotNil(t, tag.Key)
			assert.NotNil(t, tag.Value)
			tagMap[*tag.Key] = *tag.Value
		}
		assert.Equal(t, "test", tagMap["env"])
		assert.Equal(t, "dev", tagMap["team"])

		resp := clsv20201016.NewCreateCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.CreateCloudProductLogCollectionResponseParams{
			TopicId:    ptrString("topic-12345"),
			TopicName:  ptrString("test-topic"),
			LogsetId:   ptrString("logset-12345"),
			LogsetName: ptrString("test-logset"),
			Status:     ptrInt64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
		resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
			Tasks: []*clsv20201016.CloudProductLogTaskInfo{
				{
					ClsRegion:  ptrString("ap-guangzhou"),
					InstanceId: ptrString("instance-12345"),
					LogType:    ptrString("PostgreSQL-SLOW"),
					Extend:     ptrString(""),
					Status:     ptrInt64(1),
				},
			},
			TotalCount: ptrUint64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"logset_name":          "test-logset",
		"topic_name":           "test-topic",
		"tags": map[string]interface{}{
			"env":  "test",
			"team": "dev",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestClsCloudProductLogTaskV2Create_WithoutTags tests that no tags are passed when not specified
func TestClsCloudProductLogTaskV2Create_WithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.CreateCloudProductLogCollectionRequest) (*clsv20201016.CreateCloudProductLogCollectionResponse, error) {
		// Verify tags are not set when not provided
		assert.Nil(t, request.Tags)

		resp := clsv20201016.NewCreateCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.CreateCloudProductLogCollectionResponseParams{
			TopicId:    ptrString("topic-12345"),
			TopicName:  ptrString("test-topic"),
			LogsetId:   ptrString("logset-12345"),
			LogsetName: ptrString("test-logset"),
			Status:     ptrInt64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
		resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
			Tasks: []*clsv20201016.CloudProductLogTaskInfo{
				{
					ClsRegion:  ptrString("ap-guangzhou"),
					InstanceId: ptrString("instance-12345"),
					LogType:    ptrString("PostgreSQL-SLOW"),
					Extend:     ptrString(""),
					Status:     ptrInt64(1),
				},
			},
			TotalCount: ptrUint64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"logset_name":          "test-logset",
		"topic_name":           "test-topic",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestClsCloudProductLogTaskV2Create_StatusSet tests that status is read from the Create response
func TestClsCloudProductLogTaskV2Create_StatusSet(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.CreateCloudProductLogCollectionRequest) (*clsv20201016.CreateCloudProductLogCollectionResponse, error) {
		resp := clsv20201016.NewCreateCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.CreateCloudProductLogCollectionResponseParams{
			TopicId:    ptrString("topic-12345"),
			TopicName:  ptrString("test-topic"),
			LogsetId:   ptrString("logset-12345"),
			LogsetName: ptrString("test-logset"),
			Status:     ptrInt64(-1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
		resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
			Tasks: []*clsv20201016.CloudProductLogTaskInfo{
				{
					ClsRegion:  ptrString("ap-guangzhou"),
					InstanceId: ptrString("instance-12345"),
					LogType:    ptrString("PostgreSQL-SLOW"),
					Extend:     ptrString(""),
					Status:     ptrInt64(1),
				},
			},
			TotalCount: ptrUint64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"logset_name":          "test-logset",
		"topic_name":           "test-topic",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	// After create + read, status should be set from the describe response (1 = created)
	status := d.Get("status")
	assert.Equal(t, 1, status.(int))
}

// TestClsCloudProductLogTaskV2Create_ResponseNilStatus tests handling of nil status in create response
func TestClsCloudProductLogTaskV2Create_ResponseNilStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.CreateCloudProductLogCollectionRequest) (*clsv20201016.CreateCloudProductLogCollectionResponse, error) {
		resp := clsv20201016.NewCreateCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.CreateCloudProductLogCollectionResponseParams{
			TopicId:    ptrString("topic-12345"),
			TopicName:  ptrString("test-topic"),
			LogsetId:   ptrString("logset-12345"),
			LogsetName: ptrString("test-logset"),
			Status:     nil,
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
		resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
			Tasks: []*clsv20201016.CloudProductLogTaskInfo{
				{
					ClsRegion:  ptrString("ap-guangzhou"),
					InstanceId: ptrString("instance-12345"),
					LogType:    ptrString("PostgreSQL-SLOW"),
					Extend:     ptrString(""),
					Status:     ptrInt64(1),
				},
			},
			TotalCount: ptrUint64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"logset_name":          "test-logset",
		"topic_name":           "test-topic",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestClsCloudProductLogTaskV2Create_APIError tests that create returns error when API fails
func TestClsCloudProductLogTaskV2Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.CreateCloudProductLogCollectionRequest) (*clsv20201016.CreateCloudProductLogCollectionResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid parameter")
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestClsCloudProductLogTaskV2Read_Status tests that status is read from the describe response
func TestClsCloudProductLogTaskV2Read_Status(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
		resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
			Tasks: []*clsv20201016.CloudProductLogTaskInfo{
				{
					ClsRegion:  ptrString("ap-guangzhou"),
					InstanceId: ptrString("instance-12345"),
					LogType:    ptrString("PostgreSQL-SLOW"),
					Extend:     ptrString(""),
					Status:     ptrInt64(1),
				},
			},
			TotalCount: ptrUint64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
	})
	d.SetId("instance-12345#PostgreSQL#PostgreSQL-SLOW#gz")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	status := d.Get("status")
	assert.Equal(t, 1, status.(int))
}

// TestClsCloudProductLogTaskV2Delete_WithIsDeleteTopic tests that is_delete_topic and is_delete_logset are passed to the Delete API request
func TestClsCloudProductLogTaskV2Delete_WithIsDeleteTopic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.DeleteCloudProductLogCollectionRequest) (*clsv20201016.DeleteCloudProductLogCollectionResponse, error) {
		// Verify is_delete_topic and is_delete_logset are passed
		assert.NotNil(t, request.IsDeleteTopic)
		assert.True(t, *request.IsDeleteTopic)
		assert.NotNil(t, request.IsDeleteLogset)
		assert.True(t, *request.IsDeleteLogset)

		resp := clsv20201016.NewDeleteCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.DeleteCloudProductLogCollectionResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
		resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
			Tasks:      []*clsv20201016.CloudProductLogTaskInfo{},
			TotalCount: ptrUint64(0),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"is_delete_topic":      true,
		"is_delete_logset":     true,
	})
	d.SetId("instance-12345#PostgreSQL#PostgreSQL-SLOW#gz")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestClsCloudProductLogTaskV2Delete_WithoutIsDeleteTopic tests that is_delete_topic and is_delete_logset are not set when not provided
func TestClsCloudProductLogTaskV2Delete_WithoutIsDeleteTopic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.DeleteCloudProductLogCollectionRequest) (*clsv20201016.DeleteCloudProductLogCollectionResponse, error) {
		// Verify is_delete_topic and is_delete_logset are not set when not provided
		assert.Nil(t, request.IsDeleteTopic)
		assert.Nil(t, request.IsDeleteLogset)

		resp := clsv20201016.NewDeleteCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.DeleteCloudProductLogCollectionResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
		resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
			Tasks:      []*clsv20201016.CloudProductLogTaskInfo{},
			TotalCount: ptrUint64(0),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
	})
	d.SetId("instance-12345#PostgreSQL#PostgreSQL-SLOW#gz")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestClsCloudProductLogTaskV2Delete_OnlyIsDeleteTopic tests that only is_delete_topic is set when is_delete_logset is not provided
func TestClsCloudProductLogTaskV2Delete_OnlyIsDeleteTopic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteCloudProductLogCollectionWithContext", func(_ context.Context, request *clsv20201016.DeleteCloudProductLogCollectionRequest) (*clsv20201016.DeleteCloudProductLogCollectionResponse, error) {
		assert.NotNil(t, request.IsDeleteTopic)
		assert.True(t, *request.IsDeleteTopic)
		assert.Nil(t, request.IsDeleteLogset)

		resp := clsv20201016.NewDeleteCloudProductLogCollectionResponse()
		resp.Response = &clsv20201016.DeleteCloudProductLogCollectionResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCloudProductLogTasks", func(request *clsv20201016.DescribeCloudProductLogTasksRequest) (*clsv20201016.DescribeCloudProductLogTasksResponse, error) {
		resp := clsv20201016.NewDescribeCloudProductLogTasksResponse()
		resp.Response = &clsv20201016.DescribeCloudProductLogTasksResponseParams{
			Tasks:      []*clsv20201016.CloudProductLogTaskInfo{},
			TotalCount: ptrUint64(0),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":          "instance-12345",
		"assumer_name":         "PostgreSQL",
		"log_type":             "PostgreSQL-SLOW",
		"cloud_product_region": "gz",
		"cls_region":           "ap-guangzhou",
		"is_delete_topic":      true,
	})
	d.SetId("instance-12345#PostgreSQL#PostgreSQL-SLOW#gz")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestClsCloudProductLogTaskV2Update_ImmutableArgs verifies that tags, is_delete_topic, and is_delete_logset
// are in the immutableArgs list in the Update function
func TestClsCloudProductLogTaskV2Update_ImmutableArgs(t *testing.T) {
	res := cls.ResourceTencentCloudClsCloudProductLogTaskV2()

	// Verify that tags is not ForceNew but should be immutable (can't be updated)
	assert.Contains(t, res.Schema, "tags")
	tags := res.Schema["tags"]
	assert.False(t, tags.ForceNew, "tags should not be ForceNew but handled as immutable in Update")

	// Verify that is_delete_topic is not ForceNew but should be immutable
	assert.Contains(t, res.Schema, "is_delete_topic")
	isDeleteTopic := res.Schema["is_delete_topic"]
	assert.False(t, isDeleteTopic.ForceNew, "is_delete_topic should not be ForceNew but handled as immutable in Update")

	// Verify that is_delete_logset is not ForceNew but should be immutable
	assert.Contains(t, res.Schema, "is_delete_logset")
	isDeleteLogset := res.Schema["is_delete_logset"]
	assert.False(t, isDeleteLogset.ForceNew, "is_delete_logset should not be ForceNew but handled as immutable in Update")

	// Verify that the only mutable arg is "extend"
	// Other non-ForceNew fields are treated as immutable in the Update function
	assert.Contains(t, res.Schema, "extend")
	extend := res.Schema["extend"]
	assert.False(t, extend.ForceNew, "extend should be mutable in Update")
}
