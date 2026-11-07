package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

const (
	teoReplicaV3ZoneId      = "zone-2qtuhspy7cr6"
	teoReplicaV3FunctionId  = "ef-2qlxy8s7o96e"
	teoReplicaV3ReplicaName = "replica-a"
	teoReplicaV3Id          = "zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-a"
)

func ptrStringTeoReplicaV3(s string) *string {
	return &s
}

func ptrInt64TeoReplicaV3(i int64) *int64 {
	return &i
}

// errTeoReplicaV3 builds a plain error carrying the given message.
func errTeoReplicaV3(msg string) error {
	return &teoReplicaV3Error{msg: msg}
}

type teoReplicaV3Error struct {
	msg string
}

func (e *teoReplicaV3Error) Error() string {
	return e.msg
}

// teoReplicaV3AttrTypes describes the object shape of the resource model. It is
// checked against the real schema by teoReplicaV3Schema so the raw values built
// by the tests cannot silently diverge from the schema.
func teoReplicaV3AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.StringType,
		"zone_id":      types.StringType,
		"function_id":  types.StringType,
		"replica_name": types.StringType,
		"content":      types.StringType,
		"remark":       types.StringType,
		"created_on":   types.StringType,
		"modified_on":  types.StringType,
	}
}

// teoReplicaV3Schema returns the framework resource schema.
func teoReplicaV3Schema(t *testing.T) resource_schema.Schema {
	t.Helper()

	r := teo.NewTeoFunctionReplicaV3Resource()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}

	s := schemaResp.Schema
	attrTypes := teoReplicaV3AttrTypes()
	if len(s.Attributes) != len(attrTypes) {
		t.Fatalf("schema has %d attributes, model has %d", len(s.Attributes), len(attrTypes))
	}
	for name, wantType := range attrTypes {
		a, ok := s.Attributes[name]
		if !ok {
			t.Fatalf("schema is missing attribute %s", name)
		}
		if !a.GetType().Equal(wantType) {
			t.Fatalf("attribute %s type mismatch: schema=%v model=%v", name, a.GetType(), wantType)
		}
	}
	return s
}

// teoReplicaV3Model builds a complete model value. Attributes absent from
// values are null.
func teoReplicaV3Model(t *testing.T, values map[string]attr.Value) types.Object {
	t.Helper()

	attrTypes := teoReplicaV3AttrTypes()
	full := make(map[string]attr.Value, len(attrTypes))
	for name := range attrTypes {
		if v, ok := values[name]; ok {
			full[name] = v
		} else {
			full[name] = types.StringNull()
		}
	}
	obj, diags := types.ObjectValue(attrTypes, full)
	if diags.HasError() {
		t.Fatalf("failed to build model object: %v", diags)
	}
	return obj
}

// teoReplicaV3Raw converts a model object into the raw value expected by
// tfsdk.State / tfsdk.Plan.
func teoReplicaV3Raw(t *testing.T, obj types.Object) tftypes.Value {
	t.Helper()

	raw, err := obj.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("failed to convert model to terraform value: %v", err)
	}
	return raw
}

// teoReplicaV3StateAttr reads a single attribute from a raw state value.
func teoReplicaV3StateAttr(t *testing.T, raw tftypes.Value, name string) types.String {
	t.Helper()

	got, remaining, err := tftypes.WalkAttributePath(raw, tftypes.NewAttributePath().WithAttributeName(name))
	if err != nil {
		t.Fatalf("failed to walk attribute %s: %v (remaining %v)", name, err, remaining)
	}
	v, ok := got.(tftypes.Value)
	if !ok {
		t.Fatalf("attribute %s is not a tftypes.Value: %T", name, got)
	}
	if v.IsNull() {
		return types.StringNull()
	}
	var s string
	if err := v.As(&s); err != nil {
		t.Fatalf("failed to decode attribute %s: %v", name, err)
	}
	return types.StringValue(s)
}

// teoReplicaV3ModelFromState reads the model back from a state value.
func teoReplicaV3ModelFromState(t *testing.T, s resource_schema.Schema, raw tftypes.Value) teo.TeoFunctionReplicaV3ResourceModel {
	t.Helper()

	state := tfsdk.State{Raw: raw, Schema: s}
	var model teo.TeoFunctionReplicaV3ResourceModel
	if diags := state.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	return model
}

