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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringTeoDeploymentRecords(s string) *string { return &s }
func ptrInt64TeoDeploymentRecords(i int64) *int64    { return &i }
func ptrFloat64TeoDeploymentRecords(f float64) *float64 {
	return &f
}

// Test-local structs mirroring the data source schema for asserting every
// leaf field read back from state.
type teoDeploymentRecordItemTestModel struct {
	RecordId               types.String `tfsdk:"record_id"`
	Operation              types.String `tfsdk:"operation"`
	Status                 types.String `tfsdk:"status"`
	Duration               types.Int64  `tfsdk:"duration"`
	InferenceServiceConfig types.Object `tfsdk:"inference_service_config"`
	CreateTime             types.String `tfsdk:"create_time"`
	ActiveStatus           types.String `tfsdk:"active_status"`
}

type teoInferenceServiceConfigTestModel struct {
	ListenPort     types.Int64  `tfsdk:"listen_port"`
	RequestPaths   types.Set    `tfsdk:"request_paths"`
	Containers     types.List   `tfsdk:"containers"`
	ResourceConfig types.Object `tfsdk:"resource_config"`
	AffinityConfig types.Object `tfsdk:"affinity_config"`
}

type teoContainerItemTestModel struct {
	ImageType            types.String `tfsdk:"image_type"`
	TcrRepositoryConfig  types.Object `tfsdk:"tcr_repository_config"`
	StartupCommand       types.String `tfsdk:"startup_command"`
	EnvironmentVariables types.List   `tfsdk:"environment_variables"`
}

type teoTCRRepositoryConfigTestModel struct {
	TCRType    types.String `tfsdk:"tcr_type"`
	Image      types.String `tfsdk:"image"`
	RegistryId types.String `tfsdk:"registry_id"`
	RegionName types.String `tfsdk:"region_name"`
}

type teoEnvironmentVariableTestModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type teoResourceConfigTestModel struct {
	ScalingMode          types.String `tfsdk:"scaling_mode"`
	HardwareSpec         types.String `tfsdk:"hardware_spec"`
	HardwareSpecId       types.String `tfsdk:"hardware_spec_id"`
	HardwareConfig       types.Object `tfsdk:"hardware_config"`
	AutoScalingConfig    types.Object `tfsdk:"auto_scaling_config"`
	ManualInstanceConfig types.Object `tfsdk:"manual_instance_config"`
	Concurrency          types.Int64  `tfsdk:"concurrency"`
}

type teoHardwareConfigTestModel struct {
	GPUNum   types.Float64 `tfsdk:"gpu_num"`
	CPUNum   types.Float64 `tfsdk:"cpu_num"`
	MemSize  types.Int64   `tfsdk:"mem_size"`
	DiskSize types.Int64   `tfsdk:"disk_size"`
}

type teoAutoScalingConfigTestModel struct {
	MinInstanceCount types.Int64 `tfsdk:"min_instance_count"`
	ScalingPolicies  types.List  `tfsdk:"scaling_policies"`
}

type teoScalingPolicyTestModel struct {
	PolicyName             types.String `tfsdk:"policy_name"`
	PolicyType             types.String `tfsdk:"policy_type"`
	ScheduledScalingPolicy types.Object `tfsdk:"scheduled_scaling_policy"`
}

type teoScheduledScalingPolicyTestModel struct {
	ScheduledActions types.List   `tfsdk:"scheduled_actions"`
	EffectiveRange   types.Object `tfsdk:"effective_range"`
	TimeZone         types.String `tfsdk:"time_zone"`
}

type teoScheduledScalingActionTestModel struct {
	CronExpression   types.String `tfsdk:"cron_expression"`
	MinInstanceCount types.Int64  `tfsdk:"min_instance_count"`
}

type teoScheduledScalingEffectiveRangeTestModel struct {
	EffectiveType types.String `tfsdk:"effective_type"`
	StartDate     types.String `tfsdk:"start_date"`
	EndDate       types.String `tfsdk:"end_date"`
}

type teoAffinityConfigTestModel struct {
	Switch                  types.String `tfsdk:"switch"`
	AffinityMode            types.String `tfsdk:"affinity_mode"`
	SessionIdAffinityConfig types.Object `tfsdk:"session_id_affinity_config"`
}

