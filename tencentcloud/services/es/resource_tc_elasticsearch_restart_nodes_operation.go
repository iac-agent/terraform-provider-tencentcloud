package es

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudElasticsearchRestartNodesOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchRestartNodesOperationCreate,
		Read:   resourceTencentCloudElasticsearchRestartNodesOperationRead,
		Delete: resourceTencentCloudElasticsearchRestartNodesOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"node_names": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 节点 names。",
			},

			"force_restart": {
				Optional:    true,
				Default:     0,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否force restart。",
			},

			"restart_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "可选 restart 模式 在-place,blue-green，其中 表示 restart 和 blue-green restart，respectively. 默认为 在-place。",
			},

			"is_offline": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Node 状态，使用 在 blue-green 模式; 关闭-line 节点 blue-green 是 risky。",
			},
		},
	}
}

func resourceTencentCloudElasticsearchRestartNodesOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_nodes_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = elasticsearch.NewRestartNodesRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("node_names"); ok {
		nodeNamesSet := v.(*schema.Set).List()
		for i := range nodeNamesSet {
			nodeNames := nodeNamesSet[i].(string)
			request.NodeNames = append(request.NodeNames, &nodeNames)
		}
	}

	if v, ok := d.GetOkExists("force_restart"); ok {
		request.ForceRestart = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("restart_mode"); ok {
		request.RestartMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_offline"); ok {
		request.IsOffline = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEsClient().RestartNodes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate elasticsearch RestartNodesOperation failed, reason:%+v", logId, err)
		return err
	}

	elasticsearchService := ElasticsearchService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, 10*tccommon.ReadRetryTimeout, time.Second, elasticsearchService.ElasticsearchInstanceRefreshFunc(instanceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	d.SetId(instanceId)

	return resourceTencentCloudElasticsearchRestartNodesOperationRead(d, meta)
}

func resourceTencentCloudElasticsearchRestartNodesOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_nodes_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudElasticsearchRestartNodesOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_restart_nodes_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
