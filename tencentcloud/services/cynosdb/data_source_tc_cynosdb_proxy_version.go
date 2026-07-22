package cynosdb

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbProxyVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbProxyVersionRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "集群 ID。",
			},
			"proxy_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "数据库代理组 ID。",
			},
			"support_proxy_versions": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "支持的数据库代理版本集合 注意：该字段可能返回null，表示无法获取有效值。",
			},
			"current_proxy_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "当前代理版本号 注意：该字段可能返回null，表示无法获取有效值。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudCynosdbProxyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_proxy_version.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		proxyVersion *cynosdb.DescribeSupportProxyVersionResponseParams
		clusterId    string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("proxy_group_id"); ok {
		paramMap["ProxyGroupId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbProxyVersionByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			e = fmt.Errorf("cynosdb proxy version not exists")
			return resource.NonRetryableError(e)
		}

		proxyVersion = result
		return nil
	})

	if err != nil {
		return err
	}

	if proxyVersion.SupportProxyVersions != nil {
		_ = d.Set("support_proxy_versions", proxyVersion.SupportProxyVersions)
	}

	if proxyVersion.CurrentProxyVersion != nil {
		_ = d.Set("current_proxy_version", proxyVersion.CurrentProxyVersion)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