// setTeoReplicaV3Client injects a mock connectivity client into the resource
// instance through the Configure method.
func setTeoReplicaV3Client(client *connectivity.TencentCloudClient) *teo.TeoFunctionReplicaV3Resource {
	r, ok := teo.NewTeoFunctionReplicaV3Resource().(*teo.TeoFunctionReplicaV3Resource)
	if !ok {
		panic("factory did not return *teo.TeoFunctionReplicaV3Resource")
	}
	meta := &sharedmeta.ProviderMeta{Client: client}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: meta}, &resource.ConfigureResponse{})
	return r
}

// teoReplicaV3Mock wires the four cloud APIs of a mocked TEO client. A nil
// describe callback yields an empty successful response body.
type teoReplicaV3Mock struct {
	createErr   error
	modifyErr   error
	deleteErr   error
	describe    func(*teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error)
	gotCreate   *teov20220901.CreateFunctionReplicaRequest
	gotModify   *teov20220901.ModifyFunctionReplicaRequest
	gotDelete   *teov20220901.DeleteFunctionReplicaRequest
	gotDescribe *teov20220901.DescribeFunctionReplicasRequest
}

// apply installs the mock onto the given patches and TEO client.
func (m *teoReplicaV3Mock) apply(p *gomonkey.Patches, client *teov20220901.Client) {
	p.ApplyMethodFunc(client, "CreateFunctionReplicaWithContext",
		func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
			if m.createErr != nil {
				return nil, m.createErr
			}
			m.gotCreate = request
			resp := teov20220901.NewCreateFunctionReplicaResponse()
			resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
				RequestId: ptrStringTeoReplicaV3("fake-request-id"),
			}
			return resp, nil
		})

	p.ApplyMethodFunc(client, "ModifyFunctionReplicaWithContext",
		func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
			if m.modifyErr != nil {
				return nil, m.modifyErr
			}
			m.gotModify = request
			resp := teov20220901.NewModifyFunctionReplicaResponse()
			resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
				RequestId: ptrStringTeoReplicaV3("fake-request-id"),
			}
			return resp, nil
		})

	p.ApplyMethodFunc(client, "DeleteFunctionReplicaWithContext",
		func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
			if m.deleteErr != nil {
				return nil, m.deleteErr
			}
			m.gotDelete = request
			resp := teov20220901.NewDeleteFunctionReplicaResponse()
			resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
				RequestId: ptrStringTeoReplicaV3("fake-request-id"),
			}
			return resp, nil
		})

	p.ApplyMethodFunc(client, "DescribeFunctionReplicasWithContext",
		func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			m.gotDescribe = request
			if m.describe != nil {
				return m.describe(request)
			}
			resp := teov20220901.NewDescribeFunctionReplicasResponse()
			resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{}
			return resp, nil
		})
}

// teoReplicaV3DescribeOne builds a describe response holding the given replica.
func teoReplicaV3DescribeOne(content, remark string) *teov20220901.DescribeFunctionReplicasResponse {
	resp := teov20220901.NewDescribeFunctionReplicasResponse()
	resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
		TotalCount: ptrInt64TeoReplicaV3(1),
		FunctionReplicas: []*teov20220901.FunctionReplica{
			{
				FunctionId:  ptrStringTeoReplicaV3(teoReplicaV3FunctionId),
				ReplicaName: ptrStringTeoReplicaV3(teoReplicaV3ReplicaName),
				Content:     ptrStringTeoReplicaV3(content),
				Remark:      ptrStringTeoReplicaV3(remark),
				CreatedOn:   ptrStringTeoReplicaV3("2025-01-01 00:00:00"),
				ModifiedOn:  ptrStringTeoReplicaV3("2025-01-02 00:00:00"),
			},
		},
		RequestId: ptrStringTeoReplicaV3("fake-request-id"),
	}
	return resp
}

// teoReplicaV3FullAttrs returns the attributes of an existing replica.
func teoReplicaV3FullAttrs(content, remark string) map[string]attr.Value {
	return map[string]attr.Value{
		"id":           types.StringValue(teoReplicaV3Id),
		"zone_id":      types.StringValue(teoReplicaV3ZoneId),
		"function_id":  types.StringValue(teoReplicaV3FunctionId),
		"replica_name": types.StringValue(teoReplicaV3ReplicaName),
		"content":      types.StringValue(content),
		"remark":       types.StringValue(remark),
	}
}

