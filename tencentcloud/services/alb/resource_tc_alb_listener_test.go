package alb

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const (
	albTestLoadBalancerId = "alb-11111111"
	albTestListenerId     = "lst-22222222"
	albTestTargetGroupId  = "lbtg-33333333"
	albTestCreateTime     = "2026-01-01T08:30:00+08:00"
)

// albListenerTestContext wires a resource instance to a patched SDK client.
type albListenerTestContext struct {
	resource resource.Resource
	client   *albv20251030.Client
	schema   schema.Schema
}

func newAlbListenerTestContext(t *testing.T, patches *gomonkey.Patches) *albListenerTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseAlbV20251030Client", albClient)

	res := NewAlbListenerResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp,
	)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	return &albListenerTestContext{resource: res, client: albClient, schema: schemaResp.Schema}
}

func (c *albListenerTestContext) plan(t *testing.T, model AlbListenerResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *albListenerTestContext) state(t *testing.T, model AlbListenerResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

// buildDefaultActionsList constructs a default_actions list for tests.
func buildDefaultActionsList(ctx context.Context, t *testing.T) types.List {
	t.Helper()

	targetGroups, diags := types.ListValueFrom(ctx, albTargetGroupObjectType, []albTargetGroupModel{
		{
			TargetGroupId: types.StringValue(albTestTargetGroupId),
			Weight:        types.Int64Value(10),
		},
	})
	if diags.HasError() {
		t.Fatalf("build target_groups list failed: %v", diags)
	}

	stickySessions, diags := types.ListValueFrom(ctx, albTargetGroupStickySessionObjectType, []albTargetGroupStickySessionModel{
		{
			Enabled: types.BoolValue(true),
			Timeout: types.Int64Value(1000),
		},
	})
	if diags.HasError() {
		t.Fatalf("build target_group_sticky_session list failed: %v", diags)
	}

	configs, diags := types.ListValueFrom(ctx, albTargetGroupConfigObjectType, []albTargetGroupConfigModel{
		{
			TargetGroups:             targetGroups,
			TargetGroupStickySession: stickySessions,
		},
	})
	if diags.HasError() {
		t.Fatalf("build target_group_config list failed: %v", diags)
	}

	actions, diags := types.ListValueFrom(ctx, albDefaultActionObjectType, []albDefaultActionModel{
		{
			TargetGroupConfig: configs,
			Type:              types.StringValue("TargetGroup"),
		},
	})
	if diags.HasError() {
		t.Fatalf("build default_actions list failed: %v", diags)
	}

	return actions
}

// sampleAlbListenerModel returns a fully populated model for the CRUD tests.
func sampleAlbListenerModel(ctx context.Context, t *testing.T) AlbListenerResourceModel {
	t.Helper()

	tags, diags := types.ListValueFrom(ctx, albTagObjectType, []albTagModel{
		{TagKey: types.StringValue("env"), TagValue: types.StringValue("test")},
	})
	if diags.HasError() {
		t.Fatalf("build tags list failed: %v", diags)
	}

	xff, diags := types.ListValueFrom(ctx, albXForwardedForConfigObjectType, []albXForwardedForConfigModel{
		{
			XForwardedForAlbIdEnabled:         types.BoolValue(true),
			XForwardedForClientSrcPortEnabled: types.BoolValue(false),
			XForwardedForHostEnabled:          types.BoolValue(true),
			XForwardedForMode:                 types.StringValue("append"),
			XForwardedForPortEnabled:          types.BoolValue(true),
			XForwardedForProtoEnabled:         types.BoolValue(true),
			XTencentClientIdnEnabled:          types.BoolValue(false),
			XTencentClientSdnEnabled:          types.BoolValue(false),
			XTencentClientSerialEnabled:       types.BoolValue(false),
			XTencentClientVerifyEnabled:       types.BoolValue(false),
		},
	})
	if diags.HasError() {
		t.Fatalf("build x_forwarded_for_config list failed: %v", diags)
	}

	certificateIds, diags := types.SetValueFrom(ctx, types.StringType, []string{"cert-1"})
	if diags.HasError() {
		t.Fatalf("build certificate_ids set failed: %v", diags)
	}
	caCertificateIds, diags := types.SetValueFrom(ctx, types.StringType, []string{"ca-cert-1"})
	if diags.HasError() {
		t.Fatalf("build ca_certificate_ids set failed: %v", diags)
	}

	return AlbListenerResourceModel{
		Id:                  types.StringValue(albListenerId(albTestLoadBalancerId, albTestListenerId)),
		LoadBalancerId:      types.StringValue(albTestLoadBalancerId),
		ListenerId:          types.StringValue(albTestListenerId),
		ListenerName:        types.StringValue("tf-test-listener"),
		ListenerPort:        types.Int64Value(80),
		ListenerProtocol:    types.StringValue("HTTP"),
		DefaultActions:      buildDefaultActionsList(ctx, t),
		CaEnabled:           types.BoolValue(true),
		CaCertificateIds:    caCertificateIds,
		CertificateIds:      certificateIds,
		ClientToken:         types.StringValue("token-1"),
		GzipEnabled:         types.BoolValue(true),
		Http2Enabled:        types.BoolValue(false),
		IdleTimeout:         types.Int64Value(15),
		RequestTimeout:      types.Int64Value(60),
		SecurityPolicyId:    types.StringValue("tls-44444444"),
		Tags:                tags,
		XForwardedForConfig: xff,
		ListenerStatus:      types.StringValue(ALB_LISTENER_STATUS_ACTIVE),
		CreateTime:          types.StringValue(albTestCreateTime),
		ModifyTime:          types.StringValue(albTestCreateTime),
	}
}

// sampleAlbListenerOutput returns the cloud representation of the listener
// above, as DescribeListeners reports it.
func sampleAlbListenerOutput() *albv20251030.ListenerOutput {
	return &albv20251030.ListenerOutput{
		CaEnable:            helper.Bool(true),
		CreateTime:          helper.String(albTestCreateTime),
		GzipEnabled:         helper.Bool(true),
		Http2Enable:         helper.Bool(false),
		IdleTimeout:         helper.Uint64(15),
		ListenerId:          helper.String(albTestListenerId),
		ListenerName:        helper.String("tf-test-listener"),
		ListenerPort:        helper.Uint64(80),
		ListenerProtocol:    helper.String("HTTP"),
		ListenerStatus:      helper.String(ALB_LISTENER_STATUS_ACTIVE),
		ModifyTime:          helper.String(albTestCreateTime),
		RequestTimeout:      helper.Uint64(60),
		TlsSecurityPolicyId: helper.String("tls-44444444"),
		Tags: []*albv20251030.TagInfo{
			{TagKey: helper.String("env"), TagValue: helper.String("test")},
		},
		XForwardedForConfig: &albv20251030.XForwardedForConfig{
			XForwardedForAlbIdEnabled:         helper.Bool(true),
			XForwardedForClientSrcPortEnabled: helper.Bool(false),
			XForwardedForHostEnabled:          helper.Bool(true),
			XForwardedForMode:                 helper.String("append"),
			XForwardedForPortEnabled:          helper.Bool(true),
			XForwardedForProtoEnabled:         helper.Bool(true),
			XTencentClientIDNEnabled:          helper.Bool(false),
			XTencentClientSDNEnabled:          helper.Bool(false),
			XTencentClientSerialEnabled:       helper.Bool(false),
			XTencentClientVerifyEnabled:       helper.Bool(false),
		},
	}
}

func describeAlbListenersResponse(listeners ...*albv20251030.ListenerOutput) *albv20251030.DescribeListenersResponse {
	return &albv20251030.DescribeListenersResponse{
		Response: &albv20251030.DescribeListenersResponseParams{
			Listeners:  listeners,
			TotalCount: helper.Uint64(uint64(len(listeners))),
		},
	}
}

// nullStateRaw builds the empty (null) state value the framework hands to
// ImportState before any attribute is written.
func nullStateRaw(t *testing.T, s schema.Schema) tftypes.Value {
	t.Helper()

	return tftypes.NewValue(s.Type().TerraformType(context.Background()), nil)
}

// --- Metadata / Schema / Factory ---

// TestAlbListenerResourceMetadata validates the type name.
func TestAlbListenerResourceMetadata(t *testing.T) {
	r := &AlbListenerResource{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)

	if resp.TypeName != "tencentcloud_alb_listener" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_alb_listener")
	}
}

