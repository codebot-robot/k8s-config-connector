# ComputeServiceAttachment KRM Implementation Journal

## Observations & Quirks

- **API Compatibility & Nested Types**: 
  When migrating the existing TF-based `ComputeServiceAttachment` (v1beta1) to the direct KRM types, we had to preserve exact field names and types under `spec` and `status` to remain 100% backward compatible and prevent breaking changes.
  
- **Mapping Slices of Structs**:
  The direct code-generation tool expects sub-struct types to match the protobuf-defined type names (e.g. `ServiceAttachmentConsumerProjectLimit` instead of `ServiceattachmentConsumerAcceptLists`) for automated generation of slice mapping helpers. Rather than changing the nested type names in the CRD (which would alter schema names, though K8s inlines them), we chose to hand-write custom conversion functions (`ServiceattachmentConsumerAcceptLists_v1beta1_FromProto`, `ServiceattachmentConsumerAcceptLists_v1beta1_ToProto`, etc.) in a new `/pkg/controller/direct/compute/computeserviceattachment_mapper.go` file.

- **Handling References**:
  - `consumerRejectLists` previously used a list of raw resource references. We mapped them to `[]refsv1beta1.ProjectRef` to follow direct KRM standards while maintaining structural compatibility.
  - `natSubnets` was mapped to `[]ComputeSubnetworkRef`. Since its name doesn't contain a `Ref` suffix and cannot be changed without breaking existing configs, we added it to `tests/apichecks/testdata/exceptions/missingrefs.txt`.
  - `targetServiceRef` was mapped as a pointer `*refsv1beta1.ComputeForwardingRuleRef` to match the pointer mapping pattern expected by the direct mapper, while keeping CRD OpenAPI schema fully compatible.

- **Compilation Safeguard**:
  The direct code generator (`apis/compute/v1beta1/generate.sh`) seamlessly detects hand-written conversion helpers in the same package and comments out redundant generated functions in `mapper.generated.go`. This keeps the auto-generated code and the custom hand-written mappings completely synchronized and compiled flawlessly.
