package privatedns

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPrivateDnsRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPrivateDnsRecordsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Private 可用区 ID: 可用区-xxxxxx。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 参数 (值 和 RecordType filtering 是 支持)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter 名称",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Parameter 值。",
						},
					},
				},
			},

			"record_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Parse 记录 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record sid。",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private 可用区 ID: 可用区-xxxxxx。",
						},
						"sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subdomain 名称",
						},
						"record_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 类型，可选 记录 类型 是: A，AAAA，CNAME，MX，TXT，PTR。",
						},
						"record_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 值",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record 缓存 时间， smaller 值， faster 它 takes effect. 值 是 1-86400s. 默认为 600。",
						},
						"mx": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX 优先级: 必填 如果 记录 类型 是 MX. 取值范围：5,10,15,20,30,40,50。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 状态",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record 权重，值 是 1-100。",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 创建时间。",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 更新时间。",
						},
						"extra": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Additional 信息。",
						},
						"enabled": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "已启用 0 meaning paused，1 meaning senabled。",
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

func dataSourceTencentCloudPrivateDnsRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_private_dns_records.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	zoneId := d.Get("zone_id").(string)
	filterList := make([]*privatedns.Filter, 0)
	if v, ok := d.GetOk("filters"); ok {
		filters := v.([]interface{})
		for _, item := range filters {
			filter := privatedns.Filter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			filterList = append(filterList, &filter)
		}
	}

	var recordSet []*privatedns.PrivateZoneRecord
	recordSet, err := service.DescribePrivateDnsRecordByFilter(ctx, zoneId, filterList)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(recordSet))
	tmpList := make([]map[string]interface{}, 0, len(recordSet))
	if recordSet != nil {
		for _, privateZoneRecord := range recordSet {
			privateZoneRecordMap := map[string]interface{}{}
			if privateZoneRecord.RecordId != nil {
				privateZoneRecordMap["record_id"] = privateZoneRecord.RecordId
			}

			if privateZoneRecord.ZoneId != nil {
				privateZoneRecordMap["zone_id"] = privateZoneRecord.ZoneId
			}

			if privateZoneRecord.SubDomain != nil {
				privateZoneRecordMap["sub_domain"] = privateZoneRecord.SubDomain
			}

			if privateZoneRecord.RecordType != nil {
				privateZoneRecordMap["record_type"] = privateZoneRecord.RecordType
			}

			if privateZoneRecord.RecordValue != nil {
				privateZoneRecordMap["record_value"] = privateZoneRecord.RecordValue
			}

			if privateZoneRecord.TTL != nil {
				privateZoneRecordMap["ttl"] = privateZoneRecord.TTL
			}

			if privateZoneRecord.MX != nil {
				privateZoneRecordMap["mx"] = privateZoneRecord.MX
			}

			if privateZoneRecord.Status != nil {
				privateZoneRecordMap["status"] = privateZoneRecord.Status
			}

			if privateZoneRecord.Weight != nil {
				privateZoneRecordMap["weight"] = privateZoneRecord.Weight
			}

			if privateZoneRecord.CreatedOn != nil {
				privateZoneRecordMap["created_on"] = privateZoneRecord.CreatedOn
			}

			if privateZoneRecord.UpdatedOn != nil {
				privateZoneRecordMap["updated_on"] = privateZoneRecord.UpdatedOn
			}

			if privateZoneRecord.Extra != nil {
				privateZoneRecordMap["extra"] = privateZoneRecord.Extra
			}

			if privateZoneRecord.Enabled != nil {
				privateZoneRecordMap["enabled"] = privateZoneRecord.Enabled
			}

			ids = append(ids, *privateZoneRecord.RecordId)
			tmpList = append(tmpList, privateZoneRecordMap)
		}

		_ = d.Set("record_set", tmpList)
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