// TestAlbListenerResourceFactory validates the registry factory.
func TestAlbListenerResourceFactory(t *testing.T) {
	got := NewAlbListenerResource()
	if got == nil {
		t.Fatal("NewAlbListenerResource returned nil")
	}
	if _, ok := got.(*AlbListenerResource); !ok {
		t.Fatalf("factory returned %T, want *AlbListenerResource", got)
	}
}

// TestAlbListenerResourceSchema validates required/optional/computed flags on
// every top-level attribute.
func TestAlbListenerResourceSchema(t *testing.T) {
	r := &AlbListenerResource{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":                     {computed: true},
		"load_balancer_id":       {required: true},
		"listener_id":            {computed: true},
		"listener_name":          {optional: true, computed: true},
		"listener_port":          {required: true},
		"listener_protocol":      {required: true},
		"default_actions":        {required: true},
		"ca_enabled":             {optional: true, computed: true},
		"ca_certificate_ids":     {optional: true, computed: true},
		"certificate_ids":        {optional: true, computed: true},
		"client_token":           {optional: true},
		"gzip_enabled":           {optional: true, computed: true},
		"http2enabled":           {optional: true, computed: true},
		"idle_timeout":           {optional: true, computed: true},
		"request_timeout":        {optional: true, computed: true},
		"security_policy_id":     {optional: true, computed: true},
		"tags":                   {optional: true, computed: true},
		"x_forwarded_for_config": {optional: true, computed: true},
		"listener_status":        {computed: true},
		"create_time":            {computed: true},
		"modify_time":            {computed: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}

		switch a := attr.(type) {
		case schema.StringAttribute:
			assertAttributeFlags(t, name, a.Required, a.Optional, a.Computed, want.required, want.optional, want.computed)
		case schema.Int64Attribute:
			assertAttributeFlags(t, name, a.Required, a.Optional, a.Computed, want.required, want.optional, want.computed)
		case schema.BoolAttribute:
			assertAttributeFlags(t, name, a.Required, a.Optional, a.Computed, want.required, want.optional, want.computed)
		case schema.SetAttribute:
			assertAttributeFlags(t, name, a.Required, a.Optional, a.Computed, want.required, want.optional, want.computed)
		case schema.ListNestedAttribute:
			assertAttributeFlags(t, name, a.Required, a.Optional, a.Computed, want.required, want.optional, want.computed)
		default:
			t.Errorf("attribute %q unexpected type %T", name, attr)
		}
	}

	// create-only attributes must force replacement, because the modify API
	// does not accept them.
	for _, name := range []string{"load_balancer_id", "listener_port", "listener_protocol"} {
		var modifiers int
		switch a := resp.Schema.Attributes[name].(type) {
		case schema.StringAttribute:
			modifiers = len(a.PlanModifiers)
		case schema.Int64Attribute:
			modifiers = len(a.PlanModifiers)
		default:
			t.Fatalf("attribute %q unexpected type %T", name, a)
		}
		if modifiers == 0 {
			t.Errorf("attribute %q should carry a RequiresReplace plan modifier", name)
		}
	}
}

