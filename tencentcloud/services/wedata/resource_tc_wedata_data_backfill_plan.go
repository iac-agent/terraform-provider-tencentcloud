package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataDataBackfillPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataDataBackfillPlanCreate,
		Read:   resourceTencentCloudWedataDataBackfillPlanRead,
		Delete: resourceTencentCloudWedataDataBackfillPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "项目 ID",
			},

			"task_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "Backfill 任务 collection。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"data_backfill_range_list": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "指定data 时间 配置 对于 backfill 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_date": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Start date 在 yyyy-MM-dd 格式 表示start 从 00:00:00 在 指定 date。",
						},
						"end_date": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "End date 在 格式 yyyy-MM-dd，表示ending 在 23:59:59 的 指定 date。",
						},
						"execution_start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "开始时间 的 each day between [StartDate，EndDate] 在 HH:mm 格式 effective 对于 tasks 使用 周期 的 hours 或 less。",
						},
						"execution_end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "结束时间 point between [StartDate，EndDate] 在 HH:mm 格式 effective 对于 tasks 使用 周期 的 hours 或 less。",
						},
					},
				},
			},

			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "时区，默认值 UTC+8。",
			},

			"data_backfill_plan_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Backfill plan 名称 如果为空， 字符串 的 字符 是 randomly generated 通过 系统。",
			},

			"check_parent_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Check parent 任务 类型 有效值：NONE (do 不 check ALL)，ALL (check ALL upstream parent tasks)，MAKE_SCOPE (仅 check 在 currently selected tasks 的 backfill plan). 默认值：NONE (do 不 check)。",
			},

			"skip_event_listening": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "指定是否ignore 事件 dependency 对于 backfill. 默认值 true。",
			},

			"redefine_self_workflow_dependency": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Custom 工作流 self-dependency. 有效值：yes 或 无. 如果 不 已配置，使用 original 工作流 self-dependency。",
			},

			"redefine_parallel_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Customizes degree 的 并发 对于 实例 running. 如果 without configuring，使用 existing self-dependent 的 任务。",
			},

			"scheduler_resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Scheduling 资源 组 ID. 如果为空，表示usage 的 original 任务 scheduling execution 资源 组。",
			},

			"integration_resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Integration 任务 资源 组 ID. 表示usage 的 original 任务 scheduling execution 资源 组 如果为空。",
			},

			"redefine_param_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom 参数. re-指定task's 参数 到 facilitate execution 的 new logic 通过 replenished 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"k": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "键 名称",
						},
						"v": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "值 do 不 pass SQL ( 请求 将 是 deemed 作为 attack 在 api). 如果 needed，transcode SQL 使用 Base64 和 decode 它。",
						},
					},
				},
			},

			"data_time_order": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Backfill Execution 顺序 - execution 顺序 对于 backfill 实例 based 在 their 数据 时间. Effective 仅 当 both conditions 是 met:\n\n1. Must 是 same cycle 任务.\n\n2. 优先级 是 given 到 dependency 顺序 如果 无 dependencies apply，execution follows 已配置 顺序\n\nValid 值:\n\n-NORMAL: No 特定 顺序 (默认值)\n\n-ORDER: Execute 在 chronological 顺序\n\n-REVERSE: Execute 在 reverse chronological 顺序",
			},

			"redefine_cycle_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backfill 实例 Regeneration Cycle - 如果 集合，此 将 redefine generation cycle 的 backfill 任务 实例. Currently，仅 daily 实例 可以 是 converted into 实例 generated 在 first day 的 each month.\n\nValid 值:\n\nMONTH_CYCLE: Monthly。",
			},

			// computed
			"data_backfill_plan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data backfill plan ID。",
			},
		},
	}
}

