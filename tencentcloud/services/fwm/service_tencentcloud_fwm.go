package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type FwmService struct {
	client *connectivity.TencentCloudClient
}

func NewFwmService(client *connectivity.TencentCloudClient) *FwmService {
	return &FwmService{client: client}
}

func (me *FwmService) DescribeFwmAccountGroupById(ctx context.Context, groupId string) (accountGroup *fwmv20250611.AccountGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeAccountGroupsRequest()
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("GroupId"),
			Values:       []*string{helper.String(groupId)},
			OperatorType: helper.Int64(1),
		},
	}
	request.Limit = helper.Int64(20)
	request.Offset = helper.Int64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response := fwmv20250611.NewDescribeAccountGroupsResponse()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeAccountGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm account group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.AccountGroups) < 1 {
		return
	}

	accountGroup = response.Response.AccountGroups[0]
	return
}

func (me *FwmService) DeleteFwmAccountGroupById(ctx context.Context, groupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDeleteAccountGroupRequest()
	request.GroupId = helper.String(groupId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DeleteAccountGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm account group failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

const (
	// FWM_VPC_ACL_RULE_DESCRIBE_LIMIT is the max page size allowed by DescribeVpcAclRules.
	FWM_VPC_ACL_RULE_DESCRIBE_LIMIT  = int64(100)
	FWM_VPC_ACL_RULE_DESCRIBE_OFFSET = int64(0)
)

func (me *FwmService) CreateFwmVpcAclRule(ctx context.Context, groupId string, rules []*fwmv20250611.VpcAclRule) (ruleIds []*string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewCreateVpcAclRuleRequest()
	request.GroupId = helper.String(groupId)
	request.Rules = rules

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response := fwmv20250611.NewCreateVpcAclRuleResponse()
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().CreateVpcAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm vpc acl rule failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ruleIds = response.Response.RuleIds
	return
}

func (me *FwmService) DescribeFwmVpcAclRulesByGroupId(ctx context.Context, groupId string) (rules []*fwmv20250611.VpcAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeVpcAclRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Offset = helper.Int64(FWM_VPC_ACL_RULE_DESCRIBE_OFFSET)
	request.Limit = helper.Int64(FWM_VPC_ACL_RULE_DESCRIBE_LIMIT)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response := fwmv20250611.NewDescribeVpcAclRulesResponse()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeVpcAclRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm vpc acl rules failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	rules = response.Response.Rules
	return
}

func (me *FwmService) DescribeFwmVpcAclRuleById(ctx context.Context, groupId, ruleId string) (rule *fwmv20250611.VpcAclRuleResp, errRet error) {
	rules, err := me.DescribeFwmVpcAclRulesByGroupId(ctx, groupId)
	if err != nil {
		errRet = err
		return
	}

	for _, item := range rules {
		if item == nil || item.RuleId == nil {
			continue
		}
		if *item.RuleId == ruleId {
			rule = item
			return
		}
	}

	return
}

func (me *FwmService) DeleteFwmVpcAclRule(ctx context.Context, groupId string, ruleIds []string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDeleteVpcAclRuleRequest()
	request.GroupId = helper.String(groupId)
	request.RuleIds = helper.Strings(ruleIds)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DeleteVpcAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm vpc acl rule failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
