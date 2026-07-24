package vpc

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// mockMeta implements tccommon.ProviderMeta
type vpcReplaceRoutesMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *vpcReplaceRoutesMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &vpcReplaceRoutesMockMeta{}

func newVpcReplaceRoutesMockMeta() *vpcReplaceRoutesMockMeta {
	return &vpcReplaceRoutesMockMeta{client: &connectivity.TencentCloudClient{}}
}

// buildDescribeRouteTablesResponse builds a minimal successful DescribeRouteTables
// response containing a single route table so that the read helper returns found=true.
func buildDescribeRouteTablesResponse() *vpcv20170312.DescribeRouteTablesResponse {
	resp := vpcv20170312.NewDescribeRouteTablesResponse()
	resp.Response = &vpcv20170312.DescribeRouteTablesResponseParams{
		TotalCount: helper.Uint64(1),
		RouteTableSet: []*vpcv20170312.RouteTable{
			{
				CreatedTime:    helper.String("2024-01-01 00:00:00"),
				Main:           helper.Bool(false),
				RouteTableName: helper.String("rtb-name"),
				RouteTableId:   helper.String("rtb-olsbhnyc"),
				VpcId:          helper.String("vpc-xxxxxxxx"),
				AssociationSet: []*vpcv20170312.RouteTableAssociation{},
				RouteSet:       []*vpcv20170312.Route{},
			},
		},
	}
	return resp
}

// TestReplaceRoutesWithRoutePolicyConfig_ReadWithNameValues verifies that when the
// read helper is invoked with a filter name and values, the DescribeRouteTables
// request carries a Filter{Name, Values}.
func TestReplaceRoutesWithRoutePolicyConfig_ReadWithNameValues(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newVpcReplaceRoutesMockMeta().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcv20170312.DescribeRouteTablesRequest
	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		capturedRequest = request
		return buildDescribeRouteTablesResponse(), nil
	})

	service := VpcService{client: newVpcReplaceRoutesMockMeta().client}
	info, found, err := service.DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig(
		context.Background(),
		"rtb-olsbhnyc",
		"route-table-name",
		[]*string{helper.String("my-rtb")},
		nil,
		nil,
		0,
	)

	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "rtb-olsbhnyc", info.RouteTableId())

	// Verify a Filter{Name, Values} is present alongside the route-table-id filter
	assert.NotNil(t, capturedRequest.Filters)
	assert.Len(t, capturedRequest.Filters, 2)

	var nameFilter *vpcv20170312.Filter
	for _, f := range capturedRequest.Filters {
		if f.Name != nil && *f.Name == "route-table-name" {
			nameFilter = f
		}
	}
	assert.NotNil(t, nameFilter, "expected a Filter with name route-table-name")
	assert.Equal(t, []*string{helper.String("my-rtb")}, nameFilter.Values)
	assert.Nil(t, capturedRequest.RouteTableIds)
}

// TestReplaceRoutesWithRoutePolicyConfig_ReadWithRouteTableIds verifies that when
// route_table_ids is set, the request uses RouteTableIds and NO Filters.
func TestReplaceRoutesWithRoutePolicyConfig_ReadWithRouteTableIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newVpcReplaceRoutesMockMeta().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcv20170312.DescribeRouteTablesRequest
	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		capturedRequest = request
		return buildDescribeRouteTablesResponse(), nil
	})

	service := VpcService{client: newVpcReplaceRoutesMockMeta().client}
	routeTableIds := []*string{helper.String("rtb-aaa"), helper.String("rtb-bbb")}
	_, found, err := service.DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig(
		context.Background(),
		"rtb-olsbhnyc",
		"route-table-name",
		[]*string{helper.String("my-rtb")},
		nil,
		routeTableIds,
		0,
	)

	assert.NoError(t, err)
	assert.True(t, found)

	// Verify RouteTableIds is set and Filters is NOT set (mutually exclusive)
	assert.Equal(t, routeTableIds, capturedRequest.RouteTableIds)
	assert.Nil(t, capturedRequest.Filters)
}

