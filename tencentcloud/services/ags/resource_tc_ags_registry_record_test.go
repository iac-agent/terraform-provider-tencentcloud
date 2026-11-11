package ags_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

func agsPtrString(s string) *string { return &s }
func agsPtrInt64(i int64) *int64    { return &i }

// ---- tftypes builders mirroring the framework schema ----

var agsApprovalActionTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"action_id":             tftypes.String,
	"action_type":           tftypes.String,
	"actor_type":            tftypes.String,
	"actor_uin":             tftypes.String,
	"actor_sub_account_uin": tftypes.String,
	"comment":               tftypes.String,
	"create_time":           tftypes.String,
	"request_id":            tftypes.String,
}}

var agsHarnessVersionTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"harness_engine":         tftypes.String,
	"harness_engine_version": tftypes.String,
}}

var agsMetadataEntryTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"name":  tftypes.String,
	"value": tftypes.String,
}}

var agsRecordTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"record_id":               tftypes.String,
	"registry_id":             tftypes.String,
	"name":                    tftypes.String,
	"description":             tftypes.String,
	"descriptor_type":         tftypes.String,
	"lifecycle_status":        tftypes.String,
	"app_id":                  tftypes.Number,
	"creator_uin":             tftypes.String,
	"creator_sub_account_uin": tftypes.String,
	"create_time":             tftypes.String,
	"update_time":             tftypes.String,
	"version_count":           tftypes.Number,
	"label_set":               tftypes.Set{ElementType: tftypes.String},
}}

var agsVersionTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"version_id":                 tftypes.String,
	"record_id":                  tftypes.String,
	"revision":                   tftypes.Number,
	"status":                     tftypes.String,
	"approval_mode":              tftypes.String,
	"app_id":                     tftypes.Number,
	"creator_uin":                tftypes.String,
	"create_time":                tftypes.String,
	"update_time":                tftypes.String,
	"version_name":               tftypes.String,
	"descriptors":                tftypes.String,
	"source_type":                tftypes.String,
	"source_config":              tftypes.String,
	"content_status":             tftypes.String,
	"content_sha256":             tftypes.String,
	"content_size_bytes":         tftypes.Number,
	"config_sha256":              tftypes.String,
	"creator_sub_account_uin":    tftypes.String,
	"approval_actions":           tftypes.List{ElementType: agsApprovalActionTfType},
	"content_ready_time":         tftypes.String,
	"change_log":                 tftypes.String,
	"label_set":                  tftypes.Set{ElementType: tftypes.String},
	"supported_harness_versions": tftypes.List{ElementType: agsHarnessVersionTfType},
	"metadata":                   tftypes.List{ElementType: agsMetadataEntryTfType},
	"name":                       tftypes.String,
	"description":                tftypes.String,
	"descriptor_type":            tftypes.String,
}}

var agsMCPSourceTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"type":         tftypes.String,
	"descriptors":  tftypes.String,
	"endpoint_url": tftypes.String,
}}

var agsAgentSourceTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"type":         tftypes.String,
	"descriptors":  tftypes.String,
	"endpoint_url": tftypes.String,
}}

var agsSkillSourceTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"type":     tftypes.String,
	"skill_md": tftypes.String,
}}

var agsLabelMutationTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"operation":  tftypes.String,
	"name":       tftypes.String,
	"version_id": tftypes.String,
	"reason":     tftypes.String,
}}

func agsRegistryRecordTfType() tftypes.Type {
	return tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":                 tftypes.String,
		"registry_id":        tftypes.String,
		"name":               tftypes.String,
		"descriptor_type":    tftypes.String,
		"description":        tftypes.String,
		"version_name":       tftypes.String,
		"mcp_source":         agsMCPSourceTfType,
		"agent_source":       agsAgentSourceTfType,
		"skill_source":       agsSkillSourceTfType,
		"custom_descriptors": tftypes.String,
		"label_mutations":    tftypes.List{ElementType: agsLabelMutationTfType},
		"reason":             tftypes.String,
		"record_id":          tftypes.String,
		"record":             agsRecordTfType,
		"version":            agsVersionTfType,
		"resolved_by":        tftypes.String,
		"resolved_label":     tftypes.String,
	}}
}

// agsNullValue returns the null value for any type in the schema tree.
func agsNullValue(t tftypes.Type) tftypes.Value {
	return tftypes.NewValue(t, nil)
}

