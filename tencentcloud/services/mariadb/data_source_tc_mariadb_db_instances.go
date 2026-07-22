package mariadb

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbDbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDbInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "实例 ids。",
			},

			"project_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "项目 ids。",
			},

			"search_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例名称 或 VIP",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "私有网络 ID",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "子网 ID",
			},

			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "实例 info。",
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
							Description: "实例名称",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用 可用区",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "meory 的 实例。",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "存储 的 实例。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID",
						},
						"db_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db 版本 ID。",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "内网 IP 地址",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Intranet 端口",
						},
						"internet_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public 网络 访问 域名 名称",
						},
						"internet_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "公网 IP 地址",
						},
						"internet_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Public 网络 端口",
						},
						"resource_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "资源 标签",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
								},
							},
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

func dataSourceTencentCloudMariadbDbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_db_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ids := make([]string, 0)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		var instanceIds []*string
		for i := range instanceIdsSet {
			instanceId := instanceIdsSet[i].(string)
			instanceIds = append(instanceIds, &instanceId)
		}
		paramMap["instance_ids"] = instanceIds
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsSet := v.(*schema.Set).List()
		var projectIds []*int64
		for i := range projectIdsSet {
			projectIds = append(projectIds, helper.IntInt64(projectIdsSet[i].(int)))
		}
		paramMap["project_ids"] = projectIds
	}

	if v, ok := d.GetOk("search_name"); ok {
		paramMap["search_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["vpc_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		paramMap["subnet_id"] = helper.String(v.(string))
	}

	var instances []*mariadb.DBInstance
	mariadbService := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := mariadbService.DescribeMariadbDbInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Mariadb instances failed, reason:%+v", logId, err)
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(instances))
	if instances != nil {
		for _, instance := range instances {
			instanceMap := map[string]interface{}{}
			if instance.InstanceId != nil {
				instanceMap["instance_id"] = instance.InstanceId
			}
			if instance.InstanceName != nil {
				instanceMap["instance_name"] = instance.InstanceName
			}
			if instance.ProjectId != nil {
				instanceMap["project_id"] = *instance.ProjectId
			}
			if instance.Region != nil {
				instanceMap["region"] = instance.Region
			}
			if instance.Zone != nil {
				instanceMap["zone"] = instance.Zone
			}
			if instance.Memory != nil {
				instanceMap["memory"] = instance.Memory
			}
			if instance.Storage != nil {
				instanceMap["storage"] = instance.Storage
			}
			if instance.VpcId != nil {
				instanceMap["vpc_id"] = strconv.Itoa(int(*instance.VpcId))
			}
			if instance.SubnetId != nil {
				instanceMap["subnet_id"] = strconv.Itoa(int(*instance.SubnetId))
			}
			if instance.DbVersionId != nil {
				instanceMap["db_version_id"] = instance.DbVersionId
			}
			if instance.Vip != nil {
				instanceMap["vip"] = instance.Vip
			}
			if instance.Vport != nil {
				instanceMap["vport"] = instance.Vport
			}
			if instance.WanDomain != nil {
				instanceMap["internet_domain"] = instance.WanDomain
			}
			if instance.WanVip != nil {
				instanceMap["internet_ip"] = instance.WanVip
			}
			if instance.WanPort != nil {
				instanceMap["internet_port"] = instance.WanPort
			}
			if instance.ResourceTags != nil {
				resourceTagsList := []interface{}{}
				for _, resourceTags := range instance.ResourceTags {
					resourceTagsMap := map[string]interface{}{}
					if resourceTags.TagKey != nil {
						resourceTagsMap["tag_key"] = resourceTags.TagKey
					}
					if resourceTags.TagValue != nil {
						resourceTagsMap["tag_value"] = resourceTags.TagValue
					}

					resourceTagsList = append(resourceTagsList, resourceTagsMap)
				}
				instanceMap["resource_tags"] = resourceTagsList
			}

			instanceList = append(instanceList, instanceMap)
			ids = append(ids, *instance.InstanceId)
		}
		err = d.Set("instances", instanceList)
		if err != nil {
			log.Printf("[CRITAL]%s provider set instances list fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil
}
