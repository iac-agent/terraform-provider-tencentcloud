package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssXp2pDetailInfoList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssXp2pDetailInfoListRead,
		Schema: map[string]*schema.Schema{
			"query_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "UTC minute granularity 查询 时间 对于 querying usage 数据 对于 特定 minute 是 在 格式: yyyy-mm-ddTHH:MM:00Z. Please refer 到 link https://云.tencent.com/document/product/266/11732#I.For 示例，如果 本地 时间 是 2019-01-08 10:00:00 在 Beijing， corresponding UTC 时间 would 是 2019-01-08T10:00:00+08:00.此 查询 支持 数据 从 past six months。",
			},

			"type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "类型 数组 可以 是 用于指定type 的 media 内容 到 是 queried. two 可用 options 是 live 对于 live streaming 和 vod 对于 视频 在 demand. 如果 无 类型 是 指定， 查询 将 include both live 和 VOD 内容 通过 默认值。",
			},

			"stream_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "流 数组 可以 是 用于指定streams 到 是 queried. 如果 无 流 是 指定， 查询 将 include all streams 通过 默认值。",
			},

			"dimension": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "dimension 参数 可以 是 用于指定dimension 对于 查询. 如果 此 参数 是 不 passed， 查询 将 默认为 流-级别 数据. 如果 您 pass 此 参数，它 将 仅 retrieve 数据 对于 指定 dimension. 可用 dimension currently 支持 是 AppId dimension，其中 allows 您 到 查询 数据 based 在 应用 ID. Please note 该 返回 字段 将 是 related 到 指定 dimension。",
			},

			"data_info_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "P2P streaming statistical 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cdn_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CDN 流量。",
						},
						"p2p_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "P2P 流量。",
						},
						"stuck_people": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "People count。",
						},
						"stuck_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Count。",
						},
						"online_people": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Online numbers。",
						},
						"request": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Request numbers。",
						},
						"request_success": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Success numbers。",
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "requested 格式 对于 时间 在 UTC 使用 一个-minute granularity 是 yyyy-mm-ddTHH:MM:SSZ. 此 格式 follows ISO 8601 standard 和 是 commonly 用于representing timestamps 在 UTC. For more 信息 和 examples，您 可以 refer 到 link 提供: https://云.tencent.com/document/product/266/11732#I。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型，divided into two categories: live 和 vod.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stream ID.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AppId. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudCssXp2pDetailInfoListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_xp2p_detail_info_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("query_time"); ok {
		paramMap["QueryTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		typeSet := v.(*schema.Set).List()
		paramMap["Type"] = helper.InterfacesStringsPoint(typeSet)
	}

	if v, ok := d.GetOk("stream_names"); ok {
		streamNamesSet := v.(*schema.Set).List()
		paramMap["StreamNames"] = helper.InterfacesStringsPoint(streamNamesSet)
	}

	if v, ok := d.GetOk("dimension"); ok {
		dimensionSet := v.(*schema.Set).List()
		paramMap["Dimension"] = helper.InterfacesStringsPoint(dimensionSet)
	}

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var dataInfoList []*css.XP2PDetailInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssXp2pDetailInfoListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dataInfoList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dataInfoList))
	tmpList := make([]map[string]interface{}, 0, len(dataInfoList))

	if dataInfoList != nil {
		for _, xP2PDetailInfo := range dataInfoList {
			xP2PDetailInfoMap := map[string]interface{}{}

			if xP2PDetailInfo.CdnBytes != nil {
				xP2PDetailInfoMap["cdn_bytes"] = xP2PDetailInfo.CdnBytes
			}

			if xP2PDetailInfo.P2pBytes != nil {
				xP2PDetailInfoMap["p2p_bytes"] = xP2PDetailInfo.P2pBytes
			}

			if xP2PDetailInfo.StuckPeople != nil {
				xP2PDetailInfoMap["stuck_people"] = xP2PDetailInfo.StuckPeople
			}

			if xP2PDetailInfo.StuckTimes != nil {
				xP2PDetailInfoMap["stuck_times"] = xP2PDetailInfo.StuckTimes
			}

			if xP2PDetailInfo.OnlinePeople != nil {
				xP2PDetailInfoMap["online_people"] = xP2PDetailInfo.OnlinePeople
			}

			if xP2PDetailInfo.Request != nil {
				xP2PDetailInfoMap["request"] = xP2PDetailInfo.Request
			}

			if xP2PDetailInfo.RequestSuccess != nil {
				xP2PDetailInfoMap["request_success"] = xP2PDetailInfo.RequestSuccess
			}

			if xP2PDetailInfo.Time != nil {
				xP2PDetailInfoMap["time"] = xP2PDetailInfo.Time
			}

			if xP2PDetailInfo.Type != nil {
				xP2PDetailInfoMap["type"] = xP2PDetailInfo.Type
			}

			if xP2PDetailInfo.StreamName != nil {
				xP2PDetailInfoMap["stream_name"] = xP2PDetailInfo.StreamName
			}

			if xP2PDetailInfo.AppId != nil {
				xP2PDetailInfoMap["app_id"] = xP2PDetailInfo.AppId
			}

			ids = append(ids, *xP2PDetailInfo.StreamName)
			tmpList = append(tmpList, xP2PDetailInfoMap)
		}

		_ = d.Set("data_info_list", tmpList)
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
