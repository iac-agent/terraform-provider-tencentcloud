package dlc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcDetachUserPolicyOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcDetachUserPolicyOperationCreate,
		Read:   resourceTencentCloudDlcDetachUserPolicyOperationRead,
		Delete: resourceTencentCloudDlcDetachUserPolicyOperationDelete,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "用户 ID，which matches Uin on the CAM side。",
			},

			"policy_set": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Collection of unbound permissions。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The 名称 target database. `*` represents all databases in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any database。",
						},
						"catalog": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The 名称 target data 来源 To grant admin permission，it must be `*` (all resources at this 级别); to grant data 来源 and database permissions，it must be `COSDataCatalog` or `*`; to grant table permissions，it can be a custom data 来源; if it is left empty，`DataLakeCatalog` is used. Note: To grant permissions on a custom data 来源，the permissions that can be managed in the Data Lake Compute console are subsets of the 账号 permissions granted when you connect the data 来源 to the console。",
						},
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The 名称 target table. `*` represents all tables in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any table。",
						},
						"operation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target permissions，which vary by permission 级别 Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`，`CREATE`，`ALTER`，and `DROP`; table: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，and `UPDATE`. Note: For table permissions，if a data 来源 other than `COSDataCatalog` is specified，only the `SELECT` permission can be granted here。",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The permission 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，and `ENGINE`. Note: If it is left empty，`ADMIN` is used。",
						},
						"function": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The 名称 target function. `*` represents all functions in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any function。",
						},
						"view": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The 名称 target view. `*` represents all views in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any view。",
						},
						"column": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The 名称 target column. `*` represents all columns. To grant admin permissions，it must be `*`。",
						},
						"data_engine": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The 名称 target data engine. `*` represents all engines. To grant admin permissions，it must be `*`。",
						},
						"re_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否grantee is allowed to further grant the permissions. 有效值：`false` (default) and `true` (the grantee can grant permissions gained here to other sub-users)。",
						},
						"source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The permission 来源，which 不是必填项 when input parameters are passed in. 有效值：`USER` (from the 用户) and `WORKGROUP` (from one or more associated work groups)。",
						},
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The grant 模式，which 不是必填项 as an input parameter. 有效值：`COMMON` and `SENIOR`。",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The 操作者，which 不是必填项 as an input parameter。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The permission policy 创建时间，which 不是必填项 as an input parameter。",
						},
						"source_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID work group，which applies only when the 值 of the `来源` field is `WORKGROUP`。",
						},
						"source_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The 名称 work group，which applies only when the 值 of the `来源` field is `WORKGROUP`。",
						},
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The policy ID。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcDetachUserPolicyOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_user_policy_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = dlc.NewDetachUserPolicyRequest()
		userId  string
	)

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		request.UserId = helper.String(v.(string))
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DetachUserPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate dlc detachUserPolicyOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userId)
	return resourceTencentCloudDlcDetachUserPolicyOperationRead(d, meta)
}

func resourceTencentCloudDlcDetachUserPolicyOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_user_policy_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcDetachUserPolicyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_detach_user_policy_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
