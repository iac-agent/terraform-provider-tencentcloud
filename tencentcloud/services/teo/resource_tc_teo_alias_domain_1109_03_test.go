package teo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

func ptrStrTeoAliasDomain110903(s string) *string { return &s }

// teoAliasDomain110903ResourceSchema returns the resource schema for tests.
func teoAliasDomain110903ResourceSchema(t *testing.T) schema.Schema {
	r := &TeoAliasDomain110903Resource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// teoAliasDomain110903State builds a tfsdk.State carrying the model.
func teoAliasDomain110903State(t *testing.T, s schema.Schema, m *TeoAliasDomain110903ResourceModel) tfsdk.State {
	return tfsdk.State{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, m)}
}

// teoAliasDomain110903Plan builds a tfsdk.Plan carrying the model.
func teoAliasDomain110903Plan(t *testing.T, s schema.Schema, m *TeoAliasDomain110903ResourceModel) tfsdk.Plan {
	return tfsdk.Plan{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, m)}
}

// teoAliasDomain110903RawValue converts the model into a raw tftypes object value.
func teoAliasDomain110903RawValue(t *testing.T, s schema.Schema, m *TeoAliasDomain110903ResourceModel) tftypes.Value {
	objType := s.Type().TerraformType(context.Background())
	attrs, ok := objType.(tftypes.Object)
	if !ok {
		t.Fatalf("resource type is not an object: %T", objType)
	}
	vals := make(map[string]tftypes.Value, len(attrs.AttributeTypes))
	for name := range attrs.AttributeTypes {
		var v tftypes.Value
		switch name {
		case "id":
			v = tftypes.NewValue(tftypes.String, teoAliasDomain110903StringOrNull(m.Id))
		case "zone_id":
			v = tftypes.NewValue(tftypes.String, teoAliasDomain110903StringOrNull(m.ZoneId))
		case "alias_name":
			v = tftypes.NewValue(tftypes.String, teoAliasDomain110903StringOrNull(m.AliasName))
		case "target_name":
			v = tftypes.NewValue(tftypes.String, teoAliasDomain110903StringOrNull(m.TargetName))
		case "cert_type":
			v = tftypes.NewValue(tftypes.String, teoAliasDomain110903StringOrNull(m.CertType))
		case "status":
			v = tftypes.NewValue(tftypes.String, teoAliasDomain110903StringOrNull(m.Status))
		case "created_on":
			v = tftypes.NewValue(tftypes.String, teoAliasDomain110903StringOrNull(m.CreatedOn))
		case "modified_on":
			v = tftypes.NewValue(tftypes.String, teoAliasDomain110903StringOrNull(m.ModifiedOn))
		case "forbid_mode":
			if m.ForbidMode.IsNull() {
				v = tftypes.NewValue(tftypes.Number, nil)
			} else {
				v = tftypes.NewValue(tftypes.Number, m.ForbidMode.ValueInt64())
			}
		case "cert_id":
			if m.CertId.IsNull() {
				v = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
			} else {
				elems := make([]tftypes.Value, 0, len(m.CertId.Elements()))
				for _, e := range m.CertId.Elements() {
					s, ok := e.(types.String)
					if !ok || s.IsNull() || s.IsUnknown() {
						continue
					}
					elems = append(elems, tftypes.NewValue(tftypes.String, s.ValueString()))
				}
				v = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, elems)
			}
		default:
			t.Fatalf("unexpected attribute %q in raw value builder", name)
		}
		vals[name] = v
	}
	return tftypes.NewValue(attrs, vals)
}

func teoAliasDomain110903StringOrNull(v types.String) interface{} {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return v.ValueString()
}

// setTeoAliasDomain110903Client injects a mock connectivity client via Configure.
func setTeoAliasDomain110903Client(r *TeoAliasDomain110903Resource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
}

// mockTeoAliasDomain110903Describe installs a gomonkey patch that replaces
// DescribeTeoAliasDomainById on the TeoService value received by the resource.
func mockTeoAliasDomain110903Describe(patches *gomonkey.Patches, item *teov20220901.AliasDomain, err error) {
	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		return item, err
	})
}

// mockTeoAliasDomain110903Create installs a gomonkey patch that replaces
// CreateTeoAliasDomain on the TeoService value received by the resource.
func mockTeoAliasDomain110903Create(patches *gomonkey.Patches, request **teov20220901.CreateAliasDomainRequest, err error) {
	patches.ApplyMethodFunc(&TeoService{}, "CreateTeoAliasDomain", func(_ context.Context, req *teov20220901.CreateAliasDomainRequest) error {
		*request = req
		return err
	})
}

// mockTeoAliasDomain110903Modify installs a gomonkey patch that replaces
// ModifyTeoAliasDomain on the TeoService value received by the resource.
func mockTeoAliasDomain110903Modify(patches *gomonkey.Patches, request **teov20220901.ModifyAliasDomainRequest, err error) {
	patches.ApplyMethodFunc(&TeoService{}, "ModifyTeoAliasDomain", func(_ context.Context, req *teov20220901.ModifyAliasDomainRequest) error {
		*request = req
		return err
	})
}