// teoReplicaV3CreatePlanAttrs returns the attributes of a create plan.
func teoReplicaV3CreatePlanAttrs(content string, remark *string) map[string]attr.Value {
	attrs := map[string]attr.Value{
		"zone_id":      types.StringValue(teoReplicaV3ZoneId),
		"function_id":  types.StringValue(teoReplicaV3FunctionId),
		"replica_name": types.StringValue(teoReplicaV3ReplicaName),
		"content":      types.StringValue(content),
	}
	if remark != nil {
		attrs["remark"] = types.StringValue(*remark)
	}
	return attrs
}

// TestTeoReplicaV3_MetadataAndSchema validates the type name and every
// attribute flag of the schema.
func TestTeoReplicaV3_MetadataAndSchema(t *testing.T) {
	r := teo.NewTeoFunctionReplicaV3Resource()

	metaResp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, metaResp)
	if metaResp.TypeName != "tencentcloud_teo_function_replica_v3" {
		t.Fatalf("unexpected type name: %s", metaResp.TypeName)
	}

	s := teoReplicaV3Schema(t)

	for _, name := range []string{"zone_id", "function_id", "replica_name", "content"} {
		if !s.Attributes[name].IsRequired() {
			t.Fatalf("attribute %s must be required", name)
		}
	}

	for _, name := range []string{"id", "created_on", "modified_on"} {
		if !s.Attributes[name].IsComputed() || s.Attributes[name].IsRequired() || s.Attributes[name].IsOptional() {
			t.Fatalf("attribute %s must be computed only", name)
		}
	}

	if !s.Attributes["remark"].IsOptional() || s.Attributes["remark"].IsRequired() {
		t.Fatal("attribute remark must be optional")
	}
	if s.Attributes["remark"].IsComputed() {
		t.Fatal("attribute remark must not be computed: a cleared remark has to stay null")
	}
}

// TestTeoReplicaV3_ReplaceOnIdentityFields validates that the three identity
// attributes and the computed id force a replacement while content and remark
// are updated in place.
func TestTeoReplicaV3_ReplaceOnIdentityFields(t *testing.T) {
	s := teoReplicaV3Schema(t)
	for _, name := range []string{"id", "zone_id", "function_id", "replica_name"} {
		a, ok := s.Attributes[name].(resource_schema.StringAttribute)
		if !ok {
			t.Fatalf("attribute %s is not a StringAttribute", name)
		}
		if len(a.PlanModifiers) != 1 {
			t.Fatalf("attribute %s must carry exactly 1 plan modifier, got %d", name, len(a.PlanModifiers))
		}
	}

	for _, name := range []string{"content", "remark"} {
		a, ok := s.Attributes[name].(resource_schema.StringAttribute)
		if !ok {
			t.Fatalf("attribute %s is not a StringAttribute", name)
		}
		if len(a.PlanModifiers) != 0 {
			t.Fatalf("attribute %s must not carry plan modifiers, got %d", name, len(a.PlanModifiers))
		}
	}
}

