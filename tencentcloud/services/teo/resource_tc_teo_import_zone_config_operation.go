package teo

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoImportZoneConfigOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoImportZoneConfigOperationCreate,
		Read:   resourceTencentCloudTeoImportZoneConfigOperationRead,
		Delete: resourceTencentCloudTeoImportZoneConfigOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The configuration content to import. It must be in JSON format and UTF-8 encoded. The content can be obtained via the ExportZoneConfig API.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task ID returned by the API, used for polling the import result.",
			},
		},
	}
}

func resourceTencentCloudTeoImportZoneConfigOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_import_zone_config_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	content := d.Get("content").(string)

	request := teov20220901.NewImportZoneConfigRequest()
	request.ZoneId = helper.String(zoneId)
	request.Content = helper.String(content)

	var taskId string
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ImportZoneConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result != nil && result.Response != nil && result.Response.TaskId != nil {
			taskId = *result.Response.TaskId
		}
		return nil
	})
	if err != nil {
		return err
	}

	if taskId == "" {
		return fmt.Errorf("ImportZoneConfig returned empty TaskId")
	}

	// Poll DescribeZoneConfigImportResult until status is success or failure
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{"doing"}, []string{"success"}, 10*tccommon.ReadRetryTimeout, time.Second, service.TeoImportZoneConfigStateRefreshFunc(zoneId, taskId, []string{"failure"}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	_ = d.Set("task_id", taskId)

	d.SetId(zoneId + tccommon.FILED_SP + taskId)

	return resourceTencentCloudTeoImportZoneConfigOperationRead(d, meta)
}

func resourceTencentCloudTeoImportZoneConfigOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_import_zone_config_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoImportZoneConfigOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_import_zone_config_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
