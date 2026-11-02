package fwm

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// TestFwmVpcAclRule2_Metadata validates the resource type name.
func TestFwmVpcAclRule2_Metadata(t *testing.T) {
	r := &FwmVpcAclRule2Resource{}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	if resp.TypeName != "tencentcloud_fwm_vpc_acl_rule2" {
		t.Fatalf("TypeName = %q, want %q", resp.TypeName, "tencentcloud_fwm_vpc_acl_rule2")
	}
}

// TestFwmVpcAclRule2_Factory validates the registry factory.
func TestFwmVpcAclRule2_Factory(t *testing.T) {
	got := NewFwmVpcAclRule2Resource()
	if got == nil {
		t.Fatal("NewFwmVpcAclRule2Resource returned nil")
	}
	if _, ok := got.(*FwmVpcAclRule2Resource); !ok {
		t.Fatalf("factory returned %T, want *FwmVpcAclRule2Resource", got)
	}
}

// TestFwmVpcAclRule2_ImplementsInterfaces checks the advertised interfaces.
func TestFwmVpcAclRule2_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &FwmVpcAclRule2Resource{}
	var _ resource.ResourceWithImportState = &FwmVpcAclRule2Resource{}
}

// TestFwmVpcAclRule2_ConfigureNilProviderData ensures a nil ProviderData is a no-op.
func TestFwmVpcAclRule2_ConfigureNilProviderData(t *testing.T) {
	r := &FwmVpcAclRule2Resource{}
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestFwmVpcAclRule2_Schema checks the required / optional / computed flags of
// every attribute declared by the resource.
func TestFwmVpcAclRule2_Schema(t *testing.T) {
	r := &FwmVpcAclRule2Resource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	stringAttrs := map[string]struct{ required, optional, computed bool }{
		"id":                  {computed: true},
		"group_id":            {required: true},
		"rule_id":             {computed: true},
		"source_content":      {required: true},
		"source_type":         {required: true},
		"dest_content":        {required: true},
		"dest_type":           {required: true},
		"protocol":            {required: true},
		"port":                {optional: true},
		"rule_action":         {required: true},
		"edge_id":             {required: true},
		"fw_group_id":         {required: true},
		"description":         {optional: true},
		"param_template_id":   {optional: true},
		"belong_member_id":    {optional: true},
		"ip_version":          {computed: true},
		"source_name":         {computed: true},
		"dest_name":           {computed: true},
		"param_template_name": {computed: true},
		"app_id":              {computed: true},
		"uin":                 {computed: true},
		"nickname":            {computed: true},
		"node_name":           {computed: true},
		"create_time":         {computed: true},
		"update_time":         {computed: true},
		"fw_group_name":       {computed: true},
		"last_hit_time":       {computed: true},
	}
	for name, want := range stringAttrs {
		attr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q missing or not a StringAttribute", name)
			continue
		}
		if attr.Required != want.required || attr.Optional != want.optional || attr.Computed != want.computed {
			t.Errorf("attribute %q flags = (required=%v optional=%v computed=%v), want (%v %v %v)",
				name, attr.Required, attr.Optional, attr.Computed, want.required, want.optional, want.computed)
		}
	}

	intAttrs := map[string]struct{ required, optional, computed bool }{
		"order_index":   {required: true},
		"sequence":      {computed: true},
		"invalid":       {computed: true},
		"dns_parse_cnt": {computed: true},
		"hit_count":     {computed: true},
	}
	for name, want := range intAttrs {
		attr, ok := resp.Schema.Attributes[name].(schema.Int64Attribute)
		if !ok {
			t.Errorf("attribute %q missing or not an Int64Attribute", name)
			continue
		}
		if attr.Required != want.required || attr.Optional != want.optional || attr.Computed != want.computed {
			t.Errorf("attribute %q flags = (required=%v optional=%v computed=%v), want (%v %v %v)",
				name, attr.Required, attr.Optional, attr.Computed, want.required, want.optional, want.computed)
		}
	}

	if len(resp.Schema.Attributes) != len(stringAttrs)+len(intAttrs) {
		t.Errorf("schema has %d attributes, want %d", len(resp.Schema.Attributes), len(stringAttrs)+len(intAttrs))
	}
}

// TestSplitFwmVpcAclRule2Id covers the valid and invalid composite id forms.
func TestSplitFwmVpcAclRule2Id(t *testing.T) {
	good := "fwmrg-9z03vakjj8" + tccommon.FILED_SP + "cfwrule-1"
	g, rid, ok := splitFwmVpcAclRule2Id(good)
	if !ok || g != "fwmrg-9z03vakjj8" || rid != "cfwrule-1" {
		t.Errorf("splitFwmVpcAclRule2Id(%q) = (%q, %q, %v)", good, g, rid, ok)
	}

	for _, bad := range []string{
		"",
		"only-group",
		tccommon.FILED_SP + "rule",
		"group" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	} {
		if _, _, ok := splitFwmVpcAclRule2Id(bad); ok {
			t.Errorf("splitFwmVpcAclRule2Id(%q) ok=true, want false", bad)
		}
	}
}

