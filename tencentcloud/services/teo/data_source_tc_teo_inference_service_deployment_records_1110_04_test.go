package teo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceServiceDeploymentRecords111004DataSource" -v -count=1 -gcflags="all=-l"

func ptrStringDeploymentRecords(s string) *string {
	return &s
}

func ptrInt64DeploymentRecords(n int64) *int64 {
	return &n
}

func ptrFloat64DeploymentRecords(f float64) *float64 {
	return &f
}

// binaryFloat is the tftypes representation of a numeric config value.
func binaryFloat(f float64) tftypes.Value {
	return tftypes.NewValue(tftypes.Number, f)
}

func newDeploymentRecordsTestDataSource(t *testing.T) datasource.DataSourceWithConfigure {
	t.Helper()

	ds := teo.NewTeoInferenceServiceDeploymentRecords111004DataSource().(datasource.DataSourceWithConfigure)

	cfgResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}},
	}, cfgResp)
	if cfgResp.Diagnostics.HasError() {
		t.Fatalf("configure returned diagnostics: %v", cfgResp.Diagnostics)
	}

	return ds
}

func deploymentRecordsSchema(t *testing.T, ds datasource.DataSourceWithConfigure) datasource.SchemaResponse {
	t.Helper()

	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema returned diagnostics: %v", schemaResp.Diagnostics)
	}
	return *schemaResp
}

// newDeploymentRecordsReadRequest builds a config value that only sets the
// attributes listed in values; every other attribute stays null.
func newDeploymentRecordsReadRequest(t *testing.T, schemaResp datasource.SchemaResponse, values map[string]tftypes.Value) datasource.ReadRequest {
	t.Helper()

	objectType, ok := schemaResp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("data source schema must map to a tftypes.Object")
	}

	attrs := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, attrType := range objectType.AttributeTypes {
		if v, ok := values[name]; ok {
			attrs[name] = v
			continue
		}
		attrs[name] = tftypes.NewValue(attrType, nil)
	}

	return datasource.ReadRequest{
		Config: tfsdk.Config{
			Raw:    tftypes.NewValue(objectType, attrs),
			Schema: schemaResp.Schema,
		},
	}
}

func newDeploymentRecordsReadResponse(schemaResp datasource.SchemaResponse) *datasource.ReadResponse {
	return &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: schemaResp.Schema,
		},
	}
}

func readRequestValues() map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"zone_id":    tftypes.NewValue(tftypes.String, "zone-2qtuhspy7cr6"),
		"service_id": tftypes.NewValue(tftypes.String, "service-abc"),
	}
}

