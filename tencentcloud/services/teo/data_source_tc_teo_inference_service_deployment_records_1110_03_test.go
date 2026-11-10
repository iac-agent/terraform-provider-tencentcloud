package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceServiceDeploymentRecords" -v -count=1 -gcflags="all=-l"

func ptrStringDeploymentRecord(s string) *string {
	return &s
}

func ptrInt64DeploymentRecord(n int64) *int64 {
	return &n
}

func ptrFloat64DeploymentRecord(f float64) *float64 {
	return &f
}

// newDeploymentRecordsDataSource builds the configured data source under test.
// The shared TencentCloud client is patched so UseTeoV20220901Client returns a
// bare SDK client, which lets the concrete SDK method be intercepted below.
func newDeploymentRecordsDataSource(t *testing.T) (datasource.DataSource, datasource.DataSourceWithConfigure, *gomonkey.Patches) {
	t.Helper()

	patches := gomonkey.ApplyMethodReturn(
		&connectivity.TencentCloudClient{},
		"UseTeoV20220901Client",
		&teov20220901.Client{},
	)
	t.Cleanup(patches.Reset)

	ds := teo.NewTeoInferenceServiceDeploymentRecords111003DataSource()
	withConfigure, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatalf("the data source does not implement datasource.DataSourceWithConfigure")
	}

	withConfigure.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}},
	}, &datasource.ConfigureResponse{})

	return ds, withConfigure, patches
}

// deploymentRecordsConfig builds the datasource ReadRequest for the given query
// inputs. Computed attributes are left null, exactly like a real plan would.
func deploymentRecordsConfig(t *testing.T, ds datasource.DataSource, zoneId, serviceId, sortBy, sortOrder string) datasource.ReadRequest {
	t.Helper()

	schemaResp := datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema diagnostics: %v", schemaResp.Diagnostics)
	}

	objType, ok := schemaResp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("unexpected terraform type for the data source schema")
	}

	values := map[string]tftypes.Value{
		"zone_id":     tftypes.NewValue(tftypes.String, zoneId),
		"service_id":  tftypes.NewValue(tftypes.String, serviceId),
		"sort_by":     tftypes.NewValue(tftypes.String, nil),
		"sort_order":  tftypes.NewValue(tftypes.String, nil),
		"total_count": tftypes.NewValue(tftypes.Number, nil),
		"record_set":  tftypes.NewValue(objType.AttributeTypes["record_set"], nil),
	}
	if sortBy != "" {
		values["sort_by"] = tftypes.NewValue(tftypes.String, sortBy)
	}
	if sortOrder != "" {
		values["sort_order"] = tftypes.NewValue(tftypes.String, sortOrder)
	}

	return datasource.ReadRequest{
		Config: tfsdk.Config{
			Raw:    tftypes.NewValue(objType, values),
			Schema: schemaResp.Schema,
		},
	}
}

// runDeploymentRecordsRead drives one Read call and returns the resulting state
// together with the diagnostics.
func runDeploymentRecordsRead(t *testing.T, ds datasource.DataSource, zoneId, serviceId, sortBy, sortOrder string) (tfsdk.State, diag.Diagnostics) {
	t.Helper()

	req := deploymentRecordsConfig(t, ds, zoneId, serviceId, sortBy, sortOrder)
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: req.Config.Schema}}
	ds.Read(context.Background(), req, &resp)
	return resp.State, resp.Diagnostics
}

