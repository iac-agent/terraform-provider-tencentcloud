package oceanus

import (
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOceanusRunJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusRunJobCreate,
		Read:   resourceTencentCloudOceanusRunJobRead,
		Delete: resourceTencentCloudOceanusRunJobDelete,

		Schema: map[string]*schema.Schema{
			"run_job_descriptions": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "The 描述 information for batch job startup。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "作业 ID",
						},
						"run_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The 类型 run. 1 表示start，and 2 表示resume。",
						},
						"start_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Compatible with the startup parameters of the old SQL 类型 job: 指定start time point of data 来源 consumption (recommended to pass the 值)Ensure that the parameter is LATEST，EARLIEST，T+时间戳 (example: T1557394288000)。",
						},
						"job_config_version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "A certain 版本 of the current job(Not passed by default as a non-draft job 版本)。",
						},
						"savepoint_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Savepoint 路径",
						},
						"savepoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Savepoint ID。",
						},
						"use_old_system_connector": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Use the historical 版本 of the system dependency。",
						},
						"custom_timestamp": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Custom 时间戳。",
						},
					},
				},
			},
			"work_space_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId。",
			},
		},
	}
}

func resourceTencentCloudOceanusRunJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_run_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = oceanus.NewRunJobsRequest()
		jobIds  []string
	)

	if v, ok := d.GetOk("run_job_descriptions"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			runJobDescription := oceanus.RunJobDescription{}
			if v, ok := dMap["job_id"]; ok {
				runJobDescription.JobId = helper.String(v.(string))
				jobIds = append(jobIds, v.(string))
			}

			if v, ok := dMap["run_type"]; ok {
				runJobDescription.RunType = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["start_mode"]; ok {
				runJobDescription.StartMode = helper.String(v.(string))
			}

			if v, ok := dMap["job_config_version"]; ok {
				runJobDescription.JobConfigVersion = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["savepoint_path"]; ok {
				runJobDescription.SavepointPath = helper.String(v.(string))
			}

			if v, ok := dMap["savepoint_id"]; ok {
				runJobDescription.SavepointId = helper.String(v.(string))
			}

			if v, ok := dMap["use_old_system_connector"]; ok {
				runJobDescription.UseOldSystemConnector = helper.Bool(v.(bool))
			}

			if v, ok := dMap["custom_timestamp"]; ok {
				runJobDescription.CustomTimestamp = helper.IntInt64(v.(int))
			}

			request.RunJobDescriptions = append(request.RunJobDescriptions, &runJobDescription)
		}
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().RunJobs(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate oceanus runJob failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join(jobIds, tccommon.FILED_SP))

	return resourceTencentCloudOceanusRunJobRead(d, meta)
}

func resourceTencentCloudOceanusRunJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_run_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudOceanusRunJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_run_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
