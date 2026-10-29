package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

func ptrStringFSV2(s string) *string {
	return &s
}

func ptrInt64FSV2(i int64) *int64 {
	return &i
}

// strategyGroupAttrTypes is the tftypes.Object type of one pre_strategy or
// post_strategy entry.
var strategyGroupAttrTypes = map[string]tftypes.Type{
	"group_id": tftypes.String,
	"sequence": tftypes.Number,
}

var strategyMemberAttrTypes = map[string]tftypes.Type{
	"app_id":    tftypes.String,
	"uin":       tftypes.String,
	"nickname":  tftypes.String,
	"member_id": tftypes.String,
	"node_name": tftypes.String,
}

var strategyReceiveAccountAttrTypes = map[string]tftypes.Type{
	"uin":           tftypes.String,
	"nickname":      tftypes.String,
	"receiver_type": tftypes.Number,
	"members":       tftypes.List{ElementType: tftypes.Object{AttributeTypes: strategyMemberAttrTypes}},
}

var strategyDetailAttrTypes = map[string]tftypes.Type{
	"strategy_id":     tftypes.String,
	"group_id":        tftypes.String,
	"group_name":      tftypes.String,
	"rule_count":      tftypes.Number,
	"rule_status":     tftypes.Number,
	"receive_account": tftypes.List{ElementType: tftypes.Object{AttributeTypes: strategyReceiveAccountAttrTypes}},
	"sequence":        tftypes.Number,
	"err_msg":         tftypes.String,
	"error_type":      tftypes.String,
	"update_time":     tftypes.String,
	"create_by":       tftypes.String,
	"update_by":       tftypes.String,
	"exec_area":       tftypes.String,
	"create_name":     tftypes.String,
	"update_name":     tftypes.String,
	"create_time":     tftypes.String,
}

// newFwmStrategyV2Resource builds a resource bound to a bare TencentCloud
// client, so that the FWM SDK client calls can be mocked with gomonkey.
func newFwmStrategyV2Resource() *fwm.FwmStrategyV2Resource {
	sharedmeta.SetSharedMeta(&connectivity.TencentCloudClient{})
	r := fwm.NewFwmStrategyV2Resource().(*fwm.FwmStrategyV2Resource)
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}},
	}, &resource.ConfigureResponse{})
	return r
}

func mockFwmClient() (*fwmv20250611.Client, *gomonkey.Patches) {
	fwmClient := &fwmv20250611.Client{}
	patches := gomonkey.NewPatches()
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", fwmClient)
	return fwmClient, patches
}

func resourceSchema(t *testing.T, r *fwm.FwmStrategyV2Resource) schema.Schema {
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	assert.False(t, resp.Diagnostics.HasError(), "Schema returned diagnostics: %v", resp.Diagnostics)
	return resp.Schema
}

// stateValue builds a full resource object value. A nil preStrategy means the
// attribute is null; entries are given as {group_id, sequence} pairs.
func stateValue(s schema.Schema, id string, preStrategy map[string]interface{}) tftypes.Value {
	preType := tftypes.Object{AttributeTypes: strategyGroupAttrTypes}

	var pre tftypes.Value
	if preStrategy == nil {
		pre = tftypes.NewValue(tftypes.List{ElementType: preType}, nil)
	} else {
		pre = tftypes.NewValue(tftypes.List{ElementType: preType}, []tftypes.Value{
			tftypes.NewValue(preType, map[string]tftypes.Value{
				"group_id": tftypes.NewValue(tftypes.String, preStrategy["group_id"]),
				"sequence": tftypes.NewValue(tftypes.Number, preStrategy["sequence"]),
			}),
		})
	}

	return tftypes.NewValue(s.Type().TerraformType(context.Background()), map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, id),
		"product":         tftypes.NewValue(tftypes.String, "cfw"),
		"receive_account": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"pre_strategy":    pre,
		"post_strategy":   tftypes.NewValue(tftypes.List{ElementType: preType}, nil),
		"receive_group":   tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"strategy":        tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: strategyDetailAttrTypes}}, nil),
	})
}