// go test ./tencentcloud/services/teo/ -run "TestTeoAliasDomain110903" -v -count=1 -gcflags="all=-l"

// TestTeoAliasDomain110903Resource_Metadata validates the resource type name.
func TestTeoAliasDomain110903Resource_Metadata(t *testing.T) {
	r := &TeoAliasDomain110903Resource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_teo_alias_domain_1109_03" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_alias_domain_1109_03")
	}
}

// TestTeoAliasDomain110903Resource_Schema validates the required / optional /
// computed flags on every attribute.
func TestTeoAliasDomain110903Resource_Schema(t *testing.T) {
	r := &TeoAliasDomain110903Resource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	stringFlags := map[string][3]bool{
		"id":          {false, false, true},
		"zone_id":     {true, false, false},
		"alias_name":  {true, false, false},
		"target_name": {true, false, false},
		"cert_type":   {false, true, false},
		"status":      {false, false, true},
		"created_on":  {false, false, true},
		"modified_on": {false, false, true},
	}
	for name, want := range stringFlags {
		attr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q missing or not a StringAttribute", name)
			continue
		}
		if attr.Required != want[0] {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.Required, want[0])
		}
		if attr.Optional != want[1] {
			t.Errorf("attribute %q Optional=%v, want %v", name, attr.Optional, want[1])
		}
		if attr.Computed != want[2] {
			t.Errorf("attribute %q Computed=%v, want %v", name, attr.Computed, want[2])
		}
	}

	certId, ok := resp.Schema.Attributes["cert_id"].(schema.SetAttribute)
	if !ok {
		t.Fatalf("attribute %q missing or not a SetAttribute", "cert_id")
	}
	if !certId.Optional || certId.Required || certId.Computed {
		t.Errorf("attribute cert_id flags Optional=%v Required=%v Computed=%v, want optional only", certId.Optional, certId.Required, certId.Computed)
	}
	forbidMode, ok := resp.Schema.Attributes["forbid_mode"].(schema.Int64Attribute)
	if !ok {
		t.Fatalf("attribute %q missing or not an Int64Attribute", "forbid_mode")
	}
	if !forbidMode.Computed || forbidMode.Optional || forbidMode.Required {
		t.Errorf("attribute forbid_mode flags Optional=%v Required=%v Computed=%v, want computed only", forbidMode.Optional, forbidMode.Required, forbidMode.Computed)
	}
}

// TestTeoAliasDomain110903Resource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces it advertises in package-level asserts.
func TestTeoAliasDomain110903Resource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &TeoAliasDomain110903Resource{}
	var _ resource.ResourceWithImportState = &TeoAliasDomain110903Resource{}
}

// TestTeoAliasDomain110903Resource_Factory validates the registry factory.
func TestTeoAliasDomain110903Resource_Factory(t *testing.T) {
	got := NewTeoAliasDomain110903Resource()
	if got == nil {
		t.Fatal("NewTeoAliasDomain110903Resource returned nil")
	}
	if _, ok := got.(*TeoAliasDomain110903Resource); !ok {
		t.Fatalf("factory returned %T, want *TeoAliasDomain110903Resource", got)
	}
}

// TestTeoAliasDomain110903Resource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestTeoAliasDomain110903Resource_ConfigureNilProviderData(t *testing.T) {
	r := &TeoAliasDomain110903Resource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestSplitTeoAliasDomain110903Id validates parsing of `<zone_id>#<alias_name>`.
func TestSplitTeoAliasDomain110903Id(t *testing.T) {
	good := "zone-2qtuhspy7cr6" + tccommon.FILED_SP + "alias.demo.com"
	zone, alias, ok := splitTeoAliasDomain110903Id(good)
	if !ok || zone != "zone-2qtuhspy7cr6" || alias != "alias.demo.com" {
		t.Errorf("splitTeoAliasDomain110903Id(%q) = (%q, %q, %v), want (zone, alias, true)", good, zone, alias, ok)
	}

	bad := []string{
		"",
		"only-zone",
		tccommon.FILED_SP + "alias.demo.com",
		"zone-2qtuhspy7cr6" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	}
	for _, in := range bad {
		if _, _, ok := splitTeoAliasDomain110903Id(in); ok {
			t.Errorf("splitTeoAliasDomain110903Id(%q) ok=true, want false", in)
		}
	}
}

