package fwm

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

func newTestFwmStrategyResource() *FwmStrategyResource {
	return &FwmStrategyResource{
		client: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}},
	}
}

func ptrFwmString(s string) *string { return &s }
func ptrFwmInt64(i int64) *int64    { return &i }
func ptrFwmStrings(s ...string) []*string {
	out := make([]*string, 0, len(s))
	for _, v := range s {
		c := v
		out = append(out, &c)
	}
	return out
}

func TestFwmStrategyResourceRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &sdk.Client{}
	patches.ApplyMethodReturn(newTestFwmStrategyResource().client.Client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, request *sdk.DescribeStrategyRequest) (*sdk.DescribeStrategyResponse, error) {
		assert.NotNil(t, request.StrategyId)
		assert.Equal(t, "strategy-xxxx", *request.StrategyId)

		resp := sdk.NewDescribeStrategyResponse()
		resp.Response = &sdk.DescribeStrategyResponseParams{
			Strategy: &sdk.StrategyResp{
				StrategyId: ptrFwmString("strategy-xxxx"),
				GroupId:    ptrFwmString("group-xxxx"),
				GroupName:  ptrFwmString("group-name"),
				RuleCount:  ptrFwmInt64(3),
				RuleStatus: ptrFwmInt64(1),
				Sequence:   ptrFwmInt64(10),
				ReceiveAccount: []*sdk.ReceiveAccount{
					{
						Uin:          ptrFwmString("100000000001"),
						Nickname:     ptrFwmString("tenant-a"),
						ReceiverType: ptrFwmInt64(1),
						Members: []*sdk.MemberInfo{
							{
								AppId:    ptrFwmString("1250000001"),
								Uin:      ptrFwmString("100000000001"),
								Nickname: ptrFwmString("member-a"),
								MemberId: ptrFwmString("member-xxxx"),
								NodeName: ptrFwmString("dept-a"),
							},
						},
					},
				},
				ErrMsg:     ptrFwmString(""),
				ErrorType:  ptrFwmString(""),
				UpdateTime: ptrFwmString("2025-06-01 12:00:00"),
				CreateBy:   ptrFwmString("user-a"),
				UpdateBy:   ptrFwmString("user-b"),
				ExecArea:   ptrFwmString("all"),
				CreateName: ptrFwmString("creator-a"),
				UpdateName: ptrFwmString("updater-b"),
				CreateTime: ptrFwmString("2025-06-01 10:00:00"),
			},
		}
		return resp, nil
	})

	res := newTestFwmStrategyResource()
	item, err := res.read(context.Background(), "strategy-xxxx")
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "strategy-xxxx", *item.StrategyId)
	assert.Equal(t, int64(3), *item.RuleCount)
	assert.Equal(t, 1, len(item.ReceiveAccount))
	assert.Equal(t, "member-a", *item.ReceiveAccount[0].Members[0].Nickname)
}

func TestFwmStrategyResourceReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &sdk.Client{}
	patches.ApplyMethodReturn(newTestFwmStrategyResource().client.Client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, request *sdk.DescribeStrategyRequest) (*sdk.DescribeStrategyResponse, error) {
		resp := sdk.NewDescribeStrategyResponse()
		resp.Response = &sdk.DescribeStrategyResponseParams{}
		return resp, nil
	})

	res := newTestFwmStrategyResource()
	item, err := res.read(context.Background(), "strategy-missing")
	assert.NoError(t, err)
	assert.Nil(t, item)
}

func TestFwmStrategyResourceReadError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &sdk.Client{}
	patches.ApplyMethodReturn(newTestFwmStrategyResource().client.Client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, request *sdk.DescribeStrategyRequest) (*sdk.DescribeStrategyResponse, error) {
		return nil, assert.AnError
	})

	res := newTestFwmStrategyResource()
	_, err := res.read(context.Background(), "strategy-xxxx")
	assert.Error(t, err)
}

func TestFwmStrategyResourceRefreshModel(t *testing.T) {
	m := &FwmStrategyResourceModel{}
	item := &sdk.StrategyResp{
		StrategyId: ptrFwmString("strategy-xxxx"),
		GroupId:    ptrFwmString("group-xxxx"),
		GroupName:  ptrFwmString("group-name"),
		RuleCount:  ptrFwmInt64(3),
		RuleStatus: ptrFwmInt64(1),
		Sequence:   ptrFwmInt64(10),
		ErrMsg:     ptrFwmString("failed"),
		ErrorType:  ptrFwmString("type-a"),
		UpdateTime: ptrFwmString("2025-06-01 12:00:00"),
		CreateBy:   ptrFwmString("user-a"),
		UpdateBy:   ptrFwmString("user-b"),
		ExecArea:   ptrFwmString("all"),
		CreateName: ptrFwmString("creator-a"),
		UpdateName: ptrFwmString("updater-b"),
		CreateTime: ptrFwmString("2025-06-01 10:00:00"),
	}

	res := &FwmStrategyResource{}
	res.refreshModel(m, item)

	assert.Equal(t, "strategy-xxxx", m.Id.ValueString())
	assert.Equal(t, "strategy-xxxx", m.StrategyId.ValueString())
	assert.Equal(t, "group-xxxx", m.GroupId.ValueString())
	assert.Equal(t, int64(10), m.Sequence.ValueInt64())
	assert.Equal(t, "group-name", m.GroupName.ValueString())
	assert.Equal(t, int64(3), m.RuleCount.ValueInt64())
	assert.Equal(t, int64(1), m.RuleStatus.ValueInt64())
	assert.Equal(t, "failed", m.ErrMsg.ValueString())
	assert.Equal(t, "type-a", m.ErrorType.ValueString())
	assert.Equal(t, "all", m.ExecArea.ValueString())
	assert.Equal(t, "2025-06-01 10:00:00", m.CreateTime.ValueString())
	assert.True(t, m.ReceiveAccountDetail.IsNull())
}

func TestFwmStrategyResourceRefreshModelNil(t *testing.T) {
	m := &FwmStrategyResourceModel{}
	res := &FwmStrategyResource{}
	res.refreshModel(m, nil)
	assert.True(t, m.Id.IsNull())
}

func TestFwmStrategyResourceMetadata(t *testing.T) {
	res := &FwmStrategyResource{}
	mreq := resource.MetadataRequest{ProviderTypeName: "tencentcloud"}
	mresp := &resource.MetadataResponse{}
	res.Metadata(context.Background(), mreq, mresp)
	assert.Equal(t, "tencentcloud_fwm_strategy", mresp.TypeName)
}

func TestFwmStrategyResourceDeleteRequest(t *testing.T) {
	request := sdk.NewDeleteStrategyRequest()
	request.StrategyIds = ptrFwmStrings("strategy-xxxx")
	assert.Equal(t, 1, len(request.StrategyIds))
	assert.Equal(t, "strategy-xxxx", *request.StrategyIds[0])
}
