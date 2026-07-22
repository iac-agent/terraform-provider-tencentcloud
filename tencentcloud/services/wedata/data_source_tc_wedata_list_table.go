package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataListTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataListTableRead,
		Schema: map[string]*schema.Schema{
			"catalog_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Directory 名称",
			},

			"datasource_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "数据源 ID",
			},

			"database_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database 名称",
			},

			"schema_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database schema 名称",
			},

			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Table search keyword。",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Schema record list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data table GUID。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data table 名称",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data table 描述",
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
							Description: "Technical metadata of the table。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"owner": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "所有者",
									},
									"location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data table location。",
									},
									"storage_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Storage size。",
									},
								},
							},
						},
						"business_metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Business metadata of the table。",
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

func dataSourceTencentCloudWedataListTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_list_table.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("catalog_name"); ok {
		paramMap["CatalogName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("datasource_id"); ok {
		paramMap["DatasourceId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("database_name"); ok {
		paramMap["DatabaseName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("schema_name"); ok {
		paramMap["SchemaName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.TableInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataListTableByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	itemsList := make([]map[string]interface{}, 0, len(respData))
	for _, items := range respData {
		itemsMap := map[string]interface{}{}
		if items.Guid != nil {
			itemsMap["guid"] = items.Guid
		}

		if items.Name != nil {
			itemsMap["name"] = items.Name
		}

		if items.Description != nil {
			itemsMap["description"] = items.Description
		}

		if items.DatabaseName != nil {
			itemsMap["database_name"] = items.DatabaseName
		}

		if items.SchemaName != nil {
			itemsMap["schema_name"] = items.SchemaName
		}

		if items.TableType != nil {
			itemsMap["table_type"] = items.TableType
		}

		if items.CreateTime != nil {
			itemsMap["create_time"] = items.CreateTime
		}

		if items.UpdateTime != nil {
			itemsMap["update_time"] = items.UpdateTime
		}

		technicalMetadataMap := map[string]interface{}{}
		if items.TechnicalMetadata != nil {
			if items.TechnicalMetadata.Owner != nil {
				technicalMetadataMap["owner"] = items.TechnicalMetadata.Owner
			}

			if items.TechnicalMetadata.Location != nil {
				technicalMetadataMap["location"] = items.TechnicalMetadata.Location
			}

			if items.TechnicalMetadata.StorageSize != nil {
				technicalMetadataMap["storage_size"] = items.TechnicalMetadata.StorageSize
			}

			itemsMap["technical_metadata"] = []interface{}{technicalMetadataMap}
		}

		businessMetadataMap := map[string]interface{}{}
		if items.BusinessMetadata != nil {
			if items.BusinessMetadata.TagNames != nil {
				businessMetadataMap["tag_names"] = items.BusinessMetadata.TagNames
			}

			itemsMap["business_metadata"] = []interface{}{businessMetadataMap}
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("items", itemsList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