func assertAttributeFlags(t *testing.T, name string, required, optional, computed, wantRequired, wantOptional, wantComputed bool) {
	t.Helper()

	if required != wantRequired {
		t.Errorf("attribute %q Required=%v, want %v", name, required, wantRequired)
	}
	if optional != wantOptional {
		t.Errorf("attribute %q Optional=%v, want %v", name, optional, wantOptional)
	}
	if computed != wantComputed {
		t.Errorf("attribute %q Computed=%v, want %v", name, computed, wantComputed)
	}
}

// --- ID ---

// TestSplitAlbListenerId validates parsing of the composite id.
func TestSplitAlbListenerId(t *testing.T) {
	good := albTestLoadBalancerId + tccommon.FILED_SP + albTestListenerId
	lb, listener, ok := splitAlbListenerId(good)
	if !ok || lb != albTestLoadBalancerId || listener != albTestListenerId {
		t.Errorf("splitAlbListenerId(%q) = (%q, %q, %v)", good, lb, listener, ok)
	}

	for _, bad := range []string{
		"",
		albTestLoadBalancerId,
		tccommon.FILED_SP + albTestListenerId,
		albTestLoadBalancerId + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	} {
		if _, _, ok := splitAlbListenerId(bad); ok {
			t.Errorf("splitAlbListenerId(%q) ok=true, want false", bad)
		}
	}
}

