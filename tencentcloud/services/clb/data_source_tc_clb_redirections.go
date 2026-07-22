package clb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbRedirections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbRedirectionsRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "需要查询的CLB ID。",
			},
			"source_listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "需要查询的源监听ID。",
			},
			"target_listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的目标监听ID。",
			},
			"source_rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "需要查询的源监听的规则ID。",
			},
			"target_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的目标监听的规则ID。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"redirection_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "云负载均衡器重定向配置列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB的ID。",
						},
						"source_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "源监听器ID。",
						},
						"target_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标监听者ID。",
						},
						"source_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "源监听器的规则ID。",
						},
						"target_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标监听器的规则ID。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbRedirectionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_redirections.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]string)
	params["clb_id"] = d.Get("clb_id").(string)
	params["source_listener_id"] = d.Get("source_listener_id").(string)
	checkErr := ListenerIdCheck(params["source_listener_id"])
	if checkErr != nil {
		return checkErr
	}
	params["source_rule_id"] = d.Get("source_rule_id").(string)
	checkErr = RuleIdCheck(params["source_rule_id"])
	if checkErr != nil {
		return checkErr
	}
	if v, ok := d.GetOk("target_listener_id"); ok {
		params["target_listener_id"] = v.(string)
		checkErr := ListenerIdCheck(params["target_listener_id"])
		if checkErr != nil {
			return checkErr
		}
	}
	if v, ok := d.GetOk("target_rule_id"); ok {
		params["target_rule_id"] = v.(string)
		checkErr = RuleIdCheck(params["target_rule_id"])
		if checkErr != nil {
			return checkErr
		}
	}

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var redirections []*map[string]string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeRedirectionsByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		redirections = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB redirections failed, reason:%+v", logId, err)
		return err
	}
	redirectionList := make([]map[string]interface{}, 0, len(redirections))
	ids := make([]string, 0, len(redirections))
	for _, r := range redirections {
		mapping := map[string]interface{}{
			"clb_id":             (*r)["clb_id"],
			"source_listener_id": (*r)["source_listener_id"],
			"target_listener_id": (*r)["target_listener_id"],
			"source_rule_id":     (*r)["source_rule_id"],
			"target_rule_id":     (*r)["target_rule_id"],
		}

		redirectionList = append(redirectionList, mapping)
		ids = append(ids, (*r)["source_rule_id"]+"#"+(*r)["target_rule_id"]+(*r)["source_listener_id"]+"#"+(*r)["target_listener_id"]+"#"+(*r)["clb_id"])
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("redirection_list", redirectionList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB redirection list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), redirectionList); e != nil {
			return e
		}
	}

	return nil
}
