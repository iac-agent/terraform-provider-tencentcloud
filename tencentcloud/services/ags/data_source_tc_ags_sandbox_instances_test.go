package ags_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

func agsSandboxInstancesStr(s string) *string { return &s }
func agsSandboxInstancesInt64(v int64) *int64 { return &v }
func agsSandboxInstancesUint64(v uint64) *uint64 {
	return &v
}
func agsSandboxInstancesBool(v bool) *bool { return &v }

// agsSandboxInstancesSchema returns the framework data source schema used to
// build the raw config/state value.
func agsSandboxInstancesSchema(t *testing.T) schema.Schema {
	d := svags.NewAgsSandboxInstancesDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// agsSandboxInstancesRootType builds the tftypes.Object type for the data
// source root, mirroring every schema attribute.
func agsSandboxInstancesRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":           tftypes.String,
			"instance_ids": tftypes.List{ElementType: tftypes.String},
			"instance_set": tftypes.List{ElementType: agsSandboxInstancesItemType()},
		},
	}
}

// agsSandboxInstancesItemType builds the tftypes.Object type for an
// instance_set element.
func agsSandboxInstancesItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"instance_id":            tftypes.String,
			"tool_id":                tftypes.String,
			"tool_name":              tftypes.String,
			"status":                 tftypes.String,
			"persistent":             tftypes.Bool,
			"timeout_seconds":        tftypes.Number,
			"expires_at":             tftypes.String,
			"stop_reason":            tftypes.String,
			"create_time":            tftypes.String,
			"update_time":            tftypes.String,
			"mount_options":          tftypes.List{ElementType: agsSandboxInstancesMountOptionType()},
			"custom_configuration":   agsSandboxInstancesCustomConfigurationType(),
			"computer_configuration": agsSandboxInstancesComputerConfigurationType(),
			"network_mode":           tftypes.String,
			"metadata":               tftypes.List{ElementType: agsSandboxInstancesMetadataVarType()},
			"auth_mode":              tftypes.String,
		},
	}
}

func agsSandboxInstancesMountOptionType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":       tftypes.String,
			"mount_path": tftypes.String,
			"sub_path":   tftypes.String,
			"read_only":  tftypes.Bool,
		},
	}
}

func agsSandboxInstancesEnvVarType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":  tftypes.String,
			"value": tftypes.String,
		},
	}
}

func agsSandboxInstancesPortType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":     tftypes.String,
			"port":     tftypes.Number,
			"protocol": tftypes.String,
		},
	}
}

func agsSandboxInstancesResourcesType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"cpu":     tftypes.String,
			"memory":  tftypes.String,
			"storage": tftypes.String,
		},
	}
}

func agsSandboxInstancesHttpGetType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"path":   tftypes.String,
			"port":   tftypes.Number,
			"scheme": tftypes.String,
		},
	}
}

func agsSandboxInstancesProbeType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"http_get":          agsSandboxInstancesHttpGetType(),
			"ready_timeout_ms":  tftypes.Number,
			"probe_timeout_ms":  tftypes.Number,
			"probe_period_ms":   tftypes.Number,
			"success_threshold": tftypes.Number,
			"failure_threshold": tftypes.Number,
		},
	}
}

func agsSandboxInstancesDnsConfigType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"servers":  tftypes.Set{ElementType: tftypes.String},
			"searches": tftypes.Set{ElementType: tftypes.String},
			"options":  tftypes.Set{ElementType: tftypes.String},
		},
	}
}

func agsSandboxInstancesCapabilitiesType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"add":  tftypes.Set{ElementType: tftypes.String},
			"drop": tftypes.Set{ElementType: tftypes.String},
		},
	}
}

func agsSandboxInstancesSecurityContextType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"privileged":   tftypes.Bool,
			"capabilities": agsSandboxInstancesCapabilitiesType(),
		},
	}
}

func agsSandboxInstancesCustomConfigurationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"image":               tftypes.String,
			"image_registry_type": tftypes.String,
			"image_digest":        tftypes.String,
			"command":             tftypes.Set{ElementType: tftypes.String},
			"args":                tftypes.Set{ElementType: tftypes.String},
			"env":                 tftypes.List{ElementType: agsSandboxInstancesEnvVarType()},
			"ports":               tftypes.List{ElementType: agsSandboxInstancesPortType()},
			"resources":           agsSandboxInstancesResourcesType(),
			"probe":               agsSandboxInstancesProbeType(),
			"dns_config":          agsSandboxInstancesDnsConfigType(),
			"security_context":    agsSandboxInstancesSecurityContextType(),
		},
	}
}

func agsSandboxInstancesWaaConfigurationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"image_id": tftypes.String,
		},
	}
}

func agsSandboxInstancesOsWorldConfigurationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"version": tftypes.String,
		},
	}
}

func agsSandboxInstancesComputerConfigurationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"waa_configuration":      agsSandboxInstancesWaaConfigurationType(),
			"os_world_configuration": agsSandboxInstancesOsWorldConfigurationType(),
		},
	}
}

func agsSandboxInstancesMetadataVarType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":  tftypes.String,
			"value": tftypes.String,
		},
	}
}

// buildAgsSandboxInstancesRawValue builds the raw config value with the given
// instance id filter; every other attribute is unset (nil).
func buildAgsSandboxInstancesRawValue(instanceIds []string) tftypes.Value {
	var instanceIdsValue tftypes.Value
	if instanceIds == nil {
		instanceIdsValue = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil)
	} else {
		elements := make([]tftypes.Value, 0, len(instanceIds))
		for _, id := range instanceIds {
			elements = append(elements, tftypes.NewValue(tftypes.String, id))
		}
		instanceIdsValue = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, elements)
	}

	attrs := map[string]tftypes.Value{
		"id":           tftypes.NewValue(tftypes.String, nil),
		"instance_ids": instanceIdsValue,
		"instance_set": tftypes.NewValue(tftypes.List{ElementType: agsSandboxInstancesItemType()}, nil),
	}

	return tftypes.NewValue(agsSandboxInstancesRootType(), attrs)
}

