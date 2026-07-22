package apigateway

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudApiGatewayServiceReleaseVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApiGatewayServiceReleaseVersionsRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID service to be queried。",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 service releases.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 number.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"version_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 描述注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudApiGatewayServiceReleaseVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_service_release_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		versionList []*apigateway.DescribeServiceReleaseVersionResultVersionListInfo
		serviceId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
		serviceId = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApiGatewayServiceReleaseVersionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		versionList = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(versionList))
	if versionList != nil {
		for _, version := range versionList {
			versionListMap := map[string]interface{}{}
			if version.VersionName != nil {
				versionListMap["version_name"] = version.VersionName
			}

			if version.VersionDesc != nil {
				versionListMap["version_desc"] = version.VersionDesc
			}

			tmpList = append(tmpList, versionListMap)
		}

		_ = d.Set("result", tmpList)
	}

	d.SetId(serviceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