// TestSetToStringPtrs covers null / unknown / empty / value sets.
func TestSetToStringPtrs(t *testing.T) {
	if got := setToStringPtrs(types.SetNull(types.StringType)); got != nil {
		t.Errorf("setToStringPtrs(null) = %v, want nil", got)
	}

	empty, diags := types.SetValueFrom(context.Background(), types.StringType, []string{})
	if diags.HasError() {
		t.Fatalf("empty set construction failed: %v", diags)
	}
	if got := setToStringPtrs(empty); len(got) != 0 {
		t.Errorf("setToStringPtrs(empty) = %v, want empty", got)
	}

	set, diags := types.SetValueFrom(context.Background(), types.StringType, []string{"cert-1", "cert-2"})
	if diags.HasError() {
		t.Fatalf("set construction failed: %v", diags)
	}
	got := setToStringPtrs(set)
	if len(got) != 2 || *got[0] != "cert-1" || *got[1] != "cert-2" {
		t.Errorf("setToStringPtrs(set) = %v, want [cert-1 cert-2]", got)
	}
}

// TestTeoAliasDomain110903Resource_Create_Success covers a successful create:
// request conversion, id persistence, status wait and read-back.
func TestTeoAliasDomain110903Resource_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	var createReq *teov20220901.CreateAliasDomainRequest
	mockTeoAliasDomain110903Create(patches, &createReq, nil)

	describeCalls := 0
	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		describeCalls++
		return &teov20220901.AliasDomain{
			AliasName:  ptrStrTeoAliasDomain110903("alias.demo.com"),
			ZoneId:     ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6"),
			TargetName: ptrStrTeoAliasDomain110903("target.demo.com"),
			Status:     ptrStrTeoAliasDomain110903("active"),
			ForbidMode: func() *int64 { v := int64(0); return &v }(),
			CreatedOn:  ptrStrTeoAliasDomain110903("2024-05-20 10:00:00"),
			ModifiedOn: ptrStrTeoAliasDomain110903("2024-05-20 10:00:00"),
		}, nil
	})

	certId, diags := types.SetValueFrom(context.Background(), types.StringType, []string{"cert-1"})
	if diags.HasError() {
		t.Fatalf("cert_id construction failed: %v", diags)
	}

	planModel := &TeoAliasDomain110903ResourceModel{
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
		CertType:   types.StringValue("hosting"),
		CertId:     certId,
	}

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	req := resource.CreateRequest{Plan: teoAliasDomain110903Plan(t, s, planModel)}
	resp := &resource.CreateResponse{State: teoAliasDomain110903State(t, s, &TeoAliasDomain110903ResourceModel{})}
	r.Create(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.True(t, describeCalls >= 1, "expected the resource to be read back after create")
	assert.NotNil(t, createReq)
	assert.Equal(t, "zone-2qtuhspy7cr6", *createReq.ZoneId)
	assert.Equal(t, "alias.demo.com", *createReq.AliasName)
	assert.Equal(t, "target.demo.com", *createReq.TargetName)
	assert.Equal(t, "hosting", *createReq.CertType)
	assert.Equal(t, []*string{ptrStrTeoAliasDomain110903("cert-1")}, createReq.CertId)

	var got TeoAliasDomain110903ResourceModel
	diags = resp.State.Get(context.Background(), &got)
	if diags.HasError() {
		t.Fatalf("reading back the state failed: %v", diags)
	}
	assert.Equal(t, "zone-2qtuhspy7cr6#alias.demo.com", got.Id.ValueString())
	assert.Equal(t, "target.demo.com", got.TargetName.ValueString())
	assert.Equal(t, "active", got.Status.ValueString())
	assert.Equal(t, int64(0), got.ForbidMode.ValueInt64())
	assert.Equal(t, "2024-05-20 10:00:00", got.CreatedOn.ValueString())
	assert.Equal(t, "2024-05-20 10:00:00", got.ModifiedOn.ValueString())
	// cert_type / cert_id are not returned by DescribeAliasDomains and keep the planned values.
	assert.Equal(t, "hosting", got.CertType.ValueString())
	assert.Equal(t, 1, len(got.CertId.Elements()))
}

// TestTeoAliasDomain110903Resource_Create_ApiError covers the API error branch.
func TestTeoAliasDomain110903Resource_Create_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	var createReq *teov20220901.CreateAliasDomainRequest
	mockTeoAliasDomain110903Create(patches, &createReq, fmt.Errorf("InternalError"))

	planModel := &TeoAliasDomain110903ResourceModel{
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
	}

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	req := resource.CreateRequest{Plan: teoAliasDomain110903Plan(t, s, planModel)}
	resp := &resource.CreateResponse{State: teoAliasDomain110903State(t, s, &TeoAliasDomain110903ResourceModel{})}
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics from a failed CreateAliasDomain")
}

