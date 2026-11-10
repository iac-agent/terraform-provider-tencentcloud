package ags

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	agssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

func agsPtrString(s string) *string { return helper.String(s) }

// TestAgsSandboxToolResource_Metadata validates the resource type name.
func TestAgsSandboxToolResource_Metadata(t *testing.T) {
	r := &AgsSandboxToolResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_ags_sandbox_tool" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_ags_sandbox_tool")
	}
}

// TestAgsSandboxToolResource_Schema validates required / optional / computed
// flags on the top-level attributes.
func TestAgsSandboxToolResource_Schema(t *testing.T) {
	r := &AgsSandboxToolResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":                     {computed: true},
		"tool_id":                {computed: true},
		"tool_name":              {required: true},
		"tool_type":              {required: true},
		"network_configuration":  {required: true},
		"description":            {optional: true},
		"default_timeout":        {optional: true},
		"tags":                   {optional: true},
		"client_token":           {optional: true},
		"role_arn":               {optional: true},
		"storage_mounts":         {optional: true},
		"custom_configuration":   {optional: true},
		"computer_configuration": {optional: true},
		"log_configuration":      {optional: true},
		"persistent":             {optional: true},
		"status":                 {computed: true},
		"status_reason":          {computed: true},
		"create_time":            {computed: true},
		"update_time":            {computed: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing", name)
			continue
		}
		if attr.IsRequired() != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.IsRequired(), want.required)
		}
		if attr.IsOptional() != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, attr.IsOptional(), want.optional)
		}
		if attr.IsComputed() != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, attr.IsComputed(), want.computed)
		}
	}
}

// TestAgsSandboxToolResource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces it advertises.
func TestAgsSandboxToolResource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &AgsSandboxToolResource{}
	var _ resource.ResourceWithImportState = &AgsSandboxToolResource{}
}

// TestAgsSandboxToolResource_Factory validates the registry factory.
func TestAgsSandboxToolResource_Factory(t *testing.T) {
	got := NewAgsSandboxToolResource()
	if got == nil {
		t.Fatal("NewAgsSandboxToolResource returned nil")
	}
	if _, ok := got.(*AgsSandboxToolResource); !ok {
		t.Fatal("NewAgsSandboxToolResource returned an unexpected type")
	}
}

// TestAgsSandboxToolResource_ConfigureNilProviderData verifies that a nil
// ProviderData is tolerated (no panic) and produces an error diagnostic only
// when the CRUD methods later request the client.
func TestAgsSandboxToolResource_ConfigureNilProviderData(t *testing.T) {
	r := &AgsSandboxToolResource{}
	configureReq := resource.ConfigureRequest{ProviderData: nil}
	configureResp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), configureReq, configureResp)

	if configureResp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData should not error: %v", configureResp.Diagnostics)
	}
}

// TestDescribeAgsSandboxTool_Found covers a successful lookup by ToolId.
func TestDescribeAgsSandboxTool_Found(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &agssdk.Client{}
	var captured *agssdk.DescribeSandboxToolListRequest
	patches.ApplyMethodFunc(sdkClient, "DescribeSandboxToolListWithContext",
		func(ctx context.Context, request *agssdk.DescribeSandboxToolListRequest) (*agssdk.DescribeSandboxToolListResponse, error) {
			captured = request
			resp := agssdk.NewDescribeSandboxToolListResponse()
			resp.Response = &agssdk.DescribeSandboxToolListResponseParams{
				SandboxToolSet: []*agssdk.SandboxTool{
					{ToolId: agsPtrString("tool-other"), ToolName: agsPtrString("other")},
					{ToolId: agsPtrString("tool-1"), ToolName: agsPtrString("target")},
				},
			}
			return resp, nil
		})

	item, err := describeAgsSandboxTool(context.Background(), sdkClient, "tool-1")
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "tool-1", *item.ToolId)
	assert.Equal(t, "target", *item.ToolName)
	assert.NotNil(t, captured)
	assert.Len(t, captured.ToolIds, 1)
	assert.Equal(t, "tool-1", *captured.ToolIds[0])
}