// TestBuildVpcAclRule2 verifies every required and optional input is mapped
// onto the SDK request struct, and that unset optional fields stay nil.
func TestBuildVpcAclRule2(t *testing.T) {
	m := FwmVpcAclRule2ResourceModel{
		GroupId:         types.StringValue("fwmrg-1"),
		RuleId:          types.StringValue("rule-1"),
		SourceContent:   types.StringValue("10.0.0.0/16"),
		SourceType:      types.StringValue("ip"),
		DestContent:     types.StringValue("0.0.0.0/0"),
		DestType:        types.StringValue("ip"),
		Protocol:        types.StringValue("TCP"),
		Port:            types.StringValue("80"),
		RuleAction:      types.StringValue("accept"),
		OrderIndex:      types.Int64Value(3),
		EdgeId:          types.StringValue("ALL"),
		FwGroupId:       types.StringValue("ALL"),
		Description:     types.StringValue("desc"),
		ParamTemplateId: types.StringValue("pt-1"),
		BelongMemberId:  types.StringValue("mem-1"),
	}

	rule := buildVpcAclRule2(m)
	if rule == nil {
		t.Fatal("buildVpcAclRule2 returned nil")
	}
	checks := map[string]struct{ got, want *string }{
		"SourceContent":   {rule.SourceContent, ptrOf("10.0.0.0/16")},
		"SourceType":      {rule.SourceType, ptrOf("ip")},
		"DestContent":     {rule.DestContent, ptrOf("0.0.0.0/0")},
		"DestType":        {rule.DestType, ptrOf("ip")},
		"Protocol":        {rule.Protocol, ptrOf("TCP")},
		"Port":            {rule.Port, ptrOf("80")},
		"RuleAction":      {rule.RuleAction, ptrOf("accept")},
		"EdgeId":          {rule.EdgeId, ptrOf("ALL")},
		"FwGroupId":       {rule.FwGroupId, ptrOf("ALL")},
		"RuleId":          {rule.RuleId, ptrOf("rule-1")},
		"Description":     {rule.Description, ptrOf("desc")},
		"ParamTemplateId": {rule.ParamTemplateId, ptrOf("pt-1")},
		"BelongMemberId":  {rule.BelongMemberId, ptrOf("mem-1")},
	}
	for name, c := range checks {
		if c.got == nil || *c.got != *c.want {
			t.Errorf("%s = %v, want %q", name, c.got, *c.want)
		}
	}
	if rule.OrderIndex == nil || *rule.OrderIndex != 3 {
		t.Errorf("OrderIndex = %v, want 3", rule.OrderIndex)
	}

	empty := buildVpcAclRule2(FwmVpcAclRule2ResourceModel{
		SourceContent: types.StringValue("a"),
		SourceType:    types.StringValue("ip"),
		DestContent:   types.StringValue("b"),
		DestType:      types.StringValue("ip"),
		Protocol:      types.StringValue("ANY"),
		RuleAction:    types.StringValue("drop"),
		EdgeId:        types.StringValue("ALL"),
		FwGroupId:     types.StringValue("ALL"),
		OrderIndex:    types.Int64Value(1),
		Port:          types.StringNull(),
		RuleId:        types.StringNull(),
	})
	if empty.Port != nil {
		t.Errorf("Port = %v, want nil for a null attribute", empty.Port)
	}
	if empty.RuleId != nil {
		t.Errorf("RuleId = %v, want nil for a null attribute", empty.RuleId)
	}
	if empty.Description != nil || empty.ParamTemplateId != nil || empty.BelongMemberId != nil {
		t.Errorf("unset optional attributes must map to nil, got %+v", empty)
	}
}

// buildVpcAclRule2Resp returns a fully populated DescribeVpcAclRules item.
func buildVpcAclRule2Resp() *fwmv20250611.VpcAclRuleResp {
	return &fwmv20250611.VpcAclRuleResp{
		RuleId:            ptrOf("rule-1"),
		Sequence:          ptrOfInt64(5),
		IpVersion:         ptrOf("ipv4"),
		SourceContent:     ptrOf("10.0.0.0/16"),
		SourceType:        ptrOf("ip"),
		SourceName:        ptrOf("src-asset"),
		DestContent:       ptrOf("0.0.0.0/0"),
		DestType:          ptrOf("ip"),
		DestName:          ptrOf("dst-asset"),
		Protocol:          ptrOf("TCP"),
		Port:              ptrOf("80"),
		ParamTemplateId:   ptrOf("pt-1"),
		ParamTemplateName: ptrOf("pt-name"),
		RuleAction:        ptrOf("accept"),
		Description:       ptrOf("desc"),
		EdgeId:            ptrOf("ALL"),
		FwGroupId:         ptrOf("ALL"),
		Invalid:           ptrOfInt64(0),
		BelongMember: &fwmv20250611.MemberInfo{
			AppId:    ptrOf("app-1"),
			Uin:      ptrOf("uin-1"),
			Nickname: ptrOf("nick"),
			MemberId: ptrOf("mem-1"),
			NodeName: ptrOf("dept"),
		},
		CreateTime:  ptrOf("2025-01-01 00:00:00"),
		UpdateTime:  ptrOf("2025-01-02 00:00:00"),
		DnsParseCnt: ptrOfInt64(2),
		FwGroupName: ptrOf("fw-group-name"),
		HitCount:    ptrOfInt64(42),
		LastHitTime: ptrOf("2025-01-03 00:00:00"),
	}
}

