package teo

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

func ptrStringInferenceRecords(s string) *string { return &s }

func ptrInt64InferenceRecords(i int64) *int64 { return &i }

func ptrFloat64InferenceRecords(f float64) *float64 { return &f }

// inferenceRecordsDataSourceSchema returns the framework schema of the data
// source, used to build raw config values for ReadRequest.
func inferenceRecordsDataSourceSchema(t *testing.T) *datasource.SchemaResponse {
	t.Helper()
	ds := NewTeoInferenceServiceDeploymentRecords111002DataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp
}

// newInferenceRecordsReadRequest builds a datasource.ReadRequest whose config
// carries the four query inputs and a null record_set.
func newInferenceRecordsReadRequest(t *testing.T) datasource.ReadRequest {
	t.Helper()
	schemaResp := inferenceRecordsDataSourceSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":    tftypes.String,
			"service_id": tftypes.String,
			"sort_by":    tftypes.String,
			"sort_order": tftypes.String,
			"record_set": tftypes.List{ElementType: tftypes.Object{
				AttributeTypes: inferenceRecordsRecordSetAttributeTypes(),
			}},
		},
	}, map[string]tftypes.Value{
		"zone_id":    tftypes.NewValue(tftypes.String, "zone-2qtuhspy7cr6"),
		"service_id": tftypes.NewValue(tftypes.String, "svc-2qtuhspyeq8r6"),
		"sort_by":    tftypes.NewValue(tftypes.String, "create-time"),
		"sort_order": tftypes.NewValue(tftypes.String, "desc"),
		"record_set": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{
			AttributeTypes: inferenceRecordsRecordSetAttributeTypes(),
		}}, nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = schemaResp.Schema
	return req
}

// inferenceRecordsRecordSetAttributeTypes mirrors the record_set nested object
// types so the raw config value matches the schema.
func inferenceRecordsRecordSetAttributeTypes() map[string]tftypes.Type {
	container := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"image_type": tftypes.String,
		"tcr_repository_config": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"tcr_type":    tftypes.String,
			"image":       tftypes.String,
			"registry_id": tftypes.String,
			"region_name": tftypes.String,
		}}},
		"startup_command": tftypes.String,
		"environment_variables": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"key":   tftypes.String,
			"value": tftypes.String,
		}}},
	}}

	scalingPolicy := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"policy_name": tftypes.String,
		"policy_type": tftypes.String,
		"scheduled_scaling_policy": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"scheduled_actions": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"cron_expression":    tftypes.String,
				"min_instance_count": tftypes.Number,
			}}},
			"effective_range": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"effective_type": tftypes.String,
				"start_date":     tftypes.String,
				"end_date":       tftypes.String,
			}}},
			"time_zone": tftypes.String,
		}}},
	}}

	resourceConfig := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"scaling_mode":     tftypes.String,
		"hardware_spec":    tftypes.String,
		"hardware_spec_id": tftypes.String,
		"hardware_config": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"gpu_num":   tftypes.Number,
			"cpu_num":   tftypes.Number,
			"mem_size":  tftypes.Number,
			"disk_size": tftypes.Number,
		}}},
		"auto_scaling_config": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"min_instance_count": tftypes.Number,
			"scaling_policies":   tftypes.List{ElementType: scalingPolicy},
		}}},
		"manual_instance_config": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"fixed_instance_count": tftypes.Number,
		}}},
		"concurrency": tftypes.Number,
	}}

	affinityConfig := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"switch":        tftypes.String,
		"affinity_mode": tftypes.String,
		"session_id_affinity_config": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"source":      tftypes.String,
			"header_name": tftypes.String,
		}}},
	}}

	serviceConfig := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"listen_port":     tftypes.Number,
		"request_paths":   tftypes.Set{ElementType: tftypes.String},
		"containers":      tftypes.List{ElementType: container},
		"resource_config": tftypes.List{ElementType: resourceConfig},
		"affinity_config": tftypes.List{ElementType: affinityConfig},
	}}

	return map[string]tftypes.Type{
		"record_id":                tftypes.String,
		"operation":                tftypes.String,
		"status":                   tftypes.String,
		"duration":                 tftypes.Number,
		"inference_service_config": tftypes.List{ElementType: serviceConfig},
		"create_time":              tftypes.String,
		"active_status":            tftypes.String,
	}
}

