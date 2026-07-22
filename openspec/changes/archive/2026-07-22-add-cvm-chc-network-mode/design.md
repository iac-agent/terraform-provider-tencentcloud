## Context

Terraform Provider for TencentCloud需要新增CVM产品的CHC服务器网络模式切换资源。当前CVM产品提供了`ModifyChcNetworkMode`接口，支持将CHC物理服务器在部署网络模式（DEPLOY）和业务网络模式（BUSINESS）之间切换。

关键约束：
- 该资源仅有Modify接口（`ModifyChcNetworkMode`），没有独立的Create/Delete接口
- `DescribeChcHosts`接口返回的`ChcHost`结构体中没有`NetworkMode`字段，无法通过Read接口获取当前网络模式
- 这属于只有CRD接口的资源类型，需要按照CRD模式处理

参考实现：`tencentcloud_cvm_chc_config`资源（同为CVM CHC相关资源）

## Goals / Non-Goals

**Goals:**
- 封装`ModifyChcNetworkMode`接口为Terraform资源，支持CHC服务器网络模式的切换
- 资源Create操作调用`ModifyChcNetworkMode`切换到目标网络模式
- 资源Update操作再次调用`ModifyChcNetworkMode`切换到新的网络模式
- 资源Read操作使用`DescribeChcHosts`查询CHC主机状态，网络模式字段从terraform state中读取（因云API不返回该字段）
- 资源Delete操作从terraform state中移除，不做实际API调用（因为网络模式是状态切换，无法真正"删除"）
- 在provider.go中注册新资源

**Non-Goals:**
- 不修改现有`cvm_chc_config`资源的行为
- 不实现网络模式的查询数据源（`ChcHost`结构体中无`NetworkMode`字段）
- 不支持import（因为Read无法确认当前网络模式状态）

## Decisions

### 1. 资源ID方案：使用chc_ids拼接作为资源ID
**选择**: 将`chc_ids`列表中的ID用`tccommon.FILED_SP`拼接作为资源ID
**理由**: 该资源没有独立的资源ID，`ModifyChcNetworkMode`接口也不返回任何ID。使用chc_ids作为标识符最合理。
**替代方案**: 使用UUID作为ID - 但这样无法从ID反推资源对应的CHC服务器

### 2. CRD模式处理：Id字段ForceNew + immutableArgs
**选择**: 将id字段设置为ForceNew，update方法中将network_mode加入immutableArgs
**理由**: 根据规则，若一个资源只有CRD接口，则只将Id()字段设置成ForceNew，并在资源update方法中将其余顶层字段加入immutableArgs数组。chc_ids作为ForceNew是因为ID由它决定；network_mode在update中检查是否变化，若变化则调用ModifyChcNetworkMode。
**修正**: 实际上，由于该资源的update操作就是调用同一个ModifyChcNetworkMode接口，network_mode的变化应该被允许并调用API更新，而非作为immutable参数。重新评估：该资源只有Modify接口，Create/Update/Delete都是调用同一个接口，因此network_mode的变化应允许更新。

### 3. Read操作策略：从state读取network_mode
**选择**: Read操作中，`network_mode`字段不从云API获取（因为`DescribeChcHosts`不返回此字段），而是从terraform state中保持不变
**理由**: `ChcHost`结构体中没有`NetworkMode`字段，无法从云API确认当前网络模式。只能在Read时保留state中的值。

### 4. Delete操作策略：仅移除state
**选择**: Delete操作仅调用`d.SetId("")`移除资源，不调用云API
**理由**: 网络模式切换是一个状态操作，没有对应的"恢复"或"删除"API。删除Terraform资源时只需要从state中移除。

## Risks / Trade-offs

- [云API不返回NetworkMode] → Read操作无法检测外部变更。如果用户在Terraform之外修改了网络模式，Terraform不会检测到漂移。这是云API的限制，无法规避。
- [Delete不做实际操作] → 资源删除不会恢复网络模式。用户需要注意destroy操作不会切换回原始模式。这是符合预期的行为，在文档中说明即可。