// agsSandboxInstancesSampleInstance builds a fully populated fake sandbox
// instance returned by the stubbed SDK method.
func agsSandboxInstancesSampleInstance() *agsv20250920.SandboxInstance {
	return &agsv20250920.SandboxInstance{
		InstanceId:     agsSandboxInstancesStr("sandbox-12345"),
		ToolId:         agsSandboxInstancesStr("tool-67890"),
		ToolName:       agsSandboxInstancesStr("tf-test-tool"),
		Status:         agsSandboxInstancesStr("RUNNING"),
		Persistent:     agsSandboxInstancesBool(false),
		TimeoutSeconds: agsSandboxInstancesUint64(300),
		ExpiresAt:      agsSandboxInstancesStr("2024-01-01T01:00:00Z"),
		CreateTime:     agsSandboxInstancesStr("2024-01-01T00:00:00Z"),
		UpdateTime:     agsSandboxInstancesStr("2024-01-01T00:30:00Z"),
		MountOptions: []*agsv20250920.MountOption{
			{
				Name:      agsSandboxInstancesStr("storage-config-1"),
				MountPath: agsSandboxInstancesStr("/mnt/data"),
				SubPath:   agsSandboxInstancesStr("sub"),
				ReadOnly:  agsSandboxInstancesBool(true),
			},
		},
		CustomConfiguration: &agsv20250920.CustomConfigurationDetail{
			Image:             agsSandboxInstancesStr("registry.example.com/image:latest"),
			ImageRegistryType: agsSandboxInstancesStr("enterprise"),
			ImageDigest:       agsSandboxInstancesStr("sha256:abcdef"),
			Command:           []*string{agsSandboxInstancesStr("/bin/sh")},
			Args:              []*string{agsSandboxInstancesStr("-c")},
			Env: []*agsv20250920.EnvVar{
				{
					Name:  agsSandboxInstancesStr("TF_TEST_VAR"),
					Value: agsSandboxInstancesStr("tf-test-value"),
				},
			},
			Ports: []*agsv20250920.PortConfiguration{
				{
					Name:     agsSandboxInstancesStr("http"),
					Port:     agsSandboxInstancesInt64(8080),
					Protocol: agsSandboxInstancesStr("HTTP"),
				},
			},
			Resources: &agsv20250920.ResourceConfiguration{
				CPU:     agsSandboxInstancesStr("2"),
				Memory:  agsSandboxInstancesStr("4Gi"),
				Storage: agsSandboxInstancesStr("10Gi"),
			},
			Probe: &agsv20250920.ProbeConfiguration{
				HttpGet: &agsv20250920.HttpGetAction{
					Path:   agsSandboxInstancesStr("/healthz"),
					Port:   agsSandboxInstancesInt64(8080),
					Scheme: agsSandboxInstancesStr("HTTP"),
				},
				ReadyTimeoutMs:   agsSandboxInstancesInt64(5000),
				ProbeTimeoutMs:   agsSandboxInstancesInt64(1000),
				ProbePeriodMs:    agsSandboxInstancesInt64(2000),
				SuccessThreshold: agsSandboxInstancesInt64(1),
				FailureThreshold: agsSandboxInstancesInt64(3),
			},
			DNSConfig: &agsv20250920.DNSConfig{
				Servers:  []*string{agsSandboxInstancesStr("10.0.0.1")},
				Searches: []*string{agsSandboxInstancesStr("example.internal")},
				Options:  []*string{agsSandboxInstancesStr("ndots:2")},
			},
			SecurityContext: &agsv20250920.SecurityContext{
				Privileged: agsSandboxInstancesBool(false),
				Capabilities: &agsv20250920.Capabilities{
					Add:  []*string{agsSandboxInstancesStr("SYS_ADMIN")},
					Drop: []*string{agsSandboxInstancesStr("ALL")},
				},
			},
		},
		ComputerConfiguration: &agsv20250920.ComputerConfiguration{
			WAAConfiguration: &agsv20250920.WAAConfiguration{
				ImageId: agsSandboxInstancesStr("waa-image-001"),
			},
			OSWorldConfiguration: &agsv20250920.OSWorldConfiguration{
				Version: agsSandboxInstancesStr("osworld1"),
			},
		},
		NetworkMode: agsSandboxInstancesStr("PUBLIC"),
		Metadata: []*agsv20250920.MetadataVar{
			{
				Name:  agsSandboxInstancesStr("team"),
				Value: agsSandboxInstancesStr("infra"),
			},
		},
		AuthMode: agsSandboxInstancesStr("TOKEN"),
	}
}

// agsSandboxInstancesDescribeResponse wraps the given instances into a stub
// response.
func agsSandboxInstancesDescribeResponse(instances ...*agsv20250920.SandboxInstance) *agsv20250920.DescribeSandboxInstancesResponse {
	return &agsv20250920.DescribeSandboxInstancesResponse{
		Response: &agsv20250920.DescribeSandboxInstancesResponseParams{
			InstanceSet: instances,
			RequestId:   agsSandboxInstancesStr("fake-request-id"),
		},
	}
}

// setAgsSandboxInstancesClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setAgsSandboxInstancesClient(d *svags.AgsSandboxInstancesDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func newAgsSandboxInstancesReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

func newAgsSandboxInstancesReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

// Test-local mirror structs used to decode nested state objects with
// ElementsAs. They mirror the generated schema attributes one to one.

type agsSandboxInstancesMountOptionTestModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

type agsSandboxInstancesEnvVarTestModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type agsSandboxInstancesPortTestModel struct {
	Name     types.String `tfsdk:"name"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

type agsSandboxInstancesResourcesTestModel struct {
	Cpu     types.String `tfsdk:"cpu"`
	Memory  types.String `tfsdk:"memory"`
	Storage types.String `tfsdk:"storage"`
}

type agsSandboxInstancesHttpGetTestModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

type agsSandboxInstancesProbeTestModel struct {
	HttpGet          types.Object `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64  `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64  `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64  `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64  `tfsdk:"success_threshold"`
	FailureThreshold types.Int64  `tfsdk:"failure_threshold"`
}

type agsSandboxInstancesDnsConfigTestModel struct {
	Servers  types.Set `tfsdk:"servers"`
	Searches types.Set `tfsdk:"searches"`
	Options  types.Set `tfsdk:"options"`
}

type agsSandboxInstancesCapabilitiesTestModel struct {
	Add  types.Set `tfsdk:"add"`
	Drop types.Set `tfsdk:"drop"`
}

type agsSandboxInstancesSecurityContextTestModel struct {
	Privileged   types.Bool   `tfsdk:"privileged"`
	Capabilities types.Object `tfsdk:"capabilities"`
}

type agsSandboxInstancesCustomConfigTestModel struct {
	Image             types.String `tfsdk:"image"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	ImageDigest       types.String `tfsdk:"image_digest"`
	Command           types.Set    `tfsdk:"command"`
	Args              types.Set    `tfsdk:"args"`
	Env               types.List   `tfsdk:"env"`
	Ports             types.List   `tfsdk:"ports"`
	Resources         types.Object `tfsdk:"resources"`
	Probe             types.Object `tfsdk:"probe"`
	DnsConfig         types.Object `tfsdk:"dns_config"`
	SecurityContext   types.Object `tfsdk:"security_context"`
}

type agsSandboxInstancesWaaConfigTestModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

type agsSandboxInstancesOsWorldConfigTestModel struct {
	Version types.String `tfsdk:"version"`
}

type agsSandboxInstancesComputerConfigTestModel struct {
	WaaConfiguration     types.Object `tfsdk:"waa_configuration"`
	OsWorldConfiguration types.Object `tfsdk:"os_world_configuration"`
}

type agsSandboxInstancesMetadataVarTestModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// TestAgsSandboxInstances_MetadataAndSchema validates metadata and schema
// definition.
func TestAgsSandboxInstances_MetadataAndSchema(t *testing.T) {
	d := svags.NewAgsSandboxInstancesDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_ags_sandbox_instances", metaResp.TypeName)

	sch := agsSandboxInstancesSchema(t)
	attrs := sch.Attributes

	// The only query input is the optional instance id list filter; the
	// outputs are computed.
	assert.True(t, attrs["instance_ids"].(schema.ListAttribute).Optional)
	assert.False(t, attrs["instance_ids"].(schema.ListAttribute).Required)
	assert.True(t, attrs["instance_set"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["id"].(schema.StringAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
}

// TestAgsSandboxInstances_ReadBasic tests a successful read with one record
// and asserts every leaf field is read back.
func TestAgsSandboxInstances_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxInstancesWithContext", func(ctx context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
		if request.InstanceIds == nil {
			t.Errorf("expected instance_ids to be passed to the API request")
		} else if len(request.InstanceIds) != 1 || *request.InstanceIds[0] != "sandbox-12345" {
			t.Errorf("expected instance_ids [sandbox-12345], got %v", request.InstanceIds)
		}

		return agsSandboxInstancesDescribeResponse(agsSandboxInstancesSampleInstance()), nil
	})

	ds := svags.NewAgsSandboxInstancesDataSource().(*svags.AgsSandboxInstancesDataSource)
	setAgsSandboxInstancesClient(ds, client)

	sch := agsSandboxInstancesSchema(t)
	raw := buildAgsSandboxInstancesRawValue([]string{"sandbox-12345"})
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svags.AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	assert.False(t, result.InstanceSet.IsNull())

	var items []svags.AgsSandboxInstancesItemModel
	resp.Diagnostics.Append(result.InstanceSet.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "sandbox-12345", item.InstanceId.ValueString())
	assert.Equal(t, "tool-67890", item.ToolId.ValueString())
	assert.Equal(t, "tf-test-tool", item.ToolName.ValueString())
	assert.Equal(t, "RUNNING", item.Status.ValueString())
	assert.Equal(t, false, item.Persistent.ValueBool())
	assert.Equal(t, int64(300), item.TimeoutSeconds.ValueInt64())
	assert.Equal(t, "2024-01-01T01:00:00Z", item.ExpiresAt.ValueString())
	assert.True(t, item.StopReason.IsNull())
	assert.Equal(t, "2024-01-01T00:00:00Z", item.CreateTime.ValueString())
	assert.Equal(t, "2024-01-01T00:30:00Z", item.UpdateTime.ValueString())
	assert.Equal(t, "PUBLIC", item.NetworkMode.ValueString())
	assert.Equal(t, "TOKEN", item.AuthMode.ValueString())

	// mount_options
	assert.False(t, item.MountOptions.IsNull())
	var mounts []agsSandboxInstancesMountOptionTestModel
	resp.Diagnostics.Append(item.MountOptions.ElementsAs(context.Background(), &mounts, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	if assert.Len(t, mounts, 1) {
		assert.Equal(t, "storage-config-1", mounts[0].Name.ValueString())
		assert.Equal(t, "/mnt/data", mounts[0].MountPath.ValueString())
		assert.Equal(t, "sub", mounts[0].SubPath.ValueString())
		assert.Equal(t, true, mounts[0].ReadOnly.ValueBool())
	}

	// custom_configuration
	assert.False(t, item.CustomConfiguration.IsNull())
	var custom agsSandboxInstancesCustomConfigTestModel
	resp.Diagnostics.Append(item.CustomConfiguration.As(context.Background(), &custom, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "registry.example.com/image:latest", custom.Image.ValueString())
	assert.Equal(t, "enterprise", custom.ImageRegistryType.ValueString())
	assert.Equal(t, "sha256:abcdef", custom.ImageDigest.ValueString())

	var commands []string
	resp.Diagnostics.Append(custom.Command.ElementsAs(context.Background(), &commands, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"/bin/sh"}, commands)

	var argsList []string
	resp.Diagnostics.Append(custom.Args.ElementsAs(context.Background(), &argsList, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"-c"}, argsList)

	var envVars []agsSandboxInstancesEnvVarTestModel
	resp.Diagnostics.Append(custom.Env.ElementsAs(context.Background(), &envVars, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	if assert.Len(t, envVars, 1) {
		assert.Equal(t, "TF_TEST_VAR", envVars[0].Name.ValueString())
		assert.Equal(t, "tf-test-value", envVars[0].Value.ValueString())
	}

	var ports []agsSandboxInstancesPortTestModel
	resp.Diagnostics.Append(custom.Ports.ElementsAs(context.Background(), &ports, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	if assert.Len(t, ports, 1) {
		assert.Equal(t, "http", ports[0].Name.ValueString())
		assert.Equal(t, int64(8080), ports[0].Port.ValueInt64())
		assert.Equal(t, "HTTP", ports[0].Protocol.ValueString())
	}

	var resources agsSandboxInstancesResourcesTestModel
	resp.Diagnostics.Append(custom.Resources.As(context.Background(), &resources, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "2", resources.Cpu.ValueString())
	assert.Equal(t, "4Gi", resources.Memory.ValueString())
	assert.Equal(t, "10Gi", resources.Storage.ValueString())

	var probe agsSandboxInstancesProbeTestModel
	resp.Diagnostics.Append(custom.Probe.As(context.Background(), &probe, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, int64(5000), probe.ReadyTimeoutMs.ValueInt64())
	assert.Equal(t, int64(1000), probe.ProbeTimeoutMs.ValueInt64())
	assert.Equal(t, int64(2000), probe.ProbePeriodMs.ValueInt64())
	assert.Equal(t, int64(1), probe.SuccessThreshold.ValueInt64())
	assert.Equal(t, int64(3), probe.FailureThreshold.ValueInt64())

	var httpGet agsSandboxInstancesHttpGetTestModel
	resp.Diagnostics.Append(probe.HttpGet.As(context.Background(), &httpGet, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "/healthz", httpGet.Path.ValueString())
	assert.Equal(t, int64(8080), httpGet.Port.ValueInt64())
	assert.Equal(t, "HTTP", httpGet.Scheme.ValueString())

	var dnsConfig agsSandboxInstancesDnsConfigTestModel
	resp.Diagnostics.Append(custom.DnsConfig.As(context.Background(), &dnsConfig, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	var servers []string
	resp.Diagnostics.Append(dnsConfig.Servers.ElementsAs(context.Background(), &servers, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"10.0.0.1"}, servers)
	var options []string
	resp.Diagnostics.Append(dnsConfig.Options.ElementsAs(context.Background(), &options, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"ndots:2"}, options)

	var securityContext agsSandboxInstancesSecurityContextTestModel
	resp.Diagnostics.Append(custom.SecurityContext.As(context.Background(), &securityContext, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, false, securityContext.Privileged.ValueBool())

	var capabilities agsSandboxInstancesCapabilitiesTestModel
	resp.Diagnostics.Append(securityContext.Capabilities.As(context.Background(), &capabilities, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	var capAdd []string
	resp.Diagnostics.Append(capabilities.Add.ElementsAs(context.Background(), &capAdd, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"SYS_ADMIN"}, capAdd)
	var capDrop []string
	resp.Diagnostics.Append(capabilities.Drop.ElementsAs(context.Background(), &capDrop, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"ALL"}, capDrop)

	// computer_configuration
	assert.False(t, item.ComputerConfiguration.IsNull())
	var computer agsSandboxInstancesComputerConfigTestModel
	resp.Diagnostics.Append(item.ComputerConfiguration.As(context.Background(), &computer, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())

	var waa agsSandboxInstancesWaaConfigTestModel
	resp.Diagnostics.Append(computer.WaaConfiguration.As(context.Background(), &waa, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "waa-image-001", waa.ImageId.ValueString())

	var osWorld agsSandboxInstancesOsWorldConfigTestModel
	resp.Diagnostics.Append(computer.OsWorldConfiguration.As(context.Background(), &osWorld, basetypes.ObjectAsOptions{})...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "osworld1", osWorld.Version.ValueString())

	// metadata
	assert.False(t, item.Metadata.IsNull())
	var metadataVars []agsSandboxInstancesMetadataVarTestModel
	resp.Diagnostics.Append(item.Metadata.ElementsAs(context.Background(), &metadataVars, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	if assert.Len(t, metadataVars, 1) {
		assert.Equal(t, "team", metadataVars[0].Name.ValueString())
		assert.Equal(t, "infra", metadataVars[0].Value.ValueString())
	}
}

// TestAgsSandboxInstances_ReadNoFilter tests a successful read without any
// filter (query all instances).
func TestAgsSandboxInstances_ReadNoFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxInstancesWithContext", func(ctx context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
		if request.InstanceIds != nil {
			t.Errorf("expected no instance_ids filter, got %v", request.InstanceIds)
		}

		return agsSandboxInstancesDescribeResponse(agsSandboxInstancesSampleInstance()), nil
	})

	ds := svags.NewAgsSandboxInstancesDataSource().(*svags.AgsSandboxInstancesDataSource)
	setAgsSandboxInstancesClient(ds, client)

	sch := agsSandboxInstancesSchema(t)
	raw := buildAgsSandboxInstancesRawValue(nil)
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svags.AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.False(t, result.InstanceSet.IsNull())

	var items []svags.AgsSandboxInstancesItemModel
	resp.Diagnostics.Append(result.InstanceSet.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)
}

// TestAgsSandboxInstances_ReadEmpty covers an empty result: a legitimate
// answer, not a failure — the data source must succeed with an empty list.
func TestAgsSandboxInstances_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxInstancesWithContext", func(ctx context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
		return &agsv20250920.DescribeSandboxInstancesResponse{Response: &agsv20250920.DescribeSandboxInstancesResponseParams{}}, nil
	})

	ds := svags.NewAgsSandboxInstancesDataSource().(*svags.AgsSandboxInstancesDataSource)
	setAgsSandboxInstancesClient(ds, client)

	sch := agsSandboxInstancesSchema(t)
	raw := buildAgsSandboxInstancesRawValue(nil)
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svags.AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.InstanceSet.IsNull())

	var items []svags.AgsSandboxInstancesItemModel
	resp.Diagnostics.Append(result.InstanceSet.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestAgsSandboxInstances_ReadApiError covers an API-level failure: the error
// must surface through diagnostics instead of silently persisting empty
// state.
func TestAgsSandboxInstances_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxInstancesWithContext", func(ctx context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svags.NewAgsSandboxInstancesDataSource().(*svags.AgsSandboxInstancesDataSource)
	setAgsSandboxInstancesClient(ds, client)

	sch := agsSandboxInstancesSchema(t)
	raw := buildAgsSandboxInstancesRawValue(nil)
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}
