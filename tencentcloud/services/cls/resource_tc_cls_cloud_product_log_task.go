package cls

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsCloudProductLogTask() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.81.188. Please use `tencentcloud_cls_cloud_product_log_task_v2` instead.",
		Create:             resourceTencentCloudClsCloudProductLogTaskCreate,
		Read:               resourceTencentCloudClsCloudProductLogTaskRead,
		Update:             resourceTencentCloudClsCloudProductLogTaskUpdate,
		Delete:             resourceTencentCloudClsCloudProductLogTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例ID。",
			},

			"assumer_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "云产品标识，值：CDS、CWP、CDB、TDSQL-C、MongoDB、TDStore、DCDB、MariaDB、PostgreSQL、BH、APIS。",
			},

			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "日志类型，值：CDS-AUDIT、CDS-RISK、CDB-AUDIT、TDSQL-C-AUDIT、MongoDB-AUDIT、MongoDB-SlowLog、MongoDB-ErrorLog、TDMYSQL-SLOW、DCDB-AUDIT、DCDB-SLOW、DCDB-ERROR、MariaDB-AUDIT、MariaDB-SLOW、MariaDB-ERROR、PostgreSQL-SLOW、PostgreSQL-ERROR、 PostgreSQL-审计、BH-文件日志、BH-命令日志、APIS-访问。",
			},

			"cloud_product_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "云产品区域。不同地区不同日志类型的输入格式存在差异。请参考以下示例：\n- CDS(所有日志类型): ap-guangzhou\n- CDB-AUDIT: gz\n- TDSQL-C-AUDIT: gz\n- MongoDB-AUDIT: gz\n- MongoDB-SlowLog: ap-guangzhou\n- MongoDB-ErrorLog: ap-guangzhou\n- TDMYSQL-SLOW: gz\n- DCDB(所有日志type): gz\n- MariaDB(所有日志类型): gz\n- PostgreSQL(所有日志类型): gz\n- BH(所有日志类型): Overseas-polaris(海外国内站点)/fsi-polaris(国内站点金融)/general-polaris(国内站点)/intl-sg-prod(国际站点)\n- APIS(所有日志类型): gz。",
			},

			"cls_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CLS 目标区域。",
			},

			"logset_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "日志集名称，会自动创建。",
			},

			"logset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "日志集 ID。",
			},

			"topic_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "日志主题的名称，系统会自动创建。",
			},

			"topic_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "记录主题 ID。",
			},

			"extend": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "日志配置扩展信息，一般用于存储额外的日志传递配置。",
			},
		},
	}
}

func resourceTencentCloudClsCloudProductLogTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cloud_product_log_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId         string
		assumerName        string
		logType            string
		cloudProductRegion string
	)
	var (
		request  = clsv20201016.NewCreateCloudProductLogCollectionRequest()
		response = clsv20201016.NewCreateCloudProductLogCollectionResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("assumer_name"); ok {
		assumerName = v.(string)
	}
	if v, ok := d.GetOk("log_type"); ok {
		logType = v.(string)
	}
	if v, ok := d.GetOk("cloud_product_region"); ok {
		cloudProductRegion = v.(string)
	}

	request.InstanceId = helper.String(instanceId)

	request.AssumerName = helper.String(assumerName)

	request.LogType = helper.String(logType)

	request.CloudProductRegion = helper.String(cloudProductRegion)

	if v, ok := d.GetOk("cls_region"); ok {
		request.ClsRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("extend"); ok {
		request.Extend = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().CreateCloudProductLogCollectionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls cloud product log task failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 10*tccommon.ReadRetryTimeout, time.Second, service.ClsCloudProductLogTaskStateRefreshFunc(ctx, instanceId, assumerName, logType, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(strings.Join([]string{instanceId, assumerName, logType, cloudProductRegion}, tccommon.FILED_SP))

	return resourceTencentCloudClsCloudProductLogTaskRead(d, meta)
}

func resourceTencentCloudClsCloudProductLogTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cloud_product_log_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	assumerName := idSplit[1]
	logType := idSplit[2]
	cloudProductRegion := idSplit[3]

	_ = d.Set("instance_id", instanceId)

	_ = d.Set("assumer_name", assumerName)

	_ = d.Set("log_type", logType)

	_ = d.Set("cloud_product_region", cloudProductRegion)

	respData, err := service.DescribeClsCloudProductLogTaskById(ctx, instanceId, assumerName, logType)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `cls_cloud_product_log_task` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if len(respData.Tasks) > 0 {
		_ = d.Set("topic_id", respData.Tasks[0].TopicId)
		_ = d.Set("logset_id", respData.Tasks[0].LogsetId)
		_ = d.Set("extend", respData.Tasks[0].Extend)
		_ = d.Set("cls_region", respData.Tasks[0].ClsRegion)
	}
	_ = instanceId
	_ = assumerName
	_ = logType
	_ = cloudProductRegion
	return nil
}

func resourceTencentCloudClsCloudProductLogTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cloud_product_log_task.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"instance_id", "assumer_name", "log_type", "cloud_product_region", "cls_region", "logset_name", "topic_name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	assumerName := idSplit[1]
	logType := idSplit[2]
	cloudProductRegion := idSplit[3]

	needChange := false
	mutableArgs := []string{"extend"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := clsv20201016.NewModifyCloudProductLogCollectionRequest()

		request.InstanceId = helper.String(instanceId)

		request.AssumerName = helper.String(assumerName)

		request.LogType = helper.String(logType)

		request.CloudProductRegion = helper.String(cloudProductRegion)

		if v, ok := d.GetOk("extend"); ok {
			request.Extend = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().ModifyCloudProductLogCollectionWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cls cloud product log task failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClsCloudProductLogTaskRead(d, meta)
}

func resourceTencentCloudClsCloudProductLogTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cloud_product_log_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	assumerName := idSplit[1]
	logType := idSplit[2]
	cloudProductRegion := idSplit[3]

	var (
		request  = clsv20201016.NewDeleteCloudProductLogCollectionRequest()
		response = clsv20201016.NewDeleteCloudProductLogCollectionResponse()
	)

	request.InstanceId = helper.String(instanceId)

	request.AssumerName = helper.String(assumerName)

	request.LogType = helper.String(logType)

	request.CloudProductRegion = helper.String(cloudProductRegion)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteCloudProductLogCollectionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cls cloud product log task failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"3"}, 10*tccommon.ReadRetryTimeout, time.Second, service.ClsCloudProductLogTaskStateRefreshFunc(ctx, instanceId, assumerName, logType, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	var (
		request1  = clsv20201016.NewDeleteTopicRequest()
		response1 = clsv20201016.NewDeleteTopicResponse()
	)

	if v, ok := d.GetOk("topic_id"); ok {
		request1.TopicId = helper.String(v.(string))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteTopicWithContext(ctx, request1)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
		}
		response1 = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cls cloud product log task failed, reason:%+v", logId, err)
		return err
	}

	_ = response1
	var (
		request2  = clsv20201016.NewDeleteLogsetRequest()
		response2 = clsv20201016.NewDeleteLogsetResponse()
	)

	if v, ok := d.GetOk("logset_id"); ok {
		request2.LogsetId = helper.String(v.(string))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsV20201016Client().DeleteLogsetWithContext(ctx, request2)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request2.GetAction(), request2.ToJsonString(), result.ToJsonString())
		}
		response2 = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete cls cloud product log task failed, reason:%+v", logId, err)
		return err
	}

	_ = response2
	return nil
}
