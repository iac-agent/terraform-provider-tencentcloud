package dnspod

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodDomainList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodDomainListRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Get 域名 names based 在 域名 组 类型 Available 值 是 ALL，MINE，SHARE，RECENT. ALL: All MINE: My 域名 names SHARE: 域名 names shared 使用 me RECENT: Recently operated 域名 names。",
			},

			"group_id": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Get 域名 names based 在 域名 组 ID，其中 可以 是 获取 through GroupId 字段 在 DescribeDomain 或 DescribeDomainList interface。",
			},

			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get 域名 names based 在 keywords。",
			},

			"sort_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 字段. Available 值 是 NAME，STATUS，RECORDS，GRADE，UPDATED_ON. NAME: 域名 名称 STATUS: 域名 状态 RECORDS: 数量 records GRADE: Package 级别 UPDATED_ON: 更新时间。",
			},

			"sort_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 类型，ascending: ASC，descending: DESC。",
			},

			"status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Get 域名 names based 在 域名 状态 Available 值 是 ENABLE，LOCK，PAUSE，SPAM. ENABLE: Normal LOCK: Locked PAUSE: Paused SPAM: Banned。",
			},

			"package": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Get 域名 names based 在 包，其中 可以 是 获取 through Grade 字段 在 DescribeDomain 或 DescribeDomainList interface。",
			},

			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Get 域名 names based 在 备注 信息。",
			},

			"updated_at_begin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "开始时间 的 域名 名称&amp;#39;s 更新时间 到 是 获取，such 作为 &amp;#39;2021-05-01 03:00:00&amp;#39;。",
			},

			"updated_at_end": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "结束时间 的 域名 名称&amp;#39;s 更新时间 到 是 获取，such 作为 &amp;#39;2021-05-10 20:00:00&amp;#39;。",
			},

			"record_count_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "start point 的 域名 名称&amp;#39;s 记录 count 查询 范围。",
			},

			"record_count_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "end point 的 域名 名称&amp;#39;s 记录 count 查询 范围。",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "项目 ID",
			},

			"tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "标签描述列表",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "待过滤字段",
						},
						"tag_value": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "过滤值 的 字段。",
						},
					},
				},
			},

			"domain_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "域名 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier assigned 到 域名 通过 系统。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original 格式 的 域名",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 状态，normal: ENABLE，paused: PAUSE，banned: SPAM。",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Default TTL 值 对于 域名 resolution records。",
						},
						"cname_speedup": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否enable CNAME acceleration，已启用: ENABLE，已禁用: DISABLE。",
						},
						"dns_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS settings 状态，错误: DNSERROR，normal: 空 字符串。",
						},
						"grade": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 包 级别 代码",
						},
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Group ID 域名 belongs 到。",
						},
						"search_engine_push": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否enable search 引擎 push optimization，YES: YES，NO: NO。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 备注 描述",
						},
						"punycode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Punycode encoded 域名 格式",
						},
						"effective_dns": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "有效 DNS assigned 到 域名 通过 系统。",
						},
						"grade_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sequence 数量 corresponding 到 域名 包 级别",
						},
						"grade_title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package 名称",
						},
						"is_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否为a paid 包。",
						},
						"vip_start_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Paid 包 activation 时间。",
						},
						"vip_end_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Paid 包 过期时间。",
						},
						"vip_auto_renew": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否domain has VIP auto-renewal 已启用，YES: YES，NO: NO，DEFAULT: DEFAULT。",
						},
						"record_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 records under 域名",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 addition 时间。",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 更新时间。",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 所有者 账号",
						},
						"tag_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "域名-related 标签列表 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值 注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudDnspodDomainListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnspod_domain_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupIds := make([]*int64, 0)
		for _, item := range v.(*schema.Set).List() {
			groupIds = append(groupIds, helper.IntInt64(item.(int)))
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

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		paramMap["Status"] = helper.InterfacesStringsPoint(statusSet)
	}

	if v, ok := d.GetOk("package"); ok {
		packageSet := v.(*schema.Set).List()
		paramMap["Package"] = helper.InterfacesStringsPoint(packageSet)
	}

	if v, ok := d.GetOk("remark"); ok {
		paramMap["Remark"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("updated_at_begin"); ok {
		paramMap["UpdatedAtBegin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("updated_at_end"); ok {
		paramMap["UpdatedAtEnd"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("record_count_begin"); ok {
		paramMap["RecordCountBegin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("record_count_end"); ok {
		paramMap["RecordCountEnd"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	// tags := helper.GetTagsFilter(d, "tags")
	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.([]interface{})
		tmpSet := make([]*dnspod.TagItemFilter, 0, len(tagsSet))

		for _, item := range tagsSet {
			filter := dnspod.TagItemFilter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["tag_key"]; ok {
				filter.TagKey = helper.String(v.(string))
			}

			if v, ok := filterMap["tag_value"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.TagValue = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Tags"] = tmpSet
	}

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var domainList []*dnspod.DomainListItem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDnspodDomainListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		domainList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(domainList))
	tmpList := make([]map[string]interface{}, 0, len(domainList))
	if domainList != nil {
		for _, domainListItem := range domainList {
			domainListItemMap := map[string]interface{}{}

			if domainListItem.DomainId != nil {
				domainListItemMap["domain_id"] = domainListItem.DomainId
			}

			if domainListItem.Name != nil {
				domainListItemMap["name"] = domainListItem.Name
			}

			if domainListItem.Status != nil {
				domainListItemMap["status"] = domainListItem.Status
			}

			if domainListItem.TTL != nil {
				domainListItemMap["ttl"] = domainListItem.TTL
			}

			if domainListItem.CNAMESpeedup != nil {
				domainListItemMap["cname_speedup"] = domainListItem.CNAMESpeedup
			}

			if domainListItem.DNSStatus != nil {
				domainListItemMap["dns_status"] = domainListItem.DNSStatus
			}

			if domainListItem.Grade != nil {
				domainListItemMap["grade"] = domainListItem.Grade
			}

			if domainListItem.GroupId != nil {
				domainListItemMap["group_id"] = domainListItem.GroupId
			}

			if domainListItem.SearchEnginePush != nil {
				domainListItemMap["search_engine_push"] = domainListItem.SearchEnginePush
			}

			if domainListItem.Remark != nil {
				domainListItemMap["remark"] = domainListItem.Remark
			}

			if domainListItem.Punycode != nil {
				domainListItemMap["punycode"] = domainListItem.Punycode
			}

			if domainListItem.EffectiveDNS != nil {
				domainListItemMap["effective_dns"] = domainListItem.EffectiveDNS
			}

			if domainListItem.GradeLevel != nil {
				domainListItemMap["grade_level"] = domainListItem.GradeLevel
			}

			if domainListItem.GradeTitle != nil {
				domainListItemMap["grade_title"] = domainListItem.GradeTitle
			}

			if domainListItem.IsVip != nil {
				domainListItemMap["is_vip"] = domainListItem.IsVip
			}

			if domainListItem.VipStartAt != nil {
				domainListItemMap["vip_start_at"] = domainListItem.VipStartAt
			}

			if domainListItem.VipEndAt != nil {
				domainListItemMap["vip_end_at"] = domainListItem.VipEndAt
			}

			if domainListItem.VipAutoRenew != nil {
				domainListItemMap["vip_auto_renew"] = domainListItem.VipAutoRenew
			}

			if domainListItem.RecordCount != nil {
				domainListItemMap["record_count"] = domainListItem.RecordCount
			}

			if domainListItem.CreatedOn != nil {
				domainListItemMap["created_on"] = domainListItem.CreatedOn
			}

			if domainListItem.UpdatedOn != nil {
				domainListItemMap["updated_on"] = domainListItem.UpdatedOn
			}

			if domainListItem.Owner != nil {
				domainListItemMap["owner"] = domainListItem.Owner
			}

			if domainListItem.TagList != nil {
				tagListList := []interface{}{}
				for _, tagList := range domainListItem.TagList {
					tagListMap := map[string]interface{}{}

					if tagList.TagKey != nil {
						tagListMap["tag_key"] = tagList.TagKey
					}

					if tagList.TagValue != nil {
						tagListMap["tag_value"] = tagList.TagValue
					}

					tagListList = append(tagListList, tagListMap)
				}

				domainListItemMap["tag_list"] = tagListList
			}

			ids = append(ids, strconv.FormatUint(*domainListItem.DomainId, 10))
			tmpList = append(tmpList, domainListItemMap)
		}

		_ = d.Set("domain_list", tmpList)
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
