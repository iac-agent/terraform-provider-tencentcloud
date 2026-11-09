package teo_test

import (
	"context"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/stretchr/testify/assert"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	providerTeo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// ptrStringTeoFunctionReplica110901 is a helper to build a *string.
func ptrStringTeoFunctionReplica110901(s string) *string { return &s }

// ptrInt64TeoFunctionReplica110901 is a helper to build a *int64.
func ptrInt64TeoFunctionReplica110901(i int64) *int64 { return &i }

// teoFunctionReplica110901ObjectType is the tftypes object type of the
// resource (mirrors the framework schema).
func teoFunctionReplica110901ObjectType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                tftypes.String,
			"zone_id":           tftypes.String,
			"function_id":       tftypes.String,
			"replica_name":      tftypes.String,
			"content":           tftypes.String,
			"remark":            tftypes.String,
			"function_replicas": tftypes.List{ElementType: teoFunctionReplica110901ItemObjectType()},
		},
	}
}

// teoFunctionReplica110901ItemObjectType is the tftypes object type of one
// function_replicas element.
func teoFunctionReplica110901ItemObjectType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"function_id":  tftypes.String,
			"replica_name": tftypes.String,
			"content":      tftypes.String,
			"remark":       tftypes.String,
			"created_on":   tftypes.String,
			"modified_on":  tftypes.String,
		},
	}
}

// teoFunctionReplica110901RawValue builds the raw tftypes value of the whole
// resource object. A nil string becomes a null attribute.
func teoFunctionReplica110901RawValue(id, zoneId, functionId, replicaName, content, remark *string) tftypes.Value {
	str := func(p *string) tftypes.Value {
		if p == nil {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, *p)
	}
	return tftypes.NewValue(teoFunctionReplica110901ObjectType(), map[string]tftypes.Value{
		"id":                str(id),
		"zone_id":           str(zoneId),
		"function_id":       str(functionId),
		"replica_name":      str(replicaName),
		"content":           str(content),
		"remark":            str(remark),
		"function_replicas": tftypes.NewValue(tftypes.List{ElementType: teoFunctionReplica110901ItemObjectType()}, nil),
	})
}

// newTeoFunctionReplica110901Resource returns a configured resource instance.
func newTeoFunctionReplica110901Resource(t *testing.T, client *connectivity.TencentCloudClient) resource.Resource {
	r := providerTeo.NewTeoFunctionReplica110901Resource()
	configureReq := resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("Configure returned diagnostics: %v", configureResp.Diagnostics)
	}
	return r
}

// newTeoFunctionReplica110901DescribeResponse builds a Describe response
// carrying one replica with the given fields.
func newTeoFunctionReplica110901DescribeResponse(replicaName, content, remark string) *teov20220901.DescribeFunctionReplicasResponse {
	resp := teov20220901.NewDescribeFunctionReplicasResponse()
	resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
		TotalCount: ptrInt64TeoFunctionReplica110901(1),
		FunctionReplicas: []*teov20220901.FunctionReplica{
			{
				FunctionId:  ptrStringTeoFunctionReplica110901("ef-test456"),
				ReplicaName: ptrStringTeoFunctionReplica110901(replicaName),
				Content:     ptrStringTeoFunctionReplica110901(content),
				Remark:      ptrStringTeoFunctionReplica110901(remark),
				CreatedOn:   ptrStringTeoFunctionReplica110901("2024-01-01T00:00:00Z"),
				ModifiedOn:  ptrStringTeoFunctionReplica110901("2024-01-02T00:00:00Z"),
			},
		},
		RequestId: ptrStringTeoFunctionReplica110901("fake-request-id"),
	}
	return resp
}

// mockTeoFunctionReplica110901Client wires a mocked teo client into a fake
// connectivity client.
func mockTeoFunctionReplica110901Client(patches *gomonkey.Patches) (*connectivity.TencentCloudClient, *teov20220901.Client) {
	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)
	return client, teoClient
}

