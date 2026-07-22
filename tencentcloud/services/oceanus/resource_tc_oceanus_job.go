package oceanus

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOceanusJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusJobCreate,
		Read:   resourceTencentCloudOceanusJobRead,
		Update: resourceTencentCloudOceanusJobUpdate,
		Delete: resourceTencentCloudOceanusJobDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "名称 作业. It 可以 是 composed 的 Chinese，English，numbers，hyphens (-)，underscores (_)，和 periods (.)，和 长度 不能 exceed 50 字符. 注意 该 作业名称 不能 是 same 作为 existing 作业。",
			},
			"job_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(JOB_TYPE),
				Description:  "类型 作业. 1 表示SQL 作业，和 2 表示JAR 作业。",
			},
			"cluster_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CLUSTER_TYPE),
				Description:  "类型 集群. 1 表示shared 集群，和 2 表示exclusive 集群。",
			},
			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "当 ClusterType=2，它 为必填项 到 指定ID 的 exclusive 集群 到 其中 作业 是 submitted。",
			},
			"cu_mem": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      CU_MEM_4,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CU_MEM),
				Description:  "Set 内存 规格 的 each CU，（GB）。 It 支持 2，4，8，和 16 (其中 needs 到 apply 对于 whitelist before 使用). 默认为 4，该 是，1 CU corresponds 到 4 GB 的 running 内存。",
			},
			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "备注 信息 的 作业. It 可以 是 集合 arbitrarily。",
			},
			"folder_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "文件夹 ID 到 其中 作业名称 belongs. root directory 是 root。",
			},
			"flink_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Flink 版本 该 作业 runs。",
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "workspace SerialId。",
			},
		},
	}
}

func resourceTencentCloudOceanusJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = oceanus.NewCreateJobRequest()
		response = oceanus.NewCreateJobResponse()
		jobId    string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("job_type"); ok {
		request.JobType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("cluster_type"); ok {
		request.ClusterType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cu_mem"); ok {
		request.CuMem = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("folder_id"); ok {
		request.FolderId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("flink_version"); ok {
		request.FlinkVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().CreateJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("oceanus Job not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus Job failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	return resourceTencentCloudOceanusJobRead(d, meta)
}

func resourceTencentCloudOceanusJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		jobId   = d.Id()
	)

	Job, err := service.DescribeOceanusJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if Job == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if Job.Name != nil {
		_ = d.Set("name", Job.Name)
	}

	if Job.JobType != nil {
		_ = d.Set("job_type", Job.JobType)
	}

	if Job.ClusterId != nil {
		_ = d.Set("cluster_id", Job.ClusterId)
	}

	if Job.CuMem != nil {
		_ = d.Set("cu_mem", Job.CuMem)
	}

	if Job.Remark != nil {
		_ = d.Set("remark", Job.Remark)
	}

	if Job.FlinkVersion != nil {
		_ = d.Set("flink_version", Job.FlinkVersion)
	}

	if Job.WorkSpaceId != nil {
		_ = d.Set("work_space_id", Job.WorkSpaceId)
	}

	return nil
}

func resourceTencentCloudOceanusJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_job.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = oceanus.NewModifyJobRequest()
		jobId   = d.Id()
	)

	immutableArgs := []string{"job_type", "cluster_type", "cluster_id", "cu_mem", "folder_id", "flink_version"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.JobId = &jobId

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	if d.HasChange("work_space_id") {
		if v, ok := d.GetOk("work_space_id"); ok {
			request.WorkSpaceId = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().ModifyJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update oceanus Job failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudOceanusJobRead(d, meta)
}

func resourceTencentCloudOceanusJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		jobId   = d.Id()
	)

	if err := service.DeleteOceanusJobById(ctx, jobId); err != nil {
		return err
	}

	return nil
}