// TestAlbListenerScope validates the diagnostic identity rendering.
func TestAlbListenerScope(t *testing.T) {
	if got := albListenerScope(albTestLoadBalancerId, ""); got != `load_balancer_id="alb-11111111"` {
		t.Errorf("albListenerScope without listener id = %q", got)
	}
	if got := albListenerScope(albTestLoadBalancerId, albTestListenerId); got != `load_balancer_id="alb-11111111", listener_id="lst-22222222"` {
		t.Errorf("albListenerScope = %q", got)
	}
}

// --- expand / flatten ---

// TestExpandAlbDefaultActions validates the plan list -> SDK struct conversion.
func TestExpandAlbDefaultActions(t *testing.T) {
	ctx := context.Background()

	actions, diags := expandAlbDefaultActions(ctx, buildDefaultActionsList(ctx, t))
	if diags.HasError() {
		t.Fatalf("expandAlbDefaultActions returned diagnostics: %v", diags)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0].Type == nil || *actions[0].Type != "TargetGroup" {
		t.Errorf("Type = %v, want TargetGroup", actions[0].Type)
	}
	if actions[0].TargetGroupConfig == nil {
		t.Fatal("TargetGroupConfig is nil")
	}
	groups := actions[0].TargetGroupConfig.TargetGroups
	if len(groups) != 1 {
		t.Fatalf("expected 1 target group, got %d", len(groups))
	}
	if groups[0].TargetGroupId == nil || *groups[0].TargetGroupId != albTestTargetGroupId {
		t.Errorf("TargetGroupId = %v, want %s", groups[0].TargetGroupId, albTestTargetGroupId)
	}
	if groups[0].Weight == nil || *groups[0].Weight != 10 {
		t.Errorf("Weight = %v, want 10", groups[0].Weight)
	}
	sticky := actions[0].TargetGroupConfig.TargetGroupStickySession
	if sticky == nil {
		t.Fatal("TargetGroupStickySession is nil")
	}
	if sticky.Enabled == nil || !*sticky.Enabled {
		t.Errorf("StickySession Enabled = %v, want true", sticky.Enabled)
	}
	if sticky.Timeout == nil || *sticky.Timeout != 1000 {
		t.Errorf("StickySession Timeout = %v, want 1000", sticky.Timeout)
	}
}

// TestExpandAlbDefaultActionsNull validates a null list yields no actions.
func TestExpandAlbDefaultActionsNull(t *testing.T) {
	actions, diags := expandAlbDefaultActions(context.Background(), types.ListNull(albDefaultActionObjectType))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if actions != nil {
		t.Errorf("expected nil actions for a null list, got %v", actions)
	}
}

// TestExpandStringSetNull validates a null set is sent as "absent".
func TestExpandStringSetNull(t *testing.T) {
	values, diags := expandStringSet(context.Background(), types.SetNull(types.StringType))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if values != nil {
		t.Errorf("expected nil values for a null set, got %v", values)
	}
}

// TestFlattenAlbTags verifies cloud tags are read back and the cloud-managed
// `qcs:` tags are skipped.
func TestFlattenAlbTags(t *testing.T) {
	ctx := context.Background()

	list, diags := flattenAlbTags(ctx, []*albv20251030.TagInfo{
		{TagKey: helper.String("env"), TagValue: helper.String("prod")},
		{TagKey: helper.String("qcs:cost"), TagValue: helper.String("x")},
		nil,
	}, types.ListNull(albTagObjectType))
	if diags.HasError() {
		t.Fatalf("flattenAlbTags returned diagnostics: %v", diags)
	}

	var out []albTagModel
	if diags := list.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("elements as failed: %v", diags)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(out))
	}
	if out[0].TagKey.ValueString() != "env" || out[0].TagValue.ValueString() != "prod" {
		t.Errorf("tag = %q=%q, want env=prod", out[0].TagKey.ValueString(), out[0].TagValue.ValueString())
	}
}

// TestFlattenAlbTagsKeepsPrev verifies an empty response keeps the configured
// tags, so the planned value never flips to null after apply.
func TestFlattenAlbTagsKeepsPrev(t *testing.T) {
	ctx := context.Background()

	prev, diags := types.ListValueFrom(ctx, albTagObjectType, []albTagModel{
		{TagKey: types.StringValue("env"), TagValue: types.StringValue("test")},
	})
	if diags.HasError() {
		t.Fatalf("build tags list failed: %v", diags)
	}

	list, diags := flattenAlbTags(ctx, nil, prev)
	if diags.HasError() {
		t.Fatalf("flattenAlbTags returned diagnostics: %v", diags)
	}
	if !list.Equal(prev) {
		t.Errorf("flattenAlbTags dropped the configured tags")
	}
}