// TestRefreshVpcAclRule2Model verifies that every leaf attribute is written
// back from the API response and that absent optional values are preserved.
func TestRefreshVpcAclRule2Model(t *testing.T) {
	m := FwmVpcAclRule2ResourceModel{
		Id:              types.StringValue("fwmrg-1#rule-1"),
		GroupId:         types.StringValue("fwmrg-1"),
		RuleId:          types.StringValue("rule-1"),
		OrderIndex:      types.Int64Value(5),
		Port:            types.StringValue("80"),
		Description:     types.StringValue("desc"),
		ParamTemplateId: types.StringValue("pt-1"),
		BelongMemberId:  types.StringValue("mem-1"),
	}

	refreshVpcAclRule2Model(&m, buildVpcAclRule2Resp())

	strChecks := map[string]struct {
		got  types.String
		want string
	}{
		"SourceContent":     {m.SourceContent, "10.0.0.0/16"},
		"SourceType":        {m.SourceType, "ip"},
		"SourceName":        {m.SourceName, "src-asset"},
		"DestContent":       {m.DestContent, "0.0.0.0/0"},
		"DestType":          {m.DestType, "ip"},
		"DestName":          {m.DestName, "dst-asset"},
		"Protocol":          {m.Protocol, "TCP"},
		"Port":              {m.Port, "80"},
		"ParamTemplateId":   {m.ParamTemplateId, "pt-1"},
		"ParamTemplateName": {m.ParamTemplateName, "pt-name"},
		"RuleAction":        {m.RuleAction, "accept"},
		"Description":       {m.Description, "desc"},
		"EdgeId":            {m.EdgeId, "ALL"},
		"FwGroupId":         {m.FwGroupId, "ALL"},
		"BelongMemberId":    {m.BelongMemberId, "mem-1"},
		"IpVersion":         {m.IpVersion, "ipv4"},
		"AppId":             {m.AppId, "app-1"},
		"Uin":               {m.Uin, "uin-1"},
		"Nickname":          {m.Nickname, "nick"},
		"NodeName":          {m.NodeName, "dept"},
		"CreateTime":        {m.CreateTime, "2025-01-01 00:00:00"},
		"UpdateTime":        {m.UpdateTime, "2025-01-02 00:00:00"},
		"FwGroupName":       {m.FwGroupName, "fw-group-name"},
		"LastHitTime":       {m.LastHitTime, "2025-01-03 00:00:00"},
	}
	for name, c := range strChecks {
		if c.got.IsNull() || c.got.ValueString() != c.want {
			t.Errorf("%s = %v, want %q", name, c.got, c.want)
		}
	}

	intChecks := map[string]struct {
		got  types.Int64
		want int64
	}{
		"OrderIndex":  {m.OrderIndex, 5},
		"Sequence":    {m.Sequence, 5},
		"Invalid":     {m.Invalid, 0},
		"DnsParseCnt": {m.DnsParseCnt, 2},
		"HitCount":    {m.HitCount, 42},
	}
	for name, c := range intChecks {
		if c.got.IsNull() || c.got.ValueInt64() != c.want {
			t.Errorf("%s = %v, want %d", name, c.got, c.want)
		}
	}

	// A response without BelongMember must keep belong_member_id and null out
	// the member-only computed fields.
	sparse := FwmVpcAclRule2ResourceModel{
		BelongMemberId: types.StringValue("mem-1"),
		OrderIndex:     types.Int64Value(9),
		Port:           types.StringValue("443"),
		Description:    types.StringValue("keep-me"),
	}
	refreshVpcAclRule2Model(&sparse, &fwmv20250611.VpcAclRuleResp{RuleId: ptrOf("rule-1")})
	if sparse.BelongMemberId.IsNull() || sparse.BelongMemberId.ValueString() != "mem-1" {
		t.Errorf("BelongMemberId = %v, want preserved mem-1", sparse.BelongMemberId)
	}
	if sparse.OrderIndex.IsNull() || sparse.OrderIndex.ValueInt64() != 9 {
		t.Errorf("OrderIndex = %v, want preserved 9", sparse.OrderIndex)
	}
	if sparse.Port.IsNull() || sparse.Port.ValueString() != "443" {
		t.Errorf("Port = %v, want preserved 443", sparse.Port)
	}
	if sparse.Description.IsNull() || sparse.Description.ValueString() != "keep-me" {
		t.Errorf("Description = %v, want preserved keep-me", sparse.Description)
	}
	for name, v := range map[string]types.String{
		"AppId":    sparse.AppId,
		"Uin":      sparse.Uin,
		"Nickname": sparse.Nickname,
		"NodeName": sparse.NodeName,
	} {
		if !v.IsNull() {
			t.Errorf("%s = %v, want null", name, v)
		}
	}

	// A nil item must be a no-op.
	untouched := FwmVpcAclRule2ResourceModel{Port: types.StringValue("8080")}
	refreshVpcAclRule2Model(&untouched, nil)
	if untouched.Port.ValueString() != "8080" {
		t.Errorf("refreshVpcAclRule2Model(nil) modified the model")
	}
}

// ---- resource lifecycle tests (mock cloud API via gomonkey) ----

// newVpcAclRule2ResourceWithClient builds a resource whose shared client is
// replaced by a patched accessor returning the mock API client.
func newVpcAclRule2ResourceWithClient(t *testing.T, patches *gomonkey.Patches, client *fwmv20250611.Client) *FwmVpcAclRule2Resource {
	t.Helper()
	shared := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(shared, "UseFwmV20250611Client", client)

	r := &FwmVpcAclRule2Resource{}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: shared}}, &resource.ConfigureResponse{})
	if r.Client() != shared {
		t.Fatal("Configure did not inject the shared client")
	}
	return r
}

// schemaTypes converts the resource schema into the tftypes used to build
// framework plan/state values.
func schemaTypes(t *testing.T, r *FwmVpcAclRule2Resource) (schema.Schema, tftypes.Object) {
	t.Helper()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}
	objType := resp.Schema.Type().TerraformType(context.Background())
	obj, ok := objType.(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", objType)
	}
	return resp.Schema, obj
}

func ptrOf(s string) *string { return &s }

func ptrOfInt64(i int64) *int64 { return &i }