// TestDescribeAgsSandboxTool_NotFound covers the absent-from-list case, which
// must return nil without an error so the caller can drop the state entry.
func TestDescribeAgsSandboxTool_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &agssdk.Client{}
	patches.ApplyMethodFunc(sdkClient, "DescribeSandboxToolListWithContext",
		func(ctx context.Context, request *agssdk.DescribeSandboxToolListRequest) (*agssdk.DescribeSandboxToolListResponse, error) {
			resp := agssdk.NewDescribeSandboxToolListResponse()
			resp.Response = &agssdk.DescribeSandboxToolListResponseParams{}
			return resp, nil
		})

	item, err := describeAgsSandboxTool(context.Background(), sdkClient, "tool-1")
	assert.NoError(t, err)
	assert.Nil(t, item)
}

// TestDescribeAgsSandboxTool_APIError covers the transport/API failure case.
func TestDescribeAgsSandboxTool_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &agssdk.Client{}
	patches.ApplyMethodFunc(sdkClient, "DescribeSandboxToolListWithContext",
		func(ctx context.Context, request *agssdk.DescribeSandboxToolListRequest) (*agssdk.DescribeSandboxToolListResponse, error) {
			return nil, assert.AnError
		})

	item, err := describeAgsSandboxTool(context.Background(), sdkClient, "tool-1")
	assert.Error(t, err)
	assert.Nil(t, item)
}

// TestExpandAgsNetworkConfiguration covers network_configuration expansion,
// including a null list and a populated VPC config.
func TestExpandAgsNetworkConfiguration(t *testing.T) {
	ctx := context.Background()

	// null list
	got, diags := expandAgsNetworkConfiguration(ctx, types.ListNull(agsNetworkConfigurationObjectType()))
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Nil(t, got)

	// populated
	networkModel := []agsNetworkConfigurationModel{{
		NetworkMode: types.StringValue("VPC"),
		VpcConfig: mustList(ctx, []agsVpcConfigModel{{
			SubnetIds:        mustSet(ctx, []string{"subnet-1", "subnet-2"}),
			SecurityGroupIds: mustSet(ctx, []string{"sg-1"}),
		}}, agsVpcConfigObjectType()),
	}}
	list, listDiags := types.ListValueFrom(ctx, agsNetworkConfigurationObjectType(), networkModel)
	assert.False(t, listDiags.HasError(), "unexpected diagnostics: %v", listDiags)

	got, diags = expandAgsNetworkConfiguration(ctx, list)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.NotNil(t, got)
	assert.Equal(t, "VPC", *got.NetworkMode)
	assert.NotNil(t, got.VpcConfig)
	assert.Len(t, got.VpcConfig.SubnetIds, 2)
	assert.Len(t, got.VpcConfig.SecurityGroupIds, 1)
}

// TestExpandAgsTags covers tag expansion including the null case.
func TestExpandAgsTags(t *testing.T) {
	ctx := context.Background()

	got, diags := expandAgsTags(ctx, types.ListNull(agsTagObjectType()))
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Nil(t, got)

	tagModels := []agsTagModel{
		{Key: types.StringValue("env"), Value: types.StringValue("prod")},
		{Key: types.StringValue("team"), Value: types.StringValue("dev")},
	}
	list, listDiags := types.ListValueFrom(ctx, agsTagObjectType(), tagModels)
	assert.False(t, listDiags.HasError(), "unexpected diagnostics: %v", listDiags)

	got, diags = expandAgsTags(ctx, list)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Len(t, got, 2)
	assert.Equal(t, "env", *got[0].Key)
	assert.Equal(t, "prod", *got[0].Value)
}

