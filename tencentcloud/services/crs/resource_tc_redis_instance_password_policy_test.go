package crs_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccrs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"
)

type mockMetaRedisInstancePasswordPolicy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaRedisInstancePasswordPolicy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaRedisInstancePasswordPolicy{}

func newMockMetaRedisInstancePasswordPolicy() *mockMetaRedisInstancePasswordPolicy {
	return &mockMetaRedisInstancePasswordPolicy{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/crs/ -run "TestRedisInstancePasswordPolicy" -v -count=1 -gcflags="all=-l"

// TestRedisInstancePasswordPolicy_Create tests the Create function
func TestRedisInstancePasswordPolicy_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisInstancePasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyInstancePasswordPolicy", func(request *redis.ModifyInstancePasswordPolicyRequest) (*redis.ModifyInstancePasswordPolicyResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "crs-test123", *request.InstanceId)
		assert.NotNil(t, request.PasswordPolicy)
		assert.NotNil(t, request.PasswordPolicy.Enabled)
		assert.Equal(t, true, *request.PasswordPolicy.Enabled)
		assert.NotNil(t, request.PasswordPolicy.MinLetterCount)
		assert.Equal(t, int64(2), *request.PasswordPolicy.MinLetterCount)
		assert.NotNil(t, request.PasswordPolicy.MinDigitCount)
		assert.Equal(t, int64(1), *request.PasswordPolicy.MinDigitCount)
		assert.NotNil(t, request.PasswordPolicy.MinSpecialCount)
		assert.Equal(t, int64(1), *request.PasswordPolicy.MinSpecialCount)
		assert.NotNil(t, request.PasswordPolicy.MinLength)
		assert.Equal(t, int64(8), *request.PasswordPolicy.MinLength)

		resp := redis.NewModifyInstancePasswordPolicyResponse()
		resp.Response = &redis.ModifyInstancePasswordPolicyResponseParams{
			RequestId: stringPtrRedisPolicy("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "crs-test123", *request.InstanceId)

		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis.PasswordPolicy{
				Enabled:         boolPtrRedisPolicy(true),
				MinLetterCount:  int64PtrRedisPolicy(2),
				MinDigitCount:   int64PtrRedisPolicy(1),
				MinSpecialCount: int64PtrRedisPolicy(1),
				MinLength:       int64PtrRedisPolicy(8),
			},
			RequestId: stringPtrRedisPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRedisInstancePasswordPolicy()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":       "crs-test123",
		"enabled":           true,
		"min_letter_count":  2,
		"min_digit_count":   1,
		"min_special_count": 1,
		"min_length":        8,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-test123", d.Id())
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, 2, d.Get("min_letter_count"))
	assert.Equal(t, 1, d.Get("min_digit_count"))
	assert.Equal(t, 1, d.Get("min_special_count"))
	assert.Equal(t, 8, d.Get("min_length"))
}

// TestRedisInstancePasswordPolicy_Read tests the Read function with a successful response
func TestRedisInstancePasswordPolicy_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisInstancePasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "crs-test123", *request.InstanceId)

		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis.PasswordPolicy{
				Enabled:         boolPtrRedisPolicy(true),
				MinLetterCount:  int64PtrRedisPolicy(2),
				MinDigitCount:   int64PtrRedisPolicy(1),
				MinSpecialCount: int64PtrRedisPolicy(1),
				MinLength:       int64PtrRedisPolicy(8),
			},
			RequestId: stringPtrRedisPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRedisInstancePasswordPolicy()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":       "crs-test123",
		"enabled":           true,
		"min_letter_count":  2,
		"min_digit_count":   1,
		"min_special_count": 1,
		"min_length":        8,
	})
	d.SetId("crs-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-test123", d.Id())
	assert.Equal(t, "crs-test123", d.Get("instance_id"))
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, 2, d.Get("min_letter_count"))
	assert.Equal(t, 1, d.Get("min_digit_count"))
	assert.Equal(t, 1, d.Get("min_special_count"))
	assert.Equal(t, 8, d.Get("min_length"))
}

