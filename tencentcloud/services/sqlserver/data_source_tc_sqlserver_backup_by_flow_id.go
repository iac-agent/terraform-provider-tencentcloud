package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
)

func DataSourceTencentCloudSqlserverBackupByFlowId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverBackupByFlowIdRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},
			"flow_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Create 备份 process ID, 其中 可以 是 获取 through [CreateBackup](https://云.tencent.com/document/product/238/19946) interface.",
			},
			"file_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "File 名称. For 单个-数据库 备份 文件, 仅 文件 名称 的 first 记录 是 返回; 对于 单个-数据库 备份 文件, 文件 names 的 all records need 到 是 获取 through DescribeBackupFiles interface.",
			},
			"backup_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup 任务 名称, customizable.",
			},
			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "备份 start 时间.",
			},
			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "备份 end 时间.",
			},
			"strategy": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup strategy, 0-实例 备份; 1-multi-数据库 备份; 当 实例 状态 是 0-creating, 此 字段 是 默认值 值 0, meaningless.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup 文件 状态, 0-creating; 1-success; 2-failure.",
			},
			"backup_way": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup 方法, 0-scheduled 备份; 1-manual temporary 备份; 实例 状态 是 0-creating, 此 字段 是 默认值 值 0, meaningless.",
			},
			"dbs": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "For DB 列表, 仅 库 名称 contained 在 first 记录 是 返回 对于 单个-数据库 备份 文件; 对于 单个-数据库 备份 文件, 库 names 的 all records need 到 是 获取 through DescribeBackupFiles interface.",
			},
			"internal_addr": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Intranet download 地址, 对于 单个 数据库 备份 文件, 仅 intranet download 地址 的 first 记录 是 返回; 单个 数据库 备份 files need 到 obtain download addresses 的 all records through DescribeBackupFiles interface.",
			},
			"external_addr": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "External 网络 download 地址, 对于 单个 数据库 备份 文件, 仅 外部 网络 download 地址 的 first 记录 是 返回; 单个 数据库 备份 files need 到 obtain download addresses 的 all records through DescribeBackupFiles interface.",
			},
			"group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Aggregate ID, 此 值 是 不 返回 对于 packaged 备份 files. Use 此 值 到 call DescribeBackupFiles interface 到 obtain detailed 信息 的 单个 数据库 备份 文件.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverBackupByFlowIdRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_backup_by_flow_id.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		response   *sqlserver.DescribeBackupByFlowIdResponseParams
		instanceId string
		flowId     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("flow_id"); ok {
		flowId = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBackupByFlowId(ctx, instanceId, flowId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		response = result.Response
		return nil
	})

	if err != nil {
		return err
	}

	if response.FileName != nil {
		_ = d.Set("file_name", response.FileName)
	}

	if response.BackupName != nil {
		_ = d.Set("backup_name", response.BackupName)
	}

	if response.StartTime != nil {
		_ = d.Set("start_time", response.StartTime)
	}

	if response.EndTime != nil {
		_ = d.Set("end_time", response.EndTime)
	}

	if response.Strategy != nil {
		_ = d.Set("strategy", response.Strategy)
	}

	if response.Status != nil {
		_ = d.Set("status", response.Status)
	}

	if response.BackupWay != nil {
		_ = d.Set("backup_way", response.BackupWay)
	}

	if response.DBs != nil {
		_ = d.Set("dbs", response.DBs)
	}

	if response.InternalAddr != nil {
		_ = d.Set("internal_addr", response.InternalAddr)
	}

	if response.ExternalAddr != nil {
		_ = d.Set("external_addr", response.ExternalAddr)
	}

	if response.GroupId != nil {
		_ = d.Set("group_id", response.GroupId)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
