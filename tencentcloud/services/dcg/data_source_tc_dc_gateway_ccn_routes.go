package dcg

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudDcGatewayCCNRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcGatewayCCNRoutesRead,
		Schema: map[string]*schema.Schema{
			"dcg_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID DCG to be queried。",
			},
			"ccn_route_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud networking routing learning 类型，可选 values: BGP - Automatic Learning; STATIC - 用户 configured. 默认为 STATIC。",
			},
			"address_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "地址 类型，supports: IPv4，IPv6. 默认为 IPv4。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 the DCG route entries。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID DCG。",
						},
						"route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID DCG route。",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A network 地址 segment of IDC。",
						},
						"as_path": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "As 路径 列表 the BGP。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDcGatewayCCNRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dc_gateway_ccn_routes.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		id           string
		ccnRouteType string
		addressType  string
	)

	if v, ok := d.GetOk("dcg_id"); ok {
		id = v.(string)
	}

	if v, ok := d.GetOk("ccn_route_type"); ok {
		ccnRouteType = v.(string)
	}

	if v, ok := d.GetOk("address_type"); ok {
		addressType = v.(string)
	}

	var infos, err = service.DescribeDirectConnectGatewayCcnRoutes(ctx, id, ccnRouteType, addressType)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["dcg_id"] = item.dcgId
		infoMap["route_id"] = item.routeId
		infoMap["cidr_block"] = item.cidrBlock
		infoMap["as_path"] = item.asPaths
		infoList = append(infoList, infoMap)
	}
	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  dcg  ccn routes fail, reason:%s\n ",
			logId,
			err.Error())
		return err
	}

	d.SetId(id)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId,
				output.(string),
				err.Error())
			return err
		}
	}
	return nil

}
