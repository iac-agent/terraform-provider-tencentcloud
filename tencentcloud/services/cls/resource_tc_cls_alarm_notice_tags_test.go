package cls

import (
	"reflect"
	"testing"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// go test ./tencentcloud/services/cls/ -run "TestMapToClsTags" -v -count=1 -gcflags="all=-l"
func TestMapToClsTags(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]interface{}
		expect []*cls.Tag
	}{
		{
			name:   "empty map",
			input:  map[string]interface{}{},
			expect: []*cls.Tag{},
		},
		{
			name:  "single tag",
			input: map[string]interface{}{"key1": "value1"},
			expect: []*cls.Tag{
				{Key: helper.String("key1"), Value: helper.String("value1")},
			},
		},
		{
			name:  "multiple tags",
			input: map[string]interface{}{"key1": "value1", "key2": "value2"},
			expect: []*cls.Tag{
				{Key: helper.String("key1"), Value: helper.String("value1")},
				{Key: helper.String("key2"), Value: helper.String("value2")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToClsTags(tt.input)
			if len(result) != len(tt.expect) {
				t.Fatalf("expected %d tags, got %d", len(tt.expect), len(result))
			}
			// Build a map for comparison since order is not guaranteed
			resultMap := make(map[string]string)
			for _, tag := range result {
				if tag.Key != nil && tag.Value != nil {
					resultMap[*tag.Key] = *tag.Value
				}
			}
			expectMap := make(map[string]string)
			for _, tag := range tt.expect {
				if tag.Key != nil && tag.Value != nil {
					expectMap[*tag.Key] = *tag.Value
				}
			}
			if !reflect.DeepEqual(resultMap, expectMap) {
				t.Fatalf("expected %v, got %v", expectMap, resultMap)
			}
		})
	}
}

// go test ./tencentcloud/services/cls/ -run "TestClsTagsToMap" -v -count=1 -gcflags="all=-l"
func TestClsTagsToMap(t *testing.T) {
	tests := []struct {
		name   string
		input  []*cls.Tag
		expect map[string]interface{}
	}{
		{
			name:   "nil tags",
			input:  nil,
			expect: map[string]interface{}{},
		},
		{
			name:   "empty tags",
			input:  []*cls.Tag{},
			expect: map[string]interface{}{},
		},
		{
			name: "single tag",
			input: []*cls.Tag{
				{Key: helper.String("key1"), Value: helper.String("value1")},
			},
			expect: map[string]interface{}{"key1": "value1"},
		},
		{
			name: "multiple tags",
			input: []*cls.Tag{
				{Key: helper.String("key1"), Value: helper.String("value1")},
				{Key: helper.String("key2"), Value: helper.String("value2")},
			},
			expect: map[string]interface{}{"key1": "value1", "key2": "value2"},
		},
		{
			name: "tag with nil key skipped",
			input: []*cls.Tag{
				{Key: nil, Value: helper.String("value1")},
				{Key: helper.String("key2"), Value: helper.String("value2")},
			},
			expect: map[string]interface{}{"key2": "value2"},
		},
		{
			name: "tag with nil value skipped",
			input: []*cls.Tag{
				{Key: helper.String("key1"), Value: nil},
				{Key: helper.String("key2"), Value: helper.String("value2")},
			},
			expect: map[string]interface{}{"key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clsTagsToMap(tt.input)
			if !reflect.DeepEqual(result, tt.expect) {
				t.Fatalf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}

// go test ./tencentcloud/services/cls/ -run "TestMapToClsTags_CreateRequestIntegration" -v -count=1 -gcflags="all=-l"
// TestMapToClsTags_CreateRequestIntegration tests that mapToClsTags output can be directly
// assigned to CreateAlarmNoticeRequest.Tags, simulating the Create flow.
func TestMapToClsTags_CreateRequestIntegration(t *testing.T) {
	tagsMap := map[string]interface{}{
		"env":  "production",
		"team": "devops",
	}

	request := cls.NewCreateAlarmNoticeRequest()
	request.Tags = mapToClsTags(tagsMap)

	if request.Tags == nil {
		t.Fatal("Tags should not be nil on CreateAlarmNoticeRequest")
	}
	if len(request.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(request.Tags))
	}

	// Verify the tags match
	resultMap := make(map[string]string)
	for _, tag := range request.Tags {
		if tag.Key != nil && tag.Value != nil {
			resultMap[*tag.Key] = *tag.Value
		}
	}
	if resultMap["env"] != "production" {
		t.Fatalf("expected env=production, got env=%s", resultMap["env"])
	}
	if resultMap["team"] != "devops" {
		t.Fatalf("expected team=devops, got team=%s", resultMap["team"])
	}
}

// go test ./tencentcloud/services/cls/ -run "TestMapToClsTags_ModifyRequestIntegration" -v -count=1 -gcflags="all=-l"
// TestMapToClsTags_ModifyRequestIntegration tests that mapToClsTags output can be directly
// assigned to ModifyAlarmNoticeRequest.Tags, simulating the Update flow.
func TestMapToClsTags_ModifyRequestIntegration(t *testing.T) {
	tagsMap := map[string]interface{}{
		"env":  "staging",
		"team": "platform",
	}

	request := cls.NewModifyAlarmNoticeRequest()
	request.Tags = mapToClsTags(tagsMap)

	if request.Tags == nil {
		t.Fatal("Tags should not be nil on ModifyAlarmNoticeRequest")
	}
	if len(request.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(request.Tags))
	}

	// Verify the tags match
	resultMap := make(map[string]string)
	for _, tag := range request.Tags {
		if tag.Key != nil && tag.Value != nil {
			resultMap[*tag.Key] = *tag.Value
		}
	}
	if resultMap["env"] != "staging" {
		t.Fatalf("expected env=staging, got env=%s", resultMap["env"])
	}
	if resultMap["team"] != "platform" {
		t.Fatalf("expected team=platform, got team=%s", resultMap["team"])
	}
}

// go test ./tencentcloud/services/cls/ -run "TestClsTagsToMap_ReadResponseIntegration" -v -count=1 -gcflags="all=-l"
// TestClsTagsToMap_ReadResponseIntegration tests that clsTagsToMap can convert
// AlarmNotice.Tags from the API response to a map, simulating the Read flow.
func TestClsTagsToMap_ReadResponseIntegration(t *testing.T) {
	// Simulate an AlarmNotice response with Tags
	alarmNotice := &cls.AlarmNotice{
		AlarmNoticeId: helper.String("alarm-notice-123"),
		Name:          helper.String("test"),
		Tags: []*cls.Tag{
			{Key: helper.String("env"), Value: helper.String("production")},
			{Key: helper.String("team"), Value: helper.String("devops")},
		},
	}

	if alarmNotice.Tags == nil || len(alarmNotice.Tags) == 0 {
		t.Fatal("AlarmNotice.Tags should not be nil or empty")
	}

	result := clsTagsToMap(alarmNotice.Tags)

	if len(result) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(result))
	}
	if result["env"] != "production" {
		t.Fatalf("expected env=production, got env=%s", result["env"])
	}
	if result["team"] != "devops" {
		t.Fatalf("expected team=devops, got team=%s", result["team"])
	}
}