func newDescribeStrategyResponse(strategyId string) *fwmv20250611.DescribeStrategyResponse {
	resp := fwmv20250611.NewDescribeStrategyResponse()
	resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
		RequestId: ptrStringFSV2("fake-request-id"),
		Strategy: &fwmv20250611.StrategyResp{
			StrategyId: ptrStringFSV2(strategyId),
			GroupId:    ptrStringFSV2("grp-1a2b3c4d"),
			GroupName:  ptrStringFSV2("test-group"),
			RuleCount:  ptrInt64FSV2(3),
			RuleStatus: ptrInt64FSV2(1),
			Sequence:   ptrInt64FSV2(1),
			ErrMsg:     ptrStringFSV2(""),
			ErrorType:  ptrStringFSV2(""),
			UpdateTime: ptrStringFSV2("2026-10-22 10:00:00"),
			CreateBy:   ptrStringFSV2("100000000001"),
			UpdateBy:   ptrStringFSV2("100000000001"),
			ExecArea:   ptrStringFSV2("ap-guangzhou"),
			CreateName: ptrStringFSV2("tester"),
			UpdateName: ptrStringFSV2("tester"),
			CreateTime: ptrStringFSV2("2026-10-22 09:00:00"),
			ReceiveAccount: []*fwmv20250611.ReceiveAccount{
				{
					Uin:          ptrStringFSV2("100000000002"),
					Nickname:     ptrStringFSV2("tenant"),
					ReceiverType: ptrInt64FSV2(0),
				},
			},
		},
	}
	return resp
}

// TestFwmStrategyV2Resource_Metadata validates the resource type name.
func TestFwmStrategyV2Resource_Metadata(t *testing.T) {
	r := fwm.NewFwmStrategyV2Resource()
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)

	assert.Equal(t, "tencentcloud_fwm_strategy_v2", resp.TypeName)
}

// TestFwmStrategyV2Resource_Schema validates the attribute flags of the schema.
func TestFwmStrategyV2Resource_Schema(t *testing.T) {
	r := fwm.NewFwmStrategyV2Resource()
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	assert.False(t, resp.Diagnostics.HasError(), "Schema returned diagnostics: %v", resp.Diagnostics)

	product, ok := resp.Schema.Attributes["product"].(schema.StringAttribute)
	assert.True(t, ok, "product must be a StringAttribute")
	assert.True(t, product.Required, "product must be required")
	assert.False(t, product.Optional, "product must not be optional")
	assert.False(t, product.Computed, "product must not be computed")

	receiveAccount, ok := resp.Schema.Attributes["receive_account"].(schema.SetAttribute)
	assert.True(t, ok, "receive_account must be a SetAttribute")
	assert.True(t, receiveAccount.Optional, "receive_account must be optional")
	assert.False(t, receiveAccount.Computed, "receive_account must not be computed")

	receiveGroup, ok := resp.Schema.Attributes["receive_group"].(schema.SetAttribute)
	assert.True(t, ok, "receive_group must be a SetAttribute")
	assert.True(t, receiveGroup.Optional, "receive_group must be optional")

	preStrategy, ok := resp.Schema.Attributes["pre_strategy"].(schema.ListNestedAttribute)
	assert.True(t, ok, "pre_strategy must be a ListNestedAttribute")
	assert.True(t, preStrategy.Optional, "pre_strategy must be optional")

	postStrategy, ok := resp.Schema.Attributes["post_strategy"].(schema.ListNestedAttribute)
	assert.True(t, ok, "post_strategy must be a ListNestedAttribute")
	assert.True(t, postStrategy.Optional, "post_strategy must be optional")

	strategy, ok := resp.Schema.Attributes["strategy"].(schema.ListNestedAttribute)
	assert.True(t, ok, "strategy must be a ListNestedAttribute")
	assert.True(t, strategy.Computed, "strategy must be computed")
	assert.False(t, strategy.Optional, "strategy must not be optional")

	nestedReceiveAccount, ok := strategy.NestedObject.Attributes["receive_account"].(schema.ListNestedAttribute)
	assert.True(t, ok, "strategy.receive_account must be a ListNestedAttribute")
	assert.True(t, nestedReceiveAccount.Computed, "strategy.receive_account must be computed")

	_, exists := resp.Schema.Attributes["create_time"]
	assert.False(t, exists, "create_time must only exist inside the computed strategy list")

	for _, name := range []string{"strategy_id", "group_id", "group_name", "rule_count", "rule_status",
		"sequence", "err_msg", "error_type", "update_time", "create_by", "update_by", "exec_area",
		"create_name", "update_name", "create_time"} {
		attr, exists := strategy.NestedObject.Attributes[name]
		assert.True(t, exists, "strategy.%s is missing", name)
		assert.True(t, attr.IsComputed(), "strategy.%s must be computed", name)
	}
}

