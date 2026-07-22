package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsExportCreate,
		Read:   resourceTencentCloudClsExportRead,
		Delete: resourceTencentCloudClsExportDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "主题 ID。",
			},

			"query": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "导出查询规则。",
			},

			"log_count": {
				Required:    true,
				Type:        schema.TypeInt,
				ForceNew:    true,
				Description: "导出日志量。",
			},

			"from": {
				Required:    true,
				Type:        schema.TypeInt,
				ForceNew:    true,
				Description: "导出开始时间。",
			},

			"to": {
				Required:    true,
				Type:        schema.TypeInt,
				ForceNew:    true,
				Description: "导出结束时间。",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "日志导出时间排序。降序或升序。",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "日志导出格式。",
			},
		},
	}
}

func resourceTencentCloudClsExportCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_export.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cls.NewCreateExportRequest()
		response = cls.NewCreateExportResponse()
		topicId  string
		exportId string
	)
	if v, ok := d.GetOk("topic_id"); ok {
		topicId = v.(string)
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("log_count"); ok {
		request.Count = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("query"); ok {
		request.Query = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("from"); ok {
		request.From = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("to"); ok {
		request.To = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("order"); ok {
		request.Order = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateExport(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls export failed, reason:%+v", logId, err)
		return err
	}

	exportId = *response.Response.ExportId
	d.SetId(topicId + tccommon.FILED_SP + exportId)

	return resourceTencentCloudClsExportRead(d, meta)
}

func resourceTencentCloudClsExportRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_export.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	exportId := idSplit[1]

	export, err := service.DescribeClsExportById(ctx, topicId, exportId)
	if err != nil {
		return err
	}

	if export == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsExport` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if export.TopicId != nil {
		_ = d.Set("topic_id", export.TopicId)
	}

	if export.Count != nil {
		_ = d.Set("log_count", export.Count)
	}

	if export.Query != nil {
		_ = d.Set("query", export.Query)
	}

	if export.From != nil {
		_ = d.Set("from", export.From)
	}

	if export.To != nil {
		_ = d.Set("to", export.To)
	}

	if export.Order != nil {
		_ = d.Set("order", export.Order)
	}

	if export.Format != nil {
		_ = d.Set("format", export.Format)
	}

	return nil
}

func resourceTencentCloudClsExportDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_export.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)

	exportId := idSplit[1]

	if err := service.DeleteClsExportById(ctx, exportId); err != nil {
		return err
	}

	return nil
}
