package alb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// NewAlbService returns the service wrapper of the ALB (application load
// balancer) product APIs used by the framework resources of this package.
func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

// AlbService wraps the ALB SDK calls of the framework resources. Every method
// takes the fully built SDK request and owns the retry, rate limit and logging
// behaviour, so that lifecycle methods never talk to the SDK client directly.
type AlbService struct {
	client *connectivity.TencentCloudClient
}

// CreateRules creates forwarding rules in one listener and returns the id of
// the first created rule. The create API is synchronous and batch shaped: the
// caller always sends a single rule, so the first id of the response belongs to
// it.
func (me *AlbService) CreateRules(ctx context.Context, request *albv20251030.CreateRulesRequest) (ruleId string, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.CreateRulesResponse
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAlbV20251030Client().CreateRules(request)
		if err != nil {
			return tccommon.RetryError(err)
		}

		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateRules returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RuleIds) < 1 || response.Response.RuleIds[0] == nil {
		errRet = fmt.Errorf("api %s returns empty rule id", request.GetAction())
		return
	}

	ruleId = *response.Response.RuleIds[0]
	return
}

// DescribeRules returns the forwarding rules matching the given request. The
// describe API pages with MaxResults/NextToken, so the pagination loop lives
// here: a request addressed by a unique rule id terminates after the first
// round, while a broader filter keeps paging until the cloud reports the last
// page.
func (me *AlbService) DescribeRules(ctx context.Context, request *albv20251030.DescribeRulesRequest) (rules []*albv20251030.RuleOutput, errRet error) {
	var (
		logId = tccommon.GetLogId(ctx)
		limit = int64(100)
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	rules = make([]*albv20251030.RuleOutput, 0)

	for {
		var response *albv20251030.DescribeRulesResponse

		request.MaxResults = &limit

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, err := me.client.UseAlbV20251030Client().DescribeRules(request)
			if err != nil {
				return tccommon.RetryError(err)
			}

			response = result
			return nil
		})
		if err != nil {
			return nil, err
		}

		if response == nil || response.Response == nil {
			return nil, fmt.Errorf("DescribeRules returned an empty response")
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if len(response.Response.Rules) < 1 {
			break
		}

		rules = append(rules, response.Response.Rules...)

		if response.Response.NextToken == nil || *response.Response.NextToken == "" {
			break
		}

		nextToken := *response.Response.NextToken
		request.NextToken = &nextToken
	}

	return
}

// ModifyRulesAttributes updates the attributes of the given forwarding rules.
func (me *AlbService) ModifyRulesAttributes(ctx context.Context, request *albv20251030.ModifyRulesAttributesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *albv20251030.ModifyRulesAttributesResponse

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAlbV20251030Client().ModifyRulesAttributes(request)
		if err != nil {
			return tccommon.RetryError(err)
		}

		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("ModifyRulesAttributes returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// DeleteRules deletes the given forwarding rules.
func (me *AlbService) DeleteRules(ctx context.Context, request *albv20251030.DeleteRulesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *albv20251030.DeleteRulesResponse

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAlbV20251030Client().DeleteRules(request)
		if err != nil {
			return tccommon.RetryError(err)
		}

		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DeleteRules returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
