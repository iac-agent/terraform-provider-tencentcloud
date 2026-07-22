package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsMachineGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsMachineGroupCreate,
		Read:   resourceTencentCloudClsMachineGroupRead,
		Delete: resourceTencentCloudClsMachineGroupDelete,
		Update: resourceTencentCloudClsMachineGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "机器组名称，必须唯一。",
			},
			"machine_group_type": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "要创建的机器组的类型。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							Description: "机组类型。有效值： ip：机器组Values中保存采集机器的IP地址；" +
								"label: the tags of the machines are stored in Values of the machine group.",
						},
						"values": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "机器描述列表。",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签描述列表。最多支持 10 个标签键值对，并且必须是唯一的。",
			},
			"auto_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "是否启用机器组自动更新。",
			},
			"update_start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "更新开始时间。我们建议您在非高峰时段更新LogListener。",
			},
			"update_end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "更新结束时间。我们建议您在非高峰时段更新LogListener。",
			},
			"service_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "是否开启服务日志记录LogListener服务本身产生的日志。启用后，内部日志集cls_service_logging和loglistener_status、loglistener_alarm、" +
					"and loglistener_business log topics will be created, which will not incur fees.",
			},
		},
	}
}

func resourceTencentCloudClsMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_machine_group.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cls.NewCreateMachineGroupRequest()
		response *cls.CreateMachineGroupResponse
	)

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("machine_group_type"); ok {
		machineGroupTypes := make([]*cls.MachineGroupTypeInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			machineGroupType := cls.MachineGroupTypeInfo{}
			if v, ok := dMap["type"]; ok {
				machineGroupType.Type = helper.String(v.(string))
			}
			if v, ok := dMap["values"]; ok {
				for _, u := range v.([]interface{}) {
					machineGroupType.Values = append(machineGroupType.Values, helper.String(u.(string)))
				}
			}
			machineGroupTypes = append(machineGroupTypes, &machineGroupType)
		}
		request.MachineGroupType = machineGroupTypes[0]
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &key,
				Value: &value,
			})
		}
	}

	if v, ok := d.GetOk("auto_update"); ok {
		request.AutoUpdate = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("update_start_time"); ok {
		request.UpdateStartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("update_end_time"); ok {
		request.UpdateEndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_logging"); ok {
		request.ServiceLogging = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateMachineGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls machine group failed, reason:%+v", logId, err)
		return err
	}

	id := *response.Response.GroupId
	d.SetId(id)
	return resourceTencentCloudClsMachineGroupRead(d, meta)
}

func resourceTencentCloudClsMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_machine_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	machineGroup, err := service.DescribeClsMachineGroupById(ctx, id)

	if err != nil {
		return err
	}

	if machineGroup == nil {
		d.SetId("")
		return fmt.Errorf("resource `MachineGroup` %s does not exist", id)
	}

	_ = d.Set("group_name", machineGroup.GroupName)

	machineGroupType := make([]map[string]interface{}, 0)

	machineGroupType = append(machineGroupType, map[string]interface{}{
		"type":   *machineGroup.MachineGroupType.Type,
		"values": machineGroup.MachineGroupType.Values,
	})

	_ = d.Set("machine_group_type", machineGroupType)

	tags := make(map[string]string, len(machineGroup.Tags))
	for _, tag := range machineGroup.Tags {
		tags[*tag.Key] = *tag.Value
	}
	_ = d.Set("tags", tags)

	_ = d.Set("auto_update", helper.StrToBool(*machineGroup.AutoUpdate))
	_ = d.Set("update_start_time", machineGroup.UpdateStartTime)
	_ = d.Set("update_end_time", machineGroup.UpdateEndTime)
	_ = d.Set("service_logging", machineGroup.ServiceLogging)

	return nil
}

func resourceTencentCloudClsMachineGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_machine_group.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cls.NewModifyMachineGroupRequest()

	request.GroupId = helper.String(d.Id())
	request.GroupName = helper.String(d.Get("group_name").(string))

	if d.HasChange("machine_group_type") {
		if v, ok := d.GetOk("machine_group_type"); ok {
			machineGroupTypes := make([]*cls.MachineGroupTypeInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				machineGroupType := cls.MachineGroupTypeInfo{}
				if v, ok := dMap["type"]; ok {
					machineGroupType.Type = helper.String(v.(string))
				}
				if v, ok := dMap["values"]; ok {
					for _, u := range v.([]interface{}) {
						machineGroupType.Values = append(machineGroupType.Values, helper.String(u.(string)))
					}
					machineGroupTypes = append(machineGroupTypes, &machineGroupType)
				}
			}
			request.MachineGroupType = machineGroupTypes[0]
		}
	}

	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})
		request.Tags = make([]*cls.Tag, 0, len(tags))
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &key,
				Value: helper.String(value.(string)),
			})
		}
	}

	if d.HasChange("auto_update") || d.HasChange("update_start_time") || d.HasChange("update_end_time") {
		request.AutoUpdate = helper.Bool(d.Get("auto_update").(bool))
		request.UpdateStartTime = helper.String(d.Get("update_start_time").(string))
		request.UpdateEndTime = helper.String(d.Get("update_end_time").(string))
	}

	if d.HasChange("service_logging") {
		request.ServiceLogging = helper.Bool(d.Get("service_logging").(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyMachineGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsMachineGroupRead(d, meta)
}

func resourceTencentCloudClsMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_machine_group.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	if err := service.DeleteClsMachineGroup(ctx, id); err != nil {
		return err
	}

	return nil
}