// TestFwmVpcAclRule2_CreateReadDelete covers the full lifecycle against a
// mocked cloud client and asserts every leaf attribute is read back.
func TestFwmVpcAclRule2_CreateReadDelete(t *testing.T) {
	client := &fwmv20250611.Client{}
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var createCalls, describeCalls, deleteCalls int

	patches.ApplyMethodFunc(client, "CreateVpcAclRuleWithContext",
		func(_ context.Context, request *fwmv20250611.CreateVpcAclRuleRequest) (*fwmv20250611.CreateVpcAclRuleResponse, error) {
			createCalls++
			if request.GroupId == nil || *request.GroupId != "fwmrg-1" {
				t.Errorf("Create GroupId = %v, want fwmrg-1", request.GroupId)
			}
			if len(request.Rules) != 1 {
				t.Fatalf("Create Rules len = %d, want 1", len(request.Rules))
			}
			rule := request.Rules[0]
			if rule.SourceContent == nil || *rule.SourceContent != "10.0.0.0/16" {
				t.Errorf("Create SourceContent = %v", rule.SourceContent)
			}
			if rule.OrderIndex == nil || *rule.OrderIndex != 5 {
				t.Errorf("Create OrderIndex = %v, want 5", rule.OrderIndex)
			}
			if rule.EdgeId == nil || *rule.EdgeId != "ALL" {
				t.Errorf("Create EdgeId = %v", rule.EdgeId)
			}
			return &fwmv20250611.CreateVpcAclRuleResponse{
				Response: &fwmv20250611.CreateVpcAclRuleResponseParams{
					RuleIds:   []*string{ptrOf("rule-1")},
					RequestId: ptrOf("req-1"),
				},
			}, nil
		})

	patches.ApplyMethodFunc(client, "DescribeVpcAclRulesWithContext",
		func(_ context.Context, request *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
			describeCalls++
			if request.GroupId == nil || *request.GroupId != "fwmrg-1" {
				t.Errorf("Describe GroupId = %v, want fwmrg-1", request.GroupId)
			}
			return &fwmv20250611.DescribeVpcAclRulesResponse{
				Response: &fwmv20250611.DescribeVpcAclRulesResponseParams{
					Rules:         []*fwmv20250611.VpcAclRuleResp{buildVpcAclRule2Resp()},
					TotalCount:    ptrOfInt64(1),
					AllTotalCount: ptrOfInt64(1),
					RequestId:     ptrOf("req-2"),
				},
			}, nil
		})

	patches.ApplyMethodFunc(client, "DeleteVpcAclRuleWithContext",
		func(_ context.Context, request *fwmv20250611.DeleteVpcAclRuleRequest) (*fwmv20250611.DeleteVpcAclRuleResponse, error) {
			deleteCalls++
			if request.GroupId == nil || *request.GroupId != "fwmrg-1" {
				t.Errorf("Delete GroupId = %v, want fwmrg-1", request.GroupId)
			}
			if len(request.RuleIds) != 1 || request.RuleIds[0] == nil || *request.RuleIds[0] != "rule-1" {
				t.Errorf("Delete RuleIds = %v, want [rule-1]", request.RuleIds)
			}
			return &fwmv20250611.DeleteVpcAclRuleResponse{
				Response: &fwmv20250611.DeleteVpcAclRuleResponseParams{RequestId: ptrOf("req-3")},
			}, nil
		})

	r := newVpcAclRule2ResourceWithClient(t, patches, client)

	tfSchema, objType := schemaTypes(t, r)

	ctx := context.Background()
	plan := tfsdk.Plan{Schema: tfSchema, Raw: tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, nil),
		"group_id":            tftypes.NewValue(tftypes.String, "fwmrg-1"),
		"rule_id":             tftypes.NewValue(tftypes.String, nil),
		"source_content":      tftypes.NewValue(tftypes.String, "10.0.0.0/16"),
		"source_type":         tftypes.NewValue(tftypes.String, "ip"),
		"dest_content":        tftypes.NewValue(tftypes.String, "0.0.0.0/0"),
		"dest_type":           tftypes.NewValue(tftypes.String, "ip"),
		"protocol":            tftypes.NewValue(tftypes.String, "TCP"),
		"port":                tftypes.NewValue(tftypes.String, "80"),
		"rule_action":         tftypes.NewValue(tftypes.String, "accept"),
		"order_index":         tftypes.NewValue(tftypes.Number, 5),
		"edge_id":             tftypes.NewValue(tftypes.String, "ALL"),
		"fw_group_id":         tftypes.NewValue(tftypes.String, "ALL"),
		"description":         tftypes.NewValue(tftypes.String, "desc"),
		"param_template_id":   tftypes.NewValue(tftypes.String, nil),
		"belong_member_id":    tftypes.NewValue(tftypes.String, nil),
		"sequence":            tftypes.NewValue(tftypes.Number, nil),
		"ip_version":          tftypes.NewValue(tftypes.String, nil),
		"source_name":         tftypes.NewValue(tftypes.String, nil),
		"dest_name":           tftypes.NewValue(tftypes.String, nil),
		"param_template_name": tftypes.NewValue(tftypes.String, nil),
		"invalid":             tftypes.NewValue(tftypes.Number, nil),
		"app_id":              tftypes.NewValue(tftypes.String, nil),
		"uin":                 tftypes.NewValue(tftypes.String, nil),
		"nickname":            tftypes.NewValue(tftypes.String, nil),
		"node_name":           tftypes.NewValue(tftypes.String, nil),
		"create_time":         tftypes.NewValue(tftypes.String, nil),
		"update_time":         tftypes.NewValue(tftypes.String, nil),
		"dns_parse_cnt":       tftypes.NewValue(tftypes.Number, nil),
		"fw_group_name":       tftypes.NewValue(tftypes.String, nil),
		"hit_count":           tftypes.NewValue(tftypes.Number, nil),
		"last_hit_time":       tftypes.NewValue(tftypes.String, nil),
	})}

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: tfSchema}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", createResp.Diagnostics)
	}
	if createCalls != 1 {
		t.Errorf("CreateVpcAclRule calls = %d, want 1", createCalls)
	}

	var created FwmVpcAclRule2ResourceModel
	createResp.Diagnostics.Append(createResp.State.Get(ctx, &created)...)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("state read failed: %v", createResp.Diagnostics)
	}
	assertCreated := map[string]struct{ got, want string }{
		"Id":              {created.Id.ValueString(), "fwmrg-1#rule-1"},
		"GroupId":         {created.GroupId.ValueString(), "fwmrg-1"},
		"RuleId":          {created.RuleId.ValueString(), "rule-1"},
		"SourceContent":   {created.SourceContent.ValueString(), "10.0.0.0/16"},
		"SourceName":      {created.SourceName.ValueString(), "src-asset"},
		"DestName":        {created.DestName.ValueString(), "dst-asset"},
		"ParamTemplate":   {"", ""},
		"IpVersion":       {created.IpVersion.ValueString(), "ipv4"},
		"AppId":           {created.AppId.ValueString(), "app-1"},
		"Uin":             {created.Uin.ValueString(), "uin-1"},
		"Nickname":        {created.Nickname.ValueString(), "nick"},
		"NodeName":        {created.NodeName.ValueString(), "dept"},
		"CreateTime":      {created.CreateTime.ValueString(), "2025-01-01 00:00:00"},
		"UpdateTime":      {created.UpdateTime.ValueString(), "2025-01-02 00:00:00"},
		"FwGroupName":     {created.FwGroupName.ValueString(), "fw-group-name"},
		"LastHitTime":     {created.LastHitTime.ValueString(), "2025-01-03 00:00:00"},
		"BelongMemberId":  {created.BelongMemberId.ValueString(), "mem-1"},
		"ParamTemplateId": {created.ParamTemplateId.ValueString(), "pt-1"},
		"ParamTemplateNm": {created.ParamTemplateName.ValueString(), "pt-name"},
	}
	for name, c := range assertCreated {
		if c.got != c.want {
			t.Errorf("created %s = %q, want %q", name, c.got, c.want)
		}
	}
	if created.OrderIndex.ValueInt64() != 5 || created.Sequence.ValueInt64() != 5 {
		t.Errorf("created OrderIndex/Sequence = %v/%v, want 5/5", created.OrderIndex, created.Sequence)
	}
	if created.HitCount.ValueInt64() != 42 {
		t.Errorf("created HitCount = %v, want 42", created.HitCount)
	}
	if created.Invalid.ValueInt64() != 0 || created.Invalid.IsNull() {
		t.Errorf("created Invalid = %v, want 0", created.Invalid)
	}
	if created.DnsParseCnt.ValueInt64() != 2 {
		t.Errorf("created DnsParseCnt = %v, want 2", created.DnsParseCnt)
	}

	// Read: refreshes state from the cloud side.
	state := createResp.State
	readResp := &resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
	}
	if describeCalls != 2 {
		t.Errorf("DescribeVpcAclRules calls = %d, want 2 (create read-back + read)", describeCalls)
	}

	// Delete: only the referenced rule id is sent.
	deleteResp := &resource.DeleteResponse{State: state}
	r.Delete(ctx, resource.DeleteRequest{State: state}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete returned diagnostics: %v", deleteResp.Diagnostics)
	}
	if deleteCalls != 1 {
		t.Errorf("DeleteVpcAclRule calls = %d, want 1", deleteCalls)
	}
}

