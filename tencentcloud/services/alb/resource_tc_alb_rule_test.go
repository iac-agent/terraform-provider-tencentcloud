package alb

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tftypes "github.com/hashicorp/terraform-plugin-go/tftypes"

	alb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"
	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const (
	albRuleTestLbId       = "alb-12345678"
	albRuleTestListenerId = "lst-12345678"
	albRuleTestRuleId     = "rule-12345678"
	albRuleTestTargetId   = "lbtg-12345678"
	albRuleTestCreateTime = "2026-01-01T00:00:00Z"
)

func albRuleTestString(s string) *string { return &s }
func albRuleTestInt64(i int64) *int64    { return &i }
func albRuleTestBool(b bool) *bool       { return &b }

// albRuleTestContext wires a resource instance to a patched ALB SDK client.
type albRuleTestContext struct {
	resource resource.Resource
	client   *alb.Client
	schema   schema.Schema
}

func newAlbRuleTestContext(t *testing.T, patches *gomonkey.Patches) *albRuleTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	albClient := &alb.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseAlbV20251030Client", albClient)

	res := NewAlbRuleResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp,
	)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	return &albRuleTestContext{resource: res, client: albClient, schema: schemaResp.Schema}
}

func (c *albRuleTestContext) plan(t *testing.T, model AlbRuleResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *albRuleTestContext) state(t *testing.T, model AlbRuleResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

// runReadWithState mirrors the framework server: the response state starts as
// a copy of the current state, so a failing Read leaves the prior state intact.
func (c *albRuleTestContext) runReadWithState(ctx context.Context, currentState tfsdk.State) *resource.ReadResponse {
	resp := &resource.ReadResponse{State: tfsdk.State{
		Schema: currentState.Schema,
		Raw:    currentState.Raw.Copy(),
	}}
	c.resource.Read(ctx, resource.ReadRequest{State: currentState}, resp)
	return resp
}

func (c *albRuleTestContext) readModel(t *testing.T, state tfsdk.State) AlbRuleResourceModel {
	t.Helper()

	var model AlbRuleResourceModel
	if diags := state.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags.Errors())
	}
	return model
}

func albRuleTestObject(t *testing.T, attrTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	t.Helper()

	object, diags := types.ObjectValue(attrTypes, values)
	if diags.HasError() {
		t.Fatalf("failed to build object value: %v", diags.Errors())
	}
	return object
}

func albRuleTestList(t *testing.T, attrTypes map[string]attr.Type, elements ...types.Object) types.List {
	t.Helper()

	objType := types.ObjectType{AttrTypes: attrTypes}
	values := make([]attr.Value, 0, len(elements))
	for _, element := range elements {
		values = append(values, element)
	}
	list, diags := types.ListValue(objType, values)
	if diags.HasError() {
		t.Fatalf("failed to build list value: %v", diags.Errors())
	}
	return list
}

func albRuleTestSet(t *testing.T, values ...string) types.Set {
	t.Helper()

	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}
	set, diags := types.SetValue(types.StringType, elements)
	if diags.HasError() {
		t.Fatalf("failed to build set value: %v", diags.Errors())
	}
	return set
}

// albRuleTestFullRule builds one fully populated rule object.
func albRuleTestFullRule(t *testing.T) types.Object {
	t.Helper()

	fixedResponse := albRuleTestObject(t, albRuleFixedResponseAttrTypes(), map[string]attr.Value{
		"http_code":    types.Int64Value(200),
		"content":      types.StringValue("ok"),
		"content_type": types.StringValue("text/plain"),
	})
	insertHeader := albRuleTestObject(t, albRuleInsertHeaderAttrTypes(), map[string]attr.Value{
		"key":        types.StringValue("X-Test"),
		"value":      types.StringValue("ClientIp"),
		"value_type": types.StringValue("SystemDefined"),
	})
	redirect := albRuleTestObject(t, albRuleRedirectAttrTypes(), map[string]attr.Value{
		"http_code": types.Int64Value(301),
		"host":      types.StringValue("example.com"),
		"path":      types.StringValue("/new"),
		"port":      types.StringValue("443"),
		"protocol":  types.StringValue("HTTPS"),
		"query":     types.StringValue("a=b"),
	})
	removeHeader := albRuleTestObject(t, albRuleRemoveHeaderAttrTypes(), map[string]attr.Value{
		"key": types.StringValue("X-Remove"),
	})
	rewrite := albRuleTestObject(t, albRuleRewriteAttrTypes(), map[string]attr.Value{
		"host":  types.StringValue("example.com"),
		"path":  types.StringValue("/rewrite"),
		"query": types.StringValue("k=v"),
	})
	targetGroup := albRuleTestObject(t, albRuleTargetGroupAttrTypes(), map[string]attr.Value{
		"target_group_id": types.StringValue(albRuleTestTargetId),
		"weight":          types.Int64Value(30),
	})
	stickySession := albRuleTestObject(t, albRuleStickySessionAttrTypes(), map[string]attr.Value{
		"enabled": types.BoolValue(true),
		"timeout": types.Int64Value(1000),
	})
	targetGroupConfig := albRuleTestObject(t, albRuleTargetGroupConfigAttrTypes(), map[string]attr.Value{
		"target_groups":              albRuleTestList(t, albRuleTargetGroupAttrTypes(), targetGroup),
		"target_group_sticky_session": stickySession,
	})
	action := albRuleTestObject(t, albRuleActionAttrTypes(), map[string]attr.Value{
		"order":                 types.Int64Value(1),
		"type":                  types.StringValue("TargetGroup"),
		"fixed_response_config": fixedResponse,
		"insert_header_config":  insertHeader,
		"redirect_config":       redirect,
		"remove_header_config":  removeHeader,
		"rewrite_config":        rewrite,
		"target_group_config":   targetGroupConfig,
	})

	cookie := albRuleTestObject(t, albRuleCookieAttrTypes(), map[string]attr.Value{
		"key":   types.StringValue("session"),
		"value": types.StringValue("abc"),
	})
	header := albRuleTestObject(t, albRuleHeaderAttrTypes(), map[string]attr.Value{
		"key":    types.StringValue("X-Env"),
		"values": albRuleTestSet(t, "prod", "staging"),
	})
	queryString := albRuleTestObject(t, albRuleQueryStringAttrTypes(), map[string]attr.Value{
		"key":   types.StringValue("version"),
		"value": types.StringValue("v1"),
	})
	condition := albRuleTestObject(t, albRuleConditionAttrTypes(), map[string]attr.Value{
		"type":                types.StringValue("Path"),
		"cookie_config":       albRuleTestList(t, albRuleCookieAttrTypes(), cookie),
		"header_config":       header,
		"host_config":         albRuleTestSet(t, "example.com"),
		"method_config":       albRuleTestSet(t, "GET"),
		"path_config":         albRuleTestSet(t, "/api"),
		"query_string_config": albRuleTestList(t, albRuleQueryStringAttrTypes(), queryString),
		"source_ip_config":    albRuleTestSet(t, "10.0.0.1/32"),
	})
	tag := albRuleTestObject(t, albRuleTagAttrTypes(), map[string]attr.Value{
		"tag_key":   types.StringValue("env"),
		"tag_value": types.StringValue("test"),
	})

	return albRuleTestObject(t, albRuleInputAttrTypes(), map[string]attr.Value{
		"actions":    albRuleTestList(t, albRuleActionAttrTypes(), action),
		"conditions": albRuleTestList(t, albRuleConditionAttrTypes(), condition),
		"priority":   types.Int64Value(10),
		"direction":  types.StringValue("Request"),
		"rule_name":  types.StringValue("tf-example"),
		"tags":       albRuleTestList(t, albRuleTagAttrTypes(), tag),
	})
}