// setInferenceRecordsDataSourceClient injects the shared client through
// Configure, mirroring how the framework provider wires data sources.
func setInferenceRecordsDataSourceClient(ds *TeoInferenceServiceDeploymentRecords111002DataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	ds.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: meta}, &datasource.ConfigureResponse{})
}

// buildInferenceDeploymentRecordResponse returns a fully populated API
// response so every leaf attribute can be asserted after Read.
func buildInferenceDeploymentRecordResponse() *teov20220901.DescribeInferenceServiceDeploymentRecordsResponse {
	resp := &teov20220901.DescribeInferenceServiceDeploymentRecordsResponse{}
	resp.Response = &teov20220901.DescribeInferenceServiceDeploymentRecordsResponseParams{
		TotalCount: ptrInt64InferenceRecords(1),
		RecordSet: []*teov20220901.InferenceServiceDeploymentRecord{
			{
				RecordId:     ptrStringInferenceRecords("record-1"),
				Operation:    ptrStringInferenceRecords("create"),
				Status:       ptrStringInferenceRecords("succeeded"),
				Duration:     ptrInt64InferenceRecords(120),
				CreateTime:   ptrStringInferenceRecords("2025-01-02T03:04:05+08:00"),
				ActiveStatus: ptrStringInferenceRecords("active"),
				InferenceServiceConfig: &teov20220901.InferenceServiceConfig{
					ListenPort:   ptrInt64InferenceRecords(8080),
					RequestPaths: []*string{ptrStringInferenceRecords("/v1/predict"), ptrStringInferenceRecords("/v1/chat")},
					Containers: []*teov20220901.InferenceContainerConfig{
						{
							ImageType:      ptrStringInferenceRecords("TCR"),
							StartupCommand: ptrStringInferenceRecords("python serve.py"),
							TcrRepositoryConfig: &teov20220901.InferenceTCRRepositoryConfig{
								TCRType:    ptrStringInferenceRecords("Enterprise"),
								Image:      ptrStringInferenceRecords("ccr.ccs.tencentyun.com/ns/image:v1"),
								RegistryId: ptrStringInferenceRecords("tcr-1234"),
								RegionName: ptrStringInferenceRecords("ap-guangzhou"),
							},
							EnvironmentVariables: []*teov20220901.InferenceEnvironmentVariable{
								{Key: ptrStringInferenceRecords("MODEL_PATH"), Value: ptrStringInferenceRecords("/models/llm")},
							},
						},
					},
					ResourceConfig: &teov20220901.InferenceResourceConfig{
						ScalingMode:    ptrStringInferenceRecords("Auto"),
						HardwareSpec:   ptrStringInferenceRecords("legacy-spec"),
						HardwareSpecId: ptrStringInferenceRecords("spec-1"),
						HardwareConfig: &teov20220901.InferenceHardwareConfig{
							GPUNum:   ptrFloat64InferenceRecords(1),
							CPUNum:   ptrFloat64InferenceRecords(4),
							MemSize:  ptrInt64InferenceRecords(8192),
							DiskSize: ptrInt64InferenceRecords(10240),
						},
						AutoScalingConfig: &teov20220901.InferenceAutoScalingConfig{
							MinInstanceCount: ptrInt64InferenceRecords(1),
							ScalingPolicies: []*teov20220901.InferenceScalingPolicy{
								{
									PolicyName: ptrStringInferenceRecords("night"),
									PolicyType: ptrStringInferenceRecords("ScheduledScaling"),
									ScheduledScalingPolicy: &teov20220901.InferenceScheduledScalingPolicy{
										ScheduledActions: []*teov20220901.InferenceScheduledScalingAction{
											{
												CronExpression:   ptrStringInferenceRecords("0 20 * * *"),
												MinInstanceCount: ptrInt64InferenceRecords(2),
											},
										},
										EffectiveRange: &teov20220901.InferenceScheduledScalingEffectiveRange{
											EffectiveType: ptrStringInferenceRecords("Custom"),
											StartDate:     ptrStringInferenceRecords("2025-01-01"),
											EndDate:       ptrStringInferenceRecords("2025-12-31"),
										},
										TimeZone: ptrStringInferenceRecords("Asia/Shanghai"),
									},
								},
							},
						},
						ManualInstanceConfig: &teov20220901.InferenceManualInstanceConfig{
							FixedInstanceCount: ptrInt64InferenceRecords(3),
						},
						Concurrency: ptrInt64InferenceRecords(8),
					},
					AffinityConfig: &teov20220901.InferenceAffinityConfig{
						Switch:       ptrStringInferenceRecords("On"),
						AffinityMode: ptrStringInferenceRecords("SessionId"),
						SessionIdAffinityConfig: &teov20220901.SessionIdAffinityConfig{
							Source:     ptrStringInferenceRecords("Header"),
							HeaderName: ptrStringInferenceRecords("EO-Infer-Session-Id"),
						},
					},
				},
			},
		},
	}
	return resp
}

