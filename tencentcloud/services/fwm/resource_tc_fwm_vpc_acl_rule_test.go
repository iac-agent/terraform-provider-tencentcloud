package fwm

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

func fwmStrPtr(s string) *string { return &s }

func fwmI64Ptr(i int64) *int64 { return &i }

// newFwmVpcAclRuleResourceForTest builds a resource whose cloud API calls are
// all mocked through gomonkey, so the bare client is never actually used.
func newFwmVpcAclRuleResourceForTest() *FwmVpcAclRuleResource {
	r := NewFwmVpcAclRuleResource().(*FwmVpcAclRuleResource)
	r.ResourceWithConfigure = fw.ResourceWithConfigure{}
	return r
}

func TestFwmVpcAclRuleResource_Metadata(t *testing.T) {
	r := NewFwmVpcAclRuleResource()
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	if resp.TypeName != "tencentcloud_fwm_vpc_acl_rule" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_fwm_vpc_acl_rule")
	}
}

func TestFwmVpcAclRuleResource_Schema(t *testing.T) {
	r := NewFwmVpcAclRuleResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	cases := map[string]struct {
		required, optional, computed bool
	}{
		"id":              {computed: true},
		"group_id":        {required: true},
		"fw_group_id":     {required: true},
		"rule_action":     {required: true},
		"description":     {optional: true, computed: true},
		"order_index":     {optional: true, computed: true},
		"rule_id":         {computed: true},
		"rule_ids":        {computed: true},
		"all_total_count": {computed: true},
	}
	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q is missing from schema", name)
			continue
		}
		if attr.IsRequired() != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.IsRequired(), want.required)
		}
		if attr.IsOptional() != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, attr.IsOptional(), want.optional)
		}
		if attr.IsComputed() != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, attr.IsComputed(), want.computed)
		}
	}

	rules, ok := resp.Schema.Attributes["rules"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute rules is %T, want schema.ListNestedAttribute", resp.Schema.Attributes["rules"])
	}
	if !rules.Computed {
		t.Error("attribute rules must be Computed")
	}
	ruleAttrs := rules.NestedObject.Attributes
	expectedRuleAttrs := []string{
		"source_content", "source_type", "source_name", "dest_content", "dest_type", "dest_name",
		"protocol", "port", "param_template_id", "param_template_name", "rule_action", "order_index",
		"edge_id", "fw_group_id", "fw_group_name", "rule_id", "description", "belong_member_id",
		"belong_member", "sequence", "ip_version", "invalid", "create_time", "update_time",
		"dns_parse_cnt", "hit_count", "last_hit_time",
	}
	for _, name := range expectedRuleAttrs {
		if _, ok := ruleAttrs[name]; !ok {
			t.Errorf("nested attribute rules.%s is missing", name)
		}
	}

	belongMember, ok := ruleAttrs["belong_member"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("rules.belong_member is %T, want schema.SingleNestedAttribute", ruleAttrs["belong_member"])
	}
	for _, name := range []string{"app_id", "uin", "nickname", "member_id", "node_name"} {
		if _, ok := belongMember.Attributes[name]; !ok {
			t.Errorf("nested attribute rules.belong_member.%s is missing", name)
		}
	}
}

func TestSplitFwmVpcAclRuleId(t *testing.T) {
	group, rule, ok := splitFwmVpcAclRuleId("acg-1" + tccommon.FILED_SP + "rule-1")
	if !ok || group != "acg-1" || rule != "rule-1" {
		t.Errorf("splitFwmVpcAclRuleId = (%q, %q, %v), want (\"acg-1\", \"rule-1\", true)", group, rule, ok)
	}

	for _, in := range []string{
		"",
		"only-group",
		tccommon.FILED_SP + "rule-1",
		"acg-1" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	} {
		if _, _, ok := splitFwmVpcAclRuleId(in); ok {
			t.Errorf("splitFwmVpcAclRuleId(%q) ok=true, want false", in)
		}
	}
}

