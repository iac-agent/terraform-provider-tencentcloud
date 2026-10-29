package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

func newFwmNatAclRuleTestResource(t *testing.T) *fwm.FwmNatAclRuleResource {
	t.Helper()

	res := fwm.NewFwmNatAclRuleResource()

	ctx := context.Background()
	schemaResp := &resource.SchemaResponse{}
	res.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	configureResp := &resource.ConfigureResponse{}
	res.Configure(ctx, resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: newMockMetaFwmAccountGroup().client},
	}, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	return res.(*fwm.FwmNatAclRuleResource)
}

// fwmNatAclRuleStateModel reads the whole state back into a typed struct.
func fwmNatAclRuleStateModel(t *testing.T, state tfsdk.State) fwm.FwmNatAclRuleResourceModel {
	t.Helper()

	var model fwm.FwmNatAclRuleResourceModel
	diags := state.Get(context.Background(), &model)
	assert.False(t, diags.HasError(), diags.Errors())
	return model
}

// fwmNatAclRuleRawValue builds a tftypes.Value for the resource schema where
// every attribute is null except the ones provided in values.
func fwmNatAclRuleRawValue(t *testing.T, res *fwm.FwmNatAclRuleResource, values map[string]interface{}) tftypes.Value {
	t.Helper()

	ctx := context.Background()

	schemaResp := &resource.SchemaResponse{}
	res.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	resType := schemaResp.Schema.Type()
	objType, ok := resType.TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("unexpected schema type %T", resType.TerraformType(ctx))
	}

	attrs := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		attrs[name] = tftypes.NewValue(attrType, nil)
	}
	for name, v := range values {
		attrType, ok := objType.AttributeTypes[name]
		if !ok {
			t.Fatalf("unknown attribute %q", name)
		}
		attrs[name] = tftypes.NewValue(attrType, v)
	}

	return tftypes.NewValue(objType, attrs)
}

func fwmNatAclRulePlan(t *testing.T, res *fwm.FwmNatAclRuleResource, values map[string]interface{}) tfsdk.Plan {
	t.Helper()

	schemaResp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

	return tfsdk.Plan{
		Raw:    fwmNatAclRuleRawValue(t, res, values),
		Schema: schemaResp.Schema,
	}
}

func fwmNatAclRuleState(t *testing.T, res *fwm.FwmNatAclRuleResource, values map[string]interface{}) tfsdk.State {
	t.Helper()

	schemaResp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

	return tfsdk.State{
		Raw:    fwmNatAclRuleRawValue(t, res, values),
		Schema: schemaResp.Schema,
	}
}

func fwmNatAclRuleResponse(ruleId string) *fwmv20250611.DescribeNatAclRulesResponse {
	sequence := int64(1)
	direction := int64(0)
	allTotalCount := int64(7)
	hitCount := int64(0)

	resp := fwmv20250611.NewDescribeNatAclRulesResponse()
	resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
		RequestId:     ptrStringFAG("fake-request-id"),
		AllTotalCount: &allTotalCount,
		Rules: []*fwmv20250611.NatAclRuleResp{
			{
				RuleId:        ptrStringFAG(ruleId),
				Sequence:      &sequence,
				Direction:     &direction,
				SourceContent: ptrStringFAG("10.0.0.0/8"),
				SourceType:    ptrStringFAG("ip"),
				SourceName:    ptrStringFAG("src-asset"),
				TargetContent: ptrStringFAG("192.168.1.0/24"),
				TargetType:    ptrStringFAG("ip"),
				TargetName:    ptrStringFAG("dst-asset"),
				Port:          ptrStringFAG("80,443"),
				Protocol:      ptrStringFAG("TCP"),
				RuleAction:    ptrStringFAG("accept"),
				Description:   ptrStringFAG("tf-example"),
				Scope:         ptrStringFAG("ALL"),
				ScopeDesc:     ptrStringFAG("global"),
				FwInsId:       ptrStringFAG("cfwnat-1a2b3c4d"),
				Invalid:       ptrInt64FAG(0),
				BelongMember: &fwmv20250611.MemberInfo{
					AppId:    ptrStringFAG("1300000000"),
					Uin:      ptrStringFAG("100000000001"),
					Nickname: ptrStringFAG("member-nick"),
					MemberId: ptrStringFAG("mem-1a2b3c4d"),
					NodeName: ptrStringFAG("root"),
				},
				CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				UpdateTime:  ptrStringFAG("2026-10-22 11:00:00"),
				DnsParseCnt: ptrInt64FAG(0),
				HitCount:    &hitCount,
				LastHitTime: ptrStringFAG("-"),
			},
		},
	}
	return resp
}

