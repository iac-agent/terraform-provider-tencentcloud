package fwm

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

func fwmVpcAclRule2StrPtr(s string) *string { return &s }
func fwmVpcAclRule2Int64Ptr(v int64) *int64 { return &v }

// fwmVpcAclRule2ResourceSchema returns the schema of the resource under test.
func fwmVpcAclRule2ResourceSchema() schema.Schema {
	resp := &resource.SchemaResponse{}
	NewFwmVpcAclRule2Resource().Schema(context.Background(), resource.SchemaRequest{}, resp)
	return resp.Schema
}

// fwmVpcAclRule2ObjectType builds the tftypes.Object matching the schema, used
// to materialize plan / state values for the CRUD callbacks.
func fwmVpcAclRule2ObjectType(s schema.Schema) tftypes.Object {
	attrs := make(map[string]tftypes.Type, len(s.Attributes))
	for name, attr := range s.Attributes {
		switch attr.(type) {
		case schema.StringAttribute:
			attrs[name] = tftypes.String
		default:
			attrs[name] = tftypes.Number
		}
	}
	return tftypes.Object{AttributeTypes: attrs}
}

// fwmVpcAclRule2NullValues returns a fully-null attribute value map.
func fwmVpcAclRule2NullValues(s schema.Schema) map[string]tftypes.Value {
	attrs := make(map[string]tftypes.Value, len(s.Attributes))
	objectType := fwmVpcAclRule2ObjectType(s)
	for name := range s.Attributes {
		attrs[name] = tftypes.NewValue(objectType.AttributeTypes[name], nil)
	}
	return attrs
}

// fwmVpcAclRule2ModelValues builds the attribute value map for a full model.
func fwmVpcAclRule2ModelValues(s schema.Schema, m FwmVpcAclRule2ResourceModel) map[string]tftypes.Value {
	objectType := fwmVpcAclRule2ObjectType(s)
	str := func(v types.String) tftypes.Value {
		if v.IsNull() || v.IsUnknown() {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, v.ValueString())
	}
	num := func(v types.Int64) tftypes.Value {
		if v.IsNull() || v.IsUnknown() {
			return tftypes.NewValue(tftypes.Number, nil)
		}
		return tftypes.NewValue(tftypes.Number, v.ValueInt64())
	}

	values := fwmVpcAclRule2NullValues(s)
	values["id"] = str(m.Id)
	values["group_id"] = str(m.GroupId)
	values["rule_id"] = str(m.RuleId)
	values["source_content"] = str(m.SourceContent)
	values["source_type"] = str(m.SourceType)
	values["dest_content"] = str(m.DestContent)
	values["dest_type"] = str(m.DestType)
	values["protocol"] = str(m.Protocol)
	values["port"] = str(m.Port)
	values["rule_action"] = str(m.RuleAction)
	values["order_index"] = num(m.OrderIndex)
	values["edge_id"] = str(m.EdgeId)
	values["fw_group_id"] = str(m.FwGroupId)
	values["description"] = str(m.Description)
	values["param_template_id"] = str(m.ParamTemplateId)
	values["belong_member_id"] = str(m.BelongMemberId)
	values["ip_version"] = str(m.IpVersion)
	values["invalid"] = num(m.Invalid)
	values["param_template_name"] = str(m.ParamTemplateName)
	values["fw_group_name"] = str(m.FwGroupName)
	values["source_name"] = str(m.SourceName)
	values["dest_name"] = str(m.DestName)
	values["create_time"] = str(m.CreateTime)
	values["update_time"] = str(m.UpdateTime)
	values["dns_parse_cnt"] = num(m.DnsParseCnt)
	values["hit_count"] = num(m.HitCount)
	values["last_hit_time"] = str(m.LastHitTime)

	_ = objectType
	return values
}

// newFwmVpcAclRule2Resource wires the resource under test to a fake SDK client.
func newFwmVpcAclRule2Resource(client *connectivity.TencentCloudClient) *FwmVpcAclRule2Resource {
	r := &FwmVpcAclRule2Resource{}
	r.Configure(context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}},
		&resource.ConfigureResponse{})
	return r
}

func fwmVpcAclRule2NewConnectivityClient(patches *gomonkey.Patches, sdkClient *fwmv20250611.Client) *connectivity.TencentCloudClient {
	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", sdkClient)
	return client
}

// buildFwmVpcAclRule2Resp builds the API item returned by DescribeVpcAclRules.
func buildFwmVpcAclRule2Resp(ruleId, groupId string) *fwmv20250611.VpcAclRuleResp {
	return &fwmv20250611.VpcAclRuleResp{
		RuleId:            fwmVpcAclRule2StrPtr(ruleId),
		Sequence:          fwmVpcAclRule2Int64Ptr(5),
		IpVersion:         fwmVpcAclRule2StrPtr("ipv4"),
		SourceContent:     fwmVpcAclRule2StrPtr("10.0.0.0/8"),
		SourceType:        fwmVpcAclRule2StrPtr("ip"),
		SourceName:        fwmVpcAclRule2StrPtr("source-asset"),
		DestContent:       fwmVpcAclRule2StrPtr("192.168.0.0/16"),
		DestType:          fwmVpcAclRule2StrPtr("ip"),
		DestName:          fwmVpcAclRule2StrPtr("dest-asset"),
		Protocol:          fwmVpcAclRule2StrPtr("TCP"),
		Port:              fwmVpcAclRule2StrPtr("80,443"),
		ParamTemplateId:   fwmVpcAclRule2StrPtr("ptb-123"),
		ParamTemplateName: fwmVpcAclRule2StrPtr("port-template"),
		RuleAction:        fwmVpcAclRule2StrPtr("accept"),
		Description:       fwmVpcAclRule2StrPtr("drift-description"),
		EdgeId:            fwmVpcAclRule2StrPtr("CFWS-abc"),
		FwGroupId:         fwmVpcAclRule2StrPtr("cfwg-def"),
		Invalid:           fwmVpcAclRule2Int64Ptr(0),
		BelongMember: &fwmv20250611.MemberInfo{
			AppId:    fwmVpcAclRule2StrPtr("1300000000"),
			Uin:      fwmVpcAclRule2StrPtr("100000000001"),
			Nickname: fwmVpcAclRule2StrPtr("member-nick"),
			MemberId: fwmVpcAclRule2StrPtr("mem-0001"),
			NodeName: fwmVpcAclRule2StrPtr("root-dept"),
		},
		CreateTime:   fwmVpcAclRule2StrPtr("2025-06-11 10:00:00"),
		UpdateTime:   fwmVpcAclRule2StrPtr("2025-06-11 11:00:00"),
		DnsParseCnt:  fwmVpcAclRule2Int64Ptr(3),
		FwGroupName:  fwmVpcAclRule2StrPtr("fw-group-name"),
		HitCount:     fwmVpcAclRule2Int64Ptr(42),
		LastHitTime:  fwmVpcAclRule2StrPtr("2025-06-11 12:00:00"),
	}
}

