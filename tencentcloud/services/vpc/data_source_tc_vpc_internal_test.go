package vpc

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type mockMetaVpcDSInternal struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaVpcDSInternal) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaVpcDSInternal{}

func newMockMetaVpcDSInternal() *mockMetaVpcDSInternal {
	return &mockMetaVpcDSInternal{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/vpc/ -run "TestDataSourceTencentCloudVpc_DomainName" -v -count=1 -gcflags="all=-l"

// TestDataSourceTencentCloudVpc_DomainNameFilled covers the scenario where the queried
// VPC has a DHCP domain name configured: the data source Read MUST populate the
// `domain_name` field with the value returned by DescribeVpcs.
func TestDataSourceTencentCloudVpc_DomainNameFilled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	expectedDomainName := "example.com"
	expectedVpcId := "vpc-abcdefgh"

	patches.ApplyMethodFunc(&VpcService{}, "DescribeVpcs",
		func(_ context.Context, _ string, _ string, _ map[string]string, _ *bool, _ string, _ string) ([]VpcBasicInfo, error) {
			info := VpcBasicInfo{
				vpcId:      expectedVpcId,
				name:       "tf-ci-test",
				cidr:       "10.0.0.0/16",
				isDefault:  false,
				domainName: expectedDomainName,
			}
			return []VpcBasicInfo{info}, nil
		})

	meta := newMockMetaVpcDSInternal()
	res := DataSourceTencentCloudVpc()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.Set("id", expectedVpcId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, expectedVpcId, d.Id())
	assert.Equal(t, expectedDomainName, d.Get("domain_name").(string))
}

// TestDataSourceTencentCloudVpc_DomainNameEmpty covers the scenario where the queried
// VPC has no DHCP domain name configured: DescribeVpcs returns an empty domainName,
// and the data source Read MUST keep `domain_name` as an empty string without panic.
func TestDataSourceTencentCloudVpc_DomainNameEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	expectedVpcId := "vpc-empty-domain"

	patches.ApplyMethodFunc(&VpcService{}, "DescribeVpcs",
		func(_ context.Context, _ string, _ string, _ map[string]string, _ *bool, _ string, _ string) ([]VpcBasicInfo, error) {
			info := VpcBasicInfo{
				vpcId:      expectedVpcId,
				name:       "tf-ci-test",
				cidr:       "10.0.0.0/16",
				isDefault:  false,
				domainName: "",
			}
			return []VpcBasicInfo{info}, nil
		})

	meta := newMockMetaVpcDSInternal()
	res := DataSourceTencentCloudVpc()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.Set("id", expectedVpcId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, expectedVpcId, d.Id())
	assert.Equal(t, "", d.Get("domain_name").(string))
}