// TestTeoInferenceServiceDeploymentRecords111002DataSourceRead covers the happy
// path and asserts every leaf attribute is read back into state.
func TestTeoInferenceServiceDeploymentRecords111002DataSourceRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotRequest *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext",
		func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
			gotRequest = request
			return buildInferenceDeploymentRecordResponse(), nil
		})

	ds := NewTeoInferenceServiceDeploymentRecords111002DataSource().(*TeoInferenceServiceDeploymentRecords111002DataSource)
	setInferenceRecordsDataSourceClient(ds, client)

	req := newInferenceRecordsReadRequest(t)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	if gotRequest == nil {
		t.Fatal("DescribeInferenceServiceDeploymentRecords should have been called")
	}
	if gotRequest.ZoneId == nil || *gotRequest.ZoneId != "zone-2qtuhspy7cr6" {
		t.Errorf("zone_id in request = %v", gotRequest.ZoneId)
	}
	if gotRequest.ServiceId == nil || *gotRequest.ServiceId != "svc-2qtuhspyeq8r6" {
		t.Errorf("service_id in request = %v", gotRequest.ServiceId)
	}
	if gotRequest.SortBy == nil || *gotRequest.SortBy != "create-time" {
		t.Errorf("sort_by in request = %v", gotRequest.SortBy)
	}
	if gotRequest.SortOrder == nil || *gotRequest.SortOrder != "desc" {
		t.Errorf("sort_order in request = %v", gotRequest.SortOrder)
	}
	if gotRequest.Offset == nil || *gotRequest.Offset != 0 {
		t.Errorf("offset in request = %v, want 0", gotRequest.Offset)
	}
	if gotRequest.Limit == nil || *gotRequest.Limit != teoInferenceDeploymentRecordsPageSize {
		t.Errorf("limit in request = %v, want %d", gotRequest.Limit, teoInferenceDeploymentRecordsPageSize)
	}

	var state TeoInferenceServiceDeploymentRecords111002DataSourceModel
	resp.State.Get(context.Background(), &state)

	if got := state.ZoneId.ValueString(); got != "zone-2qtuhspy7cr6" {
		t.Errorf("zone_id = %q, want %q", got, "zone-2qtuhspy7cr6")
	}
	if got := state.ServiceId.ValueString(); got != "svc-2qtuhspyeq8r6" {
		t.Errorf("service_id = %q, want %q", got, "svc-2qtuhspyeq8r6")
	}
	if got := state.SortBy.ValueString(); got != "create-time" {
		t.Errorf("sort_by = %q, want %q", got, "create-time")
	}
	if got := state.SortOrder.ValueString(); got != "desc" {
		t.Errorf("sort_order = %q, want %q", got, "desc")
	}

	if state.RecordSet.IsNull() || state.RecordSet.IsUnknown() {
		t.Fatal("record_set must be a known list")
	}
	var records []teoInferenceDeploymentRecordModel
	if diags := state.RecordSet.ElementsAs(context.Background(), &records, false); diags.HasError() {
		t.Fatalf("failed to decode record_set: %v", diags)
	}
	if len(records) != 1 {
		t.Fatalf("record_set length = %d, want 1", len(records))
	}

	record := records[0]
	if got := record.RecordId.ValueString(); got != "record-1" {
		t.Errorf("record_id = %q", got)
	}
	if got := record.Operation.ValueString(); got != "create" {
		t.Errorf("operation = %q", got)
	}
	if got := record.Status.ValueString(); got != "succeeded" {
		t.Errorf("status = %q", got)
	}
	if got := record.Duration.ValueInt64(); got != 120 {
		t.Errorf("duration = %d, want 120", got)
	}
	if got := record.CreateTime.ValueString(); got != "2025-01-02T03:04:05+08:00" {
		t.Errorf("create_time = %q", got)
	}
	if got := record.ActiveStatus.ValueString(); got != "active" {
		t.Errorf("active_status = %q", got)
	}

	var configs []teoInferenceServiceConfigModel
	if diags := record.InferenceServiceConfig.ElementsAs(context.Background(), &configs, false); diags.HasError() {
		t.Fatalf("failed to decode inference_service_config: %v", diags)
	}
	if len(configs) != 1 {
		t.Fatalf("inference_service_config length = %d, want 1", len(configs))
	}
	serviceConfig := configs[0]
	if got := serviceConfig.ListenPort.ValueInt64(); got != 8080 {
		t.Errorf("listen_port = %d, want 8080", got)
	}
	if len(serviceConfig.RequestPaths.Elements()) != 2 {
		t.Errorf("request_paths length = %d, want 2", len(serviceConfig.RequestPaths.Elements()))
	}

	var containers []teoInferenceContainerConfigModel
	if diags := serviceConfig.Containers.ElementsAs(context.Background(), &containers, false); diags.HasError() {
		t.Fatalf("failed to decode containers: %v", diags)
	}
	if len(containers) != 1 {
		t.Fatalf("containers length = %d, want 1", len(containers))
	}
	container := containers[0]
	if got := container.ImageType.ValueString(); got != "TCR" {
		t.Errorf("image_type = %q", got)
	}
	if got := container.StartupCommand.ValueString(); got != "python serve.py" {
		t.Errorf("startup_command = %q", got)
	}

	var tcrConfigs []teoInferenceTcrRepositoryConfigModel
	if diags := container.TcrRepositoryConfig.ElementsAs(context.Background(), &tcrConfigs, false); diags.HasError() {
		t.Fatalf("failed to decode tcr_repository_config: %v", diags)
	}
	if len(tcrConfigs) != 1 {
		t.Fatalf("tcr_repository_config length = %d, want 1", len(tcrConfigs))
	}
	if got := tcrConfigs[0].TcrType.ValueString(); got != "Enterprise" {
		t.Errorf("tcr_type = %q", got)
	}
	if got := tcrConfigs[0].Image.ValueString(); got != "ccr.ccs.tencentyun.com/ns/image:v1" {
		t.Errorf("image = %q", got)
	}
	if got := tcrConfigs[0].RegistryId.ValueString(); got != "tcr-1234" {
		t.Errorf("registry_id = %q", got)
	}
	if got := tcrConfigs[0].RegionName.ValueString(); got != "ap-guangzhou" {
		t.Errorf("region_name = %q", got)
	}

	var envVars []teoInferenceEnvironmentVariableModel
	if diags := container.EnvironmentVariables.ElementsAs(context.Background(), &envVars, false); diags.HasError() {
		t.Fatalf("failed to decode environment_variables: %v", diags)
	}
	if len(envVars) != 1 {
		t.Fatalf("environment_variables length = %d, want 1", len(envVars))
	}
	if got := envVars[0].Key.ValueString(); got != "MODEL_PATH" {
		t.Errorf("environment variable key = %q", got)
	}
	if got := envVars[0].Value.ValueString(); got != "/models/llm" {
		t.Errorf("environment variable value = %q", got)
	}

	var resourceConfigs []teoInferenceResourceConfigModel
	if diags := serviceConfig.ResourceConfig.ElementsAs(context.Background(), &resourceConfigs, false); diags.HasError() {
		t.Fatalf("failed to decode resource_config: %v", diags)
	}
	if len(resourceConfigs) != 1 {
		t.Fatalf("resource_config length = %d, want 1", len(resourceConfigs))
	}
	resourceConfig := resourceConfigs[0]
	if got := resourceConfig.ScalingMode.ValueString(); got != "Auto" {
		t.Errorf("scaling_mode = %q", got)
	}
	if got := resourceConfig.HardwareSpec.ValueString(); got != "legacy-spec" {
		t.Errorf("hardware_spec = %q", got)
	}
	if got := resourceConfig.HardwareSpecId.ValueString(); got != "spec-1" {
		t.Errorf("hardware_spec_id = %q", got)
	}
	if got := resourceConfig.Concurrency.ValueInt64(); got != 8 {
		t.Errorf("concurrency = %d, want 8", got)
	}

	var hardwareConfigs []teoInferenceHardwareConfigModel
	if diags := resourceConfig.HardwareConfig.ElementsAs(context.Background(), &hardwareConfigs, false); diags.HasError() {
		t.Fatalf("failed to decode hardware_config: %v", diags)
	}
	if len(hardwareConfigs) != 1 {
		t.Fatalf("hardware_config length = %d, want 1", len(hardwareConfigs))
	}
	if got := hardwareConfigs[0].GpuNum.ValueFloat64(); got != 1 {
		t.Errorf("gpu_num = %v, want 1", got)
	}
	if got := hardwareConfigs[0].CpuNum.ValueFloat64(); got != 4 {
		t.Errorf("cpu_num = %v, want 4", got)
	}
	if got := hardwareConfigs[0].MemSize.ValueInt64(); got != 8192 {
		t.Errorf("mem_size = %d, want 8192", got)
	}
	if got := hardwareConfigs[0].DiskSize.ValueInt64(); got != 10240 {
		t.Errorf("disk_size = %d, want 10240", got)
	}

	var autoScalingConfigs []teoInferenceAutoScalingConfigModel
	if diags := resourceConfig.AutoScalingConfig.ElementsAs(context.Background(), &autoScalingConfigs, false); diags.HasError() {
		t.Fatalf("failed to decode auto_scaling_config: %v", diags)
	}
	if len(autoScalingConfigs) != 1 {
		t.Fatalf("auto_scaling_config length = %d, want 1", len(autoScalingConfigs))
	}
	if got := autoScalingConfigs[0].MinInstanceCount.ValueInt64(); got != 1 {
		t.Errorf("min_instance_count = %d, want 1", got)
	}

	var policies []teoInferenceScalingPolicyModel
	if diags := autoScalingConfigs[0].ScalingPolicies.ElementsAs(context.Background(), &policies, false); diags.HasError() {
		t.Fatalf("failed to decode scaling_policies: %v", diags)
	}
	if len(policies) != 1 {
		t.Fatalf("scaling_policies length = %d, want 1", len(policies))
	}
	if got := policies[0].PolicyName.ValueString(); got != "night" {
		t.Errorf("policy_name = %q", got)
	}
	if got := policies[0].PolicyType.ValueString(); got != "ScheduledScaling" {
		t.Errorf("policy_type = %q", got)
	}

	var scheduledPolicies []teoInferenceScheduledScalingPolicyModel
	if diags := policies[0].ScheduledScalingPolicy.ElementsAs(context.Background(), &scheduledPolicies, false); diags.HasError() {
		t.Fatalf("failed to decode scheduled_scaling_policy: %v", diags)
	}
	if len(scheduledPolicies) != 1 {
		t.Fatalf("scheduled_scaling_policy length = %d, want 1", len(scheduledPolicies))
	}
	if got := scheduledPolicies[0].TimeZone.ValueString(); got != "Asia/Shanghai" {
		t.Errorf("time_zone = %q", got)
	}

	var actions []teoInferenceScheduledScalingActionModel
	if diags := scheduledPolicies[0].ScheduledActions.ElementsAs(context.Background(), &actions, false); diags.HasError() {
		t.Fatalf("failed to decode scheduled_actions: %v", diags)
	}
	if len(actions) != 1 {
		t.Fatalf("scheduled_actions length = %d, want 1", len(actions))
	}
	if got := actions[0].CronExpression.ValueString(); got != "0 20 * * *" {
		t.Errorf("cron_expression = %q", got)
	}
	if got := actions[0].MinInstanceCount.ValueInt64(); got != 2 {
		t.Errorf("scheduled action min_instance_count = %d, want 2", got)
	}

	var effectiveRanges []teoInferenceScheduledScalingEffectiveRangeModel
	if diags := scheduledPolicies[0].EffectiveRange.ElementsAs(context.Background(), &effectiveRanges, false); diags.HasError() {
		t.Fatalf("failed to decode effective_range: %v", diags)
	}
	if len(effectiveRanges) != 1 {
		t.Fatalf("effective_range length = %d, want 1", len(effectiveRanges))
	}
	if got := effectiveRanges[0].EffectiveType.ValueString(); got != "Custom" {
		t.Errorf("effective_type = %q", got)
	}
	if got := effectiveRanges[0].StartDate.ValueString(); got != "2025-01-01" {
		t.Errorf("start_date = %q", got)
	}
	if got := effectiveRanges[0].EndDate.ValueString(); got != "2025-12-31" {
		t.Errorf("end_date = %q", got)
	}

	var manualInstanceConfigs []teoInferenceManualInstanceConfigModel
	if diags := resourceConfig.ManualInstanceConfig.ElementsAs(context.Background(), &manualInstanceConfigs, false); diags.HasError() {
		t.Fatalf("failed to decode manual_instance_config: %v", diags)
	}
	if len(manualInstanceConfigs) != 1 {
		t.Fatalf("manual_instance_config length = %d, want 1", len(manualInstanceConfigs))
	}
	if got := manualInstanceConfigs[0].FixedInstanceCount.ValueInt64(); got != 3 {
		t.Errorf("fixed_instance_count = %d, want 3", got)
	}

	var affinityConfigs []teoInferenceAffinityConfigModel
	if diags := serviceConfig.AffinityConfig.ElementsAs(context.Background(), &affinityConfigs, false); diags.HasError() {
		t.Fatalf("failed to decode affinity_config: %v", diags)
	}
	if len(affinityConfigs) != 1 {
		t.Fatalf("affinity_config length = %d, want 1", len(affinityConfigs))
	}
	if got := affinityConfigs[0].Switch.ValueString(); got != "On" {
		t.Errorf("affinity switch = %q", got)
	}
	if got := affinityConfigs[0].AffinityMode.ValueString(); got != "SessionId" {
		t.Errorf("affinity_mode = %q", got)
	}

	var sessionIdConfigs []teoInferenceSessionIdAffinityConfigModel
	if diags := affinityConfigs[0].SessionIdAffinityConfig.ElementsAs(context.Background(), &sessionIdConfigs, false); diags.HasError() {
		t.Fatalf("failed to decode session_id_affinity_config: %v", diags)
	}
	if len(sessionIdConfigs) != 1 {
		t.Fatalf("session_id_affinity_config length = %d, want 1", len(sessionIdConfigs))
	}
	if got := sessionIdConfigs[0].Source.ValueString(); got != "Header" {
		t.Errorf("session_id_affinity_config source = %q", got)
	}
	if got := sessionIdConfigs[0].HeaderName.ValueString(); got != "EO-Infer-Session-Id" {
		t.Errorf("session_id_affinity_config header_name = %q", got)
	}
}