func TestFwmVpcAclRuleBelongMemberAttrTypes(t *testing.T) {
	attrTypes := fwmVpcAclRuleBelongMemberAttrTypes()
	for _, name := range []string{"app_id", "uin", "nickname", "member_id", "node_name"} {
		if _, ok := attrTypes[name]; !ok {
			t.Errorf("belong_member attribute %q is missing", name)
		}
	}
	if len(attrTypes) != 5 {
		t.Errorf("belong_member attribute count = %d, want 5", len(attrTypes))
	}
}

func TestFwmVpcAclRuleRuleObjectType(t *testing.T) {
	objType, ok := fwmVpcAclRuleRuleObjectType().(attr.TypeWithAttributeTypes)
	if !ok {
		t.Fatalf("fwmVpcAclRuleRuleObjectType returned %T, want attr.TypeWithAttributeTypes", fwmVpcAclRuleRuleObjectType())
	}
	attrTypes := objType.AttributeTypes()
	if len(attrTypes) != 27 {
		t.Errorf("rules element attribute count = %d, want 27", len(attrTypes))
	}

	_, nestedOk := attrTypes["belong_member"].(attr.TypeWithAttributeTypes)
	if !nestedOk {
		t.Errorf("rules.belong_member type = %T, want attr.TypeWithAttributeTypes", attrTypes["belong_member"])
	}
}

func TestFwmVpcAclRuleResource_ConfigureNilProviderData(t *testing.T) {
	r := NewFwmVpcAclRuleResource()
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

func TestFwmVpcAclRuleResource_ImportState(t *testing.T) {
	r := NewFwmVpcAclRuleResource()

	resp := &resource.ImportStateResponse{}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "acg-1" + tccommon.FILED_SP + "rule-1"}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState returned diagnostics: %v", resp.Diagnostics)
	}

	resp = &resource.ImportStateResponse{}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "broken-id"}, resp)
	if !resp.Diagnostics.HasError() {
		t.Error("ImportState with a broken id should return a diagnostic error")
	}
}