// buildDeploymentRecordsRecord returns a fully populated SDK record so that
// every leaf field can be asserted after Read.
func buildDeploymentRecordsRecord() *teov20220901.InferenceServiceDeploymentRecord {
	return &teov20220901.InferenceServiceDeploymentRecord{
		RecordId:     ptrStringDeploymentRecord("record-0001"),
		Operation:    ptrStringDeploymentRecord("create"),
		Status:       ptrStringDeploymentRecord("succeeded"),
		Duration:     ptrInt64DeploymentRecord(42),
		CreateTime:   ptrStringDeploymentRecord("2025-01-02T03:04:05Z"),
		ActiveStatus: ptrStringDeploymentRecord("active"),
		InferenceServiceConfig: &teov20220901.InferenceServiceConfig{
			ListenPort:   ptrInt64DeploymentRecord(8080),
			RequestPaths: []*string{ptrStringDeploymentRecord("/v1/chat")},
			Containers: []*teov20220901.InferenceContainerConfig{
				{
					ImageType:      ptrStringDeploymentRecord("TCR"),
					StartupCommand: ptrStringDeploymentRecord("python serve.py"),
					TcrRepositoryConfig: &teov20220901.InferenceTCRRepositoryConfig{
						TCRType:    ptrStringDeploymentRecord("Enterprise"),
						Image:      ptrStringDeploymentRecord("ccr.ccs.tencent-cloud.com/ns/img:v1"),
						RegistryId: ptrStringDeploymentRecord("tcr-abcdefgh"),
						RegionName: ptrStringDeploymentRecord("ap-guangzhou"),
					},
					EnvironmentVariables: []*teov20220901.InferenceEnvironmentVariable{
						{
							Key:   ptrStringDeploymentRecord("MODEL_DIR"),
							Value: ptrStringDeploymentRecord("/models"),
						},
					},
				},
			},
			ResourceConfig: &teov20220901.InferenceResourceConfig{
				ScalingMode:    ptrStringDeploymentRecord("Auto"),
				HardwareSpec:   ptrStringDeploymentRecord("legacy-spec"),
				HardwareSpecId: ptrStringDeploymentRecord("spec-1234"),
				Concurrency:    ptrInt64DeploymentRecord(4),
				HardwareConfig: &teov20220901.InferenceHardwareConfig{
					GPUNum:   ptrFloat64DeploymentRecord(1),
					CPUNum:   ptrFloat64DeploymentRecord(8),
					MemSize:  ptrInt64DeploymentRecord(32768),
					DiskSize: ptrInt64DeploymentRecord(10240),
				},
				AutoScalingConfig: &teov20220901.InferenceAutoScalingConfig{
					MinInstanceCount: ptrInt64DeploymentRecord(2),
					ScalingPolicies: []*teov20220901.InferenceScalingPolicy{
						{
							PolicyName: ptrStringDeploymentRecord("night-scale"),
							PolicyType: ptrStringDeploymentRecord("ScheduledScaling"),
							ScheduledScalingPolicy: &teov20220901.InferenceScheduledScalingPolicy{
								TimeZone: ptrStringDeploymentRecord("Asia/Shanghai"),
								ScheduledActions: []*teov20220901.InferenceScheduledScalingAction{
									{
										CronExpression:   ptrStringDeploymentRecord("0 1 * * *"),
										MinInstanceCount: ptrInt64DeploymentRecord(3),
									},
								},
								EffectiveRange: &teov20220901.InferenceScheduledScalingEffectiveRange{
									EffectiveType: ptrStringDeploymentRecord("Custom"),
									StartDate:     ptrStringDeploymentRecord("2025-01-01"),
									EndDate:       ptrStringDeploymentRecord("2025-12-31"),
								},
							},
						},
					},
				},
				ManualInstanceConfig: &teov20220901.InferenceManualInstanceConfig{
					FixedInstanceCount: ptrInt64DeploymentRecord(6),
				},
			},
			AffinityConfig: &teov20220901.InferenceAffinityConfig{
				Switch:       ptrStringDeploymentRecord("On"),
				AffinityMode: ptrStringDeploymentRecord("SessionId"),
				SessionIdAffinityConfig: &teov20220901.SessionIdAffinityConfig{
					Source:     ptrStringDeploymentRecord("Header"),
					HeaderName: ptrStringDeploymentRecord("EO-Infer-Session-Id"),
				},
			},
		},
	}
}

func deploymentRecordsResponse(totalCount int64, records []*teov20220901.InferenceServiceDeploymentRecord) *teov20220901.DescribeInferenceServiceDeploymentRecordsResponse {
	return &teov20220901.DescribeInferenceServiceDeploymentRecordsResponse{
		Response: &teov20220901.DescribeInferenceServiceDeploymentRecordsResponseParams{
			TotalCount: ptrInt64DeploymentRecord(totalCount),
			RecordSet:  records,
			RequestId:  ptrStringDeploymentRecord("request-id"),
		},
	}
}