// agsNullRaw builds a full null top-level value matching the resource schema.
func agsNullRaw() tftypes.Value {
	objType := agsRegistryRecordTfType().(tftypes.Object)
	attrs := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		attrs[name] = agsNullValue(typ)
	}
	return tftypes.NewValue(objType, attrs)
}

// agsRawWith returns the null raw value with the given top-level attributes
// replaced.
func agsRawWith(overrides map[string]tftypes.Value) tftypes.Value {
	objType := agsRegistryRecordTfType().(tftypes.Object)
	attrs := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		if v, ok := overrides[name]; ok {
			attrs[name] = v
			continue
		}
		attrs[name] = agsNullValue(typ)
	}
	return tftypes.NewValue(objType, attrs)
}

func agsResourceSchema(t *testing.T) tfsdk.Schema {
	res := ags.NewAgsRegistryRecordResource().(*ags.AgsRegistryRecordResource)
	schemaResp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

func agsSetClient(res *ags.AgsRegistryRecordResource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	res.Configure(context.Background(), configureReq, configureResp)
}

func agsAttr(t *testing.T, raw tftypes.Value, name string) tftypes.Value {
	v, err := raw.Attribute(tftypes.NewAttributePath().WithAttributeName(name))
	if err != nil {
		t.Fatalf("failed to get attribute %s: %v", name, err)
	}
	return v
}

func agsAssertStringAttr(t *testing.T, raw tftypes.Value, name, want string) {
	v := agsAttr(t, raw, name)
	var s string
	if err := v.As(&s); err != nil {
		t.Fatalf("attribute %s is not a string: %v", name, err)
	}
	assert.Equal(t, want, s)
}

func agsMockDescribeResponse() *agsv20250920.DescribeRegistryRecordResponse {
	resp := agsv20250920.NewDescribeRegistryRecordResponse()
	resp.Response = &agsv20250920.DescribeRegistryRecordResponseParams{
		Record: &agsv20250920.CloudRecord{
			RecordId:             agsPtrString("rec-12345678"),
			RegistryId:           agsPtrString("reg-abcdef01"),
			Name:                 agsPtrString("my-record"),
			Description:          agsPtrString("record description"),
			DescriptorType:       agsPtrString("MCP"),
			LifecycleStatus:      agsPtrString("ACTIVE"),
			AppId:                agsPtrInt64(1250000000),
			CreatorUin:           agsPtrString("100000000001"),
			CreatorSubAccountUin: agsPtrString(""),
			CreateTime:           agsPtrString("2025-01-01T00:00:00Z"),
			UpdateTime:           agsPtrString("2025-01-02T00:00:00Z"),
			VersionCount:         agsPtrInt64(1),
			LabelSet:             []*string{agsPtrString("stable"), agsPtrString("latest")},
		},
		Version: &agsv20250920.CloudRecordVersion{
			VersionId:            agsPtrString("rv-87654321"),
			RecordId:             agsPtrString("rec-12345678"),
			Revision:             agsPtrInt64(1),
			Status:               agsPtrString("READY"),
			ApprovalMode:         agsPtrString("AUTO"),
			AppId:                agsPtrInt64(1250000000),
			CreatorUin:           agsPtrString("100000000001"),
			CreateTime:           agsPtrString("2025-01-01T00:00:00Z"),
			UpdateTime:           agsPtrString("2025-01-02T00:00:00Z"),
			VersionName:          agsPtrString("v1"),
			Descriptors:          agsPtrString("{\"name\":\"mcp-server\"}"),
			SourceType:           agsPtrString("MANUAL"),
			SourceConfig:         agsPtrString("{}"),
			ContentStatus:        agsPtrString("READY"),
			ContentSHA256:        agsPtrString("abc123"),
			ContentSizeBytes:     agsPtrInt64(1024),
			ConfigSHA256:         agsPtrString("def456"),
			CreatorSubAccountUin: agsPtrString(""),
			ApprovalActions: []*agsv20250920.CloudVersionApprovalAction{
				{
					ActionId:           agsPtrString("act-1"),
					ActionType:         agsPtrString("APPROVE"),
					ActorType:          agsPtrString("SYSTEM"),
					ActorUin:           agsPtrString("100000000001"),
					ActorSubAccountUin: agsPtrString(""),
					Comment:            agsPtrString("auto approved"),
					CreateTime:         agsPtrString("2025-01-01T00:00:00Z"),
					RequestId:          agsPtrString("req-1"),
				},
			},
			ContentReadyTime: agsPtrString("2025-01-01T00:00:00Z"),
			ChangeLog:        agsPtrString(""),
			LabelSet:         []*string{agsPtrString("stable")},
			SupportedHarnessVersions: []*agsv20250920.HarnessVersionReference{
				{
					HarnessEngine:        agsPtrString("pi"),
					HarnessEngineVersion: agsPtrString("1.3.0"),
				},
			},
			Metadata: []*agsv20250920.MetadataEntry{
				{Name: agsPtrString("publisher"), Value: agsPtrString("tencent")},
			},
		},
		ResolvedBy:    agsPtrString("DEFAULT_STABLE"),
		ResolvedLabel: agsPtrString("stable"),
	}
	return resp
}

// go test ./tencentcloud/services/ags/ -run "TestAgsRegistryRecordResource" -v -count=1 -gcflags="all=-l"

func TestAgsRegistryRecordResource_Schema(t *testing.T) {
	s := agsResourceSchema(t)
	for _, name := range []string{
		"id", "registry_id", "name", "descriptor_type", "description",
		"version_name", "mcp_source", "agent_source", "skill_source",
		"custom_descriptors", "label_mutations", "reason", "record_id",
		"record", "version", "resolved_by", "resolved_label",
	} {
		assert.Contains(t, s.Attributes, name)
	}
	assert.True(t, s.Attributes["registry_id"].IsRequired())
	assert.True(t, s.Attributes["name"].IsRequired())
	assert.True(t, s.Attributes["descriptor_type"].IsRequired())
	assert.True(t, s.Attributes["description"].IsOptional())
	assert.True(t, s.Attributes["description"].IsComputed())
	assert.True(t, s.Attributes["record"].IsComputed())
	assert.True(t, s.Attributes["version"].IsComputed())
	assert.True(t, s.Attributes["record_id"].IsComputed())
}

func TestAgsRegistryRecordResource_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	var gotCreate *agsv20250920.CreateRegistryRecordRequest
	patches.ApplyMethodFunc(agsClient, "CreateRegistryRecordWithContext", func(_ context.Context, request *agsv20250920.CreateRegistryRecordRequest) (*agsv20250920.CreateRegistryRecordResponse, error) {
		gotCreate = request
		resp := agsv20250920.NewCreateRegistryRecordResponse()
		resp.Response = &agsv20250920.CreateRegistryRecordResponseParams{
			RecordId: agsPtrString("rec-12345678"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(agsClient, "DescribeRegistryRecordWithContext", func(_ context.Context, _ *agsv20250920.DescribeRegistryRecordRequest) (*agsv20250920.DescribeRegistryRecordResponse, error) {
		return agsMockDescribeResponse(), nil
	})

	res := ags.NewAgsRegistryRecordResource().(*ags.AgsRegistryRecordResource)
	agsSetClient(res, client)

	s := agsResourceSchema(t)
	raw := agsRawWith(map[string]tftypes.Value{
		"registry_id":     tftypes.NewValue(tftypes.String, "reg-abcdef01"),
		"name":            tftypes.NewValue(tftypes.String, "my-record"),
		"descriptor_type": tftypes.NewValue(tftypes.String, "MCP"),
		"description":     tftypes.NewValue(tftypes.String, "record description"),
		"mcp_source": tftypes.NewValue(agsMCPSourceTfType, map[string]tftypes.Value{
			"type":         tftypes.NewValue(tftypes.String, "MANUAL"),
			"descriptors":  tftypes.NewValue(tftypes.String, "{\"name\":\"mcp-server\"}"),
			"endpoint_url": agsNullValue(tftypes.String),
		}),
	})

	createReq := resource.CreateRequest{Plan: tfsdk.Plan{Schema: s, Raw: raw}}
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: s, Raw: agsNullRaw()}}
	res.Create(context.Background(), createReq, createResp)

	assert.False(t, createResp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", createResp.Diagnostics)

	assert.NotNil(t, gotCreate)
	assert.Equal(t, "reg-abcdef01", *gotCreate.RegistryId)
	assert.Equal(t, "my-record", *gotCreate.Name)
	assert.Equal(t, "MCP", *gotCreate.DescriptorType)
	assert.Equal(t, "record description", *gotCreate.Description)
	assert.NotNil(t, gotCreate.MCPSource)
	assert.Equal(t, "MANUAL", *gotCreate.MCPSource.Type)
	assert.Nil(t, gotCreate.AgentSource)
	assert.Nil(t, gotCreate.SkillSource)
	assert.Nil(t, gotCreate.CustomDescriptors)

	agsAssertStringAttr(t, createResp.State.Raw, "id", "reg-abcdef01#rec-12345678")
	agsAssertStringAttr(t, createResp.State.Raw, "record_id", "rec-12345678")
	agsAssertStringAttr(t, createResp.State.Raw, "resolved_by", "DEFAULT_STABLE")

	recordVal := agsAttr(t, createResp.State.Raw, "record")
	recordIdVal, err := recordVal.Attribute(tftypes.NewAttributePath().WithAttributeName("record_id"))
	assert.NoError(t, err)
	var recordId string
	assert.NoError(t, recordIdVal.As(&recordId))
	assert.Equal(t, "rec-12345678", recordId)

	labelSetVal, err := recordVal.Attribute(tftypes.NewAttributePath().WithAttributeName("label_set"))
	assert.NoError(t, err)
	var labelSet []string
	assert.NoError(t, labelSetVal.As(&labelSet))
	assert.ElementsMatch(t, []string{"stable", "latest"}, labelSet)

	versionVal := agsAttr(t, createResp.State.Raw, "version")
	versionIdVal, err := versionVal.Attribute(tftypes.NewAttributePath().WithAttributeName("version_id"))
	assert.NoError(t, err)
	var versionId string
	assert.NoError(t, versionIdVal.As(&versionId))
	assert.Equal(t, "rv-87654321", versionId)
}

func TestAgsRegistryRecordResource_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "CreateRegistryRecordWithContext", func(_ context.Context, _ *agsv20250920.CreateRegistryRecordRequest) (*agsv20250920.CreateRegistryRecordResponse, error) {
		return nil, fmt.Errorf("InternalError")
	})

	res := ags.NewAgsRegistryRecordResource().(*ags.AgsRegistryRecordResource)
	agsSetClient(res, client)

	s := agsResourceSchema(t)
	raw := agsRawWith(map[string]tftypes.Value{
		"registry_id":     tftypes.NewValue(tftypes.String, "reg-abcdef01"),
		"name":            tftypes.NewValue(tftypes.String, "my-record"),
		"descriptor_type": tftypes.NewValue(tftypes.String, "MCP"),
	})

	createReq := resource.CreateRequest{Plan: tfsdk.Plan{Schema: s, Raw: raw}}
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: s, Raw: agsNullRaw()}}
	res.Create(context.Background(), createReq, createResp)

	assert.True(t, createResp.Diagnostics.HasError())
}

