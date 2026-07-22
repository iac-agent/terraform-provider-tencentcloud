package wedata

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataRuleTemplateCreate,
		Read:   resourceTencentCloudWedataRuleTemplateRead,
		Update: resourceTencentCloudWedataRuleTemplateUpdate,
		Delete: resourceTencentCloudWedataRuleTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "项目 ID",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "模板 类型 `1` 表示 System template，`2` 表示 Custom template。",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板名称",
			},

			"quality_dim": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Quality inspection dimensions. `1`: Accuracy，`2`: Uniqueness，`3`: Completeness，`4`: Consistency，`5`: Timeliness，`6`: Effectiveness。",
			},

			"source_object_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "来源 数据 对象 类型 `1`: Constant，`2`: Offline 表 级别，`3`: Offline 字段 级别",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "描述 模板。",
			},

			"source_engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "引擎 类型 corresponding 到 来源 `2`: hive,`4`: spark，`16`: dlc。",
			},

			"multi_source_flag": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否associate other 库 tables。",
			},

			"sql_expression": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "SQL Expression。",
			},

			"where_flag": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "如果 add 其中。",
			},
		},
	}
}

func resourceTencentCloudWedataRuleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_rule_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = wedata.NewCreateRuleTemplateRequest()
		response       = wedata.NewCreateRuleTemplateResponse()
		ruleTemplateId uint64
		projectId      string
	)

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("quality_dim"); ok {
		request.QualityDim = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("source_object_type"); ok {
		request.SourceObjectType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_engine_types"); ok {
		sourceEngineTypesSet := v.(*schema.Set).List()
		for i := range sourceEngineTypesSet {
			sourceEngineTypes := sourceEngineTypesSet[i].(int)
			request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
		}
	}

	if v, ok := d.GetOkExists("multi_source_flag"); ok {
		request.MultiSourceFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("sql_expression"); ok {
		request.SqlExpression = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("where_flag"); ok {
		request.WhereFlag = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().CreateRuleTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata ruleTemplate failed, reason:%+v", logId, err)
		return err
	}

	ruleTemplateId = *response.Response.Data
	d.SetId(projectId + tccommon.FILED_SP + helper.UInt64ToStr(ruleTemplateId))

	return resourceTencentCloudWedataRuleTemplateRead(d, meta)
}

func resourceTencentCloudWedataRuleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_rule_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	ruleTemplateId := idSplit[1]

	ruleTemplate, err := service.DescribeWedataRuleTemplateById(ctx, projectId, ruleTemplateId)
	if err != nil {
		return err
	}

	if ruleTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataRuleTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ruleTemplate.Type != nil {
		_ = d.Set("type", ruleTemplate.Type)
	}

	if ruleTemplate.Name != nil {
		_ = d.Set("name", ruleTemplate.Name)
	}

	if ruleTemplate.QualityDim != nil {
		_ = d.Set("quality_dim", ruleTemplate.QualityDim)
	}

	if ruleTemplate.SourceObjectType != nil {
		_ = d.Set("source_object_type", ruleTemplate.SourceObjectType)
	}

	if ruleTemplate.Description != nil {
		_ = d.Set("description", ruleTemplate.Description)
	}

	if ruleTemplate.SourceEngineTypes != nil {
		_ = d.Set("source_engine_types", ruleTemplate.SourceEngineTypes)
	}

	if ruleTemplate.MultiSourceFlag != nil {
		_ = d.Set("multi_source_flag", ruleTemplate.MultiSourceFlag)
	}

	if ruleTemplate.SqlExpression != nil {
		_ = d.Set("sql_expression", base64.StdEncoding.EncodeToString([]byte(*ruleTemplate.SqlExpression)))
	}

	if ruleTemplate.WhereFlag != nil {
		_ = d.Set("where_flag", ruleTemplate.WhereFlag)
	}

	return nil
}

func resourceTencentCloudWedataRuleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_rule_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := wedata.NewModifyRuleTemplateRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	ruleTemplateId := idSplit[1]

	needChange := false
	mutableArgs := []string{
		"type", "name", "quality_dim", "source_object_type",
		"description", "source_engine_types", "multi_source_flag",
		"sql_expression", "where_flag",
	}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request.ProjectId = helper.String(projectId)
		request.TemplateId = helper.StrToUint64Point(ruleTemplateId)

		if v, ok := d.GetOkExists("type"); ok {
			request.Type = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("quality_dim"); ok {
			request.QualityDim = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("source_object_type"); ok {
			request.SourceObjectType = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("source_engine_types"); ok {
			sourceEngineTypesSet := v.(*schema.Set).List()
			for i := range sourceEngineTypesSet {
				sourceEngineTypes := sourceEngineTypesSet[i].(int)
				request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
			}
		}

		if v, ok := d.GetOkExists("multi_source_flag"); ok {
			request.MultiSourceFlag = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("sql_expression"); ok {
			request.SqlExpression = helper.String(v.(string))
		}

		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("where_flag"); ok {
			request.WhereFlag = helper.Bool(v.(bool))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().ModifyRuleTemplate(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update wedata ruleTemplate failed, reason:%+v", logId, err)
			return err
		}
	}
	return resourceTencentCloudWedataRuleTemplateRead(d, meta)
}

func resourceTencentCloudWedataRuleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_rule_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	ruleTemplateId := idSplit[1]

	if err := service.DeleteWedataRuleTemplateById(ctx, projectId, ruleTemplateId); err != nil {
		return err
	}

	return nil
}
