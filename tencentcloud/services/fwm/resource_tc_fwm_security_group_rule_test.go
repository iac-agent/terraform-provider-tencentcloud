package fwm_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

// go test ./tencentcloud/services/fwm/ -run "TestFwmSecurityGroupRule" -v -count=1 -gcflags="all=-l"

const (
	fwmSgRuleGroupId = "cfwg-1a2b3c4d"
	fwmSgRuleRuleId  = "rule-1a2b3c4d"
)

func ptrStringFwmSgRule(s string) *string {
	return &s
}

func ptrInt64FwmSgRule(i int64) *int64 {
	return &i
}

// fwmSgRuleAttrTypes returns the attribute types of the resource schema, keyed
// by schema attribute name. It is used to build raw tftypes values.
func fwmSgRuleAttrTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"id":       tftypes.String,
		"group_id": tftypes.String,
		"rules": tftypes.List{ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"ip_version":          tftypes.String,
				"source_content":      tftypes.String,
				"source_type":         tftypes.String,
				"dest_content":        tftypes.String,
				"dest_type":           tftypes.String,
				"protocol":            tftypes.String,
				"port":                tftypes.String,
				"service_template_id": tftypes.String,
				"rule_action":         tftypes.String,
				"description":         tftypes.String,
				"order_index":         tftypes.Int64,
				"rule_id":             tftypes.String,
				"scope":               tftypes.String,
				"protocol_port_type":  tftypes.Int64,
				"belong_member_id":    tftypes.String,
			},
		}},
		"rule": tftypes.List{ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"rule_id":             tftypes.String,
				"sequence":            tftypes.Int64,
				"region":              tftypes.String,
				"ip_version":          tftypes.String,
				"src_content":         tftypes.String,
				"src_type":            tftypes.String,
				"dst_content":         tftypes.String,
				"dst_type":            tftypes.String,
				"protocol":            tftypes.String,
				"protocol_port_type":  tftypes.Int64,
				"service_template_id": tftypes.String,
				"dst_port":            tftypes.String,
				"rule_action":         tftypes.String,
				"detail":              tftypes.String,
				"rule_source":         tftypes.String,
				"scope":               tftypes.String,
			},
		}},
	}
}

// fwmSgRuleRawRule builds one raw `rules` element.
func fwmSgRuleRawRule() tftypes.Value {
	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: fwmSgRuleAttrTypes()["rules"].(tftypes.List).ElementType.(tftypes.Object).AttributeTypes,
	}, map[string]tftypes.Value{
		"ip_version":          tftypes.NewValue(tftypes.String, "ipv4"),
		"source_content":      tftypes.NewValue(tftypes.String, "10.0.0.0/16"),
		"source_type":         tftypes.NewValue(tftypes.String, "net"),
		"dest_content":        tftypes.NewValue(tftypes.String, "10.0.1.0/24"),
		"dest_type":           tftypes.NewValue(tftypes.String, "net"),
		"protocol":            tftypes.NewValue(tftypes.String, "TCP"),
		"port":                tftypes.NewValue(tftypes.String, "80"),
		"service_template_id": tftypes.NewValue(tftypes.String, nil),
		"rule_action":         tftypes.NewValue(tftypes.String, "2"),
		"description":         tftypes.NewValue(tftypes.String, "allow http"),
		"order_index":         tftypes.NewValue(tftypes.Int64, int64(1)),
		"rule_id":             tftypes.NewValue(tftypes.String, nil),
		"scope":               tftypes.NewValue(tftypes.String, "1"),
		"protocol_port_type":  tftypes.NewValue(tftypes.Int64, int64(0)),
		"belong_member_id":    tftypes.NewValue(tftypes.String, nil),
	})
}

// fwmSgRuleResourceSchema returns the framework resource schema.
func fwmSgRuleResourceSchema(t *testing.T) *resource.SchemaResponse {
	r := fwm.NewFwmSecurityGroupRuleResource()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp
}

