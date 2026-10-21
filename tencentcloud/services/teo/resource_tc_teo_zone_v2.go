package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoZoneV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoZoneV2Create,
		Read:   resourceTencentCloudTeoZoneV2Read,
		Update: resourceTencentCloudTeoZoneV2Update,
		Delete: resourceTencentCloudTeoZoneV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Site access type. Valid values: `partial` (CNAME access), `full` (NS access), `noDomainAccess` (no-domain access), `dnsPodAccess` (DNSPod hosted access), `ai` (edge inference access). Default value: `partial`.",
			},
			"zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Site name. For CNAME/NS access, pass in a second-level domain (example.com) as the site name; for no-domain access, leave this value empty.",
			},
			"area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Acceleration area of the site when Type is partial/full. Valid values: `global`, `mainland`, `overseas`. Default value: `overseas`. Leave it empty when Type is noDomainAccess.",
			},
			"plan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the target plan to be bound. If this parameter is left empty, a site in init status (not activated) will be created, which will not be displayed in the console until it is bound to a plan via BindZoneToPlan.",
			},
			"alias_zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alias zone name, which can contain up to 200 characters, including digits, letters, `.`, `- and `_`. Leave it empty if not needed.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag list of the site. The tags need to be created in the tag console in advance.",
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
			"allow_duplicates": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether to allow repeated access. Default value: `false`.",
			},
			"jump_start": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether to skip scanning the existing DNS records of the site. Default value: `false`.",
			},
			"vanity_name_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Custom name servers of the site, which replace the default name servers assigned by the system. Only effective when the site is not in no-domain access mode.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable custom name servers. Valid values: `on`, `off`.",
						},
						"servers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Custom name server list.",
						},
					},
				},
			},

			"zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site ID.",
			},
			"paused": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the site is closed.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site status. Valid values: `active`, `pending`, `moved`, `deactivated`, `initializing`.",
			},
			"cname_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME access status. Valid values: `finished`, `pending`.",
			},
			"active_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display status. Valid values: `active`, `inactive`, `paused`.",
			},
			"lock_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lock status. Valid values: `enable`, `disable`, `plan_migrate`.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site creation time.",
			},
			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site modification time.",
			},
			"name_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of NS servers assigned by Tencent Cloud.",
			},
			"original_name_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of NS servers currently used by the site.",
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Billing resource list of the site.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing mode. Valid values: `0` (pay-as-you-go).",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"enable_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Effective time.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan status. Valid values: `normal`, `isolated`, `destroyed`.",
						},
						"sv": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inquiry parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Inquiry parameter key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Inquiry parameter value.",
									},
									"pack": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Quota mapped by the inquiry parameter.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Layer-4 proxy instance ID mapped by the inquiry parameter.",
									},
									"protection_specs": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protection level of the inquiry parameter.",
									},
								},
							},
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto renewal flag. Valid values: `0` (default), `1` (auto-renewal), `2` (no auto-renewal).",
						},
						"plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan-associated resource ID.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region. Valid values: `mainland`, `overseas`, `global`.",
						},
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type.",
						},
						"zone_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of sites bound to the current resource.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource tag type.",
						},
					},
				},
			},
			"ns_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details of the NS access type site, which is returned only when Type is full.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cname_speed_up": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to enable CNAME acceleration. Valid values: `enabled`, `disabled`.",
						},
						"is_fake": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether a site with the same name exists. Valid values: `0` (no), `1` (yes).",
						},
						"ownership_verification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Ownership verification info of the NS access type site.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Info required for DNS resolution verification.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subdomain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Host record.",
												},
												"record_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Record type.",
												},
												"record_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Record value.",
												},
											},
										},
									},
									"file_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Info required for file verification.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "URL Path part of the verification file.",
												},
												"content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Content of the verification file.",
												},
											},
										},
									},
									"ns_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Info required to switch DNS servers for NS access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name_servers": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "DNS server address assigned for NS access.",
												},
											},
										},
									},
								},
							},
						},
						"original_name_servers": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of NS servers currently used by the site.",
						},
						"name_servers": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of NS servers assigned by Tencent Cloud EdgeOne.",
						},
						"vanity_name_servers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Custom name server domain info of the user.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to enable custom name servers.",
									},
									"servers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Custom name server list.",
									},
								},
							},
						},
						"vanity_name_servers_ips": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP address info of custom name servers of the user.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom name server name.",
									},
									"i_pv4": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom name server IPv4 address.",
									},
								},
							},
						},
					},
				},
			},
			"cname_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details of the CNAME access type site, which is returned only when Type is partial.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_fake": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is a fake site. Valid values: `0` (no), `1` (yes).",
						},
						"ownership_verification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Ownership verification info of the CNAME access type site.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Info required for DNS resolution verification.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subdomain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Host record.",
												},
												"record_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Record type.",
												},
												"record_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Record value.",
												},
											},
										},
									},
									"file_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Info required for file verification.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "URL Path part of the verification file.",
												},
												"content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Content of the verification file.",
												},
											},
										},
									},
									"ns_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Info required to switch DNS servers for NS access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name_servers": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "DNS server address assigned for NS access.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"dns_pod_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details of the DNSPod hosted access type site, which is returned only when Type is dnsPodAccess.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_fake": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is a fake site. Valid values: `0` (no), `1` (yes).",
						},
					},
				},
			},
			"ownership_verification": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Site ownership verification info. After the site is created, ownership verification must be completed before the site can serve properly.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_verification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Info required for DNS resolution verification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subdomain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Host record.",
									},
									"record_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Record type.",
									},
									"record_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Record value.",
									},
								},
							},
						},
						"file_verification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Info required for file verification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL Path part of the verification file.",
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Content of the verification file.",
									},
								},
							},
						},
						"ns_verification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Info required to switch DNS servers for NS access.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name_servers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "DNS server address assigned for NS access.",
									},
								},
							},
						},
					},
				},
			},
			"vanity_name_servers_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IP address info of custom name servers of the user.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom name server name.",
						},
						"i_pv4": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom name server IPv4 address.",
						},
					},
				},
			},
			"work_mode_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Version management config group working mode of the site.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_group_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Config group type. Valid values: `l7_acceleration`, `edge_functions`, `web_security`.",
						},
						"work_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Working mode. Valid values: `immediate_effect`, `version_control`.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoZoneV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_zone_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewCreateZoneRequest()
	)

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_name"); ok {
		request.ZoneName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plan_id"); ok {
		request.PlanId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alias_zone_name"); ok {
		request.AliasZoneName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsList := v.([]interface{})
		if len(tagsList) > 0 {
			request.Tags = make([]*teov20220901.Tag, 0, len(tagsList))
			for _, item := range tagsList {
				tagMap := item.(map[string]interface{})
				tag := &teov20220901.Tag{}
				if v, ok := tagMap["tag_key"]; ok {
					tag.TagKey = helper.String(v.(string))
				}
				if v, ok := tagMap["tag_value"]; ok {
					tag.TagValue = helper.String(v.(string))
				}
				request.Tags = append(request.Tags, tag)
			}
		}
	}

	if v, ok := d.GetOk("allow_duplicates"); ok {
		request.AllowDuplicates = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("jump_start"); ok {
		request.JumpStart = helper.Bool(v.(bool))
	}

	var response *teov20220901.CreateZoneResponse
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateZoneWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo_zone_v2 failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo_zone_v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create teo_zone_v2, logId: %s, id: %s", logId, logId, d.Id())

	if response.Response.ZoneId == nil || *response.Response.ZoneId == "" {
		return fmt.Errorf("Create teo_zone_v2 failed, ZoneId is empty")
	}

	zoneId := *response.Response.ZoneId

	d.SetId(zoneId)

	return resourceTencentCloudTeoZoneV2Read(d, meta)
}

func resourceTencentCloudTeoZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_zone_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDescribeZonesRequest()
	)

	zoneId := d.Id()

	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("zone-id"),
			Values: []*string{helper.String(zoneId)},
		},
	}

	var response *teov20220901.DescribeZonesResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeZonesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read teo_zone_v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || len(response.Response.Zones) == 0 {
		log.Printf("[CRUD] teo_zone_v2 id=%s", d.Id())
		d.SetId("")
		return nil
	}

	zone := response.Response.Zones[0]

	if zone.ZoneId != nil {
		_ = d.Set("zone_id", zone.ZoneId)
	}

	if zone.ZoneName != nil {
		_ = d.Set("zone_name", zone.ZoneName)
	}

	if zone.AliasZoneName != nil {
		_ = d.Set("alias_zone_name", zone.AliasZoneName)
	}

	if zone.Area != nil {
		_ = d.Set("area", zone.Area)
	}

	if zone.Type != nil {
		_ = d.Set("type", zone.Type)
	}

	if zone.Paused != nil {
		_ = d.Set("paused", zone.Paused)
	}

	if zone.Status != nil {
		_ = d.Set("status", zone.Status)
	}

	if zone.CnameStatus != nil {
		_ = d.Set("cname_status", zone.CnameStatus)
	}

	if zone.ActiveStatus != nil {
		_ = d.Set("active_status", zone.ActiveStatus)
	}

	if zone.LockStatus != nil {
		_ = d.Set("lock_status", zone.LockStatus)
	}

	if zone.CreatedOn != nil {
		_ = d.Set("created_on", zone.CreatedOn)
	}

	if zone.ModifiedOn != nil {
		_ = d.Set("modified_on", zone.ModifiedOn)
	}

	if zone.NameServers != nil {
		_ = d.Set("name_servers", zone.NameServers)
	}

	if zone.OriginalNameServers != nil {
		_ = d.Set("original_name_servers", zone.OriginalNameServers)
	}

	if zone.Tags != nil {
		tagsList := make([]map[string]interface{}, 0, len(zone.Tags))
		for _, tag := range zone.Tags {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = *tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = *tag.TagValue
			}
			tagsList = append(tagsList, tagMap)
		}
		_ = d.Set("tags", tagsList)
	}

	if zone.VanityNameServers != nil {
		vanityNameServersList := make([]map[string]interface{}, 0, 1)
		vanityNameServersMap := map[string]interface{}{}
		if zone.VanityNameServers.Switch != nil {
			vanityNameServersMap["switch"] = *zone.VanityNameServers.Switch
		}
		if zone.VanityNameServers.Servers != nil {
			servers := make([]string, 0, len(zone.VanityNameServers.Servers))
			for _, server := range zone.VanityNameServers.Servers {
				if server != nil {
					servers = append(servers, *server)
				}
			}
			vanityNameServersMap["servers"] = servers
		}
		vanityNameServersList = append(vanityNameServersList, vanityNameServersMap)
		_ = d.Set("vanity_name_servers", vanityNameServersList)
	}

	if zone.Resources != nil {
		resourcesList := make([]map[string]interface{}, 0, len(zone.Resources))
		for _, resource := range zone.Resources {
			resourceMap := map[string]interface{}{}
			if resource.Id != nil {
				resourceMap["id"] = *resource.Id
			}
			if resource.PayMode != nil {
				resourceMap["pay_mode"] = *resource.PayMode
			}
			if resource.CreateTime != nil {
				resourceMap["create_time"] = *resource.CreateTime
			}
			if resource.EnableTime != nil {
				resourceMap["enable_time"] = *resource.EnableTime
			}
			if resource.ExpireTime != nil {
				resourceMap["expire_time"] = *resource.ExpireTime
			}
			if resource.Status != nil {
				resourceMap["status"] = *resource.Status
			}
			if resource.Sv != nil {
				svList := make([]map[string]interface{}, 0, len(resource.Sv))
				for _, sv := range resource.Sv {
					svMap := map[string]interface{}{}
					if sv.Key != nil {
						svMap["key"] = *sv.Key
					}
					if sv.Value != nil {
						svMap["value"] = *sv.Value
					}
					if sv.Pack != nil {
						svMap["pack"] = *sv.Pack
					}
					if sv.InstanceId != nil {
						svMap["instance_id"] = *sv.InstanceId
					}
					if sv.ProtectionSpecs != nil {
						svMap["protection_specs"] = *sv.ProtectionSpecs
					}
					svList = append(svList, svMap)
				}
				resourceMap["sv"] = svList
			}
			if resource.AutoRenewFlag != nil {
				resourceMap["auto_renew_flag"] = *resource.AutoRenewFlag
			}
			if resource.PlanId != nil {
				resourceMap["plan_id"] = *resource.PlanId
			}
			if resource.Area != nil {
				resourceMap["area"] = *resource.Area
			}
			if resource.Group != nil {
				resourceMap["group"] = *resource.Group
			}
			if resource.ZoneNumber != nil {
				resourceMap["zone_number"] = *resource.ZoneNumber
			}
			if resource.Type != nil {
				resourceMap["type"] = *resource.Type
			}
			resourcesList = append(resourcesList, resourceMap)
		}
		_ = d.Set("resources", resourcesList)
	}

	if zone.NSDetail != nil {
		nsDetailList := make([]map[string]interface{}, 0, 1)
		nsDetailMap := map[string]interface{}{}
		if zone.NSDetail.CnameSpeedUp != nil {
			nsDetailMap["cname_speed_up"] = *zone.NSDetail.CnameSpeedUp
		}
		if zone.NSDetail.IsFake != nil {
			nsDetailMap["is_fake"] = *zone.NSDetail.IsFake
		}
		if zone.NSDetail.OwnershipVerification != nil {
			nsDetailMap["ownership_verification"] = flattenTeoZoneV2OwnershipVerification(zone.NSDetail.OwnershipVerification)
		}
		if zone.NSDetail.OriginalNameServers != nil {
			nsDetailMap["original_name_servers"] = flattenTeoZoneV2StringList(zone.NSDetail.OriginalNameServers)
		}
		if zone.NSDetail.NameServers != nil {
			nsDetailMap["name_servers"] = flattenTeoZoneV2StringList(zone.NSDetail.NameServers)
		}
		if zone.NSDetail.VanityNameServers != nil {
			vanityNameServersMap := map[string]interface{}{}
			if zone.NSDetail.VanityNameServers.Switch != nil {
				vanityNameServersMap["switch"] = *zone.NSDetail.VanityNameServers.Switch
			}
			if zone.NSDetail.VanityNameServers.Servers != nil {
				vanityNameServersMap["servers"] = flattenTeoZoneV2StringList(zone.NSDetail.VanityNameServers.Servers)
			}
			nsDetailMap["vanity_name_servers"] = []interface{}{vanityNameServersMap}
		}
		if zone.NSDetail.VanityNameServersIps != nil {
			vanityNameServersIpsList := make([]map[string]interface{}, 0, len(zone.NSDetail.VanityNameServersIps))
			for _, ip := range zone.NSDetail.VanityNameServersIps {
				ipMap := map[string]interface{}{}
				if ip.Name != nil {
					ipMap["name"] = *ip.Name
				}
				if ip.IPv4 != nil {
					ipMap["i_pv4"] = *ip.IPv4
				}
				vanityNameServersIpsList = append(vanityNameServersIpsList, ipMap)
			}
			nsDetailMap["vanity_name_servers_ips"] = vanityNameServersIpsList
		}
		nsDetailList = append(nsDetailList, nsDetailMap)
		_ = d.Set("ns_detail", nsDetailList)
	}

	if zone.CNAMEDetail != nil {
		cnameDetailList := make([]map[string]interface{}, 0, 1)
		cnameDetailMap := map[string]interface{}{}
		if zone.CNAMEDetail.IsFake != nil {
			cnameDetailMap["is_fake"] = *zone.CNAMEDetail.IsFake
		}
		if zone.CNAMEDetail.OwnershipVerification != nil {
			cnameDetailMap["ownership_verification"] = flattenTeoZoneV2OwnershipVerification(zone.CNAMEDetail.OwnershipVerification)
		}
		cnameDetailList = append(cnameDetailList, cnameDetailMap)
		_ = d.Set("cname_detail", cnameDetailList)
	}

	if zone.DNSPodDetail != nil {
		dnsPodDetailList := make([]map[string]interface{}, 0, 1)
		dnsPodDetailMap := map[string]interface{}{}
		if zone.DNSPodDetail.IsFake != nil {
			dnsPodDetailMap["is_fake"] = *zone.DNSPodDetail.IsFake
		}
		dnsPodDetailList = append(dnsPodDetailList, dnsPodDetailMap)
		_ = d.Set("dns_pod_detail", dnsPodDetailList)
	}

	if zone.OwnershipVerification != nil {
		_ = d.Set("ownership_verification", flattenTeoZoneV2OwnershipVerification(zone.OwnershipVerification))
	}

	if zone.VanityNameServersIps != nil {
		vanityNameServersIpsList := make([]map[string]interface{}, 0, len(zone.VanityNameServersIps))
		for _, ip := range zone.VanityNameServersIps {
			ipMap := map[string]interface{}{}
			if ip.Name != nil {
				ipMap["name"] = *ip.Name
			}
			if ip.IPv4 != nil {
				ipMap["i_pv4"] = *ip.IPv4
			}
			vanityNameServersIpsList = append(vanityNameServersIpsList, ipMap)
		}
		_ = d.Set("vanity_name_servers_ips", vanityNameServersIpsList)
	}

	if zone.WorkModeInfos != nil {
		workModeInfosList := make([]map[string]interface{}, 0, len(zone.WorkModeInfos))
		for _, workModeInfo := range zone.WorkModeInfos {
			workModeInfoMap := map[string]interface{}{}
			if workModeInfo.ConfigGroupType != nil {
				workModeInfoMap["config_group_type"] = *workModeInfo.ConfigGroupType
			}
			if workModeInfo.WorkMode != nil {
				workModeInfoMap["work_mode"] = *workModeInfo.WorkMode
			}
			workModeInfosList = append(workModeInfosList, workModeInfoMap)
		}
		_ = d.Set("work_mode_infos", workModeInfosList)
	}

	return nil
}

