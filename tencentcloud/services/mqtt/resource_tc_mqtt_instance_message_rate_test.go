package mqtt_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mqtt"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

// mockMqttMeta implements tccommon.ProviderMeta
type mockMqttMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMqttMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMqttMeta{}

func newMockMqttMeta() *mockMqttMeta {
	return &mockMqttMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string { return &s }
func ptrInt64(i int64) *int64    { return &i }
func ptrBool(b bool) *bool       { return &b }

// go test ./tencentcloud/services/mqtt/ -run "TestMqttInstanceMessageRate" -v -count=1 -gcflags="all=-l"

// TestMqttInstanceMessageRate_Schema validates message_rate field in schema
func TestMqttInstanceMessageRate_Schema(t *testing.T) {
	res := mqtt.ResourceTencentCloudMqttInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "message_rate")

	messageRate := res.Schema["message_rate"]
	assert.Equal(t, schema.TypeInt, messageRate.Type)
	assert.True(t, messageRate.Optional)
	assert.True(t, messageRate.Computed)
	assert.False(t, messageRate.Required)
	assert.False(t, messageRate.ForceNew)
}

// TestMqttInstanceMessageRate_Read_Success tests Read retrieves MessageRate from API response
func TestMqttInstanceMessageRate_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	patches.ApplyMethodReturn(newMockMqttMeta().client, "UseMqttV20240516Client", mqttClient)

	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:        ptrString("PRO"),
			InstanceId:          ptrString("mqtt-test123"),
			InstanceName:        ptrString("test-instance"),
			Remark:              ptrString("test remark"),
			SkuCode:             ptrString("pro_6k_1"),
			PayMode:             ptrString("POSTPAID"),
			InstanceStatus:      ptrString("RUNNING"),
			RenewFlag:           ptrInt64(0),
			AutomaticActivation: ptrBool(false),
			AuthorizationPolicy: ptrBool(false),
			MessageRate:         ptrInt64(100),
		}
		return resp, nil
	})

	// Mock tag service
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMqttMeta()
	res := mqtt.ResourceTencentCloudMqttInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "PRO",
		"name":          "test-instance",
		"sku_code":      "pro_6k_1",
	})
	d.SetId("mqtt-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 100, d.Get("message_rate"))
}

// TestMqttInstanceMessageRate_Read_NilMessageRate tests Read when MessageRate is nil in response
func TestMqttInstanceMessageRate_Read_NilMessageRate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	patches.ApplyMethodReturn(newMockMqttMeta().client, "UseMqttV20240516Client", mqttClient)

	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:   ptrString("PRO"),
			InstanceId:     ptrString("mqtt-test123"),
			InstanceName:   ptrString("test-instance"),
			SkuCode:        ptrString("pro_6k_1"),
			PayMode:        ptrString("POSTPAID"),
			InstanceStatus: ptrString("RUNNING"),
			// MessageRate is nil
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMqttMeta()
	res := mqtt.ResourceTencentCloudMqttInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "PRO",
		"name":          "test-instance",
		"sku_code":      "pro_6k_1",
	})
	d.SetId("mqtt-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, d.Get("message_rate"))
}

