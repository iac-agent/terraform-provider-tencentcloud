package cls_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
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

// go test ./tencentcloud/services/cls/ -run "TestClsCkafkaConsumerCreate_WithNewParams" -v -count=1 -gcflags="all=-l"
func TestClsCkafkaConsumerCreate_WithNewParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateConsumer", func(request *clsv20201016.CreateConsumerRequest) (*clsv20201016.CreateConsumerResponse, error) {
		assert.Equal(t, "topic-test-123", *request.TopicId)
		assert.NotNil(t, request.RoleArn)
		assert.Equal(t, "qcs::cam::uin/123456789:roleName/MyRole", *request.RoleArn)
		assert.NotNil(t, request.ExternalId)
		assert.Equal(t, "my-external-id", *request.ExternalId)
		assert.NotNil(t, request.AdvancedConfig)
		assert.NotNil(t, request.AdvancedConfig.PartitionHashStatus)
		assert.True(t, *request.AdvancedConfig.PartitionHashStatus)
		assert.NotNil(t, request.AdvancedConfig.PartitionFields)
		assert.Equal(t, 2, len(request.AdvancedConfig.PartitionFields))
		// TypeSet is unordered, so check both values exist regardless of order
		fieldValues := make(map[string]bool)
		for _, f := range request.AdvancedConfig.PartitionFields {
			fieldValues[*f] = true
		}
		assert.True(t, fieldValues["__SOURCE__"])
		assert.True(t, fieldValues["__HOSTNAME__"])

		resp := clsv20201016.NewCreateConsumerResponse()
		resp.Response = &clsv20201016.CreateConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	// Mock ModifyConsumer for the update call after create
	patches.ApplyMethodFunc(clsClient, "ModifyConsumer", func(request *clsv20201016.ModifyConsumerRequest) (*clsv20201016.ModifyConsumerResponse, error) {
		resp := clsv20201016.NewModifyConsumerResponse()
		resp.Response = &clsv20201016.ModifyConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeConsumer for the read call
	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *clsv20201016.DescribeConsumerRequest) (*clsv20201016.DescribeConsumerResponse, error) {
		resp := clsv20201016.NewDescribeConsumerResponse()
		resp.Response = &clsv20201016.DescribeConsumerResponseParams{
			Effective: boolPtr(true),
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := cls.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":    "topic-test-123",
		"role_arn":    "qcs::cam::uin/123456789:roleName/MyRole",
		"external_id": "my-external-id",
		"advanced_config": []interface{}{
			map[string]interface{}{
				"partition_hash_status": true,
				"partition_fields":      []interface{}{"__SOURCE__", "__HOSTNAME__"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "topic-test-123", d.Id())
}

// go test ./tencentcloud/services/cls/ -run "TestClsCkafkaConsumerCreate_WithoutNewParams" -v -count=1 -gcflags="all=-l"
func TestClsCkafkaConsumerCreate_WithoutNewParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateConsumer", func(request *clsv20201016.CreateConsumerRequest) (*clsv20201016.CreateConsumerResponse, error) {
		assert.Nil(t, request.RoleArn)
		assert.Nil(t, request.ExternalId)
		assert.Nil(t, request.AdvancedConfig)

		resp := clsv20201016.NewCreateConsumerResponse()
		resp.Response = &clsv20201016.CreateConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "ModifyConsumer", func(request *clsv20201016.ModifyConsumerRequest) (*clsv20201016.ModifyConsumerResponse, error) {
		resp := clsv20201016.NewModifyConsumerResponse()
		resp.Response = &clsv20201016.ModifyConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *clsv20201016.DescribeConsumerRequest) (*clsv20201016.DescribeConsumerResponse, error) {
		resp := clsv20201016.NewDescribeConsumerResponse()
		resp.Response = &clsv20201016.DescribeConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := cls.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-456",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "topic-test-456", d.Id())
}

// go test ./tencentcloud/services/cls/ -run "TestClsCkafkaConsumerRead_Effective" -v -count=1 -gcflags="all=-l"
func TestClsCkafkaConsumerRead_Effective(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *clsv20201016.DescribeConsumerRequest) (*clsv20201016.DescribeConsumerResponse, error) {
		resp := clsv20201016.NewDescribeConsumerResponse()
		resp.Response = &clsv20201016.DescribeConsumerResponseParams{
			Effective: boolPtr(true),
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := cls.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-789",
	})
	d.SetId("topic-test-789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, true, d.Get("effective"))
}

// go test ./tencentcloud/services/cls/ -run "TestClsCkafkaConsumerRead_EffectiveFalse" -v -count=1 -gcflags="all=-l"
func TestClsCkafkaConsumerRead_EffectiveFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *clsv20201016.DescribeConsumerRequest) (*clsv20201016.DescribeConsumerResponse, error) {
		resp := clsv20201016.NewDescribeConsumerResponse()
		resp.Response = &clsv20201016.DescribeConsumerResponseParams{
			Effective: boolPtr(false),
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := cls.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-false",
	})
	d.SetId("topic-test-false")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, false, d.Get("effective"))
}

// go test ./tencentcloud/services/cls/ -run "TestClsCkafkaConsumerUpdate_WithNewParams" -v -count=1 -gcflags="all=-l"
func TestClsCkafkaConsumerUpdate_WithNewParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newCkafkaConsumerMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifyConsumer", func(request *clsv20201016.ModifyConsumerRequest) (*clsv20201016.ModifyConsumerResponse, error) {
		assert.Equal(t, "topic-update-123", *request.TopicId)
		assert.NotNil(t, request.Effective)
		assert.True(t, *request.Effective)
		assert.NotNil(t, request.RoleArn)
		assert.Equal(t, "qcs::cam::uin/123456789:roleName/UpdatedRole", *request.RoleArn)
		assert.NotNil(t, request.ExternalId)
		assert.Equal(t, "updated-external-id", *request.ExternalId)
		assert.NotNil(t, request.AdvancedConfig)
		assert.NotNil(t, request.AdvancedConfig.PartitionHashStatus)
		assert.False(t, *request.AdvancedConfig.PartitionHashStatus)

		resp := clsv20201016.NewModifyConsumerResponse()
		resp.Response = &clsv20201016.ModifyConsumerResponseParams{
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeConsumer", func(request *clsv20201016.DescribeConsumerRequest) (*clsv20201016.DescribeConsumerResponse, error) {
		resp := clsv20201016.NewDescribeConsumerResponse()
		resp.Response = &clsv20201016.DescribeConsumerResponseParams{
			Effective: boolPtr(true),
			RequestId: strPtr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCkafkaConsumerMockMeta()
	res := cls.ResourceTencentCloudClsCkafkaConsumer()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":    "topic-update-123",
		"effective":   true,
		"role_arn":    "qcs::cam::uin/123456789:roleName/UpdatedRole",
		"external_id": "updated-external-id",
		"advanced_config": []interface{}{
			map[string]interface{}{
				"partition_hash_status": false,
				"partition_fields":      []interface{}{},
			},
		},
	})
	d.SetId("topic-update-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// go test ./tencentcloud/services/cls/ -run "TestClsCkafkaConsumerSchema_NewParams" -v -count=1 -gcflags="all=-l"
func TestClsCkafkaConsumerSchema_NewParams(t *testing.T) {
	res := cls.ResourceTencentCloudClsCkafkaConsumer()

	assert.NotNil(t, res)

	// Check effective schema
	effectiveSchema, ok := res.Schema["effective"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeBool, effectiveSchema.Type)
	assert.True(t, effectiveSchema.Optional)
	assert.True(t, effectiveSchema.Computed)

	// Check role_arn schema
	roleArnSchema, ok := res.Schema["role_arn"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, roleArnSchema.Type)
	assert.True(t, roleArnSchema.Optional)

	// Check external_id schema
	externalIdSchema, ok := res.Schema["external_id"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, externalIdSchema.Type)
	assert.True(t, externalIdSchema.Optional)

	// Check advanced_config schema
	advancedConfigSchema, ok := res.Schema["advanced_config"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeList, advancedConfigSchema.Type)
	assert.True(t, advancedConfigSchema.Optional)
	assert.Equal(t, 1, advancedConfigSchema.MaxItems)

	// Check nested fields
	elem := advancedConfigSchema.Elem.(*schema.Resource)
	assert.Contains(t, elem.Schema, "partition_hash_status")
	assert.Equal(t, schema.TypeBool, elem.Schema["partition_hash_status"].Type)
	assert.True(t, elem.Schema["partition_hash_status"].Optional)
	assert.Contains(t, elem.Schema, "partition_fields")
	assert.Equal(t, schema.TypeSet, elem.Schema["partition_fields"].Type)
	assert.True(t, elem.Schema["partition_fields"].Optional)
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
