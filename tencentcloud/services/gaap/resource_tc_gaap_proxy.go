package gaap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var gaapActionMu = &sync.Mutex{}

func ResourceTencentCloudGaapProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapProxyCreate,
		Read:   resourceTencentCloudGaapProxyRead,
		Update: resourceTencentCloudGaapProxyUpdate,
		Delete: resourceTencentCloudGaapProxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 30),
				Description:  "名称 GAAP proxy， 最大 长度 是 30。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID 项目 within GAAP proxy，`0` 表示 是 默认值 项目。",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum 带宽 的 GAAP proxy，单位 是 Mbps. 有效 值: `10`，`20`，`50`，`100`，`200`，`500`，`1000`，`2000`，`5000` 和 `10000`. To 集合 `2000`，`5000` 或 `10000`，您 need 到 apply 对于 whitelist 从 Tencent Cloud。",
			},
			"concurrent": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum 并发 的 GAAP proxy，单位 是 10k. 有效 值: `2`，`5`，`10`，`20`，`30`，`40`，`50`，`60`，`70`，`80`，`90`，`100`，`150`，`200`，`250` 和 `300`. To 集合 `150`，`200`，`250` 或 `300`，您 need 到 apply 对于 whitelist 从 Tencent Cloud。",
			},
			"access_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Access 地域 的 GAAP proxy. 有效 值: `Hongkong`，`SoutheastAsia`，`Korea`，`Europe`，`NorthAmerica`，`Canada`，`WestIndia`，`Thailand`，`Virginia`，`Japan`，`Taipei`，`SL_AZURE_NorthUAE`，`SL_AZURE_EastAUS`，`SL_AZURE_NorthCentralUSA`，`SL_AZURE_SouthIndia`，`SL_AZURE_SouthBrazil`，`SL_AZURE_NorthZAF`，`SL_AZURE_SoutheastAsia`，`SL_AZURE_CentralFrance`，`SL_AZURE_SouthEngland`，`SL_AZURE_EastUS`，`SL_AZURE_WestUS`，`SL_AZURE_SouthCentralUSA`，`Jakarta`，`Beijing`，`Shanghai`，`Guangzhou`，`Chengdu`，`SL_AZURE_NorwayEast`，`Chongqing`，`Nanjing`，`SaoPaulo`，`SL_AZURE_JapanEast`，`Changsha`，`Xian`，`Wuhan`，`Fuzhou`，`Shenyang`，`Zhengzhou`，`Jinan`，`Hangzhou`，`Shijiazhuang`，`Hefei`。",
			},
			"realserver_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "地域 的 GAAP realserver. 有效 值: `Hongkong`，`SoutheastAsia`，`Korea`，`Europe`，`NorthAmerica`，`Canada`，`WestIndia`，`Thailand`，`Virginia`，`Japan`，`Taipei`，`SL_AZURE_NorthUAE`，`SL_AZURE_EastAUS`，`SL_AZURE_NorthCentralUSA`，`SL_AZURE_SouthIndia`，`SL_AZURE_SouthBrazil`，`SL_AZURE_NorthZAF`，`SL_AZURE_SoutheastAsia`，`SL_AZURE_CentralFrance`，`SL_AZURE_SouthEngland`，`SL_AZURE_EastUS`，`SL_AZURE_WestUS`，`SL_AZURE_SouthCentralUSA`，`Jakarta`，`Beijing`，`Shanghai`，`Guangzhou`，`Chengdu`，`SL_AZURE_NorwayEast`，`Chongqing`，`Nanjing`，`SaoPaulo`，`SL_AZURE_JapanEast`。",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "表示是否GAAP proxy 是 已启用，默认值为 `true`。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 GAAP proxy. 标签 该 do 不 exist 是 不 创建 automatically。",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(PROXY_NETWORK_TYPE),
				Description:  "Network 类型 `normal`: regular BGP，`cn2`: boutique BGP，`triple`: triple play。",
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 GAAP proxy。",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "状态 GAAP proxy。",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access 域名 的 GAAP proxy。",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access IP 的 GAAP proxy。",
			},
			"scalable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "表示是否GAAP proxy 可以 scalable。",
			},
			"support_protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Supported protocols 的 GAAP proxy。",
			},
			"forward_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Forwarding IP 的 GAAP proxy。",
			},
		},
	}
}

func resourceTencentCloudGaapProxyCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_proxy.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]interface{})
	name := d.Get("name").(string)
	projectId := d.Get("project_id").(int)
	bandwidth := d.Get("bandwidth").(int)
	concurrent := d.Get("concurrent").(int)
	accessRegion := d.Get("access_region").(string)
	realserverRegion := d.Get("realserver_region").(string)
	enable := d.Get("enable").(bool)
	tags := helper.GetTags(d, "tags")

	if v, ok := d.GetOk("network_type"); ok {
		params["network_type"] = v.(string)
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id, err := service.CreateProxy(ctx, name, accessRegion, realserverRegion, bandwidth, concurrent, projectId, tags, params)
	if err != nil {
		return err
	}

	d.SetId(id)

	if !enable {
		if err := service.DisableProxy(ctx, id); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapProxyRead(d, m)
}

func resourceTencentCloudGaapProxyRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_proxy.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	proxies, err := service.DescribeProxies(ctx, []string{id}, nil, nil, nil, nil)
	if err != nil {
		return err
	}

	var proxy *gaap.ProxyInfo
	for _, p := range proxies {
		if p.ProxyId == nil {
			return errors.New("proxy id is nil")
		}
		if *p.ProxyId == id {
			proxy = p
			break
		}
	}

	if proxy == nil {
		d.SetId("")
		return nil
	}

	if proxy.ProxyName == nil {
		return errors.New("proxy name is nil")
	}
	_ = d.Set("name", proxy.ProxyName)

	if proxy.ProjectId == nil {
		return errors.New("proxy project id is nil")
	}
	_ = d.Set("project_id", proxy.ProjectId)

	if proxy.Bandwidth == nil {
		return errors.New("proxy bandwidth is nil")
	}
	_ = d.Set("bandwidth", proxy.Bandwidth)

	if proxy.Concurrent == nil {
		return errors.New("proxy concurrent is nil")
	}
	_ = d.Set("concurrent", proxy.Concurrent)

	if proxy.AccessRegion == nil {
		return errors.New("proxy access region is nil")
	}
	_ = d.Set("access_region", proxy.AccessRegion)

	if proxy.RealServerRegion == nil {
		return errors.New("proxy realserver region is nil")
	}
	_ = d.Set("realserver_region", proxy.RealServerRegion)

	if proxy.Status == nil {
		return errors.New("proxy status is nil")
	}
	_ = d.Set("enable", *proxy.Status == GAAP_PROXY_RUNNING)
	_ = d.Set("status", proxy.Status)

	if len(proxy.TagSet) > 0 {
		tags := make(map[string]string, len(proxy.TagSet))
		for _, tag := range proxy.TagSet {
			tags[*tag.TagKey] = *tag.TagValue
		}
		_ = d.Set("tags", tags)
	}

	if proxy.CreateTime == nil {
		return errors.New("proxy create time is nil")
	}
	_ = d.Set("create_time", helper.FormatUnixTime(*proxy.CreateTime))

	if proxy.Domain == nil {
		return errors.New("proxy access domain is nil")
	}
	_ = d.Set("domain", proxy.Domain)

	if proxy.IP == nil {
		return errors.New("proxy access IP is nil")
	}
	_ = d.Set("ip", proxy.IP)

	if proxy.Scalarable == nil {
		return errors.New("proxy scalable is nil")
	}
	_ = d.Set("scalable", *proxy.Scalarable == 1)

	if len(proxy.SupportProtocols) == 0 {
		return errors.New("proxy support protocols is empty")
	}
	supportProtocols := make([]string, 0, len(proxy.SupportProtocols))
	for _, sp := range proxy.SupportProtocols {
		supportProtocols = append(supportProtocols, *sp)
	}
	_ = d.Set("support_protocols", supportProtocols)

	if proxy.ForwardIP == nil {
		return errors.New("proxy forward ip is nil")
	}
	_ = d.Set("forward_ip", proxy.ForwardIP)

	if proxy.NetworkType != nil {
		_ = d.Set("network_type", proxy.NetworkType)
	}
	return nil
}

func resourceTencentCloudGaapProxyUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_proxy.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	gaapService := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		if err := gaapService.ModifyProxyName(ctx, id, name); err != nil {
			return err
		}

	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		if err := gaapService.ModifyProxyProjectId(ctx, id, projectId); err != nil {
			return err
		}

	}

	if d.HasChange("bandwidth") || d.HasChange("concurrent") {
		var (
			bandwidth  *int
			concurrent *int
		)
		if d.HasChange("bandwidth") {
			bandwidth = common.IntPtr(d.Get("bandwidth").(int))
		}
		if d.HasChange("concurrent") {
			concurrent = common.IntPtr(d.Get("concurrent").(int))
		}
		if err := gaapService.ModifyProxyConfiguration(ctx, id, bandwidth, concurrent); err != nil {
			return err
		}
		//deal with sync delay
		time.Sleep(time.Duration(10) * time.Second)
	}

	if d.HasChange("enable") {
		enable := d.Get("enable").(bool)
		if enable {
			if err := gaapService.EnableProxy(ctx, id); err != nil {
				return err
			}
		} else {
			if err := gaapService.DisableProxy(ctx, id); err != nil {
				return err
			}
		}

	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := svctag.NewTagService(m.(tccommon.ProviderMeta).GetAPIV3Conn())

		region := m.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::gaap:%s:uin/:proxy/%s", region, id)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudGaapProxyRead(d, m)
}

func resourceTencentCloudGaapProxyDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_proxy.update")()
	//gaapActionMu.Lock()
	//defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	createTimeStr := d.Get("create_time").(string)

	if createTime, err := helper.ParseTime(createTimeStr); err == nil {
		if !time.Now().After(createTime.Add(2 * time.Minute)) {
			log.Printf("[DEBUG]%s proxy can't be deleted unless it has lived 2 minutes", logId)
			time.Sleep(time.Until(createTime.Add(2 * time.Minute)))
		}
	} else {
		log.Printf("[WARN]%s parse create time failed, delete immediately", logId)
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteProxy(ctx, id)
}
