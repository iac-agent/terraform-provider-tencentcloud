package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverMigration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverMigrationCreate,
		Read:   resourceTencentCloudSqlserverMigrationRead,
		Update: resourceTencentCloudSqlserverMigrationUpdate,
		Delete: resourceTencentCloudSqlserverMigrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"migrate_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name 的 迁移 任务.",
			},

			"migrate_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Migration 类型 (1 structure 迁移 2 数据 迁移 3 incremental synchronization).",
			},

			"source_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Type 的 迁移 source 1 TencentDB 对于 SQLServer 2 Cloud 服务器 self-built SQLServer 数据库 4 SQLServer 备份 和 恢复 5 SQLServer 备份 和 恢复 (COS 模式).",
			},

			"source": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Migration source.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID 的 迁移 source 实例, 其中 是 使用 当 MigrateType=1 (TencentDB 对于 SQLServers). 格式 是 mssql-si2823jyl.",
						},
						"cvm_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID 的 迁移 source Cvm, 使用 当 MigrateType=2 (云 服务器 self-built SQL Server 数据库).",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Vpc 网络 ID 的 迁移 source Cvm 是 使用 当 MigrateType=2 (云 服务器 self-built SQL Server 数据库). 格式 是 作为 follows vpc-6ys9ont9.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "子网 ID under Vpc 的 source Cvm 是 使用 当 MigrateType=2 (ECS self-built SQL Server 数据库). 格式 是 作为 follows 子网-h9extioi.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User 名称, MigrateType=1 或 MigrateType=2.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Password, MigrateType=1 或 MigrateType=2.",
						},
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Migrate intranet IP 的 self-built 数据库 的 source Cvm, 和 使用 它 当 MigrateType=2 (self-built SQL Server 数据库 的 云 服务器).",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "端口 数量 的 self-built 数据库 的 迁移 source Cvm, 其中 是 使用 当 MigrateType=2 (self-built SQL Server 数据库 的 云 服务器).",
						},
						"url": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "source 备份 地址 对于 offline 迁移. MigrateType=4 或 MigrateType=5.",
						},
						"url_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "source 备份 密码 对于 offline 迁移, MigrateType=4 或 MigrateType=5.",
						},
					},
				},
			},

			"target": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Migration 目标.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID 的 迁移 目标 实例, 在 格式 mssql-si2823jyl.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User 名称 的 迁移 目标 实例.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Password 的 迁移 目标 实例.",
						},
					},
				},
			},

			"migrate_db_set": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Migrate DB objects. Offline 迁移 是 不 使用 (SourceType=4 或 SourceType=5).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name 的 迁移 数据库.",
						},
					},
				},
			},

			"rename_restore": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Restore 和 rename 数据库 在 ReNameRestoreDatabase. 如果 它 是 不 filled 在, restored 数据库 将 是 named 通过 默认值 和 all databases 将 是 restored. 有效 如果 SourceType=5.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"old_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 的 库. 如果 oldName does 不 exist, failure 是 返回.It 可以 是 left blank 当 使用 对于 offline 迁移 tasks.",
						},
						"new_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "当 new 名称 的 库 是 使用 对于 offline 迁移, 如果 它 是 不 filled 在, 它 将 是 named according 到 OldName. OldName 和 NewName 不能 是 filled 在 在 same 时间. OldName 和 NewName 必须 是 filled 在 和 不能 是 duplicate 当 使用 对于 cloning 数据库.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverMigrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_migration.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = sqlserver.NewCreateMigrationRequest()
		response  = sqlserver.NewCreateMigrationResponse()
		migrateId string
	)
	if v, ok := d.GetOk("migrate_name"); ok {
		request.MigrateName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("migrate_type"); v != nil {
		request.MigrateType = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("source_type"); v != nil {
		request.SourceType = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "source"); ok {
		migrateSource := sqlserver.MigrateSource{}
		if v, ok := dMap["instance_id"]; ok {
			migrateSource.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["cvm_id"]; ok {
			migrateSource.CvmId = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_id"]; ok {
			migrateSource.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			migrateSource.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			migrateSource.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			migrateSource.Password = helper.String(v.(string))
		}
		if v, ok := dMap["ip"]; ok {
			migrateSource.Ip = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			migrateSource.Port = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["url"]; ok {
			urlSet := v.(*schema.Set).List()
			for i := range urlSet {
				url := urlSet[i].(string)
				migrateSource.Url = append(migrateSource.Url, &url)
			}
		}
		if v, ok := dMap["url_password"]; ok {
			migrateSource.UrlPassword = helper.String(v.(string))
		}
		request.Source = &migrateSource
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "target"); ok {
		migrateTarget := sqlserver.MigrateTarget{}
		if v, ok := dMap["instance_id"]; ok {
			migrateTarget.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			migrateTarget.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			migrateTarget.Password = helper.String(v.(string))
		}
		request.Target = &migrateTarget
	}

	if v, ok := d.GetOk("migrate_db_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			migrateDB := sqlserver.MigrateDB{}
			if v, ok := dMap["db_name"]; ok {
				migrateDB.DBName = helper.String(v.(string))
			}
			request.MigrateDBSet = append(request.MigrateDBSet, &migrateDB)
		}
	}

	if v, ok := d.GetOk("rename_restore"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			renameRestoreDatabase := sqlserver.RenameRestoreDatabase{}
			if v, ok := dMap["old_name"]; ok {
				renameRestoreDatabase.OldName = helper.String(v.(string))
			}
			if v, ok := dMap["new_name"]; ok {
				renameRestoreDatabase.NewName = helper.String(v.(string))
			}
			request.RenameRestore = append(request.RenameRestore, &renameRestoreDatabase)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().CreateMigration(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver migration failed, reason:%+v", logId, err)
		return err
	}

	migrateId = helper.Int64ToStr(*response.Response.MigrateId)
	d.SetId(migrateId)

	return resourceTencentCloudSqlserverMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverMigrationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_migration.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	migrationId := d.Id()

	migration, err := service.DescribeSqlserverMigrationById(ctx, migrationId)
	if err != nil {
		return err
	}

	if migration == nil {
		d.SetId("")
		return fmt.Errorf("resource `SqlserverMigration` %s does not exist", d.Id())
	}

	if migration.MigrateName != nil {
		_ = d.Set("migrate_name", migration.MigrateName)
	}

	if migration.MigrateType != nil {
		_ = d.Set("migrate_type", migration.MigrateType)
	}

	if migration.SourceType != nil {
		_ = d.Set("source_type", migration.SourceType)
	}

	if migration.Source != nil {
		sourceMap := map[string]interface{}{}

		if migration.Source.InstanceId != nil {
			sourceMap["instance_id"] = migration.Source.InstanceId
		}

		if migration.Source.CvmId != nil {
			sourceMap["cvm_id"] = migration.Source.CvmId
		}

		if migration.Source.VpcId != nil {
			sourceMap["vpc_id"] = migration.Source.VpcId
		}

		if migration.Source.SubnetId != nil {
			sourceMap["subnet_id"] = migration.Source.SubnetId
		}

		if migration.Source.UserName != nil {
			sourceMap["user_name"] = migration.Source.UserName
		}

		if migration.Source.Password != nil {
			sourceMap["password"] = migration.Source.Password
		}

		if migration.Source.Ip != nil {
			sourceMap["ip"] = migration.Source.Ip
		}

		if migration.Source.Port != nil {
			sourceMap["port"] = migration.Source.Port
		}

		if migration.Source.Url != nil {
			sourceMap["url"] = migration.Source.Url
		}

		if migration.Source.UrlPassword != nil {
			sourceMap["url_password"] = migration.Source.UrlPassword
		}

		_ = d.Set("source", []interface{}{sourceMap})
	}

	if migration.Target != nil {
		targetMap := map[string]interface{}{}

		if migration.Target.InstanceId != nil {
			targetMap["instance_id"] = migration.Target.InstanceId
		}

		if migration.Target.UserName != nil {
			targetMap["user_name"] = migration.Target.UserName
		}

		if migration.Target.Password != nil {
			targetMap["password"] = migration.Target.Password
		}

		_ = d.Set("target", []interface{}{targetMap})
	}

	if migration.MigrateDBSet != nil {
		migrateDBSetList := []interface{}{}
		for _, migrateDB := range migration.MigrateDBSet {
			migrateDBSetMap := map[string]interface{}{}

			if migrateDB.DBName != nil {
				migrateDBSetMap["db_name"] = migrateDB.DBName
			}

			migrateDBSetList = append(migrateDBSetList, migrateDBSetMap)
		}

		_ = d.Set("migrate_db_set", migrateDBSetList)

	}

	// omit rename_restore because read api doesn't return it
	// _ = d.Set("rename_restore", renameRestoreList)

	return nil
}

func resourceTencentCloudSqlserverMigrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_migration.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := sqlserver.NewModifyMigrationRequest()
	migrateId := d.Id()

	request.MigrateId = helper.StrToUint64Point(migrateId)
	if d.HasChange("rename_restore") {
		o, _ := d.GetChange("rename_restore")
		_ = d.Set("rename_restore", o)
		return fmt.Errorf("argument `%s` cannot be changed", d.Id())
	}
	if d.HasChange("migrate_name") {
		if v, ok := d.GetOk("migrate_name"); ok {
			request.MigrateName = helper.String(v.(string))
		}
	}

	if d.HasChange("migrate_type") {
		if v, _ := d.GetOk("migrate_type"); v != nil {
			request.MigrateType = helper.IntUint64(v.(int))
		}
		if v, _ := d.GetOk("source_type"); v != nil {
			request.SourceType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("source_type") {
		if v, _ := d.GetOk("source_type"); v != nil {
			request.SourceType = helper.IntUint64(v.(int))
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "source"); ok {
		migrateSource := sqlserver.MigrateSource{}
		if v, ok := dMap["instance_id"]; ok {
			migrateSource.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["cvm_id"]; ok {
			migrateSource.CvmId = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_id"]; ok {
			migrateSource.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			migrateSource.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			migrateSource.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			migrateSource.Password = helper.String(v.(string))
		}
		if v, ok := dMap["ip"]; ok {
			migrateSource.Ip = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			migrateSource.Port = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["url"]; ok {
			urlSet := v.(*schema.Set).List()
			for i := range urlSet {
				url := urlSet[i].(string)
				migrateSource.Url = append(migrateSource.Url, &url)
			}
		}
		if v, ok := dMap["url_password"]; ok {
			migrateSource.UrlPassword = helper.String(v.(string))
		}
		request.Source = &migrateSource
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "target"); ok {
		migrateTarget := sqlserver.MigrateTarget{}
		if v, ok := dMap["instance_id"]; ok {
			migrateTarget.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			migrateTarget.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			migrateTarget.Password = helper.String(v.(string))
		}
		request.Target = &migrateTarget
	}

	if v, ok := d.GetOk("migrate_db_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			migrateDB := sqlserver.MigrateDB{}
			if v, ok := dMap["db_name"]; ok {
				migrateDB.DBName = helper.String(v.(string))
			}
			request.MigrateDBSet = append(request.MigrateDBSet, &migrateDB)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyMigration(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver migration failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverMigrationRead(d, meta)
}

func resourceTencentCloudSqlserverMigrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_migration.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	migrateId := d.Id()

	if err := service.DeleteSqlserverMigrationById(ctx, migrateId); err != nil {
		return err
	}

	return nil
}
