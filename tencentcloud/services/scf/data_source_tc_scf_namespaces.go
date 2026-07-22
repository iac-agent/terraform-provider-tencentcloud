package scf

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfNamespacesRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 SCF 命名空间 到 是 queried。",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "描述 SCF 命名空间 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"namespaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 命名空间. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 SCF 命名空间。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 SCF 命名空间。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 SCF 命名空间。",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "修改时间 的 SCF 命名空间。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 SCF 命名空间。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudScfNamespacesRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_namespaces.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		namespace *string
		desc      *string
	)

	if raw, ok := d.GetOk("namespace"); ok {
		namespace = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("description"); ok {
		desc = helper.String(raw.(string))
	}

	nss, err := service.DescribeNamespaces(ctx)
	if err != nil {
		log.Printf("[CRITAL]%s read namespace list failed: %+v", logId, err)
		return err
	}

	namespaces := make([]map[string]*string, 0, len(nss))
	ids := make([]string, 0, len(nss))

	for _, ns := range nss {
		if namespace != nil && !strings.Contains(*ns.Name, *namespace) {
			continue
		}
		if desc != nil && !strings.Contains(*ns.Description, *desc) {
			continue
		}

		ids = append(ids, *ns.Name)

		namespaces = append(namespaces, map[string]*string{
			"namespace":   ns.Name,
			"description": ns.Description,
			"create_time": ns.AddTime,
			"modify_time": ns.ModTime,
			"type":        ns.Type,
		})
	}

	_ = d.Set("namespaces", namespaces)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), namespaces); err != nil {
			err = errors.WithStack(err)
			log.Printf("[CRITAL]%s output file[%s] fail, reason: %+v", logId, output.(string), err)
			return err
		}
	}

	return nil
}