// fwmVpcAclRule2PlanModel is the configured model used by the CRUD tests.
func fwmVpcAclRule2PlanModel() FwmVpcAclRule2ResourceModel {
	return FwmVpcAclRule2ResourceModel{
		Id:              types.StringUnknown(),
		GroupId:         types.StringValue("fwmrg-group"),
		RuleId:          types.StringUnknown(),
		SourceContent:   types.StringValue("10.1.0.0/16"),
		SourceType:      types.StringValue("ip"),
		DestContent:     types.StringValue("10.2.0.0/16"),
		DestType:        types.StringValue("ip"),
		Protocol:        types.StringValue("TCP"),
		Port:            types.StringValue("8080"),
		RuleAction:      types.StringValue("drop"),
		OrderIndex:      types.Int64Value(5),
		EdgeId:          types.StringValue("CFWS-abc"),
		FwGroupId:       types.StringValue("cfwg-def"),
		Description:     types.StringValue("managed by terraform"),
		ParamTemplateId: types.StringValue("ptb-123"),
		BelongMemberId:  types.StringValue("mem-0001"),
	}
}

// TestFwmVpcAclRule2Resource_Metadata validates the resource type name.
func TestFwmVpcAclRule2Resource_Metadata(t *testing.T) {
	resp := &resource.MetadataResponse{}
	(&FwmVpcAclRule2Resource{}).Metadata(context.Background(), resource.MetadataRequest{}, resp)

	if resp.TypeName != "tencentcloud_fwm_vpc_acl_rule2" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_fwm_vpc_acl_rule2")
	}
}

// TestFwmVpcAclRule2Resource_Schema validates the required / optional /
// computed flags of every attribute.
func TestFwmVpcAclRule2Resource_Schema(t *testing.T) {
	resp := &resource.SchemaResponse{}
	NewFwmVpcAclRule2Resource().Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required bool
		optional bool
		computed bool
	}
	strCases := map[string]expect{
		"id":                  {computed: true},
		"group_id":            {required: true},
		"rule_id":             {computed: true},
		"source_content":      {required: true},
		"source_type":         {required: true},
		"dest_content":        {required: true},
		"dest_type":           {required: true},
		"protocol":            {required: true},
		"port":                {optional: true, computed: true},
		"rule_action":         {required: true},
		"edge_id":             {required: true},
		"fw_group_id":         {required: true},
		"param_template_id":   {optional: true},
		"belong_member_id":    {optional: true},
		"description":         {optional: true},
		"ip_version":          {computed: true},
		"param_template_name": {computed: true},
		"fw_group_name":       {computed: true},
		"source_name":         {computed: true},
		"dest_name":           {computed: true},
		"create_time":         {computed: true},
		"update_time":         {computed: true},
		"last_hit_time":       {computed: true},
	}
	for name, want := range strCases {
		attr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q missing or not a StringAttribute", name)
			continue
		}
		if attr.Required != want.required || attr.Optional != want.optional || attr.Computed != want.computed {
			t.Errorf("attribute %q flags = (required=%v, optional=%v, computed=%v), want (required=%v, optional=%v, computed=%v)",
				name, attr.Required, attr.Optional, attr.Computed, want.required, want.optional, want.computed)
		}
	}

	intCases := map[string]expect{
		"order_index":   {required: true},
		"invalid":       {computed: true},
		"dns_parse_cnt": {computed: true},
		"hit_count":     {computed: true},
	}
	for name, want := range intCases {
		attr, ok := resp.Schema.Attributes[name].(schema.Int64Attribute)
		if !ok {
			t.Errorf("attribute %q missing or not an Int64Attribute", name)
			continue
		}
		if attr.Required != want.required || attr.Optional != want.optional || attr.Computed != want.computed {
			t.Errorf("attribute %q flags = (required=%v, optional=%v, computed=%v), want (required=%v, optional=%v, computed=%v)",
				name, attr.Required, attr.Optional, attr.Computed, want.required, want.optional, want.computed)
		}
	}

	expectedCount := len(strCases) + len(intCases)
	if len(resp.Schema.Attributes) != expectedCount {
		t.Errorf("schema has %d attributes, want %d", len(resp.Schema.Attributes), expectedCount)
	}
}

// TestFwmVpcAclRule2Resource_ImplementsInterfaces verifies the advertised
// framework interfaces.
func TestFwmVpcAclRule2Resource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &FwmVpcAclRule2Resource{}
	var _ resource.ResourceWithImportState = &FwmVpcAclRule2Resource{}
}

// TestFwmVpcAclRule2Resource_Factory validates the registry factory.
func TestFwmVpcAclRule2Resource_Factory(t *testing.T) {
	got := NewFwmVpcAclRule2Resource()
	if got == nil {
		t.Fatal("NewFwmVpcAclRule2Resource returned nil")
	}
	if _, ok := got.(*FwmVpcAclRule2Resource); !ok {
		t.Fatalf("factory returned %T, want *FwmVpcAclRule2Resource", got)
	}
}

// TestSplitFwmVpcAclRule2Id validates parsing of the `<group_id>#<rule_id>` id.
func TestSplitFwmVpcAclRule2Id(t *testing.T) {
	good := "fwmrg-group" + tccommon.FILED_SP + "rule-0001"
	groupId, ruleId, ok := splitFwmVpcAclRule2Id(good)
	if !ok || groupId != "fwmrg-group" || ruleId != "rule-0001" {
		t.Errorf("splitFwmVpcAclRule2Id(%q) = (%q, %q, %v), want (fwmrg-group, rule-0001, true)", good, groupId, ruleId, ok)
	}

	bad := []string{
		"",
		"no-separator",
		tccommon.FILED_SP + "rule-0001",
		"fwmrg-group" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	}
	for _, in := range bad {
		if _, _, ok := splitFwmVpcAclRule2Id(in); ok {
			t.Errorf("splitFwmVpcAclRule2Id(%q) ok=true, want false", in)
		}
	}
}

