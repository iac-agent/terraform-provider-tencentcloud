package teo

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

func ptrStringFunctionReplicaV4(s string) *string {
	return &s
}

func ptrInt64FunctionReplicaV4(i int64) *int64 {
	return &i
}

// teoFunctionReplicaV4ResourceSchema builds the schema used to construct the
// raw plan/state values passed to the framework lifecycle methods.
func teoFunctionReplicaV4ResourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := NewTeoFunctionReplicaV4Resource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// teoFunctionReplicaV4ObjectType derives the tftypes object used for raw
// plan/state construction.
func teoFunctionReplicaV4ObjectType(t *testing.T) tftypes.Object {
	t.Helper()
	return teoFunctionReplicaV4ResourceSchema(t).Type().TerraformType(context.Background()).(tftypes.Object)
}

// newTeoFunctionReplicaV4PlanValue builds a fully known raw value.
func newTeoFunctionReplicaV4PlanValue(t *testing.T, values map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	return tftypes.NewValue(teoFunctionReplicaV4ObjectType(t), values)
}

// newTeoFunctionReplicaV4Plan constructs a resource.CreateRequest/UpdateRequest
// style plan carrying the given attribute values.
func newTeoFunctionReplicaV4Plan(t *testing.T, id string, m TeoFunctionReplicaV4ResourceModel) tfsdk.Plan {
	t.Helper()

	idVal := tftypes.NewValue(tftypes.String, nil)
	if id != "" {
		idVal = tftypes.NewValue(tftypes.String, id)
	}

	stringOrNull := func(v types.String) tftypes.Value {
		if v.IsNull() {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, v.ValueString())
	}

	raw := newTeoFunctionReplicaV4PlanValue(t, map[string]tftypes.Value{
		"id":           idVal,
		"zone_id":      tftypes.NewValue(tftypes.String, m.ZoneId.ValueString()),
		"function_id":  tftypes.NewValue(tftypes.String, m.FunctionId.ValueString()),
		"replica_name": tftypes.NewValue(tftypes.String, m.ReplicaName.ValueString()),
		"content":      stringOrNull(m.Content),
		"remark":       stringOrNull(m.Remark),
	})

	return tfsdk.Plan{Schema: teoFunctionReplicaV4ResourceSchema(t), Raw: raw}
}

// newTeoFunctionReplicaV4State constructs the state passed to Read/Update/Delete.
func newTeoFunctionReplicaV4State(t *testing.T, id string, m TeoFunctionReplicaV4ResourceModel) tfsdk.State {
	t.Helper()
	plan := newTeoFunctionReplicaV4Plan(t, id, m)
	return tfsdk.State{Schema: plan.Schema, Raw: plan.Raw}
}

// newTeoFunctionReplicaV4EmptyState builds the state value written by the
// lifecycle methods before any apply, so a state response can be constructed.
func newTeoFunctionReplicaV4EmptyState(t *testing.T) tfsdk.State {
	t.Helper()

	s := teoFunctionReplicaV4ResourceSchema(t)
	nullString := tftypes.NewValue(tftypes.String, nil)
	raw := newTeoFunctionReplicaV4PlanValue(t, map[string]tftypes.Value{
		"id":           nullString,
		"zone_id":      nullString,
		"function_id":  nullString,
		"replica_name": nullString,
		"content":      nullString,
		"remark":       nullString,
	})

	return tfsdk.State{Schema: s, Raw: raw}
}

// setTeoFunctionReplicaV4Client injects the mock connectivity client through
// the Configure method, exactly like the framework provider does at runtime.
func setTeoFunctionReplicaV4Client(r *TeoFunctionReplicaV4Resource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: meta}, &resource.ConfigureResponse{})
}

// mockTeoFunctionReplicaV4Client returns a client whose teo sdk client is patched.
func mockTeoFunctionReplicaV4Client(patches *gomonkey.Patches) (*connectivity.TencentCloudClient, *teov20220901.Client) {
	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)
	return client, teoClient
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV4" -v -count=1 -gcflags="all=-l"

