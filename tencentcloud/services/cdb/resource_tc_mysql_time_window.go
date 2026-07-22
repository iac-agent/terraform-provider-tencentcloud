package cdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlTimeWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlTimeWindowCreate,
		Read:   resourceTencentCloudMysqlTimeWindowRead,
		Update: resourceTencentCloudMysqlTimeWindowUpdate,
		Delete: resourceTencentCloudMysqlTimeWindowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID，格式为 cdb-c1nl9rpv 或 cdbro-c1nl9rpv。与腾讯数据库控制台页面显示的实例ID相同。",
			},

			"time_ranges": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "修改后可进行维护的时间段格式为10:00-12:00。每个时段持续半小时至三小时，开始时间和结束时间按半小时对齐。最多可以设置两个时间段。开始和结束时间范围：[00:00，24:00]。",
			},

			"weekdays": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "指定修改时间段的哪一天。取值范围：周一、周二、周三、周四、周五、周六、周日。如果不指定或留空，则默认每天修改时间段。",
			},

			"max_delay_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "数据延迟阈值。仅对源实例和灾备实例生效。默认值：10。",
			},
		},
	}
}

func resourceTencentCloudMysqlTimeWindowCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_time_window.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(d.Get("instance_id").(string))

	return resourceTencentCloudMysqlTimeWindowUpdate(d, meta)
}

func resourceTencentCloudMysqlTimeWindowRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_time_window.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	timeWindow, err := service.DescribeMysqlTimeWindowById(ctx, instanceId)
	if err != nil {
		return err
	}

	if timeWindow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_mysql_time_window` [%s] not found, please check if it has been deleted.",
			logId, instanceId,
		)
		return nil
	}

	var timeRanges []*string
	var weekdays []*string

	_ = d.Set("instance_id", instanceId)

	if *timeWindow.Response.Monday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Monday
		weekdays = append(weekdays, helper.String("monday"))
	}
	if *timeWindow.Response.Tuesday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Tuesday
		weekdays = append(weekdays, helper.String("tuesday"))
	}
	if *timeWindow.Response.Wednesday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Wednesday
		weekdays = append(weekdays, helper.String("wednesday"))
	}
	if *timeWindow.Response.Thursday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Thursday
		weekdays = append(weekdays, helper.String("thursday"))
	}
	if *timeWindow.Response.Friday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Friday
		weekdays = append(weekdays, helper.String("friday"))
	}
	if *timeWindow.Response.Saturday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Saturday
		weekdays = append(weekdays, helper.String("saturday"))
	}
	if *timeWindow.Response.Wednesday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Wednesday
		weekdays = append(weekdays, helper.String("wednesday"))
	}

	if timeRanges != nil {
		_ = d.Set("time_ranges", timeRanges)
	}

	if weekdays != nil {
		_ = d.Set("weekdays", weekdays)
	}

	if timeWindow.Response.MaxDelayTime != nil {
		_ = d.Set("max_delay_time", timeWindow.Response.MaxDelayTime)
	}

	return nil
}

func resourceTencentCloudMysqlTimeWindowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_time_window.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mysql.NewModifyTimeWindowRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("time_ranges"); ok {
		timeRangesSet := v.(*schema.Set).List()
		for i := range timeRangesSet {
			timeRange := timeRangesSet[i].(string)
			request.TimeRanges = append(request.TimeRanges, &timeRange)
		}
	}

	if v, ok := d.GetOk("weekdays"); ok {
		weekdaysSet := v.(*schema.Set).List()
		for i := range weekdaysSet {
			weekday := weekdaysSet[i].(string)
			request.Weekdays = append(request.Weekdays, &weekday)
		}
	}

	if d.HasChange("max_delay_time") {
		if v, _ := d.GetOk("max_delay_time"); v != nil {
			request.MaxDelayTime = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyTimeWindow(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		d.SetId("")
		log.Printf("[CRITAL]%s update mysql timeWindow failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlTimeWindowRead(d, meta)
}

func resourceTencentCloudMysqlTimeWindowDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_time_window.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId := d.Id()

	if err := service.DeleteMysqlTimeWindowById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