func resourceTencentCloudTeoZoneV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_zone_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewModifyZoneRequest()
	)

	zoneId := d.Id()
	request.ZoneId = helper.String(zoneId)

	needChange := false
	mutableArgs := []string{"type", "zone_name", "area", "alias_zone_name", "vanity_name_servers"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}

		if v, ok := d.GetOk("zone_name"); ok {
			request.ZoneName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("area"); ok {
			request.Area = helper.String(v.(string))
		}

		if v, ok := d.GetOk("alias_zone_name"); ok {
			request.AliasZoneName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("vanity_name_servers"); ok {
			vanityNameServersList := v.([]interface{})
			if len(vanityNameServersList) > 0 {
				vanityNameServersMap := vanityNameServersList[0].(map[string]interface{})
				vanityNameServers := &teov20220901.VanityNameServers{}
				if v, ok := vanityNameServersMap["switch"]; ok {
					vanityNameServers.Switch = helper.String(v.(string))
				}
				if v, ok := vanityNameServersMap["servers"]; ok {
					serversList := v.([]interface{})
					if len(serversList) > 0 {
						vanityNameServers.Servers = make([]*string, 0, len(serversList))
						for _, server := range serversList {
							vanityNameServers.Servers = append(vanityNameServers.Servers, helper.String(server.(string)))
						}
					}
				}
				request.VanityNameServers = vanityNameServers
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyZoneWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update teo_zone_v2 failed, Response is nil"))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo_zone_v2 failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoZoneV2Read(d, meta)
}

func resourceTencentCloudTeoZoneV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_zone_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteZoneRequest()
	)

	zoneId := d.Id()
	request.ZoneId = helper.String(zoneId)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteZoneWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo_zone_v2 failed, Response is nil"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo_zone_v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func flattenTeoZoneV2OwnershipVerification(ownershipVerification *teov20220901.OwnershipVerification) []interface{} {
	if ownershipVerification == nil {
		return []interface{}{}
	}

	ownershipVerificationMap := map[string]interface{}{}

	if ownershipVerification.DnsVerification != nil {
		dnsVerificationMap := map[string]interface{}{}
		if ownershipVerification.DnsVerification.Subdomain != nil {
			dnsVerificationMap["subdomain"] = *ownershipVerification.DnsVerification.Subdomain
		}
		if ownershipVerification.DnsVerification.RecordType != nil {
			dnsVerificationMap["record_type"] = *ownershipVerification.DnsVerification.RecordType
		}
		if ownershipVerification.DnsVerification.RecordValue != nil {
			dnsVerificationMap["record_value"] = *ownershipVerification.DnsVerification.RecordValue
		}
		ownershipVerificationMap["dns_verification"] = []interface{}{dnsVerificationMap}
	}

	if ownershipVerification.FileVerification != nil {
		fileVerificationMap := map[string]interface{}{}
		if ownershipVerification.FileVerification.Path != nil {
			fileVerificationMap["path"] = *ownershipVerification.FileVerification.Path
		}
		if ownershipVerification.FileVerification.Content != nil {
			fileVerificationMap["content"] = *ownershipVerification.FileVerification.Content
		}
		ownershipVerificationMap["file_verification"] = []interface{}{fileVerificationMap}
	}

	if ownershipVerification.NsVerification != nil {
		nsVerificationMap := map[string]interface{}{}
		if ownershipVerification.NsVerification.NameServers != nil {
			nsVerificationMap["name_servers"] = flattenTeoZoneV2StringList(ownershipVerification.NsVerification.NameServers)
		}
		ownershipVerificationMap["ns_verification"] = []interface{}{nsVerificationMap}
	}

	return []interface{}{ownershipVerificationMap}
}

func flattenTeoZoneV2StringList(stringList []*string) []string {
	result := make([]string, 0, len(stringList))
	for _, item := range stringList {
		if item != nil {
			result = append(result, *item)
		}
	}
	return result
}
