package cls_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	clsService "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

// ckafkaConsumerMockMeta implements tccommon.ProviderMeta
type ckafkaConsumerMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *ckafkaConsumerMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &ckafkaConsumerMockMeta{}

func newCkafkaConsumerMockMeta() *ckafkaConsumerMockMeta {
	return &ckafkaConsumerMockMeta{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/cls/ -run "TestClsCkafkaConsumer" -v -count=1 -gcflags="all=-l"

// TestClsCkafkaConsumer_Schema validates new parameters in schema
func TestClsCkafkaConsumer_Schema(t *testing.T) {
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)

	// Check effective parameter
	assert.Contains(t, res.Schema, "effective")
	effectiveSchema := res.Schema["effective"]
	assert.Equal(t, schema.TypeBool, effectiveSchema.Type)
	assert.True(t, effectiveSchema.Optional)

	// Check role_arn parameter
	assert.Contains(t, res.Schema, "role_arn")
	roleArnSchema := res.Schema["role_arn"]
	assert.Equal(t, schema.TypeString, roleArnSchema.Type)
	assert.True(t, roleArnSchema.Optional)

	// Check external_id parameter
	assert.Contains(t, res.Schema, "external_id")
	externalIdSchema := res.Schema["external_id"]
	assert.Equal(t, schema.TypeString, externalIdSchema.Type)
	assert.True(t, externalIdSchema.Optional)

	// Check advanced_config parameter
	assert.Contains(t, res.Schema, "advanced_config")
	advancedConfigSchema := res.Schema["advanced_config"]
	assert.Equal(t, schema.TypeList, advancedConfigSchema.Type)
	assert.True(t, advancedConfigSchema.Optional)
	assert.Equal(t, 1, advancedConfigSchema.MaxItems)
}

// TestClsCkafkaConsumer_Create_WithRoleArn tests role_arn in create flow
func TestClsCkafkaConsumer_Create_WithRoleArn(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateConsumer", func(request *cls.CreateConsumerRequest) (*cls.CreateConsumerResponse, error) {
		assert.NotNil(t, request.RoleArn)
		assert.Equal(t, "qcs::cam::uin/100000000001:roleName/MyRole", *request.RoleArn)
		resp := cls.NewCreateConsumerResponse()
		resp.Response = &cls.CreateConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			Effective: boolPtr(true),
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-123456",
		"role_arn": "qcs::cam::uin/100000000001:roleName/MyRole",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "topic-123456", d.Id())
}

// TestClsCkafkaConsumer_Create_WithExternalId tests external_id in create flow
func TestClsCkafkaConsumer_Create_WithExternalId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateConsumer", func(request *cls.CreateConsumerRequest) (*cls.CreateConsumerResponse, error) {
		assert.NotNil(t, request.ExternalId)
		assert.Equal(t, "myExternalId", *request.ExternalId)
		resp := cls.NewCreateConsumerResponse()
		resp.Response = &cls.CreateConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":    "topic-123456",
		"external_id": "myExternalId",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "topic-123456", d.Id())
}

// TestClsCkafkaConsumer_Create_WithAdvancedConfig tests advanced_config in create flow
func TestClsCkafkaConsumer_Create_WithAdvancedConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateConsumer", func(request *cls.CreateConsumerRequest) (*cls.CreateConsumerResponse, error) {
		assert.NotNil(t, request.AdvancedConfig)
		assert.NotNil(t, request.AdvancedConfig.PartitionHashStatus)
		assert.True(t, *request.AdvancedConfig.PartitionHashStatus)
		assert.NotNil(t, request.AdvancedConfig.PartitionFields)
		assert.Equal(t, 2, len(request.AdvancedConfig.PartitionFields))
		resp := cls.NewCreateConsumerResponse()
		resp.Response = &cls.CreateConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-123456",
		"advanced_config": []interface{}{
			map[string]interface{}{
				"partition_hash_status": true,
				"partition_fields":      []interface{}{"__SOURCE__", "__FILENAME__"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "topic-123456", d.Id())
}

// TestClsCkafkaConsumer_Read_Effective tests effective field in read flow
func TestClsCkafkaConsumer_Read_Effective(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		assert.Equal(t, "topic-123456", *request.TopicId)
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			Effective: boolPtr(true),
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-123456",
	})
	d.SetId("topic-123456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, true, d.Get("effective"))
}

