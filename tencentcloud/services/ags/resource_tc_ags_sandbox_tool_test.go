package ags

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	tcsdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const (
	agsTestToolId   = "sdt-test-0001"
	agsTestToolName = "tf-test-tool"
)

func agsString(s string) *string { return &s }
func agsBool(b bool) *bool       { return &b }
func agsInt64(i int64) *int64    { return &i }
func agsUint64(u uint64) *uint64 { return &u }

// agsSandboxToolSchema builds the framework resource schema once per test.
func agsSandboxToolSchema(t *testing.T) schema.Schema {
	t.Helper()
	r := NewAgsSandboxToolResource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// agsSandboxToolObjectType returns the tftypes object type of the schema.
func agsSandboxToolObjectType(t *testing.T) tftypes.Object {
	t.Helper()
	return agsSandboxToolSchema(t).Type().TerraformType(context.Background()).(tftypes.Object)
}

// agsNestedObjectValue builds a nested object value, using null for every
// attribute that the caller does not provide.
func agsNestedObjectValue(objType tftypes.Object, values map[string]tftypes.Value) tftypes.Value {
	filled := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		if v, ok := values[name]; ok {
			filled[name] = v
			continue
		}
		filled[name] = agsNullValue(attrType)
	}
	return tftypes.NewValue(objType, filled)
}

// agsNullValue builds a null tftypes value of the given type.
func agsNullValue(attrType tftypes.Type) tftypes.Value {
	switch t := attrType.(type) {
	case tftypes.List:
		return tftypes.NewValue(t, nil)
	case tftypes.Set:
		return tftypes.NewValue(t, nil)
	case tftypes.Map:
		return tftypes.NewValue(t, nil)
	case tftypes.Object:
		attrs := make(map[string]tftypes.Value, len(t.AttributeTypes))
		for name, nestedType := range t.AttributeTypes {
			attrs[name] = agsNullValue(nestedType)
		}
		return tftypes.NewValue(t, attrs)
	default:
		return tftypes.NewValue(attrType, nil)
	}
}

// agsSandboxToolValue builds a raw object value for the resource schema. A
// non-nil sandboxTool overrides the tool related attributes; empty strings keep
// the corresponding attribute null.
type agsSandboxToolValue struct {
	// toolId is written to both `id` and `tool_id`; leave it empty to keep them
	// null (the create case).
	toolId           string
	description      string
	defaultTimeout   string
	roleArn          string
	clientToken      string
	networkMode      string
	subnetId         string
	securityGroupId  string
	customImage      string
	includeComputer  bool
	includeLogs      bool
	includeStorage   bool
	includeTags      bool
	persistent       *bool
	includeAllNested bool
}

