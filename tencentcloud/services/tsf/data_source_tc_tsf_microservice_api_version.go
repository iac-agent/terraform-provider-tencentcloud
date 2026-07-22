package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfMicroserviceApiVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfMicroserviceApiVersionRead,
		Schema: map[string]*schema.Schema{
			"microservice_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Microservice ID。",
			},

			"path": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "api 路径",
			},

			"method": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "请求 方法。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API 版本 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application ID。",
						},
						"application_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application 名称",
						},
						"pkg_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "应用 pkg 版本",
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

func dataSourceTencentCloudTsfMicroserviceApiVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_microservice_api_version.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("microservice_id"); ok {
		paramMap["MicroserviceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		paramMap["Path"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("method"); ok {
		paramMap["Method"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var apiVersion []*tsf.ApiVersionArray
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfMicroserviceApiVersionByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		apiVersion = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(apiVersion))
	tmpList := make([]map[string]interface{}, 0, len(apiVersion))
	if apiVersion != nil {
		for _, apiVersionArray := range apiVersion {
			apiVersionArrayMap := map[string]interface{}{}

			if apiVersionArray.ApplicationId != nil {
				apiVersionArrayMap["application_id"] = apiVersionArray.ApplicationId
			}

			if apiVersionArray.ApplicationName != nil {
				apiVersionArrayMap["application_name"] = apiVersionArray.ApplicationName
			}

			if apiVersionArray.PkgVersion != nil {
				apiVersionArrayMap["pkg_version"] = apiVersionArray.PkgVersion
			}

			ids = append(ids, *apiVersionArray.ApplicationId)
			tmpList = append(tmpList, apiVersionArrayMap)
		}

		_ = d.Set("result", tmpList)
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