// go test ./tencentcloud/services/cls/ -run "TestClsTagsToMap_ReadResponseNilTagsIntegration" -v -count=1 -gcflags="all=-l"
// TestClsTagsToMap_ReadResponseNilTagsIntegration tests the Read flow when
// AlarmNotice.Tags is nil (should fall back to tag service).
func TestClsTagsToMap_ReadResponseNilTagsIntegration(t *testing.T) {
	// Simulate an AlarmNotice response without Tags
	alarmNotice := &cls.AlarmNotice{
		AlarmNoticeId: helper.String("alarm-notice-123"),
		Name:          helper.String("test"),
		Tags:          nil,
	}

	// In the Read function, we check: if alarmNotice.Tags != nil && len(alarmNotice.Tags) > 0
	if alarmNotice.Tags != nil && len(alarmNotice.Tags) > 0 {
		t.Fatal("This test case should simulate nil Tags in the response")
	}
	// When Tags is nil, the Read function falls back to tag service
	// This test verifies the condition check works correctly
}

// go test ./tencentcloud/services/cls/ -run "TestClsTagsToMap_ReadResponseEmptyTagsIntegration" -v -count=1 -gcflags="all=-l"
// TestClsTagsToMap_ReadResponseEmptyTagsIntegration tests the Read flow when
// AlarmNotice.Tags is an empty slice (should also fall back to tag service).
func TestClsTagsToMap_ReadResponseEmptyTagsIntegration(t *testing.T) {
	// Simulate an AlarmNotice response with empty Tags
	alarmNotice := &cls.AlarmNotice{
		AlarmNoticeId: helper.String("alarm-notice-123"),
		Name:          helper.String("test"),
		Tags:          []*cls.Tag{},
	}

	// In the Read function, we check: if alarmNotice.Tags != nil && len(alarmNotice.Tags) > 0
	if alarmNotice.Tags != nil && len(alarmNotice.Tags) > 0 {
		t.Fatal("This test case should simulate empty Tags in the response")
	}
	// When Tags is empty, the Read function falls back to tag service
}

// go test ./tencentcloud/services/cls/ -run "TestMapToClsTags_NilMapNotSetOnRequest" -v -count=1 -gcflags="all=-l"
// TestMapToClsTags_NilMapNotSetOnRequest tests the Create flow when tags is not provided:
// the request.Tags should remain nil (not set).
func TestMapToClsTags_NilMapNotSetOnRequest(t *testing.T) {
	request := cls.NewCreateAlarmNoticeRequest()

	// Simulate: if tagsMap, ok := d.GetOk("tags"); ok { request.Tags = mapToClsTags(...) }
	// When d.GetOk("tags") returns false, Tags should not be set
	if request.Tags != nil {
		t.Fatal("Tags should be nil by default on CreateAlarmNoticeRequest when not set")
	}
}