// albRuleTestCloudRule returns a full DescribeRules entry for the sample rule.
func albRuleTestCloudRule() *alb.RuleOutput {
	return &alb.RuleOutput{
		Actions: []*alb.RuleAction{
			{
				Order: albRuleTestInt64(1),
				Type:  albRuleTestString("TargetGroup"),
				FixedResponseConfig: &alb.FixedResponseInfo{
					HttpCode:    albRuleTestInt64(200),
					Content:     albRuleTestString("ok"),
					ContentType: albRuleTestString("text/plain"),
				},
				InsertHeaderConfig: &alb.InsertHTTPHeaderInfo{
					Key:       albRuleTestString("X-Test"),
					Value:     albRuleTestString("ClientIp"),
					ValueType: albRuleTestString("SystemDefined"),
				},
				RedirectConfig: &alb.HTTPRedirectInfo{
					HttpCode: albRuleTestInt64(301),
					Host:     albRuleTestString("example.com"),
					Path:     albRuleTestString("/new"),
					Port:     albRuleTestString("443"),
					Protocol: albRuleTestString("HTTPS"),
					Query:    albRuleTestString("a=b"),
				},
				RemoveHeaderConfig: &alb.RemoveHTTPHeaderInfo{Key: albRuleTestString("X-Remove")},
				RewriteConfig: &alb.HTTPRewriteInfo{
					Host:  albRuleTestString("example.com"),
					Path:  albRuleTestString("/rewrite"),
					Query: albRuleTestString("k=v"),
				},
				TargetGroupConfig: &alb.TargetGroupConfig{
					TargetGroups: []*alb.TargetGroupTuple{
						{TargetGroupId: albRuleTestString(albRuleTestTargetId), Weight: albRuleTestInt64(30)},
					},
					TargetGroupStickySession: &alb.TargetGroupStickySession{
						Enabled: albRuleTestBool(true),
						Timeout: albRuleTestInt64(1000),
					},
				},
			},
		},
		Conditions: []*alb.RuleCondition{
			{
				Type:              albRuleTestString("Path"),
				CookieConfig:      []*alb.HTTPCookieInfo{{Key: albRuleTestString("session"), Value: albRuleTestString("abc")}},
				HeaderConfig:      &alb.HTTPHeaderInfo{Key: albRuleTestString("X-Env"), Values: []*string{albRuleTestString("prod"), albRuleTestString("staging")}},
				HostConfig:        []*string{albRuleTestString("example.com")},
				MethodConfig:      []*string{albRuleTestString("GET")},
				PathConfig:        []*string{albRuleTestString("/api")},
				QueryStringConfig: []*alb.HTTPQueryStringInfo{{Key: albRuleTestString("version"), Value: albRuleTestString("v1")}},
				SourceIpConfig:    []*string{albRuleTestString("10.0.0.1/32")},
			},
		},
		CreateTime: albRuleTestString(albRuleTestCreateTime),
		Direction:  albRuleTestString("Request"),
		Priority:   albRuleTestInt64(10),
		RuleId:     albRuleTestString(albRuleTestRuleId),
		RuleName:   albRuleTestString("tf-example"),
		Status:     albRuleTestString("Active"),
		Tags:       []*alb.TagInfo{{TagKey: albRuleTestString("env"), TagValue: albRuleTestString("test")}},
	}
}

func albRuleTestDescribeResponse(rules ...*alb.RuleOutput) *alb.DescribeRulesResponse {
	return &alb.DescribeRulesResponse{
		Response: &alb.DescribeRulesResponseParams{
			Rules:      rules,
			NextToken:  albRuleTestString(""),
			TotalCount: albRuleTestInt64(int64(len(rules))),
		},
	}
}

func albRuleTestModel(t *testing.T) AlbRuleResourceModel {
	t.Helper()

	return AlbRuleResourceModel{
		Id:             types.StringNull(),
		ListenerId:     types.StringValue(albRuleTestListenerId),
		LoadBalancerId: types.StringValue(albRuleTestLbId),
		Rules:          albRuleTestList(t, albRuleInputAttrTypes(), albRuleTestFullRule(t)),
		ClientToken:    types.StringValue("tf-client-token"),
		DryRun:         types.BoolValue(false),
		NextToken:      types.StringNull(),
	}
}

func albRuleTestManagedModel(t *testing.T) AlbRuleResourceModel {
	t.Helper()

	model := albRuleTestModel(t)
	model.Id = types.StringValue(albRuleComposeId(albRuleTestLbId, albRuleTestListenerId, albRuleTestRuleId))
	return model
}

func albRuleTestNoErrorDiagnostics(t *testing.T, diags interface{ HasError() bool }) {
	t.Helper()

	if diags.HasError() {
		t.Fatalf("unexpected error diagnostics")
	}
}

func albRuleTestAssertString(t *testing.T, name string, value types.String, want string) {
	t.Helper()

	if value.IsNull() || value.IsUnknown() {
		t.Fatalf("%s is null or unknown, want %q", name, want)
	}
	if value.ValueString() != want {
		t.Fatalf("%s is %q, want %q", name, value.ValueString(), want)
	}
}

func albRuleTestAssertInt64(t *testing.T, name string, value types.Int64, want int64) {
	t.Helper()

	if value.IsNull() || value.IsUnknown() {
		t.Fatalf("%s is null or unknown, want %d", name, want)
	}
	if value.ValueInt64() != want {
		t.Fatalf("%s is %d, want %d", name, value.ValueInt64(), want)
	}
}

func albRuleTestAssertBool(t *testing.T, name string, value types.Bool, want bool) {
	t.Helper()

	if value.IsNull() || value.IsUnknown() {
		t.Fatalf("%s is null or unknown, want %t", name, want)
	}
	if value.ValueBool() != want {
		t.Fatalf("%s is %t, want %t", name, value.ValueBool(), want)
	}
}