// TestTeoFunctionReplicaV4_MetadataAndSchema validates the type name and the
// Required/Optional/Computed flags of every attribute.
func TestTeoFunctionReplicaV4_MetadataAndSchema(t *testing.T) {
	r := NewTeoFunctionReplicaV4Resource()

	metaResp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, metaResp)
	if metaResp.TypeName != "tencentcloud_teo_function_replica_v4" {
		t.Errorf("Metadata.TypeName = %q, want %q", metaResp.TypeName, "tencentcloud_teo_function_replica_v4")
	}

	s := teoFunctionReplicaV4ResourceSchema(t)

	type expect struct {
		required bool
		optional bool
		computed bool
	}
	cases := map[string]expect{
		"id":           {computed: true},
		"zone_id":      {required: true},
		"function_id":  {required: true},
		"replica_name": {required: true},
		"content":      {required: true},
		"remark":       {optional: true},
	}

	if len(s.Attributes) != len(cases) {
		t.Errorf("attribute count = %d, want %d", len(s.Attributes), len(cases))
	}

	for name, want := range cases {
		attr, ok := s.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q missing or not a StringAttribute", name)
			continue
		}
		if attr.Required != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.Required, want.required)
		}
		if attr.Optional != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, attr.Optional, want.optional)
		}
		if attr.Computed != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, attr.Computed, want.computed)
		}
	}
}

// TestTeoFunctionReplicaV4_FactoryAndConfigure covers the factory and the
// nil-ProviderData Configure path.
func TestTeoFunctionReplicaV4_FactoryAndConfigure(t *testing.T) {
	r, ok := NewTeoFunctionReplicaV4Resource().(*TeoFunctionReplicaV4Resource)
	if !ok {
		t.Fatalf("factory returned %T, want *TeoFunctionReplicaV4Resource", NewTeoFunctionReplicaV4Resource())
	}

	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
	if r.Client() != nil {
		t.Errorf("Client() = %v, want nil before Configure", r.Client())
	}
}

