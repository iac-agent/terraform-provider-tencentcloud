package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoZoneCustomVariablesDataSource" -v -count=1 -gcflags="all=-l"

func ptrZoneCustomVariablesString(s string) *string {
	return &s
}

// zoneCustomVariablesDataSourceSchema returns the framework data source schema
// used to build the raw config value for ReadRequest.
func zoneCustomVariablesDataSourceSchema(t *testing.T) schema.Schema {
	d := svcteo.NewTeoZoneCustomVariablesDataSource()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newZoneCustomVariablesReadRequest builds a datasource.ReadRequest carrying
// the given zone_id as raw tftypes value. Computed attributes are set to null.
func newZoneCustomVariablesReadRequest(t *testing.T, zoneId string) (datasource.ReadRequest, *datasource.ReadResponse) {
	s := zoneCustomVariablesDataSourceSchema(t)

	objType, ok := s.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is not tftypes.Object, got: %T", s.Type().TerraformType(context.Background()))
	}

	raw := tftypes.NewValue(objType, map[string]tftypes.Value{
		"zone_id":                    tftypes.NewValue(tftypes.String, zoneId),
		"custom_variables":           tftypes.NewValue(objType.AttributeTypes["custom_variables"], nil),
		"custom_variable_operations": tftypes.NewValue(objType.AttributeTypes["custom_variable_operations"], nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s

	resp := &datasource.ReadResponse{}
	resp.State.Schema = s
	return req, resp
}

// setZoneCustomVariablesClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setZoneCustomVariablesClient(d datasource.DataSourceWithConfigure, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newZoneCustomVariablesResponse builds a full DescribeZoneCustomVariables
// response covering every nested level of the schema.
func newZoneCustomVariablesResponse() *teov20220901.DescribeZoneCustomVariablesResponse {
	return &teov20220901.DescribeZoneCustomVariablesResponse{
		Response: &teov20220901.DescribeZoneCustomVariablesResponseParams{
			CustomVariables: []*teov20220901.CustomVariable{
				{
					Name:         ptrZoneCustomVariablesString("user.zone.example_var"),
					InitialValue: ptrZoneCustomVariablesString("initial"),
					Description:  ptrZoneCustomVariablesString("example variable"),
				},
			},
			CustomVariableOperations: []*teov20220901.CustomVariableOperation{
				{
					Branches: []*teov20220901.CustomVariableOperationRuleBranch{
						{
							Condition: ptrZoneCustomVariablesString("http.request.host eq \"example.com\""),
							Actions: []*teov20220901.CustomVariableOperationRuleAction{
								{
									Name: ptrZoneCustomVariablesString("Set"),
									SetParameters: &teov20220901.SetParameters{
										Name:  ptrZoneCustomVariablesString("user.zone.example_var"),
										Value: ptrZoneCustomVariablesString("computed-value"),
									},
								},
							},
							SubRules: []*teov20220901.CustomVariableOperationSubRule{
								{
									Branches: []*teov20220901.CustomVariableOperationRuleBranch{
										{
											Condition: ptrZoneCustomVariablesString("http.request.uri eq \"/path\""),
											Actions: []*teov20220901.CustomVariableOperationRuleAction{
												{
													Name: ptrZoneCustomVariablesString("Set"),
													SetParameters: &teov20220901.SetParameters{
														Name:  ptrZoneCustomVariablesString("user.zone.example_var"),
														Value: ptrZoneCustomVariablesString("deeper-value"),
													},
												},
											},
											SubRules: []*teov20220901.CustomVariableOperationSubRule{
												{
													Description: []*string{
														ptrZoneCustomVariablesString("deep sub rule comment"),
													},
												},
											},
										},
									},
									Description: []*string{
										ptrZoneCustomVariablesString("sub rule comment"),
									},
								},
							},
						},
					},
					Description: []*string{
						ptrZoneCustomVariablesString("operation comment"),
					},
				},
			},
			RequestId: ptrZoneCustomVariablesString("fake-request-id"),
		},
	}
}

// TestTeoZoneCustomVariablesDataSource_Read_Success verifies Read converts a
// full API response into state, covering every nested level of the schema.
func TestTeoZoneCustomVariablesDataSource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotRequest *teov20220901.DescribeZoneCustomVariablesRequest
	patches.ApplyMethodFunc(teoClient, "DescribeZoneCustomVariablesWithContext", func(_ context.Context, request *teov20220901.DescribeZoneCustomVariablesRequest) (*teov20220901.DescribeZoneCustomVariablesResponse, error) {
		gotRequest = request
		return newZoneCustomVariablesResponse(), nil
	})

	d := svcteo.NewTeoZoneCustomVariablesDataSource()
	setZoneCustomVariablesClient(d.(datasource.DataSourceWithConfigure), client)

	req, resp := newZoneCustomVariablesReadRequest(t, "zone-29371sfdvge2")
	d.Read(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}
	if gotRequest == nil {
		t.Fatal("DescribeZoneCustomVariablesWithContext should have been called")
	}
	if gotRequest.ZoneId == nil || *gotRequest.ZoneId != "zone-29371sfdvge2" {
		t.Fatalf("expected ZoneId zone-29371sfdvge2, got: %v", gotRequest.ZoneId)
	}

	var state svcteo.TeoZoneCustomVariablesDataSourceModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}

	if state.ZoneId.ValueString() != "zone-29371sfdvge2" {
		t.Errorf("state zone_id = %q, want %q", state.ZoneId.ValueString(), "zone-29371sfdvge2")
	}

	// custom_variables
	if state.CustomVariables.IsNull() {
		t.Fatal("custom_variables should not be null")
	}
	var variables []struct {
		Name         string `tfsdk:"name"`
		InitialValue string `tfsdk:"initial_value"`
		Description  string `tfsdk:"description"`
	}
	if diags := state.CustomVariables.ElementsAs(context.Background(), &variables, false); diags.HasError() {
		t.Fatalf("failed to convert custom_variables: %v", diags)
	}
	if len(variables) != 1 {
		t.Fatalf("expected 1 custom variable, got %d", len(variables))
	}
	if variables[0].Name != "user.zone.example_var" {
		t.Errorf("custom variable name = %q, want %q", variables[0].Name, "user.zone.example_var")
	}
	if variables[0].InitialValue != "initial" {
		t.Errorf("custom variable initial_value = %q, want %q", variables[0].InitialValue, "initial")
	}
	if variables[0].Description != "example variable" {
		t.Errorf("custom variable description = %q, want %q", variables[0].Description, "example variable")
	}

	// custom_variable_operations
	if state.CustomVariableOperation.IsNull() {
		t.Fatal("custom_variable_operations should not be null")
	}
	var operations []struct {
		Branches []struct {
			Condition string `tfsdk:"condition"`
			Actions   []struct {
				Name string `tfsdk:"name"`
			} `tfsdk:"actions"`
			SubRules []struct {
				Branches []struct {
					Condition string `tfsdk:"condition"`
					Actions   []struct {
						Name string `tfsdk:"name"`
					} `tfsdk:"actions"`
					SubRules []struct {
						Description []string `tfsdk:"description"`
					} `tfsdk:"sub_rules"`
				} `tfsdk:"branches"`
				Description []string `tfsdk:"description"`
			} `tfsdk:"sub_rules"`
		} `tfsdk:"branches"`
		Description []string `tfsdk:"description"`
	}
	if diags := state.CustomVariableOperation.ElementsAs(context.Background(), &operations, false); diags.HasError() {
		t.Fatalf("failed to convert custom_variable_operations: %v", diags)
	}
	if len(operations) != 1 {
		t.Fatalf("expected 1 custom variable operation, got %d", len(operations))
	}
	op := operations[0]
	if len(op.Description) != 1 || op.Description[0] != "operation comment" {
		t.Errorf("operation description = %v, want [operation comment]", op.Description)
	}
	if len(op.Branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(op.Branches))
	}
	branch := op.Branches[0]
	if branch.Condition != "http.request.host eq \"example.com\"" {
		t.Errorf("branch condition = %q, want %q", branch.Condition, "http.request.host eq \"example.com\"")
	}
	if len(branch.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(branch.Actions))
	}
	if branch.Actions[0].Name != "Set" {
		t.Errorf("action name = %q, want %q", branch.Actions[0].Name, "Set")
	}
	if len(branch.SubRules) != 1 {
		t.Fatalf("expected 1 sub rule, got %d", len(branch.SubRules))
	}
	subRule := branch.SubRules[0]
	if len(subRule.Description) != 1 || subRule.Description[0] != "sub rule comment" {
		t.Errorf("sub rule description = %v, want [sub rule comment]", subRule.Description)
	}
	if len(subRule.Branches) != 1 {
		t.Fatalf("expected 1 sub rule branch, got %d", len(subRule.Branches))
	}
	subBranch := subRule.Branches[0]
	if subBranch.Condition != "http.request.uri eq \"/path\"" {
		t.Errorf("sub branch condition = %q, want %q", subBranch.Condition, "http.request.uri eq \"/path\"")
	}
	if len(subBranch.Actions) != 1 || subBranch.Actions[0].Name != "Set" {
		t.Errorf("sub branch actions = %v, want one Set action", subBranch.Actions)
	}
	if len(subBranch.SubRules) != 1 {
		t.Fatalf("expected 1 deep sub rule, got %d", len(subBranch.SubRules))
	}
	if len(subBranch.SubRules[0].Description) != 1 || subBranch.SubRules[0].Description[0] != "deep sub rule comment" {
		t.Errorf("deep sub rule description = %v, want [deep sub rule comment]", subBranch.SubRules[0].Description)
	}
}

