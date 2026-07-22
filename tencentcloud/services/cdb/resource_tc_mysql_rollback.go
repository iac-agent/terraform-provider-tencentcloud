package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlRollback() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRollbackCreate,
		Read:   resourceTencentCloudMysqlRollbackRead,
		Delete: resourceTencentCloudMysqlRollbackDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "云数据库实例ID。",
			},

			"strategy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "回滚策略。可用值有：表、db、full；默认值是满的。 表- 极速回滚模式，仅导入所选表级别的备份和binlog，如果存在跨表操作，且未同时选择关联表，则回滚将失败。该模式下，参数Databases必须为空； db- Quick模式，只导入所选库级别的备份和binlog，如果有跨库操作，且没有同时选择关联库，则回滚会失败； full-普通回滚模式，会导入整个实例的备份和binlog，速度较慢。",
			},

			"rollback_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "数据库回滚时间，时间格式为：yyyy-mm-dd hh:mm:ss。",
			},

			"databases": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "待归档的数据库信息，表示整个数据库已归档。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "回滚前的原始数据库名称。",
						},
						"new_database_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "回滚后的新数据库名称。",
						},
					},
				},
			},

			"tables": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "待回滚的数据库表信息，表示该文件是按表回滚的。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "数据库名称。",
						},
						"table": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "数据库表详细信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"table_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "回滚前的原始数据库表名。",
									},
									"new_table_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "回滚后新的数据库表名。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMysqlRollbackCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_rollback.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request               = mysql.NewStartBatchRollbackRequest()
		response              = mysql.NewStartBatchRollbackResponse()
		rollbackInstancesInfo = mysql.RollbackInstancesInfo{}
		instanceId            string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		rollbackInstancesInfo.InstanceId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("strategy"); ok {
		rollbackInstancesInfo.Strategy = helper.String(v.(string))
	}
	if v, ok := d.GetOk("rollback_time"); ok {
		rollbackInstancesInfo.RollbackTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("databases"); ok {
		for _, item := range v.([]interface{}) {
			databasesMap := item.(map[string]interface{})
			rollbackDBName := mysql.RollbackDBName{}
			if v, ok := databasesMap["database_name"]; ok {
				rollbackDBName.DatabaseName = helper.String(v.(string))
			}
			if v, ok := databasesMap["new_database_name"]; ok {
				rollbackDBName.NewDatabaseName = helper.String(v.(string))
			}
			rollbackInstancesInfo.Databases = append(rollbackInstancesInfo.Databases, &rollbackDBName)
		}
	}
	if v, ok := d.GetOk("tables"); ok {
		for _, item := range v.([]interface{}) {
			tablesMap := item.(map[string]interface{})
			rollbackTables := mysql.RollbackTables{}
			if v, ok := tablesMap["database"]; ok {
				rollbackTables.Database = helper.String(v.(string))
			}
			if v, ok := tablesMap["table"]; ok {
				for _, item := range v.([]interface{}) {
					tableMap := item.(map[string]interface{})
					rollbackTableName := mysql.RollbackTableName{}
					if v, ok := tableMap["table_name"]; ok {
						rollbackTableName.TableName = helper.String(v.(string))
					}
					if v, ok := tableMap["new_table_name"]; ok {
						rollbackTableName.NewTableName = helper.String(v.(string))
					}
					rollbackTables.Table = append(rollbackTables.Table, &rollbackTableName)
				}
			}
			rollbackInstancesInfo.Tables = append(rollbackInstancesInfo.Tables, &rollbackTables)
		}
	}
	request.Instances = append(request.Instances, &rollbackInstancesInfo)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().StartBatchRollback(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql rollback failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId := *response.Response.AsyncRequestId
	d.SetId(instanceId + tccommon.FILED_SP + asyncRequestId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s create mysql rollback status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s create mysql rollback status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql rollback fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlRollbackRead(d, meta)
}

func resourceTencentCloudMysqlRollbackRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_rollback.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	asyncRequestId := idSplit[1]

	rollbacks, err := service.DescribeMysqlRollbackById(ctx, instanceId, asyncRequestId)
	if err != nil {
		return err
	}

	if rollbacks == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlRollback` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rollbacks != nil {
		instances := rollbacks[0]
		if instances.InstanceId != nil {
			_ = d.Set("instance_id", instances.InstanceId)
		}

		if instances.Strategy != nil {
			_ = d.Set("strategy", instances.Strategy)
		}

		if instances.RollbackTime != nil {
			_ = d.Set("rollback_time", instances.RollbackTime)
		}

		if instances.Databases != nil {
			databasesList := []interface{}{}
			for _, databases := range instances.Databases {
				databasesMap := map[string]interface{}{}

				if databases.DatabaseName != nil {
					databasesMap["database_name"] = databases.DatabaseName
				}

				if databases.NewDatabaseName != nil {
					databasesMap["new_database_name"] = databases.NewDatabaseName
				}

				databasesList = append(databasesList, databasesMap)
			}

			_ = d.Set("databases", databasesList)
		}

		if instances.Tables != nil {
			tablesList := []interface{}{}
			for _, tables := range instances.Tables {
				tablesMap := map[string]interface{}{}

				if tables.Database != nil {
					tablesMap["database"] = tables.Database
				}

				if tables.Table != nil {
					tableList := []interface{}{}
					for _, table := range tables.Table {
						tableMap := map[string]interface{}{}

						if table.TableName != nil {
							tableMap["table_name"] = table.TableName
						}

						if table.NewTableName != nil {
							tableMap["new_table_name"] = table.NewTableName
						}

						tableList = append(tableList, tableMap)
					}

					tablesMap["table"] = tableList
				}

				tablesList = append(tablesList, tablesMap)
			}

			_ = d.Set("tables", tablesList)
		}
	}

	return nil
}

func resourceTencentCloudMysqlRollbackDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_rollback.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
