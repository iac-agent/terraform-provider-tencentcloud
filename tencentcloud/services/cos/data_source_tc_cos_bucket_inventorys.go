package cos

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCosBucketInventorys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBucketInventorysRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "存储桶",
			},
			"inventorys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multiple batch processing task information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否enable the inventory. true or false。",
						},
						"is_enabled": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否enable the inventory. true or false。",
						},
						"included_object_versions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否include object versions in the inventory. All or No。",
						},
						"filter": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Filters objects prefixed with the specified 值 to analyze。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix of the objects to analyze。",
									},
									"period": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "创建时间 range of the objects to analyze。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Creation 开始时间 of the objects to analyze. The parameter is a 时间戳 （秒）， for example，1568688761。",
												},
												"end_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Creation 结束时间 of the objects to analyze. The parameter is a 时间戳 （秒）， for example，1568688762。",
												},
											},
										},
									},
								},
							},
						},
						"optional_fields": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Analysis items to include in the inventory 结果	。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fields": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "可选 analysis items to include in the inventory 结果 The 可选 fields include Size，LastModifiedDate，StorageClass，ETag，IsMultipartUploaded，ReplicationStatus，标签，Crc64，and x-cos-meta-*。",
									},
								},
							},
						},
						"schedule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inventory job cycle。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Frequency of the inventory job. Enumerated values: Daily，Weekly。",
									},
								},
							},
						},
						"destination": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information about the inventory 结果 destination。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "存储桶名称",
									},
									"account_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID 存储桶 所有者",
									},
									"prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix of the inventory 结果",
									},
									"format": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "格式 of the inventory 结果 Valid 值: CSV。",
									},
									"encryption": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Server-side encryption for the inventory 结果",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sse_cos": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Encryption with COS-managed 键 This field can be left empty。",
												},
											},
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

func dataSourceTencentCloudCosBucketInventorysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cos_bucket_inventorys.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	bucket := d.Get("bucket").(string)
	inventoryConfigurations := make([]map[string]interface{}, 0)
	token := ""
	ids := make([]string, 0)
	for {
		result, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.ListInventoryConfigurations(ctx, token)
		responseBody, _ := json.Marshal(response.Body)
		log.Printf("[DEBUG]%s api[ListInventoryConfigurations] success, response body [%s]\n", logId, responseBody)
		if err != nil {
			return err
		}

		for _, item := range result.InventoryConfigurations {
			itemMap := make(map[string]interface{})
			itemMap["id"] = item.ID
			itemMap["is_enabled"] = item.IsEnabled
			itemMap["included_object_versions"] = item.IncludedObjectVersions

			filterMap := make(map[string]interface{})
			if item.Filter != nil {
				filterMap["prefix"] = item.Filter.Prefix
				periodMap := make(map[string]interface{})
				if item.Filter.Period != nil {
					if item.Filter.Period.StartTime != 0 {
						periodMap["start_time"] = strconv.FormatInt(item.Filter.Period.StartTime, 10)
					}
					if item.Filter.Period.EndTime != 0 {
						periodMap["end_time"] = strconv.FormatInt(item.Filter.Period.EndTime, 10)
					}
					filterMap["period"] = []interface{}{periodMap}
				}
				itemMap["filter"] = []interface{}{filterMap}
			}
			if item.OptionalFields != nil {
				optionalFieldsMap := make(map[string]interface{})
				fields := make([]string, 0)
				if item.OptionalFields.BucketInventoryFields != nil {
					fields = append(fields, item.OptionalFields.BucketInventoryFields...)
					optionalFieldsMap["fields"] = fields
				}
				itemMap["optional_fields"] = []interface{}{optionalFieldsMap}
			}

			if item.Schedule != nil {
				scheduleMap := make(map[string]interface{})
				scheduleMap["frequency"] = item.Schedule.Frequency
				itemMap["schedule"] = []interface{}{scheduleMap}
			}

			if item.Destination != nil {
				destinationMap := make(map[string]interface{})
				destinationMap["bucket"] = item.Destination.Bucket
				destinationMap["account_id"] = item.Destination.AccountId
				destinationMap["prefix"] = item.Destination.Prefix
				destinationMap["format"] = item.Destination.Format
				if item.Destination.Encryption != nil && item.Destination.Encryption.SSECOS != "" {
					encryptionMap := make(map[string]interface{})

					encryptionMap["sse_cos"] = item.Destination.Encryption.SSECOS
					destinationMap["encryption"] = []interface{}{encryptionMap}

				}
				itemMap["destination"] = []interface{}{destinationMap}
			}
			ids = append(ids, item.ID)
			inventoryConfigurations = append(inventoryConfigurations, itemMap)
		}
		if result.NextContinuationToken != "" {
			token = result.NextContinuationToken
		} else {
			break
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("inventorys", inventoryConfigurations)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), inventoryConfigurations); err != nil {
			return err
		}
	}

	return nil
}
