package dcg

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudDcGatewayInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcGatewayInstancesRead,
		Schema: map[string]*schema.Schema{
			"dcg_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID DCG 到 是 queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 DCG 到 是 queried。",
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
				Description: "Information 列表 DCG。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID DCG。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 DCG。",
						},
						"dcg_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 的 DCG。",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 的 DCG。",
						},
						"network_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 associated 网络. 有效值：`VPC` 和 `CCN`。",
						},
						"gateway_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 网关. 有效值：`NORMAL` 和 `NAT`. 默认为 `NORMAL`。",
						},
						"cnn_route_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 CCN 路由. 有效值：`BGP` 和 `STATIC`。",
						},
						"enable_bgp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否BGP 是 已启用",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 资源。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDcGatewayInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dc_gateway_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		id   = ""
		name = ""
	)

	if temp, ok := d.GetOk("dcg_id"); ok {
		if tempStr := temp.(string); tempStr != "" {
			id = tempStr
		}
	}

	if temp, ok := d.GetOk("name"); ok {
		if tempStr := temp.(string); tempStr != "" {
			name = tempStr
		}
	}
	var infos, err = service.DescribeDirectConnectGateways(ctx, id, name)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["dcg_id"] = item.dcgId
		infoMap["name"] = item.name
		infoMap["dcg_ip"] = item.dcgIp
		infoMap["network_type"] = item.networkType
		infoMap["network_instance_id"] = item.networkInstanceId
		infoMap["gateway_type"] = item.gatewayType
		infoMap["cnn_route_type"] = item.cnnRouteType
		infoMap["create_time"] = item.createTime
		infoMap["enable_bgp"] = item.enableBGP
		infoList = append(infoList, infoMap)
	}
	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  dcg instances fail, reason:%s\n ",
			logId,
			err.Error())
		return err
	}

	m := md5.New()
	_, err = m.Write([]byte("dcg_instances" + id + "_" + name))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

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