// TestTeoAliasDomain110903Resource_Create_WaitTimeout covers the branch where
// the alias domain stays in pending state past the wait budget.
func TestTeoAliasDomain110903Resource_Create_WaitTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	var createReq *teov20220901.CreateAliasDomainRequest
	mockTeoAliasDomain110903Create(patches, &createReq, nil)

	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		return &teov20220901.AliasDomain{
			AliasName:  ptrStrTeoAliasDomain110903("alias.demo.com"),
			ZoneId:     ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6"),
			TargetName: ptrStrTeoAliasDomain110903("target.demo.com"),
			Status:     ptrStrTeoAliasDomain110903("pending"),
		}, nil
	})
	// Shrink the wait budget and the poll interval so the test terminates fast.
	// Save the original values and restore them afterwards.
	origTimeout := tccommon.WriteRetryTimeout
	origInterval := teoAliasDomain110903PollInterval
	tccommon.WriteRetryTimeout = 20 * time.Millisecond
	teoAliasDomain110903PollInterval = 5 * time.Millisecond
	defer func() {
		tccommon.WriteRetryTimeout = origTimeout
		teoAliasDomain110903PollInterval = origInterval
	}()

	planModel := &TeoAliasDomain110903ResourceModel{
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
	}

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	req := resource.CreateRequest{Plan: teoAliasDomain110903Plan(t, s, planModel)}
	resp := &resource.CreateResponse{State: teoAliasDomain110903State(t, s, &TeoAliasDomain110903ResourceModel{})}
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics when the alias domain never settles")
	var got TeoAliasDomain110903ResourceModel
	diags := resp.State.Get(context.Background(), &got)
	if diags.HasError() {
		t.Fatalf("reading back the state failed: %v", diags)
	}
	// The identity must be persisted before the wait so destroy stays possible.
	assert.Equal(t, "zone-2qtuhspy7cr6#alias.demo.com", got.Id.ValueString())
}

// TestTeoAliasDomain110903Resource_Read_NotFound covers the confirmed
// not-found branch: the resource is removed from state without any error.
func TestTeoAliasDomain110903Resource_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mockTeoAliasDomain110903Describe(patches, nil, nil)

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
	}

	req := resource.ReadRequest{State: teoAliasDomain110903State(t, s, stateModel)}
	resp := &resource.ReadResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.True(t, resp.State.Raw.IsNull(), "expected the resource to be removed from state")
}

// TestTeoAliasDomain110903Resource_Read_ApiError covers the API error branch:
// an API failure must not remove the resource from state.
func TestTeoAliasDomain110903Resource_Read_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mockTeoAliasDomain110903Describe(patches, nil, fmt.Errorf("UnauthorizedOperation"))

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
	}

	req := resource.ReadRequest{State: teoAliasDomain110903State(t, s, stateModel)}
	resp := &resource.ReadResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics from a failed DescribeAliasDomains")
	assert.False(t, resp.State.Raw.IsNull(), "state must be kept when the read fails")
}

// TestTeoAliasDomain110903Resource_Read_Success covers the read-back of every
// leaf field, including the keep-semantics of cert_type / cert_id.
func TestTeoAliasDomain110903Resource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mockTeoAliasDomain110903Describe(patches, &teov20220901.AliasDomain{
		AliasName:  ptrStrTeoAliasDomain110903("alias.demo.com"),
		ZoneId:     ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6"),
		TargetName: ptrStrTeoAliasDomain110903("target.demo.com"),
		Status:     ptrStrTeoAliasDomain110903("active"),
		ForbidMode: func() *int64 { v := int64(11); return &v }(),
		CreatedOn:  ptrStrTeoAliasDomain110903("2024-05-20 10:00:00"),
		ModifiedOn: ptrStrTeoAliasDomain110903("2024-05-21 11:00:00"),
	}, nil)

	certId, diags := types.SetValueFrom(context.Background(), types.StringType, []string{"cert-9"})
	if diags.HasError() {
		t.Fatalf("cert_id construction failed: %v", diags)
	}

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
		CertType:   types.StringValue("hosting"),
		CertId:     certId,
	}

	req := resource.ReadRequest{State: teoAliasDomain110903State(t, s, stateModel)}
	resp := &resource.ReadResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var got TeoAliasDomain110903ResourceModel
	diags = resp.State.Get(context.Background(), &got)
	if diags.HasError() {
		t.Fatalf("reading back the state failed: %v", diags)
	}
	assert.Equal(t, "zone-2qtuhspy7cr6", got.ZoneId.ValueString())
	assert.Equal(t, "alias.demo.com", got.AliasName.ValueString())
	assert.Equal(t, "target.demo.com", got.TargetName.ValueString())
	assert.Equal(t, "active", got.Status.ValueString())
	assert.Equal(t, int64(11), got.ForbidMode.ValueInt64())
	assert.Equal(t, "2024-05-20 10:00:00", got.CreatedOn.ValueString())
	assert.Equal(t, "2024-05-21 11:00:00", got.ModifiedOn.ValueString())
	// cert_type / cert_id are not returned by DescribeAliasDomains and keep the state values.
	assert.Equal(t, "hosting", got.CertType.ValueString())
	assert.Equal(t, 1, len(got.CertId.Elements()))
	elem, ok := got.CertId.Elements()[0].(types.String)
	if !ok {
		t.Fatalf("cert_id element is not a types.String: %T", got.CertId.Elements()[0])
	}
	assert.Equal(t, "cert-9", elem.ValueString())
}