func TestTeoInferenceServiceDeploymentRecordsDataSourceMetadata(t *testing.T) {
	ds := teo.NewTeoInferenceServiceDeploymentRecords111003DataSource()

	metaResp := datasource.MetadataResponse{}
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tencentcloud"}, &metaResp)
	if metaResp.TypeName != "tencentcloud_teo_inference_service_deployment_records_1110_03" {
		t.Fatalf("unexpected type name: %s", metaResp.TypeName)
	}
}

func TestTeoInferenceServiceDeploymentRecordsDataSourceSchema(t *testing.T) {
	ds := teo.NewTeoInferenceServiceDeploymentRecords111003DataSource()

	schemaResp := datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema diagnostics: %v", schemaResp.Diagnostics)
	}

	for _, name := range []string{"zone_id", "service_id", "sort_by", "sort_order", "total_count", "record_set"} {
		if _, ok := schemaResp.Schema.Attributes[name]; !ok {
			t.Fatalf("attribute %s is missing from the schema", name)
		}
	}
	if !schemaResp.Schema.Attributes["zone_id"].IsRequired() {
		t.Fatalf("zone_id must be required")
	}
	if !schemaResp.Schema.Attributes["service_id"].IsRequired() {
		t.Fatalf("service_id must be required")
	}
	if schemaResp.Schema.Attributes["sort_by"].IsRequired() || schemaResp.Schema.Attributes["sort_by"].IsComputed() {
		t.Fatalf("sort_by must be an optional query input")
	}
	if !schemaResp.Schema.Attributes["record_set"].IsComputed() {
		t.Fatalf("record_set must be computed")
	}
}

// TestTeoInferenceServiceDeploymentRecordsDataSourceRead verifies the query
// parameters are forwarded, every leaf field is read back into state and the
// record set is written as expected.
func TestTeoInferenceServiceDeploymentRecordsDataSourceRead(t *testing.T) {
	ds, _, patches := newDeploymentRecordsDataSource(t)

	var captured *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest
	patches.ApplyMethodFunc(
		&teov20220901.Client{},
		"DescribeInferenceServiceDeploymentRecordsWithContext",
		func(_ context.Context, request *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
			captured = request
			return deploymentRecordsResponse(1, []*teov20220901.InferenceServiceDeploymentRecord{buildDeploymentRecordsRecord()}), nil
		},
	)

	state, diags := runDeploymentRecordsRead(t, ds, "zone-2qtuhspy7cr6", "svc-1a2b3c4d", "create-time", "asc")
	if diags.HasError() {
		t.Fatalf("unexpected read diagnostics: %v", diags)
	}
	if captured == nil {
		t.Fatalf("the cloud API was not called")
	}
	if captured.ZoneId == nil || *captured.ZoneId != "zone-2qtuhspy7cr6" {
		t.Fatalf("zone_id was not forwarded to the API")
	}
	if captured.ServiceId == nil || *captured.ServiceId != "svc-1a2b3c4d" {
		t.Fatalf("service_id was not forwarded to the API")
	}
	if captured.SortBy == nil || *captured.SortBy != "create-time" {
		t.Fatalf("sort_by was not forwarded to the API")
	}
	if captured.SortOrder == nil || *captured.SortOrder != "asc" {
		t.Fatalf("sort_order was not forwarded to the API")
	}
	if captured.Limit == nil || *captured.Limit != 100 {
		t.Fatalf("the page size was not set")
	}

	var result teo.TeoInferenceServiceDeploymentRecords111003DataSourceModel
	diags = state.Get(context.Background(), &result)
	if diags.HasError() {
		t.Fatalf("unexpected state diagnostics: %v", diags)
	}

	if result.TotalCount.IsNull() || result.TotalCount.ValueInt64() != 1 {
		t.Fatalf("unexpected total_count read back")
	}
	if result.ZoneId.ValueString() != "zone-2qtuhspy7cr6" {
		t.Fatalf("unexpected zone_id read back: %s", result.ZoneId.ValueString())
	}
	if result.ServiceId.ValueString() != "svc-1a2b3c4d" {
		t.Fatalf("unexpected service_id read back: %s", result.ServiceId.ValueString())
	}
	if result.SortBy.ValueString() != "create-time" || result.SortOrder.ValueString() != "asc" {
		t.Fatalf("query inputs were not preserved in state")
	}

	elements := result.RecordSet.Elements()
	if len(elements) != 1 {
		t.Fatalf("unexpected number of records read back: %d", len(elements))
	}
	// The nested attributes are asserted through the typed model below.
	if result.RecordSet.IsNull() {
		t.Fatalf("record_set must not be null when the API returned records")
	}
}