// TestFwmVpcAclRule2Resource_Create tests a full create, asserting the request
// payload and every attribute read back into state.
func TestFwmVpcAclRule2Resource_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &fwmv20250611.Client{}
	var gotCreate *fwmv20250611.CreateVpcAclRuleRequest
	patches.ApplyMethodFunc(sdkClient, "CreateVpcAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateVpcAclRuleRequest) (*fwmv20250611.CreateVpcAclRuleResponse, error) {
		gotCreate = request
		resp := fwmv20250611.NewCreateVpcAclRuleResponse()
		resp.Response = &fwmv20250611.CreateVpcAclRuleResponseParams{
			RuleIds:   []*string{fwmVpcAclRule2StrPtr("rule-0001")},
			RequestId: fwmVpcAclRule2StrPtr("fake-request-id"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(sdkClient, "DescribeVpcAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{
			Rules:         []*fwmv20250611.VpcAclRuleResp{buildFwmVpcAclRule2Resp("rule-0001", "fwmrg-group")},
			TotalCount:    fwmVpcAclRule2Int64Ptr(1),
			AllTotalCount: fwmVpcAclRule2Int64Ptr(1),
			RequestId:     fwmVpcAclRule2StrPtr("fake-request-id"),
		}
		return resp, nil
	})

	r := newFwmVpcAclRule2Resource(fwmVpcAclRule2NewConnectivityClient(patches, sdkClient))

	s := fwmVpcAclRule2ResourceSchema()
	planModel := fwmVpcAclRule2PlanModel()
	createReq := resource.CreateRequest{}
	createReq.Plan.Schema = s
	createReq.Plan.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, planModel))
	createResp := &resource.CreateResponse{}
	createResp.State.Schema = s
	createResp.State.Raw = createReq.Plan.Raw

	r.Create(context.Background(), createReq, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", createResp.Diagnostics)
	}

	if gotCreate == nil {
		t.Fatal("CreateVpcAclRule was not called")
	}
	if gotCreate.GroupId == nil || *gotCreate.GroupId != "fwmrg-group" {
		t.Errorf("request GroupId = %v, want fwmrg-group", gotCreate.GroupId)
	}
	if len(gotCreate.Rules) != 1 {
		t.Fatalf("request carries %d rules, want exactly 1", len(gotCreate.Rules))
	}
	sent := gotCreate.Rules[0]
	if sent.SourceContent == nil || *sent.SourceContent != "10.1.0.0/16" {
		t.Errorf("Rules.SourceContent = %v", sent.SourceContent)
	}
	if sent.SourceType == nil || *sent.SourceType != "ip" {
		t.Errorf("Rules.SourceType = %v", sent.SourceType)
	}
	if sent.DestContent == nil || *sent.DestContent != "10.2.0.0/16" {
		t.Errorf("Rules.DestContent = %v", sent.DestContent)
	}
	if sent.DestType == nil || *sent.DestType != "ip" {
		t.Errorf("Rules.DestType = %v", sent.DestType)
	}
	if sent.Protocol == nil || *sent.Protocol != "TCP" {
		t.Errorf("Rules.Protocol = %v", sent.Protocol)
	}
	if sent.Port == nil || *sent.Port != "8080" {
		t.Errorf("Rules.Port = %v", sent.Port)
	}
	if sent.RuleAction == nil || *sent.RuleAction != "drop" {
		t.Errorf("Rules.RuleAction = %v", sent.RuleAction)
	}
	if sent.OrderIndex == nil || *sent.OrderIndex != 5 {
		t.Errorf("Rules.OrderIndex = %v", sent.OrderIndex)
	}
	if sent.EdgeId == nil || *sent.EdgeId != "CFWS-abc" {
		t.Errorf("Rules.EdgeId = %v", sent.EdgeId)
	}
	if sent.FwGroupId == nil || *sent.FwGroupId != "cfwg-def" {
		t.Errorf("Rules.FwGroupId = %v", sent.FwGroupId)
	}
	if sent.Description == nil || *sent.Description != "managed by terraform" {
		t.Errorf("Rules.Description = %v", sent.Description)
	}
	if sent.ParamTemplateId == nil || *sent.ParamTemplateId != "ptb-123" {
		t.Errorf("Rules.ParamTemplateId = %v", sent.ParamTemplateId)
	}
	if sent.BelongMemberId == nil || *sent.BelongMemberId != "mem-0001" {
		t.Errorf("Rules.BelongMemberId = %v", sent.BelongMemberId)
	}

	fwmVpcAclRule2AssertState(t, s, createResp.State, FwmVpcAclRule2ResourceModel{
		Id:                types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001"),
		GroupId:           types.StringValue("fwmrg-group"),
		RuleId:            types.StringValue("rule-0001"),
		SourceContent:     types.StringValue("10.0.0.0/8"),
		SourceType:        types.StringValue("ip"),
		DestContent:       types.StringValue("192.168.0.0/16"),
		DestType:          types.StringValue("ip"),
		Protocol:          types.StringValue("TCP"),
		Port:              types.StringValue("80,443"),
		RuleAction:        types.StringValue("accept"),
		OrderIndex:        types.Int64Value(5),
		EdgeId:            types.StringValue("CFWS-abc"),
		FwGroupId:         types.StringValue("cfwg-def"),
		Description:       types.StringValue("drift-description"),
		ParamTemplateId:   types.StringValue("ptb-123"),
		BelongMemberId:    types.StringValue("mem-0001"),
		IpVersion:         types.StringValue("ipv4"),
		Invalid:           types.Int64Value(0),
		ParamTemplateName: types.StringValue("port-template"),
		FwGroupName:       types.StringValue("fw-group-name"),
		SourceName:        types.StringValue("source-asset"),
		DestName:          types.StringValue("dest-asset"),
		CreateTime:        types.StringValue("2025-06-11 10:00:00"),
		UpdateTime:        types.StringValue("2025-06-11 11:00:00"),
		DnsParseCnt:       types.Int64Value(3),
		HitCount:          types.Int64Value(42),
		LastHitTime:       types.StringValue("2025-06-11 12:00:00"),
	})
}