func TestAgsRegistryRecordResource_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	var gotDescribe *agsv20250920.DescribeRegistryRecordRequest
	patches.ApplyMethodFunc(agsClient, "DescribeRegistryRecordWithContext", func(_ context.Context, request *agsv20250920.DescribeRegistryRecordRequest) (*agsv20250920.DescribeRegistryRecordResponse, error) {
		gotDescribe = request
		resp := agsv20250920.NewDescribeRegistryRecordResponse()
		resp.Response = &agsv20250920.DescribeRegistryRecordResponseParams{}
		return resp, nil
	})

	res := ags.NewAgsRegistryRecordResource().(*ags.AgsRegistryRecordResource)
	agsSetClient(res, client)

	s := agsResourceSchema(t)
	raw := agsRawWith(map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, "reg-abcdef01#rec-12345678"),
		"registry_id": tftypes.NewValue(tftypes.String, "reg-abcdef01"),
		"record_id":   tftypes.NewValue(tftypes.String, "rec-12345678"),
	})

	readReq := resource.ReadRequest{State: tfsdk.State{Schema: s, Raw: raw}}
	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: s, Raw: raw}}
	res.Read(context.Background(), readReq, readResp)

	assert.False(t, readResp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", readResp.Diagnostics)
	assert.NotNil(t, gotDescribe)
	assert.Equal(t, "reg-abcdef01", *gotDescribe.RegistryId)
	assert.Equal(t, "rec-12345678", *gotDescribe.RecordId)
	assert.Nil(t, gotDescribe.VersionId)
	assert.Nil(t, gotDescribe.Label)
	assert.True(t, readResp.State.Raw.IsNull(), "expected state to be removed")
}