// readTeoFunctionReplica110901State reads the model back from resp.State.Raw.
func readTeoFunctionReplica110901State(t *testing.T, raw tftypes.Value) (id, zoneId, functionId, replicaName, content, remark string, replicas []map[string]string) {
	val, ok := raw.Value.(map[string]tftypes.Value)
	if !ok {
		t.Fatalf("state raw value is not an object: %v", raw)
	}
	readStr := func(name string) string {
		v, ok := val[name]
		if !ok {
			t.Fatalf("state has no attribute %q", name)
		}
		s, ok := v.Value.(string)
		if !ok {
			return ""
		}
		return s
	}
	id = readStr("id")
	zoneId = readStr("zone_id")
	functionId = readStr("function_id")
	replicaName = readStr("replica_name")
	content = readStr("content")
	remark = readStr("remark")

	listVal, ok := val["function_replicas"]
	if !ok {
		t.Fatalf("state has no attribute function_replicas")
	}
	items, ok := listVal.Value.([]tftypes.Value)
	if !ok {
		if listVal.Value == nil {
			return
		}
		t.Fatalf("function_replicas is not a list: %v", listVal)
	}
	for _, item := range items {
		m, ok := item.Value.(map[string]tftypes.Value)
		if !ok {
			t.Fatalf("function_replicas element is not an object: %v", item)
		}
		entry := make(map[string]string)
		for k, v := range m {
			if s, ok := v.Value.(string); ok {
				entry[k] = s
			} else {
				entry[k] = ""
			}
		}
		replicas = append(replicas, entry)
	}
	return
}

// TestTeoFunctionReplica110901Resource_Metadata validates the type name.
func TestTeoFunctionReplica110901Resource_Metadata(t *testing.T) {
	r := providerTeo.NewTeoFunctionReplica110901Resource()
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)
	if resp.TypeName != "tencentcloud_teo_function_replica_1109_01" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_function_replica_1109_01")
	}
}

// TestTeoFunctionReplica110901Resource_Schema validates required / optional /
// computed flags on every attribute.
func TestTeoFunctionReplica110901Resource_Schema(t *testing.T) {
	r := providerTeo.NewTeoFunctionReplica110901Resource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	stringAttrs := map[string][3]bool{
		// name -> {required, optional, computed}
		"id":           {false, false, true},
		"zone_id":      {true, false, false},
		"function_id":  {true, false, false},
		"replica_name": {true, false, false},
		"content":      {true, false, false},
		"remark":       {false, true, false},
	}
	for name, want := range stringAttrs {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		if attr.IsRequired() != want[0] {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.IsRequired(), want[0])
		}
		if attr.IsOptional() != want[1] {
			t.Errorf("attribute %q Optional=%v, want %v", name, attr.IsOptional(), want[1])
		}
		if attr.IsComputed() != want[2] {
			t.Errorf("attribute %q Computed=%v, want %v", name, attr.IsComputed(), want[2])
		}
	}

	listAttr, ok := resp.Schema.Attributes["function_replicas"]
	if !ok {
		t.Fatal("attribute function_replicas missing from schema")
	}
	if !listAttr.IsComputed() {
		t.Error("attribute function_replicas should be Computed")
	}
	if listAttr.IsRequired() || listAttr.IsOptional() {
		t.Error("attribute function_replicas should not be Required or Optional")
	}
}

// TestTeoFunctionReplica110901Resource_Factory validates the registry factory.
func TestTeoFunctionReplica110901Resource_Factory(t *testing.T) {
	got := providerTeo.NewTeoFunctionReplica110901Resource()
	if got == nil {
		t.Fatal("NewTeoFunctionReplica110901Resource returned nil")
	}
	if _, ok := got.(*providerTeo.TeoFunctionReplica110901Resource); !ok {
		t.Fatalf("factory returned %T, want *TeoFunctionReplica110901Resource", got)
	}
}

// TestTeoFunctionReplica110901SplitId validates the composite id helpers.
func TestTeoFunctionReplica110901SplitId(t *testing.T) {
	good := strings.Join([]string{"zone-1", "ef-1", "replica-1"}, tccommon.FILED_SP)
	z, f, r, ok := providerTeo.SplitTeoFunctionReplica110901Id(good)
	if !ok || z != "zone-1" || f != "ef-1" || r != "replica-1" {
		t.Errorf("split id = (%q, %q, %q, %v), want (zone-1, ef-1, replica-1, true)", z, f, r, ok)
	}

	bad := []string{
		"",
		"only-one",
		"a" + tccommon.FILED_SP + "b",
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP,
		tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c" + tccommon.FILED_SP + "d",
	}
	for _, in := range bad {
		if _, _, _, ok := providerTeo.SplitTeoFunctionReplica110901Id(in); ok {
			t.Errorf("SplitTeoFunctionReplica110901Id(%q) ok=true, want false", in)
		}
	}
}

