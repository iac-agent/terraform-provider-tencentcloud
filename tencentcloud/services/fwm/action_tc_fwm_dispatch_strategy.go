package fwm

import (
	"context"
	"fmt"
	"time"

	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ action.Action is intentionally omitted: the factory function
// NewFwmDispatchStrategy returns action.Action, which already performs the
// implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) and the offline checks at
// compile time.
var (
	_ action.ActionWithConfigure      = &FwmDispatchStrategy{}
	_ action.ActionWithValidateConfig = &FwmDispatchStrategy{}
)

// fwmDispatchStrategyActionPollInterval is the interval between two
// DescribeStrategyDispatchStatus polls. It is a package-level variable so unit
// tests can shorten it.
var fwmDispatchStrategyActionPollInterval = 5 * time.Second

// SetFwmDispatchStrategyPollIntervalForTest overrides the polling interval and
// returns the previous value so callers can restore it. It only exists for unit
// tests that must not sleep for the production interval.
func SetFwmDispatchStrategyPollIntervalForTest(interval time.Duration) time.Duration {
	previous := fwmDispatchStrategyActionPollInterval
	fwmDispatchStrategyActionPollInterval = interval
	return previous
}

// defaultFwmDispatchStrategyInvokeTimeout is the built-in maximum duration to wait
// for the asynchronous strategy dispatch to finish. Practitioners can override it
// through the `timeouts` block, e.g. `timeouts { invoke = "10m" }`.
const defaultFwmDispatchStrategyInvokeTimeout = 15 * time.Minute

// NewFwmDispatchStrategy is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewFwmDispatchStrategy() action.Action {
	return &FwmDispatchStrategy{}
}

// FwmDispatchStrategy implements action.Action for
// tencentcloud_fwm_dispatch_strategy.
type FwmDispatchStrategy struct {
	fw.ActionWithConfigure
}

// FwmDispatchStrategyModel maps the schema attributes.
//
// An action schema has no Computed / Sensitive / PlanModifiers marks: only
// Required / Optional / WriteOnly / Validators are available. The dispatch result
// (status, progress, rule group count ...) therefore cannot be exposed as an
// attribute referenceable from HCL; it is reported through progress events.
type FwmDispatchStrategyModel struct {
	Product types.String `tfsdk:"product"`
	Status  types.Int64  `tfsdk:"status"`

	// Timeouts maps the optional `timeouts` block, which only has the `invoke`
	// stage because an action has a single lifecycle stage.
	Timeouts actiontimeouts.Value `tfsdk:"timeouts"`
}

func (a *FwmDispatchStrategy) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_dispatch_strategy"
}

func (a *FwmDispatchStrategy) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to dispatch FWM (Firewall Manager) rule configurations to the cloud via the " +
			"DispatchStrategy API, and waits until the dispatch reaches a terminal status. " +
			"'Dispatch' pushes the strategy definitions that are already configured (for example through " +
			"`tencentcloud_fwm_strategy` and `tencentcloud_fwm_rule_group`) so that they take effect. " +
			"This is a one-time invocation: it creates no Terraform-managed resource, exposes no attributes " +
			"referenceable from other resources, and destroying or removing it does not revert anything that was " +
			"already dispatched. Every successful invocation triggers a new cloud-side dispatch, so repeating the " +
			"invocation repeats the side effect. Triggering this action through `lifecycle.action_trigger` runs it " +
			"automatically during apply, including on drift-driven `after_update` events. A failed invocation does not " +
			"roll back the part that was already dispatched. `on_failure` accepts `halt` (default), `continue` and " +
			"`taint`; it is not a destroy guard, because a failing action during a destroy walk is downgraded to a " +
			"warning.",
		Attributes: map[string]schema.Attribute{
			"product": schema.StringAttribute{
				Optional: true,
				Description: "Product whose rule configurations are dispatched, for example `cfw_edge_acl`. " +
					"Omit it to dispatch every product. Only strategies of the given scope are affected; " +
					"strategies of other products are left untouched.",
				Validators: []validator.String{
					nonEmptyFwmDispatchString{},
				},
			},
			"status": schema.Int64Attribute{
				Optional: true,
				Description: "Dispatch operation to send: `1` dispatches the configured strategies, `2` aborts an " +
					"ongoing dispatch. Default is `1`. Aborting is not a rollback: strategies that were already " +
					"dispatched stay in effect. When `2` is used the action does not wait for a successful dispatch.",
				Validators: []validator.Int64{
					fwmDispatchStatusValidator{},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": fwmDispatchStrategyTimeoutsBlock(ctx),
		},
	}
}

// fwmDispatchStrategyTimeoutsBlock builds the `timeouts` block from the official
// terraform-plugin-framework-timeouts module, whose single optional `invoke`
// attribute is validated as a time.Duration at plan time. That module leaves the
// block itself undescribed, so the description is restored here to keep the schema
// self-explanatory.
func fwmDispatchStrategyTimeoutsBlock(ctx context.Context) schema.Block {
	block := actiontimeouts.BlockWithOpts(ctx, actiontimeouts.Opts{
		InvokeDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. " +
			"`10m`. It limits how long the action waits for the dispatch to finish. Default is `15m`.",
	})

	if nested, ok := block.(schema.SingleNestedBlock); ok {
		nested.Description = "The timeouts block allows you to specify the timeout for the invoke operation."
		return nested
	}

	return block
}