func agsSandboxToolRawValue(t *testing.T, spec agsSandboxToolValue) tftypes.Value {
	t.Helper()
	objType := agsSandboxToolObjectType(t)

	nullString := tftypes.NewValue(tftypes.String, nil)
	str := func(s string) tftypes.Value {
		if s == "" {
			return nullString
		}
		return tftypes.NewValue(tftypes.String, s)
	}
	nullList := func(elem tftypes.Type) tftypes.Value { return tftypes.NewValue(tftypes.List{ElementType: elem}, nil) }
	nullSet := func(elem tftypes.Type) tftypes.Value { return tftypes.NewValue(tftypes.Set{ElementType: elem}, nil) }

	// network_configuration
	networkTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"network_mode": tftypes.String,
		"vpc_config": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"subnet_ids":         tftypes.Set{ElementType: tftypes.String},
			"security_group_ids": tftypes.Set{ElementType: tftypes.String},
		}}},
	}}

	var networkConfiguration tftypes.Value
	if spec.networkMode == "" {
		networkConfiguration = nullList(networkTypes)
	} else {
		vpcTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"subnet_ids":         tftypes.Set{ElementType: tftypes.String},
			"security_group_ids": tftypes.Set{ElementType: tftypes.String},
		}}
		vpcValues := map[string]tftypes.Value{
			"subnet_ids":         nullSet(tftypes.String),
			"security_group_ids": nullSet(tftypes.String),
		}
		if spec.subnetId != "" {
			vpcValues["subnet_ids"] = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, spec.subnetId)})
		}
		if spec.securityGroupId != "" {
			vpcValues["security_group_ids"] = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, spec.securityGroupId)})
		}

		vpcList := tftypes.NewValue(tftypes.List{ElementType: vpcTypes}, []tftypes.Value{agsNestedObjectValue(vpcTypes, vpcValues)})
		networkConfiguration = tftypes.NewValue(tftypes.List{ElementType: networkTypes}, []tftypes.Value{
			agsNestedObjectValue(networkTypes, map[string]tftypes.Value{
				"network_mode": tftypes.NewValue(tftypes.String, spec.networkMode),
				"vpc_config":   vpcList,
			}),
		})
	}

	// tags
	tagTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"key": tftypes.String, "value": tftypes.String}}
	tags := nullList(tagTypes)
	if spec.includeTags {
		tags = tftypes.NewValue(tftypes.List{ElementType: tagTypes}, []tftypes.Value{
			agsNestedObjectValue(tagTypes, map[string]tftypes.Value{
				"key":   tftypes.NewValue(tftypes.String, "env"),
				"value": tftypes.NewValue(tftypes.String, "test"),
			}),
		})
	}

	// custom_configuration
	customTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"image":               tftypes.String,
		"image_registry_type": tftypes.String,
		"command":             tftypes.Set{ElementType: tftypes.String},
		"args":                tftypes.Set{ElementType: tftypes.String},
		"env":                 tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"name": tftypes.String, "value": tftypes.String}}},
		"ports":               tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"name": tftypes.String, "port": tftypes.Number, "protocol": tftypes.String}}},
		"resources":           tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"cpu": tftypes.String, "memory": tftypes.String, "storage": tftypes.String}}},
		"probe": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"http_get":          tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"path": tftypes.String, "port": tftypes.Number, "scheme": tftypes.String}}},
			"ready_timeout_ms":  tftypes.Number,
			"probe_timeout_ms":  tftypes.Number,
			"probe_period_ms":   tftypes.Number,
			"success_threshold": tftypes.Number,
			"failure_threshold": tftypes.Number,
		}}},
		"dns_config": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"servers":  tftypes.Set{ElementType: tftypes.String},
			"searches": tftypes.Set{ElementType: tftypes.String},
			"options":  tftypes.Set{ElementType: tftypes.String},
		}}},
		"security_context": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"privileged":   tftypes.Bool,
			"capabilities": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{"add": tftypes.Set{ElementType: tftypes.String}, "drop": tftypes.Set{ElementType: tftypes.String}}}},
		}}},
	}}

	customConfiguration := nullList(customTypes)
	if spec.customImage != "" {
		envTypes := customTypes.AttributeTypes["env"].(tftypes.List).ElementType.(tftypes.Object)
		portTypes := customTypes.AttributeTypes["ports"].(tftypes.List).ElementType.(tftypes.Object)
		resourcesTypes := customTypes.AttributeTypes["resources"].(tftypes.List).ElementType.(tftypes.Object)
		probeTypes := customTypes.AttributeTypes["probe"].(tftypes.List).ElementType.(tftypes.Object)
		httpGetTypes := probeTypes.AttributeTypes["http_get"].(tftypes.List).ElementType.(tftypes.Object)
		dnsTypes := customTypes.AttributeTypes["dns_config"].(tftypes.List).ElementType.(tftypes.Object)
		securityTypes := customTypes.AttributeTypes["security_context"].(tftypes.List).ElementType.(tftypes.Object)
		capabilityTypes := securityTypes.AttributeTypes["capabilities"].(tftypes.List).ElementType.(tftypes.Object)

		customConfiguration = tftypes.NewValue(tftypes.List{ElementType: customTypes}, []tftypes.Value{
			agsNestedObjectValue(customTypes, map[string]tftypes.Value{
				"image":               tftypes.NewValue(tftypes.String, spec.customImage),
				"image_registry_type": tftypes.NewValue(tftypes.String, "enterprise"),
				"command":             tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "/bin/sh")}),
				"args":                tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "-c")}),
				"env": tftypes.NewValue(tftypes.List{ElementType: envTypes}, []tftypes.Value{
					agsNestedObjectValue(envTypes, map[string]tftypes.Value{
						"name":  tftypes.NewValue(tftypes.String, "FOO"),
						"value": tftypes.NewValue(tftypes.String, "bar"),
					}),
				}),
				"ports": tftypes.NewValue(tftypes.List{ElementType: portTypes}, []tftypes.Value{
					agsNestedObjectValue(portTypes, map[string]tftypes.Value{
						"name":     tftypes.NewValue(tftypes.String, "http"),
						"port":     tftypes.NewValue(tftypes.Number, 8080),
						"protocol": tftypes.NewValue(tftypes.String, "TCP"),
					}),
				}),
				"resources": tftypes.NewValue(tftypes.List{ElementType: resourcesTypes}, []tftypes.Value{
					agsNestedObjectValue(resourcesTypes, map[string]tftypes.Value{
						"cpu":     tftypes.NewValue(tftypes.String, "1"),
						"memory":  tftypes.NewValue(tftypes.String, "2Gi"),
						"storage": tftypes.NewValue(tftypes.String, "10Gi"),
					}),
				}),
				"probe": tftypes.NewValue(tftypes.List{ElementType: probeTypes}, []tftypes.Value{
					agsNestedObjectValue(probeTypes, map[string]tftypes.Value{
						"http_get": tftypes.NewValue(tftypes.List{ElementType: httpGetTypes}, []tftypes.Value{
							agsNestedObjectValue(httpGetTypes, map[string]tftypes.Value{
								"path":   tftypes.NewValue(tftypes.String, "/healthz"),
								"port":   tftypes.NewValue(tftypes.Number, 8080),
								"scheme": tftypes.NewValue(tftypes.String, "HTTP"),
							}),
						}),
						"ready_timeout_ms":  tftypes.NewValue(tftypes.Number, 30000),
						"probe_timeout_ms":  tftypes.NewValue(tftypes.Number, 1000),
						"probe_period_ms":   tftypes.NewValue(tftypes.Number, 5000),
						"success_threshold": tftypes.NewValue(tftypes.Number, 1),
						"failure_threshold": tftypes.NewValue(tftypes.Number, 3),
					}),
				}),
				"dns_config": tftypes.NewValue(tftypes.List{ElementType: dnsTypes}, []tftypes.Value{
					agsNestedObjectValue(dnsTypes, map[string]tftypes.Value{
						"servers":  tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "10.0.0.1")}),
						"searches": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "svc.local")}),
						"options":  tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "ndots:5")}),
					}),
				}),
				"security_context": tftypes.NewValue(tftypes.List{ElementType: securityTypes}, []tftypes.Value{
					agsNestedObjectValue(securityTypes, map[string]tftypes.Value{
						"privileged": tftypes.NewValue(tftypes.Bool, true),
						"capabilities": tftypes.NewValue(tftypes.List{ElementType: capabilityTypes}, []tftypes.Value{
							agsNestedObjectValue(capabilityTypes, map[string]tftypes.Value{
								"add":  tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "SYS_ADMIN")}),
								"drop": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "ALL")}),
							}),
						}),
					}),
				}),
			}),
		})
	}

	// storage_mounts
	cosTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"endpoint": tftypes.String, "bucket_name": tftypes.String, "bucket_path": tftypes.String}}
	imageTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"reference": tftypes.String, "image_registry_type": tftypes.String, "sub_path": tftypes.String}}
	cfsTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"file_system_id": tftypes.String, "path": tftypes.String}}
	bucketTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"library_id": tftypes.String, "space_id": tftypes.String, "access_domain": tftypes.String}}
	sourceTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"cos":          tftypes.List{ElementType: cosTypes},
		"image":        tftypes.List{ElementType: imageTypes},
		"cfs":          tftypes.List{ElementType: cfsTypes},
		"agent_bucket": tftypes.List{ElementType: bucketTypes},
	}}
	volumeTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"volume_name": tftypes.String, "volume_id": tftypes.String}}
	volumeTemplateTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"volume_template_name": tftypes.String, "reuse_key": tftypes.String}}
	storageMountTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name":            tftypes.String,
		"storage_source":  tftypes.List{ElementType: sourceTypes},
		"mount_path":      tftypes.String,
		"read_only":       tftypes.Bool,
		"volume":          tftypes.List{ElementType: volumeTypes},
		"volume_template": tftypes.List{ElementType: volumeTemplateTypes},
	}}

	storageMounts := nullList(storageMountTypes)
	if spec.includeStorage {
		storageMounts = tftypes.NewValue(tftypes.List{ElementType: storageMountTypes}, []tftypes.Value{
			agsNestedObjectValue(storageMountTypes, map[string]tftypes.Value{
				"name":       tftypes.NewValue(tftypes.String, "rootfs"),
				"mount_path": tftypes.NewValue(tftypes.String, "/mnt/rootfs"),
				"read_only":  tftypes.NewValue(tftypes.Bool, true),
				"storage_source": tftypes.NewValue(tftypes.List{ElementType: sourceTypes}, []tftypes.Value{
					agsNestedObjectValue(sourceTypes, map[string]tftypes.Value{
						"cos": tftypes.NewValue(tftypes.List{ElementType: cosTypes}, []tftypes.Value{
							agsNestedObjectValue(cosTypes, map[string]tftypes.Value{
								"endpoint":    tftypes.NewValue(tftypes.String, "cos.ap-guangzhou.myqcloud.com"),
								"bucket_name": tftypes.NewValue(tftypes.String, "bucket-1250000000"),
								"bucket_path": tftypes.NewValue(tftypes.String, "/data"),
							}),
						}),
						"image": tftypes.NewValue(tftypes.List{ElementType: imageTypes}, []tftypes.Value{
							agsNestedObjectValue(imageTypes, map[string]tftypes.Value{
								"reference":           tftypes.NewValue(tftypes.String, "ccr.ccs.tencentcloudcr.com/demo:1.0"),
								"image_registry_type": tftypes.NewValue(tftypes.String, "personal"),
								"sub_path":            tftypes.NewValue(tftypes.String, "/opt"),
							}),
						}),
						"cfs": tftypes.NewValue(tftypes.List{ElementType: cfsTypes}, []tftypes.Value{
							agsNestedObjectValue(cfsTypes, map[string]tftypes.Value{
								"file_system_id": tftypes.NewValue(tftypes.String, "cfs-12345678"),
								"path":           tftypes.NewValue(tftypes.String, "/shared"),
							}),
						}),
						"agent_bucket": tftypes.NewValue(tftypes.List{ElementType: bucketTypes}, []tftypes.Value{
							agsNestedObjectValue(bucketTypes, map[string]tftypes.Value{
								"library_id":    tftypes.NewValue(tftypes.String, "lib-1"),
								"space_id":      tftypes.NewValue(tftypes.String, "space-1"),
								"access_domain": tftypes.NewValue(tftypes.String, "bucket.example.com"),
							}),
						}),
					}),
				}),
				"volume": tftypes.NewValue(tftypes.List{ElementType: volumeTypes}, []tftypes.Value{
					agsNestedObjectValue(volumeTypes, map[string]tftypes.Value{
						"volume_name": tftypes.NewValue(tftypes.String, "vol-name"),
						"volume_id":   tftypes.NewValue(tftypes.String, "vol-12345678"),
					}),
				}),
				"volume_template": tftypes.NewValue(tftypes.List{ElementType: volumeTemplateTypes}, []tftypes.Value{
					agsNestedObjectValue(volumeTemplateTypes, map[string]tftypes.Value{
						"volume_template_name": tftypes.NewValue(tftypes.String, "vtpl-1"),
						"reuse_key":            tftypes.NewValue(tftypes.String, "${metadata.userId}"),
					}),
				}),
			}),
		})
	}

	// computer_configuration
	waaTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"image_id": tftypes.String}}
	osWorldTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"version": tftypes.String}}
	computerTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"waa_configuration":      tftypes.List{ElementType: waaTypes},
		"os_world_configuration": tftypes.List{ElementType: osWorldTypes},
	}}
	computerConfiguration := nullList(computerTypes)
	if spec.includeComputer {
		computerConfiguration = tftypes.NewValue(tftypes.List{ElementType: computerTypes}, []tftypes.Value{
			agsNestedObjectValue(computerTypes, map[string]tftypes.Value{
				"waa_configuration": tftypes.NewValue(tftypes.List{ElementType: waaTypes}, []tftypes.Value{
					agsNestedObjectValue(waaTypes, map[string]tftypes.Value{
						"image_id": tftypes.NewValue(tftypes.String, "img-waa-1"),
					}),
				}),
				"os_world_configuration": tftypes.NewValue(tftypes.List{ElementType: osWorldTypes}, []tftypes.Value{
					agsNestedObjectValue(osWorldTypes, map[string]tftypes.Value{
						"version": tftypes.NewValue(tftypes.String, "osworld2"),
					}),
				}),
			}),
		})
	}

	// log_configuration
	clsTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"topic_id": tftypes.String}}
	logSourcesTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"files": tftypes.Set{ElementType: tftypes.String}}}
	logTypes := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"cls_config":  tftypes.List{ElementType: clsTypes},
		"log_sources": tftypes.List{ElementType: logSourcesTypes},
	}}
	logConfiguration := nullList(logTypes)
	if spec.includeLogs {
		logConfiguration = tftypes.NewValue(tftypes.List{ElementType: logTypes}, []tftypes.Value{
			agsNestedObjectValue(logTypes, map[string]tftypes.Value{
				"cls_config": tftypes.NewValue(tftypes.List{ElementType: clsTypes}, []tftypes.Value{
					agsNestedObjectValue(clsTypes, map[string]tftypes.Value{
						"topic_id": tftypes.NewValue(tftypes.String, "topic-1"),
					}),
				}),
				"log_sources": tftypes.NewValue(tftypes.List{ElementType: logSourcesTypes}, []tftypes.Value{
					agsNestedObjectValue(logSourcesTypes, map[string]tftypes.Value{
						"files": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "/logs/app.log")}),
					}),
				}),
			}),
		})
	}

	persistent := tftypes.NewValue(tftypes.Bool, nil)
	if spec.persistent != nil {
		persistent = tftypes.NewValue(tftypes.Bool, *spec.persistent)
	}

	return agsNestedObjectValue(objType, map[string]tftypes.Value{
		"id":                     str(spec.toolId),
		"tool_id":                str(spec.toolId),
		"tool_name":              tftypes.NewValue(tftypes.String, agsTestToolName),
		"tool_type":              tftypes.NewValue(tftypes.String, "custom"),
		"network_configuration":  networkConfiguration,
		"description":            str(spec.description),
		"default_timeout":        str(spec.defaultTimeout),
		"tags":                   tags,
		"client_token":           str(spec.clientToken),
		"role_arn":               str(spec.roleArn),
		"storage_mounts":         storageMounts,
		"custom_configuration":   customConfiguration,
		"computer_configuration": computerConfiguration,
		"log_configuration":      logConfiguration,
		"persistent":             persistent,
	})
}

