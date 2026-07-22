package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlRoGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRoGroupCreate,
		Read:   resourceTencentCloudMysqlRoGroupRead,
		Update: resourceTencentCloudMysqlRoGroupUpdate,
		Delete: resourceTencentCloudMysqlRoGroupDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID，格式为：cdbro-3i70uj0k。",
			},

			"ro_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "RO组的ID。",
			},

			"ro_group_info": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "RO 组的详细信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ro_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "RO 组名。",
						},
						"ro_max_delay_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "RO实例最大延迟阈值。单位为秒，最小值为1。注意RO组必须启用实例延迟剔除策略才能使该值有效。",
						},
						"ro_offline_delay": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "是否启用实例的延迟剔除。支持的值为： 1 - 开启； 0 - 未开启。请注意，如果启用实例延迟剔除，则必须设置延迟阈值（RoMaxDelayTime）参数。",
						},
						"min_ro_in_group": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "保留实例的最小数量。可以设置为小于或等于该RO组下RO实例数量的任意值。注意，如果设置值大于RO实例数量，则不会被移除；如果设置为0，则所有延迟超过限制的实例将被删除。",
						},
						"weight_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "重量模式。支持的值包括： `system` - 由系统自动分配； `custom` - 用户定义的设置。请注意，如果设置了“自定义”模式，则必须设置RO实例权重配置（RoWeightValues）参数。",
						},
						"replication_delay_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "复制时间延迟。",
						},
					},
				},
			},

			"ro_weight_values": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "RO组内实例的权重。如果RO组的权重模式改为用户自定义模式（自定义），则必须设置该参数，并且需要设置每个RO实例的权重值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "RO 实例 ID。",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "重量。取值范围为[0，100]。",
						},
					},
				},
			},

			"is_balance_ro_load": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否重新平衡RO组中RO实例的负载。支持的值包括： 1 - 重新平衡负载； 0 - 不重新平衡负载。默认值为0。注意，当设置为重新平衡负载时，RO组中的RO实例会出现数据库连接瞬间断开的情况，请确保应用程序可以重新连接数据库。",
			},
		},
	}
}

func resourceTencentCloudMysqlRoGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	var roGroupId string
	if v, ok := d.GetOk("ro_group_id"); ok {
		roGroupId = v.(string)
	}

	d.SetId(instanceId + tccommon.FILED_SP + roGroupId)

	return resourceTencentCloudMysqlRoGroupUpdate(d, meta)
}

func resourceTencentCloudMysqlRoGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	roGroupId := idSplit[1]

	roGroup, err := service.DescribeMysqlRoGroupById(ctx, instanceId, roGroupId)
	if err != nil {
		return err
	}

	if roGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlRoGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if roGroup.RoGroupId != nil {
		_ = d.Set("ro_group_id", roGroup.RoGroupId)
	}

	if roGroup != nil {
		roGroupInfoMap := map[string]interface{}{}

		if roGroup.RoGroupName != nil {
			roGroupInfoMap["ro_group_name"] = roGroup.RoGroupName
		}

		if roGroup.RoMaxDelayTime != nil {
			roGroupInfoMap["ro_max_delay_time"] = roGroup.RoMaxDelayTime
		}

		if roGroup.RoOfflineDelay != nil {
			roGroupInfoMap["ro_offline_delay"] = roGroup.RoOfflineDelay
		}

		if roGroup.MinRoInGroup != nil {
			roGroupInfoMap["min_ro_in_group"] = roGroup.MinRoInGroup
		}

		if roGroup.WeightMode != nil {
			roGroupInfoMap["weight_mode"] = roGroup.WeightMode
		}

		_ = d.Set("ro_group_info", []interface{}{roGroupInfoMap})
	}

	if roGroup.RoInstances != nil {
		roWeightValuesList := []interface{}{}
		for _, roWeightValues := range roGroup.RoInstances {
			roWeightValuesMap := map[string]interface{}{}

			if roWeightValues.InstanceId != nil {
				roWeightValuesMap["instance_id"] = roWeightValues.InstanceId
			}

			if roWeightValues.Weight != nil {
				roWeightValuesMap["weight"] = roWeightValues.Weight
			}

			roWeightValuesList = append(roWeightValuesList, roWeightValuesMap)
		}

		_ = d.Set("ro_weight_values", roWeightValuesList)

	}

	return nil
}

func resourceTencentCloudMysqlRoGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := mysql.NewModifyRoGroupInfoRequest()
	response := mysql.NewModifyRoGroupInfoResponse()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	roGroupId := idSplit[1]

	request.RoGroupId = &roGroupId

	if d.HasChange("ro_group_info") {
		if dMap, ok := helper.InterfacesHeadMap(d, "ro_group_info"); ok {
			roGroupAttr := mysql.RoGroupAttr{}
			if v, ok := dMap["ro_group_name"]; ok {
				roGroupAttr.RoGroupName = helper.String(v.(string))
			}
			if v, ok := dMap["ro_max_delay_time"]; ok {
				roGroupAttr.RoMaxDelayTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["ro_offline_delay"]; ok {
				roGroupAttr.RoOfflineDelay = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["min_ro_in_group"]; ok {
				roGroupAttr.MinRoInGroup = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["weight_mode"]; ok {
				roGroupAttr.WeightMode = helper.String(v.(string))
			}
			if v, ok := dMap["replication_delay_time"]; ok {
				roGroupAttr.ReplicationDelayTime = helper.IntInt64(v.(int))
			}
			request.RoGroupInfo = &roGroupAttr
		}
	}

	if d.HasChange("ro_weight_values") {
		if v, ok := d.GetOk("ro_weight_values"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				roWeightValue := mysql.RoWeightValue{}
				if v, ok := dMap["instance_id"]; ok {
					roWeightValue.InstanceId = helper.String(v.(string))
				}
				if v, ok := dMap["weight"]; ok {
					roWeightValue.Weight = helper.IntInt64(v.(int))
				}
				request.RoWeightValues = append(request.RoWeightValues, &roWeightValue)
			}
		}
	}

	if d.HasChange("is_balance_ro_load") {
		if v, ok := d.GetOkExists("is_balance_ro_load"); ok {
			request.IsBalanceRoLoad = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyRoGroupInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql roGroup failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s create mysql rollback status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s create mysql rollback status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql rollback fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlRoGroupRead(d, meta)
}

func resourceTencentCloudMysqlRoGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