// TestTeoReplicaV3_Create_Success covers a create with a full read back.
func TestTeoReplicaV3_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			return teoReplicaV3DescribeOne("console.log(1)", "note a"), nil
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	plan := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3CreatePlanAttrs("console.log(1)", ptrStringTeoReplicaV3("note a"))))

	resp := &resource.CreateResponse{State: tfsdk.State{Raw: plan, Schema: s}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Raw: plan, Schema: s}}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}

	if mock.gotCreate == nil {
		t.Fatal("expected CreateFunctionReplica to be called")
	}
	if v := mock.gotCreate.ZoneId; v == nil || *v != teoReplicaV3ZoneId {
		t.Fatalf("unexpected ZoneId: %v", v)
	}
	if v := mock.gotCreate.FunctionId; v == nil || *v != teoReplicaV3FunctionId {
		t.Fatalf("unexpected FunctionId: %v", v)
	}
	if v := mock.gotCreate.ReplicaName; v == nil || *v != teoReplicaV3ReplicaName {
		t.Fatalf("unexpected ReplicaName: %v", v)
	}
	if v := mock.gotCreate.Content; v == nil || *v != "console.log(1)" {
		t.Fatalf("unexpected Content: %v", v)
	}
	if v := mock.gotCreate.Remark; v == nil || *v != "note a" {
		t.Fatalf("unexpected Remark: %v", v)
	}

	if mock.gotDescribe == nil {
		t.Fatal("expected DescribeFunctionReplicas to be called")
	}
	if v := mock.gotDescribe.ZoneId; v == nil || *v != teoReplicaV3ZoneId {
		t.Fatalf("unexpected describe ZoneId: %v", v)
	}
	if v := mock.gotDescribe.FunctionId; v == nil || *v != teoReplicaV3FunctionId {
		t.Fatalf("unexpected describe FunctionId: %v", v)
	}
	if len(mock.gotDescribe.Filters) != 1 || mock.gotDescribe.Filters[0] == nil {
		t.Fatalf("expected exactly 1 describe filter, got %d", len(mock.gotDescribe.Filters))
	}
	if v := mock.gotDescribe.Filters[0].Name; v == nil || *v != "replica-name" {
		t.Fatalf("unexpected filter name: %v", v)
	}
	if v := mock.gotDescribe.Filters[0].Fuzzy; v == nil || *v {
		t.Fatal("expected the fuzzy flag to be explicitly disabled")
	}

	model := teoReplicaV3ModelFromState(t, s, resp.State.Raw)
	if got := model.Id.ValueString(); got != teoReplicaV3Id {
		t.Fatalf("unexpected id: %s", got)
	}
	if got := model.ZoneId.ValueString(); got != teoReplicaV3ZoneId {
		t.Fatalf("unexpected zone_id: %s", got)
	}
	if got := model.FunctionId.ValueString(); got != teoReplicaV3FunctionId {
		t.Fatalf("unexpected function_id: %s", got)
	}
	if got := model.ReplicaName.ValueString(); got != teoReplicaV3ReplicaName {
		t.Fatalf("unexpected replica_name: %s", got)
	}
	if got := model.Content.ValueString(); got != "console.log(1)" {
		t.Fatalf("unexpected content: %s", got)
	}
	if got := model.Remark.ValueString(); got != "note a" {
		t.Fatalf("unexpected remark: %s", got)
	}
	if got := model.CreatedOn.ValueString(); got != "2025-01-01 00:00:00" {
		t.Fatalf("unexpected created_on: %s", got)
	}
	if got := model.ModifiedOn.ValueString(); got != "2025-01-02 00:00:00" {
		t.Fatalf("unexpected modified_on: %s", got)
	}
}

// TestTeoReplicaV3_Create_NullRemarkAndNullTimestamps covers the null branches
// of the read back.
func TestTeoReplicaV3_Create_NullRemarkAndNullTimestamps(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			resp := teov20220901.NewDescribeFunctionReplicasResponse()
			resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
				TotalCount: ptrInt64TeoReplicaV3(1),
				FunctionReplicas: []*teov20220901.FunctionReplica{
					{
						FunctionId:  ptrStringTeoReplicaV3(teoReplicaV3FunctionId),
						ReplicaName: ptrStringTeoReplicaV3(teoReplicaV3ReplicaName),
						Content:     ptrStringTeoReplicaV3("console.log(1)"),
					},
				},
			}
			return resp, nil
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	plan := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3CreatePlanAttrs("console.log(1)", nil)))

	resp := &resource.CreateResponse{State: tfsdk.State{Raw: plan, Schema: s}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Raw: plan, Schema: s}}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}

	if mock.gotCreate == nil {
		t.Fatal("expected CreateFunctionReplica to be called")
	}
	if mock.gotCreate.Remark != nil {
		t.Fatalf("expected a nil remark to be sent, got %q", *mock.gotCreate.Remark)
	}

	model := teoReplicaV3ModelFromState(t, s, resp.State.Raw)
	if !model.Remark.IsNull() {
		t.Fatalf("expected a null remark in state, got %q", model.Remark.ValueString())
	}
	if !model.CreatedOn.IsNull() {
		t.Fatalf("expected a null created_on, got %q", model.CreatedOn.ValueString())
	}
	if !model.ModifiedOn.IsNull() {
		t.Fatalf("expected a null modified_on, got %q", model.ModifiedOn.ValueString())
	}
}

// TestTeoReplicaV3_Create_APIError covers the create failure branch.
func TestTeoReplicaV3_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{createErr: errTeoReplicaV3("ResourceInUse.ReplicaNameDuplicated: replica already exists")}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	plan := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3CreatePlanAttrs("console.log(1)", nil)))

	resp := &resource.CreateResponse{State: tfsdk.State{Raw: plan, Schema: s}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Raw: plan, Schema: s}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Error creating TEO function replica" {
		t.Fatalf("unexpected error summary: %s", got)
	}
	if mock.gotDescribe != nil {
		t.Fatal("no read back must happen when the create failed")
	}
}

