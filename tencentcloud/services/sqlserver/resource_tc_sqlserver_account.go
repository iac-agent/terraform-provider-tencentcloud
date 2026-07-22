package sqlserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudSqlserverAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverAccountCreate,
		Read:   resourceTencentCloudSqlserverAccountRead,
		Update: resourceTencentCloudSqlserverAccountUpdate,
		Delete: resourceTencentCLoudSqlserverAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name 的 SQL Server account.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password 的 SQL Server account.",
			},
			"is_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate 该 account 是 root account 或 不.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Remark 的 SQL Server account.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "实例 ID 该 account belongs 到.",
			},
			//computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status 的 SQL Server account. 有效 值: 1, 2, 3, 4. 1 对于 creating, 2 对于 running, 3 对于 modifying, 4 对于 resetting 密码, -1 对于 deleting.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create 时间 的 SQL Server account.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last 更新 时间 的 SQL Server account.",
			},
		},
	}
}

func resourceTencentCloudSqlserverAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_account.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		name       = d.Get("name").(string)
		password   = d.Get("password").(string)
		remark     = d.Get("remark").(string)
		isAdmin    = d.Get("is_admin").(bool)
		instanceId = d.Get("instance_id").(string)
	)

	var outErr, inErr error

	var flowId int64
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		flowId, inErr = sqlserverService.CreateSqlserverAccountReturnFlowId(ctx, instanceId, name, password, remark, isAdmin)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskStatus, inErr := sqlserverService.DescribeCloneStatusByFlowId(ctx, flowId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if taskStatus == nil || taskStatus.Status == nil {
			return resource.RetryableError(fmt.Errorf("sqlserver account flow %d status is nil, retrying", flowId))
		}
		if *taskStatus.Status == int64(SQLSERVER_TASK_RUNNING) {
			return resource.RetryableError(fmt.Errorf("sqlserver account flow %d is still running, retrying", flowId))
		}
		if *taskStatus.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(fmt.Errorf("sqlserver account flow %d failed", flowId))
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId + tccommon.FILED_SP + name)

	return resourceTencentCloudSqlserverAccountRead(d, meta)
}

func resourceTencentCloudSqlserverAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_account.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()
	idStrs := strings.Split(id, tccommon.FILED_SP)
	if len(idStrs) != 2 {
		return fmt.Errorf("invalid SQL Server account id %s", id)
	}
	instanceId := idStrs[0]
	name := idStrs[1]
	d.Partial(true)

	var outErr, inErr error

	//update is_admin
	if d.HasChange("is_admin") {
		return fmt.Errorf("is_admin is not allowed to change")
	}

	//update password
	if d.HasChange("password") {
		password := d.Get("password").(string)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ResetSqlserverAccountPassword(ctx, instanceId, name, password)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	//update remark
	if d.HasChange("remark") {
		remark := d.Get("remark").(string)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverAccountRemark(ctx, instanceId, name, remark)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	d.Partial(false)

	return resourceTencentCloudSqlserverAccountRead(d, meta)
}

func resourceTencentCloudSqlserverAccountRead(d *schema.ResourceData, meta interface{}) error {

	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_account.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	idStrs := strings.Split(id, tccommon.FILED_SP)
	if len(idStrs) != 2 {
		return fmt.Errorf("invalid SQL Server account id %s", id)
	}
	instanceId := idStrs[0]
	name := idStrs[1]

	var outErr, inErr error
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	account, has, outErr := sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			account, has, inErr = sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("name", account.Name)
	if *account.Remark == "--" {
		_ = d.Set("remark", "")
	} else {
		_ = d.Set("remark", account.Remark)
	}
	_ = d.Set("create_time", account.CreateTime)
	_ = d.Set("update_time", account.UpdateTime)
	_ = d.Set("status", account.Status)

	return nil
}

func resourceTencentCLoudSqlserverAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_account.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	idStrs := strings.Split(id, tccommon.FILED_SP)
	if len(idStrs) != 2 {
		return fmt.Errorf("invalid SQL Server account id %s", id)
	}
	instanceId := idStrs[0]
	name := idStrs[1]

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var outErr, inErr error
	var has bool

	// Block 1: confirm account exists
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr = sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	if !has {
		return nil
	}

	// Block 2: call delete API only, no wait
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.DeleteSqlserverAccountOnly(ctx, instanceId, name)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	// Block 3: wait for async delete to complete (status -1 → gone)
	outErr = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		instance, has, inErr := sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if !has {
			return nil
		}
		if int(*instance.Status) == -1 {
			return resource.RetryableError(fmt.Errorf("deleting SQL Server account %s, status %d, retrying", id, *instance.Status))
		}
		return resource.NonRetryableError(fmt.Errorf("delete SQL Server account %s failed, unexpected status %d", id, *instance.Status))
	})
	if outErr != nil {
		return outErr
	}

	// Block 4: final verification that account is gone
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr := sqlserverService.DescribeSqlserverAccountById(ctx, instanceId, name)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("delete SQL Server account %s fail, account still exists", id))
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	return nil
}