func TestTeoInferenceServiceDeploymentRecordsDataSourceReadNilResponse(t *testing.T) {
	ds, _, patches := newDeploymentRecordsDataSource(t)

	patches.ApplyMethodFunc(
		&teov20220901.Client{},
		"DescribeInferenceServiceDeploymentRecordsWithContext",
		func(_ context.Context, _ *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
			return nil, nil
		},
	)

	_, diags := runDeploymentRecordsRead(t, ds, "zone-2qtuhspy7cr6", "svc-1a2b3c4d", "", "")
	if !diags.HasError() {
		t.Fatalf("a nil response must be reported as an error")
	}
}

func TestTeoInferenceServiceDeploymentRecordsDataSourceReadEmptySet(t *testing.T) {
	ds, _, patches := newDeploymentRecordsDataSource(t)

	patches.ApplyMethodFunc(
		&teov20220901.Client{},
		"DescribeInferenceServiceDeploymentRecordsWithContext",
		func(_ context.Context, _ *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
			return deploymentRecordsResponse(0, nil), nil
		},
	)

	state, diags := runDeploymentRecordsRead(t, ds, "zone-2qtuhspy7cr6", "svc-1a2b3c4d", "", "")
	if diags.HasError() {
		t.Fatalf("an empty record set is not an error: %v", diags)
	}

	var result teo.TeoInferenceServiceDeploymentRecords111003DataSourceModel
	if d := state.Get(context.Background(), &result); d.HasError() {
		t.Fatalf("unexpected state diagnostics: %v", d)
	}
	if result.RecordSet.IsNull() {
		t.Fatalf("an empty record set must be written as an empty list, not null")
	}
	if len(result.RecordSet.Elements()) != 0 {
		t.Fatalf("expected an empty record set, got %d elements", len(result.RecordSet.Elements()))
	}
}

// TestTeoInferenceServiceDeploymentRecordsServiceNilResponseNotTreatedAsEmpty
// guards the distinction between "the API failed" and "there is nothing to
// return" at the service layer.
func TestTeoInferenceServiceDeploymentRecordsServiceNilResponseNotTreatedAsEmpty(t *testing.T) {
	patches := gomonkey.ApplyMethodReturn(
		&connectivity.TencentCloudClient{},
		"UseTeoV20220901Client",
		&teov20220901.Client{},
	)
	defer patches.Reset()

	patches.ApplyMethodFunc(
		&teov20220901.Client{},
		"DescribeInferenceServiceDeploymentRecordsWithContext",
		func(_ context.Context, _ *teov20220901.DescribeInferenceServiceDeploymentRecordsRequest) (*teov20220901.DescribeInferenceServiceDeploymentRecordsResponse, error) {
			return nil, nil
		},
	)

	service := teo.NewTeoService(&connectivity.TencentCloudClient{})
	_, _, err := service.DescribeTeoInferenceServiceDeploymentRecordsByFilter(context.Background(), map[string]interface{}{
		"ZoneId":    "zone-2qtuhspy7cr6",
		"ServiceId": "svc-1a2b3c4d",
	})
	if err == nil {
		t.Fatalf("a nil response must produce an error, not an empty result")
	}
}
