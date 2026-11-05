package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

type mockMetaFwmAccountGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmAccountGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmAccountGroup{}

func newMockMetaFwmAccountGroup() *mockMetaFwmAccountGroup {
	return &mockMetaFwmAccountGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFAG(s string) *string {
	return &s
}

func TestFwmAccountGroup_Schema(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmAccountGroup()

	expected := map[string]struct {
		required bool
		optional bool
		computed bool
		forceNew bool
	}{
		"group_name":   {required: true},
		"description":  {optional: true, forceNew: true},
		"group_id":     {computed: true},
		"member_count": {computed: true},
		"can_delete":   {computed: true},
		"create_time":  {computed: true},
	}

	if len(res.Schema) != len(expected) {
		t.Fatalf("unexpected schema field count: got %d, want %d", len(res.Schema), len(expected))
	}

	for name, want := range expected {
		s, ok := res.Schema[name]
		if !ok {
			t.Fatalf("schema field %q is missing", name)
		}
		if s.Required != want.required {
			t.Errorf("schema field %q Required = %v, want %v", name, s.Required, want.required)
		}
		if s.Optional != want.optional {
			t.Errorf("schema field %q Optional = %v, want %v", name, s.Optional, want.optional)
		}
		if s.Computed != want.computed {
			t.Errorf("schema field %q Computed = %v, want %v", name, s.Computed, want.computed)
		}
		if s.ForceNew != want.forceNew {
			t.Errorf("schema field %q ForceNew = %v, want %v", name, s.ForceNew, want.forceNew)
		}
	}
}

func TestFwmAccountGroup_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(0)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group"),
					Description: ptrStringFAG("test account group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
	assert.Equal(t, "test-account-group", d.Get("group_name"))
	assert.Equal(t, "test account group", d.Get("description"))
	assert.Equal(t, 0, d.Get("member_count"))
	assert.Equal(t, true, d.Get("can_delete"))
	assert.Equal(t, "2026-10-22 10:00:00", d.Get("create_time"))
}

func TestFwmAccountGroup_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(0)
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:     ptrStringFAG("fake-request-id"),
			TotalCount:    &total,
			AccountGroups: []*fwmv20250611.AccountGroup{},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmAccountGroup_Create_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateAccountGroupRequest) (*fwmv20250611.CreateAccountGroupResponse, error) {
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "test-account-group", *request.GroupName)
		assert.NotNil(t, request.Description)
		assert.Equal(t, "test account group", *request.Description)

		resp := fwmv20250611.NewCreateAccountGroupResponse()
		resp.Response = &fwmv20250611.CreateAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-create"),
			GroupId:   ptrStringFAG("acg-1a2b3c4d"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(0)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group"),
					Description: ptrStringFAG("test account group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name":  "test-account-group",
		"description": "test account group",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
	assert.Equal(t, "test-account-group", d.Get("group_name"))
	assert.Equal(t, "test account group", d.Get("description"))
}

func TestFwmAccountGroup_Create_EmptyGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateAccountGroupRequest) (*fwmv20250611.CreateAccountGroupResponse, error) {
		resp := fwmv20250611.NewCreateAccountGroupResponse()
		resp.Response = &fwmv20250611.CreateAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-create"),
			GroupId:   nil,
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name":  "test-account-group",
		"description": "test account group",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

func TestFwmAccountGroup_Update_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyAccountGroupRequest) (*fwmv20250611.ModifyAccountGroupResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "test-account-group-updated", *request.GroupName)

		resp := fwmv20250611.NewModifyAccountGroupResponse()
		resp.Response = &fwmv20250611.ModifyAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(0)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group-updated"),
					Description: ptrStringFAG("test account group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name":  "test-account-group-updated",
		"description": "test account group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-account-group-updated", d.Get("group_name"))
}

func TestFwmAccountGroup_Delete_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteAccountGroupRequest) (*fwmv20250611.DeleteAccountGroupResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)

		resp := fwmv20250611.NewDeleteAccountGroupResponse()
		resp.Response = &fwmv20250611.DeleteAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestFwmAccountGroup_Delete_ErrorPropagated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteAccountGroupRequest) (*fwmv20250611.DeleteAccountGroupResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}

func TestFwmAccountGroup_Read_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
}

func TestFwmAccountGroup_Read_NilDescription(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(3)
		canDelete := false
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group"),
					Description: nil,
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name":  "test-account-group",
		"description": "will not be overwritten when API returns null",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
	assert.Equal(t, "test-account-group", d.Get("group_name"))
	assert.Equal(t, 3, d.Get("member_count"))
	assert.Equal(t, false, d.Get("can_delete"))
	assert.Equal(t, "2026-10-22 10:00:00", d.Get("create_time"))
}

