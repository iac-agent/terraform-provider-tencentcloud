package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
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

func (me *FwmService) DescribeFwmEdgeAclRuleById(ctx context.Context, groupId, ruleId string, direction int64) (fwmEdgeAclRule *fwmv20250611.EdgeAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeEdgeAclRulesRequest()
	response := fwmv20250611.NewDescribeEdgeAclRulesResponse()

	request.GroupId = helper.String(groupId)
	request.Direction = helper.Int64(direction)
	request.Limit = helper.Int64(1000)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("RuleId"),
			Values:       []*string{helper.String(ruleId)},
			OperatorType: helper.Int64(1),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeEdgeAclRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm edge acl rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Rules) == 0 {
		return nil, nil
	}

	fwmEdgeAclRule = response.Response.Rules[0]
	return
}

func (me *FwmService) DescribeFwmOrganMemberById(ctx context.Context, appId, memberId string) (organMember *fwmv20250611.OrganMemberItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeOrganMembersRequest()
	response := fwmv20250611.NewDescribeOrganMembersResponse()

	operatorType := int64(1)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("AppId"),
			Values:       helper.Strings([]string{appId}),
			OperatorType: helper.Int64(operatorType),
		},
	}
	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = memberId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmClient(iacExtInfo).DescribeOrganMembersWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Members == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm organ members failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Members) < 1 {
		return
	}

	for _, member := range response.Response.Members {
		if member.MemberId != nil && *member.MemberId == memberId {
			organMember = member
			return
		}
	}

	return
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

