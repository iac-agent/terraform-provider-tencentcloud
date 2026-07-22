package vpc

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcRouteTablesRead,

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID routing 表 到 是 queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 routing 表 到 是 queried。",
			},
			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 如果 routing 表 has 此 标签",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID VPC 到 是 queried。",
			},
			"association_main": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "过滤器 main routing 表。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 routing 表 到 是 queried。",
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
				Description: "信息 列表 VPC 路由 表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID routing 表。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 routing 表。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC。",
						},
						"subnet_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "列表 子网 IDs bound 到 路由 表。",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否it 是 默认值 routing 表。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 routing 表。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 routing 表。",
						},
						"route_entry_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed 信息 的 each entry 的 路由 表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"route_entry_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID 路由 表 entry。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述 信息 用户 defined 对于 路由 表 规则。",
									},
									"destination_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "destination 地址 block。",
									},
									"next_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 next-hop，和 可用 值 include `CVM`，`VPN`，`DIRECTCONNECT`，`PEERCONNECTION`，`SSLVPN`，`NAT`，`NORMAL_CVM`，`EIP` 和 `CCN`。",
									},
									"next_hub": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID next-hop 网关. 注意: 当 'next_type' 是 EIP，GatewayId 将 fix 值 `0`。",
									},
									"route_item_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "唯一 策略 ID 对于 路由。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpcRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_route_tables.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region

	var (
		routeTableId    string
		vpcId           string
		name            string
		associationMain *bool
		tagKey          string
	)
	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcId = temp.(string)
	}

	if temp, ok := d.GetOk("route_table_id"); ok {
		routeTableId = temp.(string)
	}

	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}

	if temp, ok := d.GetOkExists("association_main"); ok {
		associationMain = helper.Bool(temp.(bool))
	}

	if temp, ok := d.GetOk("tag_key"); ok {
		tagKey = temp.(string)
	}

	var (
		tags  = helper.GetTags(d, "tags")
		infos []VpcRouteTableBasicInfo
		err   error
	)

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, err = service.DescribeRouteTables(ctx, routeTableId, name, vpcId, tags, associationMain, tagKey)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		routeEntryInfos := make([]map[string]string, len(item.entryInfos))

		for _, v := range item.entryInfos {
			routeEntryInfo := make(map[string]string)
			routeEntryInfo["route_entry_id"] = fmt.Sprintf("%d.%s",
				v.routeEntryId, item.routeTableId)
			routeEntryInfo["description"] = v.description
			routeEntryInfo["destination_cidr_block"] = v.destinationCidr
			routeEntryInfo["next_type"] = v.nextType
			routeEntryInfo["next_hub"] = v.nextBub
			routeEntryInfo["route_item_id"] = v.routeItemId
			routeEntryInfos = append(routeEntryInfos, routeEntryInfo)
		}

		respTags, err := tagService.DescribeResourceTags(ctx, "vpc", "rtb", region, item.routeTableId)
		if err != nil {
			return err
		}

		var infoMap = make(map[string]interface{})

		infoMap["route_table_id"] = item.routeTableId
		infoMap["name"] = item.name
		infoMap["vpc_id"] = item.vpcId
		infoMap["is_default"] = item.isDefault
		infoMap["subnet_ids"] = item.subnetIds
		infoMap["route_entry_infos"] = routeEntryInfos
		infoMap["create_time"] = item.createTime
		infoMap["tags"] = respTags

		infoList = append(infoList, infoMap)
	}

	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  route table instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	idBytes, err := json.Marshal(map[string]interface{}{
		"routeTableId":    routeTableId,
		"associationMain": associationMain,
		"vpcId":           vpcId,
		"name":            name,
		"tagKey":          tagKey,
		"tags":            tags,
	})
	if err != nil {
		log.Printf("[CRITAL]%s create data source id error, reason:%s\n ", logId, err.Error())
		return err
	}

	md := md5.New()
	_, _ = md.Write(idBytes)
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