// TestExpandAgsStorageMounts covers storage_mounts expansion with a COS source
// and a volume template.
func TestExpandAgsStorageMounts(t *testing.T) {
	ctx := context.Background()

	mountModels := []agsStorageMountModel{{
		Name:      types.StringValue("mount-1"),
		MountPath: types.StringValue("/data"),
		ReadOnly:  types.BoolValue(true),
		StorageSource: mustList(ctx, []agsStorageSourceModel{{
			Cos: mustList(ctx, []agsCosStorageSourceModel{{
				Endpoint:   types.StringValue("cos.ap-guangzhou.myqcloud.com"),
				BucketName: types.StringValue("bucket-1"),
				BucketPath: types.StringValue("/logs"),
			}}, agsCosStorageSourceObjectType()),
		}}, agsStorageSourceObjectType()),
		VolumeTemplate: mustList(ctx, []agsVolumeTemplateRefModel{{
			VolumeTemplateName: types.StringValue("tpl-1"),
			ReuseKey:           types.StringValue("session"),
		}}, agsVolumeTemplateRefObjectType()),
	}}
	list, listDiags := types.ListValueFrom(ctx, agsStorageMountObjectType(), mountModels)
	assert.False(t, listDiags.HasError(), "unexpected diagnostics: %v", listDiags)

	got, diags := expandAgsStorageMounts(ctx, list)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Len(t, got, 1)
	assert.Equal(t, "mount-1", *got[0].Name)
	assert.True(t, *got[0].ReadOnly)
	assert.NotNil(t, got[0].StorageSource)
	assert.NotNil(t, got[0].StorageSource.Cos)
	assert.Equal(t, "bucket-1", *got[0].StorageSource.Cos.BucketName)
	assert.NotNil(t, got[0].VolumeTemplate)
	assert.Equal(t, "tpl-1", *got[0].VolumeTemplate.VolumeTemplateName)
	assert.Nil(t, got[0].Volume)
}

// TestExpandAgsCustomConfiguration covers custom_configuration expansion with
// command/args sets, env, dns_config and security_context.
func TestExpandAgsCustomConfiguration(t *testing.T) {
	ctx := context.Background()

	customModels := []agsCustomConfigurationModel{{
		Image:             types.StringValue("registry/image:1.0"),
		ImageRegistryType: types.StringValue("personal"),
		Command:           mustSet(ctx, []string{"/bin/sh"}),
		Args:              mustSet(ctx, []string{"-c", "run"}),
		Env: mustList(ctx, []agsEnvVarModel{
			{Name: types.StringValue("FOO"), Value: types.StringValue("bar")},
		}, agsEnvVarObjectType()),
		DNSConfig: mustList(ctx, []agsDNSConfigModel{{
			Servers:  mustSet(ctx, []string{"10.0.0.1"}),
			Searches: mustSet(ctx, []string{"example.com"}),
		}}, agsDNSConfigObjectType()),
		SecurityContext: mustList(ctx, []agsSecurityContextModel{{
			Privileged: types.BoolValue(true),
			Capabilities: mustList(ctx, []agsCapabilitiesModel{{
				Add: mustSet(ctx, []string{"SYS_ADMIN"}),
			}}, agsCapabilitiesObjectType()),
		}}, agsSecurityContextObjectType()),
	}}
	list, listDiags := types.ListValueFrom(ctx, agsCustomConfigurationObjectType(), customModels)
	assert.False(t, listDiags.HasError(), "unexpected diagnostics: %v", listDiags)

	got, diags := expandAgsCustomConfiguration(ctx, list)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.NotNil(t, got)
	assert.Equal(t, "registry/image:1.0", *got.Image)
	assert.Len(t, got.Command, 1)
	assert.Len(t, got.Args, 2)
	assert.Len(t, got.Env, 1)
	assert.Equal(t, "FOO", *got.Env[0].Name)
	assert.NotNil(t, got.DNSConfig)
	assert.Len(t, got.DNSConfig.Servers, 1)
	assert.NotNil(t, got.SecurityContext)
	assert.True(t, *got.SecurityContext.Privileged)
	assert.Len(t, got.SecurityContext.Capabilities.Add, 1)
}