// TestTeoZoneCustomVariablesDataSource_Read_EmptyResult verifies an empty
// configuration (both lists empty) is a valid result, not an error.
func TestTeoZoneCustomVariablesDataSource_Read_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZoneCustomVariablesWithContext", func(_ context.Context, request *teov20220901.DescribeZoneCustomVariablesRequest) (*teov20220901.DescribeZoneCustomVariablesResponse, error) {
		resp := teov20220901.NewDescribeZoneCustomVariablesResponse()
		resp.Response = &teov20220901.DescribeZoneCustomVariablesResponseParams{
			CustomVariables:          []*teov20220901.CustomVariable{},
			CustomVariableOperations: []*teov20220901.CustomVariableOperation{},
			RequestId:                ptrZoneCustomVariablesString("fake-request-id"),
		}
		return resp, nil
	})

	d := svcteo.NewTeoZoneCustomVariablesDataSource()
	setZoneCustomVariablesClient(d.(datasource.DataSourceWithConfigure), client)

	req, resp := newZoneCustomVariablesReadRequest(t, "zone-29371sfdvge2")
	d.Read(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}

	var state svcteo.TeoZoneCustomVariablesDataSourceModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}

	if state.CustomVariables.IsNull() {
		t.Error("custom_variables should be an empty list, not null")
	} else if len(state.CustomVariables.Elements()) != 0 {
		t.Errorf("expected empty custom_variables, got %d elements", len(state.CustomVariables.Elements()))
	}
	if state.CustomVariableOperation.IsNull() {
		t.Error("custom_variable_operations should be an empty list, not null")
	} else if len(state.CustomVariableOperation.Elements()) != 0 {
		t.Errorf("expected empty custom_variable_operations, got %d elements", len(state.CustomVariableOperation.Elements()))
	}
}

