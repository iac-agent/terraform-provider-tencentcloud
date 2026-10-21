package crs

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func newBackupEncryptTestResource() *CrsRedisBackupEncryptConfigResource {
	return &CrsRedisBackupEncryptConfigResource{
		client: NewRedisService(&connectivity.TencentCloudClient{}),
	}
}

func backupEncryptStringPtr(v string) *string { return &v }
func backupEncryptInt64Ptr(v int64) *int64    { return &v }
func backupEncryptBoolPtr(v bool) *bool       { return &v }

// TestCrsRedisBackupEncryptConfig_Metadata verifies the registered type name.
func TestCrsRedisBackupEncryptConfig_Metadata(t *testing.T) {
	var resp resource.MetadataResponse
	newBackupEncryptTestResource().Metadata(context.Background(), resource.MetadataRequest{}, &resp)
	assert.Equal(t, "tencentcloud_redis_backup_encrypt", resp.TypeName)
}

// TestCrsRedisBackupEncryptConfig_Schema verifies the schema shape required by the change.
func TestCrsRedisBackupEncryptConfig_Schema(t *testing.T) {
	var resp resource.SchemaResponse
	newBackupEncryptTestResource().Schema(context.Background(), resource.SchemaRequest{}, &resp)
	assert.False(t, resp.Diagnostics.HasError())

	instanceId, ok := resp.Schema.Attributes["instance_id"]
	assert.True(t, ok)
	assert.True(t, instanceId.IsRequired())

	backupEncrypt, ok := resp.Schema.Attributes["backup_encrypt"]
	assert.True(t, ok)
	assert.True(t, backupEncrypt.IsRequired())

	encryptKeyType, ok := resp.Schema.Attributes["encrypt_key_type"]
	assert.True(t, ok)
	assert.True(t, encryptKeyType.IsOptional())
	assert.True(t, encryptKeyType.IsComputed())

	for _, name := range []string{"kms_key_id", "kms_region"} {
		attr, exists := resp.Schema.Attributes[name]
		assert.True(t, exists, name)
		assert.True(t, attr.IsOptional(), name)
		assert.True(t, attr.IsComputed(), name)
	}

	id, ok := resp.Schema.Attributes["id"]
	assert.True(t, ok)
	assert.True(t, id.IsComputed())
}

// TestCrsRedisBackupEncryptConfig_ModifyBackupEncrypt_RequestMapping verifies that every
// configured attribute is forwarded to ModifyBackupEncrypt.
func TestCrsRedisBackupEncryptConfig_ModifyBackupEncrypt_RequestMapping(t *testing.T) {
	r := newBackupEncryptTestResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var gotInstanceId string
	var gotBackupEncrypt bool
	var gotEncryptKeyType *int64
	var gotKmsKeyId *string
	var gotKmsRegion *string

	patches.ApplyMethodFunc(&r.client, "ModifyBackupEncrypt",
		func(_ context.Context, instanceId string, backupEncrypt bool, encryptKeyType *int64, kmsKeyId, kmsRegion *string) error {
			gotInstanceId = instanceId
			gotBackupEncrypt = backupEncrypt
			gotEncryptKeyType = encryptKeyType
			gotKmsKeyId = kmsKeyId
			gotKmsRegion = kmsRegion
			return nil
		},
	)

	m := CrsRedisBackupEncryptConfigResourceModel{
		Id:             types.StringUnknown(),
		InstanceId:     types.StringValue("crs-abc12345"),
		BackupEncrypt:  types.BoolValue(true),
		EncryptKeyType: types.Int64Value(1),
		KmsKeyId:       types.StringValue("kms-key-id"),
		KmsRegion:      types.StringValue("ap-guangzhou"),
	}

	err := r.modifyBackupEncrypt(context.Background(), &m)
	assert.NoError(t, err)

	assert.Equal(t, "crs-abc12345", gotInstanceId)
	assert.True(t, gotBackupEncrypt)
	assert.NotNil(t, gotEncryptKeyType)
	assert.Equal(t, int64(1), *gotEncryptKeyType)
	assert.NotNil(t, gotKmsKeyId)
	assert.Equal(t, "kms-key-id", *gotKmsKeyId)
	assert.NotNil(t, gotKmsRegion)
	assert.Equal(t, "ap-guangzhou", *gotKmsRegion)
}

// TestCrsRedisBackupEncryptConfig_ModifyBackupEncrypt_OmitUnset verifies that unknown
// optional attributes are not forwarded to the cloud API.
func TestCrsRedisBackupEncryptConfig_ModifyBackupEncrypt_OmitUnset(t *testing.T) {
	r := newBackupEncryptTestResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var gotEncryptKeyType *int64
	var gotKmsKeyId *string
	var gotKmsRegion *string

	patches.ApplyMethodFunc(&r.client, "ModifyBackupEncrypt",
		func(_ context.Context, _ string, _ bool, encryptKeyType *int64, kmsKeyId, kmsRegion *string) error {
			gotEncryptKeyType = encryptKeyType
			gotKmsKeyId = kmsKeyId
			gotKmsRegion = kmsRegion
			return nil
		},
	)

	m := CrsRedisBackupEncryptConfigResourceModel{
		Id:             types.StringUnknown(),
		InstanceId:     types.StringValue("crs-abc12345"),
		BackupEncrypt:  types.BoolValue(false),
		EncryptKeyType: types.Int64Unknown(),
		KmsKeyId:       types.StringUnknown(),
		KmsRegion:      types.StringUnknown(),
	}

	err := r.modifyBackupEncrypt(context.Background(), &m)
	assert.NoError(t, err)
	assert.Nil(t, gotEncryptKeyType)
	assert.Nil(t, gotKmsKeyId)
	assert.Nil(t, gotKmsRegion)
}

