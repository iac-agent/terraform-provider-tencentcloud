package ags_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

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
			"custom_configuration":   tftypes.Object{AttributeTypes: agsSandboxInstancesCustomConfigurationTypes()},
			"computer_configuration": tftypes.Object{AttributeTypes: agsSandboxInstancesComputerConfigurationTypes()},
			"network_mode":           tftypes.String,
			"metadata":               tftypes.List{ElementType: agsSandboxInstancesMetadataType()},
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

func agsSandboxInstancesCustomConfigurationTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"image":               tftypes.String,
		"image_registry_type": tftypes.String,
		"image_digest":        tftypes.String,
		"command":             tftypes.List{ElementType: tftypes.String},
		"args":                tftypes.List{ElementType: tftypes.String},
		"env":                 tftypes.List{ElementType: agsSandboxInstancesEnvVarType()},
		"ports":               tftypes.List{ElementType: agsSandboxInstancesPortType()},
		"resources": tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"cpu":     tftypes.String,
			"memory":  tftypes.String,
			"storage": tftypes.String,
		}},
		"probe":            tftypes.Object{AttributeTypes: agsSandboxInstancesProbeTypes()},
		"dns_config":       tftypes.Object{AttributeTypes: agsSandboxInstancesDNSConfigTypes()},
		"security_context": tftypes.Object{AttributeTypes: agsSandboxInstancesSecurityContextTypes()},
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

func agsSandboxInstancesProbeTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"http_get":          tftypes.Object{AttributeTypes: agsSandboxInstancesHttpGetTypes()},
		"ready_timeout_ms":  tftypes.Number,
		"probe_timeout_ms":  tftypes.Number,
		"probe_period_ms":   tftypes.Number,
		"success_threshold": tftypes.Number,
		"failure_threshold": tftypes.Number,
	}
}

func agsSandboxInstancesHttpGetTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"path":   tftypes.String,
		"port":   tftypes.Number,
		"scheme": tftypes.String,
	}
}

func agsSandboxInstancesDNSConfigTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"servers":  tftypes.List{ElementType: tftypes.String},
		"searches": tftypes.List{ElementType: tftypes.String},
		"options":  tftypes.List{ElementType: tftypes.String},
	}
}

func agsSandboxInstancesSecurityContextTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"privileged":   tftypes.Bool,
		"capabilities": tftypes.Object{AttributeTypes: agsSandboxInstancesCapabilitiesTypes()},
	}
}

func agsSandboxInstancesCapabilitiesTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"add":  tftypes.List{ElementType: tftypes.String},
		"drop": tftypes.List{ElementType: tftypes.String},
	}
}

func agsSandboxInstancesComputerConfigurationTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"waa_configuration":      tftypes.Object{AttributeTypes: map[string]tftypes.Type{"image_id": tftypes.String}},
		"os_world_configuration": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"version": tftypes.String}},
	}
}

func agsSandboxInstancesMetadataType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":  tftypes.String,
			"value": tftypes.String,
		},
	}
}

// buildAgsSandboxInstancesRawValue builds the raw config value with the given
// query input; every other attribute is unset (nil).
func buildAgsSandboxInstancesRawValue(instanceIds tftypes.Value) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":           tftypes.NewValue(tftypes.String, nil),
		"instance_ids": instanceIds,
		"instance_set": tftypes.NewValue(tftypes.List{ElementType: agsSandboxInstancesItemType()}, nil),
	}

	return tftypes.NewValue(agsSandboxInstancesRootType(), attrs)
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