// TestSplitTeoFunctionReplicaV4Id validates composite id building and parsing.
func TestSplitTeoFunctionReplicaV4Id(t *testing.T) {
	id := buildTeoFunctionReplicaV4Id("zone-1", "ef-2", "replica-3")
	want := "zone-1" + tccommon.FILED_SP + "ef-2" + tccommon.FILED_SP + "replica-3"
	if id != want {
		t.Fatalf("buildTeoFunctionReplicaV4Id = %q, want %q", id, want)
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV4Id(id)
	if !ok || zoneId != "zone-1" || functionId != "ef-2" || replicaName != "replica-3" {
		t.Errorf("splitTeoFunctionReplicaV4Id(%q) = (%q, %q, %q, %v)", id, zoneId, functionId, replicaName, ok)
	}

	bad := []string{
		"",
		"zone-1",
		"zone-1" + tccommon.FILED_SP + "ef-2",
		tccommon.FILED_SP + "ef-2" + tccommon.FILED_SP + "replica-3",
		"zone-1" + tccommon.FILED_SP + tccommon.FILED_SP + "replica-3",
		"zone-1" + tccommon.FILED_SP + "ef-2" + tccommon.FILED_SP + "replica-3" + tccommon.FILED_SP + "extra",
	}
	for _, in := range bad {
		if _, _, _, ok := splitTeoFunctionReplicaV4Id(in); ok {
			t.Errorf("splitTeoFunctionReplicaV4Id(%q) ok=true, want false", in)
		}
	}
}

// TestTeoFunctionReplicaV4_Create covers a successful create, asserting the
// request mapping, the composite id and every leaf field read back into state.
func TestTeoFunctionReplicaV4_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)

	var gotCreate *teov20220901.CreateFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		gotCreate = request
		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	var gotDescribe *teov20220901.DescribeFunctionReplicasRequest
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		gotDescribe = request
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV4(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV4("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	plan := newTeoFunctionReplicaV4Plan(t, "", TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringValue("test remark"),
	})

	resp := &resource.CreateResponse{State: newTeoFunctionReplicaV4EmptyState(t)}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	if gotCreate == nil {
		t.Fatal("CreateFunctionReplica was not called")
	}
	if gotCreate.ZoneId == nil || *gotCreate.ZoneId != "zone-test123" {
		t.Errorf("request ZoneId = %v, want zone-test123", gotCreate.ZoneId)
	}
	if gotCreate.FunctionId == nil || *gotCreate.FunctionId != "ef-test456" {
		t.Errorf("request FunctionId = %v, want ef-test456", gotCreate.FunctionId)
	}
	if gotCreate.ReplicaName == nil || *gotCreate.ReplicaName != "test-replica" {
		t.Errorf("request ReplicaName = %v, want test-replica", gotCreate.ReplicaName)
	}
	if gotCreate.Content == nil || *gotCreate.Content != "console.log(123)" {
		t.Errorf("request Content = %v, want console.log(123)", gotCreate.Content)
	}
	if gotCreate.Remark == nil || *gotCreate.Remark != "test remark" {
		t.Errorf("request Remark = %v, want test remark", gotCreate.Remark)
	}

	if gotDescribe == nil {
		t.Fatal("DescribeFunctionReplicas was not called")
	}
	if gotDescribe.ZoneId == nil || *gotDescribe.ZoneId != "zone-test123" {
		t.Errorf("describe ZoneId = %v, want zone-test123", gotDescribe.ZoneId)
	}
	if gotDescribe.FunctionId == nil || *gotDescribe.FunctionId != "ef-test456" {
		t.Errorf("describe FunctionId = %v, want ef-test456", gotDescribe.FunctionId)
	}
	if len(gotDescribe.Filters) != 1 || gotDescribe.Filters[0].Name == nil || *gotDescribe.Filters[0].Name != "replica-name" {
		t.Fatalf("describe Filters = %v, want a single replica-name filter", gotDescribe.Filters)
	}
	if len(gotDescribe.Filters[0].Values) != 1 || gotDescribe.Filters[0].Values[0] == nil || *gotDescribe.Filters[0].Values[0] != "test-replica" {
		t.Errorf("describe filter values = %v, want [test-replica]", gotDescribe.Filters[0].Values)
	}

	var state TeoFunctionReplicaV4ResourceModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}

	wantId := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	if state.Id.ValueString() != wantId {
		t.Errorf("state id = %q, want %q", state.Id.ValueString(), wantId)
	}
	if state.ZoneId.ValueString() != "zone-test123" {
		t.Errorf("state zone_id = %q, want %q", state.ZoneId.ValueString(), "zone-test123")
	}
	if state.FunctionId.ValueString() != "ef-test456" {
		t.Errorf("state function_id = %q, want %q", state.FunctionId.ValueString(), "ef-test456")
	}
	if state.ReplicaName.ValueString() != "test-replica" {
		t.Errorf("state replica_name = %q, want %q", state.ReplicaName.ValueString(), "test-replica")
	}
	if state.Content.ValueString() != "console.log(123)" {
		t.Errorf("state content = %q, want %q", state.Content.ValueString(), "console.log(123)")
	}
	if state.Remark.IsNull() || state.Remark.ValueString() != "test remark" {
		t.Errorf("state remark = %v, want %q", state.Remark, "test remark")
	}
}

// TestTeoFunctionReplicaV4_Create_WithoutRemark covers an omitted remark: it is
// not sent on create and a null remark returned by the cloud stays null.
func TestTeoFunctionReplicaV4_Create_WithoutRemark(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)

	var gotCreate *teov20220901.CreateFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		gotCreate = request
		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(123)"),
				},
			},
		}
		return resp, nil
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	plan := newTeoFunctionReplicaV4Plan(t, "", TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.CreateResponse{State: newTeoFunctionReplicaV4EmptyState(t)}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	if gotCreate == nil {
		t.Fatal("CreateFunctionReplica was not called")
	}
	if gotCreate.Remark != nil {
		t.Errorf("request Remark = %v, want nil for an omitted remark", gotCreate.Remark)
	}

	var state TeoFunctionReplicaV4ResourceModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if !state.Remark.IsNull() {
		t.Errorf("state remark = %v, want null", state.Remark)
	}
}