// TestTeoAliasDomain110903Resource_Read_BrokenId covers the invalid-id branch.
func TestTeoAliasDomain110903Resource_Read_BrokenId(t *testing.T) {
	r := &TeoAliasDomain110903Resource{}

	s := teoAliasDomain110903ResourceSchema(t)
	stateModel := &TeoAliasDomain110903ResourceModel{Id: types.StringValue("broken-id")}

	req := resource.ReadRequest{State: teoAliasDomain110903State(t, s, stateModel)}
	resp := &resource.ReadResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for a broken id")
}

// TestTeoAliasDomain110903Resource_Update_Success covers a successful update:
// TargetName is always carried; a set CertType is sent; a null CertType keeps
// the cloud-side configuration (nothing sent).
func TestTeoAliasDomain110903Resource_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	var modifyReq *teov20220901.ModifyAliasDomainRequest
	mockTeoAliasDomain110903Modify(patches, &modifyReq, nil)

	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		return &teov20220901.AliasDomain{
			AliasName:  ptrStrTeoAliasDomain110903("alias.demo.com"),
			ZoneId:     ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6"),
			TargetName: ptrStrTeoAliasDomain110903("new-target.demo.com"),
			Status:     ptrStrTeoAliasDomain110903("active"),
		}, nil
	})

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	certId, diags := types.SetValueFrom(context.Background(), types.StringType, []string{"cert-1"})
	if diags.HasError() {
		t.Fatalf("cert_id construction failed: %v", diags)
	}

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
		CertType:   types.StringValue("hosting"),
		CertId:     certId,
	}
	planModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("new-target.demo.com"),
		CertType:   types.StringValue("hosting"),
		CertId:     certId,
	}

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, planModel)},
		Plan:   teoAliasDomain110903Plan(t, s, planModel),
		State:  teoAliasDomain110903State(t, s, stateModel),
	}
	resp := &resource.UpdateResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, modifyReq)
	assert.Equal(t, "zone-2qtuhspy7cr6", *modifyReq.ZoneId)
	assert.Equal(t, "alias.demo.com", *modifyReq.AliasName)
	assert.Equal(t, "new-target.demo.com", *modifyReq.TargetName)
	assert.Equal(t, "hosting", *modifyReq.CertType)
	assert.Equal(t, []*string{ptrStrTeoAliasDomain110903("cert-1")}, modifyReq.CertId)

	var got TeoAliasDomain110903ResourceModel
	diags = resp.State.Get(context.Background(), &got)
	if diags.HasError() {
		t.Fatalf("reading back the state failed: %v", diags)
	}
	assert.Equal(t, "zone-2qtuhspy7cr6#alias.demo.com", got.Id.ValueString())
	assert.Equal(t, "new-target.demo.com", got.TargetName.ValueString())
}

// TestTeoAliasDomain110903Resource_Update_NullCertKeepsCloudConfig covers the
// confirmed keep-semantics: removing cert_type / cert_id from the configuration
// sends nothing for them, and the state keeps the previous values.
func TestTeoAliasDomain110903Resource_Update_NullCertKeepsCloudConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	var modifyReq *teov20220901.ModifyAliasDomainRequest
	mockTeoAliasDomain110903Modify(patches, &modifyReq, nil)

	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		return &teov20220901.AliasDomain{
			AliasName:  ptrStrTeoAliasDomain110903("alias.demo.com"),
			ZoneId:     ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6"),
			TargetName: ptrStrTeoAliasDomain110903("target.demo.com"),
			Status:     ptrStrTeoAliasDomain110903("active"),
		}, nil
	})

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateCertId, diags := types.SetValueFrom(context.Background(), types.StringType, []string{"cert-1"})
	if diags.HasError() {
		t.Fatalf("cert_id construction failed: %v", diags)
	}

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
		CertType:   types.StringValue("hosting"),
		CertId:     stateCertId,
	}
	// The plan removed cert_type / cert_id from the configuration.
	planModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
		CertType:   types.StringNull(),
		CertId:     types.SetNull(types.StringType),
	}

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, planModel)},
		Plan:   teoAliasDomain110903Plan(t, s, planModel),
		State:  teoAliasDomain110903State(t, s, stateModel),
	}
	resp := &resource.UpdateResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, modifyReq)
	// Keep semantics: nothing is sent for the removed attributes.
	assert.Nil(t, modifyReq.CertType)
	assert.Nil(t, modifyReq.CertId)
	assert.Equal(t, "target.demo.com", *modifyReq.TargetName)

	var got TeoAliasDomain110903ResourceModel
	diags = resp.State.Get(context.Background(), &got)
	if diags.HasError() {
		t.Fatalf("reading back the state failed: %v", diags)
	}
	assert.True(t, got.CertType.IsNull(), "cert_type must stay null in the refreshed state")
	assert.True(t, got.CertId.IsNull(), "cert_id must stay null in the refreshed state")
}

