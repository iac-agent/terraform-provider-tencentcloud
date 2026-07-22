package dbbrain

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDbbrainModifyDiagDbInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbbrainModifyDiagDbInstanceOperationCreate,
		Read:   resourceTencentCloudDbbrainModifyDiagDbInstanceOperationRead,
		Delete: resourceTencentCloudDbbrainModifyDiagDbInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_confs": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "实例 配置, 包括 inspection, overview switch, etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily_inspection": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database inspection switch, Yes/No.",
						},
						"overview_display": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "实例 overview switch, Yes/No.",
						},
					},
				},
			},

			"regions": {
				Optional:    true,
				ForceNew:    true,
				Default:     "All",
				Type:        schema.TypeString,
				Description: "Effective 实例 地域, 值 是 All, 其中 表示 all regions.",
			},

			"product": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: mysql - 云 数据库 MySQL, cynosdb - 云 数据库 CynosDB 对于 MySQL.",
			},

			"instance_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Specifies ID 的 实例 whose inspection 状态 是 changed.",
			},
		},
	}
}

func resourceTencentCloudDbbrainModifyDiagDbInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_modify_diag_db_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request     = dbbrain.NewModifyDiagDBInstanceConfRequest()
		operationId string
	)

	instanceConfs := dbbrain.InstanceConfs{}
	if dMap, ok := helper.InterfacesHeadMap(d, "instance_confs"); ok {
		if v, ok := dMap["daily_inspection"]; ok {
			instanceConfs.DailyInspection = helper.String(v.(string))
		}
		if v, ok := dMap["overview_display"]; ok {
			instanceConfs.OverviewDisplay = helper.String(v.(string))
		}
		request.InstanceConfs = &instanceConfs
	}

	if v, ok := d.GetOk("regions"); ok {
		request.Regions = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbbrainClient().ModifyDiagDBInstanceConf(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dbbrain modifyDiagDbInstanceConf failed, reason:%+v", logId, err)
		return err
	}

	operationId = helper.ResourceIdsHash([]string{*instanceConfs.DailyInspection, *instanceConfs.OverviewDisplay})
	d.SetId(operationId)

	return resourceTencentCloudDbbrainModifyDiagDbInstanceOperationRead(d, meta)
}

func resourceTencentCloudDbbrainModifyDiagDbInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_modify_diag_db_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDbbrainModifyDiagDbInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_modify_diag_db_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
