package ags_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

func ptrStrDSAgs(s string) *string    { return &s }
func ptrUint64DSAgs(v uint64) *uint64 { return &v }

// agsGPUResourceSpecsSchema returns the framework data source schema used to
// build the raw config/state value.
func agsGPUResourceSpecsSchema(t *testing.T) schema.Schema {
	d := svcags.NewAgsGPUResourceSpecsDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// agsGPUResourceSpecItemType builds the tftypes.Object type for a specs element.
func agsGPUResourceSpecItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"gpu_type":        tftypes.String,
			"gpu_type_alias":  tftypes.String,
			"gpu_memory_gib":  tftypes.Number,
			"cpu":             tftypes.Number,
			"memory":          tftypes.Number,
			"gpu_num":         tftypes.Number,
			"available_count": tftypes.Number,
		},
	}
}

// agsGPUResourceSpecsRootType builds the tftypes.Object type for the data
// source root, mirroring every schema attribute.
func agsGPUResourceSpecsRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"gpu_type":    tftypes.String,
			"cpu":         tftypes.Number,
			"memory":      tftypes.Number,
			"gpu_num":     tftypes.Number,
			"offset":      tftypes.Number,
			"limit":       tftypes.Number,
			"total_count": tftypes.Number,
			"specs":       tftypes.List{ElementType: agsGPUResourceSpecItemType()},
		},
	}
}

// buildAgsGPUResourceSpecsRawValue builds the raw config value with the given
// query inputs (all inputs are optional).
func buildAgsGPUResourceSpecsRawValue(gpuType string, cpu, memory, gpuNum, offset, limit int64) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"gpu_type":    tftypes.NewValue(tftypes.String, nil),
		"cpu":         tftypes.NewValue(tftypes.Number, nil),
		"memory":      tftypes.NewValue(tftypes.Number, nil),
		"gpu_num":     tftypes.NewValue(tftypes.Number, nil),
		"offset":      tftypes.NewValue(tftypes.Number, nil),
		"limit":       tftypes.NewValue(tftypes.Number, nil),
		"total_count": tftypes.NewValue(tftypes.Number, nil),
		"specs":       tftypes.NewValue(tftypes.List{ElementType: agsGPUResourceSpecItemType()}, nil),
	}

	if gpuType != "" {
		attrs["gpu_type"] = tftypes.NewValue(tftypes.String, gpuType)
	}
	for name, v := range map[string]int64{"cpu": cpu, "memory": memory, "gpu_num": gpuNum, "offset": offset, "limit": limit} {
		if v != -1 {
			attrs[name] = tftypes.NewValue(tftypes.Number, v)
		}
	}

	return tftypes.NewValue(agsGPUResourceSpecsRootType(), attrs)
}

// setAgsGPUResourceSpecsClient injects a mock connectivity client into the data
// source instance via the Configure method.
func setAgsGPUResourceSpecsClient(d *svcags.AgsGPUResourceSpecsDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newAgsGPUResourceSpecsReadRequest builds a ReadRequest with the given config
// value and schema.
func newAgsGPUResourceSpecsReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	return req
}

// newAgsGPUResourceSpecsReadResponse builds a ReadResponse whose State carries
// the schema, so resp.State.Set can reflect the model into a raw value.
func newAgsGPUResourceSpecsReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch
	return resp
}

// go test ./tencentcloud/services/ags/ -run "TestAgsGPUResourceSpecs" -v -count=1 -gcflags="all=-l"