// TestFwmVpcAclRule2Resource_CreateErrors covers the failure branches of Create.
func TestFwmVpcAclRule2Resource_CreateErrors(t *testing.T) {
	cases := []struct {
		name       string
		ruleIds    []*string
		createErr  error
		describeNil bool
		describeErr error
	}{
		{name: "empty-rule-ids", ruleIds: nil},
		{name: "two-rule-ids", ruleIds: []*string{fwmVpcAclRule2StrPtr("a"), fwmVpcAclRule2StrPtr("b")}},
		{name: "empty-rule-id-string", ruleIds: []*string{fwmVpcAclRule2StrPtr("")}},
		{name: "create-api-error", createErr: errors.New("boom")},
		{name: "describe-not-found", ruleIds: []*string{fwmVpcAclRule2StrPtr("rule-0001")}, describeNil: true},
		{name: "describe-api-error", ruleIds: []*string{fwmVpcAclRule2StrPtr("rule-0001")}, describeErr: errors.New("describe boom")},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			sdkClient := &fwmv20250611.Client{}
			patches.ApplyMethodFunc(sdkClient, "CreateVpcAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.CreateVpcAclRuleRequest) (*fwmv20250611.CreateVpcAclRuleResponse, error) {
				if c.createErr != nil {
					return nil, c.createErr
				}
				resp := fwmv20250611.NewCreateVpcAclRuleResponse()
				resp.Response = &fwmv20250611.CreateVpcAclRuleResponseParams{RuleIds: c.ruleIds}
				return resp, nil
			})
			patches.ApplyMethodFunc(sdkClient, "DescribeVpcAclRulesWithContext", func(_ context.Context, _ *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
				if c.describeErr != nil {
					return nil, c.describeErr
				}
				resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
				resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{}
				if !c.describeNil {
					resp.Response.Rules = []*fwmv20250611.VpcAclRuleResp{buildFwmVpcAclRule2Resp("rule-0001", "fwmrg-group")}
				}
				return resp, nil
			})

			r := newFwmVpcAclRule2Resource(fwmVpcAclRule2NewConnectivityClient(patches, sdkClient))

			s := fwmVpcAclRule2ResourceSchema()
			createReq := resource.CreateRequest{}
			createReq.Plan.Schema = s
			createReq.Plan.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, fwmVpcAclRule2PlanModel()))
			createResp := &resource.CreateResponse{}
			createResp.State.Schema = s
			createResp.State.Raw = createReq.Plan.Raw

			r.Create(context.Background(), createReq, createResp)
			if !createResp.Diagnostics.HasError() {
				t.Fatalf("Create expected an error diagnostic, got none")
			}
		})
	}
}

// TestFwmVpcAclRule2Resource_Read tests a successful refresh (full model round
// trip) plus the "removed from the cloud" and "API error" branches.
func TestFwmVpcAclRule2Resource_Read(t *testing.T) {
	cases := []struct {
		name          string
		rules         []*fwmv20250611.VpcAclRuleResp
		describeErr   error
		expectRemoved bool
	}{
		{
			name:  "found",
			rules: []*fwmv20250611.VpcAclRuleResp{buildFwmVpcAclRule2Resp("rule-0001", "fwmrg-group")},
		},
		{
			name:          "not-found-removes-state",
			rules:         []*fwmv20250611.VpcAclRuleResp{},
			expectRemoved: true,
		},
		{
			name:  "ignores-other-rule-in-the-same-group",
			rules: []*fwmv20250611.VpcAclRuleResp{buildFwmVpcAclRule2Resp("other-rule", "fwmrg-group")},
			expectRemoved: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			sdkClient := &fwmv20250611.Client{}
			patches.ApplyMethodFunc(sdkClient, "DescribeVpcAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
				if c.describeErr != nil {
					return nil, c.describeErr
				}
				if request.GroupId == nil || *request.GroupId != "fwmrg-group" {
					return nil, errors.New("unexpected group id")
				}
				resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
				resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{Rules: c.rules}
				return resp, nil
			})

			r := newFwmVpcAclRule2Resource(fwmVpcAclRule2NewConnectivityClient(patches, sdkClient))

			s := fwmVpcAclRule2ResourceSchema()
			stateModel := fwmVpcAclRule2PlanModel()
			stateModel.Id = types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001")
			stateModel.RuleId = types.StringValue("rule-0001")

			readReq := resource.ReadRequest{}
			readReq.State.Schema = s
			readReq.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, stateModel))
			readResp := &resource.ReadResponse{}
			readResp.State.Schema = s
			readResp.State.Raw = readReq.State.Raw

			r.Read(context.Background(), readReq, readResp)
			if readResp.Diagnostics.HasError() {
				t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
			}

			if c.expectRemoved {
				if !readResp.State.Raw.IsNull() {
					t.Fatalf("state should be removed for a rule that no longer exists")
				}
				return
			}

			fwmVpcAclRule2AssertState(t, s, readResp.State, FwmVpcAclRule2ResourceModel{
				Id:                types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001"),
				GroupId:           types.StringValue("fwmrg-group"),
				RuleId:            types.StringValue("rule-0001"),
				SourceContent:     types.StringValue("10.0.0.0/8"),
				SourceType:        types.StringValue("ip"),
				DestContent:       types.StringValue("192.168.0.0/16"),
				DestType:          types.StringValue("ip"),
				Protocol:          types.StringValue("TCP"),
				Port:              types.StringValue("80,443"),
				RuleAction:        types.StringValue("accept"),
				OrderIndex:        types.Int64Value(5),
				EdgeId:            types.StringValue("CFWS-abc"),
				FwGroupId:         types.StringValue("cfwg-def"),
				Description:       types.StringValue("drift-description"),
				ParamTemplateId:   types.StringValue("ptb-123"),
				BelongMemberId:    types.StringValue("mem-0001"),
				IpVersion:         types.StringValue("ipv4"),
				Invalid:           types.Int64Value(0),
				ParamTemplateName: types.StringValue("port-template"),
				FwGroupName:       types.StringValue("fw-group-name"),
				SourceName:        types.StringValue("source-asset"),
				DestName:          types.StringValue("dest-asset"),
				CreateTime:        types.StringValue("2025-06-11 10:00:00"),
				UpdateTime:        types.StringValue("2025-06-11 11:00:00"),
				DnsParseCnt:       types.Int64Value(3),
				HitCount:          types.Int64Value(42),
				LastHitTime:       types.StringValue("2025-06-11 12:00:00"),
			})
		})
	}

	t.Run("api-error-keeps-state", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		sdkClient := &fwmv20250611.Client{}
		patches.ApplyMethodFunc(sdkClient, "DescribeVpcAclRulesWithContext", func(_ context.Context, _ *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
			return nil, errors.New("permission denied")
		})

		r := newFwmVpcAclRule2Resource(fwmVpcAclRule2NewConnectivityClient(patches, sdkClient))

		s := fwmVpcAclRule2ResourceSchema()
		stateModel := fwmVpcAclRule2PlanModel()
		stateModel.Id = types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001")
		stateModel.RuleId = types.StringValue("rule-0001")

		readReq := resource.ReadRequest{}
		readReq.State.Schema = s
		readReq.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, stateModel))
		readResp := &resource.ReadResponse{}
		readResp.State.Schema = s
		readResp.State.Raw = readReq.State.Raw

		r.Read(context.Background(), readReq, readResp)
		if !readResp.Diagnostics.HasError() {
			t.Fatal("Read expected an error diagnostic for an API failure")
		}
		if readResp.State.Raw.IsNull() {
			t.Fatal("an API error must not remove the resource from state")
		}
	})

	t.Run("broken-id", func(t *testing.T) {
		r := newFwmVpcAclRule2Resource(&connectivity.TencentCloudClient{})
		s := fwmVpcAclRule2ResourceSchema()
		stateModel := fwmVpcAclRule2PlanModel()
		stateModel.Id = types.StringValue("broken-id")
		stateModel.RuleId = types.StringValue("rule-0001")

		readReq := resource.ReadRequest{}
		readReq.State.Schema = s
		readReq.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, stateModel))
		readResp := &resource.ReadResponse{}
		readResp.State.Schema = s
		readResp.State.Raw = readReq.State.Raw

		r.Read(context.Background(), readReq, readResp)
		if !readResp.Diagnostics.HasError() {
			t.Fatal("Read expected an error diagnostic for a broken resource id")
		}
	})
}

