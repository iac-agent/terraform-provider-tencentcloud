package fwm_test

import (
	"context"
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

func ptrStringFSGR(s string) *string {
	return &s
}

func ptrInt64FSGR(i int64) *int64 {
	return &i
}

// buildFwmSecurityGroupRuleSchemaType builds the tftypes.Object describing the
// tencentcloud_fwm_security_group_rule schema, used to construct tfsdk.State /
// tfsdk.Plan values for the framework unit tests.
func buildFwmSecurityGroupRuleSchemaType() tftypes.Object {
	str := tftypes.String
	i64 := tftypes.Number

	ruleDetailAttr := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"rule_id":             str,
		"sequence":            i64,
		"region":              str,
		"ip_version":          str,
		"src_content":         str,
		"src_type":            str,
		"dst_content":         str,
		"dst_type":            str,
		"protocol":            str,
		"protocol_port_type":  i64,
		"service_template_id": str,
		"dst_port":            str,
		"rule_action":         str,
		"detail":              str,
		"rule_source":         str,
		"scope":               str,
	}}

	rulesAttr := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"ip_version":          str,
		"source_content":      str,
		"source_type":         str,
		"dest_content":        str,
		"dest_type":           str,
		"protocol":            str,
		"port":                str,
		"service_template_id": str,
		"rule_action":         str,
		"description":         str,
		"order_index":         i64,
		"rule_id":             str,
		"scope":               str,
		"protocol_port_type":  i64,
		"belong_member_id":    str,
	}}

	return tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":       str,
		"group_id": str,
		"rules":    tftypes.List{ElementType: rulesAttr},
		"rule":     tftypes.List{ElementType: ruleDetailAttr},
	}}
}

func newFwmSecurityGroupRuleResourceWithClient(client *connectivity.TencentCloudClient) resource.Resource {
	res := fwm.NewFwmSecurityGroupRuleResource()
	res.(resource.ResourceWithConfigure).Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}},
		&resource.ConfigureResponse{},
	)
	return res
}

func newFwmSecurityGroupRuleState() tfsdk.State {
	return tfsdk.State{
		Raw:    tftypes.NewValue(buildFwmSecurityGroupRuleSchemaType(), nil),
		Schema: getFwmSecurityGroupRuleSchema(),
	}
}

func newFwmSecurityGroupRulePlan() tfsdk.Plan {
	return tfsdk.Plan{
		Raw:    tftypes.NewValue(buildFwmSecurityGroupRuleSchemaType(), nil),
		Schema: getFwmSecurityGroupRuleSchema(),
	}
}

func getFwmSecurityGroupRuleSchema() tfsdk.Schema {
	var schemaResp resource.SchemaResponse
	fwm.NewFwmSecurityGroupRuleResource().Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	return schemaResp.Schema
}

// buildFwmSecurityGroupRuleValue builds a raw object value for the resource schema.
func buildFwmSecurityGroupRuleValue(
	id interface{},
	groupId interface{},
	rules []map[string]interface{},
	details []map[string]interface{},
) tftypes.Value {
	return tftypes.NewValue(buildFwmSecurityGroupRuleSchemaType(), map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, id),
		"group_id": tftypes.NewValue(tftypes.String, groupId),
		"rules":    buildFwmSecurityGroupRuleListValue(rules, fwmSecurityGroupRuleRulesAttributes()),
		"rule":     buildFwmSecurityGroupRuleListValue(details, fwmSecurityGroupRuleDetailAttributes()),
	})
}

func buildFwmSecurityGroupRuleListValue(items []map[string]interface{}, attributeTypes map[string]tftypes.Type) tftypes.Value {
	elementType := tftypes.Object{AttributeTypes: attributeTypes}
	values := make([]tftypes.Value, 0, len(items))
	for _, item := range items {
		attrValues := make(map[string]tftypes.Value, len(attributeTypes))
		for name, attrType := range attributeTypes {
			var raw interface{}
			if v, ok := item[name]; ok {
				raw = v
			}
			attrValues[name] = tftypes.NewValue(attrType, raw)
		}
		values = append(values, tftypes.NewValue(elementType, attrValues))
	}
	return tftypes.NewValue(tftypes.List{ElementType: elementType}, values)
}

func fwmSecurityGroupRuleRulesAttributes() map[string]tftypes.Type {
	out := make(map[string]tftypes.Type)
	rulesType := buildFwmSecurityGroupRuleSchemaType().AttributeTypes["rules"].(tftypes.List)
	objType := rulesType.ElementType.(tftypes.Object)
	for name, attrType := range objType.AttributeTypes {
		out[name] = attrType
	}
	return out
}

