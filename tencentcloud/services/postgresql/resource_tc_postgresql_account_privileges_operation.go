package postgresql

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlAccountPrivilegesOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlAccountPrivilegesOperationCreate,
		Read:   resourceTencentCloudPostgresqlAccountPrivilegesOperationRead,
		Update: resourceTencentCloudPostgresqlAccountPrivilegesOperationUpdate,
		Delete: resourceTencentCloudPostgresqlAccountPrivilegesOperationDelete,

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例 ID in the 格式 of postgres-4wdeb0zv。",
			},
			"user_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance 用户名",
			},
			"modify_privilege_set": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Privileges to modify. Batch modification supported，up to 50 entries at a time。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_privilege": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Database objects and the 用户 permissions on these objects. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Database object.If ObjectType is database，DatabaseName/SchemaName/TableName can be null.If ObjectType is schema，SchemaName/TableName can be null.If ObjectType is table，TableName can be null.If ObjectType is column，DatabaseName/SchemaName/TableName can&amp;#39;t be null.In all other cases，DatabaseName/SchemaName/TableName can be null. 注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"object_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Supported database object types: 账号，database，schema，sequence，procedure，类型，function，table，view，matview，column. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"object_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Database object 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"database_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database 名称 to which the database object belongs. This parameter is mandatory when ObjectType is not database. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema 名称 to which the database object belongs. This parameter is mandatory when ObjectType is not database or schema. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"table_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Table 名称 to which the database object belongs. This parameter is mandatory when ObjectType is column. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"privilege_set": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "Privileges the specific 账号 has on database object. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"modify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Supported modification method: grantObject，revokeObject，alterRole. grantObject represents granting permissions on object，revokeObject represents revoking permissions on object，and alterRole represents modifying the 账号 类型",
						},
						"is_cascade": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "必填 only when ModifyType is revokeObject. When the parameter is true，revoking permissions will cascade. The 默认值为 false。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPostgresqlAccountPrivilegesOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_account_privileges_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = postgresql.NewModifyAccountPrivilegesRequest()
		dBInstanceId string
		userName     string
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = &dBInstanceId
		dBInstanceId = v.(string)
	}

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = &userName
		userName = v.(string)
	}

	if v, ok := d.GetOk("modify_privilege_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			modifyPrivilege := postgresql.ModifyPrivilege{}
			if databasePrivilegeMap, ok := helper.InterfaceToMap(dMap, "database_privilege"); ok {
				databasePrivilege := postgresql.DatabasePrivilege{}
				if objectMap, ok := helper.InterfaceToMap(databasePrivilegeMap, "object"); ok {
					databaseObject := postgresql.DatabaseObject{}
					if v, ok := objectMap["object_type"]; ok {
						databaseObject.ObjectType = helper.String(v.(string))
					}

					if v, ok := objectMap["object_name"]; ok {
						databaseObject.ObjectName = helper.String(v.(string))
					}

					if v, ok := objectMap["database_name"]; ok {
						databaseObject.DatabaseName = helper.String(v.(string))
					}

					if v, ok := objectMap["schema_name"]; ok {
						databaseObject.SchemaName = helper.String(v.(string))
					}

					if v, ok := objectMap["table_name"]; ok {
						databaseObject.TableName = helper.String(v.(string))
					}

					databasePrivilege.Object = &databaseObject
				}

				if v, ok := databasePrivilegeMap["privilege_set"]; ok {
					privilegeSetSet := v.(*schema.Set).List()
					for i := range privilegeSetSet {
						if privilegeSetSet[i] != nil {
							privilegeSet := privilegeSetSet[i].(string)
							databasePrivilege.PrivilegeSet = append(databasePrivilege.PrivilegeSet, &privilegeSet)
						}
					}
				}

				modifyPrivilege.DatabasePrivilege = &databasePrivilege
			}

			if v, ok := dMap["modify_type"]; ok {
				modifyPrivilege.ModifyType = helper.String(v.(string))
			}

			if v, ok := dMap["is_cascade"]; ok {
				modifyPrivilege.IsCascade = helper.Bool(v.(bool))
			}

			request.ModifyPrivilegeSet = append(request.ModifyPrivilegeSet, &modifyPrivilege)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyAccountPrivileges(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update postgres accountPrivileges failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{dBInstanceId, userName}, tccommon.FILED_SP))

	return resourceTencentCloudPostgresqlAccountPrivilegesOperationUpdate(d, meta)
}

func resourceTencentCloudPostgresqlAccountPrivilegesOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_account_privileges_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlAccountPrivilegesOperationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_account_privileges_operation.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudPostgresqlAccountPrivilegesOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlAccountPrivilegesOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_account_privileges_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
