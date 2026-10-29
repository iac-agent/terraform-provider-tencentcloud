package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

func ptrStringFRG(s string) *string {
	return &s
}

func ptrInt64FRG(i int64) *int64 {
	return &i
}

// ruleGroupTestClient is the mocked TencentCloud client shared by the service
// layer calls in this test file.
var ruleGroupTestClient *connectivity.TencentCloudClient

// newRuleGroupResource builds a FwmRuleGroupResource with the mocked client
// already injected, as the framework provider would do in Configure.
func newRuleGroupResource() (*fwm.FwmRuleGroupResource, *fwmv20250611.Client) {
	fwmClient := &fwmv20250611.Client{}
	client := &connectivity.TencentCloudClient{}
	ruleGroupTestClient = client
	patches := gomonkey.NewPatches()
	// The patches are intentionally not reset here: the fake client methods are
	// replaced by each test case, and the client accessor stays mocked for the
	// duration of the test process.
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)

	res := fwm.NewFwmRuleGroupResource().(*fwm.FwmRuleGroupResource)
	res.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &resource.ConfigureResponse{})

	return res, fwmClient
}

func TestFwmRuleGroup_Metadata(t *testing.T) {
	res, _ := newRuleGroupResource()

	metaResp := resource.MetadataResponse{}
	res.Metadata(context.Background(), resource.MetadataRequest{}, &metaResp)
	assert.Equal(t, "tencentcloud_fwm_rule_group", metaResp.TypeName)
}

func TestFwmRuleGroup_Schema_HasExpectedAttributes(t *testing.T) {
	res, _ := newRuleGroupResource()

	schemaResp := resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	expected := []string{
		"id",
		"group_name",
		"product",
		"filters",
		"order",
		"by",
		"rule_groups",
		"group_id",
		"group_ids",
	}
	for _, name := range expected {
		_, ok := schemaResp.Schema.Attributes[name]
		assert.True(t, ok, "attribute %q should exist in the schema", name)
	}
	assert.Len(t, schemaResp.Schema.Attributes, len(expected))

	ruleGroups, ok := schemaResp.Schema.Attributes["rule_groups"]
	assert.True(t, ok)
	assert.True(t, ruleGroups.IsComputed())
	assert.False(t, ruleGroups.IsRequired())
	assert.False(t, ruleGroups.IsOptional())
}

func TestFwmRuleGroup_Create_EmptyGroupIdReturnsError(t *testing.T) {
	res, fwmClient := newRuleGroupResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext",
		func(_ context.Context, request *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
			resp := fwmv20250611.NewCreateRuleGroupResponse()
			resp.Response = &fwmv20250611.CreateRuleGroupResponseParams{
				RequestId: ptrStringFRG("fake-request-id"),
				GroupId:   nil,
			}
			return resp, nil
		})

	service := fwm.NewFwmService(ruleGroupTestClient)
	groupId, err := service.CreateFwmRuleGroup(context.Background(), "tf-rule-group", "cfw_vpc_acl")
	assert.Error(t, err)
	assert.Equal(t, "", groupId)
}

func TestFwmRuleGroup_Create_ApiErrorPropagated(t *testing.T) {
	res, fwmClient := newRuleGroupResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext",
		func(_ context.Context, request *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
			return nil, assert.AnError
		})

	service := fwm.NewFwmService(ruleGroupTestClient)
	groupId, err := service.CreateFwmRuleGroup(context.Background(), "tf-rule-group", "cfw_vpc_acl")
	assert.Error(t, err)
	assert.Equal(t, "", groupId)
}

func TestFwmRuleGroup_Describe_NotFound(t *testing.T) {
	res, fwmClient := newRuleGroupResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext",
		func(_ context.Context, request *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
			resp := fwmv20250611.NewDescribeRuleGroupsResponse()
			resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
				RequestId:  ptrStringFRG("fake-request-id"),
				TotalCount: ptrInt64FRG(0),
				RuleGroups: []*fwmv20250611.RuleGroupResp{},
			}
			return resp, nil
		})

	service := fwm.NewFwmService(ruleGroupTestClient)
	ruleGroup, err := service.DescribeFwmRuleGroupById(context.Background(), "fwmrg-notexist")
	assert.NoError(t, err)
	assert.Nil(t, ruleGroup)
}

