package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudAsScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingGroupCreate,
		Read:   resourceTencentCloudAsScalingGroupRead,
		Update: resourceTencentCloudAsScalingGroupUpdate,
		Delete: resourceTencentCloudAsScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 55),
				Description:  "名称 scaling 组。",
			},
			"configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An 可用 ID 对于 launch 配置。",
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 2000),
				Description:  "最大CVM 实例. 有效 值 ranges: (0~2000)。",
			},
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 2000),
				Description:  "最小CVM 实例. 有效 值 ranges: (0~2000)。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID VPC 网络。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "指定to 其中 项目 scaling 组 belongs。",
			},
			"subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID 列表 子网，和 对于 VPC 它 为必填项。",
			},
			"zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "列表 可用 zones，对于 Basic 网络 它 为必填项。",
			},
			"default_cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Default cooldown 时间 在 second，和 默认值为 `300`。",
			},
			"desired_capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Desired 卷 的 CVM 实例，其中 是 between `max_size` 和 `min_size`。",
			},
			"load_balancer_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"forward_balancer_ids"},
				Description:   "ID 列表 traditional load balancers。",
			},
			"forward_balancer_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"load_balancer_ids"},
				Description:   "列表 应用 load balancers，其中 可以't 是 指定 使用 `load_balancer_ids` together。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID 可用 load balancers。",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Listener ID 对于 应用 load balancers。",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID forwarding 规则。",
						},
						"target_attribute": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Attribute 列表 目标 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "端口 数量。",
									},
									"weight": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "权重",
									},
								},
							},
						},
					},
				},
			},
			"termination_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Available 值 对于 termination policies. 有效值：OLDEST_INSTANCE 和 NEWEST_INSTANCE。",
				Elem: &schema.Schema{
					Type:    schema.TypeString,
					Default: SCALING_GROUP_TERMINATION_POLICY_OLDEST_INSTANCE,
					ValidateFunc: tccommon.ValidateAllowedStringValue([]string{SCALING_GROUP_TERMINATION_POLICY_OLDEST_INSTANCE,
						SCALING_GROUP_TERMINATION_POLICY_NEWEST_INSTANCE}),
				},
			},
			"retry_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Available 值 对于 retry policies. 有效值：IMMEDIATE_RETRY 和 INCREMENTAL_INTERVALS。",
				Default:     SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
					SCALING_GROUP_RETRY_POLICY_INCREMENTAL_INTERVALS}),
			},
			// Service Settings
			"replace_monitor_unhealthy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Enables unhealthy 实例 replacement. 如果 集合 到 `true`，AS 将 replace 实例 该 是 flagged 作为 unhealthy 通过 Cloud Monitor。",
			},
			"scaling_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "表示scaling 模式 其中 creates 和 terminates 实例 (classic 方法)，或 方法 first tries 到 start stopped 实例 (wake up stopped) 到 perform scaling operations. 可用值：`CLASSIC_SCALING`，`WAKE_UP_STOPPED_SCALING`. 默认值：`CLASSIC_SCALING`。",
			},
			"replace_load_balancer_unhealthy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Enable unhealthy 实例 replacement. 如果 集合 到 `true`，AS 将 replace 实例 该 是 found unhealthy 在 CLB health check。",
			},
			"replace_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Replace 模式 的 unhealthy replacement 服务. 有效值：RECREATE: Rebuild 实例 到 replace original unhealthy 实例. RESET: Performing 系统 reinstallation 在 unhealthy 实例 到 keep 信息 such 作为 数据 disks，私有 IP addresses，和 实例 IDs unchanged. 实例 login settings，HostName，enhanced services，和 UserData 将 remain consistent 使用 当前 launch 配置. 默认值：RECREATE. 注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"desired_capacity_sync_with_max_min_size": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "expected 数量 实例 是 synchronized 使用 最大 和 最小 值. 默认值为 `False`. 此 参数 是 effective 仅 在 scenario 其中 expected 数量 是 不 passed 在 当 modifying scaling 组 interface. True: 当 modifying 最大 或 最小 值，如果 there 是 conflict 使用 当前 expected 数量， expected 数量 是 adjusted synchronously. For 示例，当 modifying，如果 最小 值 2 是 passed 在 和 当前 expected 数量 是 1， expected 数量 是 adjusted synchronously 到 2; False: 当 modifying 最大 或 最小 值，如果 there 是 conflict 使用 当前 expected 数量， 错误信息 是 displayed indicating 该 modification 是 不 allowed。",
			},
			"priority_scale_in_unhealthy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "是否enable 优先级 对于 unhealthy 实例 during scale-在 operations. 如果 集合 到 `true`，unhealthy 实例 将 是 removed first 当 scaling 在。",
			},
			"health_check_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Health check 类型 实例 在 scaling 组.<br><li>CVM: confirm whether 实例 是 healthy based 在 网络 状态 如果 pinged 实例 是 unreachable， 实例 将 是 considered unhealthy. For more 信息，see [实例 Health Check](https://intl.云.tencent.com/document/product/377/8553?from_cn_redirect=1)<br><li>CLB: confirm whether 实例 是 healthy based 在 CLB health check 状态 For more 信息，see [Health Check Overview](https://intl.云.tencent.com/document/product/214/6097?from_cn_redirect=1).<br>如果 参数 是 集合 到 `CLB`， scaling 组 将 check both 网络 状态 和 CLB health check 状态 如果 网络 check 表示unhealthy， `HealthStatus` 字段 将 返回 `UNHEALTHY`. 如果 CLB health check 表示unhealthy， `HealthStatus` 字段 将 返回 `CLB_UNHEALTHY`. 如果 both checks indicate unhealthy， `HealthStatus` 字段 将 返回 `UNHEALTHY|CLB_UNHEALTHY`. 默认值：`CLB`。",
			},
			"lb_health_check_grace_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Grace 周期 的 CLB health check during 其中 `IN_SERVICE` 实例 added 将 不 是 marked 作为 `CLB_UNHEALTHY`.<br>有效 范围: 0-7200，（秒）。 默认值：`0`。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 scaling 组。",
			},

			// computed value
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current 状态 scaling 组。",
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "实例 数量 scaling 组。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "时间 当 AS 组 是 创建。",
			},
			"multi_zone_subnet_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{MultiZoneSubnetPolicyPriority,
					MultiZoneSubnetPolicyEquality}),
				Description: "Multi 可用区 或 子网 strategy，有效值：PRIORITY 和 EQUALITY。",
			},
		},
	}
}

func resourceTencentCloudAsScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	request := as.NewCreateAutoScalingGroupRequest()

	request.AutoScalingGroupName = helper.String(d.Get("scaling_group_name").(string))
	request.LaunchConfigurationId = helper.String(d.Get("configuration_id").(string))
	request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	request.MinSize = helper.IntUint64(d.Get("min_size").(int))
	request.VpcId = helper.String(d.Get("vpc_id").(string))
	if v, ok := d.GetOk("default_cooldown"); ok {
		request.DefaultCooldown = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("desired_capacity"); ok {
		request.DesiredCapacity = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("retry_policy"); ok {
		request.RetryPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIds := v.([]interface{})
		request.SubnetIds = make([]*string, 0, len(subnetIds))
		for i := range subnetIds {
			subnetId := subnetIds[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetId)
		}
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := v.([]interface{})
		request.Zones = make([]*string, 0, len(zones))
		for i := range zones {
			zone := zones[i].(string)
			request.Zones = append(request.Zones, &zone)
		}
	}

	if v, ok := d.GetOk("load_balancer_ids"); ok {
		loadBalancerIds := v.([]interface{})
		request.LoadBalancerIds = make([]*string, 0, len(loadBalancerIds))
		for i := range loadBalancerIds {
			loadBalancerId := loadBalancerIds[i].(string)
			request.LoadBalancerIds = append(request.LoadBalancerIds, &loadBalancerId)
		}
	}

	if v, ok := d.GetOk("forward_balancer_ids"); ok {
		forwardBalancers := v.(*schema.Set).List()
		request.ForwardLoadBalancers = make([]*as.ForwardLoadBalancer, 0, len(forwardBalancers))
		for _, v := range forwardBalancers {
			vv := v.(map[string]interface{})
			targets := vv["target_attribute"].([]interface{})
			forwardBalancer := as.ForwardLoadBalancer{
				LoadBalancerId: helper.String(vv["load_balancer_id"].(string)),
				ListenerId:     helper.String(vv["listener_id"].(string)),
				LocationId:     helper.String(vv["rule_id"].(string)),
			}
			forwardBalancer.TargetAttributes = make([]*as.TargetAttribute, 0, len(targets))
			for _, target := range targets {
				t := target.(map[string]interface{})
				targetAttribute := as.TargetAttribute{
					Port:   helper.IntUint64(t["port"].(int)),
					Weight: helper.IntUint64(t["weight"].(int)),
				}
				forwardBalancer.TargetAttributes = append(forwardBalancer.TargetAttributes, &targetAttribute)
			}

			request.ForwardLoadBalancers = append(request.ForwardLoadBalancers, &forwardBalancer)
		}
	}

	if v, ok := d.GetOk("termination_policies"); ok {
		terminationPolicies := v.([]interface{})
		request.TerminationPolicies = make([]*string, 0, len(terminationPolicies))
		for i := range terminationPolicies {
			terminationPolicy := terminationPolicies[i].(string)
			request.TerminationPolicies = append(request.TerminationPolicies, &terminationPolicy)
		}
	}

	if v, ok := d.GetOk("multi_zone_subnet_policy"); ok {
		request.MultiZoneSubnetPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("health_check_type"); ok {
		request.HealthCheckType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("lb_health_check_grace_period"); ok {
		request.LoadBalancerHealthCheckGracePeriod = helper.IntUint64(v.(int))
	}

	var (
		replaceMonitorUnhealthy           = d.Get("replace_monitor_unhealthy").(bool)
		scalingMode                       = d.Get("scaling_mode").(string)
		replaceLBUnhealthy                = d.Get("replace_load_balancer_unhealthy").(bool)
		replaceMode                       = d.Get("replace_mode").(string)
		desiredCapacitySyncWithMaxMinSize = d.Get("desired_capacity_sync_with_max_min_size").(bool)
		priorityScaleInUnhealthy          = d.Get("priority_scale_in_unhealthy").(bool)
	)

	if replaceMonitorUnhealthy || scalingMode != "" || replaceLBUnhealthy || replaceMode != "" || desiredCapacitySyncWithMaxMinSize || priorityScaleInUnhealthy {
		if scalingMode == "" {
			scalingMode = SCALING_MODE_CLASSIC
		}

		if replaceMode == "" {
			replaceMode = REPLACE_MODE_RECREATE
		}

		request.ServiceSettings = &as.ServiceSettings{
			ReplaceMonitorUnhealthy:           &replaceMonitorUnhealthy,
			ScalingMode:                       &scalingMode,
			ReplaceLoadBalancerUnhealthy:      &replaceLBUnhealthy,
			ReplaceMode:                       &replaceMode,
			DesiredCapacitySyncWithMaxMinSize: &desiredCapacitySyncWithMaxMinSize,
			PriorityScaleInUnhealthy:          &priorityScaleInUnhealthy,
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := as.Tag{
				ResourceType: helper.String("auto-scaling-group"),
				Key:          helper.String(tagKey),
				Value:        helper.String(tagValue),
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	var id string
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateAutoScalingGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.AutoScalingGroupId == nil {
			err = fmt.Errorf("Create auto scaling group failed, Auto scaling group id is nil.")
			return resource.NonRetryableError(err)
		}

		id = *response.Response.AutoScalingGroupId

		return nil
	}); err != nil {
		return err
	}
	d.SetId(id)

	// wait for status
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		scalingGroup, _, errRet := asService.DescribeAutoScalingGroupById(ctx, id)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if scalingGroup != nil && scalingGroup.InActivityStatus != nil && *scalingGroup.InActivityStatus == SCALING_GROUP_NOT_IN_ACTIVITY_STATUS {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("scaling group status is %s, retry...", *scalingGroup.InActivityStatus))
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudAsScalingGroupRead(d, meta)
}

func resourceTencentCloudAsScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var (
		scalingGroup *as.AutoScalingGroup
		e            error
		has          int
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		scalingGroup, has, e = asService.DescribeAutoScalingGroupById(ctx, scalingGroupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if has == 0 {
		d.SetId("")
		return nil
	}

	_ = d.Set("scaling_group_name", scalingGroup.AutoScalingGroupName)
	_ = d.Set("configuration_id", scalingGroup.LaunchConfigurationId)
	_ = d.Set("status", scalingGroup.AutoScalingGroupStatus)
	_ = d.Set("instance_count", scalingGroup.InstanceCount)
	_ = d.Set("max_size", scalingGroup.MaxSize)
	_ = d.Set("min_size", scalingGroup.MinSize)
	_ = d.Set("vpc_id", scalingGroup.VpcId)
	_ = d.Set("project_id", scalingGroup.ProjectId)
	_ = d.Set("subnet_ids", helper.StringsInterfaces(scalingGroup.SubnetIdSet))
	_ = d.Set("zones", helper.StringsInterfaces(scalingGroup.ZoneSet))
	_ = d.Set("default_cooldown", scalingGroup.DefaultCooldown)
	_ = d.Set("desired_capacity", scalingGroup.DesiredCapacity)
	_ = d.Set("load_balancer_ids", helper.StringsInterfaces(scalingGroup.LoadBalancerIdSet))
	_ = d.Set("termination_policies", helper.StringsInterfaces(scalingGroup.TerminationPolicySet))
	_ = d.Set("retry_policy", scalingGroup.RetryPolicy)
	_ = d.Set("create_time", scalingGroup.CreatedTime)
	_ = d.Set("retry_policy", scalingGroup.RetryPolicy)
	_ = d.Set("health_check_type", scalingGroup.HealthCheckType)
	_ = d.Set("lb_health_check_grace_period", scalingGroup.LoadBalancerHealthCheckGracePeriod)
	if v, ok := d.GetOk("multi_zone_subnet_policy"); ok && v.(string) != "" {
		_ = d.Set("multi_zone_subnet_policy", scalingGroup.MultiZoneSubnetPolicy)
	}

	if scalingGroup.ServiceSettings != nil {
		_ = d.Set("replace_monitor_unhealthy", scalingGroup.ServiceSettings.ReplaceMonitorUnhealthy)
		_ = d.Set("scaling_mode", scalingGroup.ServiceSettings.ScalingMode)
		_ = d.Set("replace_load_balancer_unhealthy", scalingGroup.ServiceSettings.ReplaceLoadBalancerUnhealthy)
		_ = d.Set("replace_mode", scalingGroup.ServiceSettings.ReplaceMode)
		_ = d.Set("desired_capacity_sync_with_max_min_size", scalingGroup.ServiceSettings.DesiredCapacitySyncWithMaxMinSize)
		_ = d.Set("priority_scale_in_unhealthy", scalingGroup.ServiceSettings.PriorityScaleInUnhealthy)
	}

	if scalingGroup.ForwardLoadBalancerSet != nil && len(scalingGroup.ForwardLoadBalancerSet) > 0 {
		forwardLoadBalancers := make([]map[string]interface{}, 0, len(scalingGroup.ForwardLoadBalancerSet))
		for _, v := range scalingGroup.ForwardLoadBalancerSet {
			targetAttributes := make([]map[string]interface{}, 0, len(v.TargetAttributes))
			for _, vv := range v.TargetAttributes {
				targetAttribute := map[string]interface{}{
					"port":   *vv.Port,
					"weight": *vv.Weight,
				}
				targetAttributes = append(targetAttributes, targetAttribute)
			}
			forwardLoadBalancer := map[string]interface{}{
				"load_balancer_id": *v.LoadBalancerId,
				"listener_id":      *v.ListenerId,
				"target_attribute": targetAttributes,
				"rule_id":          *v.LocationId,
			}
			forwardLoadBalancers = append(forwardLoadBalancers, forwardLoadBalancer)
		}
		_ = d.Set("forward_balancer_ids", forwardLoadBalancers)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "as", "auto-scaling-group", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudAsScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region

	request := as.NewModifyAutoScalingGroupRequest()
	scalingGroupId := d.Id()

	d.Partial(true)

	var updateAttrs []string

	request.AutoScalingGroupId = &scalingGroupId
	if d.HasChange("scaling_group_name") {
		updateAttrs = append(updateAttrs, "scaling_group_name")
		request.AutoScalingGroupName = helper.String(d.Get("scaling_group_name").(string))
	}
	if d.HasChange("configuration_id") {
		updateAttrs = append(updateAttrs, "configuration_id")
		request.LaunchConfigurationId = helper.String(d.Get("configuration_id").(string))
	}
	if d.HasChange("max_size") {
		updateAttrs = append(updateAttrs, "max_size")
		request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	}
	if d.HasChange("min_size") {
		updateAttrs = append(updateAttrs, "max_size")
		request.MinSize = helper.IntUint64(d.Get("min_size").(int))
	}
	if d.HasChange("vpc_id") {
		updateAttrs = append(updateAttrs, "vpc_id")
		request.VpcId = helper.String(d.Get("vpc_id").(string))
	}
	if d.HasChange("project_id") {
		updateAttrs = append(updateAttrs, "project_id")
		request.ProjectId = helper.IntUint64(d.Get("project_id").(int))
	}
	if d.HasChange("default_cooldown") {
		updateAttrs = append(updateAttrs, "default_cooldown")
		request.DefaultCooldown = helper.IntUint64(d.Get("default_cooldown").(int))
	}
	if d.HasChange("desired_capacity") {
		updateAttrs = append(updateAttrs, "desired_capacity")
		request.DesiredCapacity = helper.IntUint64(d.Get("desired_capacity").(int))
	}
	if d.HasChange("retry_policy") {
		updateAttrs = append(updateAttrs, "retry_policy")
		request.RetryPolicy = helper.String(d.Get("retry_policy").(string))
	}
	if d.HasChange("subnet_ids") {
		updateAttrs = append(updateAttrs, "subnet_ids")
		subnetIds := d.Get("subnet_ids").([]interface{})
		request.SubnetIds = make([]*string, 0, len(subnetIds))
		for i := range subnetIds {
			subnetId := subnetIds[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetId)
		}
	}
	if d.HasChange("zones") {
		updateAttrs = append(updateAttrs, "zones")
		zones := d.Get("zones").([]interface{})
		request.Zones = make([]*string, 0, len(zones))
		for i := range zones {
			zone := zones[i].(string)
			request.Zones = append(request.Zones, &zone)
		}
	}
	if d.HasChange("termination_policies") {
		updateAttrs = append(updateAttrs, "termination_policies")
		terminationPolicies := d.Get("termination_policies").([]interface{})
		request.TerminationPolicies = make([]*string, 0, len(terminationPolicies))
		for i := range terminationPolicies {
			terminationPolicy := terminationPolicies[i].(string)
			request.TerminationPolicies = append(request.TerminationPolicies, &terminationPolicy)
		}
	}

	if d.HasChange("multi_zone_subnet_policy") {
		updateAttrs = append(updateAttrs, "multi_zone_subnet_policy")
		request.MultiZoneSubnetPolicy = helper.String(d.Get("multi_zone_subnet_policy").(string))
	}

	if d.HasChange("replace_monitor_unhealthy") ||
		d.HasChange("scaling_mode") ||
		d.HasChange("replace_load_balancer_unhealthy") ||
		d.HasChange("replace_mode") ||
		d.HasChange("desired_capacity_sync_with_max_min_size") ||
		d.HasChange("priority_scale_in_unhealthy") {
		updateAttrs = append(updateAttrs, "replace_monitor_unhealthy", "scaling_mode", "replace_load_balancer_unhealthy", "replace_mode", "desired_capacity_sync_with_max_min_size", "priority_scale_in_unhealthy")
		scalingMode := d.Get("scaling_mode").(string)
		replaceMode := d.Get("replace_mode").(string)
		if scalingMode == "" {
			scalingMode = SCALING_MODE_CLASSIC
		}
		if replaceMode == "" {
			replaceMode = REPLACE_MODE_RECREATE
		}
		replaceMonitor := d.Get("replace_monitor_unhealthy").(bool)
		replaceLB := d.Get("replace_load_balancer_unhealthy").(bool)
		desiredCapacitySyncWithMaxMinSize := d.Get("desired_capacity_sync_with_max_min_size").(bool)
		priorityScaleInUnhealthy := d.Get("priority_scale_in_unhealthy").(bool)
		request.ServiceSettings = &as.ServiceSettings{
			ReplaceMonitorUnhealthy:           &replaceMonitor,
			ScalingMode:                       &scalingMode,
			ReplaceLoadBalancerUnhealthy:      &replaceLB,
			ReplaceMode:                       &replaceMode,
			DesiredCapacitySyncWithMaxMinSize: &desiredCapacitySyncWithMaxMinSize,
			PriorityScaleInUnhealthy:          &priorityScaleInUnhealthy,
		}
	}

	if d.HasChange("health_check_type") || d.HasChange("lb_health_check_grace_period") {
		request.HealthCheckType = helper.String(d.Get("health_check_type").(string))
		if v, ok := d.GetOkExists("lb_health_check_grace_period"); ok {
			request.LoadBalancerHealthCheckGracePeriod = helper.IntUint64(v.(int))
		}
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.UseAsClient().ModifyAutoScalingGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		return nil
	}); err != nil {
		return err
	}

	updateAttrs = updateAttrs[:0]

	balancerRequest := as.NewModifyLoadBalancersRequest()
	balancerRequest.AutoScalingGroupId = &scalingGroupId
	if d.HasChange("load_balancer_ids") {
		updateAttrs = append(updateAttrs, "load_balancer_ids")

		loadBalancerIds := d.Get("load_balancer_ids").([]interface{})
		balancerRequest.LoadBalancerIds = make([]*string, 0, len(loadBalancerIds))
		for i := range loadBalancerIds {
			loadBalancerId := loadBalancerIds[i].(string)
			balancerRequest.LoadBalancerIds = append(balancerRequest.LoadBalancerIds, &loadBalancerId)
		}
	}

	if d.HasChange("forward_balancer_ids") {
		updateAttrs = append(updateAttrs, "forward_balancer_ids")

		forwardBalancers := d.Get("forward_balancer_ids").(*schema.Set).List()
		balancerRequest.ForwardLoadBalancers = make([]*as.ForwardLoadBalancer, 0, len(forwardBalancers))
		for _, v := range forwardBalancers {
			vv := v.(map[string]interface{})
			targets := vv["target_attribute"].([]interface{})
			forwardBalancer := as.ForwardLoadBalancer{
				LoadBalancerId: helper.String(vv["load_balancer_id"].(string)),
				ListenerId:     helper.String(vv["listener_id"].(string)),
				LocationId:     helper.String(vv["rule_id"].(string)),
			}
			forwardBalancer.TargetAttributes = make([]*as.TargetAttribute, 0, len(targets))
			for _, target := range targets {
				t := target.(map[string]interface{})
				targetAttribute := as.TargetAttribute{
					Port:   helper.IntUint64(t["port"].(int)),
					Weight: helper.IntUint64(t["weight"].(int)),
				}
				forwardBalancer.TargetAttributes = append(forwardBalancer.TargetAttributes, &targetAttribute)
			}

			balancerRequest.ForwardLoadBalancers = append(balancerRequest.ForwardLoadBalancers, &forwardBalancer)
		}
	}

	if len(updateAttrs) > 0 {
		if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(balancerRequest.GetAction())

			balancerResponse, err := client.UseAsClient().ModifyLoadBalancers(balancerRequest)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, balancerRequest.GetAction(), balancerRequest.ToJsonString(), err.Error())
				return tccommon.RetryError(err)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, balancerRequest.GetAction(), balancerRequest.ToJsonString(), balancerResponse.ToJsonString())

			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("as", "auto-scaling-group", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudAsScalingGroupRead(d, meta)
}

func resourceTencentCloudAsScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_group.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scalingGroupId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	// We need read the scaling group in order to check if there are instances.
	// If so, we need to remove those first.
	scalingGroup, has, err := asService.DescribeAutoScalingGroupById(ctx, scalingGroupId)
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}
	if *scalingGroup.InstanceCount > 0 || *scalingGroup.DesiredCapacity > 0 {
		if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := asService.ClearScalingGroupInstance(ctx, scalingGroupId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		if errRet := asService.DeleteScalingGroup(ctx, scalingGroupId); errRet != nil {
			if sdkErr, ok := errRet.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == AsScalingGroupNotFound {
					return nil
				} else if sdkErr.Code == AsScalingGroupInProgress || sdkErr.Code == AsScalingGroupInstanceInGroup {
					return resource.RetryableError(sdkErr)
				}
			}
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
