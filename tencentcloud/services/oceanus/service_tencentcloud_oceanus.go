package oceanus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type OceanusService struct {
	client *connectivity.TencentCloudClient
}

func (me *OceanusService) DescribeOceanusResourceRelatedJobByFilter(ctx context.Context, param map[string]interface{}) (ResourceRelatedJob []*oceanus.ResourceRefJobInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeResourceRelatedJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ResourceId" {
			request.ResourceId = v.(*string)
		}

		if k == "DESCByJobConfigCreateTime" {
			request.DESCByJobConfigCreateTime = v.(*int64)
		}

		if k == "ResourceConfigVersion" {
			request.ResourceConfigVersion = v.(*int64)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeResourceRelatedJobs(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RefJobInfos) < 1 {
			break
		}

		ResourceRelatedJob = append(ResourceRelatedJob, response.Response.RefJobInfos...)
		if len(response.Response.RefJobInfos) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusSavepointListByFilter(ctx context.Context, param map[string]interface{}) (savepointList []*oceanus.Savepoint, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeJobSavepointRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "JobId" {
			request.JobId = v.(*string)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeJobSavepoint(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Savepoint) < 1 {
			break
		}

		savepointList = append(savepointList, response.Response.Savepoint...)
		if len(response.Response.Savepoint) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusSystemResourceByFilter(ctx context.Context, param map[string]interface{}) (SystemResource []*oceanus.SystemResourceItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeSystemResourcesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ResourceIds" {
			request.ResourceIds = v.([]*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*oceanus.Filter)
		}

		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}

		if k == "FlinkVersion" {
			request.FlinkVersion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeSystemResources(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ResourceSet) < 1 {
			break
		}

		SystemResource = append(SystemResource, response.Response.ResourceSet...)
		if len(response.Response.ResourceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusWorkSpacesByFilter(ctx context.Context, param map[string]interface{}) (WorkSpace []*oceanus.WorkSpaceSetItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeWorkSpacesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}

		if k == "Filters" {
			request.Filters = v.([]*oceanus.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeWorkSpaces(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.WorkSpaceSetItem) < 1 {
			break
		}

		WorkSpace = append(WorkSpace, response.Response.WorkSpaceSetItem...)
		if len(response.Response.WorkSpaceSetItem) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusJobById(ctx context.Context, jobId string) (Job *oceanus.JobV1, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeJobsRequest()
	request.JobIds = common.StringPtrs([]string{jobId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeJobs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	Job = response.Response.JobSet[0]
	return
}

func (me *OceanusService) DeleteOceanusJobById(ctx context.Context, jobId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteJobsRequest()
	request.JobIds = common.StringPtrs([]string{jobId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteJobs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusResourceById(ctx context.Context, resourceId string) (resource *oceanus.ResourceItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeResourcesRequest()
	request.ResourceIds = common.StringPtrs([]string{resourceId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeResources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	resource = response.Response.ResourceSet[0]
	return
}

func (me *OceanusService) DeleteOceanusResourceById(ctx context.Context, resourceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteResourcesRequest()
	request.ResourceIds = common.StringPtrs([]string{resourceId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteResources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusResourceConfigById(ctx context.Context, resourceId, version string) (resourceConfig *oceanus.ResourceConfigItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeResourceConfigsRequest()
	request.ResourceId = &resourceId
	versionInt, _ := strconv.ParseInt(version, 10, 64)
	request.ResourceConfigVersions = common.Int64Ptrs([]int64{versionInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeResourceConfigs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	resourceConfig = response.Response.ResourceConfigSet[0]
	return
}

func (me *OceanusService) DeleteOceanusResourceConfigById(ctx context.Context, resourceId, version string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteResourceConfigsRequest()
	request.ResourceId = &resourceId
	versionInt, _ := strconv.ParseInt(version, 10, 64)
	request.ResourceConfigVersions = common.Int64Ptrs([]int64{versionInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteResourceConfigs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusClustersByFilter(ctx context.Context, param map[string]interface{}) (Clusters []*oceanus.Cluster, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeClustersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterIds" {
			request.ClusterIds = v.([]*string)
		}

		if k == "OrderType" {
			request.OrderType = v.(*int64)
		}

		if k == "Filters" {
			request.Filters = v.([]*oceanus.Filter)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOceanusClient().DescribeClusters(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClusterSet) < 1 {
			break
		}

		Clusters = append(Clusters, response.Response.ClusterSet...)
		if len(response.Response.ClusterSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OceanusService) DescribeOceanusTreeJobsByFilter(ctx context.Context, param map[string]interface{}) (treeJobs *oceanus.DescribeTreeJobsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeTreeJobsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*oceanus.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeTreeJobs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	treeJobs = response.Response
	return
}

func (me *OceanusService) DescribeOceanusTreeResourcesByFilter(ctx context.Context, param map[string]interface{}) (treeResources *oceanus.DescribeTreeResourcesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeTreeResourcesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeTreeResources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	treeResources = response.Response
	return
}

func (me *OceanusService) DescribeOceanusJobSubmissionLogByFilter(ctx context.Context, param map[string]interface{}) (jobSubmissionLog *oceanus.DescribeJobSubmissionLogResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeJobSubmissionLogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "JobId" {
			request.JobId = v.(*string)
		}

		if k == "StartTime" {
			request.StartTime = v.(*int64)
		}

		if k == "EndTime" {
			request.EndTime = v.(*int64)
		}

		if k == "RunningOrderId" {
			request.RunningOrderId = v.(*int64)
		}

		if k == "Keyword" {
			request.Keyword = v.(*string)
		}

		if k == "OrderType" {
			request.OrderType = v.(*string)
		}

		if k == "Cursor" {
			request.Cursor = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	request.Limit = common.Int64Ptr(100)
	response, err := me.client.UseOceanusClient().DescribeJobSubmissionLog(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	jobSubmissionLog = response.Response

	return
}

func (me *OceanusService) DescribeOceanusCheckSavepointByFilter(ctx context.Context, param map[string]interface{}) (CheckSavepoint *oceanus.CheckSavepointResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewCheckSavepointRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "JobId" {
			request.JobId = v.(*string)
		}

		if k == "SerialId" {
			request.SerialId = v.(*string)
		}

		if k == "RecordType" {
			request.RecordType = v.(*int64)
		}

		if k == "SavepointPath" {
			request.SavepointPath = v.(*string)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().CheckSavepoint(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	CheckSavepoint = response.Response
	return
}

func (me *OceanusService) DescribeOceanusWorkSpaceById(ctx context.Context, workSpaceName string) (WorkSpace *oceanus.WorkSpaceSetItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeWorkSpacesRequest()
	request.Filters = []*oceanus.Filter{
		{
			Name:   common.StringPtr("WorkSpaceName"),
			Values: common.StringPtrs([]string{workSpaceName}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeWorkSpaces(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.WorkSpaceSetItem) != 1 {
		return
	}

	WorkSpace = response.Response.WorkSpaceSetItem[0]
	return
}

func (me *OceanusService) DeleteOceanusWorkSpaceById(ctx context.Context, workSpaceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteWorkSpaceRequest()
	request.WorkSpaceId = &workSpaceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteWorkSpace(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusJobConfigById(ctx context.Context, jobId, version string) (JobConfig *oceanus.JobConfig, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeJobConfigsRequest()
	request.JobId = &jobId
	versionInt, _ := strconv.ParseUint(version, 10, 64)
	request.JobConfigVersions = common.Uint64Ptrs([]uint64{versionInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeJobConfigs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.JobConfigSet) != 1 {
		return
	}

	JobConfig = response.Response.JobConfigSet[0]
	return
}

func (me *OceanusService) DeleteOceanusJobConfigById(ctx context.Context, jobId, version string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteJobConfigsRequest()
	request.JobId = &jobId
	versionInt, _ := strconv.ParseInt(version, 10, 64)
	request.JobConfigVersions = common.Int64Ptrs([]int64{versionInt})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteJobConfigs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusFolderById(ctx context.Context, workSpaceId, folderId, folderType string) (Folder *oceanus.DescribeFolderResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeFolderRequest()
	request.FolderId = common.StringPtr(folderId)
	folderTypeInt, _ := strconv.ParseInt(folderType, 10, 64)
	request.FolderType = common.Int64Ptr(folderTypeInt)
	request.WorkSpaceId = common.StringPtr(workSpaceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeFolder(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	Folder = response.Response
	return
}

func (me *OceanusService) DeleteOceanusFolderById(ctx context.Context, workSpaceId, folderId, folderType string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDeleteFoldersRequest()
	request.FolderIds = common.StringPtrs([]string{folderId})
	folderTypeInt, _ := strconv.ParseInt(folderType, 10, 64)
	request.FolderType = common.Int64Ptr(folderTypeInt)
	request.WorkSpaceId = common.StringPtr(workSpaceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DeleteFolders(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OceanusService) DescribeOceanusJobEventsByFilter(ctx context.Context, param map[string]interface{}) (JobEvents *oceanus.DescribeJobEventsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewDescribeJobEventsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "JobId" {
			request.JobId = v.(*string)
		}

		if k == "StartTimestamp" {
			request.StartTimestamp = v.(*uint64)
		}

		if k == "EndTimestamp" {
			request.EndTimestamp = v.(*uint64)
		}

		if k == "Types" {
			request.Types = v.([]*string)
		}

		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().DescribeJobEvents(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	JobEvents = response.Response
	return
}

func (me *OceanusService) DescribeOceanusMetaTableByFilter(ctx context.Context, param map[string]interface{}) (MetaTable *oceanus.GetMetaTableResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = oceanus.NewGetMetaTableRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "WorkSpaceId" {
			request.WorkSpaceId = v.(*string)
		}

		if k == "Catalog" {
			request.Catalog = v.(*string)
		}

		if k == "Database" {
			request.Database = v.(*string)
		}

		if k == "Table" {
			request.Table = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOceanusClient().GetMetaTable(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	MetaTable = response.Response
	return
}

// CreateOceanusClusterFramework creates an Oceanus cluster and returns the
// created cluster id (framework reference implementation; the SDKv2-style
// filter methods above share the same service struct).
func (me *OceanusService) CreateOceanusCluster(ctx context.Context, request *oceanus.CreateOceanusClusterRequest) (clusterId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseOceanusClient().CreateOceanusCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create oceanus cluster failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result.Response.ClusterId == nil || *result.Response.ClusterId == "" {
			return resource.NonRetryableError(fmt.Errorf("create oceanus cluster failed, ClusterId is empty"))
		}

		clusterId = *result.Response.ClusterId
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// DescribeClustersFramework lists Oceanus clusters with built-in paging. A
// unique ClusterIds filter terminates after the first round; a broad query
// keeps paging until the full set is fetched. The response TotalCount is a
// paging detail and is never part of the return value.
func (me *OceanusService) DescribeClustersFramework(ctx context.Context, request *oceanus.DescribeClustersRequest) (Clusters []*oceanus.Cluster, errRet error) {
	logId := tccommon.GetLogId(ctx)

	Clusters = make([]*oceanus.Cluster, 0)

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		pageClusters, err := func() ([]*oceanus.Cluster, error) {
			var page []*oceanus.Cluster
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				ratelimit.Check(request.GetAction())
				result, e := me.client.UseOceanusClient().DescribeClusters(request)
				if e != nil {
					return tccommon.RetryError(e)
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("describe oceanus clusters failed, Response is nil"))
				}

				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

				page = result.Response.ClusterSet
				return nil
			})
			return page, err
		}()
		if err != nil {
			log.Printf("[CRITAL]%s describe oceanus clusters failed, reason:%+v", logId, err)
			errRet = err
			return
		}

		Clusters = append(Clusters, pageClusters...)

		if len(pageClusters) < 1 {
			break
		}

		if len(pageClusters) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

// ScaleOceanusClusterFramework resizes an Oceanus cluster (framework reference
// implementation). ScaleOceanusCluster is synchronous from the API's point of
// view: the cluster CU change is applied asynchronously on the cloud side and
// is observed through DescribeClusters, so no task id is returned here.
func (me *OceanusService) ScaleOceanusClusterFramework(ctx context.Context, request *oceanus.ScaleOceanusClusterRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseOceanusClient().ScaleOceanusCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("scale oceanus cluster failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// DeleteOceanusClusterFramework destroys an Oceanus cluster (framework
// reference implementation).
func (me *OceanusService) DeleteOceanusClusterFramework(ctx context.Context, request *oceanus.DeleteOceanusClusterRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseOceanusClient().DeleteOceanusCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete oceanus cluster failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

// WaitOceanusClusterDestroyed polls DescribeClusters until the cluster is no
// longer returned: an empty describe answer for the tracked cluster id is the
// terminal "destroyed" state. The helper talks to the SDK client directly, so
// an EMPTY answer is judged explicitly here.
func (me *OceanusService) WaitOceanusClusterDestroyed(ctx context.Context, clusterId string, timeout time.Duration) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeClustersRequest()
	request.ClusterIds = []*string{&clusterId}

	err := resource.Retry(timeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseOceanusClient().DescribeClusters(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if response == nil || response.Response == nil || len(response.Response.ClusterSet) == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("oceanus cluster %s still exists, waiting for destroyed", clusterId))
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, cluster[%s], reason[%s]\n", logId, request.GetAction(), clusterId, err.Error())
		errRet = err
	}

	return
}

// WaitOceanusClusterCreated polls DescribeClusters until the cluster created by
// CreateOceanusCluster becomes visible (initially the cluster may not yet be
// returned by the describe API). A "cluster visible" signal is judged
// structurally; no SDK error code is special-cased here.
func (me *OceanusService) WaitOceanusClusterCreated(ctx context.Context, clusterId string, timeout time.Duration) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeClustersRequest()
	request.ClusterIds = []*string{&clusterId}

	err := resource.Retry(timeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseOceanusClient().DescribeClusters(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if response == nil || response.Response == nil || len(response.Response.ClusterSet) == 0 {
			return resource.RetryableError(fmt.Errorf("oceanus cluster %s not visible yet, waiting for created", clusterId))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, cluster[%s], reason[%s]\n", logId, request.GetAction(), clusterId, err.Error())
		errRet = err
	}

	return
}

// WaitOceanusClusterScaled polls DescribeClusters until the cluster CU count
// differs from the pre-scale value. The SDK exposes no terminal state signal
// for ScaleOceanusCluster beyond the describe answer itself.
// TODO(human): confirm the cloud-side terminal semantics of a CU resize; the
// current signal ("CuNum changed") cannot detect a resize that lands on the
// same value and keeps polling until the timeout.
func (me *OceanusService) WaitOceanusClusterScaled(ctx context.Context, clusterId string, previousCU int64, timeout time.Duration) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := oceanus.NewDescribeClustersRequest()
	request.ClusterIds = []*string{&clusterId}

	err := resource.Retry(timeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseOceanusClient().DescribeClusters(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if response == nil || response.Response == nil || len(response.Response.ClusterSet) == 0 {
			return resource.RetryableError(fmt.Errorf("oceanus cluster %s not visible during scaling, waiting for scale", clusterId))
		}

		cluster := response.Response.ClusterSet[0]
		if cluster != nil && cluster.CuNum != nil && *cluster.CuNum != uint64(previousCU) {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("oceanus cluster %s still reports %d CUs, waiting for scale", clusterId, previousCU))
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, cluster[%s], reason[%s]\n", logId, request.GetAction(), clusterId, err.Error())
		errRet = err
	}

	return
}
