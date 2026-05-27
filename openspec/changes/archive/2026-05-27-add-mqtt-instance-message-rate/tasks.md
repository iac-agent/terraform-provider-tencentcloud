## 1. Schema & CRUD Implementation

- [x] 1.1 Add `message_rate` parameter to the `tencentcloud_mqtt_instance` resource schema in `resource_tc_mqtt_instance.go` (TypeInt, Optional+Computed, with description)
- [x] 1.2 In the Create function, add `message_rate` to the post-creation ModifyInstance call (alongside `automatic_activation` and `authorization_policy`), setting `modifyRequest.MessageRate` when the user has specified it
- [x] 1.3 In the Read function, add logic to read `MessageRate` from the DescribeInstance response and set it in state (with nil check before calling d.Set)
- [x] 1.4 In the Update function, add `message_rate` to the `mutableArgs` list and include it in the ModifyInstance request when changed

## 2. Unit Tests

- [x] 2.1 Add unit test cases for `message_rate` parameter in `resource_tc_mqtt_instance_test.go` covering create, read, and update scenarios using gomonkey mocks

## 3. Documentation

- [x] 3.1 Update the `resource_tc_mqtt_instance.md` file to include `message_rate` in the example usage
