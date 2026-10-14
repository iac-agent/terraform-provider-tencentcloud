package config

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigAggregator() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigAggregatorCreate,
		Read:   resourceTencentCloudConfigAggregatorRead,
		Update: resourceTencentCloudConfigAggregatorUpdate,
		Delete: resourceTencentCloudConfigAggregatorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account group name.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account group description.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account group type: `RD` global account group, `CUSTOM` custom account group.",
			},
			"owner_uin": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Account group creator ID.",
			},
			"aggregator_accounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Account group member information list, up to 100.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_uin": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Member ID.",
						},
						"member_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Member name.",
						},
					},
				},
			},
			// computed
			"account_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account group ID.",
			},
			"aggregator_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation status.",
			},
		},
	}
}

func buildConfigAggregatorAccounts(reqAccounts *[]*configv20220802.AggregatorAccount, v interface{}) {
	for _, item := range v.([]interface{}) {
		accountMap := item.(map[string]interface{})
		account := configv20220802.AggregatorAccount{}
		if v, ok := accountMap["member_uin"].(int); ok {
			account.MemberUin = helper.IntUint64(v)
		}

		if v, ok := accountMap["member_name"].(string); ok && v != "" {
			account.MemberName = helper.String(v)
		}

		*reqAccounts = append(*reqAccounts, &account)
	}
}

func resourceTencentCloudConfigAggregatorCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request        = configv20220802.NewCreateAggregatorRequest()
		response       = configv20220802.NewCreateAggregatorResponse()
		accountGroupId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("aggregator_accounts"); ok {
		buildConfigAggregatorAccounts(&request.AggregatorAccounts, v)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().CreateAggregatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create config aggregator failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create config aggregator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.AccountGroupId == nil || *response.Response.AccountGroupId == "" {
		log.Printf("[CRITAL]%s create config aggregator failed, logId=%s, AccountGroupId is empty.", logId, d.Id())
		return fmt.Errorf("Create config aggregator failed, AccountGroupId is empty.")
	}

	accountGroupId = *response.Response.AccountGroupId
	ownerUin := d.Get("owner_uin").(int)
	d.SetId(strings.Join([]string{helper.IntToStr(ownerUin), accountGroupId}, tccommon.FILED_SP))
	return resourceTencentCloudConfigAggregatorRead(d, meta)
}

func resourceTencentCloudConfigAggregatorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = configv20220802.NewDescribeAggregatorRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	ownerUin := idSplit[0]
	accountGroupId := idSplit[1]

	request.AccountGroupId = &accountGroupId
	request.OwnerUin = helper.StrToUint64Point(ownerUin)
	var response *configv20220802.DescribeAggregatorResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().DescribeAggregatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe config aggregator failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read config aggregator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Name == nil && response.Response.Type == nil {
		log.Printf("[CRUD] config_aggregator id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if response.Response.Name != nil {
		_ = d.Set("name", response.Response.Name)
	}

	if response.Response.Description != nil {
		_ = d.Set("description", response.Response.Description)
	}

	if response.Response.Type != nil {
		_ = d.Set("type", response.Response.Type)
	}

	if response.Response.AggregatorAccounts != nil {
		accountsList := make([]map[string]interface{}, 0, len(response.Response.AggregatorAccounts))
		for _, account := range response.Response.AggregatorAccounts {
			accountMap := map[string]interface{}{}
			if account.MemberUin != nil {
				accountMap["member_uin"] = account.MemberUin
			}

			if account.MemberName != nil {
				accountMap["member_name"] = account.MemberName
			}

			accountsList = append(accountsList, accountMap)
		}

		_ = d.Set("aggregator_accounts", accountsList)
	}

	if response.Response.AggregatorStatus != nil {
		_ = d.Set("aggregator_status", response.Response.AggregatorStatus)
	}

	return nil
}

func resourceTencentCloudConfigAggregatorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	ownerUin := idSplit[0]
	accountGroupId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "description", "aggregator_accounts", "owner_uin"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := configv20220802.NewUpdateAggregatorRequest()
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("aggregator_accounts"); ok {
			buildConfigAggregatorAccounts(&request.AggregatorAccounts, v)
		}

		request.AccountGroupId = &accountGroupId
		request.OwnerUin = helper.StrToUint64Point(ownerUin)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateAggregatorWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update config aggregator failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update config aggregator failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudConfigAggregatorRead(d, meta)
}

func resourceTencentCloudConfigAggregatorDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = configv20220802.NewDeleteAggregatorsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	ownerUin := idSplit[0]
	accountGroupId := idSplit[1]

	request.AccountGroupId = &accountGroupId
	request.OwnerUin = helper.StrToUint64Point(ownerUin)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().DeleteAggregatorsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete config aggregator failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete config aggregator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
