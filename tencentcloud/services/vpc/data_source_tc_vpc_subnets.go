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

func DataSourceTencentCloudVpcSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSubnetsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID VPC 到 是 queried。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 子网 到 是 queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 子网 到 是 queried。",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "可用区 的 子网 到 是 queried。",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "过滤器 默认值 或 无 默认值 subnets。",
			},
			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 如果 子网 has 此 标签",
			},
			"cidr_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 子网 使用 此 CIDR。",
			},
			"cdc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID CDC 实例。",
			},
			"is_remote_vpc_snat": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "过滤器 VPC SNAT 地址 池 子网。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 子网 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"instance_list": {Type: schema.TypeList,
				Computed:    true,
				Description: "列表 subnets。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "availability 可用区 的 子网。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 子网。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 子网。",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A 网络 地址 block 的 子网。",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CDC 实例。",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否it 是 默认值 子网 的 VPC 对于 此 地域",
						},
						"is_multicast": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否multicast 是 已启用",
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID routing 表。",
						},
						"available_ip_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 可用 IPs。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 子网 资源。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 子网 资源。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpcSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_subnets.read")()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		vpcService       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		tagService       = svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region           = meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		vpcId            string
		subnetId         string
		name             string
		availabilityZone string
		isDefault        *bool
		isRemoteVpcSNAT  *bool
		tagKey           string
		cidrBlock        string
		cdcId            string
	)

	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcId = temp.(string)
	}

	if temp, ok := d.GetOk("subnet_id"); ok {
		subnetId = temp.(string)
	}

	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}

	if temp, ok := d.GetOk("availability_zone"); ok {
		availabilityZone = temp.(string)
	}

	if temp, ok := d.GetOkExists("is_default"); ok {
		isDefault = helper.Bool(temp.(bool))
	}

	if temp, ok := d.GetOkExists("is_remote_vpc_snat"); ok {
		isRemoteVpcSNAT = helper.Bool(temp.(bool))
	}

	if temp, ok := d.GetOk("tag_key"); ok {
		tagKey = temp.(string)
	}

	if temp, ok := d.GetOk("cidr_block"); ok {
		cidrBlock = temp.(string)
	}

	if temp, ok := d.GetOk("cdc_id"); ok {
		cdcId = temp.(string)
	}

	var (
		tags  = helper.GetTags(d, "tags")
		infos []VpcSubnetBasicInfo
		err   error
	)

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, err = vpcService.DescribeSubnets(ctx, subnetId, vpcId,
			name, availabilityZone, tags,
			isDefault, isRemoteVpcSNAT, tagKey,
			cidrBlock, cdcId)

		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}

		return nil
	})

	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))
	for _, item := range infos {
		respTags, err := tagService.DescribeResourceTags(ctx, "vpc", "subnet", region, item.subnetId)
		if err != nil {
			return err
		}

		var infoMap = make(map[string]interface{})
		infoMap["availability_zone"] = item.zone
		infoMap["vpc_id"] = item.vpcId
		infoMap["subnet_id"] = item.subnetId
		infoMap["name"] = item.name
		infoMap["cidr_block"] = item.cidr
		infoMap["cdc_id"] = item.cdcId
		infoMap["is_default"] = item.isDefault
		infoMap["is_multicast"] = item.isMulticast
		infoMap["route_table_id"] = item.routeTableId
		infoMap["available_ip_count"] = item.availableIpCount
		infoMap["create_time"] = item.createTime
		infoMap["tags"] = respTags
		infoList = append(infoList, infoMap)
	}

	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  subnet instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	idBytes, err := json.Marshal(map[string]interface{}{
		"vpcId":            vpcId,
		"subnetId":         subnetId,
		"availabilityZone": availabilityZone,
		"name":             name,
		"isDefault":        isDefault,
		"tagKey":           tagKey,
		"isRemoteVpcSnat":  isRemoteVpcSNAT,
		"cidrBlock":        cidrBlock,
		"tags":             tags,
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
