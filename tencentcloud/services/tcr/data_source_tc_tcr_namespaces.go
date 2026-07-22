package tcr

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTCRNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTCRNamespacesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID 实例 该 命名空间 belongs 到。",
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID TCR 命名空间 到 查询。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			// Computed values
			"namespace_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 dedicated TCR namespaces。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 TCR 命名空间。",
						},
						"is_public": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate 该 命名空间 是 公有 或 不。",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID TCR 命名空间。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTCRNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_namespaces.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var name, instanceId string
	instanceId = d.Get("instance_id").(string)
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	tcrService := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error
	namespaces, outErr := tcrService.DescribeTCRNameSpaces(ctx, instanceId, name)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			namespaces, inErr = tcrService.DescribeTCRNameSpaces(ctx, instanceId, name)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}

	ids := make([]string, 0, len(namespaces))
	namespaceList := make([]map[string]interface{}, 0, len(namespaces))
	for _, namespace := range namespaces {
		mapping := map[string]interface{}{
			"name":      namespace.Name,
			"is_public": namespace.Public,
			"id":        namespace.NamespaceId,
		}

		namespaceList = append(namespaceList, mapping)
		ids = append(ids, instanceId+tccommon.FILED_SP+*namespace.Name)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("namespace_list", namespaceList); e != nil {
		log.Printf("[CRITAL]%s provider set TCR namespace list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), namespaceList); e != nil {
			return e
		}
	}

	return nil

}