// TestFwmStrategyV2Resource_Create_ErrorIsReported verifies that a failing
// CreateStrategy call surfaces as an error diagnostic.
func TestFwmStrategyV2Resource_Create_ErrorIsReported(t *testing.T) {
	fwmClient, patches := mockFwmClient()
	defer patches.Reset()

	patches.ApplyMethodFunc(fwmClient, "CreateStrategyWithContext", func(_ context.Context, request *fwmv20250611.CreateStrategyRequest) (*fwmv20250611.CreateStrategyResponse, error) {
		return nil, assert.AnError
	})

	r := newFwmStrategyV2Resource()
	s := resourceSchema(t, r)
	planRaw := stateValue(s, "", nil)

	createResp := &resource.CreateResponse{
		State: tfsdk.State{Schema: s, Raw: planRaw},
	}
	r.Create(context.Background(), resource.CreateRequest{
		Plan: tfsdk.Plan{Schema: s, Raw: planRaw},
	}, createResp)

	assert.True(t, createResp.Diagnostics.HasError(), "Create must report the CreateStrategy failure")
}

// TestFwmStrategyV2Resource_Create_EmptyStrategyIdsIsError verifies that an
// empty StrategyIds is rejected instead of writing an empty id to state.
func TestFwmStrategyV2Resource_Create_EmptyStrategyIdsIsError(t *testing.T) {
	fwmClient, patches := mockFwmClient()
	defer patches.Reset()

	patches.ApplyMethodFunc(fwmClient, "CreateStrategyWithContext", func(_ context.Context, request *fwmv20250611.CreateStrategyRequest) (*fwmv20250611.CreateStrategyResponse, error) {
		resp := fwmv20250611.NewCreateStrategyResponse()
		resp.Response = &fwmv20250611.CreateStrategyResponseParams{RequestId: ptrStringFSV2("fake-request-id")}
		return resp, nil
	})

	r := newFwmStrategyV2Resource()
	s := resourceSchema(t, r)
	planRaw := stateValue(s, "", nil)

	createResp := &resource.CreateResponse{
		State: tfsdk.State{Schema: s, Raw: planRaw},
	}
	r.Create(context.Background(), resource.CreateRequest{
		Plan: tfsdk.Plan{Schema: s, Raw: planRaw},
	}, createResp)

	assert.True(t, createResp.Diagnostics.HasError(), "Create must reject an empty StrategyIds")
}

// TestFwmStrategyV2Resource_Read_ErrorIsReported verifies that a failing
// DescribeStrategy call surfaces as an error diagnostic.
func TestFwmStrategyV2Resource_Read_ErrorIsReported(t *testing.T) {
	fwmClient, patches := mockFwmClient()
	defer patches.Reset()

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		return nil, assert.AnError
	})

	r := newFwmStrategyV2Resource()
	s := resourceSchema(t, r)
	stateRaw := stateValue(s, "sg-1a2b3c4d", nil)

	readResp := &resource.ReadResponse{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}
	r.Read(context.Background(), resource.ReadRequest{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}, readResp)

	assert.True(t, readResp.Diagnostics.HasError(), "Read must report the DescribeStrategy failure")
}

// TestFwmStrategyV2Resource_Read_NotFoundRemovesState verifies that a strategy
// missing on the cloud side is removed from state.
func TestFwmStrategyV2Resource_Read_NotFoundRemovesState(t *testing.T) {
	fwmClient, patches := mockFwmClient()
	defer patches.Reset()

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{RequestId: ptrStringFSV2("fake-request-id")}
		return resp, nil
	})

	r := newFwmStrategyV2Resource()
	s := resourceSchema(t, r)
	stateRaw := stateValue(s, "sg-1a2b3c4d", nil)

	readResp := &resource.ReadResponse{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}
	r.Read(context.Background(), resource.ReadRequest{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}, readResp)

	assert.False(t, readResp.Diagnostics.HasError(), "Read returned diagnostics: %v", readResp.Diagnostics)
	assert.True(t, readResp.State.Raw.IsNull(), "state must be removed when the strategy no longer exists")
}

// TestFwmStrategyV2Resource_Read_FlattensStrategy verifies that the describe
// result is written into the computed strategy attribute.
func TestFwmStrategyV2Resource_Read_FlattensStrategy(t *testing.T) {
	fwmClient, patches := mockFwmClient()
	defer patches.Reset()

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		return newDescribeStrategyResponse("sg-1a2b3c4d"), nil
	})

	r := newFwmStrategyV2Resource()
	s := resourceSchema(t, r)
	stateRaw := stateValue(s, "sg-1a2b3c4d", nil)

	readResp := &resource.ReadResponse{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}
	r.Read(context.Background(), resource.ReadRequest{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}, readResp)

	assert.False(t, readResp.Diagnostics.HasError(), "Read returned diagnostics: %v", readResp.Diagnostics)

	var strategies types.List
	getDiags := readResp.State.GetAttribute(context.Background(), path.Root("strategy"), &strategies)
	if getDiags.HasError() {
		t.Fatalf("failed to read strategy from state: %v", getDiags)
	}
	assert.Len(t, strategies.Elements(), 1, "exactly one strategy is expected")
}