func TestFwmVpcAclRuleResource_ModifyPlan(t *testing.T) {
	r := NewFwmVpcAclRuleResource()

	plan := FwmVpcAclRuleResourceModel{
		Id:          types.StringValue("acg-1" + tccommon.FILED_SP + "rule-1"),
		GroupId:     types.StringValue("acg-1"),
		FwGroupId:   types.StringValue("cfwg-1"),
		RuleAction:  types.StringValue("accept"),
		Description: types.StringValue("updated"),
		OrderIndex:  types.Int64Value(1),
		RuleId:      types.StringValue("rule-1"),
	}
	state := plan
	state.Description = types.StringValue("original")

	resp := &resource.ModifyPlanResponse{}
	r.checkImmutablePlanChanges(plan, state, resp)
	if !resp.Diagnostics.HasError() {
		t.Error("ModifyPlan should report an error when description changes")
	}

	resp = &resource.ModifyPlanResponse{}
	plan.Description = state.Description
	plan.OrderIndex = types.Int64Value(2)
	r.checkImmutablePlanChanges(plan, state, resp)
	if !resp.Diagnostics.HasError() {
		t.Error("ModifyPlan should report an error when order_index changes")
	}

	resp = &resource.ModifyPlanResponse{}
	plan.OrderIndex = state.OrderIndex
	r.checkImmutablePlanChanges(plan, state, resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("ModifyPlan returned unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestFwmVpcAclRuleResource_refreshModel(t *testing.T) {
	r := NewFwmVpcAclRuleResource()
	var model FwmVpcAclRuleResourceModel

	item := &fwmv20250611.VpcAclRuleResp{
		RuleId:            fwmStrPtr("rule-1"),
		Sequence:          fwmI64Ptr(3),
		IpVersion:         fwmStrPtr("ipv4"),
		SourceContent:     fwmStrPtr("10.0.0.0/8"),
		SourceType:        fwmStrPtr("ip"),
		SourceName:        fwmStrPtr("source-asset"),
		DestContent:       fwmStrPtr("0.0.0.0/0"),
		DestType:          fwmStrPtr("ip"),
		DestName:          fwmStrPtr("dest-asset"),
		Protocol:          fwmStrPtr("TCP"),
		Port:              fwmStrPtr("80"),
		ParamTemplateId:   fwmStrPtr("ppm-1"),
		ParamTemplateName: fwmStrPtr("port-template"),
		RuleAction:        fwmStrPtr("accept"),
		Description:       fwmStrPtr("desc"),
		EdgeId:            fwmStrPtr("ALL"),
		FwGroupId:         fwmStrPtr("cfwg-1"),
		Invalid:           fwmI64Ptr(0),
		BelongMember: &fwmv20250611.MemberInfo{
			AppId:    fwmStrPtr("1250000000"),
			Uin:      fwmStrPtr("100000000001"),
			Nickname: fwmStrPtr("member"),
			MemberId: fwmStrPtr("member-1"),
			NodeName: fwmStrPtr("dept"),
		},
		CreateTime:  fwmStrPtr("2026-01-01 00:00:00"),
		UpdateTime:  fwmStrPtr("2026-01-02 00:00:00"),
		DnsParseCnt: fwmI64Ptr(2),
		FwGroupName: fwmStrPtr("fw-group"),
		HitCount:    fwmI64Ptr(7),
		LastHitTime: fwmStrPtr("2026-01-03 00:00:00"),
	}

	model.FwGroupId = types.StringValue("cfwg-1")
	model.RuleAction = types.StringValue("accept")
	r.refreshModel(context.Background(), &model, item)

	if model.Description.ValueString() != "desc" {
		t.Errorf("Description = %q, want \"desc\"", model.Description.ValueString())
	}
	if model.OrderIndex.ValueInt64() != 3 {
		t.Errorf("OrderIndex = %d, want 3", model.OrderIndex.ValueInt64())
	}
	if model.FwGroupId.ValueString() != "cfwg-1" {
		t.Errorf("FwGroupId = %q, want \"cfwg-1\"", model.FwGroupId.ValueString())
	}
	if model.RuleAction.ValueString() != "accept" {
		t.Errorf("RuleAction = %q, want \"accept\"", model.RuleAction.ValueString())
	}
	if model.Rules.IsNull() || len(model.Rules.Elements()) != 1 {
		t.Fatalf("Rules should contain exactly one element")
	}

	var rules []FwmVpcAclRuleRuleModel
	if diags := model.Rules.ElementsAs(context.Background(), &rules, false); diags.HasError() {
		t.Fatalf("Rules.ElementsAs failed: %v", diags)
	}
	if rules[0].RuleId.ValueString() != "rule-1" {
		t.Errorf("rules[0].RuleId = %q, want \"rule-1\"", rules[0].RuleId.ValueString())
	}
	if rules[0].HitCount.ValueInt64() != 7 {
		t.Errorf("rules[0].HitCount = %d, want 7", rules[0].HitCount.ValueInt64())
	}
	if rules[0].BelongMember.IsNull() {
		t.Error("rules[0].BelongMember should not be null")
	}
}

func TestFwmVpcAclRuleResource_refreshModelNullMember(t *testing.T) {
	r := NewFwmVpcAclRuleResource()
	var model FwmVpcAclRuleResourceModel

	r.refreshModel(context.Background(), &model, &fwmv20250611.VpcAclRuleResp{
		RuleId: fwmStrPtr("rule-1"),
	})

	var rules []FwmVpcAclRuleRuleModel
	if diags := model.Rules.ElementsAs(context.Background(), &rules, false); diags.HasError() {
		t.Fatalf("Rules.ElementsAs failed: %v", diags)
	}
	if !rules[0].BelongMember.IsNull() {
		t.Error("rules[0].BelongMember should be null when the API omits it")
	}
}

func TestFwmVpcAclRuleResource_Read_NotFound(t *testing.T) {
	r := newFwmVpcAclRuleResourceForTest()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(r.Client(), "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeVpcAclRulesWithContext",
		func(_ context.Context, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
			allTotal := int64(0)
			resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
			resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{
				RequestId:     fwmStrPtr("request-id"),
				Rules:         []*fwmv20250611.VpcAclRuleResp{},
				AllTotalCount: &allTotal,
			}
			return resp, nil
		})

	found, err := r.read(context.Background(), "acg-1", "rule-not-exist")
	if err != nil {
		t.Fatalf("read returned unexpected error: %v", err)
	}
	if found != nil {
		t.Errorf("read should return nil when the rule does not exist, got %+v", found)
	}
}

func TestFwmVpcAclRuleResource_Read_Found(t *testing.T) {
	r := newFwmVpcAclRuleResourceForTest()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(r.Client(), "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeVpcAclRulesWithContext",
		func(_ context.Context, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
			if request.GroupId == nil || *request.GroupId != "acg-1" {
				t.Errorf("DescribeVpcAclRules GroupId = %v, want \"acg-1\"", request.GroupId)
			}
			resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
			resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{
				RequestId: fwmStrPtr("request-id"),
				Rules: []*fwmv20250611.VpcAclRuleResp{
					{RuleId: fwmStrPtr("rule-other"), Sequence: fwmI64Ptr(1)},
					{RuleId: fwmStrPtr("rule-1"), Sequence: fwmI64Ptr(2), RuleAction: fwmStrPtr("drop")},
				},
			}
			return resp, nil
		})

	found, err := r.read(context.Background(), "acg-1", "rule-1")
	if err != nil {
		t.Fatalf("read returned unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("read returned nil, want the matching rule")
	}
	if found.RuleId == nil || *found.RuleId != "rule-1" {
		t.Errorf("read returned rule id %v, want \"rule-1\"", found.RuleId)
	}
}

func TestFwmVpcAclRuleResource_Read_EmptyResponse(t *testing.T) {
	r := newFwmVpcAclRuleResourceForTest()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(r.Client(), "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeVpcAclRulesWithContext",
		func(_ context.Context, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
			return fwmv20250611.NewDescribeVpcAclRulesResponse(), nil
		})

	if _, err := r.read(context.Background(), "acg-1", "rule-1"); err == nil {
		t.Error("read should return an error when the API response body is nil")
	}
}

func TestFwmService_DeleteFwmVpcAclRule(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	var gotGroupId string
	var gotRuleIds []string
	patches.ApplyMethodFunc(fwmClient, "DeleteVpcAclRuleWithContext",
		func(_ context.Context, request *fwmv20250611.DeleteVpcAclRuleRequest) (*fwmv20250611.DeleteVpcAclRuleResponse, error) {
			if request.GroupId != nil {
				gotGroupId = *request.GroupId
			}
			for _, id := range request.RuleIds {
				if id != nil {
					gotRuleIds = append(gotRuleIds, *id)
				}
			}
			resp := fwmv20250611.NewDeleteVpcAclRuleResponse()
			resp.Response = &fwmv20250611.DeleteVpcAclRuleResponseParams{RequestId: fwmStrPtr("request-id")}
			return resp, nil
		})

	service := NewFwmService(client)
	if err := service.DeleteFwmVpcAclRule(context.Background(), "acg-1", []string{"rule-1"}); err != nil {
		t.Fatalf("DeleteFwmVpcAclRule returned unexpected error: %v", err)
	}
	if gotGroupId != "acg-1" {
		t.Errorf("DeleteVpcAclRule GroupId = %q, want \"acg-1\"", gotGroupId)
	}
	if len(gotRuleIds) != 1 || gotRuleIds[0] != "rule-1" {
		t.Errorf("DeleteVpcAclRule RuleIds = %v, want [\"rule-1\"]", gotRuleIds)
	}
}

func TestFwmService_CreateFwmVpcAclRule_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "CreateVpcAclRuleWithContext",
		func(_ context.Context, request *fwmv20250611.CreateVpcAclRuleRequest) (*fwmv20250611.CreateVpcAclRuleResponse, error) {
			return fwmv20250611.NewCreateVpcAclRuleResponse(), nil
		})

	service := NewFwmService(client)
	if _, err := service.CreateFwmVpcAclRule(context.Background(), "acg-1", []*fwmv20250611.VpcAclRule{}); err == nil {
		t.Error("CreateFwmVpcAclRule should return an error when the API response body is nil")
	}
}
