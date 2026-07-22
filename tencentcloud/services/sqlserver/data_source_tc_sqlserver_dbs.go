package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudSqlserverDBs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverDBRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 store results.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SQL Server 实例 ID 其中 DB belongs 到.",
			},
			// Computed
			"db_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 的 dbs belong 到 特定 实例. Each element contains following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL Server 实例 ID 其中 DB belongs 到.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name 的 DB.",
						},
						"charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Character 集合 DB uses, could 是 `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, 和 `SQL_Latin1_General_CP1_CS_AS`.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark 的 DB.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database creation 时间.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database 状态. 有效 值 是 `creating`, `running`, `modifying`, `dropping`.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverDBRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencent_sqlserver_dbs.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// precheck
	instanceId := d.Get("instance_id").(string)
	_, has, err := sqlserverService.DescribeSqlserverInstanceById(ctx, instanceId)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s DescribeSqlserverInstanceById fail, reason:%s\n", logId, err)
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s SQL Server instance %s dose not exist", logId, instanceId)
	}

	dbInfos, err := sqlserverService.DescribeDBsOfInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	var dbList []map[string]interface{}
	for _, item := range dbInfos {
		var dbInfo = make(map[string]interface{})
		dbInfo["name"] = item.Name
		dbInfo["charset"] = item.Charset
		dbInfo["remark"] = item.Remark
		dbInfo["create_time"] = item.CreateTime
		dbInfo["status"] = SQLSERVER_DB_STATUS[*item.Status]
		dbList = append(dbList, dbInfo)
	}
	_ = d.Set("db_list", dbList)
	d.SetId(instanceId)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), dbList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}
	}
	return nil
}