// TestTeoReplicaV3_Create_NotReadable covers an empty read back after a
// successful create: an error is reported but the identity is already in state.
func TestTeoReplicaV3_Create_NotReadable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	plan := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3CreatePlanAttrs("console.log(1)", nil)))

	resp := &resource.CreateResponse{State: tfsdk.State{Raw: plan, Schema: s}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Raw: plan, Schema: s}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "TEO function replica not readable after create" {
		t.Fatalf("unexpected error summary: %s", got)
	}
	if got := teoReplicaV3StateAttr(t, resp.State.Raw, "id").ValueString(); got != teoReplicaV3Id {
		t.Fatalf("the identity must be persisted before the read back, got %q", got)
	}
}

// TestTeoReplicaV3_Create_DescribeAPIError covers the read back error branch.
func TestTeoReplicaV3_Create_DescribeAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			return nil, errTeoReplicaV3("InternalError: describe failed")
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	plan := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3CreatePlanAttrs("console.log(1)", nil)))

	resp := &resource.CreateResponse{State: tfsdk.State{Raw: plan, Schema: s}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Raw: plan, Schema: s}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Error reading TEO function replica after create" {
		t.Fatalf("unexpected error summary: %s", got)
	}
	if got := teoReplicaV3StateAttr(t, resp.State.Raw, "id").ValueString(); got != teoReplicaV3Id {
		t.Fatalf("the identity must be persisted before the read back, got %q", got)
	}
}

// TestTeoReplicaV3_Read_Success covers a refresh that rewrites every leaf.
func TestTeoReplicaV3_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			return teoReplicaV3DescribeOne("console.log(2)", "note b"), nil
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(1)", "note a")))

	resp := &resource.ReadResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Read(context.Background(), resource.ReadRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}

	model := teoReplicaV3ModelFromState(t, s, resp.State.Raw)
	if got := model.Id.ValueString(); got != teoReplicaV3Id {
		t.Fatalf("unexpected id: %s", got)
	}
	if got := model.Content.ValueString(); got != "console.log(2)" {
		t.Fatalf("unexpected content: %s", got)
	}
	if got := model.Remark.ValueString(); got != "note b" {
		t.Fatalf("unexpected remark: %s", got)
	}
	if got := model.CreatedOn.ValueString(); got != "2025-01-01 00:00:00" {
		t.Fatalf("unexpected created_on: %s", got)
	}
	if got := model.ModifiedOn.ValueString(); got != "2025-01-02 00:00:00" {
		t.Fatalf("unexpected modified_on: %s", got)
	}
}

// TestTeoReplicaV3_Read_UsesExactNameMatch ensures the read ignores replicas
// whose name is only a partial match of the requested name.
func TestTeoReplicaV3_Read_UsesExactNameMatch(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			resp := teov20220901.NewDescribeFunctionReplicasResponse()
			resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
				TotalCount: ptrInt64TeoReplicaV3(1),
				FunctionReplicas: []*teov20220901.FunctionReplica{
					{
						FunctionId:  ptrStringTeoReplicaV3(teoReplicaV3FunctionId),
						ReplicaName: ptrStringTeoReplicaV3("replica-a-extra"),
						Content:     ptrStringTeoReplicaV3("other"),
					},
				},
			}
			return resp, nil
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(1)", "note a")))

	resp := &resource.ReadResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Read(context.Background(), resource.ReadRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("expected the resource to be removed from state: no exact name match")
	}
}

// TestTeoReplicaV3_Read_NotFound covers the confirmed "replica is gone" path.
func TestTeoReplicaV3_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			resp := teov20220901.NewDescribeFunctionReplicasResponse()
			resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
				TotalCount:       ptrInt64TeoReplicaV3(0),
				FunctionReplicas: []*teov20220901.FunctionReplica{},
			}
			return resp, nil
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(1)", "note a")))

	resp := &resource.ReadResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Read(context.Background(), resource.ReadRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("expected the resource to be removed from state when the replica no longer exists")
	}
}