// TestFwmVpcAclRule2Resource_Update tests the modify request and the refreshed
// state, including the "optional field removed from configuration" path.
func TestFwmVpcAclRule2Resource_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &fwmv20250611.Client{}
	var gotModify *fwmv20250611.ModifyVpcAclRuleRequest
	patches.ApplyMethodFunc(sdkClient, "ModifyVpcAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyVpcAclRuleRequest) (*fwmv20250611.ModifyVpcAclRuleResponse, error) {
		gotModify = request
		resp := fwmv20250611.NewModifyVpcAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyVpcAclRuleResponseParams{RequestId: fwmVpcAclRule2StrPtr("fake-request-id")}
		return resp, nil
	})
	patches.ApplyMethodFunc(sdkClient, "DescribeVpcAclRulesWithContext", func(_ context.Context, _ *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{
			Rules: []*fwmv20250611.VpcAclRuleResp{buildFwmVpcAclRule2Resp("rule-0001", "fwmrg-group")},
		}
		return resp, nil
	})

	r := newFwmVpcAclRule2Resource(fwmVpcAclRule2NewConnectivityClient(patches, sdkClient))

	s := fwmVpcAclRule2ResourceSchema()

	stateModel := fwmVpcAclRule2PlanModel()
	stateModel.Id = types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001")
	stateModel.RuleId = types.StringValue("rule-0001")

	// The user removes description/param_template_id from the configuration:
	// the plan carries Null for them, the request must omit them and the
	// refreshed state must keep the previous values.
	planModel := fwmVpcAclRule2PlanModel()
	planModel.Id = types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001")
	planModel.RuleId = types.StringValue("rule-0001")
	planModel.Description = types.StringNull()
	planModel.ParamTemplateId = types.StringNull()
	planModel.RuleAction = types.StringValue("log")
	planModel.OrderIndex = types.Int64Value(6)

	updateReq := resource.UpdateRequest{}
	updateReq.Plan.Schema = s
	updateReq.Plan.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, planModel))
	updateReq.State.Schema = s
	updateReq.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, stateModel))
	updateResp := &resource.UpdateResponse{}
	updateResp.State.Schema = s
	updateResp.State.Raw = updateReq.State.Raw

	r.Update(context.Background(), updateReq, updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", updateResp.Diagnostics)
	}

	if gotModify == nil {
		t.Fatal("ModifyVpcAclRule was not called")
	}
	if gotModify.GroupId == nil || *gotModify.GroupId != "fwmrg-group" {
		t.Errorf("request GroupId = %v, want fwmrg-group", gotModify.GroupId)
	}
	sent := gotModify.Rule
	if sent == nil {
		t.Fatal("request Rule is nil")
	}
	if sent.RuleId == nil || *sent.RuleId != "rule-0001" {
		t.Errorf("Rule.RuleId = %v, want rule-0001", sent.RuleId)
	}
	if sent.RuleAction == nil || *sent.RuleAction != "log" {
		t.Errorf("Rule.RuleAction = %v, want log", sent.RuleAction)
	}
	if sent.OrderIndex == nil || *sent.OrderIndex != 6 {
		t.Errorf("Rule.OrderIndex = %v, want 6", sent.OrderIndex)
	}
	// Removing an attribute from the configuration must keep the cloud value:
	// the previous state value is sent again instead of clearing it.
	if sent.Description == nil || *sent.Description != "managed by terraform" {
		t.Errorf("Rule.Description = %v, want the previous state value", sent.Description)
	}
	if sent.ParamTemplateId == nil || *sent.ParamTemplateId != "ptb-123" {
		t.Errorf("Rule.ParamTemplateId = %v, want the previous state value", sent.ParamTemplateId)
	}
	if sent.SourceContent == nil || *sent.SourceContent != "10.1.0.0/16" {
		t.Errorf("Rule.SourceContent = %v", sent.SourceContent)
	}
	if sent.SourceType == nil || *sent.SourceType != "ip" {
		t.Errorf("Rule.SourceType = %v", sent.SourceType)
	}
	if sent.DestContent == nil || *sent.DestContent != "10.2.0.0/16" {
		t.Errorf("Rule.DestContent = %v", sent.DestContent)
	}
	if sent.DestType == nil || *sent.DestType != "ip" {
		t.Errorf("Rule.DestType = %v", sent.DestType)
	}
	if sent.Protocol == nil || *sent.Protocol != "TCP" {
		t.Errorf("Rule.Protocol = %v", sent.Protocol)
	}
	if sent.Port == nil || *sent.Port != "8080" {
		t.Errorf("Rule.Port = %v", sent.Port)
	}
	if sent.EdgeId == nil || *sent.EdgeId != "CFWS-abc" {
		t.Errorf("Rule.EdgeId = %v", sent.EdgeId)
	}
	if sent.FwGroupId == nil || *sent.FwGroupId != "cfwg-def" {
		t.Errorf("Rule.FwGroupId = %v", sent.FwGroupId)
	}
	if sent.BelongMemberId == nil || *sent.BelongMemberId != "mem-0001" {
		t.Errorf("Rule.BelongMemberId = %v", sent.BelongMemberId)
	}

	fwmVpcAclRule2AssertState(t, s, updateResp.State, FwmVpcAclRule2ResourceModel{
		Id:                types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001"),
		GroupId:           types.StringValue("fwmrg-group"),
		RuleId:            types.StringValue("rule-0001"),
		SourceContent:     types.StringValue("10.0.0.0/8"),
		SourceType:        types.StringValue("ip"),
		DestContent:       types.StringValue("192.168.0.0/16"),
		DestType:          types.StringValue("ip"),
		Protocol:          types.StringValue("TCP"),
		Port:              types.StringValue("80,443"),
		RuleAction:        types.StringValue("accept"),
		OrderIndex:        types.Int64Value(5),
		EdgeId:            types.StringValue("CFWS-abc"),
		FwGroupId:         types.StringValue("cfwg-def"),
		Description:       types.StringValue("drift-description"),
		ParamTemplateId:   types.StringValue("ptb-123"),
		BelongMemberId:    types.StringValue("mem-0001"),
		IpVersion:         types.StringValue("ipv4"),
		Invalid:           types.Int64Value(0),
		ParamTemplateName: types.StringValue("port-template"),
		FwGroupName:       types.StringValue("fw-group-name"),
		SourceName:        types.StringValue("source-asset"),
		DestName:          types.StringValue("dest-asset"),
		CreateTime:        types.StringValue("2025-06-11 10:00:00"),
		UpdateTime:        types.StringValue("2025-06-11 11:00:00"),
		DnsParseCnt:       types.Int64Value(3),
		HitCount:          types.Int64Value(42),
		LastHitTime:       types.StringValue("2025-06-11 12:00:00"),
	})
}