// TestTeoFunctionReplicaV4_Create_APIError covers the create error branch.
func TestTeoFunctionReplicaV4_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)
	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("FailedOperation.InvalidReplicaName")
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	plan := newTeoFunctionReplicaV4Plan(t, "", TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.CreateResponse{State: newTeoFunctionReplicaV4EmptyState(t)}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic, got none")
	}
	if summary := resp.Diagnostics.Errors()[0].Summary(); summary != "Error creating TEO function replica" {
		t.Errorf("error summary = %q, want %q", summary, "Error creating TEO function replica")
	}
	if detail := resp.Diagnostics.Errors()[0].Detail(); detail == "" {
		t.Error("error detail is empty, want the API failure reason")
	}
}

// TestTeoFunctionReplicaV4_Read covers a successful read of every leaf field.
func TestTeoFunctionReplicaV4_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV4(2),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica-other"),
					Content:     ptrStringFunctionReplicaV4("other content"),
					Remark:      ptrStringFunctionReplicaV4("other remark"),
				},
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV4("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV4("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("stale content"),
		Remark:      types.StringValue("stale remark"),
	})

	resp := &resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	var got TeoFunctionReplicaV4ResourceModel
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}

	if got.Id.ValueString() != id {
		t.Errorf("state id = %q, want %q", got.Id.ValueString(), id)
	}
	if got.ZoneId.ValueString() != "zone-test123" {
		t.Errorf("state zone_id = %q, want %q", got.ZoneId.ValueString(), "zone-test123")
	}
	if got.FunctionId.ValueString() != "ef-test456" {
		t.Errorf("state function_id = %q, want %q", got.FunctionId.ValueString(), "ef-test456")
	}
	if got.ReplicaName.ValueString() != "test-replica" {
		t.Errorf("state replica_name = %q, want %q", got.ReplicaName.ValueString(), "test-replica")
	}
	if got.Content.ValueString() != "console.log(123)" {
		t.Errorf("state content = %q, want %q", got.Content.ValueString(), "console.log(123)")
	}
	if got.Remark.IsNull() || got.Remark.ValueString() != "test remark" {
		t.Errorf("state remark = %v, want %q", got.Remark, "test remark")
	}
}

// TestTeoFunctionReplicaV4_Read_NotFound covers a successful query that returns
// no exact replica-name match: the resource is removed from state instead of
// being reported as an API failure.
func TestTeoFunctionReplicaV4_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64FunctionReplicaV4(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}
	if len(resp.Diagnostics.Warnings()) != 1 {
		t.Errorf("warning count = %d, want 1", len(resp.Diagnostics.Warnings()))
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("state was not removed for a missing replica: %v", resp.State.Raw)
	}
}

// TestTeoFunctionReplicaV4_Read_NilResponse covers the "nil response is an API
// failure" rule: it must not be treated as "not found".
func TestTeoFunctionReplicaV4_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return &teov20220901.DescribeFunctionReplicasResponse{}, nil
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for a nil response")
	}
	if resp.Diagnostics.Errors()[0].Summary() != "Error reading TEO function replica" {
		t.Errorf("error summary = %q, want %q", resp.Diagnostics.Errors()[0].Summary(), "Error reading TEO function replica")
	}
	if resp.State.Raw.IsNull() {
		t.Error("state must be kept when the query failed")
	}
}

// TestTeoFunctionReplicaV4_Read_InvalidId covers a broken composite id.
func TestTeoFunctionReplicaV4_Read_InvalidId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, _ := mockTeoFunctionReplicaV4Client(patches)
	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	state := newTeoFunctionReplicaV4State(t, "broken-id", TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for a broken id")
	}
	if resp.Diagnostics.Errors()[0].Summary() != "Invalid resource id" {
		t.Errorf("error summary = %q, want %q", resp.Diagnostics.Errors()[0].Summary(), "Invalid resource id")
	}
}

