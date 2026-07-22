package vpc

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcInstancesRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID VPC to be queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 VPC to be queried。",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter default or no default VPC。",
			},
			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter if VPC has this 标签",
			},
			"cidr_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter VPC with this CIDR。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 of the VPC to be queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"instance_list": {Type: schema.TypeList,
				Computed:    true,
				Description: "The information 列表 the VPC。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 VPC。",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A network 地址 block of a VPC CIDR。",
						},
						"common_assistant_cidr": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "common assistant CIDR block。",
						},
						"container_assistant_cidr": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "container assistant CIDR block。",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否it is the default VPC for this 地域",
						},
						"is_multicast": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否VPC multicast is 已启用",
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "A 列表 DNS servers which can be used within the VPC。",
						},
						"subnet_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "A ID 列表 subnets within this VPC。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of VPC。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 of the VPC。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpcInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		vpcId     string
		name      string
		isDefault *bool
		tagKey    string
		cidrBlock string
	)
	if temp, ok := d.GetOk("vpc_id"); ok {
		vpcId = temp.(string)
	}
	if temp, ok := d.GetOk("name"); ok {
		name = temp.(string)
	}
	if temp, ok := d.GetOkExists("is_default"); ok {
		isDefault = helper.Bool(temp.(bool))
	}
	if temp, ok := d.GetOkExists("tag_key"); ok {
		tagKey = temp.(string)
	}
	if temp, ok := d.GetOkExists("cidr_block"); ok {
		cidrBlock = temp.(string)
	}

	var (
		tags     = helper.GetTags(d, "tags")
		vpcInfos []VpcBasicInfo
		err      error
	)
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		vpcInfos, err = service.DescribeVpcs(ctx, vpcId, name, tags, isDefault, tagKey, cidrBlock)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	var vpcInfoList = make([]map[string]interface{}, 0, len(vpcInfos))

	for _, item := range vpcInfos {
		var infoMap = make(map[string]interface{})
		infoMap["vpc_id"] = item.vpcId
		infoMap["name"] = item.name
		infoMap["cidr_block"] = item.cidr
		infoMap["is_default"] = item.isDefault
		infoMap["is_multicast"] = item.isMulticast
		infoMap["dns_servers"] = item.dnsServers
		infoMap["create_time"] = item.createTime
		infoMap["common_assistant_cidr"] = item.assistantCidrs
		infoMap["container_assistant_cidr"] = item.dockerAssistantCidrs

		respTags := make(map[string]string, len(item.tags))
		for _, tag := range item.tags {
			if tag.Key == nil {
				return errors.New("vpc tag key is nil")
			}
			if tag.Value == nil {
				return errors.New("vpc tag value is nil")
			}

			respTags[*tag.Key] = *tag.Value
		}
		infoMap["tags"] = respTags

		subnetInfos, err := service.DescribeSubnets(ctx, "", item.vpcId, "", "", nil, nil, nil, "", "", "")
		if err != nil {
			return err
		}
		subnetIds := make([]string, 0, len(subnetInfos))
		for _, v := range subnetInfos {
			subnetIds = append(subnetIds, v.subnetId)
		}

		infoMap["subnet_ids"] = subnetIds
		vpcInfoList = append(vpcInfoList, infoMap)
	}

	if err := d.Set("instance_list", vpcInfoList); err != nil {
		log.Printf("[CRITAL]%s provider set  vpc instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	idBytes, err := json.Marshal(map[string]interface{}{
		"vpcId":     vpcId,
		"name":      name,
		"isDefault": isDefault,
		"tagKey":    tagKey,
		"cidrBlock": cidrBlock,
		"tags":      tags,
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
		if err := tccommon.WriteToFile(output.(string), vpcInfoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
