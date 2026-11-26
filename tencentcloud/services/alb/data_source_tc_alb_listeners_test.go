package alb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svab "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

func albListenersStr(s string) *string    { return &s }
func albListenersBool(b bool) *bool       { return &b }
func albListenersUint64(v uint64) *uint64 { return &v }

func albListenersFilterItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":   tftypes.String,
			"values": tftypes.List{ElementType: tftypes.String},
		},
	}
}

func albListenersTagItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"tag_key":   tftypes.String,
			"tag_value": tftypes.String,
		},
	}
}

func albListenersXffType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"x_forwarded_for_alb_id_enabled":          tftypes.Bool,
			"x_forwarded_for_client_src_port_enabled": tftypes.Bool,
			"x_forwarded_for_host_enabled":            tftypes.Bool,
			"x_forwarded_for_mode":                    tftypes.String,
			"x_forwarded_for_port_enabled":            tftypes.Bool,
			"x_forwarded_for_proto_enabled":           tftypes.Bool,
			"x_tencent_client_idn_enabled":            tftypes.Bool,
			"x_tencent_client_sdn_enabled":            tftypes.Bool,
			"x_tencent_client_serial_enabled":         tftypes.Bool,
			"x_tencent_client_verify_enabled":         tftypes.Bool,
		},
	}
}

func albListenersItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"ca_enable":              tftypes.Bool,
			"create_time":            tftypes.String,
			"gzip_enabled":           tftypes.Bool,
			"http2_enable":           tftypes.Bool,
			"idle_timeout":           tftypes.Number,
			"listener_id":            tftypes.String,
			"listener_name":          tftypes.String,
			"listener_port":          tftypes.Number,
			"listener_protocol":      tftypes.String,
			"listener_status":        tftypes.String,
			"modify_time":            tftypes.String,
			"request_timeout":        tftypes.Number,
			"tags":                   tftypes.List{ElementType: albListenersTagItemType()},
			"tls_security_policy_id": tftypes.String,
			"x_forwarded_for_config": albListenersXffType(),
		},
	}
}

func albListenersRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":               tftypes.String,
			"load_balancer_id": tftypes.String,
			"listener_ids":     tftypes.List{ElementType: tftypes.String},
			"filters":          tftypes.List{ElementType: albListenersFilterItemType()},
			"listeners":        tftypes.List{ElementType: albListenersItemType()},
		},
	}
}

func buildAlbListenersRawValue(loadBalancerId string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":               tftypes.NewValue(tftypes.String, nil),
		"load_balancer_id": tftypes.NewValue(tftypes.String, loadBalancerId),
		"listener_ids":     tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
		"filters":          tftypes.NewValue(tftypes.List{ElementType: albListenersFilterItemType()}, nil),
		"listeners":        tftypes.NewValue(tftypes.List{ElementType: albListenersItemType()}, nil),
	}

	return tftypes.NewValue(albListenersRootType(), attrs)
}

