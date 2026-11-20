package alb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcalb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

type mockMetaAlbListener struct {
	client *connectivity.TencentCloudClient
}

var _ tccommon.ProviderMeta = &mockMetaAlbListener{}

func (m *mockMetaAlbListener) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

func newMockMetaAlbListener() *mockMetaAlbListener {
	return &mockMetaAlbListener{client: &connectivity.TencentCloudClient{}}
}

func TestAlbListenerSchema(t *testing.T) {
	res := svcalb.ResourceTencentCloudAlbListener()
	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	for _, key := range []string{
		"load_balancer_id",
		"listener_port",
		"listener_protocol",
		"default_actions",
		"ca_certificate_ids",
		"ca_enabled",
		"certificate_ids",
		"client_token",
		"gzip_enabled",
		"http2enabled",
		"idle_timeout",
		"listener_name",
		"request_timeout",
		"security_policy_id",
		"tags",
		"x_forwarded_for_config",
	} {
		assert.Contains(t, res.Schema, key)
	}

	for _, key := range []string{
		"load_balancer_id",
		"listener_port",
		"listener_protocol",
	} {
		assert.True(t, res.Schema[key].Required)
		assert.True(t, res.Schema[key].ForceNew)
	}

	assert.True(t, res.Schema["default_actions"].Required)
	assert.Equal(t, 1, res.Schema["default_actions"].MaxItems)
	assert.Equal(t, 1, res.Schema["x_forwarded_for_config"].MaxItems)
}

func TestBuildAlbListenerDefaultActions(t *testing.T) {
	res := svcalb.ResourceTencentCloudAlbListener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"default_actions": []interface{}{
			map[string]interface{}{
				"type": "Forward",
				"target_group_config": []interface{}{
					map[string]interface{}{
						"target_groups": []interface{}{
							map[string]interface{}{
								"target_group_id": "lbtg-12345678",
								"weight":          80,
							},
						},
						"target_group_sticky_session": []interface{}{
							map[string]interface{}{
								"enabled": true,
								"timeout": 1200,
							},
						},
					},
				},
			},
		},
	})

	actions := svcalb.BuildAlbListenerDefaultActionsForTest(d)
	assert.Len(t, actions, 1)
	action := actions[0].(map[string]interface{})
	assert.Equal(t, "Forward", action["Type"])
	tgc := action["TargetGroupConfig"].(map[string]interface{})
	tgs := tgc["TargetGroups"].([]interface{})
	assert.Len(t, tgs, 1)
	assert.Equal(t, "lbtg-12345678", tgs[0].(map[string]interface{})["TargetGroupId"])
	assert.Equal(t, 80, tgs[0].(map[string]interface{})["Weight"])
	sticky := tgc["TargetGroupStickySession"].(map[string]interface{})
	assert.Equal(t, true, sticky["Enabled"])
	assert.Equal(t, 1200, sticky["Timeout"])
}

func TestBuildAlbListenerXForwardedForConfig(t *testing.T) {
	res := svcalb.ResourceTencentCloudAlbListener()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"x_forwarded_for_config": []interface{}{
			map[string]interface{}{
				"x_forwarded_for_mode":                    "append",
				"x_forwarded_for_host_enabled":            true,
				"x_forwarded_for_client_src_port_enabled": true,
			},
		},
	})

	xff := svcalb.BuildAlbListenerXForwardedForConfigForTest(d)
	assert.NotNil(t, xff)
	assert.Equal(t, "append", xff["XForwardedForMode"])
	assert.Equal(t, true, xff["XForwardedForHostEnabled"])
	assert.Equal(t, true, xff["XForwardedForClientSrcPortEnabled"])
}