// newFwmSgRulePlan builds a raw plan value.
func newFwmSgRulePlan(groupId string, rules interface{}, ruleDetail interface{}) tftypes.Value {
	rulesVal := tftypes.NewValue(fwmSgRuleAttrTypes()["rules"], nil)
	if rules != nil {
		rulesVal = rules.(tftypes.Value)
	}

	ruleVal := tftypes.NewValue(fwmSgRuleAttrTypes()["rule"], nil)
	if ruleDetail != nil {
		ruleVal = ruleDetail.(tftypes.Value)
	}

	groupIdVal := tftypes.NewValue(tftypes.String, nil)
	if groupId != "" {
		groupIdVal = tftypes.NewValue(tftypes.String, groupId)
	}

	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: fwmSgRuleAttrTypes(),
	}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, nil),
		"group_id": groupIdVal,
		"rules":    rulesVal,
		"rule":     ruleVal,
	})
}

// fwmSgRuleRawList wraps rules into a list value.
func fwmSgRuleRawList(rules []tftypes.Value) tftypes.Value {
	if len(rules) == 0 {
		return tftypes.NewValue(fwmSgRuleAttrTypes()["rules"], nil)
	}
	return tftypes.NewValue(fwmSgRuleAttrTypes()["rules"], rules)
}

// fwmSgRuleRawDetailList wraps one rule detail into a list value.
func fwmSgRuleRawDetailList() tftypes.Value {
	detail := tftypes.NewValue(tftypes.Object{
		AttributeTypes: fwmSgRuleAttrTypes()["rule"].(tftypes.List).ElementType.(tftypes.Object).AttributeTypes,
	}, map[string]tftypes.Value{
		"rule_id":             tftypes.NewValue(tftypes.String, fwmSgRuleRuleId),
		"sequence":            tftypes.NewValue(tftypes.Int64, int64(1)),
		"region":              tftypes.NewValue(tftypes.String, "ap-guangzhou"),
		"ip_version":          tftypes.NewValue(tftypes.String, "ipv4"),
		"src_content":         tftypes.NewValue(tftypes.String, "10.0.0.0/16"),
		"src_type":            tftypes.NewValue(tftypes.String, "net"),
		"dst_content":         tftypes.NewValue(tftypes.String, "10.0.1.0/24"),
		"dst_type":            tftypes.NewValue(tftypes.String, "net"),
		"protocol":            tftypes.NewValue(tftypes.String, "TCP"),
		"protocol_port_type":  tftypes.NewValue(tftypes.Int64, int64(0)),
		"service_template_id": tftypes.NewValue(tftypes.String, nil),
		"dst_port":            tftypes.NewValue(tftypes.String, "80"),
		"rule_action":         tftypes.NewValue(tftypes.String, "2"),
		"detail":              tftypes.NewValue(tftypes.String, "allow http"),
		"rule_source":         tftypes.NewValue(tftypes.String, "0"),
		"scope":               tftypes.NewValue(tftypes.String, "1"),
	})
	return tftypes.NewValue(fwmSgRuleAttrTypes()["rule"], []tftypes.Value{detail})
}

// setFwmSgRuleClient injects a mock connectivity client into the resource.
func setFwmSgRuleClient(t *testing.T, client *connectivity.TencentCloudClient) resource.Resource {
	r := fwm.NewFwmSecurityGroupRuleResource()
	configureReq := resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics)
	}
	return r
}

// newFwmSgRuleCreateRequest builds a CreateRequest carrying the given raw plan.
func newFwmSgRuleCreateRequest(t *testing.T, raw tftypes.Value) resource.CreateRequest {
	schemaResp := fwmSgRuleResourceSchema(t)
	req := resource.CreateRequest{}
	req.Plan = tfsdk.Plan{Raw: raw, Schema: schemaResp.Schema}
	return req
}

// newFwmSgRuleCreateResponse builds a CreateResponse whose state mirrors the
// plan, as the framework does before invoking Create.
func newFwmSgRuleCreateResponse(raw tftypes.Value) *resource.CreateResponse {
	schemaResp := fwmSgRuleResourceSchema(t)
	resp := &resource.CreateResponse{}
	resp.State = tfsdk.State{Raw: raw, Schema: schemaResp.Schema}
	return resp
}

