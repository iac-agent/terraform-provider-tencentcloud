package vpc

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcACLRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "ID VPC instance。",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 60),
				Description:  "名称 network ACL。",
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "ID network ACL instance。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"acl_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The information 列表 the VPC. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC instance。",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID network ACL instance。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 network ACL。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subnets associated with the network ACL。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID VPC instance。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网实例 ID",
									},
									"subnet_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet 名称",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IPv4 CIDR of the subnet。",
									},
									"tags": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "标签 of the subnet。",
									},
								},
							},
						},
						"ingress": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inbound rules of the network ACL。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 IP 协议",
									},
									"port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Range of the 端口",
									},
									"policy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule policy of Network ACL。",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An IP 地址 network or segment。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule 描述",
									},
								},
							},
						},
						"egress": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Outbound rules of the network ACL。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 IP 协议",
									},
									"port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Range of the 端口",
									},
									"policy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule policy of Network ACL。",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An IP 地址 network or segment。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule 描述",
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

func dataSourceTencentCloudVpcACLRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_acls.read")()
	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		vpcID = d.Get("vpc_id").(string)
		name  = d.Get("name").(string)
		id    = d.Get("id").(string)
	)

	networkAcls, err := service.DescribeNetWorkAcls(ctx, id, vpcID, name)
	if err != nil {
		return err
	}

	aclList := make([]map[string]interface{}, 0, len(networkAcls))
	ids := make([]string, 0, len(networkAcls))

	for _, info := range networkAcls {
		subnetInfo := info.SubnetSet
		subnets := make([]map[string]interface{}, 0, len(subnetInfo))
		for i := range subnetInfo {
			v := subnetInfo[i]
			subnet := make(map[string]interface{}, 5)
			subnet["vpc_id"] = v.VpcId
			subnet["subnet_id"] = v.SubnetId
			subnet["subnet_name"] = v.SubnetName
			subnet["cidr_block"] = v.CidrBlock

			tag := make(map[string]interface{}, len(v.TagSet))
			for t := range v.TagSet {
				tagValue := v.TagSet[t]
				tag[*tagValue.Key] = tagValue.Value
			}
			subnet["tags"] = tag

			subnets = append(subnets, subnet)
		}

		ingressInfo := info.IngressEntries
		ingress := make([]map[string]interface{}, 0, len(ingressInfo))
		for i := range ingressInfo {
			v := ingressInfo[i]
			egressMap := map[string]interface{}{
				"protocol":    v.Protocol,
				"port":        v.Port,
				"cidr_block":  v.CidrBlock,
				"policy":      v.Action,
				"description": v.Description,
			}
			ingress = append(ingress, egressMap)
		}

		egressInfo := info.EgressEntries
		egress := make([]map[string]interface{}, 0, len(egressInfo))
		for i := range egressInfo {
			v := egressInfo[i]
			egressMap := map[string]interface{}{
				"protocol":    v.Protocol,
				"port":        v.Port,
				"cidr_block":  v.CidrBlock,
				"policy":      v.Action,
				"description": v.Description,
			}
			egress = append(egress, egressMap)
		}

		aclResult := map[string]interface{}{
			"vpc_id":      info.VpcId,
			"id":          info.NetworkAclId,
			"name":        info.NetworkAclName,
			"create_time": info.CreatedTime,
			"subnets":     subnets,
			"ingress":     ingress,
			"egress":      egress,
		}
		aclList = append(aclList, aclResult)
		ids = append(ids, *info.NetworkAclId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("acl_list", aclList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set ACL list fail, reason:%v \n ", logId, err)
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), aclList); err != nil {
			return err
		}
	}
	return nil
}