// TestMqttInstanceMessageRate_Create_WithMessageRate tests Create with message_rate set triggers ModifyInstance
func TestMqttInstanceMessageRate_Create_WithMessageRate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	patches.ApplyMethodReturn(newMockMqttMeta().client, "UseMqttV20240516Client", mqttClient)

	// Mock CreateInstanceWithContext
	patches.ApplyMethodFunc(mqttClient, "CreateInstanceWithContext", func(ctx context.Context, request *mqttv20240516.CreateInstanceRequest) (*mqttv20240516.CreateInstanceResponse, error) {
		resp := mqttv20240516.NewCreateInstanceResponse()
		resp.Response = &mqttv20240516.CreateInstanceResponseParams{
			InstanceId: ptrString("mqtt-test123"),
		}
		return resp, nil
	})

	// Mock DescribeInstanceWithContext for wait loop (return RUNNING immediately)
	patches.ApplyMethodFunc(mqttClient, "DescribeInstanceWithContext", func(ctx context.Context, request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceStatus: ptrString("RUNNING"),
		}
		return resp, nil
	})

	// Mock ModifyInstanceWithContext (called for message_rate after creation)
	modifyCalled := false
	patches.ApplyMethodFunc(mqttClient, "ModifyInstanceWithContext", func(ctx context.Context, request *mqttv20240516.ModifyInstanceRequest) (*mqttv20240516.ModifyInstanceResponse, error) {
		modifyCalled = true
		assert.NotNil(t, request.MessageRate)
		assert.Equal(t, int64(50), *request.MessageRate)
		resp := mqttv20240516.NewModifyInstanceResponse()
		resp.Response = &mqttv20240516.ModifyInstanceResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeInstance for Read call at end of Create
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceStatus:      ptrString("RUNNING"),
			InstanceType:        ptrString("BASIC"),
			InstanceName:        ptrString("test-instance"),
			SkuCode:             ptrString("basic_2k"),
			PayMode:             ptrString("POSTPAID"),
			AutomaticActivation: ptrBool(false),
			AuthorizationPolicy: ptrBool(false),
			MessageRate:         ptrInt64(50),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMqttMeta()
	res := mqtt.ResourceTencentCloudMqttInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "BASIC",
		"name":          "test-instance",
		"sku_code":      "basic_2k",
		"pay_mode":      0,
		"message_rate":  50,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled, "ModifyInstanceWithContext should be called when message_rate is set")
	assert.Equal(t, "mqtt-test123", d.Id())
}

// TestMqttInstanceMessageRate_Create_WithoutMessageRate tests Create without message_rate does not trigger ModifyInstance
func TestMqttInstanceMessageRate_Create_WithoutMessageRate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	patches.ApplyMethodReturn(newMockMqttMeta().client, "UseMqttV20240516Client", mqttClient)

	patches.ApplyMethodFunc(mqttClient, "CreateInstanceWithContext", func(ctx context.Context, request *mqttv20240516.CreateInstanceRequest) (*mqttv20240516.CreateInstanceResponse, error) {
		resp := mqttv20240516.NewCreateInstanceResponse()
		resp.Response = &mqttv20240516.CreateInstanceResponseParams{
			InstanceId: ptrString("mqtt-test456"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mqttClient, "DescribeInstanceWithContext", func(ctx context.Context, request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceStatus: ptrString("RUNNING"),
		}
		return resp, nil
	})

	// ModifyInstanceWithContext should NOT be called when message_rate, automatic_activation, and authorization_policy are all unset
	modifyCalled := false
	patches.ApplyMethodFunc(mqttClient, "ModifyInstanceWithContext", func(ctx context.Context, request *mqttv20240516.ModifyInstanceRequest) (*mqttv20240516.ModifyInstanceResponse, error) {
		modifyCalled = true
		resp := mqttv20240516.NewModifyInstanceResponse()
		return resp, nil
	})

	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceStatus:      ptrString("RUNNING"),
			InstanceType:        ptrString("BASIC"),
			InstanceName:        ptrString("test-instance"),
			SkuCode:             ptrString("basic_2k"),
			PayMode:             ptrString("POSTPAID"),
			AutomaticActivation: ptrBool(false),
			AuthorizationPolicy: ptrBool(false),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMqttMeta()
	res := mqtt.ResourceTencentCloudMqttInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "BASIC",
		"name":          "test-instance",
		"sku_code":      "basic_2k",
		"pay_mode":      0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled, "ModifyInstanceWithContext should NOT be called when message_rate is not set")
}