func TestAgsRegistryRecordResource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	patches.ApplyMethodFunc(agsClient, "DescribeRegistryRecordWithContext", func(_ context.Context, _ *agsv20250920.DescribeRegistryRecordRequest) (*agsv20250920.DescribeRegistryRecordResponse, error) {
		return agsMockDescribeResponse(), nil
	})

	res := ags.NewAgsRegistryRecordResource().(*ags.AgsRegistryRecordResource)
	agsSetClient(res, client)

	s := agsResourceSchema(t)
	raw := agsRawWith(map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, "reg-abcdef01#rec-12345678"),
		"registry_id":     tftypes.NewValue(tftypes.String, "reg-abcdef01"),
		"record_id":       tftypes.NewValue(tftypes.String, "rec-12345678"),
		"name":            tftypes.NewValue(tftypes.String, "my-record"),
		"descriptor_type": tftypes.NewValue(tftypes.String, "MCP"),
		"mcp_source": tftypes.NewValue(agsMCPSourceTfType, map[string]tftypes.Value{
			"type":         tftypes.NewValue(tftypes.String, "MANUAL"),
			"descriptors":  tftypes.NewValue(tftypes.String, "{\"name\":\"mcp-server\"}"),
			"endpoint_url": agsNullValue(tftypes.String),
		}),
		"version_name": tftypes.NewValue(tftypes.String, "v1"),
	})

	readReq := resource.ReadRequest{State: tfsdk.State{Schema: s, Raw: raw}}
	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: s, Raw: raw}}
	res.Read(context.Background(), readReq, readResp)

	assert.False(t, readResp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", readResp.Diagnostics)

	agsAssertStringAttr(t, readResp.State.Raw, "name", "my-record")
	agsAssertStringAttr(t, readResp.State.Raw, "description", "record description")
	agsAssertStringAttr(t, readResp.State.Raw, "record_id", "rec-12345678")
	agsAssertStringAttr(t, readResp.State.Raw, "resolved_label", "stable")
	// Configuration-only fields keep their state values even without readback.
	assert.False(t, agsAttr(t, readResp.State.Raw, "mcp_source").IsNull())
}