func fullDeploymentRecord() *teov20220901.InferenceServiceDeploymentRecord {
	return &teov20220901.InferenceServiceDeploymentRecord{
		RecordId:     ptrStringDeploymentRecords("record-abc123"),
		Operation:    ptrStringDeploymentRecords("create"),
		Status:       ptrStringDeploymentRecords("succeeded"),
		Duration:     ptrInt64DeploymentRecords(120),
		CreateTime:   ptrStringDeploymentRecords("2025-01-01T00:00:00+08:00"),
		ActiveStatus: ptrStringDeploymentRecords("active"),
		InferenceServiceConfig: &teov20220901.InferenceServiceConfig{
			ListenPort:   ptrInt64DeploymentRecords(8080),
			RequestPaths: []*string{ptrStringDeploymentRecords("/v1/chat"), ptrStringDeploymentRecords("/v1/embed")},
			Containers: []*teov20220901.InferenceContainerConfig{
				{
					ImageType:      ptrStringDeploymentRecords("TCR"),
					StartupCommand: ptrStringDeploymentRecords("python app.py"),
					TcrRepositoryConfig: &teov20220901.InferenceTCRRepositoryConfig{
						TCRType:    ptrStringDeploymentRecords("Enterprise"),
						Image:      ptrStringDeploymentRecords("registry.example.com/ns/image:v1"),
						RegistryId: ptrStringDeploymentRecords("tcr-instance-1"),
						RegionName: ptrStringDeploymentRecords("ap-guangzhou"),
					},
					EnvironmentVariables: []*teov20220901.InferenceEnvironmentVariable{
						{Key: ptrStringDeploymentRecords("MODEL_NAME"), Value: ptrStringDeploymentRecords("llama")},
						{Key: ptrStringDeploymentRecords("LOG_LEVEL"), Value: ptrStringDeploymentRecords("info")},
					},
				},
			},
			ResourceConfig: &teov20220901.InferenceResourceConfig{
				ScalingMode:    ptrStringDeploymentRecords("Auto"),
				HardwareSpec:   ptrStringDeploymentRecords("legacy-spec"),
				HardwareSpecId: ptrStringDeploymentRecords("spec-gpu-1"),
				HardwareConfig: &teov20220901.InferenceHardwareConfig{
					GPUNum:   ptrFloat64DeploymentRecords(1),
					CPUNum:   ptrFloat64DeploymentRecords(8),
					MemSize:  ptrInt64DeploymentRecords(32768),
					DiskSize: ptrInt64DeploymentRecords(10240),
				},
				AutoScalingConfig: &teov20220901.InferenceAutoScalingConfig{
					MinInstanceCount: ptrInt64DeploymentRecords(1),
					ScalingPolicies: []*teov20220901.InferenceScalingPolicy{
						{
							PolicyName: ptrStringDeploymentRecords("night"),
							PolicyType: ptrStringDeploymentRecords("ScheduledScaling"),
							ScheduledScalingPolicy: &teov20220901.InferenceScheduledScalingPolicy{
								ScheduledActions: []*teov20220901.InferenceScheduledScalingAction{
									{
										CronExpression:   ptrStringDeploymentRecords("0 1 * * *"),
										MinInstanceCount: ptrInt64DeploymentRecords(3),
									},
								},
								EffectiveRange: &teov20220901.InferenceScheduledScalingEffectiveRange{
									EffectiveType: ptrStringDeploymentRecords("Custom"),
									StartDate:     ptrStringDeploymentRecords("2025-01-01"),
									EndDate:       ptrStringDeploymentRecords("2025-12-31"),
								},
								TimeZone: ptrStringDeploymentRecords("Asia/Shanghai"),
							},
						},
					},
				},
				ManualInstanceConfig: &teov20220901.InferenceManualInstanceConfig{
					FixedInstanceCount: ptrInt64DeploymentRecords(2),
				},
				Concurrency: ptrInt64DeploymentRecords(4),
			},
			AffinityConfig: &teov20220901.InferenceAffinityConfig{
				Switch:       ptrStringDeploymentRecords("On"),
				AffinityMode: ptrStringDeploymentRecords("SessionId"),
				SessionIdAffinityConfig: &teov20220901.SessionIdAffinityConfig{
					Source:     ptrStringDeploymentRecords("Header"),
					HeaderName: ptrStringDeploymentRecords("EO-Infer-Session-Id"),
				},
			},
		},
	}
}

