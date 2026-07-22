package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeHostCosInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostCosInstanceListRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "证书 ID 到 是 deployed。",
			},

			"resource_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy 资源类型 cos。",
			},

			"is_cache": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否query 缓存，1: Yes; 0: No， 默认为 查询 缓存， 缓存 是 half hour。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "列表 过滤器 参数。",
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

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "COS 实例 listNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
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
							Type:        schema.TypeString,
							Computed:    true,
							Description: "已启用: 域名 名称 online statusDisabled: 域名 名称 offline 状态",
						},
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reserve 存储桶 nameNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Barrel areaNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
					},
				},
			},

			"async_total_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "总数 数量 asynchronous refreshNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"async_offset": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Asynchronous refresh 当前 execution numberNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"async_cache_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current 缓存 read timeNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeHostCosInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_host_cos_instance_list.read")()
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

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceList *ssl.DescribeHostCosInstanceListResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeHostCosInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	tmpList := make([]map[string]interface{}, 0)

	if instanceList != nil && instanceList.InstanceList != nil {

		for _, cosInstanceDetail := range instanceList.InstanceList {
			cosInstanceDetailMap := map[string]interface{}{}

			if cosInstanceDetail.Domain != nil {
				cosInstanceDetailMap["domain"] = cosInstanceDetail.Domain
			}

			if cosInstanceDetail.CertId != nil {
				cosInstanceDetailMap["cert_id"] = cosInstanceDetail.CertId
			}

			if cosInstanceDetail.Status != nil {
				cosInstanceDetailMap["status"] = cosInstanceDetail.Status
			}

			if cosInstanceDetail.Bucket != nil {
				cosInstanceDetailMap["bucket"] = cosInstanceDetail.Bucket
			}

			if cosInstanceDetail.Region != nil {
				cosInstanceDetailMap["region"] = cosInstanceDetail.Region
			}

			ids = append(ids, *cosInstanceDetail.CertId)
			tmpList = append(tmpList, cosInstanceDetailMap)
		}

		_ = d.Set("instance_list", tmpList)
	}

	if instanceList != nil && instanceList.AsyncTotalNum != nil {
		_ = d.Set("async_total_num", instanceList.AsyncTotalNum)
	}

	if instanceList != nil && instanceList.AsyncOffset != nil {
		_ = d.Set("async_offset", instanceList.AsyncOffset)
	}

	if instanceList != nil && instanceList.AsyncCacheTime != nil {
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