func ptrInt64FAG(v int64) *int64 {
	return &v
}

func fwmNatAclRuleBasePlanValues() map[string]interface{} {
	return map[string]interface{}{
		"group_id":       "acg-1a2b3c4d",
		"source_content": "10.0.0.0/8",
		"source_type":    "ip",
		"target_content": "192.168.1.0/24",
		"target_type":    "ip",
		"protocol":       "TCP",
		"rule_action":    "accept",
		"order_index":    int64(1),
		"scope":          "ALL",
		"port":           "80,443",
		"description":    "tf-example",
		"direction":      int64(0),
	}
}

func TestFwmNatAclRule_Create_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)
		assert.Len(t, request.Rules, 1)
		assert.Equal(t, "ip", *request.Rules[0].SourceType)
		assert.Equal(t, "TCP", *request.Rules[0].Protocol)

		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RequestId: ptrStringFAG("fake-request-id-create"),
			RuleIds:   []*string{ptrStringFAG("aclr-9x8y7z6w")},
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)
		return fwmNatAclRuleResponse("aclr-9x8y7z6w"), nil
	})

	res := newFwmNatAclRuleTestResource(t)
	resp := &resource.CreateResponse{
		State: fwmNatAclRuleState(t, res, nil),
	}
	res.Create(context.Background(), resource.CreateRequest{
		Plan: fwmNatAclRulePlan(t, res, fwmNatAclRuleBasePlanValues()),
	}, resp)

	assert.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())

	gotState := fwmNatAclRuleStateModel(t, resp.State)
	assert.Equal(t, "acg-1a2b3c4d"+tccommon.FILED_SP+"aclr-9x8y7z6w", gotState.Id.ValueString())
	assert.Equal(t, "aclr-9x8y7z6w", gotState.RuleId.ValueString())
	assert.Equal(t, int64(1), gotState.OrderIndex.ValueInt64())
	assert.Equal(t, "global", gotState.ScopeDesc.ValueString())
	assert.Equal(t, int64(7), gotState.AllTotalCount.ValueInt64())
}

func TestFwmNatAclRule_Create_EmptyRuleId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RequestId: ptrStringFAG("fake-request-id-create"),
			RuleIds:   []*string{},
		}
		return resp, nil
	})

	res := newFwmNatAclRuleTestResource(t)
	resp := &resource.CreateResponse{
		State: fwmNatAclRuleState(t, res, nil),
	}
	res.Create(context.Background(), resource.CreateRequest{
		Plan: fwmNatAclRulePlan(t, res, fwmNatAclRuleBasePlanValues()),
	}, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Empty(t, fwmNatAclRuleStateModel(t, resp.State).Id.ValueString())
}

