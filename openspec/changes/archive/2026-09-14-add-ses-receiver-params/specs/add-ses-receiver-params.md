## ADDED Requirements

After thorough analysis of the vendor SDK and existing resource code, all specified parameters are already implemented in `resource_tc_ses_receiver.go`. No new parameters need to be added.

### Summary of findings

- `receivers_name`: Already implemented in schema and CreateReceiver/ListReceivers
- `desc`: Already implemented in schema and CreateReceiver/ListReceivers
- `data`: Already implemented in schema
  - `email`: Already implemented in schema and CreateReceiverDetailWithData/CreateReceiverDetail
  - `template_data`: Already implemented in schema and CreateReceiverDetailWithData
- `InvalidCount`: Does NOT exist in the vendor SDK `ReceiverData` struct - cannot be implemented