// TestFwmVpcAclRule2Resource_UpdateErrors covers the failure branches of Update.
func TestFwmVpcAclRule2Resource_UpdateErrors(t *testing.T) {
	cases := []struct {
		name        string
		modifyErr   error
		describeNil bool
		describeErr error
		brokenId    bool
	}{
		{name: "modify-api-error", modifyErr: errors.New("modify boom")},
		{name: "describe-not-found", describeNil: true},
		{name: "describe-api-error", describeErr: errors.New("describe boom")},
		{name: "broken-id", brokenId: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			sdkClient := &fwmv20250611.Client{}
			patches.ApplyMethodFunc(sdkClient, "ModifyVpcAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.ModifyVpcAclRuleRequest) (*fwmv20250611.ModifyVpcAclRuleResponse, error) {
				if c.modifyErr != nil {
					return nil, c.modifyErr
				}
				resp := fwmv20250611.NewModifyVpcAclRuleResponse()
				resp.Response = &fwmv20250611.ModifyVpcAclRuleResponseParams{}
				return resp, nil
			})
			patches.ApplyMethodFunc(sdkClient, "DescribeVpcAclRulesWithContext", func(_ context.Context, _ *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
				if c.describeErr != nil {
					return nil, c.describeErr
				}
				resp := fwmv20250611.NewDescribeVpcAclRulesResponse()
				resp.Response = &fwmv20250611.DescribeVpcAclRulesResponseParams{}
				if !c.describeNil {
					resp.Response.Rules = []*fwmv20250611.VpcAclRuleResp{buildFwmVpcAclRule2Resp("rule-0001", "fwmrg-group")}
				}
				return resp, nil
			})

			r := newFwmVpcAclRule2Resource(fwmVpcAclRule2NewConnectivityClient(patches, sdkClient))

			s := fwmVpcAclRule2ResourceSchema()
			stateModel := fwmVpcAclRule2PlanModel()
			stateModel.Id = types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001")
			stateModel.RuleId = types.StringValue("rule-0001")
			if c.brokenId {
				stateModel.Id = types.StringValue("broken-id")
			}

			planModel := stateModel

			updateReq := resource.UpdateRequest{}
			updateReq.Plan.Schema = s
			updateReq.Plan.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, planModel))
			updateReq.State.Schema = s
			updateReq.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, stateModel))
			updateResp := &resource.UpdateResponse{}
			updateResp.State.Schema = s
			updateResp.State.Raw = updateReq.State.Raw

			r.Update(context.Background(), updateReq, updateResp)
			if !updateResp.Diagnostics.HasError() {
				t.Fatal("Update expected an error diagnostic, got none")
			}
		})
	}
}

