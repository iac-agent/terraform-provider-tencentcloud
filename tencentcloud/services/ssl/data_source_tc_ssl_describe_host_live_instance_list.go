package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeHostLiveInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostLiveInstanceListRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "证书 ID 到 是 deployed。",
			},

			"resource_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy 资源类型",
			},

			"is_cache": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否query 缓存，1: Yes; 0: No， 默认为 查询 缓存， 缓存 是 half hour。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "列表 filtering 参数; Filterkey: domainmatch。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤参数键",
						},
						"filter_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤参数值",
						},
					},
				},
			},

			"old_certificate_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Deployed 证书 ID",
			},

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Live 实例 listNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Binded 证书 IDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "-1: Unrelated 证书 的 域名 名称1: 域名 名称 HTTPS has been opened.0: 域名 名称 HTTPS has been closed。",
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

func dataSourceTencentCloudSslDescribeHostLiveInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_host_live_instance_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("certificate_id"); ok {
		paramMap["CertificateId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		paramMap["ResourceType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("is_cache"); v != nil {
		paramMap["IsCache"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*ssl.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := ssl.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["filter_key"]; ok {
				filter.FilterKey = helper.String(v.(string))
			}
			if v, ok := filterMap["filter_value"]; ok {
				filter.FilterValue = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("old_certificate_id"); ok {
		paramMap["OldCertificateId"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceList []*ssl.LiveInstanceDetail

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeHostLiveInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	tmpList := make([]map[string]interface{}, 0, len(instanceList))

	if instanceList != nil {
		for _, liveInstanceDetail := range instanceList {
			liveInstanceDetailMap := map[string]interface{}{}

			if liveInstanceDetail.Domain != nil {
				liveInstanceDetailMap["domain"] = liveInstanceDetail.Domain
			}

			if liveInstanceDetail.CertId != nil {
				liveInstanceDetailMap["cert_id"] = liveInstanceDetail.CertId
			}

			if liveInstanceDetail.Status != nil {
				liveInstanceDetailMap["status"] = liveInstanceDetail.Status
			}

			ids = append(ids, *liveInstanceDetail.CertId)
			tmpList = append(tmpList, liveInstanceDetailMap)
		}

		_ = d.Set("instance_list", tmpList)
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
