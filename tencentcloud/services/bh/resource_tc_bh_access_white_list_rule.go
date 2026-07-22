package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhAccessWhiteListRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhAccessWhiteListRuleCreate,
		Read:   resourceTencentCloudBhAccessWhiteListRuleRead,
		Update: resourceTencentCloudBhAccessWhiteListRuleUpdate,
		Delete: resourceTencentCloudBhAccessWhiteListRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP 地址 10.10.10.1 或 网络 segment 10.10.10.0/24，最小 长度 4 bytes，最大 长度 40 bytes。",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "备注 信息，最小 长度 0 字符，最大 长度 40 字符。",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID 访问 white 列表 规则。",
			},
		},
	}
}

func resourceTencentCloudBhAccessWhiteListRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_access_white_list_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bhv20230418.NewCreateAccessWhiteListRuleRequest()
		response = bhv20230418.NewCreateAccessWhiteListRuleResponse()
		ruleId   string
	)

	if v, ok := d.GetOk("source"); ok {
		request.Source = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().CreateAccessWhiteListRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bh access white list rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh access white list rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Id == nil {
		return fmt.Errorf("Id is nil.")
	}

	ruleId = helper.UInt64ToStr(*response.Response.Id)
	d.SetId(ruleId)
	return resourceTencentCloudBhAccessWhiteListRuleRead(d, meta)
}

func resourceTencentCloudBhAccessWhiteListRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_access_white_list_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ruleId  = d.Id()
	)

	respData, err := service.DescribeBhAccessWhiteListRuleById(ctx, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_access_white_list_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Source != nil {
		_ = d.Set("source", respData.Source)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.Id != nil {
		_ = d.Set("rule_id", respData.Id)
	}

	return nil
}

func resourceTencentCloudBhAccessWhiteListRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_access_white_list_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		ruleId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"source", "remark"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := bhv20230418.NewModifyAccessWhiteListRuleRequest()
		if v, ok := d.GetOk("source"); ok {
			request.Source = helper.String(v.(string))
		}

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		request.Id = helper.StrToUint64Point(ruleId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyAccessWhiteListRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh access white list rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhAccessWhiteListRuleRead(d, meta)
}

func resourceTencentCloudBhAccessWhiteListRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_access_white_list_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewDeleteAccessWhiteListRulesRequest()
		ruleId  = d.Id()
	)

	request.IdSet = append(request.IdSet, helper.StrToUint64Point(ruleId))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DeleteAccessWhiteListRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bh access white list rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