// TestTeoAliasDomain110903Resource_Update_ExplicitNoneClearsCert covers the
// explicit clear path: cert_type = "none" is sent to the cloud side.
func TestTeoAliasDomain110903Resource_Update_ExplicitNoneClearsCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	var modifyReq *teov20220901.ModifyAliasDomainRequest
	mockTeoAliasDomain110903Modify(patches, &modifyReq, nil)

	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		return &teov20220901.AliasDomain{
			AliasName:  ptrStrTeoAliasDomain110903("alias.demo.com"),
			ZoneId:     ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6"),
			TargetName: ptrStrTeoAliasDomain110903("target.demo.com"),
			Status:     ptrStrTeoAliasDomain110903("active"),
		}, nil
	})

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
		CertType:   types.StringValue("hosting"),
	}
	planModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
		CertType:   types.StringValue("none"),
	}

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, planModel)},
		Plan:   teoAliasDomain110903Plan(t, s, planModel),
		State:  teoAliasDomain110903State(t, s, stateModel),
	}
	resp := &resource.UpdateResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, modifyReq)
	assert.Equal(t, "none", *modifyReq.CertType)
}

// TestTeoAliasDomain110903Resource_Update_ApiError covers the API error branch.
func TestTeoAliasDomain110903Resource_Update_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	var modifyReq *teov20220901.ModifyAliasDomainRequest
	mockTeoAliasDomain110903Modify(patches, &modifyReq, fmt.Errorf("InternalError"))

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
	}
	planModel := stateModel

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, planModel)},
		Plan:   teoAliasDomain110903Plan(t, s, planModel),
		State:  teoAliasDomain110903State(t, s, stateModel),
	}
	resp := &resource.UpdateResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics from a failed ModifyAliasDomain")
}

// TestTeoAliasDomain110903Resource_Update_BrokenId covers the invalid-id branch.
func TestTeoAliasDomain110903Resource_Update_BrokenId(t *testing.T) {
	r := &TeoAliasDomain110903Resource{}

	s := teoAliasDomain110903ResourceSchema(t)
	stateModel := &TeoAliasDomain110903ResourceModel{Id: types.StringValue("broken-id")}

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, stateModel)},
		Plan:   teoAliasDomain110903Plan(t, s, stateModel),
		State:  teoAliasDomain110903State(t, s, stateModel),
	}
	resp := &resource.UpdateResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for a broken id")
}

// TestTeoAliasDomain110903Resource_Delete_Success covers a successful delete:
// query first, then delete with a single-element alias list.
func TestTeoAliasDomain110903Resource_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	describeCalls := 0
	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		describeCalls++
		return &teov20220901.AliasDomain{
			AliasName: ptrStrTeoAliasDomain110903("alias.demo.com"),
			ZoneId:    ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6"),
			Status:    ptrStrTeoAliasDomain110903("active"),
		}, nil
	})

	var deleteErr error
	deleteCalls := 0
	patches.ApplyMethodFunc(&TeoService{}, "DeleteTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) error {
		deleteCalls++
		return deleteErr
	})

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
	}

	req := resource.DeleteRequest{State: teoAliasDomain110903State(t, s, stateModel)}
	resp := &resource.DeleteResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.Equal(t, 1, describeCalls)
	assert.Equal(t, 1, deleteCalls)
}

// TestTeoAliasDomain110903Resource_Delete_AlreadyGone covers the confirmed
// already-gone branch: no delete call is issued and no error is raised.
func TestTeoAliasDomain110903Resource_Delete_AlreadyGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		return nil, nil
	})

	deleteCalls := 0
	patches.ApplyMethodFunc(&TeoService{}, "DeleteTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) error {
		deleteCalls++
		return nil
	})

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
	}

	req := resource.DeleteRequest{State: teoAliasDomain110903State(t, s, stateModel)}
	resp := &resource.DeleteResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.Equal(t, 0, deleteCalls, "no delete call is expected when the alias domain is already gone")
}

