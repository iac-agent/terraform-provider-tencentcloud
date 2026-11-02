package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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

// FWM strategy dispatch status values returned by DescribeStrategyDispatchStatus.
const (
	// FwmStrategyDispatchStatusNoChange means no strategy change is pending.
	FwmStrategyDispatchStatusNoChange int64 = 0
	// FwmStrategyDispatchStatusDispatching means the dispatch is in progress.
	FwmStrategyDispatchStatusDispatching int64 = 1
	// FwmStrategyDispatchStatusSuccess means the dispatch finished successfully.
	FwmStrategyDispatchStatusSuccess int64 = 2
	// FwmStrategyDispatchStatusFailed means the dispatch failed.
	FwmStrategyDispatchStatusFailed int64 = 3
	// FwmStrategyDispatchStatusPendingUpdate means an update is waiting to be dispatched.
	FwmStrategyDispatchStatusPendingUpdate int64 = 4
)

// FwmStrategyDispatchAbortStatus is the DispatchStrategy request value that aborts
// an ongoing dispatch instead of starting a new one.
const FwmStrategyDispatchAbortStatus int64 = 2

// DispatchFwmStrategy triggers a rule dispatch through the DispatchStrategy API.
// product limits the dispatch scope; an empty product means every product.
// status is the raw DispatchStrategy status value (1: dispatch, 2: abort).
func (me *FwmService) DispatchFwmStrategy(ctx context.Context, product string, status int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDispatchStrategyRequest()
	request.Status = helper.Int64(status)
	if product != "" {
		request.Product = helper.String(product)
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseFwmClient().DispatchStrategyWithContext(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("Dispatch fwm strategy failed, Response is nil.")
		return
	}

	return
}

// DescribeFwmStrategyDispatchStatus returns the raw status progress of the
// strategy dispatch scoped by product. An empty product means every product.
func (me *FwmService) DescribeFwmStrategyDispatchStatus(ctx context.Context, product string) (ret *fwmv20250611.DescribeStrategyDispatchStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeStrategyDispatchStatusRequest()
	if product != "" {
		request.Product = helper.String(product)
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmClient().DescribeStrategyDispatchStatusWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm strategy dispatch status failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		ret = result.Response
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// WaitFwmStrategyDispatch polls DescribeStrategyDispatchStatus until the dispatch
// reaches a terminal status (success) or a terminal failure status, or until the
// timeout is exceeded. Only status 2 is treated as completion: 1 (dispatching), 4
// (pending update) and 0 (no change) are intermediate and keep the poll running.
// For status 3 (failed) the failure reason reported by the cloud is returned.
func (me *FwmService) WaitFwmStrategyDispatch(ctx context.Context, product string, timeout time.Duration, interval time.Duration) (errRet error) {
	var last *fwmv20250611.DescribeStrategyDispatchStatusResponseParams
	deadline := time.Now().Add(timeout)

	for {
		statusParams, e := me.DescribeFwmStrategyDispatchStatus(ctx, product)
		if e != nil {
			errRet = fmt.Errorf("querying dispatch status failed: %v", e)
			return
		}
		last = statusParams

		if endErr, done := fwmStrategyDispatchTerminalState(statusParams); done {
			errRet = endErr
			return
		}

		// The remaining budget may be shorter than one interval; never sleep past it.
		remaining := time.Until(deadline)
		if remaining <= 0 {
			errRet = fmt.Errorf(
				"dispatch did not converge within %s: last observed %s",
				timeout,
				fwmStrategyDispatchSnapshot(last),
			)
			return
		}
		if interval < remaining {
			remaining = interval
		}
		select {
		case <-ctx.Done():
			errRet = fmt.Errorf("dispatch wait was cancelled: %v; last observed %s", ctx.Err(), fwmStrategyDispatchSnapshot(last))
			return
		case <-time.After(remaining):
		}
	}
}

// fwmStrategyDispatchTerminalState reports whether the observed dispatch status is
// terminal, and the resulting error (nil when the dispatch succeeded). Only status 2
// counts as success; status 3 counts as a failure carrying the cloud-side reason.
func fwmStrategyDispatchTerminalState(p *fwmv20250611.DescribeStrategyDispatchStatusResponseParams) (error, bool) {
	if p == nil || p.Status == nil {
		return nil, false
	}

	switch *p.Status {
	case FwmStrategyDispatchStatusSuccess:
		return nil, true
	case FwmStrategyDispatchStatusFailed:
		reason := "the cloud did not report an error message"
		if p.ErrorMsg != nil && *p.ErrorMsg != "" {
			reason = *p.ErrorMsg
		}
		return fmt.Errorf("dispatch failed: %s", reason), true
	default:
		return nil, false
	}
}

// fwmStrategyDispatchSnapshot renders the last observed dispatch status so a timeout
// can be told apart from a wrong completion predicate.
func fwmStrategyDispatchSnapshot(p *fwmv20250611.DescribeStrategyDispatchStatusResponseParams) string {
	if p == nil {
		return "status=<nil>"
	}

	value := func(v *string) string {
		if v == nil {
			return "<nil>"
		}
		return *v
	}
	number := func(v *int64) string {
		if v == nil {
			return "<nil>"
		}
		return fmt.Sprintf("%d", *v)
	}

	progress := "<nil>"
	if p.Progress != nil {
		progress = fmt.Sprintf("%v", *p.Progress)
	}

	return fmt.Sprintf(
		"status=%s, progress=%s, start_time=%s, end_time=%s, rule_group_num=%s, error_msg=%s",
		number(p.Status),
		progress,
		value(p.StartTime),
		value(p.EndTime),
		number(p.RuleGroupNum),
		value(p.ErrorMsg),
	)
}
