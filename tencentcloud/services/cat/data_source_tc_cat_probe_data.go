package cat

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCatProbeData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatProbedataRead,
		Schema: map[string]*schema.Schema{
			"begin_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Start 时间戳 (在 milliseconds)。",
			},

			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "End 时间戳 (在 milliseconds)。",
			},

			"task_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "任务 类型 在 AnalyzeTaskType_Network,AnalyzeTaskType_Browse,AnalyzeTaskType_UploadDownload,AnalyzeTaskType_Transport,AnalyzeTaskType_MediaStream。",
			},

			"sort_field": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Fields 到 是 sorted ProbeTime dial 测试 时间 sorting 可以 是 filled 在 You 可以 also fill 在 selected 字段 在 SelectedFields。",
			},

			"ascending": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "true 是 Ascending。",
			},

			"selected_fields": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Selected Fields。",
			},

			"offset": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "偏移量",
			},

			"limit": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "限制",
			},

			"task_id": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "TaskID 列表。",
			},

			"operators": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Operators 列表。",
			},

			"districts": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Districts 列表。",
			},

			"error_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "ErrorTypes 列表。",
			},

			"city": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "City 列表。",
			},

			"code": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "代码 列表。",
			},

			"detailed_single_data_define": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Probe 节点 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"probe_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Probe 时间。",
						},
						"labels": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Labels。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom Field 名称/描述",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "值",
									},
								},
							},
						},
						"fields": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Fields。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom Field 名称/描述",
									},
									"value": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "值",
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

func dataSourceTencentCloudCatProbedataRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cat_probedata.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("begin_time"); v != nil {
		paramMap["begin_time"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["end_time"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("task_type"); v != nil {
		paramMap["task_type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_field"); ok {
		paramMap["sort_field"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ascending"); v != nil {
		paramMap["ascending"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("selected_fields"); ok {
		selectedFieldsSet := v.(*schema.Set).List()
		selectedFields := make([]*string, 0, len(selectedFieldsSet))
		for _, vv := range selectedFieldsSet {
			selectedFields = append(selectedFields, helper.String(vv.(string)))
		}
		paramMap["selected_fields"] = selectedFields
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("task_id"); ok {
		taskIdSet := v.(*schema.Set).List()
		taskId := make([]*string, 0, len(taskIdSet))
		for _, vv := range taskIdSet {
			taskId = append(taskId, helper.String(vv.(string)))
		}
		paramMap["task_id"] = taskId
	}

	if v, ok := d.GetOk("operators"); ok {
		operatorsSet := v.(*schema.Set).List()
		operators := make([]*string, 0, len(operatorsSet))
		for _, vv := range operatorsSet {
			operators = append(operators, helper.String(vv.(string)))
		}
		paramMap["operators"] = operators
	}

	if v, ok := d.GetOk("districts"); ok {
		districtsSet := v.(*schema.Set).List()
		districts := make([]*string, 0, len(districtsSet))
		for _, vv := range districtsSet {
			districts = append(districts, helper.String(vv.(string)))
		}
		paramMap["districts"] = districts
	}

	if v, ok := d.GetOk("error_types"); ok {
		errorTypesSet := v.(*schema.Set).List()
		errorTypes := make([]*string, 0, len(errorTypesSet))
		for _, vv := range errorTypesSet {
			errorTypes = append(errorTypes, helper.String(vv.(string)))
		}
		paramMap["error_types"] = errorTypes
	}

	if v, ok := d.GetOk("city"); ok {
		citySet := v.(*schema.Set).List()
		city := make([]*string, 0, len(citySet))
		for _, vv := range citySet {
			city = append(city, helper.String(vv.(string)))
		}
		paramMap["city"] = city
	}

	catService := CatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var dataSets []*cat.DetailedSingleDataDefine
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := catService.DescribeCatProbeDataByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dataSets = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Cat dataSet failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(dataSets))
	dataSetList := make([]map[string]interface{}, 0, len(dataSets))

	if dataSets != nil {
		for _, dataSet := range dataSets {
			dataSetMap := map[string]interface{}{}
			if dataSet.ProbeTime != nil {
				dataSetMap["probe_time"] = dataSet.ProbeTime
			}
			if dataSet.Labels != nil {
				labelsList := []interface{}{}
				for _, labels := range dataSet.Labels {
					labelsMap := map[string]interface{}{}
					if labels.ID != nil {
						labelsMap["id"] = labels.ID
					}
					if labels.Name != nil {
						labelsMap["name"] = labels.Name
					}
					if labels.Value != nil {
						labelsMap["value"] = labels.Value
					}

					labelsList = append(labelsList, labelsMap)
				}
				dataSetMap["labels"] = labelsList
			}
			if dataSet.Fields != nil {
				fieldsList := []interface{}{}
				for _, fields := range dataSet.Fields {
					fieldsMap := map[string]interface{}{}
					if fields.ID != nil {
						fieldsMap["id"] = fields.ID
					}
					if fields.Name != nil {
						fieldsMap["name"] = fields.Name
					}
					if fields.Value != nil {
						fieldsMap["value"] = fields.Value
					}

					fieldsList = append(fieldsList, fieldsMap)
				}
				dataSetMap["fields"] = fieldsList
			}
			ids = append(ids, helper.UInt64ToStr(*dataSet.ProbeTime))
			dataSetList = append(dataSetList, dataSetMap)
		}

		_ = d.Set("detailed_single_data_define", dataSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataSetList); e != nil {
			return e
		}
	}

	return nil
}