// TestCrsRedisBackupEncryptConfig_ModifyBackupEncrypt_Error verifies cloud errors are returned.
func TestCrsRedisBackupEncryptConfig_ModifyBackupEncrypt_Error(t *testing.T) {
	r := newBackupEncryptTestResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&r.client, "ModifyBackupEncrypt",
		func(_ context.Context, _ string, _ bool, _ *int64, _, _ *string) error {
			return sdkErrors.NewTencentCloudSDKError("InternalError.SystemError", "boom", "req-modify")
		},
	)

	m := CrsRedisBackupEncryptConfigResourceModel{
		Id:            types.StringUnknown(),
		InstanceId:    types.StringValue("crs-abc12345"),
		BackupEncrypt: types.BoolValue(true),
	}

	err := r.modifyBackupEncrypt(context.Background(), &m)
	assert.Error(t, err)
}

// TestCrsRedisBackupEncryptConfig_QueryBackupEncryptConfig verifies the state mapping helper
// used by Read/Create/Update fills every attribute from the cloud response.
func TestCrsRedisBackupEncryptConfig_QueryBackupEncryptConfig(t *testing.T) {
	r := newBackupEncryptTestResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var gotInstanceId string
	patches.ApplyMethodFunc(&r.client, "DescribeBackupEncryptConfig",
		func(_ context.Context, instanceId string) (*redis.DescribeBackupEncryptConfigResponseParams, error) {
			gotInstanceId = instanceId
			return &redis.DescribeBackupEncryptConfigResponseParams{
				BackupEncrypt:  backupEncryptBoolPtr(true),
				EncryptKeyType: backupEncryptInt64Ptr(1),
				KmsKeyId:       backupEncryptStringPtr("kms-key-id"),
				KmsRegion:      backupEncryptStringPtr("ap-guangzhou"),
			}, nil
		},
	)

	var diags diag.Diagnostics
	m := CrsRedisBackupEncryptConfigResourceModel{
		Id:            types.StringValue("crs-abc12345"),
		InstanceId:    types.StringNull(),
		BackupEncrypt: types.BoolValue(false),
	}

	config, err := r.queryBackupEncryptConfig(context.Background(), m.InstanceId.ValueString(), &diags)
	assert.NoError(t, err)
	assert.False(t, diags.HasError())
	assert.Equal(t, "crs-abc12345", gotInstanceId)

	applyBackupEncryptConfig(&m, "crs-abc12345", config)

	assert.Equal(t, "crs-abc12345", m.Id.ValueString())
	assert.Equal(t, "crs-abc12345", m.InstanceId.ValueString())
	assert.True(t, m.BackupEncrypt.ValueBool())
	assert.Equal(t, int64(1), m.EncryptKeyType.ValueInt64())
	assert.Equal(t, "kms-key-id", m.KmsKeyId.ValueString())
	assert.Equal(t, "ap-guangzhou", m.KmsRegion.ValueString())
}

// TestCrsRedisBackupEncryptConfig_ApplyBackupEncryptConfig_NilFields verifies that nil
// response pointers become null state values instead of zero values.
func TestCrsRedisBackupEncryptConfig_ApplyBackupEncryptConfig_NilFields(t *testing.T) {
	m := CrsRedisBackupEncryptConfigResourceModel{
		Id:            types.StringValue("crs-abc12345"),
		InstanceId:    types.StringValue("crs-abc12345"),
		BackupEncrypt: types.BoolValue(true),
	}

	applyBackupEncryptConfig(&m, "crs-abc12345", &redis.DescribeBackupEncryptConfigResponseParams{
		EncryptKeyType: backupEncryptInt64Ptr(0),
	})

	assert.True(t, m.EncryptKeyType.Equal(types.Int64Value(0)))
	assert.True(t, m.KmsKeyId.IsNull())
	assert.True(t, m.KmsRegion.IsNull())
}

// TestCrsRedisBackupEncryptConfig_RetrieveError verifies the not-found detection helper.
func TestCrsRedisBackupEncryptConfig_RetrieveError(t *testing.T) {
	notFound := sdkErrors.NewTencentCloudSDKError(RedisInstanceNotFound, "instance not exists", "req-1")
	assert.True(t, isRedisInstanceNotFound(notFound))

	other := sdkErrors.NewTencentCloudSDKError("InternalError.SystemError", "boom", "req-2")
	assert.False(t, isRedisInstanceNotFound(other))

	assert.False(t, isRedisInstanceNotFound(nil))
}

// TestCrsRedisBackupEncryptConfig_ImportState_EmptyId verifies an empty import id is rejected.
func TestCrsRedisBackupEncryptConfig_ImportState_EmptyId(t *testing.T) {
	r := newBackupEncryptTestResource()

	resp := resource.ImportStateResponse{}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: ""}, &resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestCrsRedisBackupEncryptConfig_Int64OneOf verifies the encrypt key type validator description.
func TestCrsRedisBackupEncryptConfig_Int64OneOf(t *testing.T) {
	v := int64OneOf{allowed: []int64{0, 1}}
	assert.Contains(t, v.Description(context.Background()), "value must be one of")
}
