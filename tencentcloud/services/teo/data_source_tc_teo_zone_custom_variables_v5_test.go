package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringZoneCustomVariables(s string) *string {
	return &s
}

// teoZoneCustomVariablesV5Schema returns the framework data source schema
// used to build the raw config value for ReadRequest.
func teoZoneCustomVariablesV5Schema(t *testing.T) datasource_schema.Schema {
	ds, ok := svcteo.NewTeoZoneCustomVariablesV5DataSource().(*svcteo.TeoZoneCustomVariablesV5DataSource)
	if !ok {
		t.Fatalf("factory did not return *TeoZoneCustomVariablesV5DataSource")
	}
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newTeoZoneCustomVariablesV5Raw builds a raw object value with the given
// zone_id and every other attribute null, matching the data source schema
// type exactly.
func newTeoZoneCustomVariablesV5Raw(s datasource_schema.Schema, zoneIdVal *string) tftypes.Value {
	objType, ok := s.Type().(types.ObjectType)
	if !ok {
		panic("schema type is not types.ObjectType")
	}
	attrs := objType.AttributeTypes()

	var zoneId tftypes.Value
	if zoneIdVal == nil {
		zoneId = tftypes.NewValue(tftypes.String, nil)
	} else {
		zoneId = tftypes.NewValue(tftypes.String, *zoneIdVal)
	}

	attrValues := make(map[string]tftypes.Value, len(attrs))
	for name, attrType := range attrs {
		attrValues[name] = tftypes.NewValue(attrType.TerraformType(context.Background()), nil)
	}
	attrValues["zone_id"] = zoneId

	return tftypes.NewValue(s.Type().TerraformType(context.Background()), attrValues)
}

// newTeoZoneCustomVariablesV5ReadRequest builds a datasource.ReadRequest
// carrying the given zone_id as a raw tftypes value.
func newTeoZoneCustomVariablesV5ReadRequest(t *testing.T, zoneId string) datasource.ReadRequest {
	s := teoZoneCustomVariablesV5Schema(t)

	req := datasource.ReadRequest{}
	req.Config = tfsdk.Config{
		Raw:    newTeoZoneCustomVariablesV5Raw(s, &zoneId),
		Schema: s,
	}
	return req
}

// newTeoZoneCustomVariablesV5ReadRequestNullZoneId builds a ReadRequest whose
// zone_id is null.
func newTeoZoneCustomVariablesV5ReadRequestNullZoneId(t *testing.T) datasource.ReadRequest {
	s := teoZoneCustomVariablesV5Schema(t)

	req := datasource.ReadRequest{}
	req.Config = tfsdk.Config{
		Raw:    newTeoZoneCustomVariablesV5Raw(s, nil),
		Schema: s,
	}
	return req
}

// newTeoZoneCustomVariablesV5ReadResponse builds an empty ReadResponse with a
// pre-populated null state matching the data source schema.
func newTeoZoneCustomVariablesV5ReadResponse(t *testing.T) *datasource.ReadResponse {
	s := teoZoneCustomVariablesV5Schema(t)

	objType, ok := s.Type().(types.ObjectType)
	if !ok {
		t.Fatalf("schema type is not types.ObjectType: %T", s.Type())
	}
	attrs := objType.AttributeTypes()
	stateAttrs := make(map[string]tftypes.Value, len(attrs))
	for name, attrType := range attrs {
		stateAttrs[name] = tftypes.NewValue(attrType.TerraformType(context.Background()), nil)
	}

	stateRaw := tftypes.NewValue(s.Type().TerraformType(context.Background()), stateAttrs)

	return &datasource.ReadResponse{
		State: tfsdk.State{
			Raw:    stateRaw,
			Schema: s,
		},
	}
}

// setTeoZoneCustomVariablesV5Client injects a mock connectivity client into
// the data source instance via the Configure method.
func setTeoZoneCustomVariablesV5Client(d *svcteo.TeoZoneCustomVariablesV5DataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoZoneCustomVariablesV5" -v -count=1 -gcflags="all=-l"

// TestTeoZoneCustomVariablesV5_Read_Success tests a successful Read with
// one custom variable and one operation rule covering all nesting levels.
func TestTeoZoneCustomVariablesV5_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotRequest *teov20220901.DescribeZoneCustomVariablesRequest
	patches.ApplyMethodFunc(teoClient, "DescribeZoneCustomVariables", func(request *teov20220901.DescribeZoneCustomVariablesRequest) (*teov20220901.DescribeZoneCustomVariablesResponse, error) {
		gotRequest = request
		resp := teov20220901.NewDescribeZoneCustomVariablesResponse()
		resp.Response = &teov20220901.DescribeZoneCustomVariablesResponseParams{
			CustomVariables: []*teov20220901.CustomVariable{
				{
					Name:         ptrStringZoneCustomVariables("user.zone.var_a"),
					InitialValue: ptrStringZoneCustomVariables("init-a"),
					Description:  ptrStringZoneCustomVariables("variable a"),
				},
			},
			CustomVariableOperations: []*teov20220901.CustomVariableOperation{
				{
					Branches: []*teov20220901.CustomVariableOperationRuleBranch{
						{
							Condition: ptrStringZoneCustomVariables("http.request.uri extend \"/path\""),
							Actions: []*teov20220901.CustomVariableOperationRuleAction{
								{
									Name: ptrStringZoneCustomVariables("Set"),
									SetParameters: &teov20220901.SetParameters{
										Name:  ptrStringZoneCustomVariables("user.zone.var_a"),
										Value: ptrStringZoneCustomVariables("value-a"),
									},
								},
							},
							SubRules: []*teov20220901.CustomVariableOperationSubRule{
								{
									Branches: []*teov20220901.CustomVariableOperationRuleBranch{
										{
											Condition: ptrStringZoneCustomVariables("http.request.host equal \"example.com\""),
											Actions: []*teov20220901.CustomVariableOperationRuleAction{
												{
													Name: ptrStringZoneCustomVariables("Set"),
													SetParameters: &teov20220901.SetParameters{
														Name:  ptrStringZoneCustomVariables("user.zone.var_a"),
														Value: ptrStringZoneCustomVariables("value-l2"),
													},
												},
											},
											SubRules: []*teov20220901.CustomVariableOperationSubRule{
												{
													Branches: []*teov20220901.CustomVariableOperationRuleBranch{
														{
															Condition: ptrStringZoneCustomVariables("http.request.uri extend \"/deep\""),
															Actions: []*teov20220901.CustomVariableOperationRuleAction{
																{
																	Name: ptrStringZoneCustomVariables("Set"),
																	SetParameters: &teov20220901.SetParameters{
																		Name:  ptrStringZoneCustomVariables("user.zone.var_a"),
																		Value: ptrStringZoneCustomVariables("value-leaf"),
																	},
																},
															},
														},
													},
													Description: []*string{
														ptrStringZoneCustomVariables("l2 sub rule remark"),
													},
												},
											},
										},
									},
									Description: []*string{
										ptrStringZoneCustomVariables("l1 sub rule remark"),
									},
								},
							},
						},
					},
					Description: []*string{
						ptrStringZoneCustomVariables("operation remark"),
					},
				},
			},
			RequestId: ptrStringZoneCustomVariables("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoZoneCustomVariablesV5DataSource()
	dataSource, ok := ds.(*svcteo.TeoZoneCustomVariablesV5DataSource)
	if !ok {
		t.Fatalf("factory did not return *TeoZoneCustomVariablesV5DataSource")
	}
	setTeoZoneCustomVariablesV5Client(dataSource, client)

	req := newTeoZoneCustomVariablesV5ReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoZoneCustomVariablesV5ReadResponse(t)
	dataSource.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, "zone-2qtuhspy7cr6", *gotRequest.ZoneId)

	var state svcteo.TeoZoneCustomVariablesV5DataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "expected no state read errors, got: %v", diags)

	// zone_id is echoed back into state
	assert.Equal(t, "zone-2qtuhspy7cr6", state.ZoneId.ValueString())

	// custom_variables
	assert.False(t, state.CustomVariables.IsNull())
	var customVariables []svcteo.ZoneCustomVariableModelForTest
	diags = state.CustomVariables.ElementsAs(context.Background(), &customVariables, false)
	assert.False(t, diags.HasError())
	if assert.Len(t, customVariables, 1) {
		assert.Equal(t, "user.zone.var_a", customVariables[0].Name.ValueString())
		assert.Equal(t, "init-a", customVariables[0].InitialValue.ValueString())
		assert.Equal(t, "variable a", customVariables[0].Description.ValueString())
	}

	// custom_variable_operations
	assert.False(t, state.CustomVariableOperations.IsNull())
	var operations []svcteo.CustomVariableOperationModelForTest
	diags = state.CustomVariableOperations.ElementsAs(context.Background(), &operations, false)
	assert.False(t, diags.HasError())
	if !assert.Len(t, operations, 1) {
		return
	}

	// operation description
	var opDescriptions []types.String
	diags = operations[0].Description.ElementsAs(context.Background(), &opDescriptions, false)
	assert.False(t, diags.HasError())
	if assert.Len(t, opDescriptions, 1) {
		assert.Equal(t, "operation remark", opDescriptions[0].ValueString())
	}

	// level-1 branches
	var branches []svcteo.CustomVariableOperationBranchModelForTest
	diags = operations[0].Branches.ElementsAs(context.Background(), &branches, false)
	assert.False(t, diags.HasError())
	if !assert.Len(t, branches, 1) {
		return
	}
	assert.Equal(t, "http.request.uri extend \"/path\"", branches[0].Condition.ValueString())

	// level-1 actions
	var actions []svcteo.CustomVariableOperationActionModelForTest
	diags = branches[0].Actions.ElementsAs(context.Background(), &actions, false)
	assert.False(t, diags.HasError())
	if assert.Len(t, actions, 1) {
		assert.Equal(t, "Set", actions[0].Name.ValueString())
		assert.False(t, actions[0].SetParameters.IsNull())
		assert.Equal(t, "user.zone.var_a", actions[0].SetParameters.Attributes()["name"].(types.String).ValueString())
		assert.Equal(t, "value-a", actions[0].SetParameters.Attributes()["value"].(types.String).ValueString())
	}

	// level-1 sub_rules
	var subRules []svcteo.CustomVariableOperationSubRuleModelForTest
	diags = branches[0].SubRules.ElementsAs(context.Background(), &subRules, false)
	assert.False(t, diags.HasError())
	if !assert.Len(t, subRules, 1) {
		return
	}
	var subRuleDescriptions []types.String
	diags = subRules[0].Description.ElementsAs(context.Background(), &subRuleDescriptions, false)
	assert.False(t, diags.HasError())
	if assert.Len(t, subRuleDescriptions, 1) {
		assert.Equal(t, "l1 sub rule remark", subRuleDescriptions[0].ValueString())
	}

	// level-2 branches
	var branchesL2 []svcteo.CustomVariableOperationBranchModelForTest
	diags = subRules[0].Branches.ElementsAs(context.Background(), &branchesL2, false)
	assert.False(t, diags.HasError())
	if !assert.Len(t, branchesL2, 1) {
		return
	}
	assert.Equal(t, "http.request.host equal \"example.com\"", branchesL2[0].Condition.ValueString())

	// level-2 actions
	var actionsL2 []svcteo.CustomVariableOperationActionModelForTest
	diags = branchesL2[0].Actions.ElementsAs(context.Background(), &actionsL2, false)
	assert.False(t, diags.HasError())
	if assert.Len(t, actionsL2, 1) {
		assert.Equal(t, "Set", actionsL2[0].Name.ValueString())
		assert.Equal(t, "value-l2", actionsL2[0].SetParameters.Attributes()["value"].(types.String).ValueString())
	}

	// level-2 sub_rules
	var subRulesL2 []svcteo.CustomVariableOperationSubRuleModelForTest
	diags = branchesL2[0].SubRules.ElementsAs(context.Background(), &subRulesL2, false)
	assert.False(t, diags.HasError())
	if !assert.Len(t, subRulesL2, 1) {
		return
	}
	var subRuleL2Descriptions []types.String
	diags = subRulesL2[0].Description.ElementsAs(context.Background(), &subRuleL2Descriptions, false)
	assert.False(t, diags.HasError())
	if assert.Len(t, subRuleL2Descriptions, 1) {
		assert.Equal(t, "l2 sub rule remark", subRuleL2Descriptions[0].ValueString())
	}

	// leaf branches (truncated recursion level)
	var branchLeaves []svcteo.CustomVariableOperationBranchLeafModelForTest
	diags = subRulesL2[0].Branches.ElementsAs(context.Background(), &branchLeaves, false)
	assert.False(t, diags.HasError())
	if assert.Len(t, branchLeaves, 1) {
		assert.Equal(t, "http.request.uri extend \"/deep\"", branchLeaves[0].Condition.ValueString())
		var leafActions []svcteo.CustomVariableOperationActionModelForTest
		diags = branchLeaves[0].Actions.ElementsAs(context.Background(), &leafActions, false)
		assert.False(t, diags.HasError())
		if assert.Len(t, leafActions, 1) {
			assert.Equal(t, "Set", leafActions[0].Name.ValueString())
			assert.Equal(t, "value-leaf", leafActions[0].SetParameters.Attributes()["value"].(types.String).ValueString())
		}
	}
}

// TestTeoZoneCustomVariablesV5_Read_EmptyResult tests that an empty result
// (no custom variables and no operations configured) is a normal outcome.
func TestTeoZoneCustomVariablesV5_Read_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZoneCustomVariables", func(request *teov20220901.DescribeZoneCustomVariablesRequest) (*teov20220901.DescribeZoneCustomVariablesResponse, error) {
		resp := teov20220901.NewDescribeZoneCustomVariablesResponse()
		resp.Response = &teov20220901.DescribeZoneCustomVariablesResponseParams{
			CustomVariables:          []*teov20220901.CustomVariable{},
			CustomVariableOperations: []*teov20220901.CustomVariableOperation{},
			RequestId:                ptrStringZoneCustomVariables("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoZoneCustomVariablesV5DataSource()
	dataSource, ok := ds.(*svcteo.TeoZoneCustomVariablesV5DataSource)
	if !ok {
		t.Fatalf("factory did not return *TeoZoneCustomVariablesV5DataSource")
	}
	setTeoZoneCustomVariablesV5Client(dataSource, client)

	req := newTeoZoneCustomVariablesV5ReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoZoneCustomVariablesV5ReadResponse(t)
	dataSource.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var state svcteo.TeoZoneCustomVariablesV5DataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError())

	assert.False(t, state.CustomVariables.IsNull())
	assert.Equal(t, 0, len(state.CustomVariables.Elements()))
	assert.False(t, state.CustomVariableOperations.IsNull())
	assert.Equal(t, 0, len(state.CustomVariableOperations.Elements()))
}

// TestTeoZoneCustomVariablesV5_Read_NilLists tests that nil response lists
// (the cloud API "may return null" semantics) are preserved as null lists.
func TestTeoZoneCustomVariablesV5_Read_NilLists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZoneCustomVariables", func(request *teov20220901.DescribeZoneCustomVariablesRequest) (*teov20220901.DescribeZoneCustomVariablesResponse, error) {
		resp := teov20220901.NewDescribeZoneCustomVariablesResponse()
		resp.Response = &teov20220901.DescribeZoneCustomVariablesResponseParams{
			RequestId: ptrStringZoneCustomVariables("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoZoneCustomVariablesV5DataSource()
	dataSource, ok := ds.(*svcteo.TeoZoneCustomVariablesV5DataSource)
	if !ok {
		t.Fatalf("factory did not return *TeoZoneCustomVariablesV5DataSource")
	}
	setTeoZoneCustomVariablesV5Client(dataSource, client)

	req := newTeoZoneCustomVariablesV5ReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoZoneCustomVariablesV5ReadResponse(t)
	dataSource.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var state svcteo.TeoZoneCustomVariablesV5DataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError())

	// nil slices are surfaced as empty lists (never an error)
	assert.False(t, state.CustomVariables.IsNull())
	assert.Equal(t, 0, len(state.CustomVariables.Elements()))
	assert.False(t, state.CustomVariableOperations.IsNull())
	assert.Equal(t, 0, len(state.CustomVariableOperations.Elements()))
}

// TestTeoZoneCustomVariablesV5_Read_APIError tests API error handling.
func TestTeoZoneCustomVariablesV5_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZoneCustomVariables", func(request *teov20220901.DescribeZoneCustomVariablesRequest) (*teov20220901.DescribeZoneCustomVariablesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=zone not found")
	})

	ds := svcteo.NewTeoZoneCustomVariablesV5DataSource()
	dataSource, ok := ds.(*svcteo.TeoZoneCustomVariablesV5DataSource)
	if !ok {
		t.Fatalf("factory did not return *TeoZoneCustomVariablesV5DataSource")
	}
	setTeoZoneCustomVariablesV5Client(dataSource, client)

	req := newTeoZoneCustomVariablesV5ReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoZoneCustomVariablesV5ReadResponse(t)
	dataSource.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "querying TEO zone custom variables")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ResourceNotFound")
}