// agsSandboxToolPlan builds a tfsdk.Plan carrying the given raw value.
func agsSandboxToolPlan(t *testing.T, raw tftypes.Value) tfsdk.Plan {
	t.Helper()
	return tfsdk.Plan{Raw: raw, Schema: agsSandboxToolSchema(t)}
}

// agsSandboxToolState builds a tfsdk.State carrying the given raw value.
func agsSandboxToolState(t *testing.T, raw tftypes.Value) tfsdk.State {
	t.Helper()
	return tfsdk.State{Raw: raw, Schema: agsSandboxToolSchema(t)}
}

// agsSandboxToolCreateRequest builds a CreateRequest.
func agsSandboxToolCreateRequest(t *testing.T, spec agsSandboxToolValue) resource.CreateRequest {
	raw := agsSandboxToolRawValue(t, spec)
	return resource.CreateRequest{Plan: agsSandboxToolPlan(t, raw)}
}

// agsSandboxToolUpdateRequest builds an UpdateRequest from two specs.
func agsSandboxToolUpdateRequest(t *testing.T, stateSpec, planSpec agsSandboxToolValue) resource.UpdateRequest {
	return resource.UpdateRequest{
		Plan:  agsSandboxToolPlan(t, agsSandboxToolRawValue(t, planSpec)),
		State: agsSandboxToolState(t, agsSandboxToolRawValue(t, stateSpec)),
	}
}

// agsSandboxToolReadRequest builds a ReadRequest.
func agsSandboxToolReadRequest(t *testing.T, spec agsSandboxToolValue) resource.ReadRequest {
	return resource.ReadRequest{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
}

// agsSandboxToolDeleteRequest builds a DeleteRequest.
func agsSandboxToolDeleteRequest(t *testing.T, spec agsSandboxToolValue) resource.DeleteRequest {
	return resource.DeleteRequest{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
}

// agsSandboxToolTestSetup patches the shared client so that the resource talks
// to a mock ags client. The poll interval is shrunk so that multi-step waits do
// not slow the test suite down.
func agsSandboxToolTestSetup(t *testing.T) (*gomonkey.Patches, *ags.Client) {
	t.Helper()

	patches := gomonkey.NewPatches()
	client := &connectivity.TencentCloudClient{}
	agsClient := &ags.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)

	originalInterval := AGS_SANDBOX_TOOL_STATE_POLL_INTERVAL
	AGS_SANDBOX_TOOL_STATE_POLL_INTERVAL = time.Millisecond
	t.Cleanup(func() {
		patches.Reset()
		AGS_SANDBOX_TOOL_STATE_POLL_INTERVAL = originalInterval
	})

	return patches, agsClient
}

// agsSandboxToolResourceWithClient builds the resource and injects the client.
func agsSandboxToolResourceWithClient(client *connectivity.TencentCloudClient) *AgsSandboxToolResource {
	r := NewAgsSandboxToolResource().(*AgsSandboxToolResource)
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}}, configureResp)
	return r
}

// agsSandboxToolResponse builds a DescribeSandboxToolList response holding the
// given sandbox tool.
func agsSandboxToolResponse(status string, item *ags.SandboxTool) *ags.DescribeSandboxToolListResponse {
	if item == nil {
		return &ags.DescribeSandboxToolListResponse{Response: &ags.DescribeSandboxToolListResponseParams{}}
	}
	if item.ToolId == nil {
		item.ToolId = agsString(agsTestToolId)
	}
	if item.Status == nil {
		item.Status = agsString(status)
	}
	return &ags.DescribeSandboxToolListResponse{
		Response: &ags.DescribeSandboxToolListResponseParams{
			SandboxToolSet: []*ags.SandboxTool{item},
			TotalCount:     agsInt64(1),
			RequestId:      agsString("fake-request-id"),
		},
	}
}

// agsSandboxToolFullItem builds a fully populated sandbox tool as returned by
// DescribeSandboxToolList.
func agsSandboxToolFullItem() *ags.SandboxTool {
	return &ags.SandboxTool{
		ToolId:                agsString(agsTestToolId),
		ToolName:              agsString(agsTestToolName),
		ToolType:              agsString("custom"),
		Status:                agsString("ACTIVE"),
		StatusReason:          agsString("ready"),
		Description:           agsString("managed by terraform"),
		Persistent:            agsBool(true),
		DefaultTimeoutSeconds: agsUint64(300),
		CreateTime:            agsString("2025-09-20T10:00:00+08:00"),
		UpdateTime:            agsString("2025-09-20T10:05:00+08:00"),
		RoleArn:               agsString("qcs::cam::uin/100000000001:uin/100000000001"),
		NetworkConfiguration: &ags.NetworkConfiguration{
			NetworkMode: agsString("VPC"),
			VpcConfig: &ags.VPCConfig{
				SubnetIds:        []*string{agsString("subnet-1")},
				SecurityGroupIds: []*string{agsString("sg-1")},
			},
		},
		Tags: []*ags.Tag{
			{Key: agsString("env"), Value: agsString("test")},
		},
		StorageMounts: []*ags.StorageMount{
			{
				Name:      agsString("rootfs"),
				MountPath: agsString("/mnt/rootfs"),
				ReadOnly:  agsBool(true),
				StorageSource: &ags.StorageSource{
					Cos: &ags.CosStorageSource{
						Endpoint:   agsString("cos.ap-guangzhou.myqcloud.com"),
						BucketName: agsString("bucket-1250000000"),
						BucketPath: agsString("/data"),
					},
					Image: &ags.ImageStorageSource{
						Reference:         agsString("ccr.ccs.tencentcloudcr.com/demo:1.0"),
						ImageRegistryType: agsString("personal"),
						SubPath:           agsString("/opt"),
					},
					Cfs: &ags.CfsStorageSource{
						FileSystemId: agsString("cfs-12345678"),
						Path:         agsString("/shared"),
					},
					AgentBucket: &ags.AgentBucketStorageSource{
						LibraryId:    agsString("lib-1"),
						SpaceId:      agsString("space-1"),
						AccessDomain: agsString("bucket.example.com"),
					},
				},
				Volume: &ags.VolumeRef{
					VolumeName: agsString("vol-name"),
					VolumeId:   agsString("vol-12345678"),
				},
				VolumeTemplate: &ags.VolumeTemplateRef{
					VolumeTemplateName: agsString("vtpl-1"),
					ReuseKey:           agsString("${metadata.userId}"),
				},
			},
		},
		CustomConfiguration: &ags.CustomConfigurationDetail{
			Image:             agsString("ccr.ccs.tencentcloudcr.com/demo:1.0"),
			ImageRegistryType: agsString("enterprise"),
			ImageDigest:       agsString("sha256:deadbeef"),
			Command:           []*string{agsString("/bin/sh")},
			Args:              []*string{agsString("-c")},
			Env: []*ags.EnvVar{
				{Name: agsString("FOO"), Value: agsString("bar")},
			},
			Ports: []*ags.PortConfiguration{
				{Name: agsString("http"), Port: agsInt64(8080), Protocol: agsString("TCP")},
			},
			Resources: &ags.ResourceConfiguration{
				CPU:     agsString("1"),
				Memory:  agsString("2Gi"),
				Storage: agsString("10Gi"),
			},
			Probe: &ags.ProbeConfiguration{
				HttpGet: &ags.HttpGetAction{
					Path:   agsString("/healthz"),
					Port:   agsInt64(8080),
					Scheme: agsString("HTTP"),
				},
				ReadyTimeoutMs:   agsInt64(30000),
				ProbeTimeoutMs:   agsInt64(1000),
				ProbePeriodMs:    agsInt64(5000),
				SuccessThreshold: agsInt64(1),
				FailureThreshold: agsInt64(3),
			},
			DNSConfig: &ags.DNSConfig{
				Servers:  []*string{agsString("10.0.0.1")},
				Searches: []*string{agsString("svc.local")},
				Options:  []*string{agsString("ndots:5")},
			},
			SecurityContext: &ags.SecurityContext{
				Privileged: agsBool(true),
				Capabilities: &ags.Capabilities{
					Add:  []*string{agsString("SYS_ADMIN")},
					Drop: []*string{agsString("ALL")},
				},
			},
		},
		ComputerConfiguration: &ags.ComputerConfiguration{
			WAAConfiguration:     &ags.WAAConfiguration{ImageId: agsString("img-waa-1")},
			OSWorldConfiguration: &ags.OSWorldConfiguration{Version: agsString("osworld2")},
		},
		LogConfiguration: &ags.LogConfiguration{
			CLSConfig:  &ags.CLSConfig{TopicId: agsString("topic-1")},
			LogSources: &ags.LogSources{Files: []*string{agsString("/logs/app.log")}},
		},
	}
}