// ValidateConfig runs during `terraform validate` as well, where the provider is
// not configured yet, so it never touches a.Client().
func (a *FwmDispatchStrategy) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data FwmDispatchStrategyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values may be unknown during offline validation.
	if data.Product.IsNull() || data.Product.IsUnknown() {
		return
	}

	if data.Product.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("product"),
			"Invalid product",
			"product must not be an empty string; omit it to dispatch every product.",
		)
	}
}

// Invoke is called to run the logic of the action.
func (a *FwmDispatchStrategy) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data FwmDispatchStrategyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	progress := func(msg string) {
		if resp.SendProgress != nil {
			resp.SendProgress(action.InvokeProgressEvent{Message: msg})
		}
	}

	// Re-check the inputs at invoke time: config validation does not guarantee the
	// values are still valid at apply time.
	product := ""
	if !data.Product.IsNull() && !data.Product.IsUnknown() {
		product = data.Product.ValueString()
		if product == "" {
			resp.Diagnostics.AddError(
				"Invalid product",
				"product must not be an empty string; omit it to dispatch every product.",
			)
			return
		}
	}

	status := int64(1)
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status = data.Status.ValueInt64()
	}
	if status != 1 && status != FwmStrategyDispatchAbortStatus {
		resp.Diagnostics.AddError(
			"Invalid status",
			fmt.Sprintf("status must be 1 (dispatch) or 2 (abort), got %d.", status),
		)
		return
	}
	if status == FwmStrategyDispatchAbortStatus && product == "" {
		resp.Diagnostics.AddError(
			"Missing product",
			"aborting a dispatch (status = 2) needs the product scope; omit it only when dispatching.",
		)
		return
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	scope := product
	if scope == "" {
		scope = "all products"
	}

	service := NewFwmService(a.Client())
	progress(fmt.Sprintf("dispatching FWM strategy configuration for %s (status = %d)", scope, status))

	if err := service.DispatchFwmStrategy(ctx, product, status); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("dispatching FWM strategy configuration for %s", scope),
			err.Error(),
		)
		return
	}
	progress("DispatchStrategy request accepted; waiting for the dispatch to finish")

	if status == FwmStrategyDispatchAbortStatus {
		// An abort request stops the ongoing dispatch; it has no successful dispatch
		// status to poll and is therefore not awaited here.
		progress("abort request submitted")
		return
	}

	// An unset/null timeouts block falls back to the built-in default; a malformed
	// duration is reported as an error diagnostic (the schema validator normally
	// rejects it at plan time, this is the runtime safety net).
	timeout, timeoutDiags := data.Timeouts.Invoke(ctx, defaultFwmDispatchStrategyInvokeTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := service.WaitFwmStrategyDispatch(ctx, product, timeout, fwmDispatchStrategyActionPollInterval); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("waiting for FWM strategy dispatch of %s to finish", scope),
			err.Error(),
		)
		return
	}

	// Report the final snapshot through progress; an action has no state and no
	// attribute referenceable from HCL.
	if last, err := service.DescribeFwmStrategyDispatchStatus(ctx, product); err == nil {
		progress(fmt.Sprintf("dispatch finished: %s", fwmStrategyDispatchSnapshot(last)))
	} else {
		progress("dispatch finished; the final status could not be re-read")
	}
}

// nonEmptyFwmDispatchString rejects empty strings so that "omit the attribute"
// stays the only way to express "no value".
var _ validator.String = nonEmptyFwmDispatchString{}

type nonEmptyFwmDispatchString struct{}

func (nonEmptyFwmDispatchString) Description(_ context.Context) string {
	return "value must not be an empty string; omit the attribute instead"
}

func (v nonEmptyFwmDispatchString) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (nonEmptyFwmDispatchString) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if req.ConfigValue.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid empty value",
			"an empty string is not a valid value; omit the attribute instead",
		)
	}
}

// fwmDispatchStatusValidator accepts only the two documented DispatchStrategy
// status values.
var _ validator.Int64 = fwmDispatchStatusValidator{}

type fwmDispatchStatusValidator struct{}

func (fwmDispatchStatusValidator) Description(_ context.Context) string {
	return "value must be 1 (dispatch) or 2 (abort)"
}

func (v fwmDispatchStatusValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (fwmDispatchStatusValidator) ValidateInt64(_ context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	switch req.ConfigValue.ValueInt64() {
	case 1, FwmStrategyDispatchAbortStatus:
		return
	default:
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid status",
			fmt.Sprintf("status must be 1 (dispatch) or 2 (abort), got %d", req.ConfigValue.ValueInt64()),
		)
	}
}

// fwmDispatchDurationValidator rejects duration strings that are not parseable as a
// positive Go duration.
var _ validator.String = fwmDispatchDurationValidator{}

type fwmDispatchDurationValidator struct{}

func (fwmDispatchDurationValidator) Description(_ context.Context) string {
	return "value must be a positive Go duration string such as 10m"
}

func (v fwmDispatchDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (fwmDispatchDurationValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	parsed, err := time.ParseDuration(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid duration",
			fmt.Sprintf("%q is not a valid Go duration: %v", req.ConfigValue.ValueString(), err),
		)
		return
	}
	if parsed <= 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid duration",
			"the duration must be greater than zero",
		)
	}
}