// TestMqttInstanceMessageRate_Update_ChangeMessageRate tests Update with message_rate change
func TestMqttInstanceMessageRate_Update_ChangeMessageRate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	patches.ApplyMethodReturn(newMockMqttMeta().client, "UseMqttV20240516Client", mqttClient)

	// Patch HasChange to simulate only message_rate changing
	immutableArgs := map[string]bool{
		"instance_type": true,
		"vpc_list":      true,
		"renew_flag":    true,
		"time_span":     true,
		"pay_mode":      true,
	}
	patches.ApplyMethodFunc(new(schema.ResourceData), "HasChange", func(key string) bool {
		if immutableArgs[key] {
			return false
		}
		if key == "message_rate" {
			return true
		}
		return false
	})

	// Mock ModifyInstanceWithContext
	modifyCalled := false
	patches.ApplyMethodFunc(mqttClient, "ModifyInstanceWithContext", func(ctx context.Context, request *mqttv20240516.ModifyInstanceRequest) (*mqttv20240516.ModifyInstanceResponse, error) {
		modifyCalled = true
		assert.NotNil(t, request.MessageRate)
		assert.Equal(t, int64(200), *request.MessageRate)
		resp := mqttv20240516.NewModifyInstanceResponse()
		resp.Response = &mqttv20240516.ModifyInstanceResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeInstance for Read call at end of Update
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceStatus:      ptrString("RUNNING"),
			InstanceType:        ptrString("BASIC"),
			InstanceName:        ptrString("test-instance"),
			SkuCode:             ptrString("basic_2k"),
			PayMode:             ptrString("POSTPAID"),
			AutomaticActivation: ptrBool(false),
			AuthorizationPolicy: ptrBool(false),
			MessageRate:         ptrInt64(200),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMqttMeta()
	res := mqtt.ResourceTencentCloudMqttInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "BASIC",
		"name":          "test-instance",
		"sku_code":      "basic_2k",
		"pay_mode":      0,
		"message_rate":  200,
	})
	d.SetId("mqtt-test123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled, "ModifyInstanceWithContext should be called when message_rate changes")
}

// TestMqttInstanceMessageRate_Update_ImmutableArgs tests that message_rate is not in immutableArgs
func TestMqttInstanceMessageRate_Update_ImmutableArgs(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	patches.ApplyMethodReturn(newMockMqttMeta().client, "UseMqttV20240516Client", mqttClient)

	// Patch HasChange to simulate only message_rate changing
	immutableArgs := map[string]bool{
		"instance_type": true,
		"vpc_list":      true,
		"renew_flag":    true,
		"time_span":     true,
		"pay_mode":      true,
	}
	patches.ApplyMethodFunc(new(schema.ResourceData), "HasChange", func(key string) bool {
		if immutableArgs[key] {
			return false
		}
		if key == "message_rate" {
			return true
		}
		return false
	})

	// Mock ModifyInstanceWithContext - should be called since message_rate is mutable
	patches.ApplyMethodFunc(mqttClient, "ModifyInstanceWithContext", func(ctx context.Context, request *mqttv20240516.ModifyInstanceRequest) (*mqttv20240516.ModifyInstanceResponse, error) {
		resp := mqttv20240516.NewModifyInstanceResponse()
		resp.Response = &mqttv20240516.ModifyInstanceResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceStatus:      ptrString("RUNNING"),
			InstanceType:        ptrString("BASIC"),
			InstanceName:        ptrString("test-instance"),
			SkuCode:             ptrString("basic_2k"),
			PayMode:             ptrString("POSTPAID"),
			AutomaticActivation: ptrBool(false),
			AuthorizationPolicy: ptrBool(false),
			MessageRate:         ptrInt64(300),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMqttMeta()
	res := mqtt.ResourceTencentCloudMqttInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "BASIC",
		"name":          "test-instance",
		"sku_code":      "basic_2k",
		"pay_mode":      0,
		"message_rate":  300,
	})
	d.SetId("mqtt-test123")

	// Changing message_rate should not return an immutable error
	err := res.Update(d, meta)
	assert.NoError(t, err)
}