// TestFwmVpcAclRule2_Update verifies the update path sends the full rule and
// preserves optional values the user removed from the configuration.
func TestFwmVpcAclRule2_Update(t *testing.T) {
	client := &fwmv20250611.Client{}
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var modified *fwmv20250611.VpcAclRule
	patches.ApplyMethodFunc(client, "ModifyVpcAclRuleWithContext",
		func(_ context.Context, request *fwmv20250611.ModifyVpcAclRuleRequest) (*fwmv20250611.ModifyVpcAclRuleResponse, error) {
			if request.GroupId == nil || *request.GroupId != "fwmrg-1" {
				t.Errorf("Modify GroupId = %v, want fwmrg-1", request.GroupId)
			}
			modified = request.Rule
			return &fwmv20250611.ModifyVpcAclRuleResponse{
				Response: &fwmv20250611.ModifyVpcAclRuleResponseParams{RequestId: ptrOf("req-4")},
			}, nil
		})
	patches.ApplyMethodFunc(client, "DescribeVpcAclRulesWithContext",
		func(_ context.Context, _ *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
			return &fwmv20250611.DescribeVpcAclRulesResponse{
				Response: &fwmv20250611.DescribeVpcAclRulesResponseParams{
					Rules: []*fwmv20250611.VpcAclRuleResp{buildVpcAclRule2Resp()},
				},
			}, nil
		})

	r := newVpcAclRule2ResourceWithClient(t, patches, client)

	tfSchema, objType := schemaTypes(t, r)
	ctx := context.Background()

	base := func(port, desc, ptID, memberID tftypes.Value) tftypes.Value {
		return tftypes.NewValue(objType, map[string]tftypes.Value{
			"id":                  tftypes.NewValue(tftypes.String, "fwmrg-1#rule-1"),
			"group_id":            tftypes.NewValue(tftypes.String, "fwmrg-1"),
			"rule_id":             tftypes.NewValue(tftypes.String, "rule-1"),
			"source_content":      tftypes.NewValue(tftypes.String, "10.0.0.0/16"),
			"source_type":         tftypes.NewValue(tftypes.String, "ip"),
			"dest_content":        tftypes.NewValue(tftypes.String, "0.0.0.0/0"),
			"dest_type":           tftypes.NewValue(tftypes.String, "ip"),
			"protocol":            tftypes.NewValue(tftypes.String, "TCP"),
			"port":                port,
			"rule_action":         tftypes.NewValue(tftypes.String, "drop"),
			"order_index":         tftypes.NewValue(tftypes.Number, 2),
			"edge_id":             tftypes.NewValue(tftypes.String, "ALL"),
			"fw_group_id":         tftypes.NewValue(tftypes.String, "ALL"),
			"description":         desc,
			"param_template_id":   ptID,
			"belong_member_id":    memberID,
			"sequence":            tftypes.NewValue(tftypes.Number, 2),
			"ip_version":          tftypes.NewValue(tftypes.String, "ipv4"),
			"source_name":         tftypes.NewValue(tftypes.String, "src-asset"),
			"dest_name":           tftypes.NewValue(tftypes.String, "dst-asset"),
			"param_template_name": tftypes.NewValue(tftypes.String, "pt-name"),
			"invalid":             tftypes.NewValue(tftypes.Number, 0),
			"app_id":              tftypes.NewValue(tftypes.String, "app-1"),
			"uin":                 tftypes.NewValue(tftypes.String, "uin-1"),
			"nickname":            tftypes.NewValue(tftypes.String, "nick"),
			"node_name":           tftypes.NewValue(tftypes.String, "dept"),
			"create_time":         tftypes.NewValue(tftypes.String, "2025-01-01 00:00:00"),
			"update_time":         tftypes.NewValue(tftypes.String, "2025-01-02 00:00:00"),
			"dns_parse_cnt":       tftypes.NewValue(tftypes.Number, 2),
			"fw_group_name":       tftypes.NewValue(tftypes.String, "fw-group-name"),
			"hit_count":           tftypes.NewValue(tftypes.Number, 42),
			"last_hit_time":       tftypes.NewValue(tftypes.String, "2025-01-03 00:00:00"),
		})
	}

	state := tfsdk.State{Schema: tfSchema, Raw: base(
		tftypes.NewValue(tftypes.String, "80"),
		tftypes.NewValue(tftypes.String, "old-desc"),
		tftypes.NewValue(tftypes.String, "pt-old"),
		tftypes.NewValue(tftypes.String, "mem-1"),
	)}
	plan := tfsdk.Plan{Schema: tfSchema, Raw: base(
		tftypes.NewValue(tftypes.String, nil),
		tftypes.NewValue(tftypes.String, nil),
		tftypes.NewValue(tftypes.String, nil),
		tftypes.NewValue(tftypes.String, nil),
	)}

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: tfSchema}}
	r.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", updateResp.Diagnostics)
	}
	if modified == nil {
		t.Fatal("ModifyVpcAclRule was not called")
	}
	if modified.RuleId == nil || *modified.RuleId != "rule-1" {
		t.Errorf("Modify RuleId = %v, want rule-1", modified.RuleId)
	}
	if modified.RuleAction == nil || *modified.RuleAction != "drop" {
		t.Errorf("Modify RuleAction = %v, want drop", modified.RuleAction)
	}
	if modified.OrderIndex == nil || *modified.OrderIndex != 2 {
		t.Errorf("Modify OrderIndex = %v, want 2", modified.OrderIndex)
	}
	// Optional attributes removed from the configuration keep their prior value
	// because the clear-vs-keep semantics of this API is not documented.
	if modified.Port == nil || *modified.Port != "80" {
		t.Errorf("Modify Port = %v, want the previous value 80", modified.Port)
	}
	if modified.Description == nil || *modified.Description != "old-desc" {
		t.Errorf("Modify Description = %v, want the previous value old-desc", modified.Description)
	}
	if modified.ParamTemplateId == nil || *modified.ParamTemplateId != "pt-old" {
		t.Errorf("Modify ParamTemplateId = %v, want the previous value pt-old", modified.ParamTemplateId)
	}
	if modified.BelongMemberId == nil || *modified.BelongMemberId != "mem-1" {
		t.Errorf("Modify BelongMemberId = %v, want the previous value mem-1", modified.BelongMemberId)
	}

	var updated FwmVpcAclRule2ResourceModel
	updateResp.Diagnostics.Append(updateResp.State.Get(ctx, &updated)...)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("state read failed: %v", updateResp.Diagnostics)
	}
	if updated.Id.ValueString() != "fwmrg-1#rule-1" {
		t.Errorf("updated Id = %q, want fwmrg-1#rule-1", updated.Id.ValueString())
	}
	if updated.RuleAction.ValueString() != "accept" {
		t.Errorf("updated RuleAction = %q, want the read-back value accept", updated.RuleAction.ValueString())
	}
}

