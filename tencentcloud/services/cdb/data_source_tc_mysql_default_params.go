package cdb

import (
	"context"
	"fmt"
	"math/rand"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlDefaultParams() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTencentCloudMysqlDefaultParamsRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MySQL 数据库版本。",
			},
			//"template_type": {
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//	Description: "",
			//},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"param_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "参数详细信息列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数当前值。",
						},
						"default": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数默认值。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数说明。",
						},
						"enum_value": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "如果参数类型为枚举，则参数可用值。",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "如果参数类型为整数，则参数最大值。",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "如果参数类型为整数，则参数最小值。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数键名称。",
						},
						"need_reboot": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "指示如果参数修改，数据库实例是否会重新启动。",
						},
						"param_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数类型",
						},
					},
				},
			},
		},
	}
}

func datasourceTencentCloudMysqlDefaultParamsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("datasource.tencentcloud_mysql_default_params.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := MysqlService{client: client}

	engineVersion := d.Get("db_version").(string)

	params, err := service.DescribeDefaultParameters(ctx, engineVersion)

	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(fmt.Sprintf("mysql-params-%d", rand.Intn(1000)))

	paramList := make([]map[string]interface{}, 0, len(params))

	for i := range params {
		item := params[i]
		param := map[string]interface{}{
			"current_value": item.CurrentValue,
			"default":       item.Default,
			"description":   item.Description,
			"max":           item.Max,
			"min":           item.Min,
			"name":          item.Name,
			"need_reboot":   item.NeedReboot,
			"param_type":    item.ParamType,
		}

		if item.EnumValue != nil {
			param["enum_value"] = helper.StringsInterfaces(item.EnumValue)
		}
		paramList = append(paramList, param)
	}

	if len(paramList) > 0 {
		err = d.Set("param_list", paramList)
		if err != nil {
			return err
		}
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), params); e != nil {
			return e
		}
	}

	return nil
}