func albListenersSchema(t *testing.T) schema.Schema {
	d := svab.NewAlbListenersDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

func setAlbListenersClient(d *svab.AlbListenersDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func albListenersFakeListener() *albv20251030.ListenerOutput {
	return &albv20251030.ListenerOutput{
		CaEnable:         albListenersBool(true),
		CreateTime:       albListenersStr("2025-01-01T08:30:00+08:00"),
		GzipEnabled:      albListenersBool(false),
		Http2Enable:      albListenersBool(true),
		IdleTimeout:      albListenersUint64(60),
		ListenerId:       albListenersStr("lst-a1b2c3d4"),
		ListenerName:     albListenersStr("tf-test-listener"),
		ListenerPort:     albListenersUint64(80),
		ListenerProtocol: albListenersStr("HTTP"),
		ListenerStatus:   albListenersStr("Active"),
		ModifyTime:       albListenersStr("2025-01-02T08:30:00+08:00"),
		RequestTimeout:   albListenersUint64(30),
		Tags: []*albv20251030.TagInfo{
			{TagKey: albListenersStr("env"), TagValue: albListenersStr("tf-test")},
		},
		TlsSecurityPolicyId: albListenersStr("policy-12345"),
		XForwardedForConfig: &albv20251030.XForwardedForConfig{
			XForwardedForAlbIdEnabled:         albListenersBool(true),
			XForwardedForClientSrcPortEnabled: albListenersBool(false),
			XForwardedForHostEnabled:          albListenersBool(true),
			XForwardedForMode:                 albListenersStr("append"),
			XForwardedForPortEnabled:          albListenersBool(true),
			XForwardedForProtoEnabled:         albListenersBool(true),
			XTencentClientIDNEnabled:          albListenersBool(false),
			XTencentClientSDNEnabled:          albListenersBool(false),
			XTencentClientSerialEnabled:       albListenersBool(false),
			XTencentClientVerifyEnabled:       albListenersBool(false),
		},
	}
}

func TestAlbListeners_MetadataAndSchema(t *testing.T) {
	d := svab.NewAlbListenersDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_alb_listeners", metaResp.TypeName)

	sch := albListenersSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["load_balancer_id"].(schema.StringAttribute).Required)
	assert.True(t, attrs["listener_ids"].(schema.ListAttribute).Optional)
	assert.True(t, attrs["filters"].(schema.ListNestedAttribute).Optional)
	assert.True(t, attrs["listeners"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["id"].(schema.StringAttribute).Computed)

	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
	assert.NotContains(t, attrs, "max_results")
	assert.NotContains(t, attrs, "next_token")
}

func TestAlbListeners_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(request *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		if request.LoadBalancerId == nil || *request.LoadBalancerId != "alb-f8q2xk9m" {
			t.Errorf("expected load balancer id to be passed to the API request")
		}
		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: []*albv20251030.ListenerOutput{albListenersFakeListener()},
			RequestId: albListenersStr("fake-request-id"),
		}}, nil
	})

	ds := svab.NewAlbListenersDataSource().(*svab.AlbListenersDataSource)
	setAlbListenersClient(ds, client)

	sch := albListenersSchema(t)
	raw := buildAlbListenersRawValue("alb-f8q2xk9m")
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svab.AlbListenersDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.NotEmpty(t, result.Id.ValueString())
	assert.False(t, result.Listeners.IsNull())

	var items []svab.AlbListenersItemModel
	resp.Diagnostics.Append(result.Listeners.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	if assert.Len(t, items, 1) {
		item := items[0]
		assert.True(t, item.CaEnable.ValueBool())
		assert.Equal(t, "2025-01-01T08:30:00+08:00", item.CreateTime.ValueString())
		assert.False(t, item.GzipEnabled.ValueBool())
		assert.True(t, item.Http2Enable.ValueBool())
		assert.Equal(t, int64(60), item.IdleTimeout.ValueInt64())
		assert.Equal(t, "lst-a1b2c3d4", item.ListenerId.ValueString())
		assert.Equal(t, "tf-test-listener", item.ListenerName.ValueString())
		assert.Equal(t, int64(80), item.ListenerPort.ValueInt64())
		assert.Equal(t, "HTTP", item.ListenerProtocol.ValueString())
		assert.Equal(t, "Active", item.ListenerStatus.ValueString())
		assert.Equal(t, "2025-01-02T08:30:00+08:00", item.ModifyTime.ValueString())
		assert.Equal(t, int64(30), item.RequestTimeout.ValueInt64())
		assert.Equal(t, "policy-12345", item.TlsSecurityPolicyId.ValueString())

		var tags []svab.AlbListenersTagItemModel
		resp.Diagnostics.Append(item.Tags.ElementsAs(context.Background(), &tags, false)...)
		assert.False(t, resp.Diagnostics.HasError())
		if assert.Len(t, tags, 1) {
			assert.Equal(t, "env", tags[0].TagKey.ValueString())
			assert.Equal(t, "tf-test", tags[0].TagValue.ValueString())
		}

		var xff svab.AlbListenersXForwardedForModel
		resp.Diagnostics.Append(item.XForwardedForConfig.As(context.Background(), &xff, basetypes.ObjectAsOptions{})...)
		assert.False(t, resp.Diagnostics.HasError())
		assert.True(t, xff.XForwardedForAlbIdEnabled.ValueBool())
		assert.Equal(t, "append", xff.XForwardedForMode.ValueString())
	}
}

func TestAlbListeners_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(request *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{}}, nil
	})

	ds := svab.NewAlbListenersDataSource().(*svab.AlbListenersDataSource)
	setAlbListenersClient(ds, client)

	sch := albListenersSchema(t)
	raw := buildAlbListenersRawValue("alb-empty")
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svab.AlbListenersDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.False(t, result.Listeners.IsNull())

	var items []svab.AlbListenersItemModel
	resp.Diagnostics.Append(result.Listeners.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

func TestAlbListeners_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(request *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svab.NewAlbListenersDataSource().(*svab.AlbListenersDataSource)
	setAlbListenersClient(ds, client)

	sch := albListenersSchema(t)
	raw := buildAlbListenersRawValue("alb-f8q2xk9m")
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

func TestAlbListenersService_ReadAllPages(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)

	callCount := 0
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(request *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		callCount++
		if request.MaxResults == nil || *request.MaxResults != 100 {
			t.Errorf("expected MaxResults 100 on call %d, got %v", callCount, request.MaxResults)
		}
		switch callCount {
		case 1:
			if request.NextToken != nil && *request.NextToken != "" {
				t.Errorf("expected empty NextToken on first call, got %s", *request.NextToken)
			}
			return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
				Listeners: []*albv20251030.ListenerOutput{albListenersFakeListener()},
				NextToken: albListenersStr("next-page"),
			}}, nil
		case 2:
			if request.NextToken == nil || *request.NextToken != "next-page" {
				t.Errorf("expected NextToken next-page on second call")
			}
			second := albListenersFakeListener()
			second.ListenerId = albListenersStr("lst-e5f6g7h8")
			return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
				Listeners: []*albv20251030.ListenerOutput{second},
			}}, nil
		default:
			t.Errorf("unexpected extra API call %d", callCount)
			return &albv20251030.DescribeListenersResponse{}, nil
		}
	})

	service := svab.NewAlbService(client)
	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = albListenersStr("alb-f8q2xk9m")

	listeners, err := service.DescribeListeners(context.Background(), describeListenersRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
	if len(listeners) != 2 {
		t.Fatalf("expected 2 listeners across pages, got %d", len(listeners))
	}
	if listeners[0].ListenerId == nil || *listeners[0].ListenerId != "lst-a1b2c3d4" {
		t.Errorf("unexpected first listener id: %v", listeners[0].ListenerId)
	}
	if listeners[1].ListenerId == nil || *listeners[1].ListenerId != "lst-e5f6g7h8" {
		t.Errorf("unexpected second listener id: %v", listeners[1].ListenerId)
	}
}
