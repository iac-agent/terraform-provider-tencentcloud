package dnspod

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Description:  "域名 对于 其中 DNS records 是 到 是 获取。",
				Optional:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"domain", "domain_id"},
			},
			"domain_id": {
				Description:  "ID 域名 对于 其中 DNS records 是 到 是 获取. 如果 DomainId 是 passed 在， 系统 将 omit 参数 域名",
				Optional:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"domain", "domain_id"},
			},
			"group_id": {
				Description: "组 ID",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"keyword": {
				Description: "keyword 对于 searching 对于 DNS records. 主机 headers 和 记录 值 是 支持。",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"limit": {
				Description: "限制 It 默认为 100 和 可以 是 up 到 3,000。",
				Optional:    true,
				Type:        schema.TypeInt,
			},
			"offset": {
				Description: "偏移量 默认值：0。",
				Optional:    true,
				Type:        schema.TypeInt,
			},
			"record_count_info": {
				Computed:    true,
				Description: "Count info 的 queried 记录 列表。",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"list_count": {
						Computed:    true,
						Description: "count 的 records 返回 在 列表。",
						Type:        schema.TypeInt,
					},
					"subdomain_count": {
						Computed:    true,
						Description: "subdomain count。",
						Type:        schema.TypeInt,
					},
					"total_count": {
						Computed:    true,
						Description: "总数 记录 count。",
						Type:        schema.TypeInt,
					},
				}},
				Type: schema.TypeList,
			},
			"record_line": {
				Description: "split 可用区 名称",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"record_line_id": {
				Description: "split 可用区 ID. 如果 `record_line_id` 是 passed 在， 系统 将 omit 参数 `record_line`。",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"record_type": {
				Description: "类型 DNS 记录，such 作为 A，CNAME，NS，AAAA，explicit URL，implicit URL，CAA，或 SPF 记录。",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"result": {
				Computed:    true,
				Description: "记录 列表 结果",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"line": {
						Computed:    true,
						Description: "记录 split 可用区",
						Type:        schema.TypeString,
					},
					"line_id": {
						Computed:    true,
						Description: "split 可用区 ID。",
						Type:        schema.TypeString,
					},
					"monitor_status": {
						Computed:    true,
						Description: "监控 状态 记录. 有效值：OK (normal)，WARN (警告)，和 DOWN (downtime). It 是 空 如果 无 监控 是 集合 或 监控 是 suspended。",
						Type:        schema.TypeString,
					},
					"mx": {
						Computed:    true,
						Description: "MX 值，applicable 到 MX 记录 仅.\n注意：此字段可能返回 null，表示无法获取有效值。",
						Type:        schema.TypeInt,
					},
					"name": {
						Computed:    true,
						Description: "主机名",
						Type:        schema.TypeString,
					},
					"record_id": {
						Computed:    true,
						Description: "Record ID。",
						Type:        schema.TypeInt,
					},
					"remark": {
						Computed:    true,
						Description: "记录 备注",
						Type:        schema.TypeString,
					},
					"status": {
						Computed:    true,
						Description: "记录 状态 有效值：ENABLE (已启用)，DISABLE (已禁用)。",
						Type:        schema.TypeString,
					},
					"ttl": {
						Computed:    true,
						Description: "记录 缓存 时间。",
						Type:        schema.TypeInt,
					},
					"type": {
						Computed:    true,
						Description: "记录 类型",
						Type:        schema.TypeString,
					},
					"updated_on": {
						Computed:    true,
						Description: "更新时间。",
						Type:        schema.TypeString,
					},
					"value": {
						Computed:    true,
						Description: "记录 值",
						Type:        schema.TypeString,
					},
					"weight": {
						Computed:    true,
						Description: "记录 权重，其中 为必填项 对于 round-robin DNS records。",
						Type:        schema.TypeInt,
					},
				}},
				Type: schema.TypeList,
			},
			"result_output_file": {
				Description: "用于store 查询 结果 作为 JSON。",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"sort_field": {
				Description: "sorting 字段. 可用值：名称，line，类型，值，权重，mx，和 ttl,updated_on。",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"sort_type": {
				Description: "sorting 类型 有效值：ASC (ascending，默认值)，DESC (descending)。",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"subdomain": {
				Description: "主机 头部 的 DNS 记录. 如果 此 参数 是 passed 在，仅 DNS 记录 corresponding 到 此 主机 头部 将 是 返回。",
				Optional:    true,
				Type:        schema.TypeString,
			},
		},
	}
}
func dataSourceTencentCloudDnspodRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("tencentcloud_dnspod_records.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := DnspodService{client}

	request := dnspod.NewDescribeRecordListRequest()
	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		request.DomainId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("subdomain"); ok {
		request.Subdomain = helper.String(v.(string))
	}
	if v, ok := d.GetOk("record_type"); ok {
		request.RecordType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("record_line"); ok {
		request.RecordLine = helper.String(v.(string))
	}
	if v, ok := d.GetOk("record_line"); ok {
		request.RecordLineId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("keyword"); ok {
		request.Keyword = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sort_field"); ok {
		request.SortField = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sort_type"); ok {
		request.SortType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("offset"); ok {
		request.Offset = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		request.Limit = helper.IntUint64(v.(int))
	}

	list, info, err := service.DescribeRecordList(ctx, request)

	if err != nil {
		return err
	}

	d.SetId("dnspod_records" + helper.DataResourceIdHash(request.ToJsonString()))

	result := make([]map[string]interface{}, 0, len(list))
	for i := range list {
		record := list[i]
		result = append(result, map[string]interface{}{
			"line":           record.Line,
			"line_id":        record.LineId,
			"monitor_status": record.MonitorStatus,
			"mx":             record.MX,
			"name":           record.Name,
			"record_id":      record.RecordId,
			"remark":         record.Remark,
			"status":         record.Status,
			"ttl":            record.TTL,
			"type":           record.Type,
			"updated_on":     record.UpdatedOn,
			"value":          record.Value,
			"weight":         record.Weight,
		})
	}

	err = helper.SetMapInterfaces(d, "record_count_info", map[string]interface{}{
		"list_count":      info.ListCount,
		"subdomain_count": info.SubdomainCount,
		"total_count":     info.TotalCount,
	})
	if err != nil {
		return err
	}

	err = d.Set("result", result)
	if err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		err = tccommon.WriteToFile(output.(string), map[string]interface{}{
			"record_count_info": info,
			"result":            result,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
