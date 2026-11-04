package teo_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringOriginGroupV2(s string) *string { return &s }

func ptrUint64OriginGroupV2(v uint64) *uint64 { return &v }

// originGroupV2Resource returns the resource under test.
func originGroupV2Resource() resource.Resource {
	return svcteo.NewTeoOriginGroupV2Resource()
}

// originGroupV2ResourceSchema returns the framework schema of the resource
// under test, failing the test when the Schema method reports errors.
func originGroupV2ResourceSchema(t *testing.T) resourceschema.Schema {
	r := originGroupV2Resource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// originGroupV2PlanAttrs returns the raw plan attribute values used to seed
// the plan/config in Create and Update tests. zoneId drives zone_id; name,
// type and host_header keep the values asserted on the modify path.
func originGroupV2PlanAttrs(zoneId string) map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, nil),
		"zone_id":     tftypes.NewValue(tftypes.String, zoneId),
		"name":        tftypes.NewValue(tftypes.String, "tf-test-origin-group"),
		"type":        tftypes.NewValue(tftypes.String, "GENERAL"),
		"host_header": tftypes.NewValue(tftypes.String, nil),
		"create_time": tftypes.NewValue(tftypes.String, nil),
		"update_time": tftypes.NewValue(tftypes.String, nil),
		"records": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"record":    tftypes.String,
			"type":      tftypes.String,
			"record_id": tftypes.String,
			"weight":    tftypes.Number,
			"private":   tftypes.Bool,
			"private_parameters": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"name":  tftypes.String,
				"value": tftypes.String,
			}}},
		}}}, []tftypes.Value{
			tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"record":    tftypes.String,
				"type":      tftypes.String,
				"record_id": tftypes.String,
				"weight":    tftypes.Number,
				"private":   tftypes.Bool,
				"private_parameters": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
					"name":  tftypes.String,
					"value": tftypes.String,
				}}},
			}}, map[string]tftypes.Value{
				"record":            tftypes.NewValue(tftypes.String, "tf-teo.xyz"),
				"type":              tftypes.NewValue(tftypes.String, "IP_DOMAIN"),
				"record_id":         tftypes.NewValue(tftypes.String, nil),
				"weight":            tftypes.NewValue(tftypes.Number, big.NewFloat(100)),
				"private":           tftypes.NewValue(tftypes.Bool, false),
				"private_parameters": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
					"name":  tftypes.String,
					"value": tftypes.String,
				}}}, nil),
			}),
		}),
		"references": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"instance_type":   tftypes.String,
			"instance_id":     tftypes.String,
			"instance_name":   tftypes.String,
			"zone_id":         tftypes.String,
			"zone_name":       tftypes.String,
			"alias_zone_name": tftypes.String,
		}}}, nil),
	}
}

// originGroupV2Plan builds a tfsdk.Plan carrying the default test attributes.
func originGroupV2Plan(t *testing.T, zoneId string) tfsdk.Plan {
	return originGroupV2RawAsPlan(t, originGroupV2PlanAttrs(zoneId))
}

// originGroupV2Config builds a tfsdk.Config carrying the default test attributes.
func originGroupV2Config(t *testing.T, zoneId string) tfsdk.Config {
	plan := originGroupV2Plan(t, zoneId)
	return tfsdk.Config{Raw: plan.Raw, Schema: plan.Schema}
}

// originGroupV2RawAsPlan wraps raw attribute values into a tfsdk.Plan using
// the resource schema.
func originGroupV2RawAsPlan(t *testing.T, attrs map[string]tftypes.Value) tfsdk.Plan {
	return tfsdk.Plan{
		Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: originGroupV2AttrTypes()}, attrs),
		Schema: originGroupV2ResourceSchema(t),
	}
}

