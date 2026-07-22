package tcaplusdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func ResourceTencentCloudTcaplusCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusClusterCreate,
		Read:   resourceTencentCloudTcaplusClusterRead,
		Update: resourceTencentCloudTcaplusClusterUpdate,
		Delete: resourceTencentCloudTcaplusClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"idl_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(TCAPLUS_IDL_TYPES),
				Description:  "IDL 类型 TcaplusDB 集群. 有效值：`PROTO` 和 `TDR`。",
			},
			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 30),
				Description:  "名称 TcaplusDB 集群. 名称 长度 should 是 between 1 和 30。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID TcaplusDB 集群。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "子网 ID TcaplusDB 集群。",
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 12 || len(value) > 16 {
						errors = append(errors, fmt.Errorf("invalid password, length should between 12 and 16"))
						return
					}
					var match = make(map[string]bool)
					for i := 0; i < len(value); i++ {
						if len(match) >= 2 {
							break
						}
						if value[i] >= '0' && value[i] <= '9' {
							match["number"] = true
							continue
						}
						if value[i] >= 'a' && value[i] <= 'z' {
							match["low"] = true
							continue
						}
						if value[i] >= 'A' && value[i] <= 'Z' {
							match["up"] = true
							continue
						}
					}
					if len(match) < 2 {
						errors = append(errors, fmt.Errorf("invalid password, a-z and 0-9 and A-Z must contain"))
					}
					return
				},
				Description: "密码 的 TcaplusDB 集群. 密码 长度 should 是 between 12 和 16. 密码 必须 是 *mix* 的 uppercase letters (A-Z)，lowercase *letters* (-z) 和 *numbers* (0-9)。",
			},
			"old_password_expire_last": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: tccommon.ValidateIntegerMin(300),
				Description:  "过期时间 的 old 密码 after 密码 update，单位: second。",
			},

			// Computed values.
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
				Description: "密码 状态 TcaplusDB 集群. 有效值：`unmodifiable`，`modifiable`. `unmodifiable`. 其中 表示 密码 可以 不 是 changed 在 此 moment; `modifiable`，其中 表示 密码 可以 是 changed 在 此 moment。",
			},
			"api_access_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access ID TcaplusDB 集群.For TcaplusDB SDK connect。",
			},
			"api_access_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access IP 的 TcaplusDB 集群.For TcaplusDB SDK connect。",
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
	}
}

func resourceTencentCloudTcaplusClusterCreate(d *schema.ResourceData, meta interface{}) error {

	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		idlType     = d.Get("idl_type").(string)
		clusterName = d.Get("cluster_name").(string)
		vpcId       = d.Get("vpc_id").(string)
		subnetId    = d.Get("subnet_id").(string)
		password    = d.Get("password").(string)
	)

	var clusterId string
	var inErr, outErr error

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		clusterId, inErr = tcaplusService.CreateCluster(ctx, idlType, clusterName, vpcId, subnetId, password)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	d.SetId(clusterId)
	time.Sleep(3 * time.Second)
	return resourceTencentCloudTcaplusClusterRead(d, meta)
}

func resourceTencentCloudTcaplusClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterInfo, has, err := tcaplusService.DescribeCluster(ctx, d.Id())
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			clusterInfo, has, err = tcaplusService.DescribeCluster(ctx, d.Id())
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("idl_type", clusterInfo.IdlType)
	_ = d.Set("cluster_name", clusterInfo.ClusterName)
	_ = d.Set("vpc_id", clusterInfo.VpcId)
	_ = d.Set("subnet_id", clusterInfo.SubnetId)
	_ = d.Set("network_type", clusterInfo.NetworkType)
	_ = d.Set("create_time", clusterInfo.CreatedTime)
	_ = d.Set("password_status", clusterInfo.PasswordStatus)
	_ = d.Set("api_access_id", clusterInfo.ApiAccessId)
	_ = d.Set("api_access_ip", clusterInfo.ApiAccessIp)
	_ = d.Set("api_access_port", clusterInfo.ApiAccessPort)

	if clusterInfo.OldPasswordExpireTime == nil || *clusterInfo.OldPasswordExpireTime == "" {
		_ = d.Set("old_password_expire_time", "-")
	} else {
		_ = d.Set("old_password_expire_time", clusterInfo.OldPasswordExpireTime)
	}

	return nil
}

func resourceTencentCloudTcaplusClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_clusterupdate")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	if d.HasChange("cluster_name") {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyClusterName(ctx, d.Id(), d.Get("cluster_name").(string))
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("password") {
		oldPwd, newPwd := d.GetChange("password")
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyClusterPassword(ctx, d.Id(),
				oldPwd.(string),
				newPwd.(string),
				int64(d.Get("old_password_expire_last").(int)))

			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "FailedOperation.OldPasswordInUse" {
					err = fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.OldPasswordInUse,`password_status` is unmodifiable now, can modify after `old_password_expire_time`")
					return resource.NonRetryableError(err)
				}
			}
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudTcaplusClusterRead(d, meta)
}

func resourceTencentCloudTcaplusClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	_, err := tcaplusService.DeleteCluster(ctx, d.Id())

	if err != nil {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, err = tcaplusService.DeleteCluster(ctx, d.Id())
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeCluster(ctx, d.Id())
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeCluster(ctx, d.Id())
			if err != nil {
				return tccommon.RetryError(err)
			}

			if has {
				err = fmt.Errorf("delete cluster fail, cluster still exist from sdk DescribeClusters")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete cluster fail, cluster still exist from sdk DescribeClusters")
	}
}