func (me *FwmService) DescribeFwmMemberGroup(ctx context.Context, userMemberId string) (member *fwmv20250611.OrganMemberItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeOrganMembersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		ref := []*fwmv20250611.OrganMemberItem{}
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseFwmV20250611Client().DescribeOrganMembers(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Members == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe organ members failed, Response is nil."))
			}

			ref = result.Response.Members
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		for _, m := range ref {
			if m != nil && m.MemberId != nil && *m.MemberId == userMemberId {
				member = m
				return
			}
		}

		if len(ref) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *FwmService) ModifyFwmMemberGroup(ctx context.Context, userMemberId, groupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewModifyMemberGroupRequest()
	request.UserMemberId = helper.String(userMemberId)
	request.GroupId = helper.String(groupId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseFwmV20250611Client().ModifyMemberGroupWithContext(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *FwmService) FwmMemberGroupIdToAppIdMemberId(id string) (appId, memberId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid fwm_member_group id: %s", id)
		return
	}
	appId = parts[0]
	memberId = parts[1]
	return
}

func (me *FwmService) DescribeFwmRuleGroupById(ctx context.Context, groupId string) (ret *fwmv20250611.RuleGroupResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeRuleGroupsRequest()
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("GroupId"),
			Values:       []*string{helper.String(groupId)},
			OperatorType: helper.Int64(1),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ruleGroups, err := me.DescribeFwmRuleGroupsByFilter(ctx, map[string]interface{}{
		"Filters": request.Filters,
	})
	if err != nil {
		errRet = err
		return
	}

	for _, ruleGroup := range ruleGroups {
		if ruleGroup == nil || ruleGroup.GroupId == nil {
			continue
		}
		if *ruleGroup.GroupId == groupId {
			ret = ruleGroup
			return
		}
	}

	return
}

func (me *FwmService) DescribeFwmRuleGroupsByFilter(ctx context.Context, param map[string]interface{}) (ret []*fwmv20250611.RuleGroupResp, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = fwmv20250611.NewDescribeRuleGroupsRequest()
		response = fwmv20250611.NewDescribeRuleGroupsResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Product" {
			request.Product = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*fwmv20250611.CommonFilter)
		}
		if k == "Order" {
			request.Order = v.(*string)
		}
		if k == "By" {
			request.By = v.(*string)
		}
	}

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = helper.Int64(offset)
		request.Limit = helper.Int64(limit)
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseFwmV20250611Client().DescribeRuleGroupsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe fwm rule groups failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		ret = append(ret, response.Response.RuleGroups...)
		if len(response.Response.RuleGroups) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *FwmService) DescribeFwmSecurityGroupRuleById(ctx context.Context, groupId, ruleId string) (ret *fwmv20250611.SecGroupRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeSecurityGroupRulesRequest()
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
		result, e := me.client.UseFwmV20250611Client().DescribeSecurityGroupRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm_security_group_rule failed, Response is nil."))
		}

		if len(result.Response.Rules) == 0 {
			return nil
		}

		ret = result.Response.Rules[0]
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *FwmService) DeleteFwmSecurityGroupRuleById(ctx context.Context, groupId, ruleId string) error {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDeleteSecurityGroupRuleRequest()
	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DeleteSecurityGroupRuleWithContext(ctx, request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				if strings.Contains(errCode, "ResourceNotFound") {
					return nil
				}
			}
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm_security_group_rule failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete fwm_security_group_rule failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func (me *FwmService) DescribeFwmNatAclRuleById(ctx context.Context, groupId, ruleId string, direction int64) (rule *fwmv20250611.NatAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeNatAclRulesRequest()
	request.GroupId = helper.String(groupId)
	request.Direction = helper.Int64(direction)
	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("RuleId"),
			Values:       []*string{helper.String(ruleId)},
			OperatorType: helper.Int64(1),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseFwmV20250611Client().DescribeNatAclRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm_nat_acl_rule failed, Response is nil."))
		}

		if len(response.Response.Rules) == 0 {
			return nil
		}

		rule = response.Response.Rules[0]
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// EdgeAclRuleGroupProduct is the fixed Product value of the internet border
// (cfw_edge_acl) rule group. CreateEdgeAclRuleGroup and DescribeRuleGroups
// both require it, so the provider hard-codes it instead of exposing an
// attribute the user could set to an unsupported value.
const EdgeAclRuleGroupProduct = "cfw_edge_acl"

// CreateFwmEdgeAclRuleGroup creates one internet border rule group and
// returns the group id assigned by the cloud.
func (me *FwmService) CreateFwmEdgeAclRuleGroup(ctx context.Context, groupName string, rules []*fwmv20250611.EdgeAclRuleInfo) (groupId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewCreateEdgeAclRuleGroupRequest()
	request.Product = helper.String(EdgeAclRuleGroupProduct)
	if groupName != "" {
		request.GroupName = helper.String(groupName)
	}
	request.Rules = rules

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().CreateEdgeAclRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm edge acl rule group failed, Response is nil."))
		}

		if result.Response.GroupId == nil || *result.Response.GroupId == "" {
			return resource.NonRetryableError(fmt.Errorf("Create fwm edge acl rule group failed, GroupId is empty."))
		}

		groupId = *result.Response.GroupId
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// DescribeFwmEdgeAclRuleGroupById reads one internet border rule group by id.
//
// DescribeRuleGroups is a list API; GroupId is passed as an exact-match filter
// and the returned items are matched again locally, so a caller never has to
// trust position 0 of the list.
func (me *FwmService) DescribeFwmEdgeAclRuleGroupById(ctx context.Context, groupId string) (ruleGroup *fwmv20250611.RuleGroupResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeRuleGroupsRequest()
	request.Product = helper.String(EdgeAclRuleGroupProduct)
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("GroupId"),
			Values:       []*string{helper.String(groupId)},
			OperatorType: helper.Int64(1),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = helper.Int64(offset)
		request.Limit = helper.Int64(limit)

		response := fwmv20250611.NewDescribeRuleGroupsResponse()
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseFwmV20250611Client().DescribeRuleGroupsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe fwm rule groups failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		for _, item := range response.Response.RuleGroups {
			if item != nil && item.GroupId != nil && *item.GroupId == groupId {
				ruleGroup = item
				return
			}
		}

		if len(response.Response.RuleGroups) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

// DescribeFwmEdgeAclRulesByGroupId reads every rule of one rule group.
//
// DescribeEdgeAclRules requires Direction, so both directions are queried and
// merged; the direction filter is the only stable way to page through the
// rules of a group.
func (me *FwmService) DescribeFwmEdgeAclRulesByGroupId(ctx context.Context, groupId string) (rules []*fwmv20250611.EdgeAclRuleResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s describe fwm edge acl rules of group [%s] fail, reason[%s]\n", logId, groupId, errRet.Error())
		}
	}()

	directions := []int64{0, 1}
	for _, direction := range directions {
		var (
			offset int64 = 0
			limit  int64 = 100
		)
		for {
			request := fwmv20250611.NewDescribeEdgeAclRulesRequest()
			request.GroupId = helper.String(groupId)
			request.Direction = helper.Int64(direction)
			request.Offset = helper.Int64(offset)
			request.Limit = helper.Int64(limit)
			request.Order = helper.String("asc")
			request.By = helper.String("Sequence")

			response := fwmv20250611.NewDescribeEdgeAclRulesResponse()
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				ratelimit.Check(request.GetAction())
				result, e := me.client.UseFwmV20250611Client().DescribeEdgeAclRulesWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				}
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe fwm edge acl rules failed, Response is nil."))
				}

				response = result
				return nil
			})

			if err != nil {
				errRet = err
				return
			}

			rules = append(rules, response.Response.Rules...)
			if len(response.Response.Rules) < int(limit) {
				break
			}

			offset += limit
		}
	}

	return
}

// ModifyFwmEdgeAclRuleGroup renames one rule group.
//
// ModifyRuleGroup only accepts GroupName; the rules of a group cannot be
// edited through it, which is why the caller turns rule changes into a
// replacement.
func (me *FwmService) ModifyFwmEdgeAclRuleGroup(ctx context.Context, groupId, groupName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewModifyRuleGroupRequest()
	request.GroupId = helper.String(groupId)
	request.GroupName = helper.String(groupName)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseFwmV20250611Client().ModifyRuleGroupWithContext(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// DeleteFwmEdgeAclRuleGroup deletes exactly one rule group.
//
// DeleteRuleGroup accepts a batch of ids; only the single managed id is ever
// passed so a destroy can never remove a rule group that Terraform does not
// own.
func (me *FwmService) DeleteFwmEdgeAclRuleGroup(ctx context.Context, groupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDeleteRuleGroupRequest()
	request.GroupIds = []*string{helper.String(groupId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DeleteRuleGroupWithContext(ctx, request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if strings.Contains(ee.GetCode(), "ResourceNotFound") {
					return nil
				}
			}
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm edge acl rule group failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