// originGroupV2ReadState builds a tfsdk.State carrying only the id attribute,
// mirroring a state freshly imported or after a create.
func originGroupV2ReadState(t *testing.T, id string) tfsdk.State {
	attrs := map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, id),
		"zone_id":     tftypes.NewValue(tftypes.String, nil),
		"name":        tftypes.NewValue(tftypes.String, nil),
		"type":        tftypes.NewValue(tftypes.String, nil),
		"host_header": tftypes.NewValue(tftypes.String, nil),
		"create_time": tftypes.NewValue(tftypes.String, nil),
		"update_time": tftypes.NewValue(tftypes.String, nil),
		"records":     tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"record":    tftypes.String,
			"type":      tftypes.String,
			"record_id": tftypes.String,
			"weight":    tftypes.Number,
			"private":   tftypes.Bool,
			"private_parameters": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"name":  tftypes.String,
				"value": tftypes.String,
			}}},
		}}}, nil),
		"references": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"instance_type":   tftypes.String,
			"instance_id":     tftypes.String,
			"instance_name":   tftypes.String,
			"zone_id":         tftypes.String,
			"zone_name":       tftypes.String,
			"alias_zone_name": tftypes.String,
		}}}, nil),
	}
	return tfsdk.State{
		Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: originGroupV2AttrTypes()}, attrs),
		Schema: originGroupV2ResourceSchema(t),
	}
}

// originGroupV2CreateResponse builds a CreateResponse whose State is
// pre-populated from the plan, mirroring the framework protocol behavior.
func originGroupV2CreateResponse(plan tfsdk.Plan) *resource.CreateResponse {
	resp := &resource.CreateResponse{}
	resp.State = tfsdk.State{Raw: plan.Raw, Schema: plan.Schema}
	return resp
}

// originGroupV2ReadResponse builds a ReadResponse whose State is pre-populated
// from the request state, mirroring the framework protocol behavior.
func originGroupV2ReadResponse(state tfsdk.State) *resource.ReadResponse {
	resp := &resource.ReadResponse{}
	resp.State = tfsdk.State{Raw: state.Raw, Schema: state.Schema}
	return resp
}

// originGroupV2StateFromRaw unmarshals a raw state value into the resource model.
func originGroupV2StateFromRaw(t *testing.T, raw tftypes.Value) (svcteo.TeoOriginGroupV2ResourceModel, bool) {
	var model svcteo.TeoOriginGroupV2ResourceModel
	state := tfsdk.State{
		Raw:    raw,
		Schema: originGroupV2ResourceSchema(t),
	}
	diags := state.Get(context.Background(), &model)
	if diags.HasError() {
		return model, false
	}
	return model, true
}

// originGroupV2RecordsFromState unmarshals the records list of the model.
func originGroupV2RecordsFromState(t *testing.T, list types.List) []svcteo.TeoOriginGroupV2RecordModel {
	records := make([]svcteo.TeoOriginGroupV2RecordModel, 0, len(list.Elements()))
	_ = list.ElementsAs(context.Background(), &records, false)
	return records
}

// originGroupV2PrivateParametersFromState unmarshals the private_parameters
// list of a record model.
func originGroupV2PrivateParametersFromState(t *testing.T, list types.List) []svcteo.TeoOriginGroupV2PrivateParameterModel {
	privateParameters := make([]svcteo.TeoOriginGroupV2PrivateParameterModel, 0, len(list.Elements()))
	_ = list.ElementsAs(context.Background(), &privateParameters, false)
	return privateParameters
}

// originGroupV2ReferencesFromState unmarshals the references list of the model.
func originGroupV2ReferencesFromState(t *testing.T, list types.List) []svcteo.TeoOriginGroupV2ReferenceModel {
	references := make([]svcteo.TeoOriginGroupV2ReferenceModel, 0, len(list.Elements()))
	_ = list.ElementsAs(context.Background(), &references, false)
	return references
}