// TestTeoInferenceServiceDeploymentRecords111002DataSourceRead_emptyRecordSet
// verifies that an empty RecordSet is a valid business state, not an error.
func TestTeoInferenceServiceDeploymentRecords111002DataSourceRead_emptyRecordSet(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext",
		func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
			resp := &teov20220901.DescribeInferenceServiceDeploymentRecordsResponse{}
			resp.Response = &teov20220901.DescribeInferenceServiceDeploymentRecordsResponseParams{
				TotalCount: ptrInt64InferenceRecords(0),
				RecordSet:  []*teov20220901.InferenceServiceDeploymentRecord{},
			}
			return resp, nil
		})

	ds := NewTeoInferenceServiceDeploymentRecords111002DataSource().(*TeoInferenceServiceDeploymentRecords111002DataSource)
	setInferenceRecordsDataSourceClient(ds, client)

	req := newInferenceRecordsReadRequest(t)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	var state TeoInferenceServiceDeploymentRecords111002DataSourceModel
	resp.State.Get(context.Background(), &state)
	if state.RecordSet.IsNull() || state.RecordSet.IsUnknown() {
		t.Fatal("record_set must be a known empty list")
	}
	if len(state.RecordSet.Elements()) != 0 {
		t.Errorf("record_set length = %d, want 0", len(state.RecordSet.Elements()))
	}
}