// TestTeoAliasDomain110903Resource_Delete_ApiError covers the API error branch.
func TestTeoAliasDomain110903Resource_Delete_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}

	patches.ApplyMethodFunc(&TeoService{}, "DescribeTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
		return &teov20220901.AliasDomain{
			AliasName: ptrStrTeoAliasDomain110903("alias.demo.com"),
			ZoneId:    ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6"),
			Status:    ptrStrTeoAliasDomain110903("active"),
		}, nil
	})
	patches.ApplyMethodFunc(&TeoService{}, "DeleteTeoAliasDomainById", func(_ context.Context, zoneId, aliasName string) error {
		return fmt.Errorf("InternalError")
	})

	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}
	setTeoAliasDomain110903Client(r, client)

	stateModel := &TeoAliasDomain110903ResourceModel{
		Id:         types.StringValue("zone-2qtuhspy7cr6#alias.demo.com"),
		ZoneId:     types.StringValue("zone-2qtuhspy7cr6"),
		AliasName:  types.StringValue("alias.demo.com"),
		TargetName: types.StringValue("target.demo.com"),
	}

	req := resource.DeleteRequest{State: teoAliasDomain110903State(t, s, stateModel)}
	resp := &resource.DeleteResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Delete(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics from a failed DeleteAliasDomain")
}

// TestTeoAliasDomain110903Resource_Delete_BrokenId covers the invalid-id branch.
func TestTeoAliasDomain110903Resource_Delete_BrokenId(t *testing.T) {
	r := &TeoAliasDomain110903Resource{}

	s := teoAliasDomain110903ResourceSchema(t)
	stateModel := &TeoAliasDomain110903ResourceModel{Id: types.StringValue("broken-id")}

	req := resource.DeleteRequest{State: teoAliasDomain110903State(t, s, stateModel)}
	resp := &resource.DeleteResponse{State: teoAliasDomain110903State(t, s, stateModel)}
	r.Delete(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for a broken id")
}

// TestTeoAliasDomain110903Resource_ImportState covers the composite import id:
// valid ids fill zone_id / alias_name / id; broken ids raise an error.
func TestTeoAliasDomain110903Resource_ImportState(t *testing.T) {
	s := teoAliasDomain110903ResourceSchema(t)
	r := &TeoAliasDomain110903Resource{}

	// Valid id.
	req := resource.ImportStateRequest{ID: "zone-2qtuhspy7cr6#alias.demo.com"}
	resp := &resource.ImportStateResponse{
		State: tfsdk.State{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, &TeoAliasDomain110903ResourceModel{})},
	}
	r.ImportState(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var got TeoAliasDomain110903ResourceModel
	diags := resp.State.Get(context.Background(), &got)
	if diags.HasError() {
		t.Fatalf("reading back the state failed: %v", diags)
	}
	assert.Equal(t, "zone-2qtuhspy7cr6#alias.demo.com", got.Id.ValueString())
	assert.Equal(t, "zone-2qtuhspy7cr6", got.ZoneId.ValueString())
	assert.Equal(t, "alias.demo.com", got.AliasName.ValueString())

	// Broken ids.
	for _, id := range []string{"", "only-zone", "zone#alias#extra"} {
		badReq := resource.ImportStateRequest{ID: id}
		badResp := &resource.ImportStateResponse{
			State: tfsdk.State{Schema: s, Raw: teoAliasDomain110903RawValue(t, s, &TeoAliasDomain110903ResourceModel{})},
		}
		r.ImportState(context.Background(), badReq, badResp)
		assert.True(t, badResp.Diagnostics.HasError(), "expected error diagnostics for import id %q", id)
	}
}

// TestTeoAliasDomain110903Service_Create covers the service-layer conversion
// for CreateAliasDomain: null optional fields are omitted.
func TestTeoAliasDomain110903Service_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotRequest *teov20220901.CreateAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teov20220901.CreateAliasDomainRequest) (*teov20220901.CreateAliasDomainResponse, error) {
		gotRequest = request
		resp := teov20220901.NewCreateAliasDomainResponse()
		resp.Response = &teov20220901.CreateAliasDomainResponseParams{
			RequestId: ptrStrTeoAliasDomain110903("fake-request-id"),
		}
		return resp, nil
	})

	request := teov20220901.NewCreateAliasDomainRequest()
	request.ZoneId = ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6")
	request.AliasName = ptrStrTeoAliasDomain110903("alias.demo.com")
	request.TargetName = ptrStrTeoAliasDomain110903("target.demo.com")
	// cert_type / cert_id omitted on purpose.

	service := NewTeoService(client)
	err := service.CreateTeoAliasDomain(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, gotRequest)
	assert.Equal(t, "zone-2qtuhspy7cr6", *gotRequest.ZoneId)
	assert.Equal(t, "alias.demo.com", *gotRequest.AliasName)
	assert.Equal(t, "target.demo.com", *gotRequest.TargetName)
	assert.Nil(t, gotRequest.CertType)
	assert.Nil(t, gotRequest.CertId)
}

// TestTeoAliasDomain110903Service_Create_ApiError covers the service-layer
// API error branch.
func TestTeoAliasDomain110903Service_Create_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teov20220901.CreateAliasDomainRequest) (*teov20220901.CreateAliasDomainResponse, error) {
		return nil, fmt.Errorf("InternalError")
	})

	request := teov20220901.NewCreateAliasDomainRequest()
	request.ZoneId = ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6")
	request.AliasName = ptrStrTeoAliasDomain110903("alias.demo.com")
	request.TargetName = ptrStrTeoAliasDomain110903("target.demo.com")

	service := NewTeoService(client)
	err := service.CreateTeoAliasDomain(context.Background(), request)

	assert.Error(t, err)
}

