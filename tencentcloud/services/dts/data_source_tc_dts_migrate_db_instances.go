package dts

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDtsMigrateDbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDtsMigrateDbInstancesRead,
		Schema: map[string]*schema.Schema{
			"database_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database 类型",
			},

			"migrate_role": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "是否instance is the migration 来源 or destination,src(for 来源)，dst(for destination)。",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database 实例 ID",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database 实例名称",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "限制",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "偏移量",
			},

			"account_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The owning 账号 of the resource is null or self(resources in the self 账号)，other(resources in the other 账号)。",
			},

			"tmp_secret_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "temporary secret id，used across 账号",
			},

			"tmp_secret_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "temporary secret 键，used across 账号",
			},

			"tmp_token": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "temporary 令牌，used across 账号",
			},

			"instances": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database 实例名称",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance VIP",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例端口",
						},
						"usable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Can used in migration，1-yes，0-no。",
						},
						"hint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason of can&#39;t used in migration。",
						},
					},
				},
			},

			"request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Unique request id，provide this when encounter a problem。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudDtsMigrateDbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dts_migrate_db_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("database_type"); ok {
		paramMap["DatabaseType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("migrate_role"); ok {
		paramMap["MigrateRole"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		paramMap["InstanceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_mode"); ok {
		paramMap["AccountMode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmp_secret_id"); ok {
		paramMap["TmpSecretId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmp_secret_key"); ok {
		paramMap["TmpSecretKey"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmp_token"); ok {
		paramMap["TmpToken"] = helper.String(v.(string))
	}

	service := DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instances []*dts.MigrateDBItem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDtsMigrateDbInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instances))
	tmpList := make([]map[string]interface{}, 0, len(instances))

	if instances != nil {
		for _, migrateDBItem := range instances {
			migrateDBItemMap := map[string]interface{}{}

			if migrateDBItem.InstanceId != nil {
				migrateDBItemMap["instance_id"] = migrateDBItem.InstanceId
			}

			if migrateDBItem.InstanceName != nil {
				migrateDBItemMap["instance_name"] = migrateDBItem.InstanceName
			}

			if migrateDBItem.Vip != nil {
				migrateDBItemMap["vip"] = migrateDBItem.Vip
			}

			if migrateDBItem.Vport != nil {
				migrateDBItemMap["vport"] = migrateDBItem.Vport
			}

			if migrateDBItem.Usable != nil {
				migrateDBItemMap["usable"] = migrateDBItem.Usable
			}

			if migrateDBItem.Hint != nil {
				migrateDBItemMap["hint"] = migrateDBItem.Hint
			}

			ids = append(ids, *migrateDBItem.InstanceId)
			tmpList = append(tmpList, migrateDBItemMap)
		}

		_ = d.Set("instances", tmpList)
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