func TestFwmAccountGroup_Create_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateAccountGroupRequest) (*fwmv20250611.CreateAccountGroupResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmAccountGroup_Create_WithoutDescription(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateAccountGroupRequest) (*fwmv20250611.CreateAccountGroupResponse, error) {
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "test-account-group", *request.GroupName)
		assert.Nil(t, request.Description)

		resp := fwmv20250611.NewCreateAccountGroupResponse()
		resp.Response = &fwmv20250611.CreateAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-create"),
			GroupId:   ptrStringFAG("acg-1a2b3c4d"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(0)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
	assert.Equal(t, "test-account-group", d.Get("group_name"))
}

func TestFwmAccountGroup_Update_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyAccountGroupRequest) (*fwmv20250611.ModifyAccountGroupResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	if err := d.Set("group_name", "test-account-group-updated"); err != nil {
		t.Fatalf("failed to set group_name: %v", err)
	}

	err := res.Update(d, meta)
	assert.Error(t, err)
}

func TestFwmAccountGroup_Update_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyAccountGroupRequest) (*fwmv20250611.ModifyAccountGroupResponse, error) {
		return &fwmv20250611.ModifyAccountGroupResponse{}, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	if err := d.Set("group_name", "test-account-group-updated"); err != nil {
		t.Fatalf("failed to set group_name: %v", err)
	}

	err := res.Update(d, meta)
	assert.Error(t, err)
}

func TestFwmAccountGroup_Delete_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteAccountGroupRequest) (*fwmv20250611.DeleteAccountGroupResponse, error) {
		return &fwmv20250611.DeleteAccountGroupResponse{}, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}

func TestFwmAccountGroup_Read_PicksMatchingGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(2)
		memberCount := int64(2)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-other"),
					GroupName:   ptrStringFAG("other-account-group"),
					Description: ptrStringFAG("other group description"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-20 10:00:00"),
				},
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group"),
					Description: ptrStringFAG("test account group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
	assert.Equal(t, "acg-1a2b3c4d", d.Get("group_id"))
	assert.Equal(t, "test-account-group", d.Get("group_name"))
	assert.Equal(t, "test account group", d.Get("description"))
}

func TestFwmAccountGroup_Read_MismatchedGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(1)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-other"),
					GroupName:   ptrStringFAG("other-account-group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-20 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
	assert.Equal(t, "", d.Get("group_id"))
}

func TestFwmAccountGroup_Read_Pagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	var offsets []int64
	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		if request.Offset == nil || request.Limit == nil {
			t.Fatalf("offset and limit must be set on the describe request")
		}
		offsets = append(offsets, *request.Offset)

		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(21)
		memberCount := int64(1)
		canDelete := true

		items := make([]*fwmv20250611.AccountGroup, 0, 20)
		if *request.Offset == 0 {
			for i := 0; i < 20; i++ {
				items = append(items, &fwmv20250611.AccountGroup{
					GroupId:     ptrStringFAG("acg-filler"),
					GroupName:   ptrStringFAG("filler-account-group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-20 10:00:00"),
				})
			}
		} else {
			items = append(items, &fwmv20250611.AccountGroup{
				GroupId:     ptrStringFAG("acg-1a2b3c4d"),
				GroupName:   ptrStringFAG("test-account-group"),
				Description: ptrStringFAG("test account group"),
				MemberCount: &memberCount,
				CanDelete:   &canDelete,
				CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
			})
		}

		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:     ptrStringFAG("fake-request-id"),
			TotalCount:    &total,
			AccountGroups: items,
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
	assert.Equal(t, "acg-1a2b3c4d", d.Get("group_id"))
	assert.Equal(t, "test-account-group", d.Get("group_name"))
	assert.Equal(t, "test account group", d.Get("description"))
	if len(offsets) != 2 {
		t.Fatalf("expected 2 describe calls while paginating, got %d", len(offsets))
	}
	assert.Equal(t, int64(0), offsets[0])
	assert.Equal(t, int64(20), offsets[1])
}

func TestFwmAccountGroup_Read_StopsOnTotalCount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	calls := 0
	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		calls++
		if calls > 3 {
			t.Fatalf("pagination did not stop, describe called %d times", calls)
		}

		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(20)
		memberCount := int64(1)
		canDelete := true

		items := make([]*fwmv20250611.AccountGroup, 0, 20)
		for i := 0; i < 20; i++ {
			items = append(items, &fwmv20250611.AccountGroup{
				GroupId:     ptrStringFAG("acg-filler"),
				GroupName:   ptrStringFAG("filler-account-group"),
				MemberCount: &memberCount,
				CanDelete:   &canDelete,
				CreateTime:  ptrStringFAG("2026-10-20 10:00:00"),
			})
		}

		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:     ptrStringFAG("fake-request-id"),
			TotalCount:    &total,
			AccountGroups: items,
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
	if calls != 1 {
		t.Fatalf("expected pagination to stop after the first page when TotalCount is reached, got %d call(s)", calls)
	}
}