// TestReplaceRoutesWithRoutePolicyConfig_ReadWithNeedRouterInfoFalse verifies that
// when need_router_info = false, request.NeedRouterInfo == false.
func TestReplaceRoutesWithRoutePolicyConfig_ReadWithNeedRouterInfoFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newVpcReplaceRoutesMockMeta().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcv20170312.DescribeRouteTablesRequest
	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		capturedRequest = request
		return buildDescribeRouteTablesResponse(), nil
	})

	service := VpcService{client: newVpcReplaceRoutesMockMeta().client}
	_, found, err := service.DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig(
		context.Background(),
		"rtb-olsbhnyc",
		"",
		nil,
		helper.Bool(false),
		nil,
		0,
	)

	assert.NoError(t, err)
	assert.True(t, found)
	assert.NotNil(t, capturedRequest.NeedRouterInfo)
	assert.False(t, *capturedRequest.NeedRouterInfo)
}

// TestReplaceRoutesWithRoutePolicyConfig_ReadWithLimit verifies that when limit = 50,
// request.Limit == "50".
func TestReplaceRoutesWithRoutePolicyConfig_ReadWithLimit(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newVpcReplaceRoutesMockMeta().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcv20170312.DescribeRouteTablesRequest
	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		capturedRequest = request
		return buildDescribeRouteTablesResponse(), nil
	})

	service := VpcService{client: newVpcReplaceRoutesMockMeta().client}
	_, found, err := service.DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig(
		context.Background(),
		"rtb-olsbhnyc",
		"",
		nil,
		nil,
		nil,
		50,
	)

	assert.NoError(t, err)
	assert.True(t, found)
	assert.NotNil(t, capturedRequest.Limit)
	assert.Equal(t, "50", *capturedRequest.Limit)
}

// TestReplaceRoutesWithRoutePolicyConfig_ReadEmptyResult verifies that an empty
// DescribeRouteTables response results in found=false without error.
func TestReplaceRoutesWithRoutePolicyConfig_ReadEmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newVpcReplaceRoutesMockMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		resp := vpcv20170312.NewDescribeRouteTablesResponse()
		resp.Response = &vpcv20170312.DescribeRouteTablesResponseParams{
			TotalCount:    helper.Uint64(0),
			RouteTableSet: []*vpcv20170312.RouteTable{},
		}
		return resp, nil
	})

	service := VpcService{client: newVpcReplaceRoutesMockMeta().client}
	_, found, err := service.DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig(
		context.Background(),
		"rtb-olsbhnyc",
		"",
		nil,
		nil,
		nil,
		0,
	)

	assert.NoError(t, err)
	assert.False(t, found)
}

// TestReplaceRoutesWithRoutePolicyConfig_ReadWithLimitClamp verifies that a limit
// greater than 100 is clamped to 100.
func TestReplaceRoutesWithRoutePolicyConfig_ReadWithLimitClamp(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newVpcReplaceRoutesMockMeta().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcv20170312.DescribeRouteTablesRequest
	patches.ApplyMethodFunc(vpcClient, "DescribeRouteTables", func(request *vpcv20170312.DescribeRouteTablesRequest) (*vpcv20170312.DescribeRouteTablesResponse, error) {
		capturedRequest = request
		return buildDescribeRouteTablesResponse(), nil
	})

	service := VpcService{client: newVpcReplaceRoutesMockMeta().client}
	_, found, err := service.DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig(
		context.Background(),
		"rtb-olsbhnyc",
		"",
		nil,
		nil,
		nil,
		150,
	)

	assert.NoError(t, err)
	assert.True(t, found)
	assert.NotNil(t, capturedRequest.Limit)
	assert.Equal(t, "100", *capturedRequest.Limit)
}