func fwmSecurityGroupRuleDetailAttributes() map[string]tftypes.Type {
	out := make(map[string]tftypes.Type)
	detailType := buildFwmSecurityGroupRuleSchemaType().AttributeTypes["rule"].(tftypes.List)
	objType := detailType.ElementType.(tftypes.Object)
	for name, attrType := range objType.AttributeTypes {
		out[name] = attrType
	}
	return out
}

func fwmMockRuleDetail(ruleId string) *fwmv20250611.SgRuleResp {
	return &fwmv20250611.SgRuleResp{
		RuleId:            ptrStringFSGR(ruleId),
		Sequence:          ptrInt64FSGR(1),
		Region:            ptrStringFSGR("ap-guangzhou"),
		IpVersion:         ptrStringFSGR("v4"),
		SrcContent:        ptrStringFSGR("10.0.0.0/8"),
		SrcType:           ptrStringFSGR("net"),
		DstContent:        ptrStringFSGR("192.168.0.0/16"),
		DstType:           ptrStringFSGR("net"),
		Protocol:          ptrStringFSGR("TCP"),
		ProtocolPortType:  ptrInt64FSGR(0),
		ServiceTemplateId: ptrStringFSGR(""),
		DstPort:           ptrStringFSGR("80"),
		RuleAction:        ptrStringFSGR("2"),
		Detail:            ptrStringFSGR("allow http"),
		RuleSource:        ptrStringFSGR("0"),
		Scope:             ptrStringFSGR("1"),
	}
}

func TestFwmSecurityGroupRule_Metadata(t *testing.T) {
	var resp resource.MetadataResponse
	fwm.NewFwmSecurityGroupRuleResource().Metadata(context.Background(), resource.MetadataRequest{}, &resp)
	assert.Equal(t, "tencentcloud_fwm_security_group_rule", resp.TypeName)
}

func TestFwmSecurityGroupRule_Schema_Attributes(t *testing.T) {
	s := getFwmSecurityGroupRuleSchema()

	idAttr, ok := s.Attributes["id"]
	assert.True(t, ok)
	assert.True(t, idAttr.IsComputed())

	groupIdAttr, ok := s.Attributes["group_id"]
	assert.True(t, ok)
	assert.True(t, groupIdAttr.IsOptional())
	assert.True(t, groupIdAttr.IsComputed())

	rulesAttr, ok := s.Attributes["rules"]
	assert.True(t, ok)
	assert.True(t, rulesAttr.IsOptional())

	ruleAttr, ok := s.Attributes["rule"]
	assert.True(t, ok)
	assert.True(t, ruleAttr.IsComputed())
}

