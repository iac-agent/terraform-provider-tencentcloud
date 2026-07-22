package sqlserver

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverBackups() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentSqlserverBackupsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例 ID.",
			},
			"backup_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 通过 备份 名称, do 不 过滤器 如果 left blank.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start 时间 的 实例 列表, like yyyy-MM-dd HH:mm:ss.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End 时间 的 实例 列表, like yyyy-MM-dd HH:mm:ss.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 的 SQL Server 备份. Each element contains following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 备份.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID.",
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File 名称 的 备份.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start 时间 的 备份.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End 时间 的 备份.",
						},
						"db_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Database 名称 列表 的 备份.",
						},
						"strategy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy 的 备份. `0` 对于 实例 备份, `1` 对于 multi-databases 备份.",
						},
						"trigger_model": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "way 到 触发器 备份. `0` 对于 timed 触发器, `1` 对于 manual 触发器.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status 的 备份. `1` 对于 creating, `2` 对于 successfully 创建, 3 对于 failed.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "大小 的 备份 文件. Unit 是 KB.",
						},
						"intranet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL 对于 downloads internally.",
						},
						"internet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL 对于 downloads externally.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverBackupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_backups.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	backupName := d.Get("backup_name").(string)
	startTime := d.Get("start_time").(string)
	endTime := d.Get("end_time").(string)
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	backInfoItems, err := sqlserverService.DescribeSqlserverBackups(ctx, instanceId, backupName, startTime, endTime)

	if err != nil {
		return fmt.Errorf("api[DescribeBackups]fail, return %s", err.Error())
	}

	var list []map[string]interface{}
	var ids = make([]string, len(backInfoItems))

	for _, item := range backInfoItems {
		mapping := map[string]interface{}{
			"start_time":    item.StartTime,
			"end_time":      item.EndTime,
			"size":          item.Size,
			"trigger_model": item.BackupWay,
			"intranet_url":  item.InternalAddr,
			"internet_url":  item.ExternalAddr,
			"status":        item.Status,
			"file_name":     item.FileName,
			"instance_id":   instanceId,
			"id":            strconv.Itoa(int(*item.Id)),
			"db_list":       helper.StringsInterfaces(item.DBs),
		}
		list = append(list, mapping)
		ids = append(ids, fmt.Sprintf("%d", *item.Id))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil
}