// TestRedisInstancePasswordPolicy_ReadEmpty tests the Read function when the response is empty
func TestRedisInstancePasswordPolicy_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisInstancePasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{}
		return resp, nil
	})

	meta := newMockMetaRedisInstancePasswordPolicy()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-test123",
		"enabled":     true,
	})
	d.SetId("crs-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestRedisInstancePasswordPolicy_Update tests the Update function
func TestRedisInstancePasswordPolicy_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisInstancePasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyInstancePasswordPolicy", func(request *redis.ModifyInstancePasswordPolicyRequest) (*redis.ModifyInstancePasswordPolicyResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "crs-test123", *request.InstanceId)
		assert.NotNil(t, request.PasswordPolicy)
		assert.NotNil(t, request.PasswordPolicy.Enabled)
		assert.Equal(t, true, *request.PasswordPolicy.Enabled)
		assert.NotNil(t, request.PasswordPolicy.MinLetterCount)
		assert.Equal(t, int64(2), *request.PasswordPolicy.MinLetterCount)
		assert.NotNil(t, request.PasswordPolicy.MinDigitCount)
		assert.Equal(t, int64(1), *request.PasswordPolicy.MinDigitCount)
		assert.NotNil(t, request.PasswordPolicy.MinSpecialCount)
		assert.Equal(t, int64(1), *request.PasswordPolicy.MinSpecialCount)
		assert.NotNil(t, request.PasswordPolicy.MinLength)
		assert.Equal(t, int64(8), *request.PasswordPolicy.MinLength)

		resp := redis.NewModifyInstancePasswordPolicyResponse()
		resp.Response = &redis.ModifyInstancePasswordPolicyResponseParams{
			RequestId: stringPtrRedisPolicy("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis.PasswordPolicy{
				Enabled:         boolPtrRedisPolicy(true),
				MinLetterCount:  int64PtrRedisPolicy(2),
				MinDigitCount:   int64PtrRedisPolicy(1),
				MinSpecialCount: int64PtrRedisPolicy(1),
				MinLength:       int64PtrRedisPolicy(8),
			},
			RequestId: stringPtrRedisPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRedisInstancePasswordPolicy()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":       "crs-test123",
		"enabled":           true,
		"min_letter_count":  2,
		"min_digit_count":   1,
		"min_special_count": 1,
		"min_length":        8,
	})
	d.SetId("crs-test123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-test123", d.Id())
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, 2, d.Get("min_letter_count"))
	assert.Equal(t, 1, d.Get("min_digit_count"))
	assert.Equal(t, 1, d.Get("min_special_count"))
	assert.Equal(t, 8, d.Get("min_length"))
}

// TestRedisInstancePasswordPolicy_UpdateDisabled tests Update sends Enabled=false without optional fields
func TestRedisInstancePasswordPolicy_UpdateDisabled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisInstancePasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyInstancePasswordPolicy", func(request *redis.ModifyInstancePasswordPolicyRequest) (*redis.ModifyInstancePasswordPolicyResponse, error) {
		assert.NotNil(t, request.PasswordPolicy)
		assert.NotNil(t, request.PasswordPolicy.Enabled)
		assert.Equal(t, false, *request.PasswordPolicy.Enabled)
		// optional min_* fields should not be sent when unset
		assert.Nil(t, request.PasswordPolicy.MinLetterCount)
		assert.Nil(t, request.PasswordPolicy.MinDigitCount)
		assert.Nil(t, request.PasswordPolicy.MinSpecialCount)
		assert.Nil(t, request.PasswordPolicy.MinLength)

		resp := redis.NewModifyInstancePasswordPolicyResponse()
		resp.Response = &redis.ModifyInstancePasswordPolicyResponseParams{
			RequestId: stringPtrRedisPolicy("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis.PasswordPolicy{
				Enabled: boolPtrRedisPolicy(false),
			},
			RequestId: stringPtrRedisPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaRedisInstancePasswordPolicy()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-test123",
		"enabled":     false,
	})
	d.SetId("crs-test123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-test123", d.Id())
	assert.Equal(t, false, d.Get("enabled"))
}

// TestRedisInstancePasswordPolicy_Delete tests the Delete function
func TestRedisInstancePasswordPolicy_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisInstancePasswordPolicy().client, "UseRedisClient", redisClient)

	meta := newMockMetaRedisInstancePasswordPolicy()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-test123",
		"enabled":     true,
	})
	d.SetId("crs-test123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-test123", d.Id())
}

func stringPtrRedisPolicy(s string) *string {
	return &s
}

func boolPtrRedisPolicy(b bool) *bool {
	return &b
}

func int64PtrRedisPolicy(i int64) *int64 {
	return &i
}