// TestTeoFunctionReplica110901Resource_CreateSuccess covers a successful
// create: the id is persisted before the read-back and all leaf fields are
// refreshed from the Describe response.
func TestTeoFunctionReplica110901Resource_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	var gotCreate *teov20220901.CreateFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		gotCreate = request
		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	var gotDescribe *teov20220901.DescribeFunctionReplicasRequest
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		gotDescribe = request
		return newTeoFunctionReplica110901DescribeResponse("test-replica", "console.log(123)", "test remark"), nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.CreateRequest{}
	req.Plan.Raw = teoFunctionReplica110901RawValue(
		nil,
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		ptrStringTeoFunctionReplica110901("test remark"),
	)
	req.Plan.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.CreateResponse{}
	resp.State.Schema = req.Plan.Schema
	r.Create(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// Verify the create request fields.
	assert.NotNil(t, gotCreate)
	assert.Equal(t, "zone-test123", *gotCreate.ZoneId)
	assert.Equal(t, "ef-test456", *gotCreate.FunctionId)
	assert.Equal(t, "test-replica", *gotCreate.ReplicaName)
	assert.Equal(t, "console.log(123)", *gotCreate.Content)
	assert.Equal(t, "test remark", *gotCreate.Remark)

	// Verify the read-back query uses the replica-name filter.
	assert.NotNil(t, gotDescribe)
	assert.Equal(t, "zone-test123", *gotDescribe.ZoneId)
	assert.Equal(t, "ef-test456", *gotDescribe.FunctionId)
	assert.Len(t, gotDescribe.Filters, 1)
	assert.Equal(t, "replica-name", *gotDescribe.Filters[0].Name)
	assert.Equal(t, []*string{ptrStringTeoFunctionReplica110901("test-replica")}, gotDescribe.Filters[0].Values)

	// Verify every leaf field in the resulting state.
	id, zoneId, functionId, replicaName, content, remark, replicas := readTeoFunctionReplica110901State(t, resp.State.Raw)
	assert.Equal(t, "zone-test123#ef-test456#test-replica", id)
	assert.Equal(t, "zone-test123", zoneId)
	assert.Equal(t, "ef-test456", functionId)
	assert.Equal(t, "test-replica", replicaName)
	assert.Equal(t, "console.log(123)", content)
	assert.Equal(t, "test remark", remark)
	if assert.Len(t, replicas, 1) {
		item := replicas[0]
		assert.Equal(t, "ef-test456", item["function_id"])
		assert.Equal(t, "test-replica", item["replica_name"])
		assert.Equal(t, "console.log(123)", item["content"])
		assert.Equal(t, "test remark", item["remark"])
		assert.Equal(t, "2024-01-01T00:00:00Z", item["created_on"])
		assert.Equal(t, "2024-01-02T00:00:00Z", item["modified_on"])
	}
}

// TestTeoFunctionReplica110901Resource_CreateApiError verifies that an API
// error on create surfaces as an error diagnostic and no state is written.
func TestTeoFunctionReplica110901Resource_CreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		return nil, sdkErrors.NewTencentCloudSDKError("InvalidParameter", "invalid replica name", "fake-request-id")
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.CreateRequest{}
	req.Plan.Raw = teoFunctionReplica110901RawValue(
		nil,
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.Plan.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.CreateResponse{}
	resp.State.Schema = req.Plan.Schema
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
}

// TestTeoFunctionReplica110901Resource_CreateReadBackNotFound verifies that
// a created replica which cannot be read back fails with an error while the id
// stays in state.
func TestTeoFunctionReplica110901Resource_CreateReadBackNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64TeoFunctionReplica110901(0),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.CreateRequest{}
	req.Plan.Raw = teoFunctionReplica110901RawValue(
		nil,
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.Plan.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.CreateResponse{}
	resp.State.Schema = req.Plan.Schema
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics when the replica cannot be read back")
}

