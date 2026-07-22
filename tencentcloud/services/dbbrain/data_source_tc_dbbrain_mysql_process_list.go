package dbbrain

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainMysqlProcessList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainMysqlProcessListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "thread ID, 使用 到 过滤器 thread 列表.",
			},

			"user": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "operating account 名称 的 thread, 使用 到 过滤器 thread 列表.",
			},

			"host": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "operating 主机 地址 的 thread, 使用 到 过滤器 thread 列表.",
			},

			"db": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "threads operations 数据库, 使用 到 过滤器 thread 列表.",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "operational state 的 thread, 使用 到 过滤器 thread 列表.",
			},

			"command": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "execution 类型 的 thread, 使用 到 过滤器 thread 列表.",
			},

			"time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "最小 值 的 operation 时长 的 thread, 在 秒, 使用 到 过滤器 列表 的 threads whose operation 时长 是 longer 比 此 值.",
			},

			"info": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "threads operation statement 是 使用 到 过滤器 thread 列表.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值: `mysql` - 云 数据库 MySQL; `cynosdb` - 云 数据库 TDSQL-C 对于 MySQL, 默认值 是 `mysql`.",
			},

			"process_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Live thread 列表.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "thread ID.",
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operating account 名称 的 thread.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operating 主机 地址 的 thread.",
						},
						"db": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "thread 该 operates 数据库.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operational state 的 thread.",
						},
						"command": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "execution 类型 的 thread.",
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation 时长 的 thread, 在 秒.",
						},
						"info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation statement 对于 thread.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainMysqlProcessListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_mysql_process_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var (
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOkExists("id"); ok {
		paramMap["ID"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("user"); ok {
		paramMap["User"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		paramMap["Host"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db"); ok {
		paramMap["DB"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		paramMap["State"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command"); ok {
		paramMap["Command"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("time"); ok {
		paramMap["Time"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("info"); ok {
		paramMap["Info"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var processList []*dbbrain.MySqlProcess

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainMysqlProcessListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		processList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(processList))
	tmpList := make([]map[string]interface{}, 0, len(processList))

	if processList != nil {
		for _, mySqlProcess := range processList {
			mySqlProcessMap := map[string]interface{}{}

			if mySqlProcess.ID != nil {
				mySqlProcessMap["id"] = mySqlProcess.ID
			}

			if mySqlProcess.User != nil {
				mySqlProcessMap["user"] = mySqlProcess.User
			}

			if mySqlProcess.Host != nil {
				mySqlProcessMap["host"] = mySqlProcess.Host
			}

			if mySqlProcess.DB != nil {
				mySqlProcessMap["db"] = mySqlProcess.DB
			}

			if mySqlProcess.State != nil {
				mySqlProcessMap["state"] = mySqlProcess.State
			}

			if mySqlProcess.Command != nil {
				mySqlProcessMap["command"] = mySqlProcess.Command
			}

			if mySqlProcess.Time != nil {
				mySqlProcessMap["time"] = mySqlProcess.Time
			}

			if mySqlProcess.Info != nil {
				mySqlProcessMap["info"] = mySqlProcess.Info
			}

			ids = append(ids, strings.Join([]string{instanceId, *mySqlProcess.ID}, tccommon.FILED_SP))
			tmpList = append(tmpList, mySqlProcessMap)
		}

		_ = d.Set("process_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