// TestTeoZoneCustomVariablesDataSource_Read_APIError verifies API errors are
// surfaced as error diagnostics.
func TestTeoZoneCustomVariablesDataSource_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZoneCustomVariablesWithContext", func(_ context.Context, request *teov20220901.DescribeZoneCustomVariablesRequest) (*teov20220901.DescribeZoneCustomVariablesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=UnauthorizedOperation.Unknown, Message=unauthorized")
	})

	d := svcteo.NewTeoZoneCustomVariablesDataSource()
	setZoneCustomVariablesClient(d.(datasource.DataSourceWithConfigure), client)

	req, resp := newZoneCustomVariablesReadRequest(t, "zone-invalid")
	d.Read(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics for API error, got none")
	}
}

// TestTeoZoneCustomVariablesDataSource_Metadata validates Metadata.TypeName.
func TestTeoZoneCustomVariablesDataSource_Metadata(t *testing.T) {
	d := svcteo.NewTeoZoneCustomVariablesDataSource()
	req := datasource.MetadataRequest{}
	resp := &datasource.MetadataResponse{}

	d.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_teo_zone_custom_variables_1109_01" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_zone_custom_variables_1109_01")
	}
}

// TestTeoZoneCustomVariablesDataSource_Schema validates the schema defines
// the required input attribute and computed output attributes.
func TestTeoZoneCustomVariablesDataSource_Schema(t *testing.T) {
	d := svcteo.NewTeoZoneCustomVariablesDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	attrKeys := make(map[string]bool)
	for k := range resp.Schema.Attributes {
		attrKeys[k] = true
	}

	for _, required := range []string{"zone_id"} {
		if !attrKeys[required] {
			t.Errorf("schema missing required attribute: %s", required)
			continue
		}
		attr, ok := resp.Schema.Attributes[required].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %s is not a StringAttribute", required)
			continue
		}
		if !attr.Required {
			t.Errorf("attribute %s should be Required", required)
		}
	}

	for _, computed := range []string{"custom_variables", "custom_variable_operations"} {
		if !attrKeys[computed] {
			t.Errorf("schema missing computed attribute: %s", computed)
			continue
		}
		attr, ok := resp.Schema.Attributes[computed].(schema.ListNestedAttribute)
		if !ok {
			t.Errorf("attribute %s is not a ListNestedAttribute", computed)
			continue
		}
		if !attr.Computed {
			t.Errorf("attribute %s should be Computed", computed)
		}
	}
}
