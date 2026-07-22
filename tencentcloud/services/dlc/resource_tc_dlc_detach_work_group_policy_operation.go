package dlc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcDetachWorkGroupPolicyOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcDetachWorkGroupPolicyOperationCreate,
		Read:   resourceTencentCloudDlcDetachWorkGroupPolicyOperationRead,
		Delete: resourceTencentCloudDlcDetachWorkGroupPolicyOperationDelete,
		Schema: map[string]*schema.Schema{
			"work_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Work 组 ID",
			},

			"policy_set": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Collection 的 policies 到 是 bound。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 目标 数据库. `*` 表示 all databases 在 当前 catalog. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 数据库。",
						},
						"catalog": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 目标 数据 来源 To grant admin 权限，它 必须 是 `*` (all resources 在 此 级别); 到 grant 数据 来源 和 数据库 permissions，它 必须 是 `COSDataCatalog` 或 `*`; 到 grant 表 permissions，它 可以 是 自定义 数据 来源; 如果 它 是 left 空，`DataLakeCatalog` 是 使用. 注意: To grant permissions 在 自定义 数据 来源， permissions 该 可以 是 managed 在 Data Lake Compute console 是 subsets 的 账号 permissions granted 当 您 connect 数据 来源 到 console。",
						},
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 目标 表. `*` 表示 all tables 在 当前 数据库. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 和 数据库 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 表。",
						},
						"operation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "目标 permissions，其中 vary 通过 权限 级别 Admin: `ALL` (默认值); 数据 连接: `CREATE`; 数据库: `ALL`，`CREATE`，`ALTER`，和 `DROP`; 表: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，和 `UPDATE`. 注意: For 表 permissions，如果 数据 来源 other 比 `COSDataCatalog` 是 指定，仅 `SELECT` 权限 可以 是 granted here。",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "权限 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，和 `ENGINE`. 注意: 如果 它 是 left 空，`ADMIN` 是 使用。",
						},
						"function": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 目标 函数. `*` 表示 all functions 在 当前 catalog. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 函数。",
						},
						"view": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 目标 view. `*` 表示 all views 在 当前 数据库. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 和 数据库 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any view。",
						},
						"column": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 目标 列. `*` 表示 all columns. To grant admin permissions，它 必须 是 `*`。",
						},
						"data_engine": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 目标 数据 引擎. `*` 表示 all engines. To grant admin permissions，它 必须 是 `*`。",
						},
						"re_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否grantee 是 allowed 到 further grant permissions. 有效值：`false` (默认值) 和 `true` ( grantee 可以 grant permissions gained here 到 other sub-users)。",
						},
						"source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "权限 来源，其中 不是必填项 当 input 参数 是 passed 在. 有效值：`USER` (从 用户) 和 `WORKGROUP` (从 一个 或 more associated work groups)。",
						},
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "grant 模式，其中 不是必填项 作为 input 参数. 有效值：`COMMON` 和 `SENIOR`。",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "操作者，其中 不是必填项 作为 input 参数。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "权限 策略 创建时间，其中 不是必填项 作为 input 参数。",
						},
						"source_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ID work 组，其中 applies 仅 当 值 的 `来源` 字段 是 `WORKGROUP`。",
						},
						"source_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 work 组，其中 applies 仅 当 值 的 `来源` 字段 是 `WORKGROUP`。",
						},
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "策略 ID。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcDetachWorkGroupPolicyOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_work_group_policy_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = dlc.NewDetachWorkGroupPolicyRequest()
		workGroupId string
	)

	if v, _ := d.GetOkExists("work_group_id"); v != nil {
		workGroupId = helper.IntToStr(v.(int))
		request.WorkGroupId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("policy_set"); ok {
		for _, item := range v.([]interface{}) {
			policy := dlc.Policy{}
			dMap := item.(map[string]interface{})
			if v, ok := dMap["database"]; ok {
				policy.Database = helper.String(v.(string))
			}

			if v, ok := dMap["catalog"]; ok {
				policy.Catalog = helper.String(v.(string))
			}

			if v, ok := dMap["table"]; ok {
				policy.Table = helper.String(v.(string))
			}

			if v, ok := dMap["operation"]; ok {
				policy.Operation = helper.String(v.(string))
			}

			if v, ok := dMap["policy_type"]; ok {
				policy.PolicyType = helper.String(v.(string))
			}

			if v, ok := dMap["function"]; ok {
				policy.Function = helper.String(v.(string))
			}

			if v, ok := dMap["view"]; ok {
				policy.View = helper.String(v.(string))
			}

			if v, ok := dMap["column"]; ok {
				policy.Column = helper.String(v.(string))
			}

			if v, ok := dMap["data_engine"]; ok {
				policy.DataEngine = helper.String(v.(string))
			}

			if v, ok := dMap["re_auth"]; ok {
				policy.ReAuth = helper.Bool(v.(bool))
			}

			if v, ok := dMap["source"]; ok {
				policy.Source = helper.String(v.(string))
			}

			if v, ok := dMap["mode"]; ok {
				policy.Mode = helper.String(v.(string))
			}

			if v, ok := dMap["operator"]; ok {
				policy.Operator = helper.String(v.(string))
			}

			if v, ok := dMap["create_time"]; ok {
				policy.CreateTime = helper.String(v.(string))
			}

			if v, ok := dMap["source_id"]; ok {
				policy.SourceId = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["source_name"]; ok {
				policy.SourceName = helper.String(v.(string))
			}

			if v, ok := dMap["id"]; ok {
				policy.Id = helper.IntInt64(v.(int))
			}

			request.PolicySet = append(request.PolicySet, &policy)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DetachWorkGroupPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate dlc detachWorkGroupPolicyOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(workGroupId)
	return resourceTencentCloudDlcDetachWorkGroupPolicyOperationRead(d, meta)
}

func resourceTencentCloudDlcDetachWorkGroupPolicyOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_work_group_policy_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcDetachWorkGroupPolicyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_work_group_policy_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
