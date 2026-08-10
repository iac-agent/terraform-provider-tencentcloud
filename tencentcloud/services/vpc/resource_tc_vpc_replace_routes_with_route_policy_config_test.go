package vpc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
)

// mockVpcReplaceRoutesMeta implements tccommon.ProviderMeta for vpc replace routes tests.
type mockVpcReplaceRoutesMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockVpcReplaceRoutesMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockVpcReplaceRoutesMeta{}

func newMockVpcReplaceRoutesMeta() *mockVpcReplaceRoutesMeta {
	return &mockVpcReplaceRoutesMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrStringVpcRR(s string) *string { return &s }
func ptrUint64VpcRR(i uint64) *uint64 { return &i }
func ptrBoolVpcRR(b bool) *bool       { return &b }

// go test ./tencentcloud/services/vpc/ -run "TestVpcReplaceRoutesWithRoutePolicyConfig" -v -count=1 -gcflags="all=-l"

// TestVpcReplaceRoutesWithRoutePolicyConfig_Schema validates the new schema fields
func TestVpcReplaceRoutesWithRoutePolicyConfig_Schema(t *testing.T) {
	res := svcvpc.ResourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfig()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Schema)

	assert.Contains(t, res.Schema, "route_table_id")
	assert.Contains(t, res.Schema, "routes")

	// new fields
	needRouterInfo := res.Schema["need_router_info"]
	assert.NotNil(t, needRouterInfo)
	assert.Equal(t, schema.TypeBool, needRouterInfo.Type)
	assert.True(t, needRouterInfo.Optional)
	assert.False(t, needRouterInfo.Required)

	nameField := res.Schema["name"]
	assert.NotNil(t, nameField)
	assert.Equal(t, schema.TypeString, nameField.Type)
	assert.True(t, nameField.Optional)
	assert.False(t, nameField.Required)

	valuesField := res.Schema["values"]
	assert.NotNil(t, valuesField)
	assert.Equal(t, schema.TypeList, valuesField.Type)
	assert.True(t, valuesField.Optional)
	assert.False(t, valuesField.Required)
	assert.NotNil(t, valuesField.Elem)

	totalCountField := res.Schema["total_count"]
	assert.NotNil(t, totalCountField)
	assert.Equal(t, schema.TypeInt, totalCountField.Type)
	assert.True(t, totalCountField.Computed)
	assert.False(t, totalCountField.Required)
	assert.False(t, totalCountField.Optional)
}

// TestVpcReplaceRoutesWithRoutePolicyConfig_Read_Success tests Read retrieves route table and total_count
func TestVpcReplaceRoutesWithRoutePolicyConfig_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockVpcReplaceRoutesMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		resp := vpcv20170312.NewDescribeRouteTablesResponse()
		resp.Response = &vpcv20170312.DescribeRouteTablesResponseParams{
			TotalCount: ptrUint64VpcRR(1),
			RouteTableSet: []*vpcv20170312.RouteTable{
				{
					VpcId:          ptrStringVpcRR("vpc-xxxxxxxx"),
					RouteTableId:   ptrStringVpcRR("rtb-olsbhnyc"),
					RouteTableName: ptrStringVpcRR("test-route-table"),
					Main:           ptrBoolVpcRR(false),
					AssociationSet: []*vpcv20170312.RouteTableAssociation{},
					RouteSet:       []*vpcv20170312.Route{},
					CreatedTime:    ptrStringVpcRR("2024-01-01 00:00:00"),
				},
			},
			RequestId: ptrStringVpcRR("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockVpcReplaceRoutesMeta()
	res := svcvpc.ResourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"route_table_id": "rtb-olsbhnyc",
		"routes":         []interface{}{},
	})
	d.SetId("rtb-olsbhnyc")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "rtb-olsbhnyc", d.Id())
	assert.Equal(t, "rtb-olsbhnyc", d.Get("route_table_id"))
	assert.Equal(t, 1, d.Get("total_count"))
}

