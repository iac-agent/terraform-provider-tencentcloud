package crs_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"
)

type mockMetaRedisBackupEncrypt struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaRedisBackupEncrypt) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaRedisBackupEncrypt{}

func newMockMetaRedisBackupEncrypt() *mockMetaRedisBackupEncrypt {
	return &mockMetaRedisBackupEncrypt{client: &connectivity.TencentCloudClient{}}
}

func ptrStringRBE(s string) *string {
	return &s
}

func TestRedisBackupEncryptConfig_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(_ context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "crs-xxxxxxxx", *request.InstanceId)

		backupEncrypt := true
		encryptKeyType := int64(1)
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = &redis.DescribeBackupEncryptConfigResponseParams{
			RequestId:      ptrStringRBE("fake-request-id"),
			BackupEncrypt:  &backupEncrypt,
			EncryptKeyType: &encryptKeyType,
			KmsKeyId:       ptrStringRBE("23e80852-1e26-46e5-9940-251abb4028e8"),
			KmsRegion:      ptrStringRBE("ap-guangzhou"),
		}
		return resp, nil
	})

	meta := newMockMetaRedisBackupEncrypt()
	res := crs.ResourceTencentCloudRedisBackupEncryptConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":      "crs-xxxxxxxx",
		"backup_encrypt":   false,
		"encrypt_key_type": 0,
	})
	d.SetId("crs-xxxxxxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-xxxxxxxx", d.Id())
	assert.Equal(t, true, d.Get("backup_encrypt"))
	assert.Equal(t, 1, d.Get("encrypt_key_type"))
	assert.Equal(t, "23e80852-1e26-46e5-9940-251abb4028e8", d.Get("kms_key_id"))
	assert.Equal(t, "ap-guangzhou", d.Get("kms_region"))
}

func TestRedisBackupEncryptConfig_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(_ context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{
			Code:    "ResourceNotFound.InstanceNotExists",
			Message: "instance not found",
		}
	})

	meta := newMockMetaRedisBackupEncrypt()
	res := crs.ResourceTencentCloudRedisBackupEncryptConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "crs-notexist",
		"backup_encrypt": true,
	})
	d.SetId("crs-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestRedisBackupEncryptConfig_Read_PartialFieldsNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(_ context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		backupEncrypt := true
		encryptKeyType := int64(0)
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = &redis.DescribeBackupEncryptConfigResponseParams{
			RequestId:      ptrStringRBE("fake-request-id"),
			BackupEncrypt:  &backupEncrypt,
			EncryptKeyType: &encryptKeyType,
		}
		return resp, nil
	})

	meta := newMockMetaRedisBackupEncrypt()
	res := crs.ResourceTencentCloudRedisBackupEncryptConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "crs-xxxxxxxx",
		"backup_encrypt": true,
	})
	d.SetId("crs-xxxxxxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-xxxxxxxx", d.Id())
	assert.Equal(t, true, d.Get("backup_encrypt"))
	assert.Equal(t, 0, d.Get("encrypt_key_type"))
	assert.Equal(t, "", d.Get("kms_key_id"))
	assert.Equal(t, "", d.Get("kms_region"))
}

func TestRedisBackupEncryptConfig_Update_EnableWithTencentCloudManagedKey(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyBackupEncryptWithContext", func(_ context.Context, request *redis.ModifyBackupEncryptRequest) (*redis.ModifyBackupEncryptResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "crs-xxxxxxxx", *request.InstanceId)
		assert.NotNil(t, request.BackupEncrypt)
		assert.Equal(t, true, *request.BackupEncrypt)

		resp := redis.NewModifyBackupEncryptResponse()
		resp.Response = &redis.ModifyBackupEncryptResponseParams{
			RequestId: ptrStringRBE("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(_ context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		backupEncrypt := true
		encryptKeyType := int64(0)
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = &redis.DescribeBackupEncryptConfigResponseParams{
			RequestId:      ptrStringRBE("fake-request-id"),
			BackupEncrypt:  &backupEncrypt,
			EncryptKeyType: &encryptKeyType,
		}
		return resp, nil
	})

	meta := newMockMetaRedisBackupEncrypt()
	res := crs.ResourceTencentCloudRedisBackupEncryptConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "crs-xxxxxxxx",
		"backup_encrypt": true,
	})
	d.SetId("crs-xxxxxxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-xxxxxxxx", d.Id())
	assert.Equal(t, true, d.Get("backup_encrypt"))
}

