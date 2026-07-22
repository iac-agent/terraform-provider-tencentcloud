package dbbrain

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainTopSpaceTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainTopSpaceTablesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     20,
				Description: "数量 的 Top tables 返回, 最大 值 是 100, 和 默认值 是 20.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sorting 字段 使用 到 过滤器 Top 表. 可选 字段 include DataLength, IndexLength, TotalLength, DataFree, FragRatio, TableRows, 和 PhysicalFileSize (仅 支持 通过 ApsaraDB 对于 MySQL 实例). 默认值 对于 ApsaraDB 对于 MySQL 实例 是 PhysicalFileSize, 和 默认值 对于 other product 实例 是 TotalLength.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: mysql - 云 数据库 MySQL, cynosdb - 云 数据库 CynosDB 对于 MySQL, 默认值 是 mysql.",
			},

			"top_space_tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 的 Top tablespace 统计 返回.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "表 名称.",
						},
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库 名称.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage 引擎 对于 数据库 tables.",
						},
						"data_length": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "数据 space (MB).",
						},
						"index_length": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Index space (MB).",
						},
						"data_free": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Fragmentation space (MB).",
						},
						"total_length": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total space 使用 (MB).",
						},
						"frag_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Fragmentation 速率 (%).",
						},
						"table_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number 的 lines.",
						},
						"physical_file_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "independent physical 文件 大小 (MB) corresponding 到 表.",
						},
					},
				},
			},

			"timestamp": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "timestamp (在 秒) 的 collecting tablespace 数据.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainTopSpaceTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_top_space_tables.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var instanceId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var topSpaceTables []*dbbrain.TableSpaceData
	var timestamp *int64

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, ts, e := service.DescribeDbbrainTopSpaceTablesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		topSpaceTables = result
		timestamp = ts
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topSpaceTables))
	tmpList := make([]map[string]interface{}, 0, len(topSpaceTables))

	if topSpaceTables != nil {
		for _, tableSpaceData := range topSpaceTables {
			tableSpaceDataMap := map[string]interface{}{}

			if tableSpaceData.TableName != nil {
				tableSpaceDataMap["table_name"] = tableSpaceData.TableName
			}

			if tableSpaceData.TableSchema != nil {
				tableSpaceDataMap["table_schema"] = tableSpaceData.TableSchema
			}

			if tableSpaceData.Engine != nil {
				tableSpaceDataMap["engine"] = tableSpaceData.Engine
			}

			if tableSpaceData.DataLength != nil {
				tableSpaceDataMap["data_length"] = tableSpaceData.DataLength
			}

			if tableSpaceData.IndexLength != nil {
				tableSpaceDataMap["index_length"] = tableSpaceData.IndexLength
			}

			if tableSpaceData.DataFree != nil {
				tableSpaceDataMap["data_free"] = tableSpaceData.DataFree
			}

			if tableSpaceData.TotalLength != nil {
				tableSpaceDataMap["total_length"] = tableSpaceData.TotalLength
			}

			if tableSpaceData.FragRatio != nil {
				tableSpaceDataMap["frag_ratio"] = tableSpaceData.FragRatio
			}

			if tableSpaceData.TableRows != nil {
				tableSpaceDataMap["table_rows"] = tableSpaceData.TableRows
			}

			if tableSpaceData.PhysicalFileSize != nil {
				tableSpaceDataMap["physical_file_size"] = tableSpaceData.PhysicalFileSize
			}

			ids = append(ids, strings.Join([]string{instanceId, *tableSpaceData.TableSchema, *tableSpaceData.TableName}, tccommon.FILED_SP))
			tmpList = append(tmpList, tableSpaceDataMap)
		}

		_ = d.Set("top_space_tables", tmpList)
	}

	if timestamp != nil {
		_ = d.Set("timestamp", timestamp)
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