func TestFwmSecurityGroupRule_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwg-1a2b3c4d", *request.GroupId)
		assert.Len(t, request.Rules, 1)
		assert.Equal(t, "TCP", *request.Rules[0].Protocol)

		resp := fwmv20250611.NewCreateSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.CreateSecurityGroupRuleResponseParams{
			RequestId: ptrStringFSGR("fake-request-id-create"),
			RuleIds:   []*string{ptrStringFSGR("fwr-1a2b3c4d")},
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwg-1a2b3c4d", *request.GroupId)
		resp := fwmv20250611.NewDescribeSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.DescribeSecurityGroupRuleResponseParams{
			RequestId: ptrStringFSGR("fake-request-id-describe"),
			Rule:      fwmMockRuleDetail(*request.RuleId),
		}
		return resp, nil
	})

	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	plan := newFwmSecurityGroupRulePlan()
	plan.Raw = buildFwmSecurityGroupRuleValue(nil, "fwg-1a2b3c4d", []map[string]interface{}{
		{
			"ip_version":         "v4",
			"source_type":        "net",
			"source_content":     "10.0.0.0/8",
			"dest_type":          "net",
			"dest_content":       "192.168.0.0/16",
			"protocol":           "TCP",
			"port":               "80",
			"rule_action":        "2",
			"description":        "allow http",
			"scope":              "1",
			"order_index":        float64(1),
			"protocol_port_type": float64(0),
		},
	}, nil)

	createResp := resource.CreateResponse{State: newFwmSecurityGroupRuleState()}
	res.Create(context.Background(), resource.CreateRequest{Plan: plan}, &createResp)

	assert.False(t, createResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Create_EmptyRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateSecurityGroupRuleRequest) (*fwmv20250611.CreateSecurityGroupRuleResponse, error) {
		resp := fwmv20250611.NewCreateSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.CreateSecurityGroupRuleResponseParams{
			RequestId: ptrStringFSGR("fake-request-id-create"),
			RuleIds:   []*string{},
		}
		return resp, nil
	})

	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	plan := newFwmSecurityGroupRulePlan()
	plan.Raw = buildFwmSecurityGroupRuleValue(nil, "fwg-1a2b3c4d", []map[string]interface{}{
		{"ip_version": "v4", "protocol": "TCP", "port": "80", "rule_action": "2"},
	}, nil)

	createResp := resource.CreateResponse{State: newFwmSecurityGroupRuleState()}
	res.Create(context.Background(), resource.CreateRequest{Plan: plan}, &createResp)

	assert.True(t, createResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Create_MissingGroupId(t *testing.T) {
	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	plan := newFwmSecurityGroupRulePlan()
	plan.Raw = buildFwmSecurityGroupRuleValue(nil, nil, []map[string]interface{}{
		{"ip_version": "v4", "protocol": "TCP", "port": "80", "rule_action": "2"},
	}, nil)

	createResp := resource.CreateResponse{State: newFwmSecurityGroupRuleState()}
	res.Create(context.Background(), resource.CreateRequest{Plan: plan}, &createResp)

	assert.True(t, createResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Create_MissingRules(t *testing.T) {
	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	plan := newFwmSecurityGroupRulePlan()
	plan.Raw = buildFwmSecurityGroupRuleValue(nil, "fwg-1a2b3c4d", nil, nil)

	createResp := resource.CreateResponse{State: newFwmSecurityGroupRuleState()}
	res.Create(context.Background(), resource.CreateRequest{Plan: plan}, &createResp)

	assert.True(t, createResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		resp := fwmv20250611.NewDescribeSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.DescribeSecurityGroupRuleResponseParams{
			RequestId: ptrStringFSGR("fake-request-id-describe"),
			Rule:      fwmMockRuleDetail(*request.RuleId),
		}
		return resp, nil
	})

	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	state := newFwmSecurityGroupRuleState()
	state.Raw = buildFwmSecurityGroupRuleValue("fwg-1a2b3c4d#fwr-1a2b3c4d", "fwg-1a2b3c4d", []map[string]interface{}{
		{
			"ip_version":         "v4",
			"source_type":        "net",
			"source_content":     "10.0.0.0/8",
			"dest_type":          "net",
			"dest_content":       "192.168.0.0/16",
			"protocol":           "TCP",
			"port":               "80",
			"rule_action":        "2",
			"description":        "allow http",
			"scope":              "1",
			"rule_id":            "fwr-1a2b3c4d",
			"order_index":        float64(1),
			"protocol_port_type": float64(0),
		},
	}, nil)

	readResp := resource.ReadResponse{State: state}
	res.Read(context.Background(), resource.ReadRequest{State: state}, &readResp)

	assert.False(t, readResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		resp := fwmv20250611.NewDescribeSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.DescribeSecurityGroupRuleResponseParams{
			RequestId: ptrStringFSGR("fake-request-id-describe"),
			Rule:      nil,
		}
		return resp, nil
	})

	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	state := newFwmSecurityGroupRuleState()
	state.Raw = buildFwmSecurityGroupRuleValue("fwg-1a2b3c4d#fwr-notexist", "fwg-1a2b3c4d", nil, nil)

	readResp := resource.ReadResponse{State: state}
	res.Read(context.Background(), resource.ReadRequest{State: state}, &readResp)

	assert.False(t, readResp.Diagnostics.HasError())
	assert.True(t, readResp.State.Raw.IsNull())
}

func TestFwmSecurityGroupRule_Read_InvalidId(t *testing.T) {
	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	state := newFwmSecurityGroupRuleState()
	state.Raw = buildFwmSecurityGroupRuleValue("broken-id", "fwg-1a2b3c4d", nil, nil)

	readResp := resource.ReadResponse{State: state}
	res.Read(context.Background(), resource.ReadRequest{State: state}, &readResp)

	assert.True(t, readResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifySecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifySecurityGroupRuleRequest) (*fwmv20250611.ModifySecurityGroupRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwg-1a2b3c4d", *request.GroupId)
		assert.NotNil(t, request.Rule)
		assert.Equal(t, "fwr-1a2b3c4d", *request.Rule.RuleId)

		resp := fwmv20250611.NewModifySecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.ModifySecurityGroupRuleResponseParams{
			RequestId: ptrStringFSGR("fake-request-id-modify"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DescribeSecurityGroupRuleRequest) (*fwmv20250611.DescribeSecurityGroupRuleResponse, error) {
		resp := fwmv20250611.NewDescribeSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.DescribeSecurityGroupRuleResponseParams{
			RequestId: ptrStringFSGR("fake-request-id-describe"),
			Rule:      fwmMockRuleDetail(*request.RuleId),
		}
		return resp, nil
	})

	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	state := newFwmSecurityGroupRuleState()
	state.Raw = buildFwmSecurityGroupRuleValue("fwg-1a2b3c4d#fwr-1a2b3c4d", "fwg-1a2b3c4d", []map[string]interface{}{
		{
			"ip_version":  "v4",
			"protocol":    "TCP",
			"port":        "80",
			"rule_action": "2",
			"rule_id":     "fwr-1a2b3c4d",
			"scope":       "1",
		},
	}, nil)

	plan := newFwmSecurityGroupRulePlan()
	plan.Raw = buildFwmSecurityGroupRuleValue("fwg-1a2b3c4d#fwr-1a2b3c4d", "fwg-1a2b3c4d", []map[string]interface{}{
		{
			"ip_version":  "v4",
			"protocol":    "TCP",
			"port":        "8080",
			"rule_action": "2",
			"rule_id":     "fwr-1a2b3c4d",
			"scope":       "1",
		},
	}, nil)

	updateResp := resource.UpdateResponse{State: state}
	res.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, &updateResp)

	assert.False(t, updateResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Update_GroupIdChanged(t *testing.T) {
	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	state := newFwmSecurityGroupRuleState()
	state.Raw = buildFwmSecurityGroupRuleValue("fwg-1a2b3c4d#fwr-1a2b3c4d", "fwg-1a2b3c4d", []map[string]interface{}{
		{"ip_version": "v4", "protocol": "TCP", "port": "80", "rule_action": "2", "rule_id": "fwr-1a2b3c4d"},
	}, nil)

	plan := newFwmSecurityGroupRulePlan()
	plan.Raw = buildFwmSecurityGroupRuleValue("fwg-1a2b3c4d#fwr-1a2b3c4d", "fwg-other", []map[string]interface{}{
		{"ip_version": "v4", "protocol": "TCP", "port": "80", "rule_action": "2", "rule_id": "fwr-1a2b3c4d"},
	}, nil)

	updateResp := resource.UpdateResponse{State: state}
	res.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, &updateResp)

	assert.True(t, updateResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteSecurityGroupRuleRequest) (*fwmv20250611.DeleteSecurityGroupRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwg-1a2b3c4d", *request.GroupId)
		assert.Len(t, request.RuleIds, 2)
		assert.Equal(t, "fwr-1a2b3c4d", *request.RuleIds[0])
		assert.Equal(t, "fwr-5e6f7g8h", *request.RuleIds[1])

		resp := fwmv20250611.NewDeleteSecurityGroupRuleResponse()
		resp.Response = &fwmv20250611.DeleteSecurityGroupRuleResponseParams{
			RequestId: ptrStringFSGR("fake-request-id-delete"),
		}
		return resp, nil
	})

	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	state := newFwmSecurityGroupRuleState()
	state.Raw = buildFwmSecurityGroupRuleValue("fwg-1a2b3c4d#fwr-1a2b3c4d#fwr-5e6f7g8h", "fwg-1a2b3c4d", nil, nil)

	deleteResp := resource.DeleteResponse{State: state}
	res.Delete(context.Background(), resource.DeleteRequest{State: state}, &deleteResp)

	assert.False(t, deleteResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Delete_Error(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteSecurityGroupRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteSecurityGroupRuleRequest) (*fwmv20250611.DeleteSecurityGroupRuleResponse, error) {
		return nil, assert.AnError
	})

	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	state := newFwmSecurityGroupRuleState()
	state.Raw = buildFwmSecurityGroupRuleValue("fwg-1a2b3c4d#fwr-1a2b3c4d", "fwg-1a2b3c4d", nil, nil)

	deleteResp := resource.DeleteResponse{State: state}
	res.Delete(context.Background(), resource.DeleteRequest{State: state}, &deleteResp)

	assert.True(t, deleteResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_Delete_InvalidId(t *testing.T) {
	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	state := newFwmSecurityGroupRuleState()
	state.Raw = buildFwmSecurityGroupRuleValue("broken-id", "fwg-1a2b3c4d", nil, nil)

	deleteResp := resource.DeleteResponse{State: state}
	res.Delete(context.Background(), resource.DeleteRequest{State: state}, &deleteResp)

	assert.True(t, deleteResp.Diagnostics.HasError())
}

func TestFwmSecurityGroupRule_ImportState(t *testing.T) {
	res := newFwmSecurityGroupRuleResourceWithClient(&connectivity.TencentCloudClient{})

	importResp := resource.ImportStateResponse{State: newFwmSecurityGroupRuleState()}
	res.(resource.ResourceWithImportState).ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: "fwg-1a2b3c4d#fwr-1a2b3c4d"},
		&importResp,
	)
	assert.False(t, importResp.Diagnostics.HasError())
}