// agsSandboxInstancesSampleItem builds a fake sandbox instance covering the
// nested structures of the schema.
func agsSandboxInstancesSampleItem() *agsv20250920.SandboxInstance {
	command := "go-httpbin"
	arg := "-port"
	image := "ccr.ccs.tencentyun.com/tf-test/go-httpbin:v2.16.1"
	imageDigest := "sha256:fake-digest"
	envName := "TF_TEST_ENV"
	envValue := "tf-test-value"
	portName := "http"
	var portNumber int64 = 8080
	protocol := "TCP"
	cpu := "1"
	memory := "1Gi"
	storage := "5Gi"
	var readyTimeoutMs int64 = 30000
	var probeTimeoutMs int64 = 1000
	var probePeriodMs int64 = 1000
	var successThreshold int64 = 1
	var failureThreshold int64 = 60
	probePath := "/status/200"
	var probePort int64 = 8080
	scheme := "HTTP"
	server := "10.0.0.1"
	search := "tf-test.local"
	option := "ndots:2"
	capability := "SYS_ADMIN"
	waaImageId := "img-fake-waa"
	osWorldVersion := "osworld1"
	mountName := "tf-test-mount"
	mountPath := "/mnt/tf-test"
	subPath := "sub"
	metadataName := "tf-test-metadata"
	metadataValue := "tf-test-value"
	networkMode := "PUBLIC"
	authMode := "DEFAULT"
	stopReason := "manual"
	createTime := "2025-09-20T10:00:00Z"

	return &agsv20250920.SandboxInstance{
		InstanceId:     agsString("ags-fake-instance-id"),
		ToolId:         agsString("tool-12345"),
		ToolName:       agsString("tf-test-tool"),
		Status:         agsString("RUNNING"),
		Persistent:     agsBool(false),
		TimeoutSeconds: agsUint64(600),
		ExpiresAt:      nil,
		StopReason:     agsString(stopReason),
		CreateTime:     agsString(createTime),
		UpdateTime:     nil,
		MountOptions: []*agsv20250920.MountOption{
			{
				Name:      agsString(mountName),
				MountPath: agsString(mountPath),
				SubPath:   agsString(subPath),
				ReadOnly:  agsBool(true),
			},
		},
		CustomConfiguration: &agsv20250920.CustomConfigurationDetail{
			Image:             agsString(image),
			ImageRegistryType: agsString("personal"),
			ImageDigest:       agsString(imageDigest),
			Command:           []*string{&command},
			Args:              []*string{&arg},
			Env: []*agsv20250920.EnvVar{
				{
					Name:  agsString(envName),
					Value: agsString(envValue),
				},
			},
			Ports: []*agsv20250920.PortConfiguration{
				{
					Name:     agsString(portName),
					Port:     agsInt64(portNumber),
					Protocol: agsString(protocol),
				},
			},
			Resources: &agsv20250920.ResourceConfiguration{
				CPU:     agsString(cpu),
				Memory:  agsString(memory),
				Storage: agsString(storage),
			},
			Probe: &agsv20250920.ProbeConfiguration{
				HttpGet: &agsv20250920.HttpGetAction{
					Path:   agsString(probePath),
					Port:   agsInt64(probePort),
					Scheme: agsString(scheme),
				},
				ReadyTimeoutMs:   agsInt64(readyTimeoutMs),
				ProbeTimeoutMs:   agsInt64(probeTimeoutMs),
				ProbePeriodMs:    agsInt64(probePeriodMs),
				SuccessThreshold: agsInt64(successThreshold),
				FailureThreshold: agsInt64(failureThreshold),
			},
			DNSConfig: &agsv20250920.DNSConfig{
				Servers:  []*string{&server},
				Searches: []*string{&search},
				Options:  []*string{&option},
			},
			SecurityContext: &agsv20250920.SecurityContext{
				Privileged: agsBool(true),
				Capabilities: &agsv20250920.Capabilities{
					Add:  []*string{&capability},
					Drop: nil,
				},
			},
		},
		ComputerConfiguration: &agsv20250920.ComputerConfiguration{
			WAAConfiguration: &agsv20250920.WAAConfiguration{
				ImageId: agsString(waaImageId),
			},
			OSWorldConfiguration: &agsv20250920.OSWorldConfiguration{
				Version: agsString(osWorldVersion),
			},
		},
		NetworkMode: agsString(networkMode),
		Metadata: []*agsv20250920.MetadataVar{
			{
				Name:  agsString(metadataName),
				Value: agsString(metadataValue),
			},
		},
		AuthMode: agsString(authMode),
	}
}

func agsUint64(value uint64) *uint64 { return &value }

