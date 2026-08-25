package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestInferenceAPIToken" -v -count=1 -gcflags="all=-l"

// TestInferenceAPIToken_Create_Success tests Create calls API and sets composite ID
func TestInferenceAPIToken_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "test-token", *request.Name)
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrString("token-abc123"),
			Content:   ptrString("my-secret-token-content"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ interface{}, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64_IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-abc123"),
					Name:       ptrString("test-token"),
					Content:    ptrString("my-secret-token-content"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#token-abc123", d.Id())
	assert.Equal(t, "test-token", d.Get("name"))
	assert.Equal(t, "my-secret-token-content", d.Get("content"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
}

// TestInferenceAPIToken_Create_APIError tests Create handles API error
func TestInferenceAPIToken_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "test-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestInferenceAPIToken_Create_EmptyTokenId tests Create handles empty TokenId
func TestInferenceAPIToken_Create_EmptyTokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrString(""),
			Content:   ptrString("my-secret-token-content"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty TokenId")
}

// TestInferenceAPIToken_Read_Success tests Read retrieves token data
func TestInferenceAPIToken_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ interface{}, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64_IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-abc123"),
					Name:       ptrString("test-token"),
					Content:    ptrString("my-secret-token-content"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})
	d.SetId("zone-test123#token-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#token-abc123", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "test-token", d.Get("name"))
	assert.Equal(t, "my-secret-token-content", d.Get("content"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
}

// TestInferenceAPIToken_Read_NotFound tests Read handles token not found
func TestInferenceAPIToken_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ interface{}, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64_IAT(0),
			Tokens:     []*teov20220901.InferenceAPIToken{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})
	d.SetId("zone-test123#token-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestInferenceAPIToken_Update_Immutable tests Update rejects field changes
func TestInferenceAPIToken_Update_Immutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ interface{}, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64_IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-abc123"),
					Name:       ptrString("test-token"),
					Content:    ptrString("my-secret-token-content"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "modified-name",
	})
	d.SetId("zone-test123#token-abc123")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be changed")
}

// TestInferenceAPIToken_Delete_Success tests Delete removes token
func TestInferenceAPIToken_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "token-abc123", *request.TokenId)
		resp := teov20220901.NewDeleteInferenceAPITokenResponse()
		resp.Response = &teov20220901.DeleteInferenceAPITokenResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})
	d.SetId("zone-test123#token-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestInferenceAPIToken_Delete_APIError tests Delete handles API error
func TestInferenceAPIToken_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Token not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})
	d.SetId("zone-test123#token-abc123")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestInferenceAPIToken_Import tests Import via zone_id#token_id
func TestInferenceAPIToken_Import(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ interface{}, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64_IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-abc123"),
					Name:       ptrString("imported-token"),
					Content:    ptrString("imported-content"),
					CreateTime: ptrString("2024-06-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-test123#token-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#token-abc123", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "imported-token", d.Get("name"))
}

// TestInferenceAPIToken_Schema validates schema definition
func TestInferenceAPIToken_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check zone_id
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check name
	assert.Contains(t, res.Schema, "name")
	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.True(t, name.ForceNew)

	// Check computed fields
	assert.Contains(t, res.Schema, "token_id")
	tokenId := res.Schema["token_id"]
	assert.Equal(t, schema.TypeString, tokenId.Type)
	assert.True(t, tokenId.Computed)

	assert.Contains(t, res.Schema, "content")
	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Computed)
	assert.True(t, content.Sensitive)

	assert.Contains(t, res.Schema, "create_time")
	createTime := res.Schema["create_time"]
	assert.Equal(t, schema.TypeString, createTime.Type)
	assert.True(t, createTime.Computed)
}

func ptrInt64_IAT(i int64) *int64 {
	return &i
}
