package dnspod

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodRecordList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRecordListRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The 域名 to which the resolution record belongs。",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The 域名 ID to which the resolution record belongs. If DomainId is provided，the system will ignore the 域名 parameter. You can find all 域名 and DomainId through the DescribeDomainList interface。",
			},

			"sub_domain": {
				Optional:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"sub_domains"},
				Description:   "Retrieve resolution records based on the 主机 header of the resolution record. Fuzzy matching is used by default. You can set the IsExactSubdomain parameter to true for precise searching。",
			},
			"sub_domains": {
				Optional:      true,
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"sub_domain"},
				Description:   "Sub domains。",
			},

			"record_type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Retrieve certain types of resolution records，such as A，CNAME，NS，AAAA，explicit URL，implicit URL，CAA，SPF，etc。",
			},

			"record_line": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Retrieve resolution records for certain line IDs. You can view the allowed line information for the current 域名 through the DescribeRecordLineList interface。",
			},

			"group_id": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "When retrieving resolution records under certain groups，pass this 组 ID You can obtain the GroupId field through the DescribeRecordGroupList interface。",
			},

			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search for resolution records by keyword，currently supporting searching 主机 headers and record values。",
			},

			"sort_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field，supporting NAME，LINE，TYPE，VALUE，WEIGHT，MX，TTL，UPDATED_ON fields. NAME: The 主机 header of the resolution record LINE: The resolution record line TYPE: The resolution record 类型 VALUE: The resolution record 值 WEIGHT: The 权重 MX: MX 优先级 TTL: The resolution record cache time UPDATED_ON: The resolution record 更新时间。",
			},

			"sort_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting method，ascending: ASC，descending: DESC. The 默认值为 ASC。",
			},

			"record_value": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get the resolution record based on the resolution record 值",
			},

			"record_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Get the resolution record based on the resolution record 状态 The possible values are ENABLE and DISABLE. ENABLE: Normal DISABLE: Paused。",
			},

			"weight_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The starting point of the resolution record 权重 query interval。",
			},

			"weight_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The endpoint of the resolution record 权重 query interval。",
			},

			"mx_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The starting point of the resolution record MX 优先级 query interval。",
			},

			"mx_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The endpoint of the resolution record MX 优先级 query interval。",
			},

			"ttl_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The starting point of the resolution record TTL query interval。",
			},

			"ttl_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The endpoint of the resolution record TTL query interval。",
			},

			"updated_at_begin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The starting point of the resolution record 更新时间 query interval。",
			},

			"updated_at_end": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The endpoint of the resolution record 更新时间 query interval。",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get the resolution record based on the resolution record 备注",
			},

			"is_exact_sub_domain": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否perform an exact search based on the SubDomain parameter。",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "项目 ID",
			},

			"filter_at_ns": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Filter @ 类型 NS records. 默认为 false。",
				Default:     false,
			},

			"record_count_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Statistics of the 数量 records。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subdomain_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 subdomains。",
						},
						"list_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 records returned in the list。",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 records。",
						},
					},
				},
			},

			"record_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 records。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record ID。",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 值",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 状态，已启用: ENABLE，paused: DISABLE。",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 header。",
						},
						"line": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record line。",
						},
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line ID。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 类型",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record 权重，用于load balancing records. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"monitor_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record monitoring 状态，normal: OK，alarm: WARN，downtime: DOWN，empty if monitoring is not set or paused。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 备注 描述",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record cache time。",
						},
						"mx": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX 值，only available for MX records 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"default_ns": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为the default NS record。",
						},
					},
				},
			},
			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 records。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "ID。",
						},
						"domain": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "域名",
						},
						"record_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record ID。",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 值",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 状态，已启用: ENABLE，paused: DISABLE。",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 header。",
						},
						"record_line": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record line。",
						},
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line ID。",
						},
						"record_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 类型",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record 权重，用于load balancing records. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"monitor_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record monitoring 状态，normal: OK，alarm: WARN，downtime: DOWN，empty if monitoring is not set or paused。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record 备注 描述",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record cache time。",
						},
						"mx": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX 值，only available for MX records 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"default_ns": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为the default NS record。",
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

func dataSourceTencentCloudDnspodRecordListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnspod_record_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var domain string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		paramMap["Domain"] = helper.String(domain)
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		paramMap["DomainId"] = helper.IntUint64(v.(int))
	}

	subDomains := make([]string, 0)
	if v, ok := d.GetOk("sub_domain"); ok {
		subDomains = append(subDomains, v.(string))
	}

	if v, ok := d.GetOk("sub_domains"); ok {
		subDomainList := v.(*schema.Set).List()
		for _, subDomain := range subDomainList {
			subDomains = append(subDomains, subDomain.(string))
		}
	}

	if v, ok := d.GetOk("record_type"); ok {
		recordTypeSet := v.(*schema.Set).List()
		paramMap["RecordType"] = helper.InterfacesStringsPoint(recordTypeSet)
	}

	if v, ok := d.GetOk("record_line"); ok {
		recordLineSet := v.(*schema.Set).List()
		paramMap["RecordLine"] = helper.InterfacesStringsPoint(recordLineSet)
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupIds := make([]*uint64, 0)
		for _, item := range v.(*schema.Set).List() {
			groupIds = append(groupIds, helper.IntUint64(item.(int)))
		}
		paramMap["GroupId"] = groupIds
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_field"); ok {
		paramMap["SortField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_type"); ok {
		paramMap["SortType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_value"); ok {
		paramMap["RecordValue"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_status"); ok {
		recordStatusSet := v.(*schema.Set).List()
		paramMap["RecordStatus"] = helper.InterfacesStringsPoint(recordStatusSet)
	}

	if v, ok := d.GetOkExists("weight_begin"); ok {
		paramMap["WeightBegin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("weight_end"); ok {
		paramMap["WeightEnd"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("mx_begin"); ok {
		paramMap["MXBegin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("mx_end"); ok {
		paramMap["MXEnd"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("ttl_begin"); ok {
		paramMap["TTLBegin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("ttl_end"); ok {
		paramMap["TTLEnd"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("updated_at_begin"); ok {
		paramMap["UpdatedAtBegin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("updated_at_end"); ok {
		paramMap["UpdatedAtEnd"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		paramMap["Remark"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_exact_sub_domain"); ok {
		paramMap["IsExactSubDomain"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}
	var filterAtNS bool
	if v, ok := d.GetOkExists("filter_at_ns"); ok {
		filterAtNS = v.(bool)
	}
	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var recordList []*dnspod.RecordListItem

	if len(subDomains) == 0 {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := service.DescribeDnspodRecordListByFilter(ctx, paramMap)
			if e != nil {
				return tccommon.RetryError(e)
			}
			recordList = append(recordList, result...)
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		recordIds := map[uint64]struct{}{}

		for _, subDomain := range subDomains {
			paramMap["SubDomain"] = helper.String(subDomain)
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				result, e := service.DescribeDnspodRecordListByFilter(ctx, paramMap)
				if e != nil {
					return tccommon.RetryError(e)
				}
				for _, resultItem := range result {
					if resultItem.RecordId == nil {
						return resource.NonRetryableError(fmt.Errorf("record id is nil"))
					}
					if _, ok := recordIds[*resultItem.RecordId]; ok {
						continue
					} else {
						recordIds[*resultItem.RecordId] = struct{}{}
						recordList = append(recordList, resultItem)
					}
				}

				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	ids := make([]string, 0, len(recordList))
	tmpList := make([]map[string]interface{}, 0, len(recordList))
	instanceList := make([]map[string]interface{}, 0, len(recordList))
	if recordList != nil {
		for _, recordListItem := range recordList {
			if filterAtNS && recordListItem.Name != nil && *recordListItem.Name == DNSPOD_RECORD_NAME_AT && recordListItem.Type != nil && *recordListItem.Type == DNSPOD_RECORD_TYPE_NS {
				continue
			}
			recordListItemMap := map[string]interface{}{}
			instanceListItemMap := map[string]interface{}{}
			instanceListItemMap["domain"] = domain
			if recordListItem.RecordId != nil {
				recordListItemMap["record_id"] = recordListItem.RecordId
				instanceListItemMap["record_id"] = recordListItem.RecordId
				instanceListItemMap["id"] = domain + tccommon.FILED_SP + helper.UInt64ToStr(*recordListItem.RecordId)
			}

			if recordListItem.Value != nil {
				recordListItemMap["value"] = recordListItem.Value
				instanceListItemMap["value"] = recordListItem.Value
			}

			if recordListItem.Status != nil {
				recordListItemMap["status"] = recordListItem.Status
				instanceListItemMap["status"] = recordListItem.Status
			}

			if recordListItem.UpdatedOn != nil {
				recordListItemMap["updated_on"] = recordListItem.UpdatedOn
				instanceListItemMap["updated_on"] = recordListItem.UpdatedOn
			}

			if recordListItem.Name != nil {
				recordListItemMap["name"] = recordListItem.Name
				instanceListItemMap["sub_domain"] = recordListItem.Name
			}

			if recordListItem.Line != nil {
				recordListItemMap["line"] = recordListItem.Line
				instanceListItemMap["record_line"] = recordListItem.Line
			}

			if recordListItem.LineId != nil {
				recordListItemMap["line_id"] = recordListItem.LineId
				instanceListItemMap["line_id"] = recordListItem.LineId
			}

			if recordListItem.Type != nil {
				recordListItemMap["type"] = recordListItem.Type
				instanceListItemMap["record_type"] = recordListItem.Type
			}

			if recordListItem.Weight != nil {
				recordListItemMap["weight"] = recordListItem.Weight
				instanceListItemMap["weight"] = recordListItem.Weight
			}

			if recordListItem.MonitorStatus != nil {
				recordListItemMap["monitor_status"] = recordListItem.MonitorStatus
				instanceListItemMap["monitor_status"] = recordListItem.MonitorStatus
			}

			if recordListItem.Remark != nil {
				recordListItemMap["remark"] = recordListItem.Remark
				instanceListItemMap["remark"] = recordListItem.Remark
			}

			if recordListItem.TTL != nil {
				recordListItemMap["ttl"] = recordListItem.TTL
				instanceListItemMap["ttl"] = recordListItem.TTL
			}

			if recordListItem.MX != nil {
				recordListItemMap["mx"] = recordListItem.MX
				instanceListItemMap["mx"] = recordListItem.MX
			}

			if recordListItem.DefaultNS != nil {
				recordListItemMap["default_ns"] = recordListItem.DefaultNS
				instanceListItemMap["default_ns"] = recordListItem.DefaultNS
			}

			ids = append(ids, helper.UInt64ToStr(*recordListItem.RecordId))
			tmpList = append(tmpList, recordListItemMap)
			instanceList = append(instanceList, instanceListItemMap)
		}

		_ = d.Set("record_list", tmpList)
		_ = d.Set("instance_list", instanceList)
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