func TestFwmRuleGroup_Describe_Normal(t *testing.T) {
	res, fwmClient := newRuleGroupResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext",
		func(_ context.Context, request *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
			assert.NotNil(t, request.Filters)
			assert.Len(t, request.Filters, 1)
			assert.Equal(t, "GroupId", *request.Filters[0].Name)
			assert.Equal(t, "fwmrg-1a2b3c4d", *request.Filters[0].Values[0])

			resp := fwmv20250611.NewDescribeRuleGroupsResponse()
			resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
				RequestId:  ptrStringFRG("fake-request-id"),
				TotalCount: ptrInt64FRG(1),
				RuleGroups: []*fwmv20250611.RuleGroupResp{
					{
						AppId:      ptrStringFRG("1300000001"),
						GroupId:    ptrStringFRG("fwmrg-1a2b3c4d"),
						GroupName:  ptrStringFRG("tf-rule-group"),
						RuleCount:  ptrInt64FRG(0),
						Product:    ptrStringFRG("cfw_vpc_acl"),
						CreateTime: ptrStringFRG("2026-10-22 10:00:00"),
						UpdateTime: ptrStringFRG("2026-10-22 10:00:00"),
						CreateBy:   ptrStringFRG("100000000001"),
						CreateName: ptrStringFRG("creator"),
						UpdateBy:   ptrStringFRG("100000000001"),
						UpdateName: ptrStringFRG("creator"),
						ReceiveAccount: []*fwmv20250611.ReceiveAccount{
							{
								Uin:          ptrStringFRG("100000000001"),
								Nickname:     ptrStringFRG("tenant-a"),
								ReceiverType: ptrInt64FRG(0),
							},
						},
						ExecAreas: []*string{ptrStringFRG("post")},
						Direction: ptrInt64FRG(1),
					},
				},
			}
			return resp, nil
		})

	service := fwm.NewFwmService(ruleGroupTestClient)
	ruleGroup, err := service.DescribeFwmRuleGroupById(context.Background(), "fwmrg-1a2b3c4d")
	assert.NoError(t, err)
	assert.NotNil(t, ruleGroup)
	assert.Equal(t, "tf-rule-group", *ruleGroup.GroupName)
	assert.Equal(t, "cfw_vpc_acl", *ruleGroup.Product)

	model := fwm.BuildRuleGroupModelForTest(ruleGroup)
	assert.Equal(t, "tf-rule-group", model.GroupName.ValueString())
	assert.Equal(t, "cfw_vpc_acl", model.Product.ValueString())
	assert.Len(t, model.RuleGroups, 1)
	assert.Equal(t, "fwmrg-1a2b3c4d", model.RuleGroups[0].GroupId.ValueString())
	assert.Equal(t, int64(0), model.RuleGroups[0].RuleCount.ValueInt64())
	assert.Len(t, model.RuleGroups[0].ExecAreas, 1)
	assert.Equal(t, "post", model.RuleGroups[0].ExecAreas[0].ValueString())
	assert.Equal(t, int64(1), model.RuleGroups[0].Direction.ValueInt64())
	assert.Len(t, model.RuleGroups[0].ReceiveAccount, 1)
}

func TestFwmRuleGroup_BuildReceiveAccounts_Empty(t *testing.T) {
	assert.Nil(t, fwm.BuildReceiveAccountsForTest(nil))
	assert.Nil(t, fwm.BuildReceiveAccountsForTest([]*fwmv20250611.ReceiveAccount{}))
}