// agsSandboxToolFullSpec is the configuration matching agsSandboxToolFullItem.
func agsSandboxToolFullSpec() agsSandboxToolValue {
	persistent := true
	return agsSandboxToolValue{
		toolId:          agsTestToolId,
		description:     "managed by terraform",
		defaultTimeout:  "300s",
		roleArn:         "qcs::cam::uin/100000000001:uin/100000000001",
		clientToken:     "token-1",
		networkMode:     "VPC",
		subnetId:        "subnet-1",
		securityGroupId: "sg-1",
		customImage:     "ccr.ccs.tencentcloudcr.com/demo:1.0",
		includeComputer: true,
		includeLogs:     true,
		includeStorage:  true,
		includeTags:     true,
		persistent:      &persistent,
	}
}

// assertAgsSandboxToolState verifies every leaf attribute written into state.
func assertAgsSandboxToolState(t *testing.T, state tfsdk.State) {
	t.Helper()
	ctx := context.Background()

	var m AgsSandboxToolResourceModel
	if diags := state.Get(ctx, &m); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}

	if got := m.Id.ValueString(); got != agsTestToolId {
		t.Errorf("id = %q, want %q", got, agsTestToolId)
	}
	if got := m.ToolId.ValueString(); got != agsTestToolId {
		t.Errorf("tool_id = %q, want %q", got, agsTestToolId)
	}
	if got := m.ToolName.ValueString(); got != agsTestToolName {
		t.Errorf("tool_name = %q, want %q", got, agsTestToolName)
	}
	if got := m.ToolType.ValueString(); got != "custom" {
		t.Errorf("tool_type = %q, want %q", got, "custom")
	}
	if got := m.Status.ValueString(); got != "ACTIVE" {
		t.Errorf("status = %q, want %q", got, "ACTIVE")
	}
	if got := m.StatusReason.ValueString(); got != "ready" {
		t.Errorf("status_reason = %q, want %q", got, "ready")
	}
	if got := m.Description.ValueString(); got != "managed by terraform" {
		t.Errorf("description = %q, want %q", got, "managed by terraform")
	}
	if got := m.DefaultTimeout.ValueString(); got != "300s" {
		t.Errorf("default_timeout = %q, want %q", got, "300s")
	}
	if got := m.RoleArn.ValueString(); got != "qcs::cam::uin/100000000001:uin/100000000001" {
		t.Errorf("role_arn = %q", got)
	}
	if got := m.CreateTime.ValueString(); got != "2025-09-20T10:00:00+08:00" {
		t.Errorf("create_time = %q", got)
	}
	if got := m.UpdateTime.ValueString(); got != "2025-09-20T10:05:00+08:00" {
		t.Errorf("update_time = %q", got)
	}
	if m.Persistent.IsNull() || m.Persistent.ValueBool() != true {
		t.Errorf("persistent = %v, want true", m.Persistent)
	}

	// network_configuration
	if m.NetworkConfiguration.IsNull() {
		t.Fatal("network_configuration is null")
	}
	var networks []agsSandboxToolNetworkConfigurationModel
	if diags := m.NetworkConfiguration.ElementsAs(ctx, &networks, false); diags.HasError() {
		t.Fatalf("failed to read network_configuration: %v", diags)
	}
	if len(networks) != 1 {
		t.Fatalf("network_configuration length = %d, want 1", len(networks))
	}
	if got := networks[0].NetworkMode.ValueString(); got != "VPC" {
		t.Errorf("network_mode = %q, want VPC", got)
	}
	if len(networks[0].VpcConfig) != 1 {
		t.Fatalf("vpc_config length = %d, want 1", len(networks[0].VpcConfig))
	}
	if got := agsSetStrings(ctx, t, networks[0].VpcConfig[0].SubnetIds); len(got) != 1 || got[0] != "subnet-1" {
		t.Errorf("subnet_ids = %v, want [subnet-1]", got)
	}
	if got := agsSetStrings(ctx, t, networks[0].VpcConfig[0].SecurityGroupIds); len(got) != 1 || got[0] != "sg-1" {
		t.Errorf("security_group_ids = %v, want [sg-1]", got)
	}

	// tags
	var tags []agsSandboxToolTagModel
	if diags := m.Tags.ElementsAs(ctx, &tags, false); diags.HasError() {
		t.Fatalf("failed to read tags: %v", diags)
	}
	if len(tags) != 1 || tags[0].Key.ValueString() != "env" || tags[0].Value.ValueString() != "test" {
		t.Errorf("tags = %v, want [{env test}]", tags)
	}

	// storage_mounts
	var mounts []agsSandboxToolStorageMountModel
	if diags := m.StorageMounts.ElementsAs(ctx, &mounts, false); diags.HasError() {
		t.Fatalf("failed to read storage_mounts: %v", diags)
	}
	if len(mounts) != 1 {
		t.Fatalf("storage_mounts length = %d, want 1", len(mounts))
	}
	mount := mounts[0]
	if mount.Name.ValueString() != "rootfs" || mount.MountPath.ValueString() != "/mnt/rootfs" {
		t.Errorf("storage mount name/mount_path = %q/%q", mount.Name.ValueString(), mount.MountPath.ValueString())
	}
	if mount.ReadOnly.IsNull() || !mount.ReadOnly.ValueBool() {
		t.Errorf("read_only = %v, want true", mount.ReadOnly)
	}
	if len(mount.StorageSource) != 1 {
		t.Fatalf("storage_source length = %d, want 1", len(mount.StorageSource))
	}
	source := mount.StorageSource[0]
	if len(source.Cos) != 1 || source.Cos[0].BucketName.ValueString() != "bucket-1250000000" || source.Cos[0].BucketPath.ValueString() != "/data" {
		t.Errorf("cos = %v", source.Cos)
	}
	if len(source.Image) != 1 || source.Image[0].Reference.ValueString() != "ccr.ccs.tencentcloudcr.com/demo:1.0" || source.Image[0].SubPath.ValueString() != "/opt" {
		t.Errorf("image = %v", source.Image)
	}
	if len(source.Cfs) != 1 || source.Cfs[0].FileSystemId.ValueString() != "cfs-12345678" || source.Cfs[0].Path.ValueString() != "/shared" {
		t.Errorf("cfs = %v", source.Cfs)
	}
	if len(source.AgentBucket) != 1 || source.AgentBucket[0].LibraryId.ValueString() != "lib-1" || source.AgentBucket[0].AccessDomain.ValueString() != "bucket.example.com" {
		t.Errorf("agent_bucket = %v", source.AgentBucket)
	}
	if len(mount.Volume) != 1 || mount.Volume[0].VolumeId.ValueString() != "vol-12345678" || mount.Volume[0].VolumeName.ValueString() != "vol-name" {
		t.Errorf("volume = %v", mount.Volume)
	}
	if len(mount.VolumeTemplate) != 1 || mount.VolumeTemplate[0].VolumeTemplateName.ValueString() != "vtpl-1" || mount.VolumeTemplate[0].ReuseKey.ValueString() != "${metadata.userId}" {
		t.Errorf("volume_template = %v", mount.VolumeTemplate)
	}

	// custom_configuration
	var customs []agsSandboxToolCustomConfigurationModel
	if diags := m.CustomConfiguration.ElementsAs(ctx, &customs, false); diags.HasError() {
		t.Fatalf("failed to read custom_configuration: %v", diags)
	}
	if len(customs) != 1 {
		t.Fatalf("custom_configuration length = %d, want 1", len(customs))
	}
	custom := customs[0]
	if custom.Image.ValueString() != "ccr.ccs.tencentcloudcr.com/demo:1.0" || custom.ImageRegistryType.ValueString() != "enterprise" {
		t.Errorf("custom image = %q/%q", custom.Image.ValueString(), custom.ImageRegistryType.ValueString())
	}
	if got := agsSetStrings(ctx, t, custom.Command); len(got) != 1 || got[0] != "/bin/sh" {
		t.Errorf("command = %v", got)
	}
	if got := agsSetStrings(ctx, t, custom.Args); len(got) != 1 || got[0] != "-c" {
		t.Errorf("args = %v", got)
	}
	if len(custom.Env) != 1 || custom.Env[0].Name.ValueString() != "FOO" || custom.Env[0].Value.ValueString() != "bar" {
		t.Errorf("env = %v", custom.Env)
	}
	if len(custom.Ports) != 1 || custom.Ports[0].Port.ValueInt64() != 8080 || custom.Ports[0].Name.ValueString() != "http" || custom.Ports[0].Protocol.ValueString() != "TCP" {
		t.Errorf("ports = %v", custom.Ports)
	}
	if len(custom.Resources) != 1 || custom.Resources[0].Cpu.ValueString() != "1" || custom.Resources[0].Memory.ValueString() != "2Gi" || custom.Resources[0].Storage.ValueString() != "10Gi" {
		t.Errorf("resources = %v", custom.Resources)
	}
	if len(custom.Probe) != 1 {
		t.Fatalf("probe length = %d, want 1", len(custom.Probe))
	}
	probe := custom.Probe[0]
	if probe.ReadyTimeoutMs.ValueInt64() != 30000 || probe.ProbeTimeoutMs.ValueInt64() != 1000 || probe.ProbePeriodMs.ValueInt64() != 5000 ||
		probe.SuccessThreshold.ValueInt64() != 1 || probe.FailureThreshold.ValueInt64() != 3 {
		t.Errorf("probe = %v", probe)
	}
	if len(probe.HttpGet) != 1 || probe.HttpGet[0].Path.ValueString() != "/healthz" || probe.HttpGet[0].Port.ValueInt64() != 8080 || probe.HttpGet[0].Scheme.ValueString() != "HTTP" {
		t.Errorf("http_get = %v", probe.HttpGet)
	}
	if len(custom.DnsConfig) != 1 {
		t.Fatalf("dns_config length = %d, want 1", len(custom.DnsConfig))
	}
	if got := agsSetStrings(ctx, t, custom.DnsConfig[0].Servers); len(got) != 1 || got[0] != "10.0.0.1" {
		t.Errorf("dns servers = %v", got)
	}
	if got := agsSetStrings(ctx, t, custom.DnsConfig[0].Searches); len(got) != 1 || got[0] != "svc.local" {
		t.Errorf("dns searches = %v", got)
	}
	if got := agsSetStrings(ctx, t, custom.DnsConfig[0].Options); len(got) != 1 || got[0] != "ndots:5" {
		t.Errorf("dns options = %v", got)
	}
	if len(custom.SecurityContext) != 1 {
		t.Fatalf("security_context length = %d, want 1", len(custom.SecurityContext))
	}
	securityContext := custom.SecurityContext[0]
	if securityContext.Privileged.IsNull() || !securityContext.Privileged.ValueBool() {
		t.Errorf("privileged = %v, want true", securityContext.Privileged)
	}
	if len(securityContext.Capabilities) != 1 {
		t.Fatalf("capabilities length = %d, want 1", len(securityContext.Capabilities))
	}
	if got := agsSetStrings(ctx, t, securityContext.Capabilities[0].Add); len(got) != 1 || got[0] != "SYS_ADMIN" {
		t.Errorf("capabilities add = %v", got)
	}
	if got := agsSetStrings(ctx, t, securityContext.Capabilities[0].Drop); len(got) != 1 || got[0] != "ALL" {
		t.Errorf("capabilities drop = %v", got)
	}

	// computer_configuration
	var computers []agsSandboxToolComputerConfigurationModel
	if diags := m.ComputerConfiguration.ElementsAs(ctx, &computers, false); diags.HasError() {
		t.Fatalf("failed to read computer_configuration: %v", diags)
	}
	if len(computers) != 1 {
		t.Fatalf("computer_configuration length = %d, want 1", len(computers))
	}
	if len(computers[0].WaaConfiguration) != 1 || computers[0].WaaConfiguration[0].ImageId.ValueString() != "img-waa-1" {
		t.Errorf("waa_configuration = %v", computers[0].WaaConfiguration)
	}
	if len(computers[0].OsWorldConfiguration) != 1 || computers[0].OsWorldConfiguration[0].Version.ValueString() != "osworld2" {
		t.Errorf("os_world_configuration = %v", computers[0].OsWorldConfiguration)
	}

	// log_configuration
	var logConfigs []agsSandboxToolLogConfigurationModel
	if diags := m.LogConfiguration.ElementsAs(ctx, &logConfigs, false); diags.HasError() {
		t.Fatalf("failed to read log_configuration: %v", diags)
	}
	if len(logConfigs) != 1 {
		t.Fatalf("log_configuration length = %d, want 1", len(logConfigs))
	}
	if len(logConfigs[0].ClsConfig) != 1 || logConfigs[0].ClsConfig[0].TopicId.ValueString() != "topic-1" {
		t.Errorf("cls_config = %v", logConfigs[0].ClsConfig)
	}
	if len(logConfigs[0].LogSources) != 1 {
		t.Fatalf("log_sources length = %d, want 1", len(logConfigs[0].LogSources))
	}
	if got := agsSetStrings(ctx, t, logConfigs[0].LogSources[0].Files); len(got) != 1 || got[0] != "/logs/app.log" {
		t.Errorf("log files = %v", got)
	}
}