func TestFwmNatAclRule_Create_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		return nil, assert.AnError
	})

	res := newFwmNatAclRuleTestResource(t)
	resp := &resource.CreateResponse{
		State: fwmNatAclRuleState(t, res, nil),
	}
	res.Create(context.Background(), resource.CreateRequest{
		Plan: fwmNatAclRulePlan(t, res, fwmNatAclRuleBasePlanValues()),
	}, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

func TestFwmNatAclRule_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "RuleId", *request.Filters[0].Name)
		assert.Equal(t, "aclr-9x8y7z6w", *request.Filters[0].Values[0])
		return fwmNatAclRuleResponse("aclr-9x8y7z6w"), nil
	})

	res := newFwmNatAclRuleTestResource(t)
	stateValues := fwmNatAclRuleBasePlanValues()
	stateValues["id"] = "acg-1a2b3c4d" + tccommon.FILED_SP + "aclr-9x8y7z6w"
	stateValues["rule_id"] = "aclr-9x8y7z6w"

	resp := &resource.ReadResponse{
		State: fwmNatAclRuleState(t, res, stateValues),
	}
	res.Read(context.Background(), resource.ReadRequest{
		State: fwmNatAclRuleState(t, res, stateValues),
	}, resp)

	assert.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())

	gotState := fwmNatAclRuleStateModel(t, resp.State)
	assert.Equal(t, "acg-1a2b3c4d"+tccommon.FILED_SP+"aclr-9x8y7z6w", gotState.Id.ValueString())
	assert.Equal(t, "aclr-9x8y7z6w", gotState.RuleId.ValueString())
	assert.Equal(t, "src-asset", gotState.SourceName.ValueString())
	assert.Equal(t, "dst-asset", gotState.TargetName.ValueString())
	assert.Equal(t, "cfwnat-1a2b3c4d", gotState.FwInsId.ValueString())
	assert.False(t, gotState.BelongMember.IsNull())
}

func TestFwmNatAclRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		total := int64(0)
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			Rules:      []*fwmv20250611.NatAclRuleResp{},
		}
		return resp, nil
	})

	res := newFwmNatAclRuleTestResource(t)
	stateValues := fwmNatAclRuleBasePlanValues()
	stateValues["id"] = "acg-1a2b3c4d" + tccommon.FILED_SP + "aclr-notexist"
	stateValues["rule_id"] = "aclr-notexist"

	resp := &resource.ReadResponse{
		State: fwmNatAclRuleState(t, res, stateValues),
	}
	res.Read(context.Background(), resource.ReadRequest{
		State: fwmNatAclRuleState(t, res, stateValues),
	}, resp)

	assert.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
	assert.True(t, resp.State.Raw.IsNull())
}

func TestFwmNatAclRule_Read_InvalidId(t *testing.T) {
	res := newFwmNatAclRuleTestResource(t)

	stateValues := fwmNatAclRuleBasePlanValues()
	stateValues["id"] = "broken-id"

	resp := &resource.ReadResponse{
		State: fwmNatAclRuleState(t, res, stateValues),
	}
	res.Read(context.Background(), resource.ReadRequest{
		State: fwmNatAclRuleState(t, res, stateValues),
	}, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

func TestFwmNatAclRule_Update_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "ModifyNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyNatAclRuleRequest) (*fwmv20250611.ModifyNatAclRuleResponse, error) {
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)
		assert.NotNil(t, request.Rule)
		assert.Equal(t, "aclr-9x8y7z6w", *request.Rule.RuleId)
		assert.Equal(t, "drop", *request.Rule.RuleAction)

		resp := fwmv20250611.NewModifyNatAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyNatAclRuleResponseParams{
			RequestId: ptrStringFAG("fake-request-id-update"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		return fwmNatAclRuleResponse("aclr-9x8y7z6w"), nil
	})

	res := newFwmNatAclRuleTestResource(t)

	stateValues := fwmNatAclRuleBasePlanValues()
	stateValues["id"] = "acg-1a2b3c4d" + tccommon.FILED_SP + "aclr-9x8y7z6w"
	stateValues["rule_id"] = "aclr-9x8y7z6w"

	planValues := fwmNatAclRuleBasePlanValues()
	planValues["rule_action"] = "drop"

	resp := &resource.UpdateResponse{
		State: fwmNatAclRuleState(t, res, stateValues),
	}
	res.Update(context.Background(), resource.UpdateRequest{
		Plan:  fwmNatAclRulePlan(t, res, planValues),
		State: fwmNatAclRuleState(t, res, stateValues),
	}, resp)

	assert.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
	assert.Equal(t, "acg-1a2b3c4d"+tccommon.FILED_SP+"aclr-9x8y7z6w", fwmNatAclRuleStateModel(t, resp.State).Id.ValueString())
}

