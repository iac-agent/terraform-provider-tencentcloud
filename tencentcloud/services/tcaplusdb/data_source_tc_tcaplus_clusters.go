package tcaplusdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudTcaplusClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 TcaplusDB 集群 到 是 查询。",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID TcaplusDB 集群 到 是 查询。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File 对于 saving results。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 TcaplusDB 集群. Each element 包含following attributes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 TcaplusDB 集群。",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID TcaplusDB 集群。",
						},
						"idl_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IDL 类型 TcaplusDB 集群。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID TcaplusDB 集群。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID TcaplusDB 集群。",
						},
						"password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access 密码 的 TcaplusDB 集群。",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network 类型 TcaplusDB 集群。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 TcaplusDB 集群。",
						},
						"password_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "密码 状态 TcaplusDB 集群. 有效值：`unmodifiable`，`modifiable`. `unmodifiable` 表示 密码 可以 不 是 changed 在 此 moment; `modifiable` 表示 密码 可以 是 changed 在 此 moment。",
						},
						"api_access_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access ID TcaplusDB 集群.For TcaplusDB SDK connect。",
						},
						"api_access_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access ip 的 TcaplusDB 集群.For TcaplusDB SDK connect。",
						},
						"api_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Access 端口 的 TcaplusDB 集群.For TcaplusDB SDK connect。",
						},
						"old_password_expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间 的 old 密码 如果 `password_status` 是 `unmodifiable`，它 表示 old 密码 has 不 yet expired。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcaplus_clusters.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcaplusService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	clusterId := d.Get("cluster_id").(string)
	clusterName := d.Get("cluster_name").(string)

	clusters, err := service.DescribeClusters(ctx, clusterId, clusterName)
	if err != nil {
		clusters, err = service.DescribeClusters(ctx, clusterId, clusterName)
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(clusters))

	for _, cluster := range clusters {
		listItem := make(map[string]interface{})
		listItem["cluster_name"] = cluster.ClusterName
		listItem["cluster_id"] = cluster.ClusterId
		listItem["idl_type"] = cluster.IdlType
		listItem["vpc_id"] = cluster.VpcId
		listItem["subnet_id"] = cluster.SubnetId
		listItem["password"] = cluster.Password
		listItem["network_type"] = cluster.NetworkType
		listItem["create_time"] = cluster.CreatedTime
		listItem["password_status"] = cluster.PasswordStatus
		listItem["api_access_id"] = cluster.ApiAccessId
		listItem["api_access_ip"] = cluster.ApiAccessIp
		listItem["api_access_port"] = cluster.ApiAccessPort
		listItem["old_password_expire_time"] = cluster.OldPasswordExpireTime
		list = append(list, listItem)
	}

	d.SetId("cluster." + clusterId + "." + clusterName)
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil

}