// agsSetStrings converts a types.Set of strings into a plain slice.
func agsSetStrings(ctx context.Context, t *testing.T, set types.Set) []string {
	t.Helper()
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	var out []string
	if diags := set.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("failed to read string set: %v", diags)
	}
	return out
}

// TestAgsSandboxTool_MetadataAndSchema verifies the type name, the schema shape
// and that read only outputs are not user configurable.
func TestAgsSandboxTool_MetadataAndSchema(t *testing.T) {
	r := NewAgsSandboxToolResource()

	metadataResp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, metadataResp)
	if metadataResp.TypeName != "tencentcloud_ags_sandbox_tool" {
		t.Fatalf("TypeName = %q, want tencentcloud_ags_sandbox_tool", metadataResp.TypeName)
	}

	s := agsSandboxToolSchema(t)

	required := map[string]bool{
		"tool_name":             true,
		"tool_type":             true,
		"network_configuration": true,
	}
	for name := range required {
		attrSchema, ok := s.Attributes[name]
		if !ok {
			t.Fatalf("missing attribute %q", name)
		}
		if !attrSchema.IsRequired() {
			t.Errorf("attribute %q should be required", name)
		}
	}

	computedOnly := []string{"id", "tool_id", "status", "status_reason", "create_time", "update_time"}
	for _, name := range computedOnly {
		attrSchema, ok := s.Attributes[name]
		if !ok {
			t.Fatalf("missing attribute %q", name)
		}
		if !attrSchema.IsComputed() || attrSchema.IsRequired() || attrSchema.IsOptional() {
			t.Errorf("attribute %q should be computed only", name)
		}
	}

	optionalNames := []string{"description", "tags", "storage_mounts", "custom_configuration", "computer_configuration", "log_configuration", "client_token"}
	for _, name := range optionalNames {
		attrSchema, ok := s.Attributes[name]
		if !ok {
			t.Fatalf("missing attribute %q", name)
		}
		if !attrSchema.IsOptional() {
			t.Errorf("attribute %q should be optional", name)
		}
	}

	if _, ok := s.Attributes["sandbox_tool_set"]; ok {
		t.Error("sandbox_tool_set is a Describe only output and must not be part of the resource schema")
	}
}