// fwmSgRuleDescribeResponse builds a describe response carrying one rule.
func fwmSgRuleDescribeResponse() *fwmv20250611.DescribeSecurityGroupRuleResponse {
	resp := &fwmv20250611.DescribeSecurityGroupRuleResponse{}
	resp.Response = &fwmv20250611.DescribeSecurityGroupRuleResponseParams{
		RequestId: ptrStringFwmSgRule("fake-request-id"),
		Rule: &fwmv20250611.SgRuleResp{
			RuleId:           ptrStringFwmSgRule(fwmSgRuleRuleId),
			Sequence:         ptrInt64FwmSgRule(1),
			Region:           ptrStringFwmSgRule("ap-guangzhou"),
			IpVersion:        ptrStringFwmSgRule("ipv4"),
			SrcContent:       ptrStringFwmSgRule("10.0.0.0/16"),
			SrcType:          ptrStringFwmSgRule("net"),
			DstContent:       ptrStringFwmSgRule("10.0.1.0/24"),
			DstType:          ptrStringFwmSgRule("net"),
			Protocol:         ptrStringFwmSgRule("TCP"),
			ProtocolPortType: ptrInt64FwmSgRule(0),
			DstPort:          ptrStringFwmSgRule("80"),
			RuleAction:       ptrStringFwmSgRule("2"),
			Detail:           ptrStringFwmSgRule("allow http"),
			RuleSource:       ptrStringFwmSgRule("0"),
			Scope:            ptrStringFwmSgRule("1"),
		},
	}
	return resp
}

// TestFwmSecurityGroupRule_Create_Success covers the happy path of Create,
// including the composed resource identifier and the read-back of the rule detail.
func TestFwmSecurityGroupRule_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	var gotRequest *fwmv20250611.CreateSecurityGroupRuleRequest
	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		gotRequest = request
		result := &fwmv20250611.CreateSecurityGroupRuleResponse{}
		result.Response = &fwmv20250611.CreateSecurityGroupRuleResponseParams{
			RequestId: ptrStringFwmSgRule("fake-request-id"),
			RuleIds:   []*string{ptrStringFwmSgRule(fwmSgRuleRuleId)},
		}
		return result, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		assert.Equal(t, fwmSgRuleGroupId, *request.GroupId)
		assert.Equal(t, fwmSgRuleRuleId, *request.RuleId)
		return fwmSgRuleDescribeResponse(), nil
	})

	r := setFwmSgRuleClient(t, client)

	raw := newFwmSgRulePlan(fwmSgRuleGroupId, fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}), nil)
	req := newFwmSgRuleCreateRequest(t, raw)
	resp := newFwmSgRuleCreateResponse(raw)
	r.Create(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, fwmSgRuleGroupId, *gotRequest.GroupId)
	assert.Len(t, gotRequest.Rules, 1)
	assert.Equal(t, "10.0.0.0/16", *gotRequest.Rules[0].SourceContent)
	assert.Equal(t, "10.0.1.0/24", *gotRequest.Rules[0].DestContent)
	assert.Equal(t, "TCP", *gotRequest.Rules[0].Protocol)
	assert.Equal(t, "80", *gotRequest.Rules[0].Port)
	assert.Equal(t, "2", *gotRequest.Rules[0].RuleAction)
	assert.Equal(t, int64(1), *gotRequest.Rules[0].OrderIndex)

	var state fwm.FwmSecurityGroupRuleResourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "failed to read state: %v", diags)
	assert.Equal(t, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId), state.Id.ValueString())
	assert.Len(t, state.Rule, 1)
	assert.Equal(t, fwmSgRuleRuleId, state.Rule[0].RuleId.ValueString())
	assert.Equal(t, "80", state.Rule[0].DstPort.ValueString())
}

// TestFwmSecurityGroupRule_Create_EmptyRuleIds covers the case where the cloud
// API returns no rule ID.
func TestFwmSecurityGroupRule_Create_EmptyRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		result := &fwmv20250611.CreateSecurityGroupRuleResponse{}
		result.Response = &fwmv20250611.CreateSecurityGroupRuleResponseParams{
			RequestId: ptrStringFwmSgRule("fake-request-id"),
			RuleIds:   []*string{},
		}
		return result, nil
	})

	r := setFwmSgRuleClient(t, client)

	raw := newFwmSgRulePlan(fwmSgRuleGroupId, fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}), nil)
	req := newFwmSgRuleCreateRequest(t, raw)
	resp := newFwmSgRuleCreateResponse(raw)
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error creating FWM security group rule")
}

