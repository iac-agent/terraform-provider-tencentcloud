package controlcenter

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	controlcenterv20230110 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/controlcenter/v20230110"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudControlcenterAccountFactoryBaselineConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudControlcenterAccountFactoryBaselineConfigCreate,
		Read:   resourceTencentCloudControlcenterAccountFactoryBaselineConfigRead,
		Update: resourceTencentCloudControlcenterAccountFactoryBaselineConfigUpdate,
		Delete: resourceTencentCloudControlcenterAccountFactoryBaselineConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Baseline 名称，其中 必须 是 唯一. Supports 仅 English letters，numbers，Chinese 字符，和 symbols @，&，_，[]，-. Combination 的 1-25 Chinese 或 English 字符。",
			},

			"baseline_config_items": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Baseline 配置，overwrite update. You 可以 查询 existing baseline configurations via controlcenter:GetAccountFactoryBaseline. You 可以 查询 支持 baseline lists via controlcenter:ListAccountFactoryBaselineItems。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "指定unique identifier 对于 账号 factory baseline item，可以 仅 contain `english letters`，`digits`，和 `@,._[]-:()()[]+=.`，使用 长度 的 2-128 字符。",
						},
						"configuration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "账号 factory baseline item 配置，different baseline items have different 配置 参数。",
						},
						"apply_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "指定number 的 accounts 对于 baseline applications。",
						},
					},
				},
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间。",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "更新时间。",
			},
		},
	}
}

func resourceTencentCloudControlcenterAccountFactoryBaselineConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_controlcenter_account_factory_baseline_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudControlcenterAccountFactoryBaselineConfigUpdate(d, meta)
}

func resourceTencentCloudControlcenterAccountFactoryBaselineConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_controlcenter_account_factory_baseline_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ControlcenterService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeControlcenterAccountFactoryBaselineConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_controlcenter_account_factory_baseline_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.BaselineConfigItems != nil && len(respData.BaselineConfigItems) > 0 {
		baselineConfigItemsList := make([]map[string]interface{}, 0, len(respData.BaselineConfigItems))
		for _, baselineConfigItems := range respData.BaselineConfigItems {
			baselineConfigItemsMap := map[string]interface{}{}
			if baselineConfigItems.Identifier != nil {
				baselineConfigItemsMap["identifier"] = baselineConfigItems.Identifier
			}

			if baselineConfigItems.Configuration != nil {
				baselineConfigItemsMap["configuration"] = baselineConfigItems.Configuration
			}

			if baselineConfigItems.ApplyCount != nil {
				baselineConfigItemsMap["apply_count"] = baselineConfigItems.ApplyCount
			}

			baselineConfigItemsList = append(baselineConfigItemsList, baselineConfigItemsMap)
		}

		_ = d.Set("baseline_config_items", baselineConfigItemsList)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudControlcenterAccountFactoryBaselineConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_controlcenter_account_factory_baseline_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = controlcenterv20230110.NewUpdateAccountFactoryBaselineRequest()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("baseline_config_items"); ok {
		for _, item := range v.(*schema.Set).List() {
			baselineConfigItemsMap := item.(map[string]interface{})
			baselineConfigItem := controlcenterv20230110.BaselineConfigItem{}
			if v, ok := baselineConfigItemsMap["identifier"].(string); ok && v != "" {
				baselineConfigItem.Identifier = helper.String(v)
			}

			if v, ok := baselineConfigItemsMap["configuration"].(string); ok && v != "" {
				baselineConfigItem.Configuration = helper.String(v)
			}

			request.BaselineConfigItems = append(request.BaselineConfigItems, &baselineConfigItem)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseControlcenterV20230110Client().UpdateAccountFactoryBaselineWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update controlcenter account factory baseline config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudControlcenterAccountFactoryBaselineConfigRead(d, meta)
}

func resourceTencentCloudControlcenterAccountFactoryBaselineConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_controlcenter_account_factory_baseline_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
