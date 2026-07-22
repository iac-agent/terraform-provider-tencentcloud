package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlBackupList() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentMysqlBackupListRead,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例 ID，例如“cdb-c1nl9rpv”。它与数据库控制台页面中显示的实例 ID 相同。",
			},
			"max_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 10000),
				Description: "列出最新的文件，范围从 1 到 10000。默认值为“10”。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "MySQL 备份列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备份开始的最早时间。例如，“2”表示凌晨 2:00。",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备份完成的时间。",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "备份文件的大小。",
						},
						"backup_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "备份任务ID。",
						},
						"backup_model": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备份方法。支持的值包括：“physical”- 物理备份和“逻辑”- 逻辑备份。",
						},
						"intranet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "内部下载的 URL。",
						},
						"internet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "外部下载的 URL。",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备份文件的所有者。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMysqlBackupListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_backup_list.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	max_number, _ := d.Get("max_number").(int)
	backInfoItems, err := mysqlService.DescribeBackupsByMysqlId(ctx, d.Get("mysql_id").(string), int64(max_number))

	if err != nil {
		return fmt.Errorf("api[DescribeBackups]fail, return %s", err.Error())
	}

	var itemShemas []map[string]interface{}
	var ids = make([]string, len(backInfoItems))

	for index, item := range backInfoItems {
		mapping := map[string]interface{}{
			"time":         *item.Date,
			"finish_time":  *item.FinishTime,
			"size":         *item.Size,
			"backup_id":    *item.BackupId,
			"backup_model": *item.Type,
			"intranet_url": strings.Replace(*item.IntranetUrl, "\u0026", "&", -1),
			"internet_url": strings.Replace(*item.InternetUrl, "\u0026", "&", -1),
			"creator":      *item.Creator,
		}
		ids[index] = fmt.Sprintf("%d", *item.BackupId)
		itemShemas = append(itemShemas, mapping)
	}

	if err := d.Set("list", itemShemas); err != nil {
		log.Printf("[CRITAL]%s provider set itemShemas fail, reason:%s\n ", logId, err.Error())
	}
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {

		if err := tccommon.WriteToFile(output.(string), itemShemas); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail,  reason[%s]\n",
				logId, output.(string), err.Error())
		}

	}
	return nil
}