// TestFwmSecurityGroupRule_Create_MissingGroupId covers the guard that rejects
// a plan without group_id before any cloud API call.
func TestFwmSecurityGroupRule_Create_MissingGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	apiCalled := false
	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		apiCalled = true
		return &fwmv20250611.CreateSecurityGroupRuleResponse{}, nil
	})

	r := setFwmSgRuleClient(t, client)

	raw := newFwmSgRulePlan("", fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}), nil)
	req := newFwmSgRuleCreateRequest(t, raw)
	resp := newFwmSgRuleCreateResponse(raw)
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing group_id")
	assert.False(t, apiCalled)
}

// TestFwmSecurityGroupRule_Create_InvalidRuleCount covers the guard that only
// one rule may be managed per resource instance.
func TestFwmSecurityGroupRule_Create_InvalidRuleCount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	apiCalled := false
	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		apiCalled = true
		return &fwmv20250611.CreateSecurityGroupRuleResponse{}, nil
	})

	r := setFwmSgRuleClient(t, client)

	rules := []tftypes.Value{
		fwmSgRuleRawRule(),
		fwmSgRuleRawRule(),
	}
	raw := newFwmSgRulePlan(fwmSgRuleGroupId, fwmSgRuleRawList(rules), nil)
	req := newFwmSgRuleCreateRequest(t, raw)
	resp := newFwmSgRuleCreateResponse(raw)
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid rules")
	assert.False(t, apiCalled)
}

// TestFwmSecurityGroupRule_Create_APIError covers cloud API failure handling.
func TestFwmSecurityGroupRule_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=invalid rule")
	})

	r := setFwmSgRuleClient(t, client)

	raw := newFwmSgRulePlan(fwmSgRuleGroupId, fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}), nil)
	req := newFwmSgRuleCreateRequest(t, raw)
	resp := newFwmSgRuleCreateResponse(raw)
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error creating FWM security group rule")
	assert.Regexp(t, regexp.MustCompile("InvalidParameter"), resp.Diagnostics.Errors()[0].Detail())
}

// TestFwmSecurityGroupRule_Read_Success covers a normal refresh.
func TestFwmSecurityGroupRule_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		assert.Equal(t, fwmSgRuleGroupId, *request.GroupId)
		assert.Equal(t, fwmSgRuleRuleId, *request.RuleId)
		return fwmSgRuleDescribeResponse(), nil
	})

	r := setFwmSgRuleClient(t, client)

	schemaResp := fwmSgRuleResourceSchema(t)
	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId)),
		"group_id": tftypes.NewValue(tftypes.String, fwmSgRuleGroupId),
		"rules":    fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}),
		"rule":     tftypes.NewValue(fwmSgRuleAttrTypes()["rule"], nil),
	})

	req := resource.ReadRequest{}
	req.State = tfsdk.State{Raw: raw, Schema: schemaResp.Schema}
	resp := &resource.ReadResponse{}
	resp.State = tfsdk.State{Raw: raw, Schema: schemaResp.Schema}

	r.Read(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	var state fwm.FwmSecurityGroupRuleResourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "failed to read state: %v", diags)
	assert.Equal(t, fwmSgRuleGroupId, state.GroupId.ValueString())
	assert.Len(t, state.Rule, 1)
	assert.Equal(t, fwmSgRuleRuleId, state.Rule[0].RuleId.ValueString())
	assert.Equal(t, "TCP", state.Rule[0].Protocol.ValueString())
	assert.Equal(t, "2", state.Rule[0].RuleAction.ValueString())
}

