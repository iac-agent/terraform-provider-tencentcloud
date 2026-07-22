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

func DataSourceTencentCloudDbbrainTopSpaceSchemas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainTopSpaceSchemasRead,
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
				Description: "数量 的 Top libraries 到 返回, 最大 值 是 100, 和 默认值 是 20.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sorting 字段 使用 到 过滤器 Top 库. 可选 字段 include DataLength, IndexLength, TotalLength, DataFree, FragRatio, TableRows, 和 PhysicalFileSize (仅 支持 通过 ApsaraDB 对于 MySQL 实例). 默认值 对于 ApsaraDB 对于 MySQL 实例 是 PhysicalFileSize, 和 默认值 对于 other product 实例 是 TotalLength.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: mysql - 云 数据库 MySQL, cynosdb - 云 数据库 CynosDB 对于 MySQL, 默认值 是 mysql.",
			},

			"top_space_schemas": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "返回 列表 的 top 库 space 统计.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "库 名称.",
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
							Description: "sum (MB) 的 independent physical 文件 sizes corresponding 到 all tables 在 库. 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
					},
				},
			},

			"timestamp": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Timestamp (在 秒) 当 库 space 数据 是 collected.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainTopSpaceSchemasRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_top_space_schemas.read")()
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

	var topSpaceSchemas []*dbbrain.SchemaSpaceData
	var timestamp *int64

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, ts, e := service.DescribeDbbrainTopSpaceSchemasByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		topSpaceSchemas = result
		timestamp = ts
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topSpaceSchemas))
	tmpList := make([]map[string]interface{}, 0, len(topSpaceSchemas))

	if topSpaceSchemas != nil {
		for _, schemaSpaceData := range topSpaceSchemas {
			schemaSpaceDataMap := map[string]interface{}{}

			if schemaSpaceData.TableSchema != nil {
				schemaSpaceDataMap["table_schema"] = schemaSpaceData.TableSchema
			}

			if schemaSpaceData.DataLength != nil {
				schemaSpaceDataMap["data_length"] = schemaSpaceData.DataLength
			}

			if schemaSpaceData.IndexLength != nil {
				schemaSpaceDataMap["index_length"] = schemaSpaceData.IndexLength
			}

			if schemaSpaceData.DataFree != nil {
				schemaSpaceDataMap["data_free"] = schemaSpaceData.DataFree
			}

			if schemaSpaceData.TotalLength != nil {
				schemaSpaceDataMap["total_length"] = schemaSpaceData.TotalLength
			}

			if schemaSpaceData.FragRatio != nil {
				schemaSpaceDataMap["frag_ratio"] = schemaSpaceData.FragRatio
			}

			if schemaSpaceData.TableRows != nil {
				schemaSpaceDataMap["table_rows"] = schemaSpaceData.TableRows
			}

			if schemaSpaceData.PhysicalFileSize != nil {
				schemaSpaceDataMap["physical_file_size"] = schemaSpaceData.PhysicalFileSize
			}

			ids = append(ids, strings.Join([]string{instanceId, *schemaSpaceData.TableSchema}, tccommon.FILED_SP))
			tmpList = append(tmpList, schemaSpaceDataMap)
		}

		_ = d.Set("top_space_schemas", tmpList)
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