// TestTeoReplicaV3_Read_NilResponseBody ensures a nil response body is treated
// as an API failure and never removes the resource from state.
func TestTeoReplicaV3_Read_NilResponseBody(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			return &teov20220901.DescribeFunctionReplicasResponse{}, nil
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(1)", "note a")))

	resp := &resource.ReadResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Read(context.Background(), resource.ReadRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics for a nil response body")
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("state must not be removed when the API returns a nil response body")
	}
}

// TestTeoReplicaV3_Read_APIError ensures an API error is reported and keeps the
// resource in state.
func TestTeoReplicaV3_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			return nil, errTeoReplicaV3("UnauthorizedOperation: no permission")
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(1)", "note a")))

	resp := &resource.ReadResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Read(context.Background(), resource.ReadRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Error reading TEO function replica" {
		t.Fatalf("unexpected error summary: %s", got)
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("state must not be removed on an API error")
	}
}

// TestTeoReplicaV3_Read_BrokenId covers the id parsing error branch.
func TestTeoReplicaV3_Read_BrokenId(t *testing.T) {
	r := setTeoReplicaV3Client(&connectivity.TencentCloudClient{})

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, map[string]attr.Value{
		"id": types.StringValue("broken-id"),
	}))

	resp := &resource.ReadResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Read(context.Background(), resource.ReadRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics for a broken id")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Invalid resource ID" {
		t.Fatalf("unexpected error summary: %s", got)
	}
}

// TestTeoReplicaV3_Update_Success covers an in-place update of content and
// remark. The plan carries different values on purpose: the request must be
// built from the plan, not from the previous state.
func TestTeoReplicaV3_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			return teoReplicaV3DescribeOne("console.log(3)", "note c"), nil
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(2)", "note b")))
	plan := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(3)", "note c")))

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Update(context.Background(), resource.UpdateRequest{
		State: tfsdk.State{Raw: state, Schema: s},
		Plan:  tfsdk.Plan{Raw: plan, Schema: s},
	}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}

	if mock.gotModify == nil {
		t.Fatal("expected ModifyFunctionReplica to be called")
	}
	if v := mock.gotModify.ZoneId; v == nil || *v != teoReplicaV3ZoneId {
		t.Fatalf("unexpected ZoneId: %v", v)
	}
	if v := mock.gotModify.FunctionId; v == nil || *v != teoReplicaV3FunctionId {
		t.Fatalf("unexpected FunctionId: %v", v)
	}
	if v := mock.gotModify.ReplicaName; v == nil || *v != teoReplicaV3ReplicaName {
		t.Fatalf("unexpected ReplicaName: %v", v)
	}
	if v := mock.gotModify.Content; v == nil || *v != "console.log(3)" {
		t.Fatalf("unexpected Content: %v", v)
	}
	if v := mock.gotModify.Remark; v == nil || *v != "note c" {
		t.Fatalf("unexpected Remark: %v", v)
	}

	model := teoReplicaV3ModelFromState(t, s, resp.State.Raw)
	if got := model.Id.ValueString(); got != teoReplicaV3Id {
		t.Fatalf("unexpected id: %s", got)
	}
	if got := model.Content.ValueString(); got != "console.log(3)" {
		t.Fatalf("unexpected content: %s", got)
	}
	if got := model.Remark.ValueString(); got != "note c" {
		t.Fatalf("unexpected remark: %s", got)
	}
}

// TestTeoReplicaV3_Update_ClearedRemark covers removing remark from the
// configuration: an explicit empty value is sent and the state settles to null,
// so the next plan shows no diff.
func TestTeoReplicaV3_Update_ClearedRemark(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{
		describe: func(_ *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
			return teoReplicaV3DescribeOne("console.log(2)", ""), nil
		},
	}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(2)", "note b")))
	plan := teoReplicaV3Raw(t, teoReplicaV3Model(t, map[string]attr.Value{
		"id":           types.StringValue(teoReplicaV3Id),
		"zone_id":      types.StringValue(teoReplicaV3ZoneId),
		"function_id":  types.StringValue(teoReplicaV3FunctionId),
		"replica_name": types.StringValue(teoReplicaV3ReplicaName),
		"content":      types.StringValue("console.log(2)"),
	}))

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Update(context.Background(), resource.UpdateRequest{
		State: tfsdk.State{Raw: state, Schema: s},
		Plan:  tfsdk.Plan{Raw: plan, Schema: s},
	}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}

	if mock.gotModify == nil {
		t.Fatal("expected ModifyFunctionReplica to be called")
	}
	if mock.gotModify.Remark == nil || *mock.gotModify.Remark != "" {
		t.Fatalf("expected an explicit empty remark to be sent, got %v", mock.gotModify.Remark)
	}

	model := teoReplicaV3ModelFromState(t, s, resp.State.Raw)
	if !model.Remark.IsNull() {
		t.Fatalf("expected a null remark after clearing, got %q", model.Remark.ValueString())
	}
}