// TestFwmSecurityGroupRule_Read_NotFound covers a rule removed on the cloud side.
func TestFwmSecurityGroupRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		result := &fwmv20250611.DescribeSecurityGroupRuleResponse{}
		result.Response = &fwmv20250611.DescribeSecurityGroupRuleResponseParams{
			RequestId: ptrStringFwmSgRule("fake-request-id"),
			Rule:      nil,
		}
		return result, nil
	})

	r := setFwmSgRuleClient(t, client)

	schemaResp := fwmSgRuleResourceSchema(t)
	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId)),
		"group_id": tftypes.NewValue(tftypes.String, fwmSgRuleGroupId),
		"rules":    fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}),
		"rule":     tftypes.NewValue(fwmSgRuleAttrTypes()["rule"], nil),
	})

	req := resource.ReadRequest{}
	req.State = tfsdk.State{Raw: raw, Schema: schemaResp.Schema}
	resp := &resource.ReadResponse{}
	resp.State = tfsdk.State{Raw: raw, Schema: schemaResp.Schema}

	r.Read(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
	assert.False(t, resp.Diagnostics.HasError())
	assert.NotEmpty(t, resp.Diagnostics.Warnings())
	assert.True(t, resp.State.Raw.IsNull(), "expected the resource to be removed from state")
}

// TestFwmSecurityGroupRule_Read_InvalidId covers a malformed resource identifier.
func TestFwmSecurityGroupRule_Read_InvalidId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	apiCalled := false
	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		apiCalled = true
		return fwmSgRuleDescribeResponse(), nil
	})

	r := setFwmSgRuleClient(t, client)

	schemaResp := fwmSgRuleResourceSchema(t)
	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, "invalid-id-without-separator"),
		"group_id": tftypes.NewValue(tftypes.String, nil),
		"rules":    tftypes.NewValue(fwmSgRuleAttrTypes()["rules"], nil),
		"rule":     tftypes.NewValue(fwmSgRuleAttrTypes()["rule"], nil),
	})

	req := resource.ReadRequest{}
	req.State = tfsdk.State{Raw: raw, Schema: schemaResp.Schema}
	resp := &resource.ReadResponse{}
	resp.State = tfsdk.State{Raw: raw, Schema: schemaResp.Schema}

	r.Read(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid resource identifier")
	assert.False(t, apiCalled)
}

// TestFwmSecurityGroupRule_Update_Success covers an in-place update, where the
// rule ID is taken from the resource identifier when it is not configured.
func TestFwmSecurityGroupRule_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	var gotRequest *fwmv20250611.ModifySecurityGroupRuleRequest
	patches.ApplyMethodFunc(fwmClient, "ModifySecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifySecurityGroupRuleRequest) (*fwmv20250611.ModifySecurityGroupRuleResponse, error) {
		gotRequest = request
		return &fwmv20250611.ModifySecurityGroupRuleResponse{
			Response: &fwmv20250611.ModifySecurityGroupRuleResponseParams{
				RequestId: ptrStringFwmSgRule("fake-request-id"),
			},
		}, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		return fwmSgRuleDescribeResponse(), nil
	})

	r := setFwmSgRuleClient(t, client)

	schemaResp := fwmSgRuleResourceSchema(t)
	stateRaw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId)),
		"group_id": tftypes.NewValue(tftypes.String, fwmSgRuleGroupId),
		"rules":    fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}),
		"rule":     fwmSgRuleRawDetailList(),
	})
	planRaw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId)),
		"group_id": tftypes.NewValue(tftypes.String, fwmSgRuleGroupId),
		"rules":    fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}),
		"rule":     fwmSgRuleRawDetailList(),
	})

	req := resource.UpdateRequest{}
	req.State = tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema}
	req.Plan = tfsdk.Plan{Raw: planRaw, Schema: schemaResp.Schema}
	resp := &resource.UpdateResponse{}
	resp.State = tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema}

	r.Update(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, fwmSgRuleGroupId, *gotRequest.GroupId)
	assert.Equal(t, fwmSgRuleRuleId, *gotRequest.Rule.RuleId)
	assert.Equal(t, "TCP", *gotRequest.Rule.Protocol)
}

