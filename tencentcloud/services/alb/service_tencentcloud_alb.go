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

// NewAlbService creates the service wrapper for ALB API calls.
func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

// AlbService wraps ALB API calls with retry and logging.
type AlbService struct {
	client *connectivity.TencentCloudClient
}

// CreateRules creates forwarding rules and returns their rule IDs.
func (me *AlbService) CreateRules(ctx context.Context, request *albv20251030.CreateRulesRequest) (ruleIds []string, err error) {
	logId := tccommon.GetLogId(ctx)
	var response *albv20251030.CreateRulesResponse

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().CreateRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create rules failed, Response is nil."))
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create rules failed, reason:%+v", logId, err)
		return nil, err
	}

	for _, id := range response.Response.RuleIds {
		if id != nil {
			ruleIds = append(ruleIds, *id)
		}
	}

	return ruleIds, nil
}

// DescribeRules returns all forwarding rules matched by the request.
func (me *AlbService) DescribeRules(ctx context.Context, request *albv20251030.DescribeRulesRequest) (rules []*albv20251030.RuleOutput, err error) {
	logId := tccommon.GetLogId(ctx)
	var response *albv20251030.DescribeRulesResponse
	limit := int64(100)
	nextToken := ""

	rules = make([]*albv20251030.RuleOutput, 0)

	for {
		request.MaxResults = &limit
		if nextToken != "" {
			token := nextToken
			request.NextToken = &token
		}

		errRet := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseAlbV20251030Client().DescribeRules(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("describe rules failed, Response is nil."))
			}

			response = result
			return nil
		})
		if errRet != nil {
			log.Printf("[CRITAL]%s describe rules failed, reason:%+v", logId, errRet)
			return nil, errRet
		}

		if len(response.Response.Rules) < 1 {
			break
		}

		rules = append(rules, response.Response.Rules...)
		if response.Response.NextToken == nil || *response.Response.NextToken == "" {
			break
		}

		nextToken = *response.Response.NextToken
	}

	return rules, nil
}

// ModifyRulesAttributes updates forwarding rule attributes.
func (me *AlbService) ModifyRulesAttributes(ctx context.Context, request *albv20251030.ModifyRulesAttributesRequest) error {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().ModifyRulesAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s modify rules attributes failed, reason:%+v", logId, err)
	}

	return err
}

// DeleteRules deletes forwarding rules.
func (me *AlbService) DeleteRules(ctx context.Context, request *albv20251030.DeleteRulesRequest) error {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DeleteRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete rules failed, reason:%+v", logId, err)
	}

	return err
}
