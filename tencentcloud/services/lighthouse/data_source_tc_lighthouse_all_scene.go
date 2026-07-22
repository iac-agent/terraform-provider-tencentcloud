package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseAllScene() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseAllSceneRead,
		Schema: map[string]*schema.Schema{
			"scene_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 scene IDs。",
			},

			"offset": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "偏移量 默认值为 0。",
			},

			"limit": {
				Optional:    true,
				Default:     20,
				Type:        schema.TypeInt,
				Description: "数量 返回 results. 默认值为 20. Maximum 值 是 100。",
			},

			"scene_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 scene info。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scene_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use scene ID。",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use scene presentation 名称",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use scene 描述",
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

func dataSourceTencentCloudLighthouseAllSceneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_scene.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("scene_ids"); ok {
		sceneIdsSet := v.(*schema.Set).List()
		sceneIds := make([]string, 0)
		for _, sceneId := range sceneIdsSet {
			sceneIds = append(sceneIds, sceneId.(string))
		}
		paramMap["scene_ids"] = sceneIds
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var sceneSet []*lighthouse.SceneInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseAllSceneByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		sceneSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(sceneSet))
	tmpList := make([]map[string]interface{}, 0, len(sceneSet))
	for _, scene := range sceneSet {
		ids = append(ids, *scene.SceneId)
		tmpList = append(tmpList, map[string]interface{}{
			"scene_id":     *scene.SceneId,
			"display_name": *scene.DisplayName,
			"description":  *scene.Description,
		})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("scene_set", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