// TestFwmStrategyV2Resource_Delete_ErrorIsReported verifies that a failing
// DeleteStrategy call surfaces as an error diagnostic, and that every strategy
// id from state is sent to the API.
func TestFwmStrategyV2Resource_Delete_ErrorIsReported(t *testing.T) {
	fwmClient, patches := mockFwmClient()
	defer patches.Reset()

	var captured []*string
	patches.ApplyMethodFunc(fwmClient, "DeleteStrategyWithContext", func(_ context.Context, request *fwmv20250611.DeleteStrategyRequest) (*fwmv20250611.DeleteStrategyResponse, error) {
		captured = request.StrategyIds
		return nil, assert.AnError
	})

	r := newFwmStrategyV2Resource()
	s := resourceSchema(t, r)
	stateRaw := stateValue(s, "sg-1a2b3c4d#sg-5e6f7a8b", nil)

	deleteResp := &resource.DeleteResponse{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}
	r.Delete(context.Background(), resource.DeleteRequest{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}, deleteResp)

	assert.True(t, deleteResp.Diagnostics.HasError(), "Delete must report the DeleteStrategy failure")
	assert.Len(t, captured, 2, "Delete must send one StrategyIds entry per strategy id in state")
}

// TestFwmStrategyV2Resource_Update_ModifyRequest verifies that ModifyStrategy
// receives the strategy id and the mutable arguments from the plan.
func TestFwmStrategyV2Resource_Update_ModifyRequest(t *testing.T) {
	fwmClient, patches := mockFwmClient()
	defer patches.Reset()

	var modifyRequest *fwmv20250611.ModifyStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "ModifyStrategyWithContext", func(_ context.Context, request *fwmv20250611.ModifyStrategyRequest) (*fwmv20250611.ModifyStrategyResponse, error) {
		modifyRequest = request
		resp := fwmv20250611.NewModifyStrategyResponse()
		resp.Response = &fwmv20250611.ModifyStrategyResponseParams{RequestId: ptrStringFSV2("fake-request-id")}
		return resp, nil
	})
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		return newDescribeStrategyResponse("sg-1a2b3c4d"), nil
	})

	r := newFwmStrategyV2Resource()
	s := resourceSchema(t, r)
	planRaw := stateValue(s, "sg-1a2b3c4d", map[string]interface{}{"group_id": "grp-1a2b3c4d", "sequence": 2})
	stateRaw := stateValue(s, "sg-1a2b3c4d", map[string]interface{}{"group_id": "grp-1a2b3c4d", "sequence": 1})

	updateResp := &resource.UpdateResponse{
		State: tfsdk.State{Schema: s, Raw: planRaw},
	}
	r.Update(context.Background(), resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: s, Raw: planRaw},
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}, updateResp)

	assert.False(t, updateResp.Diagnostics.HasError(), "Update returned diagnostics: %v", updateResp.Diagnostics)
	if assert.NotNil(t, modifyRequest, "ModifyStrategy must be called") {
		assert.Equal(t, "sg-1a2b3c4d", *modifyRequest.StrategyId)
		if assert.NotNil(t, modifyRequest.GroupId) {
			assert.Equal(t, "grp-1a2b3c4d", *modifyRequest.GroupId)
		}
		if assert.NotNil(t, modifyRequest.Sequence) {
			assert.Equal(t, int64(2), *modifyRequest.Sequence)
		}
	}
}

// TestFwmStrategyV2Resource_Update_MultipleStrategiesIsError verifies that a
// resource composed of more than one strategy cannot be updated in place.
func TestFwmStrategyV2Resource_Update_MultipleStrategiesIsError(t *testing.T) {
	_, patches := mockFwmClient()
	defer patches.Reset()

	r := newFwmStrategyV2Resource()
	s := resourceSchema(t, r)
	id := "sg-1a2b3c4d#sg-5e6f7a8b"
	stateRaw := stateValue(s, id, nil)

	updateResp := &resource.UpdateResponse{
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}
	r.Update(context.Background(), resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: s, Raw: stateRaw},
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}, updateResp)

	assert.True(t, updateResp.Diagnostics.HasError(), "Update must reject multi strategy resources")
}

// TestFwmStrategyV2Resource_ConfigureWithoutProviderData makes sure a missing
// provider data does not panic and simply leaves the client unset.
func TestFwmStrategyV2Resource_ConfigureWithoutProviderData(t *testing.T) {
	r := fwm.NewFwmStrategyV2Resource().(*fwm.FwmStrategyV2Resource)
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{}, configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())
}
