package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TencentCloudMysqlParameterDetail() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"parameter_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "参数名称。",
		},
		"parameter_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "参数类型。",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "参数规格说明。",
		},
		"current_value": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "当前值。",
		},
		"default_value": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "默认值。",
		},
		"enum_value": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "枚举值。",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"max": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "参数的最大值。",
		},
		"min": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "参数的最小值。",
		},
		"need_reboot": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "指示是否需要重新启动才能启用新参数。",
		},
	}
}

func DataSourceTencentCloudMysqlParameterList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMysqlParameterListRead,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例ID。",
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"5.1", "5.5", "5.6", "5.7", "8.0"}),
				Description: "要使用的数据库引擎的版本号。支持的版本包括5.5/5.6/5.7/8.0。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			"parameter_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "参数列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: TencentCloudMysqlParameterDetail(),
				},
			},
		},
	}
}

func dataSourceTencentMysqlParameterListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_parameter_list.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	mysqlService := MysqlService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var parameterDetails []*cdb.ParameterDetail
	var err error
	instanceIdString := ""
	engineVersionString := ""
	if instanceId, ok := d.GetOk("mysql_id"); ok {
		instanceIdString = instanceId.(string)
		parameterDetails, err = mysqlService.DescribeInstanceParameters(ctx, instanceIdString)
	} else if engineVersion, ok := d.GetOk("engine_version"); ok {
		engineVersionString = engineVersion.(string)
		parameterDetails, err = mysqlService.DescribeDefaultParameters(ctx, engineVersionString)
	} else {
		return fmt.Errorf("mysql_id and engine_version cannot be empty at the same time")
	}
	if err != nil {
		return fmt.Errorf("api[DescribeParameters]fail, return %s", err.Error())
	}

	parameterList := make([]map[string]interface{}, 0, len(parameterDetails))
	for _, item := range parameterDetails {
		mapping := map[string]interface{}{
			"parameter_name": *item.Name,
			"parameter_type": *item.ParamType,
			"description":    *item.Description,
			"current_value":  *item.CurrentValue,
			"default_value":  *item.Default,
			"enum_value":     item.EnumValue,
			"max":            *item.Max,
			"min":            *item.Min,
			"need_reboot":    *item.NeedReboot,
		}
		parameterList = append(parameterList, mapping)
	}
	ids := make([]string, 3)
	ids[0] = "DescribeParameter"
	ids[1] = instanceIdString
	ids[2] = engineVersionString
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("parameter_list", parameterList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set parameter list fail, reason:%s\n ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), parameterList); err != nil {
			return err
		}
	}
	return nil
}