// TestTeoReplicaV3_Update_APIError covers the update failure branch.
func TestTeoReplicaV3_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{modifyErr: errTeoReplicaV3("InternalError: modify failed")}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(2)", "note b")))
	plan := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(3)", "note c")))

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Update(context.Background(), resource.UpdateRequest{
		State: tfsdk.State{Raw: state, Schema: s},
		Plan:  tfsdk.Plan{Raw: plan, Schema: s},
	}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Error updating TEO function replica" {
		t.Fatalf("unexpected error summary: %s", got)
	}
	if mock.gotDescribe != nil {
		t.Fatal("no read back must happen when the update failed")
	}
}

// TestTeoReplicaV3_Update_NotReadable covers the read back failure after a
// successful update.
func TestTeoReplicaV3_Update_NotReadable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	raw := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(2)", "note b")))

	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: raw, Schema: s}}
	r.Update(context.Background(), resource.UpdateRequest{
		State: tfsdk.State{Raw: raw, Schema: s},
		Plan:  tfsdk.Plan{Raw: raw, Schema: s},
	}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "TEO function replica not readable after update" {
		t.Fatalf("unexpected error summary: %s", got)
	}
}

// TestTeoReplicaV3_Delete_Success asserts that exactly one replica name is sent
// so that destroy never affects the sibling replicas of the same function.
func TestTeoReplicaV3_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(1)", "note a")))

	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Delete(context.Background(), resource.DeleteRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}

	if mock.gotDelete == nil {
		t.Fatal("expected DeleteFunctionReplica to be called")
	}
	if v := mock.gotDelete.ZoneId; v == nil || *v != teoReplicaV3ZoneId {
		t.Fatalf("unexpected ZoneId: %v", v)
	}
	if v := mock.gotDelete.FunctionId; v == nil || *v != teoReplicaV3FunctionId {
		t.Fatalf("unexpected FunctionId: %v", v)
	}
	if len(mock.gotDelete.ReplicaNames) != 1 {
		t.Fatalf("expected exactly 1 replica name, got %d", len(mock.gotDelete.ReplicaNames))
	}
	if v := mock.gotDelete.ReplicaNames[0]; v == nil || *v != teoReplicaV3ReplicaName {
		t.Fatalf("unexpected replica name: %v", v)
	}
}

// TestTeoReplicaV3_Delete_APIError covers the delete failure branch.
func TestTeoReplicaV3_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	mock := &teoReplicaV3Mock{deleteErr: errTeoReplicaV3("UnauthorizedOperation: not allowed")}
	mock.apply(patches, teoClient)

	r := setTeoReplicaV3Client(client)

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(1)", "note a")))

	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Delete(context.Background(), resource.DeleteRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Error deleting TEO function replica" {
		t.Fatalf("unexpected error summary: %s", got)
	}
}

// TestTeoReplicaV3_Delete_BrokenId covers the id parsing error branch of delete.
func TestTeoReplicaV3_Delete_BrokenId(t *testing.T) {
	r := setTeoReplicaV3Client(&connectivity.TencentCloudClient{})

	s := teoReplicaV3Schema(t)
	state := teoReplicaV3Raw(t, teoReplicaV3Model(t, map[string]attr.Value{
		"id": types.StringValue("broken-id"),
	}))

	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: state, Schema: s}}
	r.Delete(context.Background(), resource.DeleteRequest{State: tfsdk.State{Raw: state, Schema: s}}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics for a broken id")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Invalid resource ID" {
		t.Fatalf("unexpected error summary: %s", got)
	}
}

