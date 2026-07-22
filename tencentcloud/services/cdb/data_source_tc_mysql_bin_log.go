package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlBinLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlBinLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例ID，格式为：cdb-c1nl9rpv。与云数据库控制台页面显示的实例ID相同。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "符合查询条件的二进制日志文件的详细信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "binlog日志备份文件名。",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "备份文件大小，单位：Byte。",
						},
						"date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "文件存储时间，时间格式：2016-03-17 02:10:37。",
						},
						"intranet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "下载链接。",
						},
						"internet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "下载链接。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "具体的日志类型，可能的值有： binlog - 二进制日志。",
						},
						"binlog_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Binlog文件开始时间。",
						},
						"binlog_finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "binlog 文件截止时间。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "本地binlog文件所在区域。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备份任务状态。可能的值为“SUCCESS”：备份成功，“FAILED”：备份失败，“RUNNING”：备份正在进行。",
						},
						"remote_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Binlog远程备份详细信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub_backup_id": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "远程备份子任务ID。",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "远程备份所在区域。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "备份任务状态。可能的值为“SUCCESS”：备份成功，“FAILED”：备份失败，“RUNNING”：备份正在进行。",
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "远程备份任务的开始时间。",
									},
									"finish_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "远程备份任务结束时间。",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "下载链接。",
									},
								},
							},
						},
						"cos_storage_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "存储方式，0-常规存储，1-归档存储，默认0。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例ID，格式为：cdb-c1nl9rpv。与云数据库控制台页面显示的实例ID相同。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlBinLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_bin_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var binLog []*cdb.BinlogInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlBinLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		binLog = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(binLog))
	tmpList := make([]map[string]interface{}, 0, len(binLog))
	if binLog != nil {
		for _, binlogInfo := range binLog {
			binlogInfoMap := map[string]interface{}{}

			if binlogInfo.Name != nil {
				binlogInfoMap["name"] = binlogInfo.Name
			}

			if binlogInfo.Size != nil {
				binlogInfoMap["size"] = binlogInfo.Size
			}

			if binlogInfo.Date != nil {
				binlogInfoMap["date"] = binlogInfo.Date
			}

			if binlogInfo.IntranetUrl != nil {
				binlogInfoMap["intranet_url"] = binlogInfo.IntranetUrl
			}

			if binlogInfo.InternetUrl != nil {
				binlogInfoMap["internet_url"] = binlogInfo.InternetUrl
			}

			if binlogInfo.Type != nil {
				binlogInfoMap["type"] = binlogInfo.Type
			}

			if binlogInfo.BinlogStartTime != nil {
				binlogInfoMap["binlog_start_time"] = binlogInfo.BinlogStartTime
			}

			if binlogInfo.BinlogFinishTime != nil {
				binlogInfoMap["binlog_finish_time"] = binlogInfo.BinlogFinishTime
			}

			if binlogInfo.Region != nil {
				binlogInfoMap["region"] = binlogInfo.Region
			}

			if binlogInfo.Status != nil {
				binlogInfoMap["status"] = binlogInfo.Status
			}

			if binlogInfo.RemoteInfo != nil {
				remoteInfoList := []interface{}{}
				for _, remoteInfo := range binlogInfo.RemoteInfo {
					remoteInfoMap := map[string]interface{}{}

					if remoteInfo.SubBackupId != nil {
						remoteInfoMap["sub_backup_id"] = []interface{}{remoteInfo.SubBackupId}
					}

					if remoteInfo.Region != nil {
						remoteInfoMap["region"] = remoteInfo.Region
					}

					if remoteInfo.Status != nil {
						remoteInfoMap["status"] = remoteInfo.Status
					}

					if remoteInfo.StartTime != nil {
						remoteInfoMap["start_time"] = remoteInfo.StartTime
					}

					if remoteInfo.FinishTime != nil {
						remoteInfoMap["finish_time"] = remoteInfo.FinishTime
					}

					if remoteInfo.Url != nil {
						remoteInfoMap["url"] = remoteInfo.Url
					}

					remoteInfoList = append(remoteInfoList, remoteInfoMap)
				}

				binlogInfoMap["remote_info"] = remoteInfoList
			}

			if binlogInfo.CosStorageType != nil {
				binlogInfoMap["cos_storage_type"] = binlogInfo.CosStorageType
			}

			if binlogInfo.InstanceId != nil {
				binlogInfoMap["instance_id"] = binlogInfo.InstanceId
			}

			ids = append(ids, *binlogInfo.Name)
			tmpList = append(tmpList, binlogInfoMap)
		}

		err = d.Set("items", tmpList)
		if err != nil {
			return err
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
