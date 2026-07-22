package crs

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisParamTemplate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudRedisParamTemplateRead,
		Create: resourceTencentCloudRedisParamTemplateCreate,
		Update: resourceTencentCloudRedisParamTemplateUpdate,
		Delete: resourceTencentCloudRedisParamTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Parameter 模板名称",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parameter 模板描述",
			},
			"product_type": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"template_id"},
				Description:   "指定product 类型 有效值：1 (Redis 2.8 Memory Edition 在 集群 architecture)，2 (Redis 2.8 Memory Edition 在 standard architecture)，3 (CKV 3.2 Memory Edition 在 standard architecture)，4 (CKV 3.2 Memory Edition 在 集群 architecture)，5 (Redis 2.8 Memory Edition 在 standalone architecture)，6 (Redis 4.0 Memory Edition 在 standard architecture)，7 (Redis 4.0 Memory Edition 在 集群 architecture)，8 (Redis 5.0 Memory Edition 在 standard architecture)，9 (Redis 5.0 Memory Edition 在 集群 architecture). 如果 `template_id` 是 指定，此 参数 可以 是 left blank; otherwise，它 为必填项。",
			},
			"template_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"product_type"},
				Description:   "指定which existed template import 从。",
			},
			"params_override": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "指定override 参数 列表，NOTE: Do 不 remove override params once 集合，removing 将 不 take effects 到 当前值",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter 键 e.g. `超时`，check https://www.tencentcloud.com/document/product/239/39796 对于 more reference。",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter 值，check https://www.tencentcloud.com/document/product/239/39796 对于 more reference。",
						},
					},
				},
			},
			"param_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Readonly full 参数 列表 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 键 名称",
						},
						"param_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 类型",
						},
						"default": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "默认值",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 描述",
						},
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "当前值",
						},
						"need_reboot": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "表示是否reboot redis 实例 如果 modified。",
						},
						"max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Maximum 值",
						},
						"min": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Minimum 值",
						},
						"enum_value": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Enum 值。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudRedisParamTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_param_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := RedisService{client: client}

	var (
		request = redis.NewCreateParamTemplateRequest()
		id      string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_type"); ok {
		request.ProductType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("template_id"); ok {
		request.TemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("params_override"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			instanceParam := redis.InstanceParam{}
			if v, ok := dMap["key"]; ok {
				instanceParam.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				instanceParam.Value = helper.String(v.(string))
			}

			request.ParamList = append(request.ParamList, &instanceParam)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		resTemplateId, e := service.CreateParamTemplate(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		id = resTemplateId
		return nil
	})

	if err != nil {
		return err
	}

	if id == "" {
		return fmt.Errorf("cannot get redis template id")
	}

	d.SetId(id)
	return resourceTencentCloudRedisParamTemplateRead(d, meta)
}

func resourceTencentCloudRedisParamTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_param_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	paramTemplate, err := service.DescribeParamTemplateInfo(ctx, id)

	if err != nil {
		return err
	}

	if paramTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `param_template` %s does not exist", id)
	}

	d.SetId(id)

	if paramTemplate.Name != nil {
		_ = d.Set("name", paramTemplate.Name)
	}

	if paramTemplate.Description != nil {
		_ = d.Set("description", paramTemplate.Description)
	}

	if _, ok := d.GetOk("product_type"); ok && paramTemplate.ProductType != nil {
		_ = d.Set("product_type", paramTemplate.ProductType)
	}

	if len(paramTemplate.Items) > 0 {
		result := make([]interface{}, 0)
		for i := range paramTemplate.Items {
			item := paramTemplate.Items[i]
			listMap := map[string]interface{}{
				"name":          item.Name,
				"param_type":    item.ParamType,
				"default":       item.Default,
				"description":   item.Description,
				"current_value": item.CurrentValue,
				"need_reboot":   item.NeedReboot,
				"max":           item.Max,
				"min":           item.Min,
				"enum_value":    item.EnumValue,
			}

			result = append(result, listMap)
		}
		_ = d.Set("param_details", result)
	}

	return nil
}

func resourceTencentCloudRedisParamTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_param_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	request := redis.NewModifyParamTemplateRequest()
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := RedisService{client: client}

	templateId := d.Id()

	request.TemplateId = &templateId

	if d.HasChange("name") {
		request.Name = helper.String(d.Get("name").(string))
	}

	if d.HasChange("description") {
		request.Description = helper.String(d.Get("description").(string))
	}

	if d.HasChange("product_type") {
		return fmt.Errorf("`product_type` do not support change now.")
	}

	if d.HasChange("template_id") {
		return fmt.Errorf("`template_id` do not support change now.")
	}

	if d.HasChange("params_override") {
		if v, ok := d.GetOk("params_override"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				request.ParamList = append(request.ParamList, &redis.InstanceParam{
					Key:   helper.String(dMap["key"].(string)),
					Value: helper.String(dMap["value"].(string)),
				})
			}
		}
	}

	err := service.ModifyParamTemplate(ctx, request)

	if err != nil {
		return err
	}

	return resourceTencentCloudRedisParamTemplateRead(d, meta)
}

func resourceTencentCloudRedisParamTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_param_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	request := redis.NewDeleteParamTemplateRequest()
	request.TemplateId = &id
	if err := service.DeleteParamTemplate(ctx, request); err != nil {
		return err
	}

	return nil
}
