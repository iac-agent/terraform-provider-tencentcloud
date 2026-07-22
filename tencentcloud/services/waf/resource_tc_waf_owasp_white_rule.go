package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafOwaspWhiteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafOwaspWhiteRuleCreate,
		Read:   resourceTencentCloudWafOwaspWhiteRuleRead,
		Update: resourceTencentCloudWafOwaspWhiteRuleUpdate,
		Delete: resourceTencentCloudWafOwaspWhiteRuleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule 名称",
			},

			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "域名 名称",
			},

			"strategies": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule-Based matching 策略 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies matching 字段.\n\nDifferent matching 字段 结果 在 different matching 参数, logical operators, 和 matching contents. details 是 作为 follows:.\n<表><thead><tr><th>Matching Field</th> <th>Matching Parameter</th> <th>Logical Symbol</th> <th>Matching Content</th></tr></thead> <tbody><tr><td>IP (source IP)</td> <td>Parameters 是 不 支持.</td> <td>ipmatch (match)<br/>ipnmatch (mismatch)</td> <td>Multiple IP addresses 是 separated 通过 commas. A 最大 的 20 IP addresses 是 allowed.</td></tr> <tr><td>IPv6 (source IPv6)</td> <td>Parameters 是 不 支持.</td> <td>ipmatch (match)<br/>ipnmatch (mismatch)</td> <td>A 单个 IPv6 地址 是 支持.</td></tr> <tr><td>Referer (referer)</td> <td>Parameters 是 不 支持.</td> <td>空 (Content 是 空.)<br/>null (do 不 exist)<br/>eq (equal 到)<br/>neq (不 equal 到)<br/>contains (contain)<br/>ncontains (do 不 contain)<br/>len_eq (长度 equals 到)<br/>len_gt (长度 是 greater 比)<br/>len_lt (长度 是 less 比)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter 内容, 使用 最大 的 512 字符.</td></tr> <tr><td>URL (请求 路径)</td> <td>Parameters 是 不 支持.</td> <td>eq (equal 到)<br/>neq (不 equal 到)<br/>contains (contain)<br/>ncontains (do 不 contain)<br/>len_eq (长度 equals 到)<br/>len_gt (长度 是 greater 比)<br/>len_lt (长度 是 \n less 比)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter 内容 starting 使用 /, 使用 最大 的 512 字符.</td></tr> <tr><td>UserAgent (UserAgent)</td> <td>Parameters 是 不 支持.</td><td>Same logical symbols 作为 matching 字段 <font color=\"Red\">Referer</font></td> <td>Enter 内容 使用 最大 的 512 字符.</td></tr> <tr><td>HTTP_METHOD (HTTP 请求 方法)</td> <td>Parameters 是 不 支持.</td> <td>eq (equal 到)<br/>neq (不 equal 到)</td> <td>Enter 方法 名称. uppercase 是 recommended.</td></tr> <tr><td>QUERY_STRING (请求 字符串)</td> <td>Parameters 是 不 支持.</td> <td>Same logical symbol 作为 matching 字段 <font color=\"Red\">Request Path</font></td><td>Enter 内容 使用 最大 的 512 字符.</td></tr> <tr><td>GET (GET 参数 值)</td> <td>Parameter entry 是 支持.</td> <td>contains (contain)<br/>ncontains (do 不 contain)<br/>len_eq (长度 equals 到)<br/>len_gt (长度 是 greater 比)<br/>len_lt (长度 是 less 比)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)</td> <td>Enter 内容 使用 最大 的 512 字符.</td></tr> <tr><td>GET_PARAMS_NAMES (GET 参数 名称)</td> <td>Parameters 是 不 支持.</td> <td>exist (Parameter exists.)<br/>nexist (Parameter does 不 exist.)<br/>len_eq (长度 equals 到)<br/>len_gt (长度 是 greater 比)<br/>len_lt (长度 是 less 比)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)</td><td>Enter 内容 使用 最大 的 512 字符.</td></tr> <tr><td>POST (POST 参数 值)</td> <td>Parameter entry 是 支持.</td> <td>Same logical symbol 作为 matching 字段 <font color=\"Red\">GET Parameter Value</font></td> <td>Enter 内容 使用 最大 的 512 字符.</td></tr> <tr><td>GET_POST_NAMES (POST 参数 名称)</td> <td>Parameters 是 不 支持.</td> <td>Same logical symbol 作为 matching 字段 <font color=\"Red\">GET Parameter Name</font></td> <td>Enter 内容 使用 最大 的 512 字符.</td></tr> <tr><td>POST_BODY (完整 正文)</td> <td>Parameters 是 不 支持.</td> <td>Same logical symbol 作为 matching 字段 <font color=\"Red\">Request Path</font></td><td>Enter 正文 内容 使用 最大 的 512 字符.</td></tr> <tr><td>COOKIE (cookie)</td> <td>Parameters 是 不 支持.</td> <td>空 (Content 是 空.)<br/>null (do 不 exist)<br/>rematch (regular expression matching)</td> <td><font color=\"Red\">Unsupported currently</font></td></tr> <tr><td>GET_COOKIES_NAMES (cookie 参数 名称)</td> <td>Parameters 是 不 支持.</td> <td>Same logical symbol 作为 matching 字段 <font color=\"Red\">GET Parameter Name</font></td> <td>Enter 内容 使用 最大 的 512 字符.</td></tr> <tr><td>ARGS_COOKIE (cookie 参数 值)</td> <td>Parameter entry 是 支持.</td> <td>Same logical symbol 作为 matching 字段 <font color=\"Red\">GET Parameter Value</font></td> <td>Enter content512 字符 限制</td></tr><tr><td>GET_HEADERS_NAMES (Header 参数 名称)</td><td>参数 不 支持</td><td>exsit (参数 exists)<br/>nexsit (参数 does 不 exist)<br/>len_eq (LENGTH equal)<br/>len_gt (LENGTH greater 比)<br/>len_lt (LENGTH less 比)<br/>strprefix (prefix match)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td><td>enter CONTENT, lowercase 是 recommended, up 到 512 字符</td></tr><tr><td>ARGS_Header (Header 参数 值)</td><td>support 参数 entry</td><td>contains (include)<br/>ncontains (does 不 include)<br/>len_eq (LENGTH equal)<br/>len_gt (LENGTH greater 比)<br/>len_lt (LENGTH less 比)<br/>strprefix (prefix match)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td><td>enter CONTENT, up 到 512 字符</td></tr><tr><td>CONTENT_LENGTH (CONTENT-LENGTH)</td><td>support 参数 entry</td><td>numgt (值 greater 比)<br/>numlt (值 smaller 比)<br/>numeq (值 equal 到)<br/></td><td>enter 整数 between 0-9999999999999</td></tr><tr><td>IP_GEO (source IP geolocation)</td><td>support 参数 entry</td><td>GEO_in (belong)<br/>GEO_not_in (not_in)<br/></td><td>enter CONTENT, up 到 10240 字符, 格式: serialized JSON, 格式: [{\"Country\":\"china\",\"Region\":\"guangdong\",\"City\":\"shenzhen\"}]</td></tr><tr><td>CAPTCHA_RISK (CAPTCHA RISK)</td><td>参数 不 支持</td><td>eq (equal)<br/>neq (不 equal 到)<br/>belong (belong)<br/>not_belong (不 belong 到)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter RISK 级别 值, 值 范围 0-255</td></tr><tr><td>CAPTCHA_DEVICE_RISK (CAPTCHA DEVICE RISK)</td><td>参数 不 支持</td><td>eq (equal)<br/>neq (不 equal 到)<br/>belong (belong)<br/>not_belong (不 belong 到)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter DEVICE RISK 代码, 有效 值: 101, 201, 301, 401, 501, 601, 701</td></tr><tr><td>CAPTCHAR_SCORE (CAPTCHA RISK assessment SCORE)</td><td>参数 不 支持</td><td>numeq (值 equal 到)<br/>numgt (值 greater 比)<br/>numlt (值 smaller 比)<br/>numle (less 比 或 equal 到)<br/>numge (值 是 greater 比 或 equal 到)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter assessment SCORE, 值 范围 0-100</td></tr>.\n</tbody></表>.",
						},
						"compare_func": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "指定logic symbol. \n\nLogical symbols 是 divided into following types:.\nEmpty (内容 是 空).\nnull (不 found).\nEq (equal 到).\nneq (不 equal 到).\n包含(contain).\nn包含(do 不 contain).\nstrprefix (prefix matching).\nstrsuffix (suffix matching).\nLen_eq (长度 equals 到).\nLen_gt (长度 greater 比).\nLen_lt (长度 less 比).\nipmatch (belong).\nipnmatch (not_in).\nnumgt (值 greater 比).\nNumValue smaller 比].\nValue equal 到.\nnumneq (值 不 equal 到).\nnumle (less 比 或 equal 到).\nnumge (值 是 greater 比 或 equal 到).\ngeo_in (IP geographic belong).\ngeo_not_in (IP geographic not_in).\n指定different logical operators 对于 matching 字段. 对于 details，see matching 字段 表 above。",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "指定match 内容\n\nCurrently，当 match 字段 是 COOKIE (COOKIE)，match 内容 不是必填项. all others 是 needed。",
						},
						"arg": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "指定matching 参数.\n\nConfiguration 参数 是 divided into two 数据 types: 参数 不 支持 和 support 参数.\nWhen match 字段 是 一个 的 following four， matching 参数 可以 是 entered，otherwise 不 支持.\nGET (get 参数 值).\v\t\t\nPOST (post 参数 值).\v\t\t\nARGS_COOKIE (COOKIE 参数 值).\v\t\t\nARGS_HEADER (HEADER 参数 值)。",
						},
						"case_not_sensitive": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Case-Sensitive.\nCase-Insensitive。",
						},
					},
				},
			},

			"ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "ID 列表 allowlisted 规则。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Allowlist 类型 有效值：0 (allowlisting 通过 特定 规则 ID)，1 (allowlisting 通过 规则 类型)。",
			},

			"job_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule execution 模式: TimedJob 表示scheduled execution. CronJob 表示periodic execution。",
			},

			"job_date_time": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Scheduled 任务 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timed": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Time 参数 对于 scheduled execution。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_date_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Start 时间戳，（秒）。",
									},
									"end_date_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "End 时间戳，（秒）。",
									},
								},
							},
						},
						"cron": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Time 参数 对于 periodic execution。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Execution day 的 each month。",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"w_days": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Execution day 的 each week。",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"start_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "开始时间。",
									},
									"end_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "结束时间。",
									},
								},
							},
						},
						"time_t_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "指定time 可用区",
						},
					},
				},
			},

			"expire_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "如果 JobDateTime 字段 是 不 集合，此 字段 是 使用. 0 表示 permanent，other 值 indicate cutoff 时间 对于 scheduled effect (单位: 秒)。",
			},

			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Rule 状态 有效值：0 (已禁用)，1 (已启用). 已启用 通过 默认值。",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule ID。",
			},
		},
	}
}

func resourceTencentCloudWafOwaspWhiteRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_white_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = wafv20180125.NewCreateOwaspWhiteRuleRequest()
		response = wafv20180125.NewCreateOwaspWhiteRuleResponse()
		domain   string
		ruleId   string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("strategies"); ok {
		for _, item := range v.([]interface{}) {
			strategiesMap := item.(map[string]interface{})
			strategy := wafv20180125.Strategy{}
			if v, ok := strategiesMap["field"].(string); ok && v != "" {
				strategy.Field = helper.String(v)
			}

			if v, ok := strategiesMap["compare_func"].(string); ok && v != "" {
				strategy.CompareFunc = helper.String(v)
			}

			if v, ok := strategiesMap["content"].(string); ok && v != "" {
				strategy.Content = helper.String(v)
			}

			if v, ok := strategiesMap["arg"].(string); ok && v != "" {
				strategy.Arg = helper.String(v)
			}

			if v, ok := strategiesMap["case_not_sensitive"].(int); ok {
				strategy.CaseNotSensitive = helper.IntUint64(v)
			}

			request.Strategies = append(request.Strategies, &strategy)
		}
	}

	if v, ok := d.GetOk("ids"); ok {
		idsSet := v.(*schema.Set).List()
		for i := range idsSet {
			ids := idsSet[i].(int)
			request.Ids = append(request.Ids, helper.IntUint64(ids))
		}
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("job_type"); ok {
		request.JobType = helper.String(v.(string))
	}

	if jobDateTimeMap, ok := helper.InterfacesHeadMap(d, "job_date_time"); ok {
		jobDateTime := wafv20180125.JobDateTime{}
		if v, ok := jobDateTimeMap["timed"]; ok {
			for _, item := range v.([]interface{}) {
				timedMap := item.(map[string]interface{})
				timedJob := wafv20180125.TimedJob{}
				if v, ok := timedMap["start_date_time"].(int); ok {
					timedJob.StartDateTime = helper.IntUint64(v)
				}

				if v, ok := timedMap["end_date_time"].(int); ok {
					timedJob.EndDateTime = helper.IntUint64(v)
				}

				jobDateTime.Timed = append(jobDateTime.Timed, &timedJob)
			}
		}

		if v, ok := jobDateTimeMap["cron"]; ok {
			for _, item := range v.([]interface{}) {
				cronMap := item.(map[string]interface{})
				cronJob := wafv20180125.CronJob{}
				if v, ok := cronMap["days"]; ok {
					daysSet := v.(*schema.Set).List()
					for i := range daysSet {
						days := daysSet[i].(int)
						cronJob.Days = append(cronJob.Days, helper.IntUint64(days))
					}
				}

				if v, ok := cronMap["w_days"]; ok {
					wDaysSet := v.(*schema.Set).List()
					for i := range wDaysSet {
						wDays := wDaysSet[i].(int)
						cronJob.WDays = append(cronJob.WDays, helper.IntUint64(wDays))
					}
				}

				if v, ok := cronMap["start_time"].(string); ok && v != "" {
					cronJob.StartTime = helper.String(v)
				}

				if v, ok := cronMap["end_time"].(string); ok && v != "" {
					cronJob.EndTime = helper.String(v)
				}

				jobDateTime.Cron = append(jobDateTime.Cron, &cronJob)
			}
		}

		if v, ok := jobDateTimeMap["time_t_zone"].(string); ok && v != "" {
			jobDateTime.TimeTZone = helper.String(v)
		}

		request.JobDateTime = &jobDateTime
	}

	if v, ok := d.GetOkExists("expire_time"); ok {
		request.ExpireTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().CreateOwaspWhiteRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create waf owasp white rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create waf owasp white rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.RuleId == nil {
		return fmt.Errorf("RuleId is nil.")
	}

	ruleId = helper.UInt64ToStr(*response.Response.RuleId)
	d.SetId(strings.Join([]string{domain, ruleId}, tccommon.FILED_SP))
	return resourceTencentCloudWafOwaspWhiteRuleRead(d, meta)
}

func resourceTencentCloudWafOwaspWhiteRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_white_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	respData, err := service.DescribeWafOwaspWhiteRuleById(ctx, domain, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_owasp_white_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Strategies != nil {
		tmpList := make([]map[string]interface{}, 0, len(respData.Strategies))
		for _, strategies := range respData.Strategies {
			strategiesMap := map[string]interface{}{}
			if strategies.Field != nil {
				strategiesMap["field"] = strategies.Field
			}

			if strategies.CompareFunc != nil {
				strategiesMap["compare_func"] = strategies.CompareFunc
			}

			if strategies.Content != nil {
				strategiesMap["content"] = strategies.Content
			}

			if strategies.Arg != nil {
				strategiesMap["arg"] = strategies.Arg
			}

			if strategies.CaseNotSensitive != nil {
				strategiesMap["case_not_sensitive"] = strategies.CaseNotSensitive
			}

			tmpList = append(tmpList, strategiesMap)
		}

		_ = d.Set("strategies", tmpList)
	}

	if respData.Ids != nil {
		_ = d.Set("ids", respData.Ids)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.JobType != nil {
		_ = d.Set("job_type", respData.JobType)
	}

	if respData.JobDateTime != nil {
		jobDateTimeMap := map[string]interface{}{}
		if respData.JobDateTime.Timed != nil {
			timedList := make([]map[string]interface{}, 0, len(respData.JobDateTime.Timed))
			for _, timed := range respData.JobDateTime.Timed {
				timedMap := map[string]interface{}{}
				if timed.StartDateTime != nil {
					timedMap["start_date_time"] = timed.StartDateTime
				}

				if timed.EndDateTime != nil {
					timedMap["end_date_time"] = timed.EndDateTime
				}

				timedList = append(timedList, timedMap)
			}

			jobDateTimeMap["timed"] = timedList
		}

		if respData.JobDateTime.Cron != nil {
			cronList := make([]map[string]interface{}, 0, len(respData.JobDateTime.Cron))
			for _, cron := range respData.JobDateTime.Cron {
				cronMap := map[string]interface{}{}
				if cron.Days != nil {
					cronMap["days"] = cron.Days
				}

				if cron.WDays != nil {
					cronMap["w_days"] = cron.WDays
				}

				if cron.StartTime != nil {
					cronMap["start_time"] = cron.StartTime
				}

				if cron.EndTime != nil {
					cronMap["end_time"] = cron.EndTime
				}

				cronList = append(cronList, cronMap)
			}

			jobDateTimeMap["cron"] = cronList
		}

		if respData.JobDateTime.TimeTZone != nil {
			jobDateTimeMap["time_t_zone"] = respData.JobDateTime.TimeTZone
		}

		_ = d.Set("job_date_time", []interface{}{jobDateTimeMap})
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	return nil
}

func resourceTencentCloudWafOwaspWhiteRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_white_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "strategies", "ids", "type", "job_type", "job_date_time", "expire_time", "status"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wafv20180125.NewModifyOwaspWhiteRuleRequest()
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("strategies"); ok {
			for _, item := range v.([]interface{}) {
				strategiesMap := item.(map[string]interface{})
				strategy := wafv20180125.Strategy{}
				if v, ok := strategiesMap["field"].(string); ok && v != "" {
					strategy.Field = helper.String(v)
				}

				if v, ok := strategiesMap["compare_func"].(string); ok && v != "" {
					strategy.CompareFunc = helper.String(v)
				}

				if v, ok := strategiesMap["content"].(string); ok && v != "" {
					strategy.Content = helper.String(v)
				}

				if v, ok := strategiesMap["arg"].(string); ok && v != "" {
					strategy.Arg = helper.String(v)
				}

				if v, ok := strategiesMap["case_not_sensitive"].(int); ok {
					strategy.CaseNotSensitive = helper.IntUint64(v)
				}

				request.Strategies = append(request.Strategies, &strategy)
			}
		}

		if v, ok := d.GetOk("ids"); ok {
			idsSet := v.(*schema.Set).List()
			for i := range idsSet {
				ids := idsSet[i].(int)
				request.Ids = append(request.Ids, helper.IntUint64(ids))
			}
		}

		if v, ok := d.GetOkExists("type"); ok {
			request.Type = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("job_type"); ok {
			request.JobType = helper.String(v.(string))
		}

		if jobDateTimeMap, ok := helper.InterfacesHeadMap(d, "job_date_time"); ok {
			jobDateTime := wafv20180125.JobDateTime{}
			if v, ok := jobDateTimeMap["timed"]; ok {
				for _, item := range v.([]interface{}) {
					timedMap := item.(map[string]interface{})
					timedJob := wafv20180125.TimedJob{}
					if v, ok := timedMap["start_date_time"].(int); ok {
						timedJob.StartDateTime = helper.IntUint64(v)
					}

					if v, ok := timedMap["end_date_time"].(int); ok {
						timedJob.EndDateTime = helper.IntUint64(v)
					}

					jobDateTime.Timed = append(jobDateTime.Timed, &timedJob)
				}
			}

			if v, ok := jobDateTimeMap["cron"]; ok {
				for _, item := range v.([]interface{}) {
					cronMap := item.(map[string]interface{})
					cronJob := wafv20180125.CronJob{}
					if v, ok := cronMap["days"]; ok {
						daysSet := v.(*schema.Set).List()
						for i := range daysSet {
							days := daysSet[i].(int)
							cronJob.Days = append(cronJob.Days, helper.IntUint64(days))
						}
					}

					if v, ok := cronMap["w_days"]; ok {
						wDaysSet := v.(*schema.Set).List()
						for i := range wDaysSet {
							wDays := wDaysSet[i].(int)
							cronJob.WDays = append(cronJob.WDays, helper.IntUint64(wDays))
						}
					}

					if v, ok := cronMap["start_time"].(string); ok && v != "" {
						cronJob.StartTime = helper.String(v)
					}

					if v, ok := cronMap["end_time"].(string); ok && v != "" {
						cronJob.EndTime = helper.String(v)
					}

					jobDateTime.Cron = append(jobDateTime.Cron, &cronJob)
				}
			}

			if v, ok := jobDateTimeMap["time_t_zone"].(string); ok && v != "" {
				jobDateTime.TimeTZone = helper.String(v)
			}

			request.JobDateTime = &jobDateTime
		}

		if v, ok := d.GetOkExists("expire_time"); ok {
			request.ExpireTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}

		request.Domain = &domain
		request.RuleId = helper.StrToUint64Point(ruleId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyOwaspWhiteRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf owasp white rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWafOwaspWhiteRuleRead(d, meta)
}

func resourceTencentCloudWafOwaspWhiteRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_white_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wafv20180125.NewDeleteOwaspWhiteRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	request.Domain = &domain
	request.Ids = append(request.Ids, helper.StrToUint64Point(ruleId))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().DeleteOwaspWhiteRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf owasp white rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