// TestTeoInferenceServiceDeploymentRecords111002DataSourceRead_apiError verifies
// that a cloud API error surfaces as an error diagnostic.
func TestTeoInferenceServiceDeploymentRecords111002DataSourceRead_apiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext",
		func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
			return nil, errors.New("mock api error")
		})

	ds := NewTeoInferenceServiceDeploymentRecords111002DataSource().(*TeoInferenceServiceDeploymentRecords111002DataSource)
	setInferenceRecordsDataSourceClient(ds, client)

	req := newInferenceRecordsReadRequest(t)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics, got none")
	}
}

// TestTeoInferenceServiceDeploymentRecords111002DataSourceRead_nilResponse
// verifies that a nil API response is treated as a failure, not an empty result.
func TestTeoInferenceServiceDeploymentRecords111002DataSourceRead_nilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext",
		func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
			return nil, nil
		})

	ds := NewTeoInferenceServiceDeploymentRecords111002DataSource().(*TeoInferenceServiceDeploymentRecords111002DataSource)
	setInferenceRecordsDataSourceClient(ds, client)

	req := newInferenceRecordsReadRequest(t)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics for nil response, got none")
	}
}

// TestTeoInferenceServiceDeploymentRecords111002DataSourceMetadata verifies the
// Terraform type name and the required/optional/computed split of the inputs.
func TestTeoInferenceServiceDeploymentRecords111002DataSourceMetadata(t *testing.T) {
	ds := NewTeoInferenceServiceDeploymentRecords111002DataSource()

	metadataResp := &datasource.MetadataResponse{}
	ds.Metadata(context.Background(), datasource.MetadataRequest{}, metadataResp)
	if metadataResp.TypeName != "tencentcloud_teo_inference_service_deployment_records_1110_02" {
		t.Fatalf("TypeName = %q", metadataResp.TypeName)
	}

	schemaResp := inferenceRecordsDataSourceSchema(t)
	attributes := schemaResp.Schema.Attributes
	for _, name := range []string{"zone_id", "service_id", "sort_by", "sort_order", "record_set"} {
		if _, ok := attributes[name]; !ok {
			t.Fatalf("attribute %q is missing from schema", name)
		}
	}
	if !attributes["zone_id"].IsRequired() {
		t.Error("zone_id must be required")
	}
	if !attributes["service_id"].IsRequired() {
		t.Error("service_id must be required")
	}
	if !attributes["sort_by"].IsOptional() || attributes["sort_by"].IsComputed() {
		t.Error("sort_by must be optional and not computed")
	}
	if !attributes["sort_order"].IsOptional() || attributes["sort_order"].IsComputed() {
		t.Error("sort_order must be optional and not computed")
	}
	if !attributes["record_set"].IsComputed() {
		t.Error("record_set must be computed")
	}
}
