package alb

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

// TestAlbListenerResource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces advertised in the _ asserts.
func TestAlbListenerResource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &AlbListenerResource{}
	var _ resource.ResourceWithImportState = &AlbListenerResource{}
}

// TestAlbListenerResource_Factory validates the registry factory.
func TestAlbListenerResource_Factory(t *testing.T) {
	got := NewAlbListenerResource()
	if got == nil {
		t.Fatal("NewAlbListenerResource returned nil")
	}
	if _, ok := got.(*AlbListenerResource); !ok {
		t.Fatalf("factory returned %T, want *AlbListenerResource", got)
	}
}

// TestAlbListenerResource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestAlbListenerResource_ConfigureNilProviderData(t *testing.T) {
	r := &AlbListenerResource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// --- ID split ---

// TestSplitAlbListenerId validates parsing of the composite id.
func TestSplitAlbListenerId(t *testing.T) {
	good := "alb-xxxxxxxx" + tccommon.FILED_SP + "lst-yyyyyyyy"
	lb, lst, ok := splitAlbListenerId(good)
	if !ok || lb != "alb-xxxxxxxx" || lst != "lst-yyyyyyyy" {
		t.Errorf("splitAlbListenerId(%q) = (%q, %q, %v); want (\"alb-xxxxxxxx\", \"lst-yyyyyyyy\", true)", good, lb, lst, ok)
	}

	bad := []string{
		"",
		"only-listener",
		tccommon.FILED_SP + "lst",
		"alb" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	}
	for _, in := range bad {
		if _, _, ok := splitAlbListenerId(in); ok {
			t.Errorf("splitAlbListenerId(%q) should not be ok", in)
		}
	}
}

// TestAlbListenerId builds the composite id.
func TestAlbListenerId(t *testing.T) {
	got := albListenerId("alb-xxxxxxxx", "lst-yyyyyyyy")
	want := "alb-xxxxxxxx" + tccommon.FILED_SP + "lst-yyyyyyyy"
	if got != want {
		t.Errorf("albListenerId() = %q, want %q", got, want)
	}
}
