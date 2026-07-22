package dnspod

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodDomainAnalytics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodDomainAnalyticsRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "域名 名称 到 查询 对于 resolution 卷。",
			},

			"start_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "start date 的 查询，格式: YYYY-MM-DD。",
			},

			"end_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "end date 的 查询，格式: YYYY-MM-DD。",
			},

			"dns_format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DATE: Statistics 通过 day dimension HOUR: Statistics 通过 hour dimension。",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "域名 ID. 参数 DomainId has higher 优先级 比 参数 域名 如果 参数 DomainId 是 passed， 参数 域名 将 是 ignored. You 可以 find all Domains 和 DomainIds through DescribeDomainList interface。",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Subtotal 的 resolution 卷 对于 当前 statistical dimension。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Subtotal 的 resolution 卷 对于 当前 statistical dimension。",
						},
						"date_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "For daily 统计，它 是 statistical date。",
						},
						"hour_key": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For hourly 统计，它 是 hour 的 当前 时间 (0-23)，对于 示例，当 HourKey 是 23， statistical 周期 是 resolution 卷 从 22:00 到 23:00. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "域名 resolution 卷 统计 查询 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DATE: Statistics 通过 day dimension HOUR: Statistics 通过 hour dimension。",
						},
						"dns_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total resolution 卷 对于 当前 statistical 周期",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称 currently being queried。",
						},
						"start_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "开始时间 的 当前 statistical 周期",
						},
						"end_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间 的 当前 statistical 周期",
						},
					},
				},
			},

			"alias_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "域名 alias resolution 卷 统计 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "域名 resolution 卷 统计 查询 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_format": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DATE: Statistics 通过 day dimension HOUR: Statistics 通过 hour dimension。",
									},
									"dns_total": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total resolution 卷 对于 当前 statistical 周期",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "域名 名称 currently being queried。",
									},
									"start_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "开始时间 的 当前 statistical 周期",
									},
									"end_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "结束时间 的 当前 statistical 周期",
									},
								},
							},
						},
						"data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subtotal 的 resolution 卷 对于 当前 statistical dimension。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Subtotal 的 resolution 卷 对于 当前 statistical dimension。",
									},
									"date_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "For daily 统计，它 是 statistical date。",
									},
									"hour_key": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "For hourly 统计，它 是 hour 的 当前 时间 (0-23)，对于 示例，当 HourKey 是 23， statistical 周期 是 resolution 卷 从 22:00 到 23:00. 注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudDnspodDomainAnalyticsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnspod_domain_analytics.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var (
		domain    string
		aliasData []*dnspod.DomainAliasAnalyticsItem
		data      []*dnspod.DomainAnalyticsDetail
		info      *dnspod.DomainAnalyticsInfo
		err       error
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_date"); ok {
		paramMap["StartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dns_format"); ok {
		paramMap["DnsFormat"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		paramMap["DomainId"] = helper.IntUint64(v.(int))
	}

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	e := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		aliasData, data, info, err = service.DescribeDnspodDomainAnalyticsByFilter(ctx, paramMap)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if e != nil {
		return e
	}

	// ids := make([]string, 0, len(data))
	tmpDataList := make([]map[string]interface{}, 0, len(data))
	tmpAliasDataList := make([]map[string]interface{}, 0, len(aliasData))

	if data != nil {
		for _, domainAnalyticsDetail := range data {
			domainAnalyticsDetailMap := map[string]interface{}{}

			if domainAnalyticsDetail.Num != nil {
				domainAnalyticsDetailMap["num"] = domainAnalyticsDetail.Num
			}

			if domainAnalyticsDetail.DateKey != nil {
				domainAnalyticsDetailMap["date_key"] = domainAnalyticsDetail.DateKey
			}

			if domainAnalyticsDetail.HourKey != nil {
				domainAnalyticsDetailMap["hour_key"] = domainAnalyticsDetail.HourKey
			}

			// ids = append(ids, *domainAnalyticsDetail.Domain)
			tmpDataList = append(tmpDataList, domainAnalyticsDetailMap)
		}

		_ = d.Set("data", tmpDataList)
	}

	if info != nil {
		domainAnalyticsInfoMap := map[string]interface{}{}

		if info.DnsFormat != nil {
			domainAnalyticsInfoMap["dns_format"] = info.DnsFormat
		}

		if info.DnsTotal != nil {
			domainAnalyticsInfoMap["dns_total"] = info.DnsTotal
		}

		if info.Domain != nil {
			domainAnalyticsInfoMap["domain"] = info.Domain
		}

		if info.StartDate != nil {
			domainAnalyticsInfoMap["start_date"] = info.StartDate
		}

		if info.EndDate != nil {
			domainAnalyticsInfoMap["end_date"] = info.EndDate
		}

		// ids = append(ids, *info.Domain)
		// _ = d.Set("info", domainAnalyticsInfoMap)
		e = helper.SetMapInterfaces(d, "info", domainAnalyticsInfoMap)
		if e != nil {
			return e
		}
	}

	if aliasData != nil {
		for _, domainAliasAnalyticsItem := range aliasData {
			domainAliasAnalyticsItemMap := map[string]interface{}{}

			if domainAliasAnalyticsItem.Info != nil {
				infoMap := map[string]interface{}{}

				if domainAliasAnalyticsItem.Info.DnsFormat != nil {
					infoMap["dns_format"] = domainAliasAnalyticsItem.Info.DnsFormat
				}

				if domainAliasAnalyticsItem.Info.DnsTotal != nil {
					infoMap["dns_total"] = domainAliasAnalyticsItem.Info.DnsTotal
				}

				if domainAliasAnalyticsItem.Info.Domain != nil {
					infoMap["domain"] = domainAliasAnalyticsItem.Info.Domain
				}

				if domainAliasAnalyticsItem.Info.StartDate != nil {
					infoMap["start_date"] = domainAliasAnalyticsItem.Info.StartDate
				}

				if domainAliasAnalyticsItem.Info.EndDate != nil {
					infoMap["end_date"] = domainAliasAnalyticsItem.Info.EndDate
				}

				domainAliasAnalyticsItemMap["info"] = []interface{}{infoMap}
			}

			if domainAliasAnalyticsItem.Data != nil {
				dataList := []interface{}{}
				for _, data := range domainAliasAnalyticsItem.Data {
					dataMap := map[string]interface{}{}

					if data.Num != nil {
						dataMap["num"] = data.Num
					}

					if data.DateKey != nil {
						dataMap["date_key"] = data.DateKey
					}

					if data.HourKey != nil {
						dataMap["hour_key"] = data.HourKey
					}

					dataList = append(dataList, dataMap)
				}

				domainAliasAnalyticsItemMap["data"] = []interface{}{dataList}
			}

			// ids = append(ids, *domainAliasAnalyticsItem.Domain)
			tmpAliasDataList = append(tmpAliasDataList, domainAliasAnalyticsItemMap)
		}

		_ = d.Set("alias_data", tmpAliasDataList)
	}

	d.SetId(helper.DataResourceIdHash(domain))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		e = tccommon.WriteToFile(output.(string), map[string]interface{}{
			"info":       info,
			"data":       data,
			"alias_data": aliasData,
		})
		if e != nil {
			return e
		}
	}
	return nil
}