// TestTeoFunctionReplicaV4_Update covers the in-place update of content and
// remark.
func TestTeoFunctionReplicaV4_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)

	var gotModify *teov20220901.ModifyFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		gotModify = request
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(456)"),
					Remark:      ptrStringFunctionReplicaV4("updated remark"),
				},
			},
		}
		return resp, nil
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringValue("test remark"),
	})
	plan := newTeoFunctionReplicaV4Plan(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(456)"),
		Remark:      types.StringValue("updated remark"),
	})

	resp := &resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	if gotModify == nil {
		t.Fatal("ModifyFunctionReplica was not called")
	}
	if gotModify.ZoneId == nil || *gotModify.ZoneId != "zone-test123" {
		t.Errorf("request ZoneId = %v, want zone-test123", gotModify.ZoneId)
	}
	if gotModify.FunctionId == nil || *gotModify.FunctionId != "ef-test456" {
		t.Errorf("request FunctionId = %v, want ef-test456", gotModify.FunctionId)
	}
	if gotModify.ReplicaName == nil || *gotModify.ReplicaName != "test-replica" {
		t.Errorf("request ReplicaName = %v, want test-replica", gotModify.ReplicaName)
	}
	if gotModify.Content == nil || *gotModify.Content != "console.log(456)" {
		t.Errorf("request Content = %v, want console.log(456)", gotModify.Content)
	}
	if gotModify.Remark == nil || *gotModify.Remark != "updated remark" {
		t.Errorf("request Remark = %v, want updated remark", gotModify.Remark)
	}

	var got TeoFunctionReplicaV4ResourceModel
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if got.Content.ValueString() != "console.log(456)" {
		t.Errorf("state content = %q, want %q", got.Content.ValueString(), "console.log(456)")
	}
	if got.Remark.IsNull() || got.Remark.ValueString() != "updated remark" {
		t.Errorf("state remark = %v, want %q", got.Remark, "updated remark")
	}
}

// TestTeoFunctionReplicaV4_Update_RemarkRemovedFromConfig covers the documented
// clear contract: a null remark is not sent to the API and the cloud value is
// kept, so no permanent diff is produced.
func TestTeoFunctionReplicaV4_Update_RemarkRemovedFromConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)

	var gotModify *teov20220901.ModifyFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		gotModify = request
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(789)"),
					Remark:      ptrStringFunctionReplicaV4("kept remark"),
				},
			},
		}
		return resp, nil
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringValue("kept remark"),
	})
	plan := newTeoFunctionReplicaV4Plan(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(789)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	if gotModify == nil {
		t.Fatal("ModifyFunctionReplica was not called")
	}
	if gotModify.Remark != nil {
		t.Errorf("request Remark = %v, want nil when the attribute is removed from the configuration", gotModify.Remark)
	}

	var got TeoFunctionReplicaV4ResourceModel
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if got.Remark.IsNull() || got.Remark.ValueString() != "kept remark" {
		t.Errorf("state remark = %v, want %q", got.Remark, "kept remark")
	}
}

// TestTeoFunctionReplicaV4_Update_APIError covers the update error branch.
func TestTeoFunctionReplicaV4_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("InternalError")
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})
	plan := newTeoFunctionReplicaV4Plan(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(456)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic, got none")
	}
	if resp.Diagnostics.Errors()[0].Summary() != "Error updating TEO function replica" {
		t.Errorf("error summary = %q, want %q", resp.Diagnostics.Errors()[0].Summary(), "Error updating TEO function replica")
	}
}

// TestTeoFunctionReplicaV4_Delete covers a successful delete of the single
// replica this resource manages.
func TestTeoFunctionReplicaV4_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)

	var gotDelete *teov20220901.DeleteFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		gotDelete = request
		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}

	if gotDelete == nil {
		t.Fatal("DeleteFunctionReplica was not called")
	}
	if gotDelete.ZoneId == nil || *gotDelete.ZoneId != "zone-test123" {
		t.Errorf("request ZoneId = %v, want zone-test123", gotDelete.ZoneId)
	}
	if gotDelete.FunctionId == nil || *gotDelete.FunctionId != "ef-test456" {
		t.Errorf("request FunctionId = %v, want ef-test456", gotDelete.FunctionId)
	}
	if len(gotDelete.ReplicaNames) != 1 {
		t.Fatalf("request ReplicaNames length = %d, want 1", len(gotDelete.ReplicaNames))
	}
	if gotDelete.ReplicaNames[0] == nil || *gotDelete.ReplicaNames[0] != "test-replica" {
		t.Errorf("request ReplicaNames[0] = %v, want test-replica", gotDelete.ReplicaNames[0])
	}
}