func TestAgsRegistryRecordResource_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	var gotUpdate *agsv20250920.UpdateRegistryRecordRequest
	patches.ApplyMethodFunc(agsClient, "UpdateRegistryRecordWithContext", func(_ context.Context, request *agsv20250920.UpdateRegistryRecordRequest) (*agsv20250920.UpdateRegistryRecordResponse, error) {
		gotUpdate = request
		resp := agsv20250920.NewUpdateRegistryRecordResponse()
		resp.Response = &agsv20250920.UpdateRegistryRecordResponseParams{}
		return resp, nil
	})
	patches.ApplyMethodFunc(agsClient, "DescribeRegistryRecordWithContext", func(_ context.Context, _ *agsv20250920.DescribeRegistryRecordRequest) (*agsv20250920.DescribeRegistryRecordResponse, error) {
		return agsMockDescribeResponse(), nil
	})

	res := ags.NewAgsRegistryRecordResource().(*ags.AgsRegistryRecordResource)
	agsSetClient(res, client)

	s := agsResourceSchema(t)
	common := map[string]tftypes.Value{
		"id":              tftypes.NewValue(tftypes.String, "reg-abcdef01#rec-12345678"),
		"registry_id":     tftypes.NewValue(tftypes.String, "reg-abcdef01"),
		"record_id":       tftypes.NewValue(tftypes.String, "rec-12345678"),
		"name":            tftypes.NewValue(tftypes.String, "my-record"),
		"descriptor_type": tftypes.NewValue(tftypes.String, "MCP"),
	}
	planRaw := agsRawWith(map[string]tftypes.Value{
		"id":              common["id"],
		"registry_id":     common["registry_id"],
		"record_id":       common["record_id"],
		"name":            common["name"],
		"descriptor_type": common["descriptor_type"],
		"description":     tftypes.NewValue(tftypes.String, "new description"),
		"label_mutations": tftypes.NewValue(tftypes.List{ElementType: agsLabelMutationTfType}, []tftypes.Value{
			tftypes.NewValue(agsLabelMutationTfType, map[string]tftypes.Value{
				"operation":  tftypes.NewValue(tftypes.String, "SET"),
				"name":       tftypes.NewValue(tftypes.String, "canary"),
				"version_id": tftypes.NewValue(tftypes.String, "rv-87654321"),
				"reason":     agsNullValue(tftypes.String),
			}),
		}),
	})
	stateRaw := agsRawWith(map[string]tftypes.Value{
		"id":              common["id"],
		"registry_id":     common["registry_id"],
		"record_id":       common["record_id"],
		"name":            common["name"],
		"descriptor_type": common["descriptor_type"],
		"description":     tftypes.NewValue(tftypes.String, "record description"),
	})

	updateReq := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: s, Raw: planRaw},
		State: tfsdk.State{Schema: s, Raw: stateRaw},
	}
	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: s, Raw: stateRaw}}
	res.Update(context.Background(), updateReq, updateResp)

	assert.False(t, updateResp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", updateResp.Diagnostics)

	assert.NotNil(t, gotUpdate)
	assert.Equal(t, "reg-abcdef01", *gotUpdate.RegistryId)
	assert.Equal(t, "rec-12345678", *gotUpdate.RecordId)
	assert.Equal(t, "new description", *gotUpdate.Description)
	assert.Len(t, gotUpdate.LabelMutations, 1)
	assert.Equal(t, "SET", *gotUpdate.LabelMutations[0].Operation)
	assert.Equal(t, "canary", *gotUpdate.LabelMutations[0].Name)
	assert.Equal(t, "rv-87654321", *gotUpdate.LabelMutations[0].VersionId)
	assert.Nil(t, gotUpdate.VersionName)
	assert.Nil(t, gotUpdate.MCPSource)

	agsAssertStringAttr(t, updateResp.State.Raw, "description", "record description")
}