// TestFwmVpcAclRule2_Read_NotFound verifies a rule absent from the response
// removes the resource from state.
func TestFwmVpcAclRule2_Read_NotFound(t *testing.T) {
	client := &fwmv20250611.Client{}
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(client, "DescribeVpcAclRulesWithContext",
		func(_ context.Context, _ *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
			return &fwmv20250611.DescribeVpcAclRulesResponse{
				Response: &fwmv20250611.DescribeVpcAclRulesResponseParams{
					Rules: []*fwmv20250611.VpcAclRuleResp{},
				},
			}, nil
		})

	r := newVpcAclRule2ResourceWithClient(t, patches, client)

	tfSchema, objType := schemaTypes(t, r)
	ctx := context.Background()

	state := tfsdk.State{Schema: tfSchema, Raw: tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, "fwmrg-1#rule-1"),
		"group_id":            tftypes.NewValue(tftypes.String, "fwmrg-1"),
		"rule_id":             tftypes.NewValue(tftypes.String, "rule-1"),
		"source_content":      tftypes.NewValue(tftypes.String, "10.0.0.0/16"),
		"source_type":         tftypes.NewValue(tftypes.String, "ip"),
		"dest_content":        tftypes.NewValue(tftypes.String, "0.0.0.0/0"),
		"dest_type":           tftypes.NewValue(tftypes.String, "ip"),
		"protocol":            tftypes.NewValue(tftypes.String, "TCP"),
		"port":                tftypes.NewValue(tftypes.String, nil),
		"rule_action":         tftypes.NewValue(tftypes.String, "accept"),
		"order_index":         tftypes.NewValue(tftypes.Number, 1),
		"edge_id":             tftypes.NewValue(tftypes.String, "ALL"),
		"fw_group_id":         tftypes.NewValue(tftypes.String, "ALL"),
		"description":         tftypes.NewValue(tftypes.String, nil),
		"param_template_id":   tftypes.NewValue(tftypes.String, nil),
		"belong_member_id":    tftypes.NewValue(tftypes.String, nil),
		"sequence":            tftypes.NewValue(tftypes.Number, 1),
		"ip_version":          tftypes.NewValue(tftypes.String, "ipv4"),
		"source_name":         tftypes.NewValue(tftypes.String, nil),
		"dest_name":           tftypes.NewValue(tftypes.String, nil),
		"param_template_name": tftypes.NewValue(tftypes.String, nil),
		"invalid":             tftypes.NewValue(tftypes.Number, 0),
		"app_id":              tftypes.NewValue(tftypes.String, nil),
		"uin":                 tftypes.NewValue(tftypes.String, nil),
		"nickname":            tftypes.NewValue(tftypes.String, nil),
		"node_name":           tftypes.NewValue(tftypes.String, nil),
		"create_time":         tftypes.NewValue(tftypes.String, nil),
		"update_time":         tftypes.NewValue(tftypes.String, nil),
		"dns_parse_cnt":       tftypes.NewValue(tftypes.Number, 0),
		"fw_group_name":       tftypes.NewValue(tftypes.String, nil),
		"hit_count":           tftypes.NewValue(tftypes.Number, 0),
		"last_hit_time":       tftypes.NewValue(tftypes.String, nil),
	})}

	readResp := &resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read with an empty rule list returned diagnostics: %v", readResp.Diagnostics)
	}
	if !readResp.State.Raw.IsNull() {
		t.Errorf("Read did not remove the resource from state: %v", readResp.State.Raw)
	}
}

