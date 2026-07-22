package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeHostCdnInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostCdnInstanceListRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "证书 ID to be deployed。",
			},

			"resource_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy 资源类型",
			},

			"is_cache": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否query the cache，1: Yes; 0: No，the 默认为 the query cache，the cache is half an hour。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "列表 filtering parameters; Filterkey: domainmatch。",
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
				Description: "Original 证书 ID",
			},

			"async_cache": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether。",
			},

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "CDN instance listNote: This field may return NULL，indicating that the valid 值 cannot be obtained。",
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
							Description: "Deployment 证书 ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称",
						},
						"https_billing_switch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称 billing 状态",
						},
					},
				},
			},

			"async_total_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total 数量 asynchronous refreshNote: This field may return NULL，indicating that the valid 值 cannot be obtained。",
			},

			"async_offset": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Asynchronous refresh current execution numberNote: This field may return NULL，indicating that the valid 值 cannot be obtained。",
			},

			"async_cache_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current cache read timeNote: This field may return NULL，indicating that the valid 值 cannot be obtained。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeHostCdnInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_host_cdn_instance_list.read")()
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

	if v, _ := d.GetOk("async_cache"); v != nil {
		paramMap["AsyncCache"] = helper.IntInt64(v.(int))
	}

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceList *ssl.DescribeHostCdnInstanceListResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeHostCdnInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList.InstanceList))
	tmpList := make([]map[string]interface{}, 0, len(instanceList.InstanceList))

	if instanceList != nil && instanceList.InstanceList != nil {
		for _, cdnInstanceDetail := range instanceList.InstanceList {
			cdnInstanceDetailMap := map[string]interface{}{}

			if cdnInstanceDetail.Domain != nil {
				cdnInstanceDetailMap["domain"] = cdnInstanceDetail.Domain
			}

			if cdnInstanceDetail.CertId != nil {
				cdnInstanceDetailMap["cert_id"] = cdnInstanceDetail.CertId
			}

			if cdnInstanceDetail.Status != nil {
				cdnInstanceDetailMap["status"] = cdnInstanceDetail.Status
			}

			if cdnInstanceDetail.HttpsBillingSwitch != nil {
				cdnInstanceDetailMap["https_billing_switch"] = cdnInstanceDetail.HttpsBillingSwitch
			}

			ids = append(ids, *cdnInstanceDetail.CertId)
			tmpList = append(tmpList, cdnInstanceDetailMap)
		}

		_ = d.Set("instance_list", tmpList)
	}

	if instanceList.AsyncTotalNum != nil {
		_ = d.Set("async_total_num", instanceList.AsyncTotalNum)
	}

	if instanceList.AsyncOffset != nil {
		_ = d.Set("async_offset", instanceList.AsyncOffset)
	}

	if instanceList.AsyncCacheTime != nil {
		_ = d.Set("async_cache_time", instanceList.AsyncCacheTime)
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