// TestFlattenAlbXForwardedForConfig verifies the X-Forwarded-For object is
// mapped onto the list schema, and that omitted fields keep the previous value.
func TestFlattenAlbXForwardedForConfig(t *testing.T) {
	ctx := context.Background()

	list, diags := flattenAlbXForwardedForConfig(ctx, &albv20251030.XForwardedForConfig{
		XForwardedForAlbIdEnabled: helper.Bool(true),
		XForwardedForMode:         helper.String("remove"),
	}, types.ListNull(albXForwardedForConfigObjectType))
	if diags.HasError() {
		t.Fatalf("flattenAlbXForwardedForConfig returned diagnostics: %v", diags)
	}

	var out []albXForwardedForConfigModel
	if diags := list.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("elements as failed: %v", diags)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 config, got %d", len(out))
	}
	if !out[0].XForwardedForAlbIdEnabled.ValueBool() {
		t.Errorf("x_forwarded_for_alb_id_enabled = %v, want true", out[0].XForwardedForAlbIdEnabled)
	}
	if out[0].XForwardedForMode.ValueString() != "remove" {
		t.Errorf("x_forwarded_for_mode = %q, want remove", out[0].XForwardedForMode.ValueString())
	}
	if !out[0].XForwardedForProtoEnabled.IsNull() {
		t.Errorf("x_forwarded_for_proto_enabled = %v, want null", out[0].XForwardedForProtoEnabled)
	}

	nullList, diags := flattenAlbXForwardedForConfig(ctx, nil, types.ListNull(albXForwardedForConfigObjectType))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !nullList.IsNull() {
		t.Errorf("expected a null list when the cloud returns no config, got %v", nullList)
	}
}

// TestBoolValueOrPrev validates the previous-value fallback.
func TestBoolValueOrPrev(t *testing.T) {
	cases := []struct {
		name  string
		cloud *bool
		prev  types.Bool
		want  types.Bool
	}{
		{"cloud wins", helper.Bool(false), types.BoolValue(true), types.BoolValue(false)},
		{"prev kept", nil, types.BoolValue(true), types.BoolValue(true)},
		{"unknown prev becomes null", nil, types.BoolUnknown(), types.BoolNull()},
		{"null prev stays null", nil, types.BoolNull(), types.BoolNull()},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := boolValueOrPrev(c.cloud, c.prev); !got.Equal(c.want) {
				t.Fatalf("boolValueOrPrev = %v, want %v", got, c.want)
			}
		})
	}
}

// --- CRUD ---