// TestFwmVpcAclRule2_Errors covers the failure branches of Create, Read,
// Update and Delete.
func TestFwmVpcAclRule2_Errors(t *testing.T) {
	ctx := context.Background()

	t.Run("create-empty-rule-ids", func(t *testing.T) {
		client := &fwmv20250611.Client{}
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodFunc(client, "CreateVpcAclRuleWithContext",
			func(_ context.Context, _ *fwmv20250611.CreateVpcAclRuleRequest) (*fwmv20250611.CreateVpcAclRuleResponse, error) {
				return &fwmv20250611.CreateVpcAclRuleResponse{
					Response: &fwmv20250611.CreateVpcAclRuleResponseParams{RuleIds: []*string{}},
				}, nil
			})

		r := newVpcAclRule2ResourceWithClient(t, patches, client)

		tfSchema, objType := schemaTypes(t, r)
		plan := tfsdk.Plan{Schema: tfSchema, Raw: minimalPlanValue(objType)}
		createResp := &resource.CreateResponse{State: tfsdk.State{Schema: tfSchema}}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, createResp)
		if !createResp.Diagnostics.HasError() {
			t.Fatal("Create with an empty RuleIds must report an error")
		}
	})

	t.Run("create-multiple-rule-ids", func(t *testing.T) {
		client := &fwmv20250611.Client{}
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodFunc(client, "CreateVpcAclRuleWithContext",
			func(_ context.Context, _ *fwmv20250611.CreateVpcAclRuleRequest) (*fwmv20250611.CreateVpcAclRuleResponse, error) {
				return &fwmv20250611.CreateVpcAclRuleResponse{
					Response: &fwmv20250611.CreateVpcAclRuleResponseParams{
						RuleIds: []*string{ptrOf("rule-1"), ptrOf("rule-2")},
					},
				}, nil
			})

		r := newVpcAclRule2ResourceWithClient(t, patches, client)

		tfSchema, objType := schemaTypes(t, r)
		plan := tfsdk.Plan{Schema: tfSchema, Raw: minimalPlanValue(objType)}
		createResp := &resource.CreateResponse{State: tfsdk.State{Schema: tfSchema}}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, createResp)
		if !createResp.Diagnostics.HasError() {
			t.Fatal("Create returning several rule ids must report an error")
		}
	})

	t.Run("create-api-error", func(t *testing.T) {
		client := &fwmv20250611.Client{}
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodFunc(client, "CreateVpcAclRuleWithContext",
			func(_ context.Context, _ *fwmv20250611.CreateVpcAclRuleRequest) (*fwmv20250611.CreateVpcAclRuleResponse, error) {
				return nil, errors.New("boom")
			})

		r := newVpcAclRule2ResourceWithClient(t, patches, client)

		tfSchema, objType := schemaTypes(t, r)
		plan := tfsdk.Plan{Schema: tfSchema, Raw: minimalPlanValue(objType)}
		createResp := &resource.CreateResponse{State: tfsdk.State{Schema: tfSchema}}
		r.Create(ctx, resource.CreateRequest{Plan: plan}, createResp)
		if !createResp.Diagnostics.HasError() {
			t.Fatal("Create must report an API error")
		}
	})

	t.Run("read-nil-response-is-an-error", func(t *testing.T) {
		client := &fwmv20250611.Client{}
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodFunc(client, "DescribeVpcAclRulesWithContext",
			func(_ context.Context, _ *fwmv20250611.DescribeVpcAclRulesRequest) (*fwmv20250611.DescribeVpcAclRulesResponse, error) {
				return &fwmv20250611.DescribeVpcAclRulesResponse{}, nil
			})

		r := newVpcAclRule2ResourceWithClient(t, patches, client)

		tfSchema, objType := schemaTypes(t, r)
		state := tfsdk.State{Schema: tfSchema, Raw: minimalPlanValue(objType)}
		readResp := &resource.ReadResponse{State: state}
		r.Read(ctx, resource.ReadRequest{State: state}, readResp)
		if !readResp.Diagnostics.HasError() {
			t.Fatal("Read with a nil Response must report an error")
		}
	})

	t.Run("read-broken-id", func(t *testing.T) {
		r := &FwmVpcAclRule2Resource{}
		tfSchema, objType := schemaTypes(t, r)
		raw := minimalPlanValue(objType)
		state := tfsdk.State{Schema: tfSchema, Raw: raw}
		var model FwmVpcAclRule2ResourceModel
		if diags := state.Get(ctx, &model); diags.HasError() {
			t.Fatalf("state get failed: %v", diags)
		}
		model.Id = types.StringValue("broken-id")
		newState := tfsdk.State{Schema: tfSchema}
		if diags := newState.Set(ctx, &model); diags.HasError() {
			t.Fatalf("state set failed: %v", diags)
		}
		readResp := &resource.ReadResponse{State: newState}
		r.Read(ctx, resource.ReadRequest{State: newState}, readResp)
		if !readResp.Diagnostics.HasError() {
			t.Fatal("Read with a broken id must report an error")
		}
	})

	t.Run("update-api-error", func(t *testing.T) {
		client := &fwmv20250611.Client{}
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodFunc(client, "ModifyVpcAclRuleWithContext",
			func(_ context.Context, _ *fwmv20250611.ModifyVpcAclRuleRequest) (*fwmv20250611.ModifyVpcAclRuleResponse, error) {
				return nil, errors.New("boom")
			})

		r := newVpcAclRule2ResourceWithClient(t, patches, client)

		tfSchema, objType := schemaTypes(t, r)
		state := tfsdk.State{Schema: tfSchema, Raw: minimalPlanValue(objType)}
		var model FwmVpcAclRule2ResourceModel
		if diags := state.Get(ctx, &model); diags.HasError() {
			t.Fatalf("state get failed: %v", diags)
		}
		model.Id = types.StringValue("fwmrg-1#rule-1")
		var withId tfsdk.State
		withId.Schema = tfSchema
		if diags := withId.Set(ctx, &model); diags.HasError() {
			t.Fatalf("state set failed: %v", diags)
		}
		plan := tfsdk.Plan{Schema: tfSchema, Raw: withId.Raw}
		updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: tfSchema}}
		r.Update(ctx, resource.UpdateRequest{Plan: plan, State: withId}, updateResp)
		if !updateResp.Diagnostics.HasError() {
			t.Fatal("Update must report an API error")
		}
	})

	t.Run("delete-not-found-is-swallowed", func(t *testing.T) {
		client := &fwmv20250611.Client{}
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodFunc(client, "DeleteVpcAclRuleWithContext",
			func(_ context.Context, _ *fwmv20250611.DeleteVpcAclRuleRequest) (*fwmv20250611.DeleteVpcAclRuleResponse, error) {
				return nil, sdkErrors.NewTencentCloudSDKError("ResourceNotFound", "gone", "req-x")
			})

		r := newVpcAclRule2ResourceWithClient(t, patches, client)

		tfSchema, objType := schemaTypes(t, r)
		state := tfsdk.State{Schema: tfSchema, Raw: minimalPlanValue(objType)}
		var model FwmVpcAclRule2ResourceModel
		if diags := state.Get(ctx, &model); diags.HasError() {
			t.Fatalf("state get failed: %v", diags)
		}
		model.Id = types.StringValue("fwmrg-1#rule-1")
		var withId tfsdk.State
		withId.Schema = tfSchema
		if diags := withId.Set(ctx, &model); diags.HasError() {
			t.Fatalf("state set failed: %v", diags)
		}
		deleteResp := &resource.DeleteResponse{State: withId}
		r.Delete(ctx, resource.DeleteRequest{State: withId}, deleteResp)
		if deleteResp.Diagnostics.HasError() {
			t.Fatalf("Delete of an absent rule must succeed, got %v", deleteResp.Diagnostics)
		}
	})

	t.Run("delete-permission-error-surfaces", func(t *testing.T) {
		client := &fwmv20250611.Client{}
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodFunc(client, "DeleteVpcAclRuleWithContext",
			func(_ context.Context, _ *fwmv20250611.DeleteVpcAclRuleRequest) (*fwmv20250611.DeleteVpcAclRuleResponse, error) {
				return nil, sdkErrors.NewTencentCloudSDKError("UnauthorizedOperation", "no permission", "req-y")
			})

		r := newVpcAclRule2ResourceWithClient(t, patches, client)

		tfSchema, objType := schemaTypes(t, r)
		state := tfsdk.State{Schema: tfSchema, Raw: minimalPlanValue(objType)}
		var model FwmVpcAclRule2ResourceModel
		if diags := state.Get(ctx, &model); diags.HasError() {
			t.Fatalf("state get failed: %v", diags)
		}
		model.Id = types.StringValue("fwmrg-1#rule-1")
		var withId tfsdk.State
		withId.Schema = tfSchema
		if diags := withId.Set(ctx, &model); diags.HasError() {
			t.Fatalf("state set failed: %v", diags)
		}
		deleteResp := &resource.DeleteResponse{State: withId}
		r.Delete(ctx, resource.DeleteRequest{State: withId}, deleteResp)
		if !deleteResp.Diagnostics.HasError() {
			t.Fatal("Delete with a permission error must report an error")
		}
	})

	t.Run("delete-broken-id", func(t *testing.T) {
		r := &FwmVpcAclRule2Resource{}
		tfSchema, objType := schemaTypes(t, r)
		state := tfsdk.State{Schema: tfSchema, Raw: minimalPlanValue(objType)}
		var model FwmVpcAclRule2ResourceModel
		if diags := state.Get(ctx, &model); diags.HasError() {
			t.Fatalf("state get failed: %v", diags)
		}
		model.Id = types.StringValue("nope")
		var withId tfsdk.State
		withId.Schema = tfSchema
		if diags := withId.Set(ctx, &model); diags.HasError() {
			t.Fatalf("state set failed: %v", diags)
		}
		deleteResp := &resource.DeleteResponse{State: withId}
		r.Delete(ctx, resource.DeleteRequest{State: withId}, deleteResp)
		if !deleteResp.Diagnostics.HasError() {
			t.Fatal("Delete with a broken id must report an error")
		}
	})
}

