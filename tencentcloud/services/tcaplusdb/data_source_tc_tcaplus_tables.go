package tcaplusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudTcaplusTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusTablesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID TcaplusDB 集群 到 是 查询。",
			},
			"tablegroup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 表 组 到 是 查询。",
			},
			"table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Table ID 到 是 查询。",
			},
			"table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Table 名称 到 是 查询。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File 对于 saving results。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 TcaplusDB tables. Each element 包含following attributes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tablegroup_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table 组 ID TcaplusDB 表。",
						},
						"table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID TcaplusDB 表。",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 TcaplusDB 表。",
						},
						"table_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 TcaplusDB 表。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 TcaplusDB 表。",
						},
						"idl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IDL 文件 ID TcaplusDB 表。",
						},
						"table_idl_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IDL 类型 TcaplusDB 表。",
						},
						"reserved_read_cu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Reserved read 容量 units 的 TcaplusDB 表。",
						},
						"reserved_write_cu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Reserved write 容量 units 的 TcaplusDB 表。",
						},
						"reserved_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Reserved 存储 容量 的 TcaplusDB 表 (单位:GB)。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 TcaplusDB 表。",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "错误信息 对于 creating TcaplusDB 表。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 TcaplusDB 表。",
						},
						"table_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size 的 TcaplusDB 表。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcaplus_tables.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcaplusService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Get("tablegroup_id").(string)
	tableId := d.Get("table_id").(string)
	tableName := d.Get("table_name").(string)

	tables, err := service.DescribeTables(ctx, clusterId, groupId, tableId, tableName)
	if err != nil {
		tables, err = service.DescribeTables(ctx, clusterId, groupId, tableId, tableName)
	}
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(tables))

	for _, tableInfo := range tables {

		listItem := make(map[string]interface{})

		if tableInfo.IdlFiles != nil && len(tableInfo.IdlFiles) > 0 {
			idlFile := tableInfo.IdlFiles[0]

			var tcaplusIdlId TcaplusIdlId
			tcaplusIdlId.ClusterId = clusterId
			tcaplusIdlId.FileName = *idlFile.FileName
			tcaplusIdlId.FileType = *idlFile.FileType

			tcaplusIdlId.FileExtType = *idlFile.FileExtType
			tcaplusIdlId.FileSize = *idlFile.FileSize
			tcaplusIdlId.FileId = *idlFile.FileId
			id, err := json.Marshal(tcaplusIdlId)

			if err != nil {
				return fmt.Errorf("format idl id fail,%s", err.Error())
			}
			listItem["idl_id"] = string(id)
		}

		if tableInfo.Error != nil && tableInfo.Error.Message != nil {
			listItem["error"] = *tableInfo.Error.Message
		} else {
			listItem["error"] = ""
		}
		if tableInfo.TableGroupId != nil {
			listItem["tablegroup_id"] = fmt.Sprintf("%s:%s", clusterId, *tableInfo.TableGroupId)
		}
		if tableInfo.TableInstanceId != nil {
			listItem["table_id"] = *tableInfo.TableInstanceId
		}
		if tableInfo.TableName != nil {
			listItem["table_name"] = *tableInfo.TableName
		}
		if tableInfo.TableType != nil {
			listItem["table_type"] = *tableInfo.TableType
		}
		if tableInfo.Memo != nil {
			listItem["description"] = *tableInfo.Memo
		}
		if tableInfo.TableIdlType != nil {
			listItem["table_idl_type"] = *tableInfo.TableIdlType
		}
		if tableInfo.ReservedReadQps != nil {
			listItem["reserved_read_cu"] = *tableInfo.ReservedReadQps
		}
		if tableInfo.ReservedWriteQps != nil {
			listItem["reserved_write_cu"] = *tableInfo.ReservedWriteQps
		}
		if tableInfo.ReservedVolume != nil {
			listItem["reserved_volume"] = *tableInfo.ReservedVolume
		}
		if tableInfo.CreatedTime != nil {
			listItem["create_time"] = *tableInfo.CreatedTime
		}
		if tableInfo.Status != nil {
			listItem["status"] = *tableInfo.Status
		}
		if tableInfo.TableSize != nil {
			listItem["table_size"] = *tableInfo.TableSize
		}
		list = append(list, listItem)
	}

	d.SetId("table." + clusterId + "." + groupId + "." + tableId + "." + tableName)
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