// TestAlbListenerCreate verifies the create request and the read back state.
func TestAlbListenerCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	var capturedRequest *albv20251030.CreateListenerRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"CreateListenerWithContext",
		func(_ context.Context, request *albv20251030.CreateListenerRequest) (*albv20251030.CreateListenerResponse, error) {
			capturedRequest = request
			return &albv20251030.CreateListenerResponse{
				Response: &albv20251030.CreateListenerResponseParams{ListenerId: helper.String(albTestListenerId)},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeListenersWithContext",
		func(_ context.Context, request *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
			if request.LoadBalancerId == nil || *request.LoadBalancerId != albTestLoadBalancerId {
				t.Fatalf("DescribeListeners LoadBalancerId is not %q", albTestLoadBalancerId)
			}
			return describeAlbListenersResponse(sampleAlbListenerOutput()), nil
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, sampleAlbListenerModel(context.Background(), t))}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", createResp.Diagnostics.Errors())
	}

	if capturedRequest == nil {
		t.Fatal("CreateListener was not called")
	}
	if capturedRequest.LoadBalancerId == nil || *capturedRequest.LoadBalancerId != albTestLoadBalancerId {
		t.Errorf("CreateListener LoadBalancerId is not %q", albTestLoadBalancerId)
	}
	if capturedRequest.ListenerPort == nil || *capturedRequest.ListenerPort != 80 {
		t.Errorf("CreateListener ListenerPort is not 80")
	}
	if capturedRequest.ListenerProtocol == nil || *capturedRequest.ListenerProtocol != "HTTP" {
		t.Errorf("CreateListener ListenerProtocol is not HTTP")
	}
	if capturedRequest.ListenerName == nil || *capturedRequest.ListenerName != "tf-test-listener" {
		t.Errorf("CreateListener ListenerName is not tf-test-listener")
	}
	if capturedRequest.ClientToken == nil || *capturedRequest.ClientToken != "token-1" {
		t.Errorf("CreateListener ClientToken is not token-1")
	}
	if capturedRequest.IdleTimeout == nil || *capturedRequest.IdleTimeout != 15 {
		t.Errorf("CreateListener IdleTimeout is not 15")
	}
	if capturedRequest.RequestTimeout == nil || *capturedRequest.RequestTimeout != 60 {
		t.Errorf("CreateListener RequestTimeout is not 60")
	}
	if len(capturedRequest.DefaultActions) != 1 {
		t.Fatalf("CreateListener DefaultActions has %d elements, want 1", len(capturedRequest.DefaultActions))
	}
	if len(capturedRequest.CaCertificateIds) != 1 || *capturedRequest.CaCertificateIds[0] != "ca-cert-1" {
		t.Errorf("CreateListener CaCertificateIds is not [ca-cert-1]")
	}
	if len(capturedRequest.CertificateIds) != 1 || *capturedRequest.CertificateIds[0] != "cert-1" {
		t.Errorf("CreateListener CertificateIds is not [cert-1]")
	}
	if len(capturedRequest.Tags) != 1 || *capturedRequest.Tags[0].TagKey != "env" {
		t.Errorf("CreateListener Tags are not [env=test]")
	}
	if capturedRequest.XForwardedForConfig == nil || capturedRequest.XForwardedForConfig.XForwardedForMode == nil ||
		*capturedRequest.XForwardedForConfig.XForwardedForMode != "append" {
		t.Errorf("CreateListener XForwardedForConfig.XForwardedForMode is not append")
	}

	var state AlbListenerResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if state.Id.ValueString() != albTestLoadBalancerId+tccommon.FILED_SP+albTestListenerId {
		t.Errorf("state id = %q, want composite id", state.Id.ValueString())
	}
	if state.ListenerId.ValueString() != albTestListenerId {
		t.Errorf("state listener_id = %q", state.ListenerId.ValueString())
	}
	if state.ListenerStatus.ValueString() != ALB_LISTENER_STATUS_ACTIVE {
		t.Errorf("state listener_status = %q", state.ListenerStatus.ValueString())
	}
	if state.SecurityPolicyId.ValueString() != "tls-44444444" {
		t.Errorf("state security_policy_id = %q", state.SecurityPolicyId.ValueString())
	}
	// default_actions is not echoed back by DescribeListeners, so the planned
	// value must survive the read back.
	if state.DefaultActions.IsNull() || len(state.DefaultActions.Elements()) != 1 {
		t.Errorf("state default_actions was dropped")
	}
}

// TestAlbListenerReadNotFound verifies a removed listener is dropped from state.
func TestAlbListenerReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeListenersWithContext",
		func(_ context.Context, _ *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
			return describeAlbListenersResponse(), nil
		},
	)

	model := sampleAlbListenerModel(context.Background(), t)
	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics.Errors())
	}

	if !readResp.State.Raw.IsNull() {
		t.Errorf("expected the resource to be removed from state, got %s", readResp.State.Raw.String())
	}
}

// TestAlbListenerReadInvalidId verifies a malformed id is rejected.
func TestAlbListenerReadInvalidId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	model := sampleAlbListenerModel(context.Background(), t)
	model.Id = types.StringValue(albTestListenerId)

	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read accepted a malformed id")
	}
}

