package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

/*
Use this data source to query Docker activities for a Lighthouse instance.

# Example Usage

# Query Docker activities by instance ID

```hcl

	data "tencentcloud_lighthouse_docker_activitie" "example" {
	  instance_id = "lhins-12345678"
	}

```

# Query Docker activities by instance ID and activity IDs

```hcl

	data "tencentcloud_lighthouse_docker_activitie" "example" {
	  instance_id  = "lhins-12345678"
	  activity_ids = ["lhda-12345678", "lhda-87654321"]
	}

```

# Query Docker activities by time range

```hcl

	data "tencentcloud_lighthouse_docker_activitie" "example" {
	  instance_id        = "lhins-12345678"
	  created_time_begin = 1717200000
	  created_time_end   = 1719800000
	}

```
*/
func DataSourceTencentCloudLighthouseDockerActivitie() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseDockerActivitieRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例 ID Can be obtained from the 实例 ID field returned by the DescribeInstances interface。",
			},

			"activity_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Docker activity ID list. Can be obtained from the ActivityId field returned by the DescribeDockerActivities interface。",
			},

			"created_time_begin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The start 值 of the activity 创建时间，时间戳 （秒）。",
			},

			"created_time_end": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The end 值 of the activity 创建时间，时间戳 （秒）。",
			},

			"docker_activity_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Docker activity list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"activity_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Activity ID。",
						},
						"activity_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Activity 名称",
						},
						"activity_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Activity state. 有效值：INIT，OPERATING，SUCCESS，FAILED。",
						},
						"activity_command_output": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Activity command output，base64 encoded。",
						},
						"container_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Container ID list。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 according to ISO8601 standard. UTC time is used. 格式 is YYYY-MM-DDThh:mm:ssZ。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间 according to ISO8601 standard. UTC time is used. 格式 is YYYY-MM-DDThh:mm:ssZ。",
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

func dataSourceTencentCloudLighthouseDockerActivitieRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_docker_activitie.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})

	// Parse instance_id
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = v.(string)
	}

	// Parse activity_ids
	if v, ok := d.GetOk("activity_ids"); ok {
		activityIdsSet := v.(*schema.Set).List()
		paramMap["ActivityIds"] = helper.InterfacesStringsPoint(activityIdsSet)
	}

	// Parse created_time_begin
	if v, ok := d.GetOk("created_time_begin"); ok {
		paramMap["CreatedTimeBegin"] = int64(v.(int))
	}

	// Parse created_time_end
	if v, ok := d.GetOk("created_time_end"); ok {
		paramMap["CreatedTimeEnd"] = int64(v.(int))
	}

	// Call service layer
	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var dockerActivitySet []*lighthouse.DockerActivity

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseDockerActivitiesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dockerActivitySet = result
		return nil
	})

	if err != nil {
		return err
	}

	// Map response to schema
	ids := make([]string, 0, len(dockerActivitySet))
	tmpList := make([]map[string]interface{}, 0, len(dockerActivitySet))

	if dockerActivitySet != nil {
		for _, activity := range dockerActivitySet {
			activityMap := map[string]interface{}{}

			if activity.ActivityId != nil {
				activityMap["activity_id"] = activity.ActivityId
				ids = append(ids, *activity.ActivityId)
			}

			if activity.ActivityName != nil {
				activityMap["activity_name"] = activity.ActivityName
			}

			if activity.ActivityState != nil {
				activityMap["activity_state"] = activity.ActivityState
			}

			if activity.ActivityCommandOutput != nil {
				activityMap["activity_command_output"] = activity.ActivityCommandOutput
			}

			if activity.ContainerIds != nil {
				activityMap["container_ids"] = activity.ContainerIds
			}

			if activity.CreatedTime != nil {
				activityMap["created_time"] = activity.CreatedTime
			}

			if activity.EndTime != nil {
				activityMap["end_time"] = activity.EndTime
			}

			tmpList = append(tmpList, activityMap)
		}

		_ = d.Set("docker_activity_set", tmpList)
	}

	// Set resource ID
	d.SetId(helper.DataResourceIdsHash(ids))

	// Handle result_output_file
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