// TestClsCkafkaConsumer_Read_EffectiveFalse tests effective=false in read flow
func TestClsCkafkaConsumer_Read_EffectiveFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			Effective: boolPtr(false),
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-123456",
	})
	d.SetId("topic-123456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, false, d.Get("effective"))
}

// TestClsCkafkaConsumer_Update_Effective tests effective parameter in update flow
func TestClsCkafkaConsumer_Update_Effective(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifyConsumer", func(request *cls.ModifyConsumerRequest) (*cls.ModifyConsumerResponse, error) {
		assert.NotNil(t, request.Effective)
		assert.True(t, *request.Effective)
		resp := cls.NewModifyConsumerResponse()
		resp.Response = &cls.ModifyConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			Effective: boolPtr(true),
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":  "topic-123456",
		"effective": true,
	})
	d.SetId("topic-123456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestClsCkafkaConsumer_Update_RoleArn tests role_arn in update flow
func TestClsCkafkaConsumer_Update_RoleArn(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifyConsumer", func(request *cls.ModifyConsumerRequest) (*cls.ModifyConsumerResponse, error) {
		assert.NotNil(t, request.RoleArn)
		assert.Equal(t, "qcs::cam::uin/100000000001:roleName/MyRole", *request.RoleArn)
		resp := cls.NewModifyConsumerResponse()
		resp.Response = &cls.ModifyConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-123456",
		"role_arn": "qcs::cam::uin/100000000001:roleName/MyRole",
	})
	d.SetId("topic-123456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestClsCkafkaConsumer_Update_ExternalId tests external_id in update flow
func TestClsCkafkaConsumer_Update_ExternalId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifyConsumer", func(request *cls.ModifyConsumerRequest) (*cls.ModifyConsumerResponse, error) {
		assert.NotNil(t, request.ExternalId)
		assert.Equal(t, "myExternalId", *request.ExternalId)
		resp := cls.NewModifyConsumerResponse()
		resp.Response = &cls.ModifyConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":    "topic-123456",
		"external_id": "myExternalId",
	})
	d.SetId("topic-123456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestClsCkafkaConsumer_Update_AdvancedConfig tests advanced_config in update flow
func TestClsCkafkaConsumer_Update_AdvancedConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifyConsumer", func(request *cls.ModifyConsumerRequest) (*cls.ModifyConsumerResponse, error) {
		assert.NotNil(t, request.AdvancedConfig)
		assert.NotNil(t, request.AdvancedConfig.PartitionHashStatus)
		assert.True(t, *request.AdvancedConfig.PartitionHashStatus)
		assert.NotNil(t, request.AdvancedConfig.PartitionFields)
		assert.Equal(t, 2, len(request.AdvancedConfig.PartitionFields))
		resp := cls.NewModifyConsumerResponse()
		resp.Response = &cls.ModifyConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		resp := cls.NewDescribeConsumerResponse()
		resp.Response = &cls.DescribeConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-123456",
		"advanced_config": []interface{}{
			map[string]interface{}{
				"partition_hash_status": true,
				"partition_fields":      []interface{}{"__SOURCE__", "__FILENAME__"},
			},
		},
	})
	d.SetId("topic-123456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestClsCkafkaConsumer_Create_APIError tests Create handles API error
func TestClsCkafkaConsumer_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateConsumer", func(request *cls.CreateConsumerRequest) (*cls.CreateConsumerResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid topic_id")
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-invalid",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestClsCkafkaConsumer_Read_NotFound tests Read handles resource not found
func TestClsCkafkaConsumer_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *cls.DescribeConsumerRequest) (*cls.DescribeConsumerResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Consumer not found")
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-not-found",
	})
	d.SetId("topic-not-found")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestClsCkafkaConsumer_Delete_Success tests Delete flow
func TestClsCkafkaConsumer_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteConsumer", func(request *cls.DeleteConsumerRequest) (*cls.DeleteConsumerResponse, error) {
		assert.Equal(t, "topic-123456", *request.TopicId)
		resp := cls.NewDeleteConsumerResponse()
		resp.Response = &cls.DeleteConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := clsService.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-123456",
	})
	d.SetId("topic-123456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