// TestFwmSecurityGroupRule_Update_APIError covers a cloud API failure during update.
func TestFwmSecurityGroupRule_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifySecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifySecurityGroupRuleRequest) (*fwmv20250611.ModifySecurityGroupRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound.RuleNotFound, Message=rule not found")
	})

	r := setFwmSgRuleClient(t, client)

	schemaResp := fwmSgRuleResourceSchema(t)
	stateRaw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId)),
		"group_id": tftypes.NewValue(tftypes.String, fwmSgRuleGroupId),
		"rules":    fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}),
		"rule":     fwmSgRuleRawDetailList(),
	})
	planRaw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId)),
		"group_id": tftypes.NewValue(tftypes.String, fwmSgRuleGroupId),
		"rules":    fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}),
		"rule":     fwmSgRuleRawDetailList(),
	})

	req := resource.UpdateRequest{}
	req.State = tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema}
	req.Plan = tfsdk.Plan{Raw: planRaw, Schema: schemaResp.Schema}
	resp := &resource.UpdateResponse{}
	resp.State = tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema}

	r.Update(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error updating FWM security group rule")
}

// TestFwmSecurityGroupRule_Delete_Success covers a normal delete.
func TestFwmSecurityGroupRule_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	var gotRequest *fwmv20250611.DeleteSecurityGroupRuleRequest
	patches.ApplyMethodFunc(fwmClient, "DeleteSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteSecurityGroupRuleRequest) (*fwmv20250611.DeleteSecurityGroupRuleResponse, error) {
		gotRequest = request
		return &fwmv20250611.DeleteSecurityGroupRuleResponse{
			Response: &fwmv20250611.DeleteSecurityGroupRuleResponseParams{
				RequestId: ptrStringFwmSgRule("fake-request-id"),
			},
		}, nil
	})

	r := setFwmSgRuleClient(t, client)

	schemaResp := fwmSgRuleResourceSchema(t)
	stateRaw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId)),
		"group_id": tftypes.NewValue(tftypes.String, fwmSgRuleGroupId),
		"rules":    fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}),
		"rule":     fwmSgRuleRawDetailList(),
	})

	req := resource.DeleteRequest{}
	req.State = tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema}
	resp := &resource.DeleteResponse{}
	resp.State = tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema}

	r.Delete(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, fwmSgRuleGroupId, *gotRequest.GroupId)
	assert.Len(t, gotRequest.RuleIds, 1)
	assert.Equal(t, fwmSgRuleRuleId, *gotRequest.RuleIds[0])
}

// TestFwmSecurityGroupRule_Delete_APIError covers a cloud API failure during delete.
func TestFwmSecurityGroupRule_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteSecurityGroupRuleRequest) (*fwmv20250611.DeleteSecurityGroupRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=delete failed")
	})

	r := setFwmSgRuleClient(t, client)

	schemaResp := fwmSgRuleResourceSchema(t)
	stateRaw := tftypes.NewValue(tftypes.Object{AttributeTypes: fwmSgRuleAttrTypes()}, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, fmt.Sprintf("%s#%s", fwmSgRuleGroupId, fwmSgRuleRuleId)),
		"group_id": tftypes.NewValue(tftypes.String, fwmSgRuleGroupId),
		"rules":    fwmSgRuleRawList([]tftypes.Value{fwmSgRuleRawRule()}),
		"rule":     fwmSgRuleRawDetailList(),
	})

	req := resource.DeleteRequest{}
	req.State = tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema}
	resp := &resource.DeleteResponse{}
	resp.State = tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema}

	r.Delete(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error deleting FWM security group rule")
}

// TestFwmSecurityGroupRule_Metadata covers the registered type name.
func TestFwmSecurityGroupRule_Metadata(t *testing.T) {
	r := fwm.NewFwmSecurityGroupRuleResource()
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	assert.Equal(t, "tencentcloud_fwm_security_group_rule", resp.TypeName)
}

// TestFwmSecurityGroupRule_Schema covers the expected attribute shape.
func TestFwmSecurityGroupRule_Schema(t *testing.T) {
	schemaResp := fwmSgRuleResourceSchema(t)

	rulesAttr, ok := schemaResp.Schema.Attributes["rules"]
	assert.True(t, ok, "rules attribute is missing")
	assert.True(t, rulesAttr.IsOptional())

	ruleAttr, ok := schemaResp.Schema.Attributes["rule"]
	assert.True(t, ok, "rule attribute is missing")
	assert.True(t, ruleAttr.IsComputed())

	assert.NotNil(t, rulesAttr.GetType())
	assert.NotNil(t, ruleAttr.GetType())
}
