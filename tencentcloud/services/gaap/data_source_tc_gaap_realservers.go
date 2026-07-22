package gaap

import (
	"context"
	"errors"
	"log"
	"net"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapRealservers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapRealserversRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeInt,
				Default:     -1,
				Optional:    true,
				Description: "ID 项目 within GAAP realserver 到 是 queried，默认值为 `-1`，无 集合 表示 all projects。",
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ip"},
				Description:   "域名 的 GAAP realserver 到 是 queried，conflict 使用 `ip`。",
			},
			"ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"domain"},
				Description:   "IP 的 GAAP realserver 到 是 queried，conflict 使用 `域名`。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 GAAP realserver 到 是 queried， 最大 长度 是 30。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 GAAP proxy 到 是 queried. Support up 到 5，display 信息 作为 long 作为 它 matches 一个。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"realservers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 GAAP realserver. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID GAAP realserver。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 GAAP realserver。",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 的 GAAP realserver。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 的 GAAP realserver。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 项目 within GAAP realserver。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 GAAP realserver。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapRealserversRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_realservers.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	projectId := d.Get("project_id").(int)

	var (
		address *string
		name    *string
	)
	if raw, ok := d.GetOk("ip"); ok {
		address = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("domain"); ok {
		address = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = helper.String(raw.(string))
	}

	tags := helper.GetTags(d, "tags")

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	realservers, err := service.DescribeRealservers(ctx, address, name, tags, projectId)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(realservers))
	realserverList := make([]map[string]interface{}, 0, len(realservers))
	for _, rs := range realservers {
		if rs.RealServerId == nil {
			return errors.New("realserver id is nil")
		}
		if rs.RealServerName == nil {
			return errors.New("realserver name is nil")
		}
		if rs.RealServerIP == nil {
			return errors.New("realserver name is nil")
		}
		if rs.ProjectId == nil {
			return errors.New("realserver project id is nil")
		}

		ids = append(ids, *rs.RealServerId)

		m := map[string]interface{}{
			"id":         *rs.RealServerId,
			"name":       *rs.RealServerName,
			"project_id": *rs.ProjectId,
		}

		if net.ParseIP(*rs.RealServerIP) == nil {
			m["domain"] = *rs.RealServerIP
		} else {
			m["ip"] = *rs.RealServerIP
		}

		if len(rs.TagSet) > 0 {
			tags := make(map[string]string, len(rs.TagSet))
			for _, tag := range rs.TagSet {
				if tag.TagKey == nil {
					return errors.New("tag key is nil")
				}
				if tag.TagValue == nil {
					return errors.New("tag value is nil")
				}
				tags[*tag.TagKey] = *tag.TagValue
			}
			m["tags"] = tags
		}

		realserverList = append(realserverList, m)
	}

	_ = d.Set("realservers", realserverList)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), realserverList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
