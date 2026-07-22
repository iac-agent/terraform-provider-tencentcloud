package clb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbTargetGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbTargetGroupRead,

		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_id", "target_group_name"},
				AtLeastOneOf:  []string{"vpc_id", "target_group_name"},
				Description: "目标组的 ID。与“vpc_id”和“target_group_name”互斥。首选“target_group_id”。",
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"target_group_id", "target_group_name"},
				Description: "目标组VPC ID。与“target_group_id”互斥。首选“target_group_id”。",
			},
			"target_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"target_group_id", "vpc_id"},
				Description: "目标群体的名称。与“target_group_id”互斥。首选“target_group_id”。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "目标群体信息列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标组的 ID。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标群体的名称。",
						},
						"target_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标组VPC ID。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "目标组的端口。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标组的创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标组的修改时间。",
						},
						"associated_rule_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "相关规则列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "负载均衡ID。",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "听众 ID。",
									},
									"location_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "转发规则ID。",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "侦听器协议类型。",
									},
									"listener_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "侦听器端口。",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "转发规则域。",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "转发规则 URL。",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "负载平衡名称。",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "听众姓名。",
									},
								},
							},
						},
						"target_group_instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "目标组绑定的后端服务器列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "后端服务的类型。",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "后端服务ID。",
									},
									"server_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "后端服务端口。",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "后端服务的转发权重。",
									},
									"public_ip_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        schema.TypeString,
										Description: "后端服务的外网IP列表。",
									},
									"private_ip_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        schema.TypeString,
										Description: "后端服务的内网IP列表。",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "后端服务的实例名称。",
									},
									"registered_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "绑定后端服务的时间。",
									},
									"eni_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "弹性网络接口ID。",
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

func dataSourceTencentCloudClbTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_target_groups.read")()

	var (
		clbService = ClbService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		instances            []*clb.TargetGroupBackend
		targetInfos          []*clb.TargetGroupInfo
		filters              = make(map[string]string, 2)
		targetGroupInstances []map[string]interface{}
		targetGroupId        string
		err                  error
	)

	if id, ok := d.GetOk("target_group_id"); ok {
		targetGroupId = id.(string)
	}
	if id, ok := d.GetOk("vpc_id"); ok {
		filters["TargetGroupVpcId"] = id.(string)
	}
	if name, ok := d.GetOk("target_group_name"); ok {
		filters["TargetGroupName"] = name.(string)
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		targetInfos, err = clbService.DescribeTargetGroups(ctx, targetGroupId, filters)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	var (
		list    = make([]map[string]interface{}, 0, len(targetInfos))
		ids     = make([]string, 0, len(targetInfos))
		isExist = make(map[string]bool)
	)

	for _, info := range targetInfos {
		targetId := *info.TargetGroupId
		ids = append(ids, targetId)
		if _, ok := isExist[targetId]; !ok {
			instances = []*clb.TargetGroupBackend{}
			targetGroupInstances = []map[string]interface{}{}

			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instances, err = clbService.DescribeTargetGroupInstances(ctx, map[string]string{
					"TargetGroupId": *info.TargetGroupId,
				})
				if err != nil {
					return tccommon.RetryError(err, tccommon.InternalError)
				}
				return nil
			})
			if err != nil {
				return err
			}

			isExist[targetId] = true
			for _, instance := range instances {
				targetGroupInstances = append(targetGroupInstances, map[string]interface{}{
					"server_type":          instance.Type,
					"instance_id":          instance.InstanceId,
					"server_port":          instance.Port,
					"weight":               instance.Weight,
					"public_ip_addresses":  instance.PublicIpAddresses,
					"private_ip_addresses": instance.PrivateIpAddresses,
					"instance_name":        instance.InstanceName,
					"registered_time":      instance.RegisteredTime,
					"eni_id":               instance.EniId,
				})
			}
		}

		ruleInfo := make([]map[string]interface{}, 0, len(info.AssociatedRule))
		for _, rule := range info.AssociatedRule {
			ruleInfo = append(ruleInfo, map[string]interface{}{
				"load_balancer_id":   rule.LoadBalancerId,
				"listener_id":        rule.ListenerId,
				"location_id":        rule.LocationId,
				"protocol":           rule.Protocol,
				"listener_port":      rule.Port,
				"domain":             rule.Domain,
				"url":                rule.Url,
				"load_balancer_name": rule.LoadBalancerName,
				"listener_name":      rule.ListenerName,
			})
		}

		list = append(list, map[string]interface{}{
			"target_group_id":            targetId,
			"vpc_id":                     info.VpcId,
			"target_group_name":          info.TargetGroupName,
			"port":                       info.Port,
			"create_time":                info.CreatedTime,
			"update_time":                info.UpdatedTime,
			"associated_rule_list":       ruleInfo,
			"target_group_instance_list": targetGroupInstances,
		})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set target group list fail, reason:%s ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), list); err != nil {
			return err
		}
	}

	return nil
}