// TestTeoFunctionReplica110901Resource_ReadSuccess covers a successful
// refresh of every leaf field.
func TestTeoFunctionReplica110901Resource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		return newTeoFunctionReplica110901DescribeResponse("test-replica", "console.log(123)", "test remark"), nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.ReadRequest{}
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(old)"),
		ptrStringTeoFunctionReplica110901("old remark"),
	)
	req.State.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.ReadResponse{}
	resp.State.Schema = req.State.Schema
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	id, zoneId, functionId, replicaName, content, remark, replicas := readTeoFunctionReplica110901State(t, resp.State.Raw)
	assert.Equal(t, "zone-test123#ef-test456#test-replica", id)
	assert.Equal(t, "zone-test123", zoneId)
	assert.Equal(t, "ef-test456", functionId)
	assert.Equal(t, "test-replica", replicaName)
	assert.Equal(t, "console.log(123)", content)
	assert.Equal(t, "test remark", remark)
	if assert.Len(t, replicas, 1) {
		item := replicas[0]
		assert.Equal(t, "ef-test456", item["function_id"])
		assert.Equal(t, "test-replica", item["replica_name"])
		assert.Equal(t, "console.log(123)", item["content"])
		assert.Equal(t, "test remark", item["remark"])
		assert.Equal(t, "2024-01-01T00:00:00Z", item["created_on"])
		assert.Equal(t, "2024-01-02T00:00:00Z", item["modified_on"])
	}
}

// TestTeoFunctionReplica110901Resource_ReadNotFound verifies that an empty
// (but successful) query result removes the resource from state.
func TestTeoFunctionReplica110901Resource_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64TeoFunctionReplica110901(0),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.ReadRequest{}
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.State.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.ReadResponse{}
	resp.State.Schema = req.State.Schema
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.Nil(t, resp.State.Raw.Value, "expected the resource to be removed from state")
}

// TestTeoFunctionReplica110901Resource_ReadApiError verifies that an API
// error during Read surfaces as an error diagnostic and keeps the state.
func TestTeoFunctionReplica110901Resource_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return nil, sdkErrors.NewTencentCloudSDKError("InternalError", "backend error", "fake-request-id")
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.ReadRequest{}
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.State.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.ReadResponse{}
	resp.State.Schema = req.State.Schema
	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
	assert.NotNil(t, resp.State.Raw.Value, "state must not be cleared on an API failure")
}

// TestTeoFunctionReplica110901Resource_ReadBrokenId verifies that a broken
// composite id fails loudly instead of silently proceeding.
func TestTeoFunctionReplica110901Resource_ReadBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, _ := mockTeoFunctionReplica110901Client(patches)

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.ReadRequest{}
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("broken-id"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.State.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.ReadResponse{}
	resp.State.Schema = req.State.Schema
	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for a broken id")
}

// TestTeoFunctionReplica110901Resource_UpdateSuccess covers an in-place
// update of content and remark followed by a full read-back.
func TestTeoFunctionReplica110901Resource_UpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	var gotModify *teov20220901.ModifyFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		gotModify = request
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return newTeoFunctionReplica110901DescribeResponse("test-replica", "console.log(456)", "updated remark"), nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.UpdateRequest{}
	req.Plan.Raw = teoFunctionReplica110901RawValue(
		nil,
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(456)"),
		ptrStringTeoFunctionReplica110901("updated remark"),
	)
	req.Plan.Schema = schemaOfTeoFunctionReplica110901(t, r)
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		ptrStringTeoFunctionReplica110901("test remark"),
	)
	req.State.Schema = req.Plan.Schema
	resp := &resource.UpdateResponse{}
	resp.State.Schema = req.Plan.Schema
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// Verify the modify request fields.
	assert.NotNil(t, gotModify)
	assert.Equal(t, "zone-test123", *gotModify.ZoneId)
	assert.Equal(t, "ef-test456", *gotModify.FunctionId)
	assert.Equal(t, "test-replica", *gotModify.ReplicaName)
	assert.Equal(t, "console.log(456)", *gotModify.Content)
	assert.Equal(t, "updated remark", *gotModify.Remark)

	// Verify every refreshed leaf field.
	id, zoneId, functionId, replicaName, content, remark, replicas := readTeoFunctionReplica110901State(t, resp.State.Raw)
	assert.Equal(t, "zone-test123#ef-test456#test-replica", id)
	assert.Equal(t, "zone-test123", zoneId)
	assert.Equal(t, "ef-test456", functionId)
	assert.Equal(t, "test-replica", replicaName)
	assert.Equal(t, "console.log(456)", content)
	assert.Equal(t, "updated remark", remark)
	if assert.Len(t, replicas, 1) {
		item := replicas[0]
		assert.Equal(t, "ef-test456", item["function_id"])
		assert.Equal(t, "test-replica", item["replica_name"])
		assert.Equal(t, "console.log(456)", item["content"])
		assert.Equal(t, "updated remark", item["remark"])
		assert.Equal(t, "2024-01-01T00:00:00Z", item["created_on"])
		assert.Equal(t, "2024-01-02T00:00:00Z", item["modified_on"])
	}
}

