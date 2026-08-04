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

// go test ./tencentcloud/services/teo/ -run "TestAccTeoInferenceAPIToken" -v -count=1 -gcflags="all=-l"
// TestAccTeoInferenceAPIToken_Success tests Create calls API and sets ID and computed attributes
func TestAccTeoInferenceAPIToken_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPIToken", func(request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "my-token", *request.Name)
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrString("token-abc123"),
			Content:   ptrString("dGhpcyBpcyBhIHRva2Vu"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())

	tokenId := d.Get("token_id").(string)
	assert.Equal(t, "token-abc123", tokenId)

	content := d.Get("content").(string)
	assert.Equal(t, "dGhpcyBpcyBhIHRva2Vu", content)
}

// TestAccTeoInferenceAPIToken_APIError tests Create handles API error
func TestAccTeoInferenceAPIToken_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPIToken", func(request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestAccTeoInferenceAPIToken_EmptyResponse tests Create handles empty API response
func TestAccTeoInferenceAPIToken_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPIToken", func(request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

// TestAccTeoInferenceAPIToken_MissingZoneId tests error handling when zone_id is missing
func TestAccTeoInferenceAPIToken_MissingZoneId(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "my-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone_id is required")
}

// TestAccTeoInferenceAPIToken_MissingName tests error handling when name is missing
func TestAccTeoInferenceAPIToken_MissingName(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceAPITokenOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

// TestAccTeoInferenceAPIToken_Read tests Read is no-op
func TestAccTeoInferenceAPIToken_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceAPITokenOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestAccTeoInferenceAPIToken_Delete tests Delete is no-op
func TestAccTeoInferenceAPIToken_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceAPITokenOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})
	d.SetId("zone-12345678")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestAccTeoInferenceAPIToken_Schema validates schema definition
func TestAccTeoInferenceAPIToken_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceAPITokenOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "token_id")
	assert.Contains(t, res.Schema, "content")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.True(t, name.ForceNew)

	tokenId := res.Schema["token_id"]
	assert.Equal(t, schema.TypeString, tokenId.Type)
	assert.True(t, tokenId.Computed)

	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Computed)
	assert.True(t, content.Sensitive)
}