func resourceTencentCloudWedataDataBackfillPlanCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_data_backfill_plan.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request            = wedatav20250806.NewCreateDataBackfillPlanRequest()
		response           = wedatav20250806.NewCreateDataBackfillPlanResponse()
		projectId          string
		dataBackfillPlanId string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("task_ids"); ok {
		taskIdsSet := v.(*schema.Set).List()
		for i := range taskIdsSet {
			taskIds := taskIdsSet[i].(string)
			request.TaskIds = append(request.TaskIds, helper.String(taskIds))
		}
	}

	if v, ok := d.GetOk("data_backfill_range_list"); ok {
		for _, item := range v.([]interface{}) {
			dataBackfillRangeListMap := item.(map[string]interface{})
			dataBackfillRange := wedatav20250806.DataBackfillRange{}
			if v, ok := dataBackfillRangeListMap["start_date"].(string); ok && v != "" {
				dataBackfillRange.StartDate = helper.String(v)
			}

			if v, ok := dataBackfillRangeListMap["end_date"].(string); ok && v != "" {
				dataBackfillRange.EndDate = helper.String(v)
			}

			if v, ok := dataBackfillRangeListMap["execution_start_time"].(string); ok && v != "" {
				dataBackfillRange.ExecutionStartTime = helper.String(v)
			}

			if v, ok := dataBackfillRangeListMap["execution_end_time"].(string); ok && v != "" {
				dataBackfillRange.ExecutionEndTime = helper.String(v)
			}

			request.DataBackfillRangeList = append(request.DataBackfillRangeList, &dataBackfillRange)
		}
	}

	if v, ok := d.GetOk("time_zone"); ok {
		request.TimeZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_backfill_plan_name"); ok {
		request.DataBackfillPlanName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("check_parent_type"); ok {
		request.CheckParentType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("skip_event_listening"); ok {
		request.SkipEventListening = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("redefine_self_workflow_dependency"); ok {
		request.RedefineSelfWorkflowDependency = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("redefine_parallel_num"); ok {
		request.RedefineParallelNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("scheduler_resource_group_id"); ok {
		request.SchedulerResourceGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("integration_resource_group_id"); ok {
		request.IntegrationResourceGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("redefine_param_list"); ok {
		for _, item := range v.([]interface{}) {
			redefineParamListMap := item.(map[string]interface{})
			kVPair := wedatav20250806.KVPair{}
			if v, ok := redefineParamListMap["k"].(string); ok && v != "" {
				kVPair.K = helper.String(v)
			}

			if v, ok := redefineParamListMap["v"].(string); ok && v != "" {
				kVPair.V = helper.String(v)
			}

			request.RedefineParamList = append(request.RedefineParamList, &kVPair)
		}
	}

	if v, ok := d.GetOk("data_time_order"); ok {
		request.DataTimeOrder = helper.String(v.(string))
	}

	if v, ok := d.GetOk("redefine_cycle_type"); ok {
		request.RedefineCycleType = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateDataBackfillPlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata data backfill plan operation failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata data backfill plan operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.DataBackfillPlanId == nil {
		return fmt.Errorf("DataBackfillPlanId is nil.")
	}

	dataBackfillPlanId = *response.Response.Data.DataBackfillPlanId
	d.SetId(strings.Join([]string{projectId, dataBackfillPlanId}, tccommon.FILED_SP))
	return resourceTencentCloudWedataDataBackfillPlanRead(d, meta)
}

func resourceTencentCloudWedataDataBackfillPlanRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_data_backfill_plan.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	dataBackfillPlanId := idSplit[1]

	respData, err := service.DescribeWedataDataBackfillPlanById(ctx, projectId, dataBackfillPlanId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_data_backfill_plan` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("data_backfill_plan_id", dataBackfillPlanId)

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	if respData.TaskIds != nil {
		_ = d.Set("task_ids", respData.TaskIds)
	}

	if respData.DataBackfillRangeList != nil {
		dataBackfillRangeListList := []interface{}{}
		for _, dataBackfillRangeList := range respData.DataBackfillRangeList {
			dataBackfillRangeListMap := map[string]interface{}{}
			if dataBackfillRangeList.StartDate != nil {
				dataBackfillRangeListMap["start_date"] = dataBackfillRangeList.StartDate
			}

			if dataBackfillRangeList.EndDate != nil {
				dataBackfillRangeListMap["end_date"] = dataBackfillRangeList.EndDate
			}

			if dataBackfillRangeList.ExecutionStartTime != nil {
				dataBackfillRangeListMap["execution_start_time"] = dataBackfillRangeList.ExecutionStartTime
			}

			if dataBackfillRangeList.ExecutionEndTime != nil {
				dataBackfillRangeListMap["execution_end_time"] = dataBackfillRangeList.ExecutionEndTime
			}

			dataBackfillRangeListList = append(dataBackfillRangeListList, dataBackfillRangeListMap)
		}

		_ = d.Set("data_backfill_range_list", dataBackfillRangeListList)
	}

	if respData.DataBackfillPlanName != nil {
		_ = d.Set("data_backfill_plan_name", respData.DataBackfillPlanName)
	}

	if respData.CheckParentType != nil {
		_ = d.Set("check_parent_type", respData.CheckParentType)
	}

	if respData.SkipEventListening != nil {
		_ = d.Set("skip_event_listening", respData.SkipEventListening)
	}

	if respData.RedefineSelfWorkflowDependency != nil {
		_ = d.Set("redefine_self_workflow_dependency", respData.RedefineSelfWorkflowDependency)
	}

	if respData.RedefineParallelNum != nil {
		_ = d.Set("redefine_parallel_num", respData.RedefineParallelNum)
	}

	if respData.SchedulerResourceGroupId != nil {
		_ = d.Set("scheduler_resource_group_id", respData.SchedulerResourceGroupId)
	}

	if respData.IntegrationResourceGroupId != nil {
		_ = d.Set("integration_resource_group_id", respData.IntegrationResourceGroupId)
	}

	if respData.RedefineParamList != nil {
		redefineParamListList := []interface{}{}
		for _, redefineParamList := range respData.RedefineParamList {
			redefineParamListMap := map[string]interface{}{}
			if redefineParamList.K != nil {
				redefineParamListMap["k"] = redefineParamList.K
			}

			if redefineParamList.V != nil {
				redefineParamListMap["v"] = redefineParamList.V
			}

			redefineParamListList = append(redefineParamListList, redefineParamListMap)
		}

		_ = d.Set("redefine_param_list", redefineParamListList)
	}

	if respData.DataTimeOrder != nil {
		_ = d.Set("data_time_order", respData.DataTimeOrder)
	}

	if respData.RedefineCycleType != nil {
		_ = d.Set("redefine_cycle_type", respData.RedefineCycleType)
	}

	if respData.DataBackfillPlanId != nil {
		_ = d.Set("data_backfill_plan_id", respData.DataBackfillPlanId)
	}

	return nil
}

func resourceTencentCloudWedataDataBackfillPlanDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_data_backfill_plan.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request  = wedatav20250806.NewDeleteDataBackfillPlanAsyncRequest()
		response = wedatav20250806.NewDeleteDataBackfillPlanAsyncResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	dataBackfillPlanId := idSplit[1]

	request.ProjectId = &projectId
	request.DataBackfillPlanId = &dataBackfillPlanId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteDataBackfillPlanAsyncWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete data backfill plan async failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata workflow permissions failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.AsyncId == nil {
		return fmt.Errorf("Delete wedata workflow permissions failed, AsyncId is false")
	}

	// wait
	asyncId := *response.Response.Data.AsyncId
	if _, err := (&resource.StateChangeConf{
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{"INIT", "RUNNING"},
		Refresh:    service.WedataOpsAsyncJobRefresh(ctx, projectId, asyncId),
		Target:     []string{"SUCCESS"},
		Timeout:    900 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}
