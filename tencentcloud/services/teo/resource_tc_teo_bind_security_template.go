package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoBindSecurityTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoBindSecurityTemplateCreate,
		Read:   resourceTencentCloudTeoBindSecurityTemplateRead,
		Delete: resourceTencentCloudTeoBindSecurityTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "站点 ID 策略 template 到 是 bound 到 或 unbound 从。",
			},

			"entity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "列表 域名 names 到 bind 到/unbind 从 策略 template。",
			},

			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "指定ID 的 策略 template 或 site 全局 策略 到 是 bound 或 unbound.\n<li>To bind 到 策略 template，或 unbind 从 它，指定policy 模板 ID</li>.\n<li>To bind 到 site's 全局 策略，或 unbind 从 它，使用 @ZoneLevel@域名 参数 值</li>.\n\nNote: After unbinding， 域名 名称 将 使用 independent 策略 和 规则 配额 将 是 calculated separately. Please make sure there 是 sufficient 规则 配额 before unbinding。",
			},

			"operate": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Unbind operation 选项. 有效值：`unbind-keep-策略`: unbind 域名 名称 从 策略 template while retaining 当前 策略. `unbind-使用-默认值`: unbind 域名 名称 从 策略 template 和 使用 默认值 blank 策略. 默认值：`unbind-keep-策略`。",
			},

			"over_write": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "如果 passed-在 域名 是 already bound 到 策略 template (包括 site-级别 protection policies)，setting 此 参数 表示是否replace 该 template. 默认值为 true. Supported 值 是: `true`: Replace currently bound template 对于 域名 `false`: Do 不 replace currently bound template 对于 域名 注意: 当 集合 到 false，如果 passed-在 域名 是 already bound 到 策略 template， API 将 返回 错误; site-级别 protection policies 是 also 类型 策略 template。",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "实例 配置 delivery 状态， possible 值 是: `online`: 配置 has taken effect; `fail`: 配置 failed; `process`: 配置 是 being delivered。",
			},
		},
	}
}

func resourceTencentCloudTeoBindSecurityTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_security_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	var (
		zoneId     string
		templateId string
		entity     string
	)

	request := teov20220901.NewBindSecurityTemplateToEntityRequest()

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
		request.Entities = append(request.Entities, helper.String(v.(string)))
	}

	if v, ok := d.GetOk("template_id"); ok {
		templateId = v.(string)
		request.TemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("over_write"); ok {
		request.OverWrite = helper.Bool(v.(bool))
	} else {
		request.OverWrite = helper.Bool(true)
	}

	request.Operate = helper.String("bind")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().BindSecurityTemplateToEntityWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo bind security template failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourceTeoBindSecurityTemplateCreateStateRefreshFunc_0_0(ctx, zoneId, templateId, entity),
		Target:     []string{"online"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{zoneId, templateId, entity}, tccommon.FILED_SP))

	return resourceTencentCloudTeoBindSecurityTemplateRead(d, meta)
}

func resourceTencentCloudTeoBindSecurityTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_security_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	templateId := idSplit[1]
	entity := idSplit[2]

	_ = d.Set("zone_id", zoneId)

	_ = d.Set("template_id", templateId)

	_ = d.Set("entity", entity)

	respData, err := service.DescribeTeoBindSecurityTemplateById(ctx, zoneId, templateId, entity)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_bind_security_template` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if v, ok := d.GetOk("operate"); ok {
		_ = d.Set("operate", v.(string))
	} else {
		_ = d.Set("operate", "unbind-keep-policy")
	}

	if v, ok := d.GetOkExists("over_write"); ok {
		_ = d.Set("over_write", v.(bool))
	} else {
		_ = d.Set("over_write", true)
	}

	return nil
}

func resourceTencentCloudTeoBindSecurityTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_security_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	templateId := idSplit[1]
	entity := idSplit[2]

	request := teov20220901.NewBindSecurityTemplateToEntityRequest()
	request.ZoneId = &zoneId
	request.Entities = append(request.Entities, &entity)
	request.TemplateId = &templateId

	if v, ok := d.GetOk("operate"); ok {
		request.Operate = helper.String(v.(string))
	} else {
		request.Operate = helper.String("unbind-keep-policy")
	}

	if v, ok := d.GetOkExists("over_write"); ok {
		request.OverWrite = helper.Bool(v.(bool))
	} else {
		request.OverWrite = helper.Bool(true)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().BindSecurityTemplateToEntityWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo bind security template failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