// TestFlattenAgsTags_KeepsPrevWhenCloudOmits verifies that a configured tag
// list survives a Describe response that omits tags.
func TestFlattenAgsTags_KeepsPrevWhenCloudOmits(t *testing.T) {
	ctx := context.Background()

	tagModels := []agsTagModel{{Key: types.StringValue("env"), Value: types.StringValue("prod")}}
	prev, prevDiags := types.ListValueFrom(ctx, agsTagObjectType(), tagModels)
	assert.False(t, prevDiags.HasError(), "unexpected diagnostics: %v", prevDiags)

	got, diags := flattenAgsTags(ctx, nil, prev)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.True(t, got.Equal(prev), "expected the previous tag list to be preserved")
}

// TestFlattenAgsCustomConfiguration verifies the round trip from a Describe
// response detail back into the model, including the image_digest echo.
func TestFlattenAgsCustomConfiguration(t *testing.T) {
	ctx := context.Background()

	detail := &agssdk.CustomConfigurationDetail{
		Image:             agsPtrString("registry/image:1.0"),
		ImageRegistryType: agsPtrString("personal"),
		ImageDigest:       agsPtrString("sha256:abc"),
		Command:           []*string{agsPtrString("/bin/sh")},
		Args:              []*string{agsPtrString("-c")},
		Env: []*agssdk.EnvVar{
			{Name: agsPtrString("FOO"), Value: agsPtrString("bar")},
		},
		Ports: []*agssdk.PortConfiguration{
			{Name: agsPtrString("http"), Port: helper.Int64(8080), Protocol: agsPtrString("TCP")},
		},
		Resources: &agssdk.ResourceConfiguration{
			CPU:    agsPtrString("2"),
			Memory: agsPtrString("4Gi"),
		},
	}

	prevList := types.ListNull(agsCustomConfigurationObjectType())
	got, diags := flattenAgsCustomConfiguration(ctx, detail, prevList)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	var models []agsCustomConfigurationModel
	diags = got.ElementsAs(ctx, &models, false)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Len(t, models, 1)
	assert.Equal(t, "registry/image:1.0", models[0].Image.ValueString())
	assert.Len(t, models[0].Command.Elements(), 1)
	var ports []agsPortConfigurationModel
	diags = models[0].Ports.ElementsAs(ctx, &ports, false)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Len(t, ports, 1)
	assert.Equal(t, "http", ports[0].Name.ValueString())
	assert.Equal(t, int64(8080), ports[0].Port.ValueInt64())
	var resources []agsResourceConfigurationModel
	diags = models[0].Resources.ElementsAs(ctx, &resources, false)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Len(t, resources, 1)
	assert.Equal(t, "2", resources[0].Cpu.ValueString())
	assert.True(t, models[0].DNSConfig.IsNull(), "dns_config should stay null when the cloud omits it")
}

// TestFlattenAgsNetworkConfiguration verifies the network configuration
// round trip.
func TestFlattenAgsNetworkConfiguration(t *testing.T) {
	ctx := context.Background()

	config := &agssdk.NetworkConfiguration{
		NetworkMode: agsPtrString("PUBLIC"),
	}
	prevList := types.ListNull(agsNetworkConfigurationObjectType())

	got, diags := flattenAgsNetworkConfiguration(ctx, config, prevList)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	var models []agsNetworkConfigurationModel
	diags = got.ElementsAs(ctx, &models, false)
	assert.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
	assert.Len(t, models, 1)
	assert.Equal(t, "PUBLIC", models[0].NetworkMode.ValueString())
	assert.True(t, models[0].VpcConfig.IsNull(), "vpc_config should stay null when the cloud omits it")
}