// TestTeoReplicaV3_ImportState_Success covers the composite id import.
func TestTeoReplicaV3_ImportState_Success(t *testing.T) {
	r := teo.NewTeoFunctionReplicaV3Resource()

	s := teoReplicaV3Schema(t)
	empty := teoReplicaV3Raw(t, teoReplicaV3Model(t, nil))

	resp := &resource.ImportStateResponse{State: tfsdk.State{Raw: empty, Schema: s}}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: teoReplicaV3Id}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}

	model := teoReplicaV3ModelFromState(t, s, resp.State.Raw)
	if got := model.Id.ValueString(); got != teoReplicaV3Id {
		t.Fatalf("unexpected id: %s", got)
	}
	if got := model.ZoneId.ValueString(); got != teoReplicaV3ZoneId {
		t.Fatalf("unexpected zone_id: %s", got)
	}
	if got := model.FunctionId.ValueString(); got != teoReplicaV3FunctionId {
		t.Fatalf("unexpected function_id: %s", got)
	}
	if got := model.ReplicaName.ValueString(); got != teoReplicaV3ReplicaName {
		t.Fatalf("unexpected replica_name: %s", got)
	}
}

// TestTeoReplicaV3_ImportState_InvalidId covers the id parsing error branch.
func TestTeoReplicaV3_ImportState_InvalidId(t *testing.T) {
	r := teo.NewTeoFunctionReplicaV3Resource()

	s := teoReplicaV3Schema(t)
	empty := teoReplicaV3Raw(t, teoReplicaV3Model(t, nil))

	cases := []string{
		"",
		"zone-2qtuhspy7cr6",
		"zone-2qtuhspy7cr6#ef-2qlxy8s7o96e",
		"zone-2qtuhspy7cr6#ef-2qlxy8s7o96e#replica-a#extra",
		"zone-2qtuhspy7cr6##replica-a",
		"#ef-2qlxy8s7o96e#replica-a",
	}
	for _, c := range cases {
		resp := &resource.ImportStateResponse{State: tfsdk.State{Raw: empty, Schema: s}}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: c}, resp)
		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected error diagnostics for invalid import id %q", c)
		}
		if got := resp.Diagnostics.Errors()[0].Summary(); got != "Invalid import ID" {
			t.Fatalf("unexpected error summary for %q: %s", c, got)
		}
	}
}

// TestTeoReplicaV3_ClientNotConfigured covers the nil client guard of every
// lifecycle callback.
func TestTeoReplicaV3_ClientNotConfigured(t *testing.T) {
	r := teo.NewTeoFunctionReplicaV3Resource()

	s := teoReplicaV3Schema(t)
	raw := teoReplicaV3Raw(t, teoReplicaV3Model(t, teoReplicaV3FullAttrs("console.log(1)", "note a")))

	createResp := &resource.CreateResponse{State: tfsdk.State{Raw: raw, Schema: s}}
	r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Raw: raw, Schema: s}}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on create without a configured client")
	}
	if got := createResp.Diagnostics.Errors()[0].Summary(); got != "Provider not configured" {
		t.Fatalf("unexpected error summary: %s", got)
	}

	readResp := &resource.ReadResponse{State: tfsdk.State{Raw: raw, Schema: s}}
	r.Read(context.Background(), resource.ReadRequest{State: tfsdk.State{Raw: raw, Schema: s}}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on read without a configured client")
	}
	if got := readResp.Diagnostics.Errors()[0].Summary(); got != "Provider not configured" {
		t.Fatalf("unexpected error summary: %s", got)
	}

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Raw: raw, Schema: s}}
	r.Update(context.Background(), resource.UpdateRequest{
		State: tfsdk.State{Raw: raw, Schema: s},
		Plan:  tfsdk.Plan{Raw: raw, Schema: s},
	}, updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on update without a configured client")
	}
	if got := updateResp.Diagnostics.Errors()[0].Summary(); got != "Provider not configured" {
		t.Fatalf("unexpected error summary: %s", got)
	}

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Raw: raw, Schema: s}}
	r.Delete(context.Background(), resource.DeleteRequest{State: tfsdk.State{Raw: raw, Schema: s}}, deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on delete without a configured client")
	}
	if got := deleteResp.Diagnostics.Errors()[0].Summary(); got != "Provider not configured" {
		t.Fatalf("unexpected error summary: %s", got)
	}
}

// TestTeoReplicaV3_Factory validates the registry factory.
func TestTeoReplicaV3_Factory(t *testing.T) {
	got := teo.NewTeoFunctionReplicaV3Resource()
	if got == nil {
		t.Fatal("NewTeoFunctionReplicaV3Resource returned nil")
	}
	if _, ok := got.(*teo.TeoFunctionReplicaV3Resource); !ok {
		t.Fatalf("factory returned %T, want *teo.TeoFunctionReplicaV3Resource", got)
	}
}
