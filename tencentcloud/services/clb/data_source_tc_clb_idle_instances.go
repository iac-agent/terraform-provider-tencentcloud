package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbIdleInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbIdleInstancesRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CLB实例区域。",
			},

			"idle_load_balancers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "空闲 CLB 列表。注意：该字段可能返回null，表示取不到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例ID。",
						},
						"load_balancer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB 实例名称。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例区域。",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例VIP。",
						},
						"idle_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "负载均衡器被视为空闲的原因。 NO_RULES：未配置规则。 NO_RS：规则不与服务器关联。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB实例状态，包括：0：正在创建； 1：跑步。",
						},
						"forward": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB型。取值范围：1（CLB）； 0（经典 CLB）。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "负载均衡主机名。注意：该字段可能返回null，表示取不到有效值。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudClbIdleInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_idle_loadbalancers.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_region"); ok {
		paramMap["LoadBalancerRegion"] = helper.String(v.(string))
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var idleLoadBalancers []*clb.IdleLoadBalancer

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbIdleInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		idleLoadBalancers = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(idleLoadBalancers))
	tmpList := make([]map[string]interface{}, 0, len(idleLoadBalancers))

	if idleLoadBalancers != nil {
		for _, idleLoadBalancer := range idleLoadBalancers {
			idleLoadBalancerMap := map[string]interface{}{}

			if idleLoadBalancer.LoadBalancerId != nil {
				idleLoadBalancerMap["load_balancer_id"] = idleLoadBalancer.LoadBalancerId
			}

			if idleLoadBalancer.LoadBalancerName != nil {
				idleLoadBalancerMap["load_balancer_name"] = idleLoadBalancer.LoadBalancerName
			}

			if idleLoadBalancer.Region != nil {
				idleLoadBalancerMap["region"] = idleLoadBalancer.Region
			}

			if idleLoadBalancer.Vip != nil {
				idleLoadBalancerMap["vip"] = idleLoadBalancer.Vip
			}

			if idleLoadBalancer.IdleReason != nil {
				idleLoadBalancerMap["idle_reason"] = idleLoadBalancer.IdleReason
			}

			if idleLoadBalancer.Status != nil {
				idleLoadBalancerMap["status"] = idleLoadBalancer.Status
			}

			if idleLoadBalancer.Forward != nil {
				idleLoadBalancerMap["forward"] = idleLoadBalancer.Forward
			}

			if idleLoadBalancer.Domain != nil {
				idleLoadBalancerMap["domain"] = idleLoadBalancer.Domain
			}

			ids = append(ids, *idleLoadBalancer.LoadBalancerId)
			tmpList = append(tmpList, idleLoadBalancerMap)
		}

		_ = d.Set("idle_load_balancers", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
