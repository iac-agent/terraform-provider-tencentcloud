package gaap

import (
	"context"
	"errors"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapProxies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapProxiesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"project_id", "access_region", "realserver_region"},
				Description:   "ID GAAP proxy to be queried. Conflict with `project_id`，`access_region` and `realserver_region`。",
			},
			"project_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "项目 ID GAAP proxy to be queried. Conflict with `ids`。",
			},
			"access_region": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				ValidateFunc:  tccommon.ValidateAllowedStringValue([]string{"NorthChina", "Beijing", "EastChina", "Shanghai", "SouthChina", "Guangzhou", "SouthwestChina", "Chengdu", "Hongkong", "SL_TAIWAN", "SoutheastAsia", "Korea", "SL_India", "SL_Australia", "Europe", "SL_UK", "SL_SouthAmerica", "NorthAmerica", "SL_MiddleUSA", "Canada", "SL_VIET", "WestIndia", "Thailand", "Virginia", "Russia", "Japan", "SL_Indonesia"}),
				Description:   "Access 地域 of the GAAP proxy to be queried. Conflict with `ids`。",
			},
			"realserver_region": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				ValidateFunc:  tccommon.ValidateAllowedStringValue([]string{"NorthChina", "Beijing", "EastChina", "Shanghai", "SouthChina", "Guangzhou", "SouthwestChina", "Chengdu", "Hongkong", "SL_TAIWAN", "SoutheastAsia", "Korea", "SL_India", "SL_Australia", "Europe", "SL_UK", "SL_SouthAmerica", "NorthAmerica", "SL_MiddleUSA", "Canada", "SL_VIET", "WestIndia", "Thailand", "Virginia", "Russia", "Japan", "SL_Indonesia"}),
				Description:   "地域 of the GAAP realserver to be queried. Conflict with `ids`。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 of the GAAP proxy to be queried. Support up to 5，display the information as long as it matches one。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"proxies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information 列表 GAAP proxy. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID GAAP proxy。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 GAAP proxy。",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security policy ID GAAP proxy。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access 域名 of the GAAP proxy。",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access 域名 of the GAAP proxy。",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum bandwidth of the GAAP proxy，unit is Mbps。",
						},
						"concurrent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 并发 of the GAAP proxy，unit is 10k。",
						},
						"access_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access 地域 of the GAAP proxy。",
						},
						"realserver_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 of the GAAP realserver。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID project within the GAAP proxy，'0' means is default project。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of the GAAP proxy。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 GAAP proxy。",
						},
						"scalable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否GAAP proxy can scalable。",
						},
						"support_protocols": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Supported protocols of the GAAP proxy。",
						},
						"forward_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding IP of the GAAP proxy。",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 of the GAAP proxy。",
						},
						"is_auto_scale_proxy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "表示是否auto scale channel is 已启用，with 0 for no and 1 for yes。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 of the GAAP proxy。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapProxiesRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_proxies.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		ids              []string
		projectId        *int
		accessRegion     *string
		realserverRegion *string
	)

	if raw, ok := d.GetOk("ids"); ok {
		set := raw.(*schema.Set).List()
		ids = make([]string, 0, len(set))
		for _, id := range set {
			ids = append(ids, id.(string))
		}
	}

	if raw, ok := d.GetOkExists("project_id"); ok {
		projectId = common.IntPtr(raw.(int))
	}

	if raw, ok := d.GetOk("access_region"); ok {
		accessRegion = helper.String(raw.(string))
	}

	if raw, ok := d.GetOk("realserver_region"); ok {
		realserverRegion = helper.String(raw.(string))
	}

	tags := helper.GetTags(d, "tags")

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	proxies, err := service.DescribeProxies(ctx, ids, projectId, accessRegion, realserverRegion, tags)
	if err != nil {
		return err
	}

	proxyIds := make([]string, 0, len(proxies))
	respProxies := make([]map[string]interface{}, 0, len(proxies))
	for _, proxy := range proxies {
		if proxy.ProxyId == nil {
			return errors.New("proxy id is nil")
		}
		if proxy.ProxyName == nil {
			return errors.New("proxy name is nil")
		}
		if proxy.Domain == nil {
			return errors.New("proxy domain is nil")
		}
		if proxy.IP == nil {
			return errors.New("proxy ip is nil")
		}
		if proxy.Bandwidth == nil {
			return errors.New("proxy bandwidth is nil")
		}
		if proxy.Concurrent == nil {
			return errors.New("proxy concurrent is nil")
		}
		if proxy.AccessRegion == nil {
			return errors.New("proxy access region is nil")
		}
		if proxy.RealServerRegion == nil {
			return errors.New("proxy realserver region is nil")
		}
		if proxy.ProjectId == nil {
			return errors.New("proxy project id is nil")
		}
		if proxy.CreateTime == nil {
			return errors.New("proxy create time is nil")
		}
		if proxy.Status == nil {
			return errors.New("proxy status is nil")
		}
		if proxy.Scalarable == nil {
			return errors.New("proxy scalable is nil")
		}
		if proxy.SupportProtocols == nil {
			return errors.New("proxy support protocols is nil")
		}
		if proxy.ForwardIP == nil {
			return errors.New("proxy forward ip is nil")
		}
		if proxy.Version == nil {
			return errors.New("proxy version is nil")
		}
		if proxy.IsAutoScaleProxy == nil {
			return errors.New("proxy IsAutoScaleProxy is nil")
		}

		proxyIds = append(proxyIds, *proxy.ProxyId)

		m := map[string]interface{}{
			"id":                  *proxy.ProxyId,
			"name":                *proxy.ProxyName,
			"domain":              *proxy.Domain,
			"ip":                  *proxy.IP,
			"bandwidth":           *proxy.Bandwidth,
			"concurrent":          *proxy.Concurrent,
			"access_region":       *proxy.AccessRegion,
			"realserver_region":   *proxy.RealServerRegion,
			"project_id":          *proxy.ProjectId,
			"create_time":         helper.FormatUnixTime(*proxy.CreateTime),
			"status":              *proxy.Status,
			"scalable":            *proxy.Scalarable == 1,
			"forward_ip":          *proxy.ForwardIP,
			"version":             *proxy.Version,
			"is_auto_scale_proxy": *proxy.IsAutoScaleProxy,
		}

		if proxy.PolicyId != nil {
			m["policy_id"] = *proxy.PolicyId
		}

		supportProtocols := make([]string, 0, len(proxy.SupportProtocols))
		for _, v := range proxy.SupportProtocols {
			supportProtocols = append(supportProtocols, *v)
		}
		m["support_protocols"] = supportProtocols

		tags := make(map[string]string, len(proxy.TagSet))
		for _, t := range proxy.TagSet {
			if t.TagKey == nil {
				return errors.New("tag key is nil")
			}
			if t.TagValue == nil {
				return errors.New("tag value is nil")
			}
			tags[*t.TagKey] = *t.TagValue
		}
		m["tags"] = tags

		respProxies = append(respProxies, m)
	}

	_ = d.Set("proxies", respProxies)
	d.SetId(helper.DataResourceIdsHash(proxyIds))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), respProxies); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
