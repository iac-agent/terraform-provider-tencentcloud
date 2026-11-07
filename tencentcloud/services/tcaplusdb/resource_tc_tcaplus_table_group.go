package tcaplusdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcaplusdbv20190823 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudTcaplusTableGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusTableGroupCreate,
		Read:   resourceTencentCloudTcaplusTableGroupRead,
		Update: resourceTencentCloudTcaplusTableGroupUpdate,
		Delete: resourceTencentCloudTcaplusTableGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the TcaplusDB cluster to which the table group belongs.",
			},

			"table_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Table group name, which can consist of Chinese characters, English letters or digits, with a maximum length of 32 characters.",
			},

			"table_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Table group ID. It can be specified by the user but must be unique within the cluster. If not specified, the API assigns an auto-incremented value. Changing this forces a new resource to be created.",
			},

			"resource_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of tags of the table group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"table_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tables contained in the table group.",
			},

			"total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total storage size of the tables contained in the table group, in MB.",
			},

			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the table group.",
			},

			"txh_backup_expire_day": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of days after which the Txh backup file of the tables expires and is deleted.",
			},

			"enable_mysql": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether MySQL load balancing is enabled. 0: not enabled, 1: enabling, 2: enabled.",
			},

			"mysql_conn_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VIP of the MySQL load balancing.",
			},

			"mysql_conn_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VPort of the MySQL load balancing.",
			},
		},
	}
}

func resourceTencentCloudTcaplusTableGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_table_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = tcaplusdbv20190823.NewCreateTableGroupRequest()
		response  = tcaplusdbv20190823.NewCreateTableGroupResponse()
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table_group_name"); ok {
		request.TableGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table_group_id"); ok {
		request.TableGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.(*schema.Set).List() {
			resourceTagsMap := item.(map[string]interface{})
			tagInfoUnit := tcaplusdbv20190823.TagInfoUnit{}
			if val, ok := resourceTagsMap["tag_key"].(string); ok && val != "" {
				tagInfoUnit.TagKey = helper.String(val)
			}

			if val, ok := resourceTagsMap["tag_value"].(string); ok && val != "" {
				tagInfoUnit.TagValue = helper.String(val)
			}

			request.ResourceTags = append(request.ResourceTags, &tagInfoUnit)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().CreateTableGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tencentcloud_tcaplus_table_group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tencentcloud_tcaplus_table_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TableGroupId == nil || *response.Response.TableGroupId == "" {
		log.Printf("[CRITAL]%s create tencentcloud_tcaplus_table_group failed, logId=%s, TableGroupId is empty", logId, logId)
		return fmt.Errorf("Create tencentcloud_tcaplus_table_group failed, TableGroupId is empty.")
	}

	d.SetId(strings.Join([]string{clusterId, *response.Response.TableGroupId}, tccommon.FILED_SP))
	return resourceTencentCloudTcaplusTableGroupRead(d, meta)
}

func resourceTencentCloudTcaplusTableGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_table_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tcaplusdbv20190823.NewDescribeTableGroupsRequest()
		service = TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	clusterId, tableGroupId, err := parseTcaplusTableGroupId(d.Id())
	if err != nil {
		return err
	}

	request.ClusterId = helper.String(clusterId)
	request.TableGroupIds = []*string{helper.String(tableGroupId)}
	request.Offset = helper.Int64(0)
	request.Limit = helper.Int64(20)

	response := tcaplusdbv20190823.NewDescribeTableGroupsResponse()
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().DescribeTableGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tencentcloud_tcaplus_table_group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read tencentcloud_tcaplus_table_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.TableGroups) == 0 {
		log.Printf("[CRUD] tencentcloud_tcaplus_table_group id=%s", d.Id())
		d.SetId("")
		return nil
	}

	tableGroupInfo := response.Response.TableGroups[0]
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("table_group_id", tableGroupId)

	if tableGroupInfo.TableGroupName != nil {
		_ = d.Set("table_group_name", tableGroupInfo.TableGroupName)
	}

	if tableGroupInfo.TableCount != nil {
		_ = d.Set("table_count", tableGroupInfo.TableCount)
	}

	if tableGroupInfo.TotalSize != nil {
		_ = d.Set("total_size", tableGroupInfo.TotalSize)
	}

	if tableGroupInfo.CreatedTime != nil {
		_ = d.Set("created_time", tableGroupInfo.CreatedTime)
	}

	if tableGroupInfo.TxhBackupExpireDay != nil {
		_ = d.Set("txh_backup_expire_day", tableGroupInfo.TxhBackupExpireDay)
	}

	if tableGroupInfo.EnableMysql != nil {
		_ = d.Set("enable_mysql", tableGroupInfo.EnableMysql)
	}

	if tableGroupInfo.MysqlConnIp != nil {
		_ = d.Set("mysql_conn_ip", tableGroupInfo.MysqlConnIp)
	}

	if tableGroupInfo.MysqlConnPort != nil {
		_ = d.Set("mysql_conn_port", tableGroupInfo.MysqlConnPort)
	}

	tags, err := service.DescribeTableGroupTags(ctx, clusterId, tableGroupId)
	if err != nil {
		return err
	}

	if len(tags) > 0 {
		resourceTagsList := make([]map[string]interface{}, 0, len(tags))
		for _, tag := range tags {
			resourceTagsMap := map[string]interface{}{}
			if tag.TagKey != nil {
				resourceTagsMap["tag_key"] = *tag.TagKey
			}

			if tag.TagValue != nil {
				resourceTagsMap["tag_value"] = *tag.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)
	} else {
		_ = d.Set("resource_tags", nil)
	}

	return nil
}

func resourceTencentCloudTcaplusTableGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_table_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	clusterId, tableGroupId, err := parseTcaplusTableGroupId(d.Id())
	if err != nil {
		return err
	}

	immutableArgs := []string{"cluster_id", "table_group_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` of tencentcloud_tcaplus_table_group cannot be changed", v)
		}
	}

	if d.HasChange("table_group_name") {
		request := tcaplusdbv20190823.NewModifyTableGroupNameRequest()
		request.ClusterId = helper.String(clusterId)
		request.TableGroupId = helper.String(tableGroupId)
		if v, ok := d.GetOk("table_group_name"); ok {
			request.TableGroupName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().ModifyTableGroupNameWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tencentcloud_tcaplus_table_group name failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("resource_tags") {
		service := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		o, n := d.GetChange("resource_tags")
		oldTags := o.(*schema.Set).List()
		newTags := n.(*schema.Set).List()

		oldMap := make(map[string]string)
		for _, item := range oldTags {
			m := item.(map[string]interface{})
			key, _ := m["tag_key"].(string)
			value, _ := m["tag_value"].(string)
			oldMap[key] = value
		}

		newMap := make(map[string]string)
		for _, item := range newTags {
			m := item.(map[string]interface{})
			key, _ := m["tag_key"].(string)
			value, _ := m["tag_value"].(string)
			newMap[key] = value
		}

		var replaceTags []*tcaplusdbv20190823.TagInfoUnit
		var deleteTags []*tcaplusdbv20190823.TagInfoUnit

		for key, value := range newMap {
			if oldValue, ok := oldMap[key]; !ok || oldValue != value {
				replaceTags = append(replaceTags, &tcaplusdbv20190823.TagInfoUnit{
					TagKey:   helper.String(key),
					TagValue: helper.String(value),
				})
			}
		}

		for key := range oldMap {
			if _, ok := newMap[key]; !ok {
				deleteTags = append(deleteTags, &tcaplusdbv20190823.TagInfoUnit{
					TagKey: helper.String(key),
				})
			}
		}

		if len(replaceTags) > 0 || len(deleteTags) > 0 {
			var taskId string
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := service.ModifyTableGroupTags(ctx, clusterId, tableGroupId, replaceTags, deleteTags)
				if e != nil {
					return tccommon.RetryError(e)
				}

				log.Printf("[DEBUG]%s ModifyTableGroupTags success, taskId:%s", logId, result)
				taskId = result
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s modify tencentcloud_tcaplus_table_group tags failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			waitReq := tcaplusdbv20190823.NewDescribeTasksRequest()
			waitReq.ClusterIds = []*string{helper.String(clusterId)}
			waitReq.TaskIds = []*string{helper.String(taskId)}

			reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				ratelimit.Check(waitReq.GetAction())
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().DescribeTasksWithContext(ctx, waitReq)
				if e != nil {
					return tccommon.RetryError(e)
				}

				if result == nil || result.Response == nil || len(result.Response.TaskInfos) == 0 {
					return resource.NonRetryableError(fmt.Errorf("Describe tencentcloud_tcaplus_table_group tag task failed, task %s not found", taskId))
				}

				taskInfo := result.Response.TaskInfos[0]
				if taskInfo.Progress == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe tencentcloud_tcaplus_table_group tag task failed, task %s progress is nil", taskId))
				}

				if *taskInfo.Progress == 100 {
					return nil
				}

				if *taskInfo.Progress >= 0 {
					return resource.RetryableError(fmt.Errorf("modifying tencentcloud_tcaplus_table_group tags, task %s is in progress", taskId))
				}

				return resource.NonRetryableError(fmt.Errorf("modify tencentcloud_tcaplus_table_group tags failed, task %s returns progress %d", taskId, *taskInfo.Progress))
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s wait tencentcloud_tcaplus_table_group tag task failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	return resourceTencentCloudTcaplusTableGroupRead(d, meta)
}

func resourceTencentCloudTcaplusTableGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_table_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tcaplusdbv20190823.NewDeleteTableGroupRequest()
	)

	clusterId, tableGroupId, err := parseTcaplusTableGroupId(d.Id())
	if err != nil {
		return err
	}

	request.ClusterId = helper.String(clusterId)
	request.TableGroupId = helper.String(tableGroupId)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcaplusClient().DeleteTableGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete tencentcloud_tcaplus_table_group failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tencentcloud_tcaplus_table_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func parseTcaplusTableGroupId(id string) (string, string, error) {
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return "", "", fmt.Errorf("id is broken,%s", id)
	}

	return idSplit[0], idSplit[1], nil
}