func TestFwmNatAclRule_Update_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "ModifyNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyNatAclRuleRequest) (*fwmv20250611.ModifyNatAclRuleResponse, error) {
		return nil, assert.AnError
	})

	res := newFwmNatAclRuleTestResource(t)

	stateValues := fwmNatAclRuleBasePlanValues()
	stateValues["id"] = "acg-1a2b3c4d" + tccommon.FILED_SP + "aclr-9x8y7z6w"
	stateValues["rule_id"] = "aclr-9x8y7z6w"

	resp := &resource.UpdateResponse{
		State: fwmNatAclRuleState(t, res, stateValues),
	}
	res.Update(context.Background(), resource.UpdateRequest{
		Plan:  fwmNatAclRulePlan(t, res, fwmNatAclRuleBasePlanValues()),
		State: fwmNatAclRuleState(t, res, stateValues),
	}, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

func TestFwmNatAclRule_Delete_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DeleteNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteNatAclRuleRequest) (*fwmv20250611.DeleteNatAclRuleResponse, error) {
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)
		assert.Len(t, request.RuleIds, 1)
		assert.Equal(t, "aclr-9x8y7z6w", *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteNatAclRuleResponse()
		resp.Response = &fwmv20250611.DeleteNatAclRuleResponseParams{
			RequestId: ptrStringFAG("fake-request-id-delete"),
		}
		return resp, nil
	})

	res := newFwmNatAclRuleTestResource(t)
	stateValues := fwmNatAclRuleBasePlanValues()
	stateValues["id"] = "acg-1a2b3c4d" + tccommon.FILED_SP + "aclr-9x8y7z6w"
	stateValues["rule_id"] = "aclr-9x8y7z6w"

	resp := &resource.DeleteResponse{
		State: fwmNatAclRuleState(t, res, stateValues),
	}
	res.Delete(context.Background(), resource.DeleteRequest{
		State: fwmNatAclRuleState(t, res, stateValues),
	}, resp)

	assert.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())
}

func TestFwmNatAclRule_Delete_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DeleteNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteNatAclRuleRequest) (*fwmv20250611.DeleteNatAclRuleResponse, error) {
		return nil, assert.AnError
	})

	res := newFwmNatAclRuleTestResource(t)
	stateValues := fwmNatAclRuleBasePlanValues()
	stateValues["id"] = "acg-1a2b3c4d" + tccommon.FILED_SP + "aclr-9x8y7z6w"
	stateValues["rule_id"] = "aclr-9x8y7z6w"

	resp := &resource.DeleteResponse{
		State: fwmNatAclRuleState(t, res, stateValues),
	}
	res.Delete(context.Background(), resource.DeleteRequest{
		State: fwmNatAclRuleState(t, res, stateValues),
	}, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

func TestFwmNatAclRule_ImportState(t *testing.T) {
	res := newFwmNatAclRuleTestResource(t)

	importId := "acg-1a2b3c4d" + tccommon.FILED_SP + "aclr-9x8y7z6w"
	resp := &resource.ImportStateResponse{
		State: fwmNatAclRuleState(t, res, nil),
	}
	res.ImportState(context.Background(), resource.ImportStateRequest{ID: importId}, resp)

	assert.False(t, resp.Diagnostics.HasError(), resp.Diagnostics.Errors())

	var id types.String
	diags := resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
	assert.False(t, diags.HasError(), diags.Errors())
	assert.Equal(t, importId, id.ValueString())
}

func TestFwmNatAclRule_Metadata(t *testing.T) {
	res := newFwmNatAclRuleTestResource(t)

	resp := &resource.MetadataResponse{}
	res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "tencentcloud"}, resp)
	assert.Equal(t, "tencentcloud_fwm_nat_acl_rule", resp.TypeName)
}