type teoSessionIdAffinityConfigTestModel struct {
	Source     types.String `tfsdk:"source"`
	HeaderName types.String `tfsdk:"header_name"`
}

// teoDeploymentRecordsSchema returns the framework data source schema used to
// build the raw config value for the Read request.
func teoDeploymentRecordsSchema(t *testing.T) datasource_schema.Schema {
	d := svcteo.NewTeoInferenceServiceDeploymentRecords111001DataSource()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newTeoDeploymentRecordsReadRequest builds a datasource.ReadRequest carrying
// the given query attributes as raw tftypes values.
func newTeoDeploymentRecordsReadRequest(t *testing.T, zoneId, serviceId, sortBy, sortOrder interface{}) datasource.ReadRequest {
	s := teoDeploymentRecordsSchema(t)

	strOrNull := func(v interface{}) tftypes.Value {
		if v == nil {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, v)
	}

	// Build the config from the schema-derived tftypes type so the raw value
	// matches the data source schema exactly. record_set is Computed-only and
	// therefore always null in the config.
	schemaType, ok := s.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is not tftypes.Object, got %T", s.Type().TerraformType(context.Background()))
	}

	raw := tftypes.NewValue(schemaType, map[string]tftypes.Value{
		"zone_id":    strOrNull(zoneId),
		"service_id": strOrNull(serviceId),
		"sort_by":    strOrNull(sortBy),
		"sort_order": strOrNull(sortOrder),
		"record_set": tftypes.NewValue(schemaType.AttributeTypes["record_set"], nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// newTeoDeploymentRecordsState builds an empty tfsdk.State sharing the data
// source schema so that resp.State.Set writes the read result.
func newTeoDeploymentRecordsState(t *testing.T) tfsdk.State {
	s := teoDeploymentRecordsSchema(t)
	tftypesSchemaType := s.Type().TerraformType(context.Background())
	return tfsdk.State{
		Raw:    tftypes.NewValue(tftypesSchemaType, nil),
		Schema: s,
	}
}

// setTeoDeploymentRecordsClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setTeoDeploymentRecordsClient(d datasource.DataSource, client *connectivity.TencentCloudClient) {
	ds, ok := d.(datasource.DataSourceWithConfigure)
	if !ok {
		panic("data source does not implement DataSourceWithConfigure")
	}
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceServiceDeploymentRecords111001" -v -count=1 -gcflags="all=-l"

// TestTeoInferenceServiceDeploymentRecords111001_Read_Success verifies that a
// successful API response is fully mapped into state, including every nested
// leaf field.
func TestTeoInferenceServiceDeploymentRecords111001_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotRequests []*teov20220901.DescribeInferenceServiceDeploymentRecordsRequest
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecords", func(request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		gotRequests = append(gotRequests, request)
		return buildTeoDeploymentRecordsResponse(), nil
	})

	d := svcteo.NewTeoInferenceServiceDeploymentRecords111001DataSource()
	setTeoDeploymentRecordsClient(d, client)

	req := newTeoDeploymentRecordsReadRequest(t, "zone-2qtuhspy7cr6", "isvc-2qtuhspy7cr6", "create-time", "desc")
	resp := &datasource.ReadResponse{State: newTeoDeploymentRecordsState(t)}
	d.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// Verify the request carried the query attributes.
	assert.Equal(t, 1, len(gotRequests))
	gotReq := gotRequests[0]
	assert.Equal(t, "zone-2qtuhspy7cr6", *gotReq.ZoneId)
	assert.Equal(t, "isvc-2qtuhspy7cr6", *gotReq.ServiceId)
	assert.Equal(t, "create-time", *gotReq.SortBy)
	assert.Equal(t, "desc", *gotReq.SortOrder)

	// Verify the record_set list was written into state.
	var state svcteo.TeoInferenceServiceDeploymentRecords111001DataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "failed to read state back: %v", diags)
	assert.False(t, state.RecordSet.IsNull())

	var items []teoDeploymentRecordItemTestModel
	itemsDiags := state.RecordSet.ElementsAs(context.Background(), &items, false)
	assert.False(t, itemsDiags.HasError(), "failed to unmarshal record_set: %v", itemsDiags)
	assert.Equal(t, 1, len(items))
	if len(items) != 1 {
		t.Fatal("expected exactly one record in record_set")
	}
	item := items[0]

	assert.Equal(t, "record-1111", item.RecordId.ValueString())
	assert.Equal(t, "create", item.Operation.ValueString())
	assert.Equal(t, "succeeded", item.Status.ValueString())
	assert.Equal(t, int64(60), item.Duration.ValueInt64())
	assert.Equal(t, "2025-01-01T00:00:00+08:00", item.CreateTime.ValueString())
	assert.Equal(t, "active", item.ActiveStatus.ValueString())

	// inference_service_config leaf fields
	assert.False(t, item.InferenceServiceConfig.IsNull())
	var cfg teoInferenceServiceConfigTestModel
	assert.False(t, item.InferenceServiceConfig.As(context.Background(), &cfg, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal object attribute")
	assert.Equal(t, int64(8080), cfg.ListenPort.ValueInt64())

	var requestPaths []types.String
	requestPathsDiags := cfg.RequestPaths.ElementsAs(context.Background(), &requestPaths, false)
	assert.False(t, requestPathsDiags.HasError(), "failed to unmarshal request_paths: %v", requestPathsDiags)
	assert.Equal(t, 1, len(requestPaths))
	assert.Equal(t, "/v1/chat", requestPaths[0].ValueString())

	var containers []teoContainerItemTestModel
	containersDiags := cfg.Containers.ElementsAs(context.Background(), &containers, false)
	assert.False(t, containersDiags.HasError(), "failed to unmarshal containers: %v", containersDiags)
	assert.Equal(t, 1, len(containers))
	container := containers[0]
	assert.Equal(t, "TCR", container.ImageType.ValueString())
	assert.Equal(t, "python serving.py", container.StartupCommand.ValueString())

	assert.False(t, container.TcrRepositoryConfig.IsNull())
	var tcr teoTCRRepositoryConfigTestModel
	assert.False(t, container.TcrRepositoryConfig.As(context.Background(), &tcr, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal object attribute")
	assert.Equal(t, "Enterprise", tcr.TCRType.ValueString())
	assert.Equal(t, "test.example.com/repo/image:latest", tcr.Image.ValueString())
	assert.Equal(t, "tcr-1111", tcr.RegistryId.ValueString())
	assert.Equal(t, "ap-guangzhou", tcr.RegionName.ValueString())

	var envs []teoEnvironmentVariableTestModel
	envsDiags := container.EnvironmentVariables.ElementsAs(context.Background(), &envs, false)
	assert.False(t, envsDiags.HasError(), "failed to unmarshal environment_variables: %v", envsDiags)
	assert.Equal(t, 1, len(envs))
	assert.Equal(t, "MODEL", envs[0].Key.ValueString())
	assert.Equal(t, "llama3", envs[0].Value.ValueString())

	// resource_config leaf fields
	assert.False(t, cfg.ResourceConfig.IsNull())
	var resourceConfig teoResourceConfigTestModel
	assert.False(t, cfg.ResourceConfig.As(context.Background(), &resourceConfig, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal object attribute")
	assert.Equal(t, "Auto", resourceConfig.ScalingMode.ValueString())
	assert.Equal(t, "spec-legacy", resourceConfig.HardwareSpec.ValueString())
	assert.Equal(t, "spec-1111", resourceConfig.HardwareSpecId.ValueString())
	assert.Equal(t, int64(1), resourceConfig.Concurrency.ValueInt64())

	assert.False(t, resourceConfig.HardwareConfig.IsNull())
	var hardwareConfig teoHardwareConfigTestModel
	assert.False(t, resourceConfig.HardwareConfig.As(context.Background(), &hardwareConfig, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal object attribute")
	assert.Equal(t, float64(1), hardwareConfig.GPUNum.ValueFloat64())
	assert.Equal(t, float64(8), hardwareConfig.CPUNum.ValueFloat64())
	assert.Equal(t, int64(32768), hardwareConfig.MemSize.ValueInt64())
	assert.Equal(t, int64(10240), hardwareConfig.DiskSize.ValueInt64())

	assert.False(t, resourceConfig.AutoScalingConfig.IsNull())
	var autoScalingConfig teoAutoScalingConfigTestModel
	assert.False(t, resourceConfig.AutoScalingConfig.As(context.Background(), &autoScalingConfig, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal object attribute")
	assert.Equal(t, int64(1), autoScalingConfig.MinInstanceCount.ValueInt64())

	var policies []teoScalingPolicyTestModel
	policiesDiags := autoScalingConfig.ScalingPolicies.ElementsAs(context.Background(), &policies, false)
	assert.False(t, policiesDiags.HasError(), "failed to unmarshal scaling_policies: %v", policiesDiags)
	assert.Equal(t, 1, len(policies))
	assert.Equal(t, "policy-1", policies[0].PolicyName.ValueString())
	assert.Equal(t, "ScheduledScaling", policies[0].PolicyType.ValueString())

	assert.False(t, policies[0].ScheduledScalingPolicy.IsNull())
	var scheduledPolicy teoScheduledScalingPolicyTestModel
	assert.False(t, policies[0].ScheduledScalingPolicy.As(context.Background(), &scheduledPolicy, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal scheduled_scaling_policy")
	assert.Equal(t, "UTC", scheduledPolicy.TimeZone.ValueString())

	var actions []teoScheduledScalingActionTestModel
	actionsDiags := scheduledPolicy.ScheduledActions.ElementsAs(context.Background(), &actions, false)
	assert.False(t, actionsDiags.HasError(), "failed to unmarshal scheduled_actions: %v", actionsDiags)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "0 0 * * *", actions[0].CronExpression.ValueString())
	assert.Equal(t, int64(2), actions[0].MinInstanceCount.ValueInt64())

	assert.False(t, scheduledPolicy.EffectiveRange.IsNull())
	var effectiveRange teoScheduledScalingEffectiveRangeTestModel
	assert.False(t, scheduledPolicy.EffectiveRange.As(context.Background(), &effectiveRange, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal object attribute")
	assert.Equal(t, "LongTerm", effectiveRange.EffectiveType.ValueString())
	assert.Equal(t, "2025-01-01", effectiveRange.StartDate.ValueString())
	assert.Equal(t, "2099-12-31", effectiveRange.EndDate.ValueString())

	// manual_instance_config is nil in the API response, so it must be null.
	assert.True(t, resourceConfig.ManualInstanceConfig.IsNull())

	// affinity_config leaf fields
	assert.False(t, cfg.AffinityConfig.IsNull())
	var affinityConfig teoAffinityConfigTestModel
	assert.False(t, cfg.AffinityConfig.As(context.Background(), &affinityConfig, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal object attribute")
	assert.Equal(t, "On", affinityConfig.Switch.ValueString())
	assert.Equal(t, "SessionId", affinityConfig.AffinityMode.ValueString())

	assert.False(t, affinityConfig.SessionIdAffinityConfig.IsNull())
	var sessionIdConfig teoSessionIdAffinityConfigTestModel
	assert.False(t, affinityConfig.SessionIdAffinityConfig.As(context.Background(), &sessionIdConfig, basetypes.ObjectAsOptions{}).HasError(), "failed to unmarshal object attribute")
	assert.Equal(t, "Header", sessionIdConfig.Source.ValueString())
	assert.Equal(t, "EO-Infer-Session-Id", sessionIdConfig.HeaderName.ValueString())
}

// TestTeoInferenceServiceDeploymentRecords111001_Read_EmptyResult verifies
// that an empty record list is a normal business result and produces an empty
// (non-null) record_set list in state.
func TestTeoInferenceServiceDeploymentRecords111001_Read_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecords", func(request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		resp := teov20220901.NewDescribeInferenceServiceDeploymentRecordsResponse()
		resp.Response = &teov20220901.DescribeInferenceServiceDeploymentRecordsResponseParams{
			TotalCount: ptrInt64TeoDeploymentRecords(0),
			RecordSet:  []*teov20220901.InferenceServiceDeploymentRecord{},
		}
		return resp, nil
	})

	d := svcteo.NewTeoInferenceServiceDeploymentRecords111001DataSource()
	setTeoDeploymentRecordsClient(d, client)

	req := newTeoDeploymentRecordsReadRequest(t, "zone-2qtuhspy7cr6", "isvc-none", nil, nil)
	resp := &datasource.ReadResponse{State: newTeoDeploymentRecordsState(t)}
	d.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var state svcteo.TeoInferenceServiceDeploymentRecords111001DataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "failed to read state back: %v", diags)
	assert.False(t, state.RecordSet.IsNull())
	assert.Equal(t, 0, state.RecordSet.Elements().Len())
}

// TestTeoInferenceServiceDeploymentRecords111001_Read_APIError verifies that
// an API failure is reported as an error diagnostic.
func TestTeoInferenceServiceDeploymentRecords111001_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecords", func(request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=inference service not found")
	})

	d := svcteo.NewTeoInferenceServiceDeploymentRecords111001DataSource()
	setTeoDeploymentRecordsClient(d, client)

	req := newTeoDeploymentRecordsReadRequest(t, "zone-2qtuhspy7cr6", "isvc-none", nil, nil)
	resp := &datasource.ReadResponse{State: newTeoDeploymentRecordsState(t)}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error querying TEO inference service deployment records")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ResourceNotFound")
}

// TestTeoInferenceServiceDeploymentRecords111001_Read_ClientNotConfigured
// verifies the nil-client guard.
func TestTeoInferenceServiceDeploymentRecords111001_Read_ClientNotConfigured(t *testing.T) {
	d := svcteo.NewTeoInferenceServiceDeploymentRecords111001DataSource()

	req := newTeoDeploymentRecordsReadRequest(t, "zone-2qtuhspy7cr6", "isvc-2qtuhspy7cr6", nil, nil)
	resp := &datasource.ReadResponse{State: newTeoDeploymentRecordsState(t)}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestTeoInferenceServiceDeploymentRecords111001_MetadataAndSchema validates
// metadata and schema definition.
func TestTeoInferenceServiceDeploymentRecords111001_MetadataAndSchema(t *testing.T) {
	d := svcteo.NewTeoInferenceServiceDeploymentRecords111001DataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_inference_service_deployment_records_1110_01", metaResp.TypeName)

	schemaResp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "zone_id")
	assert.Contains(t, attrs, "service_id")
	assert.Contains(t, attrs, "sort_by")
	assert.Contains(t, attrs, "sort_order")
	assert.Contains(t, attrs, "record_set")

	zoneAttr, ok := attrs["zone_id"].(datasource_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, zoneAttr.Required)

	serviceAttr, ok := attrs["service_id"].(datasource_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, serviceAttr.Required)

	sortByAttr, ok := attrs["sort_by"].(datasource_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, sortByAttr.Optional)

	sortOrderAttr, ok := attrs["sort_order"].(datasource_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, sortOrderAttr.Optional)

	recordSetAttr, ok := attrs["record_set"].(datasource_schema.ListNestedAttribute)
	assert.True(t, ok)
	assert.True(t, recordSetAttr.Computed)

	nested := recordSetAttr.NestedObject.Attributes
	assert.Contains(t, nested, "record_id")
	assert.Contains(t, nested, "operation")
	assert.Contains(t, nested, "status")
	assert.Contains(t, nested, "duration")
	assert.Contains(t, nested, "inference_service_config")
	assert.Contains(t, nested, "create_time")
	assert.Contains(t, nested, "active_status")
}

// TestTeoInferenceServiceDeploymentRecords111001_ConfigureNilProviderData
// verifies that Configure handles nil ProviderData gracefully.
func TestTeoInferenceServiceDeploymentRecords111001_ConfigureNilProviderData(t *testing.T) {
	d := &svcteo.TeoInferenceServiceDeploymentRecords111001DataSource{}
	req := datasource.ConfigureRequest{}
	resp := &datasource.ConfigureResponse{}

	d.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}

// buildTeoDeploymentRecordsResponse builds a fully populated response with one
// deployment record covering every nested attribute.
func buildTeoDeploymentRecordsResponse() *teov20220901.DescribeInferenceServiceDeploymentRecordsResponse {
	resp := teov20220901.NewDescribeInferenceServiceDeploymentRecordsResponse()
	resp.Response = &teov20220901.DescribeInferenceServiceDeploymentRecordsResponseParams{
		TotalCount: ptrInt64TeoDeploymentRecords(1),
		RecordSet: []*teov20220901.InferenceServiceDeploymentRecord{
			{
				RecordId:     ptrStringTeoDeploymentRecords("record-1111"),
				Operation:    ptrStringTeoDeploymentRecords("create"),
				Status:       ptrStringTeoDeploymentRecords("succeeded"),
				Duration:     ptrInt64TeoDeploymentRecords(60),
				CreateTime:   ptrStringTeoDeploymentRecords("2025-01-01T00:00:00+08:00"),
				ActiveStatus: ptrStringTeoDeploymentRecords("active"),
				InferenceServiceConfig: &teov20220901.InferenceServiceConfig{
					ListenPort:   ptrInt64TeoDeploymentRecords(8080),
					RequestPaths: []*string{ptrStringTeoDeploymentRecords("/v1/chat")},
					Containers: []*teov20220901.InferenceContainerConfig{
						{
							ImageType: ptrStringTeoDeploymentRecords("TCR"),
							TcrRepositoryConfig: &teov20220901.InferenceTCRRepositoryConfig{
								TCRType:    ptrStringTeoDeploymentRecords("Enterprise"),
								Image:      ptrStringTeoDeploymentRecords("test.example.com/repo/image:latest"),
								RegistryId: ptrStringTeoDeploymentRecords("tcr-1111"),
								RegionName: ptrStringTeoDeploymentRecords("ap-guangzhou"),
							},
							StartupCommand: ptrStringTeoDeploymentRecords("python serving.py"),
							EnvironmentVariables: []*teov20220901.InferenceEnvironmentVariable{
								{
									Key:   ptrStringTeoDeploymentRecords("MODEL"),
									Value: ptrStringTeoDeploymentRecords("llama3"),
								},
							},
						},
					},
					ResourceConfig: &teov20220901.InferenceResourceConfig{
						ScalingMode:    ptrStringTeoDeploymentRecords("Auto"),
						HardwareSpec:   ptrStringTeoDeploymentRecords("spec-legacy"),
						HardwareSpecId: ptrStringTeoDeploymentRecords("spec-1111"),
						HardwareConfig: &teov20220901.InferenceHardwareConfig{
							GPUNum:   ptrFloat64TeoDeploymentRecords(1),
							CPUNum:   ptrFloat64TeoDeploymentRecords(8),
							MemSize:  ptrInt64TeoDeploymentRecords(32768),
							DiskSize: ptrInt64TeoDeploymentRecords(10240),
						},
						AutoScalingConfig: &teov20220901.InferenceAutoScalingConfig{
							MinInstanceCount: ptrInt64TeoDeploymentRecords(1),
							ScalingPolicies: []*teov20220901.InferenceScalingPolicy{
								{
									PolicyName: ptrStringTeoDeploymentRecords("policy-1"),
									PolicyType: ptrStringTeoDeploymentRecords("ScheduledScaling"),
									ScheduledScalingPolicy: &teov20220901.InferenceScheduledScalingPolicy{
										ScheduledActions: []*teov20220901.InferenceScheduledScalingAction{
											{
												CronExpression:   ptrStringTeoDeploymentRecords("0 0 * * *"),
												MinInstanceCount: ptrInt64TeoDeploymentRecords(2),
											},
										},
										EffectiveRange: &teov20220901.InferenceScheduledScalingEffectiveRange{
											EffectiveType: ptrStringTeoDeploymentRecords("LongTerm"),
											StartDate:     ptrStringTeoDeploymentRecords("2025-01-01"),
											EndDate:       ptrStringTeoDeploymentRecords("2099-12-31"),
										},
										TimeZone: ptrStringTeoDeploymentRecords("UTC"),
									},
								},
							},
						},
						ManualInstanceConfig: nil,
						Concurrency:          ptrInt64TeoDeploymentRecords(1),
					},
					AffinityConfig: &teov20220901.InferenceAffinityConfig{
						Switch:       ptrStringTeoDeploymentRecords("On"),
						AffinityMode: ptrStringTeoDeploymentRecords("SessionId"),
						SessionIdAffinityConfig: &teov20220901.SessionIdAffinityConfig{
							Source:     ptrStringTeoDeploymentRecords("Header"),
							HeaderName: ptrStringTeoDeploymentRecords("EO-Infer-Session-Id"),
						},
					},
				},
			},
		},
	}
	return resp
}