// TestFwmVpcAclRule2Resource_Delete covers the delete request payload and the
// error branches.
func TestFwmVpcAclRule2Resource_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		sdkClient := &fwmv20250611.Client{}
		var gotDelete *fwmv20250611.DeleteVpcAclRuleRequest
		patches.ApplyMethodFunc(sdkClient, "DeleteVpcAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteVpcAclRuleRequest) (*fwmv20250611.DeleteVpcAclRuleResponse, error) {
			gotDelete = request
			resp := fwmv20250611.NewDeleteVpcAclRuleResponse()
			resp.Response = &fwmv20250611.DeleteVpcAclRuleResponseParams{RequestId: fwmVpcAclRule2StrPtr("fake-request-id")}
			return resp, nil
		})

		r := newFwmVpcAclRule2Resource(fwmVpcAclRule2NewConnectivityClient(patches, sdkClient))

		s := fwmVpcAclRule2ResourceSchema()
		stateModel := fwmVpcAclRule2PlanModel()
		stateModel.Id = types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001")
		stateModel.RuleId = types.StringValue("rule-0001")

		deleteReq := resource.DeleteRequest{}
		deleteReq.State.Schema = s
		deleteReq.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, stateModel))
		deleteResp := &resource.DeleteResponse{}
		deleteResp.State.Schema = s
		deleteResp.State.Raw = deleteReq.State.Raw

		r.Delete(context.Background(), deleteReq, deleteResp)
		if deleteResp.Diagnostics.HasError() {
			t.Fatalf("Delete returned diagnostics: %v", deleteResp.Diagnostics)
		}

		if gotDelete == nil {
			t.Fatal("DeleteVpcAclRule was not called")
		}
		if gotDelete.GroupId == nil || *gotDelete.GroupId != "fwmrg-group" {
			t.Errorf("request GroupId = %v, want fwmrg-group", gotDelete.GroupId)
		}
		if len(gotDelete.RuleIds) != 1 || gotDelete.RuleIds[0] == nil || *gotDelete.RuleIds[0] != "rule-0001" {
			t.Errorf("request RuleIds = %v, want [rule-0001]", gotDelete.RuleIds)
		}
	})

	t.Run("api-error", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		sdkClient := &fwmv20250611.Client{}
		patches.ApplyMethodFunc(sdkClient, "DeleteVpcAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.DeleteVpcAclRuleRequest) (*fwmv20250611.DeleteVpcAclRuleResponse, error) {
			return nil, errors.New("delete boom")
		})

		r := newFwmVpcAclRule2Resource(fwmVpcAclRule2NewConnectivityClient(patches, sdkClient))

		s := fwmVpcAclRule2ResourceSchema()
		stateModel := fwmVpcAclRule2PlanModel()
		stateModel.Id = types.StringValue("fwmrg-group" + tccommon.FILED_SP + "rule-0001")
		stateModel.RuleId = types.StringValue("rule-0001")

		deleteReq := resource.DeleteRequest{}
		deleteReq.State.Schema = s
		deleteReq.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, stateModel))
		deleteResp := &resource.DeleteResponse{}
		deleteResp.State.Schema = s
		deleteResp.State.Raw = deleteReq.State.Raw

		r.Delete(context.Background(), deleteReq, deleteResp)
		if !deleteResp.Diagnostics.HasError() {
			t.Fatal("Delete expected an error diagnostic for an API failure")
		}
	})

	t.Run("broken-id", func(t *testing.T) {
		r := newFwmVpcAclRule2Resource(&connectivity.TencentCloudClient{})
		s := fwmVpcAclRule2ResourceSchema()
		stateModel := fwmVpcAclRule2PlanModel()
		stateModel.Id = types.StringValue("broken-id")
		stateModel.RuleId = types.StringValue("rule-0001")

		deleteReq := resource.DeleteRequest{}
		deleteReq.State.Schema = s
		deleteReq.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), fwmVpcAclRule2ModelValues(s, stateModel))
		deleteResp := &resource.DeleteResponse{}
		deleteResp.State.Schema = s
		deleteResp.State.Raw = deleteReq.State.Raw

		r.Delete(context.Background(), deleteReq, deleteResp)
		if !deleteResp.Diagnostics.HasError() {
			t.Fatal("Delete expected an error diagnostic for a broken resource id")
		}
	})
}

// TestFwmVpcAclRule2Resource_ImportState covers the composite import id.
func TestFwmVpcAclRule2Resource_ImportState(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		r := &FwmVpcAclRule2Resource{}
		s := fwmVpcAclRule2ResourceSchema()

		importResp := &resource.ImportStateResponse{}
		importResp.State.Schema = s
		importResp.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), nil)

		r.ImportState(context.Background(), resource.ImportStateRequest{ID: "fwmrg-group" + tccommon.FILED_SP + "rule-0001"}, importResp)
		if importResp.Diagnostics.HasError() {
			t.Fatalf("ImportState returned diagnostics: %v", importResp.Diagnostics)
		}

		var idVal, groupVal, ruleVal types.String
		if diags := importResp.State.GetAttribute(context.Background(), path.Root("id"), &idVal); diags.HasError() {
			t.Fatalf("reading imported id failed: %v", diags)
		}
		if diags := importResp.State.GetAttribute(context.Background(), path.Root("group_id"), &groupVal); diags.HasError() {
			t.Fatalf("reading imported group_id failed: %v", diags)
		}
		if diags := importResp.State.GetAttribute(context.Background(), path.Root("rule_id"), &ruleVal); diags.HasError() {
			t.Fatalf("reading imported rule_id failed: %v", diags)
		}

		if idVal.ValueString() != "fwmrg-group"+tccommon.FILED_SP+"rule-0001" {
			t.Errorf("id = %q", idVal.ValueString())
		}
		if groupVal.ValueString() != "fwmrg-group" {
			t.Errorf("group_id = %q, want fwmrg-group", groupVal.ValueString())
		}
		if ruleVal.ValueString() != "rule-0001" {
			t.Errorf("rule_id = %q, want rule-0001", ruleVal.ValueString())
		}
	})

	t.Run("invalid", func(t *testing.T) {
		r := &FwmVpcAclRule2Resource{}
		s := fwmVpcAclRule2ResourceSchema()

		importResp := &resource.ImportStateResponse{}
		importResp.State.Schema = s
		importResp.State.Raw = tftypes.NewValue(fwmVpcAclRule2ObjectType(s), nil)

		r.ImportState(context.Background(), resource.ImportStateRequest{ID: "broken-id"}, importResp)
		if !importResp.Diagnostics.HasError() {
			t.Fatal("ImportState expected an error diagnostic for a broken import id")
		}
	})
}

