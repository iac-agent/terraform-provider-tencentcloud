## 1. Schema 与 CRUD 函数实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.go`，定义 ResourceTencentCloudTeoJustInTimeTranscodeTemplate() 函数，包含完整 schema（zone_id、template_name、comment、video_stream_switch、audio_stream_switch、video_template、audio_template、template_id、create_time、update_time、type）及嵌套结构（video_template 含 codec/fps/bitrate/resolution_adaptive/width/height/fill_type，audio_template 含 codec/audio_channel），所有非 Computed 字段设置 ForceNew，无 Update 函数
- [x] 1.2 实现 Create 函数：调用 CreateJustInTimeTranscodeTemplate API，构造 VideoTemplateInfo 和 AudioTemplateInfo 子结构，使用 tccommon.ReadRetryTimeout 重试，创建成功后设置复合 ID（zoneId#templateId），然后调用 Read 刷新状态
- [x] 1.3 实现 Read 函数：拆分复合 ID 获取 zoneId 和 templateId，调用 DescribeJustInTimeTranscodeTemplates API 并使用 Filter(template-id) 过滤查询指定模板，Limit 设为 1000，从 TemplateSet 中匹配 templateId，若未找到则设置 d.SetId("")，将响应字段映射到 state
- [x] 1.4 实现 Delete 函数：拆分复合 ID 获取 zoneId 和 templateId，调用 DeleteJustInTimeTranscodeTemplates API，传入 ZoneId 和 TemplateIds=[templateId]，使用 tccommon.ReadRetryTimeout 重试
- [x] 1.5 配置 Importer：使用 schema.ImportStatePassthrough 支持资源导入

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册资源 `tencentcloud_teo_just_in_time_transcode_template`，引用 `teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()`
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_just_in_time_transcode_template` 资源条目

## 3. 单元测试

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template_test.go`，使用 gomonkey mock 云 API 调用，编写 Create/Read/Delete 操作的单元测试
- [x] 3.2 使用 `go test` 运行单元测试，确保所有测试通过

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.md`，包含一句话描述（提及 TEO）、Example Usage（含 video_template 和 audio_template 配置）和 Import 部分（格式：zoneId#templateId）