func TestTeoInferenceServiceDeploymentRecords111004DataSource_ReadFullRecord(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	var capturedRequest *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		capturedRequest = request
		return &teov20220901.DescribeInferenceServiceDeploymentRecordsResponse{
			Response: &teov20220901.DescribeInferenceServiceDeploymentRecordsResponseParams{
				TotalCount: ptrInt64DeploymentRecords(1),
				RecordSet:  []*teov20220901.InferenceServiceDeploymentRecord{fullDeploymentRecord()},
				RequestId:  ptrStringDeploymentRecords("fake-request-id"),
			},
		}, nil
	})

	ds := newDeploymentRecordsTestDataSource(t)
	schemaResp := deploymentRecordsSchema(t, ds)

	values := readRequestValues()
	values["sort_by"] = tftypes.NewValue(tftypes.String, "create-time")
	values["sort_order"] = tftypes.NewValue(tftypes.String, "asc")
	values["offset"] = binaryFloat(10)
	values["limit"] = binaryFloat(20)

	readReq := newDeploymentRecordsReadRequest(t, schemaResp, values)
	readResp := newDeploymentRecordsReadResponse(schemaResp)

	ds.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("read returned diagnostics: %v", readResp.Diagnostics)
	}
	if capturedRequest == nil {
		t.Fatal("expected DescribeInferenceServiceDeploymentRecordsWithContext to be called")
	}
	if capturedRequest.ZoneId == nil || *capturedRequest.ZoneId != "zone-2qtuhspy7cr6" {
		t.Errorf("unexpected ZoneId: %v", capturedRequest.ZoneId)
	}
	if capturedRequest.ServiceId == nil || *capturedRequest.ServiceId != "service-abc" {
		t.Errorf("unexpected ServiceId: %v", capturedRequest.ServiceId)
	}
	if capturedRequest.SortBy == nil || *capturedRequest.SortBy != "create-time" {
		t.Errorf("unexpected SortBy: %v", capturedRequest.SortBy)
	}
	if capturedRequest.SortOrder == nil || *capturedRequest.SortOrder != "asc" {
		t.Errorf("unexpected SortOrder: %v", capturedRequest.SortOrder)
	}
	if capturedRequest.Offset == nil || *capturedRequest.Offset != 10 {
		t.Errorf("unexpected Offset: %v", capturedRequest.Offset)
	}
	if capturedRequest.Limit == nil || *capturedRequest.Limit != 20 {
		t.Errorf("unexpected Limit: %v", capturedRequest.Limit)
	}

	var state teo.TeoInferenceServiceDeploymentRecords111004DataSourceModel
	readResp.Diagnostics.Append(readResp.State.Get(context.Background(), &state)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to read state: %v", readResp.Diagnostics)
	}

	if state.TotalCount.IsNull() || state.TotalCount.ValueInt64() != 1 {
		t.Errorf("unexpected total_count: %v", state.TotalCount)
	}
	if state.RecordSet.IsNull() {
		t.Fatal("record_set must not be null when records are returned")
	}
	if len(state.RecordSet.Elements()) != 1 {
		t.Fatalf("expected 1 record, got %d", len(state.RecordSet.Elements()))
	}

	records := make([]teo.TeoInferenceServiceDeploymentRecordState, 0, 1)
	readResp.Diagnostics.Append(state.RecordSet.ElementsAs(context.Background(), &records, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode record_set: %v", readResp.Diagnostics)
	}
	record := records[0]

	if record.RecordId.ValueString() != "record-abc123" {
		t.Errorf("unexpected record_id: %s", record.RecordId.ValueString())
	}
	if record.Operation.ValueString() != "create" {
		t.Errorf("unexpected operation: %s", record.Operation.ValueString())
	}
	if record.Status.ValueString() != "succeeded" {
		t.Errorf("unexpected status: %s", record.Status.ValueString())
	}
	if record.Duration.ValueInt64() != 120 {
		t.Errorf("unexpected duration: %d", record.Duration.ValueInt64())
	}
	if record.CreateTime.ValueString() != "2025-01-01T00:00:00+08:00" {
		t.Errorf("unexpected create_time: %s", record.CreateTime.ValueString())
	}
	if record.ActiveStatus.ValueString() != "active" {
		t.Errorf("unexpected active_status: %s", record.ActiveStatus.ValueString())
	}

	configs := make([]teo.TeoInferenceServiceConfigState, 0, 1)
	readResp.Diagnostics.Append(record.InferenceServiceConfig.ElementsAs(context.Background(), &configs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode inference_service_config: %v", readResp.Diagnostics)
	}
	if len(configs) != 1 {
		t.Fatalf("expected 1 inference_service_config, got %d", len(configs))
	}
	config := configs[0]

	if config.ListenPort.ValueInt64() != 8080 {
		t.Errorf("unexpected listen_port: %d", config.ListenPort.ValueInt64())
	}
	paths := make([]string, 0, 2)
	readResp.Diagnostics.Append(config.RequestPaths.ElementsAs(context.Background(), &paths, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode request_paths: %v", readResp.Diagnostics)
	}
	if len(paths) != 2 {
		t.Errorf("expected 2 request paths, got %d", len(paths))
	}

	containers := make([]teo.TeoInferenceContainerConfigState, 0, 1)
	readResp.Diagnostics.Append(config.Containers.ElementsAs(context.Background(), &containers, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode containers: %v", readResp.Diagnostics)
	}
	if len(containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(containers))
	}
	container := containers[0]
	if container.ImageType.ValueString() != "TCR" {
		t.Errorf("unexpected image_type: %s", container.ImageType.ValueString())
	}
	if container.StartupCommand.ValueString() != "python app.py" {
		t.Errorf("unexpected startup_command: %s", container.StartupCommand.ValueString())
	}

	tcrConfigs := make([]teo.TeoInferenceTcrRepositoryConfigState, 0, 1)
	readResp.Diagnostics.Append(container.TcrRepositoryConfig.ElementsAs(context.Background(), &tcrConfigs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode tcr_repository_config: %v", readResp.Diagnostics)
	}
	if len(tcrConfigs) != 1 {
		t.Fatalf("expected 1 tcr_repository_config, got %d", len(tcrConfigs))
	}
	if tcrConfigs[0].TcrType.ValueString() != "Enterprise" {
		t.Errorf("unexpected tcr_type: %s", tcrConfigs[0].TcrType.ValueString())
	}
	if tcrConfigs[0].Image.ValueString() != "registry.example.com/ns/image:v1" {
		t.Errorf("unexpected image: %s", tcrConfigs[0].Image.ValueString())
	}
	if tcrConfigs[0].RegistryId.ValueString() != "tcr-instance-1" {
		t.Errorf("unexpected registry_id: %s", tcrConfigs[0].RegistryId.ValueString())
	}
	if tcrConfigs[0].RegionName.ValueString() != "ap-guangzhou" {
		t.Errorf("unexpected region_name: %s", tcrConfigs[0].RegionName.ValueString())
	}

	envs := make([]teo.TeoInferenceEnvironmentVariableState, 0, 2)
	readResp.Diagnostics.Append(container.EnvironmentVariables.ElementsAs(context.Background(), &envs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode environment_variables: %v", readResp.Diagnostics)
	}
	if len(envs) != 2 {
		t.Fatalf("expected 2 environment variables, got %d", len(envs))
	}
	if envs[0].Key.ValueString() != "MODEL_NAME" || envs[0].Value.ValueString() != "llama" {
		t.Errorf("unexpected first environment variable: %s=%s", envs[0].Key.ValueString(), envs[0].Value.ValueString())
	}
	if envs[1].Key.ValueString() != "LOG_LEVEL" || envs[1].Value.ValueString() != "info" {
		t.Errorf("unexpected second environment variable: %s=%s", envs[1].Key.ValueString(), envs[1].Value.ValueString())
	}

	resourceConfigs := make([]teo.TeoInferenceResourceConfigState, 0, 1)
	readResp.Diagnostics.Append(config.ResourceConfig.ElementsAs(context.Background(), &resourceConfigs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode resource_config: %v", readResp.Diagnostics)
	}
	if len(resourceConfigs) != 1 {
		t.Fatalf("expected 1 resource_config, got %d", len(resourceConfigs))
	}
	resourceConfig := resourceConfigs[0]
	if resourceConfig.ScalingMode.ValueString() != "Auto" {
		t.Errorf("unexpected scaling_mode: %s", resourceConfig.ScalingMode.ValueString())
	}
	if resourceConfig.HardwareSpec.ValueString() != "legacy-spec" {
		t.Errorf("unexpected hardware_spec: %s", resourceConfig.HardwareSpec.ValueString())
	}
	if resourceConfig.HardwareSpecId.ValueString() != "spec-gpu-1" {
		t.Errorf("unexpected hardware_spec_id: %s", resourceConfig.HardwareSpecId.ValueString())
	}
	if resourceConfig.Concurrency.ValueInt64() != 4 {
		t.Errorf("unexpected concurrency: %d", resourceConfig.Concurrency.ValueInt64())
	}

	hardwareConfigs := make([]teo.TeoInferenceHardwareConfigState, 0, 1)
	readResp.Diagnostics.Append(resourceConfig.HardwareConfig.ElementsAs(context.Background(), &hardwareConfigs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode hardware_config: %v", readResp.Diagnostics)
	}
	if len(hardwareConfigs) != 1 {
		t.Fatalf("expected 1 hardware_config, got %d", len(hardwareConfigs))
	}
	if hardwareConfigs[0].GpuNum.ValueFloat64() != 1 {
		t.Errorf("unexpected gpu_num: %v", hardwareConfigs[0].GpuNum.ValueFloat64())
	}
	if hardwareConfigs[0].CpuNum.ValueFloat64() != 8 {
		t.Errorf("unexpected cpu_num: %v", hardwareConfigs[0].CpuNum.ValueFloat64())
	}
	if hardwareConfigs[0].MemSize.ValueInt64() != 32768 {
		t.Errorf("unexpected mem_size: %d", hardwareConfigs[0].MemSize.ValueInt64())
	}
	if hardwareConfigs[0].DiskSize.ValueInt64() != 10240 {
		t.Errorf("unexpected disk_size: %d", hardwareConfigs[0].DiskSize.ValueInt64())
	}

	autoScalingConfigs := make([]teo.TeoInferenceAutoScalingConfigState, 0, 1)
	readResp.Diagnostics.Append(resourceConfig.AutoScalingConfig.ElementsAs(context.Background(), &autoScalingConfigs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode auto_scaling_config: %v", readResp.Diagnostics)
	}
	if len(autoScalingConfigs) != 1 {
		t.Fatalf("expected 1 auto_scaling_config, got %d", len(autoScalingConfigs))
	}
	if autoScalingConfigs[0].MinInstanceCount.ValueInt64() != 1 {
		t.Errorf("unexpected min_instance_count: %d", autoScalingConfigs[0].MinInstanceCount.ValueInt64())
	}

	policies := make([]teo.TeoInferenceScalingPolicyState, 0, 1)
	readResp.Diagnostics.Append(autoScalingConfigs[0].ScalingPolicies.ElementsAs(context.Background(), &policies, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode scaling_policies: %v", readResp.Diagnostics)
	}
	if len(policies) != 1 {
		t.Fatalf("expected 1 scaling policy, got %d", len(policies))
	}
	if policies[0].PolicyName.ValueString() != "night" {
		t.Errorf("unexpected policy_name: %s", policies[0].PolicyName.ValueString())
	}
	if policies[0].PolicyType.ValueString() != "ScheduledScaling" {
		t.Errorf("unexpected policy_type: %s", policies[0].PolicyType.ValueString())
	}

	scheduledPolicies := make([]teo.TeoInferenceScheduledScalingPolicyState, 0, 1)
	readResp.Diagnostics.Append(policies[0].ScheduledScalingPolicy.ElementsAs(context.Background(), &scheduledPolicies, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode scheduled_scaling_policy: %v", readResp.Diagnostics)
	}
	if len(scheduledPolicies) != 1 {
		t.Fatalf("expected 1 scheduled_scaling_policy, got %d", len(scheduledPolicies))
	}
	if scheduledPolicies[0].TimeZone.ValueString() != "Asia/Shanghai" {
		t.Errorf("unexpected time_zone: %s", scheduledPolicies[0].TimeZone.ValueString())
	}

	actions := make([]teo.TeoInferenceScheduledScalingActionState, 0, 1)
	readResp.Diagnostics.Append(scheduledPolicies[0].ScheduledActions.ElementsAs(context.Background(), &actions, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode scheduled_actions: %v", readResp.Diagnostics)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 scheduled action, got %d", len(actions))
	}
	if actions[0].CronExpression.ValueString() != "0 1 * * *" {
		t.Errorf("unexpected cron_expression: %s", actions[0].CronExpression.ValueString())
	}
	if actions[0].MinInstanceCount.ValueInt64() != 3 {
		t.Errorf("unexpected min_instance_count: %d", actions[0].MinInstanceCount.ValueInt64())
	}

	effectiveRanges := make([]teo.TeoInferenceScheduledScalingEffectiveRangeState, 0, 1)
	readResp.Diagnostics.Append(scheduledPolicies[0].EffectiveRange.ElementsAs(context.Background(), &effectiveRanges, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode effective_range: %v", readResp.Diagnostics)
	}
	if len(effectiveRanges) != 1 {
		t.Fatalf("expected 1 effective_range, got %d", len(effectiveRanges))
	}
	if effectiveRanges[0].EffectiveType.ValueString() != "Custom" {
		t.Errorf("unexpected effective_type: %s", effectiveRanges[0].EffectiveType.ValueString())
	}
	if effectiveRanges[0].StartDate.ValueString() != "2025-01-01" {
		t.Errorf("unexpected start_date: %s", effectiveRanges[0].StartDate.ValueString())
	}
	if effectiveRanges[0].EndDate.ValueString() != "2025-12-31" {
		t.Errorf("unexpected end_date: %s", effectiveRanges[0].EndDate.ValueString())
	}

	manualConfigs := make([]teo.TeoInferenceManualInstanceConfigState, 0, 1)
	readResp.Diagnostics.Append(resourceConfig.ManualInstanceConfig.ElementsAs(context.Background(), &manualConfigs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode manual_instance_config: %v", readResp.Diagnostics)
	}
	if len(manualConfigs) != 1 {
		t.Fatalf("expected 1 manual_instance_config, got %d", len(manualConfigs))
	}
	if manualConfigs[0].FixedInstanceCount.ValueInt64() != 2 {
		t.Errorf("unexpected fixed_instance_count: %d", manualConfigs[0].FixedInstanceCount.ValueInt64())
	}

	affinityConfigs := make([]teo.TeoInferenceAffinityConfigState, 0, 1)
	readResp.Diagnostics.Append(config.AffinityConfig.ElementsAs(context.Background(), &affinityConfigs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode affinity_config: %v", readResp.Diagnostics)
	}
	if len(affinityConfigs) != 1 {
		t.Fatalf("expected 1 affinity_config, got %d", len(affinityConfigs))
	}
	if affinityConfigs[0].Switch.ValueString() != "On" {
		t.Errorf("unexpected switch: %s", affinityConfigs[0].Switch.ValueString())
	}
	if affinityConfigs[0].AffinityMode.ValueString() != "SessionId" {
		t.Errorf("unexpected affinity_mode: %s", affinityConfigs[0].AffinityMode.ValueString())
	}

	sessionConfigs := make([]teo.TeoSessionIdAffinityConfigState, 0, 1)
	readResp.Diagnostics.Append(affinityConfigs[0].SessionIdAffinityConfig.ElementsAs(context.Background(), &sessionConfigs, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode session_id_affinity_config: %v", readResp.Diagnostics)
	}
	if len(sessionConfigs) != 1 {
		t.Fatalf("expected 1 session_id_affinity_config, got %d", len(sessionConfigs))
	}
	if sessionConfigs[0].Source.ValueString() != "Header" {
		t.Errorf("unexpected source: %s", sessionConfigs[0].Source.ValueString())
	}
	if sessionConfigs[0].HeaderName.ValueString() != "EO-Infer-Session-Id" {
		t.Errorf("unexpected header_name: %s", sessionConfigs[0].HeaderName.ValueString())
	}
}

func TestTeoInferenceServiceDeploymentRecords111004DataSource_ReadEmptyRecordSet(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		if request.SortBy != nil || request.SortOrder != nil || request.Offset != nil || request.Limit != nil {
			t.Errorf("optional sorting/pagination fields must not be sent when omitted in the config")
		}
		return &teov20220901.DescribeInferenceServiceDeploymentRecordsResponse{
			Response: &teov20220901.DescribeInferenceServiceDeploymentRecordsResponseParams{
				TotalCount: ptrInt64DeploymentRecords(0),
				RecordSet:  []*teov20220901.InferenceServiceDeploymentRecord{},
			},
		}, nil
	})

	ds := newDeploymentRecordsTestDataSource(t)
	schemaResp := deploymentRecordsSchema(t, ds)

	readReq := newDeploymentRecordsReadRequest(t, schemaResp, readRequestValues())
	readResp := newDeploymentRecordsReadResponse(schemaResp)

	ds.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("read returned diagnostics: %v", readResp.Diagnostics)
	}

	var state teo.TeoInferenceServiceDeploymentRecords111004DataSourceModel
	readResp.Diagnostics.Append(readResp.State.Get(context.Background(), &state)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to read state: %v", readResp.Diagnostics)
	}

	if state.TotalCount.IsNull() || state.TotalCount.ValueInt64() != 0 {
		t.Errorf("unexpected total_count: %v", state.TotalCount)
	}
	if state.RecordSet.IsNull() {
		t.Fatal("an empty record set must be a valid empty list, not null")
	}
	if len(state.RecordSet.Elements()) != 0 {
		t.Errorf("expected empty record set, got %d elements", len(state.RecordSet.Elements()))
	}
}

func TestTeoInferenceServiceDeploymentRecords111004DataSource_ReadNilNestedConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		return &teov20220901.DescribeInferenceServiceDeploymentRecordsResponse{
			Response: &teov20220901.DescribeInferenceServiceDeploymentRecordsResponseParams{
				TotalCount: ptrInt64DeploymentRecords(1),
				RecordSet: []*teov20220901.InferenceServiceDeploymentRecord{
					{
						RecordId:               ptrStringDeploymentRecords("record-minimal"),
						Operation:              ptrStringDeploymentRecords("stop"),
						Status:                 ptrStringDeploymentRecords("processing"),
						InferenceServiceConfig: nil,
					},
				},
			},
		}, nil
	})

	ds := newDeploymentRecordsTestDataSource(t)
	schemaResp := deploymentRecordsSchema(t, ds)

	readReq := newDeploymentRecordsReadRequest(t, schemaResp, readRequestValues())
	readResp := newDeploymentRecordsReadResponse(schemaResp)

	ds.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("read returned diagnostics: %v", readResp.Diagnostics)
	}

	var state teo.TeoInferenceServiceDeploymentRecords111004DataSourceModel
	readResp.Diagnostics.Append(readResp.State.Get(context.Background(), &state)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to read state: %v", readResp.Diagnostics)
	}

	records := make([]teo.TeoInferenceServiceDeploymentRecordState, 0, 1)
	readResp.Diagnostics.Append(state.RecordSet.ElementsAs(context.Background(), &records, false)...)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("failed to decode record_set: %v", readResp.Diagnostics)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].RecordId.ValueString() != "record-minimal" {
		t.Errorf("unexpected record_id: %s", records[0].RecordId.ValueString())
	}
	if records[0].Operation.ValueString() != "stop" {
		t.Errorf("unexpected operation: %s", records[0].Operation.ValueString())
	}
	if records[0].Status.ValueString() != "processing" {
		t.Errorf("unexpected status: %s", records[0].Status.ValueString())
	}
	if !records[0].Duration.IsNull() {
		t.Errorf("expected null duration, got %v", records[0].Duration)
	}
	if !records[0].CreateTime.IsNull() {
		t.Errorf("expected null create_time, got %v", records[0].CreateTime)
	}
	if !records[0].ActiveStatus.IsNull() {
		t.Errorf("expected null active_status, got %v", records[0].ActiveStatus)
	}
	if !records[0].InferenceServiceConfig.IsNull() {
		t.Errorf("expected null inference_service_config, got %v", records[0].InferenceServiceConfig)
	}
}

func TestTeoInferenceServiceDeploymentRecords111004DataSource_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		return nil, nil
	})

	ds := newDeploymentRecordsTestDataSource(t)
	schemaResp := deploymentRecordsSchema(t, ds)

	readReq := newDeploymentRecordsReadRequest(t, schemaResp, readRequestValues())
	readResp := newDeploymentRecordsReadResponse(schemaResp)

	ds.Read(context.Background(), readReq, readResp)

	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics error when the API returns a nil response")
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state must not be written when the API returns a nil response")
	}
}

func TestTeoInferenceServiceDeploymentRecords111004DataSource_ReadNilResponseParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		return &teov20220901.DescribeInferenceServiceDeploymentRecordsResponse{}, nil
	})

	ds := newDeploymentRecordsTestDataSource(t)
	schemaResp := deploymentRecordsSchema(t, ds)

	readReq := newDeploymentRecordsReadRequest(t, schemaResp, readRequestValues())
	readResp := newDeploymentRecordsReadResponse(schemaResp)

	ds.Read(context.Background(), readReq, readResp)

	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics error when the API returns a response without params")
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state must not be written when the API returns a response without params")
	}
}

func TestTeoInferenceServiceDeploymentRecords111004DataSource_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServiceDeploymentRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
		return nil, errors.New("api error")
	})

	ds := newDeploymentRecordsTestDataSource(t)
	schemaResp := deploymentRecordsSchema(t, ds)

	readReq := newDeploymentRecordsReadRequest(t, schemaResp, readRequestValues())
	readResp := newDeploymentRecordsReadResponse(schemaResp)

	ds.Read(context.Background(), readReq, readResp)

	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics error when the API call fails")
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state must not be written when the API call fails")
	}
}
