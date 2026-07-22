package postgresql

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlInstanceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 postgresql 实例 到 是 查询。",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID postgresql 实例 到 是 查询。",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "项目 ID postgresql 实例 到 是 查询。",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Deprecated:  "It has been deprecated from version 1.82.64. Please use `db_instance_set` instead.",
				Description: "A 列表 postgresql 实例. Each element 包含following attributes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID postgresql 实例。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 postgresql 实例。",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay 类型 postgresql 实例。",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto 续费标识",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 的 postgresql 数据库 引擎。",
						},
						"db_kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PostgreSQL kernel 版本 数量。",
						},
						"db_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PostgreSQL major 版本 数量。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 子网。",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume 大小(在 GB)。",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage 类型 有效值：`PHYSICAL_LOCAL_SSD` (本地 SSD)，`CLOUD_PREMIUM` (premium 云 磁盘)，`CLOUD_SSD` (云 SSD)，`CLOUD_HSSD` (enhanced 云 SSD)。",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小(在 GB)。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID，默认值为 0。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability 可用区",
						},
						"root_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 root 账号 名称，默认值为 `root`。",
						},
						"public_access_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否enable 访问 到 实例 从 公有 网络 或 不。",
						},
						"public_access_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 对于 公有 访问。",
						},
						"public_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "端口 对于 公有 访问。",
						},
						"private_access_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 地址 对于 私有 访问。",
						},
						"private_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "端口 对于 私有 访问。",
						},
						"charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charset 的 postgresql 实例。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 postgresql 实例。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "可用 标签 within 此 postgresql。",
						},
					},
				},
			},
			"db_instance_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "实例 details 集合。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 地域 such 作为 ap-guangzhou，其中 corresponds 到 `地域` 字段 在 `RegionSet`。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 AZ such 作为 ap-guangzhou-3，其中 corresponds 到 `可用区` 字段 的 `ZoneSet`。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpc ID，such 作为 vpc-e6w23k31. 有效 vpc ID 可以 是 获取 通过 日志记录 在 到 console 到 查询 或 通过 calling api [DescribeVpcs](https://www.tencentcloud.comom/document/api/215/15778?from_cn_redirect=1) 和 acquiring unVpcId 字段 在 api 返回。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC 子网 ID，such 作为 子网-51lcif9y. effective VPC 子网 ids 可以 是 获取 通过 日志记录 在 到 console 或 calling api [DescribeSubnets](https://www.tencentcloud.comom/document/api/215/15784?from_cn_redirect=1) 到 acquire unSubnetId 字段 在 api 返回。",
						},
						"db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"db_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称",
						},
						"db_instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例状态，包括: `applying` (applying)，`init` (到 是 initialized)，`initing` (initializing)，`running` (running)，`limited run` (restricted operation)，`isolating` (isolating)，`isolated` (isolated)，`disisolating` (de-isolating)，`recycling` (recycling)，`recycled` (recycled)，`作业 running` (任务 executing)，`offline` (offline)，`migrating` (migrating)，`expanding` (scaling out)，`waitSwitch` (waiting 到 switch)，`switching` (switching)，`readonly` (readonly)，`restarting` (restarting)，`网络 changing` (网络 modification 在 progress)，`upgrading` (kernel 版本 upgrading)，`audit-switching` (审计状态 changing)，`primary-switching` (primary-secondary switching)，`offlining` (offline)，`部署 changing` (modify az)，`cloning` (restoring 数据)，`参数 modifying` (参数 modification 在 progress)，`日志-switching` (日志 状态 change)，`restoring` (recovering)，和 `expanding` (scaling out)。",
						},
						"db_instance_memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Assigned 实例 内存 大小 （GB）。",
						},
						"db_instance_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Assigned 实例 存储 容量 （GB）。",
						},
						"db_instance_storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage 类型 有效值：`PHYSICAL_LOCAL_SSD` (本地 SSD)，`CLOUD_PREMIUM` (premium 云 磁盘)，`CLOUD_SSD` (云 SSD)，`CLOUD_HSSD` (enhanced 云 SSD)。",
						},
						"db_instance_cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 assigned CPUs。",
						},
						"db_instance_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purchasable 规格 ID。",
						},
						"db_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PostgreSQL major 版本 数量. 指定version 信息 该 可以 是 获取 从 [DescribeDBVersions](https://www.tencentcloud.comom/document/api/409/89018?from_cn_redirect=1) api. currently 支持 major versions 10，11，12，13，14，15。",
						},
						"db_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数量 major PostgreSQL community 版本 和 minor 版本，such 作为 12.4，其中 可以 是 queried 通过 [DescribeDBVersions](https://intl.云.tencent.com/document/api/409/89018?from_cn_redirect=1) API。",
						},
						"db_kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PostgreSQL kernel 版本，such 作为 v12.7_r1.8. 版本 信息 可以 是 获取 从 [DescribeDBVersions](https://www.tencentcloud.comom/document/api/409/89018?from_cn_redirect=1)。",
						},
						"db_instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例类型，其中 includes:\n<li>primary: primary 实例 </li>\n<li>readonly: read-仅 实例</li>\n<li>guard: disaster recovery 实例</li>\n<li>temp: temporary 实例</li>。",
						},
						"db_instance_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 版本 有效 值: `standard` (dual-服务器 high-availability; 一个-primary-一个-standby)。",
						},
						"db_charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 character 集合，其中 currently 支持 仅:\n<li>UTF8</li>\n<li>LATIN1</li>。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last 更新 时间 的 实例 attribute。",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 过期时间。",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 isolation 时间。",
						},
						"pay_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing 模式:\n<li>prepaid: monthly subscription，prepaid</li>\n<li>postpaid: pay-作为-您-go，postpaid</li>。",
						},
						"auto_renew": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto-renewal 或 不:\n<li>`0`: manual renewal</li>\n<li>`1`: auto-renewal</li>\n默认值：0。",
						},
						"db_instance_net_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 网络 连接 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS 域名 名称",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ip。",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Connection 端口 地址",
									},
									"net_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 类型 1: inner (私有 网络 地址)，2: 公有 (公有 网络 地址)。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 连接 状态 有效值：`initing` (never 已启用 before)，`opened` (已启用)，`closed` (已禁用)，`opening` (enabling)，`closing` (disabling)。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有网络 ID 指定ID 的 virtual 私有 云。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID",
									},
									"protocol_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "指定protocol 类型 到 connect 到 数据库. currently 支持: postgresql，mssql (mssql compatible syntax)。",
									},
								},
							},
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine 类型",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "用户 `AppId`。",
						},
						"uid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 `Uid`。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"tag_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Describes 标签 信息 associated 使用 实例。",
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
						"master_db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Primary 实例 信息. 返回 仅 当 实例 是 read-仅 实例。",
						},
						"read_only_instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "指定number 的 read-仅 实例。",
						},
						"status_in_readonly_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describes state 的 read-仅 实例 在 read-仅 组。",
						},
						"offline_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Decommissioning 时间。",
						},
						"db_node_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 节点 信息\n注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node 类型 有效 值:\n`Primary`;\n`Standby`。",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AZ 其中 节点 resides，such 作为 ap-guangzhou-1。",
									},
									"dedicated_cluster_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CDC ID。",
									},
								},
							},
						},
						"is_support_tde": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否instance 支持 TDE 数据 加密.\n<Li>0: 不 支持</li>.\n<Li>1: 支持.</li>.\n默认值：0\n\nFor TDE 数据 加密，see [overview 的 数据 transparent 加密](https://www.tencentcloud.comom/document/product/409/71748?from_cn_redirect=1)。",
						},
						"db_engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database 引擎，其中 支持:\n<li>`postgresql`: tencentdb 对于 postgresql</li>.\n<li>`mssql_compatible`: 指定mssql compatible - tencentdb 对于 PostgreSQL.</li>.\n默认值：`postgresql`。",
						},
						"db_engine_config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration 信息 对于 数据库 引擎，和 配置 格式 是 作为 follows:.\n{`$key1`:`$value1`，`$key2`:`$value2`}\nSupported engines include:.\nmssql_compatible 引擎:.\n<li>migrationMode: 指定database 模式 可选 参数. 有效值：单个-db (单个-数据库 schema) 和 multi-db (多个 数据库 schemas). 默认为 单个-db.</li>.\n<li>defaultLocale: 指定sorting area 规则， 可选 参数 该 不能 是 modified after initialization. 默认值为 en_US. 有效 值 include:.\n`af_ZA`，`sq_AL`，`ar_DZ`，`ar_BH`，`ar_EG`，`ar_IQ`，`ar_JO`，`ar_KW`，`ar_LB`，`ar_LY`，`ar_MA`，`ar_OM`，`ar_QA`，`ar_SA`，`ar_SY`，`ar_TN`，`ar_AE`，`ar_YE`，`hy_AM`，`az_Cyrl_AZ`，`az_Latn_AZ`，`eu_ES`，`be_BY`，`bg_BG`，`ca_ES`，`zh_HK`，`zh_MO`，`zh_CN`，`zh_SG`，`zh_TW`，`hr_HR`，`cs_CZ`，`da_DK`，`nl_BE`，`nl_NL`，`en_AU`，`en_BZ`，`en_CA`，`en_IE`，`en_JM`，`en_NZ`，`en_PH`，`en_ZA`，`en_TT`，`en_GB`，`en_US`，`en_ZW`，`et_EE`，`fo_FO`，`fa_IR`，`fi_FI`，`fr_BE`，`fr_CA`，`fr_FR`，`fr_LU`，`fr_MC`，`fr_CH`，`mk_MK`，`ka_GE`，`de_AT`，`de_DE`，`de_LI`，`de_LU`，`de_CH`，`el_GR`，`gu_IN`，`he_IL`，`hi_IN`，`hu_HU`，`is_IS`，`id_ID`，`it_IT`，`it_CH`，`ja_JP`，`kn_IN`，`kok_IN`，`ko_KR`，`ky_KG`，`lv_LV`，`lt_LT`，`ms_BN`，`ms_MY`，`mr_IN`，`mn_MN`，`nb_NO`，`nn_NO`，`pl_PL`，`pt_BR`，`pt_PT`，`pa_IN`，`ro_RO`，`ru_RU`，`sa_IN`，`sr_Cyrl_RS`，`sr_Latn_RS`，`sk_SK`，`sl_SI`，`es_AR`，`es_BO`，`es_CL`，`es_CO`，`es_CR`，`es_DO`，`es_EC`，`es_SV`，`es_GT`，`es_HN`，`es_MX`，`es_NI`，`es_PA`，`es_PY`,`es_PE`，`es_PR`，`es_ES`，`es_TRADITIONAL`，`es_UY`，`es_VE`，`sw_KE`，`sv_FI`，`sv_SE`，`tt_RU`，`te_IN`，`th_TH`，`tr_TR`，`uk_UA`，`ur_IN`，`ur_PK`，`uz_Cyrl_UZ`，`uz_Latn_UZ`，`vi_VN`.</li>\n<li>serverCollationName: Sorting 规则 名称， 可选 参数，其中 不能 是 modified after initialization，its 默认值为 sql_latin1_general_cp1_ci_as，和 its 有效 值 include: `bbf_unicode_general_ci_as`，`bbf_unicode_cp1_ci_as`，`bbf_unicode_CP1250_ci_as`，`bbf_unicode_CP1251_ci_as`，`bbf_unicode_cp1253_ci_as`，`bbf_unicode_cp1254_ci_as`，`bbf_unicode_cp1255_ci_as`，`bbf_unicode_cp1256_ci_as`，`bbf_unicode_cp1257_ci_as`，`bbf_unicode_cp1258_ci_as`，`bbf_unicode_cp874_ci_as`，`sql_latin1_general_cp1250_ci_as`，`sql_latin1_general_cp1251_ci_as`，`sql_latin1_general_cp1_ci_as`，`sql_latin1_general_cp1253_ci_as`，`sql_latin1_general_cp1254_ci_as`，`sql_latin1_general_cp1255_ci_as`，`sql_latin1_general_cp1256_ci_as`，`sql_latin1_general_cp1257_ci_as`，`sql_latin1_general_cp1258_ci_as`，`chinese_prc_ci_as`，`cyrillic_general_ci_as`，`finnish_swedish_ci_as`，`french_ci_as`，`japanese_ci_as`，`korean_wansung_ci_as`，`latin1_general_ci_as`，`modern_spanish_ci_as`，`polish_ci_as`，`thai_ci_as`，`traditional_spanish_ci_as`，`turkish_ci_as`，`ukrainian_ci_as`，和 `vietnamese_ci_as`.</li>。",
						},
						"network_access_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Network 访问 列表 实例 (此 字段 has been 已弃用)\nNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 资源 ID，实例 ID，或 RO 组 ID。",
									},
									"resource_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "资源类型 有效值：1 (实例)，2 (RO 组)。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有网络 ID 指定ID 的 virtual 私有 云。",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPv4 地址",
									},
									"vip6": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPv6 地址",
									},
									"vport": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "指定access 端口",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID",
									},
									"vpc_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network 状态 有效值：1-applying，2-活跃，3-deleting，4-删除。",
									},
								},
							},
						},
						"support_ipv6": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否instance 支持 IPv6:\n<li>`0`: 无</li>\n<li>`1`: yes</li>\n默认值：0。",
						},
						"expanded_cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 CPU 核数 该 have been elastically scaled out。",
						},
						"deletion_protection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "指定是否enable 删除保护 对于 实例. 有效 值 作为 follows:.\n-指定是否enable 删除保护 有效值：true (启用 删除保护).\n-指定是否disable 删除保护 有效值：false (disable 删除保护)。",
						},
						"root_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 root 账号 名称，默认值为 `root`。",
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

func dataSourceTencentCloudPostgresqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_instances.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	filter := make([]*postgresql.Filter, 0)
	if v, ok := d.GetOk("name"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-instance-name"), Values: []*string{helper.String(v.(string))}})
	}

	if v, ok := d.GetOk("id"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-instance-id"), Values: []*string{helper.String(v.(string))}})
	}

	if v, ok := d.GetOk("project_id"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-project-id"), Values: []*string{helper.String(v.(string))}})
	}

	instanceList, err := service.DescribePostgresqlInstances(ctx, filter)
	if err != nil {
		instanceList, err = service.DescribePostgresqlInstances(ctx, filter)
	}

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	list := make([]map[string]interface{}, 0, len(instanceList))
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)

	// old
	for _, v := range instanceList {
		listItem := make(map[string]interface{})
		listItem["id"] = v.DBInstanceId
		listItem["name"] = v.DBInstanceName
		listItem["auto_renew_flag"] = v.AutoRenew
		listItem["project_id"] = v.ProjectId
		listItem["storage"] = v.DBInstanceStorage
		if v.DBInstanceStorageType != nil {
			listItem["storage_type"] = v.DBInstanceStorageType
		}
		listItem["memory"] = v.DBInstanceMemory
		listItem["availability_zone"] = v.Zone
		listItem["create_time"] = v.CreateTime
		listItem["vpc_id"] = v.VpcId
		listItem["subnet_id"] = v.SubnetId
		listItem["engine_version"] = v.DBVersion
		if v.DBKernelVersion != nil {
			listItem["db_kernel_version"] = v.DBKernelVersion
		}

		if v.DBMajorVersion != nil {
			listItem["db_major_version"] = v.DBMajorVersion
		}

		listItem["public_access_switch"] = false
		listItem["charset"] = v.DBCharset
		listItem["public_access_host"] = ""

		// rootUser
		if v.DBInstanceId != nil && strings.HasPrefix(*v.DBInstanceId, "postgres-") {
			accounts, _ := service.DescribeRootUser(ctx, *v.DBInstanceId)
			if len(accounts) > 0 {
				listItem["root_user"] = accounts[0].UserName
			}
		}

		for _, netInfo := range v.DBInstanceNetInfo {
			if *netInfo.NetType == "public" {
				if *netInfo.Status == "opened" || *netInfo.Status == "1" {
					listItem["public_access_switch"] = true
				}
				listItem["public_access_host"] = netInfo.Address
				listItem["public_access_port"] = netInfo.Port
			}
			if (*netInfo.NetType == "private" || *netInfo.NetType == "inner") && *netInfo.Ip != "" {
				listItem["private_access_ip"] = netInfo.Ip
				listItem["private_access_port"] = netInfo.Port
			}
		}

		if *v.PayType == POSTGRESQL_PAYTYPE_PREPAID || *v.PayType == COMMON_PAYTYPE_PREPAID {
			listItem["charge_type"] = COMMON_PAYTYPE_PREPAID
		} else {
			listItem["charge_type"] = COMMON_PAYTYPE_POSTPAID
		}

		//the describe list API is delayed with argument `tag`
		tagList, err := tagService.DescribeResourceTags(ctx, "postgres", "DBInstanceId", tcClient.Region, *v.DBInstanceId)
		if err != nil {
			return err
		}

		listItem["tags"] = tagList

		list = append(list, listItem)
		ids = append(ids, *v.DBInstanceId)
	}

	// new
	dBInstanceSetList := make([]map[string]interface{}, 0, len(instanceList))
	for _, dBInstanceSet := range instanceList {
		dBInstanceSetMap := map[string]interface{}{}
		if dBInstanceSet.Region != nil {
			dBInstanceSetMap["region"] = dBInstanceSet.Region
		}

		if dBInstanceSet.Zone != nil {
			dBInstanceSetMap["zone"] = dBInstanceSet.Zone
		}

		if dBInstanceSet.VpcId != nil {
			dBInstanceSetMap["vpc_id"] = dBInstanceSet.VpcId
		}

		if dBInstanceSet.SubnetId != nil {
			dBInstanceSetMap["subnet_id"] = dBInstanceSet.SubnetId
		}

		if dBInstanceSet.DBInstanceId != nil {
			dBInstanceSetMap["db_instance_id"] = dBInstanceSet.DBInstanceId
		}

		if dBInstanceSet.DBInstanceName != nil {
			dBInstanceSetMap["db_instance_name"] = dBInstanceSet.DBInstanceName
		}

		if dBInstanceSet.DBInstanceStatus != nil {
			dBInstanceSetMap["db_instance_status"] = dBInstanceSet.DBInstanceStatus
		}

		if dBInstanceSet.DBInstanceMemory != nil {
			dBInstanceSetMap["db_instance_memory"] = dBInstanceSet.DBInstanceMemory
		}

		if dBInstanceSet.DBInstanceStorage != nil {
			dBInstanceSetMap["db_instance_storage"] = dBInstanceSet.DBInstanceStorage
		}

		if dBInstanceSet.DBInstanceStorageType != nil {
			dBInstanceSetMap["db_instance_storage_type"] = dBInstanceSet.DBInstanceStorageType
		}

		if dBInstanceSet.DBInstanceCpu != nil {
			dBInstanceSetMap["db_instance_cpu"] = dBInstanceSet.DBInstanceCpu
		}

		if dBInstanceSet.DBInstanceClass != nil {
			dBInstanceSetMap["db_instance_class"] = dBInstanceSet.DBInstanceClass
		}

		if dBInstanceSet.DBMajorVersion != nil {
			dBInstanceSetMap["db_major_version"] = dBInstanceSet.DBMajorVersion
		}

		if dBInstanceSet.DBVersion != nil {
			dBInstanceSetMap["db_version"] = dBInstanceSet.DBVersion
		}

		if dBInstanceSet.DBKernelVersion != nil {
			dBInstanceSetMap["db_kernel_version"] = dBInstanceSet.DBKernelVersion
		}

		if dBInstanceSet.DBInstanceType != nil {
			dBInstanceSetMap["db_instance_type"] = dBInstanceSet.DBInstanceType
		}

		if dBInstanceSet.DBInstanceVersion != nil {
			dBInstanceSetMap["db_instance_version"] = dBInstanceSet.DBInstanceVersion
		}

		if dBInstanceSet.DBCharset != nil {
			dBInstanceSetMap["db_charset"] = dBInstanceSet.DBCharset
		}

		if dBInstanceSet.CreateTime != nil {
			dBInstanceSetMap["create_time"] = dBInstanceSet.CreateTime
		}

		if dBInstanceSet.UpdateTime != nil {
			dBInstanceSetMap["update_time"] = dBInstanceSet.UpdateTime
		}

		if dBInstanceSet.ExpireTime != nil {
			dBInstanceSetMap["expire_time"] = dBInstanceSet.ExpireTime
		}

		if dBInstanceSet.IsolatedTime != nil {
			dBInstanceSetMap["isolated_time"] = dBInstanceSet.IsolatedTime
		}

		if dBInstanceSet.PayType != nil {
			dBInstanceSetMap["pay_type"] = dBInstanceSet.PayType
		}

		if dBInstanceSet.AutoRenew != nil {
			dBInstanceSetMap["auto_renew"] = dBInstanceSet.AutoRenew
		}

		dBInstanceNetInfoList := make([]map[string]interface{}, 0, len(dBInstanceSet.DBInstanceNetInfo))
		if dBInstanceSet.DBInstanceNetInfo != nil {
			for _, dBInstanceNetInfo := range dBInstanceSet.DBInstanceNetInfo {
				dBInstanceNetInfoMap := map[string]interface{}{}
				if dBInstanceNetInfo.Address != nil {
					dBInstanceNetInfoMap["address"] = dBInstanceNetInfo.Address
				}

				if dBInstanceNetInfo.Ip != nil {
					dBInstanceNetInfoMap["ip"] = dBInstanceNetInfo.Ip
				}

				if dBInstanceNetInfo.Port != nil {
					dBInstanceNetInfoMap["port"] = dBInstanceNetInfo.Port
				}

				if dBInstanceNetInfo.NetType != nil {
					dBInstanceNetInfoMap["net_type"] = dBInstanceNetInfo.NetType
				}

				if dBInstanceNetInfo.Status != nil {
					dBInstanceNetInfoMap["status"] = dBInstanceNetInfo.Status
				}

				if dBInstanceNetInfo.VpcId != nil {
					dBInstanceNetInfoMap["vpc_id"] = dBInstanceNetInfo.VpcId
				}

				if dBInstanceNetInfo.SubnetId != nil {
					dBInstanceNetInfoMap["subnet_id"] = dBInstanceNetInfo.SubnetId
				}

				if dBInstanceNetInfo.ProtocolType != nil {
					dBInstanceNetInfoMap["protocol_type"] = dBInstanceNetInfo.ProtocolType
				}

				dBInstanceNetInfoList = append(dBInstanceNetInfoList, dBInstanceNetInfoMap)
			}

			dBInstanceSetMap["db_instance_net_info"] = dBInstanceNetInfoList
		}

		if dBInstanceSet.Type != nil {
			dBInstanceSetMap["type"] = dBInstanceSet.Type
		}

		if dBInstanceSet.AppId != nil {
			dBInstanceSetMap["app_id"] = dBInstanceSet.AppId
		}

		if dBInstanceSet.Uid != nil {
			dBInstanceSetMap["uid"] = dBInstanceSet.Uid
		}

		if dBInstanceSet.ProjectId != nil {
			dBInstanceSetMap["project_id"] = dBInstanceSet.ProjectId
		}

		tagListList := make([]map[string]interface{}, 0, len(dBInstanceSet.TagList))
		if dBInstanceSet.TagList != nil {
			for _, tagList := range dBInstanceSet.TagList {
				tagListMap := map[string]interface{}{}
				if tagList.TagKey != nil {
					tagListMap["tag_key"] = tagList.TagKey
				}

				if tagList.TagValue != nil {
					tagListMap["tag_value"] = tagList.TagValue
				}

				tagListList = append(tagListList, tagListMap)
			}

			dBInstanceSetMap["tag_list"] = tagListList
		}

		if dBInstanceSet.MasterDBInstanceId != nil {
			dBInstanceSetMap["master_db_instance_id"] = dBInstanceSet.MasterDBInstanceId
		}

		if dBInstanceSet.ReadOnlyInstanceNum != nil {
			dBInstanceSetMap["read_only_instance_num"] = dBInstanceSet.ReadOnlyInstanceNum
		}

		if dBInstanceSet.StatusInReadonlyGroup != nil {
			dBInstanceSetMap["status_in_readonly_group"] = dBInstanceSet.StatusInReadonlyGroup
		}

		if dBInstanceSet.OfflineTime != nil {
			dBInstanceSetMap["offline_time"] = dBInstanceSet.OfflineTime
		}

		dBNodeSetList := make([]map[string]interface{}, 0, len(dBInstanceSet.DBNodeSet))
		if dBInstanceSet.DBNodeSet != nil {
			for _, dBNodeSet := range dBInstanceSet.DBNodeSet {
				dBNodeSetMap := map[string]interface{}{}
				if dBNodeSet.Role != nil {
					dBNodeSetMap["role"] = dBNodeSet.Role
				}

				if dBNodeSet.Zone != nil {
					dBNodeSetMap["zone"] = dBNodeSet.Zone
				}

				if dBNodeSet.DedicatedClusterId != nil {
					dBNodeSetMap["dedicated_cluster_id"] = dBNodeSet.DedicatedClusterId
				}

				dBNodeSetList = append(dBNodeSetList, dBNodeSetMap)
			}

			dBInstanceSetMap["db_node_set"] = dBNodeSetList
		}

		if dBInstanceSet.IsSupportTDE != nil {
			dBInstanceSetMap["is_support_tde"] = dBInstanceSet.IsSupportTDE
		}

		if dBInstanceSet.DBEngine != nil {
			dBInstanceSetMap["db_engine"] = dBInstanceSet.DBEngine
		}

		if dBInstanceSet.DBEngineConfig != nil {
			dBInstanceSetMap["db_engine_config"] = dBInstanceSet.DBEngineConfig
		}

		networkAccessListList := make([]map[string]interface{}, 0, len(dBInstanceSet.NetworkAccessList))
		if dBInstanceSet.NetworkAccessList != nil {
			for _, networkAccessList := range dBInstanceSet.NetworkAccessList {
				networkAccessListMap := map[string]interface{}{}
				if networkAccessList.ResourceId != nil {
					networkAccessListMap["resource_id"] = networkAccessList.ResourceId
				}

				if networkAccessList.ResourceType != nil {
					networkAccessListMap["resource_type"] = networkAccessList.ResourceType
				}

				if networkAccessList.VpcId != nil {
					networkAccessListMap["vpc_id"] = networkAccessList.VpcId
				}

				if networkAccessList.Vip != nil {
					networkAccessListMap["vip"] = networkAccessList.Vip
				}

				if networkAccessList.Vip6 != nil {
					networkAccessListMap["vip6"] = networkAccessList.Vip6
				}

				if networkAccessList.Vport != nil {
					networkAccessListMap["vport"] = networkAccessList.Vport
				}

				if networkAccessList.SubnetId != nil {
					networkAccessListMap["subnet_id"] = networkAccessList.SubnetId
				}

				if networkAccessList.VpcStatus != nil {
					networkAccessListMap["vpc_status"] = networkAccessList.VpcStatus
				}

				networkAccessListList = append(networkAccessListList, networkAccessListMap)
			}

			dBInstanceSetMap["network_access_list"] = networkAccessListList
		}

		if dBInstanceSet.SupportIpv6 != nil {
			dBInstanceSetMap["support_ipv6"] = dBInstanceSet.SupportIpv6
		}

		if dBInstanceSet.ExpandedCpu != nil {
			dBInstanceSetMap["expanded_cpu"] = dBInstanceSet.ExpandedCpu
		}

		if dBInstanceSet.DeletionProtection != nil {
			dBInstanceSetMap["deletion_protection"] = dBInstanceSet.DeletionProtection
		}

		// rootUser
		if dBInstanceSet.DBInstanceId != nil && strings.HasPrefix(*dBInstanceSet.DBInstanceId, "postgres-") {
			accounts, outErr := service.DescribeRootUser(ctx, *dBInstanceSet.DBInstanceId)
			if outErr != nil {
				continue
			}

			if len(accounts) > 0 {
				dBInstanceSetMap["root_user"] = accounts[0].UserName
			}
		}

		dBInstanceSetList = append(dBInstanceSetList, dBInstanceSetMap)
	}

	_ = d.Set("db_instance_set", dBInstanceSetList)

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("instance_list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil
}
