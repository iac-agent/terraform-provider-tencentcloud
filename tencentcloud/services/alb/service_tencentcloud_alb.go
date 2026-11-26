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

// DescribeListenerCertificates queries the certificates bound to an ALB
// listener via alb v20251030. The API pages with MaxResults/NextToken, so
// this method loops internally and returns the FULL list; the API's
// TotalCount is an internal paging detail and never part of the return value.
// Every API call inside the loop is wrapped in resource.Retry, with
// ratelimit.Check inside the retry closure and a [DEBUG] success log.
func (me *AlbService) DescribeListenerCertificates(ctx context.Context, request *albv20251030.DescribeListenerCertificatesRequest) (certificates []*albv20251030.CertificateInfo, err error) {
	var (
		logId     = tccommon.GetLogId(ctx)
		response  *albv20251030.DescribeListenerCertificatesResponse
		pageSize  = int64(100)
		nextToken string
	)

	// The named result is initialised here (NOT in the var block above): a
	// named result and a body-level var of the same name do not compile.
	certificates = make([]*albv20251030.CertificateInfo, 0)

	for {
		// The API caps MaxResults at 100; the SDK field is *uint64.
		maxResults := uint64(pageSize)
		request.MaxResults = &maxResults
		if nextToken != "" {
			request.NextToken = &nextToken
		} else {
			request.NextToken = nil
		}

		errRet := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseAlbV20251030Client().DescribeListenerCertificates(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("describe listener certificates failed, Response is nil."))
			}

			response = result
			return nil
		})

		if errRet != nil {
			log.Printf("[CRITAL]%s describe listener certificates failed, reason:%+v", logId, errRet)
			return nil, errRet
		}

		if len(response.Response.Certificates) < 1 {
			break
		}

		certificates = append(certificates, response.Response.Certificates...)

		// Token-based paging: a missing/empty NextToken is the last page;
		// a short page is also terminal.
		if response.Response.NextToken == nil || *response.Response.NextToken == "" {
			break
		}

		nextToken = *response.Response.NextToken

		if len(response.Response.Certificates) < int(pageSize) {
			break
		}
	}

	return
}