// TestTeoAliasDomain110903Service_Modify covers the service-layer conversion
// for ModifyAliasDomain.
func TestTeoAliasDomain110903Service_Modify(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotRequest *teov20220901.ModifyAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *teov20220901.ModifyAliasDomainRequest) (*teov20220901.ModifyAliasDomainResponse, error) {
		gotRequest = request
		resp := teov20220901.NewModifyAliasDomainResponse()
		resp.Response = &teov20220901.ModifyAliasDomainResponseParams{
			RequestId: ptrStrTeoAliasDomain110903("fake-request-id"),
		}
		return resp, nil
	})

	request := teov20220901.NewModifyAliasDomainRequest()
	request.ZoneId = ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6")
	request.AliasName = ptrStrTeoAliasDomain110903("alias.demo.com")
	request.TargetName = ptrStrTeoAliasDomain110903("new-target.demo.com")
	request.CertType = ptrStrTeoAliasDomain110903("none")

	service := NewTeoService(client)
	err := service.ModifyTeoAliasDomain(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, gotRequest)
	assert.Equal(t, "none", *gotRequest.CertType)
	assert.Nil(t, gotRequest.CertId)
}

// TestTeoAliasDomain110903Service_Modify_ApiError covers the service-layer API
// error branch.
func TestTeoAliasDomain110903Service_Modify_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *teov20220901.ModifyAliasDomainRequest) (*teov20220901.ModifyAliasDomainResponse, error) {
		return nil, fmt.Errorf("InternalError")
	})

	request := teov20220901.NewModifyAliasDomainRequest()
	request.ZoneId = ptrStrTeoAliasDomain110903("zone-2qtuhspy7cr6")
	request.AliasName = ptrStrTeoAliasDomain110903("alias.demo.com")
	request.TargetName = ptrStrTeoAliasDomain110903("new-target.demo.com")

	service := NewTeoService(client)
	err := service.ModifyTeoAliasDomain(context.Background(), request)

	assert.Error(t, err)
}

// TestTeoAliasDomain110903Service_Delete covers the service-layer conversion
// for DeleteAliasDomain: the alias name list carries exactly one element.
func TestTeoAliasDomain110903Service_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotRequest *teov20220901.DeleteAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *teov20220901.DeleteAliasDomainRequest) (*teov20220901.DeleteAliasDomainResponse, error) {
		gotRequest = request
		resp := teov20220901.NewDeleteAliasDomainResponse()
		resp.Response = &teov20220901.DeleteAliasDomainResponseParams{
			RequestId: ptrStrTeoAliasDomain110903("fake-request-id"),
		}
		return resp, nil
	})

	service := NewTeoService(client)
	err := service.DeleteTeoAliasDomainById(context.Background(), "zone-2qtuhspy7cr6", "alias.demo.com")

	assert.NoError(t, err)
	assert.NotNil(t, gotRequest)
	assert.Equal(t, "zone-2qtuhspy7cr6", *gotRequest.ZoneId)
	assert.Equal(t, []*string{ptrStrTeoAliasDomain110903("alias.demo.com")}, gotRequest.AliasNames)
}

// TestTeoAliasDomain110903Service_Delete_ApiError covers the service-layer API
// error branch.
func TestTeoAliasDomain110903Service_Delete_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *teov20220901.DeleteAliasDomainRequest) (*teov20220901.DeleteAliasDomainResponse, error) {
		return nil, fmt.Errorf("InternalError")
	})

	service := NewTeoService(client)
	err := service.DeleteTeoAliasDomainById(context.Background(), "zone-2qtuhspy7cr6", "alias.demo.com")

	assert.Error(t, err)
}

// TestTeoAliasDomain110903Resource_ModelCoversAllSchemaAttributes guards
// against schema/model drift: every schema attribute has a model field.
func TestTeoAliasDomain110903Resource_ModelCoversAllSchemaAttributes(t *testing.T) {
	s := teoAliasDomain110903ResourceSchema(t)

	modelAttrs := map[string]bool{
		"id":          true,
		"zone_id":     true,
		"alias_name":  true,
		"target_name": true,
		"cert_type":   true,
		"cert_id":     true,
		"status":      true,
		"forbid_mode": true,
		"created_on":  true,
		"modified_on": true,
	}
	for name := range s.Attributes {
		if !modelAttrs[name] {
			t.Errorf("schema attribute %q has no model field", name)
		}
	}
	// The raw value builder must also know every attribute: build one value to
	// exercise the default branch guard.
	_ = teoAliasDomain110903RawValue(t, s, &TeoAliasDomain110903ResourceModel{})
}
