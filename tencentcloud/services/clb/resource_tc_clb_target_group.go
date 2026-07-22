package clb

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func ResourceTencentCloudClbTargetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTargetCreate,
		Read:   resourceTencentCloudClbTargetRead,
		Update: resourceTencentCloudClbTargetUpdate,
		Delete: resourceTencentCloudClbTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TF_target_group",
				Description: "目标群体名称。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "私有网络 ID，默认基于网络。",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description: "目标组默认端口，添加服务器后即可使用。",
			},
			"target_group_instances": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "目标组绑定的后端服务器。",
				Deprecated: "It has been deprecated from version 1.77.3. " +
					"please use `tencentcloud_clb_target_group_instance_attachment` instead.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateIp,
							Description: "目标组实例的内部IP。",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description: "目标组实例的端口。",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "目标组实例的权重。",
						},
						"new_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description: "目标组实例的新端口。",
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "目标组类型，目前支持v1（旧版本目标组）和v2（新版本目标组），默认为v1（旧版本目标组）。",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "目标群体的后端转发协议。对于新版本 (v2) 目标组，此字段是必需的。目前支持 TCP、UDP、HTTP、HTTPS、GRPC。",
			},
			"health_check": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "健康检查配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_switch": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "是否开启健康检查。 true：启用，false：禁用。",
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS", "PING", "CUSTOM", "GRPC"}, false),
							Description: "健康检查协议。有效值：TCP、HTTP、HTTPS、PING、CUSTOM、GRPC。对 v2 目标组有效。",
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description: "健康检查端口。如果不指定，则默认使用后端服务器端口。",
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(2, 60),
							Description: "健康检查响应超时（以秒为单位）。范围：[2，60]。默认值：2。",
						},
						"gap_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(2, 300),
							Description: "健康检查间隔（以秒为单位）。范围：[2，300]。默认值：5。",
						},
						"good_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(2, 10),
							Description: "健康阈值。将后端标记为健康之前所需的连续成功健康检查的次数。范围：[2，10]。默认值：3。",
						},
						"bad_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(2, 10),
							Description: "不健康的阈值。将后端标记为不健康之前所需的连续失败的健康检查次数。范围：[2，10]。默认值：3。",
						},
						"http_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "指示健康状况的 HTTP 状态代码。适用于 HTTP/HTTPS 协议。示例：1 (1xx)、2 (2xx)、4 (3xx)、8 (4xx)、16 (5xx)。可以组合多个值，例如 7（1xx、2xx、3xx）。",
						},
						"http_check_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "健康检查路径。适用于 HTTP/HTTPS 协议。必须以/开头。如果不指定，则默认使用 /。",
						},
						"http_check_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "健康检查域。适用于 HTTP/HTTPS 协议。",
						},
						"http_check_method": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"HEAD", "GET"}, false),
							Description: "健康检查 HTTP 方法。适用于 HTTP/HTTPS 协议。有效值：HEAD、GET。默认值：头部。",
						},
						"http_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"HTTP/1.0", "HTTP/1.1"}, false),
							Description: "用于健康检查的 HTTP 版本。健康检查协议为HTTP时必填。有效值：HTTP/1.0、HTTP/1.1。仅对 TCP 目标组有效。",
						},
						"extended_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用于健康检查的扩展状态代码。",
						},
					},
				},
			},
			"schedule_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"WRR", "LEAST_CONN", "IP_HASH"}, false),
				Description: "调度算法。仅对使用 HTTP/HTTPS/GRPC 协议的 v2 目标组有效。有效值：WRR（加权循环）、LEAST_CONN（最少连接）、IP_HASH（IP 哈希）。默认值：WRR。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "目标组的资源标签。",
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
				Description: "默认后端服务器权重。范围：[0，100]。仅对 v2 目标组有效。设置后，如果未指定，添加到目标组的后端服务器将使用此默认权重。",
			},
			"full_listen_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "这是否是完整的听众目标组。仅对 v2 目标组有效。 true：完整侦听器目标组，false：正常目标组。",
			},
			"keepalive_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "启用保持活动连接。仅对 HTTP/HTTPS 目标组有效。 true：启用，false：禁用。默认值：假。",
			},
			"session_expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.Any(
					validation.IntBetween(30, 3600),
					validation.IntInSlice([]int{0}),
				),
				Description: "会话持续时间（以秒为单位）。仅对使用 HTTP/HTTPS/GRPC 协议的 v2 目标组有效。范围：30-3600 或 0（禁用）。默认值：0（禁用）。",
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IP版本类型。常用值：IPv4、IPv6、IPv6FullChain。",
			},
		},
	}
}

func resourceTencentCloudClbTargetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.create")()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService      = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		vpcId           = d.Get("vpc_id").(string)
		targetGroupName = d.Get("target_group_name").(string)
		port            = uint64(d.Get("port").(int))
		targetGroupType = d.Get("type").(string)
		protocol        = d.Get("protocol").(string)
		insAttachments  = make([]*clb.TargetGroupInstance, 0)
		targetGroupId   string
		err             error
	)

	if v, ok := d.GetOk("target_group_instances"); ok {
		targetGroupInstances := v.([]interface{})
		for _, v1 := range targetGroupInstances {
			value := v1.(map[string]interface{})
			bindIP := value["bind_ip"].(string)
			port := uint64(value["port"].(int))
			weight := uint64(value["weight"].(int))
			newPort := uint64(value["new_port"].(int))
			tgtGrp := &clb.TargetGroupInstance{
				BindIP:  &bindIP,
				Port:    &port,
				Weight:  &weight,
				NewPort: &newPort,
			}
			insAttachments = append(insAttachments, tgtGrp)
		}
	}

	// Extract new parameters
	var healthCheck *clb.TargetGroupHealthCheck
	if v, ok := d.GetOk("health_check"); ok && len(v.([]interface{})) > 0 {
		healthCheck = expandHealthCheck(v.([]interface{}))
	}

	scheduleAlgorithm := d.Get("schedule_algorithm").(string)

	var tags []*clb.TagInfo
	if v, ok := d.GetOk("tags"); ok {
		tags = expandTags(v.(map[string]interface{}))
	}

	var weight *uint64
	if v, ok := d.GetOk("weight"); ok {
		w := uint64(v.(int))
		weight = &w
	}

	var fullListenSwitch *bool
	if v, ok := d.GetOkExists("full_listen_switch"); ok {
		fullListenSwitch = helper.Bool(v.(bool))
	}

	var keepaliveEnable *bool
	if v, ok := d.GetOkExists("keepalive_enable"); ok {
		keepaliveEnable = helper.Bool(v.(bool))
	}

	var sessionExpireTime *uint64
	if v, ok := d.GetOk("session_expire_time"); ok {
		s := uint64(v.(int))
		sessionExpireTime = &s
	}

	ipVersion := d.Get("ip_version").(string)

	targetGroupId, err = clbService.CreateTargetGroup(ctx, targetGroupName, vpcId, port, insAttachments, targetGroupType, protocol,
		healthCheck, scheduleAlgorithm, tags, weight, fullListenSwitch, keepaliveEnable, sessionExpireTime, ipVersion)
	if err != nil {
		return err
	}
	d.SetId(targetGroupId)

	return resourceTencentCloudClbTargetRead(d, meta)

}

func resourceTencentCloudClbTargetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id         = d.Id()
	)
	filters := make(map[string]string)
	targetGroupInfos, err := clbService.DescribeTargetGroupList(ctx, id, filters)
	if err != nil {
		return err
	}
	if len(targetGroupInfos) < 1 {
		d.SetId("")
		return nil
	}

	targetGroup := targetGroupInfos[0]

	_ = d.Set("target_group_name", targetGroup.TargetGroupName)
	_ = d.Set("vpc_id", targetGroup.VpcId)
	_ = d.Set("port", targetGroup.Port)
	_ = d.Set("type", targetGroup.TargetGroupType)
	_ = d.Set("protocol", targetGroup.Protocol)

	// Set new parameters
	if targetGroup.HealthCheck != nil {
		_ = d.Set("health_check", flattenHealthCheck(targetGroup.HealthCheck))
	}

	if targetGroup.ScheduleAlgorithm != nil {
		_ = d.Set("schedule_algorithm", targetGroup.ScheduleAlgorithm)
	}

	if targetGroup.Tag != nil && len(targetGroup.Tag) > 0 {
		_ = d.Set("tags", flattenTags(targetGroup.Tag))
	}

	if targetGroup.Weight != nil {
		_ = d.Set("weight", targetGroup.Weight)
	}

	if targetGroup.FullListenSwitch != nil {
		_ = d.Set("full_listen_switch", targetGroup.FullListenSwitch)
	}

	if targetGroup.KeepaliveEnable != nil {
		_ = d.Set("keepalive_enable", targetGroup.KeepaliveEnable)
	}

	if targetGroup.SessionExpireTime != nil {
		_ = d.Set("session_expire_time", targetGroup.SessionExpireTime)
	}

	if targetGroup.IpVersion != nil {
		_ = d.Set("ip_version", targetGroup.IpVersion)
	}

	return nil
}

func resourceTencentCloudClbTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.update")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService    = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		targetGroupId = d.Id()
	)

	immutableFields := []string{"full_listen_switch", "ip_version", "vpc_id"}
	for _, field := range immutableFields {
		if d.HasChange(field) {
			return fmt.Errorf("field %s cannot be modified after creation", field)
		}
	}

	isChanged := false
	request := clb.NewModifyTargetGroupAttributeRequest()
	request.TargetGroupId = &targetGroupId

	if d.HasChange("target_group_name") {
		request.TargetGroupName = helper.String(d.Get("target_group_name").(string))
		isChanged = true
	}

	if d.HasChange("port") {
		port := uint64(d.Get("port").(int))
		request.Port = &port
		isChanged = true
	}

	if d.HasChange("schedule_algorithm") {
		if v := d.Get("schedule_algorithm").(string); v != "" {
			request.ScheduleAlgorithm = helper.String(v)
			isChanged = true
		}
	}

	if d.HasChange("health_check") {
		if v, ok := d.GetOk("health_check"); ok && len(v.([]interface{})) > 0 {
			request.HealthCheck = expandHealthCheck(v.([]interface{}))
			isChanged = true
		}
	}

	if d.HasChange("weight") {
		if v, ok := d.GetOk("weight"); ok {
			w := uint64(v.(int))
			request.Weight = &w
			isChanged = true
		}
	}

	if d.HasChange("keepalive_enable") {
		request.KeepaliveEnable = helper.Bool(d.Get("keepalive_enable").(bool))
		isChanged = true
	}

	if d.HasChange("session_expire_time") {
		if v, ok := d.GetOk("session_expire_time"); ok {
			s := uint64(v.(int))
			request.SessionExpireTime = &s
			isChanged = true
		}
	}

	if isChanged {
		err := clbService.ModifyTargetGroupAttribute(ctx, request)
		if err != nil {
			return err
		}
	}

	// Handle tags separately if changed
	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("clb", "targetgroup", tcClient.Region, targetGroupId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudClbTargetRead(d, meta)
}

func resourceTencentCloudClbTargetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_target_group.delete")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService    = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		targetGroupId = d.Id()
	)

	err := clbService.DeleteTarget(ctx, targetGroupId)

	if err != nil {
		return err
	}
	return nil
}