// agsSandboxToolAttrTypes returns the tftypes attribute map of the resource
// schema for raw plan/state construction in tests.
func agsSandboxToolAttrTypes() map[string]tftypes.Type {
	emptyObj := tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}
	stringSet := tftypes.Set{ElementType: tftypes.String}
	return map[string]tftypes.Type{
		"id":        tftypes.String,
		"tool_id":   tftypes.String,
		"tool_name": tftypes.String,
		"tool_type": tftypes.String,
		"network_configuration": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"network_mode": tftypes.String,
			"vpc_config": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"subnet_ids":         stringSet,
				"security_group_ids": stringSet,
			}}},
		}}},
		"description":     tftypes.String,
		"default_timeout": tftypes.String,
		"tags": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"key":   tftypes.String,
			"value": tftypes.String,
		}}},
		"client_token":           tftypes.String,
		"role_arn":               tftypes.String,
		"storage_mounts":         tftypes.List{ElementType: emptyObj},
		"custom_configuration":   tftypes.List{ElementType: emptyObj},
		"computer_configuration": tftypes.List{ElementType: emptyObj},
		"log_configuration":      tftypes.List{ElementType: emptyObj},
		"persistent":             tftypes.Bool,
		"status":                 tftypes.String,
		"status_reason":          tftypes.String,
		"create_time":            tftypes.String,
		"update_time":            tftypes.String,
	}
}

// agsNullState builds an all-null state/plan object of the resource type.
func agsNullState() tftypes.Value {
	objType := tftypes.Object{AttributeTypes: agsSandboxToolAttrTypes()}
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(objType, values)
}

// TestAgsSandboxToolResource_Create_MockedAPI runs the Create flow against
// mocked cloud APIs and asserts the captured request plus the written state.
func TestAgsSandboxToolResource_Create_MockedAPI(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	sdkClient := &agssdk.Client{}
	patches.ApplyMethodReturn(client, "UseAgsV20250920Client", sdkClient)

	var capturedCreate *agssdk.CreateSandboxToolRequest
	patches.ApplyMethodFunc(sdkClient, "CreateSandboxToolWithContext",
		func(ctx context.Context, request *agssdk.CreateSandboxToolRequest) (*agssdk.CreateSandboxToolResponse, error) {
			capturedCreate = request
			resp := agssdk.NewCreateSandboxToolResponse()
			resp.Response = &agssdk.CreateSandboxToolResponseParams{
				ToolId: agsPtrString("tool-123"),
			}
			return resp, nil
		})
	var capturedDescribe *agssdk.DescribeSandboxToolListRequest
	patches.ApplyMethodFunc(sdkClient, "DescribeSandboxToolListWithContext",
		func(ctx context.Context, request *agssdk.DescribeSandboxToolListRequest) (*agssdk.DescribeSandboxToolListResponse, error) {
			capturedDescribe = request
			resp := agssdk.NewDescribeSandboxToolListResponse()
			resp.Response = &agssdk.DescribeSandboxToolListResponseParams{
				SandboxToolSet: []*agssdk.SandboxTool{{
					ToolId:   agsPtrString("tool-123"),
					ToolName: agsPtrString("my-tool"),
					ToolType: agsPtrString("browser"),
					Status:   agsPtrString("ACTIVE"),
					NetworkConfiguration: &agssdk.NetworkConfiguration{
						NetworkMode: agsPtrString("PUBLIC"),
					},
					CreateTime: agsPtrString("2025-01-01T00:00:00Z"),
				}},
			}
			return resp, nil
		})

	r := &AgsSandboxToolResource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &resource.ConfigureResponse{})
	if r.Client() == nil {
		t.Fatal("failed to configure the resource client")
	}

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", schemaResp.Diagnostics)
	}

	objType := tftypes.Object{AttributeTypes: agsSandboxToolAttrTypes()}
	nulls := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		nulls[name] = tftypes.NewValue(attrType, nil)
	}
	nulls["tool_name"] = tftypes.NewValue(tftypes.String, "my-tool")
	nulls["tool_type"] = tftypes.NewValue(tftypes.String, "browser")
	nulls["network_configuration"] = tftypes.NewValue(agsSandboxToolAttrTypes()["network_configuration"],
		[]tftypes.Value{tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"network_mode": tftypes.String,
			"vpc_config":   agsSandboxToolAttrTypes()["network_configuration"].(tftypes.List).ElementType.(tftypes.Object).AttributeTypes["vpc_config"],
		}}, map[string]tftypes.Value{
			"network_mode": tftypes.NewValue(tftypes.String, "PUBLIC"),
			"vpc_config":   tftypes.NewValue(agsSandboxToolAttrTypes()["network_configuration"].(tftypes.List).ElementType.(tftypes.Object).AttributeTypes["vpc_config"], nil),
		})})
	raw := tftypes.NewValue(objType, nulls)

	req := resource.CreateRequest{
		Config: tfsdk.Config{Raw: raw, Schema: schemaResp.Schema},
		Plan:   tfsdk.Plan{Raw: raw, Schema: schemaResp.Schema},
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{Raw: agsNullState(), Schema: schemaResp.Schema},
	}
	r.Create(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
	assert.NotNil(t, capturedCreate)
	assert.Equal(t, "my-tool", *capturedCreate.ToolName)
	assert.Equal(t, "browser", *capturedCreate.ToolType)
	assert.NotNil(t, capturedCreate.NetworkConfiguration)
	assert.Equal(t, "PUBLIC", *capturedCreate.NetworkConfiguration.NetworkMode)
	assert.Equal(t, "tool-123", *capturedDescribe.ToolIds[0])

	stateMap, err := resp.State.Raw.AsValueMap()
	assert.NoError(t, err)
	var toolId, toolName, status string
	assert.NoError(t, stateMap["tool_id"].As(&toolId))
	assert.NoError(t, stateMap["tool_name"].As(&toolName))
	assert.NoError(t, stateMap["status"].As(&status))
	assert.Equal(t, "tool-123", toolId)
	assert.Equal(t, "my-tool", toolName)
	assert.Equal(t, "ACTIVE", status)
}