// TestTeoZoneCustomVariablesV5_Read_MissingZoneId tests that a missing
// zone_id is rejected before any API call.
func TestTeoZoneCustomVariablesV5_Read_MissingZoneId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	apiCalled := false
	patches.ApplyMethodFunc(teoClient, "DescribeZoneCustomVariables", func(request *teov20220901.DescribeZoneCustomVariablesRequest) (*teov20220901.DescribeZoneCustomVariablesResponse, error) {
		apiCalled = true
		return teov20220901.NewDescribeZoneCustomVariablesResponse(), nil
	})

	ds := svcteo.NewTeoZoneCustomVariablesV5DataSource()
	dataSource, ok := ds.(*svcteo.TeoZoneCustomVariablesV5DataSource)
	if !ok {
		t.Fatalf("factory did not return *TeoZoneCustomVariablesV5DataSource")
	}
	setTeoZoneCustomVariablesV5Client(dataSource, client)

	req := newTeoZoneCustomVariablesV5ReadRequestNullZoneId(t)
	resp := newTeoZoneCustomVariablesV5ReadResponse(t)
	dataSource.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing zone_id")
	assert.False(t, apiCalled)
}

// TestTeoZoneCustomVariablesV5_Read_ClientNotConfigured tests the nil-client
// guard.
func TestTeoZoneCustomVariablesV5_Read_ClientNotConfigured(t *testing.T) {
	ds := svcteo.NewTeoZoneCustomVariablesV5DataSource()
	dataSource, ok := ds.(*svcteo.TeoZoneCustomVariablesV5DataSource)
	if !ok {
		t.Fatalf("factory did not return *TeoZoneCustomVariablesV5DataSource")
	}

	req := newTeoZoneCustomVariablesV5ReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoZoneCustomVariablesV5ReadResponse(t)
	dataSource.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestTeoZoneCustomVariablesV5_MetadataAndSchema validates metadata and
// schema definition.
func TestTeoZoneCustomVariablesV5_MetadataAndSchema(t *testing.T) {
	ds := svcteo.NewTeoZoneCustomVariablesV5DataSource()
	dataSource, ok := ds.(*svcteo.TeoZoneCustomVariablesV5DataSource)
	if !ok {
		t.Fatalf("factory did not return *TeoZoneCustomVariablesV5DataSource")
	}

	metaResp := &datasource.MetadataResponse{}
	dataSource.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_zone_custom_variables_v5", metaResp.TypeName)

	schemaResp := &datasource.SchemaResponse{}
	dataSource.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "zone_id")
	assert.Contains(t, attrs, "custom_variables")
	assert.Contains(t, attrs, "custom_variable_operations")

	zoneIdAttr, ok := attrs["zone_id"].(datasource_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, zoneIdAttr.Required)

	customVariablesAttr, ok := attrs["custom_variables"].(datasource_schema.ListNestedAttribute)
	assert.True(t, ok)
	assert.True(t, customVariablesAttr.Computed)

	operationsAttr, ok := attrs["custom_variable_operations"].(datasource_schema.ListNestedAttribute)
	assert.True(t, ok)
	assert.True(t, operationsAttr.Computed)
}