func TestFwmRuleGroup_BuildReceiveAccounts_Nested(t *testing.T) {
	in := []*fwmv20250611.ReceiveAccount{
		{
			Uin:          ptrStringFRG("100000000001"),
			Nickname:     ptrStringFRG("tenant-a"),
			ReceiverType: ptrInt64FRG(1),
			Members: []*fwmv20250611.MemberInfo{
				{
					AppId:    ptrStringFRG("1300000001"),
					Uin:      ptrStringFRG("100000000002"),
					Nickname: ptrStringFRG("member-a"),
					MemberId: ptrStringFRG("member-1"),
					NodeName: ptrStringFRG("department-a"),
				},
			},
		},
	}

	out := fwm.BuildReceiveAccountsForTest(in)
	assert.Len(t, out, 1)
	assert.Equal(t, "100000000001", out[0].Uin.ValueString())
	assert.Equal(t, int64(1), out[0].ReceiverType.ValueInt64())
	assert.Len(t, out[0].Members, 1)
	assert.Equal(t, "member-1", out[0].Members[0].MemberId.ValueString())
	assert.Equal(t, "department-a", out[0].Members[0].NodeName.ValueString())
}

func TestFwmRuleGroup_StringSliceToTypeSlice_NilElement(t *testing.T) {
	out := fwm.StringSliceToTypeSliceForTest([]*string{ptrStringFRG("inbound"), nil})
	assert.Len(t, out, 2)
	assert.Equal(t, "inbound", out[0].ValueString())
	assert.True(t, out[1].IsNull())
}

func TestFwmRuleGroup_Modify_ApiErrorPropagated(t *testing.T) {
	res, fwmClient := newRuleGroupResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(fwmClient, "ModifyRuleGroupWithContext",
		func(_ context.Context, request *fwmv20250611.ModifyRuleGroupRequest) (*fwmv20250611.ModifyRuleGroupResponse, error) {
			assert.Equal(t, "fwmrg-1a2b3c4d", *request.GroupId)
			assert.Equal(t, "tf-rule-group-updated", *request.GroupName)
			return nil, assert.AnError
		})

	service := fwm.NewFwmService(ruleGroupTestClient)
	err := service.ModifyFwmRuleGroupName(context.Background(), "fwmrg-1a2b3c4d", "tf-rule-group-updated")
	assert.Error(t, err)
}

func TestFwmRuleGroup_Modify_Normal(t *testing.T) {
	res, fwmClient := newRuleGroupResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(fwmClient, "ModifyRuleGroupWithContext",
		func(_ context.Context, request *fwmv20250611.ModifyRuleGroupRequest) (*fwmv20250611.ModifyRuleGroupResponse, error) {
			assert.Equal(t, "fwmrg-1a2b3c4d", *request.GroupId)
			assert.Equal(t, "tf-rule-group-updated", *request.GroupName)

			resp := fwmv20250611.NewModifyRuleGroupResponse()
			resp.Response = &fwmv20250611.ModifyRuleGroupResponseParams{
				RequestId: ptrStringFRG("fake-request-id"),
			}
			return resp, nil
		})

	service := fwm.NewFwmService(ruleGroupTestClient)
	err := service.ModifyFwmRuleGroupName(context.Background(), "fwmrg-1a2b3c4d", "tf-rule-group-updated")
	assert.NoError(t, err)
}

func TestFwmRuleGroup_Delete_Normal(t *testing.T) {
	res, fwmClient := newRuleGroupResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(fwmClient, "DeleteRuleGroupWithContext",
		func(_ context.Context, request *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
			assert.Len(t, request.GroupIds, 1)
			assert.Equal(t, "fwmrg-1a2b3c4d", *request.GroupIds[0])

			resp := fwmv20250611.NewDeleteRuleGroupResponse()
			resp.Response = &fwmv20250611.DeleteRuleGroupResponseParams{
				RequestId: ptrStringFRG("fake-request-id"),
			}
			return resp, nil
		})

	service := fwm.NewFwmService(ruleGroupTestClient)
	err := service.DeleteFwmRuleGroupById(context.Background(), "fwmrg-1a2b3c4d")
	assert.NoError(t, err)
}

func TestFwmRuleGroup_Delete_ApiErrorPropagated(t *testing.T) {
	res, fwmClient := newRuleGroupResource()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(fwmClient, "DeleteRuleGroupWithContext",
		func(_ context.Context, request *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
			return nil, assert.AnError
		})

	service := fwm.NewFwmService(ruleGroupTestClient)
	err := service.DeleteFwmRuleGroupById(context.Background(), "fwmrg-1a2b3c4d")
	assert.Error(t, err)
}
