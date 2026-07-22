package es

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudElasticsearchInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 实例 到 是 queried。",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 实例 到 是 queried。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 实例 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 elasticsearch 实例. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 实例。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 实例。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability 可用区",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC 网络。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC 子网。",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 的 实例。",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "charge 类型 实例。",
						},
						"deploy_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster 部署 模式",
						},
						"multi_zone_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Details 的 AZs 在 multi-AZ 部署 模式",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability 可用区",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID VPC 子网。",
									},
								},
							},
						},
						"license_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "License 类型",
						},
						"node_info_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node 信息 列表，其中 describe 规格 信息 的 various types 的 nodes 在 集群。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "节点数量",
									},
									"node_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node 规格。",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node 类型",
									},
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node 磁盘 类型",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Node 磁盘 大小。",
									},
									"encrypt": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Decides 此 磁盘 encrypted 或 不。",
									},
								},
							},
						},
						"basic_security_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否enable X-Pack 安全 authentication 在 Basic Edition 6.8 和 above。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A mapping 的 标签 到 assign 到 实例。",
						},
						"elasticsearch_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Elasticsearch 域名 名称",
						},
						"elasticsearch_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Elasticsearch VIP",
						},
						"elasticsearch_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Elasticsearch 端口",
						},
						"elasticsearch_public_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Elasticsearch 公有 URL",
						},
						"kibana_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kibana 访问 URL",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 创建时间。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudElasticsearchInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_elasticsearch_instances.read")()

	var (
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		elasticsearchService = ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId           string
		instanceName         string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		instanceName = v.(string)
	}

	tags := helper.GetTags(d, "tags")
	var instances []*es.InstanceInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instances, errRet = elasticsearchService.DescribeInstancesByFilter(ctx, instanceId, instanceName, tags)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		return nil
	})

	if err != nil {
		return nil
	}

	instanceList := make([]map[string]interface{}, 0, len(instances))
	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		tags := make(map[string]string, len(instance.TagList))
		for _, tag := range instance.TagList {
			tags[*tag.TagKey] = *tag.TagValue
		}

		mapping := map[string]interface{}{
			"instance_id":          instance.InstanceId,
			"instance_name":        instance.InstanceName,
			"availability_zone":    instance.Zone,
			"vpc_id":               instance.VpcUid,
			"subnet_id":            instance.SubnetUid,
			"version":              instance.EsVersion,
			"charge_type":          instance.ChargeType,
			"deploy_mode":          instance.DeployMode,
			"license_type":         instance.LicenseType,
			"basic_security_type":  instance.SecurityType,
			"tags":                 tags,
			"elasticsearch_domain": instance.EsDomain,
			"elasticsearch_vip":    instance.EsVip,
			"elasticsearch_port":   instance.EsPort,
			"kibana_url":           instance.KibanaUrl,
			"create_time":          instance.CreateTime,
		}

		if instance.EsPublicUrl != nil {
			mapping["elasticsearch_public_url"] = instance.EsPublicUrl
		}

		if instance.MultiZoneInfo != nil && len(instance.MultiZoneInfo) > 0 {
			infos := make([]map[string]interface{}, 0, len(instance.MultiZoneInfo))
			for _, v := range instance.MultiZoneInfo {
				info := map[string]interface{}{
					"availability_zone": v.Zone,
					"subnet_id":         v.SubnetId,
				}

				infos = append(infos, info)
			}

			mapping["multi_zone_infos"] = infos
		}

		if instance.NodeInfoList != nil && len(instance.NodeInfoList) > 0 {
			infos := make([]map[string]interface{}, 0, len(instance.NodeInfoList))
			for _, v := range instance.NodeInfoList {
				// this will not keep longer as long as cloud api response update
				if *v.Type == "kibana" {
					continue
				}

				info := map[string]interface{}{
					"node_num":  v.NodeNum,
					"node_type": v.NodeType,
					"type":      v.Type,
					"disk_type": v.DiskType,
					"disk_size": v.DiskSize,
					"encrypt":   *v.DiskEncrypt > 0,
				}

				infos = append(infos, info)
			}

			mapping["node_info_list"] = infos
		}

		instanceList = append(instanceList, mapping)
		ids = append(ids, *instance.InstanceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set elasticsearch instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceList); err != nil {
			return err
		}
	}

	return nil
}
