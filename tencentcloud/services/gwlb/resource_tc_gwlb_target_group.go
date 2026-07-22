package gwlb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gwlbv20240906 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gwlb/v20240906"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGwlbTargetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGwlbTargetGroupCreate,
		Read:   resourceTencentCloudGwlbTargetGroupRead,
		Update: resourceTencentCloudGwlbTargetGroupUpdate,
		Delete: resourceTencentCloudGwlbTargetGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target 组名称，limited 到 60 字符。",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPCID attribute 的 目标 组. 如果 此 参数 是 left blank， 默认值 VPC 将 是 使用。",
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Default 端口 的 目标 组，其中 可以 是 使用 当 servers 是 added later. Either '端口' 或 'TargetGroupInstances.N.端口' 必须 是 filled 在。",
			},

			"target_group_instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Real 服务器 bound 到 目标 组。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private 网络 IP 的 目标 组 实例。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "端口 的 目标 组 实例. Only 6081 是 支持。",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "权重 的 目标 组 实例. Only 0 或 16 是 支持，和 non-0 是 uniformly treated 作为 16。",
						},
					},
				},
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "GWLB 目标 组 protocol.\n" +
					"	- TENCENT_GENEVE: GENEVE standard protocol;\n" +
					"	- AWS_GENEVE: GENEVE compatibility protocol (a ticket is required for allowlisting).",
			},

			"health_check": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Health check settings。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_switch": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "是否enable health check。",
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: "Protocol 使用 对于 health check, 其中 支持 PING 和 TCP 和 是 PING 通过 默认值.\n" +
								"	- PING: icmp;\n" +
								"	- TCP: tcp.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Health check 端口，其中 为必填项 当 probe 协议 是 TCP。",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Health check 超时. 默认为 2 秒. 取值范围：2-30 秒。",
						},
						"interval_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Detection 间隔 时间. 默认为 5 秒. 取值范围：2-300 秒。",
						},
						"health_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Health detection 阈值. 默认为 3 times. 取值范围：2-10 times。",
						},
						"un_health_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Unhealth detection 阈值. 默认为 3 times. 取值范围：2-10 times。",
						},
					},
				},
			},

			"schedule_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Load balancing algorithm.\n" +
					"	- IP_HASH_3_ELASTIC: elastic hashing.",
			},

			"all_dead_to_alive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether \"All Dead, All Alive\" 是 支持. It 是 支持 通过 默认值.",
			},
		},
	}
}

func resourceTencentCloudGwlbTargetGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_target_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
	)
	var (
		request  = gwlbv20240906.NewCreateTargetGroupRequest()
		response = gwlbv20240906.NewCreateTargetGroupResponse()
	)

	if v, ok := d.GetOk("target_group_name"); ok {
		request.TargetGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("port"); ok {
		request.Port = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if healthCheckMap, ok := helper.InterfacesHeadMap(d, "health_check"); ok {
		targetGroupHealthCheck := gwlbv20240906.TargetGroupHealthCheck{}
		if v, ok := healthCheckMap["health_switch"]; ok {
			targetGroupHealthCheck.HealthSwitch = helper.Bool(v.(bool))
		}
		if v, ok := healthCheckMap["protocol"]; ok {
			targetGroupHealthCheck.Protocol = helper.String(v.(string))
		}
		if v, ok := healthCheckMap["port"]; ok {
			targetGroupHealthCheck.Port = helper.IntInt64(v.(int))
		}
		if v, ok := healthCheckMap["timeout"]; ok {
			targetGroupHealthCheck.Timeout = helper.IntInt64(v.(int))
		}
		if v, ok := healthCheckMap["interval_time"]; ok {
			targetGroupHealthCheck.IntervalTime = helper.IntInt64(v.(int))
		}
		if v, ok := healthCheckMap["health_num"]; ok {
			targetGroupHealthCheck.HealthNum = helper.IntInt64(v.(int))
		}
		if v, ok := healthCheckMap["un_health_num"]; ok {
			targetGroupHealthCheck.UnHealthNum = helper.IntInt64(v.(int))
		}
		request.HealthCheck = &targetGroupHealthCheck
	}

	if v, ok := d.GetOk("schedule_algorithm"); ok {
		request.ScheduleAlgorithm = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("all_dead_to_alive"); ok {
		request.AllDeadToAlive = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGwlbV20240906Client().CreateTargetGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create gwlb target group failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.TargetGroupId
	_ = response

	d.SetId(instanceId)

	return resourceTencentCloudGwlbTargetGroupRead(d, meta)
}

func resourceTencentCloudGwlbTargetGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_target_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := GwlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	respData, err := service.DescribeGwlbTargetGroupById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `gwlb_target_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.TargetGroupId != nil {
		instanceId = *respData.TargetGroupId
	}

	if respData.VpcId != nil {
		_ = d.Set("vpc_id", respData.VpcId)
	}

	if respData.TargetGroupName != nil {
		_ = d.Set("target_group_name", respData.TargetGroupName)
	}

	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
	}

	if respData.HealthCheck != nil {
		healthCheck := make(map[string]interface{})
		if respData.HealthCheck.HealthSwitch != nil {
			healthCheck["health_switch"] = respData.HealthCheck.HealthSwitch
		}
		if respData.HealthCheck.Protocol != nil {
			healthCheck["protocol"] = respData.HealthCheck.Protocol
		}
		if respData.HealthCheck.Port != nil {
			healthCheck["port"] = respData.HealthCheck.Port
		}
		if respData.HealthCheck.Timeout != nil {
			healthCheck["timeout"] = respData.HealthCheck.Timeout
		}
		if respData.HealthCheck.IntervalTime != nil {
			healthCheck["interval_time"] = respData.HealthCheck.IntervalTime
		}
		if respData.HealthCheck.HealthNum != nil {
			healthCheck["health_num"] = respData.HealthCheck.HealthNum
		}
		if respData.HealthCheck.UnHealthNum != nil {
			healthCheck["un_health_num"] = respData.HealthCheck.UnHealthNum
		}

		_ = d.Set("health_check", []interface{}{healthCheck})
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.ScheduleAlgorithm != nil {
		_ = d.Set("schedule_algorithm", respData.ScheduleAlgorithm)
	}
	if respData.AllDeadToAlive != nil {
		_ = d.Set("all_dead_to_alive", respData.AllDeadToAlive)
	}

	targetGroupBackends, err := service.DescribeTargetGroupInstancesById(ctx, instanceId)
	if err != nil {
		return err
	}
	targetGroupInstanceList := make([]interface{}, 0)
	for _, targetGroupBackend := range targetGroupBackends {
		targetGroupBackendMap := make(map[string]interface{})
		if len(targetGroupBackend.PrivateIpAddresses) > 0 {
			targetGroupBackendMap["bind_ip"] = targetGroupBackend.PrivateIpAddresses[0]
		}
		if targetGroupBackend.Port != nil {
			targetGroupBackendMap["port"] = targetGroupBackend.Port
		}
		if targetGroupBackend.Weight != nil {
			targetGroupBackendMap["weight"] = targetGroupBackend.Weight
		}
	}
	_ = d.Set("target_group_instances", targetGroupInstanceList)

	return nil
}

func resourceTencentCloudGwlbTargetGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_target_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"vpc_id", "port", "protocol", "schedule_algorithm"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	instanceId := d.Id()

	needChange1 := false
	mutableArgs1 := []string{"target_group_name", "health_check", "all_dead_to_alive"}
	for _, v := range mutableArgs1 {
		if d.HasChange(v) {
			needChange1 = true
			break
		}
	}

	if needChange1 {
		request1 := gwlbv20240906.NewModifyTargetGroupAttributeRequest()

		request1.TargetGroupId = helper.String(instanceId)

		if v, ok := d.GetOk("target_group_name"); ok {
			request1.TargetGroupName = helper.String(v.(string))
		}

		if healthCheckMap, ok := helper.InterfacesHeadMap(d, "health_check"); ok {
			targetGroupHealthCheck := gwlbv20240906.TargetGroupHealthCheck{}
			if v, ok := healthCheckMap["health_switch"]; ok {
				targetGroupHealthCheck.HealthSwitch = helper.Bool(v.(bool))
			}
			if v, ok := healthCheckMap["protocol"]; ok {
				targetGroupHealthCheck.Protocol = helper.String(v.(string))
			}
			if v, ok := healthCheckMap["port"]; ok {
				targetGroupHealthCheck.Port = helper.IntInt64(v.(int))
			}
			if v, ok := healthCheckMap["timeout"]; ok {
				targetGroupHealthCheck.Timeout = helper.IntInt64(v.(int))
			}
			if v, ok := healthCheckMap["interval_time"]; ok {
				targetGroupHealthCheck.IntervalTime = helper.IntInt64(v.(int))
			}
			if v, ok := healthCheckMap["health_num"]; ok {
				targetGroupHealthCheck.HealthNum = helper.IntInt64(v.(int))
			}
			if v, ok := healthCheckMap["un_health_num"]; ok {
				targetGroupHealthCheck.UnHealthNum = helper.IntInt64(v.(int))
			}
			request1.HealthCheck = &targetGroupHealthCheck
		}

		if v, ok := d.GetOkExists("all_dead_to_alive"); ok {
			request1.AllDeadToAlive = helper.Bool(v.(bool))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGwlbV20240906Client().ModifyTargetGroupAttributeWithContext(ctx, request1)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update gwlb target group failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudGwlbTargetGroupRead(d, meta)
}

func resourceTencentCloudGwlbTargetGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gwlb_target_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	var (
		request  = gwlbv20240906.NewDeleteTargetGroupsRequest()
		response = gwlbv20240906.NewDeleteTargetGroupsResponse()
	)

	request.TargetGroupIds = helper.Strings([]string{instanceId})
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGwlbV20240906Client().DeleteTargetGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete gwlb target group failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = instanceId
	return nil
}