// TestTeoFunctionReplica110901Resource_UpdateExplicitEmptyRemark covers
// updating with remark set to an explicit empty string (the documented way to
// ask the cloud API to clear the description).
func TestTeoFunctionReplica110901Resource_UpdateExplicitEmptyRemark(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	var gotModify *teov20220901.ModifyFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		gotModify = request
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return newTeoFunctionReplica110901DescribeResponse("test-replica", "console.log(456)", ""), nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.UpdateRequest{}
	req.Plan.Raw = teoFunctionReplica110901RawValue(
		nil,
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(456)"),
		ptrStringTeoFunctionReplica110901(""),
	)
	req.Plan.Schema = schemaOfTeoFunctionReplica110901(t, r)
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		ptrStringTeoFunctionReplica110901("test remark"),
	)
	req.State.Schema = req.Plan.Schema
	resp := &resource.UpdateResponse{}
	resp.State.Schema = req.Plan.Schema
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotModify)
	// An explicit empty string is sent as-is; the cloud API validates it.
	assert.NotNil(t, gotModify.Remark)
	assert.Equal(t, "", *gotModify.Remark)

	_, _, _, _, _, remark, _ := readTeoFunctionReplica110901State(t, resp.State.Raw)
	assert.Equal(t, "", remark)
}

// TestTeoFunctionReplica110901Resource_UpdateRemarkRemovedUnconfirmed covers
// removing the remark from the configuration while the state still carries a
// value: the clear semantics is undocumented, so the update must fail loudly
// instead of silently keeping or clearing the description.
func TestTeoFunctionReplica110901Resource_UpdateRemarkRemovedUnconfirmed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		modifyCalled = true
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.UpdateRequest{}
	req.Plan.Raw = teoFunctionReplica110901RawValue(
		nil,
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(456)"),
		nil,
	)
	req.Plan.Schema = schemaOfTeoFunctionReplica110901(t, r)
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		ptrStringTeoFunctionReplica110901("test remark"),
	)
	req.State.Schema = req.Plan.Schema
	resp := &resource.UpdateResponse{}
	resp.State.Schema = req.Plan.Schema
	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for the unconfirmed remark clear semantics")
	assert.False(t, modifyCalled, "ModifyFunctionReplica must not be called on the ambiguous transition")
}

// TestTeoFunctionReplica110901Resource_UpdateRemarkBothNull covers updating
// with remark null in both plan and state (never set): the modify request is
// sent with a nil Remark.
func TestTeoFunctionReplica110901Resource_UpdateRemarkBothNull(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	var gotModify *teov20220901.ModifyFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		gotModify = request
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return newTeoFunctionReplica110901DescribeResponse("test-replica", "console.log(456)", ""), nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.UpdateRequest{}
	req.Plan.Raw = teoFunctionReplica110901RawValue(
		nil,
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(456)"),
		nil,
	)
	req.Plan.Schema = schemaOfTeoFunctionReplica110901(t, r)
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.State.Schema = req.Plan.Schema
	resp := &resource.UpdateResponse{}
	resp.State.Schema = req.Plan.Schema
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotModify)
	assert.Nil(t, gotModify.Remark)
}

// TestTeoFunctionReplica110901Resource_UpdateApiError verifies that an API
// error on modify surfaces as an error diagnostic.
func TestTeoFunctionReplica110901Resource_UpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		return nil, sdkErrors.NewTencentCloudSDKError("InvalidParameter", "invalid content", "fake-request-id")
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.UpdateRequest{}
	req.Plan.Raw = teoFunctionReplica110901RawValue(
		nil,
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(456)"),
		nil,
	)
	req.Plan.Schema = schemaOfTeoFunctionReplica110901(t, r)
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.State.Schema = req.Plan.Schema
	resp := &resource.UpdateResponse{}
	resp.State.Schema = req.Plan.Schema
	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
}