func TestRedisBackupEncryptConfig_Update_WithCustomKmsKey(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyBackupEncryptWithContext", func(_ context.Context, request *redis.ModifyBackupEncryptRequest) (*redis.ModifyBackupEncryptResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "crs-xxxxxxxx", *request.InstanceId)
		assert.NotNil(t, request.BackupEncrypt)
		assert.Equal(t, true, *request.BackupEncrypt)
		assert.NotNil(t, request.EncryptKeyType)
		assert.Equal(t, int64(1), *request.EncryptKeyType)
		assert.NotNil(t, request.KmsKeyId)
		assert.Equal(t, "23e80852-1e26-46e5-9940-251abb4028e8", *request.KmsKeyId)
		assert.NotNil(t, request.KmsRegion)
		assert.Equal(t, "ap-guangzhou", *request.KmsRegion)

		resp := redis.NewModifyBackupEncryptResponse()
		resp.Response = &redis.ModifyBackupEncryptResponseParams{
			RequestId: ptrStringRBE("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(_ context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		backupEncrypt := true
		encryptKeyType := int64(1)
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = &redis.DescribeBackupEncryptConfigResponseParams{
			RequestId:      ptrStringRBE("fake-request-id"),
			BackupEncrypt:  &backupEncrypt,
			EncryptKeyType: &encryptKeyType,
			KmsKeyId:       ptrStringRBE("23e80852-1e26-46e5-9940-251abb4028e8"),
			KmsRegion:      ptrStringRBE("ap-guangzhou"),
		}
		return resp, nil
	})

	meta := newMockMetaRedisBackupEncrypt()
	res := crs.ResourceTencentCloudRedisBackupEncryptConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":      "crs-xxxxxxxx",
		"backup_encrypt":   true,
		"encrypt_key_type": 1,
		"kms_key_id":       "23e80852-1e26-46e5-9940-251abb4028e8",
		"kms_region":       "ap-guangzhou",
	})
	d.SetId("crs-xxxxxxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-xxxxxxxx", d.Id())
	assert.Equal(t, true, d.Get("backup_encrypt"))
	assert.Equal(t, 1, d.Get("encrypt_key_type"))
	assert.Equal(t, "23e80852-1e26-46e5-9940-251abb4028e8", d.Get("kms_key_id"))
	assert.Equal(t, "ap-guangzhou", d.Get("kms_region"))
}

func TestRedisBackupEncryptConfig_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyBackupEncryptWithContext", func(_ context.Context, request *redis.ModifyBackupEncryptRequest) (*redis.ModifyBackupEncryptResponse, error) {
		assert.NotNil(t, request.InstanceId)
		assert.Equal(t, "crs-xxxxxxxx", *request.InstanceId)
		assert.NotNil(t, request.BackupEncrypt)
		assert.Equal(t, true, *request.BackupEncrypt)

		resp := redis.NewModifyBackupEncryptResponse()
		resp.Response = &redis.ModifyBackupEncryptResponseParams{
			RequestId: ptrStringRBE("fake-request-id-create"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(_ context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		backupEncrypt := true
		encryptKeyType := int64(0)
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = &redis.DescribeBackupEncryptConfigResponseParams{
			RequestId:      ptrStringRBE("fake-request-id"),
			BackupEncrypt:  &backupEncrypt,
			EncryptKeyType: &encryptKeyType,
		}
		return resp, nil
	})

	meta := newMockMetaRedisBackupEncrypt()
	res := crs.ResourceTencentCloudRedisBackupEncryptConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "crs-xxxxxxxx",
		"backup_encrypt": true,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-xxxxxxxx", d.Id())
	assert.Equal(t, true, d.Get("backup_encrypt"))
}

func TestRedisBackupEncryptConfig_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	meta := newMockMetaRedisBackupEncrypt()
	res := crs.ResourceTencentCloudRedisBackupEncryptConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "crs-xxxxxxxx",
		"backup_encrypt": true,
	})
	d.SetId("crs-xxxxxxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-xxxxxxxx", d.Id())
}