// originGroupV2AttrTypes returns the top level attribute types of the resource
// schema (mirrors the schema declaration in the resource file).
func originGroupV2AttrTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"id":           tftypes.String,
		"zone_id":      tftypes.String,
		"name":         tftypes.String,
		"type":         tftypes.String,
		"host_header":  tftypes.String,
		"create_time":  tftypes.String,
		"update_time":  tftypes.String,
		"records":      tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"record":    tftypes.String,
			"type":      tftypes.String,
			"record_id": tftypes.String,
			"weight":    tftypes.Number,
			"private":   tftypes.Bool,
			"private_parameters": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"name":  tftypes.String,
				"value": tftypes.String,
			}}},
		}}},
		"references": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"instance_type":   tftypes.String,
			"instance_id":     tftypes.String,
			"instance_name":   tftypes.String,
			"zone_id":         tftypes.String,
			"zone_name":       tftypes.String,
			"alias_zone_name": tftypes.String,
		}}},
	}
}

// originGroupV2TerraformValue builds the raw top level tftypes.Value used to
// seed plan/state in the tests.
func originGroupV2TerraformValue(attrs map[string]tftypes.Value) tftypes.Value {
	return tftypes.NewValue(tftypes.Object{AttributeTypes: originGroupV2AttrTypes()}, attrs)
}

// configureOriginGroupV2Client injects a mock connectivity client into the
// resource instance through the framework Configure method.
func configureOriginGroupV2Client(r resource.ResourceWithConfigure, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	req := resource.ConfigureRequest{ProviderData: meta}
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), req, resp)
}

// patchTeoOriginGroupV2Client patches UseTeoClient on a mock connectivity
// client and returns the fake teo client.
func patchTeoOriginGroupV2Client(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *teov20220901.Client {
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)
	return teoClient
}

// mockOriginGroupV2Response builds a DescribeOriginGroup response holding a
// single origin group with every leaf field populated.
func mockOriginGroupV2Response(group *teov20220901.OriginGroup) *teov20220901.DescribeOriginGroupResponse {
	resp := teov20220901.NewDescribeOriginGroupResponse()
	resp.Response = &teov20220901.DescribeOriginGroupResponseParams{
		OriginGroups: []*teov20220901.OriginGroup{group},
		RequestId:    ptrStringOriginGroupV2("fake-request-id"),
	}
	return resp
}