// TestAgsSandboxTool_Create_Success verifies a full create round trip: the
// request carries every configurable field and the state is refreshed from the
// describe result.
func TestAgsSandboxTool_Create_Success(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)

	var gotCreateRequest *ags.CreateSandboxToolRequest
	patches.ApplyMethodFunc(agsClient, "CreateSandboxToolWithContext", func(_ context.Context, request *ags.CreateSandboxToolRequest) (*ags.CreateSandboxToolResponse, error) {
		gotCreateRequest = request
		return &ags.CreateSandboxToolResponse{
			Response: &ags.CreateSandboxToolResponseParams{
				ToolId:    agsString(agsTestToolId),
				RequestId: agsString("fake-request-id"),
			},
		}, nil
	})

	// First poll returns CREATING, the second one ACTIVE: the identity must be
	// written to state before the wait, and the wait must converge.
	pollCount := 0
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		pollCount++
		if pollCount == 1 {
			return agsSandboxToolResponse("CREATING", &ags.SandboxTool{ToolId: agsString(agsTestToolId)}), nil
		}
		return agsSandboxToolResponse("ACTIVE", agsSandboxToolFullItem()), nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	planSpec := agsSandboxToolFullSpec()
	createResp := &resource.CreateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, planSpec))}
	r.Create(context.Background(), agsSandboxToolCreateRequest(t, planSpec), createResp)

	if createResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", createResp.Diagnostics)
	}
	if gotCreateRequest == nil {
		t.Fatal("CreateSandboxTool was not called")
	}

	if got := gotCreateRequest.ToolName; got == nil || *got != agsTestToolName {
		t.Errorf("ToolName = %v", got)
	}
	if got := gotCreateRequest.ToolType; got == nil || *got != "custom" {
		t.Errorf("ToolType = %v", got)
	}
	if got := gotCreateRequest.Description; got == nil || *got != "managed by terraform" {
		t.Errorf("Description = %v", got)
	}
	if got := gotCreateRequest.DefaultTimeout; got == nil || *got != "300s" {
		t.Errorf("DefaultTimeout = %v", got)
	}
	if got := gotCreateRequest.ClientToken; got == nil || *got != "token-1" {
		t.Errorf("ClientToken = %v", got)
	}
	if got := gotCreateRequest.RoleArn; got == nil || *got != "qcs::cam::uin/100000000001:uin/100000000001" {
		t.Errorf("RoleArn = %v", got)
	}
	if got := gotCreateRequest.Persistent; got == nil || !*got {
		t.Errorf("Persistent = %v", got)
	}
	if gotCreateRequest.NetworkConfiguration == nil || gotCreateRequest.NetworkConfiguration.NetworkMode == nil || *gotCreateRequest.NetworkConfiguration.NetworkMode != "VPC" {
		t.Fatalf("NetworkConfiguration = %v", gotCreateRequest.NetworkConfiguration)
	}
	if got := gotCreateRequest.NetworkConfiguration.VpcConfig; got == nil || len(got.SubnetIds) != 1 || len(got.SecurityGroupIds) != 1 {
		t.Errorf("VpcConfig = %v", got)
	}
	if len(gotCreateRequest.Tags) != 1 || *gotCreateRequest.Tags[0].Key != "env" || *gotCreateRequest.Tags[0].Value != "test" {
		t.Errorf("Tags = %v", gotCreateRequest.Tags)
	}
	if len(gotCreateRequest.StorageMounts) != 1 {
		t.Fatalf("StorageMounts = %v", gotCreateRequest.StorageMounts)
	}
	mount := gotCreateRequest.StorageMounts[0]
	if mount.ReadOnly == nil || !*mount.ReadOnly || mount.StorageSource == nil || mount.StorageSource.Cos == nil ||
		mount.StorageSource.Image == nil || mount.StorageSource.Cfs == nil || mount.StorageSource.AgentBucket == nil ||
		mount.Volume == nil || mount.VolumeTemplate == nil {
		t.Errorf("StorageMounts[0] is incomplete: %v", mount)
	}
	if gotCreateRequest.CustomConfiguration == nil {
		t.Fatal("CustomConfiguration is nil")
	}
	custom := gotCreateRequest.CustomConfiguration
	if len(custom.Command) != 1 || len(custom.Args) != 1 || len(custom.Env) != 1 || len(custom.Ports) != 1 {
		t.Errorf("CustomConfiguration commands/args/env/ports are incomplete")
	}
	if custom.Resources == nil || custom.Probe == nil || custom.Probe.HttpGet == nil || custom.DNSConfig == nil || custom.SecurityContext == nil ||
		custom.SecurityContext.Capabilities == nil {
		t.Errorf("CustomConfiguration nested blocks are incomplete")
	}
	if gotCreateRequest.ComputerConfiguration == nil || gotCreateRequest.ComputerConfiguration.WAAConfiguration == nil ||
		gotCreateRequest.ComputerConfiguration.OSWorldConfiguration == nil {
		t.Errorf("ComputerConfiguration is incomplete")
	}
	if gotCreateRequest.LogConfiguration == nil || gotCreateRequest.LogConfiguration.CLSConfig == nil || gotCreateRequest.LogConfiguration.LogSources == nil {
		t.Errorf("LogConfiguration is incomplete")
	}
	if pollCount < 2 {
		t.Errorf("describe poll count = %d, want at least 2", pollCount)
	}

	assertAgsSandboxToolState(t, createResp.State)
}

// TestAgsSandboxTool_Create_MinimalConfig verifies that nullable optional
// blocks are not sent and that state stays null for them.
func TestAgsSandboxTool_Create_MinimalConfig(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)

	var gotCreateRequest *ags.CreateSandboxToolRequest
	patches.ApplyMethodFunc(agsClient, "CreateSandboxToolWithContext", func(_ context.Context, request *ags.CreateSandboxToolRequest) (*ags.CreateSandboxToolResponse, error) {
		gotCreateRequest = request
		return &ags.CreateSandboxToolResponse{
			Response: &ags.CreateSandboxToolResponseParams{ToolId: agsString(agsTestToolId)},
		}, nil
	})
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		return agsSandboxToolResponse("ACTIVE", &ags.SandboxTool{
			ToolId:   agsString(agsTestToolId),
			ToolName: agsString(agsTestToolName),
			ToolType: agsString("custom"),
			Status:   agsString("ACTIVE"),
		}), nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	createResp := &resource.CreateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Create(context.Background(), agsSandboxToolCreateRequest(t, spec), createResp)

	if createResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", createResp.Diagnostics)
	}
	if gotCreateRequest.Description != nil {
		t.Errorf("Description should not be sent, got %v", *gotCreateRequest.Description)
	}
	if gotCreateRequest.DefaultTimeout != nil || gotCreateRequest.ClientToken != nil || gotCreateRequest.RoleArn != nil {
		t.Errorf("nullable string inputs should not be sent")
	}
	if gotCreateRequest.Persistent != nil {
		t.Errorf("Persistent should not be sent, got %v", *gotCreateRequest.Persistent)
	}
	if len(gotCreateRequest.StorageMounts) != 0 || gotCreateRequest.CustomConfiguration != nil ||
		gotCreateRequest.ComputerConfiguration != nil || gotCreateRequest.LogConfiguration != nil || len(gotCreateRequest.Tags) != 0 {
		t.Errorf("unset optional blocks must not be sent")
	}

	var m AgsSandboxToolResourceModel
	if diags := createResp.State.Get(context.Background(), &m); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if !m.StorageMounts.IsNull() || !m.CustomConfiguration.IsNull() || !m.ComputerConfiguration.IsNull() || !m.LogConfiguration.IsNull() || !m.Tags.IsNull() {
		t.Errorf("unset optional blocks must stay null in state")
	}
	if !m.DefaultTimeout.IsNull() || !m.RoleArn.IsNull() {
		t.Errorf("nullable inputs without cloud value must stay null in state")
	}
}

// TestAgsSandboxTool_Create_Error verifies that an API error surfaces a
// diagnostic and leaves an empty identity out of state.
func TestAgsSandboxTool_Create_Error(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "CreateSandboxToolWithContext", func(_ context.Context, request *ags.CreateSandboxToolRequest) (*ags.CreateSandboxToolResponse, error) {
		return nil, errors.New("Code=FailedOperation, Message=create failed")
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	createResp := &resource.CreateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Create(context.Background(), agsSandboxToolCreateRequest(t, spec), createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic")
	}
	if !containsString(createResp.Diagnostics.Errors()[0].Summary(), "Error creating AGS sandbox tool") {
		t.Errorf("unexpected summary: %s", createResp.Diagnostics.Errors()[0].Summary())
	}
	if !containsString(createResp.Diagnostics.Errors()[0].Detail(), "FailedOperation") {
		t.Errorf("unexpected detail: %s", createResp.Diagnostics.Errors()[0].Detail())
	}
}

// TestAgsSandboxTool_Create_EmptyToolId verifies that an empty ToolId is
// treated as an API failure.
func TestAgsSandboxTool_Create_EmptyToolId(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "CreateSandboxToolWithContext", func(_ context.Context, request *ags.CreateSandboxToolRequest) (*ags.CreateSandboxToolResponse, error) {
		return &ags.CreateSandboxToolResponse{Response: &ags.CreateSandboxToolResponseParams{}}, nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	createResp := &resource.CreateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Create(context.Background(), agsSandboxToolCreateRequest(t, spec), createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic")
	}
	if !containsString(createResp.Diagnostics.Errors()[0].Detail(), "empty ToolId") {
		t.Errorf("unexpected detail: %s", createResp.Diagnostics.Errors()[0].Detail())
	}
}

// TestAgsSandboxTool_Create_StatusFailed verifies that a FAILED status aborts
// the wait while keeping the identity in state.
func TestAgsSandboxTool_Create_StatusFailed(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "CreateSandboxToolWithContext", func(_ context.Context, request *ags.CreateSandboxToolRequest) (*ags.CreateSandboxToolResponse, error) {
		return &ags.CreateSandboxToolResponse{
			Response: &ags.CreateSandboxToolResponseParams{ToolId: agsString(agsTestToolId)},
		}, nil
	})
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		return agsSandboxToolResponse("FAILED", &ags.SandboxTool{
			ToolId:       agsString(agsTestToolId),
			Status:       agsString("FAILED"),
			StatusReason: agsString("image pull failed"),
		}), nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	createResp := &resource.CreateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Create(context.Background(), agsSandboxToolCreateRequest(t, spec), createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic")
	}
	if !containsString(createResp.Diagnostics.Errors()[0].Detail(), "image pull failed") {
		t.Errorf("unexpected detail: %s", createResp.Diagnostics.Errors()[0].Detail())
	}

	// The identity must already be in state so that destroy can clean up.
	var m AgsSandboxToolResourceModel
	if diags := createResp.State.Get(context.Background(), &m); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if m.Id.ValueString() != agsTestToolId {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), agsTestToolId)
	}
}

