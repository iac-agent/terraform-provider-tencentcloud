package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainSlowLogUserSqlAdvice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogUserSqlAdviceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"sql_text": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SQL statements.",
			},

			"schema": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "库 名称.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值: `mysql` - 云 数据库 MySQL; `cynosdb` - 云 数据库 TDSQL-C 对于 MySQL; `dbbrain-mysql` - self-built MySQL, 默认值 是 `mysql`.",
			},

			"advices": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SQL optimization suggestion, 其中 可以 是 parsed into JSON 数组, 和 output 是 空 当 无 optimization 是 必填.",
			},

			"comments": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SQL optimization suggestion remarks, 其中 可以 是 parsed into String 数组, 和 output 是 空 当 optimization 是 不 必填.",
			},

			"tables": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "DDL 信息 的 related tables 可以 是 parsed into JSON 数组.",
			},

			"sql_plan": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SQL execution plan 可以 是 parsed into JSON, 和 output 是 空 当 无 optimization 是 必填.",
			},

			"cost": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "费用 saving details after SQL optimization 可以 是 parsed 作为 JSON, 和 output 是 空 当 无 optimization 是 必填.",
			},

			"request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Unique 请求 ID, 返回 对于 every 请求. RequestId 的 请求 needs 到 是 提供 当 locating problem.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainSlowLogUserSqlAdviceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_slow_log_user_sql_advice.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var id string
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		id = v.(string)
		paramMap["instance_id"] = helper.String(id)
	}

	if v, ok := d.GetOk("sql_text"); ok {
		paramMap["sql_text"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("schema"); ok {
		paramMap["schema"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	var result *dbbrain.DescribeUserSqlAdviceResponseParams
	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		result, e = service.DescribeDbbrainSlowLogUserSqlAdviceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result != nil {
		if result.Advices != nil {
			_ = d.Set("advices", result.Advices)
		}

		if result.Comments != nil {
			_ = d.Set("comments", result.Comments)
		}

		if result.SqlText != nil {
			_ = d.Set("sql_text", result.SqlText)
		}

		if result.Schema != nil {
			_ = d.Set("schema", result.Schema)
		}

		if result.Tables != nil {
			_ = d.Set("tables", result.Tables)
		}

		if result.SqlPlan != nil {
			_ = d.Set("sql_plan", result.SqlPlan)
		}

		if result.Cost != nil {
			_ = d.Set("cost", result.Cost)
		}

	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