// originGroupV2Full builds a fully populated origin group for read assertions.
func originGroupV2Full() *teov20220901.OriginGroup {
	return &teov20220901.OriginGroup{
		GroupId:     ptrStringOriginGroupV2("origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
		Name:        ptrStringOriginGroupV2("tf-test-origin-group"),
		Type:        ptrStringOriginGroupV2("GENERAL"),
		HostHeader:  ptrStringOriginGroupV2("origin.example.com"),
		CreateTime:  ptrStringOriginGroupV2("2024-01-01T00:00:00Z"),
		UpdateTime:  ptrStringOriginGroupV2("2024-01-02T00:00:00Z"),
		Records: []*teov20220901.OriginRecord{
			{
				Record: ptrStringOriginGroupV2("1.1.1.1"),
				Type:   ptrStringOriginGroupV2("IP_DOMAIN"),
				Weight: ptrUint64OriginGroupV2(100),
				Private: new(bool),
			},
			{
				Record:   ptrStringOriginGroupV2("tf-teo.xyz"),
				Type:     ptrStringOriginGroupV2("AWS_S3"),
				RecordId: ptrStringOriginGroupV2("record-0001"),
				Weight:   ptrUint64OriginGroupV2(0),
				Private:   ptrBoolOriginGroupV2(true),
				PrivateParameters: []*teov20220901.PrivateParameter{
					{
						Name:  ptrStringOriginGroupV2("AccessKeyId"),
						Value: ptrStringOriginGroupV2("AKIDtest"),
					},
					{
						Name:  ptrStringOriginGroupV2("SecretAccessKey"),
						Value: ptrStringOriginGroupV2("secret-test"),
					},
				},
			},
		},
		References: []*teov20220901.OriginGroupReference{
			{
				InstanceType:  ptrStringOriginGroupV2("acceleration-domain"),
				InstanceId:    ptrStringOriginGroupV2("acc-domain-001"),
				InstanceName:  ptrStringOriginGroupV2("test-domain.com"),
				ZoneId:        ptrStringOriginGroupV2("zone-abc123"),
				ZoneName:      ptrStringOriginGroupV2("test-zone.com"),
				AliasZoneName: ptrStringOriginGroupV2("alias-test-zone.com"),
			},
		},
	}
}

func ptrBoolOriginGroupV2(v bool) *bool { return &v }

// go test ./tencentcloud/services/teo/ -run "TestTeoOriginGroupV2" -v -count=1 -gcflags="all=-l"

// TestTeoOriginGroupV2Resource_Metadata validates the resource type name.
func TestTeoOriginGroupV2Resource_Metadata(t *testing.T) {
	r := originGroupV2Resource()
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	assert.Equal(t, "tencentcloud_teo_origin_group_v2", resp.TypeName)
}

// TestTeoOriginGroupV2Resource_SchemaFlags validates required / optional /
// computed flags on every attribute.
func TestTeoOriginGroupV2Resource_SchemaFlags(t *testing.T) {
	s := originGroupV2ResourceSchema(t)

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":          {computed: true},
		"zone_id":     {required: true},
		"name":        {optional: true, computed: true},
		"type":        {required: true},
		"host_header": {optional: true, computed: true},
		"records":     {required: true},
		"create_time": {computed: true},
		"update_time": {computed: true},
		"references":  {computed: true},
	}
	for name, want := range cases {
		attr, ok := s.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		if attr.IsRequired() != want.required {
			t.Errorf("attribute %q IsRequired=%v, want %v", name, attr.IsRequired(), want.required)
		}
		if attr.IsOptional() != want.optional {
			t.Errorf("attribute %q IsOptional=%v, want %v", name, attr.IsOptional(), want.optional)
		}
		if attr.IsComputed() != want.computed {
			t.Errorf("attribute %q IsComputed=%v, want %v", name, attr.IsComputed(), want.computed)
		}
	}

	records := s.Attributes["records"].(resourceschema.ListNestedAttribute)
	if _, ok := records.NestedObject.Attributes["record"]; !ok {
		t.Error("records nested attribute record missing")
	}
	privateParams := records.NestedObject.Attributes["private_parameters"].(resourceschema.ListNestedAttribute)
	if privateParams.NestedObject.Attributes["name"].IsRequired() != true {
		t.Error("private_parameters.name should be required")
	}
	if !privateParams.NestedObject.Attributes["value"].(resourceschema.StringAttribute).Sensitive {
		t.Error("private_parameters.value should be sensitive")
	}
}

// TestTeoOriginGroupV2Resource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces it advertises.
func TestTeoOriginGroupV2Resource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &svcteo.TeoOriginGroupV2Resource{}
	var _ resource.ResourceWithImportState = &svcteo.TeoOriginGroupV2Resource{}
}

// TestTeoOriginGroupV2Resource_Factory validates the registry factory.
func TestTeoOriginGroupV2Resource_Factory(t *testing.T) {
	got := svcteo.NewTeoOriginGroupV2Resource()
	if got == nil {
		t.Fatal("NewTeoOriginGroupV2Resource returned nil")
	}
	if _, ok := got.(*svcteo.TeoOriginGroupV2Resource); !ok {
		t.Fatalf("factory returned %T, want *svcteo.TeoOriginGroupV2Resource", got)
	}
}