// TestAlbListenerUpdate verifies the modify request and the refreshed state.
func TestAlbListenerUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	var capturedRequest *albv20251030.ModifyListenerAttributesRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"ModifyListenerAttributesWithContext",
		func(_ context.Context, request *albv20251030.ModifyListenerAttributesRequest) (*albv20251030.ModifyListenerAttributesResponse, error) {
			capturedRequest = request
			return &albv20251030.ModifyListenerAttributesResponse{
				Response: &albv20251030.ModifyListenerAttributesResponseParams{RequestId: helper.String("req-1")},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeListenersWithContext",
		func(_ context.Context, _ *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
			return describeAlbListenersResponse(sampleAlbListenerOutput()), nil
		},
	)

	stateModel := sampleAlbListenerModel(context.Background(), t)
	planModel := sampleAlbListenerModel(context.Background(), t)
	planModel.ListenerName = types.StringValue("tf-test-listener-updated")

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, planModel), State: testContext.state(t, stateModel)},
		updateResp,
	)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", updateResp.Diagnostics.Errors())
	}

	if capturedRequest == nil {
		t.Fatal("ModifyListenerAttributes was not called")
	}
	if capturedRequest.ListenerId == nil || *capturedRequest.ListenerId != albTestListenerId {
		t.Errorf("ModifyListenerAttributes ListenerId is not %q", albTestListenerId)
	}
	if capturedRequest.LoadBalancerId == nil || *capturedRequest.LoadBalancerId != albTestLoadBalancerId {
		t.Errorf("ModifyListenerAttributes LoadBalancerId is not %q", albTestLoadBalancerId)
	}
	if capturedRequest.ListenerName == nil || *capturedRequest.ListenerName != "tf-test-listener-updated" {
		t.Errorf("ModifyListenerAttributes ListenerName is not the planned value")
	}
	if len(capturedRequest.DefaultActions) != 1 {
		t.Errorf("ModifyListenerAttributes DefaultActions has %d elements, want 1", len(capturedRequest.DefaultActions))
	}

	var state AlbListenerResourceModel
	if diags := updateResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if state.Id.ValueString() != albTestLoadBalancerId+tccommon.FILED_SP+albTestListenerId {
		t.Errorf("state id = %q, want composite id", state.Id.ValueString())
	}
	// client_token is create-only and must be preserved.
	if state.ClientToken.ValueString() != "token-1" {
		t.Errorf("state client_token = %q, want token-1", state.ClientToken.ValueString())
	}
}

// TestAlbListenerDelete verifies the delete request carries both cloud ids and
// that an already deleted listener is a successful deletion.
func TestAlbListenerDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	var capturedRequest *albv20251030.DeleteListenerRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteListenerWithContext",
		func(_ context.Context, request *albv20251030.DeleteListenerRequest) (*albv20251030.DeleteListenerResponse, error) {
			capturedRequest = request
			return &albv20251030.DeleteListenerResponse{
				Response: &albv20251030.DeleteListenerResponseParams{RequestId: helper.String("req-1")},
			}, nil
		},
	)

	model := sampleAlbListenerModel(context.Background(), t)
	deleteResp := &resource.DeleteResponse{State: testContext.state(t, model)}
	testContext.resource.Delete(context.Background(), resource.DeleteRequest{State: testContext.state(t, model)}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete returned diagnostics: %v", deleteResp.Diagnostics.Errors())
	}

	if capturedRequest == nil {
		t.Fatal("DeleteListener was not called")
	}
	if capturedRequest.LoadBalancerId == nil || *capturedRequest.LoadBalancerId != albTestLoadBalancerId {
		t.Errorf("DeleteListener LoadBalancerId is not %q", albTestLoadBalancerId)
	}
	if len(capturedRequest.ListenerIds) != 1 || *capturedRequest.ListenerIds[0] != albTestListenerId {
		t.Errorf("DeleteListener ListenerIds is not [%q]", albTestListenerId)
	}
}

// TestAlbListenerImportState verifies the composite import id is split into its
// two state attributes.
func TestAlbListenerImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema, Raw: nullStateRaw(t, testContext.schema)}}
	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: albTestLoadBalancerId + tccommon.FILED_SP + albTestListenerId},
		importResp,
	)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState returned diagnostics: %v", importResp.Diagnostics.Errors())
	}

	var state AlbListenerResourceModel
	if diags := importResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	if state.LoadBalancerId.ValueString() != albTestLoadBalancerId {
		t.Errorf("state load_balancer_id = %q", state.LoadBalancerId.ValueString())
	}
	if state.ListenerId.ValueString() != albTestListenerId {
		t.Errorf("state listener_id = %q", state.ListenerId.ValueString())
	}

	invalidResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema, Raw: nullStateRaw(t, testContext.schema)}}
	importer.ImportState(context.Background(), resource.ImportStateRequest{ID: albTestListenerId}, invalidResp)
	if !invalidResp.Diagnostics.HasError() {
		t.Fatal("ImportState accepted a malformed import id")
	}
}