// TestRefreshFwmVpcAclRule2Model covers the "cloud stops returning a value"
// paths: configured values must be kept, never blanked out.
func TestRefreshFwmVpcAclRule2Model(t *testing.T) {
	m := FwmVpcAclRule2ResourceModel{
		Port:            types.StringValue("8080"),
		Description:     types.StringValue("keep-me"),
		ParamTemplateId: types.StringValue("ptb-123"),
		OrderIndex:      types.Int64Value(9),
		BelongMemberId:  types.StringValue("mem-0001"),
	}

	item := &fwmv20250611.VpcAclRuleResp{
		RuleId: fwmVpcAclRule2StrPtr("rule-0001"),
		// Sequence deliberately left nil: the configured priority must survive.
	}
	refreshVpcAclRule2Model(&m, item)

	if m.OrderIndex.IsNull() || m.OrderIndex.ValueInt64() != 9 {
		t.Errorf("order_index = %v, want 9 (configured value kept)", m.OrderIndex)
	}
	if m.Port.ValueString() != "8080" {
		t.Errorf("port = %v, want 8080 (configured value kept)", m.Port)
	}
	if m.Description.ValueString() != "keep-me" {
		t.Errorf("description = %v, want keep-me (configured value kept)", m.Description)
	}
	if m.ParamTemplateId.ValueString() != "ptb-123" {
		t.Errorf("param_template_id = %v, want ptb-123", m.ParamTemplateId)
	}
	if m.BelongMemberId.ValueString() != "mem-0001" {
		t.Errorf("belong_member_id = %v, want mem-0001", m.BelongMemberId)
	}

	// A nil item must not touch the model at all.
	before := m
	refreshVpcAclRule2Model(&m, nil)
	if m.Port.ValueString() != before.Port.ValueString() ||
		m.Description.ValueString() != before.Description.ValueString() ||
		m.ParamTemplateId.ValueString() != before.ParamTemplateId.ValueString() ||
		m.BelongMemberId.ValueString() != before.BelongMemberId.ValueString() ||
		m.OrderIndex.ValueInt64() != before.OrderIndex.ValueInt64() {
		t.Errorf("refresh with a nil item modified the model")
	}
}

// TestMergeFwmVpcAclRule2Plan verifies that null plan values fall back to the
// previous state while known values win.
func TestMergeFwmVpcAclRule2Plan(t *testing.T) {
	state := FwmVpcAclRule2ResourceModel{
		Port:            types.StringValue("8080"),
		Description:     types.StringValue("old-description"),
		ParamTemplateId: types.StringValue("ptb-123"),
		BelongMemberId:  types.StringValue("mem-0001"),
		CreateTime:      types.StringValue("2025-06-11 10:00:00"),
		HitCount:        types.Int64Value(7),
	}
	plan := FwmVpcAclRule2ResourceModel{
		Port:            types.StringValue("9090"),
		Description:     types.StringNull(),
		ParamTemplateId: types.StringUnknown(),
		BelongMemberId:  types.StringNull(),
		CreateTime:      types.StringNull(),
		HitCount:        types.Int64Null(),
	}

	mergeVpcAclRule2Plan(state, &plan)

	if plan.Port.ValueString() != "9090" {
		t.Errorf("port = %v, want 9090 (explicit plan value wins)", plan.Port)
	}
	if plan.Description.ValueString() != "old-description" {
		t.Errorf("description = %v, want old-description", plan.Description)
	}
	if plan.ParamTemplateId.ValueString() != "ptb-123" {
		t.Errorf("param_template_id = %v, want ptb-123", plan.ParamTemplateId)
	}
	if plan.BelongMemberId.ValueString() != "mem-0001" {
		t.Errorf("belong_member_id = %v, want mem-0001", plan.BelongMemberId)
	}
	if plan.CreateTime.ValueString() != "2025-06-11 10:00:00" {
		t.Errorf("create_time = %v, want the previous state value", plan.CreateTime)
	}
	if plan.HitCount.IsNull() || plan.HitCount.ValueInt64() != 7 {
		t.Errorf("hit_count = %v, want 7", plan.HitCount)
	}
}

// fwmVpcAclRule2AssertState asserts that the state contains exactly the
// expected model, leaf field by leaf field.
func fwmVpcAclRule2AssertState(t *testing.T, s schema.Schema, state tfsdk.State, want FwmVpcAclRule2ResourceModel) {
	t.Helper()

	var got FwmVpcAclRule2ResourceModel
	if diags := state.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("reading state failed: %v", diags)
	}

	stringFields := []struct {
		name string
		got  types.String
		want types.String
	}{
		{"id", got.Id, want.Id},
		{"group_id", got.GroupId, want.GroupId},
		{"rule_id", got.RuleId, want.RuleId},
		{"source_content", got.SourceContent, want.SourceContent},
		{"source_type", got.SourceType, want.SourceType},
		{"dest_content", got.DestContent, want.DestContent},
		{"dest_type", got.DestType, want.DestType},
		{"protocol", got.Protocol, want.Protocol},
		{"port", got.Port, want.Port},
		{"rule_action", got.RuleAction, want.RuleAction},
		{"edge_id", got.EdgeId, want.EdgeId},
		{"fw_group_id", got.FwGroupId, want.FwGroupId},
		{"description", got.Description, want.Description},
		{"param_template_id", got.ParamTemplateId, want.ParamTemplateId},
		{"belong_member_id", got.BelongMemberId, want.BelongMemberId},
		{"ip_version", got.IpVersion, want.IpVersion},
		{"param_template_name", got.ParamTemplateName, want.ParamTemplateName},
		{"fw_group_name", got.FwGroupName, want.FwGroupName},
		{"source_name", got.SourceName, want.SourceName},
		{"dest_name", got.DestName, want.DestName},
		{"create_time", got.CreateTime, want.CreateTime},
		{"update_time", got.UpdateTime, want.UpdateTime},
		{"last_hit_time", got.LastHitTime, want.LastHitTime},
	}
	for _, f := range stringFields {
		if f.got.IsNull() != f.want.IsNull() {
			t.Errorf("%s: null mismatch (got null=%v, want null=%v)", f.name, f.got.IsNull(), f.want.IsNull())
			continue
		}
		if !f.got.IsNull() && f.got.ValueString() != f.want.ValueString() {
			t.Errorf("%s = %q, want %q", f.name, f.got.ValueString(), f.want.ValueString())
		}
	}

	intFields := []struct {
		name string
		got  types.Int64
		want types.Int64
	}{
		{"order_index", got.OrderIndex, want.OrderIndex},
		{"invalid", got.Invalid, want.Invalid},
		{"dns_parse_cnt", got.DnsParseCnt, want.DnsParseCnt},
		{"hit_count", got.HitCount, want.HitCount},
	}
	for _, f := range intFields {
		if f.got.IsNull() != f.want.IsNull() {
			t.Errorf("%s: null mismatch (got null=%v, want null=%v)", f.name, f.got.IsNull(), f.want.IsNull())
			continue
		}
		if !f.got.IsNull() && f.got.ValueInt64() != f.want.ValueInt64() {
			t.Errorf("%s = %d, want %d", f.name, f.got.ValueInt64(), f.want.ValueInt64())
		}
	}
}