// TestTeoFunctionReplicaV4_Delete_APIError covers the delete error branch: the
// failure is surfaced so the id stays in state for a retry.
func TestTeoFunctionReplicaV4_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplicaV4Client(patches)
	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("ResourceInUse")
	})

	r := &TeoFunctionReplicaV4Resource{}
	setTeoFunctionReplicaV4Client(r, client)

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})

	resp := &resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic, got none")
	}
	if resp.Diagnostics.Errors()[0].Summary() != "Error deleting TEO function replica" {
		t.Errorf("error summary = %q, want %q", resp.Diagnostics.Errors()[0].Summary(), "Error deleting TEO function replica")
	}
}

// TestTeoFunctionReplicaV4_ImportState covers the composite identity import for
// both the valid and the broken id form.
func TestTeoFunctionReplicaV4_ImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r := &TeoFunctionReplicaV4Resource{}

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	importResp := &resource.ImportStateResponse{State: newTeoFunctionReplicaV4EmptyState(t)}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", importResp.Diagnostics)
	}

	var state TeoFunctionReplicaV4ResourceModel
	if diags := importResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if state.Id.ValueString() != id {
		t.Errorf("imported id = %q, want %q", state.Id.ValueString(), id)
	}
	if state.ZoneId.ValueString() != "zone-test123" {
		t.Errorf("imported zone_id = %q, want %q", state.ZoneId.ValueString(), "zone-test123")
	}
	if state.FunctionId.ValueString() != "ef-test456" {
		t.Errorf("imported function_id = %q, want %q", state.FunctionId.ValueString(), "ef-test456")
	}
	if state.ReplicaName.ValueString() != "test-replica" {
		t.Errorf("imported replica_name = %q, want %q", state.ReplicaName.ValueString(), "test-replica")
	}

	brokenResp := &resource.ImportStateResponse{State: newTeoFunctionReplicaV4EmptyState(t)}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "zone-test123"}, brokenResp)
	if !brokenResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for a broken import id")
	}
	if brokenResp.Diagnostics.Errors()[0].Summary() != "Invalid import id" {
		t.Errorf("error summary = %q, want %q", brokenResp.Diagnostics.Errors()[0].Summary(), "Invalid import id")
	}
}

// TestTeoFunctionReplicaV4_ClientNotConfigured covers the nil-client guard of
// every lifecycle method.
func TestTeoFunctionReplicaV4_ClientNotConfigured(t *testing.T) {
	r := &TeoFunctionReplicaV4Resource{}

	id := buildTeoFunctionReplicaV4Id("zone-test123", "ef-test456", "test-replica")
	state := newTeoFunctionReplicaV4State(t, id, TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})
	plan := newTeoFunctionReplicaV4Plan(t, "", TeoFunctionReplicaV4ResourceModel{
		ZoneId:      types.StringValue("zone-test123"),
		FunctionId:  types.StringValue("ef-test456"),
		ReplicaName: types.StringValue("test-replica"),
		Content:     types.StringValue("console.log(123)"),
		Remark:      types.StringNull(),
	})

	createResp := &resource.CreateResponse{State: newTeoFunctionReplicaV4EmptyState(t)}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Error("Create: expected an error diagnostic for a nil client")
	}

	readResp := &resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Error("Read: expected an error diagnostic for a nil client")
	}

	updateResp := &resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Error("Update: expected an error diagnostic for a nil client")
	}

	deleteResp := &resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Error("Delete: expected an error diagnostic for a nil client")
	}
}