// TestAgsGPUResourceSpecs_MetadataAndSchema validates metadata and schema
// definition.
func TestAgsGPUResourceSpecs_MetadataAndSchema(t *testing.T) {
	d := svcags.NewAgsGPUResourceSpecsDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_ags_gpu_resource_specs", metaResp.TypeName)

	sch := agsGPUResourceSpecsSchema(t)
	attrs := sch.Attributes

	// Optional inputs
	for _, opt := range []string{"gpu_type", "offset", "limit"} {
		attr, ok := attrs[opt].(schema.StringAttribute)
		if ok {
			assert.True(t, attr.Optional, "attribute %s should be Optional", opt)
			continue
		}
		intAttr, ok := attrs[opt].(schema.Int64Attribute)
		assert.True(t, ok, "attribute %s should be StringAttribute or Int64Attribute", opt)
		assert.True(t, intAttr.Optional, "attribute %s should be Optional", opt)
	}
	for _, opt := range []string{"cpu", "memory", "gpu_num"} {
		attr, ok := attrs[opt].(schema.Int64Attribute)
		assert.True(t, ok, "attribute %s should be Int64Attribute", opt)
		assert.True(t, attr.Optional, "attribute %s should be Optional", opt)
	}

	// Computed outputs
	assert.True(t, attrs["total_count"].(schema.Int64Attribute).Computed)
	assert.True(t, attrs["specs"].(schema.ListNestedAttribute).Computed)

	// Nested spec attributes
	nested := attrs["specs"].(schema.ListNestedAttribute).NestedObject.Attributes
	for _, leaf := range []string{"gpu_type", "gpu_type_alias"} {
		nestedAttr, ok := nested[leaf].(schema.StringAttribute)
		assert.True(t, ok, "nested attribute %s should be StringAttribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
	for _, leaf := range []string{"gpu_memory_gib", "cpu", "memory", "gpu_num", "available_count"} {
		nestedAttr, ok := nested[leaf].(schema.Int64Attribute)
		assert.True(t, ok, "nested attribute %s should be Int64Attribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
}

// TestAgsGPUResourceSpecs_ReadBasic tests a successful read with one spec and
// asserts every leaf field is read back.
func TestAgsGPUResourceSpecs_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	patches.ApplyMethodFunc(agsClient, "DescribeGPUResourceSpecsWithContext", func(_ context.Context, request *agsv20250920.DescribeGPUResourceSpecsRequest) (*agsv20250920.DescribeGPUResourceSpecsResponse, error) {
		assert.Equal(t, "5090d", *request.GPUType)
		assert.Equal(t, uint64(64), *request.CPU)
		assert.Equal(t, uint64(192), *request.Memory)
		assert.Equal(t, uint64(2), *request.GPUNum)

		resp := agsv20250920.NewDescribeGPUResourceSpecsResponse()
		resp.Response = &agsv20250920.DescribeGPUResourceSpecsResponseParams{
			GPUResourceSpecSet: []*agsv20250920.GPUResourceSpec{
				{
					GPUType:        ptrStrDSAgs("5090d"),
					GPUTypeAlias:   ptrStrDSAgs("NVIDIA RTX 5090 D"),
					GPUMemoryGiB:   ptrUint64DSAgs(64),
					CPU:            ptrUint64DSAgs(64),
					Memory:         ptrUint64DSAgs(192),
					GPUNum:         ptrUint64DSAgs(2),
					AvailableCount: ptrUint64DSAgs(2),
				},
			},
			TotalCount: ptrInt64DSAgs(1),
			RequestId:  ptrStrDSAgs("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcags.NewAgsGPUResourceSpecsDataSource().(*svcags.AgsGPUResourceSpecsDataSource)
	setAgsGPUResourceSpecsClient(ds, client)

	sch := agsGPUResourceSpecsSchema(t)
	raw := buildAgsGPUResourceSpecsRawValue("5090d", 64, 192, 2, -1, -1)
	req := newAgsGPUResourceSpecsReadRequest(sch, raw)
	resp := newAgsGPUResourceSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcags.AgsGPUResourceSpecsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(1), result.TotalCount.ValueInt64())
	assert.False(t, result.Specs.IsNull())

	var specs []svcags.AgsGPUResourceSpecModel
	resp.Diagnostics.Append(result.Specs.ElementsAs(context.Background(), &specs, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, specs, 1)

	// assert every leaf field
	assert.Equal(t, "5090d", specs[0].GPUType.ValueString())
	assert.Equal(t, "NVIDIA RTX 5090 D", specs[0].GPUTypeAlias.ValueString())
	assert.Equal(t, int64(64), specs[0].GPUMemoryGiB.ValueInt64())
	assert.Equal(t, int64(64), specs[0].CPU.ValueInt64())
	assert.Equal(t, int64(192), specs[0].Memory.ValueInt64())
	assert.Equal(t, int64(2), specs[0].GPUNum.ValueInt64())
	assert.Equal(t, int64(2), specs[0].AvailableCount.ValueInt64())
}

// TestAgsGPUResourceSpecs_ReadEmpty covers a region without sellable specs. An
// empty result is a legitimate answer, not a failure: the data source must
// succeed with total_count = 0 and an empty specs list.
func TestAgsGPUResourceSpecs_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	patches.ApplyMethodFunc(agsClient, "DescribeGPUResourceSpecsWithContext", func(_ context.Context, _ *agsv20250920.DescribeGPUResourceSpecsRequest) (*agsv20250920.DescribeGPUResourceSpecsResponse, error) {
		resp := agsv20250920.NewDescribeGPUResourceSpecsResponse()
		resp.Response = &agsv20250920.DescribeGPUResourceSpecsResponseParams{
			GPUResourceSpecSet: []*agsv20250920.GPUResourceSpec{},
			TotalCount:         ptrInt64DSAgs(0),
			RequestId:          ptrStrDSAgs("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcags.NewAgsGPUResourceSpecsDataSource().(*svcags.AgsGPUResourceSpecsDataSource)
	setAgsGPUResourceSpecsClient(ds, client)

	sch := agsGPUResourceSpecsSchema(t)
	raw := buildAgsGPUResourceSpecsRawValue("", -1, -1, -1, -1, -1)
	req := newAgsGPUResourceSpecsReadRequest(sch, raw)
	resp := newAgsGPUResourceSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "a region without sellable specs must not be reported as an error")

	var result svcags.AgsGPUResourceSpecsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(0), result.TotalCount.ValueInt64())
	assert.False(t, result.Specs.IsNull())
	assert.Equal(t, 0, len(result.Specs.Elements()))
}

// TestAgsGPUResourceSpecs_ReadNilAvailableCount covers a spec whose
// AvailableCount is returned as null (quota unknown) by the API. It must be
// read back as null, not zero.
func TestAgsGPUResourceSpecs_ReadNilAvailableCount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	patches.ApplyMethodFunc(agsClient, "DescribeGPUResourceSpecsWithContext", func(_ context.Context, _ *agsv20250920.DescribeGPUResourceSpecsRequest) (*agsv20250920.DescribeGPUResourceSpecsResponse, error) {
		resp := agsv20250920.NewDescribeGPUResourceSpecsResponse()
		resp.Response = &agsv20250920.DescribeGPUResourceSpecsResponseParams{
			GPUResourceSpecSet: []*agsv20250920.GPUResourceSpec{
				{
					GPUType:        ptrStrDSAgs("5090d"),
					GPUTypeAlias:   ptrStrDSAgs("NVIDIA RTX 5090 D"),
					GPUMemoryGiB:   ptrUint64DSAgs(64),
					CPU:            ptrUint64DSAgs(64),
					Memory:         ptrUint64DSAgs(192),
					GPUNum:         ptrUint64DSAgs(2),
					AvailableCount: nil,
				},
			},
			TotalCount: ptrInt64DSAgs(1),
			RequestId:  ptrStrDSAgs("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcags.NewAgsGPUResourceSpecsDataSource().(*svcags.AgsGPUResourceSpecsDataSource)
	setAgsGPUResourceSpecsClient(ds, client)

	sch := agsGPUResourceSpecsSchema(t)
	raw := buildAgsGPUResourceSpecsRawValue("", -1, -1, -1, -1, -1)
	req := newAgsGPUResourceSpecsReadRequest(sch, raw)
	resp := newAgsGPUResourceSpecsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	var result svcags.AgsGPUResourceSpecsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	var specs []svcags.AgsGPUResourceSpecModel
	resp.Diagnostics.Append(result.Specs.ElementsAs(context.Background(), &specs, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, specs, 1)

	assert.True(t, specs[0].AvailableCount.IsNull(), "available_count must stay null when the quota is unknown")
}

func ptrInt64DSAgs(v int64) *int64 { return &v }