// TestAgsSandboxTool_Read_Success verifies the full refresh of every leaf
// attribute from the describe result.
func TestAgsSandboxTool_Read_Success(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		if request.ToolIds == nil || len(request.ToolIds) != 1 || *request.ToolIds[0] != agsTestToolId {
			got := ""
			if len(request.ToolIds) > 0 && request.ToolIds[0] != nil {
				got = *request.ToolIds[0]
			}
			t.Errorf("ToolIds = [%s], want [%s]", got, agsTestToolId)
		}
		return agsSandboxToolResponse("ACTIVE", agsSandboxToolFullItem()), nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolFullSpec()
	readResp := &resource.ReadResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Read(context.Background(), agsSandboxToolReadRequest(t, spec), readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics)
	}
	assertAgsSandboxToolState(t, readResp.State)
}

// TestAgsSandboxTool_Read_NotFound verifies that a confirmed miss removes the
// resource from state.
func TestAgsSandboxTool_Read_NotFound(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		return &ags.DescribeSandboxToolListResponse{
			Response: &ags.DescribeSandboxToolListResponseParams{
				SandboxToolSet: []*ags.SandboxTool{},
			},
		}, nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	raw := agsSandboxToolRawValue(t, spec)
	readResp := &resource.ReadResponse{State: agsSandboxToolState(t, raw)}
	r.Read(context.Background(), agsSandboxToolReadRequest(t, spec), readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics)
	}
	if !readResp.State.Raw.IsNull() {
		var m AgsSandboxToolResourceModel
		if diags := readResp.State.Get(context.Background(), &m); diags.HasError() {
			t.Fatalf("failed to read state after removal: %v", diags)
		}
		if !m.Id.IsNull() {
			t.Errorf("id = %q, want null after removal", m.Id.ValueString())
		}
	}
}

// TestAgsSandboxTool_Read_NilResponse verifies that a nil response is an API
// failure and must not remove the resource from state.
func TestAgsSandboxTool_Read_NilResponse(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		return &ags.DescribeSandboxToolListResponse{}, nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	readResp := &resource.ReadResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Read(context.Background(), agsSandboxToolReadRequest(t, spec), readResp)

	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for a nil response")
	}
	if !containsString(readResp.Diagnostics.Errors()[0].Summary(), "Error reading AGS sandbox tool") {
		t.Errorf("unexpected summary: %s", readResp.Diagnostics.Errors()[0].Summary())
	}
}

// TestAgsSandboxTool_Read_ApiError verifies that an API error surfaces a
// diagnostic.
func TestAgsSandboxTool_Read_ApiError(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		return nil, errors.New("Code=InternalError, Message=boom")
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	readResp := &resource.ReadResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Read(context.Background(), agsSandboxToolReadRequest(t, spec), readResp)

	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic")
	}
	if !containsString(readResp.Diagnostics.Errors()[0].Detail(), "InternalError") {
		t.Errorf("unexpected detail: %s", readResp.Diagnostics.Errors()[0].Detail())
	}
}

// TestAgsSandboxTool_Update_Success verifies that the mutable attributes are
// sent in the update request and that the state converges.
func TestAgsSandboxTool_Update_Success(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)

	var gotUpdateRequest *ags.UpdateSandboxToolRequest
	patches.ApplyMethodFunc(agsClient, "UpdateSandboxToolWithContext", func(_ context.Context, request *ags.UpdateSandboxToolRequest) (*ags.UpdateSandboxToolResponse, error) {
		gotUpdateRequest = request
		return &ags.UpdateSandboxToolResponse{Response: &ags.UpdateSandboxToolResponseParams{RequestId: agsString("fake-request-id")}}, nil
	})
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		return agsSandboxToolResponse("ACTIVE", agsSandboxToolFullItem()), nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	stateSpec := agsSandboxToolFullSpec()
	planSpec := agsSandboxToolFullSpec()
	planSpec.description = "updated description"
	planSpec.customImage = "ccr.ccs.tencentcloudcr.com/demo:2.0"

	updateResp := &resource.UpdateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, planSpec))}
	r.Update(context.Background(), agsSandboxToolUpdateRequest(t, stateSpec, planSpec), updateResp)

	if updateResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", updateResp.Diagnostics)
	}
	if gotUpdateRequest == nil {
		t.Fatal("UpdateSandboxTool was not called")
	}
	if got := gotUpdateRequest.ToolId; got == nil || *got != agsTestToolId {
		t.Errorf("ToolId = %v", got)
	}
	if got := gotUpdateRequest.Description; got == nil || *got != "updated description" {
		t.Errorf("Description = %v", got)
	}
	if gotUpdateRequest.NetworkConfiguration == nil || gotUpdateRequest.CustomConfiguration == nil || gotUpdateRequest.ComputerConfiguration == nil {
		t.Errorf("update request is missing nested configurations")
	}
	if len(gotUpdateRequest.Tags) != 1 {
		t.Errorf("Tags = %v", gotUpdateRequest.Tags)
	}

	var m AgsSandboxToolResourceModel
	if diags := updateResp.State.Get(context.Background(), &m); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if m.Id.ValueString() != agsTestToolId {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), agsTestToolId)
	}
	if m.Status.ValueString() != "ACTIVE" {
		t.Errorf("status = %q, want ACTIVE", m.Status.ValueString())
	}
}

// TestAgsSandboxTool_Update_ClearedDescription verifies that removing the
// description refuses to guess the cloud semantics.
func TestAgsSandboxTool_Update_ClearedDescription(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)

	apiCalled := false
	patches.ApplyMethodFunc(agsClient, "UpdateSandboxToolWithContext", func(_ context.Context, request *ags.UpdateSandboxToolRequest) (*ags.UpdateSandboxToolResponse, error) {
		apiCalled = true
		return &ags.UpdateSandboxToolResponse{Response: &ags.UpdateSandboxToolResponseParams{}}, nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	stateSpec := agsSandboxToolFullSpec()
	planSpec := agsSandboxToolFullSpec()
	planSpec.description = ""

	updateResp := &resource.UpdateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, planSpec))}
	r.Update(context.Background(), agsSandboxToolUpdateRequest(t, stateSpec, planSpec), updateResp)

	if !updateResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic when the description is removed")
	}
	if !containsString(updateResp.Diagnostics.Errors()[0].Summary(), "Clear semantics of description is not confirmed") {
		t.Errorf("unexpected summary: %s", updateResp.Diagnostics.Errors()[0].Summary())
	}
	if apiCalled {
		t.Error("UpdateSandboxTool must not be called when the clear semantics is unknown")
	}
}

// TestAgsSandboxTool_Update_Error verifies that an API error surfaces a
// diagnostic.
func TestAgsSandboxTool_Update_Error(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "UpdateSandboxToolWithContext", func(_ context.Context, request *ags.UpdateSandboxToolRequest) (*ags.UpdateSandboxToolResponse, error) {
		return nil, errors.New("Code=FailedOperation, Message=update failed")
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	stateSpec := agsSandboxToolFullSpec()
	planSpec := agsSandboxToolFullSpec()
	planSpec.description = "updated description"

	updateResp := &resource.UpdateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, planSpec))}
	r.Update(context.Background(), agsSandboxToolUpdateRequest(t, stateSpec, planSpec), updateResp)

	if !updateResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic")
	}
	if !containsString(updateResp.Diagnostics.Errors()[0].Summary(), "Error updating AGS sandbox tool") {
		t.Errorf("unexpected summary: %s", updateResp.Diagnostics.Errors()[0].Summary())
	}
}