// assertStringList asserts a list attribute holds the given string values in
// order. It is used for ListAttribute fields (unlike assertStringSet, which
// is for SetAttribute fields of the sandbox tool resource tests).
func assertStringList(t *testing.T, value attr.Value, expected ...string) {
	t.Helper()

	list, ok := value.(types.List)
	if !ok {
		t.Fatalf("value is not a list: %T", value)
	}
	if list.IsNull() || list.IsUnknown() {
		t.Fatalf("list has no value")
	}

	var actual []string
	if diags := list.ElementsAs(context.Background(), &actual, false); diags.HasError() {
		t.Fatalf("failed to read list elements: %v", diags.Errors())
	}
	if len(actual) != len(expected) {
		t.Fatalf("list has %d elements, want %d", len(actual), len(expected))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("list element %d is %q, want %q", i, actual[i], expected[i])
		}
	}
}

// TestAgsSandboxInstancesMetadataAndSchema validates metadata and schema
// definition.
func TestAgsSandboxInstancesMetadataAndSchema(t *testing.T) {
	d := svags.NewAgsSandboxInstancesDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_ags_sandbox_instances", metaResp.TypeName)

	sch := agsSandboxInstancesSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["instance_ids"].(schema.ListAttribute).Optional)
	assert.True(t, attrs["instance_set"].(schema.ListNestedAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
}

// TestAgsSandboxInstancesReadBasic tests a successful read with one record
// and asserts every leaf field is read back.
func TestAgsSandboxInstancesReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxInstancesWithContext", func(ctx context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
		if request.InstanceIds == nil || len(request.InstanceIds) != 1 || *request.InstanceIds[0] != "ags-fake-instance-id" {
			t.Errorf("expected instance_ids filter to be passed to the API request")
		}

		return &agsv20250920.DescribeSandboxInstancesResponse{Response: &agsv20250920.DescribeSandboxInstancesResponseParams{
			InstanceSet: []*agsv20250920.SandboxInstance{agsSandboxInstancesSampleItem()},
			RequestId:   agsString("fake-request-id"),
		}}, nil
	})

	ds := svags.NewAgsSandboxInstancesDataSource().(*svags.AgsSandboxInstancesDataSource)
	setAgsSandboxInstancesClient(ds, client)

	sch := agsSandboxInstancesSchema(t)
	raw := buildAgsSandboxInstancesRawValue(tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
		tftypes.NewValue(tftypes.String, "ags-fake-instance-id"),
	}))
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assertNoErrorDiagnostics(t, &resp.Diagnostics)

	var result svags.AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assertNoErrorDiagnostics(t, &resp.Diagnostics)

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	assert.False(t, result.InstanceSet.IsNull())

	var itemObjects []types.Object
	resp.Diagnostics.Append(result.InstanceSet.ElementsAs(context.Background(), &itemObjects, false)...)
	assertNoErrorDiagnostics(t, &resp.Diagnostics)
	assert.Len(t, itemObjects, 1)

	item := attributesOf(t, itemObjects[0])
	assertString(t, item, "instance_id", "ags-fake-instance-id")
	assertString(t, item, "tool_id", "tool-12345")
	assertString(t, item, "tool_name", "tf-test-tool")
	assertString(t, item, "status", "RUNNING")
	assertBool(t, item, "persistent", false)
	assertInt64(t, item, "timeout_seconds", 600)
	assertString(t, item, "stop_reason", "manual")
	assertString(t, item, "create_time", "2025-09-20T10:00:00Z")
	assertString(t, item, "network_mode", "PUBLIC")
	assertString(t, item, "auth_mode", "DEFAULT")
	assert.True(t, item["expires_at"].IsNull())

	mountOptions := listElements(t, item, "mount_options")
	assert.Len(t, mountOptions, 1)
	assertString(t, attributesOf(t, mountOptions[0]), "name", "tf-test-mount")
	assertBool(t, attributesOf(t, mountOptions[0]), "read_only", true)

	customConfig := nestedAttributes(t, item, "custom_configuration")
	assertString(t, customConfig, "image", "ccr.ccs.tencentyun.com/tf-test/go-httpbin:v2.16.1")
	assertString(t, customConfig, "image_registry_type", "personal")
	assertString(t, customConfig, "image_digest", "sha256:fake-digest")
	assertStringList(t, customConfig["command"], "go-httpbin")
	assertStringList(t, customConfig["args"], "-port")

	env := listElements(t, customConfig, "env")
	assert.Len(t, env, 1)
	assertString(t, attributesOf(t, env[0]), "name", "TF_TEST_ENV")
	assertString(t, attributesOf(t, env[0]), "value", "tf-test-value")

	ports := listElements(t, customConfig, "ports")
	assert.Len(t, ports, 1)
	assertString(t, attributesOf(t, ports[0]), "name", "http")
	assertInt64(t, attributesOf(t, ports[0]), "port", 8080)
	assertString(t, attributesOf(t, ports[0]), "protocol", "TCP")

	resources := nestedAttributes(t, customConfig, "resources")
	assertString(t, resources, "cpu", "1")
	assertString(t, resources, "memory", "1Gi")
	assertString(t, resources, "storage", "5Gi")

	probe := nestedAttributes(t, customConfig, "probe")
	assertInt64(t, probe, "ready_timeout_ms", 30000)
	assertInt64(t, probe, "probe_timeout_ms", 1000)
	assertInt64(t, probe, "probe_period_ms", 1000)
	assertInt64(t, probe, "success_threshold", 1)
	assertInt64(t, probe, "failure_threshold", 60)
	httpGet := nestedAttributes(t, probe, "http_get")
	assertString(t, httpGet, "path", "/status/200")
	assertInt64(t, httpGet, "port", 8080)
	assertString(t, httpGet, "scheme", "HTTP")

	dnsConfig := nestedAttributes(t, customConfig, "dns_config")
	assertStringList(t, dnsConfig["servers"], "10.0.0.1")
	assertStringList(t, dnsConfig["searches"], "tf-test.local")
	assertStringList(t, dnsConfig["options"], "ndots:2")

	securityContext := nestedAttributes(t, customConfig, "security_context")
	assertBool(t, securityContext, "privileged", true)
	capabilities := nestedAttributes(t, securityContext, "capabilities")
	assertStringList(t, capabilities["add"], "SYS_ADMIN")
	assert.True(t, capabilities["drop"].IsNull())

	computerConfig := nestedAttributes(t, item, "computer_configuration")
	waaConfiguration := nestedAttributes(t, computerConfig, "waa_configuration")
	assertString(t, waaConfiguration, "image_id", "img-fake-waa")
	osWorldConfiguration := nestedAttributes(t, computerConfig, "os_world_configuration")
	assertString(t, osWorldConfiguration, "version", "osworld1")

	metadata := listElements(t, item, "metadata")
	assert.Len(t, metadata, 1)
	assertString(t, attributesOf(t, metadata[0]), "name", "tf-test-metadata")
	assertString(t, attributesOf(t, metadata[0]), "value", "tf-test-value")
}

// TestAgsSandboxInstancesReadEmpty covers an empty result: a legitimate
// answer, not a failure — the data source must succeed with an empty list.
func TestAgsSandboxInstancesReadEmpty(t *testing.T) {
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
	raw := buildAgsSandboxInstancesRawValue(tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil))
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assertNoErrorDiagnostics(t, &resp.Diagnostics)

	var result svags.AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assertNoErrorDiagnostics(t, &resp.Diagnostics)

	assert.False(t, result.InstanceSet.IsNull())

	var items []svags.AgsSandboxInstancesItemModel
	resp.Diagnostics.Append(result.InstanceSet.ElementsAs(context.Background(), &items, false)...)
	assertNoErrorDiagnostics(t, &resp.Diagnostics)
	assert.Len(t, items, 0)
}

// TestAgsSandboxInstancesReadApiError covers an API-level failure: the error
// must surface through diagnostics instead of silently persisting empty state.
func TestAgsSandboxInstancesReadApiError(t *testing.T) {
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
	raw := buildAgsSandboxInstancesRawValue(tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil))
	req := newAgsSandboxInstancesReadRequest(sch, raw)
	resp := newAgsSandboxInstancesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}
