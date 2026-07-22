package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataGetTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataGetTableRead,
		Schema: map[string]*schema.Schema{
			"table_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Table GUID。",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data 表 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data 表 GUID。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data 表 名称",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data 表 描述",
						},
						"database_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database 名称",
						},
						"schema_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database schema 名称",
						},
						"table_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table 类型",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"technical_metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Technical metadata 的 表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"owner": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Responsible person。",
									},
									"location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data 表 location。",
									},
									"storage_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Storage 大小。",
									},
								},
							},
						},
						"business_metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Business metadata 的 表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_names": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "标签 names。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
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

func dataSourceTencentCloudWedataGetTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_get_table.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(nil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		tableGuid string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("table_guid"); ok {
		paramMap["TableGuid"] = helper.String(v.(string))
		tableGuid = v.(string)
	}

	var respData *wedatav20250806.TableInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataGetTableByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})
	if reqErr != nil {
		return reqErr
	}

	tmpList := make([]map[string]interface{}, 0)
	if respData != nil {
		dataMap := map[string]interface{}{}
		if respData.Guid != nil {
			dataMap["guid"] = respData.Guid
		}

		if respData.Name != nil {
			dataMap["name"] = respData.Name
		}

		if respData.Description != nil {
			dataMap["description"] = respData.Description
		}

		if respData.DatabaseName != nil {
			dataMap["database_name"] = respData.DatabaseName
		}

		if respData.SchemaName != nil {
			dataMap["schema_name"] = respData.SchemaName
		}

		if respData.TableType != nil {
			dataMap["table_type"] = respData.TableType
		}

		if respData.CreateTime != nil {
			dataMap["create_time"] = respData.CreateTime
		}

		if respData.UpdateTime != nil {
			dataMap["update_time"] = respData.UpdateTime
		}

		technicalMetadataMap := map[string]interface{}{}
		if respData.TechnicalMetadata != nil {
			if respData.TechnicalMetadata.Owner != nil {
				technicalMetadataMap["owner"] = respData.TechnicalMetadata.Owner
			}

			if respData.TechnicalMetadata.Location != nil {
				technicalMetadataMap["location"] = respData.TechnicalMetadata.Location
			}

			if respData.TechnicalMetadata.StorageSize != nil {
				technicalMetadataMap["storage_size"] = respData.TechnicalMetadata.StorageSize
			}

			dataMap["technical_metadata"] = []interface{}{technicalMetadataMap}
		}

		businessMetadataMap := map[string]interface{}{}
		if respData.BusinessMetadata != nil {
			if respData.BusinessMetadata.TagNames != nil {
				businessMetadataMap["tag_names"] = respData.BusinessMetadata.TagNames
			}

			dataMap["business_metadata"] = []interface{}{businessMetadataMap}
		}

		tmpList = append(tmpList, dataMap)
	}

	_ = d.Set("data", tmpList)

	d.SetId(tableGuid)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
