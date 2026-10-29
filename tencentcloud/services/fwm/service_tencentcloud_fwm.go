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

func (me *FwmService) CreateFwmRuleGroup(ctx context.Context, groupName, product string) (groupId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewCreateRuleGroupRequest()
	request.GroupName = helper.String(groupName)
	request.Product = helper.String(product)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response := fwmv20250611.NewCreateRuleGroupResponse()
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().CreateRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm rule group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create fwm rule group failed, reason:%+v", logId, err)
		errRet = err
		return
	}

	if response.Response.GroupId == nil || *response.Response.GroupId == "" {
		log.Printf("[CRITAL]%s create fwm rule group return empty GroupId, logId=%s", logId, logId)
		errRet = fmt.Errorf("GroupId is nil or empty, create fwm rule group failed.")
		return
	}

	groupId = *response.Response.GroupId
	return
}

func (me *FwmService) DescribeFwmRuleGroupById(ctx context.Context, groupId string) (ruleGroup *fwmv20250611.RuleGroupResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewDescribeRuleGroupsRequest()
	request.Filters = []*fwmv20250611.CommonFilter{
		{
			Name:         helper.String("GroupId"),
			Values:       []*string{helper.String(groupId)},
			OperatorType: helper.Int64(1),
		},
	}
	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response := fwmv20250611.NewDescribeRuleGroupsResponse()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().DescribeRuleGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe fwm rule group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.RuleGroups) < 1 {
		log.Printf("[CRITAL]%s fwm rule group [%s] not found in describe result.", logId, groupId)
		return
	}

	for _, item := range response.Response.RuleGroups {
		if item == nil || item.GroupId == nil {
			continue
		}
		if *item.GroupId == groupId {
			ruleGroup = item
			return
		}
	}

	return
}

func (me *FwmService) ModifyFwmRuleGroupName(ctx context.Context, groupId, groupName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := fwmv20250611.NewModifyRuleGroupRequest()
	request.GroupId = helper.String(groupId)
	request.GroupName = helper.String(groupName)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseFwmV20250611Client().ModifyRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify fwm rule group failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *FwmService) DeleteFwmRuleGroupById(ctx context.Context, groupId string) (errRet error) {
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
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm rule group failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
