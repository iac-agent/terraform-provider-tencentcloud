package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestInferenceAPITokenV9" -v -count=1 -gcflags="all=-l"

// TestInferenceAPITokenV9_Create_Success tests Create calls API and sets ID
func TestInferenceAPITokenV9_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPIToken", func(request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrString("token-abcdefgh"),
			Content:   ptrString("my-secret-token-content"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokens", func(request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: helper.Int64(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-abcdefgh"),
					Name:       ptrString("my-token"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "token-abcdefgh", d.Id())
	assert.Equal(t, "my-secret-token-content", d.Get("content"))
}

// TestInferenceAPITokenV9_Create_APIError tests Create handles API error
func TestInferenceAPITokenV9_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPIToken", func(request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestInferenceAPITokenV9_Create_NilResponse tests Create handles nil response
func TestInferenceAPITokenV9_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPIToken", func(request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		return nil, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

// TestInferenceAPITokenV9_Read_Success tests Read retrieves token data
func TestInferenceAPITokenV9_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokens", func(request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: helper.Int64(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-abcdefgh"),
					Name:       ptrString("my-token"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-token",
	})
	d.SetId("token-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "my-token", d.Get("name"))
	assert.Equal(t, "token-abcdefgh", d.Get("token_id"))
}

// TestInferenceAPITokenV9_Read_NotFound tests Read handles token not found
func TestInferenceAPITokenV9_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokens", func(request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: helper.Int64(0),
			Tokens:     []*teov20220901.InferenceAPIToken{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-token",
	})
	d.SetId("token-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestInferenceAPITokenV9_Read_NoMatchingToken tests Read handles no matching token in list
func TestInferenceAPITokenV9_Read_NoMatchingToken(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokens", func(request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: helper.Int64(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrString("token-other"),
					Name:       ptrString("other-token"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-token",
	})
	d.SetId("token-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestInferenceAPITokenV9_Delete_Success tests Delete removes token
func TestInferenceAPITokenV9_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPIToken", func(request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		resp := teov20220901.NewDeleteInferenceAPITokenResponse()
		resp.Response = &teov20220901.DeleteInferenceAPITokenResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-token",
	})
	d.SetId("token-abcdefgh")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestInferenceAPITokenV9_Delete_APIError tests Delete handles API error
func TestInferenceAPITokenV9_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPIToken", func(request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Token not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "my-token",
	})
	d.SetId("token-abcdefgh")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestInferenceAPITokenV9_Schema validates schema definition
func TestInferenceAPITokenV9_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceAPITokenV9()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Delete)

	// No Update or Import for this resource
	assert.Nil(t, res.Update)
	assert.Nil(t, res.Importer)

	// Check required fields with ForceNew
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

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
}