// TestAgsSandboxToolResource_Read_RemovesMissingState verifies that a tool
// absent from the Describe list is removed from state.
func TestAgsSandboxToolResource_Read_RemovesMissingState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	sdkClient := &agssdk.Client{}
	patches.ApplyMethodReturn(client, "UseAgsV20250920Client", sdkClient)
	patches.ApplyMethodFunc(sdkClient, "DescribeSandboxToolListWithContext",
		func(ctx context.Context, request *agssdk.DescribeSandboxToolListRequest) (*agssdk.DescribeSandboxToolListResponse, error) {
			resp := agssdk.NewDescribeSandboxToolListResponse()
			resp.Response = &agssdk.DescribeSandboxToolListResponseParams{}
			return resp, nil
		})

	r := &AgsSandboxToolResource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &resource.ConfigureResponse{})

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

	state := agsNullState()
	stateMap, err := state.AsValueMap()
	assert.NoError(t, err)
	stateMap["id"] = tftypes.NewValue(tftypes.String, "tool-123")
	state = tftypes.NewValue(state.Type(), stateMap)

	req := resource.ReadRequest{State: tfsdk.State{Raw: state, Schema: schemaResp.Schema}}
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: state, Schema: schemaResp.Schema}}
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
	assert.True(t, resp.State.Raw.IsNull(), "expected the state entry to be removed")
}

