package dcdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbLogFiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbLogFilesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID 在 格式 的 `tdsqlshard-ow728lmc`。",
			},

			"shard_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 分片 ID 在 格式 的 `分片-rc754ljk`。",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Requested 日志 类型 有效值：1 (binlog)，2 (cold 备份)，3 (errlog)，4 (slowlog)。",
			},

			"files": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information such 作为 `uri`，`长度`，和 `mtime` (修改时间)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mtime": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最后修改时间 的 日志。",
						},
						"length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File 长度。",
						},
						"uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Uniform 资源 identifier (URI) 使用 during 日志 download。",
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Filename。",
						},
					},
				},
			},

			"vpc_prefix": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "For 实例 在 VPC，此 prefix plus URI 可以 是 使用 作为 download 地址",
			},

			"normal_prefix": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "For 实例 在 common 网络，此 prefix plus URI 可以 是 使用 作为 download 地址",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudDcdbLogFilesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_log_files.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var (
		instanceId string
		shardId    string
		ids        []string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("shard_id"); ok {
		paramMap["ShardId"] = helper.String(v.(string))
		shardId = v.(string)
	}

	if v, _ := d.GetOk("type"); v != nil {
		paramMap["Type"] = helper.IntInt64(v.(int))
	}

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *dcdb.DescribeDBLogFilesResponseParams
	var e error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e = service.DescribeDcdbLogFilesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	ids = append(ids, instanceId)
	ids = append(ids, shardId)
	if result.Files != nil {
		tmpList := make([]interface{}, 0, len(result.Files))
		for _, logFileInfo := range result.Files {
			logFileInfoMap := map[string]interface{}{}

			if logFileInfo.Mtime != nil {
				logFileInfoMap["mtime"] = logFileInfo.Mtime
			}

			if logFileInfo.Length != nil {
				logFileInfoMap["length"] = logFileInfo.Length
			}

			if logFileInfo.Uri != nil {
				logFileInfoMap["uri"] = logFileInfo.Uri
			}

			if logFileInfo.FileName != nil {
				logFileInfoMap["file_name"] = logFileInfo.FileName
			}
			ids = append(ids, *logFileInfo.FileName)
			tmpList = append(tmpList, logFileInfoMap)
		}

		_ = d.Set("files", tmpList)
	}

	if result.VpcPrefix != nil {
		_ = d.Set("vpc_prefix", result.VpcPrefix)
	}

	if result.NormalPrefix != nil {
		_ = d.Set("normal_prefix", result.NormalPrefix)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
