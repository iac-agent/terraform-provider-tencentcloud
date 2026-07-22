package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapCheckProxyCreate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapCheckProxyCreateRead,
		Schema: map[string]*schema.Schema{
			"access_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "访问 (acceleration) area 的 proxy. 值 可以 是 获取 through interface DescribeAccessRegionsByDestRegion。",
			},

			"real_server_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "源站 area 的 proxy. 值 可以 是 获取 through interface DescribeDestRegions。",
			},

			"bandwidth": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "upper 限制 的 proxy 带宽，在 Mbps。",
			},

			"concurrent": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "upper 限制 的 chanproxynel 并发，representing 数量 simultaneous online connections，在 tens 的 thousands。",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "如果 creating proxy under proxy 组，您 need 到 fill 在 ID proxy 组。",
			},

			"ip_address_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "IP 版本，可以 是 taken 作为 IPv4 或 IPv6，使用 默认值 的 IPv4。",
			},

			"network_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Network 类型，可以 take 值 &amp;#39;normal&amp;#39;，&amp;#39;cn2&amp;#39;，默认值 normal。",
			},

			"package_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Channel 包 类型 Thunder 表示 standard proxy 组，Accelerator 表示 game accelerator proxy，和 CrossBorder 表示 cross-border proxy。",
			},

			"check_flag": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Query 是否proxy 使用 given 配置 可以 是 创建，1 可以 是 创建，0 不能 是 创建。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudGaapCheckProxyCreateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_check_proxy_create.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("access_region"); ok {
		paramMap["AccessRegion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("real_server_region"); ok {
		paramMap["RealServerRegion"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("bandwidth"); v != nil {
		paramMap["Bandwidth"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("concurrent"); v != nil {
		paramMap["Concurrent"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_address_version"); ok {
		paramMap["IPAddressVersion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok {
		paramMap["NetworkType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		paramMap["PackageType"] = helper.String(v.(string))
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var checkFlag *uint64
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapCheckProxyCreate(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		checkFlag = result
		return nil
	})
	if err != nil {
		return err
	}

	result := map[string]interface{}{}

	if checkFlag != nil {
		_ = d.Set("check_flag", *checkFlag)
		result["check_flag"] = *checkFlag
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
