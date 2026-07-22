package vpc

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name", "project_id"},
				Description:   "ID 安全 组 到 是 queried. Conflict 使用 `名称` 和 `project_id`。",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  tccommon.ValidateStringLengthInRange(1, 60),
				ConflictsWith: []string{"security_group_id"},
				Description:   "名称 安全 组 到 是 queried. Conflict 使用 `security_group_id`。",
			},
			"project_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"security_group_id"},
				Description:   "项目 ID 安全 组 到 是 queried. Conflict 使用 `security_group_id`。",
			},
			"tags": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"security_group_id"},
				Description:   "标签 的 安全 组 到 是 queried. Conflict 使用 `security_group_id`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 安全 组。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 安全 组。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 安全 组。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 安全 组。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 安全 组。",
						},
						"be_associate_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 安全 组 binding resources。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID 安全 组。",
						},
						"ingress": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Ingress 规则 集合. For items like `[操作]#[cidr_ip]#[端口]#[协议]`，它 表示 regular 规则; 对于 items like `sg-XXXX`，它 表示 nested 安全 组。",
						},
						"egress": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Egress 规则 集合. For items like `[操作]#[cidr_ip]#[端口]#[协议]`，它 表示 regular 规则; 对于 items like `sg-XXXX`，它 表示 nested 安全 组。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 安全 组。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSecurityGroupsRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_security_groups.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := m.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := VpcService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	var (
		sgId      *string
		sgName    *string
		projectId *int
	)

	idBuilder := strings.Builder{}
	idBuilder.WriteString("securityGroups-")

	if raw, ok := d.GetOk("security_group_id"); ok {
		sgId = helper.String(raw.(string))
		idBuilder.WriteString(*sgId)
		idBuilder.WriteRune('-')
	}

	if raw, ok := d.GetOk("name"); ok {
		sgName = common.StringPtr(raw.(string))
		idBuilder.WriteString(*sgName)
		idBuilder.WriteRune('-')
	}

	if raw, ok := d.GetOkExists("project_id"); ok {
		projectId = helper.Int(raw.(int))
		idBuilder.WriteString(strconv.Itoa(*projectId))
	}

	tags := helper.GetTags(d, "tags")

	sgs, err := vpcService.DescribeSecurityGroups(ctx, sgId, sgName, projectId, tags)
	if err != nil {
		return err
	}

	if len(sgs) == 0 {
		_ = d.Set("security_groups", []map[string]interface{}{})
		d.SetId(idBuilder.String())
		return nil
	}

	sgMap := make(map[string]*vpc.SecurityGroup, len(sgs))
	for _, sg := range sgs {
		if sg.SecurityGroupId == nil {
			return errors.New("security group id is nil")
		}
		sgMap[*sg.SecurityGroupId] = sg
	}

	sgIds := make([]string, 0, len(sgs))
	for _, sg := range sgs {
		sgIds = append(sgIds, *sg.SecurityGroupId)
	}

	associateSet, err := vpcService.DescribeSecurityGroupsAssociate(ctx, sgIds)
	if err != nil {
		return err
	}

	sgInstances := make([]map[string]interface{}, 0, len(sgs))
	for _, associate := range associateSet {
		if associate.SecurityGroupId == nil {
			return errors.New("associate statistics security group id is nil")
		}

		if sg, ok := sgMap[*associate.SecurityGroupId]; ok {
			count := int(*associate.CVM + *associate.ENI + *associate.CDB + *associate.CLB)

			// normally projectId default value is 0
			if sg.ProjectId == nil {
				return errors.New("associate statistics project id is nil")
			}

			projectId, err := strconv.Atoi(*sg.ProjectId)
			if err != nil {
				return fmt.Errorf("securtiy group %s project id invalid: %v", *sg.SecurityGroupId, err)
			}

			respIngress, respEgress, exist, err := vpcService.DescribeSecurityGroupPolices(ctx, *sg.SecurityGroupId)
			if err != nil {
				return err
			}

			if !exist {
				// when read security group all rules, it doesn't exist, maybe delete on other places, ignore it
				continue
			}

			respTags, err := tagService.DescribeResourceTags(ctx, "cvm", "sg", region, *sg.SecurityGroupId)
			if err != nil {
				return err
			}

			ingress := make([]string, 0, len(respIngress))
			for _, in := range respIngress {
				ingress = append(ingress, in.String())
			}

			egress := make([]string, 0, len(respEgress))
			for _, eg := range respEgress {
				egress = append(egress, eg.String())
			}

			sgInstances = append(sgInstances, map[string]interface{}{
				"security_group_id":  *sg.SecurityGroupId,
				"name":               *sg.SecurityGroupName,
				"description":        *sg.SecurityGroupDesc,
				"create_time":        *sg.CreatedTime,
				"be_associate_count": count,
				"project_id":         projectId,
				"ingress":            ingress,
				"egress":             egress,
				"tags":               respTags,
			})
		}
	}

	if len(sgInstances) != len(sgs) {
		return errors.New("security group associate statistics is not enough")
	}

	_ = d.Set("security_groups", sgInstances)
	d.SetId(idBuilder.String())

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), sgInstances); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]", logId, output.(string), err)
			return err
		}
	}

	return nil
}