// albRuleTestAssertRule verifies every leaf of the sample rule round-trips.
func albRuleTestAssertRule(t *testing.T, rules types.List) {
	t.Helper()

	if rules.IsNull() || rules.IsUnknown() {
		t.Fatalf("rules is null or unknown")
	}
	elements := rules.Elements()
	if len(elements) != 1 {
		t.Fatalf("rules has %d elements, want 1", len(elements))
	}
	ruleObject, ok := elements[0].(types.Object)
	if !ok {
		t.Fatalf("rules[0] is not an object")
	}
	rule := ruleObject.Attributes()
	albRuleTestAssertInt64(t, "rules[0].priority", rule["priority"].(types.Int64), 10)
	albRuleTestAssertString(t, "rules[0].direction", rule["direction"].(types.String), "Request")
	albRuleTestAssertString(t, "rules[0].rule_name", rule["rule_name"].(types.String), "tf-example")

	tags := rule["tags"].(types.List)
	if len(tags.Elements()) != 1 {
		t.Fatalf("tags has %d elements, want 1", len(tags.Elements()))
	}
	tag := tags.Elements()[0].(types.Object).Attributes()
	albRuleTestAssertString(t, "tag_key", tag["tag_key"].(types.String), "env")
	albRuleTestAssertString(t, "tag_value", tag["tag_value"].(types.String), "test")

	actions := rule["actions"].(types.List)
	if len(actions.Elements()) != 1 {
		t.Fatalf("actions has %d elements, want 1", len(actions.Elements()))
	}
	action := actions.Elements()[0].(types.Object).Attributes()
	albRuleTestAssertInt64(t, "actions[0].order", action["order"].(types.Int64), 1)
	albRuleTestAssertString(t, "actions[0].type", action["type"].(types.String), "TargetGroup")

	fixed := action["fixed_response_config"].(types.Object).Attributes()
	albRuleTestAssertInt64(t, "fixed_response_config.http_code", fixed["http_code"].(types.Int64), 200)
	albRuleTestAssertString(t, "fixed_response_config.content", fixed["content"].(types.String), "ok")
	albRuleTestAssertString(t, "fixed_response_config.content_type", fixed["content_type"].(types.String), "text/plain")

	insert := action["insert_header_config"].(types.Object).Attributes()
	albRuleTestAssertString(t, "insert_header_config.key", insert["key"].(types.String), "X-Test")
	albRuleTestAssertString(t, "insert_header_config.value", insert["value"].(types.String), "ClientIp")
	albRuleTestAssertString(t, "insert_header_config.value_type", insert["value_type"].(types.String), "SystemDefined")

	redirect := action["redirect_config"].(types.Object).Attributes()
	albRuleTestAssertInt64(t, "redirect_config.http_code", redirect["http_code"].(types.Int64), 301)
	albRuleTestAssertString(t, "redirect_config.host", redirect["host"].(types.String), "example.com")
	albRuleTestAssertString(t, "redirect_config.path", redirect["path"].(types.String), "/new")
	albRuleTestAssertString(t, "redirect_config.port", redirect["port"].(types.String), "443")
	albRuleTestAssertString(t, "redirect_config.protocol", redirect["protocol"].(types.String), "HTTPS")
	albRuleTestAssertString(t, "redirect_config.query", redirect["query"].(types.String), "a=b")

	remove := action["remove_header_config"].(types.Object).Attributes()
	albRuleTestAssertString(t, "remove_header_config.key", remove["key"].(types.String), "X-Remove")

	rewrite := action["rewrite_config"].(types.Object).Attributes()
	albRuleTestAssertString(t, "rewrite_config.host", rewrite["host"].(types.String), "example.com")
	albRuleTestAssertString(t, "rewrite_config.path", rewrite["path"].(types.String), "/rewrite")
	albRuleTestAssertString(t, "rewrite_config.query", rewrite["query"].(types.String), "k=v")

	targetGroupConfig := action["target_group_config"].(types.Object).Attributes()
	targetGroups := targetGroupConfig["target_groups"].(types.List)
	if len(targetGroups.Elements()) != 1 {
		t.Fatalf("target_groups has %d elements, want 1", len(targetGroups.Elements()))
	}
	targetGroup := targetGroups.Elements()[0].(types.Object).Attributes()
	albRuleTestAssertString(t, "target_group_id", targetGroup["target_group_id"].(types.String), albRuleTestTargetId)
	albRuleTestAssertInt64(t, "weight", targetGroup["weight"].(types.Int64), 30)

	sticky := targetGroupConfig["target_group_sticky_session"].(types.Object).Attributes()
	albRuleTestAssertBool(t, "sticky_session.enabled", sticky["enabled"].(types.Bool), true)
	albRuleTestAssertInt64(t, "sticky_session.timeout", sticky["timeout"].(types.Int64), 1000)

	conditions := rule["conditions"].(types.List)
	if len(conditions.Elements()) != 1 {
		t.Fatalf("conditions has %d elements, want 1", len(conditions.Elements()))
	}
	condition := conditions.Elements()[0].(types.Object).Attributes()
	albRuleTestAssertString(t, "conditions[0].type", condition["type"].(types.String), "Path")

	cookies := condition["cookie_config"].(types.List)
	if len(cookies.Elements()) != 1 {
		t.Fatalf("cookie_config has %d elements, want 1", len(cookies.Elements()))
	}
	cookie := cookies.Elements()[0].(types.Object).Attributes()
	albRuleTestAssertString(t, "cookie_config.key", cookie["key"].(types.String), "session")
	albRuleTestAssertString(t, "cookie_config.value", cookie["value"].(types.String), "abc")

	header := condition["header_config"].(types.Object).Attributes()
	albRuleTestAssertString(t, "header_config.key", header["key"].(types.String), "X-Env")
	if got := header["values"].(types.Set); len(got.Elements()) != 2 {
		t.Fatalf("header_config.values has %d elements, want 2", len(got.Elements()))
	}

	for _, name := range []string{"host_config", "method_config", "path_config", "source_ip_config"} {
		value := condition[name].(types.Set)
		if value.IsNull() || len(value.Elements()) != 1 {
			t.Fatalf("%s is not a one element set", name)
		}
	}

	queryStrings := condition["query_string_config"].(types.List)
	if len(queryStrings.Elements()) != 1 {
		t.Fatalf("query_string_config has %d elements, want 1", len(queryStrings.Elements()))
	}
	queryString := queryStrings.Elements()[0].(types.Object).Attributes()
	albRuleTestAssertString(t, "query_string_config.key", queryString["key"].(types.String), "version")
	albRuleTestAssertString(t, "query_string_config.value", queryString["value"].(types.String), "v1")
}