// TestTeoOriginGroupV2Resource_ConfigureNilProviderData ensures Configure with
// nil ProviderData is a no-op without diagnostics.
func TestTeoOriginGroupV2Resource_ConfigureNilProviderData(t *testing.T) {
	r := svcteo.NewTeoOriginGroupV2Resource().(*svcteo.TeoOriginGroupV2Resource)
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{}, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

// TestTeoOriginGroupV2Resource_Create tests Create with a mocked Create and
// Describe call, verifying the request conversion and the full state write.
func TestTeoOriginGroupV2Resource_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	var gotCreate *teov20220901.CreateOriginGroupRequest
	patches.ApplyMethodFunc(teoClient, "CreateOriginGroupWithContext",
		func(_ context.Context, request *teov20220901.CreateOriginGroupRequest) (*teov20220901.CreateOriginGroupResponse, error) {
			gotCreate = request
			resp := teov20220901.NewCreateOriginGroupResponse()
			resp.Response = &teov20220901.CreateOriginGroupResponseParams{
				OriginGroupId: ptrStringOriginGroupV2("origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
				RequestId:    ptrStringOriginGroupV2("fake-request-id"),
			}
			return resp, nil
		})

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup",
		func(_ *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
			return mockOriginGroupV2Response(originGroupV2Full()), nil
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	req := resource.CreateRequest{
		Plan:   originGroupV2Plan(t, "zone-297z8rf93cfw"),
		Config: originGroupV2Config(t, "zone-297z8rf93cfw"),
	}
	resp := originGroupV2CreateResponse(req.Plan)
	r.Create(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotCreate)
	assert.Equal(t, "zone-297z8rf93cfw", *gotCreate.ZoneId)
	assert.Equal(t, "tf-test-origin-group", *gotCreate.Name)
	assert.Equal(t, "GENERAL", *gotCreate.Type)
	assert.Len(t, gotCreate.Records, 1)
	assert.Equal(t, "tf-teo.xyz", *gotCreate.Records[0].Record)
	assert.Equal(t, "IP_DOMAIN", *gotCreate.Records[0].Type)
	assert.Equal(t, uint64(100), *gotCreate.Records[0].Weight)
	assert.Equal(t, false, *gotCreate.Records[0].Private)

	// Verify the state write.
	model, ok := originGroupV2StateFromRaw(t, resp.State.Raw)
	assert.True(t, ok)
	assert.Equal(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a", model.Id.ValueString())
	assert.Equal(t, "tf-test-origin-group", model.Name.ValueString())
	assert.Equal(t, "GENERAL", model.Type.ValueString())
	assert.Equal(t, "origin.example.com", model.HostHeader.ValueString())
	assert.Equal(t, "2024-01-01T00:00:00Z", model.CreateTime.ValueString())
	assert.Equal(t, "2024-01-02T00:00:00Z", model.UpdateTime.ValueString())

	records := originGroupV2RecordsFromState(t, model.Records)
	assert.Len(t, records, 2)
	assert.Equal(t, "1.1.1.1", records[0].Record.ValueString())
	assert.Equal(t, "IP_DOMAIN", records[0].Type.ValueString())
	assert.Equal(t, int64(100), records[0].Weight.ValueInt64())
	assert.False(t, records[0].Private.ValueBool())
	assert.True(t, records[0].RecordId.IsNull())
	assert.Equal(t, "tf-teo.xyz", records[1].Record.ValueString())
	assert.Equal(t, "AWS_S3", records[1].Type.ValueString())
	assert.Equal(t, "record-0001", records[1].RecordId.ValueString())
	assert.Equal(t, int64(0), records[1].Weight.ValueInt64())
	assert.True(t, records[1].Private.ValueBool())
	privateParameters := originGroupV2PrivateParametersFromState(t, records[1].PrivateParameters)
	assert.Len(t, privateParameters, 2)
	assert.Equal(t, "AccessKeyId", privateParameters[0].Name.ValueString())
	assert.Equal(t, "AKIDtest", privateParameters[0].Value.ValueString())
	assert.Equal(t, "SecretAccessKey", privateParameters[1].Name.ValueString())
	assert.Equal(t, "secret-test", privateParameters[1].Value.ValueString())

	references := originGroupV2ReferencesFromState(t, model.References)
	assert.Len(t, references, 1)
	assert.Equal(t, "acceleration-domain", references[0].InstanceType.ValueString())
	assert.Equal(t, "acc-domain-001", references[0].InstanceId.ValueString())
	assert.Equal(t, "test-domain.com", references[0].InstanceName.ValueString())
	assert.Equal(t, "zone-abc123", references[0].ZoneId.ValueString())
	assert.Equal(t, "test-zone.com", references[0].ZoneName.ValueString())
	assert.Equal(t, "alias-test-zone.com", references[0].AliasZoneName.ValueString())
}

// TestTeoOriginGroupV2Resource_CreateApiError tests that an API error during
// Create is reported through diagnostics and no state is written.
func TestTeoOriginGroupV2Resource_CreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	patches.ApplyMethodFunc(teoClient, "CreateOriginGroupWithContext",
		func(_ context.Context, request *teov20220901.CreateOriginGroupRequest) (*teov20220901.CreateOriginGroupResponse, error) {
			return nil, fmt.Errorf("internal error")
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	req := resource.CreateRequest{
		Plan:   originGroupV2Plan(t, "zone-297z8rf93cfw"),
		Config: originGroupV2Config(t, "zone-297z8rf93cfw"),
	}
	resp := originGroupV2CreateResponse(req.Plan)
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestTeoOriginGroupV2Resource_ReadFound tests Read when the origin group exists.
func TestTeoOriginGroupV2Resource_ReadFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup",
		func(_ *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
			return mockOriginGroupV2Response(originGroupV2Full()), nil
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	req := resource.ReadRequest{
		State: originGroupV2ReadState(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
	}
	resp := originGroupV2ReadResponse(req.State)
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	model, ok := originGroupV2StateFromRaw(t, resp.State.Raw)
	assert.True(t, ok)
	assert.Equal(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a", model.Id.ValueString())
	assert.Equal(t, "tf-test-origin-group", model.Name.ValueString())
	assert.Equal(t, "GENERAL", model.Type.ValueString())
	assert.Equal(t, "origin.example.com", model.HostHeader.ValueString())

	records := originGroupV2RecordsFromState(t, model.Records)
	assert.Len(t, records, 2)
	assert.Equal(t, int64(0), records[1].Weight.ValueInt64())
	privateParameters := originGroupV2PrivateParametersFromState(t, records[1].PrivateParameters)
	assert.Len(t, privateParameters, 2)

	references := originGroupV2ReferencesFromState(t, model.References)
	assert.Len(t, references, 1)
	assert.Equal(t, "acc-domain-001", references[0].InstanceId.ValueString())
}

// TestTeoOriginGroupV2Resource_ReadNotFound tests that an empty result removes
// the resource from state.
func TestTeoOriginGroupV2Resource_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup",
		func(_ *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
			resp := teov20220901.NewDescribeOriginGroupResponse()
			resp.Response = &teov20220901.DescribeOriginGroupResponseParams{
				OriginGroups: []*teov20220901.OriginGroup{},
				RequestId:    ptrStringOriginGroupV2("fake-request-id"),
			}
			return resp, nil
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	req := resource.ReadRequest{
		State: originGroupV2ReadState(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
	}
	resp := originGroupV2ReadResponse(req.State)
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.True(t, resp.State.Raw.IsNull(), "expected the state to be removed")
}

// TestTeoOriginGroupV2Resource_ReadNilResponseIsError tests that a nil
// response is an API failure, not an empty result.
func TestTeoOriginGroupV2Resource_ReadNilResponseIsError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup",
		func(_ *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
			return nil, fmt.Errorf("api failure")
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	req := resource.ReadRequest{
		State: originGroupV2ReadState(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
	}
	resp := originGroupV2ReadResponse(req.State)
	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestTeoOriginGroupV2Resource_ReadBrokenId tests that a broken id errors out.
func TestTeoOriginGroupV2Resource_ReadBrokenId(t *testing.T) {
	r := originGroupV2Resource()

	req := resource.ReadRequest{
		State: originGroupV2ReadState(t, "broken-id-without-separator"),
	}
	resp := &resource.ReadResponse{}
	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestTeoOriginGroupV2Resource_Update tests Update with name and host_header
// changes: ModifyOriginGroup receives the planned values and the state is
// refreshed from the read-back.
func TestTeoOriginGroupV2Resource_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	var gotModify *teov20220901.ModifyOriginGroupRequest
	patches.ApplyMethodFunc(teoClient, "ModifyOriginGroupWithContext",
		func(_ context.Context, request *teov20220901.ModifyOriginGroupRequest) (*teov20220901.ModifyOriginGroupResponse, error) {
			gotModify = request
			resp := teov20220901.NewModifyOriginGroupResponse()
			resp.Response = &teov20220901.ModifyOriginGroupResponseParams{
				RequestId: ptrStringOriginGroupV2("fake-request-id"),
			}
			return resp, nil
		})

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup",
		func(_ *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
			group := originGroupV2Full()
			group.Name = ptrStringOriginGroupV2("tf-test-origin-group-renamed")
			group.HostHeader = ptrStringOriginGroupV2("renamed.example.com")
			return mockOriginGroupV2Response(group), nil
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	plan := originGroupV2Plan(t, "zone-297z8rf93cfw")
	// Simulate a name + host_header change in the plan.
	planAttrs := originGroupV2PlanAttrs("zone-297z8rf93cfw")
	planAttrs["name"] = tftypes.NewValue(tftypes.String, "tf-test-origin-group-renamed")
	planAttrs["host_header"] = tftypes.NewValue(tftypes.String, "renamed.example.com")
	plan = originGroupV2RawAsPlan(t, planAttrs)

	req := resource.UpdateRequest{
		Plan:   plan,
		Config: tfsdk.Config{Raw: plan.Raw, Schema: plan.Schema},
		State:  originGroupV2ReadState(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema}}
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotModify)
	assert.Equal(t, "zone-297z8rf93cfw", *gotModify.ZoneId)
	assert.Equal(t, "origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a", *gotModify.GroupId)
	assert.Equal(t, "tf-test-origin-group-renamed", *gotModify.Name)
	assert.Equal(t, "GENERAL", *gotModify.Type)
	assert.Equal(t, "renamed.example.com", *gotModify.HostHeader)
	// The whole records list is sent as the new record set on update.
	assert.Len(t, gotModify.Records, 1)
	assert.Equal(t, "tf-teo.xyz", *gotModify.Records[0].Record)
	assert.Equal(t, "IP_DOMAIN", *gotModify.Records[0].Type)
	assert.Equal(t, uint64(100), *gotModify.Records[0].Weight)

	model, ok := originGroupV2StateFromRaw(t, resp.State.Raw)
	assert.True(t, ok)
	assert.Equal(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a", model.Id.ValueString())
	assert.Equal(t, "tf-test-origin-group-renamed", model.Name.ValueString())
	assert.Equal(t, "renamed.example.com", model.HostHeader.ValueString())
}

// TestTeoOriginGroupV2Resource_UpdateApiError tests that an API error during
// ModifyOriginGroup is reported through diagnostics.
func TestTeoOriginGroupV2Resource_UpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	patches.ApplyMethodFunc(teoClient, "ModifyOriginGroupWithContext",
		func(_ context.Context, request *teov20220901.ModifyOriginGroupRequest) (*teov20220901.ModifyOriginGroupResponse, error) {
			return nil, fmt.Errorf("internal error")
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	planAttrs := originGroupV2PlanAttrs("zone-297z8rf93cfw")
	planAttrs["name"] = tftypes.NewValue(tftypes.String, "tf-test-origin-group-renamed")
	plan := originGroupV2RawAsPlan(t, planAttrs)

	req := resource.UpdateRequest{
		Plan:   plan,
		Config: tfsdk.Config{Raw: plan.Raw, Schema: plan.Schema},
		State:  originGroupV2ReadState(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema}}
	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestTeoOriginGroupV2Resource_Delete tests a successful delete.
func TestTeoOriginGroupV2Resource_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	var gotDelete *teov20220901.DeleteOriginGroupRequest
	patches.ApplyMethodFunc(teoClient, "DeleteOriginGroupWithContext",
		func(_ context.Context, request *teov20220901.DeleteOriginGroupRequest) (*teov20220901.DeleteOriginGroupResponse, error) {
			gotDelete = request
			resp := teov20220901.NewDeleteOriginGroupResponse()
			resp.Response = &teov20220901.DeleteOriginGroupResponseParams{
				RequestId: ptrStringOriginGroupV2("fake-request-id"),
			}
			return resp, nil
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	req := resource.DeleteRequest{
		State: originGroupV2ReadState(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
	}
	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: req.State.Raw, Schema: req.State.Schema}}
	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotDelete)
	assert.Equal(t, "zone-297z8rf93cfw", *gotDelete.ZoneId)
	assert.Equal(t, "origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a", *gotDelete.GroupId)
}

// TestTeoOriginGroupV2Resource_DeleteApiError tests that an API error during
// delete is reported and the resource stays in state.
func TestTeoOriginGroupV2Resource_DeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := patchTeoOriginGroupV2Client(patches, client)

	patches.ApplyMethodFunc(teoClient, "DeleteOriginGroupWithContext",
		func(_ context.Context, request *teov20220901.DeleteOriginGroupRequest) (*teov20220901.DeleteOriginGroupResponse, error) {
			return nil, fmt.Errorf("internal error")
		})

	r := originGroupV2Resource()
	configureOriginGroupV2Client(r.(*svcteo.TeoOriginGroupV2Resource), client)

	req := resource.DeleteRequest{
		State: originGroupV2ReadState(t, "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a"),
	}
	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: req.State.Raw, Schema: req.State.Schema}}
	r.Delete(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestTeoOriginGroupV2Resource_DeleteBrokenId tests that a broken id errors out.
func TestTeoOriginGroupV2Resource_DeleteBrokenId(t *testing.T) {
	r := originGroupV2Resource()

	req := resource.DeleteRequest{
		State: originGroupV2ReadState(t, "broken-id-without-separator"),
	}
	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: req.State.Raw, Schema: req.State.Schema}}
	r.Delete(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestTeoOriginGroupV2Resource_ImportStateGood tests importing a valid id.
func TestTeoOriginGroupV2Resource_ImportStateGood(t *testing.T) {
	r := svcteo.NewTeoOriginGroupV2Resource().(*svcteo.TeoOriginGroupV2Resource)

	resp := &resource.ImportStateResponse{State: tfsdk.State{Raw: tftypes.NewValue(
		tftypes.Object{AttributeTypes: originGroupV2AttrTypes()}, nil),
		Schema: originGroupV2ResourceSchema(t)}}
	r.ImportState(context.Background(), resource.ImportStateRequest{
		ID: "zone-297z8rf93cfw#origin-2ccgtb24-7dc5-46s2-9r3e-95825d53dwe3a",
	}, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
}

// TestTeoOriginGroupV2Resource_ImportStateBroken tests that importing a
// broken id produces an error diagnostic.
func TestTeoOriginGroupV2Resource_ImportStateBroken(t *testing.T) {
	r := svcteo.NewTeoOriginGroupV2Resource().(*svcteo.TeoOriginGroupV2Resource)

	resp := &resource.ImportStateResponse{}
	r.ImportState(context.Background(), resource.ImportStateRequest{
		ID: "broken-id-without-separator",
	}, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestTeoOriginGroupV2Resource_ImportStatePassthroughBadId verifies that the
// embedded fw.WithImportByID passthrough rejects a broken composite id.
func TestTeoOriginGroupV2Resource_ImportStatePassthroughBadId(t *testing.T) {
	r := svcteo.NewTeoOriginGroupV2Resource().(*svcteo.TeoOriginGroupV2Resource)

	resp := &resource.ImportStateResponse{}
	r.ImportState(context.Background(), resource.ImportStateRequest{
		ID: "broken-id-without-separator",
	}, resp)

	assert.True(t, resp.Diagnostics.HasError())
}