// TestAgsSandboxToolResource_Update_MockedAPI verifies the update path sends
// only the fields accepted by UpdateSandboxTool.
func TestAgsSandboxToolResource_Update_MockedAPI(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	sdkClient := &agssdk.Client{}
	patches.ApplyMethodReturn(client, "UseAgsV20250920Client", sdkClient)

	var capturedUpdate *agssdk.UpdateSandboxToolRequest
	patches.ApplyMethodFunc(sdkClient, "UpdateSandboxToolWithContext",
		func(ctx context.Context, request *agssdk.UpdateSandboxToolRequest) (*agssdk.UpdateSandboxToolResponse, error) {
			capturedUpdate = request
			resp := agssdk.NewUpdateSandboxToolResponse()
			resp.Response = &agssdk.UpdateSandboxToolResponseParams{}
			return resp, nil
		})
	patches.ApplyMethodFunc(sdkClient, "DescribeSandboxToolListWithContext",
		func(ctx context.Context, request *agssdk.DescribeSandboxToolListRequest) (*agssdk.DescribeSandboxToolListResponse, error) {
			resp := agssdk.NewDescribeSandboxToolListResponse()
			resp.Response = &agssdk.DescribeSandboxToolListResponseParams{
				SandboxToolSet: []*agssdk.SandboxTool{{
					ToolId:   agsPtrString("tool-123"),
					ToolName: agsPtrString("my-tool"),
					ToolType: agsPtrString("browser"),
					Status:   agsPtrString("ACTIVE"),
					NetworkConfiguration: &agssdk.NetworkConfiguration{
						NetworkMode: agsPtrString("PUBLIC"),
					},
				}},
			}
			return resp, nil
		})

	r := &AgsSandboxToolResource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &resource.ConfigureResponse{})

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", schemaResp.Diagnostics)
	}

	objType := tftypes.Object{AttributeTypes: agsSandboxToolAttrTypes()}
	nulls := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		nulls[name] = tftypes.NewValue(attrType, nil)
	}
	nulls["id"] = tftypes.NewValue(tftypes.String, "tool-123")
	nulls["tool_name"] = tftypes.NewValue(tftypes.String, "my-tool")
	nulls["tool_type"] = tftypes.NewValue(tftypes.String, "browser")
	nulls["description"] = tftypes.NewValue(tftypes.String, "updated")
	nulls["default_timeout"] = tftypes.NewValue(tftypes.String, "5m")
	nulls["network_configuration"] = tftypes.NewValue(agsSandboxToolAttrTypes()["network_configuration"],
		[]tftypes.Value{tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"network_mode": tftypes.String,
			"vpc_config":   agsSandboxToolAttrTypes()["network_configuration"].(tftypes.List).ElementType.(tftypes.Object).AttributeTypes["vpc_config"],
		}}, map[string]tftypes.Value{
			"network_mode": tftypes.NewValue(tftypes.String, "PUBLIC"),
			"vpc_config":   tftypes.NewValue(agsSandboxToolAttrTypes()["network_configuration"].(tftypes.List).ElementType.(tftypes.Object).AttributeTypes["vpc_config"], nil),
		})})

	planRaw := tftypes.NewValue(objType, nulls)
	stateRaw := tftypes.NewValue(objType, nulls)

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Raw: planRaw, Schema: schemaResp.Schema},
		Plan:   tfsdk.Plan{Raw: planRaw, Schema: schemaResp.Schema},
		State:  tfsdk.State{Raw: stateRaw, Schema: schemaResp.Schema},
	}
	resp := &resource.UpdateResponse{
		State: tfsdk.State{Raw: agsNullState(), Schema: schemaResp.Schema},
	}
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
	assert.NotNil(t, capturedUpdate)
	assert.Equal(t, "tool-123", *capturedUpdate.ToolId)
	assert.Equal(t, "updated", *capturedUpdate.Description)
	assert.NotNil(t, capturedUpdate.NetworkConfiguration)
	assert.Equal(t, "PUBLIC", *capturedUpdate.NetworkConfiguration.NetworkMode)
	assert.Nil(t, capturedUpdate.StorageMounts)
	assert.Nil(t, capturedUpdate.LogConfiguration)
}

