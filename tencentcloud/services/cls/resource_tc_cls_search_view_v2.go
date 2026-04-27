package cls

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsSearchViewV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsSearchViewV2Create,
		Read:   resourceTencentCloudClsSearchViewV2Read,
		Update: resourceTencentCloudClsSearchViewV2Update,
		Delete: resourceTencentCloudClsSearchViewV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"logset_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Logset ID to which the search view belongs.",
			},

			"logset_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the logset, e.g., ap-guangzhou.",
			},

			"view_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Search view name, max 255 characters.",
			},

			"view_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Search view type. Valid values: log, metric.",
			},

			"topics": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    10,
				Description: "Topics included in the search view, max 10 topics.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region of the topic.",
						},
						"logset_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Logset ID of the topic.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic ID.",
						},
					},
				},
			},

			"view_id_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom search view ID prefix.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the search view.",
			},

			"view_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Search view ID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time.",
			},
		},
	}
}

func resourceTencentCloudClsSearchViewV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_search_view_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = cls.NewCreateSearchViewRequest()
		response = cls.NewCreateSearchViewResponse()
	)

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_region"); ok {
		request.LogsetRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("view_name"); ok {
		request.ViewName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("view_type"); ok {
		request.ViewType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topics"); ok {
		for _, item := range v.([]interface{}) {
			topicMap := item.(map[string]interface{})
			viewSearchTopic := cls.ViewSearchTopic{}
			if v, ok := topicMap["region"]; ok {
				viewSearchTopic.Region = helper.String(v.(string))
			}
			if v, ok := topicMap["logset_id"]; ok {
				viewSearchTopic.LogsetId = helper.String(v.(string))
			}
			if v, ok := topicMap["topic_id"]; ok {
				viewSearchTopic.TopicId = helper.String(v.(string))
			}
			request.Topics = append(request.Topics, &viewSearchTopic)
		}
	}

	if v, ok := d.GetOk("view_id_prefix"); ok {
		request.ViewIdPrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateSearchViewWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls search view v2 failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls search view v2 failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.ViewId == nil {
		return fmt.Errorf("ViewId is nil.")
	}

	viewId := *response.Response.ViewId
	d.SetId(viewId)

	return resourceTencentCloudClsSearchViewV2Read(d, meta)
}

func resourceTencentCloudClsSearchViewV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_search_view_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	viewId := d.Id()

	request := cls.NewDescribeSearchViewsRequest()
	filter := cls.Filter{
		Key:    helper.String("viewId"),
		Values: []*string{helper.String(viewId)},
	}
	request.Filters = append(request.Filters, &filter)
	request.Limit = helper.Uint64(100)

	var infos []*cls.SearchViewInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DescribeSearchViewsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Read cls search view v2 failed, Response is nil."))
		}

		infos = result.Response.Infos
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cls search view v2 failed, reason:%+v", logId, err)
		return err
	}

	if len(infos) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_cls_search_view_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	searchViewInfo := infos[0]

	if searchViewInfo.ViewId != nil {
		_ = d.Set("view_id", searchViewInfo.ViewId)
	}

	if searchViewInfo.LogsetId != nil {
		_ = d.Set("logset_id", searchViewInfo.LogsetId)
	}

	if searchViewInfo.LogsetRegion != nil {
		_ = d.Set("logset_region", searchViewInfo.LogsetRegion)
	}

	if searchViewInfo.ViewName != nil {
		_ = d.Set("view_name", searchViewInfo.ViewName)
	}

	if searchViewInfo.ViewType != nil {
		_ = d.Set("view_type", searchViewInfo.ViewType)
	}

	if searchViewInfo.Topics != nil && len(searchViewInfo.Topics) > 0 {
		topicsList := make([]map[string]interface{}, 0, len(searchViewInfo.Topics))
		for _, topic := range searchViewInfo.Topics {
			topicMap := map[string]interface{}{}
			if topic.Region != nil {
				topicMap["region"] = topic.Region
			}
			if topic.LogsetId != nil {
				topicMap["logset_id"] = topic.LogsetId
			}
			if topic.TopicId != nil {
				topicMap["topic_id"] = topic.TopicId
			}
			topicsList = append(topicsList, topicMap)
		}
		_ = d.Set("topics", topicsList)
	}

	if searchViewInfo.Description != nil {
		_ = d.Set("description", searchViewInfo.Description)
	}

	if searchViewInfo.CreateTime != nil {
		_ = d.Set("create_time", helper.UInt64ToStr(*searchViewInfo.CreateTime))
	}

	if searchViewInfo.UpdateTime != nil {
		_ = d.Set("update_time", helper.UInt64ToStr(*searchViewInfo.UpdateTime))
	}

	return nil
}

func resourceTencentCloudClsSearchViewV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_search_view_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	viewId := d.Id()

	needChange := false
	mutableArgs := []string{"view_name", "view_type", "topics", "description"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cls.NewModifySearchViewRequest()
		request.ViewId = helper.String(viewId)

		if v, ok := d.GetOk("view_name"); ok {
			request.ViewName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("view_type"); ok {
			request.ViewType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("topics"); ok {
			for _, item := range v.([]interface{}) {
				topicMap := item.(map[string]interface{})
				viewSearchTopic := cls.ViewSearchTopic{}
				if v, ok := topicMap["region"]; ok {
					viewSearchTopic.Region = helper.String(v.(string))
				}
				if v, ok := topicMap["logset_id"]; ok {
					viewSearchTopic.LogsetId = helper.String(v.(string))
				}
				if v, ok := topicMap["topic_id"]; ok {
					viewSearchTopic.TopicId = helper.String(v.(string))
				}
				request.Topics = append(request.Topics, &viewSearchTopic)
			}
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifySearchViewWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cls search view v2 failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClsSearchViewV2Read(d, meta)
}

func resourceTencentCloudClsSearchViewV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_search_view_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cls.NewDeleteSearchViewRequest()
	)

	viewId := d.Id()
	request.ViewId = helper.String(viewId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DeleteSearchViewWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cls search view v2 failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
