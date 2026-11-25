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

// NewAlbService creates the ALB service used by the framework references.
func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

// AlbService wraps ALB SDK calls with retry, rate limiting and logging.
type AlbService struct {
	client *connectivity.TencentCloudClient
}

// CreateRules creates forwarding rules and returns the created rule id list.
func (me *AlbService) CreateRules(ctx context.Context, request *albv20251030.CreateRulesRequest) (ruleIds []*string, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.CreateRulesResponse
	)

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

	return response.Response.RuleIds, nil
}

// DescribeRules returns the full list of rules matching the request, with
// pagination handled internally. TotalCount is an internal paging detail.
func (me *AlbService) DescribeRules(ctx context.Context, request *albv20251030.DescribeRulesRequest) (rules []*albv20251030.RuleOutput, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.DescribeRulesResponse
		limit    = int64(100)
		token    string
	)

	rules = make([]*albv20251030.RuleOutput, 0)

	for {
		request.MaxResults = &limit
		if token != "" {
			request.NextToken = &token
		} else {
			request.NextToken = nil
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

		rules = append(rules, response.Response.Rules...)
		if response.Response.NextToken == nil || *response.Response.NextToken == "" {
			break
		}
		token = *response.Response.NextToken
	}

	return rules, nil
}

// ModifyRulesAttributes modifies attributes of forwarding rules.
func (me *AlbService) ModifyRulesAttributes(ctx context.Context, request *albv20251030.ModifyRulesAttributesRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().ModifyRulesAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("modify rules attributes failed, Response is nil."))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify rules attributes failed, reason:%+v", logId, err)
	}
	return err
}

// DeleteRules deletes forwarding rules by id.
func (me *AlbService) DeleteRules(ctx context.Context, request *albv20251030.DeleteRulesRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DeleteRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete rules failed, Response is nil."))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete rules failed, reason:%+v", logId, err)
	}
	return err
}
