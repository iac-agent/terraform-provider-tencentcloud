package oceanus_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/oceanus"
)

func TestOceanusOceanusClusterResourceTypeName(t *testing.T) {
	r := oceanus.NewOceanusOceanusClusterResource()
	req := resource.MetadataRequest{ProviderTypeName: "tencentcloud"}
	var resp resource.MetadataResponse
	r.Metadata(t.Context(), req, &resp)
	if resp.TypeName != "tencentcloud_oceanus_oceanus_cluster" {
		t.Fatalf("unexpected type name: %s", resp.TypeName)
	}
}

func TestOceanusOceanusClusterResourceSchemaIsValid(t *testing.T) {
	r := oceanus.NewOceanusOceanusClusterResource()
	var resp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}
	for _, name := range []string{"id", "cluster_name", "region_id", "zone_id", "login_password", "vpc_descriptions", "default_cos_bucket", "cu", "remark", "period", "instance_charge_type", "cluster_type", "renew_flag", "flink_web_ui_network_access_type", "slave_vpc_descriptions", "cu_memory", "status", "create_time", "update_time"} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("missing attribute %s in schema", name)
		}
	}
}
