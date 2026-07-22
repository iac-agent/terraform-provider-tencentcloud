package ccn

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCcnCrossBorderCompliance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnCrossBorderComplianceRead,
		Schema: map[string]*schema.Schema{
			"service_provider": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) 服务 provider，可选 值: 'UNICOM'。",
			},

			"compliance_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "(Exact match) compliance approval form: 'ID'。",
			},

			"company": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy 查询) 公司名称",
			},

			"uniform_social_credit_code": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) Uniform Social Credit 代码",
			},

			"legal_person": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy 查询) legal representative。",
			},

			"issuing_authority": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy 查询) Issuing authority。",
			},

			"business_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy 查询) business license 地址",
			},

			"post_code": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "(Exact match) post 代码",
			},

			"manager": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy 查询) Person 在 charge。",
			},

			"manager_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact 查询) ID 数量 person 在 charge。",
			},

			"manager_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Fuzzy 查询) ID card 地址 的 person 在 charge。",
			},

			"manager_telephone": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) contact 数量 person 在 charge。",
			},

			"email": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) email。",
			},

			"service_start_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) 服务 start date，such 作为: '2020-07-28'。",
			},

			"service_end_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) 服务 end date，such 作为: '2020-07-28'。",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "(Exact match) 状态 Pending: PENDING，Passed: APPROVED，Denied: DENY。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudCcnCrossBorderComplianceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ccn_cross_border_compliance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_provider"); ok {
		paramMap["service_provider"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("compliance_id"); v != nil {
		paramMap["compliance_id"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("company"); ok {
		paramMap["company"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniform_social_credit_code"); ok {
		paramMap["uniform_social_credit_code"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("legal_person"); ok {
		paramMap["legal_person"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("issuing_authority"); ok {
		paramMap["issuing_authority"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("business_address"); ok {
		paramMap["business_address"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("post_code"); v != nil {
		paramMap["post_code"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("manager"); ok {
		paramMap["manager"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_id"); ok {
		paramMap["manager_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_address"); ok {
		paramMap["manager_address"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_telephone"); ok {
		paramMap["manager_telephone"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email"); ok {
		paramMap["email"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_start_date"); ok {
		paramMap["service_start_date"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_end_date"); ok {
		paramMap["service_end_date"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		paramMap["state"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var crossBorderComplianceSet []*vpc.CrossBorderCompliance

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCcnCrossBorderComplianceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		crossBorderComplianceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(crossBorderComplianceSet))
	tmpList := make([]map[string]interface{}, 0, len(crossBorderComplianceSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
