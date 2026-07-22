package ccn

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudCcnInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnInstancesRead,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID CCN 到 是 queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 CCN 到 是 queried。",
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
				Description: "Information 列表 CCN。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ccn_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CCN。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CCN。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 CCN。",
						},
						"qos": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service quality 的 CCN，和 可用 值 include 'PT'，'AU'，'AG'. 默认为 'AU'。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "States 的 实例. 可用 值 include 'ISOLATED'(arrears) 和 'AVAILABLE'。",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing 模式",
						},
						"bandwidth_limit_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "speed 限制 类型",
						},
						"attachment_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information 列表 实例 是 attached。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 attached 实例 网络，和 可用 值 include VPC，DIRECTCONNECT，BMVPC 和 VPNGW。",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 该 实例 locates 在。",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID 实例 是 attached。",
									},
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "States 的 实例 是 attached，和 可用 值 include PENDING，ACTIVE，EXPIRED，REJECTED，DELETED，FAILED(asynchronous forced disassociation after 2 hours)，ATTACHING，DETACHING 和 DETACHFAILED(asynchronous forced disassociation after 2 hours)。",
									},
									"attached_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time 的 attaching。",
									},
									"cidr_block": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "A 网络 地址 block 的 实例 该 是 attached。",
									},
								},
							},
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

func dataSourceTencentCloudCcnInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ccn_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		ccnId = ""
		name  = ""
	)

	if temp, ok := d.GetOk("ccn_id"); ok {
		if tempStr := temp.(string); tempStr != "" {
			ccnId = tempStr
		}
	}

	if temp, ok := d.GetOk("name"); ok {
		if tempStr := temp.(string); tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeCcns(ctx, ccnId, name)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["ccn_id"] = item.ccnId
		infoMap["name"] = item.name
		infoMap["description"] = item.description
		infoMap["qos"] = item.qos
		infoMap["state"] = strings.ToUpper(item.state)
		infoMap["create_time"] = item.createTime
		infoMap["charge_type"] = item.chargeType
		infoMap["bandwidth_limit_type"] = item.bandWithLimitType
		infoList = append(infoList, infoMap)

		instances, err := service.DescribeCcnAttachedInstances(ctx, item.ccnId)
		if err != nil {
			return err
		}
		attachmentList := make([]interface{}, 0, len(instances))

		for _, instance := range instances {

			instanceMap := map[string]interface{}{
				"instance_type":   instance.instanceType,
				"instance_region": instance.instanceRegion,
				"instance_id":     instance.instanceId,
				"state":           strings.ToUpper(instance.state),
				"attached_time":   instance.attachedTime,
				"cidr_block":      instance.cidrBlock,
			}
			attachmentList = append(attachmentList, instanceMap)

		}

		infoMap["attachment_list"] = attachmentList

	}
	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  ccn instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	m := md5.New()
	_, err = m.Write([]byte("ccn_instances" + ccnId + "_" + name))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