// albRuleTestAssertCreateRequest verifies every leaf of the expand path.
func albRuleTestAssertCreateRequest(t *testing.T, request *alb.CreateRulesRequest) {
	t.Helper()

	if request == nil {
		t.Fatalf("CreateRules was not called")
	}
	if request.LoadBalancerId == nil || *request.LoadBalancerId != albRuleTestLbId {
		t.Fatalf("CreateRules LoadBalancerId is not %q", albRuleTestLbId)
	}
	if request.ListenerId == nil || *request.ListenerId != albRuleTestListenerId {
		t.Fatalf("CreateRules ListenerId is not %q", albRuleTestListenerId)
	}
	if request.ClientToken == nil || *request.ClientToken != "tf-client-token" {
		t.Fatalf("CreateRules ClientToken is not tf-client-token")
	}
	if request.DryRun == nil || *request.DryRun {
		t.Fatalf("CreateRules DryRun is not false")
	}
	if len(request.Rules) != 1 {
		t.Fatalf("CreateRules sent %d rules, want 1", len(request.Rules))
	}
	rule := request.Rules[0]
	if rule.Priority == nil || *rule.Priority != 10 {
		t.Fatalf("CreateRules Priority is not 10")
	}
	if rule.Direction == nil || *rule.Direction != "Request" {
		t.Fatalf("CreateRules Direction is not Request")
	}
	if rule.RuleName == nil || *rule.RuleName != "tf-example" {
		t.Fatalf("CreateRules RuleName is not tf-example")
	}
	if len(rule.Tags) != 1 || rule.Tags[0].TagKey == nil || *rule.Tags[0].TagKey != "env" ||
		rule.Tags[0].TagValue == nil || *rule.Tags[0].TagValue != "test" {
		t.Fatalf("CreateRules Tags are not [env=test]")
	}
	if len(rule.Actions) != 1 {
		t.Fatalf("CreateRules Actions has %d elements, want 1", len(rule.Actions))
	}
	action := rule.Actions[0]
	if action.Order == nil || *action.Order != 1 || action.Type == nil || *action.Type != "TargetGroup" {
		t.Fatalf("CreateRules action order/type is wrong")
	}
	if action.FixedResponseConfig == nil || action.FixedResponseConfig.HttpCode == nil ||
		*action.FixedResponseConfig.HttpCode != 200 || action.FixedResponseConfig.Content == nil ||
		*action.FixedResponseConfig.Content != "ok" || action.FixedResponseConfig.ContentType == nil ||
		*action.FixedResponseConfig.ContentType != "text/plain" {
		t.Fatalf("CreateRules FixedResponseConfig is wrong")
	}
	if action.InsertHeaderConfig == nil || action.InsertHeaderConfig.Key == nil ||
		*action.InsertHeaderConfig.Key != "X-Test" || action.InsertHeaderConfig.Value == nil ||
		*action.InsertHeaderConfig.Value != "ClientIp" || action.InsertHeaderConfig.ValueType == nil ||
		*action.InsertHeaderConfig.ValueType != "SystemDefined" {
		t.Fatalf("CreateRules InsertHeaderConfig is wrong")
	}
	if action.RedirectConfig == nil || action.RedirectConfig.HttpCode == nil ||
		*action.RedirectConfig.HttpCode != 301 || action.RedirectConfig.Host == nil ||
		*action.RedirectConfig.Host != "example.com" || action.RedirectConfig.Path == nil ||
		*action.RedirectConfig.Path != "/new" || action.RedirectConfig.Port == nil ||
		*action.RedirectConfig.Port != "443" || action.RedirectConfig.Protocol == nil ||
		*action.RedirectConfig.Protocol != "HTTPS" || action.RedirectConfig.Query == nil ||
		*action.RedirectConfig.Query != "a=b" {
		t.Fatalf("CreateRules RedirectConfig is wrong")
	}
	if action.RemoveHeaderConfig == nil || action.RemoveHeaderConfig.Key == nil ||
		*action.RemoveHeaderConfig.Key != "X-Remove" {
		t.Fatalf("CreateRules RemoveHeaderConfig is wrong")
	}
	if action.RewriteConfig == nil || action.RewriteConfig.Host == nil ||
		*action.RewriteConfig.Host != "example.com" || action.RewriteConfig.Path == nil ||
		*action.RewriteConfig.Path != "/rewrite" || action.RewriteConfig.Query == nil ||
		*action.RewriteConfig.Query != "k=v" {
		t.Fatalf("CreateRules RewriteConfig is wrong")
	}
	if action.TargetGroupConfig == nil || len(action.TargetGroupConfig.TargetGroups) != 1 ||
		action.TargetGroupConfig.TargetGroups[0].TargetGroupId == nil ||
		*action.TargetGroupConfig.TargetGroups[0].TargetGroupId != albRuleTestTargetId ||
		action.TargetGroupConfig.TargetGroups[0].Weight == nil ||
		*action.TargetGroupConfig.TargetGroups[0].Weight != 30 ||
		action.TargetGroupConfig.TargetGroupStickySession == nil ||
		action.TargetGroupConfig.TargetGroupStickySession.Enabled == nil ||
		!*action.TargetGroupConfig.TargetGroupStickySession.Enabled ||
		action.TargetGroupConfig.TargetGroupStickySession.Timeout == nil ||
		*action.TargetGroupConfig.TargetGroupStickySession.Timeout != 1000 {
		t.Fatalf("CreateRules TargetGroupConfig is wrong")
	}
	if len(rule.Conditions) != 1 {
		t.Fatalf("CreateRules Conditions has %d elements, want 1", len(rule.Conditions))
	}
	condition := rule.Conditions[0]
	if condition.Type == nil || *condition.Type != "Path" {
		t.Fatalf("CreateRules condition Type is not Path")
	}
	if len(condition.CookieConfig) != 1 || condition.CookieConfig[0].Key == nil ||
		*condition.CookieConfig[0].Key != "session" || condition.CookieConfig[0].Value == nil ||
		*condition.CookieConfig[0].Value != "abc" {
		t.Fatalf("CreateRules CookieConfig is wrong")
	}
	if condition.HeaderConfig == nil || condition.HeaderConfig.Key == nil ||
		*condition.HeaderConfig.Key != "X-Env" || len(condition.HeaderConfig.Values) != 2 {
		t.Fatalf("CreateRules HeaderConfig is wrong")
	}
	if len(condition.HostConfig) != 1 || len(condition.MethodConfig) != 1 ||
		len(condition.PathConfig) != 1 || len(condition.SourceIpConfig) != 1 {
		t.Fatalf("CreateRules string set conditions are wrong")
	}
	if len(condition.QueryStringConfig) != 1 || condition.QueryStringConfig[0].Key == nil ||
		*condition.QueryStringConfig[0].Key != "version" || condition.QueryStringConfig[0].Value == nil ||
		*condition.QueryStringConfig[0].Value != "v1" {
		t.Fatalf("CreateRules QueryStringConfig is wrong")
	}
}

func TestAlbRuleSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)
	attributes := testContext.schema.Attributes

	for _, name := range []string{"id", "listener_id", "load_balancer_id", "rules", "client_token", "dry_run", "next_token"} {
		if _, ok := attributes[name]; !ok {
			t.Fatalf("schema is missing attribute %q", name)
		}
	}

	if !attributes["id"].IsComputed() {
		t.Fatalf("id must be computed")
	}
	for _, name := range []string{"listener_id", "load_balancer_id", "rules"} {
		if !attributes[name].IsRequired() {
			t.Fatalf("%s must be required", name)
		}
	}
	for _, name := range []string{"client_token", "dry_run"} {
		if !attributes[name].IsOptional() {
			t.Fatalf("%s must be optional", name)
		}
	}
	if !attributes["next_token"].IsComputed() {
		t.Fatalf("next_token must be computed")
	}

	rules, ok := attributes["rules"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("rules is not a list nested attribute")
	}
	ruleAttributes := rules.NestedObject.Attributes
	for _, name := range []string{"actions", "conditions", "priority", "direction", "rule_name", "tags"} {
		if _, ok := ruleAttributes[name]; !ok {
			t.Fatalf("rules is missing attribute %q", name)
		}
	}
	if !ruleAttributes["actions"].IsRequired() || !ruleAttributes["conditions"].IsRequired() ||
		!ruleAttributes["priority"].IsRequired() {
		t.Fatalf("actions, conditions and priority must be required")
	}

	actions := ruleAttributes["actions"].(schema.ListNestedAttribute).NestedObject.Attributes
	if len(actions) != 8 {
		t.Fatalf("actions has %d attributes, want 8", len(actions))
	}
	if !actions["order"].IsRequired() || !actions["type"].IsRequired() {
		t.Fatalf("action order and type must be required")
	}
	for _, name := range []string{"fixed_response_config", "insert_header_config", "redirect_config", "remove_header_config", "rewrite_config", "target_group_config"} {
		if !actions[name].IsOptional() {
			t.Fatalf("actions.%s must be optional", name)
		}
	}

	conditions := ruleAttributes["conditions"].(schema.ListNestedAttribute).NestedObject.Attributes
	if len(conditions) != 8 {
		t.Fatalf("conditions has %d attributes, want 8", len(conditions))
	}
	if !conditions["type"].IsRequired() {
		t.Fatalf("condition type must be required")
	}
}

func TestAlbRuleCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	var capturedRequest *alb.CreateRulesRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"CreateRulesWithContext",
		func(_ context.Context, request *alb.CreateRulesRequest) (*alb.CreateRulesResponse, error) {
			capturedRequest = request
			return &alb.CreateRulesResponse{
				Response: &alb.CreateRulesResponseParams{RuleIds: []*string{albRuleTestString(albRuleTestRuleId)}},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, request *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			if len(request.RuleIds) != 1 || request.RuleIds[0] == nil || *request.RuleIds[0] != albRuleTestRuleId {
				t.Fatalf("DescribeRules RuleIds is not [%q]", albRuleTestRuleId)
			}
			return albRuleTestDescribeResponse(albRuleTestCloudRule()), nil
		},
	)

	model := albRuleTestModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	albRuleTestNoErrorDiagnostics(t, createResp.Diagnostics)

	albRuleTestAssertCreateRequest(t, capturedRequest)

	state := testContext.readModel(t, createResp.State)
	wantId := albRuleComposeId(albRuleTestLbId, albRuleTestListenerId, albRuleTestRuleId)
	albRuleTestAssertString(t, "id", state.Id, wantId)
	albRuleTestAssertString(t, "listener_id", state.ListenerId, albRuleTestListenerId)
	albRuleTestAssertString(t, "load_balancer_id", state.LoadBalancerId, albRuleTestLbId)
	albRuleTestAssertString(t, "client_token", state.ClientToken, "tf-client-token")
	albRuleTestAssertBool(t, "dry_run", state.DryRun, false)
	albRuleTestAssertRule(t, state.Rules)
}

func TestAlbRuleCreateRejectsMultipleRules(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	model := albRuleTestModel(t)
	model.Rules = albRuleTestList(t, albRuleInputAttrTypes(), albRuleTestFullRule(t), albRuleTestFullRule(t))

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for multiple rules")
	}
}

func TestAlbRuleCreateRejectsDryRun(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	model := albRuleTestModel(t)
	model.DryRun = types.BoolValue(true)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when dry_run is true")
	}
}

func TestAlbRuleCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"CreateRulesWithContext",
		func(_ context.Context, _ *alb.CreateRulesRequest) (*alb.CreateRulesResponse, error) {
			return nil, errors.New("create failed")
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, albRuleTestModel(t))},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for a failing CreateRules call")
	}
}

func TestAlbRuleCreateEmptyResponses(t *testing.T) {
	cases := []struct {
		name     string
		response *alb.CreateRulesResponse
	}{
		{name: "nil response", response: nil},
		{name: "nil response params", response: &alb.CreateRulesResponse{}},
		{name: "no rule ids", response: &alb.CreateRulesResponse{Response: &alb.CreateRulesResponseParams{}}},
		{
			name: "two rule ids",
			response: &alb.CreateRulesResponse{Response: &alb.CreateRulesResponseParams{
				RuleIds: []*string{albRuleTestString(albRuleTestRuleId), albRuleTestString("rule-87654321")},
			}},
		},
		{
			name: "empty rule id",
			response: &alb.CreateRulesResponse{Response: &alb.CreateRulesResponseParams{
				RuleIds: []*string{albRuleTestString("")},
			}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			testContext := newAlbRuleTestContext(t, patches)
			response := tc.response
			patches.ApplyMethodFunc(
				testContext.client,
				"CreateRulesWithContext",
				func(_ context.Context, _ *alb.CreateRulesRequest) (*alb.CreateRulesResponse, error) {
					return response, nil
				},
			)

			createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
			testContext.resource.Create(
				context.Background(),
				resource.CreateRequest{Plan: testContext.plan(t, albRuleTestModel(t))},
				createResp,
			)
			if !createResp.Diagnostics.HasError() {
				t.Fatalf("expected an error diagnostic for %s", tc.name)
			}
		})
	}
}

func TestAlbRuleRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, request *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			if request.LoadBalancerId == nil || *request.LoadBalancerId != albRuleTestLbId {
				t.Fatalf("DescribeRules LoadBalancerId is not %q", albRuleTestLbId)
			}
			if request.ListenerId == nil || *request.ListenerId != albRuleTestListenerId {
				t.Fatalf("DescribeRules ListenerId is not %q", albRuleTestListenerId)
			}
			return albRuleTestDescribeResponse(
				&alb.RuleOutput{RuleId: albRuleTestString("rule-unrelated"), Status: albRuleTestString("Active")},
				albRuleTestCloudRule(),
			), nil
		},
	)

	model := albRuleTestManagedModel(t)
	readResp := testContext.runReadWithState(context.Background(), testContext.state(t, model))
	albRuleTestNoErrorDiagnostics(t, readResp.Diagnostics)

	state := testContext.readModel(t, readResp.State)
	albRuleTestAssertString(t, "id", state.Id, model.Id.ValueString())
	albRuleTestAssertString(t, "listener_id", state.ListenerId, albRuleTestListenerId)
	albRuleTestAssertString(t, "load_balancer_id", state.LoadBalancerId, albRuleTestLbId)
	albRuleTestAssertString(t, "next_token", state.NextToken, "")
	albRuleTestAssertRule(t, state.Rules)
}

func TestAlbRuleReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, _ *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			return albRuleTestDescribeResponse(), nil
		},
	)

	readResp := testContext.runReadWithState(context.Background(), testContext.state(t, albRuleTestManagedModel(t)))
	albRuleTestNoErrorDiagnostics(t, readResp.Diagnostics)
	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the rule does not exist")
	}
}

func TestAlbRuleReadNotFoundErrorCode(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, _ *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			return nil, sdkerr.NewTencentCloudSDKError("ResourceNotFound.Rule", "rule missing", "")
		},
	)

	readResp := testContext.runReadWithState(context.Background(), testContext.state(t, albRuleTestManagedModel(t)))
	albRuleTestNoErrorDiagnostics(t, readResp.Diagnostics)
	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the cloud reports ResourceNotFound")
	}
}

func TestAlbRuleReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, _ *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			return nil, errors.New("describe failed")
		},
	)

	readResp := testContext.runReadWithState(context.Background(), testContext.state(t, albRuleTestManagedModel(t)))
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for a failing DescribeRules call")
	}
	if readResp.State.Raw.IsNull() {
		t.Fatalf("an API failure must not remove state")
	}
}

func TestAlbRuleReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, _ *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			return &alb.DescribeRulesResponse{}, nil
		},
	)

	readResp := testContext.runReadWithState(context.Background(), testContext.state(t, albRuleTestManagedModel(t)))
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for a nil DescribeRules response")
	}
}

func TestAlbRuleReadInvalidId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	model := albRuleTestManagedModel(t)
	model.Id = types.StringValue("broken-id")

	readResp := testContext.runReadWithState(context.Background(), testContext.state(t, model))
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for an invalid resource id")
	}
}

func TestAlbRuleDescribePagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	calls := 0
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, request *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			calls++
			if calls == 1 {
				return &alb.DescribeRulesResponse{Response: &alb.DescribeRulesResponseParams{
					NextToken: albRuleTestString("page-2"),
				}}, nil
			}
			if calls > 1 && (request.NextToken == nil || *request.NextToken != "page-2") {
				t.Fatalf("second page request must carry NextToken page-2")
			}
			return albRuleTestDescribeResponse(albRuleTestCloudRule()), nil
		},
	)

	item, token, err := albRuleDescribe(
		context.Background(), testContext.client,
		albRuleTestLbId, albRuleTestListenerId, albRuleTestRuleId,
	)
	if err != nil {
		t.Fatalf("albRuleDescribe failed: %v", err)
	}
	if calls != 2 {
		t.Fatalf("DescribeRules was called %d times, want 2", calls)
	}
	if item == nil || item.RuleId == nil || *item.RuleId != albRuleTestRuleId {
		t.Fatalf("albRuleDescribe returned the wrong rule")
	}
	if token == nil || *token != "" {
		t.Fatalf("last page token is not empty")
	}
}

func TestAlbRuleCreateWaitsForTransitionalStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"CreateRulesWithContext",
		func(_ context.Context, _ *alb.CreateRulesRequest) (*alb.CreateRulesResponse, error) {
			return &alb.CreateRulesResponse{
				Response: &alb.CreateRulesResponseParams{RuleIds: []*string{albRuleTestString(albRuleTestRuleId)}},
			}, nil
		},
	)

	// Drive the wait helper directly with a converged response; exercising the
	// real RetryFramework timer is not appropriate for a unit test.
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, _ *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			item := albRuleTestCloudRule()
			item.Status = albRuleTestString("Active")
			return albRuleTestDescribeResponse(item), nil
		},
	)

	item, token, err := albRuleWaitAvailable(
		context.Background(), testContext.client,
		albRuleTestLbId, albRuleTestListenerId, albRuleTestRuleId,
	)
	if err != nil {
		t.Fatalf("albRuleWaitAvailable failed: %v", err)
	}
	if item == nil || item.RuleId == nil || *item.RuleId != albRuleTestRuleId {
		t.Fatalf("albRuleWaitAvailable returned the wrong rule")
	}
	if token == nil || *token != "" {
		t.Fatalf("albRuleWaitAvailable returned a non-empty last page token")
	}
}

