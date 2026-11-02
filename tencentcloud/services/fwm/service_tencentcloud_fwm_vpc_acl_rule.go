package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// FwmService methods for tencentcloud_fwm_vpc_acl_rule2 (VPC border firewall ACL rule).
//
// The cloud API expresses ACL rules as a list (CreateVpcAclRule.Rules /
// DescribeVpcAclRules.Response.Rules / DeleteVpcAclRule.RuleIds), while the
// Terraform resource describes exactly one rule. Every method below therefore
// sends a single-element list and validates that exactly one element comes back.

// CreateFwmVpcAclRule creates one VPC border firewall ACL rule and returns the
// generated rule id.
func (me *FwmService) CreateFwmVpcAclRule(ctx context.Context, groupId string, rule *fwmv20250611.VpcAclRule) (ruleId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewCreateVpcAclRuleRequest()
	request.GroupId = helper.String(groupId)
	request.Rules = []*fwmv20250611.VpcAclRule{rule}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().CreateVpcAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm vpc acl rule failed, Response is nil."))
		}

		// The API accepts a list but this resource manages exactly one rule, so
		// anything other than a single returned id is an unexpected response: it
		// cannot be mapped back onto this resource identity.
		if len(result.Response.RuleIds) != 1 || result.Response.RuleIds[0] == nil || *result.Response.RuleIds[0] == "" {
			return resource.NonRetryableError(fmt.Errorf("Create fwm vpc acl rule failed, expect exactly one non-empty rule id, got %d", len(result.Response.RuleIds)))
		}

		ruleId = *result.Response.RuleIds[0]
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// DescribeFwmVpcAclRuleById returns the ACL rule identified by groupId+ruleId,
// or (nil, nil) when the rule no longer exists in the group.
func (me *FwmService) DescribeFwmVpcAclRuleById(ctx context.Context, groupId, ruleId string) (rule *fwmv20250611.VpcAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeVpcAclRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("RuleId"),
			Values:       []*string{helper.String(ruleId)},
			OperatorType: helper.Int64(1),
		},
	}
	request.Offset = helper.Int64(0)
	request.Limit = helper.Int64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeVpcAclRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm vpc acl rules failed, Response is nil."))
		}

		// Match on the exact rule id even though the request already filters by
		// RuleId, so a filter that is ignored server side can never make this
		// method return some other rule of the same group.
		for _, item := range result.Response.Rules {
			if item != nil && item.RuleId != nil && *item.RuleId == ruleId {
				rule = item
				return nil
			}
		}

		rule = nil
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// ModifyFwmVpcAclRule updates one existing ACL rule in place.
func (me *FwmService) ModifyFwmVpcAclRule(ctx context.Context, groupId string, rule *fwmv20250611.VpcAclRule) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewModifyVpcAclRuleRequest()
	request.GroupId = helper.String(groupId)
	request.Rule = rule

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().ModifyVpcAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify fwm vpc acl rule failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// DeleteFwmVpcAclRule deletes one ACL rule. A rule that is already gone is
// treated as success so `terraform destroy` stays idempotent.
func (me *FwmService) DeleteFwmVpcAclRule(ctx context.Context, groupId, ruleId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDeleteVpcAclRuleRequest()
	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DeleteVpcAclRuleWithContext(ctx, request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				// Only a rule that is confirmed absent may be swallowed. Permission,
				// throttling and network failures must surface so the id stays in
				// state for a retry.
				if ee.GetCode() == "ResourceNotFound" || strings.HasPrefix(ee.GetCode(), "ResourceNotFound.") ||
					ee.GetCode() == "ResourceUnavailable.RuleNotFound" {
					log.Printf("[WARN]%s delete fwm vpc acl rule, rule not found, treat as success, groupId=%s, ruleId=%s", logId, groupId, ruleId)
					return nil
				}
			}
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm vpc acl rule failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
