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

// NewAlbService creates the service wrapper for ALB framework references.
func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

// AlbService wraps ALB SDK calls with retry and logging.
type AlbService struct {
	client *connectivity.TencentCloudClient
}

// DescribeListenerCertificates returns the full list of certificates bound to a
// listener, fetching all pages internally via the NextToken cursor.
func (me *AlbService) DescribeListenerCertificates(ctx context.Context, request *albv20251030.DescribeListenerCertificatesRequest) (certificates []*albv20251030.CertificateInfo, err error) {
	var (
		logId      = tccommon.GetLogId(ctx)
		response   *albv20251030.DescribeListenerCertificatesResponse
		maxResults = uint64(100)
	)

	certificates = make([]*albv20251030.CertificateInfo, 0)

	for {
		request.MaxResults = &maxResults
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

		certificates = append(certificates, response.Response.Certificates...)

		// Terminate when there is no next token, or the page is empty.
		if response.Response.NextToken == nil || *response.Response.NextToken == "" || len(response.Response.Certificates) < 1 {
			break
		}

		request.NextToken = response.Response.NextToken
	}

	return
}