// TestVpcReplaceRoutesWithRoutePolicyConfig_Read_WithParams tests Read passes need_router_info, name, values
func TestVpcReplaceRoutesWithRoutePolicyConfig_Read_WithParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockVpcReplaceRoutesMeta().client, "UseVpcClient", vpcClient)

	var capturedNeedRouterInfo *bool
	var capturedFilters []*vpcv20170312.Filter

	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		capturedNeedRouterInfo = request.NeedRouterInfo
		capturedFilters = request.Filters
		resp := vpcv20170312.NewDescribeRouteTablesResponse()
		resp.Response = &vpcv20170312.DescribeRouteTablesResponseParams{
			TotalCount: ptrUint64VpcRR(2),
			RouteTableSet: []*vpcv20170312.RouteTable{
				{
					VpcId:          ptrStringVpcRR("vpc-xxxxxxxx"),
					RouteTableId:   ptrStringVpcRR("rtb-olsbhnyc"),
					RouteTableName: ptrStringVpcRR("test-route-table"),
					Main:           ptrBoolVpcRR(false),
					AssociationSet: []*vpcv20170312.RouteTableAssociation{},
					RouteSet:       []*vpcv20170312.Route{},
					CreatedTime:    ptrStringVpcRR("2024-01-01 00:00:00"),
				},
			},
			RequestId: ptrStringVpcRR("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockVpcReplaceRoutesMeta()
	res := svcvpc.ResourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"route_table_id":   "rtb-olsbhnyc",
		"routes":           []interface{}{},
		"need_router_info": false,
		"name":             "route-table-name",
		"values":           []interface{}{"rtb-olsbhnyc"},
	})
	d.SetId("rtb-olsbhnyc")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "rtb-olsbhnyc", d.Id())
	assert.Equal(t, 2, d.Get("total_count"))

	// verify need_router_info passed through
	assert.NotNil(t, capturedNeedRouterInfo)
	assert.False(t, *capturedNeedRouterInfo)

	// verify name/values assembled into a filter
	assert.NotNil(t, capturedFilters)
	foundNameFilter := false
	for _, f := range capturedFilters {
		if f.Name != nil && *f.Name == "route-table-name" {
			foundNameFilter = true
			assert.Equal(t, []string{"rtb-olsbhnyc"}, []string{*f.Values[0]})
		}
	}
	assert.True(t, foundNameFilter, "expected a filter with name route-table-name")
}

// TestVpcReplaceRoutesWithRoutePolicyConfig_Read_Empty tests Read handles empty response by clearing id
func TestVpcReplaceRoutesWithRoutePolicyConfig_Read_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockVpcReplaceRoutesMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		resp := vpcv20170312.NewDescribeRouteTablesResponse()
		resp.Response = &vpcv20170312.DescribeRouteTablesResponseParams{
			TotalCount:    ptrUint64VpcRR(0),
			RouteTableSet: []*vpcv20170312.RouteTable{},
			RequestId:     ptrStringVpcRR("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockVpcReplaceRoutesMeta()
	res := svcvpc.ResourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"route_table_id": "rtb-olsbhnyc",
		"routes":         []interface{}{},
	})
	d.SetId("rtb-olsbhnyc")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestVpcReplaceRoutesWithRoutePolicyConfig_Read_APIError tests Read handles API error
func TestVpcReplaceRoutesWithRoutePolicyConfig_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockVpcReplaceRoutesMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockVpcReplaceRoutesMeta()
	res := svcvpc.ResourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"route_table_id": "rtb-olsbhnyc",
		"routes":         []interface{}{},
	})
	d.SetId("rtb-olsbhnyc")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestVpcReplaceRoutesWithRoutePolicyConfig_Create_Success tests Create calls ReplaceRoutesWithRoutePolicy and Read
func TestVpcReplaceRoutesWithRoutePolicyConfig_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockVpcReplaceRoutesMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "ReplaceRoutesWithRoutePolicyWithContext", func(ctx interface{}, request *vpcv20170312.ReplaceRoutesWithRoutePolicyRequest) (*vpcv20170312.ReplaceRoutesWithRoutePolicyResponse, error) {
		assert.NotNil(t, request.RouteTableId)
		assert.Equal(t, "rtb-olsbhnyc", *request.RouteTableId)
		assert.NotNil(t, request.Routes)
		resp := vpcv20170312.NewReplaceRoutesWithRoutePolicyResponse()
		resp.Response = &vpcv20170312.ReplaceRoutesWithRoutePolicyResponseParams{
			RequestId: ptrStringVpcRR("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		resp := vpcv20170312.NewDescribeRouteTablesResponse()
		resp.Response = &vpcv20170312.DescribeRouteTablesResponseParams{
			TotalCount: ptrUint64VpcRR(1),
			RouteTableSet: []*vpcv20170312.RouteTable{
				{
					VpcId:          ptrStringVpcRR("vpc-xxxxxxxx"),
					RouteTableId:   ptrStringVpcRR("rtb-olsbhnyc"),
					RouteTableName: ptrStringVpcRR("test-route-table"),
					Main:           ptrBoolVpcRR(false),
					AssociationSet: []*vpcv20170312.RouteTableAssociation{},
					RouteSet:       []*vpcv20170312.Route{},
					CreatedTime:    ptrStringVpcRR("2024-01-01 00:00:00"),
				},
			},
			RequestId: ptrStringVpcRR("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockVpcReplaceRoutesMeta()
	res := svcvpc.ResourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"route_table_id": "rtb-olsbhnyc",
		"routes": []interface{}{
			map[string]interface{}{
				"route_item_id":      "rti-araogi5t",
				"force_match_policy": true,
			},
		},
		"need_router_info": false,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "rtb-olsbhnyc", d.Id())
	assert.Equal(t, 1, d.Get("total_count"))
}
