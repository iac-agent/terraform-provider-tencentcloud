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
				Description: "名称 TcaplusDB cluster to be query。",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID TcaplusDB cluster to be query。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File for saving results。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 TcaplusDB cluster. Each element 包含following attributes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 TcaplusDB cluster。",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID TcaplusDB cluster。",
						},
						"idl_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IDL 类型 TcaplusDB cluster。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID TcaplusDB cluster。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID TcaplusDB cluster。",
						},
						"password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access 密码 of the TcaplusDB cluster。",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network 类型 TcaplusDB cluster。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of the TcaplusDB cluster。",
						},
						"password_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "密码 状态 TcaplusDB cluster. 有效值：`unmodifiable`，`modifiable`. `unmodifiable` means the 密码 can not be changed in this moment; `modifiable` means the 密码 can be changed in this moment。",
						},
						"api_access_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access ID TcaplusDB cluster.For TcaplusDB SDK connect。",
						},
						"api_access_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access ip of the TcaplusDB cluster.For TcaplusDB SDK connect。",
						},
						"api_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Access 端口 of the TcaplusDB cluster.For TcaplusDB SDK connect。",
						},
						"old_password_expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间 of the old 密码 If `password_status` is `unmodifiable`，it means the old 密码 has not yet expired。",
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