func TestAgsRegistryRecordResource_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	var gotDelete *agsv20250920.DeleteRegistryRecordRequest
	patches.ApplyMethodFunc(agsClient, "DeleteRegistryRecordWithContext", func(_ context.Context, request *agsv20250920.DeleteRegistryRecordRequest) (*agsv20250920.DeleteRegistryRecordResponse, error) {
		gotDelete = request
		resp := agsv20250920.NewDeleteRegistryRecordResponse()
		resp.Response = &agsv20250920.DeleteRegistryRecordResponseParams{}
		return resp, nil
	})

	res := ags.NewAgsRegistryRecordResource().(*ags.AgsRegistryRecordResource)
	agsSetClient(res, client)

	s := agsResourceSchema(t)
	raw := agsRawWith(map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, "reg-abcdef01#rec-12345678"),
		"registry_id": tftypes.NewValue(tftypes.String, "reg-abcdef01"),
		"record_id":   tftypes.NewValue(tftypes.String, "rec-12345678"),
		"reason":      tftypes.NewValue(tftypes.String, "cleanup"),
	})

	deleteReq := resource.DeleteRequest{State: tfsdk.State{Schema: s, Raw: raw}}
	deleteResp := &resource.DeleteResponse{}
	res.Delete(context.Background(), deleteReq, deleteResp)

	assert.False(t, deleteResp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", deleteResp.Diagnostics)

	assert.NotNil(t, gotDelete)
	assert.Equal(t, "reg-abcdef01", *gotDelete.RegistryId)
	assert.Equal(t, "rec-12345678", *gotDelete.RecordId)
	assert.Equal(t, "cleanup", *gotDelete.Reason)
	// VersionId must be omitted so the whole record is deleted.
	assert.Nil(t, gotDelete.VersionId)
}

func TestAgsRegistryRecordResource_ImportState(t *testing.T) {
	res := ags.NewAgsRegistryRecordResource().(*ags.AgsRegistryRecordResource)
	s := agsResourceSchema(t)

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: s, Raw: agsNullRaw()}}
	res.ImportState(context.Background(), resource.ImportStateRequest{ID: "reg-abcdef01#rec-12345678"}, importResp)

	assert.False(t, importResp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", importResp.Diagnostics)
	agsAssertStringAttr(t, importResp.State.Raw, "id", "reg-abcdef01#rec-12345678")
	agsAssertStringAttr(t, importResp.State.Raw, "registry_id", "reg-abcdef01")
	agsAssertStringAttr(t, importResp.State.Raw, "record_id", "rec-12345678")

	badResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: s, Raw: agsNullRaw()}}
	res.ImportState(context.Background(), resource.ImportStateRequest{ID: "invalid-id"}, badResp)
	assert.True(t, badResp.Diagnostics.HasError())
}
