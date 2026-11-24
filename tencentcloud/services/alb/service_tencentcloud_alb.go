package alb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

type AlbService struct {
	client *connectivity.TencentCloudClient
}

// CreateRules creates forwarding rules and returns the created rule ids.
func (me *AlbService) CreateRules(ctx context.Context, request *albv20251030.CreateRulesRequest) (ruleIds []string, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.CreateRulesResponse
	)

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().CreateRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create alb rules failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s create alb rules failed, reason:%+v", logId, errRet)
		return
	}

	if len(response.Response.RuleIds) < 1 {
		errRet = fmt.Errorf("api %s returns empty rule ids", request.GetAction())
		return
	}

	for _, ruleId := range response.Response.RuleIds {
		if ruleId != nil {
			ruleIds = append(ruleIds, *ruleId)
		}
	}
	return
}

// DescribeRules queries forwarding rules; the caller picks its own items
// from the returned list.
func (me *AlbService) DescribeRules(ctx context.Context, request *albv20251030.DescribeRulesRequest) (rules []*albv20251030.RuleOutput, nextToken string, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.DescribeRulesResponse
	)

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DescribeRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe alb rules failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe alb rules failed, reason:%+v", logId, errRet)
		return
	}

	if response.Response.Rules != nil {
		rules = response.Response.Rules
	}
	if response.Response.NextToken != nil {
		nextToken = *response.Response.NextToken
	}
	return
}

// ModifyRulesAttributes updates forwarding rule attributes.
func (me *AlbService) ModifyRulesAttributes(ctx context.Context, request *albv20251030.ModifyRulesAttributesRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().ModifyRulesAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify alb rules attributes failed, reason:%+v", logId, err)
	}

	return
}

// DeleteRules deletes forwarding rules.
func (me *AlbService) DeleteRules(ctx context.Context, request *albv20251030.DeleteRulesRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DeleteRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete alb rules failed, reason:%+v", logId, err)
	}

	return
}