// TestFwmVpcAclRule2_ImportState covers the composite import form and the
// broken import branch.
func TestFwmVpcAclRule2_ImportState(t *testing.T) {
	r := &FwmVpcAclRule2Resource{}
	tfSchema, objType := schemaTypes(t, r)
	ctx := context.Background()

	// An import starts from an empty configuration: every attribute is null.
	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: tfSchema, Raw: tftypes.NewValue(objType, nil)}}
	r.ImportState(ctx, resource.ImportStateRequest{ID: "fwmrg-1#rule-1"}, importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState returned diagnostics: %v", importResp.Diagnostics)
	}

	var model FwmVpcAclRule2ResourceModel
	importResp.Diagnostics.Append(importResp.State.Get(ctx, &model)...)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("state read failed: %v", importResp.Diagnostics)
	}
	if model.Id.ValueString() != "fwmrg-1#rule-1" {
		t.Errorf("imported Id = %q, want fwmrg-1#rule-1", model.Id.ValueString())
	}
	if model.GroupId.ValueString() != "fwmrg-1" {
		t.Errorf("imported GroupId = %q, want fwmrg-1", model.GroupId.ValueString())
	}
	if model.RuleId.ValueString() != "rule-1" {
		t.Errorf("imported RuleId = %q, want rule-1", model.RuleId.ValueString())
	}

	badResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: tfSchema, Raw: tftypes.NewValue(objType, nil)}}
	r.ImportState(ctx, resource.ImportStateRequest{ID: "no-separator"}, badResp)
	if !badResp.Diagnostics.HasError() {
		t.Fatal("ImportState with a broken id must report an error")
	}
}

// minimalPlanValue builds a plan value with only the identities left null and
// all required inputs filled, used by the error-branch tests.
func minimalPlanValue(objType tftypes.Object) tftypes.Value {
	return tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, nil),
		"group_id":            tftypes.NewValue(tftypes.String, "fwmrg-1"),
		"rule_id":             tftypes.NewValue(tftypes.String, nil),
		"source_content":      tftypes.NewValue(tftypes.String, "10.0.0.0/16"),
		"source_type":         tftypes.NewValue(tftypes.String, "ip"),
		"dest_content":        tftypes.NewValue(tftypes.String, "0.0.0.0/0"),
		"dest_type":           tftypes.NewValue(tftypes.String, "ip"),
		"protocol":            tftypes.NewValue(tftypes.String, "TCP"),
		"port":                tftypes.NewValue(tftypes.String, nil),
		"rule_action":         tftypes.NewValue(tftypes.String, "accept"),
		"order_index":         tftypes.NewValue(tftypes.Number, 1),
		"edge_id":             tftypes.NewValue(tftypes.String, "ALL"),
		"fw_group_id":         tftypes.NewValue(tftypes.String, "ALL"),
		"description":         tftypes.NewValue(tftypes.String, nil),
		"param_template_id":   tftypes.NewValue(tftypes.String, nil),
		"belong_member_id":    tftypes.NewValue(tftypes.String, nil),
		"sequence":            tftypes.NewValue(tftypes.Number, nil),
		"ip_version":          tftypes.NewValue(tftypes.String, nil),
		"source_name":         tftypes.NewValue(tftypes.String, nil),
		"dest_name":           tftypes.NewValue(tftypes.String, nil),
		"param_template_name": tftypes.NewValue(tftypes.String, nil),
		"invalid":             tftypes.NewValue(tftypes.Number, nil),
		"app_id":              tftypes.NewValue(tftypes.String, nil),
		"uin":                 tftypes.NewValue(tftypes.String, nil),
		"nickname":            tftypes.NewValue(tftypes.String, nil),
		"node_name":           tftypes.NewValue(tftypes.String, nil),
		"create_time":         tftypes.NewValue(tftypes.String, nil),
		"update_time":         tftypes.NewValue(tftypes.String, nil),
		"dns_parse_cnt":       tftypes.NewValue(tftypes.Number, nil),
		"fw_group_name":       tftypes.NewValue(tftypes.String, nil),
		"hit_count":           tftypes.NewValue(tftypes.Number, nil),
		"last_hit_time":       tftypes.NewValue(tftypes.String, nil),
	})
}