func TestAlbRuleExpandMinimalRuleOmitsNullFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	action := albRuleTestObject(t, albRuleActionAttrTypes(), map[string]attr.Value{
		"order":                 types.Int64Value(1),
		"type":                  types.StringValue("FixedResponse"),
		"fixed_response_config": types.ObjectNull(albRuleFixedResponseAttrTypes()),
		"insert_header_config":  types.ObjectNull(albRuleInsertHeaderAttrTypes()),
		"redirect_config":       types.ObjectNull(albRuleRedirectAttrTypes()),
		"remove_header_config":  types.ObjectNull(albRuleRemoveHeaderAttrTypes()),
		"rewrite_config":        types.ObjectNull(albRuleRewriteAttrTypes()),
		"target_group_config":   types.ObjectNull(albRuleTargetGroupConfigAttrTypes()),
	})
	condition := albRuleTestObject(t, albRuleConditionAttrTypes(), map[string]attr.Value{
		"type":                types.StringValue("Host"),
		"cookie_config":       types.ListNull(types.ObjectType{AttrTypes: albRuleCookieAttrTypes()}),
		"header_config":       types.ObjectNull(albRuleHeaderAttrTypes()),
		"host_config":         types.SetNull(types.StringType),
		"method_config":       types.SetNull(types.StringType),
		"path_config":         types.SetNull(types.StringType),
		"query_string_config": types.ListNull(types.ObjectType{AttrTypes: albRuleQueryStringAttrTypes()}),
		"source_ip_config":    types.SetNull(types.StringType),
	})
	rule := albRuleTestObject(t, albRuleInputAttrTypes(), map[string]attr.Value{
		"actions":    albRuleTestList(t, albRuleActionAttrTypes(), action),
		"conditions": albRuleTestList(t, albRuleConditionAttrTypes(), condition),
		"priority":   types.Int64Value(1),
		"direction":  types.StringNull(),
		"rule_name":  types.StringNull(),
		"tags":       types.ListNull(types.ObjectType{AttrTypes: albRuleTagAttrTypes()}),
	})
	model := albRuleTestModel(t)
	model.Rules = albRuleTestList(t, albRuleInputAttrTypes(), rule)
	model.ClientToken = types.StringNull()
	model.DryRun = types.BoolNull()

	var capturedRequest *alb.CreateRulesRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"CreateRulesWithContext",
		func(_ context.Context, request *alb.CreateRulesRequest) (*alb.CreateRulesResponse, error) {
			capturedRequest = request
			return &alb.CreateRulesResponse{
				Response: &alb.CreateRulesResponseParams{RuleIds: []*string{albRuleTestString(albRuleTestRuleId)}},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, _ *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			item := &alb.RuleOutput{
				RuleId:     albRuleTestString(albRuleTestRuleId),
				Status:     albRuleTestString("Active"),
				Priority:   albRuleTestInt64(1),
				Actions:    []*alb.RuleAction{{Order: albRuleTestInt64(1), Type: albRuleTestString("FixedResponse")}},
				Conditions: []*alb.RuleCondition{{Type: albRuleTestString("Host")}},
			}
			return albRuleTestDescribeResponse(item), nil
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	albRuleTestNoErrorDiagnostics(t, createResp.Diagnostics)

	if capturedRequest == nil || capturedRequest.ClientToken != nil || capturedRequest.DryRun != nil {
		t.Fatalf("null create-only fields must be omitted from the request")
	}
	if len(capturedRequest.Rules) != 1 {
		t.Fatalf("CreateRules sent %d rules, want 1", len(capturedRequest.Rules))
	}
	created := capturedRequest.Rules[0]
	if created.Direction != nil || created.RuleName != nil || created.Tags != nil {
		t.Fatalf("null optional rule fields must be omitted from the request")
	}
	createdAction := created.Actions[0]
	if createdAction.FixedResponseConfig != nil || createdAction.TargetGroupConfig != nil ||
		createdAction.RedirectConfig != nil || createdAction.RewriteConfig != nil ||
		createdAction.InsertHeaderConfig != nil || createdAction.RemoveHeaderConfig != nil {
		t.Fatalf("null action configs must be omitted from the request")
	}
	createdCondition := created.Conditions[0]
	if createdCondition.HeaderConfig != nil || createdCondition.HostConfig != nil ||
		createdCondition.PathConfig != nil || createdCondition.SourceIpConfig != nil ||
		createdCondition.MethodConfig != nil || createdCondition.CookieConfig != nil ||
		createdCondition.QueryStringConfig != nil {
		t.Fatalf("null condition configs must be omitted from the request")
	}

	state := testContext.readModel(t, createResp.State)
	if !state.ClientToken.IsNull() || !state.DryRun.IsNull() {
		t.Fatalf("null create-only fields must stay null in state")
	}
	if !state.NextToken.IsNull() {
		t.Fatalf("next_token must stay null on the last page")
	}
	rules := state.Rules.Elements()
	if len(rules) != 1 {
		t.Fatalf("rules has %d elements, want 1", len(rules))
	}
	ruleState := rules[0].(types.Object).Attributes()
	if !ruleState["direction"].IsNull() || !ruleState["rule_name"].IsNull() || !ruleState["tags"].IsNull() {
		t.Fatalf("null optional fields must stay null after the read-back")
	}
	actionState := ruleState["actions"].(types.List).Elements()[0].(types.Object).Attributes()
	if !actionState["fixed_response_config"].IsNull() || !actionState["target_group_config"].IsNull() {
		t.Fatalf("null action configs must stay null after the read-back")
	}
}

func TestAlbRuleUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	var capturedRequest *alb.ModifyRulesAttributesRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"ModifyRulesAttributesWithContext",
		func(_ context.Context, request *alb.ModifyRulesAttributesRequest) (*alb.ModifyRulesAttributesResponse, error) {
			capturedRequest = request
			return &alb.ModifyRulesAttributesResponse{Response: &alb.ModifyRulesAttributesResponseParams{}}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, _ *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			return albRuleTestDescribeResponse(albRuleTestCloudRule()), nil
		},
	)

	state := albRuleTestManagedModel(t)
	plan := albRuleTestManagedModel(t)
	plan.NextToken = types.StringUnknown()

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, plan), State: testContext.state(t, state)},
		updateResp,
	)
	albRuleTestNoErrorDiagnostics(t, updateResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("ModifyRulesAttributes was not called")
	}
	if capturedRequest.LoadBalancerId == nil || *capturedRequest.LoadBalancerId != albRuleTestLbId {
		t.Fatalf("ModifyRulesAttributes LoadBalancerId is not %q", albRuleTestLbId)
	}
	if capturedRequest.ListenerId == nil || *capturedRequest.ListenerId != albRuleTestListenerId {
		t.Fatalf("ModifyRulesAttributes ListenerId is not %q", albRuleTestListenerId)
	}
	if len(capturedRequest.Rules) != 1 {
		t.Fatalf("ModifyRulesAttributes sent %d rules, want 1", len(capturedRequest.Rules))
	}
	rule := capturedRequest.Rules[0]
	if rule.RuleId == nil || *rule.RuleId != albRuleTestRuleId {
		t.Fatalf("ModifyRulesAttributes RuleId is not %q", albRuleTestRuleId)
	}
	if rule.RuleName == nil || *rule.RuleName != "tf-example" {
		t.Fatalf("ModifyRulesAttributes RuleName is not tf-example")
	}
	if rule.Priority == nil || *rule.Priority != 10 {
		t.Fatalf("ModifyRulesAttributes Priority is not 10")
	}
	if len(rule.Actions) != 1 || rule.Actions[0].TargetGroupConfig == nil {
		t.Fatalf("ModifyRulesAttributes Actions is wrong")
	}
	if len(rule.Conditions) != 1 || rule.Conditions[0].HeaderConfig == nil {
		t.Fatalf("ModifyRulesAttributes Conditions is wrong")
	}

	result := testContext.readModel(t, updateResp.State)
	albRuleTestAssertString(t, "id", result.Id, state.Id.ValueString())
	albRuleTestAssertRule(t, result.Rules)
}

func TestAlbRuleUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"ModifyRulesAttributesWithContext",
		func(_ context.Context, _ *alb.ModifyRulesAttributesRequest) (*alb.ModifyRulesAttributesResponse, error) {
			return nil, errors.New("modify failed")
		},
	)

	state := albRuleTestManagedModel(t)
	plan := albRuleTestManagedModel(t)
	plan.NextToken = types.StringUnknown()

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, plan), State: testContext.state(t, state)},
		updateResp,
	)
	if !updateResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for a failing ModifyRulesAttributes call")
	}
}

func TestAlbRuleUpdateRejectsClearedField(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	state := albRuleTestManagedModel(t)
	plan := albRuleTestManagedModel(t)
	plan.NextToken = types.StringUnknown()
	plan.ClientToken = types.StringNull()
	plan.DryRun = types.BoolNull()

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{Plan: testContext.plan(t, plan), State: testContext.state(t, state)},
		updateResp,
	)
	if !updateResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when a configured attribute is removed")
	}
	joined := ""
	for _, diagnostic := range updateResp.Diagnostics.Errors() {
		joined += diagnostic.Detail() + " "
	}
	if !strings.Contains(joined, "client_token") || !strings.Contains(joined, "dry_run") {
		t.Fatalf("error must name the cleared attributes, got %q", joined)
	}
}

func TestAlbRuleDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	var capturedRequest *alb.DeleteRulesRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteRulesWithContext",
		func(_ context.Context, request *alb.DeleteRulesRequest) (*alb.DeleteRulesResponse, error) {
			capturedRequest = request
			return &alb.DeleteRulesResponse{Response: &alb.DeleteRulesResponseParams{}}, nil
		},
	)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, albRuleTestManagedModel(t))},
		deleteResp,
	)
	albRuleTestNoErrorDiagnostics(t, deleteResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("DeleteRules was not called")
	}
	if capturedRequest.LoadBalancerId == nil || *capturedRequest.LoadBalancerId != albRuleTestLbId {
		t.Fatalf("DeleteRules LoadBalancerId is not %q", albRuleTestLbId)
	}
	if capturedRequest.ListenerId == nil || *capturedRequest.ListenerId != albRuleTestListenerId {
		t.Fatalf("DeleteRules ListenerId is not %q", albRuleTestListenerId)
	}
	if len(capturedRequest.RuleIds) != 1 || capturedRequest.RuleIds[0] == nil ||
		*capturedRequest.RuleIds[0] != albRuleTestRuleId {
		t.Fatalf("DeleteRules RuleIds is not [%q]", albRuleTestRuleId)
	}
}

func TestAlbRuleDeleteAlreadyGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteRulesWithContext",
		func(_ context.Context, _ *alb.DeleteRulesRequest) (*alb.DeleteRulesResponse, error) {
			return nil, sdkerr.NewTencentCloudSDKError("ResourceNotFound.Rule", "already gone", "")
		},
	)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, albRuleTestManagedModel(t))},
		deleteResp,
	)
	albRuleTestNoErrorDiagnostics(t, deleteResp.Diagnostics)
}

func TestAlbRuleDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteRulesWithContext",
		func(_ context.Context, _ *alb.DeleteRulesRequest) (*alb.DeleteRulesResponse, error) {
			return nil, errors.New("delete failed")
		},
	)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, albRuleTestManagedModel(t))},
		deleteResp,
	)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for a failing DeleteRules call")
	}
}

func TestAlbRuleImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	wantId := albRuleComposeId(albRuleTestLbId, albRuleTestListenerId, albRuleTestRuleId)
	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema, Raw: tftypes.NewValue(testContext.schema.Type().TerraformType(context.Background()), nil)}}
	importer.ImportState(context.Background(), resource.ImportStateRequest{ID: wantId}, importResp)
	albRuleTestNoErrorDiagnostics(t, importResp.Diagnostics)

	state := testContext.readModel(t, importResp.State)
	albRuleTestAssertString(t, "id", state.Id, wantId)
	albRuleTestAssertString(t, "load_balancer_id", state.LoadBalancerId, albRuleTestLbId)
	albRuleTestAssertString(t, "listener_id", state.ListenerId, albRuleTestListenerId)
}

// TestAlbRuleImportStateVerify imports a composite id, performs the Read that
// follows an import and asserts that the resulting state converges with what
// the cloud reports. This is the mock-only equivalent of ImportStateVerify.
func TestAlbRuleImportStateVerify(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeRulesWithContext",
		func(_ context.Context, request *alb.DescribeRulesRequest) (*alb.DescribeRulesResponse, error) {
			if request.LoadBalancerId == nil || *request.LoadBalancerId != albRuleTestLbId {
				t.Fatalf("imported Read used the wrong load_balancer_id")
			}
			if request.ListenerId == nil || *request.ListenerId != albRuleTestListenerId {
				t.Fatalf("imported Read used the wrong listener_id")
			}
			return albRuleTestDescribeResponse(albRuleTestCloudRule()), nil
		},
	)

	wantId := albRuleComposeId(albRuleTestLbId, albRuleTestListenerId, albRuleTestRuleId)
	importer := testContext.resource.(resource.ResourceWithImportState)
	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema, Raw: tftypes.NewValue(testContext.schema.Type().TerraformType(context.Background()), nil)}}
	importer.ImportState(context.Background(), resource.ImportStateRequest{ID: wantId}, importResp)
	albRuleTestNoErrorDiagnostics(t, importResp.Diagnostics)

	readResp := testContext.runReadWithState(context.Background(), importResp.State)
	albRuleTestNoErrorDiagnostics(t, readResp.Diagnostics)

	state := testContext.readModel(t, readResp.State)
	albRuleTestAssertString(t, "id", state.Id, wantId)
	albRuleTestAssertString(t, "load_balancer_id", state.LoadBalancerId, albRuleTestLbId)
	albRuleTestAssertString(t, "listener_id", state.ListenerId, albRuleTestListenerId)
	albRuleTestAssertRule(t, state.Rules)

	if !state.Rules.IsNull() {
		// A second Read from the refreshed state must not report any further
		// cloud-owned change, otherwise the import is not verifiable.
		secondResp := testContext.runReadWithState(context.Background(), readResp.State)
		albRuleTestNoErrorDiagnostics(t, secondResp.Diagnostics)
		second := testContext.readModel(t, secondResp.State)
		if !second.Rules.Equal(state.Rules) {
			t.Fatalf("a second Read after import changed rules; import does not converge")
		}
	}
}

func TestAlbRuleImportStateInvalidId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbRuleTestContext(t, patches)

	for _, id := range []string{"", "one", "one#two", "one#two#three#four", "#lst-1#rule-1"} {
		t.Run("id="+id, func(t *testing.T) {
			importer := testContext.resource.(resource.ResourceWithImportState)
			importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema, Raw: tftypes.NewValue(testContext.schema.Type().TerraformType(context.Background()), nil)}}
			importer.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, importResp)
			if !importResp.Diagnostics.HasError() {
				t.Fatalf("expected an error diagnostic for import id %q", id)
			}
		})
	}
}

func TestAlbRuleHelpers(t *testing.T) {
	loadBalancerId, listenerId, ruleId, ok := albRuleParseId(
		albRuleComposeId(albRuleTestLbId, albRuleTestListenerId, albRuleTestRuleId),
	)
	if !ok {
		t.Fatalf("failed to parse a composed id")
	}
	if loadBalancerId != albRuleTestLbId || listenerId != albRuleTestListenerId || ruleId != albRuleTestRuleId {
		t.Fatalf("parsed id parts do not match")
	}
	if _, _, _, ok := albRuleParseId("broken"); ok {
		t.Fatalf("an invalid id must not parse")
	}

	model := AlbRuleResourceModel{
		ClientToken: types.StringValue("token"),
		DryRun:      types.BoolValue(true),
		Rules:       albRuleTestList(t, albRuleInputAttrTypes(), albRuleTestFullRule(t)),
	}
	cleared := albRuleUpdateClearedFields(model, model)
	if len(cleared) != 0 {
		t.Fatalf("an unchanged model reported cleared paths: %v", cleared)
	}

	emptied := model
	emptied.ClientToken = types.StringNull()
	emptied.DryRun = types.BoolNull()
	cleared = albRuleUpdateClearedFields(model, emptied)
	if len(cleared) != 2 {
		t.Fatalf("expected 2 cleared paths, got %v", cleared)
	}
}
