package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoPlan() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudTeoPlanCreate,
		Read:   ResourceTencentCloudTeoPlanRead,
		Update: ResourceTencentCloudTeoPlanUpdate,
		Delete: ResourceTencentCloudTeoPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"plan_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"personal", "basic", "standard", "enterprise"}),
				Description:  "subscription 包 类型， possible 值 是: `personal`: personal 包，prepaid 包; `basic`: basic 包，prepaid 包; `standard`: standard 包，prepaid 包; `enterprise`: enterprise 包，postpaid 包。",
			},

			"prepaid_plan_param": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Subscription prepaid 包 参数. 当 PlanType 是 personal，basic，或 standard，此 参数 为可选项 和 是 用于enter subscription 时长 的 包 和 是否enable automatic renewal. 如果 此 参数 是 不 filled 在， 默认值 subscription 时长 是 1 month 和 automatic renewal 是 不 已启用",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}),
							Description:  "subscription 周期 的 prepaid 包，在 months，使用 possible 值: 1，2，3，4，5，6，7，8，9，10，11，12，24，36. 如果未填写 在， 默认值 1 是 使用。",
						},
						"renew_flag": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"on", "off"}),
							Description:  "automatic renewal flag 的 prepaid 包， 值 是: `在`: turn 在 automatic renewal; `关闭`: do 不 turn 在 automatic renewal. 如果未填写 在， 默认值 关闭 是 使用. 当 automatic renewal occurs， 默认值 renewal 周期 是 1 month。",
						},
					},
				},
			},

			// computed
			"plan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Plan ID。",
			},

			"area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service area，possible 值 是: <li>mainland: Mainland China; </li><li>overseas: Worldwide (excluding Mainland China); </li><li>全局: Worldwide (包括 Mainland China). </li>。",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Package 状态， 值 是: <li>normal: normal 状态; </li><li>expiring-soon: about 到 expire; </li><li>expired: expired; </li><li>isolated: isolated; </li><li>overdue-isolated: overdue isolated. </li>。",
			},

			"pay_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Payment 类型，possible 值: <li>0: post-payment; </li><li>1: pre-payment. </li>。",
			},

			"enabled_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "时间 当 包 takes effect。",
			},

			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "expiration date 的 包。",
			},
		},
	}
}

func ResourceTencentCloudTeoPlanCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teov20220901.NewCreatePlanRequest()
		response = teov20220901.NewCreatePlanResponse()
		planId   string
	)

	if v, ok := d.GetOk("plan_type"); ok {
		request.PlanType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "prepaid_plan_param"); ok {
		prepaidPlanParam := teov20220901.PrepaidPlanParam{}
		if v, ok := dMap["period"].(int); ok && v != 0 {
			prepaidPlanParam.Period = helper.IntInt64(v)
		}

		if v, ok := dMap["renew_flag"].(string); ok && v != "" {
			prepaidPlanParam.RenewFlag = helper.String(v)
		}

		request.PrepaidPlanParam = &prepaidPlanParam
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreatePlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo plan failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo function failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.PlanId == nil {
		return fmt.Errorf("PlanId is nil.")
	}

	planId = *response.Response.PlanId
	d.SetId(planId)
	return ResourceTencentCloudTeoPlanRead(d, meta)
}

func ResourceTencentCloudTeoPlanRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		planId  = d.Id()
	)

	respData, err := service.DescribeTeoPlansById(ctx, planId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_plan` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.PlanType != nil {
		_ = d.Set("plan_type", respData.PlanType)
	}

	if respData.PlanId != nil {
		_ = d.Set("plan_id", respData.PlanId)
	}

	if respData.Area != nil {
		_ = d.Set("area", respData.Area)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.PayMode != nil {
		_ = d.Set("pay_mode", respData.PayMode)
	}

	if respData.EnabledTime != nil {
		_ = d.Set("enabled_time", respData.EnabledTime)
	}

	if respData.ExpiredTime != nil {
		_ = d.Set("expired_time", respData.ExpiredTime)
	}

	return nil
}

func ResourceTencentCloudTeoPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		planId = d.Id()
	)

	if d.HasChange("plan_type") {
		request := teov20220901.NewUpgradePlanRequest()
		if v, ok := d.GetOk("plan_type"); ok {
			request.PlanType = helper.String(v.(string))
		}

		request.PlanId = &planId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().UpgradePlanWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("prepaid_plan_param.0.period") {
		request := teov20220901.NewRenewPlanRequest()
		if v, ok := d.GetOk("period"); ok {
			request.Period = helper.IntInt64(v.(int))
		}

		request.PlanId = &planId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().RenewPlanWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("prepaid_plan_param.0.renew_flag") {
		request := teov20220901.NewModifyPlanRequest()
		if v, ok := d.GetOk("renew_flag"); ok {
			request.RenewFlag = &teov20220901.RenewFlag{
				Switch: helper.String(v.(string)),
			}
		}

		request.PlanId = &planId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyPlanWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return ResourceTencentCloudTeoPlanRead(d, meta)
}

func ResourceTencentCloudTeoPlanDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDestroyPlanRequest()
		planId  = d.Id()
	)

	request.PlanId = &planId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DestroyPlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo plan failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