// expandHealthCheck converts schema health_check to SDK TargetGroupHealthCheck
func expandHealthCheck(l []interface{}) *clb.TargetGroupHealthCheck {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	hcMap := l[0].(map[string]interface{})
	hc := &clb.TargetGroupHealthCheck{}

	if v, ok := hcMap["health_switch"].(bool); ok {
		hc.HealthSwitch = helper.Bool(v)
	}

	if v, ok := hcMap["protocol"].(string); ok && v != "" {
		hc.Protocol = helper.String(v)
	}

	if v, ok := hcMap["port"].(int); ok && v > 0 {
		hc.Port = helper.IntInt64(v)
	}

	if v, ok := hcMap["timeout"].(int); ok && v > 0 {
		hc.Timeout = helper.IntInt64(v)
	}

	if v, ok := hcMap["gap_time"].(int); ok && v > 0 {
		hc.GapTime = helper.IntInt64(v)
	}

	if v, ok := hcMap["good_limit"].(int); ok && v > 0 {
		hc.GoodLimit = helper.IntInt64(v)
	}

	if v, ok := hcMap["bad_limit"].(int); ok && v > 0 {
		hc.BadLimit = helper.IntInt64(v)
	}

	if v, ok := hcMap["http_code"].(int); ok && v > 0 {
		hc.HttpCode = helper.IntInt64(v)
	}

	if v, ok := hcMap["http_check_path"].(string); ok && v != "" {
		hc.HttpCheckPath = helper.String(v)
	}

	if v, ok := hcMap["http_check_domain"].(string); ok && v != "" {
		hc.HttpCheckDomain = helper.String(v)
	}

	if v, ok := hcMap["http_check_method"].(string); ok && v != "" {
		hc.HttpCheckMethod = helper.String(v)
	}

	if v, ok := hcMap["http_version"].(string); ok && v != "" {
		hc.HttpVersion = helper.String(v)
	}

	if v, ok := hcMap["extended_code"].(string); ok && v != "" {
		hc.ExtendedCode = helper.String(v)
	}

	return hc
}

// expandTags converts map[string]interface{} to []*clb.TagInfo
func expandTags(tags map[string]interface{}) []*clb.TagInfo {
	if len(tags) == 0 {
		return nil
	}

	tagInfos := make([]*clb.TagInfo, 0, len(tags))
	for k, v := range tags {
		tagInfo := &clb.TagInfo{
			TagKey:   helper.String(k),
			TagValue: helper.String(v.(string)),
		}
		tagInfos = append(tagInfos, tagInfo)
	}

	return tagInfos
}

// flattenHealthCheck converts SDK TargetGroupHealthCheck to schema health_check
func flattenHealthCheck(hc *clb.TargetGroupHealthCheck) []interface{} {
	if hc == nil {
		return nil
	}

	result := make(map[string]interface{})

	if hc.HealthSwitch != nil {
		result["health_switch"] = *hc.HealthSwitch
	}

	if hc.Protocol != nil {
		result["protocol"] = *hc.Protocol
	}

	if hc.Port != nil {
		result["port"] = *hc.Port
	}

	if hc.Timeout != nil {
		result["timeout"] = *hc.Timeout
	}

	if hc.GapTime != nil {
		result["gap_time"] = *hc.GapTime
	}

	if hc.GoodLimit != nil {
		result["good_limit"] = *hc.GoodLimit
	}

	if hc.BadLimit != nil {
		result["bad_limit"] = *hc.BadLimit
	}

	if hc.HttpCode != nil {
		result["http_code"] = *hc.HttpCode
	}

	if hc.HttpCheckPath != nil {
		result["http_check_path"] = *hc.HttpCheckPath
	}

	if hc.HttpCheckDomain != nil {
		result["http_check_domain"] = *hc.HttpCheckDomain
	}

	if hc.HttpCheckMethod != nil {
		result["http_check_method"] = *hc.HttpCheckMethod
	}

	if hc.HttpVersion != nil {
		result["http_version"] = *hc.HttpVersion
	}

	if hc.ExtendedCode != nil {
		result["extended_code"] = *hc.ExtendedCode
	}

	return []interface{}{result}
}

// flattenTags converts []*clb.TagInfo to map[string]string
func flattenTags(tagInfos []*clb.TagInfo) map[string]string {
	if len(tagInfos) == 0 {
		return nil
	}

	tags := make(map[string]string, len(tagInfos))
	for _, tagInfo := range tagInfos {
		if tagInfo.TagKey != nil && tagInfo.TagValue != nil {
			tags[*tagInfo.TagKey] = *tagInfo.TagValue
		}
	}

	return tags
}