// TestAgsSandboxToolResource_Update_APIError verifies the update failure path
// reports a diagnostic instead of touching the state.
func TestAgsSandboxToolResource_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	sdkClient := &agssdk.Client{}
	patches.ApplyMethodReturn(client, "UseAgsV20250920Client", sdkClient)
	patches.ApplyMethodFunc(sdkClient, "UpdateSandboxToolWithContext",
		func(ctx context.Context, request *agssdk.UpdateSandboxToolRequest) (*agssdk.UpdateSandboxToolResponse, error) {
			return nil, assert.AnError
		})

	r := &AgsSandboxToolResource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &resource.ConfigureResponse{})

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

	objType := tftypes.Object{AttributeTypes: agsSandboxToolAttrTypes()}
	nulls := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		nulls[name] = tftypes.NewValue(attrType, nil)
	}
	nulls["id"] = tftypes.NewValue(tftypes.String, "tool-123")
	nulls["tool_name"] = tftypes.NewValue(tftypes.String, "my-tool")
	nulls["tool_type"] = tftypes.NewValue(tftypes.String, "browser")
	raw := tftypes.NewValue(objType, nulls)

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Raw: raw, Schema: schemaResp.Schema},
		Plan:   tfsdk.Plan{Raw: raw, Schema: schemaResp.Schema},
		State:  tfsdk.State{Raw: raw, Schema: schemaResp.Schema},
	}
	resp := &resource.UpdateResponse{
		State: tfsdk.State{Raw: raw, Schema: schemaResp.Schema},
	}
	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics on update failure")
}

// TestAgsSandboxToolResource_Delete_MockedAPI verifies the delete path sends
// the tool id and reports no diagnostics.
func TestAgsSandboxToolResource_Delete_MockedAPI(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	sdkClient := &agssdk.Client{}
	patches.ApplyMethodReturn(client, "UseAgsV20250920Client", sdkClient)

	var capturedDelete *agssdk.DeleteSandboxToolRequest
	patches.ApplyMethodFunc(sdkClient, "DeleteSandboxToolWithContext",
		func(ctx context.Context, request *agssdk.DeleteSandboxToolRequest) (*agssdk.DeleteSandboxToolResponse, error) {
			capturedDelete = request
			resp := agssdk.NewDeleteSandboxToolResponse()
			resp.Response = &agssdk.DeleteSandboxToolResponseParams{}
			return resp, nil
		})

	r := &AgsSandboxToolResource{}
	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &resource.ConfigureResponse{})

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

	state := agsNullState()
	stateMap, err := state.AsValueMap()
	assert.NoError(t, err)
	stateMap["id"] = tftypes.NewValue(tftypes.String, "tool-123")
	state = tftypes.NewValue(state.Type(), stateMap)

	req := resource.DeleteRequest{State: tfsdk.State{Raw: state, Schema: schemaResp.Schema}}
	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
	assert.NotNil(t, capturedDelete)
	assert.Equal(t, "tool-123", *capturedDelete.ToolId)
}

// TestAgsSandboxToolResource_ImportState verifies the import id wiring.
func TestAgsSandboxToolResource_ImportState(t *testing.T) {
	r := &AgsSandboxToolResource{}
	req := resource.ImportStateRequest{ID: "tool-123"}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	id, err := resp.State.GetAttribute(context.Background(), path.Root("id"))
	assert.NoError(t, err)
	assert.Equal(t, "tool-123", id.(types.String).ValueString())
}

// mustList builds a typed list value in tests, failing the test on error.
func mustList[T any](ctx context.Context, elements []T, elemType attr.Type) types.List {
	list, diags := types.ListValueFrom(ctx, elemType, elements)
	if diags.HasError() {
		panic(diags)
	}
	return list
}

// mustSet builds a typed set value in tests, failing the test on error.
func mustSet(ctx context.Context, values []string) types.Set {
	set, diags := types.SetValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		panic(diags)
	}
	return set
}