// TestAgsSandboxTool_Delete_Success verifies that a delete waits until the
// tool disappears from the list.
func TestAgsSandboxTool_Delete_Success(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)

	deleteCalled := false
	patches.ApplyMethodFunc(agsClient, "DeleteSandboxToolWithContext", func(_ context.Context, request *ags.DeleteSandboxToolRequest) (*ags.DeleteSandboxToolResponse, error) {
		deleteCalled = true
		got := ""
		if request.ToolId != nil {
			got = *request.ToolId
		}
		if got != agsTestToolId {
			t.Errorf("ToolId = %q, want %q", got, agsTestToolId)
		}
		return &ags.DeleteSandboxToolResponse{Response: &ags.DeleteSandboxToolResponseParams{}}, nil
	})

	pollCount := 0
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		pollCount++
		if pollCount == 1 {
			return agsSandboxToolResponse("DELETING", &ags.SandboxTool{ToolId: agsString(agsTestToolId), Status: agsString("DELETING")}), nil
		}
		return agsSandboxToolResponse("", nil), nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	deleteResp := &resource.DeleteResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Delete(context.Background(), agsSandboxToolDeleteRequest(t, spec), deleteResp)

	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", deleteResp.Diagnostics)
	}
	if !deleteCalled {
		t.Error("DeleteSandboxTool was not called")
	}
	if pollCount < 2 {
		t.Errorf("describe poll count = %d, want at least 2", pollCount)
	}
}

// TestAgsSandboxTool_Delete_NotFound verifies that an "already gone" error is
// swallowed so that a repeated destroy is idempotent.
func TestAgsSandboxTool_Delete_NotFound(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)

	patches.ApplyMethodFunc(agsClient, "DeleteSandboxToolWithContext", func(_ context.Context, request *ags.DeleteSandboxToolRequest) (*ags.DeleteSandboxToolResponse, error) {
		return nil, tcsdk.NewTencentCloudSDKError("ResourceNotFound.SandboxTool", "sandbox tool not found", "req-1")
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	deleteResp := &resource.DeleteResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Delete(context.Background(), agsSandboxToolDeleteRequest(t, spec), deleteResp)

	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", deleteResp.Diagnostics)
	}
}

// TestAgsSandboxTool_Delete_Error verifies that other errors keep the resource
// in state for a retry.
func TestAgsSandboxTool_Delete_Error(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)

	patches.ApplyMethodFunc(agsClient, "DeleteSandboxToolWithContext", func(_ context.Context, request *ags.DeleteSandboxToolRequest) (*ags.DeleteSandboxToolResponse, error) {
		return nil, tcsdk.NewTencentCloudSDKError("FailedOperation", "delete failed", "req-1")
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}
	deleteResp := &resource.DeleteResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Delete(context.Background(), agsSandboxToolDeleteRequest(t, spec), deleteResp)

	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic")
	}
	if !containsString(deleteResp.Diagnostics.Errors()[0].Summary(), "Error deleting AGS sandbox tool") {
		t.Errorf("unexpected summary: %s", deleteResp.Diagnostics.Errors()[0].Summary())
	}
}

// TestAgsSandboxTool_ClientNotConfigured verifies the nil client guard of every
// lifecycle callback.
func TestAgsSandboxTool_ClientNotConfigured(t *testing.T) {
	r := NewAgsSandboxToolResource().(*AgsSandboxToolResource)
	ctx := context.Background()
	spec := agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}

	createResp := &resource.CreateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Create(ctx, agsSandboxToolCreateRequest(t, spec), createResp)
	if !createResp.Diagnostics.HasError() || !containsString(createResp.Diagnostics.Errors()[0].Summary(), "Provider not configured") {
		t.Errorf("create: unexpected diagnostics %v", createResp.Diagnostics)
	}

	readResp := &resource.ReadResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Read(ctx, agsSandboxToolReadRequest(t, spec), readResp)
	if !readResp.Diagnostics.HasError() || !containsString(readResp.Diagnostics.Errors()[0].Summary(), "Provider not configured") {
		t.Errorf("read: unexpected diagnostics %v", readResp.Diagnostics)
	}

	updateResp := &resource.UpdateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Update(ctx, agsSandboxToolUpdateRequest(t, spec, spec), updateResp)
	if !updateResp.Diagnostics.HasError() || !containsString(updateResp.Diagnostics.Errors()[0].Summary(), "Provider not configured") {
		t.Errorf("update: unexpected diagnostics %v", updateResp.Diagnostics)
	}

	deleteResp := &resource.DeleteResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Delete(ctx, agsSandboxToolDeleteRequest(t, spec), deleteResp)
	if !deleteResp.Diagnostics.HasError() || !containsString(deleteResp.Diagnostics.Errors()[0].Summary(), "Provider not configured") {
		t.Errorf("delete: unexpected diagnostics %v", deleteResp.Diagnostics)
	}
}

// TestAgsSandboxTool_ImportState verifies that the tool id is passed through to
// the id attribute.
func TestAgsSandboxTool_ImportState(t *testing.T) {
	r := NewAgsSandboxToolResource().(*AgsSandboxToolResource)

	importResp := &resource.ImportStateResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, agsSandboxToolValue{toolId: agsTestToolId, networkMode: "PUBLIC"}))}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: agsTestToolId}, importResp)

	if importResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", importResp.Diagnostics)
	}

	var id types.String
	if diags := importResp.State.GetAttribute(context.Background(), path.Root("id"), &id); diags.HasError() {
		t.Fatalf("failed to read id from state: %v", diags)
	}
	if id.ValueString() != agsTestToolId {
		t.Errorf("id = %q, want %q", id.ValueString(), agsTestToolId)
	}
}

// TestAgsSandboxTool_RefreshKeepsConfiguredValues verifies that attributes the
// cloud side either never echoes or normalises keep the configured value, so
// that a refreshed state does not produce a permanent diff.
func TestAgsSandboxTool_RefreshKeepsConfiguredValues(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)

	item := agsSandboxToolFullItem()
	// The cloud side normalises values compared to the configuration.
	item.DefaultTimeoutSeconds = agsUint64(600)
	item.CustomConfiguration.Command = []*string{agsString("/bin/bash")}
	item.StorageMounts[0].Name = agsString("renamed-by-cloud")
	item.NetworkConfiguration.NetworkMode = agsString("VPC")

	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		return agsSandboxToolResponse("ACTIVE", item), nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	spec := agsSandboxToolFullSpec()
	spec.defaultTimeout = "10m"

	readResp := &resource.ReadResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, spec))}
	r.Read(context.Background(), agsSandboxToolReadRequest(t, spec), readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics)
	}

	var m AgsSandboxToolResourceModel
	if diags := readResp.State.Get(context.Background(), &m); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}

	// Keep configured values even though the cloud reports something else.
	if got := m.DefaultTimeout.ValueString(); got != "10m" {
		t.Errorf("default_timeout = %q, want the configured value 10m", got)
	}
	if got := m.ClientToken.ValueString(); got != "token-1" {
		t.Errorf("client_token = %q, want the configured value token-1", got)
	}

	var mounts []agsSandboxToolStorageMountModel
	if diags := m.StorageMounts.ElementsAs(context.Background(), &mounts, false); diags.HasError() {
		t.Fatalf("failed to read storage_mounts: %v", diags)
	}
	if len(mounts) != 1 || mounts[0].Name.ValueString() != "rootfs" {
		t.Errorf("storage mount name must keep the configured value, got %v", mounts)
	}

	var customs []agsSandboxToolCustomConfigurationModel
	if diags := m.CustomConfiguration.ElementsAs(context.Background(), &customs, false); diags.HasError() {
		t.Fatalf("failed to read custom_configuration: %v", diags)
	}
	if len(customs) != 1 {
		t.Fatalf("custom_configuration length = %d, want 1", len(customs))
	}
	if got := agsSetStrings(context.Background(), t, customs[0].Command); len(got) != 1 || got[0] != "/bin/sh" {
		t.Errorf("command must keep the configured value, got %v", got)
	}

	// Cloud owned outputs are still refreshed.
	if got := m.Status.ValueString(); got != "ACTIVE" {
		t.Errorf("status = %q, want ACTIVE", got)
	}
}

// TestAgsSandboxTool_ReadImportFillsNestedBlocks verifies that an empty state
// (as produced by `terraform import`) is filled from the describe result.
func TestAgsSandboxTool_ReadImportFillsNestedBlocks(t *testing.T) {
	patches, agsClient := agsSandboxToolTestSetup(t)
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolListWithContext", func(_ context.Context, request *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
		return agsSandboxToolResponse("ACTIVE", agsSandboxToolFullItem()), nil
	})

	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := agsSandboxToolResourceWithClient(client)

	// Only the identity is present, every other attribute is null.
	imported := agsSandboxToolValue{toolId: agsTestToolId}
	readResp := &resource.ReadResponse{State: agsSandboxToolState(t, agsSandboxToolRawValue(t, imported))}
	r.Read(context.Background(), agsSandboxToolReadRequest(t, imported), readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics)
	}

	assertAgsSandboxToolState(t, readResp.State)
}

// containsString reports whether s contains substr.
func containsString(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && indexOf(s, substr) >= 0)
}

// indexOf returns the index of substr inside s, or -1.
func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