// TestTeoFunctionReplica110901Resource_DeleteSuccess covers a successful
// delete.
func TestTeoFunctionReplica110901Resource_DeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	var gotDelete *teov20220901.DeleteFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		gotDelete = request
		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrStringTeoFunctionReplica110901("fake-request-id"),
		}
		return resp, nil
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.DeleteRequest{}
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.State.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.DeleteResponse{}
	resp.State.Schema = req.State.Schema
	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotDelete)
	assert.Equal(t, "zone-test123", *gotDelete.ZoneId)
	assert.Equal(t, "ef-test456", *gotDelete.FunctionId)
	// Single-resource delete sends a one-element list.
	assert.Equal(t, []*string{ptrStringTeoFunctionReplica110901("test-replica")}, gotDelete.ReplicaNames)
}

// TestTeoFunctionReplica110901Resource_DeleteApiError verifies that an API
// error on delete surfaces as an error diagnostic (state is kept for retry).
func TestTeoFunctionReplica110901Resource_DeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		return nil, sdkErrors.NewTencentCloudSDKError("InternalError", "backend error", "fake-request-id")
	})

	r := newTeoFunctionReplica110901Resource(t, client)

	req := resource.DeleteRequest{}
	req.State.Raw = teoFunctionReplica110901RawValue(
		ptrStringTeoFunctionReplica110901("zone-test123#ef-test456#test-replica"),
		ptrStringTeoFunctionReplica110901("zone-test123"),
		ptrStringTeoFunctionReplica110901("ef-test456"),
		ptrStringTeoFunctionReplica110901("test-replica"),
		ptrStringTeoFunctionReplica110901("console.log(123)"),
		nil,
	)
	req.State.Schema = schemaOfTeoFunctionReplica110901(t, r)
	resp := &resource.DeleteResponse{}
	resp.State.Schema = req.State.Schema
	r.Delete(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
}

// TestTeoFunctionReplica110901Resource_ImportStateGood covers the import path
// with a valid composite id.
func TestTeoFunctionReplica110901Resource_ImportStateGood(t *testing.T) {
	r := providerTeo.NewTeoFunctionReplica110901Resource()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", schemaResp.Diagnostics)
	}

	req := resource.ImportStateRequest{ID: "zone-test123#ef-test456#test-replica"}
	resp := &resource.ImportStateResponse{}
	resp.State.Schema = schemaResp.Schema
	r.ImportState(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	val, ok := resp.State.Raw.Value.(map[string]tftypes.Value)
	if !ok {
		t.Fatalf("imported state raw value is not an object: %v", resp.State.Raw)
	}
	checkStr := func(name, want string) {
		v, ok := val[name]
		if !ok {
			t.Fatalf("imported state has no attribute %q", name)
		}
		s, ok := v.Value.(string)
		if !ok || s != want {
			t.Errorf("imported state %q = %v, want %q", name, v.Value, want)
		}
	}
	checkStr("id", "zone-test123#ef-test456#test-replica")
	checkStr("zone_id", "zone-test123")
	checkStr("function_id", "ef-test456")
	checkStr("replica_name", "test-replica")
}

// TestTeoFunctionReplica110901Resource_ImportStateBadId covers the import
// path with a broken composite id.
func TestTeoFunctionReplica110901Resource_ImportStateBadId(t *testing.T) {
	r := providerTeo.NewTeoFunctionReplica110901Resource()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", schemaResp.Diagnostics)
	}

	for _, in := range []string{"", "only-one", "a#b", "a#b#c#d"} {
		req := resource.ImportStateRequest{ID: in}
		resp := &resource.ImportStateResponse{}
		resp.State.Schema = schemaResp.Schema
		r.ImportState(context.Background(), req, resp)
		assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for import id %q", in)
	}
}

// TestTeoFunctionReplica110901Service_DescribeNilResponse verifies that a nil
// describe response is treated as an API failure, not as "not found".
func TestTeoFunctionReplica110901Service_DescribeNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := mockTeoFunctionReplica110901Client(patches)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		// nil response: an API failure, must not be folded into "not found".
		return nil, nil
	})

	service := providerTeo.NewTeoService(client)
	item, err := service.DescribeTeoFunctionReplica(context.Background(), "zone-test123", "ef-test456", "test-replica")
	assert.Error(t, err, "expected an error for a nil describe response")
	assert.Nil(t, item, "expected no item for a nil describe response")
	assert.True(t, strings.Contains(err.Error(), "test-replica"), "error should mention the replica name, got: %v", err)
}

// schemaOfTeoFunctionReplica110901 extracts the framework schema of the
// resource, used to construct test requests.
func schemaOfTeoFunctionReplica110901(t *testing.T, r resource.Resource) (schema schema.Schema) {
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}
